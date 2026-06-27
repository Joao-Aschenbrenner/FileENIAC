# NEXT_SPRINTS.md

# FileENIAC — Roadmap de Execução (Pós Fase 0)

**Status atual**

## Concluído

* ✅ Branding FileENIAC
* ✅ ADR-014 — Transport Layer & Engine Boundaries
* ✅ Fase 0 — Estabilização do ambiente
* ✅ Gate 0.5 — Auditoria pós estabilização
* ✅ Histórico limpo de commits
* ✅ Working tree preservado em stash

---

# Objetivo Geral

A partir deste ponto o foco deixa de ser infraestrutura.

Todo o trabalho passa a ser na **Engine de sincronização**, eliminando acoplamentos e preparando o sistema para se tornar um motor de deploy confiável.

Não serão adicionadas novas funcionalidades até que a arquitetura esteja estabilizada.

---

# Sprint 1 — Transport Layer Foundation

## Objetivo

Criar a base arquitetural definida no ADR-014.

Nenhum comportamento deve mudar.

Nenhum endpoint deve ser alterado.

Nenhuma funcionalidade nova deve surgir.

---

## Sprint 1A

### Escopo

Criar apenas:

```
internal/transports/
    transport.go
```

Contendo:

* interface Transport
* TransportConfig
* FileInfo

Nada além disso.

---

### Não fazer

* factory
* FTP
* mirror
* deploy
* sync

---

### Commit

```
feat(transports): introduce Transport interface
```

---

### Gate

Executar:

```
go build ./...

go test ./...
```

Nenhuma regressão aceita.

---

# Sprint 1B — Factory

## Objetivo

Criar mecanismo de construção dos transports.

---

Criar:

```
internal/transports/factory.go
```

Responsável apenas por:

* registry
* factory
* resolução por protocolo

Sem alterar consumidores.

---

Commit

```
feat(transports): introduce transport factory
```

---

Gate

```
go build ./...

go test ./...
```

---

# Sprint 1C — FTP Transport

## Objetivo

Mover FTP para trás da interface.

Criar:

```
internal/transports/ftp/
```

Implementando Transport.

Nesta etapa FTP continua funcionando exatamente como hoje.

---

Não modificar:

* deploy
* mirror

---

Commit

```
feat(transports/ftp): implement FTP transport
```

---

Gate

```
go test ./internal/transports/...
```

---

# Sprint 1D — Deploy Migration

Objetivo

Substituir ftpClientIface pelo Transport.

Alterar apenas:

```
internal/deploy/
```

Nenhum outro pacote.

---

Commit

```
refactor(deploy): depend on Transport
```

---

Gate

```
go test ./internal/deploy/...
```

---

# Sprint 1E — Mirror Migration

Objetivo

Eliminar import direto de FTP.

Alterar apenas:

```
internal/mirror/
```

Resultado esperado:

```
grep -R "jlaffaye/ftp" backend/internal
```

Retorna somente:

```
internal/transports/ftp
```

---

Commit

```
refactor(mirror): remove direct FTP dependency
```

---

Gate

```
go build ./...

go test ./...
```

---

# Sprint 2 — Engine Validation

Objetivo

Validar que a abstração não alterou comportamento.

Checklist

* upload
* download
* delete
* mkdir
* rename
* list
* stat

Todos funcionando.

---

Adicionar testes usando mocks de Transport.

---

Não adicionar funcionalidades.

---

# Sprint 3 — Core Reliability

Objetivo

Resolver dívidas técnicas críticas.

Itens previstos

* TD-001 Data Race (`api.go`)
* TD-002 Dependência `webui/dist`
* Revisão de concorrência
* Revisão de shutdown
* Revisão de timeout

---

Meta

```
go test -race ./...
```

Sem data races.

---

# Sprint 4 — Test Coverage

Objetivo

Aumentar cobertura dos módulos críticos.

Prioridades

1. deploy
2. transports
3. mirror
4. database
5. heartbeat

Meta

Cobertura mínima de 70% nos módulos centrais.

---

# Sprint 5 — Deploy Engine

Objetivo

Consolidar o pipeline completo.

Fluxo esperado

```
Workspace
    ↓
Scanner
    ↓
Diff
    ↓
Planner
    ↓
Executor
    ↓
Verifier
    ↓
History
    ↓
Rollback
```

Sem duplicidade entre Deploy e Sync.

---

# Sprint 6 — Desktop Integration

Objetivo

Validar integração completa.

Itens

* Tauri
* API
* Vault
* Sessões
* GitHub
* Deploy
* Logs

---

Meta

Aplicação funcionando ponta a ponta.

---

# Sprint 7 — Release Candidate

Objetivo

Preparar versão RC.

Checklist

* Docker
* CI
* Testes
* Auditoria
* Documentação
* CHANGELOG
* README
* ADRs
* Migrações
* Instalador

---

Critério

Nenhuma dívida técnica crítica aberta.

---

# Sprint 8 — Release 1.0

Objetivo

Primeira versão estável.

Entregáveis

* Instalador
* Documentação completa
* Pipeline de Deploy
* Rollback
* Histórico
* Healthcheck
* Vault
* Sessões
* Engine estabilizada

---

# Regras Gerais

## 1

Um commit resolve apenas um problema.

---

## 2

Toda sprint termina com Gate de Qualidade.

---

## 3

Nenhuma sprint inicia com testes quebrados.

---

## 4

Nenhuma funcionalidade nova antes da estabilização da arquitetura.

---

## 5

Toda decisão arquitetural deve possuir ADR correspondente.

---

## 6

Toda dívida técnica relevante deve possuir documento próprio em:

```
docs/technical-debt/
```

---

# Definição de Pronto (Definition of Done)

Uma sprint só é considerada concluída quando:

* Build verde
* Testes verdes
* Gate aprovado
* Documentação atualizada
* Commit isolado
* Sem regressões
* ADRs respeitados
* Working tree limpo
