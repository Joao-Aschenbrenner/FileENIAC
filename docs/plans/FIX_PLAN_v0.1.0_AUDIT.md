# Fix Plan — FileENIAC v0.1.0 Audit

| | |
|---|---|
| **Audit base** | [`docs/audits/FULL_CODE_AUDIT_v0.1.0.md`](../audits/FULL_CODE_AUDIT_v0.1.0.md) |
| **Versão alvo** | `v0.1.2` (novo release) ou `v0.1.1` recriado — decidir após Sprint F |
| **Prioridade** | Crítico/Alto primeiro; Médio/Baixo em paralelo quando não bloqueante |
| **Regra geral** | 1 problema = 1 commit pequeno, testável e reversível |

---

## 1. Definition of Done (release gate)

Antes de qualquer novo release, **todos** os itens abaixo devem estar verdes:

- [ ] `go build ./...` passa
- [ ] `go vet ./...` passa
- [ ] `go test ./...` passa
- [ ] `go test -race ./...` passa
- [ ] `go test -count=1 ./...` passa
- [ ] `npm run build` passa (raiz e `apps/desktop`)
- [ ] `npm run test` passa
- [ ] `npm audit --audit-level=moderate` retorna 0 vulnerabilidades **ou** aceite documentado
- [ ] `docker build .` passa
- [ ] `docker build -f docker/backend.Dockerfile .` passa
- [ ] `trivy fs --scanners vuln,secret,misconfig --severity HIGH,CRITICAL .` sem findings
- [ ] `gitleaks detect` sem leaks
- [ ] Instalador Tauri gera `.exe` e smoke test manual passa
- [ ] `RELEASE_CHECKLIST.md` e `CHANGELOG.md` atualizados
- [ ] CHECKSUMS publicados no release

---

## 2. Sprint A — Critical Security Fixes

**Objetivo:** eliminar vetores de ataque que comprometem dados ou acesso.

| # | Issue | Arquivos principais | Critério de aceitação | Estimativa |
|---|---|---|---|---|
| A-1 | Migrar queries SQL para prepared statements | `backend/internal/database/*.go` | Zero `fmt.Sprintf` em query SQL; todos os parâmetros via `$n` ou `:named`; testes de injeção passam | 2-3 dias |
| A-2 | Substituir `math/rand` por `crypto/rand` | IDs, tokens, salts | IDs gerados com `crypto/rand`; testes de unicidade e entropia | 1 dia |
| A-3 | Implementar Argon2id para senhas | `backend/internal/auth/password.go` | Hash contém salt+params; verificação compatível; benchmark de tempo ~300ms | 1 dia |
| A-4 | Adicionar expiração, refresh e revogação a JWT | `backend/internal/auth/jwt.go`, session store | Access token 15min; refresh token 7d rotativo; endpoint `/logout` revoga; blacklist em memória/DB | 2 dias |
| A-5 | Sanitizar logs (nunca logar senhas/tokens/caminhos absolutos) | `backend/internal/log/logger.go` | Testes garantem que dados sensíveis não aparecem em logs | 1 dia |
| A-6 | Validar e restringir permissões de arquivos | `backend/internal/storage/*.go` | Arquivos de dados 0o600; diretórios 0o700; testes de permissão | 1 dia |

**Gate A:** `go test ./...`, `go test -race ./...`, testes de segurança novos passam.

---

## 3. Sprint B — Pre-commit, CI & Code Quality

**Objetivo:** garantir que código quebrado/inseguro não entre na main.

| # | Issue | Arquivos principais | Critério de aceitação | Estimativa |
|---|---|---|---|---|
| B-1 | Corrigir husky: `gofmt`, `go build ./...`, `go test ./...` | `.husky/pre-commit` | Commit normal passa sem `--no-verify` | 0.5 dia |
| B-2 | Adicionar lint `golangci-lint` (sqlclosecheck, gosec, errcheck) | `.golangci.yml` | Pipeline local/CI roda lint sem erros | 1 dia |
| B-3 | Aplicar `gofmt` em todos os arquivos Go | todos `.go` | Diff de formatação zerado | 0.5 dia |
| B-4 | Remover imports não utilizados e dead code | todos `.go` | `go vet` e `golangci-lint` limpos | 1 dia |
| B-5 | Configurar GitHub Actions para build/test/lint em PR | `.github/workflows/ci.yml` | CI bloqueia merge em falha | 1-2 dias |

**Gate B:** qualquer commit em `main` passa por husky e CI verdes.

---

## 4. Sprint C — Docker & Build Pipeline

**Objetivo:** containerização funcional, segura e reprodutível.

| # | Issue | Arquivos principais | Critério de aceitação | Estimativa |
|---|---|---|---|---|
| C-1 | Alinhar Dockerfile com `go.mod` (Go 1.26) | `Dockerfile`, `docker/backend.Dockerfile` | Build passa; imagem usa `golang:1.26-alpine` ou equivalente | 0.5 dia |
| C-2 | Adicionar usuário não-root nos containers | `Dockerfile`, `docker/backend.Dockerfile` | `USER app` antes de `CMD`; trivy DS-0002 resolvido | 0.5 dia |
| C-3 | Multi-stage build otimizado | Dockerfiles | Imagem final sem toolchain Go; tamanho reduzido | 1 dia |
| C-4 | Atualizar base image Node para LTS atual | `Dockerfile` | `node:24-alpine` ou LTS recomendada; `docker scout` melhora | 0.5 dia |
| C-5 | Healthcheck e variáveis obrigatórias no startup | Dockerfiles, `backend/cmd/*.go` | Container falha fast se env ausente; endpoint `/healthz` | 1 dia |

**Gate C:** `docker build .` e `docker run --rm <image>` passam; `trivy fs` sem HIGH/CRITICAL.

---

## 5. Sprint D — Frontend Tests & Dependencies

**Objetivo:** testes verdes e dependências seguras.

| # | Issue | Arquivos principais | Critério de aceitação | Estimativa |
|---|---|---|---|---|
| D-1 | Corrigir 21 testes quebrados | `apps/desktop/src/**/*.test.{ts,tsx}` | `npm run test` passa | 2 dias |
| D-2 | Resolver `form-data` CRLF injection | `package.json`/`package-lock.json` | `npm audit` sem a vulnerabilidade HIGH | 0.5 dia |
| D-3 | Atualizar `vite`/`vitest`/`esbuild` para versões sem advisory | `apps/desktop/package.json` | `npm audit` sem moderate/critical dessas libs; build/test passam | 1-2 dias |
| D-4 | Adicionar testes para `client.ts` (auth, retry, timeout) | `apps/desktop/src/api/__tests__/*.test.ts` | Cobertura mínima 80% dos métodos críticos | 1 dia |
| D-5 | Tipagem rigorosa e props documentadas | `apps/desktop/src/components/ui/*` | Zero erros `tsc --noEmit`; testes não quebram por prop extra | 1 dia |

**Gate D:** `npm run test` e `npm run build` verdes; `npm audit --audit-level=moderate` limpo.

---

## 6. Sprint E — Hardening & Observability

**Objetivo:** endurecer runtime e dar visibilidade operacional.

| # | Issue | Arquivos principais | Critério de aceitação | Estimativa |
|---|---|---|---|---|
| E-1 | Rate limiting nos endpoints de auth/upload/sync | middleware em `backend/internal/api` | Testes de brute-force bloqueados após N tentativas | 1-2 dias |
| E-2 | Validadores centralizados de input | `backend/internal/validation` | UUID, e-mail, path, mime-type validados; testes | 1-2 dias |
| E-3 | Timeouts e context propagation consistentes | `backend/internal/api`, `database`, `sync` | Zero goroutine leaks em `-race`; context cancelado propaga | 1 dia |
| E-4 | Configurar pool de conexões do DB | `backend/internal/database/database.go` | `SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime` configuráveis | 0.5 dia |
| E-5 | Conectar métricas/tracing a backend real | `backend/internal/observability/*` | OTLP/Prometheus opcional; sem panic se backend ausente | 1-2 dias |
| E-6 | Logging estruturado e níveis configuráveis | `backend/internal/log/logger.go` | JSON em produção; nível via env; PII redacted | 1 dia |

**Gate E:** testes de integração de auth/sync/deploy passam; pprof/otel não quebra startup.

---

## 7. Sprint F — Release & Cleanup

**Objetivo:** publicar release estável e documentado.

| # | Issue | Arquivos principais | Critério de aceitação | Estimativa |
|---|---|---|---|---|
| F-1 | Atualizar CHANGELOG, RELEASE_NOTES, RELEASE_CHECKLIST | `CHANGELOG.md`, `docs/RELEASE_CHECKLIST.md` | Histórico claro do v0.1.0 → novo release | 0.5 dia |
| F-2 | Bump de versão | `Cargo.toml`, `tauri.conf.json`, `package.json`, backend | Todos os arquivos de versão alinhados | 0.5 dia |
| F-3 | Gerar instalador e smoke test manual | `npm run tauri -- build` | Instala/desinstala no Windows; app abre; login básico funciona | 1 dia |
| F-4 | Calcular e publicar CHECKSUMS | release assets | SHA256 do `.exe` e do source tarball no release | 0.5 dia |
| F-5 | Deletar release atual e recriar | GitHub Releases | Release antigo removido; novo release com código testado | 0.5 dia |
| F-6 | Remover remote `eniac-systems` obsoleto | git config | `git remote -v` mostra apenas origin | 0.1 dia |

**Gate F:** todos os gates da seção 1 passam; release publicado com checksums.

---

## 8. Suggested Order of Execution

```
A1 → A2 → A3 → A4 → A5 → A6
B1 → B2 → B3 → B4 → B5
C1 → C2 → C3 → C4 → C5
D1 → D2 → D3 → D4 → D5
E1 → E2 → E3 → E4 → E5 → E6
F1 → F2 → F3 → F4 → F5 → F6
```

Sprints A e B podem ser feitas em paralelo por pessoas diferentes. Sprint C depende de B (Dockerfile deve refletir código formatado/lintado). Sprint D pode começar após A. Sprint E após A/B. Sprint F após todos os gates.

---

## 9. Risk & Rollback

- **Risco:** correções de SQLi podem introduzir regressões em queries complexas.
  **Mitigação:** testes de integração com banco real (ou sqlite) cobrindo todos os métodos alterados.
- **Risco:** migração de hash de senha invalida usuários existentes (se houver).
  **Mitigação:** nesta versão ainda não há usuários reais; documentar estratégia de re-hash para futuro.
- **Risco:** atualização de `vite`/`vitest` quebra API de testes.
  **Mitigação:** atualizar com `npm install vite@latest vitest@latest` e rodar testes a cada minor bump.
- **Rollback:** cada commit é pequeno; reverter via `git revert` se gate falhar.

---

## 10. Success Metrics

- Zero findings HIGH/CRITICAL em `trivy fs` e `npm audit`.
- Cobertura de testes backend ≥ 60%, frontend ≥ 70% dos módulos críticos.
- Tempo de build do backend < 30s, frontend < 60s, Tauri < 5min.
- Instalador passa em smoke test Windows limpo.

---

*Plano gerado a partir da auditoria `FULL_CODE_AUDIT_v0.1.0.md`. Atualizar status conforme execução.*

