// SPDX-License-Identifier: MIT
package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
	path string
}

func Open(path string) (*DB, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	conn, err := sql.Open("sqlite3", path+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{conn: conn, path: path}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) GetSetting(key string) (string, error) {
	var value string
	err := db.conn.QueryRow("SELECT value FROM workspace_settings WHERE key = ?", key).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (db *DB) SetSetting(key, value string) error {
	_, err := db.conn.Exec(
		`INSERT INTO workspace_settings (key, value, updated_at) VALUES (?, ?, datetime('now'))
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = datetime('now')`,
		key, value,
	)
	return err
}

func (db *DB) ListSettings() (map[string]string, error) {
	rows, err := db.conn.Query("SELECT key, value FROM workspace_settings ORDER BY key")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		settings[k] = v
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return settings, nil
}

func (db *DB) QueryManifests(limit int) ([]map[string]interface{}, error) {
	rows, err := db.conn.Query(`SELECT id, project_id, manifest_id, operation_type, git_commit, deploy_id, files_count, result, created_at FROM sync_manifests ORDER BY id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id, projectID, filesCount int64
		var manifestID, opType, gitCommit, deployID, result, createdAt string
		err := rows.Scan(&id, &projectID, &manifestID, &opType, &gitCommit, &deployID, &filesCount, &result, &createdAt)
		if err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"id": id, "project_id": projectID, "manifest_id": manifestID,
			"operation_type": opType, "git_commit": gitCommit, "deploy_id": deployID,
			"files_count": filesCount, "result": result, "created_at": createdAt,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (db *DB) QueryManifestsByProject(projectName string, limit int) ([]map[string]interface{}, error) {
	rows, err := db.conn.Query(`SELECT sm.id, sm.project_id, sm.manifest_id, sm.operation_type, sm.git_commit, sm.deploy_id, sm.files_count, sm.result, sm.created_at FROM sync_manifests sm JOIN projects p ON p.id = sm.project_id WHERE p.name = ? ORDER BY sm.id DESC LIMIT ?`, projectName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id, projectID, filesCount int64
		var manifestID, opType, gitCommit, deployID, result, createdAt string
		err := rows.Scan(&id, &projectID, &manifestID, &opType, &gitCommit, &deployID, &filesCount, &result, &createdAt)
		if err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"id": id, "project_id": projectID, "manifest_id": manifestID,
			"operation_type": opType, "git_commit": gitCommit, "deploy_id": deployID,
			"files_count": filesCount, "result": result, "created_at": createdAt,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (db *DB) Conn() *sql.DB {
	return db.conn
}

func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.conn.Exec(query, args...)
}

func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.conn.Query(query, args...)
}

func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.conn.QueryRow(query, args...)
}

func (db *DB) Count(table string) (int, error) {
	allowed := map[string]bool{
		"projects":           true,
		"servers":            true,
		"events":             true,
		"deploy_logs":        true,
		"rollback_logs":      true,
		"workspace_settings": true,
		"mirror_snapshots":   true,
		"sync_manifests":     true,
		"repositories":       true,
		"sessions":           true,
	}
	if !allowed[table] {
		return 0, fmt.Errorf("table %q not allowed for count", table)
	}
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	err := db.conn.QueryRow(query).Scan(&count)
	return count, err
}

type Migration struct {
	Version int
	Name    string
	SQL     string
}

type Migrator struct {
	db *DB
}

func (db *DB) Migrator() *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) EnsureSchemaTable() error {
	_, err := m.db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at TEXT DEFAULT (datetime('now'))
	)`)
	return err
}

func (m *Migrator) AppliedVersions() (map[int]bool, error) {
	applied := make(map[int]bool)
	rows, err := m.db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return applied, err
	}
	defer rows.Close()

	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return applied, err
		}
		applied[v] = true
	}
	if err := rows.Err(); err != nil {
		return applied, err
	}
	return applied, nil
}

func (m *Migrator) Apply(migration Migration) error {
	tx, err := m.db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(migration.SQL); err != nil {
		return fmt.Errorf("migration %d failed: %w", migration.Version, err)
	}

	if _, err := tx.Exec("INSERT INTO schema_migrations (version, name) VALUES (?, ?)",
		migration.Version, migration.Name); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return tx.Commit()
}

func (m *Migrator) Rollback(version int, rollbackSQL string) error {
	tx, err := m.db.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(rollbackSQL); err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = ?", version); err != nil {
		return err
	}

	return tx.Commit()
}

// RecordEvent records an event in the events table.
func (db *DB) RecordEvent(eventType, description string, metadata map[string]interface{}) error {
	metaJSON := "{}"
	if metadata != nil {
		b, err := json.Marshal(metadata)
		if err == nil {
			metaJSON = string(b)
		}
	}
	_, err := db.Exec(
		"INSERT INTO events (event_type, description, metadata) VALUES (?, ?, ?)",
		eventType, description, metaJSON,
	)
	return err
}

// Migrate runs all pending migrations
func (db *DB) Migrate() error {
	m := db.Migrator()
	if err := m.EnsureSchemaTable(); err != nil {
		return fmt.Errorf("failed to ensure schema table: %w", err)
	}

	applied, err := m.AppliedVersions()
	if err != nil {
		return fmt.Errorf("failed to get applied versions: %w", err)
	}

	allMigrations := Migrations()
	sort.Slice(allMigrations, func(i, j int) bool {
		return allMigrations[i].Version < allMigrations[j].Version
	})

	for _, mig := range allMigrations {
		if applied[mig.Version] {
			continue
		}
		if err := m.Apply(mig); err != nil {
			return fmt.Errorf("migration %d (%s) failed: %w", mig.Version, mig.Name, err)
		}
	}

	return nil
}
