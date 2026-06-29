# SPDX-License-Identifier: MIT
# FileENIAC — LGPD Compliance Statement

** Lei Geral de Proteção de Dados (Brazilian General Data Protection Law) **

**Document version**: 1.0
**Last updated**: 2026-06-28

## Applicability

This statement applies to users of FileENIAC who are located in Brazil or
whose personal data is processed in connection with FileENIAC's services.

FileENIAC is a local-first desktop application. As described in this
document, the vast majority of data processing occurs locally on the user's
device. This statement is provided in accordance with Brazilian LGPD
(Lei nº 13.709/2018) requirements.

## Data Controller

**Controller**: Joao Aschenbrenner
**Contact**: https://github.com/Joao-Aschenbrenner/FileENIAC/issues

The controller is responsible for the personal data processing operations
described in this document.

## Legal Basis for Processing

FileENIAC processes personal data only in the following circumstances:

### 1. User Consent (Art. 7, I — LGPD)

When you store GitHub credentials or FTPS credentials for deployment,
you actively configure these in the application. This constitutes consent
to store and use these credentials for the purpose of deployment and
Git operations.

### 2. Legitimate Interest (Art. 7, IX — LGPD)

Application logs may contain operational metadata (file paths, operation
names, timestamps) for security monitoring and issue diagnosis. This
processing is minimal and local-only.

### 3. Contract Performance (Art. 7, V — LGPD)

Workspace metadata (project names, workspace paths) is processed to
provide the core functionality of workspace management you requested.

## Categories of Personal Data Processed

The following table shows what personal data may be processed by FileENIAC
and where it is stored:

| Data | Stored In | Location | Sent to Developer |
|------|-----------|----------|-------------------|
| GitHub username | Vault (encrypted) | Your device | No |
| GitHub token | Vault (encrypted) | Your device | No |
| FTPS username | Vault (encrypted) | Your device | No |
| FTPS password | Vault (encrypted) | Your device | No |
| Workspace paths | SQLite | Your device | No |
| Project names | SQLite | Your device | No |
| Server hostnames | SQLite | Your device | No |
| Deploy history | SQLite | Your device | No |
| Application logs | Log files | Your device | No |

**Important**: No financial data, government IDs, biometric data, or
sensitive categories under Art. 5, II of LGPD are processed by FileENIAC.

## Purpose Limitation (Art. 46 — LGPD)

Personal data is used only for:

- Authenticating with GitHub for repository operations
- Connecting to configured FTPS servers for deployment
- Maintaining workspace and project configuration
- Recording deploy history for audit purposes

Data is **never** used for:
- Profiling or automated decisions
- Marketing or advertising
- Sharing with third parties for commercial purposes
- Building user profiles

## Minimization and Necessity (Arts. 6, III and 46 — LGPD)

FileENIAC collects only what is strictly necessary:

- GitHub credentials: required for GitHub OAuth authentication
- FTPS credentials: required for deployment operations
- Workspace paths: required for workspace management
- Project names: required for project identification

The application does not collect any data beyond what is configured by you.

## Data Security (Art. 46 — LGPD)

FileENIAC implements the following security measures:

### Technical Safeguards

- **Encryption at rest**: Credentials in the Vault are encrypted with
  AES-256-GCM using a key derived from `FILEENIAC_VAULT_PASSWORD`
  via Argon2id
- **Path canonicalization**: Prevents path traversal attacks in workspace
  operations
- **Input validation**: Regex allowlists on names and paths
- **SQL injection protection**: Table names use an allowlist in all queries
- **URL credential stripping**: Git credentials are stripped from log output
- **HTTP security**: Timeouts (30s read, 10m write) and 1MB body limit

### Organizational Safeguards

- No hardcoded credentials in source code
- Secrets obtained only from environment variables
- Bearer tokens are ephemeral (32 bytes, generated per session)

**Note on Vault Encryption**: If `FILEENIAC_VAULT_PASSWORD` is not set,
the Vault does not apply encryption. Set this environment variable for
production use.

## Retention Period (Art. 7, II — LGPD)

Data is retained until you delete it:

- Workspace data: Retained until workspace is deleted or data directory
  is removed
- Vault credentials: Retained until explicitly cleared in the application
- Session data: Retained until session is deleted or application data is
  cleared
- Logs: Retained until log directory is deleted or log rotation occurs

There is no automatic deletion schedule. You control retention by managing
your data directory.

## Data Subject Rights (Arts. 17-22 — LGPD)

As a data subject, you have the following rights:

### Right of Access (Art. 18)

You can access all your data by examining:
- SQLite database at `{FILEENIAC_DATA_DIR}/workspaces.db`
- Vault at `{FILEENIAC_DATA_DIR}/vault.db`
- Log files at `{FILEENIAC_DATA_DIR}/logs/`

### Right of Correction (Art. 18, IV)

You can correct data by:
- Editing SQLite directly with any SQLite client
- Re-configuring credentials in the application
- Deleting and re-creating workspaces

### Right of Deletion (Art. 18, VI)

To delete all personal data:
```bash
rm -rf "$FILEENIAC_DATA_DIR"  # default: ./data/
```

This removes all workspace, credential, history, and log data.

To delete specific records, use the application's delete functions or
edit the SQLite database directly.

### Right of Portability (Art. 18, V)

Export your data by copying the data directory:
```bash
cp -r "$FILEENIAC_DATA_DIR" /path/to/backup/
```

The SQLite database can be exported to SQL format using the `.dump` command
in any SQLite client.

### Right to Revoke Consent (Art. 8, §5)

To revoke consent for credential storage:
1. Delete the credentials in the application settings
2. Revoke the GitHub token in GitHub Settings
3. Change the FTPS password on the server

### Right to Information (Art. 9)

This document fulfills the obligation to inform about data processing.
For additional questions, open an issue on the repository.

## International Data Transfers (Art. 33 — LGPD)

FileENIAC may result in international data transfers in the following cases:

1. **GitHub**: When you authenticate with GitHub or push/pull repositories,
   data is processed by GitHub (US-based) under GitHub's privacy policy.

2. **FTPS Servers**: If your FTPS server is hosted outside Brazil, data
   is transferred to that jurisdiction.

These transfers occur **only** when you explicitly configure these services
and are governed by the respective service providers' terms.

**The developer has no control over and takes no responsibility for these
third-party transfers.** You should review the privacy policies of GitHub
and your hosting provider.

## Incident Notification (Art. 48 — LGPD)

FileENIAC does not maintain any centralized database of user data. Therefore,
there is no centralized repository that could be breached by the developer.

In the event of a security incident on your local device (e.g., malware):
- The incident affects only data on your device
- The developer cannot notify you of local device breaches
- You are responsible for device security

If a support request involves sharing data (e.g., logs), you do so
voluntarily and explicitly.

**No automatic data breach notification process exists** because no
centralized user database exists.

## Children's Data (Art. 14 — LGPD)

FileENIAC is not directed at children. The application does not knowingly
collect data from users under 18 years of age. If you become aware that
a child has used the application, contact the developer to have their
data removed.

## Changes to This Statement

Material changes to this LGPD statement will be posted at
`docs/legal/LGPD.md` in the repository. Changes take effect upon posting.

## Contact

For LGPD-related requests or to exercise your data subject rights:
- Open an issue: https://github.com/Joao-Aschenbrenner/FileENIAC/issues

Response time: We aim to respond within 30 days of receiving a request.
Due to the local-first nature of the application, most requests can be
self-served by deleting the data directory.