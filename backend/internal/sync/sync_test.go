package sync

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ENIACSystems/FileENIAC/backend/internal/diff"
	"github.com/ENIACSystems/FileENIAC/backend/internal/mirror"
	"github.com/ENIACSystems/FileENIAC/backend/internal/registry"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
)

func TestPlan_NoDiff(t *testing.T) {
	e := New()
	s := e.Plan(&diff.Report{
		ProjectName: "test",
		SourceA:     "local",
		SourceB:     "mirror",
	})
	if s.Action != "none" {
		t.Errorf("expected none, got %s", s.Action)
	}
}

func TestPlan_HasDiff(t *testing.T) {
	e := New()
	s := e.Plan(&diff.Report{
		ProjectName: "test",
		SourceA:     "local",
		SourceB:     "mirror",
		Summary: struct {
			Total    int "json:\"total\""
			New      int "json:\"new\""
			Modified int "json:\"modified\""
			Removed  int "json:\"removed\""
			Synced   int "json:\"synced\""
		}{Total: 3, New: 1, Modified: 2, Removed: 0, Synced: 0},
	})
	if s.Action != "mirror_update" {
		t.Errorf("expected mirror_update, got %s", s.Action)
	}
	if s.FileCount != 3 {
		t.Errorf("expected 3 files, got %d", s.FileCount)
	}
}

func TestPlan_NilReport(t *testing.T) {
	e := New()
	s := e.Plan(nil)
	if s.Action != "none" {
		t.Errorf("expected none, got %s", s.Action)
	}
}

func TestGenerateManifest(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	_, err := workspace.Init("SyncTest", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "sync-project",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})
	registry.UpdateProjectState(workspace.Active(), 1, "abc123", "dep001", "", "divergente")

	// Create dummy diff report
	report := &diff.Report{
		ProjectName: "sync-project",
		SourceA:     "local",
		SourceB:     "mirror",
	}

	e := New()
	manifest, err := e.GenerateManifest(workspace.Active(), "sync-project", "sync_apply", report, "completed")
	if err != nil {
		t.Fatalf("GenerateManifest failed: %v", err)
	}

	if manifest.ManifestID == "" {
		t.Error("expected non-empty manifest ID")
	}
	if manifest.ProjectID != 1 {
		t.Errorf("expected project 1, got %d", manifest.ProjectID)
	}
	if manifest.OperationType != "sync_apply" {
		t.Errorf("expected sync_apply, got %s", manifest.OperationType)
	}
	if manifest.GitCommit != "abc123" {
		t.Errorf("expected abc123, got %s", manifest.GitCommit)
	}
}

func TestReconcile(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	_, err := workspace.Init("ReconcileTest", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "reconcile-project",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})

	e := New()
	err = e.Reconcile(workspace.Active(), "reconcile-project", "sincronizado")
	if err != nil {
		t.Fatalf("Reconcile failed: %v", err)
	}

	proj, err := registry.GetProject(workspace.Active(), "reconcile-project")
	if err != nil {
		t.Fatalf("GetProject failed: %v", err)
	}

	if proj.DivergenceStatus != "sincronizado" {
		t.Errorf("expected sincronizado, got %s", proj.DivergenceStatus)
	}
}

func TestReconcile_BlocksDestructive_NoConfirm(t *testing.T) {
	// Verify that reconcile only updates metadata â€” never modifies project files
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	originalContent := []byte("original content")
	os.WriteFile(filepath.Join(projectPath, "important.php"), originalContent, 0644)
	mirrorDir := mirror.MirrorPath(wsPath, "safe-project")
	os.MkdirAll(mirrorDir, 0700)
	os.WriteFile(filepath.Join(mirrorDir, "important.php"), []byte("mirror content"), 0644)

	_, err := workspace.Init("SafeTest", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "safe-project",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})

	// Reconcile only updates metadata â€” should NOT touch project files
	eSync := New()
	eSync.Reconcile(workspace.Active(), "safe-project", "sincronizado")

	// Verify project file is unchanged
	currentContent, _ := os.ReadFile(filepath.Join(projectPath, "important.php"))
	if string(currentContent) != string(originalContent) {
		t.Error("reconcile should never modify project files")
	}
}
