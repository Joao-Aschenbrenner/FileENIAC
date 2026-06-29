# SPDX-License-Identifier: MIT
# FileENIAC Privacy Policy

**Last updated**: 2026-06-28
**Version**: 1.0

## Overview

FileENIAC is a local-first desktop application. All your data stays on your
device. We do not collect, store, sell, or share any personal information.

## Data We Do NOT Collect

FileENIAC does not:

- Send telemetry or analytics data of any kind
- Track usage or behavior
- Contact any remote servers except those you explicitly configure
  (your GitHub account and your FTPS servers)
- Store data in cloud services operated by the developer
- Use cookies or tracking technologies
- Collect device identifiers
- Log keystrokes or screen content

## Data You Store Locally

When you use FileENIAC, the following data is stored **locally on your device**:

### Workspace and Project Data

- Workspace directory paths
- Project names and configurations
- Server profiles (host, port, username — password stored in Vault)
- Deploy history records
- Sync operation logs

### Authentication Data

- GitHub OAuth tokens (stored in Vault, encrypted)
- FTPS credentials (stored in Vault, encrypted)
- API session tokens (stored in SQLite, per-session)

### Application Logs

- Structured logs written locally to the configured log directory
- May contain operation details, file paths, and error messages
- Log files are retained until you delete them

## How Data Is Protected

### Encryption at Rest

Sensitive credentials (GitHub tokens, FTPS passwords) are encrypted using
AES-256-GCM. The encryption key is derived from the `FILEENIAC_VAULT_PASSWORD`
environment variable using Argon2id. **If you do not set this variable, the
Vault falls back to no encryption.**

### Local Storage Only

All workspace data is stored in the directory you configure via
`FILEENIAC_DATA_DIR` (default: `./data`). No data is automatically sent to
the developer or any third party.

### No Remote Sync

FileENIAC does not sync your data to any cloud service operated by the
developer. Your workspace data remains entirely under your control.

## Data Sharing

The only scenarios in which data leaves your device:

1. **GitHub OAuth**: When you authenticate with GitHub, the OAuth token
   is exchanged directly between your device and GitHub's servers.
   FileENIAC stores only the resulting access token locally.

2. **FTPS Operations**: When you deploy files via FTPS, file data is
   transferred directly between your device and your configured FTPS server.
   The developer has no access to this data.

3. **Git Operations**: When you clone or interact with Git repositories,
   data is transferred directly between your device and the Git provider
   (GitHub, GitLab, etc.).

## Logs

Application logs are written to the filesystem at the path configured in
the observability settings. Logs may contain:

- Operation names and steps
- File paths being processed
- Error messages and stack traces
- Correlation IDs for request tracing

Logs are **never sent automatically** to the developer or any remote server.
To submit logs as part of a bug report, you must explicitly export and attach
them.

**How to delete logs**: Delete the log directory or configure a log rotation
policy in the application settings.

## Version Updates

FileENIAC may check for new versions by making a request to the GitHub Releases
API (`api.github.com/repos/Joao-Aschenbrenner/FileENIAC/releases`). This request
does not include any personal data. It only retrieves version metadata to
determine if an update is available.

## Your Rights (LGPD / GDPR / CCPA)

Because all data is stored locally on your device, you have full control:

- **Access**: All your data is in the `FILEENIAC_DATA_DIR` directory
- **Correction**: Edit the SQLite database or configuration files directly
- **Deletion**: Delete the `FILEENIAC_DATA_DIR` directory to remove all data
- **Portability**: Export workspace data by copying the data directory
- **Revocation**: Delete the data directory or individual records

To request deletion of data associated with a support request, contact the
developer with the issue ID — no automatic data retention exists.

## Data Breach

In the event of a local data breach (e.g., malware accessing your device),
the impact is limited to whatever malware can access on your device.
FileENIAC does not maintain any remote backup or centralized database that
could be breached.

## Changes to This Policy

If this privacy policy changes, the change will be reflected in the
`docs/legal/PRIVACY_POLICY.md` file in the repository, with an updated
"Last updated" date. No personal notification will be sent for minor changes.

## Contact

For privacy-related questions or to report concerns:
- Open an issue at: https://github.com/Joao-Aschenbrenner/FileENIAC/issues

## Summary

| Data Type | Stored Where | Sent to Developer |
|-----------|-------------|-------------------|
| Workspace paths | Your device | No |
| Project config | Your device | No |
| Server profiles | Your device | No |
| GitHub tokens | Your device (Vault) | No |
| FTPS passwords | Your device (Vault) | No |
| Deploy history | Your device | No |
| Application logs | Your device | No |
| GitHub OAuth exchange | GitHub servers | No |
| FTPS file transfers | Your FTPS server | No |