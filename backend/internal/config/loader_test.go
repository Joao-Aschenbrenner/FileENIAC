package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoader_LoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.toml")

	configContent := `
[global]
secret_env = "TEST_SECRET"
history_db = "~/.test/history.db"

[projects.test-project]
name = "test-project"
working_dir = "/tmp/test"
run_migrations = true

[projects.test-project.ftps]
host = "ftp.example.com"
port = 21
user = "testuser"
pass = "testpass"

[projects.test-project.deploy]
target_path = "/public_html/test"
verify_url = "https://example.com/test"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	loader := NewLoader()
	cfg, err := loader.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Global.SecretEnv != "TEST_SECRET" {
		t.Errorf("expected secret_env 'TEST_SECRET', got '%s'", cfg.Global.SecretEnv)
	}

	if len(cfg.Projects) != 1 {
		t.Errorf("expected 1 project, got %d", len(cfg.Projects))
	}

	p, ok := cfg.Projects["test-project"]
	if !ok {
		t.Fatal("project 'test-project' not found")
	}

	if p.FTPS.Host != "ftp.example.com" {
		t.Errorf("expected host 'ftp.example.com', got '%s'", p.FTPS.Host)
	}

	if p.RunMigrations != true {
		t.Error("expected run_migrations to be true")
	}
}

func TestLoader_GetProject(t *testing.T) {
	loader := NewLoader()
	cfg := &Config{
		Projects: map[string]*Project{
			"project-a": {Name: "project-a"},
			"project-b": {Name: "project-b"},
		},
	}

	p, err := loader.GetProject(cfg, "project-a")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if p.Name != "project-a" {
		t.Errorf("expected 'project-a', got '%s'", p.Name)
	}

	_, err = loader.GetProject(cfg, "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent project")
	}
}

func TestApplyDefaults(t *testing.T) {
	p := &Project{}
	ApplyDefaults(p)

	if len(p.Excludes) == 0 {
		t.Error("excludes should not be empty after ApplyDefaults")
	}

	if p.Deploy.BackupPrefix != ".env.bak" {
		t.Errorf("expected backup prefix '.env.bak', got '%s'", p.Deploy.BackupPrefix)
	}

	if p.FTPS.Port != 21 {
		t.Errorf("expected default port 21, got %d", p.FTPS.Port)
	}
}