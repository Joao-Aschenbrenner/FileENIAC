package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ENIACSystems/FileENIAC/backend/internal/database"
	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/vault"
	"go.uber.org/zap"
)

type Workspace struct {
	Name        string    `json:"name" toml:"name"`
	Description string    `json:"description,omitempty" toml:"description,omitempty"`
	Path        string    `json:"path" toml:"-"`
	CreatedAt   time.Time `json:"created_at" toml:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" toml:"updated_at"`
}

type Context struct {
	Workspace *Workspace
	Config    *Config
	DB        *database.DB
}

type Config struct {
	Workspace WorkspaceConfig `toml:"workspace"`
	Vault     VaultConfig     `toml:"vault,omitempty"`
}

type WorkspaceConfig struct {
	Name        string `toml:"name"`
	Description string `toml:"description,omitempty"`
}

type VaultConfig struct {
	MasterKey string `toml:"master_key,omitempty"`
}

var activeContext *Context

func Active() *Context {
	return activeContext
}

func Init(name, path, desc string) (*Workspace, error) {
	wsDir := filepath.Join(path, ".eniac")
	if err := os.MkdirAll(wsDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create workspace dir: %w", err)
	}

	ws := &Workspace{
		Name:        name,
		Description: desc,
		Path:        path,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	db, err := database.Open(filepath.Join(wsDir, "workspace.db"))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	vaultKey, err := vault.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate vault key: %w", err)
	}

	cfg := &Config{
		Workspace: WorkspaceConfig{
			Name:        name,
			Description: desc,
		},
		Vault: VaultConfig{
			MasterKey: vaultKey,
		},
	}

	if err := saveConfig(cfg, wsDir); err != nil {
		return nil, fmt.Errorf("failed to save config: %w", err)
	}

	activeContext = &Context{
		Workspace: ws,
		Config:    cfg,
		DB:        db,
	}

	log.L().Info("workspace initialized",
		zap.String("name", name),
		zap.String("path", path),
	)

	return ws, nil
}

func Open(path string) (*Workspace, error) {
	wsDir := filepath.Join(path, ".eniac")
	if _, err := os.Stat(wsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("no workspace found at %s", path)
	}

	// Close previous context if any
	if activeContext != nil {
		_ = activeContext.DB.Close()
	}

	cfg, err := loadConfig(wsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	dbPath := filepath.Join(wsDir, "workspace.db")
	db, err := database.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	ws := &Workspace{
		Name:        cfg.Workspace.Name,
		Description: cfg.Workspace.Description,
		Path:        path,
	}

	activeContext = &Context{
		Workspace: ws,
		Config:    cfg,
		DB:        db,
	}

	log.L().Info("workspace opened",
		zap.String("name", ws.Name),
		zap.String("path", path),
	)

	return ws, nil
}

func (ws *Workspace) Status() map[string]interface{} {
	if activeContext == nil {
		return nil
	}

	projectCount, _ := activeContext.DB.Count("projects", "1=1")
	serverCount, _ := activeContext.DB.Count("servers", "1=1")
	deployCount, _ := activeContext.DB.Count("deploy_logs", "1=1")
	eventCount, _ := activeContext.DB.Count("events", "1=1")

	return map[string]interface{}{
		"name":        ws.Name,
		"description": ws.Description,
		"path":        ws.Path,
		"projects":    projectCount,
		"servers":     serverCount,
		"deploys":     deployCount,
		"events":      eventCount,
		"created_at":  ws.CreatedAt.Format(time.RFC3339),
		"updated_at":  ws.UpdatedAt.Format(time.RFC3339),
	}
}

func SaveConfig(cfg *Config, wsDir string) error {
	path := filepath.Join(wsDir, "config.toml")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(cfg)
}

func saveConfig(cfg *Config, wsDir string) error {
	return SaveConfig(cfg, wsDir)
}

func loadConfig(wsDir string) (*Config, error) {
	path := filepath.Join(wsDir, "config.toml")
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
