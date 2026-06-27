package diff

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ENIACSystems/FileENIAC/backend/internal/mirror"
	"github.com/ENIACSystems/FileENIAC/backend/internal/registry"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
)

func TestLocalVsMirror_Synced(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	_, err := workspace.Init("DiffTest", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "diff-project",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})

	// Create same file in both local and mirror
	os.WriteFile(filepath.Join(projectPath, "index.php"), []byte("same"), 0644)
	mirrorDir := mirror.MirrorPath(wsPath, "diff-project")
	os.MkdirAll(mirrorDir, 0700)
	os.WriteFile(filepath.Join(mirrorDir, "index.php"), []byte("same"), 0644)

	e := New()
	report, err := e.LocalVsMirror(workspace.Active(), "diff-project")
	if err != nil {
		t.Fatalf("LocalVsMirror failed: %v", err)
	}

	if len(report.Files) == 0 {
		t.Fatal("expected at least 1 file in report")
	}

	found := false
	for _, f := range report.Files {
		if f.Path == "index.php" && f.Status == StateSynced {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected index.php to be synced")
	}
}

func TestLocalVsMirror_Modified(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	_, err := workspace.Init("DiffTest2", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "diff-project2",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})

	os.WriteFile(filepath.Join(projectPath, "changed.php"), []byte("local version"), 0644)
	mirrorDir := mirror.MirrorPath(wsPath, "diff-project2")
	os.MkdirAll(mirrorDir, 0700)
	os.WriteFile(filepath.Join(mirrorDir, "changed.php"), []byte("mirror version"), 0644)

	e := New()
	report, err := e.LocalVsMirror(workspace.Active(), "diff-project2")
	if err != nil {
		t.Fatalf("LocalVsMirror failed: %v", err)
	}

	found := false
	for _, f := range report.Files {
		if f.Path == "changed.php" && f.Status == StateModified {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected changed.php to be modified")
	}
	if report.Summary.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", report.Summary.Modified)
	}
}

func TestLocalVsMirror_NewFile(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	_, err := workspace.Init("DiffNew", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "new-project",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})

	os.WriteFile(filepath.Join(projectPath, "new.php"), []byte("new file"), 0644)

	e := New()
	report, err := e.LocalVsMirror(workspace.Active(), "new-project")
	if err != nil {
		t.Fatalf("LocalVsMirror failed: %v", err)
	}

	if len(report.Files) == 0 {
		t.Fatal("expected at least 1 file in report")
	}

	found := false
	for _, f := range report.Files {
		if f.Path == "new.php" && f.Status == StateNew {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected new.php to be StateNew")
	}
	if report.Summary.New != 1 {
		t.Errorf("expected 1 new, got %d", report.Summary.New)
	}
}

func TestLocalVsMirror_RemovedFile(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	_, err := workspace.Init("DiffRemoved", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "rem-project",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})

	mirrorDir := mirror.MirrorPath(wsPath, "rem-project")
	os.MkdirAll(mirrorDir, 0700)
	os.WriteFile(filepath.Join(mirrorDir, "removed.php"), []byte("only in mirror"), 0644)

	e := New()
	report, err := e.LocalVsMirror(workspace.Active(), "rem-project")
	if err != nil {
		t.Fatalf("LocalVsMirror failed: %v", err)
	}

	if len(report.Files) == 0 {
		t.Fatal("expected at least 1 file in report")
	}

	found := false
	for _, f := range report.Files {
		if f.Path == "removed.php" && f.Status == StateRemoved {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected removed.php to be StateRemoved")
	}
	if report.Summary.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", report.Summary.Removed)
	}
}

func TestLocalVsMirror_Subdirs(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	_, err := workspace.Init("DiffSubdir", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "sub-project",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})

	os.MkdirAll(filepath.Join(projectPath, "sub"), 0700)
	os.WriteFile(filepath.Join(projectPath, "sub", "deep.txt"), []byte("deep"), 0644)
	os.WriteFile(filepath.Join(projectPath, "root.txt"), []byte("root"), 0644)
	mirrorDir := mirror.MirrorPath(wsPath, "sub-project")
	os.MkdirAll(mirrorDir, 0700)
	os.WriteFile(filepath.Join(mirrorDir, "root.txt"), []byte("root"), 0644)

	e := New()
	report, err := e.LocalVsMirror(workspace.Active(), "sub-project")
	if err != nil {
		t.Fatalf("LocalVsMirror failed: %v", err)
	}

	if report.Summary.Total != 2 {
		t.Errorf("expected 2 files total, got %d", report.Summary.Total)
	}
	if report.Summary.New != 1 {
		t.Errorf("expected 1 new (deep.txt), got %d", report.Summary.New)
	}
	if report.Summary.Synced != 1 {
		t.Errorf("expected 1 synced (root.txt), got %d", report.Summary.Synced)
	}
}

func TestStatus_Synced(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	_, err := workspace.Init("StatusTest", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "status-project",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})

	content := []byte("same content")
	os.WriteFile(filepath.Join(projectPath, "a.txt"), content, 0644)
	mirrorDir := mirror.MirrorPath(wsPath, "status-project")
	os.MkdirAll(mirrorDir, 0700)
	os.WriteFile(filepath.Join(mirrorDir, "a.txt"), content, 0644)

	e := New()
	status, err := e.Status(workspace.Active(), "status-project")
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}
	if status != "sincronizado" {
		t.Errorf("expected sincronizado, got %s", status)
	}
}

func TestLocalVsMirror_EmptyTree(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	_, err := workspace.Init("DiffEmpty", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "empty-project",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})

	// Both local and mirror empty — no files
	mirrorDir := mirror.MirrorPath(wsPath, "empty-project")
	os.MkdirAll(mirrorDir, 0700)

	e := New()
	report, err := e.LocalVsMirror(workspace.Active(), "empty-project")
	if err != nil {
		t.Fatalf("LocalVsMirror failed: %v", err)
	}

	if len(report.Files) != 0 {
		t.Errorf("expected 0 files for empty tree, got %d", len(report.Files))
	}
	if report.Summary.Total != 0 {
		t.Errorf("expected Total=0, got %d", report.Summary.Total)
	}
}

func TestStatus_Divergent(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	_, err := workspace.Init("StatusTest2", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "div-project",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})

	os.WriteFile(filepath.Join(projectPath, "only_local.txt"), []byte("only local"), 0644)

	e := New()
	status, err := e.Status(workspace.Active(), "div-project")
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}
	if status != "divergente" {
		t.Errorf("expected divergente, got %s", status)
	}
}
