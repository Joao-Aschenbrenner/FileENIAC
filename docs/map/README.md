# FileENIAC — Mapa do Projeto

## Sumário

1. [Visão Geral](#1-visão-geral)
2. [Arquitetura](#2-arquitetura)
3. [Fluxo de Funcionamento](#3-fluxo-de-funcionamento)
4. [Comandos CLI](#4-comandos-cli)
5. [Menus do Frontend](#5-menus-do-frontend)
6. [API Endpoints](#6-api-endpoints)
7. [Banco de Dados](#7-banco-de-dados)
8. [Integração GitHub](#8-integração-github)
9. [Deploy FTPS](#9-deploy-ftps)
10. [Segurança](#10-segurança)
11. [Build & Instalação](#11-build--instalação)

---

## 1. Visão Geral

**FileENIAC** é uma plataforma desktop de gerenciamento de deploys FTPS com mirror, diff, sync e histórico. Ajuda equipes a gerenciar deploys de projetos web para servidores remotos com segurança e rastreabilidade.

- **Nome**: FileENIAC (homenagem ao FileZilla + ENIAC)
- **Linguagens**: Go (backend), TypeScript/React (frontend), Rust (Tauri shell)
- **Armazenamento**: SQLite local (workspace-based)
- **Rede**: FTPS explícito (TLS), GitHub REST API
- **Criptografia**: AES-256-GCM para tokens e senhas

---

## 2. Arquitetura

```
┌──────────────────────────────────────────────────────────┐
│                     SEU COMPUTADOR                       │
│                                                          │
│  ┌──────────────┐    ┌──────────────────────────────┐   │
│  │   Go Backend  │    │      Tauri (WebView2)        │   │
│  │  fileeniac-cli│    │       FileENIAC.exe          │   │
│  │               │    │                              │   │
│  │  ┌─────────┐  │    │  ┌────────────────────────┐  │   │
│  │  │ HTTP API │◄─┼────┼──│  React Frontend       │  │   │
│  │  │ :8080    │  │    │  │  (BrowserContext)     │  │   │
│  │  └────┬────┘  │    │  └────────────────────────┘  │   │
│  │       │       │    └──────────────────────────────┘   │
│  │  ┌────▼────┐  │              ▲                        │
│  │  │ SQLite  │  │              │ redirect               │
│  │  │ .db     │  │              │                        │
│  │  └─────────┘  │    ┌─────────┴──────────┐            │
│  │               │    │  Env Var           │            │
│  │  ┌─────────┐  │    │ FILEENIAC_API_PORT │            │
│  │  │ FTPS    │──┼────┤ = 8080             │            │
│  │  │ Client  │  │    └────────────────────┘            │
│  │  └─────────┘  │                                      │
│  │               │                                      │
│  └───────────────┘                                      │
│         │                                               │
│         ▼                                               │
│  ┌──────────────┐                                       │
│  │ .eniac/      │                                       │
│  │ workspace    │                                       │
│  └──────────────┘                                       │
└──────────────────────────────────────────────────────────┘
```

### Processos

| Processo | Binário | Função |
|----------|---------|--------|
| **Backend** | `fileeniac-cli.exe` | Servidor HTTP, API REST, serve frontend, lógica de negócio |
| **Desktop** | `FileENIAC.exe` | Janela WebView2 que redireciona para o backend |

### Fluxo de Inicialização (`native`)

1. `fileeniac-cli native` inicia o backend Go na porta `:8080`
2. Seta env var `FILEENIAC_API_PORT=8080`
3. Seta env var `ENIAC_API_PORT=8080` (compatibilidade)
4. Após 800ms, spawna `FileENIAC.exe` (Tauri)
5. Tauri lê `FILEENIAC_API_PORT` do ambiente
6. WebView redireciona para `http://localhost:8080/`
7. Frontend carrega do backend (mesma origem)
8. `initApiClient()` chama `invoke("get_api_port")` → fallback `8080`
9. `BASE_URL = http://localhost:8080/api`

---

## 3. Fluxo de Funcionamento

### Ciclo de Vida Completo

```
1. Workspace Init
   └── workspace init → cria .eniac/ + SQLite DB
   
2. GitHub Auth (opcional)
   └── github/login → token PAT → criptografado no vault
   
3. Importar Repositórios (opcional)
   └── github/repos → seleciona → importa
       ├── Cria Project
       ├── git clone --depth 1
       └── Valida clone
   
4. Adicionar Projeto Manual (alternativo)
   └── project add → nome, path local, branch

5. Adicionar Servidor
   └── server add → host, porta, usuário, senha, path remoto

6. Mirror (opcional)
   └── mirror create → baixa estado remoto via FTPS
       ├── Salva em .eniac/mirror/{projeto}/
       └── snapshot no DB

7. Diff (opcional, pós-mirror)
   └── diff local-mirror → compara SHA-256
       ├── Arquivos novos / modificados / removidos
       └── Estado: divergente / sincronizado

8. Sync (opcional, pós-diff)
   └── sync plan → analisa diff
   └── sync apply → copia arquivos (confirmado pelo usuário)

9. Deploy
   └── deploy run
       ├── Pack: projeto local → TAR.GZ
       ├── Checksum SHA-256
       ├── Upload FTPS (com retry + circuit breaker)
       ├── Download + verifica integridade
       ├── Upload deploy-manifest.json
       └── Registro no histórico

10. Verify
    └── deploy verify → checa status do último deploy

11. Rollback
    └── deploy rollback → deleta artifact + manifest do servidor
```

### Heartbeat

- Frontend envia POST `/api/heartbeat` a cada 10s
- Backend zera timer de 30s a cada heartbeat
- Se heartbeat não chegar em 30s, backend faz `os.Exit(0)`
- **Por quê?** Evita processos órfãos se a janela for fechada

### Background Health

- Runner executa a cada 30s
- Coleta: validade do token, contagem de projetos/servidores, divergentes
- Disponível em `GET /api/health/background`

---

## 4. Comandos CLI

| Comando | Descrição |
|---------|-----------|
| `fileeniac` | Ajuda geral |
| `fileeniac version` | Versão e data de build |
| `fileeniac native` | **Modo principal**: backend + janela nativa WebView2 |
| `fileeniac desktop` | Fallback: backend + navegador |
| `fileeniac serve` | Apenas backend HTTP (sem frontend) |
| `fileeniac workspace init` | Inicializa workspace (.eniac/) |
| `fileeniac workspace open` | Abre workspace existente |
| `fileeniac workspace status` | Status do workspace atual |
| `fileeniac workspace scan` | Escaneia diretórios por .eniac/ |
| `fileeniac project add/remove/list/show` | CRUD de projetos |
| `fileeniac server add/remove/list/show` | CRUD de servidores |
| `fileeniac deploy run` | Executa deploy |
| `fileeniac deploy verify` | Verifica último deploy |
| `fileeniac deploy rollback` | Reverte deploy |
| `fileeniac deploy history` | Histórico de deploys |
| `fileeniac sync plan/apply/reconcile` | Sincronia local↔mirror |
| `fileeniac diff local-mirror/status` | Diferenças entre local e mirror |
| `fileeniac mirror create/status` | Cria/consulta snapshot remoto |
| `fileeniac config get/set/list` | Gerencia configurações |
| `fileeniac auth login/status/logout` | Gerencia token GitHub |
| `fileeniac repo add/remove/list/show` | Gerencia associações git |
| `fileeniac update-from` | Auto-update (substitui binários) |

---

## 5. Menus do Frontend

### Sidebar (11 itens)

| Ícone | Rota | Descrição |
|-------|------|-----------|
| ◉ | `/dashboard` | **Dashboard** — visão geral do workspace (projetos, servidores, divergentes, eventos) |
| 🚀 | `/bootstrap` | **Bootstrap** — wizard passo-a-passo para configurar workspace do zero |
| 📁 | `/projects` | **Projetos** — lista/cria/deleta projetos |
| 🖥 | `/servers` | **Servidores** — lista/cria/deleta servidores FTPS |
| 🐙 | `/github/login` | **GitHub** — login com PAT, organizações, repositórios, importação |
| 🚀 | `/deploy` | **Deploy** — seleciona projeto, executa deploy (normal/fallback) |
| ⏪ | `/rollback` | **Rollback** — seleciona projeto, confirma, reverte |
| 🔄 | `/sync` | **Sync** — analisa divergência, executa sync local↔mirror |
| 📊 | `/diff` | **Diff** — diff local-mirror arquivo por arquivo |
| 📋 | `/history` | **Histórico** — eventos e logs do workspace |
| ❤️ | `/health` | **Saúde** — checagens detalhadas de saúde |

### Onboarding (tela inicial, fora da sidebar)

| Etapa | Ação |
|-------|------|
| 1. Welcome | "Começar" → verifica backend health |
| 2. Config | Informa caminho do workspace ou "Procurar" (picker nativo) |
| 3. Ready | Confirma conexão → "Entrar no Workspace" → `/dashboard` |

### Páginas sem sidebar

| Rota | Descrição |
|------|-----------|
| `/` | Onboarding |
| `/github/orgs` | Selecionar organização GitHub |
| `/github/repos` | Selecionar e importar repositórios |

---

## 6. API Endpoints

### Health & Monitoramento

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/api/health` | Liveness probe (`{"status":"ok"}`) |
| GET | `/api/health/check` | Health check completo (projetos, servidores, divergentes, eventos) |
| GET | `/api/health/background` | Snapshot do background runner |
| POST | `/api/heartbeat` | Reseta heartbeat watchdog |

### Workspace

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/api/workspace` | Status do workspace ativo |

### Projetos

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/api/projects` | Lista projetos |
| POST | `/api/projects` | Cria projeto |
| GET | `/api/projects/{name}` | Detalhes do projeto |
| DELETE | `/api/projects/{name}` | Deleta projeto |

### Servidores

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/api/servers?project=` | Lista servidores (filtro opcional) |
| POST | `/api/servers` | Cria servidor |
| GET | `/api/servers/{id}` | Detalhes do servidor |
| DELETE | `/api/servers/{id}` | Deleta servidor |

### Deploy

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/api/deploys?project=&limit=` | Logs de deploy |
| POST | `/api/deploy` | Executa deploy |
| POST | `/api/rollback` | Executa rollback |
| POST | `/api/verify` | Verifica último deploy |

### Diff & Sync

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/api/diff?project=` | Diff local-mirror |
| GET | `/api/syncs?project=&limit=` | Manifests de sync |
| POST | `/api/sync` | Executa sync |

### Mirror

| Método | Rota | Descrição |
|--------|------|-----------|
| POST | `/api/mirror` | Cria snapshot mirror |

### Configurações

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/api/settings` | Lista configurações |
| POST | `/api/settings` | Atualiza configurações |

### Histórico

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/api/history?project=&type=&limit=&offset=` | Histórico de deploys/eventos |
| GET | `/api/events?type=&limit=&offset=` | Lista eventos |

### GitHub

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/api/github/status` | Status da autenticação |
| POST | `/api/github/login` | Armazena token |
| POST | `/api/github/logout` | Remove token |
| GET | `/api/github/organizations` | Lista organizações |
| GET | `/api/github/repositories?org=` | Lista repositórios |
| POST | `/api/github/import` | Importa repositórios |
| POST | `/api/github/clone` | Clona repositório |

### Repositories

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/api/repositories?org=` | Lista repositórios importados |
| GET | `/api/repositories/{id}` | Detalhes do repositório |

### Manutenção

| Método | Rota | Descrição |
|--------|------|-----------|
| POST | `/api/refresh/github` | Atualiza metadados GitHub |
| POST | `/api/revalidate` | Revalida todos os clones |
| GET | `/api/readiness/deploy?project=` | Pré-checagem de deploy |
| GET | `/api/readiness/sync?project=` | Pré-checagem de sync |
| GET | `/api/repair/check` | Verifica consistência |
| POST | `/api/repair/fix` | Corrige inconsistências |

---

## 7. Banco de Dados

SQLite em `.eniac/workspace.db`.

### Migrations

| V | Tabelas |
|---|---------|
| V1 | `projects`, `servers` |
| V2 | `events` |
| V3 | `deploy_logs`, `rollback_logs`, `workspace_settings` |
| V4 | Colunas em projects + `mirror_snapshots`, `sync_manifests` |
| V5 | Colunas GitHub + `repositories` |

### Tabelas Principais

- **`projects`**: name, local_path, remote_path, branch, server_id, last_deploy_id, divergence_status, import_status, clone_path, github_id, organization, repo_name, provider
- **`servers`**: name, host, port, username, password (encrypted), target_path, project_id
- **`events`**: id, project_id, event_type, description, created_at (17 tipos de evento)
- **`deploy_logs`**: id, project_id, deploy_id, status, artifact_path, checksum, error_message, created_at
- **`rollback_logs`**: id, project_id, deploy_id, status, error_message, created_at
- **`workspace_settings`**: key, value (token criptografado, timeout, etc.)
- **`mirror_snapshots`**: id, project_id, file_count, total_size, created_at
- **`sync_manifests`**: id, project_id, action, direction, file_count, created_at
- **`repositories`**: github_id, name, full_name, private, description, default_branch, language, organization, owner, url, clone_url, imported

---

## 8. Integração GitHub

### Fluxo

1. **Login**: PAT token → valida via `/user` → criptografa com AES-256-GCM → armazena em `workspace_settings`
2. **Organizações**: GET `/user/orgs` → frontend exibe lista
3. **Repositórios**: GET `/orgs/{org}/repos` ou `/user/repos` → frontend exibe com flag `imported`
4. **Importação**: Seleciona repositórios → POST `/api/github/import`
   - Cria `Project` no registry
   - Cria `Repository` no registry
   - `git clone --depth 1 --branch {default}`
   - Valida clone (fsck, branch, remote)
   - Atualiza status no project
   - Dispara eventos
5. **Refresh**: POST `/api/refresh/github` → atualiza metadados de todos os repositórios
6. **Repair**: GET`/repair/check` → encontra órfãos → POST`/repair/fix` → associa pelo nome do projeto

### Escopos do Token

- `repo` — acesso a repositórios privados
- `read:org` — ler organizações

---

## 9. Deploy FTPS

### Pipeline

```
[Local Project] → Pack (TAR.GZ) → Checksum (SHA-256) → FTPS Upload
                                                          ↓
                                            [Remote Server]
                                                          ↓
                     FTPS Download ← Verify Checksum ← [Remote File]
                                                          ↓
                                                deploy-manifest.json
```

### Características

- **FTPS explícito** (TLS na porta 21, depois AUTH TLS)
- **Retry**: até 3 tentativas com backoff exponencial
- **Circuit Breaker**: abre após falhas consecutivas
- **Verificação**: baixa o artifact de volta, compara SHA-256
- **Manifesto**: `deploy-manifest.json` com ID, checksum, timestamp
- **Rollback**: deleta artifact + manifest do servidor remoto

---

## 10. Segurança

- **Vault**: AES-256-GCM para tokens GitHub e senhas FTPS
- **Chave mestra**: gerada no `workspace init`, armazenada em `.eniac/config.toml`
- **Senhas mascaradas**: nunca retornadas nas respostas da API
- **FTPS**: `InsecureSkipVerify: false` (verifica certificado)
- **Heartbeat**: auto-shutdown se frontend parar de enviar heartbeats (30s)
- **CSP Tauri**: `connect-src 'self' http://localhost:* http://127.0.0.1:*`
- **Validação de input**: git clone valida URL (sem espaços, sem leading dash)
- **Confirmação do usuário**: deploy, rollback, sync exigem confirmação explícita

---

## 11. Build & Instalação

### Pré-requisitos

- Go 1.21+ com CGO (`gcc` em `C:\msys64\ucrt64\bin\gcc.exe`)
- Node.js 18+ 
- Rust 1.70+ (para Tauri)
- WebView2 Runtime (Windows 10+ já inclui)
- Inno Setup 6 (para instalador)

### Build

```bash
# Backend Go
cd backend
go build -trimpath -ldflags="-s -w -linkmode=external" -o ..\bin\fileeniac-cli.exe .

# Frontend
cd apps\desktop
npm run build

# Tauri
npx tauri build --no-bundle

# Instalador
& "C:\Users\USUARIO\AppData\Local\Programs\Inno Setup 6\ISCC.exe" build\installer\installer.iss
```

### Executáveis

| Arquivo | Tamanho | Função |
|---------|---------|--------|
| `bin\fileeniac-cli.exe` | ~13 MB | Go backend (CLI + servidor HTTP) |
| `bin\FileENIAC.exe` | ~20 MB | Tauri shell (WebView2) |
| `build\installer\FileENIAC_Setup.exe` | ~10 MB | Instalador Windows |

### Docker

```bash
docker build -t enisystems/fileeniac:latest .
docker run -p 8080:8080 -v meu-workspace:/workspace enisystems/fileeniac:latest
```
