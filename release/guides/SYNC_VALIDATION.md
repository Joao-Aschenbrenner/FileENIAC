# SYNC VALIDATION — EPIC 5

**Data:** _______
**Tester:** _______

---

## Pré-requisitos

- [ ] Projeto com deploy executado (EPIC 4 aprovado)
- [ ] Servidor FTPS com arquivos (após deploy)
- [ ] Acesso ao diretório local do projeto

## Passo 1 — Mirror

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 1.1 | Navegar para Sync | Página de sync | |
| 1.2 | Selecionar projeto | Projeto aparece no dropdown | |
| 1.3 | Clicar "Criar Mirror" | Download do servidor FTPS inicia | |
| 1.4 | Mirror concluído | Status "completed" com contagem de arquivos | |
| 1.5 | Verificar `{workspace}/.eniac/mirror/{projeto}/` | Arquivos do servidor estão lá | |
| 1.6 | Verificar `mirror_manifest.json` | Hashes SHA256 de cada arquivo | |

## Passo 2 — Diff (Sem Alterações)

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 2.1 | Navegar para Diff Viewer | Página de diff | |
| 2.2 | Selecionar projeto | Projeto no dropdown | |
| 2.3 | Clicar "Carregar Diff" | Status "sincronizado" para todos | |

## Passo 3 — Alteração Manual (Local)

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 3.1 | Editar um arquivo local (`index.html`) | Arquivo modificado | |
| 3.2 | Criar um novo arquivo (`novo.html`) | Arquivo adicionado | |
| 3.3 | Remover um arquivo (`old.html`) | Arquivo deletado | |

## Passo 4 — Diff (Com Alterações)

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 4.1 | Clicar "Carregar Diff" | 3 status: "modificado", "novo", "removido" | |
| 4.2 | Verificar hashes | Hashes diferentes para arquivo modificado | |
| 4.3 | Summary correto | Total = 3, New = 1, Modified = 1, Removed = 1 | |

## Passo 5 — Sync

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 5.1 | Navegar para Sync | Projeto selecionado | |
| 5.2 | Clicar "Sync Preview" | Mostra o que será sincronizado (sem executar) | |
| 5.3 | Verificar: alterações destrutivas listadas | Remoção de arquivo marcada | |
| 5.4 | Clicar "Executar Sync" | Modal de confirmação aparece | |
| 5.5 | Confirmar | Sync executado | |

## Passo 6 — Verificar Mirror

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 6.1 | Verificar `mirror/` | `novo.html` existe, `old.html` não existe | |
| 6.2 | Verificar `index.html` | Conteúdo atualizado (hash igual ao local) | |
| 6.3 | Executar Diff novamente | Todos "sincronizado" | |

---

## Checklist Final

- [ ] Mirror criado (download FTPS)
- [ ] Diff detecta: modified, added, removed
- [ ] Hashes SHA256 corretos
- [ ] Sync preview funciona (sem executar)
- [ ] Sync executado com confirmação
- [ ] Mirror atualizado corretamente
- [ ] Nenhum arquivo perdido/sobrescrito incorretamente

**Problemas encontrados:**
_______________________________________________________________________________
