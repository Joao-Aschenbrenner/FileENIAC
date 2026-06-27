# Changelog

All notable changes to FileENIAC will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Architecture

- [x] **Transport Layer (ADR-014):** Interface `Transport` + `TransportConfig` + `FileInfo` na
      `internal/transports` — abstrai protocolos de transferência do domínio (Sprint 1A).
- [x] **Factory + Registry:** `transports.New(cfg)` com registro via `init()`, padrão
      `database/sql`. Sem switch/condicionais (Sprint 1B).
- [x] **FTP Transport Adapter:** `internal/transports/ftp/` — delega para `deploy/ftp.Client`,
      registra protocolo `"ftp"` via `init()` (Sprint 1C).
- [x] **Deploy desacoplado:** `internal/deploy/service.go` substitui `ftpClientIface` por
      `transports.Transport` — sem import de `deploy/ftp` (Sprint 1D).
- [x] **Mirror desacoplado:** `internal/mirror/mirror.go` substitui `ftplib "github.com/jlaffaye/ftp"`
      por `transports.Transport` — sem import externo de FTP (Sprint 1E).
- [x] **Arquitetura validada:** `go build ./...` + `go vet ./...` + `go test ./...` + `go test -race ./...`
      verdes (exceto data race pré-existente TD-001 em api.go).

### Security

A comprehensive security audit identified **8 CRITICAL** and **15 HIGH** issues.
All CRITICAL and HIGH findings have been remediated in this hardening campaign.

#### Authentication & Authorization

- [x] **Bearer token authentication middleware** enforced on every route except
      `/api/health`. Tokens are 32-byte ephemeral secrets generated with
      `crypto/rand` and supplied to the backend via the `FILEENIAC_API_TOKEN`
      environment variable.
- [x] **One-shot handshake endpoint** (`/api/_handshake/token`) for retrieving
      the bearer token without exposing it through logging or static config.
- [x] **CORS preflight short-circuited** through the auth layer so OPTIONS
      requests are validated before any handler dispatch.
- [x] **Tauri capability permissions** made explicit (`allow-*-api-{token,port}`)
      so only the minimum required APIs are reachable from the frontend.

#### Secrets & Configuration

- [x] **Hardcoded vault password removed.** `FILEENIAC_VAULT_PASSWORD` is now a
      required environment variable; the backend refuses to boot without it.
- [x] **Environment validation** rejects characters used for shell injection
      (`;`, `&`, `|`, backtick, `$`, single/double quotes, NUL) before any
      spawn/exec of user input.

#### Input Handling

- [x] **SQL injection fixed** in the `database.Count` helper via an
      `allowedCountTables` whitelist; user-supplied identifiers can no longer
      reach the SQL layer.
- [x] **Path traversal protection** added using a regex `ValidateName` plus
      a canonical path join (`pathsafety` package) for every filesystem
      operation sourced from user input.
- [x] **Symlinks not followed** during diff and packer directory walks,
      preventing escape via crafted file-system layouts.

#### Transport & Limits

- [x] **FTP TLS hardened:** `MinVersion = TLS12` plus ECDHE-only cipher suites
      on the configuration deploy channel.
- [x] **MaxBytesReader 10MB body limit** on inbound HTTP requests.
- [x] **Per-endpoint request timeouts:** 2 minutes for sync/mirror/diff,
      10 minutes for deploy/rollback/verify, 30 seconds default.
- [x] **GitHub User-Agent** normalized to `FileENIAC/1.0` for outbound API
      calls (previously emitted Go's default identifier).

#### Concurrency

- [x] **`activeContext` race closed:** `sync.RWMutex` replaced with
      `atomic.Pointer[Context]`, eliminating the TOCTOU window where requests
      could observe a state mid-transition.
- [x] **Heartbeat shutdown no longer calls `os.Exit(0)`** — uses
      `context.WithCancel` for graceful termination.
- [x] **`Timer.Reset` pattern fixed** with explicit channel drain to prevent
      spurious fires after a stop.

### Hardening

- [x] **Sliding-window rate limiter:** 100 requests/minute per IP via a new
      `rateLimitMiddleware`.
- [x] **`BackgroundRunner.Stop` is idempotent** via `sync.Once` so double-stop
      does not close channels twice.
- [x] **Walk depth and size caps:** `MaxFileSize = 100MB`,
      `MaxDirDepth = 20`, `MaxTotalFiles = 10000–50000` in repository walks
      to bound resource usage on hostile inputs.
- [x] **Heartbeat and deploy FTP intervals** are now configurable through
      environment variables (no longer hardcoded).
- [x] **Frontend Authorization header** set on every API request with 10s/30s
      per-request timeouts.
- [x] **`AbortController` + `TimeoutError`** wired into the frontend HTTP
      client so cancelled/expired requests can be distinguished from
      application errors.
- [x] **`refreshGenerationRef`** prevents races during `switchSession`
      initiated refreshes.
- [x] **`visibilitychange` handler** pauses polling while the tab is hidden
      and resumes on return.
- [x] **`mountedRef` + `AbortController`** in `WorkspaceBootstrap` to prevent
      `setState` after unmount.
- [x] **Sidebar guard** displays a warning when no active session is selected.
- [x] **Servers page** clears password fields immediately after submit and
      uses `autoComplete="new-password"`.
- [x] **SyncCenter** requires a confirmation checkbox before destructive
      sync actions; `executeSync` accepts a `confirm` flag.

### Quality

#### Error Handling

- [x] All `_, _ = .Exec/.Query/.rand.Read` silent error discards removed
      throughout the backend.
- [x] All `rows.Scan` return values checked.
- [x] `SetActive` now returns an `error` so callers cannot miss persistence
      failures.

#### Database

- [x] WAL mode enabled consistently — `journal_mode = WAL` set on **both**
      database connections.
- [x] Migrated session active-flag handling to the canonical history table.

#### Schema Cleanup

- [x] Removed duplicate `deployments` schema; consolidated historical writes
      into `deploy_logs` in `history/db.go`.
- [x] Removed unused `history.NewDB` export.

#### Frontend

- [x] **ApiError class** with helper predicates (`isUnauthorized`,
      `isForbidden`, `isTimeout`) replaces ad-hoc regex/string matching on
      error messages.
- [x] **ErrorBoundary** now uses `instanceof ApiError` instead of fragile
      regex detection.
- [x] **SessionContext** detects 401 responses via `ApiError.isUnauthorized`
      (no more substring matching).

### Tests

- [x] **27 Go packages** pass under `go test -race -count=1` (all OK).
- [x] **61 Vitest tests across 15 frontend test files** added or stabilized.
- [x] Vitest suite runs clean of failures; pre-existing React Router
      warnings remain and are tracked separately.
- [x] All packages have race-free coverage for the previously-flagged
      concurrency paths (`activeContext`, `BackgroundRunner.Stop`,
      `Timer.Reset`).

### Removed

- Duplicate `deployments` schema (consolidated into `deploy_logs`).
- Unused `history.NewDB` public export.
- Hardcoded vault password fallback in configuration loader.
- `os.Exit(0)` from the heartbeat loop.
- Regex-based error classification in the frontend `ErrorBoundary` and
  `SessionContext`.
- Substring matching fallback in 401 detection.
- `_ =` discards on `Exec`, `Query`, `rand.Read`, and `rows.Scan` calls.

### Known Issues

- **Symlink walk tests are skipped on Windows.** `os.Symlink` requires
  Developer Mode or admin rights; CI on Windows uses a no-op stub. Full
  coverage runs on Linux/macOS runners.
- **MinGW linker warnings** on Windows builds of CGO-enabled tooling are
  cosmetic and originate from the toolchain, not FileENIAC. They do not
  affect produced binaries.
- Pre-existing React Router `Future Flag` warnings in the Vitest console
  remain — independent of this release and tracked in the frontend backlog.

[Unreleased]: https://example.com/fileeniac/compare/v0.0.0...HEAD
