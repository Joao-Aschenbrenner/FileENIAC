package history

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func needsCGO(t *testing.T) {
	if runtime.GOOS == "js" || runtime.GOOS == "wasip1" || os.Getenv("CGO_ENABLED") == "0" {
		t.Skip("CGO required for go-sqlite3")
	}
}

func TestDB_NewDB(t *testing.T) {
	needsCGO(t)
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDB(dbPath)
	if err != nil {
		t.Fatalf("NewDB failed: %v", err)
	}
	defer db.Close()

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("database file should exist")
	}
}

func TestCRUD_Insert(t *testing.T) {
	needsCGO(t)
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDB(dbPath)
	if err != nil {
		t.Fatalf("NewDB failed: %v", err)
	}
	defer db.Close()

	crud := NewCRUD(db)

	rec := NewSuccessRecord("test-project", "abc123", "3 migrations", "def456")
	id, err := crud.Insert(rec)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	if id <= 0 {
		t.Error("inserted ID should be positive")
	}
}

func TestCRUD_GetByProject(t *testing.T) {
	needsCGO(t)
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDB(dbPath)
	if err != nil {
		t.Fatalf("NewDB failed: %v", err)
	}
	defer db.Close()

	crud := NewCRUD(db)

	crud.Insert(NewSuccessRecord("test-project", "hash1", "1 migration", "commit1"))
	crud.Insert(NewFailedRecord("test-project", "hash2", "error"))
	crud.Insert(NewSuccessRecord("test-project", "hash3", "2 migrations", "commit3"))

	records, err := crud.GetByProject("test-project", 10)
	if err != nil {
		t.Fatalf("GetByProject failed: %v", err)
	}

	if len(records) != 3 {
		t.Errorf("expected 3 records, got %d", len(records))
	}

	if records[0].Status != StatusSuccess {
		t.Error("first record should be success (most recent)")
	}
}

func TestCRUD_GetLastSuccessful(t *testing.T) {
	needsCGO(t)
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDB(dbPath)
	if err != nil {
		t.Fatalf("NewDB failed: %v", err)
	}
	defer db.Close()

	crud := NewCRUD(db)

	crud.Insert(NewFailedRecord("test-project", "hash1", "error"))
	crud.Insert(NewSuccessRecord("test-project", "hash2", "2 migrations", "commit2"))
	crud.Insert(NewSuccessRecord("test-project", "hash3", "3 migrations", "commit3"))

	last, err := crud.GetLastSuccessful("test-project")
	if err != nil {
		t.Fatalf("GetLastSuccessful failed: %v", err)
	}

	if last == nil {
		t.Fatal("last successful should not be nil")
	}

	if last.ArtifactHash != "hash3" {
		t.Errorf("expected hash3, got %s", last.ArtifactHash)
	}
}

func TestCRUD_CountByProject(t *testing.T) {
	needsCGO(t)
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDB(dbPath)
	if err != nil {
		t.Fatalf("NewDB failed: %v", err)
	}
	defer db.Close()

	crud := NewCRUD(db)

	count, err := crud.CountByProject("test-project")
	if err != nil {
		t.Fatalf("CountByProject failed: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}

	crud.Insert(NewSuccessRecord("test-project", "hash1", "", ""))
	crud.Insert(NewSuccessRecord("test-project", "hash2", "", ""))

	count, err = crud.CountByProject("test-project")
	if err != nil {
		t.Fatalf("CountByProject failed: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2, got %d", count)
	}
}