# Architecture Audit — Sprint 1

## Data
2026-06-11

## Status
EM REVISAO

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
    ├── deploy/                  # Deploy Service
    │   ├── service.go           # Deploy, Rollback, Verify, Validate, GetHistory
    │   ├── ftp/                 # FTPS client (PoC migrada)
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
    ├── log/                     # Log wrapper (vazio, delegado ao logger)
    └── logger/                  # Logger baseado em Zap
        └── logger.go            # Info, Error, Level, configuração
```

### O Que Está Estável
- Workspace Registry (Init, Open, Status)
- Project Registry (Add, Remove, List, Get) + Server CRUD
- History Engine orientado a eventos (12 tipos, deploy/rollback/verify)
- Deploy Service (Deploy, Rollback, Verify, Validate, GetHistory)
- FTPS Client (Connect, Upload, Download, Delete, List)
- Packer (tar.gz com excludes)
- Token HMAC-SHA256
- Bypass com renomeação dinâmica
- SQLite migrator com rollback
- Logger Zap padronizado
- CLI completa (cobra commands)

### O Que Ainda É Experimental
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
- Pipeline completo: Pack → Connect FTPS → Upload → Manifest → Record
- Pack com excludes configuráveis
- Suporte a FTPS com TLS explícito (porta 21)
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
| DT-015 | Sem mock de FTPS para testes | Alta | Sprint 2 |
| DT-016 | Sem integração com git (commits, branches) | — | Sprint 2 |

---

## 9. Riscos

### Curto Prazo (Sprint 2)
| Risco | Probabilidade | Impacto | Mitigação |
|-------|:------------:|:-------:|-----------|
| FTPS upload falhar sem retry | Média | Alto | Implementar retry exponencial + circuit breaker |
| Rollback não reverter arquivos | Alta | Médio | Implementar restore de artefato anterior |
| Testes sem mock FTPS quebram em CI | Alta | Alto | Criar mock interface para FTPS |
| Workspace discovery pode corromper contexto | Baixa | Alto | Isolar activeContext ou torná-lo thread-safe |

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
| Linhas de código Go | 3.701 |
| Arquivos Go | 32 (23 prod + 9 test) |
| Pacotes | 13 |
| Testes unitários | 32 |
| Cobertura média | ~48% (ponderado) |
| Dependências externas | 11 diretas |
| Comandos CLI | 10 (workspace: 3 + project: 4 + deploy: 3) |
| Tabelas SQLite | 9 (3 ativas + 2 legado + 4 de sistema) |
| Migrations | 3 versionadas |

---

## 11. Conclusão

Sprint 1 produziu fundação sólida para o Workspace. Os débitos críticos estão identificados e endereçados no hardening (DT-007, DT-008, DT-009) e no planejamento da Sprint 2 (DT-006, DT-010, DT-012, DT-015).

**Próximo passo**: Aprovação deste audit + ADR-011 revisado + testes de integração + hardening do Deploy Engine + métricas para liberar Sprint 2.
