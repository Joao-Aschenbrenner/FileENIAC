# Manual Test — FileENIAC v0.2.0 RC

## Build

- Version: 0.2.0-rc.1
- Commit: 3e92c4f
- Installer: FileENIAC_0.2.0_x64-setup.exe
- SHA-256: (preencher)
- Date: 2026-07-02
- Tester: (preencher)
- OS: Windows 10/11

## Gates

| Gate | Resultado |
|------|-----------|
| npx tsc --noEmit | PASS |
| npm run test | PASS (171/171) |
| npm run build | PASS |
| go vet ./... | PASS |
| go test ./... | PASS |
| go build ./... | PASS |
| cargo fmt --check | PASS |
| cargo check | PASS |
| npm run tauri build | (preencher) |

## Teste Instalado

| Area | Resultado | Observacao |
|------|-----------|------------|
| Instalacao | | |
| Abertura | | |
| Sidecar backend | | |
| Workspace | | |
| Setup simplificado | | |
| Projetos vazio | | |
| Projetos com dados | | |
| Adicionar Repositorios | | |
| GitHub organizacoes | | |
| GitHub pessoais | | |
| Voltar interno | | |
| Importacao | | |
| Retorno para Projetos | | |
| Filtro de .github | | |
| Remocao segura | | |
| Fechamento | | |
| Reabertura | | |

## Bugs Encontrados

| ID | Severidade | Descricao | Como reproduzir | Bloqueia? |
|----|------------|-----------|-----------------|-----------|

## Decisao

- [ ] KEEP RC
- [ ] PATCH REQUIRED
- [ ] BLOCK RC

## Observacoes
