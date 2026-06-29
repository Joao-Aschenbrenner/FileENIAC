// SPDX-License-Identifier: MIT
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
		ProjectID:     1,
		DeployID:      "dep_test_001",
		Status:        "success",
		CommitHash:    "abc123",
		CommitMessage: "Initial deploy",
		Branch:        "main",
		ArtifactHash:  "sha256:...",
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

func TestRecordDeploy_InvalidProject(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	s := NewService(db)

	log := &DeployLog{
		ProjectID:  999,
		DeployID:   "invalid_proj",
		Status:     "success",
		CommitHash: "abc",
		Branch:     "main",
	}

	_, err := s.RecordDeploy(log)
	if err == nil {
		t.Error("expected error for non-existent project ID")
	}
}

func TestGetDeployHistory_NoRecords(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	s := NewService(db)

	logs, err := s.GetDeployHistory(1, 10)
	if err != nil {
		t.Fatalf("GetDeployHistory failed: %v", err)
	}

	if len(logs) != 0 {
		t.Errorf("expected 0 logs, got %d", len(logs))
	}
}

func TestGetLastDeploy_NoDeploys(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	s := NewService(db)

	_, err := s.GetLastDeploy(1)
	if err == nil {
		t.Error("expected error when no deploys exist")
	}
}

func TestGetRecentEvents_Empty(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	s := NewService(db)

	events, err := s.GetRecentEvents(10)
	if err != nil {
		t.Fatalf("GetRecentEvents failed: %v", err)
	}

	if len(events) != 0 {
		t.Errorf("expected 0 events, got %d", len(events))
	}
}

func TestRecordRollback_NonexistentDeploy(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	s := NewService(db)

	log := &RollbackLog{
		DeployID:  "nonexistent",
		ProjectID: 1,
		Reason:    "test",
		Status:    "completed",
	}

	// Current behavior: rollback logs are recorded independently of deploy existence
	id, err := s.RecordRollback(log)
	if err != nil {
		t.Fatalf("RecordRollback should succeed: %v", err)
	}
	if id <= 0 {
		t.Error("expected positive ID")
	}
}

func TestRecordDeploy_ClosedDB(t *testing.T) {
	db := setupDB(t)
	db.Close() // Close the DB to test unavailable database

	s := NewService(db)

	log := &DeployLog{
		ProjectID:     1,
		DeployID:      "dep_closed",
		Status:        "success",
		CommitHash:    "abc",
		CommitMessage: "closed db test",
	}

	_, err := s.RecordDeploy(log)
	if err == nil {
		t.Error("expected error when database is closed")
	}
}

func TestGetDeployHistory_ClosedDB(t *testing.T) {
	db := setupDB(t)
	db.Close()

	s := NewService(db)

	_, err := s.GetDeployHistory(1, 10)
	if err == nil {
		t.Error("expected error when database is closed")
	}
}

func TestGetEventList(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	s := NewService(db)

	s.RecordEvent(EventDeployStarted, "First deploy", nil)
	s.RecordEvent(EventDeploySuccess, "Deploy completed", map[string]interface{}{"deploy_id": "dep_001"})

	// Filter by event type
	events, err := s.GetEventList(EventDeploySuccess, 10, 0)
	if err != nil {
		t.Fatalf("GetEventList failed: %v", err)
	}

	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
	if events[0].EventType != EventDeploySuccess {
		t.Errorf("expected %s, got %s", EventDeploySuccess, events[0].EventType)
	}

	// Without filter (all events)
	all, err := s.GetEventList("", 10, 0)
	if err != nil {
		t.Fatalf("GetEventList (no filter) failed: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 events, got %d", len(all))
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
