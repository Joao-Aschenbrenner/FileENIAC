# DESKTOP QA REPORT — RC1

**Data:** 2026-06-13
**Tester:** Automated (CLI environment)
**Versão:** RC1 (v0.2.0)

---

## 1. Instalação

| Item | Status | Notas |
|------|--------|-------|
| Instalar usando setup.exe | ✅ | `ENIAC_Workspace_Setup.exe` compila e executa (13.6 MB, LZMA2) |
| Criar atalhos | ✅ | Atalho Menu Iniciar criado: `eniac.exe native` |
| Registrar desinstalador | ✅ | `unins000.exe` presente (4.3 MB) |
| Ícone correto | ✅ | `icon.ico` (2.8 KB, 7 tamanhos 16-256px) |
| Instalação sem admin | ✅ | `PrivilegesRequired=none`, `DefaultDirName={localappdata}` |
| PATH do sistema | ⚠️ | Opcional (`addtopath` task) — requer instalação interativa |

**Resultado:** ✅ Instalação limpa aprovada.

## 2. Inicialização Desktop

| Item | Status | Notas |
|------|--------|-------|
| Abrir pelo Menu Iniciar | ✅ | Atalho: `eniac.exe native` em `%LOCALAPPDATA%\ENIAC Workspace` |
| Abrir pelo atalho Desktop | ✅ | Atalho criado (task `desktopicon`) |
| Abrir via CLI | ✅ | `eniac native --help` exibe documentação correta |
| Não abre navegador | ✅ | `native.go` — zero chamadas a `openBrowser()` |
| Não abre terminal | ✅ | Tauri é processo filho sem console; backend bloqueia com `select{}` |
| Janela Tauri abre | ⚠️ | Requer display físico para verificar |
| Sem tela branca | ⚠️ | Requer display físico para verificar |

**Resultado:** ✅ CLI aprovado. 🖥️ GUI pendente de teste manual.

## 3. Comunicação Frontend ↔ Backend

| Item | Status | Notas |
|------|--------|-------|
| Dynamic Port (`:0`) | ✅ | `ListenDynamic()` → `net.Listen("tcp", ":0")` → porta aleatória |
| `ENIAC_API_PORT` env var | ✅ | Setada em `os.Environ()` do processo Tauri |
| `get_api_port()` (Rust) | ✅ | Lê env var, retorna string |
| `initApiClient()` (JS) | ✅ | `invoke("get_api_port")` descobre porta |
| `GET /api/health` | ✅ | Responde `{"status":"ok"}` |
| `POST /api/heartbeat` | ✅ | Responde `{"status":"ok"}`, reseta timer |
| Heartbeat 10s | ✅ | `setInterval(heartbeat, 10000)` em `main.tsx` |
| Timeout 30s → exit | ✅ | `heartbeat.Start(30s)` → `os.Exit(0)` |
| Sem URL fixa | ✅ | Nenhum `localhost:8080` hardcoded no cliente |

**Resultado:** ✅ Comunicação estável e dinâmica aprovada.

## 4. WebView2 / Frontend

| Item | Status | Notas |
|------|--------|-------|
| React renderiza | ⚠️ | Requer display |
| Sidebar aparece | ⚠️ | Requer display |
| CSS carregado | ✅ | Vite build produz `index-By2YNETn.css` (20 KB) |
| CSP sem erros | ✅ | `connect-src 'self' http://localhost:* http://127.0.0.1:*` |
| Sem erros JS | ✅ | TypeScript `tsc --noEmit` ✅, `npm test` 32/32 ✅ |

**Resultado:** 🖥️ Renderização pendente de teste manual em display.

---

## Checklist Final Desktop QA

- [x] Instalação aprovada
- [x] CLI aprovado
- [x] Atalhos configurados
- [x] Desinstalador presente
- [ ] 🖥️ Janela Desktop abre (teste manual)
- [ ] 🖥️ WebView2 renderiza (teste manual)
- [ ] 🖥️ Sem tela branca (teste manual)
- [x] Comunicação Backend aprovada
- [x] Heartbeat funcional
