package registry

import (
	"fmt"

	"github.com/ENIACSystems/FileENIAC/backend/internal/log"
	"github.com/ENIACSystems/FileENIAC/backend/internal/vault"
	"github.com/ENIACSystems/FileENIAC/backend/internal/workspace"
	"go.uber.org/zap"
)

type Project struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	DisplayName      string `json:"display_name,omitempty"`
	LocalPath        string `json:"local_path"`
	RemotePath       string `json:"remote_path"`
	Branch           string `json:"branch"`
	GitURL           string `json:"git_url,omitempty"`
	Environment      string `json:"environment"`
	ServerID         int64  `json:"server_id,omitempty"`
	IsActive         bool   `json:"is_active"`
	LastCommitHash   string `json:"last_commit_hash,omitempty"`
	LastDeployID     string `json:"last_deploy_id,omitempty"`
	LastSyncAt       string `json:"last_sync_at,omitempty"`
	DivergenceStatus string `json:"divergence_status,omitempty"`
	LastKnownHash    string `json:"last_known_hash,omitempty"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

type Server struct {
	ID          int64  `json:"id"`
	ProjectID   int64  `json:"project_id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	User        string `json:"user"`
	Password    string `json:"-"`
	TargetPath  string `json:"target_path"`
	VerifyURL   string `json:"verify_url,omitempty"`
	IsActive    bool   `json:"is_active"`
}

func AddProject(ctx *workspace.Context, p *Project) (int64, error) {
	query := `INSERT INTO projects (name, display_name, local_path, remote_path, branch, git_url, environment, server_id, is_active, last_commit_hash, last_deploy_id, divergence_status, last_known_hash, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1, '', '', 'unknown', '', datetime('now'), datetime('now'))`

	result, err := ctx.DB.Exec(query,
		p.Name, p.DisplayName, p.LocalPath, p.RemotePath,
		p.Branch, p.GitURL, p.Environment, p.ServerID,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert project: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	ctx.DB.RecordEvent("project_added", fmt.Sprintf("Project %s added to workspace", p.Name),
		map[string]interface{}{"project_id": id, "name": p.Name})

	log.L().Info("project added", zap.String("name", p.Name), zap.Int64("id", id))
	return id, nil
}

func RemoveProject(ctx *workspace.Context, projectID int64) error {
	_, err := ctx.DB.Exec("DELETE FROM servers WHERE project_id = ?", projectID)
	if err != nil {
		return fmt.Errorf("failed to remove servers: %w", err)
	}

	_, err = ctx.DB.Exec("DELETE FROM projects WHERE id = ?", projectID)
	if err != nil {
		return fmt.Errorf("failed to remove project: %w", err)
	}

	log.L().Info("project removed", zap.Int64("id", projectID))
	return nil
}

func ListProjects(ctx *workspace.Context) ([]*Project, error) {
	rows, err := ctx.DB.Query(`SELECT id, name, COALESCE(display_name,''), local_path, remote_path, branch, COALESCE(git_url,''), environment, COALESCE(server_id, 0), is_active, COALESCE(last_commit_hash,''), COALESCE(last_deploy_id,''), COALESCE(last_sync_at,''), COALESCE(divergence_status,'unknown'), COALESCE(last_known_hash,''), created_at, updated_at FROM projects ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	var projects []*Project
	for rows.Next() {
		p := &Project{}
		err := rows.Scan(&p.ID, &p.Name, &p.DisplayName, &p.LocalPath, &p.RemotePath,
			&p.Branch, &p.GitURL, &p.Environment, &p.ServerID, &p.IsActive,
			&p.LastCommitHash, &p.LastDeployID, &p.LastSyncAt, &p.DivergenceStatus, &p.LastKnownHash,
			&p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return projects, nil
}

func GetProject(ctx *workspace.Context, name string) (*Project, error) {
	p := &Project{}
	query := `SELECT id, name, COALESCE(display_name,''), local_path, remote_path, branch, COALESCE(git_url,''), environment, COALESCE(server_id, 0), is_active, COALESCE(last_commit_hash,''), COALESCE(last_deploy_id,''), COALESCE(last_sync_at,''), COALESCE(divergence_status,'unknown'), COALESCE(last_known_hash,''), created_at, updated_at FROM projects WHERE name = ?`
	err := ctx.DB.QueryRow(query, name).Scan(
		&p.ID, &p.Name, &p.DisplayName, &p.LocalPath, &p.RemotePath,
		&p.Branch, &p.GitURL, &p.Environment, &p.ServerID, &p.IsActive,
		&p.LastCommitHash, &p.LastDeployID, &p.LastSyncAt, &p.DivergenceStatus, &p.LastKnownHash,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("project %s not found: %w", name, err)
	}
	return p, nil
}

// UpdateProjectState updates the runtime state fields of a project.
func UpdateProjectState(ctx *workspace.Context, projectID int64, commitHash, deployID, knownHash, divergence string) error {
	_, err := ctx.DB.Exec(
		`UPDATE projects SET last_commit_hash=?, last_deploy_id=?, last_known_hash=?, divergence_status=?, last_sync_at=datetime('now'), updated_at=datetime('now') WHERE id=?`,
		commitHash, deployID, knownHash, divergence, projectID,
	)
	if err != nil {
		return fmt.Errorf("failed to update project state: %w", err)
	}
	return nil
}

func AddServer(ctx *workspace.Context, s *Server) (int64, error) {
	password := s.Password
	if password != "" {
		v, err := VaultFromCtx(ctx)
		if err != nil {
			return 0, fmt.Errorf("vault init: %w", err)
		}
		enc, err := v.Encrypt(password)
		if err != nil {
			return 0, fmt.Errorf("failed to encrypt password: %w", err)
		}
		password = enc
	}

	query := `INSERT INTO servers (project_id, name, type, host, port, user, password, target_path, verify_url, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1)`
	result, err := ctx.DB.Exec(query, s.ProjectID, s.Name, s.Type, s.Host, s.Port, s.User, password, s.TargetPath, s.VerifyURL)
	if err != nil {
		return 0, fmt.Errorf("failed to add server: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	ctx.DB.RecordEvent("server_added", fmt.Sprintf("Server %s added to project %d", s.Name, s.ProjectID),
		map[string]interface{}{"server_id": id, "project_id": s.ProjectID})

	log.L().Info("server added", zap.String("host", s.Host), zap.Int64("project_id", s.ProjectID))
	return id, nil
}

func GetServer(ctx *workspace.Context, projectID int64) (*Server, error) {
	s := &Server{}
	query := `SELECT id, project_id, name, type, host, port, user, password, target_path, verify_url, is_active FROM servers WHERE project_id = ? AND is_active = 1 LIMIT 1`
	err := ctx.DB.QueryRow(query, projectID).Scan(
		&s.ID, &s.ProjectID, &s.Name, &s.Type, &s.Host, &s.Port, &s.User, &s.Password, &s.TargetPath, &s.VerifyURL, &s.IsActive,
	)
	if err != nil {
		return nil, fmt.Errorf("no active server for project %d", projectID)
	}

	if s.Password != "" {
		v, err := VaultFromCtx(ctx)
		if err == nil {
			dec, err := v.Decrypt(s.Password)
			if err == nil {
				s.Password = dec
			}
			// if decryption fails, keep password as-is (legacy plaintext)
		}
	}

	return s, nil
}

func UpdateServer(ctx *workspace.Context, s *Server) error {
	var password string
	if s.Password != "" {
		v, err := VaultFromCtx(ctx)
		if err != nil {
			return fmt.Errorf("vault init: %w", err)
		}
		enc, err := v.Encrypt(s.Password)
		if err != nil {
			return fmt.Errorf("failed to encrypt password: %w", err)
		}
		password = enc
	} else {
		var currentPassword string
		err := ctx.DB.QueryRow(`SELECT password FROM servers WHERE id = ?`, s.ID).Scan(&currentPassword)
		if err == nil {
			password = currentPassword
		}
	}

	_, err := ctx.DB.Exec(
		`UPDATE servers SET name=?, type=?, host=?, port=?, user=?, password=?, target_path=?, verify_url=?, is_active=? WHERE id=?`,
		s.Name, s.Type, s.Host, s.Port, s.User, password, s.TargetPath, s.VerifyURL, s.IsActive, s.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update server: %w", err)
	}
	return nil
}

func RemoveServer(ctx *workspace.Context, serverID int64) error {
	_, err := ctx.DB.Exec("DELETE FROM servers WHERE id = ?", serverID)
	if err != nil {
		return fmt.Errorf("failed to remove server: %w", err)
	}
	return nil
}

func ListServers(ctx *workspace.Context) ([]*Server, error) {
	rows, err := ctx.DB.Query(`SELECT id, project_id, name, type, host, port, user, password, target_path, verify_url, is_active FROM servers ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("failed to list servers: %w", err)
	}
	defer rows.Close()

	var servers []*Server
	for rows.Next() {
		s := &Server{}
		err := rows.Scan(&s.ID, &s.ProjectID, &s.Name, &s.Type, &s.Host, &s.Port, &s.User, &s.Password, &s.TargetPath, &s.VerifyURL, &s.IsActive)
		if err != nil {
			return nil, fmt.Errorf("failed to scan server: %w", err)
		}
		servers = append(servers, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return servers, nil
}

func ListServersByProject(ctx *workspace.Context, projectID int64) ([]*Server, error) {
	rows, err := ctx.DB.Query(`SELECT id, project_id, name, type, host, port, user, password, target_path, verify_url, is_active FROM servers WHERE project_id = ? ORDER BY name`, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list servers: %w", err)
	}
	defer rows.Close()

	var servers []*Server
	for rows.Next() {
		s := &Server{}
		err := rows.Scan(&s.ID, &s.ProjectID, &s.Name, &s.Type, &s.Host, &s.Port, &s.User, &s.Password, &s.TargetPath, &s.VerifyURL, &s.IsActive)
		if err != nil {
			return nil, fmt.Errorf("failed to scan server: %w", err)
		}
		servers = append(servers, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return servers, nil
}

func GetServerByID(ctx *workspace.Context, serverID int64) (*Server, error) {
	s := &Server{}
	query := `SELECT id, project_id, name, type, host, port, user, password, target_path, verify_url, is_active FROM servers WHERE id = ?`
	err := ctx.DB.QueryRow(query, serverID).Scan(
		&s.ID, &s.ProjectID, &s.Name, &s.Type, &s.Host, &s.Port, &s.User, &s.Password, &s.TargetPath, &s.VerifyURL, &s.IsActive,
	)
	if err != nil {
		return nil, fmt.Errorf("server %d not found", serverID)
	}
	return s, nil
}

func VaultFromCtx(ctx *workspace.Context) (*vault.Vault, error) {
	if ctx.Config.Vault.MasterKey == "" {
		return nil, fmt.Errorf("vault master key not set in workspace config")
	}
	return vault.New(ctx.Config.Vault.MasterKey)
}

func UpdateRepo(ctx *workspace.Context, projectID int64, gitURL, branch string) error {
	_, err := ctx.DB.Exec(
		`UPDATE projects SET git_url=?, branch=?, updated_at=datetime('now') WHERE id=?`,
		gitURL, branch, projectID,
	)
	if err != nil {
		return fmt.Errorf("failed to update repo: %w", err)
	}
	return nil
}

func UpdateProject(ctx *workspace.Context, p *Project) error {
	_, err := ctx.DB.Exec(
		`UPDATE projects SET display_name=?, local_path=?, remote_path=?, branch=?, git_url=?, environment=?, server_id=?, is_active=?, updated_at=datetime('now') WHERE id=?`,
		p.DisplayName, p.LocalPath, p.RemotePath, p.Branch, p.GitURL, p.Environment, p.ServerID, p.IsActive, p.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}
	return nil
}

// GitHub repository registry

type Repository struct {
	ID            int64  `json:"id"`
	GitHubID      int64  `json:"github_id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Description   string `json:"description,omitempty"`
	HTMLURL       string `json:"html_url"`
	CloneURL      string `json:"clone_url"`
	DefaultBranch string `json:"default_branch"`
	Language      string `json:"language,omitempty"`
	Private       bool   `json:"private"`
	Organization  string `json:"organization"`
	ImportStatus  string `json:"import_status"`
	ProjectID     int64  `json:"project_id,omitempty"`
	ClonePath     string `json:"clone_path,omitempty"`
	LastCommit    string `json:"last_commit,omitempty"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

func AddRepository(ctx *workspace.Context, r *Repository) (int64, error) {
	result, err := ctx.DB.Exec(
		`INSERT INTO repositories (github_id, name, full_name, description, html_url, clone_url, default_branch, language, private, organization, import_status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'pending', datetime('now'), datetime('now'))`,
		r.GitHubID, r.Name, r.FullName, r.Description, r.HTMLURL, r.CloneURL, r.DefaultBranch, r.Language, boolToInt(r.Private), r.Organization,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert repository: %w", err)
	}
	return result.LastInsertId()
}

func ListRepositoriesByOrg(ctx *workspace.Context, org string) ([]*Repository, error) {
	rows, err := ctx.DB.Query(
		`SELECT id, github_id, name, full_name, COALESCE(description,''), html_url, clone_url, default_branch, COALESCE(language,''), COALESCE(private,0), organization, COALESCE(import_status,'pending'), COALESCE(project_id,0), COALESCE(clone_path,''), COALESCE(last_commit,''), created_at, updated_at
		 FROM repositories WHERE organization = ? ORDER BY name`, org)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRepositories(rows)
}

func ListAllRepositories(ctx *workspace.Context) ([]*Repository, error) {
	rows, err := ctx.DB.Query(
		`SELECT id, github_id, name, full_name, COALESCE(description,''), html_url, clone_url, default_branch, COALESCE(language,''), COALESCE(private,0), organization, COALESCE(import_status,'pending'), COALESCE(project_id,0), COALESCE(clone_path,''), COALESCE(last_commit,''), created_at, updated_at
		 FROM repositories ORDER BY organization, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRepositories(rows)
}

func GetRepositoryByGitHubID(ctx *workspace.Context, githubID int64) (*Repository, error) {
	r := &Repository{}
	err := ctx.DB.QueryRow(
		`SELECT id, github_id, name, full_name, COALESCE(description,''), html_url, clone_url, default_branch, COALESCE(language,''), COALESCE(private,0), organization, COALESCE(import_status,'pending'), COALESCE(project_id,0), COALESCE(clone_path,''), COALESCE(last_commit,''), created_at, updated_at
		 FROM repositories WHERE github_id = ?`, githubID,
	).Scan(&r.ID, &r.GitHubID, &r.Name, &r.FullName, &r.Description, &r.HTMLURL, &r.CloneURL, &r.DefaultBranch, &r.Language, &r.Private, &r.Organization, &r.ImportStatus, &r.ProjectID, &r.ClonePath, &r.LastCommit, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository %d not found", githubID)
	}
	return r, nil
}

func UpdateRepositoryImport(ctx *workspace.Context, repoID int64, projectID int64, status, clonePath string) error {
	_, err := ctx.DB.Exec(
		`UPDATE repositories SET project_id=?, import_status=?, clone_path=?, updated_at=datetime('now') WHERE id=?`,
		projectID, status, clonePath, repoID,
	)
	return err
}

func UpdateRepositoryFromGitHub(ctx *workspace.Context, repoID int64, name, fullName, cloneURL, defaultBranch, description string) error {
	_, err := ctx.DB.Exec(
		`UPDATE repositories SET name=?, full_name=?, clone_url=?, default_branch=?, description=?, updated_at=datetime('now') WHERE id=?`,
		name, fullName, cloneURL, defaultBranch, description, repoID,
	)
	return err
}

func ImportedRepositories(ctx *workspace.Context) ([]*Repository, error) {
	rows, err := ctx.DB.Query(
		`SELECT id, github_id, name, full_name, COALESCE(description,''), html_url, clone_url, default_branch, COALESCE(language,''), COALESCE(private,0), organization, COALESCE(import_status,'pending'), COALESCE(project_id,0), COALESCE(clone_path,''), COALESCE(last_commit,''), created_at, updated_at
		 FROM repositories WHERE project_id > 0 ORDER BY organization, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRepositories(rows)
}

func scanRepositories(rows interface{ Next() bool; Scan(dest ...interface{}) error; Err() error }) ([]*Repository, error) {
	var repos []*Repository
	for rows.Next() {
		r := &Repository{}
		err := rows.Scan(&r.ID, &r.GitHubID, &r.Name, &r.FullName, &r.Description, &r.HTMLURL, &r.CloneURL, &r.DefaultBranch, &r.Language, &r.Private, &r.Organization, &r.ImportStatus, &r.ProjectID, &r.ClonePath, &r.LastCommit, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}
		repos = append(repos, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return repos, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
