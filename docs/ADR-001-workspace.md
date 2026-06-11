# ADR-001: Workspace — Entidade Central

## Status
APROVADO

## Data
2026-06-10

## Contexto
O sistema precisa de uma entidade central que agrupe projetos relacionados. Trabalhar com projetos isolados sem contexto de workspace gera perda de visibilidade das relações entre projetos (ex: shared-lib é dependência de SimpleFinance).

## Decisão
Workspace é a entidade central do sistema. Todo projeto pertence a um workspace. Um workspace contém:

- Nome único
- Diretório local (caminho absoluto)
- Lista de projetos registrados
- Configurações compartilhadas (FTPS padrão, secrets, etc.)
- Histórico consolidado

## Estrutura do diretório .eniac
Cada workspace possui um diretório oculto `.eniac/` na raiz:

```
.eniac/
  config.toml          # Configuração do workspace
  registry.json        # Registro de projetos
  mirror/              # Cópia do servidor para diff seguro
    project-a/
    project-b/
  history.db           # SQLite com histórico consolidado
  cache/               # Cache temporário
```

## Fluxo de bootstrap
1. Usuário cria workspace (ou importa existente)
2. Adiciona projetos (manual ou via descoberta GitHub)
3. Sistema clona repositórios
4. Registra projetos no registry
5. Workspace pronto para operações

## Consequências
- Toda operação (deploy, sync, verify) é sempre no contexto de um workspace
- Remover um projeto do workspace não exclui o repositório local nem o remoto
- Workspace settings podem ser sobrescritas por projeto individual
