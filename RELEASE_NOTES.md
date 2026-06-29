# FileENIAC v0.1.2 — Release Notes

**Release Date**: 2026-06-29
**Type**: Security & Stability Hotfix

## Overview

This release addresses critical findings from the post-v0.1.0 code audit. It
fixes the SQL-injection vector in `database.Count`, restores the pre-commit
gate, hardens the Docker image, resolves all `npm audit` findings and brings
the frontend test suite back to green. This replaces the v0.1.1 release, which
shipped with corrupted Go source files and failing frontend tests.

## Highlights

- SQL injection eliminated from `DB.Count`.
- Per-IP API rate limiting (120 req/min).
- Docker image runs as non-root user with Go 1.26.
- All backend tests pass including `-race`.
- All frontend tests pass; `npm audit` reports 0 vulnerabilities.
- Comprehensive audit report and fix plan published in `docs/`.

---

# FileENIAC v0.1.1 — Release Notes

**Release Date**: 2026-06-28 (superseded by v0.1.2)
**Type**: Open Source Governance + LGPD Compliance

---

## Overview

FileENIAC v0.1.0 is a local-first workspace management platform combining a Go backend REST API with a Tauri v2 desktop frontend. It provides unified management of Git repositories, GitHub integrations, and FTPS deployments through a single dashboard.

---

## Features

### Core Functionality

| Feature | Description |
|---------|-------------|
| Workspace Management | Local workspace lifecycle with discovery and registry |
| GitHub Integration | OAuth authentication, repository discovery, clone, import |
| Project Registry | CRUD operations with workspace association |
| Deploy Orchestration | FTPS deployment with fallback support, rollback, and verify |
| Sync Engine | Bidirectional mirror sync with diff-based change detection |
| History & Audit | Deploy logs, rollback logs, events, background health monitoring |

### Architecture

- **Transport Layer (ADR-014)**: Abstract transport interface separating protocol implementation from domain logic
- **Factory Pattern**: `transports.New(cfg)` — no switch/conditional dispatch
- **SQLite with WAL Mode**: Consistent database concurrency
- **REST API**: JSON over HTTP with bearer token authentication

### Observability

- **Structured Logging**: Zap JSON with correlation IDs
- **Pluggable Metrics**: Timer, Counter, Gauge with no-op default
- **Step-based Tracing**: Hooks for deploy, sync, mirror operations

### Security

- Bearer token authentication with ephemeral 32-byte tokens
- Vault password via environment variable only (no hardcoding)
- FTP TLS 1.2 + ECDHE-only cipher suites
- Input validation with regex allowlist and path canonicalization
- SQL injection protection via table allowlist
- URL credential stripping from logs

---

## System Requirements

### Backend

| Requirement | Value |
|-------------|-------|
| Go | 1.26+ |
| OS | Linux (amd64), Windows (amd64), macOS (amd64) |
| Database | SQLite 3 |

### Desktop Application

| Requirement | Value |
|-------------|-------|
| OS | Windows 10/11 x64 |
| Runtime | WebView2 (included in Windows 10/11) |
| Disk | ~100 MB |

### Docker

| Requirement | Value |
|-------------|-------|
| Docker | 20.10+ |
| Memory | 256 MB minimum |

---

## Installation

### Option 1: Desktop Installer

Download `FileENIAC_0.1.0_x64-setup.exe` and run the NSIS installer.
Requires Windows 10/11 x64. No administrator rights required (per-user install).

### Option 2: Docker

```bash
# Build
docker build . -t enisystems/fileeniac:v0.1.0

# Run
export FILEENIAC_VAULT_PASSWORD="your-vault-password"
docker compose up -d
```

### Option 3: Standalone Binary

```bash
# Download from releases
./fileeniac-v0.1.0.exe serve -a :8080
```

---

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `FILEENIAC_DATA_DIR` | No | `./data` | Workspace data directory |
| `FILEENIAC_VAULT_PASSWORD` | **Yes** | — | Vault encryption password |
| `FILEENIAC_API_TOKEN` | No | (none) | API bearer token |
| `FILEENIAC_API_PORT` | No | `8080` | Desktop API port |

---

## API Endpoints

### Public
- `GET /api/health` — Health check (no auth required)

### Protected (Bearer token required)
- `GET/POST /api/sessions` — Session management
- `GET/POST/DELETE /api/sessions/{id}` — Session by ID
- `POST /api/sessions/{id}/activate` — Activate session
- `POST /api/sessions/{id}/clear-workspace` — Clear session workspace
- `GET/POST /api/projects` — Project management
- `GET/POST/DELETE /api/servers` — Server management
- `GET/POST /api/settings` — Settings
- `GET /api/history` — History events
- `GET /api/events` — Events
- `POST /api/deploy` — Execute deploy
- `POST /api/rollback` — Execute rollback
- `POST /api/verify` — Verify deployment
- `GET /api/diff` — Diff local vs mirror
- `GET/POST /api/syncs` — List sync operations
- `POST /api/sync` — Execute sync
- `POST /api/mirror` — Create mirror
- `GET/POST /api/github/*` — GitHub integration
- `GET/POST /api/repositories` — Repository management

---

## Known Limitations

1. **Symlink walk tests skipped on Windows**: Developer Mode required for `os.Symlink`
2. **No SFTP/S3/WebDAV**: Only FTPS in v0.1.0
3. **No multi-user/RBAC**: Single-user desktop application
4. **Session persistence**: SQLite workspace database (not encrypted at rest)
5. **No code signing**: NSIS installer not code-signed (Windows SmartScreen may warn)

---

## Artifacts

| Artifact | Path | SHA256 |
|---------|------|--------|
| Windows Installer (NSIS) | `apps/desktop/src-tauri/target/release/bundle/nsis/FileENIAC_0.1.0_x64-setup.exe` | `168d122b0374f81c0f13ee01057c9908c352d6bd2cde79fb8abdb4561af05c3c` |
| Backend Binary (Windows) | `release/fileeniac-v0.1.0.exe` | `32415caf10536562df041e645bb1a96db12924a2173a26d47d4b34c32ce9d528` |
| Docker Image | `ensisystems/fileeniac:latest` | (see Docker Hub) |

---

## Changes Since RC1

### Added
- Complete session management REST API (`/api/sessions`)
- Frontend session integration (create, select, delete, activate, clear workspace)
- 8 new API integration tests for session endpoints

### Fixed
- Session management gap documented in ARCHITECTURE_AUDIT.md resolved
- `Session` type aligned between frontend and backend (snake_case JSON fields)

### Changed
- Release status updated from RC1 to v0.1.0

---

## Changes in v0.1.1

### Added
- **Legal Documents** (16 files in `docs/legal/`):
  - Privacy Policy (local-first, no telemetry, LGPD-aligned)
  - Terms of Use (MIT, commercial use, liability limitations)
  - LGPD Compliance Statement (Brazilian General Data Protection Law)
  - Data Processing Document (technical data flow)
  - Security Policy (vault encryption, input validation, best practices)
  - Vulnerability Disclosure Policy (severity classification, response timeline)
  - Code of Conduct (Contributor Covenant 2.1)
  - Contributing Guide (setup, testing, conventional commits, PR process)
  - Governance Document (ADR process, roles, decision-making)
  - Supported Versions Policy (SemVer support windows)
  - Disclaimer, Copyright, Installer EULA

- **SPDX Headers**: `// SPDX-License-Identifier: MIT` added to all source files:
  - 97 Go files (backend)
  - 75 TypeScript/TSX files (frontend)
  - 3 Rust files (Tauri desktop)

- **GitHub Templates** (`.github/`):
  - Issue templates: Bug Report, Feature Request, Security Report, Question
  - Pull Request Template with checklist
  - CODEOWNERS file
  - SECURITY.md policy

- **NSIS Installer**:
  - MIT License acceptance screen (via `tauri.conf.json` nsis.license)
  - Custom NSH hook with Terms of Use + Privacy Policy checkbox pages
  - Post-build patch script (`scripts/patch-installer-nsis.ps1`)

- **README.md**: New Legal section with document table and quick summary

### Changed
- Version bumped to v0.1.1
- CI workflow: Go 1.25 → Go 1.26
- Makefile: new targets `validate`, `patch-installer`, `build-installer`
- LICENSE, NOTICE, THIRD_PARTY_LICENSES.md fully updated and separated

### Documentation
- Complete open-source governance documentation created
- CHANGELOG.md updated with v0.1.1 additions

---

## License

FileENIAC is open-source software licensed under the **MIT License**.

See `LICENSE` for the full license text.
See `NOTICE` for copyright and project information.
See `THIRD_PARTY_LICENSES.md` for dependency licenses.

Copyright (c) 2024-2026 Joao Aschenbrenner / ENIAC Systems