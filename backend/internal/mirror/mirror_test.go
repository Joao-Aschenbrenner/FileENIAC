package mirror

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ENIACSystems/FileENIAC/backend/internal/registry"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
)

func TestMirrorPath(t *testing.T) {
	path := MirrorPath("/home/ws", "MyProject")
	expected := filepath.FromSlash("/home/ws/.eniac/mirror/MyProject")
	if path != expected {
		t.Errorf("expected %s, got %s", expected, path)
	}
}

func TestCreate_MirrorDirCreated(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectPath, 0700)

	ws, err := workspace.Init("MirrorTest", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	projID, err := registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "mirror-project",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}

	// Create a server entry (no real FTPS â€” Create will fail at connect, but mirror dir should be attempted)
	_, err = registry.AddServer(workspace.Active(), &registry.Server{
		ProjectID:  projID,
		Name:       "production",
		Type:       "ftps",
		Host:       "192.0.2.1", // TEST-NET â€” will fail
		Port:       21,
		User:       "test",
		Password:   "test",
		TargetPath: "/remote",
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("AddServer failed: %v", err)
	}

	e := New()
	_, err = e.Create(workspace.Active(), "mirror-project")
	if err == nil {
		t.Fatal("expected error (no real FTPS server)")
	}

	// Mirror dir should have been created despite connection failure
	mirrorPath := MirrorPath(ws.Path, "mirror-project")
	if _, statErr := os.Stat(mirrorPath); os.IsNotExist(statErr) {
		t.Error("mirror directory should exist even if download fails")
	}

	workspace.Active().DB.Close()
}

func TestStatus_NoSnapshot(t *testing.T) {
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
		Name:       "no-snapshot",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})

	e := New()
	_, err = e.Status(workspace.Active(), "no-snapshot")
	if err == nil {
		t.Error("expected error for project with no snapshot")
	}
}
