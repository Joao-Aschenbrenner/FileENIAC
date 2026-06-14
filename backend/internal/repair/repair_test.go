package repair_test

import (
	"path/filepath"
	"testing"

	"github.com/ENIACSystems/FileENIAC/backend/internal/repair"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
)

func setupRepairCtx(t *testing.T) *workspace.Context {
	t.Helper()
	dir := t.TempDir()
	wsPath := filepath.Join(dir, "testws")
	_, err := workspace.Init("TestRepair", wsPath, "")
	if err != nil {
		t.Fatalf("workspace.Init failed: %v", err)
	}
	ctx := workspace.Active()
	if ctx == nil {
		t.Fatal("Active() returned nil")
	}
	return ctx
}

func cleanupRepairCtx(t *testing.T, ctx *workspace.Context) {
	t.Helper()
	if ctx != nil && ctx.DB != nil {
		ctx.DB.Close()
	}
}

func TestCheckConsistency_EmptyDB(t *testing.T) {
	ctx := setupRepairCtx(t)
	defer cleanupRepairCtx(t, ctx)

	report := repair.CheckConsistency(ctx)
	if report == nil {
		t.Fatal("CheckConsistency returned nil")
	}
	if report.OrphanedRepositories != 0 {
		t.Errorf("expected 0 orphans, got %d", report.OrphanedRepositories)
	}
	if report.BrokenPaths != 0 {
		t.Errorf("expected 0 broken paths, got %d", report.BrokenPaths)
	}
	if len(report.Errors) != 0 {
		t.Errorf("expected 0 errors, got %d: %v", len(report.Errors), report.Errors)
	}
}

func TestCheckConsistency_WithOrphan(t *testing.T) {
	ctx := setupRepairCtx(t)
	defer cleanupRepairCtx(t, ctx)

	_, err := ctx.DB.Exec(
		`INSERT INTO repositories (github_id, name, full_name, html_url, clone_url, organization, import_status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, 'pending', datetime('now'), datetime('now'))`,
		1001, "orphan-repo", "org/orphan-repo", "https://github.com/org/orphan-repo",
		"https://github.com/org/orphan-repo.git", "org",
	)
	if err != nil {
		t.Fatalf("insert repository failed: %v", err)
	}

	report := repair.CheckConsistency(ctx)
	if report.OrphanedRepositories != 1 {
		t.Errorf("expected 1 orphan, got %d", report.OrphanedRepositories)
	}
	if len(report.Warnings) == 0 {
		t.Error("expected warnings for orphan")
	}
}

func TestCheckConsistency_WithBrokenPath(t *testing.T) {
	ctx := setupRepairCtx(t)
	defer cleanupRepairCtx(t, ctx)

	_, err := ctx.DB.Exec(
		`INSERT INTO projects (name, local_path, remote_path, branch, environment, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, 1, datetime('now'), datetime('now'))`,
		"broken-path-proj", "/nonexistent/path/project", "/remote", "main", "production",
	)
	if err != nil {
		t.Fatalf("insert project failed: %v", err)
	}

	_, err = ctx.DB.Exec(
		`INSERT INTO repositories (github_id, name, full_name, html_url, clone_url, clone_path, organization, import_status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, 'pending', datetime('now'), datetime('now'))`,
		1002, "broken-repo", "org/broken-repo", "https://github.com/org/broken-repo",
		"https://github.com/org/broken-repo.git", "/nonexistent/path/repo", "org",
	)
	if err != nil {
		t.Fatalf("insert repository failed: %v", err)
	}

	report := repair.CheckConsistency(ctx)
	if report.BrokenPaths == 0 {
		t.Error("expected broken paths > 0")
	}
}

func TestCheckConsistency_AllClean(t *testing.T) {
	ctx := setupRepairCtx(t)
	defer cleanupRepairCtx(t, ctx)

	dir := t.TempDir()

	_, err := ctx.DB.Exec(
		`INSERT INTO projects (name, local_path, remote_path, branch, environment, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, 1, datetime('now'), datetime('now'))`,
		"clean-proj", dir, "/remote", "main", "production",
	)
	if err != nil {
		t.Fatalf("insert project failed: %v", err)
	}

	report := repair.CheckConsistency(ctx)
	if report.OrphanedRepositories != 0 {
		t.Errorf("expected 0 orphans, got %d", report.OrphanedRepositories)
	}
}

func TestRepairOrphanedRepositories_NoMatchingProject(t *testing.T) {
	ctx := setupRepairCtx(t)
	defer cleanupRepairCtx(t, ctx)

	_, err := ctx.DB.Exec(
		`INSERT INTO repositories (github_id, name, full_name, html_url, clone_url, organization, import_status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, 'pending', datetime('now'), datetime('now'))`,
		2001, "lonely-repo", "org/lonely-repo", "https://github.com/org/lonely-repo",
		"https://github.com/org/lonely-repo.git", "org",
	)
	if err != nil {
		t.Fatalf("insert repository failed: %v", err)
	}

	report, err := repair.RepairOrphanedRepositories(ctx)
	if err != nil {
		t.Fatalf("RepairOrphanedRepositories failed: %v", err)
	}
	if report == nil {
		t.Fatal("RepairOrphanedRepositories returned nil report")
	}
	if report.Fixed != 0 {
		t.Errorf("expected 0 fixed (no matching project), got %d", report.Fixed)
	}
}

func TestRepairOrphanedRepositories_WithMatchingProject(t *testing.T) {
	ctx := setupRepairCtx(t)
	defer cleanupRepairCtx(t, ctx)

	_, err := ctx.DB.Exec(
		`INSERT INTO projects (name, local_path, remote_path, branch, environment, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, 1, datetime('now'), datetime('now'))`,
		"orphan-repo", t.TempDir(), "/remote", "main", "production",
	)
	if err != nil {
		t.Fatalf("insert project failed: %v", err)
	}

	_, err = ctx.DB.Exec(
		`INSERT INTO repositories (github_id, name, full_name, html_url, clone_url, organization, import_status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, 'pending', datetime('now'), datetime('now'))`,
		2002, "orphan-repo", "org/orphan-repo", "https://github.com/org/orphan-repo",
		"https://github.com/org/orphan-repo.git", "org",
	)
	if err != nil {
		t.Fatalf("insert repository failed: %v", err)
	}

	report, err := repair.RepairOrphanedRepositories(ctx)
	if err != nil {
		t.Fatalf("RepairOrphanedRepositories failed: %v", err)
	}
	if report.Fixed != 1 {
		t.Errorf("expected 1 fixed, got %d", report.Fixed)
	}
}
