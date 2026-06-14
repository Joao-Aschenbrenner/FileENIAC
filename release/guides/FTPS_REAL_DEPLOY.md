# FTPS REAL DEPLOY — EPIC 4

**Data:** _______
**Tester:** _______

---

## Pré-requisitos

- [ ] Workspace criado
- [ ] Projeto web real (HTML simples ou PHP) clonado/importado
- [ ] Servidor FTPS 1: HostGator (ou similar)
- [ ] Servidor FTPS 2: Hostinger (ou similar)
- [ ] Credenciais FTPS: host, porta, usuário, senha, target path

## Passo 1 — Configurar Servidor FTPS (HostGator)

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 1.1 | Navegar para Servidores | Página de servidores | |
| 1.2 | Clicar "+ Novo Servidor" | Formulário aparece | |
| 1.3 | Preencher HostGator: host, porta (21), user, pass, target_path | Campos preenchidos | |
| 1.4 | Associar ao projeto | Projeto selecionado no dropdown | |
| 1.5 | Clicar "Salvar Servidor" | Servidor salvo, aparece na lista | |
| 1.6 | Senha NÃO aparece na lista | Campo `password` oculto | |

## Passo 2 — Configurar Hostinger

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 2.1 | Repetir Passo 1 para Hostinger | Segundo servidor salvo | |
| 2.2 | Ambos servidores visíveis na lista | Lista mostra ambos | |

## Passo 3 — Deploy Real

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 3.1 | Navegar para Deploy | Página de deploy com projetos listados | |
| 3.2 | Selecionar projeto | Projeto com servidor associado | |
| 3.3 | Clicar "Executar Deploy" | Loading, mostra progresso | |
| 3.4 | Verificar console eniac.exe | Logs de deploy (pack, upload, verify) | |
| 3.5 | Deploy concluído | Status "success" + deploy_id | |
| 3.6 | Verificar servidor FTPS | Arquivos do projeto estão no target_path | |
| 3.7 | Verificar `deploy-manifest.json` no servidor | JSON com deploy_id, timestamp | |

## Passo 4 — Verify

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 4.1 | Clicar "Verificar" | Status do último deploy | |
| 4.2 | Verify URL configurada | HTTP request para URL (opcional) | |

## Passo 5 — Histórico

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 5.1 | Navegar para Histórico | Evento "Deploy concluído" aparece | |
| 5.2 | Filtrar por tipo "Deploy OK" | Deploy aparece no filtro | |

---

## Checklist Final

- [ ] Servidor HostGator configurado
- [ ] Servidor Hostinger configurado
- [ ] Deploy real executado (upload FTPS)
- [ ] Arquivos verificados no servidor
- [ ] Deploy-manifest.json criado
- [ ] Histórico registrado
- [ ] Senha nunca exposta na UI

**Problemas encontrados:**
_______________________________________________________________________________
