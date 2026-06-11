# ADR-009: Database Schema — SQLite

## Status
APROVADO

## Data
2026-06-10

## Contexto
O eniac-deploy (PoC) já usa SQLite para histórico de deploys. Agora o banco precisa suportar workspace, registry, health checks, sync operations e auditoria.

## Decisão
SQLite com WAL mode, schema versionado.

### Tabelas

```sql
-- Workspace settings
CREATE TABLE workspace_settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TEXT DEFAULT (datetime('now'))
);

-- Projects registry
CREATE TABLE projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    display_name TEXT,
    local_path TEXT NOT NULL,
    remote_url TEXT,
    branch TEXT DEFAULT 'main',
    server_type TEXT DEFAULT 'ftps',
    server_host TEXT,
    server_port INTEGER DEFAULT 21,
    server_user TEXT,
    server_target_path TEXT,
    verify_url TEXT,
    run_migrations INTEGER DEFAULT 0,
    backup_prefix TEXT DEFAULT '.backup',
    endpoint TEXT DEFAULT '_deploy.php',
    is_active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now'))
);

-- Project dependencies
CREATE TABLE project_dependencies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    dependency_name TEXT NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id),
    UNIQUE(project_id, dependency_name)
);

-- Servers
CREATE TABLE servers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    type TEXT NOT NULL DEFAULT 'ftps',
    host TEXT NOT NULL,
    port INTEGER DEFAULT 21,
    user TEXT,
    target_path TEXT,
    verify_url TEXT,
    is_active INTEGER DEFAULT 1,
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

-- Deploy records
CREATE TABLE deploys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    commit_hash TEXT,
    commit_message TEXT,
    branch TEXT,
    artifact_hash TEXT,
    artifact_size INTEGER,
    files_count INTEGER,
    migration_result TEXT,
    deploy_manifest TEXT,
    error_message TEXT,
    started_at TEXT,
    completed_at TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

-- Rollback tracking
CREATE TABLE rollbacks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    deploy_id INTEGER NOT NULL,
    reason TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (deploy_id) REFERENCES deploys(id)
);

-- Health checks
CREATE TABLE health_checks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    status TEXT NOT NULL,
    http_status INTEGER,
    response_time_ms INTEGER,
    error_message TEXT,
    checked_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

-- Sync operations
CREATE TABLE sync_operations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    direction TEXT NOT NULL, -- 'pull' (server→mirror), 'push' (local→server)
    status TEXT NOT NULL,
    files_count INTEGER,
    divergences_count INTEGER,
    error_message TEXT,
    started_at TEXT,
    completed_at TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

-- Audit log
CREATE TABLE audit_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    action TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id INTEGER,
    project_id INTEGER,
    details TEXT,
    metadata TEXT, -- JSON
    created_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

-- Indexes
CREATE INDEX idx_deploys_project ON deploys(project_id);
CREATE INDEX idx_deploys_status ON deploys(status);
CREATE INDEX idx_deploys_created ON deploys(created_at);
CREATE INDEX idx_health_project ON health_checks(project_id);
CREATE INDEX idx_sync_project ON sync_operations(project_id);
CREATE INDEX idx_audit_action ON audit_log(action);
CREATE INDEX idx_audit_created ON audit_log(created_at);
```

### Migrations
- Schema versionado em `backend/database/migrations/`
- Naming: `001_initial.sql`, `002_add_indexes.sql`, etc.
- Aplicado automaticamente ao iniciar o backend

## Consequências
- Schema único para todo o workspace
- Sem migrations complexas (SQLite não suporta ALTER COLUMN)
- Performance adequada para uso local (< 100k registros)
- WAL mode garante concorrência leitura/escrita
