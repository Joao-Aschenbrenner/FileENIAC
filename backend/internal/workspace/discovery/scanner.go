package discovery

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/registry"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"go.uber.org/zap"
)

type ScanResult struct {
	Path      string              `json:"path"`
	Found     bool                `json:"found"`
	Workspace string              `json:"workspace,omitempty"`
	Error     string              `json:"error,omitempty"`
	Projects  []*registry.Project `json:"projects,omitempty"`
	Servers   []*registry.Server  `json:"servers,omitempty"`
}

type Scanner struct{}

func New() *Scanner {
	return &Scanner{}
}

// ScanDir checks if the given directory contains an FileENIAC.
// It does not modify any files â€” read-only.
func (s *Scanner) ScanDir(path string) *ScanResult {
	result := &ScanResult{Path: path}

	wsDir := filepath.Join(path, ".eniac")
	if _, err := os.Stat(wsDir); os.IsNotExist(err) {
		result.Found = false
		return result
	}

	result.Found = true

	ws, err := workspace.Open(path)
	if err != nil {
		result.Error = fmt.Sprintf("failed to open workspace: %v", err)
		log.L().Warn("discovery scan failed to open",
			zap.String("path", path),
			zap.Error(err),
		)
		return result
	}

	result.Workspace = ws.Name

	projects, err := registry.ListProjects(workspace.Active())
	if err != nil {
		result.Error = fmt.Sprintf("failed to list projects: %v", err)
		workspace.Active().DB.Close()
		return result
	}
	result.Projects = projects

	for _, p := range projects {
		srv, err := registry.GetServer(workspace.Active(), p.ID)
		if err == nil {
			result.Servers = append(result.Servers, srv)
		}
	}

	log.L().Info("discovery found workspace",
		zap.String("ws", result.Workspace),
		zap.String("path", path),
		zap.Int("projects", len(result.Projects)),
	)

	workspace.Active().DB.Close()

	return result
}

// ScanDeep walks the given root directory tree looking for .eniac/ directories.
// Returns all workspaces found without modifying any files.
func (s *Scanner) ScanDeep(root string, maxDepth int) ([]*ScanResult, error) {
	var results []*ScanResult

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip inaccessible dirs
		}
		if !d.IsDir() {
			return nil
		}
		if d.Name() == ".eniac" {
			parent := filepath.Dir(path)
			result := s.ScanDir(parent)
			results = append(results, result)
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		return results, fmt.Errorf("scan deep failed: %w", err)
	}

	return results, nil
}
