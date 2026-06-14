# ADR-000: FileENIAC — Visão do Produto

## Status
APROVADO

## Data
2026-06-10

## Contexto
O projeto nasceu como uma ferramenta de deploy FTP (eniac-deploy) para resolver problemas de publicação em hosting compartilhado HostGator. Durante o desenvolvimento, ficou evidente que o problema real não é deploy, mas sim a falta de rastreabilidade e sincronização entre workspace local, GitHub e servidor.

## Decisão
O produto evolui para **FileENIAC**: uma plataforma única para gerenciamento de workspace local, repositórios Git, GitHub, deploys, FTPS, histórico, auditoria, monitoramento e integração com IA.

## Princípios
1. **Fonte da verdade = Git** — nunca utilizar o servidor como fonte da verdade
2. **Workspace First** — o sistema opera sobre workspaces, não projetos isolados
3. **Mirror seguro** — alterações do servidor vão para mirror, nunca direto no projeto local
4. **Rastreabilidade total** — toda operação é registrada com contexto (commit, deploy, data, responsável)

## Perguntas que o produto deve responder
- Qual commit está em produção?
- Qual deploy publicou esse commit?
- Qual projeto está divergente?
- O servidor está igual ao Git?
- Posso voltar para qual versão?

## Stack
- **Backend**: Go 1.21+
- **Desktop**: Tauri v2 + React + TypeScript + Mantine
- **Database**: SQLite (WAL mode)
- **Mobile**: Futuro (Flutter)

## Roadmap
- Sprint 0: Arquitetura e ADRs
- Sprint 1: FTPS Engine + Deploy + History
- Sprint 2: Workspace + Registry + Mirror
- Sprint 3: GitHub OAuth + Discovery + Bootstrap
- Sprint 4: Desktop UI
- Sprint 5: Agent API + IA

## Consequências
- eniac-deploy existente torna-se PoC FTPS, não produto final
- Toda nova funcionalidade deve respeitar esta arquitetura
- Nenhuma implementação sem ADR aprovado
