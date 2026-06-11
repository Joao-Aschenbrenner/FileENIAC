# ADR-002: Project Registry — Cadastro e Descoberta

## Status
APROVADO

## Data
2026-06-10

## Contexto
O sistema precisa gerenciar múltiplos projetos dentro de um workspace. Cada projeto tem configurações específicas (caminho, repositório, servidor, dependências) que precisam ser registradas e descobertas automaticamente.

## Decisão
Criar um módulo **Project Registry** responsável por:

### Cadastro
- Adicionar projeto ao workspace
- Associar repositório GitHub (opcional)
- Configurar servidor FTPS
- Definir dependências entre projetos

### Descoberta
- Escanear diretório do workspace em busca de repositórios Git
- Detectar projetos não registrados
- Sugerir importação automática

### Estrutura do registro
```json
{
  "workspace": "ENIAC",
  "version": 1,
  "projects": {
    "simple-finance": {
      "name": "SimpleFinance",
      "path": "C:/workspace/SimpleFinance",
      "remote": "https://github.com/ENIACSystems/SimpleFinance.git",
      "branch": "main",
      "server": {
        "type": "ftps",
        "host": "ftp.example.com",
        "port": 21,
        "user": "user@example.com",
        "target_path": "/public_html/projects/simple-finance",
        "verify_url": "https://example.com/projects/simple-finance/"
      },
      "dependencies": ["shared-lib"],
      "deploy": {
        "run_migrations": true,
        "backup_prefix": ".backup",
        "endpoint": "_deploy.php"
      }
    }
  }
}
```

### Comandos CLI
```
eniac-workspace registry list
eniac-workspace registry add <name> --path <path>
eniac-workspace registry remove <name>
eniac-workspace registry sync    # Detecta projetos não registrados
eniac-workspace registry doctor  # Verifica integridade do registro
```

## Consequências
- Registry eliminou a necessidade de config manual projeto por projeto
- Projetos podem ser descobertos e importados em lote
- Dependências entre projetos permitem ordenação de deploy
