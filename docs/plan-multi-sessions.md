# Plano: Sessões Multi-Workspace

## Problema Atual
- Workspace é singleton global (`activeContext`)
- GitHub token fica na settings do workspace (criptografado com vault key do workspace)
- Path do workspace nunca é salvo no frontend
- Onboarding pergunta workspace toda vez
- Zero suporte a múltiplos workspaces

## Arquitetura Nova

### 1. Banco Central (`sessions.db`)
Localização: `%APPDATA%/FileENIAC/sessions.db` (ou junto do binário)

Tabela `sessions`:
```sql
id              INTEGER PRIMARY KEY AUTOINCREMENT
name            TEXT NOT NULL
description     TEXT DEFAULT ''
workspace_path  TEXT NOT NULL          -- path com .eniac/ dentro
github_token    TEXT DEFAULT ''         -- criptografado (AES-256-GCM)
github_user     TEXT DEFAULT ''
created_at      TEXT DEFAULT (datetime('now'))
updated_at      TEXT DEFAULT (datetime('now'))
```

### 2. Pacote `session` (backend/internal/session/)
- `store.go` — Open/Create/List/Get/Update/Delete operations + migrations
- A sessão guarda:
  - Nome + descrição (definidos pelo usuário)
  - Path do workspace (para `workspace.Open()`)
  - GitHub token (independente do workspace vault)
- Vault próprio: chave mestra em `%APPDATA%/FileENIAC/.vault_key`

### 3. API — Novos Endpoints (sem middleware `requireWorkspace`)
| Método | Rota | Função |
|--------|------|--------|
| GET | `/api/sessions` | Lista todas |
| POST | `/api/sessions` | Cria (name, desc, workspace_path, github_token?) |
| PUT | `/api/sessions/{id}` | Atualiza nome/desc |
| DELETE | `/api/sessions/{id}` | Remove |
| POST | `/api/sessions/{id}/activate` | Ativa (abre workspace + carrega GitHub token) |
| GET | `/api/sessions/active` | Retorna sessão ativa + status do workspace |

### 4. API — GitHub modificado
- GitHub token lido da **sessão ativa** (não mais do workspace settings)
- `getGitHubToken()` busca em: sessão ativa → workspace settings (fallback)
- Login/Logout afetam a sessão ativa, não o workspace

### 5. Frontend — Novo Fluxo
```
App abre → SessionSelector (rota /)
  ├── Lista sessões salvas (cards com nome + desc + info)
  ├── [Nova Sessão] → wizard 3 passos:
  │   1. Nome + Descrição (obrigatório)
  │   2. GitHub (opcional — PAT + validação)
  │   3. Workspace (selecionar path OU init novo)
  │   → Cria + Ativa → /dashboard
  └── [Click numa sessão] → activate → /dashboard
```

### 6. Frontend — Componentes Novos
- `pages/SessionSelector.tsx` — landing page, lista sessões
- `pages/SessionWizard.tsx` — criação passo a passo
- `context/SessionContext.tsx` — React Context com { session, switchSession, list }
- `components/ui/SessionSwitcher.tsx` — dropdown no header pra trocar rápido

### 7. Frontend — Modificações
- `App.tsx`: rota `/` → SessionSelector. Se sessão ativa → /dashboard
- `Onboarding.tsx`: **removido** (substituído pelo wizard)
- `Sidebar.tsx`: nome da sessão no topo
- `Header.tsx`: badge com "Sessão: {nome}"
- `api/client.ts`: novos métodos + `ws()` lê da session context

## Ordem de Implementação

### Fase 1 — Backend: Store + API
1. `backend/internal/session/store.go` — DB, migrations, CRUD
2. `backend/internal/session/vault.go` — encrypt/decrypt GitHub token
3. `backend/internal/api/sessions.go` — handlers novos
4. `backend/internal/api/api.go` — registrar rotas + modificar `getGitHubToken()`
5. `backend/cmd/serve.go` — init sessions DB na inicialização
6. `go test ./...` ✅

### Fase 2 — Frontend: Session Context + API client
7. `apps/desktop/src/context/SessionContext.tsx`
8. `apps/desktop/src/api/client.ts` — add session methods
9. `apps/desktop/src/pages/SessionSelector.tsx`
10. `apps/desktop/src/pages/SessionWizard.tsx` (3 passos)
11. `apps/desktop/src/components/ui/SessionSwitcher.tsx`

### Fase 3 — Frontend: Integração
12. `App.tsx` — novas rotas, SessionContext provider
13. `Sidebar.tsx` — session name
14. `Header.tsx` — session badge
15. Remover `Onboarding.tsx`
16. Ajustar `WorkspaceBootstrap.tsx` para usar sessão ativa

### Fase 4 — Build + Test
17. Rebuild Go backend
18. Rebuild Tauri app
19. Teste: criar 2+ sessões, alternar, verificar dados isolados

## Arquivos Críticos
- `backend/internal/api/api.go` — registrar rotas + modificar middleware/helpers
- `backend/internal/api/sessions.go` — NOVO: handlers de sessão
- `backend/internal/session/store.go` — NOVO: store + migrations
- `backend/internal/session/vault.go` — NOVO: encrypt/decrypt
- `apps/desktop/src/App.tsx` — router
- `apps/desktop/src/context/SessionContext.tsx` — NOVO
- `apps/desktop/src/pages/SessionSelector.tsx` — NOVO
- `apps/desktop/src/pages/SessionWizard.tsx` — NOVO
- `apps/desktop/src/components/ui/SessionSwitcher.tsx` — NOVO
- `apps/desktop/src/api/client.ts` — add session API calls
- `apps/desktop/src/pages/Sidebar.tsx` — mod
- `apps/desktop/src/pages/Header.tsx` — mod

## Verificação
1. Iniciar backend → `GET /api/sessions` retorna `[]`
2. Criar sessão via API → retorna session com ID
3. Ativar sessão → workspace abre, GitHub endpoints funcionam
4. Frontend: SessionSelector mostra lista
5. Wizard cria sessão → redireciona pra dashboard
6. Trocar sessão → SessionSwitcher no header
7. Tudo roda via duplo clique no `FileENIAC.exe`
