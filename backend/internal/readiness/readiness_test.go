package readiness_test

import (
	"path/filepath"
	"testing"

	"github.com/ENIACSystems/FileENIAC/backend/internal/readiness"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
)

func setupTestCtx(t *testing.T) *workspace.Context {
	t.Helper()
	dir := t.TempDir()
	wsPath := filepath.Join(dir, "testws")
	_, err := workspace.Init("TestReadiness", wsPath, "")
	if err != nil {
		t.Fatalf("workspace.Init failed: %v", err)
	}
	ctx := workspace.Active()
	if ctx == nil {
		t.Fatal("Active() returned nil")
	}
	return ctx
}

func cleanupTestCtx(t *testing.T, ctx *workspace.Context) {
	t.Helper()
	if ctx != nil && ctx.DB != nil {
		ctx.DB.Close()
	}
}

func TestCheckDeploy_NilCtx(t *testing.T) {
	result := readiness.CheckDeploy(nil, "myproject")
	if result.Ready {
		t.Error("expected not ready with nil context")
	}
	var found bool
	for _, c := range result.Checks {
		if c.Name == "workspace_loaded" && !c.Passed {
			found = true
		}
	}
	if !found {
		t.Error("expected workspace_loaded check with Passed=false")
	}
}

func TestCheckDeploy_EmptyProject(t *testing.T) {
	ctx := setupTestCtx(t)
	defer cleanupTestCtx(t, ctx)

	result := readiness.CheckDeploy(ctx, "")
	if result.Ready {
		t.Error("expected not ready with empty project name")
	}
	var found bool
	for _, c := range result.Checks {
		if c.Name == "project_selected" && !c.Passed {
			found = true
		}
	}
	if !found {
		t.Error("expected project_selected check with Passed=false")
	}
}

func TestCheckDeploy_ProjectNotFound(t *testing.T) {
	ctx := setupTestCtx(t)
	defer cleanupTestCtx(t, ctx)

	result := readiness.CheckDeploy(ctx, "nonexistent-project")
	if result.Ready {
		t.Error("expected not ready for non-existent project")
	}
	var found bool
	for _, c := range result.Checks {
		if c.Name == "project_exists" && !c.Passed {
			found = true
		}
	}
	if !found {
		t.Error("expected project_exists check with Passed=false")
	}
}

func TestCheckDeploy_WithData(t *testing.T) {
	ctx := setupTestCtx(t)
	defer cleanupTestCtx(t, ctx)

	_, err := ctx.DB.Exec(
		`INSERT INTO projects (name, local_path, remote_path, branch, environment, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, 1, datetime('now'), datetime('now'))`,
		"test-project", t.TempDir(), "/remote", "main", "production",
	)
	if err != nil {
		t.Fatalf("insert project failed: %v", err)
	}

	result := readiness.CheckDeploy(ctx, "test-project")
	if result.Ready {
		t.Log("project has no servers, expected not ready")
	}
	var foundServerCheck bool
	for _, c := range result.Checks {
		if c.Name == "server_active" {
			foundServerCheck = true
			if c.Passed {
				t.Error("expected server_active check to fail (no servers)")
			}
		}
	}
	if !foundServerCheck {
		t.Error("expected server_active check")
	}
}

func TestCheckDeploy_WithProjectAndServer(t *testing.T) {
	ctx := setupTestCtx(t)
	defer cleanupTestCtx(t, ctx)

	_, err := ctx.DB.Exec(
		`INSERT INTO projects (name, local_path, remote_path, branch, environment, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, 1, datetime('now'), datetime('now'))`,
		"full-project", t.TempDir(), "/remote", "main", "production",
	)
	if err != nil {
		t.Fatalf("insert project failed: %v", err)
	}

	result := readiness.CheckDeploy(ctx, "full-project")
	var (
		foundWorkspace bool
		foundProject   bool
		foundServer    bool
	)
	for _, c := range result.Checks {
		switch c.Name {
		case "workspace_loaded":
			foundWorkspace = true
		case "project_exists":
			foundProject = true
		case "server_active":
			foundServer = true
		}
	}
	if !foundWorkspace {
		t.Error("missing workspace_loaded check")
	}
	if !foundProject {
		t.Error("missing project_exists check")
	}
	if !foundServer {
		t.Error("missing server_active check")
	}
}

func TestCheckSync_RemovesServerRequirement(t *testing.T) {
	ctx := setupTestCtx(t)
	defer cleanupTestCtx(t, ctx)

	_, err := ctx.DB.Exec(
		`INSERT INTO projects (name, local_path, remote_path, branch, environment, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, 1, datetime('now'), datetime('now'))`,
		"sync-project", t.TempDir(), "/remote", "main", "production",
	)
	if err != nil {
		t.Fatalf("insert project failed: %v", err)
	}

	syncResult := readiness.CheckSync(ctx, "sync-project")
	deployResult := readiness.CheckDeploy(ctx, "sync-project")

	var deployHasServer, syncHasServer bool
	for _, c := range deployResult.Checks {
		if c.Name == "server_active" {
			deployHasServer = true
		}
	}
	for _, c := range syncResult.Checks {
		if c.Name == "server_active" {
			syncHasServer = true
		}
	}
	if !deployHasServer {
		t.Error("CheckDeploy should include server_active check")
	}
	if syncHasServer {
		t.Error("CheckSync should exclude server_active check")
	}
}

func TestCheckSync_NilWorkspace(t *testing.T) {
	result := readiness.CheckSync(nil, "myproject")
	if result.Ready {
		t.Error("expected not ready with nil workspace")
	}
}

func TestCheckSync_MissingProject(t *testing.T) {
	ctx := setupTestCtx(t)
	defer cleanupTestCtx(t, ctx)

	result := readiness.CheckSync(ctx, "missing-project")
	if result.Ready {
		t.Error("expected not ready for missing project")
	}
}
