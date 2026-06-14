# ADR-006: Deploy Engine — Orquestração de Deploy

## Status
APROVADO

## Data
2026-06-10

## Contexto
O eniac-deploy (PoC) implementou o fluxo básico: pack → FTPS upload → trigger HTTP → verify. Agora o deploy precisa ser orquestrado dentro do workspace, com suporte a múltiplos projetos, dependências, manifest e rollback.

## Decisão
O Deploy Engine orquestra todo o fluxo, delegando transporte ao FTPS Engine.

### Fluxo de Deploy
1. **Pre-flight**: verifica workspace, registry, conexão FTPS
2. **Build**: gera artifact tar.gz (com excludes)
3. **Backup**: copia diretório remoto atual para `.backup.{timestamp}`
4. **Upload**: envia artifact via FTPS Engine
5. **Extract**: servidor recebe e extrai artifact
6. **Verify**: health check HTTP pós-deploy
7. **Manifest**: gera deploy-manifest.json no servidor
8. **Record**: registra no History Engine

### Manifest
Cada deploy gera `deploy-manifest.json`:
```json
{
  "deploy_id": "dep_a1b2c3d4",
  "project": "simple-finance",
  "commit": "a1b2c3d4e5f6...",
  "branch": "main",
  "timestamp": "2026-06-10T15:30:00Z",
  "files": 128,
  "artifact_hash": "sha256:...",
  "migrations": ["2026_06_10_create_users_table"],
  "status": "success"
}
```

### Rollback
- Restaura do backup `.backup.{timestamp}`
- Não desfaz migrations (requer migration rollback manual)
- Registra rollback no History Engine

### Comandos
```
fileeniac deploy push --project <name> [--fallback] [--dry-run]
fileeniac deploy rollback --project <name>
fileeniac deploy verify --project <name>
fileeniac deploy status --project <name>
```

### Segurança
- Token HMAC-SHA256 com 5min TTL
- Endpoint dinâmico para bypass de ModSecurity
- Fallback para FTPS mirror se HTTP trigger falhar

## Consequências
- Deploy Engine herda e substitui o push da PoC
- Manifest permite rastreabilidade total
- Rollback seguro preserva backups no servidor
- Integração com History Engine garante auditoria
