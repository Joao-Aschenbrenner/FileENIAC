# ADR-004: FTPS Engine

## Status
APROVADO

## Data
2026-06-10

## Contexto
O deploy para hosting compartilhado (HostGator) exige FTPS com TLS explícito na porta 21. O eniac-deploy (PoC) já implementa esse transporte. Agora a engine deve ser modular e reutilizável dentro da arquitetura ENIAC Workspace.

## Decisão
O FTPS Engine é um módulo interno responsável APENAS pela camada de transporte. Toda orquestração (pack, mirror, manifest, rollback) fica no Deploy Engine.

### Responsabilidades
- Conexão FTPS (TLS explícito)
- Upload de arquivos
- Criação de diretórios
- Listagem remota
- Renomeação remota
- Exclusão remota
- Health check básico (NoOp + PWD)

### Não responsabilidades
- Orquestração de deploy
- Geração de artifact tar.gz
- Gerenciamento de manifest
- Rollback

### Stack
- Go 1.21+
- github.com/jlaffaye/ftp (com DialWithExplicitTLS)
- DialOptions configuráveis (timeout, TLS config, debug)

### Interface
```go
type FTPSClient interface {
    Connect(cfg FTPSConfig) error
    Disconnect() error
    Upload(localPath, remotePath string) error
    Download(remotePath, localPath string) error
    Delete(remotePath string) error
    Rename(oldPath, newPath string) error
    List(path string) ([]string, error)
    EnsureDir(path string) error
    Health() error
}
```

### Configuração
```toml
[server.ftps]
host = "ftp.example.com"
port = 21
user = "user@example.com"
pass = "{{ENV:FTP_PASSWORD}}"
timeout = 120
```

## Consequências
- FTPS Engine substitui o client.go da PoC
- Compatível com HostGator e qualquer servidor FTPS
- Credenciais via environment variables ou credential manager
- Fallback para FTP plano se TLS falhar (configurável)
