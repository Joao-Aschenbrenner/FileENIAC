# GitHub Import Flow

## Visão Geral
Fluxo completo de descoberta, importação e clonagem de repositórios GitHub para o ENIAC Workspace.

## Pré-requisitos
- Token GitHub com escopos: `repo`, `read:org`
- Backend rodando (`eniac serve`)
- Workspace configurado

## Etapas

### 1. Autenticar
```
POST /github/login {"token": "ghp_xxx"}
```
- Token validado contra API GitHub
- Criptografado e armazenado no vault
- Usuário registrado em `github_user`

### 2. Listar Organizações
```
GET /github/organizations
```
- Retorna organizações que o token tem acesso
- Inclui login, id, url, avatar

### 3. Listar Repositórios
```
GET /github/repositories?org=myorg
GET /github/repositories (repos pessoais)
```
- Retorna repositórios com metadados
- Repositórios já importados marcados como `imported: true`

### 4. Importar
```
POST /github/import
{
  "repos": [
    {
      "id": 123456,
      "name": "meu-projeto",
      "full_name": "myorg/meu-projeto",
      "clone_url": "https://github.com/myorg/meu-projeto.git",
      "default_branch": "main",
      "organization": "myorg"
    }
  ]
}
```
Para cada repositório:
1. `registry.AddProject` — cria projeto no DB
2. `registry.AddRepository` — registra repositório
3. `clone.Clone` — executa `git clone`
4. `registry.UpdateRepositoryImport` — atualiza status
5. `validate.ValidateClone` — valida clone

### 5. Verificar
```
GET /repositories?org=myorg
```
- Lista repositórios importados
- Status: pending, cloned, validation_failed

## Estrutura de Diretórios
```
{workspace}/
  .eniac/
    config.toml
    workspace.db
  projects/
    meu-projeto/     ← clone do repositório
      .git/
      index.php
      ...
```

## Tratamento de Erros
- Token inválido → HTTP 401
- Repositório já existe → erro, não sobrescreve
- Clone falha → import_status = "clone_failed"
- Validação falha → import_status = "validation_failed"
