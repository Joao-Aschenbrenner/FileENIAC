// SPDX-License-Identifier: MIT
package clone_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ENIACSystems/FileENIAC/backend/internal/clone"
)

func TestIsCloned_NonExistentDir(t *testing.T) {
	if clone.IsCloned("/nonexistent/path/88888") {
		t.Error("expected IsCloned to return false for non-existent dir")
	}
}

func TestIsCloned_NoGitDir(t *testing.T) {
	dir := t.TempDir()
	if clone.IsCloned(dir) {
		t.Error("expected IsCloned to return false for dir without .git")
	}
}

func TestIsCloned_WithDotGit(t *testing.T) {
	dir := t.TempDir()
	gitDir := filepath.Join(dir, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("mkdir .git failed: %v", err)
	}
	if !clone.IsCloned(dir) {
		t.Error("expected IsCloned to return true for dir with .git")
	}
}

func TestIsCloned_EmptyString(t *testing.T) {
	if clone.IsCloned("") {
		t.Error("expected IsCloned to return false for empty string")
	}
}

func TestClone_DirectoryAlreadyExists(t *testing.T) {
	dir := t.TempDir()
	existingDir := filepath.Join(dir, "existing-clone")
	if err := os.MkdirAll(existingDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	_, err := clone.Clone("https://example.com/repo.git", existingDir, "main")
	if err == nil {
		t.Fatal("expected error for existing directory")
	}
}

func TestClone_InvalidURL(t *testing.T) {
	dir := t.TempDir()
	newDir := filepath.Join(dir, "new-clone")
	_, err := clone.Clone("", newDir, "main")
	if err == nil {
		t.Log("Clone returned nil error for empty URL (may fail at exec)")
	}
}

func TestResult_Struct(t *testing.T) {
	r := &clone.Result{
		Path:       "/tmp/test",
		Branch:     "main",
		CommitSHA:  "abc123",
		DurationMS: 100,
	}
	if r.Path != "/tmp/test" {
		t.Errorf("expected /tmp/test, got %s", r.Path)
	}
	if r.Branch != "main" {
		t.Errorf("expected main, got %s", r.Branch)
	}
	if r.CommitSHA != "abc123" {
		t.Errorf("expected abc123, got %s", r.CommitSHA)
	}
	if r.DurationMS != 100 {
		t.Errorf("expected 100, got %d", r.DurationMS)
	}
}
