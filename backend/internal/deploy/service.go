package deploy

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ENIACSystems/FileENIAC/backend/internal/database"
	"github.com/ENIACSystems/FileENIAC/backend/internal/deploy/hardening"
	"github.com/ENIACSystems/FileENIAC/backend/internal/deploy/packer"
	"github.com/ENIACSystems/FileENIAC/backend/internal/history"
	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/registry"
	"github.com/ENIACSystems/FileENIAC/backend/internal/transports"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Package-level mockable interfaces and function variables for testing.
type (
	artifactPacker interface {
		Pack(sourceDir, outputPath string) (*packer.Result, error)
		SetExcludes(excludes []string)
	}
)

var newPackerFn = func(excludes []string) artifactPacker { return packer.NewBuilder(excludes) }
var newTransportFn = func(cfg transports.TransportConfig) (transports.Transport, error) { return transports.New(cfg) }

type Result struct {
	DeployID   string `json:"deploy_id"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Artifact   string `json:"artifact,omitempty"`
	FilesCount int    `json:"files_count,omitempty"`
}

type Service struct {
	db       *database.DB
	history  *history.Service
	cb       *hardening.CircuitBreaker
	retryCfg hardening.RetryConfig
}

func NewService(db *database.DB) *Service {
	return &Service{
		db:       db,
		history:  history.NewService(db),
		cb:       hardening.NewCircuitBreaker(hardening.DefaultCircuitBreakerConfig()),
		retryCfg: hardening.DefaultRetryConfig(),
	}
}

func (s *Service) Deploy(ctx *workspace.Context, projectName string, useFallback bool) (*Result, error) {
	project, err := registry.GetProject(ctx, projectName)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	server, err := registry.GetServer(ctx, project.ID)
	if err != nil {
		return nil, fmt.Errorf("no server configured: %w", err)
	}

	deployID := uuid.New().String()[:8]
	startedAt := time.Now().Format(time.RFC3339)

	s.history.RecordEvent(history.EventDeployStarted,
		fmt.Sprintf("Deploy %s started for project %s", deployID, projectName),
		map[string]interface{}{
			"deploy_id": deployID,
			"project":   projectName,
			"fallback":  useFallback,
		})

	log.L().Info("deploy started",
		zap.String("deploy_id", deployID),
		zap.String("project", projectName),
		zap.Bool("fallback", useFallback),
	)

	tmpArtifact := filepath.Join(os.TempDir(), fmt.Sprintf("fileeniac-deploy-%s.tar.gz", deployID))

	// Cleanup temp artifact on function exit
	defer func() {
		if _, err := os.Stat(tmpArtifact); err == nil {
			os.Remove(tmpArtifact)
		}
	}()

	builder := newPackerFn(nil)
	packResult, err := builder.Pack(project.LocalPath, tmpArtifact)
	if err != nil {
		s.recordFailure(deployID, project.ID, "pack failed", err, startedAt)
		return nil, fmt.Errorf("pack failed: %w", err)
	}

	// Compute artifact hash for integrity verification
	artifactHash, err := hardening.ChecksumFile(tmpArtifact)
	if err != nil {
		s.recordFailure(deployID, project.ID, "checksum failed", err, startedAt)
		return nil, fmt.Errorf("checksum failed: %w", err)
	}

	transportCfg := transports.TransportConfig{
		Protocol: "ftp",
		Host:     server.Host,
		Port:     server.Port,
		User:     server.User,
		Pass:     server.Password,
		Timeout:  120 * time.Second,
	}

	// Upload with retry and circuit breaker
	uploadErr := s.cb.Execute(func() error {
		return hardening.DoWithRetry(func() error {
			client, err := newTransportFn(transportCfg)
			if err != nil {
				return fmt.Errorf("ftps connect: %w", err)
			}
			if err := client.Connect(context.Background()); err != nil {
				return fmt.Errorf("ftps connect: %w", err)
			}
			defer client.Disconnect()

			remoteArtifactPath := server.TargetPath + "/" + filepath.Base(tmpArtifact)
			if err := client.Upload(context.Background(), tmpArtifact, remoteArtifactPath); err != nil {
				return fmt.Errorf("upload: %w", err)
			}

			// Integrity verification: download and compare hash
			verifyPath := tmpArtifact + ".verify"
			if err := client.Download(context.Background(), remoteArtifactPath, verifyPath); err != nil {
				os.Remove(verifyPath)
				return fmt.Errorf("verify download: %w", err)
			}
			defer os.Remove(verifyPath)

			if err := hardening.VerifyIntegrity(verifyPath, artifactHash); err != nil {
				return fmt.Errorf("integrity: %w", err)
			}

			return nil
		}, s.retryCfg)
	})
	if uploadErr != nil {
		s.recordFailure(deployID, project.ID, "upload failed", uploadErr, startedAt)
		return nil, fmt.Errorf("upload failed: %w", uploadErr)
	}

	targetDir := server.TargetPath
	manifestContent := fmt.Sprintf(`{"deploy_id":"%s","project":"%s","timestamp":"%s"}`,
		deployID, projectName, time.Now().Format(time.RFC3339))
	manifestPath := targetDir + "/deploy-manifest.json"

	manifestUploadErr := s.cb.Execute(func() error {
		return hardening.DoWithRetry(func() error {
			client, err := newTransportFn(transportCfg)
			if err != nil {
				return fmt.Errorf("manifest ftps connect: %w", err)
			}
			if err := client.Connect(context.Background()); err != nil {
				return fmt.Errorf("manifest ftps connect: %w", err)
			}
			defer client.Disconnect()

			tmpManifest, err := os.CreateTemp("", "manifest-*.json")
			if err != nil {
				return fmt.Errorf("manifest temp file: %w", err)
			}
			if _, err := tmpManifest.WriteString(manifestContent); err != nil {
				tmpManifest.Close()
				os.Remove(tmpManifest.Name())
				return fmt.Errorf("manifest write: %w", err)
			}
			tmpManifest.Close()
			defer os.Remove(tmpManifest.Name())

			if err := client.Upload(context.Background(), tmpManifest.Name(), manifestPath); err != nil {
				return fmt.Errorf("manifest upload: %w", err)
			}
			return nil
		}, s.retryCfg)
	})
	if manifestUploadErr != nil {
		log.L().Error("manifest upload failed", zap.Error(manifestUploadErr))
	}

	completedAt := time.Now().Format(time.RFC3339)
	deployLog := &history.DeployLog{
		ProjectID:    project.ID,
		DeployID:     deployID,
		Status:       "success",
		ArtifactHash: artifactHash,
		FilesCount:   packResult.FileCount,
		StartedAt:    startedAt,
		CompletedAt:  completedAt,
	}
	s.history.RecordDeploy(deployLog)

	log.L().Info("deploy completed",
		zap.String("deploy_id", deployID),
		zap.String("project", projectName),
	)

	return &Result{
		DeployID:   deployID,
		Status:     "success",
		Message:    fmt.Sprintf("Deploy %s completed successfully", deployID),
		Artifact:   tmpArtifact,
		FilesCount: packResult.FileCount,
	}, nil
}

func (s *Service) Rollback(ctx *workspace.Context, projectName string) (*Result, error) {
	project, err := registry.GetProject(ctx, projectName)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	lastDeploy, err := s.history.GetLastDeploy(project.ID)
	if err != nil {
		return nil, fmt.Errorf("no deploy to rollback: %w", err)
	}

	s.history.RecordEvent(history.EventRollbackStarted,
		fmt.Sprintf("Rollback started for deploy %s (project %s)", lastDeploy.DeployID, projectName),
		map[string]interface{}{"deploy_id": lastDeploy.DeployID, "project": projectName})

	server, err := registry.GetServer(ctx, project.ID)
	if err == nil {
		transportCfg := transports.TransportConfig{
			Protocol: "ftp",
			Host:     server.Host,
			Port:     server.Port,
			User:     server.User,
			Pass:     server.Password,
			Timeout:  120 * time.Second,
		}
		client, err := newTransportFn(transportCfg)
		if err == nil {
			if err := client.Connect(context.Background()); err == nil {
				artifactName := fmt.Sprintf("fileeniac-deploy-%s.tar.gz", lastDeploy.DeployID)
				remoteArtifactPath := server.TargetPath + "/" + artifactName
				manifestPath := server.TargetPath + "/deploy-manifest.json"

				if err := client.Delete(context.Background(), remoteArtifactPath); err != nil {
					log.L().Warn("rollback: failed to delete artifact from server",
						zap.String("path", remoteArtifactPath), zap.Error(err))
				}
				if err := client.Delete(context.Background(), manifestPath); err != nil {
					log.L().Warn("rollback: failed to delete manifest from server",
						zap.String("path", manifestPath), zap.Error(err))
				}
				client.Disconnect()
			} else {
				log.L().Warn("rollback: unable to connect to FTPS server; proceeding with log-only rollback",
					zap.Error(err))
			}
		} else {
			log.L().Warn("rollback: unable to create transport; proceeding with log-only rollback",
				zap.Error(err))
		}
	} else {
		log.L().Warn("rollback: unable to look up server; proceeding with log-only rollback",
			zap.Error(err))
	}

	rollbackLog := &history.RollbackLog{
		DeployID:  lastDeploy.DeployID,
		ProjectID: project.ID,
		Reason:    "manual rollback",
		Status:    "completed",
	}
	s.history.RecordRollback(rollbackLog)

	log.L().Info("rollback completed",
		zap.String("deploy_id", lastDeploy.DeployID),
		zap.String("project", projectName),
	)

	return &Result{
		DeployID: lastDeploy.DeployID,
		Status:   "rolled_back",
		Message:  fmt.Sprintf("Rollback of deploy %s completed", lastDeploy.DeployID),
	}, nil
}

func (s *Service) Verify(ctx *workspace.Context, projectName string) (*Result, error) {
	project, err := registry.GetProject(ctx, projectName)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	lastDeploy, err := s.history.GetLastDeploy(project.ID)
	if err != nil {
		return &Result{Status: "unknown", Message: "No deploys found for this project"}, nil
	}

	s.history.RecordEvent(history.EventVerifySuccess,
		fmt.Sprintf("Verify for project %s: last deploy %s", projectName, lastDeploy.DeployID),
		map[string]interface{}{
			"project":   projectName,
			"deploy_id": lastDeploy.DeployID,
			"status":    lastDeploy.Status,
		})

	return &Result{
		DeployID: lastDeploy.DeployID,
		Status:   lastDeploy.Status,
		Message:  fmt.Sprintf("Last deploy: %s (%s)", lastDeploy.DeployID, lastDeploy.Status),
	}, nil
}

func (s *Service) Validate(ctx *workspace.Context, projectName string) error {
	_, err := registry.GetProject(ctx, projectName)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	project, err := registry.GetProject(ctx, projectName)
	if err != nil {
		return err
	}

	_, err = registry.GetServer(ctx, project.ID)
	if err != nil {
		return fmt.Errorf("no server configured for project %s: %w", projectName, err)
	}

	return nil
}

func (s *Service) GetHistory(ctx *workspace.Context, projectName string, limit int) ([]*history.DeployLog, error) {
	project, err := registry.GetProject(ctx, projectName)
	if err != nil {
		return nil, err
	}
	return s.history.GetDeployHistory(project.ID, limit)
}

func (s *Service) recordFailure(deployID string, projectID int64, stage string, err error, startedAt string) {
	completedAt := time.Now().Format(time.RFC3339)
	deployLog := &history.DeployLog{
		ProjectID:    projectID,
		DeployID:     deployID,
		Status:       "failed",
		ErrorMessage: fmt.Sprintf("%s: %v", stage, err),
		StartedAt:    startedAt,
		CompletedAt:  completedAt,
	}
	s.history.RecordDeploy(deployLog)

	s.history.RecordEvent(history.EventDeployFailed,
		fmt.Sprintf("Deploy %s failed at %s: %v", deployID, stage, err),
		map[string]interface{}{"deploy_id": deployID, "stage": stage, "error": err.Error()})
}
