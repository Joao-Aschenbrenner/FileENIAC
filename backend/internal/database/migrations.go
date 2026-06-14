package database

func Migrations() []Migration {
	return []Migration{
		{
			Version: 1,
			Name:    "initial_schema",
			SQL:     SchemaV1,
		},
		{
			Version: 2,
			Name:    "events_table",
			SQL:     SchemaV2,
		},
		{
			Version: 3,
			Name:    "deploy_logs",
			SQL:     SchemaV3,
		},
		{
			Version: 4,
			Name:    "registry_reinforce",
			SQL:     SchemaV4,
		},
		{
			Version: 5,
			Name:    "github_integration",
			SQL:     SchemaV5,
		},
	}
}

const SchemaV1 = `
CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    display_name TEXT,
    local_path TEXT NOT NULL,
    remote_path TEXT DEFAULT '/',
    branch TEXT DEFAULT 'main',
    git_url TEXT,
    environment TEXT DEFAULT 'production',
    server_id INTEGER,
    is_active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS servers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL DEFAULT 'ftps',
    host TEXT NOT NULL,
    port INTEGER DEFAULT 21,
    user TEXT,
    password TEXT,
    target_path TEXT NOT NULL,
    verify_url TEXT,
    is_active INTEGER DEFAULT 1,
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

CREATE INDEX IF NOT EXISTS idx_servers_project ON servers(project_id);
`

const SchemaV2 = `
CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_type TEXT NOT NULL,
    description TEXT,
    metadata TEXT,
    created_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_events_type ON events(event_type);
CREATE INDEX IF NOT EXISTS idx_events_created ON events(created_at);
`

const SchemaV3 = `
CREATE TABLE IF NOT EXISTS deploy_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    deploy_id TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL DEFAULT 'pending',
    commit_hash TEXT,
    commit_message TEXT,
    branch TEXT,
    artifact_hash TEXT,
    files_count INTEGER,
    manifest TEXT,
    error_message TEXT,
    started_at TEXT,
    completed_at TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

CREATE TABLE IF NOT EXISTS rollback_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    deploy_id TEXT NOT NULL,
    project_id INTEGER NOT NULL,
    reason TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

CREATE INDEX IF NOT EXISTS idx_deploy_logs_project ON deploy_logs(project_id);
CREATE INDEX IF NOT EXISTS idx_deploy_logs_status ON deploy_logs(status);
CREATE INDEX IF NOT EXISTS idx_rollback_deploy ON rollback_logs(deploy_id);

CREATE TABLE IF NOT EXISTS workspace_settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TEXT DEFAULT (datetime('now'))
);

INSERT OR IGNORE INTO workspace_settings (key, value) VALUES ('schema_version', '4');
`

const SchemaV4 = `
ALTER TABLE projects ADD COLUMN last_commit_hash TEXT DEFAULT '';
ALTER TABLE projects ADD COLUMN last_deploy_id TEXT DEFAULT '';
ALTER TABLE projects ADD COLUMN last_sync_at TEXT;
ALTER TABLE projects ADD COLUMN divergence_status TEXT DEFAULT 'unknown';
ALTER TABLE projects ADD COLUMN last_known_hash TEXT DEFAULT '';

CREATE TABLE IF NOT EXISTS mirror_snapshots (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    snapshot_id TEXT NOT NULL UNIQUE,
    files_count INTEGER DEFAULT 0,
    total_size INTEGER DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'pending',
    started_at TEXT,
    completed_at TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

CREATE INDEX IF NOT EXISTS idx_mirror_project ON mirror_snapshots(project_id);

CREATE TABLE IF NOT EXISTS sync_manifests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    manifest_id TEXT NOT NULL UNIQUE,
    operation_type TEXT NOT NULL,
    git_commit TEXT DEFAULT '',
    deploy_id TEXT DEFAULT '',
    files_count INTEGER DEFAULT 0,
    result TEXT NOT NULL DEFAULT 'pending',
    manifest_json TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

CREATE INDEX IF NOT EXISTS idx_sync_manifests_project ON sync_manifests(project_id);
CREATE INDEX IF NOT EXISTS idx_sync_manifests_created ON sync_manifests(created_at);
`

const SchemaV5 = `
ALTER TABLE projects ADD COLUMN github_id INTEGER;
ALTER TABLE projects ADD COLUMN organization TEXT DEFAULT '';
ALTER TABLE projects ADD COLUMN repo_name TEXT DEFAULT '';
ALTER TABLE projects ADD COLUMN import_status TEXT DEFAULT 'pending';
ALTER TABLE projects ADD COLUMN clone_path TEXT DEFAULT '';
ALTER TABLE projects ADD COLUMN provider TEXT DEFAULT 'manual';
ALTER TABLE projects ADD COLUMN last_sync_commit TEXT DEFAULT '';

CREATE TABLE IF NOT EXISTS repositories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    github_id INTEGER NOT NULL UNIQUE,
    name TEXT NOT NULL,
    full_name TEXT NOT NULL,
    description TEXT,
    html_url TEXT NOT NULL,
    clone_url TEXT NOT NULL,
    default_branch TEXT DEFAULT 'main',
    language TEXT,
    private INTEGER DEFAULT 0,
    organization TEXT NOT NULL,
    import_status TEXT DEFAULT 'pending',
    project_id INTEGER,
    clone_path TEXT,
    last_commit TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

CREATE INDEX IF NOT EXISTS idx_repositories_org ON repositories(organization);
CREATE INDEX IF NOT EXISTS idx_repositories_github ON repositories(github_id);
`
