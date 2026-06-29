# SPDX-License-Identifier: MIT
# FileENIAC Supported Versions

**Version**: 1.0
**Last updated**: 2026-06-28

FileENIAC uses semantic versioning (SemVer). This document defines the
support window for each version.

## Release Schedule

New minor versions are released approximately every 2-4 weeks during active
development. Patch versions are released as needed for critical bug fixes
and security patches.

## Version Support Policy

### Current Stable Release

| Version | Status | Go Version | Support Level |
|---------|--------|------------|---------------|
| v0.1.x | **Stable** | 1.26+ | Bug fixes, security patches |
| v0.2.x | Development | 1.26+ | New features, breaking changes |

### Version Support Definition

| Support Level | Description |
|--------------|-------------|
| Bug fixes | Patches for bugs that affect normal operation |
| Security patches | Fixes for disclosed vulnerabilities |
| New features | New functionality added in backward-compatible way |
| Breaking changes | Changes that require user action to upgrade |

## v0.1.x (Current Stable)

- **Status**: Stable
- **End of Life**: When v0.2.0 is released and stable
- **Support**: Bug fixes and security patches only
- **Upgrade path**: Users should plan to upgrade to v0.2.x when released

### v0.1.x Breaking Changes from v0.1.0

v0.1.x releases may include bug fixes and security patches. No breaking
changes are planned for v0.1.x.

## v0.2.x (Next Development)

- **Status**: In development
- **Planned features**: SFTP, WebDAV, S3 support
- **Breaking changes**: Expected (new transport backends, API additions)

## Unsupported Versions

Versions older than the current stable release are not actively supported.
Users on older versions should upgrade to the latest stable release.

Security patches are only applied to:
1. The current stable version
2. The most recent previous version (only for critical security issues)

## How to Upgrade

### Desktop Application

Download the latest installer from the GitHub Releases page and run it.
The installer will update the existing installation.

### Docker

Pull the latest image:
```bash
docker pull eniacsystems/fileeniac:latest
```

### Standalone Binary

Download the new binary from GitHub Releases and replace the existing one.
Configuration files are generally compatible across versions.

## Version Information in Application

Run the application with `--version` to see the current version:
```bash
./fileeniac --version
```

## Security Updates

Critical security vulnerabilities are patched as follows:

| Severity | Current Stable | Previous Stable |
|----------|---------------|-----------------|
| Critical | Yes, within 30 days | Best effort |
| High | Yes, within 60 days | No |

Users are encouraged to stay on the latest version for security coverage.

## deprecation Policy

Before a version reaches end of life:
- Deprecation notice in release notes
- Notice in the repository README
- Minimum 30-day notice before end of support

## Contact

For questions about versioning or support, open an issue:
https://github.com/Joao-Aschenbrenner/FileENIAC/issues