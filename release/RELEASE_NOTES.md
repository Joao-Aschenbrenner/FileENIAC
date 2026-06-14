# ENIAC Workspace RC1

**Versão:** 1.0.0-rc1
**Data:** 2026-06-13

---

## Artefatos

| Arquivo | Tamanho | SHA256 |
|---------|---------|--------|
| `eniac.exe` | 18.1 MB | `11ef7c322d1bdd1e3602e2c5bde7c4c5158f36c2abe7f8602072d337246ea9b7` |
| `ENIAC Workspace.exe` | 19.8 MB | `788cf85fdd9125d6cc4af74f95b83ab870fa9e6278d18bfaaa31b3dbabd3cb26` |
| `ENIAC_Workspace_Setup.exe` | 13.6 MB | `dab02fa33c068e2b499e71b0c4fe3da23bc76ae692cfdb5544badd5e0c07b672` |
| `WebView2Loader.dll` | 160 KB | `8427b1fc58ec707813e5c0a51eb5d69397bb333250a7b891be4d3b123f1e0f1c` |
| `icon.ico` | 2.8 KB | `e0eca73f7f5b71ceffed9067d107b16521022b14abb49c442689a452d4e81234` |

## O que há de novo

### P0 — Native Desktop Mode
- `eniac native` comando agora usa **exclusivamente** a janela Tauri (WebView2)
- Nenhum navegador é aberto — nem como fallback
- Tauri carrega frontend embedado, IPC (`invoke`) funciona nativamente

### P0 — Dynamic Port
- Backend usa `:0` (porta aleatória)
- Porta real propagada via `ENIAC_API_PORT` para o processo Tauri
- Frontend descobre porta via `invoke("get_api_port")`
- Nenhuma URL fixa `localhost:8080` hardcoded

### P0 — Heartbeat
- Frontend envia `POST /api/heartbeat` a cada 10 segundos
- Backend encerra (`os.Exit(0)`) após 30 segundos sem heartbeat
- Ao fechar a janela Tauri, o backend morre automaticamente

### P0 — CORS
- Middleware `Access-Control-Allow-Origin: *` no servidor Go
- Resolvido erro "Backend offline" no desktop nativo

### Sprint 9.2 — Correções de UX
- `types.ts` corrigido (`label`→`name`, `username`→`user`)
- `Servers.tsx` com loading state + confirmação + dropdown projetos
- `RollbackCenter.tsx` com retry
- `GitHubRepos.tsx` com Voltar corrigido
- `Onboarding.tsx` com loading states
- Botão "Procurar" (folder picker nativo) com `tauri-plugin-dialog`

### Sprint 9.2 — Auto-update
- `eniac version` exibe v0.2.0
- `eniac update-from <path>` — detecta `{installDir}/update/`, faz backup, substitui e reinicia

### Sprint 9.2 — Instalador sem admin
- `DefaultDirName={localappdata}\ENIAC Workspace`
- `PrivilegesRequired=none`
- Registry em HKCU

## Mudanças Técnicas

### Go Backend
- `backend/internal/heartbeat/` — novo pacote timer para heartbeat
- `backend/internal/api/api.go` — `ListenDynamic()` com `net.Listen("tcp", ":0")`
- `backend/internal/api/api.go` — `POST /api/heartbeat` endpoint
- `backend/cmd/native.go` — porta `:0`, sem `openBrowser()`, heartbeat start
- `backend/cmd/desktop.go` — mantido como fallback (browser mode)

### Tauri Desktop
- `apps/desktop/src-tauri/src/lib.rs` — comando `get_api_port()` via env var
- `apps/desktop/src/api/client.ts` — `initApiClient()`, `heartbeat()`
- `apps/desktop/src/main.tsx` — inicialização dinâmica + heartbeat loop
- CSP configurado: `connect-src 'self' http://localhost:* http://127.0.0.1:*`

## Como Usar

```powershell
# Instalar
ENIAC_Workspace_Setup.exe

# Iniciar (janela nativa)
eniac native

# Iniciar (navegador, fallback)
eniac desktop

# Apenas servidor API
eniac serve

# Verificar versão
eniac version
```

## Build

```powershell
# Backend
cd backend
$env:CGO_ENABLED=1
$env:CC='C:\msys64\ucrt64\bin\gcc.exe'
go build -ldflags="-linkmode=internal" -o ..\build\eniac.exe .

# Desktop (Tauri)
cd apps/desktop
npm run build
npx tauri build --no-bundle

# Instalador
cd build
& "$env:LOCALAPPDATA\Programs\Inno Setup 6\ISCC.exe" installer.iss
```
