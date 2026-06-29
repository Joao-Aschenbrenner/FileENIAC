# Security Policy
# SPDX-License-Identifier: MIT

## Reporting a Vulnerability

If you discover a security vulnerability in FileENIAC, please report it
responsibly using one of these methods:

**Preferred**: GitHub Security Advisories
- Go to https://github.com/Joao-Aschenbrenner/FileENIAC/security/advisories/new
- Select "Report a vulnerability"
- Provide details of the vulnerability, affected versions, and potential impact

**Alternative**: Open a regular issue labeled `security`
- Note in the issue that it should not be made public until fixed

## Response Timeline

| Severity | Initial Response | Fix Target |
|----------|-----------------|------------|
| Critical | 48 hours | 30 days |
| High | 48 hours | 45 days |
| Medium | 7 days | 60 days |
| Low | Best effort | Next release |

## Supported Versions

| Version | Supported |
|---------|-----------|
| Latest stable (v0.1.x) | Yes |
| Older versions | No |

Upgrade to the latest stable version for security coverage.

For full vulnerability disclosure policy, see `docs/legal/VULNERABILITY_DISCLOSURE.md`.

## Security Best Practices for Users

1. Protect access to your device and data directory
2. Keep the application updated
3. Use strong, unique FTPS passwords
4. Rotate credentials if you suspect compromise
5. Review logs before sharing them publicly

## Out of Scope

- Social engineering attacks
- Security of third-party services (GitHub, FTPS servers)
- Physical security of your device
- Vulnerabilities in dependencies that don't affect FileENIAC's security posture