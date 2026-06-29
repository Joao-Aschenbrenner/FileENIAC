# SPDX-License-Identifier: MIT
# FileENIAC Data Processing

This document describes exactly where and how data flows in FileENIAC.
All statements are verifiable in the source code.

## Architecture Overview

FileENIAC is a local-first desktop application. Data processing happens
entirely on the user's device, except when explicitly configured to interact
with external services (GitHub, FTPS servers).

## Data Storage Locations

### 1. SQLite Database

**Location**: `{FILEENIAC_DATA_DIR}/workspaces.db` (default: `./data/workspaces.db`)

**Contains**:
- Workspaces: name, path, created_at, updated_at
- Projects: workspace_id, name, path, git_url, mirror_path, created_at
- Servers: workspace_id, name, host, port, username (password is in Vault)
- Sessions: workspace_id, github_user, workspace_path, is_active, created_at
- History: operation type, project_id, status, details JSON, created_at
- Events: session_id, event_type, details JSON, created_at

**Retention**: Until explicitly deleted by the user (delete workspace or
delete data directory).

**Sent to Developer**: Never. Database stays entirely on the user's device.

### 2. Vault (Encrypted Credentials)

**Location**: `{FILEENIAC_DATA_DIR}/vault.db` (default: `./data/vault.db`)

**Contains**:
- GitHub personal access tokens
- FTPS passwords
- API session tokens

**Encryption**: AES-256-GCM with unique 256-bit key generated per workspace

**Retention**: Until explicitly deleted by the user.

**Sent to Developer**: Never. Not transmitted over the network.

### 3. Workspace Filesystem

**Location**: User-configured workspace directory (e.g., `C:\Projects\myapp`)

**Contains**:
- Your project source code and files
- Git repository data (.git directory)
- Mirror copies (if configured)

**Retention**: Under your full control. FileENIAC does not delete workspace
files unless you explicitly initiate a delete operation.

**Sent to Developer**: Never. The developer has no access to your workspace.

### 4. Application Logs

**Location**: `{FILEENIAC_DATA_DIR}/logs/` (default: `./data/logs/`)

**Contains**:
- Structured JSON logs (Zap)
- Correlation IDs
- Operation steps and timing
- Error messages

**Format**: JSON lines (Zap JSON encoder)

**Retention**: Until deleted by the user or log rotation policy.

**Sent to Developer**: Never, unless you manually attach them to a bug report.

## Network Data Flows
### GitHub Personal Access Token

```

Your Device  →  GitHub API (api.github.com)
           ←  API Response
```

- FileENIAC sends your GitHub personal access token to GitHub's API
- The token is stored encrypted in the Vault
- **Developer has no access** to the token

### Git Operations

```
Your Device  →  Git Provider (github.com, gitlab.com, etc.)
            ←  Repository data
```

- Git clone/push/pull use the Git protocol
- Data travels directly between your device and the Git provider
- **Developer has no access** to this data

### FTPS Operations

```
Your Device  →  Your FTPS Server
             ←  Directory listings
             ←  File transfers
```

- Files are transferred directly between your device and your FTPS server
- The developer provides the FTPS client library (jlaffaye/ftp)
- **Developer has no access** to data in transit

### Version Check

```
Your Device  →  GitHub Releases API
            ←  Version tags and release info
```

- Performed on startup (optional, configurable)
- Only checks for new version tags
- No personal data in the request
- **Developer has no access** to the request

## Data Not Collected

The following do not exist in FileENIAC:

- Telemetry collection
- Analytics or tracking
- Crash reporting services
- Remote debugging
- Cloud backup
- Multi-device sync
- User accounts on developer-operated servers

## Data Flow Diagram

```
User Configures Workspace
         ↓
Workspace Path Stored in SQLite (workspaces table)
         ↓
GitHub Personal Access Token Stored in Vault (encrypted)
         ↓
FTPS Credentials Stored in Vault (encrypted)
         ↓
Deploys/Syncs Read/Write Workspace Files Directly
         ↓
History Written to SQLite (history table)
         ↓
Logs Written to FILEENIAC_DATA_DIR/logs/
         ↓
NO DATA SENT TO DEVELOPER
```

## Data Deletion

To delete all data:

```bash
# Stop FileENIAC
rm -rf "$FILEENIAC_DATA_DIR"  # default: ./data/
```

This removes:
- SQLite database (workspaces, projects, servers, sessions, history)
- Vault (GitHub tokens, FTPS passwords)
- Log files
- Any cached data

## Credential Management

To rotate GitHub credentials:
1. Revoke the token in GitHub Settings → Developer settings → Personal access tokens
2. Delete the entry in FileENIAC (or clear the vault)
3. Re-authenticate in FileENIAC

To rotate FTPS credentials:
1. Update the password in FileENIAC server settings
2. The new password is encrypted and stored in the Vault