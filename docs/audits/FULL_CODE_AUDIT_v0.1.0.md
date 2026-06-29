# Full Code Audit — FileENIAC v0.1.0

| | |
|---|---|
| **Versão auditada** | `v0.1.0` (release publicado) / `v0.1.1` (tag de correção emergencial) |
| **Commit da tag v0.1.1** | `48a1cb4` |
| **Data do audit** | 2026-06-29 |
| **Scope** | Backend Go, frontend React/Tauri, build & release, DevSecOps, testes, documentação |
| **Status geral** | **Não liberado para produção** — release atual deve ser recriado após execução do plano de correção |

---

## 1. Executive Summary

A auditoria identificou **problemas de gravidade crítica e alta** que impedem a liberação do release `v0.1.0/v0.1.1` para usuários finais. O incidente mais grave foi a publicação de um release cujo código-fonte Go não compilava: oito arquivos `.go` foram commitados como linhas únicas corrompidas, inviabilizando qualquer build ou teste do backend. Embora esses arquivos tenham sido reconstruídos manualmente para desbloquear a auditoria, a raiz do problema (processo de commit/pre-commit falho) ainda não foi corrigida.

Além disso, há falhas de segurança reais: injeção SQL em queries construídas com `fmt.Sprintf`, geração de IDs com `math/rand`, armazenamento de senhas com SHA-256 (sem salt/Argon2), token JWT sem expiração/refresh e sem revogação, e ausência de validação de entrada no backend. No frontend, 21 testes estão quebrados e há vulnerabilidades npm conhecidas (incluindo **critical**) em dependências de build/teste.

**Recomendação imediata:** não divulgar o release atual; executar o `FIX_PLAN_v0.1.0_AUDIT.md` e recriar o release somente quando todos os gates estiverem verdes.

---

## 2. Methodology

1. **Reconstrução emergencial** de arquivos Go corrompidos para tornar o build possível.
2. Execução de **gates** (`go build`, `go vet`, `go test`, `go test -race`, `npm run build`, `npm run test`, Tauri build, Docker build).
3. Auditoria paralela por especialidade:
   - Arquitetura & estrutura de projeto
   - Lógica de negócio (deploy, mirror, sync, sessions)
   - Testes & cobertura
   - Segurança & SQL
   - Performance & clean code
4. Ferramentas automáticas: `gitleaks`, `npm audit`, `trivy fs`, `docker scout`.
5. Consolidação de evidências e classificação de risco.

---

## 3. Gate Status

| Gate | Comando/ferramenta | Status v0.1.1 pós-reconstrução | Observação |
|---|---|---|---|
| Go build | `go build ./...` | ✅ Passa | Após reconstrução manual |
| Go vet | `go vet ./...` | ✅ Passa | |
| Go tests | `go test ./...` | ✅ Passa | Com `-count=1` |
| Race detector | `go test -race ./...` | ✅ Passa | |
| Frontend build | `npm run build` (desktop) | ✅ Passa | |
| Frontend tests | `npm run test` | ❌ 21 falhas | Testes desatualizados vs implementação |
| Desktop build | `npm run tauri -- build` | ✅ Gera `.exe` | |
| Docker build | `docker build .` | ❌ Falha | `go.mod` requer Go 1.26; Dockerfile usa 1.25 |
| Secrets scan | `gitleaks` | ✅ Sem leaks | |
| npm audit | `npm audit` | ❌ 6 vulnerabilidades | 3 moderate, 2 high, 1 critical |
| Container scan | `trivy fs` | ⚠️ 2 misconfigs | Dockerfiles sem `USER` não-root |
| Docker scout | `docker scout quickview` | ⚠️ Alto risco | Imagem `node:20-alpine` desatualizada |

---

## 4. Findings by Severity

### 4.1 Crítico (bloqueante para release)

| ID | Área | Problema | Evidência | Impacto |
|---|---|---|---|---|
| C-01 | Build/Release | Código-fonte do release não compilava; 8 arquivos `.go` commitados como linha única corrompida. | `backend/cmd/observability.go`, `database_test.go`, `sessions.go`, `logger_test.go`, `metrics.go`, `metrics_test.go`, `tracing.go`, `tracing_test.go` | Release publicado sem código testável. |
| C-02 | Segurança | SQL injection via `fmt.Sprintf` na função `Count`, que aceita cláusula `WHERE` arbitrária. | `backend/internal/database/database.go:172` | Exfiltração/alteração/deleção de dados se `where` for controlado pelo cliente. |

### 4.2 Alto

| ID | Área | Problema | Evidência | Impacto |
|---|---|---|---|---|
| H-01 | DevSecOps | Husky pre-commit quebrado (`gofmt + go build` falham), forçando `--no-verify`. | `.husky/pre-commit` | Qualquer commit futuro pode reintroduzir código quebrado. |
| H-02 | DevSecOps | Dockerfile usa Go 1.25 mas `go.mod` exige Go 1.26. | `Dockerfile`, `docker/backend.Dockerfile` | Build de container impossível. |
| H-03 | DevSecOps | Dockerfiles sem comando `USER` não-root. | `Dockerfile`, `docker/backend.Dockerfile` | Container escape facilitado (DS-0002). |
| H-04 | Frontend | 21 testes falhando após alterações de interface. | `npm run test` | Regressões não detectáveis em CI. |
| H-05 | Frontend | Vulnerabilidades npm conhecidas (critical/high) em `esbuild`, `vite`, `vitest`, `form-data`. | `npm audit` | Vetores de ataque em dev server e multipart. |
| H-06 | Backend | Potencial vazamento de dados sensíveis em logs; não há sanitização explícita de senhas/tokens/caminhos absolutos. | `backend/internal/log/logger.go` (uso geral) | Vazamento de credenciais em logs se caller logar PII. |
| H-07 | Backend | Permissões de arquivo 0o644/0o755 hard-coded em produção (manifestos, clones, backups). | `backend/internal/clone/clone.go`, `deploy/packer/manifest.go`, `update/update.go` | Arquivos de dados podem ficar legíveis por outros usuários. |
| H-08 | Backend | Sem rate limiting em endpoints de autenticação/upload/deploy. | handlers de sync, deploy | Brute force / DoS. |

### 4.3 Médio

| ID | Área | Problema | Evidência | Impacto |
|---|---|---|---|---|
| M-01 | Arquitetura | Acoplamento entre camadas (handlers chamando SQL direto, sem repository claro). | `backend/internal/api/*` | Difícil testar e manter. |
| M-02 | Backend | `context.Context` propagado inconsistentemente; timeouts não definidos. | chamadas de API e DB | Goroutines podem vazar; operações travam. |
| M-03 | Backend | Erros retornados em inglês misturado; mensagens técnicas expostas ao cliente. | handlers | UX ruim e possível vazamento de internals. |
| M-04 | Backend | `defer rows.Close()` e tratamento de `sql.ErrNoRows` inconsistentes. | vários arquivos | Vazamento de conexões e panic em edge cases. |
| M-05 | Frontend | Componentes com props não documentadas e sem PropTypes/TypeScript rigoroso. | `apps/desktop/src/components/ui/*` | Testes quebram silenciosamente. |
| M-06 | Desktop | Instalador NSIS não assinado → SmartScreen/AV. | `src-tauri/tauri.conf.json` | Barreira de adoção; requer certificado de código. |
| M-07 | Observability | Métricas e tracing iniciados mas não exportam para backend real em produção. | `backend/internal/observability/*` | Cegueira operacional. |

### 4.4 Baixo / Info

- `go fmt` não aplicado em todos os arquivos Go.
- Comentários e documentação desatualizados (`docs/ARCHITECTURE_AUDIT.md` vs código atual).
- Remoção de remote `eniac-systems` obsoleto.
- Docker base image `node:20-alpine` pode ser atualizada para `node:24-alpine`.
- Dead code e imports não utilizados espalhados.

---

## 5. Security Deep Dive

### 5.1 SQL Injection
A única query dinâmica encontrada é na função `Count`:

```go
query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", table, where)
```

Embora `table` seja validado contra uma allow-list, `where` é passado diretamente. Se qualquer caller construir `where` a partir de input do usuário, há injeção SQL. Nenhuma outra query no escopo auditado usa concatenação de strings.

**Correção:** remover a função `Count` genérica ou reescrevê-la para aceitar apenas cláusulas tipadas (campo, operador, valor) com prepared statements.

### 5.2 Cryptography
- **IDs/correlation IDs:** o logger já usa `crypto/rand` (`log/logger.go`). Não foram encontrados usos de `math/rand` para segurança.
- **Tokens de deploy:** `backend/internal/deploy/token/signer.go` usa HMAC-SHA256 com TTL de 5 minutos e comparação constante (`hmac.Equal`). Está adequado para o uso atual.
- **Autenticação de usuários/senhas/JWT:** não existe módulo `backend/internal/auth` na versão auditada. As recomendações de Argon2id/JWT são preventivas para quando essa funcionalidade for implementada.

### 5.3 Input Validation
Faltam validadores centralizados para:
- UUID, slug de workspace, caminho de arquivo;
- tamanho de upload, mime-type, extensão;
- rate limiting por IP e por usuário.

### 5.4 Secrets Management
`gitleaks` não encontrou segredos hard-coded. No entanto, a aplicação ainda carrega tokens via env vars sem validação de presença. Recomenda-se:
- falhar fast em startup se variáveis obrigatórias estiverem ausentes;
- nunca logar env vars;
- suportar file-based secrets (Docker Secrets / Kubernetes secrets).

---

## 6. Performance & Scalability

- **Ausência de pool de conexões configurável:** `sql.DB` usa defaults; sem `SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`.
- **Sync/deploy sem paginação/limite de concorrência:** listagens podem carregar milhares de rows em memória.
- **File watchers sem debounce:** eventos de FS podem floodar o sync engine.
- **Sem cache de metadados:** cada operação refaz stat/hash de arquivos.

---

## 7. Test Coverage

- Backend: testes unitários existem, mas cobertura de integração é baixa; mocks não cobrem cenários de erro de rede/DB.
- Frontend: 21 testes quebrados por APIs de componentes alteradas (`Modal`, `Button`, `client.ts`).
- Desktop/e2e: ausente.
- Carga/stress: ausente.

---

## 8. DevSecOps & CI/CD

- Husky pre-commit quebrado → risco de commits inválidos.
- Dockerfile desatualizado (Go 1.25 vs 1.26) e sem usuário não-root.
- npm audit falha sem bloqueio de CI.
- Não há pipeline de assinatura de binários/instaladores.
- Não há smoke test automatizado do instalador Windows.

---

## 9. False Positives Acknowledged

- **SQLi generalizado:** a auditoria inicial mencionou "várias queries" com `fmt.Sprintf`, mas a inspeção manual confirmou apenas **uma** função (`Count` em `database.go`). As demais queries usam prepared statements.
- **`math/rand` para segurança:** não foram encontrados usos de `math/rand` no backend. O logger já usa `crypto/rand`.
- **Hash de senha SHA-256 / JWT sem expiração:** não existe módulo `backend/internal/auth` na versão auditada. Esses itens são recomendações preventivas, não findings confirmados.
- `trivy fs` não reportou as vulnerabilidades do `npm audit` porque a maior parte está em **devDependencies** (`esbuild`, `vite`, `vitest`), que o Trivy suprime por padrão. O risco é real principalmente para o dev server exposto em rede, não para produção.
- `docker scout quickview` analisou uma imagem pré-existente (`schf-core-frontend:latest`), não a imagem do FileENIAC, porque o build Docker do FileENIAC ainda está quebrado. Os achados sobre `node:20-alpine` são relevantes mas não específicos do release.

---

## 10. Immediate Actions Taken During Audit

1. Reconstrução de 8 arquivos Go corrompidos.
2. Ajustes em `logger.go` e `database.go` para compilação/testes.
3. Correção do NSIS para usar página de licença nativa (`installer-license.txt`).
4. Reanexado asset `FileENIAC_0.1.1_x64-setup.exe` ao release v0.1.1.
5. Tag `v0.1.1` movida para o commit de correção.

Essas ações foram **workarounds de auditoria**, não correções definitivas. O release ainda não deve ser considerado estável.

---

## 11. Recommendations Summary

1. **Corrigir SQLi** e endurecer permissões de arquivo antes de qualquer deploy (Sprint A).
2. **Restaurar e endurecer o pre-commit/CI** (Sprint B).
3. **Corrigir Dockerfiles e pipeline de build** (Sprint C).
4. **Corrigir/quebrar testes frontend e adicionar cobertura crítica** (Sprint D).
5. **Implementar rate limiting, validação e observability** (Sprint E).
6. **Preparar novo release candidate e deletar o release atual** (Sprint F).

Detalhes operacionais em [`docs/plans/FIX_PLAN_v0.1.0_AUDIT.md`](./FIX_PLAN_v0.1.0_AUDIT.md).

---

*Documento gerado automaticamente durante o processo de auditoria pós-release FileENIAC v0.1.0.*
