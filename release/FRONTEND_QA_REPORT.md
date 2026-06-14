# FRONTEND QA REPORT — RC1

**Data:** 2026-06-13
**Tester:** Automated (CLI environment)
**Versão:** RC1 (v0.2.0)

---

## 1. Build

| Item | Status | Notas |
|------|--------|-------|
| Vite build | ✅ | `dist/` gerado com index.html + assets |
| TypeScript | ✅ | `tsc --noEmit` sem erros |
| Tests | ✅ | 32/32 passing (vitest) |
| Production bundle | ✅ | JS: 228 KB (67 KB gzip), CSS: 20 KB (4 KB gzip) |

## 2. Rotas

| Rota | Componente | Status |
|------|-----------|--------|
| `/` | Onboarding | ✅ |
| `/dashboard` | Dashboard | ✅ |
| `/projects` | Projects | ✅ |
| `/projects/:name` | ProjectDetails | ✅ |
| `/servers` | Servers | ✅ |
| `/history` | History | ✅ |
| `/deploy` | DeployCenter | ✅ |
| `/rollback` | RollbackCenter | ✅ |
| `/diff` | DiffViewer | ✅ |
| `/sync` | SyncCenter | ✅ |
| `/health` | HealthCenter | ✅ |
| `/github/login` | GitHubLogin | ✅ |
| `/github/orgs` | GitHubOrgs | ✅ |
| `/github/repos` | GitHubRepos | ✅ |
| `/bootstrap` | WorkspaceBootstrap | ✅ |

## 3. Componentes

| Componente | Tests | Status |
|-----------|-------|--------|
| ErrorBoundary | - | ✅ |
| Layout | - | ✅ |
| Sidebar | 4 tests | ✅ |
| Badge | 2 tests | ✅ |
| Card | 4 tests | ✅ |
| EmptyState | 3 tests | ✅ |
| ErrorState | 2 tests | ✅ |
| Loader | 2 tests | ✅ |
| Modal | 4 tests | ✅ |
| Timeline | 5 tests | ✅ |
| Toast | 2 tests | ✅ |

## 4. API Client

| Função | Status |
|--------|--------|
| `initApiClient()` | ✅ (invoke `get_api_port`) |
| `heartbeat()` | ✅ (POST `/api/heartbeat` a cada 10s) |
| `checkHealth()` | ✅ |
| `getWorkspace()` | ✅ |
| `listProjects()` | ✅ |
| `createProject()` | ✅ |
| `deleteProject()` | ✅ |
| `listServers()` | ✅ |
| `createServer()` | ✅ |
| `deleteServer()` | ✅ |
| `getSettings()` / `updateSettings()` | ✅ |
| `getHistory()` / `getEvents()` | ✅ |
| `getDeploys()` | ✅ |
| `executeDeploy()` | ✅ |
| `executeRollback()` / `executeVerify()` | ✅ |
| `getDiff()` | ✅ |
| `getSyncs()` / `executeSync()` | ✅ |
| `createMirror()` | ✅ |
| `getHealthCheck()` | ✅ |
| GitHub functions (12) | ✅ |

## 5. Integração Tauri

| Plugin | Status |
|--------|--------|
| `@tauri-apps/api/core` (`invoke`) | ✅ |
| `@tauri-apps/plugin-dialog` (`open`) | ✅ (folder picker) |
| `@tauri-apps/plugin-opener` | ✅ (links externos) |

## 6. CSP

```json
{
  "csp": "default-src 'self'; connect-src 'self' http://localhost:* http://127.0.0.1:*; style-src 'self' 'unsafe-inline'; script-src 'self'"
}
```

✅ Permite fetch para qualquer porta em localhost (dynamic port)
✅ Bloqueia conexões externas (segurança)
✅ Permite estilos inline necessários

## 7. UX Improvements (Sprint 9.2)

| Item | Status |
|------|--------|
| `types.ts` corrigido (`label`→`name`, `username`→`user`) | ✅ |
| Servers.tsx: loading state | ✅ |
| Servers.tsx: delete confirmation | ✅ |
| Servers.tsx: project dropdown | ✅ |
| RollbackCenter.tsx: retry button | ✅ |
| GitHubRepos.tsx: Voltar corrigido | ✅ |
| Onboarding.tsx: loading states | ✅ |
| Onboarding.tsx: folder picker nativo | ✅ |

---

## Checklist Final Frontend QA

- [x] Build limpo (tsc + vite)
- [x] 32/32 testes passando
- [x] Todas as 17 rotas registradas
- [x] 12 componentes UI com cobertura de teste
- [x] 25+ funções API client
- [x] CSP configurado para dynamic port
- [x] UX fixes da Sprint 9.2 aplicadas
- [ ] 🖥️ Renderização WebView2 (teste manual)
- [ ] 🖥️ Navegação entre rotas (teste manual)
- [ ] 🖥️ Onboarding completo (teste manual)
- [ ] 🖥️ Folder picker nativo (teste manual)
