# ADR-011: Workspace Registry — Estrutura, Ciclo de Vida e Relacionamentos

## Status
APROVADO

## Data
2026-06-11 (Revisão pós-Sprint 1)

## Contexto
A Sprint 1 implementou o Workspace Registry, Project Registry, Deploy Service e History Engine. Esta ADR documenta oficialmente as definições, relacionamentos e decisões de design para evitar interpretações divergentes.

## Decisão

### Definições

#### Workspace
Diretório contendo um arquivo `.eniac/config.toml` e um banco SQLite `.eniac/workspace.db`. É a unidade organizacional raiz. Todo workspace tem um nome único descritivo e um caminho no sistema de arquivos.

```
{workspace_path}/
├── .eniac/
│   ├── config.toml      # Nome, descrição
│   └── workspace.db     # Projetos, servidores, deploys, eventos
└── src/                 # Código fonte dos projetos (fora do .eniac/)
```

#### Projeto
Entidade registrada no workspace que representa um repositório local. TODO projeto pertence a um workspace. Um projeto pode ter zero ou mais servidores de deploy.

**Atributos**: name, display_name, local_path, remote_path, branch, git_url, environment, is_active

#### Servidor
Configuração de conexão FTPS associada a um projeto. Um projeto pode ter múltiplos servidores (ex: produção, staging), mas apenas um ativo por vez.

**Atributos**: name, type (ftps), host, port, user, password (secret), target_path, verify_url, is_active

#### Ambiente
Tag associada ao projeto (production, staging, development). Não é uma entidade separada — é um campo em `projects.environment`.

#### Deploy
Operação que empacota o código local em .tar.gz, faz upload FTPS para o servidor e registra o resultado no history engine.

**Ciclo de vida**: STARTED → (SUCCESS | FAILED)

### Relacionamentos

```
Workspace (1)
    │
    ├── Projetos (N) — projects table
    │       │
    │       ├── Servidores (N) — servers table
    │       │       └── FK: servers.project_id → projects.id
    │       │
    │       ├── Deploy Logs (N) — deploy_logs table
    │       │       └── FK: deploy_logs.project_id → projects.id
    │       │
    │       └── Rollback Logs (N) — rollback_logs table
    │               └── FK: rollback_logs.project_id → projects.id
    │
    └── Eventos (N) — events table (global, não vinculado a projeto)
```

### Regras de Negócio

1. **Um workspace contém N projetos** — não há limite máximo teórico
2. **Um projeto pertence a exatamente um workspace** — não há compartilhamento entre workspaces
3. **Um projeto pode ter N servidores** — mas apenas 1 ativo por operação
4. **Um deploy pertence a exatamente um projeto** — via FK project_id
5. **Um rollback referencia um deploy específico** — via deploy_id
6. **Workspace é auto-contido** — pode ser copiado/transportado entre máquinas
7. **Workspace opera 100% offline** — deploy é a única operação que requer rede
8. **Fonte da verdade é Git** — workspace nunca sobrescreve repositório local

### Armazenamento

| Entidade | Storage | Localização |
|----------|---------|-------------|
| Workspace config | TOML | `.eniac/config.toml` |
| Projetos | SQLite | `.eniac/workspace.db` → `projects` |
| Servidores | SQLite | `.eniac/workspace.db` → `servers` |
| Deploy logs | SQLite | `.eniac/workspace.db` → `deploy_logs` |
| Rollback logs | SQLite | `.eniac/workspace.db` → `rollback_logs` |
| Eventos | SQLite | `.eniac/workspace.db` → `events` |

### CLI

```
eniac workspace init <name> [path]     — Cria workspace
eniac workspace open [path]            — Abre workspace
eniac workspace status                 — Exibe estado
eniac project add <name> <path>        — Registra projeto
eniac project remove <id>              — Remove projeto
eniac project list                     — Lista projetos
eniac project show <name>              — Detalhes do projeto
eniac deploy run <project>             — Executa deploy
eniac deploy verify <project>          — Verifica último deploy
eniac deploy rollback <project>        — Reverte último deploy
eniac deploy history <project>         — Histórico de deploys
```

### Modelos de Dados

```go
// Workspace (config.toml)
type Workspace struct {
    Name        string
    Description string
    Path        string
}

// Project (SQLite)
type Project struct {
    ID          int64
    Name        string    // UNIQUE
    DisplayName string
    LocalPath   string
    RemotePath  string
    Branch      string
    GitURL      string
    Environment string    // production, staging, development
    ServerID    int64
    IsActive    bool
    CreatedAt   string
    UpdatedAt   string
}

// Server (SQLite)
type Server struct {
    ID         int64
    ProjectID  int64     // FK → projects.id
    Name       string
    Type       string    // "ftps"
    Host       string
    Port       int
    User       string
    Password   string    // secret, não serializado em JSON
    TargetPath string
    VerifyURL  string
    IsActive   bool
}
```

## Consequências

- Definições claras evitam ambiguidades entre membros da equipe
- Relacionamentos documented permitem extensões futuras (ex: servidor N:N deploy)
- CLI hierarchy reflete a árvore de entidades
- SQLite como storage único simplifica backup e transporte
- Decisão de non-shared workspace simplifica a arquitetura mas impede colaboração multi-usuário no mesmo workspace (adequado para cenário atual)
