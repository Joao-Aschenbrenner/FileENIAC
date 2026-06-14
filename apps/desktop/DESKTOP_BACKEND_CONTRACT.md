# Desktop ↔ Backend Contract

## Base URL

```
http://localhost:8080
```

O backend roda `eniac serve` e expõe uma API REST JSON. O Desktop consome essa API via HTTP.

---

## Endpoints

### Health

```
GET /health
```

Response:
```json
{ "status": "ok" }
```

---

### Workspace

```
GET /workspace?workspace=<path>
```

Response:
```json
{
  "name": "MeuWorkspace",
  "description": "...",
  "path": "C:/projetos",
  "projects": 3,
  "servers": 2,
  "deploys": 12,
  "events": 45,
  "created_at": "2026-01-15T10:00:00Z",
  "updated_at": "2026-06-11T15:30:00Z"
}
```

---

### Projects

```
GET /projects?workspace=<path>
POST /projects?workspace=<path>
GET /projects/:name?workspace=<path>
DELETE /projects/:name?workspace=<path>
```

**POST body:**
```json
{
  "name": "meu-projeto",
  "display_name": "Meu Projeto",
  "local_path": "C:/projetos/meu-projeto",
  "remote_path": "/public_html",
  "branch": "main",
  "git_url": "https://github.com/user/repo.git",
  "environment": "production"
}
```

**GET /projects response:**
```json
[
  {
    "id": 1,
    "name": "meu-projeto",
    "display_name": "Meu Projeto",
    "local_path": "C:/projetos/meu-projeto",
    "remote_path": "/public_html",
    "branch": "main",
    "git_url": "https://github.com/user/repo.git",
    "environment": "production",
    "is_active": true,
    "last_commit_hash": "abc123",
    "last_deploy_id": "dep_001",
    "divergence_status": "sincronizado",
    "created_at": "2026-01-15 10:00:00",
    "updated_at": "2026-06-11 15:30:00"
  }
]
```

---

### Servers

```
GET /servers?workspace=<path>[&project=<name>]
POST /servers?workspace=<path>
GET /servers/:id?workspace=<path>
DELETE /servers/:id?workspace=<path>
```

**POST body:**
```json
{
  "project_id": 1,
  "name": "Produção",
  "type": "ftps",
  "host": "ftp.meusite.com",
  "port": 21,
  "user": "ftpuser",
  "password": "supersecret",
  "target_path": "/public_html",
  "verify_url": "https://meusite.com/verify"
}
```

**GET /servers response** (password never exposed):
```json
[
  {
    "id": 1,
    "project_id": 1,
    "name": "Produção",
    "type": "ftps",
    "host": "ftp.meusite.com",
    "port": 21,
    "user": "ftpuser",
    "target_path": "/public_html",
    "verify_url": "https://meusite.com/verify",
    "is_active": true
  }
]
```

---

### Settings

```
GET /settings?workspace=<path>
POST /settings?workspace=<path>
```

**POST body:**
```json
{
  "key1": "value1",
  "key2": "value2"
}
```

---

### History / Events

```
GET /history?workspace=<path>[&project=<name>&type=<type>&limit=<n>&offset=<n>]
GET /events?workspace=<path>[&type=<type>&limit=<n>&offset=<n>]
GET /deploys?workspace=<path>&project=<name>[&limit=<n>]
```

**GET /history response:**
```json
[
  {
    "id": 1,
    "event_type": "DEPLOY_SUCCESS",
    "description": "Deploy meu-projeto v1.2.3",
    "metadata": "...",
    "created_at": "2026-06-11 15:30:00"
  }
]
```

Event types: `DEPLOY_STARTED`, `DEPLOY_SUCCESS`, `DEPLOY_FAILED`, `ROLLBACK_STARTED`, `ROLLBACK_SUCCESS`, `ROLLBACK_FAILED`, `VERIFY_SUCCESS`, `VERIFY_FAILED`, `SYNC_STARTED`, `SYNC_COMPLETED`, `SYNC_FAILED`, `PROJECT_CREATED`, `PROJECT_REMOVED`, `SERVER_ADDED`, `SERVER_UPDATED`, `SERVER_REMOVED`, `ALERT`, `ERROR`

---

### Deploy / Rollback / Verify

```
POST /deploy?workspace=<path>   {"project":"name","use_fallback":false}
POST /rollback?workspace=<path> {"project":"name"}
POST /verify?workspace=<path>   {"project":"name"}
```

**Response:**
```json
{
  "status": "success",
  "version": "v1.2.3",
  "deploy_id": "dep_001",
  "duration_ms": 4523
}
```

---

### Diff

```
GET /diff?workspace=<path>&project=<name>
```

**Response:**
```json
{
  "status": "divergent",
  "files": [
    {"path": "index.php", "status": "modified", "local_hash": "abc", "mirror_hash": "def"},
    {"path": "style.css", "status": "identical", "local_hash": "xyz", "mirror_hash": "xyz"}
  ]
}
```

Status values: `identical`, `modified`, `added`, `removed`

---

### Sync

```
GET /syncs?workspace=<path>[&project=<name>&limit=<n>]
POST /sync?workspace=<path> {"project":"name","action":"mirror_update"}
```

**POST /sync response:**
```json
{
  "suggestion": { "action": "mirror_update", "files": 3 },
  "manifest": { "id": 1, "result": "completed" },
  "diff": { "status": "divergent", "files": [...] }
}
```

---

### Mirror

```
POST /mirror?workspace=<path> {"project":"name"}
```

**Response:**
```json
{
  "status": "created",
  "files_count": 42,
  "duration_ms": 3200
}
```

---

### Health Check

```
GET /health/check?workspace=<path>
```

**Response:**
```json
{
  "status": "healthy",
  "projects": 3,
  "servers": 2,
  "divergent": 1,
  "last_events": [...]
}
```

---

### GitHub Authentication

```
GET  /github/status?workspace=<path>
POST /github/login?workspace=<path>
POST /github/logout?workspace=<path>
```

**POST /github/login:**
```json
{ "token": "ghp_xxxxxxxxxxxxxxxxxxxx" }
```

**GET /github/status response:**
```json
{
  "authenticated": true,
  "user": "myusername"
}
```

---

### GitHub Discovery

```
GET /github/organizations?workspace=<path>
GET /github/repositories?workspace=<path>[&org=<org>]
```

**GET /github/organizations response:**
```json
[
  { "login": "myorg", "id": 123, "url": "https://api.github.com/orgs/myorg", "avatar_url": "..." }
]
```

**GET /github/repositories response:**
```json
[
  {
    "id": 456,
    "name": "my-repo",
    "full_name": "myorg/my-repo",
    "clone_url": "https://github.com/myorg/my-repo.git",
    "default_branch": "main",
    "organization": "myorg",
    "imported": false
  }
]
```

---

### GitHub Import & Clone

```
POST /github/import?workspace=<path>
POST /github/clone?workspace=<path>
```

**POST /github/import body:**
```json
{
  "repos": [
    { "id": 456, "name": "my-repo", "full_name": "myorg/my-repo", "clone_url": "https://github.com/myorg/my-repo.git", "default_branch": "main", "organization": "myorg" }
  ],
  "clone_dir": "C:/workspace/projects"
}
```

**POST /github/import response:**
```json
[
  {
    "repo": { "name": "my-repo", ... },
    "project_id": 1,
    "repository_id": 1,
    "clone_result": { "path": "C:/workspace/projects/my-repo", "branch": "main", "commit_sha": "abc123", "duration_ms": 3200 },
    "validation": { "valid": true, "checks": [...] },
    "error": ""
  }
]
```

---

### Repositories Registry

```
GET /repositories?workspace=<path>[&org=<org>]
GET /repositories/:id?workspace=<path>
```

**GET /repositories response:**
```json
[
  {
    "id": 1,
    "github_id": 456,
    "name": "my-repo",
    "full_name": "myorg/my-repo",
    "clone_url": "https://github.com/myorg/my-repo.git",
    "default_branch": "main",
    "organization": "myorg",
    "import_status": "cloned",
    "project_id": 1,
    "clone_path": "C:/workspace/projects/my-repo",
    "last_commit": "abc123"
  }
]
```

---

## Frontend Routes (Desktop)

| Route              | Component        | Description                     |
|--------------------|------------------|---------------------------------|
| `/`                | `Onboarding`     | Wizard de configuração          |
| `/dashboard`       | `Dashboard`      | Resumo do workspace (métricas)  |
| `/bootstrap`       | `WorkspaceBootstrap` | Checklist de configuração    |
| `/projects`        | `Projects`       | Lista de projetos               |
| `/projects/:name`  | `ProjectDetails` | Detalhes + Deploy/Rollback/Verify/Diff |
| `/servers`         | `Servers`        | Lista de servidores             |
| `/github/login`    | `GitHubLogin`    | Autenticação GitHub             |
| `/github/orgs`     | `GitHubOrgs`     | Seleção de organização          |
| `/github/repos`    | `GitHubRepos`    | Seleção e importação de repositórios |
| `/deploy`          | `DeployCenter`   | Central de deploy               |
| `/rollback`        | `RollbackCenter` | Central de rollback             |
| `/sync`            | `SyncCenter`     | Central de sincronização        |
| `/diff`            | `DiffViewer`     | Visualizador de diferenças      |
| `/history`         | `History`        | Timeline de eventos (filtros)   |
| `/health`          | `HealthCenter`   | Monitoramento de saúde          |

---

## Error Format

Todos os erros seguem o formato:
```json
{ "error": "mensagem descritiva" }
```

Códigos HTTP:
- `200` — OK
- `201` — Created
- `400` — Bad Request
- `404` — Not Found
- `500` — Internal Server Error

---

## CORS

O backend **não possui CORS** por ser consumido exclusivamente pelo Tauri (WebView local).  
Se necessário no futuro, adicionar header `Access-Control-Allow-Origin: *`.
