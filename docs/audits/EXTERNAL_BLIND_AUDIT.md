# External Blind Audit — FileENIAC

## Summary

FileENIAC is a local-first workspace management platform with a Go backend (REST API) and Tauri v2 desktop frontend (React/TypeScript). The project demonstrates strong architectural awareness with ADR documentation, a clean transport abstraction layer, comprehensive test coverage, and thorough legal documentation (MIT License, LGPD, Terms of Use). However, this blind audit found two critical security blockers: bearer token authentication is claimed as implemented in the changelog and relied upon by the frontend, but no auth middleware exists on the Go API server, and the handshake endpoint that the frontend calls (`/api/_handshake/token`) is not defined anywhere in the backend routes. Additional high-severity issues include duplicate `init()` functions in `root.go` that will cause runtime panics, severe version string inconsistency across the project, and a `docker-compose.yml` referencing a now-removed environment variable. The release is blocked.

## Findings

### Critical (blocking release)

1. **Bearer token authentication middleware does not exist despite being documented and depended upon**
   - The CHANGELOG (v0.1.2, lines 124-134) explicitly states: "Bearer token authentication middleware enforced on every route except `/api/health`"
   - The frontend client (`apps/desktop/src/api/client.ts:39,69,100`) sends `Authorization: Bearer ${token}` on every request and expects a 401 on invalid tokens
   - The CORS middleware (`backend/internal/api/api.go:135`) allows the `Authorization` header — but no middleware ever validates it
   - `backend/internal/api/api.go:96-97` shows the actual handler chain: `s.rateLimitMiddleware(s.corsMiddleware(s.mux))` — there is NO auth middleware between rate-limiting and the mux
   - A grep for `_handshake`, `bearer`, `Bearer`, `Authorization`, `authMiddleware` across `backend/internal/api/` and `backend/cmd/` confirms zero auth validation logic
   - **Impact**: Any process on the local machine can call the full API without any token. All deployed credentials (GitHub tokens, FTPS passwords) and workspace data are unprotected.

2. **Frontend handshake endpoint `/api/_handshake/token` does not exist in backend routes**
   - The frontend `tokenStorage.ts:31` calls `GET /api/_handshake/token` to bootstrap the bearer token
   - The backend `routes()` method in `api.go:58-94` defines 34 routes but none match `/api/_handshake/token`
   - Since the frontend handler (`frontend.go`) catches `/` only for non-`/api` paths, this falls through to an unregistered path — likely a 404
   - **Impact**: The entire token bootstrapping mechanism is non-functional. Combined with finding #1, there is no authentication path.

### High

3. **Duplicate `init()` functions in `root.go` will cause runtime panic**
   - `backend/cmd/root.go:36-52` defines an `init()` adding 15 commands to `rootCmd`
   - `backend/cmd/root.go:61-75` defines a SECOND `init()` that re-adds the same 15 commands
   - `backend/cmd/version_cmd.go:20-22` defines a THIRD `init()` adding `VersionCmd` again
   - Cobra's `AddCommand` is not idempotent when called with the same pointer — this will cause a panic at startup: `panic: command "version" already added`
   - **Evidence**: lines 36-52 and 61-75 of `root.go` contain duplicate command registrations that cannot both execute safely

4. **Version string inconsistency across 5+ files**
   - `backend/internal/version/version.go:4`: `Version = "0.2.0"`
   - `NOTICE:13`: `Version: 0.1.3`
   - `tauri.conf.json:4`: `"version": "0.1.3"`
   - `package.json:3` (root): `"version": "1.0.0"`
   - `apps/desktop/package.json:4`: `"version": "1.0.0"`
   - `RELEASE_NOTES.md:1`: `v0.1.2` in header
   - `CHECKSUMS.txt:1`: `v0.1.0`
   - `CHANGELOG.md:8`: `[0.1.3]`
   - **Impact**: Users and integrators cannot determine the actual version. Build artifacts will have mismatched identifiers.

5. **`docker-compose.yml` requires removed environment variable `FILEENIAC_VAULT_PASSWORD`**
   - `docker-compose.yml:14`: `FILEENIAC_VAULT_PASSWORD=${FILEENIAC_VAULT_PASSWORD:?required}`
   - CHANGELOG v0.1.3 (line 19-20): "Removed references to the non-existent `FILEENIAC_VAULT_PASSWORD` environment variable. Vault uses auto-generated 256-bit AES-GCM key per workspace."
   - The vault code (`vault.go` and `workspace/registry.go`) confirms keys are auto-generated via `vault.GenerateKey()` — the env var is never read
   - Docker Compose will refuse to start unless a now-ignored variable is set, and worse, users may believe they need to supply a vault password that is silently ignored

6. **`cmd_test.go` references auth subcommands that don't exist**
   - `backend/cmd/cmd_test.go:277-283` tests for subcommands `"set"`, `"get"`, `"remove"` under `auth`
   - But `backend/cmd/auth.go` defines subcommands named `"login"`, `"status"`, `"logout"`
   - **Impact**: Test `TestAuthCmd_HasSubcommands` fails for all three subcommands it checks

### Medium

7. **`native.go` uses unbounded `select {}` — no graceful shutdown path**
   - `backend/cmd/native.go:111`: `select {}` blocks forever with no signal handling, no context cancellation, and no shutdown hook
   - The server cannot be gracefully stopped (no SIGINT/SIGTERM handling)
   - **Impact**: If run as a service, killing the process may leave resources (database connections, background runners) in an inconsistent state

8. **`Access-Control-Allow-Origin: *` CORS policy on local API**
   - `backend/internal/api/api.go:133`: `w.Header().Set("Access-Control-Allow-Origin", "*")`
   - For a local desktop application binding to `localhost`, this is permissive but not directly exploitable from remote (since the server listens on localhost). However, any local webpage could make authenticated requests to the API, and combined with the missing auth middleware (finding #1), there is zero protection against CSRF from the same machine

9. **Workspace context passed via query parameter — visible in logs and referrers**
   - `backend/internal/api/api.go:188-189`: The workspace path is accepted as `r.URL.Query().Get("workspace")`
   - The frontend (client.ts:120-122) encodes the workspace path in the URL query string
   - Workspace paths may reveal usernames, project structure, and organizational information
   - While logged (via zap), the logs explicitly warn against pasting credentials; paths may still leak sensitive directory names

10. **`serve` command does not serve frontend, unlike `desktop` and `native` commands**
    - `backend/cmd/serve.go:27`: `srv := api.New(addr)` — no call to `srv.ServeFrontend()`
    - `backend/cmd/desktop_windows.go:27-28`: calls `srv.ServeFrontend()`
    - `backend/cmd/native.go:71-72`: calls `srv.ServeFrontend()`
    - **Impact**: The primary `serve` subcommand only serves the API, not the bundled frontend. Users expecting a single binary will find an API with no UI

### Low

11. **`docker/backend.Dockerfile` runs `go run` instead of building a binary**
    - `docker/backend.Dockerfile:16`: `CMD ["go", "run", "./backend"]`
    - This compiles from source at every container start, increasing startup time and requiring the Go toolchain in the runtime image
    - The primary `Dockerfile` correctly uses multi-stage build; this secondary Dockerfile appears to be a development convenience image but lacks documentation explaining the difference

12. **Frontend `client.ts` `ws()` function mixes localStorage and Tauri IPC patterns**
    - `apps/desktop/src/api/client.ts:119-123`: Directly reads `localStorage.getItem("eniac_ws_path")` instead of using the `storage.ts` abstraction
    - Inconsistent with the rest of the client code that uses `storageGet(STORAGE_KEYS.workspacePath)` via the storage module
    - Not a security issue but an architectural inconsistency

13. **`backend/main.go` and root `main.go` duplicate entry point**
    - Root `main.go` imports `backend/cmd` and calls `cmd.Execute()` (with no log sync)
    - `backend/main.go` imports `backend/cmd` and `backend/internal/log` and calls `defer log.Sync(); cmd.Execute()`
    - Having two `main.go` files in different locations with the same purpose is confusing; the root one is slightly inferior (no `log.Sync()`)

14. **`CHECKSUMS.txt` references v0.1.0 but CHECKSUMS-v0.1.2.txt and CHECKSUMS-v0.1.3.txt also exist**
    - Only the most recent checksum file should be maintained for the current release
    - The generic `CHECKSUMS.txt` (v0.1.0) is stale and could mislead users about the correct binary checksums

### False Positives (items that look wrong but are correct)

1. **`SECURITY.md` says no required environment variables, but docker-compose requires `FILEENIAC_VAULT_PASSWORD`**
    - The `SECURITY.md:15-22` correctly states no env vars are required because vault keys are auto-generated
    - The `docker-compose.yml` requiring `FILEENIAC_VAULT_PASSWORD` is itself a bug (finding #5), not a conflict in SECURITY.md
    - Explanation: SECURITY.md was updated in v0.1.3 to match the new auto-generated key design; docker-compose.yml was not updated accordingly. SECURITY.md is correct.

2. **CSP policy in tauri.conf.json allows `'unsafe-inline'` for styles**
    - `tauri.conf.json:25`: `"style-src 'self' 'unsafe-inline'"`
    - This is standard for TailwindCSS/React applications and not a vulnerability since the WebView loads only local content (`'self'`)
    - No remote script injection is possible because `script-src 'self'` restricts script execution to bundled files

3. **`backend.exe` binary committed to the repository**
    - `backend/backend.exe` exists as a compiled binary
    - This is unusual but the `.gitignore` still includes `*.exe` while the log says `backend.exe` is not in `.gitignore`
    - However, `.dockerignore` explicitly excludes `backend.exe`, so this is likely a one-time committed artifact rather than a pattern

### Inconclusive (needs human review)

1. **Vault encryption key stored in workspace config file (`.eniac/config.toml`)**
    - The vault generates a key via `crypto/rand` and stores it in plaintext in the workspace TOML config
    - `workspace/registry.go:74-87`: key generated, stored in config alongside workspace name
    - While the vault encrypts credentials (AES-256-GCM), the key's storage alongside the encrypted data means that anyone with filesystem access to the data directory can decrypt all vault contents
    - This is acknowledged in docs ("The encryption key is stored in the workspace config") and is acceptable for a single-user desktop app where the security boundary is the device itself
    - Needs human judgment on whether this is acceptable for the target use case or if OS-level keychain integration (Windows Credential Manager, macOS Keychain) should be added

2. **`go.mod` specifies `go 1.26` which does not yet exist as a stable Go release**
    - `go.mod:3`: `go 1.26`
    - Python script `check_icons.py` present but not referenced in build targets
    - The Go 1.26 specification may be forward-looking/speculative; compatibility with the latest stable Go release (1.24 as of mid-2026) should be verified

## Recommendations

- **BLOCK RELEASE** — critical security issues found

## Plan

Priority-ordered corrections:

1. **CRITICAL — Implement bearer token auth middleware** (`backend/internal/api/api.go`)
   - Add a `authMiddleware(next http.Handler) http.Handler` function that reads the `Authorization: Bearer <token>` header and validates against the token stored/supplied via `FILEENIAC_API_TOKEN` env var
   - Wrap the handler chain: `s.authMiddleware(s.rateLimitMiddleware(s.corsMiddleware(s.mux)))`
   - Exempt `/api/health` and `/api/_handshake/token` from auth
   - Implement the `/api/_handshake/token` GET endpoint that returns a one-time or ephemeral token

2. **CRITICAL — Implement `/api/_handshake/token` endpoint** (`backend/internal/api/api.go`)
   - Add route: `s.mux.HandleFunc("/api/_handshake/token", s.handleHandshakeToken())`
   - This endpoint must be unauthenticated and return a bearer token for use by the frontend

3. **HIGH — Fix duplicate `init()` functions** (`backend/cmd/root.go`)
   - Remove the second `init()` function (lines 61-75)
   - Remove the duplicate `AddCommand(VersionCmd)` from `version_cmd.go:20-22` (already in root.go's init)
   - Verify all commands are registered exactly once

4. **HIGH — Unify version string across all files**
   - Set authoritative version in `backend/internal/version/version.go` to `0.1.3`
   - Update `NOTICE`, `tauri.conf.json`, `package.json` (root), `apps/desktop/package.json`, `RELEASE_NOTES.md` to match
   - Remove or update stale `CHECKSUMS.txt` (v0.1.0) and keep only the current release checksum file

5. **HIGH — Fix `docker-compose.yml` environment variable**
   - Remove `FILEENIAC_VAULT_PASSWORD=${FILEENIAC_VAULT_PASSWORD:?required}` since the vault no longer uses it
   - Replace with documentation about auto-generated keys, or remove entirely

6. **HIGH — Fix auth test expectations** (`backend/cmd/cmd_test.go:277-283`)
   - Change expected subcommands from `"set"`, `"get"`, `"remove"` to `"login"`, `"status"`, `"logout"` to match `auth.go`

7. **MEDIUM — Add graceful shutdown to `native.go`**
   - Replace `select {}` with signal handling (`os.Signal`, `context.WithCancel`)
   - Call `srv.Close()` on shutdown to release database connections and background runners

8. **MEDIUM — Restrict CORS origin for local API**
   - Change `Access-Control-Allow-Origin: *` to `http://localhost:*` or remove the wildcard since the frontend and backend are co-located

9. **MEDIUM — Add `ServeFrontend()` call to `serve` command** (`backend/cmd/serve.go`)
   - Add `srv.ServeFrontend()` after `api.New(addr)` so the `serve` subcommand also serves the bundled UI

10. **LOW — Remove stale `CHECKSUMS.txt`** or update it to the current release

11. **LOW — Remove committed `backend.exe` binary** from the repository and add it to `.gitignore`
