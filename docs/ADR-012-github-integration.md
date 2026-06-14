# ADR-012: GitHub Integration

## Status
Aprovado

## Contexto
O FileENIAC precisa se integrar ao GitHub para descobrir, importar e configurar repositórios automaticamente. Atualmente a configuração de projetos é manual e repetitiva.

## Decisão
Adotar integração com GitHub via Personal Access Token (classic) utilizando a API REST v3.

## Detalhes Técnicos

### Autenticação
- Token armazenado no Vault (AES-256-GCM) através do `github_token` nas `workspace_settings`
- Token validado contra `GET /user` antes de salvar
- Logout remove token das settings

### Descoberta
- `GET /user/orgs` lista organizações
- `GET /orgs/{org}/repos` lista repositórios por organização
- `GET /user/repos` lista repositórios pessoais
- Repositórios já importados marcados como `imported: true`

### Importação
- Cria projeto no registry
- Registra repositório na tabela `repositories`
- Executa `git clone --depth 1 --branch {branch}` no diretório `{workspace}/projects/{name}`
- Atualiza `import_status` do projeto e repositório
- Registra evento `github_import`

### Validação Pós-Importação
- Verifica existência do diretório `.git`
- Verifica diretório do projeto
- Verifica associação com servidor (warning)

### Tabela repositories
- `github_id` — ID do GitHub (unique)
- `full_name` — `org/repo`
- `organization` — organização dona
- `import_status` — pending/cloned/validation_failed
- `project_id` — FK para projects
- `clone_path` — diretório local do clone

### Campos Adicionados no projects
- `github_id`, `organization`, `repo_name`, `import_status`, `clone_path`, `provider`, `last_sync_commit`

## Consequências
- Token fica criptografado no DB
- Clone automático elimina setup manual
- Registry expandido rastreia origem GitHub
- Projetos criados em `{workspace}/projects/{repo}`

## Endpoints da API
- `GET /github/status`
- `POST /github/login`
- `POST /github/logout`
- `GET /github/organizations`
- `GET /github/repositories?org=`
- `POST /github/import`
- `POST /github/clone`
- `GET /repositories`
- `GET /repositories/:id`
