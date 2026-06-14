# ROLLBACK VALIDATION — EPIC 6

**Data:** _______
**Tester:** _______

---

## Pré-requisitos

- [ ] Projeto com deploy executado (EPIC 4 aprovado)
- [ ] Sync validado (EPIC 5 aprovado)
- [ ] Arquivos no servidor FTPS

## Passo 1 — Deploy A (Baseline)

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 1.1 | Garantir que projeto tem 3+ arquivos | Projeto populado | |
| 1.2 | Executar Deploy | Deploy A concluído | |
| 1.3 | Anotar `deploy_id` | `DEPLOY_A = ______` | |
| 1.4 | Verificar servidor FTPS | 3+ arquivos no target_path | |

## Passo 2 — Deploy B (Nova Versão)

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 2.1 | Modificar `index.html` localmente | Conteúdo diferente | |
| 2.2 | Adicionar `nova-pagina.html` | Novo arquivo | |
| 2.3 | Executar Deploy | Deploy B concluído | |
| 2.4 | Anotar `deploy_id` | `DEPLOY_B = ______` | |
| 2.5 | Verificar servidor FTPS | `nova-pagina.html` existe | |

## Passo 3 — Rollback

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 3.1 | Navegar para Rollback Center | Projetos com deploy listados | |
| 3.2 | Selecionar projeto | Deploy B é o último | |
| 3.3 | Clicar "Rollback" | Modal de confirmação aparece | |
| 3.4 | Confirmar | Rollback executado | |
| 3.5 | Verificar servidor FTPS | `nova-pagina.html` NÃO existe mais | |
| 3.6 | Verificar `index.html` | Volta ao conteúdo do Deploy A | |
| 3.7 | Status no Rollback Center | "rolled_back" | |

## Passo 4 — Histórico de Rollback

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 4.1 | Navegar para Histórico | Evento "Rollback" aparece | |
| 4.2 | Filtrar por Rollback | Rollback visível no filtro | |
| 4.3 | Verificar detalhes | `deploy_id` correto | |

## Passo 5 — Deploy Após Rollback

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 5.1 | Executar Deploy novamente | Novo deploy_id (não reusa) | |
| 5.2 | Verificar histórico | 3 deploys + 1 rollback registrados | |

---

## Checklist Final

- [ ] Deploy A executado (baseline)
- [ ] Deploy B executado (nova versão)
- [ ] Rollback reverte para A
- [ ] Arquivos restaurados corretamente
- [ ] Histórico correto
- [ ] Pode fazer deploy novamente após rollback

**Problemas encontrados:**
_______________________________________________________________________________
