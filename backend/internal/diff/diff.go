// SPDX-License-Identifier: MIT
package diff

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/mirror"
	"github.com/ENIACSystems/FileENIAC/backend/internal/registry"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"go.uber.org/zap"
)

type FileState string

const (
	StateNew       FileState = "novo"
	StateModified  FileState = "modificado"
	StateRemoved   FileState = "removido"
	StateDivergent FileState = "divergente"
	StateSynced    FileState = "sincronizado"
	StateUnknown   FileState = "desconhecido"
)

type FileDiff struct {
	Path       string    `json:"path"`
	LocalHash  string    `json:"local_hash,omitempty"`
	RemoteHash string    `json:"remote_hash,omitempty"`
	Status     FileState `json:"status"`
	Source     string    `json:"source"`
}

type Report struct {
	ProjectName string      `json:"project_name"`
	SourceA     string      `json:"source_a"`
	SourceB     string      `json:"source_b"`
	Files       []*FileDiff `json:"files"`
	Summary     struct {
		Total    int `json:"total"`
		New      int `json:"new"`
		Modified int `json:"modified"`
		Removed  int `json:"removed"`
		Synced   int `json:"synced"`
	} `json:"summary"`
}

type Engine struct{}

func New() *Engine {
	return &Engine{}
}

// fileHash computes SHA-256 of a file.
func fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// LocalVsMirror compares the local project directory against its mirror.
func (e *Engine) LocalVsMirror(ctx *workspace.Context, projectName string) (*Report, error) {
	proj, err := registry.GetProject(ctx, projectName)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	report := &Report{
		ProjectName: projectName,
		SourceA:     "local",
		SourceB:     "mirror",
	}

	mirrorDir := mirror.MirrorPath(ctx.Workspace.Path, projectName)
	localDir := proj.LocalPath

	localFiles := make(map[string]string)
	mirrorFiles := make(map[string]string)

	// Index local files
	filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.L().Warn("diff: walk error on local file", zap.String("path", path), zap.Error(err))
			return nil
		}
		if info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(localDir, path)
		if strings.HasPrefix(rel, ".git") || strings.HasPrefix(rel, ".eniac") {
			return nil
		}
		hash, err := fileHash(path)
		if err != nil {
			log.L().Warn("diff: failed to hash local file", zap.String("path", path), zap.Error(err))
			localFiles[rel] = "unreadable"
		} else {
			localFiles[rel] = hash
		}
		return nil
	})

	// Index mirror files
	filepath.Walk(mirrorDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.L().Warn("diff: walk error on mirror file", zap.String("path", path), zap.Error(err))
			return nil
		}
		if info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(mirrorDir, path)
		hash, err := fileHash(path)
		if err != nil {
			log.L().Warn("diff: failed to hash mirror file", zap.String("path", path), zap.Error(err))
			mirrorFiles[rel] = "unreadable"
		} else {
			mirrorFiles[rel] = hash
		}
		return nil
	})

	allFiles := make(map[string]bool)
	for k := range localFiles {
		allFiles[k] = true
	}
	for k := range mirrorFiles {
		allFiles[k] = true
	}

	for f := range allFiles {
		lh, lok := localFiles[f]
		mh, mok := mirrorFiles[f]
		fd := &FileDiff{Path: f}

		if lok && !mok {
			fd.Status = StateNew
			fd.LocalHash = lh
			fd.Source = "local"
		} else if !lok && mok {
			fd.Status = StateRemoved
			fd.RemoteHash = mh
			fd.Source = "mirror"
		} else if lok && mok && lh != mh {
			fd.Status = StateModified
			fd.LocalHash = lh
			fd.RemoteHash = mh
			fd.Source = "both"
		} else {
			fd.Status = StateSynced
			fd.LocalHash = lh
		}

		report.Files = append(report.Files, fd)

		switch fd.Status {
		case StateNew:
			report.Summary.New++
		case StateModified:
			report.Summary.Modified++
		case StateRemoved:
			report.Summary.Removed++
		case StateSynced:
			report.Summary.Synced++
		}
	}
	report.Summary.Total = len(report.Files)

	log.L().Info("diff local vs mirror",
		zap.String("project", projectName),
		zap.Int("total", report.Summary.Total),
		zap.Int("modified", report.Summary.Modified),
	)

	return report, nil
}

func (e *Engine) Status(ctx *workspace.Context, projectName string) (string, error) {
	report, err := e.LocalVsMirror(ctx, projectName)
	if err != nil {
		return "unknown", err
	}
	if report.Summary.Modified > 0 || report.Summary.New > 0 || report.Summary.Removed > 0 {
		return "divergente", nil
	}
	return "sincronizado", nil
}
