package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitWorkspace(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")

	ws, err := Init("TestWS", wsPath, "Test workspace")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	//goland:noinspection GoDeferInLoop
	defer activeContext.DB.Close()

	if ws.Name != "TestWS" {
		t.Errorf("expected TestWS, got %s", ws.Name)
	}

	if _, err := os.Stat(filepath.Join(wsPath, ".eniac", "config.toml")); os.IsNotExist(err) {
		t.Error("config.toml should exist")
	}

	if _, err := os.Stat(filepath.Join(wsPath, ".eniac", "workspace.db")); os.IsNotExist(err) {
		t.Error("workspace.db should exist")
	}
}

func TestOpenWorkspace(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")

	_, err := Init("TestWS", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	activeContext.DB.Close()

	ws, err := Open(wsPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer activeContext.DB.Close()

	if ws.Name != "TestWS" {
		t.Errorf("expected TestWS, got %s", ws.Name)
	}

	if Active() == nil {
		t.Error("active context should not be nil")
	}
}

func TestWorkspaceStatus(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")

	ws, err := Init("TestWS", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer activeContext.DB.Close()

	status := ws.Status()
	if status == nil {
		t.Fatal("status should not be nil")
	}

	if status["name"] != "TestWS" {
		t.Errorf("expected TestWS, got %v", status["name"])
	}
}
