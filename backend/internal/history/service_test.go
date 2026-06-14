package history

import (
	"path/filepath"
	"testing"

	"github.com/ENIACSystems/FileENIAC/backend/internal/database"
)

func setupDB(t *testing.T) *database.DB {
	t.Helper()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	if err := db.Migrate(); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	// Insert a project to satisfy FOREIGN KEY constraints
	_, err = db.Exec("INSERT INTO projects (name, local_path, is_active) VALUES (?, ?, 1)", "test-project", tmpDir)
	if err != nil {
		t.Fatalf("Failed to insert test project: %v", err)
	}

	return db
}

func TestRecordEvent(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	s := NewService(db)

	err := s.RecordEvent(EventDeployStarted, "Test deploy", nil)
	if err != nil {
		t.Fatalf("RecordEvent failed: %v", err)
	}
}

func TestRecordDeploy(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	s := NewService(db)

	log := &DeployLog{
		ProjectID:   1,
		DeployID:    "dep_test_001",
		Status:      "success",
		CommitHash:  "abc123",
		CommitMessage: "Initial deploy",
		Branch:      "main",
		ArtifactHash: "sha256:...",
	}

	id, err := s.RecordDeploy(log)
	if err != nil {
		t.Fatalf("RecordDeploy failed: %v", err)
	}

	if id <= 0 {
		t.Error("expected positive ID")
	}
}

func TestGetDeployHistory(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	s := NewService(db)

	s.RecordDeploy(&DeployLog{ProjectID: 1, DeployID: "dep_001", Status: "success", CommitHash: "a"})
	s.RecordDeploy(&DeployLog{ProjectID: 1, DeployID: "dep_002", Status: "success", CommitHash: "b"})

	logs, err := s.GetDeployHistory(1, 10)
	if err != nil {
		t.Fatalf("GetDeployHistory failed: %v", err)
	}

	if len(logs) != 2 {
		t.Errorf("expected 2 logs, got %d", len(logs))
	}
}

func TestGetLastDeploy(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	s := NewService(db)

	s.RecordDeploy(&DeployLog{ProjectID: 1, DeployID: "dep_001", Status: "failed"})
	s.RecordDeploy(&DeployLog{ProjectID: 1, DeployID: "dep_002", Status: "success"})

	last, err := s.GetLastDeploy(1)
	if err != nil {
		t.Fatalf("GetLastDeploy failed: %v", err)
	}

	if last.DeployID != "dep_002" {
		t.Errorf("expected dep_002, got %s", last.DeployID)
	}
}

func TestRecordRollback(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	s := NewService(db)

	s.RecordDeploy(&DeployLog{ProjectID: 1, DeployID: "dep_001", Status: "success"})

	log := &RollbackLog{
		DeployID:  "dep_001",
		ProjectID: 1,
		Reason:    "test rollback",
		Status:    "completed",
	}

	id, err := s.RecordRollback(log)
	if err != nil {
		t.Fatalf("RecordRollback failed: %v", err)
	}

	if id <= 0 {
		t.Error("expected positive ID")
	}
}

func TestGetRecentEvents(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	s := NewService(db)

	s.RecordEvent(EventDeployStarted, "First", nil)
	s.RecordEvent(EventDeploySuccess, "Second", nil)

	events, err := s.GetRecentEvents(10)
	if err != nil {
		t.Fatalf("GetRecentEvents failed: %v", err)
	}

	if len(events) < 2 {
		t.Errorf("expected at least 2 events, got %d", len(events))
	}
}

func TestEventTypes(t *testing.T) {
	types := []string{
		EventDeployStarted, EventDeploySuccess, EventDeployFailed,
		EventRollbackStarted, EventRollbackSuccess, EventRollbackFailed,
		EventVerifySuccess, EventVerifyFailed,
		EventProjectCreated, EventProjectRemoved,
		EventServerAdded, EventServerUpdated,
	}

	for _, et := range types {
		if et == "" {
			t.Error("event type should not be empty")
		}
	}
}
