# Manual Test — Auto Update v0.2.1

## Build

| Item | Valor |
|---|---|
| Versão anterior instalada | |
| Versão nova publicada | |
| Installer | |
| SHA-256 | |
| Signature (.sig) | |
| latest.json | |
| Data | |
| Tester | |

## Checklist

| Teste | Resultado | Observação |
|---|---|---|
| App antigo abre normalmente | | |
| Update check roda após startup | | |
| Sem update não mostra modal | | |
| Nova versão mostra modal com versão | | |
| Botão Mais tarde fecha modal | | |
| Botão Atualizar agora inicia download | | |
| Download não trava app | | |
| Instalação conclui sem erro | | |
| App reinicia automaticamente | | |
| Versão nova aparece correta | | |
| Erro de rede mostra mensagem amigável | | |
| Release sem assinatura não instala | | |
| Assinatura inválida bloqueia update | | |

## Fluxo completo

1. Instalar release v0.2.0 (última sem updater)
2. Publicar release v0.2.1 com `.exe`, `.sig`, `latest.json`
3. Abrir app v0.2.0
4. Aguardar 5s para o check automático
5. Confirmar modal "Nova versão 0.2.1 disponível"
6. Clicar "Mais tarde" — modal fecha
7. Reabrir app — modal reaparece
8. Clicar "Atualizar agora"
9. Aguardar download
10. Aguardar instalação + reinício
11. Confirmar que app abriu na versão 0.2.1

## Hotfix — Remoção de Projeto

| Teste | Resultado | Observação |
|---|---|---|
| Remover só do workspace | | |
| Remover com apagar local | | |
| Pasta local protegida não quebra app | | |
| Repo remoto não é apagado | | |
| Erro amigável | | |

## Decisão

- [ ] APROVADO
- [ ] PATCH NECESSÁRIO
- [ ] BLOQUEADO
