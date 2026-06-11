# ADR-010: Security — Segurança do Workspace

## Status
APROVADO

## Data
2026-06-10

## Contexto
O workspace gerencia credenciais de servidores FTPS, tokens GitHub e chaves secretas. É essencial que essas informações nunca sejam armazenadas em texto puro.

## Decisão
Implementar camada de segurança com criptografia e credential manager.

### Credential Manager
Utilizar credential manager nativo do SO:

| SO | Sistema | Implementação |
|----|---------|---------------|
| Windows | Credential Manager | `go-credential-windows` |
| macOS | Keychain | `go-keychain` |
| Linux | Secret Service | `libsecret` / `keyring` |

### Criptografia em Repouso
- Algoritmo: AES-256-GCM
- Chave derivada do login do sistema (DPAPI no Windows, Keychain no macOS)
- Utilizado para cache de tokens e configurações sensíveis

### O que é criptografado
- Tokens GitHub
- Senhas FTPS
- Secrets de deploy
- Tokens de refresh

### O que NÃO é criptografado
- Configurações de workspace (não sensíveis)
- Hostnames e portas
- Cache de metadados (hashes, timestamps)

### Variáveis de Ambiente
Credenciais podem ser injetadas via environment variables:

```toml
[server.ftps]
pass = "{{ENV:FTP_PASSWORD}}"
```

O sistema substitui `{{ENV:NOME}}` pelo valor da variável de ambiente em tempo de execução.

### Permissões de Arquivo
- `.eniac/config.toml`: 0600 (apenas dono)
- `.eniac/mirror/`: 0700
- `.eniac/history.db`: 0600

### Deploy Security
- Token HMAC-SHA256 com TTL de 5 minutos
- Endpoint de deploy renomeado dinamicamente
- Validação de origem (CORS no servidor)
- Log de todas as tentativas de deploy

## Consequências
- Credenciais protegidas por criptografia do SO
- Nenhuma senha em texto puro em disco
- Ambiente configurável via variáveis de ambiente
- Compatível com padrões de segurança corporativos
