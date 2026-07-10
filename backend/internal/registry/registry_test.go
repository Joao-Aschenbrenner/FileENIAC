// SPDX-License-Identifier: MIT
package registry

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
)

func setupTestWorkspace(t *testing.T) (*workspace.Context, string) {
	t.Helper()
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")

	_, err := workspace.Init("TestWS", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	_, err = workspace.Open(wsPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	ctx := workspace.Active()
	if ctx == nil {
		t.Fatal("context should not be nil")
	}

	t.Cleanup(func() { _ = ctx.DB.Close() })

	return ctx, wsPath
}

func TestAddProject(t *testing.T) {
	ctx, _ := setupTestWorkspace(t)

	p := &Project{
		Name:        "test-project",
		LocalPath:   os.TempDir(),
		RemotePath:  "/remote/test",
		Branch:      "main",
		Environment: "production",
		IsActive:    true,
	}

	id, err := AddProject(ctx, p)
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}

	if id <= 0 {
		t.Error("expected positive ID")
	}
}

func TestListProjects(t *testing.T) {
	ctx, _ := setupTestWorkspace(t)

	p := &Project{
		Name:       "project-a",
		LocalPath:  os.TempDir(),
		RemotePath: "/remote/a",
		IsActive:   true,
	}

	_, err := AddProject(ctx, p)
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}

	projects, err := ListProjects(ctx)
	if err != nil {
		t.Fatalf("ListProjects failed: %v", err)
	}

	if len(projects) != 1 {
		t.Errorf("expected 1 project, got %d", len(projects))
	}
}

func TestGetProject(t *testing.T) {
	ctx, _ := setupTestWorkspace(t)

	p := &Project{
		Name:       "find-me",
		LocalPath:  os.TempDir(),
		RemotePath: "/remote/find",
		IsActive:   true,
	}

	_, err := AddProject(ctx, p)
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}

	found, err := GetProject(ctx, "find-me")
	if err != nil {
		t.Fatalf("GetProject failed: %v", err)
	}

	if found.Name != "find-me" {
		t.Errorf("expected find-me, got %s", found.Name)
	}
}

func TestRemoveProject(t *testing.T) {
	ctx, _ := setupTestWorkspace(t)

	p := &Project{
		Name:       "to-remove",
		LocalPath:  os.TempDir(),
		RemotePath: "/remote/rm",
		IsActive:   true,
	}

	id, err := AddProject(ctx, p)
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}

	if err := RemoveProject(ctx, id); err != nil {
		t.Fatalf("RemoveProject failed: %v", err)
	}

	_, err = GetProject(ctx, "to-remove")
	if err == nil {
		t.Error("expected error after removal")
	}
}

func TestAddServer(t *testing.T) {
	ctx, _ := setupTestWorkspace(t)

	p := &Project{
		Name:       "server-project",
		LocalPath:  os.TempDir(),
		RemotePath: "/remote/srv",
		IsActive:   true,
	}

	projID, err := AddProject(ctx, p)
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}

	s := &Server{
		ProjectID:  projID,
		Name:       "production",
		Type:       "ftps",
		Host:       "ftp.example.com",
		Port:       21,
		User:       "user",
		Password:   "pass",
		TargetPath: "/public_html",
		IsActive:   true,
	}

	id, err := AddServer(ctx, s)
	if err != nil {
		t.Fatalf("AddServer failed: %v", err)
	}

	if id <= 0 {
		t.Error("expected positive ID")
	}

	_, err = GetServer(ctx, projID)
	if err != nil {
		t.Fatalf("GetServer failed: %v", err)
	}
}

func TestCanDeleteLocalPath_InsideWorkspace(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	if err := os.MkdirAll(wsPath, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	projectDir := filepath.Join(wsPath, "my-project")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("mkdir project: %v", err)
	}

	if err := CanDeleteLocalPath(projectDir, wsPath); err != nil {
		t.Errorf("expected safe path, got: %v", err)
	}
}

func TestCanDeleteLocalPath_OutsideWorkspace(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	if err := os.MkdirAll(wsPath, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	outsideDir := filepath.Join(tmpDir, "outside")
	if err := os.MkdirAll(outsideDir, 0755); err != nil {
		t.Fatalf("mkdir outside: %v", err)
	}

	if err := CanDeleteLocalPath(outsideDir, wsPath); err == nil {
		t.Error("expected error for path outside workspace")
	}
}

func TestCanDeleteLocalPath_ReservedDir(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	if err := os.MkdirAll(wsPath, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	for _, reserved := range []string{".git", ".github", ".eniac", "node_modules"} {
		d := filepath.Join(wsPath, reserved)
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatalf("mkdir %s: %v", reserved, err)
		}
		if err := CanDeleteLocalPath(d, wsPath); err == nil {
			t.Errorf("expected error for reserved dir %s", reserved)
		}
	}
}

func TestCanDeleteLocalPath_Empty(t *testing.T) {
	if err := CanDeleteLocalPath("", ""); err == nil {
		t.Error("expected error for empty path")
	}
}

func TestCanDeleteLocalPath_NonExistentOK(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	if err := os.MkdirAll(wsPath, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	missing := filepath.Join(wsPath, "does-not-exist")
	if err := CanDeleteLocalPath(missing, wsPath); err != nil {
		t.Errorf("expected no error for non-existent path inside workspace, got: %v", err)
	}
}

func TestDeleteLocalPath_InsideWorkspace(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	if err := os.MkdirAll(wsPath, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	projectDir := filepath.Join(wsPath, "to-delete")
	if err := os.MkdirAll(filepath.Join(projectDir, "sub"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, "file.txt"), []byte("x"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	if err := DeleteLocalPath(projectDir, wsPath); err != nil {
		t.Fatalf("DeleteLocalPath failed: %v", err)
	}
	if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
		t.Errorf("expected folder deleted, got stat err: %v", err)
	}
}

func TestDeleteLocalPath_OutsideWorkspaceBlocked(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	if err := os.MkdirAll(wsPath, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	outside := filepath.Join(tmpDir, "outside")
	if err := os.MkdirAll(outside, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := DeleteLocalPath(outside, wsPath); err == nil {
		t.Error("expected error deleting path outside workspace")
	}
	if _, err := os.Stat(outside); err != nil {
		t.Errorf("outside folder should still exist: %v", err)
	}
}
