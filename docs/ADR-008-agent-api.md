# ADR-008: Agent API — Integração com IA

## Status
APROVADO

## Data
2026-06-10

## Contexto
O workspace acumula dados valiosos (deploys, histórico, divergências) que podem ser consultados por agentes de IA para diagnóstico, recomendação e automação.

## Decisão
Definir contratos para Agent API. NÃO IMPLEMENTAR AGORA.

### API REST (futura)
```yaml
openapi: 3.0.0
info:
  title: FileENIAC Agent API
  version: 0.1.0

paths:
  /api/v1/workspace/{id}/status:
    get:
      summary: Status completo do workspace
      responses:
        200:
          schema: WorkspaceStatus

  /api/v1/workspace/{id}/divergences:
    get:
      summary: Projetos divergentes (local vs GitHub vs servidor)

  /api/v1/workspace/{id}/deploys:
    get:
      summary: Últimos deploys

  /api/v1/projects/{name}/history:
    get:
      summary: Histórico completo do projeto

  /api/v1/projects/{name}/health:
    get:
      summary: Health check do projeto
```

### Endpoints Futuros
- `GET /workspace/{id}/status` — estado completo
- `GET /workspace/{id}/divergences` — projetos divergentes
- `GET /workspace/{id}/deploys` — histórico de deploys
- `GET /projects/{name}/history` — histórico do projeto
- `GET /projects/{name}/health` — health check
- `POST /deploy` — gatilho de deploy (com validação)
- `POST /rollback` — gatilho de rollback

### Casos de Uso para IA
1. "Qual commit está em produção no SimpleFinance?"
2. "O servidor está igual ao GitHub?"
3. "Quantos deploys falharam essa semana?"
4. "Preciso reverter o último deploy do SnapSelect"
5. "Quais projetos estão com dependências desatualizadas?"

### Formato de Resposta
```json
{
  "query": "Qual commit está em produção no SimpleFinance?",
  "answer": "O commit a1b2c3d4 (main) está em produção desde 10/06/2026 às 15:30.",
  "confidence": 0.95,
  "source": "deploy-manifest.json",
  "timestamp": "2026-06-10T16:00:00Z"
}
```

## Consequências
- API versionada e documentada
- Agentes de IA podem consultar sem modificar estado
- Futuro: permitir que IA execute ações com autorização
