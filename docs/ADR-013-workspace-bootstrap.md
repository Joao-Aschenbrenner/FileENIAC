# ADR-013: Workspace Bootstrap

## Status
Aprovado

## Contexto
Usuários precisam configurar uma máquina nova do zero em poucos minutos. O fluxo manual de criar workspace, cadastrar projetos e configurar servidores é lento e propenso a erro.

## Decisão
Criar fluxo Bootstrap que guia o usuário desde a abertura do Desktop até o ambiente pronto.

## Fluxo Completo

```
Abrir Desktop
  ↓
Onboarding (backend check + workspace connect)
  ↓
Bootstrap Page
  ├── Autenticar GitHub (token → vault)
  ├── Selecionar Organização
  ├── Selecionar Repositórios
  ├── Importar (criar projeto + registrar + clonar)
  ├── Associar Servidores (opcional, via Projects page)
  └── Workspace Pronto
```

## Tela de Bootstrap
- Checklist visual com 4 etapas
- Cada etapa mostra status (✓ ou pendente)
- Botão "Configurar" direciona para etapa específica
- Resumo do ambiente no final

## Validação
- Valida clone via `validate.ValidateClone` (verifica .git)
- Valida importação via `validate.ValidateImport`
- Valida associação de servidor via `validate.ValidateAssociation`

## Consequências
- Elimina configuração manual inicial
- Workspace pronto em segundos
- Rastreabilidade total via eventos `github_import`
- Usuário pode pular etapas e configurar manualmente se preferir
