# ADR-005: Sync Engine — Sincronização e Mirror

## Status
APROVADO

## Data
2026-06-10

## Contexto
A maior dor do fluxo atual é a falta de confiança sobre o estado do servidor. Não é seguro puxar arquivos do servidor diretamente para o projeto local, pois podem conter alterações não versionadas, arquivos temporários ou configurações de ambiente.

## Decisão
Criar o Sync Engine com mirror obrigatório.

### Mirror
Estrutura no diretório `.eniac/mirror/`:
```
.eniac/mirror/
  simple-finance/
    public/
    app/
    config/
    ...
```

O mirror é uma cópia exata do servidor, organizada por projeto. Ele NUNCA substitui o projeto local.

### Fluxo de Pull Seguro
```
Servidor FTPS
    ↓ (download)
Mirror (.eniac/mirror/projeto/)
    ↓ (diff)
Diff (arquivos divergentes)
    ↓ (apresenta ao usuário)
Usuário decide o que aplicar
    ↓ (merge seletivo)
Projeto local
    ↓ (commit manual)
Git
```

### Comandos
```
fileeniac sync pull --project <name>    # Atualiza mirror do servidor
fileeniac sync diff --project <name>    # Mostra diferenças mirror vs local
fileeniac sync apply --project <name>   # Aplica arquivos selecionados
```

### Regras
- Proibido: servidor → projeto local diretamente
- Obrigatório: servidor → mirror → diff → usuário decide → projeto
- Mirror nunca é commitado no repositório do projeto
- Diff ignora arquivos em .gitignore do projeto

## Consequências
- Usuário tem controle total sobre o que sincronizar
- Mirror permite operações offline (diff sem conexão)
- Apenas metadados (hashes) são armazenados no banco
