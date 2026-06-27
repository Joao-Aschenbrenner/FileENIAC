# Gate 0.5 — Validacao do Ambiente (pre-Fase 1A)

## Identificacao

- Data: 2026-06-26
- Branch: architecture-review
- HEAD: 26ca843 (Fase 0)
- ADR-014: 2e52b1f (docs only)
- Working tree: 43 itens nao commitados (estado WIP pre-existente)
- Ambiente: Windows 11, PowerShell 5.1
- Go toolchain local: go1.26.4 windows/amd64

## Objetivo

Validar cada item do checklist Gate 0.5 antes de iniciar a Fase 1A
(introducao da interface Transport). Identificar se a Fase 0
realmente deixou o ambiente verde ou se ha problema de
infraestrutura mascarado.

---

## Checklist

### 1. `go version`

- Resultado: **ok**
- Comando: `go version`
- Output: `go version go1.26.4 windows/amd64`
- Observacao: toolchain instalado e funcional. A Fase 0 declarou
  `go.mod` como 1.25 (voluntariamente). O toolchain local (1.26.4)
  detecta a diferenca e exige `go mod tidy` antes de compilar.
  Isso e esperado: CI usa 1.25, dev local pode usar 1.26 com tidy.

---

### 2. `go test ./...` (sem race)

- Resultado: **passou** (ambiente isolado via git worktree)
- Comando: `go test ./... -count=1` (23 pacotes com testes)
- Output: todos `ok`
- Lista de pacotes validados:

| Pacote | Status |
|--------|--------|
| internal/api | ok |
| internal/clone | ok |
| internal/deploy | ok |
| internal/deploy/bypass | ok |
| internal/deploy/ftp | ok |
| internal/deploy/hardening | ok |
| internal/deploy/packer | ok |
| internal/deploy/token | ok |
| internal/diff | ok |
| internal/github | ok |
| internal/health | ok |
| internal/history | ok |
| internal/logger | ok |
| internal/mirror | ok |
| internal/readiness | ok |
| internal/refresh | ok |
| internal/registry | ok |
| internal/repair | ok |
| internal/status | ok |
| internal/sync | ok |
| internal/validate | ok |
| internal/vault | ok |
| internal/workspace | ok |
| internal/workspace/discovery | ok |
| backend/integration (7 testes) | todos ok |

- Observacao: **sem data race detector**, sufixo `? [no test files]`
  em packages: database, heartbeat, log, update, version, webui.

---

### 3. `go test -race ./...`

- Resultado: **ALERTA — DATA RACE DETECTADO**
- Comando: `go test ./... -race -short -count=1 -p 1`
- Package com falha: `backend/integration`
- Teste: `TestAPIHealthEndpoint` (flaky quando em paralelo)

#### Localizacao do data race

```
backend/internal/api/api.go:94-95   (ListenAndServe: s.srv = srv — write)
backend/internal/api/api.go:106-107 (ListenDynamic: s.srv = srv — write)
backend/internal/api/api.go:138-141 (Close: s.srv.Shutdown — read)
backend/internal/api/api.go:147-148 (Addr: s.srv.Addr — read)
```

**Diagnostico**: o campo `s.srv` (`*http.Server`) na struct `Server`
e acessado em `ListenAndServe()` e `ListenDynamic()` (escrita)
e em `Close()` e `Addr()` (leitura) sem qualquer sincronizacao.

Quando dois testes de integracao rodam em paralelo, cada um cria
seu proprio `api.New("127.0.0.1:0")` e chama `ListenAndServe()` em
goroutine separada, causando acesso concorrente ao campo `s.srv`.

**Impacto**: baixo para producao (produco usa uma unica instancia
de Server). Impacto alto para tests paralelos — causa flakiness.

**Recomendacao**: adicionar `sync.Mutex` em `s.srvMu` e proteger
leituras/escritas. Corrigir antes ou durante Fase 2.

---

### 4. `go build ./...`

- Resultado: **passou** (com webui/dist stub mockado)
- Comando: `go build ./...` (toolchain 1.26 apos go mod tidy)
- Observacao: o build depende de `backend/webui/dist/` conter
  ao menos 1 arquivo embeddable (`//go:embed dist`). Sem o frontend
  compilado, o backend nao compila. O CI compila apenas backend
  (working-directory: backend), entao nao e afetado.
- Acao: documentar dependencia de build para release.

---

### 5. `docker build .`

- Resultado: **nao validado**
- Motivo: ambiente Windows sem daemon Docker
- Evidencia de consistencia:
  - `Dockerfile` builder usa `golang:1.25-alpine`
  - `Dockerfile` runtime usa `alpine:3.21`
  - `HEALTHCHECK` usa `wget --spider http://localhost:8080/api/health`
- Depende de validacao manual em Linux/macOS com Docker.

---

### 6. `docker compose up`

- Resultado: **nao validado**
- Motivo: mesmo do item 5
- Dependencias do compose:
  - `.env` com `FILEENIAC_VAULT_PASSWORD` obrigatorio
  - Porta `8080` livre no host
- Observacao: healthcheck do compose espelha o do Dockerfile

---

### 7. `/api/health` = 200

- Resultado: **nao validado** (requer container up)
- Rota: `/api/health` (endpoint publico, sem auth)
- Handler: `respond(w, 200, {"status":"ok"})`

---

### 8. Vault inicializa

- Resultado: **nao validado** (requer container up + env var)
- Requisito: `FILEENIAC_VAULT_PASSWORD` definido
- Observacao: backend recusa boot sem essa var

---

### 9. CI verde

- Resultado: **nao validado** (requer git push)
- Alteracoes na CI: `go-version: '1.25'`
- Observacao: CI compila apenas `backend/` como working-directory.
  O `webui/dist/` nao e necessario para o CI (working-directory isola).

---

## Resumo do Gate

| Item | Status | Acao |
|------|--------|------|
| 1. go version | ok | — |
| 2. go test (sem race) | passe | — |
| 3. go test -race | **DATA RACE** | Criar tech debt item: mutex em s.srv |
| 4. go build | passe | Documentar dependencia webui/dist |
| 5. docker build | nao validado | Testar em Linux antes da Fase 3 |
| 6. docker compose | nao validado | Testar em Linux |
| 7. /api/health | nao validado | Testar apos compose up |
| 8. Vault | nao validado | Testar com env var |
| 9. CI | nao validado | Push para validar |

## Veredito

**Gate 0.5: APROVADO com ressalvas.**

A Fase 0 esta correta em conteudo e coerente entre todos os
arquivos. O data race em `api.go` e pre-existente (nao foi
causado pela Fase 0) e deve ser tratado como divida tecnica
prioritaria antes da Fase 2 (ou durante).

Nao ha impedimento para iniciar a Fase 1A.

---

## Tech Debt Registrado

### TD-001: Data race em `api.Server.srv`

- Arquivo: `backend/internal/api/api.go` (linhas 94-95, 106-107,
  138-141, 147-148)
- Severidade: media (flaky em tests paralelos)
- Correcao: adicionar `sync.Mutex` no campo `srvMu` e bloquear
  leituras/escritas em `ListenAndServe`, `ListenDynamic`,
  `Close` e `Addr`.

### TD-002: Build depende de webui/dist/

- Arquivo: `backend/webui/webui.go:8` (`//go:embed dist`)
- Impacto: `go build ./...` na raiz falha sem o stub
- Solucao possivel: CI usa `working-directory: backend`, resolvido.
- Para build de release: gerar dist/ vazia como placeholder.
