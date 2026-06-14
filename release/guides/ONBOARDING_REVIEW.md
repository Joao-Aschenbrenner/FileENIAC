# ONBOARDING REVIEW — EPIC 2

**Data:** _______
**Tester:** _______

---

## Pré-requisitos

- [ ] App instalado e rodando (FIRST_LAUNCH aprovado)
- [ ] Diretório vazio para workspace (`C:\Users\...\meu-workspace`)

## Fluxo 1 — Criar Workspace (Onboarding Page)

| # | Ação | Resultado Esperado | ✅ / ❌ | Observação |
|---|------|-------------------|---------|------------|
| 1.1 | Abrir app sem workspace | Tela de onboarding aparece | | |
| 1.2 | Campo de path vazio, clicar "Iniciar" | Mensagem de erro clara ("path é obrigatório") | | |
| 1.3 | Clicar "Procurar" (folder picker) | Diálogo nativo de pasta abre | | |
| 1.4 | Selecionar diretório | Path aparece no campo | | |
| 1.5 | Clicar "Iniciar" com path válido | Loading aparece, workspace é criado | | |
| 1.6 | Após criação | Redireciona para Dashboard | | |
| 1.7 | Verificar `{workspace}/.eniac/` | `config.toml` + `workspace.db` existem | | |

## Fluxo 2 — Erros e Tratamento

| # | Ação | Resultado Esperado | ✅ / ❌ | Observação |
|---|------|-------------------|---------|------------|
| 2.1 | Path pointing to read-only dir | Mensagem de erro clara (permissão) | | |
| 2.2 | Path pointing to existing workspace | Mensagem informa que já existe | | |
| 2.3 | Clicar "Procurar" e cancelar | Campo não muda, nenhum erro | | |
| 2.4 | Fechar app, abrir de novo | Onboarding NÃO aparece (workspace já existe) | | |

## Fluxo 3 — UX Review

Perguntas para o testador:

| # | Pergunta | Resposta |
|---|----------|----------|
| 3.1 | A mensagem "Criar workspace" é clara? | |
| 3.2 | O botão "Procurar" é fácil de encontrar? | |
| 3.3 | Os erros são explicativos? | |
| 3.4 | Há algum termo técnico confuso? | |
| 3.5 | Quantos cliques até chegar ao Dashboard? | |
| 3.6 | Algum passo parece desnecessário? | |

## Pontos Confusos

Registre aqui qualquer momento em que você hesitou ou não sabia o que fazer:

1. _________________________________________________________________________
2. _________________________________________________________________________
3. _________________________________________________________________________

## Campos Mal Explicados

| Campo | Problema | Sugestão |
|-------|----------|----------|
| | | |
| | | |

---

## Checklist Final

- [ ] Onboarding aparece sem workspace
- [ ] Folder picker nativo funciona
- [ ] Workspace criado corretamente
- [ ] Erros são claros e explicativos
- [ ] Fluxo completo em ≤ 3 cliques
- [ ] Nenhum termo técnico desnecessário
