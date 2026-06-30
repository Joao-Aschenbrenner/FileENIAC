# Changelog

All notable changes to FileENIAC will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.5] - 2026-06-30

### Desktop Self-Contained Backend

- **Backend as sidecar:** Go backend is now bundled as a Tauri sidecar binary (`fileeniac-backend.exe`). No need to run `fileeniac serve` manually.
- **Auto-start backend:** Tauri spawns the backend automatically on app startup with `--port 0` (dynamic port).
- **Token injection:** Rust generates the auth token and passes it to the Go backend via `ENIAC_API_TOKEN` env var. No more hardcoded tokens.
- **`get_backend_info` IPC:** New Rust command returns `{ base_url, token }` to the frontend. Replaces separate `get_api_port` and `get_api_token`.
- **Health wait with retry:** Onboarding page retries health check for up to 15 seconds instead of showing "execute fileeniac serve".
- **Backend cleanup:** Closing the app window kills the backend sidecar process (via `on_window_event` CloseRequested).
- **Sidecar logs:** Backend stdout/stderr are captured to `%LOCALAPPDATA%/com.eniacsystems.fileeniac/logs/backend.log` with token stripping.
- **CORS dynamic origin:** CORS now reflects the request origin for `localhost` and `127.0.0.1`.
- **Go serve flags:** `serve` command now accepts `--host` and `--port` separately (instead of `--addr`). Uses `--port 0` for dynamic port allocation. Prints `FILEENIAC_READY` line for IPC.
- **Build script:** `scripts/build-sidecar.ps1` compiles the Go backend and places it in `src-tauri/binaries/` with the correct target triple suffix.

## [0.1.4] - 2026-06-30

### Security Fixes (External Blind Audit)

- **Auth middleware:** Bearer token validated on all API routes (except `/api/health` and `/api/_handshake/token`).
- **Handshake endpoint:** `GET /api/_handshake/token` returns the API token.
- **Duplicate init() removed** from `root.go`.
- **docker-compose.yml:** Removed stale `FILEENIAC_VAULT_PASSWORD` reference.
- **Auth test fixed:** Updated expectations to match current subcommands.
- **Version unified:** `version.go` aligned with Cargo.toml.
- **CORS restricted** to `http://localhost`.
- **Graceful shutdown** in `native.go` via signal.Notify.
- **`serve` command** now calls `ServeFrontend()`.
- **`ENIAC_API_TOKEN`** passed as env var to Tauri subprocess.
- **Stale artifacts removed:** `CHECKSUMS.txt`, `CHECKSUMS-v0.1.2.txt`, `backend.exe`.

## [0.1.3] - 2026-06-29

### Legal Hardening

- [x] **Legal docs separated:** `LICENSE`, `TERMS_OF_USE.md`, `PRIVACY_POLICY.md`,
      `LGPD.md`, `DATA_PROCESSING.md`, `THIRD_PARTY_SERVICES.md`,
      `SECURITY_AND_CREDENTIALS.md`, `INSTALLER_NOTICE.md` — each document now
      has a distinct purpose, no mixing of MIT license with terms of use.
- [x] **SPDX corrected:** `SPDX-LICENSE-IDENTIFIER` → `SPDX-License-Identifier`
      in all legal documents.
- [x] **Vault encryption description fixed:** Removed references to the
      non-existent `FILEENIAC_VAULT_PASSWORD` environment variable. Vault uses
      auto-generated 256-bit AES-GCM key per workspace. No unencrypted fallback.
- [x] **OAuth → Personal Access Token:** Corrected all references to "GitHub OAuth"
      to accurately describe GitHub personal access token usage.
- [x] **Removed dangerous SQLite advice:** Replaced "edit SQLite database directly"
      with "use application settings" across LGPD and Privacy docs.
- [x] **Termination/Amendments fixed:** Terms no longer claim developer can
      "terminate" MIT license rights. Changes apply to future releases only.
- [x] **CCPA removed:** Changed to "Privacy Rights (LGPD/GDPR)".
- [x] **Logs warning added:** Explicit warning not to paste credentials in logs.
- [x] **THIRD_PARTY_SERVICES.md created:** Documents all third-party interactions.
- [x] **SECURITY_AND_CREDENTIALS.md created:** Rotation guides, log safety tips.
- [x] **INSTALLER_EULA.md → INSTALLER_NOTICE.md:** Rewritten to avoid conflict
      with MIT license (no "revocable license" or "do not redistribute" clauses).
- [x] **README updated:** Added link section to all legal documents.
- [x] **DISCLAIMER.md corrected:** Removed `FILEENIAC_VAULT_PASSWORD` references,
      fixed "Continued use" clause.

## [0.1.2] - 2026-06-29

### Security

- [x] **SQL injection fixed:** removed dynamic `WHERE` clause from `database.Count`;
      prepared statements are now used throughout the backend.
- [x] **API rate limiting:** per-IP token bucket (120 req/min) applied to all routes.
- [x] **Container hardening:** Docker images now run as non-root user `fileeniac`.

### Build & Quality

- [x] **Husky restored:** pre-commit now runs `gofmt`, `go vet`, `go build` and
      `go test` from the `backend/` directory; Go tab indentation is allowed.
- [x] **Docker updated:** builder/runtime images use Go 1.26.
- [x] **npm audit clean:** upgraded `vite` to v8, `vitest` to v3, added
      `lucide-react`; 0 known vulnerabilities.

### Tests

- [x] **Frontend tests green:** updated `ErrorBoundary`, `Sidebar`, `ThemeContext`
      and added `Onboarding` tests; removed obsolete `LGPDConsent` and
      `Sidebar.workspace-guard` tests.
- [x] **Backend tests green:** all packages pass `go test ./...` and
      `go test -race ./...`.

### Documentation

- [x] **Post-release audit:** added `docs/audits/FULL_CODE_AUDIT_v0.1.0.md`.
- [x] **Fix plan:** added `docs/plans/FIX_PLAN_v0.1.0_AUDIT.md`.

### Bug Correction

- [x] **needsDelete corrigido:** `mirror_to_local` + `StateNew` tratava arquivos novos
      no mirror como deleção (deveria ser cópia). `mirror_to_local` + `StateRemoved`
      invertido para refletir deleção correta. Testes atualizados.

### Reliability

- [x] **TD-001 — Data Race eliminada:** `api.go` — acesso concorrente a `s.srv`
      entre `ListenAndServe`/`ListenDynamic` (escrita) e `Close`/`Addr` (leitura).
      Corrigido com `sync.RWMutex`. Nenhuma race detectada em `go test -race ./...`.

### Build

- [x] **TD-002 — Build reprodutível:** `backend/webui/dist/` removido do `.gitignore`
      e adicionado ao versionamento. O `//go:embed dist` funciona em checkout limpo.

### Dead Code Removal

- [x] **history/crud.go, db.go, record.go, crud_test.go removidos:**
      Comprovadamente mortos — zero referências externas a CRUD, DB, DeployRecord.
      A Service atual usa database.DB para deploy_logs, rollback_logs e events.

### Test Coverage

- [x] **Sprint 2 — Engine Validation (3ef45e8):** 40+ novos testes para Mirror,
      Sync, Diff e History. Mock Transport para Mirror (sem FTP real). Tabela de
      8 combinações para `needsDelete`. Bug documentado: `mirror_to_local` + `StateNew`
      inverte cópia com deleção.
- [x] **Sprint 4 — Transport Layer Coverage:** 14 novos testes para Registry,
      Factory e FTP Adapter — registro, lookup duplicado, protocolo inválido,
      configuração, erros de conexão, operações sem conexão.
- [x] **Sprint 4 — CLI Coverage:** 22 novos testes — help, flags inválidas,
      subcomandos inválidos, estrutura de comandos, flags obrigatórias.

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
- [x] **Arquitetura validada (Sprint 1.5):** `go build ./...` + `go vet ./...` + `go test ./...` + `go test -race ./...`
      verdes. Transport Layer auditada: `jlaffye/ftp` isolado em `deploy/ftp/` e `transports/ftp/`.

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
