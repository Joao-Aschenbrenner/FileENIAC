# FileENIAC

Plataforma única para gerenciamento de workspace local, repositórios Git, GitHub, deploys FTPS, histórico, auditoria e monitoramento.

## Arquitetura

```
FileENIAC/
├── apps/desktop/       # Tauri v2 + React + TypeScript
├── backend/            # Go backend
│   ├── cmd/            # CLI commands
│   ├── internal/       # Core engines
│   │   ├── workspace/  # Workspace management
│   │   ├── registry/   # Project registry
│   │   ├── deploy/     # Deploy orchestration
│   │   ├── sync/       # Sync engine + mirror
│   │   ├── git/        # Git operations
│   │   ├── github/     # GitHub API + personal access token
│   │   ├── history/    # History + audit
│   │   ├── health/     # Health checks
│   │   └── agents/     # Agent API (future)
│   ├── config/         # Configuration
│   └── database/       # SQLite schema + migrations
├── docs/               # Architecture Decision Records
├── scripts/            # Build and dev scripts
└── docker/             # Container support
```

## Roadmap

| Sprint | Foco | Entregas |
|--------|------|----------|
| 0 | Arquitetura | ADRs, estrutura, schema |
| 1 | Core | FTPS, Deploy, History |
| 2 | Workspace | Registry, Mirror, Diff |
| 3 | GitHub | OAuth, Discovery, Bootstrap |
| 4 | Desktop | Dashboard, UI completa |
| 5 | Agent | API, Observabilidade, IA |

## Legal

- [MIT License](LICENSE)
- [Terms of Use](docs/legal/TERMS_OF_USE.md)
- [Privacy Policy](docs/legal/PRIVACY_POLICY.md)
- [Third-Party Services](docs/legal/THIRD_PARTY_SERVICES.md)
- [LGPD Compliance](docs/legal/LGPD.md)
- [Data Processing](docs/legal/DATA_PROCESSING.md)
- [Security Policy](docs/legal/SECURITY.md)
- [Security & Credential Management](docs/legal/SECURITY_AND_CREDENTIALS.md)
- [Installer Notice](docs/legal/INSTALLER_NOTICE.md)
- [Vulnerability Disclosure](docs/legal/VULNERABILITY_DISCLOSURE.md)
