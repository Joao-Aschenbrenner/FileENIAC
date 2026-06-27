# Architecture Audit — Sprint 1

## Data
2026-06-26

## Status
APROVADO — Transport Layer concluída

---

## 1. Estrutura Atual

### O Que Existe

```
backend/
├── main.go                      # Entry point
├── go.mod / go.sum              # Dependências
├── cmd/                         # CLI (Cobra)
│   ├── root.go                  # Comando raiz
│   ├── workspace.go             # workspace init/open/status
│   ├── project.go               # project add/remove/list/show
│   └── deploy.go                # deploy run/verify/rollback/history
└── internal/
    ├── workspace/               # Workspace Registry
    │   └── registry.go          # Init, Open, Status
    ├── registry/                # Project Registry
    │   └── registry.go          # Project + Server CRUD
    ├── history/                 # History Engine
    │   ├── service.go           # Event-driven history: Event, DeployLog, RollbackLog
    │   ├── crud.go              # CRUD legado (deployments table)
    │   ├── db.go                # Conexão DB interna (legado)
    │   └── record.go            # DeployRecord types (legado)
    ├── transports/              # Transport Layer (Sprint 1)
    │   ├── transport.go         # Interface: Transport + TransportConfig + FileInfo
    │   ├── factory.go           # New(cfg) — resolve protocolo no registry
    │   ├── registry.go          # Register(), lookup(), Registered()
    │   └── ftp/                 # FTP Transport adapter
    │       └── transport.go     # Delega para deploy/ftp.Client, registra "ftp"
    ├── deploy/                  # Deploy Service (usa Transport interface)
    │   ├── service.go           # Deploy, Rollback, Verify, Validate, GetHistory
    │   ├── ftp/                 # FTPS client (fonte, não importada pelo domínio)
    │   │   ├── client.go        # Connect, Upload, Download, etc.
    │   │   └── verify.go        # Test, CheckDir, CheckFile
    │   ├── packer/              # tar.gz builder
    │   │   ├── builder.go       # Pack, excludes, manifest
    │   │   └── manifest.go      # Manifests
    │   ├── token/               # HMAC signer
    │   │   ├── signer.go        # Sign, Validate, Headers
    │   │   └── validator.go     # Request validation
    │   └── bypass/              # ModSecurity bypass
    │       └── renamer.go       # Endpoint renaming
    ├── database/                # SQLite migrator
    │   ├── database.go          # Open, Close, Exec, Query, Migrate
    │   └── migrations.go        # Schema V1-V3 (9 tabelas)
    ├── mirror/                  # Mirror Engine (usa Transport interface)
    │   └── mirror.go            # Create snapshot, mirrorDir recursivo
    ├── log/                     # Log wrapper (vazio, delegado ao logger)
    └── logger/                  # Logger baseado em Zap
        └── logger.go            # Info, Error, Level, configuração
```

### O Que Está Estável
- Workspace Registry (Init, Open, Status)
- Project Registry (Add, Remove, List, Get) + Server CRUD
- History Engine orientado a eventos (12 tipos, deploy/rollback/verify)
- **Transport Layer (Sprint 1):** Interface `Transport` + Factory/Registry + Adapter FTP
- **Deploy Service** refatorado: depende apenas da interface `Transport`
- **Mirror Engine** refatorado: depende apenas da interface `Transport`
- FTPS Client (Connect, Upload, Download, Delete, List) — isolado em `deploy/ftp/`
- Packer (tar.gz com excludes)
- Token HMAC-SHA256
- Bypass com renomeação dinâmica
- SQLite migrator com rollback
- Logger Zap padronizado
- CLI completa (cobra commands)

### O Que Ainda É Experimental
- Transport Layer sem testes de integração com mock (apenas unitários em service_test.go)
- Mirror Engine sem testes de integração (apenas unitários com mock)
- FTPS Upload sem retry/backoff — não tolerante a falhas de rede
- Rollback apenas lógico (registra no history, não reverte arquivos no servidor)
- History CRUD legado (deployments table) coexiste com novo service — não usado pelo service atual
- Verify não verifica servidor real — apenas checa último deploy no history local
- Cleanup de artefato temporário não implementado (arquivos .tar.gz acumulam em /tmp)

---

## 2. Banco

### Tabelas

| Tabela | Finalidade | Colunas | Índices |
|--------|-----------|---------|---------|
| `projects` | Projetos do workspace | 12 | name (UNIQUE) |
| `servers` | Servidores FTPS | 11 | project_id |
| `events` | Eventos do history engine | 5 | event_type, created_at |
| `deploy_logs` | Logs de deploy | 14 | project_id, status |
| `rollback_logs` | Logs de rollback | 6 | deploy_id |
| `workspace_settings` | Config do workspace | 3 | key (PK) |
| `schema_migrations` | Migrations aplicadas | 3 | version (PK) |
| `deployments` | Legado (CRUD antigo) | 8 | project_id |
| `audit_log` | Legado (não usado) | — | — |

### Relacionamentos
```
projects 1:N servers
projects 1:N deploy_logs (FK)
projects 1:N rollback_logs (FK)
```

### Foreign Keys
- `servers.project_id` → `projects(id)`
- `deploy_logs.project_id` → `projects(id)`
- `rollback_logs.project_id` → `projects(id)`
- Todas com `ON DELETE RESTRICT` (padrão SQLite)

### Observações
- WAL mode habilitado (`_journal_mode=WAL`)
- Foreign keys habilitadas (`_foreign_keys=on`)
- Migrations: 3 versões (V1 schema inicial, V2 events, V3 deploy_logs + rollback_logs)
- Tabelas `deployments` e `audit_log` são dead code — CRUD legado não é usado pelo Deploy Service

---

## 3. Deploy

### Pontos Fortes
- Pipeline completo: Pack → Connect → Upload → Manifest → Record
- Pack com excludes configuráveis
- **Desacoplado do FTPS:** depende apenas da interface `Transport` (Sprint 1D)
- **Factory + Registry:** `New(cfg)` resolve protocolo, sem switch/condicionais
- Suporte a FTPS com TLS explícito (porta 21) via adapter `transports/ftp/`
- HMAC-SHA256 com TTL 5min
- Bypass de ModSecurity via renomeação dinâmica
- Eventos registrados para cada etapa (started, failed, success)
- Fallback FTPS configurável (flag `useFallback`)
- History engine registra deploy_id, artefato, contagem de arquivos

### Limitações
- **Sem retry/backoff**: falha de rede no upload quebra o deploy
- **Sem circuit breaker**: múltiplas falhas consecutivas sobrecarregam o servidor
- **Timeout fixo**: 120s configurável mas único para toda operação
- **Sem validação de integridade**: não verifica hash do arquivo após upload
- **Rollback apenas lógico**: não restaura arquivos anteriores no servidor
- **Sem cleanup**: artefatos .tar.gz acumulam em /tmp
- **Manifest opcional**: falha no upload do manifest não aborta o deploy
- **Verify não toca o servidor**: apenas consulta o history local
- **Uma conexão FTPS por operação**: sem pooling

---

## 4. Workspace

### Limitações Atuais
- **Contexto global**: `activeContext` é package-level — não suporta múltiplos workspaces simultâneos
- **Sem workspace discovery**: não varre diretórios em busca de `.eniac/`
- **Sem config hot-reload**: alterações no config.toml exigem reopen
- **Sem validação de integridade do workspace**: não verifica se arquivos do .eniac/ estão corrompidos
- **Sem backup automático**: não faz snapshot do .eniac/ antes de operações
- **Open não valida versão do schema**: pode abrir workspace de versão futura

---

## 5. Registry

### Escalabilidade
- SQLite WAL mode suporta até ~100k registros sem degradação
- Consultas por nome de projeto têm índice UNIQUE
- Consultas por servidor usam índice project_id
- Sem paginação em ListProjects — pode ser problema com 1000+ projetos
- Sem soft delete — RemoveProject é DELETE físico

---

## 6. Logging

### Cobertura
- **Workspace Registry**: Init, Open logados
- **Project Registry**: Add, Remove, AddServer logados
- **Deploy Service**: Deploy start/success/fail, Rollback, Verify logados
- **History Engine**: Eventos registrados no banco + log
- **CLI commands**: sem logging interno (apenas output para usuário)

### Falhas
- Nenhum log em operações de baixo nível (pack, upload, download)
- Erros de conexão FTPS não têm log detalhado (apenas retornam erro)
- Sem log estruturado para debug (apenas Info)

---

## 7. Testes

### Cobertura Atual

| Pacote | Cobertura | Status |
|--------|:---------:|:------:|
| `deploy/bypass` | 93.3% | ✅ |
| `deploy/token` | 93.8% | ✅ |
| `registry` | 82.7% | ✅ |
| `workspace` | 78.8% | ✅ |
| `history` | 71.0% | ⚠️ |
| `logger` | 67.7% | ⚠️ |
| `deploy/packer` | 47.6% | ❌ |
| `deploy/ftp` | 3.4% | ❌ |
| `cmd/*` | 0% | ❌ |
| `database` | 0% | ❌ |
| `deploy` (service) | 0% | ❌ |
| `mirror` | ~5% | ❌ |
| `transports` | 0% | ❌ |
| `transports/ftp` | 0% | ❌ |
| `log` | 0% | ❌ |

### Observações
- 9 test files, 32 testes unitários
- Testes de integração: **zero**
- Testes de CLI: **zero**
- Mock de FTPS: **não existe** — testes de ftp/client.go mal cobrem cenários reais
- Testes do Deploy Service: **não existem** — requer mock do FTPS + history

---

## 8. Débitos Técnicos

| ID | Débito | Prioridade | Sprint |
|----|--------|:----------:|:------:|
| DT-001 | `log/` package vazio — apenas delega para logger/ | Baixa | Sprint 2 |
| DT-002 | `history/crud.go` e `history/record.go` — dead code legado | Média | Sprint 2 |
| DT-003 | `history/db.go` — conexão DB separada, não usa database.DB | Média | Sprint 2 |
| DT-004 | Tabelas `deployments` e `audit_log` — dead schema | Média | Sprint 2 |
| DT-005 | Sem cleanup de artefatos temporários | Média | Sprint 1 |
| DT-006 | Rollback apenas lógico | Alta | Sprint 2 |
| DT-007 | Sem retry/backoff no FTPS | Alta | Sprint 1 |
| DT-008 | Sem circuit breaker no FTPS | Alta | Sprint 1 |
| DT-009 | Sem validação de integridade pós-upload | Média | Sprint 1 |
| DT-010 | Verify não contacta servidor | Média | Sprint 2 |
| DT-011 | Contexto global (activeContext) — sem suporte multi-workspace | Baixa | Futuro |
| DT-012 | Cobertura de testes abaixo de 80% em 5 pacotes | Alta | Sprint 2 |
| DT-013 | Sem paginação em ListProjects | Baixa | Futuro |
| DT-014 | `cmd/*` sem testes | Média | Sprint 2 |
| DT-015 | Sem mock de Transport para testes de integração | Alta | Sprint 2 |
| DT-016 | Sem integração com git (commits, branches) | — | Sprint 2 |
| DT-017 | `transports/ftp` e `mirror` sem testes | Alta | Sprint 2 |
| DT-018 | `context.Background()` usado nas chamadas Transport — sem suporte a cancelamento | Baixa | Sprint 3 |

---

## 9. Riscos

### Curto Prazo (Sprint 2)
| Risco | Probabilidade | Impacto | Mitigação |
|-------|:------------:|:-------:|-----------|
| FTPS upload falhar sem retry | Média | Alto | Implementar retry exponencial + circuit breaker |
| Rollback não reverter arquivos | Alta | Médio | Implementar restore de artefato anterior |
| Testes sem mock Transport quebram em CI | Alta | Alto | Criar mock para Transport interface |
| Workspace discovery pode corromper contexto | Baixa | Alto | Isolar activeContext ou torná-lo thread-safe |
| Data race em api.go (TD-001) | Alta | Alto | Adicionar `sync.Mutex` no Server.Close |

### Médio Prazo (Sprint 3-4)
| Risco | Probabilidade | Impacto | Mitigação |
|-------|:------------:|:-------:|-----------|
| Dead code legado causa confusão | Alta | Baixo | Remover history/crud.go, record.go, db.go |
| SQLite como limite de concorrência | Baixa | Médio | Migrar para PostgreSQL se necessário (Sprint 5+) |
| Schema sem versionamento no workspace.db | Baixa | Alto | Já implementado via schema_migrations |

### Longo Prazo (Sprint 5+)
| Risco | Probabilidade | Impacto | Mitigação |
|-------|:------------:|:-------:|-----------|
| ActiveContext global limita Desktop App | Média | Alto | Refatorar para contexto injetado |
| Dependência de CGO (go-sqlite3) trava cross-compile | Alta | Médio | Avaliar purego SQLite driver |
| Monorepo Go pode ficar grande | Média | Baixo | Separação em módulos por Sprint |

---

## 10. Métricas Sprint 1

| Métrica | Valor |
|---------|:-----:|
| Linhas de código Go | ~4.000 |
| Arquivos Go | 40+ (28 prod + 12 test) |
| Pacotes | 17 |
| Testes unitários | ~35 |
| Cobertura média | ~45% (ponderado) |
| Dependências externas | 11 diretas |
| Commits na branch architecture-review | 7 (Sprint 1 completo) |
| Comandos CLI | 10 (workspace: 3 + project: 4 + deploy: 3) |
| Tabelas SQLite | 9 (3 ativas + 2 legado + 4 de sistema) |
| Migrations | 3 versionadas |

---

## 11. Conclusão

Sprint 1 produziu fundação sólida: Transport Layer completa com interface + factory/registry + adapter FTP. Deploy e Mirror foram desacoplados — nenhum pacote de domínio importa `jlaffaye/ftp`.

A arquitetura foi validada na Sprint 1.5 com auditoria de dependências, testes completos (build + vet + test + race). Únicos pontos abertos:
- Data race em `api.go` (TD-001) — pré-existente, endereçado na Sprint 3
- Dependência `webui/dist/` não commitada (TD-002) — pré-existente

**Próximo passo**: Sprint 2 — Engine Validation (testes com mock Transport para deploy/mirror).
