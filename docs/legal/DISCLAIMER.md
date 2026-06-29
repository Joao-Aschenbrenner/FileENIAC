# SPDX-License-Identifier: MIT
# FileENIAC Disclaimer

**Version**: 1.0
**Last updated**: 2026-06-28

## No Warranty

FileENIAC is provided "AS IS", without warranty of any kind, express or
implied, including but not limited to the warranties of merchantability,
fitness for a particular purpose, and noninfringement.

The entire risk arising out of the use or performance of FileENIAC remains
with you.

## Limitation of Liability

In no event shall the authors or copyright holders be liable for any claim,
damages, or other liability, whether in an action of contract, tort, or
otherwise, arising from, out of, or in connection with FileENIAC or the
use or other dealings in FileENIAC.

This includes, without limitation:

- **Data loss**: FileENIAC does not maintain backups. You are responsible
  for backing up your workspace data and configuration.

- **Deployment failures**: Deploy operations can fail for many reasons
  (network issues, server configuration, permission problems). Verify
  deployments manually after execution.

- **Credential exposure**: If you do not set `FILEENIAC_VAULT_PASSWORD`,
  credentials are stored unencrypted. You accept the risk of storing
  credentials without encryption.

- **File corruption**: Ensure you have backups of important files before
  using workspace operations that modify files.

- **Service disruption**: The developer does not guarantee uptime,
  availability, or responsiveness of any services.

## Not a Backup Solution

FileENIAC is **not** a backup solution. It is a workspace management
and deployment tool. You are responsible for maintaining backups of:

- Workspace source code
- Project configurations
- Server credentials and settings
- Deploy history and logs

## Third-Party Services

FileENIAC interacts with third-party services you configure:

- **GitHub**: Subject to GitHub's Terms of Service and privacy policy
- **FTPS Servers**: Subject to your hosting provider's terms
- **Git Providers**: Subject to the respective provider's terms

The developer is not responsible for:
- Changes to third-party service APIs
- Service outages by third parties
- Data handling by third parties
- Violations of third-party terms of service

## Security Limitations

- **Local security**: If your device is compromised, FileENIAC data
  can be accessed. The developer cannot protect against local threats.

- **Credential vault**: Encryption requires `FILEENIAC_VAULT_PASSWORD`.
  Without it, credentials are stored in plain text.

- **No code signing**: The Windows installer is not code-signed.
  Windows SmartScreen may display warnings.

## Forward-Looking Statements

Roadmap items, planned features, and future versions described in
documentation, issues, or discussions are **forward-looking statements**.
They are subject to change and should not be interpreted as commitments.

The developer makes no guarantee that any planned feature will be
implemented in any specific timeframe or at all.

## Changes to This Disclaimer

This disclaimer may be updated. The "Last updated" date reflects the
most recent revision. Continued use of FileENIAC after any revision
constitutes acceptance of the revised disclaimer.

## Governing Law

This disclaimer is governed by the laws of Brazil, as detailed in the
Terms of Use.