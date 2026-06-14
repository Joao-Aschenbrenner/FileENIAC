# FIRST LAUNCH REPORT — EPIC 1

**Data:** _______
**Tester:** _______
**Versão:** RC1.1 (v0.2.0)

---

## Pré-requisitos

- [ ] Windows 10/11 com WebView2 Runtime (Edge)
- [ ] `ENIAC_Workspace_Setup.exe` baixado
- [ ] Nenhuma instalação anterior do ENIAC (clean slate)
- [ ] Captura de tela / gravação disponível

## Passo 1 — Instalação

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 1.1 | Executar `ENIAC_Workspace_Setup.exe` | Instalador UAC não pede admin | |
| 1.2 | Clicar "Next > Next > Install" | Barra de progresso completa | |
| 1.3 | Clicar "Finish" | Janela fecha | |
| 1.4 | Verificar Iniciar Menu | Atalho "ENIAC Workspace" existe | |
| 1.5 | Verificar Desktop | Atalho "ENIAC Workspace" existe | |
| 1.6 | Verificar `%LOCALAPPDATA%\ENIAC Workspace\` | `eniac.exe` + `ENIAC Workspace.exe` + `WebView2Loader.dll` existem | |

## Passo 2 — Primeiro Launch via Atalho

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 2.1 | Dar duplo clique no atalho "ENIAC Workspace" | Janela Tauri abre em <5s | |
| 2.2 | Verificar título da janela | "ENIAC Workspace" | |
| 2.3 | Verificar tamanho | ~1100x720 (não maximizada) | |
| 2.4 | Verificar se navegador padrão NÃO abriu | Nenhuma janela/aba do Edge/Chrome | |
| 2.5 | Verificar console do WebView2 (F12) | Sem erros vermelhos | |
| 2.6 | Verificar tela branca | Conteúdo renderizado (não branco) | |

## Passo 3 — Interface Inicial

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 3.1 | Sidebar visível | Menu lateral com ícones + labels | |
| 3.2 | Sidebar contém: | Dashboard, Projetos, Servidores, GitHub, Diff, Sync, Deploy, Histórico | |
| 3.3 | Clicar em cada item da sidebar | Rota muda, página carrega (loader aparece e some) | |
| 3.4 | Nenhuma página mostra tela branca | Todas as rotas têm conteúdo | |
| 3.5 | Botão "Fechar" (X) na janela | Janela fecha, backend também fecha (heartbeat timeout) | |

## Passo 4 — Heartbeat Funciona

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 4.1 | Abrir app, aguardar 5s | App continua rodando | |
| 4.2 | Fechar janela Tauri (X) | Backend morre em até 30s | |
| 4.3 | Verificar Task Manager | `eniac.exe` não está mais rodando após 30s | |

---

## Checklist Final

- [ ] Instalação completa sem admin
- [ ] Janela Tauri abre sem navegador
- [ ] Todas as 8+ rotas carregam sem tela branca
- [ ] Heartbeat funciona (backend morre com janela)
- [ ] Sem erros no console WebView2

**Notas do testador:**
_______________________________________________________________________________
_______________________________________________________________________________
