# Desktop Onboarding Flow

## Visão Geral
Fluxo de onboarding do Desktop para configurar o ENIAC Workspace do zero.

## Fluxo

### Tela Inicial (Onboarding)
1. **Boas-vindas** — Tela com gradiente ENIAC, botão "Começar"
2. **Verificar Backend** — Checa `GET /health`, se falha mostra erro
3. **Conectar Workspace** — Input de caminho, valida `GET /workspace`
4. **Confirmação** — Resumo do workspace, botão "Entrar"

### Bootstrap (opcional, acessível via sidebar)
Após entrar, usuário pode:
1. Acessar `/bootstrap` para ver checklist
2. Conectar GitHub em `/github/login`
3. Selecionar organização em `/github/orgs`
4. Selecionar repositórios em `/github/repos`
5. Importar repositórios com clone automático
6. Ver projetos importados em `/projects`

### Telas do Desktop

| Rota | Componente | Descrição |
|------|-----------|-----------|
| `/` | Onboarding | Wizard inicial (backend + workspace) |
| `/dashboard` | Dashboard | Métricas do workspace |
| `/bootstrap` | WorkspaceBootstrap | Checklist de configuração |
| `/projects` | Projects | CRUD de projetos |
| `/projects/:name` | ProjectDetails | Detalhes + ações |
| `/servers` | Servers | CRUD de servidores |
| `/github/login` | GitHubLogin | Autenticação GitHub |
| `/github/orgs` | GitHubOrgs | Seleção de organização |
| `/github/repos` | GitHubRepos | Seleção e importação |
| `/deploy` | DeployCenter | Central de deploy |
| `/rollback` | RollbackCenter | Central de rollback |
| `/sync` | SyncCenter | Central de sincronização |
| `/diff` | DiffViewer | Visualizador de diff |
| `/history` | History | Timeline de eventos |
| `/health` | HealthCenter | Monitoramento |

## Estados

### Loading
- Loader centralizado com texto descritivo
- Usado durante carregamento de dados

### Empty
- Mensagem amigável com ação opcional
- Usado quando não há dados (ex: "Nenhum projeto")

### Error
- Card vermelho com mensagem + botão "Tentar novamente"
- Usado em todas as páginas com chamadas API

### Toast
- Notificação flutuante (sucesso/erro/info)
- Auto-dismiss em 4 segundos
