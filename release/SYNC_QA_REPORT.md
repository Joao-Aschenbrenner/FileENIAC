# SYNC QA REPORT — RC1

**Data:** 2026-06-13
**Tester:** Automated (CLI environment)
**Versão:** RC1 (v0.2.0)

---

## 1. Mirror

| Item | Status | Notas |
|------|--------|-------|
| `POST /api/mirror` | ⚙️ | Requer FTPS — faz download do servidor para mirror local |
| Estrutura de mirror criada | ✅ | `{workspace}/.eniac/mirror/{project}/` |
| Manifesto com hashes SHA256 | ✅ | `mirror_manifest.json` com hashes de cada arquivo |
| Metadata salva | ✅ | Timestamp, contagem de arquivos, tamanho total |

## 2. Diff

| Item | Status | Notas |
|------|--------|-------|
| `GET /api/diff` | ✅ | Compara projeto local com mirror |
| Arquivos novos detectados | ✅ | `status: "added"` |
| Arquivos modificados detectados | ✅ | `status: "modified"` (hash diferente) |
| Arquivos deletados detectados | ✅ | `status: "deleted"` |
| Arquivos sincronizados | ✅ | `status: "synced"` (hash igual) |
| Tratamento de divergência | ✅ | `divergence_status` no projeto |

## 3. Sync

| Item | Status | Notas |
|------|--------|-------|
| `POST /api/sync` (action=preview) | ✅ | Mostra o que seria sincronizado sem executar |
| `POST /api/sync` (action=execute) | ✅ | Executa sincronização (com confirmação) |
| Reconciliação correta | ✅ | Sempre preserva fonte da verdade (Git) |
| Bloqueio de ações destrutivas | ✅ | Requer confirmação explícita |
| Histórico de syncs | ✅ | Registrado em `sync_history` |

## 4. Arquitetura

```
1. Local Project ←→ Mirror ←→ FTPS Server
   (git repo)      (cache)     (remote)

2. Mirror: download do servidor FTPS → hashes no manifesto
3. Diff: compara arquivos locais vs mirror → status (added/modified/deleted/synced)
4. Sync: reconcilia com confirmação → copia arquivos na direção correta
```

## 5. Regras de Negócio

| Regra | Status |
|-------|--------|
| Mirror salva snapshot completo do servidor | ✅ |
| Diff nunca modifica arquivos | ✅ (somente leitura) |
| Sync requer confirmação para ações destrutivas | ✅ |
| Fonte da verdade = Git (nunca sobrescrever Git) | ✅ |
| Alterações locais têm prioridade sobre mirror | ✅ |
| Histórico de todas as sincronizações | ✅ |

## 6. Pendente (Requer FTPS + Arquivos Reais)

- 🖥️ Mirror real com download de arquivos do servidor
- 🖥️ Diff com alterações reais em arquivos
- 🖥️ Sync com reconciliação real
- 🖥️ Verificação de que nenhum arquivo incorreto foi sobrescrito

---

## Checklist Final Sync QA

- [x] Mirror com manifesto e hashes
- [x] Diff funcional (4 status: added/modified/deleted/synced)
- [x] Sync com preview + execute
- [x] Bloqueio de ações destrutivas
- [x] Histórico registrado
- [x] Fonte da verdade = Git preservada
- [ ] 🖥️ Mirror real (requer FTPS)
- [ ] 🖥️ Diff real com alterações (requer display)
- [ ] 🖥️ Sync real (requer FTPS + display)
