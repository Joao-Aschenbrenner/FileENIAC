# GITHUB REAL USAGE — EPIC 3

**Data:** _______
**Tester:** _______
**Token GitHub:** `ghp_***` (fine-grained, perms: repo + org:read)

---

## Pré-requisitos

- [ ] Workspace criado (ONBOARDING aprovado)
- [ ] Token GitHub com permissões `repo` e `read:org`
- [ ] Acesso à internet (API github.com)

## Passo 1 — Login GitHub

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 1.1 | Navegar para GitHub (sidebar) | Página de login aparece | |
| 1.2 | Clicar "Login com GitHub" | Campo de token aparece | |
| 1.3 | Inserir token inválido ("abc") | Erro: "token inválido" ou "unauthorized" | |
| 1.4 | Inserir token válido | Loading aparece, depois mostra "Conectado como {user}" | |
| 1.5 | Verificar `workspace.db` | `github_token` armazenado criptografado | |
| 1.6 | Recarregar página (F5) | Estado de login persiste | |
| 1.7 | Clicar "Desconectar" | Token removido, volta para tela de login | |

## Passo 2 — Listar Organizações

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 2.1 | Após login, clicar "Importar" | Lista de organizações carrega | |
| 2.2 | Performance: tempo de carregamento | < 3s (ou indicador de loading) | |
| 2.3 | Organizações sem repositórios | Aparecem na lista normalmente | |

## Passo 3 — Importar Repositório

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 3.1 | Clicar em uma organização | Lista de repositórios carrega | |
| 3.2 | Verificar repositórios já importados | Marcados como "Importado" | |
| 3.3 | Selecionar 1 repositório público pequeno | Checklist aparece | |
| 3.4 | Clicar "Importar" | Clone inicia (mostra progresso) | |
| 3.5 | Verificar diretório `{workspace}/projects/{repo}` | `.git` + arquivos existem | |
| 3.6 | Verificar `workspace.db` | `projects` + `repositories` populados | |
| 3.7 | Verificar Dashboard | Projeto aparece com status | |

## Passo 4 — Importar Múltiplos

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 4.1 | Selecionar 3 repositórios | Todos são importados | |
| 4.2 | Verificar performance | Clone em paralelo (não sequencial) | |
| 4.3 | Rate limit | Nenhum erro 403 da API GitHub | |
| 4.4 | Repositório já importado | Mensagem: "já importado" (não duplicado) | |

## Passo 5 — Refresh

| # | Ação | Resultado Esperado | ✅ / ❌ |
|---|------|-------------------|---------|
| 5.1 | Clicar "Atualizar" (refresh) | Dados sincronizados com GitHub | |
| 5.2 | Repositório removido do GitHub | Status atualizado | |

---

## Checklist Final

- [ ] Login com token funciona
- [ ] Token armazenado criptografado
- [ ] Organizações listadas
- [ ] Repositórios importados e clonados
- [ ] Performance aceitável (sem rate limit)
- [ ] Refresh funciona
- [ ] Logout limpa dados

**Problemas encontrados:**
_______________________________________________________________________________
