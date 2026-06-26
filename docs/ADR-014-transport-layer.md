# ADR-014 — Transport Layer & Engine Boundaries

## Status

ACEITO (para implementação)

---

## 1. Contexto

O FileENIAC atualmente executa sincronização e deploy via FTPS, com acoplamento direto ao cliente FTP em múltiplos pacotes.

Existem três pontos de uso do FTP no sistema:

- `internal/deploy/ftp/client.go` — implementação direta FTP
- `internal/mirror/mirror.go` — uso direto de `github.com/jlaffaye/ftp` (vazamento de abstração)
- `internal/deploy/service.go` — interface parcial `ftpClientIface` (abstração local não compartilhada)

---

## 2. Problema

O sistema possui múltiplas "verdades" sobre transporte:

- FTP está diretamente acoplado ao domínio
- Mirror e Deploy não compartilham abstração comum
- Existe uma interface local não reutilizável (`ftpClientIface`)
- A Engine não possui boundary formal de transporte

Isso impede:

- substituição de FTP por outros protocolos no futuro
- testabilidade consistente (mock padronizado)
- separação clara entre Engine e infraestrutura

---

## 3. Decisão Arquitetural

Será introduzida uma camada única de abstração chamada:

> **Transport Layer**

A Engine (Sync / Deploy / Mirror) não conhece FTP.

Ela conhece apenas a interface `Transport`.

---

## 4. Modelo de Arquitetura

```
Engine (Sync / Deploy / Mirror)
            │
            ▼
     Transport Interface
            │
   ┌────────┴────────┐
   ▼                 ▼
FTP Transport   (future: SFTP, S3, WebDAV)
```

---

## 5. Interface Transport (Mínima e Estável)

A interface cobre apenas o necessário para operação FTPS atual.

```go
type Transport interface {
    Connect(cfg TransportConfig) error
    Close() error
    IsConnected() bool

    Stat(remotePath string) (FileInfo, error)
    List(remoteDir string) ([]FileInfo, error)

    Upload(localPath, remotePath string) error
    Download(remotePath, localPath string) error

    Delete(remotePath string) error
    Rename(oldPath, newPath string) error
    Mkdir(remotePath string) error
}
```

### Tipos auxiliares

```go
type TransportConfig struct {
    Endpoint string
    User     string
    Pass     string
    Protocol string // "ftp"
    Timeout  time.Duration
}
```

```go
type FileInfo struct {
    Path  string
    Size  int64
    MTime time.Time
    IsDir bool
}
```

---

## 6. Regras Arquiteturais

### Regra 1 — Engine não conhece FTP

Nenhum pacote fora de `internal/transports/*` pode importar:

```
github.com/jlaffaye/ftp
```

### Regra 2 — Transporte é responsabilidade isolada

Somente `internal/transports/` pode conter:

- implementação FTP
- dependências de rede
- bibliotecas externas de protocolo

### Regra 3 — Mirror, Sync e Deploy usam Transport

Nenhum desses módulos pode acessar FTP diretamente.

### Regra 4 — FTP é detalhe de implementação

FTP não é parte do domínio. É apenas um backend do Transport.

---

## 7. Estrutura de Diretórios

```
internal/
    transports/
        transports.go        // interface + tipos base
        factory.go           // criação de transport
        ftp/
            ftp.go           // implementação FTP

    sync/
    deploy/
    mirror/
```

---

## 8. Factory de Transport

```go
type Factory func(cfg TransportConfig) (Transport, error)

var registry = map[string]Factory{
    "ftp": NewFTPTransport,
}

func New(cfg TransportConfig) (Transport, error) {
    f, ok := registry[cfg.Protocol]
    if !ok {
        return nil, fmt.Errorf("unknown transport: %s", cfg.Protocol)
    }
    return f(cfg)
}
```

---

## 9. Estratégia de Migração

A migração ocorrerá em duas fases:

### Fase 1A — Introdução da abstração

- Criar `internal/transports`
- Implementar interface Transport
- Criar factory
- Implementar FTP wrapper

Sem alterar comportamento atual.

### Fase 1B — Migração dos consumidores

- `mirror/mirror.go` usa Transport
- `deploy/service.go` usa Transport
- remover imports diretos de FTP fora de transports

---

## 10. Critério de Conclusão

Esta ADR será considerada implementada quando:

- [ ] Não existir `github.com/jlaffaye/ftp` fora de `transports/`
- [ ] Mirror usa Transport interface
- [ ] Deploy usa Transport interface
- [ ] Sync não depende de FTP diretamente
- [ ] FTP está isolado em `internal/transports/ftp`

---

## 11. Impacto

### Positivo

- Elimina acoplamento direto com FTP
- Permite introdução futura de SFTP/S3 sem rewrite
- Melhora testabilidade (mock Transport)
- Centraliza lógica de rede
- Estabiliza boundary da Engine

### Negativo

- Introduz camada adicional de abstração
- Requer refatoração incremental de dois a três módulos

---

## 12. Fora de Escopo

Este ADR NÃO inclui:

- SFTP
- S3
- WebDAV
- Plugins de transporte
- Streaming IO abstractions

Esses pertencem a ADRs futuros.

---

## 13. Princípio de Design

> O FileENIAC não abstrai o futuro. Ele estabiliza o presente.

---

## 14. Conclusão

Este ADR define a fronteira definitiva entre:

- Engine (Sync / Deploy / Mirror)
- Infraestrutura de transporte (FTP)

Ele elimina acoplamentos diretos e estabelece base estável para evolução incremental do sistema.
