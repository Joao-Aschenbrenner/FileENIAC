# FileENIAC — Updater Release Workflow

## Objetivo

Publicar uma nova versão do FileENIAC com suporte a auto-update via Tauri.

## Pré-requisitos

- working tree clean
- gates verdes (tsc, tests, build, go vet/test/build, cargo fmt/check, tauri build)
- chave privada de assinatura disponível localmente
- senha da chave fora do Git (apenas em `.env` local)
- GitHub CLI autenticado, se usado
- versão nova definida

## Segurança

Nunca commitar:

- `.env`
- chave privada (`*.key`)
- senha da chave
- tokens de acesso
- credenciais FTPS

Somente a chave pública pode ficar versionada no app (em `tauri.conf.json`).

## Artefatos do updater

O Tauri updater para Windows exige:

| Artefato | Descrição | Obrigatório |
|---|---|---|
| `FileENIAC_X.Y.Z_x64-setup.exe` | Instalador NSIS | Sim |
| `FileENIAC_X.Y.Z_x64-setup.exe.sig` | Assinatura do instalador | Sim |

> `.blockmap` **não** é usado pelo Tauri (é específico do Electron).

## Passo 1 — Bump de versão

Atualizar os seguintes arquivos:

- `apps/desktop/package.json`
- `apps/desktop/src-tauri/tauri.conf.json` (campo `version`)
- `apps/desktop/src-tauri/Cargo.toml` (campo `version`)
- `CHANGELOG.md` (mover de Unreleased para a nova versão)

Para release candidata: usar `X.Y.Z-rc.N` (ex: `0.2.1-rc.1`)
Para release final: usar `X.Y.Z` (ex: `0.2.1`)

## Passo 2 — Gates

```bash
cd apps/desktop
npx tsc --noEmit
npm run test
npm run build

cd ../../backend
go vet ./...
go test ./...
go build ./...

cd ../apps/desktop/src-tauri
cargo fmt --check
cargo check

cd ..
npm run tauri build
```

## Passo 3 — Gerar instalador

O instalador é gerado em:

```
apps/desktop/src-tauri/target/release/bundle/nsis/FileENIAC_X.Y.Z_x64-setup.exe
```

## Passo 4 — Assinar artefato do updater

> A chave privada foi gerada com `npx tauri signer generate --password "<senha>" --write-keys "~/.tauri/fileeniac.key"`.

Após o build, assinar o instalador:

```bash
npx tauri signer sign \
  --private-key ~/.tauri/fileeniac.key \
  --password "<senha>" \
  --file apps/desktop/src-tauri/target/release/bundle/nsis/FileENIAC_X.Y.Z_x64-setup.exe
```

Isso gera o arquivo `.sig` ao lado do instalador.

**Importante:** o `.sig` precisa estar presente na mesma release do GitHub que o instalador. O Tauri updater baixa ambos e verifica a assinatura.

## Passo 5 — Criar latest.json

Criar um arquivo `latest.json` com o metadata do updater:

```json
{
  "version": "X.Y.Z",
  "notes": "Descrição das mudanças desta versão",
  "pub_date": "2026-07-07T00:00:00Z",
  "platforms": {
    "windows-x86_64": {
      "signature": "<conteúdo do arquivo .sig>",
      "url": "https://github.com/Joao-Aschenbrenner/FileENIAC/releases/download/vX.Y.Z/FileENIAC_X.Y.Z_x64-setup.exe"
    }
  }
}
```

O campo `signature` deve conter o conteúdo completo do `.sig` (uma linha Base64).

## Passo 6 — Gerar checksum

```powershell
Get-FileHash .\apps\desktop\src-tauri\target\release\bundle\nsis\FileENIAC_X.Y.Z_x64-setup.exe -Algorithm SHA256
```

## Passo 7 — Criar release no GitHub

```bash
gh release create vX.Y.Z \
  --title "FileENIAC vX.Y.Z" \
  --notes "release notes aqui"
```

Para RC, adicionar `--prerelease`.

Upload dos assets:

```bash
gh release upload vX.Y.Z \
  apps/desktop/src-tauri/target/release/bundle/nsis/FileENIAC_X.Y.Z_x64-setup.exe \
  apps/desktop/src-tauri/target/release/bundle/nsis/FileENIAC_X.Y.Z_x64-setup.exe.sig \
  latest.json
```

RC: marcar como `--prerelease`, não marcar como Latest.
Release final: marcar como Latest **somente após teste manual aprovado**.

## Passo 8 — Teste de update

1. Instalar versão anterior (v0.2.0)
2. Publicar versão nova (vX.Y.Z)
3. Abrir app antigo
4. Aguardar check automático (~5s)
5. Confirmar modal de update
6. Clicar "Atualizar agora"
7. Aguardar download e instalação
8. Confirmar que app reinicia
9. Confirmar versão nova no app

## Rollback

Se update quebrar:

1. Despublicar release problemática (gh release delete)
2. Remover tag se necessário
3. Publicar hotfix seguindo este workflow
4. Documentar incidente

## Troubleshooting

**Update check falha silenciosamente:** Normal se não houver internet ou se a release não tiver `.sig`. O app não trava — apenas não mostra modal.

**Assinatura inválida:** O Tauri bloqueia a instalação e o erro aparece no modal.

**Versão não encontrada:** Verificar se o `latest.json` e o instalador foram upados na mesma release e se as URLs estão corretas.
