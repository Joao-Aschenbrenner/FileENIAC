// SPDX-License-Identifier: MIT
package clone

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"go.uber.org/zap"
)

type Result struct {
	Path       string `json:"path"`
	Branch     string `json:"branch"`
	CommitSHA  string `json:"commit_sha,omitempty"`
	DurationMS int64  `json:"duration_ms"`
}

func Clone(ctx context.Context, repoURL, cloneDir, branch string) (*Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	start := time.Now()

	if _, err := os.Stat(cloneDir); !os.IsNotExist(err) {
		return nil, fmt.Errorf("directory already exists: %s", cloneDir)
	}

	parent := filepath.Dir(cloneDir)
	if err := os.MkdirAll(parent, 0755); err != nil {
		return nil, fmt.Errorf("failed to create parent dir: %w", err)
	}

	if strings.Contains(repoURL, " ") || strings.HasPrefix(repoURL, "-") {
		return nil, fmt.Errorf("invalid repository URL")
	}
	args := []string{"clone", "--depth", "1", "--branch", branch, repoURL, cloneDir}

	cloneCtx := ctx
	if _, ok := cloneCtx.Deadline(); !ok {
		var cancel context.CancelFunc
		cloneCtx, cancel = context.WithTimeout(ctx, 120*time.Second)
		defer cancel()
	}

	cmd := exec.CommandContext(cloneCtx, "git", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if os.RemoveAll(cloneDir) == nil {
			log.L().Info("removed partial clone dir",
				zap.String("dir", cloneDir),
				zap.String("repo", repoURL),
			)
		}
		if cloneCtx.Err() != nil {
			return nil, fmt.Errorf("git clone timed out after 120s: %s", repoURL)
		}
		return nil, fmt.Errorf("git clone failed: %s: %w", string(output), err)
	}

	commitSHA, _ := getCommitSHA(cloneDir)
	duration := time.Since(start).Milliseconds()

	log.L().Info("repo cloned",
		zap.String("url", repoURL),
		zap.String("dir", cloneDir),
		zap.String("branch", branch),
		zap.Int64("duration_ms", duration),
	)

	return &Result{
		Path:       cloneDir,
		Branch:     branch,
		CommitSHA:  commitSHA,
		DurationMS: duration,
	}, nil
}

func getCommitSHA(repoDir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = repoDir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	if len(out) < 40 {
		return "", fmt.Errorf("unexpected git output length %d", len(out))
	}
	return string(out[:40]), nil
}

func IsCloned(repoDir string) bool {
	_, err := os.Stat(filepath.Join(repoDir, ".git"))
	return err == nil
}
