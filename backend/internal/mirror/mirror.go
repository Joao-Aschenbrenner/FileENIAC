package mirror

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ENIACSystems/FileENIAC/backend/internal/deploy/hardening"
	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/registry"
	"github.com/ENIACSystems/FileENIAC/backend/internal/transports"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Snapshot struct {
	ID          string `json:"id"`
	ProjectID   int64  `json:"project_id"`
	FilesCount  int    `json:"files_count"`
	TotalSize   int64  `json:"total_size"`
	Status      string `json:"status"`
	StartedAt   string `json:"started_at"`
	CompletedAt string `json:"completed_at"`
}

type Engine struct{}

func New() *Engine {
	return &Engine{}
}

// MirrorPath returns the mirror directory path for a project.
func MirrorPath(wsPath, projectName string) string {
	return filepath.Join(wsPath, ".eniac", "mirror", projectName)
}

// Create downloads the remote server content to the local mirror directory.
// It preserves the directory structure and downloads files recursively.
func (e *Engine) Create(ctx *workspace.Context, projectName string) (*Snapshot, error) {
	proj, err := registry.GetProject(ctx, projectName)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	srv, err := registry.GetServer(ctx, proj.ID)
	if err != nil {
		return nil, fmt.Errorf("no server configured: %w", err)
	}

	snapshotID := uuid.New().String()[:12]
	startedAt := time.Now().Format(time.RFC3339)

	snapshot := &Snapshot{
		ID:        snapshotID,
		ProjectID: proj.ID,
		Status:    "running",
		StartedAt: startedAt,
	}

	// Record snapshot start
	_, err = ctx.DB.Exec(
		`INSERT INTO mirror_snapshots (project_id, snapshot_id, status, started_at) VALUES (?, ?, 'running', ?)`,
		proj.ID, snapshotID, startedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to record snapshot: %w", err)
	}

	log.L().Info("mirror started",
		zap.String("project", projectName),
		zap.String("snapshot_id", snapshotID),
	)

	mirrorDir := MirrorPath(ctx.Workspace.Path, projectName)
	if err := os.MkdirAll(mirrorDir, 0700); err != nil {
		e.failSnapshot(ctx, proj.ID, snapshotID, err)
		return nil, fmt.Errorf("failed to create mirror dir: %w", err)
	}

	transportCfg := transports.TransportConfig{
		Protocol: "ftp",
		Host:     srv.Host,
		Port:     srv.Port,
		User:     srv.User,
		Pass:     srv.Password,
		Timeout:  120 * time.Second,
	}

	client, err := transports.New(transportCfg)
	if err != nil {
		e.failSnapshot(ctx, proj.ID, snapshotID, err)
		return nil, fmt.Errorf("ftps connection failed: %w", err)
	}
	if err := client.Connect(context.Background()); err != nil {
		e.failSnapshot(ctx, proj.ID, snapshotID, err)
		return nil, fmt.Errorf("ftps connection failed: %w", err)
	}
	defer client.Disconnect()

	// Mirror remote directory recursively
	downloaded, totalSize, err := e.mirrorDir(client, srv.TargetPath, mirrorDir)
	if err != nil {
		e.failSnapshot(ctx, proj.ID, snapshotID, err)
		return nil, fmt.Errorf("failed to mirror directory: %w", err)
	}

	completedAt := time.Now().Format(time.RFC3339)
	snapshot.FilesCount = downloaded
	snapshot.TotalSize = totalSize
	snapshot.Status = "completed"
	snapshot.CompletedAt = completedAt

	_, err = ctx.DB.Exec(
		`UPDATE mirror_snapshots SET files_count=?, total_size=?, status='completed', completed_at=? WHERE snapshot_id=?`,
		downloaded, totalSize, completedAt, snapshotID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update snapshot: %w", err)
	}

	log.L().Info("mirror completed",
		zap.String("project", projectName),
		zap.Int("files", downloaded),
		zap.Int64("size", totalSize),
	)

	return snapshot, nil
}

// mirrorDir recursively mirrors a remote directory to a local path.
func (e *Engine) mirrorDir(client transports.Transport, remotePath, localPath string) (int, int64, error) {
	entries, err := client.List(context.Background(), remotePath)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to list remote dir %s: %w", remotePath, err)
	}

	if err := os.MkdirAll(localPath, 0700); err != nil {
		return 0, 0, fmt.Errorf("failed to create local dir %s: %w", localPath, err)
	}

	var totalFiles int
	var totalSize int64

	for _, entry := range entries {
		if entry.Name == "." || entry.Name == ".." {
			continue
		}

		remoteChild := remotePath + "/" + entry.Name
		localChild := filepath.Join(localPath, entry.Name)

		if entry.IsDir {
			subFiles, subSize, err := e.mirrorDir(client, remoteChild, localChild)
			if err != nil {
				log.L().Warn("mirror subdir failed, skipping",
					zap.String("dir", entry.Name),
					zap.Error(err),
				)
				continue
			}
			totalFiles += subFiles
			totalSize += subSize
		} else {
			dlErr := hardening.DoWithRetry(func() error {
				return client.Download(context.Background(), remoteChild, localChild)
			}, hardening.DefaultRetryConfig())

			if dlErr != nil {
				log.L().Warn("mirror download failed, skipping",
					zap.String("file", entry.Name),
					zap.Error(dlErr),
				)
				continue
			}

			if info, statErr := os.Stat(localChild); statErr == nil {
				totalSize += info.Size()
			}
			totalFiles++
		}
	}

	return totalFiles, totalSize, nil
}

// Status returns the latest snapshot info for a project.
func (e *Engine) Status(ctx *workspace.Context, projectName string) (*Snapshot, error) {
	proj, err := registry.GetProject(ctx, projectName)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	s := &Snapshot{}
	err = ctx.DB.QueryRow(
		`SELECT snapshot_id, project_id, files_count, total_size, status, COALESCE(started_at,''), COALESCE(completed_at,'') FROM mirror_snapshots WHERE project_id = ? ORDER BY id DESC LIMIT 1`,
		proj.ID,
	).Scan(&s.ID, &s.ProjectID, &s.FilesCount, &s.TotalSize, &s.Status, &s.StartedAt, &s.CompletedAt)
	if err != nil {
		return nil, fmt.Errorf("no mirror snapshot for project %s", projectName)
	}
	return s, nil
}

func (e *Engine) failSnapshot(ctx *workspace.Context, projectID int64, snapshotID string, err error) {
	if _, dbErr := ctx.DB.Exec(
		`UPDATE mirror_snapshots SET status='failed' WHERE snapshot_id=?`,
		snapshotID,
	); dbErr != nil {
		log.L().Error("mirror failSnapshot: failed to update DB", zap.String("snapshot_id", snapshotID), zap.Error(dbErr))
	}
	log.L().Error("mirror failed",
		zap.Int64("project_id", projectID),
		zap.String("snapshot_id", snapshotID),
		zap.Error(err),
	)
}
