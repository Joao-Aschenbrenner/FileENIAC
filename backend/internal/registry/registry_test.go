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
