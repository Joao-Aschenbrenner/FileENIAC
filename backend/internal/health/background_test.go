package health_test

import (
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/ENIACSystems/FileENIAC/backend/internal/health"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
)

func setupHealthCtx(t *testing.T) *workspace.Context {
	t.Helper()
	dir := t.TempDir()
	wsPath := filepath.Join(dir, "testws")
	_, err := workspace.Init("TestHealth", wsPath, "")
	if err != nil {
		t.Fatalf("workspace.Init failed: %v", err)
	}
	ctx := workspace.Active()
	if ctx == nil {
		t.Fatal("Active() returned nil")
	}
	return ctx
}

func cleanupHealthCtx(t *testing.T, ctx *workspace.Context) {
	t.Helper()
	if ctx != nil && ctx.DB != nil {
		ctx.DB.Close()
	}
}

func TestNewBackgroundRunner(t *testing.T) {
	interval := 5 * time.Second
	br := health.NewBackgroundRunner(interval)
	if br == nil {
		t.Fatal("NewBackgroundRunner returned nil")
	}
}

func TestBackgroundRunner_StartStop(t *testing.T) {
	ctx := setupHealthCtx(t)
	defer cleanupHealthCtx(t, ctx)

	br := health.NewBackgroundRunner(100 * time.Millisecond)
	br.Start(ctx)
	br.Stop()
}

func TestBackgroundRunner_DoubleStart(t *testing.T) {
	ctx := setupHealthCtx(t)
	defer cleanupHealthCtx(t, ctx)

	br := health.NewBackgroundRunner(100 * time.Millisecond)
	br.Start(ctx)
	br.Start(ctx)
	br.Stop()
}

func TestBackgroundRunner_StopWithoutStart(t *testing.T) {
	br := health.NewBackgroundRunner(100 * time.Millisecond)
	br.Stop()
}

func TestGetSnapshot_Initial(t *testing.T) {
	snap := health.GetSnapshot()
	if snap.ProjectsCount != 0 {
		t.Errorf("expected 0 projects in initial snapshot, got %d", snap.ProjectsCount)
	}
}

func TestGetSnapshot_AfterStart(t *testing.T) {
	ctx := setupHealthCtx(t)
	defer cleanupHealthCtx(t, ctx)

	br := health.NewBackgroundRunner(50 * time.Millisecond)
	br.Start(ctx)
	defer br.Stop()

	time.Sleep(150 * time.Millisecond)

	snap := health.GetSnapshot()
	if snap.ProjectsCount < 0 {
		t.Error("expected non-negative projects count")
	}
}

func TestGetSnapshot_WithProjectIncreasesCount(t *testing.T) {
	ctx := setupHealthCtx(t)
	defer cleanupHealthCtx(t, ctx)

	_, err := ctx.DB.Exec(
		`INSERT INTO projects (name, local_path, remote_path, branch, environment, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, 1, datetime('now'), datetime('now'))`,
		"health-proj", t.TempDir(), "/remote", "main", "production",
	)
	if err != nil {
		t.Fatalf("insert project failed: %v", err)
	}

	br := health.NewBackgroundRunner(50 * time.Millisecond)
	br.Start(ctx)
	defer br.Stop()

	time.Sleep(150 * time.Millisecond)

	snap := health.GetSnapshot()
	if snap.ProjectsCount != 1 {
		t.Errorf("expected 1 project, got %d", snap.ProjectsCount)
	}
}

func TestGetSnapshot_TokenNotSet(t *testing.T) {
	ctx := setupHealthCtx(t)
	defer cleanupHealthCtx(t, ctx)

	br := health.NewBackgroundRunner(50 * time.Millisecond)
	br.Start(ctx)
	defer br.Stop()

	time.Sleep(150 * time.Millisecond)

	snap := health.GetSnapshot()
	if snap.TokenValid {
		t.Error("expected TokenValid to be false when no token is set")
	}
}

func TestSnapshot_ConcurrentReadWrite(t *testing.T) {
	ctx := setupHealthCtx(t)
	defer cleanupHealthCtx(t, ctx)

	br := health.NewBackgroundRunner(10 * time.Millisecond)
	br.Start(ctx)
	defer br.Stop()

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				health.GetSnapshot()
			}
		}()
	}
	wg.Wait()
}

func TestBackgroundRunner_RunCheckOnStart(t *testing.T) {
	ctx := setupHealthCtx(t)
	defer cleanupHealthCtx(t, ctx)

	br := health.NewBackgroundRunner(1 * time.Hour)
	br.Start(ctx)
	defer br.Stop()

	time.Sleep(50 * time.Millisecond)

	snap := health.GetSnapshot()
	if snap.Timestamp.IsZero() {
		t.Error("expected snapshot to be populated immediately after start")
	}
}


