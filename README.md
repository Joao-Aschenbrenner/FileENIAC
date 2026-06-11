# ENIAC Workspace

Plataforma única para gerenciamento de workspace local, repositórios Git, GitHub, deploys FTPS, histórico, auditoria e monitoramento.

## Arquitetura

```
eniac-workspace/
├── apps/desktop/       # Tauri v2 + React + TypeScript
├── backend/            # Go backend
│   ├── cmd/            # CLI commands
│   ├── internal/       # Core engines
│   │   ├── workspace/  # Workspace management
│   │   ├── registry/   # Project registry
│   │   ├── deploy/     # Deploy orchestration
│   │   ├── sync/       # Sync engine + mirror
│   │   ├── git/        # Git operations
│   │   ├── github/     # GitHub OAuth + API
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
