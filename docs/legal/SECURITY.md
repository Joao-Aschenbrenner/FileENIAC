# SPDX-License-Identifier: MIT
# FileENIAC Security Policy

**Version**: 1.0
**Last updated**: 2026-06-28

## Security Model

FileENIAC is a local-first desktop application. The primary security boundary
is the user's device. The developer cannot access any data stored by FileENIAC.

## Environment Variables

### Required

- `FILEENIAC_VAULT_PASSWORD`: Encryption password for the credential Vault.
  **Must be set for credential encryption to be active.** If not set, the
  Vault stores credentials with no encryption.

### Optional

- `FILEENIAC_DATA_DIR`: Data directory (default: `./data`)
- `FILEENIAC_API_TOKEN`: Bearer token for API authentication
- `FILEENIAC_API_PORT`: Desktop API port (default: `8080`)

**Never commit these values to version control.** Use environment variables
or a secrets manager in production.

## Credential Storage

- GitHub OAuth tokens and FTPS passwords are stored in `vault.db`
- The Vault uses AES-256-GCM encryption with a key derived from
  `FILEENIAC_VAULT_PASSWORD` via Argon2id
- If `FILEENIAC_VAULT_PASSWORD` is not set, credentials are stored in
  **plain text** — always set this variable

## Network Security

- All HTTP endpoints use configurable timeouts (default: 30s read, 10m write)
- Body size limit of 1MB on all HTTP requests prevents memory exhaustion
- FTP connections use TLS 1.2 with ECDHE-only cipher suites
- URL credentials are stripped from all log output
- Bearer tokens are ephemeral (32 bytes, session-scoped)

## Input Validation

- All user-provided names validated against an allowlist regex
  (`^[a-zA-Z0-9_\-\.]+$`) before storage or use
- File paths are canonicalized to prevent path traversal attacks
- Symbolic links are skipped during workspace walks (unless Developer Mode
  is enabled on Windows)
- SQL table names use a hardcoded allowlist in all queries
- No SQL injection vectors exist in the codebase

## Source Code Security Practices

- No hardcoded credentials in source code (enforced by security tests)
- Secrets loaded from environment variables only
- Log output sanitized to prevent credential leakage
- Error messages are sanitized before returning to clients

## Dependency Security

- Go dependencies: reviewed in `go.mod` and `go.sum`
- Rust dependencies: reviewed in `Cargo.toml` and `Cargo.lock`
- Node.js dependencies: reviewed in `package.json` and `pnpm-lock.yaml`
- All direct dependencies use permissive licenses (MIT, Apache-2.0, BSD)

## Vulnerability Reporting

See `VULNERABILITY_DISCLOSURE.md` for the disclosure process.

## Known Limitations

1. **Vault fallback**: If `FILEENIAC_VAULT_PASSWORD` is not set, credentials
   are stored unencrypted in `vault.db`.

2. **No code signing**: The NSIS installer is not code-signed, so Windows
   SmartScreen may display a warning.

3. **Local-only security**: The primary security boundary is the device.
   If malware has access to the device, it can read FileENIAC's data
   directory.

4. **No multi-user isolation**: FileENIAC is a single-user desktop
   application. All users on the same device share the same data.

## Security Updates

Security fixes are released as patch versions. See `SUPPORTED_VERSIONS.md`
for the current support window.

## Penetration Testing

Independent security researchers are welcome to audit the codebase.
See `VULNERABILITY_DISCLOSURE.md` for coordinated disclosure guidelines.

## Best Practices for Users

1. **Set `FILEENIAC_VAULT_PASSWORD`**: Always enable Vault encryption
2. **Protect your data directory**: Ensure `$FILEENIAC_DATA_DIR` is not
   accessible to unauthorized users
3. **Use strong FTPS passwords**: Unique per server
4. **Rotate credentials periodically**: Especially if you suspect compromise
5. **Keep the application updated**: Use the latest stable release
6. **Secure your device**: Use full-disk encryption, screen lock, and
   antivirus software