# Architecture Audit — Release Candidate v0.1.0

> Audit date: 2026-06-27
> Sprint: 8 (Release Candidate Hardening)
> Go version: 1.26.4

---

## 1. System Overview

FileENIAC is a local-first workspace management platform combining a Go backend REST API with a Tauri v2 desktop frontend (React + TypeScript). Users manage Git repositories, GitHub integrations, and FTPS deployments through a unified dashboard.

### Components

| Component | Technology | Location |
|-----------|------------|----------|
| Backend API | Go 1.26, net/http, Cobra CLI | `backend/` |
| Database | SQLite 3 (WAL mode) | `internal/database/` |
| Desktop Shell | Tauri v2, React 18 | `apps/desktop/src-tauri/` |
| Frontend UI | TypeScript, Vite, TailwindCSS | `apps/desktop/src/` |
| Container | Multi-stage Alpine + Go | `Dockerfile` |
| Installer | NSIS (Windows x64) | `apps/desktop/src-tauri/target/release/bundle/nsis/` |

---

## 2. Architecture Decisions

### ADR-014 — Transport Layer (Implemented ✅)

Transport abstraction separates protocol implementation from domain logic.

**Interface** (`internal/transports/transports.go`):
```go
type Transport interface {
    Connect(ctx context.Context, cfg *TransportConfig) error
    Disconnect() error
    Upload(ctx context.Context, path string, data io.Reader, size int64) error
    Download(ctx context.Context, path string, w io.Writer) error
    List(ctx context.Context, path string) ([]FileInfo, error)
    Delete(ctx context.Context, path string) error
    Rename(ctx context.Context, old, new string) error
}
```

**Registered protocols**: `ftp` (via `init()` registry pattern).
**Factory**: `transports.New(cfg)` — no switch/conditional dispatch.

### Engine Layer

| Engine | Package | Responsibility |
|--------|---------|----------------|
| Registry | `internal/registry/` | Project CRUD + workspace association |
| Workspace | `internal/workspace/` | Workspace lifecycle + discovery |
| Deploy | `internal/deploy/service.go` | Orchestrates deploy pipeline |
| Sync | `internal/sync/` | Bidirectional sync + mirror |
| Diff | `internal/diff/` | Local vs mirror comparison |
| History | `internal/history/` | Audit log + deploy events |
| GitHub | `internal/github/` | OAuth + API integration |

### Data Flow

```
Desktop Shell (Tauri IPC)
         │
         ▼
   REST API (:8080)
         │
    ┌────┴──────┐
    │           │
  Engine     transports.Transport
  packages      (FTP adapter)
    │
    ▼
 SQLite DB
```

---

## 3. Observability (Sprint 5)

### Structured Logging (`internal/log/`)

- Zap JSON logger with `WithContext(ctx)`, `WithCorrelationID(ctx, id)`, `NewID()`
- Correlation ID propagated via `context.Context` and exposed as `X-Correlation-ID` header
- Never used as decision logic; never global

### Metrics (`internal/observability/metrics/`)

- Pluggable interface: `Timer`, `Counter`, `Gauge`
- No-op default implementation; `Set()`/`Get()` for optional providers
- Never blocks execution; no domain coupling

### Tracing (`internal/observability/tracing/`)

- Step-based hooks: `Tracer`, `Step`, `Start()`, `End()`, `EndWithError()`
- Wired into deploy/sync/mirror commands via `cmd/observability.go`

### Command Context (`cmd/observability.go`)

```go
commandContext(cmd) // returns context with correlation ID + traceOperation + cleanup
```

---

## 4. Security Hardening (Sprint 6)

### Input Validation

| Function | Package | Protection |
|----------|---------|------------|
| `ValidateName()` | `internal/validate/` | Regex `^[a-zA-Z0-9][a-zA-Z0-9_.-]*$` (max 128) |
| `SafePath()` | `internal/validate/` | Absolute path rejection + canonical join + prefix check |
| `SafeRelativePath()` | `internal/validate/` | Rejects `..`, absolute paths |
| `StripURLCredentials()` | `internal/clone/` | Removes userinfo from URLs before logging |
| Branch validation | `internal/clone/` | `ValidateName` on branch before git exec |

### API Hardening (`internal/api/api.go`)

- HTTP server timeouts: `ReadTimeout: 30s`, `WriteTimeout: 10m`, `IdleTimeout: 60s`, `ReadHeaderTimeout: 10s`
- `bodyLimitMiddleware`: 1MB `MaxBytesReader`
- `respondError`: sanitizes 500+ errors to generic "internal server error"; logs real error server-side
- `respondErrorWithLog`: returns sanitized response, logs actual error
- `correlationMiddleware`: generates/preserves `X-Correlation-ID`
- CORS middleware, handler chain ordering

### Path Traversal

- `internal/diff/diff.go`: `filepath.Walk` skips `ModeSymlink` entries
- `internal/validate/validate.go`: `SafePath` checks `filepath.IsAbs` before join

### SQL Injection

- `internal/database/database.go`: `Count()` uses `allowedTables` whitelist (5 tables: `projects`, `events`, `servers`, `deploy_logs`, `sync_logs`)

### Secrets Protection

- `FILEENIAC_VAULT_PASSWORD`: required env var, no fallback
- GitHub tokens stripped from URLs before logging
- Frontend password fields use `autoComplete="new-password"`

### Concurrency

- `sync.RWMutex` for `activeContext` race closure (was TOCTOU)
- `sync.Once` for `BackgroundRunner.Stop` idempotency
- `Timer.Reset` with explicit channel drain
- `atomic.Pointer[Context]` replaces `sync.RWMutex` for context switching
- All 31 packages green under `go test -race ./...`

---

## 5. Packaging & Distribution (Sprint 7)

### Docker

- **Base image**: `golang:1.26-alpine` (builder), `alpine:3.21` (runtime)
- **Build**: multi-stage, `CGO_ENABLED=1`, `go build -trimpath -ldflags="-s -w"`
- **Runtime**: `apk add ca-certificates sqlite-libs libffi tzdata`
- **Healthcheck**: `wget --spider http://localhost:8080/api/health`
- **Volume**: `fileeniac-data:/data` (persists SQLite + config)
- **Restart**: `unless-stopped`

### Docker Compose

```yaml
services:
  fileeniac:
    ports: ["8080:8080"]
    volumes: [fileeniac-data:/data]
    environment:
      FILEENIAC_DATA_DIR: /data
      FILEENIAC_VAULT_PASSWORD: ${FILEENIAC_VAULT_PASSWORD:?required}
      FILEENIAC_API_TOKEN: ${FILEENIAC_API_TOKEN:-}
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/api/health"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 15s
    restart: unless-stopped
```

### Desktop (Tauri v2)

- **Binary**: `fileeniac.exe` (Rust/Tauri)
- **Installer**: `FileENIAC_0.1.0_x64-setup.exe` (NSIS, ~4 MB)
- **Installer features**: desktop shortcut, Start Menu, uninstall registration
- **IPC**: `get_api_port` command reads `ENIAC_API_PORT` env var (default 8080)

### Frontend Build

- Vite + React 18 + TypeScript 5
- TailwindCSS 3
- 65 modules, ~230 KB JS (gzipped: 67 KB)

---

## 6. Known Gaps (Release Blockers)

### ✅ Session Management API — RESOLVED (Sprint 9)

The session management REST API is now implemented in `backend/internal/api/api.go` and `backend/internal/database/sessions.go`:

| Frontend Function | Endpoint | Status |
|-------------------|----------|--------|
| `listSessions()` | `GET /api/sessions` | ✅ Implemented |
| `createSession()` | `POST /api/sessions` | ✅ Implemented |
| `updateSession()` | `POST /api/sessions/{id}` | ✅ Implemented |
| `deleteSession()` | `DELETE /api/sessions/{id}` | ✅ Implemented |
| `activateSession()` | `POST /api/sessions/{id}/activate` | ✅ Implemented |
| `clearSessionWorkspace()` | `POST /api/sessions/{id}/clear-workspace` | ✅ Implemented |

Database layer uses `sessions` table (SchemaV6) with `SessionStore` providing all CRUD operations. No longer a release blocker.

---

## 7. Known Gaps (Non-Blockers)

| Issue | Severity | Notes |
|-------|----------|-------|
| NSIS installer not code-signed | LOW | Windows SmartScreen may show warning |
| No antivirus scanning of installer | LOW | Should be validated before distribution |
| Backend WriteTimeout (10m) inconsistent with Sprint 6 spec (2m for sync) | LOW | Current 30s read timeout is conservative |
| No GitHub Actions CI/CD pipeline | LOW | Build validation done manually |
| Symlink walk tests skipped on Windows | INFO | Developer Mode required for `os.Symlink` |

---

## 8. Gate Status

| Gate | Status |
|------|--------|
| `go build ./...` | ✅ Pass |
| `go vet ./...` | ✅ Pass |
| `go test ./...` | ✅ Pass (31 packages) |
| `go test -race ./...` | ✅ Pass (0 races) |
| `docker build .` | ✅ Pass |
| `docker compose up` | ✅ Pass (container healthy) |
| `tsc && vite build` (frontend) | ✅ Pass |
| `cargo tauri build` | ✅ Pass |
| NSIS installer generated | ✅ `FileENIAC_0.1.0_x64-setup.exe` |
| Session API tests | ✅ `TestSessions_CRUD` (8 subtests) |

---

## 9. Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Session API now implemented (Sprint 9) | RESOLVED | RESOLVED | ✅ RC2 ready |
| Unsigned installer flagged by SmartScreen | MEDIUM | LOW | Code sign before release; document warning |
| Go version drift (Dockerfile vs go.mod) | LOW | MEDIUM | go 1.26 consistent across all artifacts |
| Secrets leaked via logs (URL credentials) | LOW | HIGH | `StripURLCredentials` in place; verify logs before release |

---

## 10. Next Steps (Recommended)

1. **Sprint 10 — Final Release Prep**: After RC2:
   - Release notes final
   - Tag v0.1.0
   - Checksum generation
   - Installer validation (VirusTotal)
   - User manual
   - Known limitations documented
   - v0.1.0 final release

2. **Code signing**: Acquire certificate for NSIS installer (Windows SmartScreen)
3. **CI/CD**: Add GitHub Actions workflow for automated gate validation
4. **Antivirus scan**: Submit installer to VirusTotal before distribution
2. **Code signing**: Acquire certificate for NSIS installer
3. **CI/CD**: Add GitHub Actions workflow for automated gate validation
4. **Antivirus scan**: Submit installer to VirusTotal before distribution
5. **API contract tests**: Add integration tests for all `/api/*` endpoints

---

*Documento gerado durante Sprint 8 — Release Candidate Hardening*
*Última atualização: 2026-06-27*