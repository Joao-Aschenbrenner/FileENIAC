# ADR-007: Desktop UI — Interface Principal

## Status
APROVADO

## Data
2026-06-10

## Contexto
A ferramenta CLI é funcional, mas não oferece a experiência visual necessária para adoção em equipe. É preciso uma interface desktop que torne o workspace visível e operável.

## Decisão
Implementar Desktop App utilizando Tauri v2 + React + TypeScript.

### Stack
- **Framework**: Tauri v2 (Rust backend + WebView frontend)
- **Frontend**: React 18 + TypeScript
- **UI Library**: Mantine v7
- **State**: Zustand
- **Data Fetching**: TanStack Query v5
- **Build**: Vite

### Regra Fundamental
Frontend NUNCA executa:
- FTP/FTPS
- Git
- Deploy
- Rollback
- Qualquer regra de negócio

Frontend APENAS consome APIs do backend Go.

### Telas Mínimas (Sprint 4)
1. **Dashboard** — visão geral do workspace, status dos projetos
2. **Projetos** — lista, cadastro, configuração
3. **Deploy** — push, rollback, histórico
4. **Monitor** — health checks, logs
5. **Configurações** — workspace, servidores, GitHub

### Arquitetura de Comunicação
```
Desktop UI (React)
    ↓ HTTP/JSON (localhost)
Backend Go API (porta dinâmica)
    ↓
FTPS Engine | Git Engine | Deploy Engine | History Engine
```

### Backend Sidecar
O backend Go roda como sidecar process gerenciado pelo Tauri:
- Inicia automático com a UI
- Porta aleatória (evita conflito)
- Autenticação via token local
- Logs streamados via WebSocket para UI

## Consequências
- UI independente do SO (Windows, macOS, Linux)
- Separação clara entre apresentação e negócio
- Backend pode ser usado sem UI (CLI puro)
- Tauri permite distribuição binária única
