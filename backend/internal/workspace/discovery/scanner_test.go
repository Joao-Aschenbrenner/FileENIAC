package discovery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
)

func TestScanDir_Found(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")

	_, err := workspace.Init("ScanTest", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	workspace.Active().DB.Close()

	s := New()
	result := s.ScanDir(wsPath)
	if !result.Found {
		t.Error("expected workspace to be found")
	}
	if result.Workspace != "ScanTest" {
		t.Errorf("expected ScanTest, got %s", result.Workspace)
	}
}

func TestScanDir_NotFound(t *testing.T) {
	tmpDir := t.TempDir()

	s := New()
	result := s.ScanDir(tmpDir)
	if result.Found {
		t.Error("expected workspace not to be found")
	}
}

func TestScanDir_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, ".eniac")
	os.MkdirAll(wsDir, 0700)

	s := New()
	// Has .eniac/ but no config â€” Open will fail
	result := s.ScanDir(tmpDir)
	if !result.Found {
		t.Error("expected .eniac/ to be detected as found")
	}
	if result.Error == "" {
		t.Error("expected error about failed open")
	}
}

func TestScanDeep_FindsNested(t *testing.T) {
	tmpDir := t.TempDir()

	// Create two workspaces at different levels
	ws1 := filepath.Join(tmpDir, "projects", "ws1")
	ws2 := filepath.Join(tmpDir, "other", "ws2")

	_, err := workspace.Init("Deep1", ws1, "")
	if err != nil {
		t.Fatalf("Init ws1 failed: %v", err)
	}
	workspace.Active().DB.Close()

	_, err = workspace.Init("Deep2", ws2, "")
	if err != nil {
		t.Fatalf("Init ws2 failed: %v", err)
	}
	workspace.Active().DB.Close()

	s := New()
	results, err := s.ScanDeep(tmpDir, 5)
	if err != nil {
		t.Fatalf("ScanDeep failed: %v", err)
	}

	count := 0
	for _, r := range results {
		if r.Found {
			count++
		}
	}
	if count < 2 {
		t.Errorf("expected at least 2 workspaces, got %d", count)
	}
}
