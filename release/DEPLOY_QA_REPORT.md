# DEPLOY QA REPORT — RC1

**Data:** 2026-06-13
**Tester:** Automated (CLI environment)
**Versão:** RC1 (v0.2.0)

---

## 1. API Endpoints

| Endpoint | Método | Status | Notas |
|----------|--------|--------|-------|
| `/api/deploy` | POST | ⚙️ | Requer servidor FTPS real |
| `/api/rollback` | POST | ⚙️ | Requer deploy anterior |
| `/api/verify` | POST | ⚙️ | Requer deploy anterior |
| `/api/readiness/deploy` | GET | ✅ | Verifica pré-condições (projeto, servidor) |
| `/api/readiness/sync` | GET | ✅ | Verifica pré-condições (projeto, mirror) |
| `/api/deploys` | GET | ✅ | Histórico de deploys por projeto |
| `/api/syncs` | GET | ✅ | Histórico de sincronizações |

## 2. Readiness Checks (testados via API)

| Cenário | Resultado |
|---------|-----------|
| Projeto sem servidor | ✅ Readiness retorna erros claros |
| Projeto com servidor | ✅ Readiness aprova |
| Sync sem mirror | ✅ Readiness indica necessidade de mirror |

## 3. Arquitetura do Deploy

```
1. POST /api/readiness/deploy  → verifica pré-condições
2. POST /api/deploy             → empacota projeto, faz upload FTPS, registra história
3. POST /api/verify             → verifica deploy via URL configurada
4. POST /api/rollback           → reverte deploy, restaura versão anterior no servidor
5. GET /api/deploys             → histórico completo de deploys
```

## 4. Componentes do Deploy

| Módulo | Status | Descrição |
|--------|--------|-----------|
| `internal/deploy/packer` | ✅ | Empacota projeto (zip/tar) |
| `internal/deploy/ftp` | ✅ | Conexão FTPS com suporte a TLS implícito/explicito |
| `internal/deploy/token` | ✅ | Substituição de tokens em arquivos |
| `internal/deploy/hardening` | ✅ | Hardening de segurança do deploy |
| `internal/deploy/bypass` | ✅ | Bypass de verificação quando necessário |

## 5. Pendente (Requer FTPS Server Real)

- 🖥️ Deploy real para servidor FTPS (Hostinger, HostGator, etc.)
- 🖥️ Verify após deploy (HTTP request para URL)
- 🖥️ Rollback com restauração de versão anterior
- 🖥️ Validação de empacotamento com projeto real
- 🖥️ Upload real de arquivos via FTPS

---

## Checklist Final Deploy QA

- [x] Endpoints registrados e respondendo
- [x] Readiness checks funcionais
- [x] 4 módulos de deploy implementados
- [x] Packer, FTP, Token, Hardening, Bypass operacionais
- [x] Histórico de deploys registrado
- [ ] 🖥️ Deploy real para FTPS (requer servidor real)
- [ ] 🖥️ Verify real (requer URL de verificação)
- [ ] 🖥️ Rollback real (requer deploy anterior)
