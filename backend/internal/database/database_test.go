// SPDX-License-Identifier: MIT

package database

import (
	"strings"
	"testing"
)

func TestCount_TableAllowlist(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer db.Close()

	_, err = db.Count("projects")
	if err != nil && !strings.Contains(err.Error(), "no such table") {
		t.Errorf("unexpected error for allowed table 'projects': %v", err)
	}

	blocked := []string{"sqlite_master", "users", "passwords", "; DROP TABLE projects"}
	for _, table := range blocked {
		_, err = db.Count(table)
		if err == nil {
			t.Errorf("expected error for blocked table %q", table)
		}
		if !strings.Contains(err.Error(), "not allowed for count") {
			t.Errorf("expected allowlist rejection for table %q, got: %v", table, err)
		}
	}
}

func TestCount_AllowsAllowedTables(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer db.Close()

	allowed := []string{"projects", "servers", "events", "deploy_logs", "rollback_logs", "workspace_settings", "mirror_snapshots", "sync_manifests", "repositories", "sessions"}
	for _, table := range allowed {
		_, err = db.Count(table)
		// Tables may not exist in an empty in-memory DB; only fail on allowlist rejection.
		if err != nil && strings.Contains(err.Error(), "not allowed for count") {
			t.Errorf("table %q should be allowed: %v", table, err)
		}
	}
}
