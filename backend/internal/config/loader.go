package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Loader struct{}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse TOML: %w", err)
	}

	for _, p := range cfg.Projects {
		ApplyDefaults(p)
		p.WorkingDir = expandPath(p.WorkingDir)
	}

	return &cfg, nil
}

func (l *Loader) SaveConfig(cfg *Config, path string) error {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(cfg); err != nil {
		return fmt.Errorf("failed to marshal TOML: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func (l *Loader) GetProject(cfg *Config, name string) (*Project, error) {
	p, ok := cfg.Projects[name]
	if !ok {
		return nil, fmt.Errorf("project '%s' not found in config", name)
	}
	return p, nil
}

func (l *Loader) InitProject(name, workingDir string) *Project {
	p := &Project{
		Name:          name,
		WorkingDir:    workingDir,
		Excludes:      DefaultExcludes,
		RunMigrations: false,
		FTPS: FTPSConfig{
			Port: 21,
		},
		Deploy: DeployConfig{
			BackupPrefix: DefaultBackupPrefix,
			Endpoint:    DefaultEndpoint,
		},
	}
	ApplyDefaults(p)
	return p
}

func expandPath(path string) string {
	if len(path) == 0 {
		return path
	}
	if path[0] == '~' {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, path[1:])
		}
	}
	return path
}