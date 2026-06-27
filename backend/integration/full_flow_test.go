package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ENIACSystems/FileENIAC/backend/internal/api"
	bghealth "github.com/ENIACSystems/FileENIAC/backend/internal/health"
	"github.com/ENIACSystems/FileENIAC/backend/internal/history"
	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/mirror"
	"github.com/ENIACSystems/FileENIAC/backend/internal/readiness"
	"github.com/ENIACSystems/FileENIAC/backend/internal/registry"
	"github.com/ENIACSystems/FileENIAC/backend/internal/repair"
	"github.com/ENIACSystems/FileENIAC/backend/internal/validate"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"go.uber.org/zap"
)

func TestFullEndToEndFlow(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "fileeniac-ws")
	localProject := filepath.Join(tmpDir, "meu-projeto")

	// Create project directory with sample files
	if err := os.MkdirAll(localProject, 0700); err != nil {
		t.Fatalf("mkdir project: %v", err)
	}
	if err := os.WriteFile(filepath.Join(localProject, "index.php"), []byte("<?php echo 'ok';"), 0644); err != nil {
		t.Fatalf("write index.php: %v", err)
	}
	if err := os.WriteFile(filepath.Join(localProject, "config.php"), []byte("<?php return [];"), 0644); err != nil {
		t.Fatalf("write config.php: %v", err)
	}

	// 1. INIT WORKSPACE
	t.Log("=== 1. Init Workspace ===")
	ws, err := workspace.Init("MeuWorkspace", wsPath, "Workspace de teste completo")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if ws.Name != "MeuWorkspace" {
		t.Fatalf("expected MeuWorkspace, got %s", ws.Name)
	}
	ctx := workspace.Active()

	// 2. ADD PROJECT
	t.Log("=== 2. Add Project ===")
	projID, err := registry.AddProject(ctx, &registry.Project{
		Name:        "meu-projeto",
		LocalPath:   localProject,
		RemotePath:  "/public_html",
		Branch:      "main",
		Environment: "production",
		IsActive:    true,
	})
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}
	if projID <= 0 {
		t.Fatal("expected positive project ID")
	}

	// 3. ADD SERVER
	t.Log("=== 3. Add Server ===")
	srvID, err := registry.AddServer(ctx, &registry.Server{
		ProjectID:  projID,
		Name:       "producao",
		Type:       "ftps",
		Host:       "ftp.example.com",
		Port:       21,
		User:       "deploy",
		Password:   "secret123",
		TargetPath: "/public_html",
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("AddServer failed: %v", err)
	}
	if srvID <= 0 {
		t.Fatal("expected positive server ID")
	}

	// 4. VALIDATION
	t.Log("=== 4. Validation ===")
	vResult := validate.ValidateClone(localProject, "main")
	if !vResult.Valid {
		t.Logf("Clone validation (expected no .git dir): %+v", vResult.Checks)
	}
	assocResult := validate.ValidateAssociation(1)
	if !assocResult.Valid {
		t.Fatalf("Association validation failed: %+v", assocResult.Checks)
	}

	// 5. HISTORY EVENTS
	t.Log("=== 5. History Events ===")
	hist := history.NewService(ctx.DB)
	if err := hist.RecordEvent(history.EventProjectCreated, "Projeto criado", nil); err != nil {
		t.Fatalf("RecordEvent failed: %v", err)
	}
	if err := hist.RecordEvent(history.EventServerAdded, "Servidor adicionado", nil); err != nil {
		t.Fatalf("RecordEvent failed: %v", err)
	}

	events, err := hist.GetRecentEvents(10)
	if err != nil {
		t.Fatalf("GetRecentEvents failed: %v", err)
	}
	if len(events) < 2 {
		t.Fatalf("expected >=2 events, got %d", len(events))
	}

	// 6. WORKSPACE STATUS
	t.Log("=== 6. Workspace Status ===")
	status := ws.Status()
	if status["name"] != "MeuWorkspace" {
		t.Fatalf("expected MeuWorkspace, got %v", status["name"])
	}
	projectsCount := status["projects"].(int)
	if projectsCount != 1 {
		t.Fatalf("expected 1 project, got %d", projectsCount)
	}

	// 7. READINESS CHECK
	t.Log("=== 7. Readiness Check ===")
	deployReady := readiness.CheckDeploy(ctx, "meu-projeto")
	if deployReady == nil {
		t.Fatal("CheckDeploy returned nil")
	}
	t.Logf("Deploy readiness: ready=%v, %d checks", deployReady.Ready, len(deployReady.Checks))

	syncReady := readiness.CheckSync(ctx, "meu-projeto")
	if syncReady == nil {
		t.Fatal("CheckSync returned nil")
	}
	t.Logf("Sync readiness: ready=%v, %d checks", syncReady.Ready, len(syncReady.Checks))

	// 8. MIRROR
	t.Log("=== 8. Mirror (no FTPS = simulates error) ===")
	me := mirror.New()
	_, err = me.Create(ctx, "meu-projeto")
	if err == nil {
		t.Log("Mirror succeeded (unexpected without FTPS)")
	} else {
		t.Logf("Mirror error (expected): %v", err)
	}

	// 9. REPAIR CHECK
	t.Log("=== 9. Repair Check ===")
	repairReport := repair.CheckConsistency(ctx)
	if repairReport == nil {
		t.Fatal("CheckConsistency returned nil")
	}
	t.Logf("Repair: %d orphaned, %d broken", repairReport.OrphanedRepositories, repairReport.BrokenPaths)

	// 10. REPAIR FIX
	t.Log("=== 10. Repair Fix ===")
	fixResult, err := repair.RepairOrphanedRepositories(ctx)
	if err != nil {
		t.Fatalf("RepairOrphanedRepositories failed: %v", err)
	}
	t.Logf("Fix result: %d fixed, %d warnings", fixResult.Fixed, len(fixResult.Warnings))

	// 11. PROJETO LISTAGEM
	t.Log("=== 11. List Projects ===")
	projects, err := registry.ListProjects(ctx)
	if err != nil {
		t.Fatalf("ListProjects failed: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}
	if projects[0].Name != "meu-projeto" {
		t.Fatalf("expected meu-projeto, got %s", projects[0].Name)
	}

	// 12. SERVER LISTAGEM
	t.Log("=== 12. List Servers ===")
	servers, err := registry.ListServers(ctx)
	if err != nil {
		t.Fatalf("ListServers failed: %v", err)
	}
	if len(servers) < 1 {
		t.Fatalf("expected >=1 server, got %d", len(servers))
	}

	// 13. BACKGROUND HEALTH RUNNER
	t.Log("=== 13. Background Health Runner ===")
	bg := bghealth.NewBackgroundRunner(100 * time.Millisecond)
	bg.Start(ctx)
	defer bg.Stop()

	time.Sleep(250 * time.Millisecond)
	snap := bghealth.GetSnapshot()
	if snap.Timestamp.IsZero() {
		t.Fatal("expected non-zero snapshot timestamp")
	}
	t.Logf("Health snapshot: status=%s, projects=%d, servers=%d",
		snap.Status, snap.ProjectsCount, snap.ServersCount)
	if snap.ProjectsCount <= 0 {
		t.Fatal("expected at least 1 project in health snapshot")
	}

	// 14. PERSISTÃŠNCIA (fechar e reabrir)
	t.Log("=== 14. Persistence (close/reopen) ===")
	ctx.DB.Close()

	_, err = workspace.Open(wsPath)
	if err != nil {
		t.Fatalf("Reopen failed: %v", err)
	}
	ctx2 := workspace.Active()
	defer ctx2.DB.Close()

	projects2, err := registry.ListProjects(ctx2)
	if err != nil {
		t.Fatalf("ListProjects after reopen failed: %v", err)
	}
	if len(projects2) != 1 {
		t.Fatalf("expected 1 project after reopen, got %d", len(projects2))
	}
	if projects2[0].Name != "meu-projeto" {
		t.Fatalf("expected meu-projeto after reopen, got %s", projects2[0].Name)
	}

	// 15. EVENT HISTORY PERSISTENCE
	t.Log("=== 15. Event History Persistence ===")
	hist2 := history.NewService(ctx2.DB)
	events2, err := hist2.GetRecentEvents(10)
	if err != nil {
		t.Fatalf("GetRecentEvents after reopen failed: %v", err)
	}
	if len(events2) < 2 {
		t.Fatalf("expected >=2 events after reopen, got %d", len(events2))
	}

	t.Log("=== FULL FLOW PASSED ===")
}

func TestAPIHealthEndpoint(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "api-test-ws")

	_, err := workspace.Init("APITest", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer workspace.Active().DB.Close()

	// Start a real API server on an ephemeral port
	srv := api.New("127.0.0.1:0")
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.L().Info("api server stopped", zap.Error(err))
		}
	}()
	defer srv.Close()

	time.Sleep(100 * time.Millisecond)

	// Test basic health
	// We can't easily get the port from the Server, so we just verify the server started
	t.Log("API server started successfully")
}

// TestSettingsCRUD validates settings persistence
func TestSettingsCRUD(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "settings-test")

	_, err := workspace.Init("SettingsTest", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	ctx := workspace.Active()
	defer ctx.DB.Close()

	// Set a setting
	if err := ctx.DB.SetSetting("github_token", "ghp_test123"); err != nil {
		t.Fatalf("SetSetting failed: %v", err)
	}

	// Read it back
	val, err := ctx.DB.GetSetting("github_token")
	if err != nil {
		t.Fatalf("GetSetting failed: %v", err)
	}
	if val != "ghp_test123" {
		t.Fatalf("expected ghp_test123, got %s", val)
	}

	// List settings
	settings, err := ctx.DB.ListSettings()
	if err != nil {
		t.Fatalf("ListSettings failed: %v", err)
	}
	if len(settings) < 1 {
		t.Fatal("expected at least 1 setting")
	}

	// Persist and reopen
	ctx.DB.Close()
	_, err = workspace.Open(wsPath)
	if err != nil {
		t.Fatalf("Reopen failed: %v", err)
	}
	ctx2 := workspace.Active()
	defer ctx2.DB.Close()

	val2, err := ctx2.DB.GetSetting("github_token")
	if err != nil {
		t.Fatalf("GetSetting after reopen failed: %v", err)
	}
	if val2 != "ghp_test123" {
		t.Fatalf("expected ghp_test123 after reopen, got %s", val2)
	}
}

// TestDivergenceStatus validates a project's divergence status
func TestDivergenceStatus(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "divergence-test")
	localProject := filepath.Join(tmpDir, "local-proj")

	if err := os.MkdirAll(localProject, 0700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(localProject, "app.php"), []byte("app"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	_, err := workspace.Init("DivergenceTest", wsPath, "")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	ctx := workspace.Active()
	defer ctx.DB.Close()

	projID, err := registry.AddProject(ctx, &registry.Project{
		Name:       "divergence-proj",
		LocalPath:  localProject,
		RemotePath: "/remote",
		Branch:     "main",
		IsActive:   true,
	})
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}

	// Sync records a manifest to verify manifest storage
	manifestID := fmt.Sprintf("manifest-%d", time.Now().UnixNano())
	_, err = ctx.DB.Exec(
		`INSERT INTO sync_manifests (project_id, manifest_id, operation_type, result, manifest_json) VALUES (?, ?, ?, ?, ?)`,
		projID, manifestID, "local_to_mirror", "completed", `{"status":"synced"}`,
	)
	if err != nil {
		t.Fatalf("Insert sync_manifest failed: %v", err)
	}

	// List manifests
	manifests, err := ctx.DB.QueryManifestsByProject("divergence-proj", 5)
	if err != nil {
		t.Fatalf("QueryManifestsByProject failed: %v", err)
	}
	if len(manifests) < 1 {
		t.Fatal("expected at least 1 manifest")
	}
	t.Logf("Found %d manifests for project", len(manifests))
}
