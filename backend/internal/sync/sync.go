// SPDX-License-Identifier: MIT
package sync

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ENIACSystems/FileENIAC/backend/internal/diff"
	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/mirror"
	"github.com/ENIACSystems/FileENIAC/backend/internal/registry"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Manifest struct {
	ManifestID    string           `json:"manifest_id"`
	ProjectID     int64            `json:"project_id"`
	WorkspaceID   string           `json:"workspace_id"`
	GitCommit     string           `json:"git_commit,omitempty"`
	DeployID      string           `json:"deploy_id,omitempty"`
	Timestamp     string           `json:"timestamp"`
	OperationType string           `json:"operation_type"`
	Files         []*diff.FileDiff `json:"files"`
	Result        string           `json:"result"`
}

type Suggestion struct {
	FileCount   int    `json:"file_count"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

type Engine struct{}

func New() *Engine {
	return &Engine{}
}

// Plan analyzes differences and returns a sync suggestion.
// Sync is always assisted â€” never automatic.
func (e *Engine) Plan(diffReport *diff.Report) *Suggestion {
	if diffReport == nil {
		return &Suggestion{Action: "none", Description: "Nenhum dado para analisar"}
	}

	s := &Suggestion{
		FileCount: diffReport.Summary.Modified + diffReport.Summary.New + diffReport.Summary.Removed,
	}

	switch {
	case s.FileCount == 0:
		s.Action = "none"
		s.Description = "Local e mirror estÃ£o sincronizados. Nenhuma aÃ§Ã£o necessÃ¡ria."
	case diffReport.SourceA == "local" && diffReport.SourceB == "mirror":
		s.Action = "mirror_update"
		s.Description = fmt.Sprintf("%d arquivos diferentes entre local e mirror. O mirror precisa ser atualizado.", s.FileCount)
	default:
		s.Action = "review"
		s.Description = fmt.Sprintf("%d divergÃªncias detectadas. RevisÃ£o manual necessÃ¡ria.", s.FileCount)
	}

	return s
}

// Apply executes file copy/delete operations to reconcile local and mirror directories.
// direction must be "local_to_mirror" or "mirror_to_local".
// If confirm is false, destructive operations (deletions) return an error without executing.
func (e *Engine) Apply(ctx *workspace.Context, projectName string, diffReport *diff.Report, direction string, confirm bool) error {
	proj, err := registry.GetProject(ctx, projectName)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	mirrorDir := mirror.MirrorPath(ctx.Workspace.Path, projectName)
	localDir := proj.LocalPath

	var sourceDir, targetDir string
	switch direction {
	case "local_to_mirror":
		sourceDir = localDir
		targetDir = mirrorDir
	case "mirror_to_local":
		sourceDir = mirrorDir
		targetDir = localDir
	default:
		return fmt.Errorf("invalid direction: %q, must be local_to_mirror or mirror_to_local", direction)
	}

	if !confirm {
		for _, f := range diffReport.Files {
			if needsDelete(f, direction) {
				return fmt.Errorf("confirmation required: %s would be deleted", f.Path)
			}
		}
	}

	for _, f := range diffReport.Files {
		switch {
		case f.Status == diff.StateSynced:
			continue
		case needsDelete(f, direction):
			dstPath := filepath.Join(targetDir, f.Path)
			if err := os.Remove(dstPath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove %s: %w", f.Path, err)
			}
		default:
			srcPath := filepath.Join(sourceDir, f.Path)
			dstPath := filepath.Join(targetDir, f.Path)
			if err := os.MkdirAll(filepath.Dir(dstPath), 0700); err != nil {
				return fmt.Errorf("failed to create parent dir for %s: %w", f.Path, err)
			}
			data, err := os.ReadFile(srcPath)
			if err != nil {
				return fmt.Errorf("failed to read source %s: %w", f.Path, err)
			}
			if err := os.WriteFile(dstPath, data, 0600); err != nil {
				return fmt.Errorf("failed to write target %s: %w", f.Path, err)
			}
		}
	}

	log.L().Info("sync apply completed",
		zap.String("project", projectName),
		zap.String("direction", direction),
	)
	return nil
}

// needsDelete returns true when the diff entry represents a file that must be
// deleted from the target directory for the given sync direction.
func needsDelete(f *diff.FileDiff, direction string) bool {
	if direction == "local_to_mirror" {
		return f.Status == diff.StateRemoved
	}
	return f.Status == diff.StateRemoved
}

// GenerateManifest creates a sync-manifest.json record for the operation.
func (e *Engine) GenerateManifest(ctx *workspace.Context, projectName, opType string, diffReport *diff.Report, result string) (*Manifest, error) {
	proj, err := registry.GetProject(ctx, projectName)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	manifestID := uuid.New().String()[:12]

	manifest := &Manifest{
		ManifestID:    manifestID,
		ProjectID:     proj.ID,
		WorkspaceID:   ctx.Workspace.Name,
		Timestamp:     time.Now().Format(time.RFC3339),
		OperationType: opType,
		Result:        result,
	}

	if diffReport != nil {
		manifest.Files = diffReport.Files
		manifest.GitCommit = proj.LastCommitHash
		manifest.DeployID = proj.LastDeployID
	}

	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal manifest: %w", err)
	}

	_, err = ctx.DB.Exec(
		`INSERT INTO sync_manifests (project_id, manifest_id, operation_type, git_commit, deploy_id, files_count, result, manifest_json) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		proj.ID, manifestID, opType, manifest.GitCommit, manifest.DeployID, len(manifest.Files), result, string(manifestJSON),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to record manifest: %w", err)
	}

	log.L().Info("sync manifest generated",
		zap.String("project", projectName),
		zap.String("manifest_id", manifestID),
		zap.String("op_type", opType),
		zap.String("result", result),
	)

	return manifest, nil
}

// Reconcile updates the project state after a sync operation.
// It only updates the registry, never modifies project files.
func (e *Engine) Reconcile(ctx *workspace.Context, projectName, divergence string) error {
	proj, err := registry.GetProject(ctx, projectName)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	err = registry.UpdateProjectState(ctx, proj.ID, proj.LastCommitHash, proj.LastDeployID, "", divergence)
	if err != nil {
		return fmt.Errorf("failed to reconcile: %w", err)
	}

	log.L().Info("sync reconciled",
		zap.String("project", projectName),
		zap.String("divergence", divergence),
	)

	return nil
}
