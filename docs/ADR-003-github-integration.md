# ADR-003: GitHub Integration

## Status
APROVADO

## Data
2026-06-10

## Contexto
Configurar projetos manualmente é repetitivo e propenso a erro. O workspace precisa de integração direta com GitHub para descoberta automática de repositórios, clone e sincronização.

## Decisão
Implementar integração com GitHub via OAuth 2.0.

### Autenticação
- OAuth 2.0 com escopos: `repo`, `read:org`, `workflow`
- Tokens armazenados no credential manager do SO:
  - Windows: Credential Manager
  - macOS: Keychain
  - Linux: Secret Service (libsecret)
- Refresh automático de tokens
- Criptografia em repouso (AES-256-GCM)

### Fluxo de Discovery
1. Login via OAuth
2. Listar organizações do usuário
3. Listar repositórios da organização
4. Detectar repositórios com `composer.json` ou `package.json`
5. Apresentar checklist para importação
6. Clonar selecionados
7. Registrar no Project Registry

### Fluxo de Bootstrap (máquina nova)
1. Login GitHub
2. Selecionar workspace remoto (ou criar novo)
3. Sistema clona todos os projetos do workspace
4. Configura registry
5. Workspace pronto

### Comandos CLI
```
fileeniac github login
fileeniac github status
fileeniac github discover
fileeniac github import [--all]
fileeniac github sync
fileeniac github logout
```

## Consequências
- Credenciais nunca em texto puro
- Setup de nova máquina reduzido a minutos
- Sincronização automática entre GitHub e workspace local
- Necessário gerenciar refresh de tokens OAuth
