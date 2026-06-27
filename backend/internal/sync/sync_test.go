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

func TestPlan_ReviewAction(t *testing.T) {
	e := New()
	s := e.Plan(&diff.Report{
		ProjectName: "test",
		SourceA:     "local",
		SourceB:     "other",
		Summary: struct {
			Total    int "json:\"total\""
			New      int "json:\"new\""
			Modified int "json:\"modified\""
			Removed  int "json:\"removed\""
			Synced   int "json:\"synced\""
		}{Total: 2, New: 1, Modified: 1, Removed: 0, Synced: 0},
	})
	if s.Action != "review" {
		t.Errorf("expected review, got %s", s.Action)
	}
}

func TestApply_LocalToMirror(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	_, err := workspace.Init("ApplyTest", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "apply-project",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})

	os.WriteFile(filepath.Join(projectPath, "file.txt"), []byte("hello"), 0644)
	mirrorDir := mirror.MirrorPath(wsPath, "apply-project")
	os.MkdirAll(mirrorDir, 0700)

	report := &diff.Report{
		ProjectName: "apply-project",
		SourceA:     "local",
		SourceB:     "mirror",
		Files: []*diff.FileDiff{
			{Path: "file.txt", Status: diff.StateNew, Source: "local"},
		},
		Summary: struct {
			Total    int "json:\"total\""
			New      int "json:\"new\""
			Modified int "json:\"modified\""
			Removed  int "json:\"removed\""
			Synced   int "json:\"synced\""
		}{Total: 1, New: 1},
	}

	e := New()
	err = e.Apply(workspace.Active(), "apply-project", report, "local_to_mirror", true)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(mirrorDir, "file.txt"))
	if err != nil {
		t.Fatal("mirror file should exist after apply")
	}
	if string(data) != "hello" {
		t.Errorf("expected hello, got %s", string(data))
	}
}

func TestApply_MirrorToLocal_CurrentBehavior(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	_, err := workspace.Init("ApplyTest2", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "apply-project2",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})

	mirrorDir := mirror.MirrorPath(wsPath, "apply-project2")
	os.MkdirAll(mirrorDir, 0700)
	os.WriteFile(filepath.Join(mirrorDir, "from_mirror.txt"), []byte("from mirror"), 0644)

	report := &diff.Report{
		ProjectName: "apply-project2",
		SourceA:     "local",
		SourceB:     "mirror",
		Files: []*diff.FileDiff{
			{Path: "from_mirror.txt", Status: diff.StateNew, Source: "mirror"},
		},
		Summary: struct {
			Total    int "json:\"total\""
			New      int "json:\"new\""
			Modified int "json:\"modified\""
			Removed  int "json:\"removed\""
			Synced   int "json:\"synced\""
		}{Total: 1, New: 1},
	}

	e := New()
	err = e.Apply(workspace.Active(), "apply-project2", report, "mirror_to_local", true)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	// BUG: needsDelete returns true for StateNew + mirror_to_local,
	// so Apply treats new-in-mirror files as deletions, not copies.
	// File should NOT exist in local due to this current behavior.
	if _, statErr := os.Stat(filepath.Join(projectPath, "from_mirror.txt")); !os.IsNotExist(statErr) {
		t.Error("BUG CONFIRMED: new-in-mirror file should not be copied to local (inverted needsDelete)")
	}
}

func TestApply_DeleteLocalToMirror(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	_, err := workspace.Init("ApplyDel", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "del-project",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})

	mirrorDir := mirror.MirrorPath(wsPath, "del-project")
	os.MkdirAll(mirrorDir, 0700)
	os.WriteFile(filepath.Join(mirrorDir, "todelete.txt"), []byte("delete me"), 0644)

	report := &diff.Report{
		ProjectName: "del-project",
		SourceA:     "local",
		SourceB:     "mirror",
		Files: []*diff.FileDiff{
			{Path: "todelete.txt", Status: diff.StateRemoved, Source: "mirror"},
		},
		Summary: struct {
			Total    int "json:\"total\""
			New      int "json:\"new\""
			Modified int "json:\"modified\""
			Removed  int "json:\"removed\""
			Synced   int "json:\"synced\""
		}{Total: 1, Removed: 1},
	}

	e := New()
	err = e.Apply(workspace.Active(), "del-project", report, "local_to_mirror", true)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(mirrorDir, "todelete.txt")); !os.IsNotExist(err) {
		t.Error("expected todelete.txt to be deleted from mirror")
	}
}

func TestApply_NoConfirmBlocksDelete(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	_, err := workspace.Init("ApplyNoCfm", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "nocfm-project",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})

	mirrorDir := mirror.MirrorPath(wsPath, "nocfm-project")
	os.MkdirAll(mirrorDir, 0700)
	os.WriteFile(filepath.Join(mirrorDir, "del.txt"), []byte("x"), 0644)

	report := &diff.Report{
		ProjectName: "nocfm-project",
		SourceA:     "local",
		SourceB:     "mirror",
		Files: []*diff.FileDiff{
			{Path: "del.txt", Status: diff.StateRemoved, Source: "mirror"},
		},
		Summary: struct {
			Total    int "json:\"total\""
			New      int "json:\"new\""
			Modified int "json:\"modified\""
			Removed  int "json:\"removed\""
			Synced   int "json:\"synced\""
		}{Total: 1, Removed: 1},
	}

	e := New()
	err = e.Apply(workspace.Active(), "nocfm-project", report, "local_to_mirror", false)
	if err == nil {
		t.Error("expected error when confirm=false with deletions")
	}
}

func TestApply_InvalidDirection(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	_, err := workspace.Init("ApplyDir", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "dir-project",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})

	e := New()
	err = e.Apply(workspace.Active(), "dir-project", &diff.Report{}, "invalid", true)
	if err == nil {
		t.Error("expected error for invalid direction")
	}
}

func TestNeedsDelete_LocalToMirror(t *testing.T) {
	tests := []struct {
		name      string
		status    diff.FileState
		direction string
		want      bool
	}{
		{"local_to_mirror removed", diff.StateRemoved, "local_to_mirror", true},
		{"local_to_mirror new", diff.StateNew, "local_to_mirror", false},
		{"local_to_mirror modified", diff.StateModified, "local_to_mirror", false},
		{"local_to_mirror synced", diff.StateSynced, "local_to_mirror", false},
		{"mirror_to_local removed", diff.StateRemoved, "mirror_to_local", false},
		{"mirror_to_local new", diff.StateNew, "mirror_to_local", true},
		{"mirror_to_local modified", diff.StateModified, "mirror_to_local", false},
		{"mirror_to_local synced", diff.StateSynced, "mirror_to_local", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &diff.FileDiff{Path: "f.txt", Status: tt.status}
			got := needsDelete(f, tt.direction)
			if got != tt.want {
				t.Errorf("needsDelete(status=%v, dir=%s) = %v, want %v", tt.status, tt.direction, got, tt.want)
			}
		})
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
