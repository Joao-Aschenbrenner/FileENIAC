package refresh_test

import (
	"path/filepath"
	"testing"

	"github.com/ENIACSystems/FileENIAC/backend/internal/refresh"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
)

func setupRefreshCtx(t *testing.T) *workspace.Context {
	t.Helper()
	dir := t.TempDir()
	wsPath := filepath.Join(dir, "testws")
	_, err := workspace.Init("TestRefresh", wsPath, "")
	if err != nil {
		t.Fatalf("workspace.Init failed: %v", err)
	}
	ctx := workspace.Active()
	if ctx == nil {
		t.Fatal("Active() returned nil")
	}
	return ctx
}

func cleanupRefreshCtx(t *testing.T, ctx *workspace.Context) {
	t.Helper()
	if ctx != nil && ctx.DB != nil {
		ctx.DB.Close()
	}
}

func TestRefreshResult_ZeroValue(t *testing.T) {
	r := &refresh.RefreshResult{}
	if r.Organizations != 0 {
		t.Errorf("expected 0 organizations, got %d", r.Organizations)
	}
	if r.Repositories != 0 {
		t.Errorf("expected 0 repositories, got %d", r.Repositories)
	}
	if r.ChangesFound != 0 {
		t.Errorf("expected 0 changes, got %d", r.ChangesFound)
	}
	if r.Errors != 0 {
		t.Errorf("expected 0 errors, got %d", r.Errors)
	}
}

func TestRefreshResult_Fields(t *testing.T) {
	r := &refresh.RefreshResult{
		Organizations: 3,
		Repositories:  25,
		ChangesFound:  5,
		Errors:        1,
	}
	if r.Organizations != 3 {
		t.Errorf("expected 3, got %d", r.Organizations)
	}
	if r.Repositories != 25 {
		t.Errorf("expected 25, got %d", r.Repositories)
	}
	if r.ChangesFound != 5 {
		t.Errorf("expected 5, got %d", r.ChangesFound)
	}
	if r.Errors != 1 {
		t.Errorf("expected 1, got %d", r.Errors)
	}
}

func TestRefreshGitHub_NoToken(t *testing.T) {
	ctx := setupRefreshCtx(t)
	defer cleanupRefreshCtx(t, ctx)

	result, err := refresh.RefreshGitHub(ctx)
	if err == nil {
		t.Fatal("expected error when no token set")
	}
	if result == nil {
		t.Fatal("expected non-nil result even on error")
	}
}

func TestRefreshGitHub_EmptyToken(t *testing.T) {
	ctx := setupRefreshCtx(t)
	defer cleanupRefreshCtx(t, ctx)

	if err := ctx.DB.SetSetting("github_token", ""); err != nil {
		t.Fatalf("SetSetting failed: %v", err)
	}

	result, err := refresh.RefreshGitHub(ctx)
	if err == nil {
		t.Fatal("expected error for empty token")
	}
	if result == nil {
		t.Fatal("expected non-nil result even on error")
	}
}

func TestRefreshGitHub_NoVaultConfig(t *testing.T) {
	ctx := setupRefreshCtx(t)
	defer cleanupRefreshCtx(t, ctx)

	if err := ctx.DB.SetSetting("github_token", "some-encrypted-token"); err != nil {
		t.Fatalf("SetSetting failed: %v", err)
	}
	ctx.Config.Vault.MasterKey = ""

	result, err := refresh.RefreshGitHub(ctx)
	if err == nil {
		t.Fatal("expected error when vault is not configured")
	}
	if result == nil {
		t.Fatal("expected non-nil result even on error")
	}
}
