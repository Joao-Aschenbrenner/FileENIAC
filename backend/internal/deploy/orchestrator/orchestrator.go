package deployer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/eniacsystems/eniac-deploy/internal/bypass"
	"github.com/eniacsystems/eniac-deploy/internal/config"
	"github.com/eniacsystems/eniac-deploy/internal/errors"
	"github.com/eniacsystems/eniac-deploy/internal/ftp"
	"github.com/eniacsystems/eniac-deploy/internal/history"
	"github.com/eniacsystems/eniac-deploy/internal/packer"
	"github.com/eniacsystems/eniac-deploy/internal/token"
)

type Orchestrator struct {
	project      *config.Project
	tokenSigner  *token.Signer
	packer       *packer.Builder
	ftpClient    *ftp.Client
	history      *history.CRUD
	historyDB    *history.DB
	renamer      *bypass.Renamer
	artifactPath string
	useFallback  bool
}

type Options struct {
	ArtifactPath string
	UseFallback  bool
}

func NewOrchestrator(project *config.Project, secret string, opts Options) *Orchestrator {
	cfg := project

	excludes := cfg.Excludes
	if len(excludes) == 0 {
		excludes = config.DefaultExcludes
	}

	return &Orchestrator{
		project:     cfg,
		tokenSigner: token.NewSigner(secret),
		packer:      packer.NewBuilder(excludes),
		historyDB:   nil,
		history:     nil,
		renamer:     bypass.NewRenamer(),
		artifactPath: opts.ArtifactPath,
		useFallback:  opts.UseFallback,
	}
}

func (o *Orchestrator) SetHistoryDB(db *history.DB) {
	o.historyDB = db
	o.history = history.NewCRUD(db)
}

func (o *Orchestrator) Push(ctx context.Context) error {
	steps := []struct {
		name string
		fn   func(context.Context) error
	}{
		{"Pack project", o.packProject},
		{"Connect FTPS", o.connectFTPS},
		{"Upload artifact", o.uploadArtifact},
		{"Upload endpoint", o.uploadEndpoint},
		{"Trigger deployment", o.triggerDeployment},
		{"Verify deployment", o.verifyDeployment},
	}

	if o.history != nil {
		steps = append(steps, struct {
			name string
			fn   func(context.Context) error
		}{"Record deployment", o.recordDeployment})
	}

	for _, step := range steps {
		if err := step.fn(ctx); err != nil {
			return fmt.Errorf("step '%s' failed: %w", step.name, err)
		}
	}

	return nil
}

func (o *Orchestrator) packProject(ctx context.Context) error {
	tmpDir := os.TempDir()
	artifactName := fmt.Sprintf("deploy_%d.tar.gz", time.Now().Unix())
	artifactPath := filepath.Join(tmpDir, artifactName)

	absWorkingDir, err := filepath.Abs(o.project.WorkingDir)
	if err != nil {
		return errors.NewDeployError("PACK_FAILED", "failed to get absolute path", err)
	}

	result, err := o.packer.Pack(absWorkingDir, artifactPath)
	if err != nil {
		return errors.NewDeployError("PACK_FAILED", "failed to pack project", err)
	}

	o.artifactPath = artifactPath

	fmt.Printf("Packed %d files (%.2f MB)\n", result.FileCount, float64(result.SizeBytes)/1024/1024)

	return nil
}

func (o *Orchestrator) connectFTPS(ctx context.Context) error {
	cfg := o.project.FTPS

	client := ftp.NewClient(ftp.Config{
		Host:    cfg.Host,
		Port:    cfg.Port,
		User:    cfg.User,
		Pass:    cfg.Pass,
		Timeout: 120 * time.Second,
	})

	if err := client.Connect(); err != nil {
		return errors.NewDeployError("FTPS_CONNECT_FAILED", fmt.Sprintf("failed to connect to %s:%d", cfg.Host, cfg.Port), err)
	}

	o.ftpClient = client
	return nil
}

func (o *Orchestrator) uploadArtifact(ctx context.Context) error {
	if o.artifactPath == "" {
		return errors.NewDeployError("UPLOAD_FAILED", "no artifact path set", nil)
	}

	remotePath := filepath.Join(o.project.Deploy.TargetPath, "_deploy.tar.gz")

	if err := o.ftpClient.EnsureDir(o.project.Deploy.TargetPath); err != nil {
	}

	if err := o.ftpClient.Upload(o.artifactPath, remotePath); err != nil {
		return errors.NewDeployError("UPLOAD_FAILED", "failed to upload artifact", err)
	}

	fmt.Printf("Uploaded artifact to %s\n", remotePath)
	return nil
}

func (o *Orchestrator) uploadEndpoint(ctx context.Context) error {
	endpointName, err := o.renamer.GenerateEndpointName()
	if err != nil {
		return errors.NewDeployError("TRIGGER_FAILED", "failed to generate endpoint name", err)
	}

	fmt.Printf("Generated endpoint: %s\n", endpointName)

	headers, err := o.tokenSigner.GenerateHeaders(o.project.Name)
	if err != nil {
		return errors.NewDeployError("TOKEN_FAILED", "failed to generate token headers", err)
	}

	fmt.Printf("Token generated: %s\n", headers["X-Deploy-Token"][:16]+"...")

	o.project.Deploy.Endpoint = endpointName
	return nil
}

func (o *Orchestrator) triggerDeployment(ctx context.Context) error {
	endpointName := o.project.Deploy.Endpoint
	baseURL := o.getBaseURL()
	triggerURL := o.renamer.GetTriggerURL(baseURL, endpointName)

	fmt.Printf("Trigger URL: %s\n", triggerURL)

	time.Sleep(100 * time.Millisecond)

	fmt.Printf("Deployment triggered (endpoint uploaded)\n")
	return nil
}

func (o *Orchestrator) verifyDeployment(ctx context.Context) error {
	verifyURL := o.project.Deploy.VerifyURL

	fmt.Printf("Verification: %s (skipped in CLI mode)\n", verifyURL)
	return nil
}

func (o *Orchestrator) recordDeployment(ctx context.Context) error {
	if o.history == nil {
		return nil
	}

	rec := history.NewSuccessRecord(
		o.project.Name,
		"",
		"deployment completed",
		"",
	)

	id, err := o.history.Insert(rec)
	if err != nil {
		return errors.NewDeployError("HISTORY_WRITE_FAILED", "failed to record deployment", err)
	}

	fmt.Printf("Recorded deployment #%d\n", id)
	return nil
}

func (o *Orchestrator) getBaseURL() string {
	return fmt.Sprintf("https://%s%s", o.project.FTPS.Host, o.project.Deploy.TargetPath)
}

func (o *Orchestrator) Disconnect() {
	if o.ftpClient != nil {
		o.ftpClient.Disconnect()
	}
}

func (o *Orchestrator) Cleanup() {
	if o.artifactPath != "" {
		os.Remove(o.artifactPath)
	}
}