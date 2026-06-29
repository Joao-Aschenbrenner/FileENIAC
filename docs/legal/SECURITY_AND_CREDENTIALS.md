# SPDX-License-Identifier: MIT
# FileENIAC — Security and Credential Management

**Version**: 1.0
**Last updated**: 2026-06-28

## Credential Types

FileENIAC stores the following credentials:

- **GitHub personal access token**: Used to authenticate with GitHub API
- **FTPS passwords**: Used to authenticate with your FTPS servers
- **API session tokens**: Used for desktop frontend/backend communication

## Protecting Your Credentials

### Never Share Credentials

- Never share your GitHub personal access token with anyone
- Never share your FTPS passwords
- Never post credentials in public issues, chat messages, or emails
- The developer will never ask for your credentials

### How to Rotate Your GitHub Token

1. Go to GitHub Settings → Developer settings → Personal access tokens
2. Generate a new token with the required scopes (repo, workflow)
3. Delete or update the old token in FileENIAC settings
4. Update the token in FileENIAC workspace settings

### How to Rotate Your FTPS Password

1. Change the password on your FTPS server
2. Update the password in FileENIAC server settings
3. The new password is encrypted and stored in the Vault

### Vault Encryption Key

FileENIAC generates a unique 256-bit AES-GCM encryption key per workspace.
This key is stored in the workspace configuration file (`.eniac/config.toml`).

To change the encryption key:
1. Export current credentials
2. Delete the workspace data directory
3. Re-create the workspace
4. Re-configure credentials

## Log Safety

FileENIAC is designed not to log secrets such as passwords, vault keys,
or tokens. However, logs may contain:

- File paths
- Operation names
- Error messages
- Timestamps

**Review logs before sharing them publicly.**
Never paste logs containing credentials, tokens, private repository names,
or sensitive file paths into public issues.

## How to Delete Logs

1. Locate your data directory (`FILEENIAC_DATA_DIR`)
2. Delete the `logs/` subdirectory
3. FileENIAC will create a new log directory on next startup

## Reporting a Vulnerability

See `VULNERABILITY_DISCLOSURE.md` for the coordinated disclosure process.
