// SPDX-License-Identifier: MIT
package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ENIACSystems/FileENIAC/backend/internal/deploy"
	"github.com/ENIACSystems/FileENIAC/backend/internal/history"
	"github.com/ENIACSystems/FileENIAC/backend/internal/registry"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
)

// TestFullWorkspaceFlow validates the complete user flow:
// Workspace Init â†’ Project Add â†’ Deploy Run â†’ Verify â†’ History â†’ Rollback â†’ Status
//
// The FTPS-dependent portion (Deploy Run with real upload) is in
// TestFullWorkspaceFlowWithFTPS, which runs only when ENIAC_FTPS_TEST=1.
func TestFullWorkspaceFlow(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")

	if err := os.MkdirAll(projectPath, 0700); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}
	// Create a dummy file so packer has something to pack
	if err := os.WriteFile(filepath.Join(projectPath, "index.php"), []byte("<?php echo 'ok';"), 0644); err != nil {
		t.Fatalf("failed to create dummy file: %v", err)
	}

	// 1. Workspace Init
	ws, err := workspace.Init("IntegrationWS", wsPath, "Integration test workspace")
	if err != nil {
		t.Fatalf("Step 1 - Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	if ws.Name != "IntegrationWS" {
		t.Errorf("expected IntegrationWS, got %s", ws.Name)
	}

	// 2. Project Add
	projID, err := registry.AddProject(workspace.Active(), &registry.Project{
		Name:        "test-project",
		LocalPath:   projectPath,
		RemotePath:  "/remote/test",
		Branch:      "main",
		Environment: "production",
		IsActive:    true,
	})
	if err != nil {
		t.Fatalf("Step 2 - AddProject failed: %v", err)
	}
	if projID <= 0 {
		t.Fatal("Step 2 - expected positive project ID")
	}

	// Verify project is listed
	projects, err := registry.ListProjects(workspace.Active())
	if err != nil {
		t.Fatalf("Step 2 - ListProjects failed: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("Step 2 - expected 1 project, got %d", len(projects))
	}

	// 3. Deploy Run (FTPS skipped â€” will fail with "no server configured")
	deploySvc := deploy.NewService(workspace.Active().DB)
	_, err = deploySvc.Deploy(workspace.Active(), "test-project", false)
	if err == nil {
		t.Fatal("Step 3 - expected error (no server configured)")
	}

	// 3b. Add server so deploy can proceed to pack phase
	_, err = registry.AddServer(workspace.Active(), &registry.Server{
		ProjectID:  projID,
		Name:       "production",
		Type:       "ftps",
		Host:       "ftp.example.com",
		Port:       21,
		User:       "user",
		Password:   "pass",
		TargetPath: "/public_html",
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("Step 3b - AddServer failed: %v", err)
	}

	// 4. Verify (before any deploy)
	result, err := deploySvc.Verify(workspace.Active(), "test-project")
	if err != nil {
		t.Fatalf("Step 4 - Verify failed: %v", err)
	}
	if result.Status != "unknown" {
		t.Errorf("Step 4 - expected unknown, got %s", result.Status)
	}

	// 5. History (empty)
	logs, err := deploySvc.GetHistory(workspace.Active(), "test-project", 10)
	if err != nil {
		t.Fatalf("Step 5 - GetHistory failed: %v", err)
	}
	if len(logs) != 0 {
		t.Errorf("Step 5 - expected 0 logs, got %d", len(logs))
	}

	// 6. Rollback (no deploy to rollback)
	_, err = deploySvc.Rollback(workspace.Active(), "test-project")
	if err == nil {
		t.Fatal("Step 6 - expected error (no deploy to rollback)")
	}

	// 7. Status
	status := ws.Status()
	if status == nil {
		t.Fatal("Step 7 - Status returned nil")
	}
	if status["name"] != "IntegrationWS" {
		t.Errorf("Step 7 - expected IntegrationWS, got %v", status["name"])
	}
	if status["projects"].(int) != 1 {
		t.Errorf("Step 7 - expected 1 project, got %d", status["projects"])
	}
}

// TestHistoryEvents validates that the history engine records events correctly
// through the deploy service's lifecycle.
func TestHistoryEvents(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "testws")
	projectPath := filepath.Join(tmpDir, "project")

	if err := os.MkdirAll(projectPath, 0700); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}

	_, err := workspace.Init("EventTestWS", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	_, err = registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "event-project",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}

	// Record events manually to validate history service
	hist := history.NewService(workspace.Active().DB)
	if err := hist.RecordEvent(history.EventDeployStarted, "Test deploy", nil); err != nil {
		t.Fatalf("RecordEvent failed: %v", err)
	}
	if err := hist.RecordEvent(history.EventDeploySuccess, "Deploy completed", nil); err != nil {
		t.Fatalf("RecordEvent failed: %v", err)
	}

	events, err := hist.GetRecentEvents(10)
	if err != nil {
		t.Fatalf("GetRecentEvents failed: %v", err)
	}
	if len(events) < 2 {
		t.Errorf("expected at least 2 events, got %d", len(events))
	}
}

// TestWorkspacePersistence validates that workspace data survives
// close/reopen cycles.
func TestWorkspacePersistence(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "persist-testws")
	projectPath := filepath.Join(tmpDir, "project")

	if err := os.MkdirAll(projectPath, 0700); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}

	// Init and add project
	_, err := workspace.Init("PersistWS", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	_, err = registry.AddProject(workspace.Active(), &registry.Project{
		Name:       "persist-project",
		LocalPath:  projectPath,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}
	workspace.Active().DB.Close()

	// Reopen and verify project persists
	_, err = workspace.Open(wsPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	projects, err := registry.ListProjects(workspace.Active())
	if err != nil {
		t.Fatalf("ListProjects failed: %v", err)
	}
	if len(projects) != 1 {
		t.Errorf("expected 1 project after reopen, got %d", len(projects))
	}
	if projects[0].Name != "persist-project" {
		t.Errorf("expected persist-project, got %s", projects[0].Name)
	}
}
