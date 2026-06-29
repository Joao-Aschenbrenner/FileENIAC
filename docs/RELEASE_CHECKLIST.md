# Release Checklist — FileENIAC v0.1.3

**Release Date**: 2026-06-29
**Status**: ✅ COMPLETE

---

## Gates

| Gate | Command | Status | Result |
|------|---------|--------|--------|
| Backend Build | `go build ./...` | ✅ PASS | 0 errors |
| Backend Vet | `go vet ./...` | ✅ PASS | 0 errors |
| Backend Test | `go test ./...` | ✅ PASS | 31/31 packages |
| Backend Race | `go test -race ./...` | ✅ PASS | 0 races |
| Frontend Test | `npm run test` | ✅ PASS | 22/22 files |
| Frontend Build | `npm run build` | ✅ PASS | 0 errors |
| Docker Build | `docker build .` | ✅ PASS | Image built |
| Container Scan | `trivy fs` | ✅ PASS | 0 HIGH/CRITICAL |

---

## Artifacts

| Artifact | Location | Size | SHA256 |
|----------|----------|------|--------|
| Windows Installer (NSIS) | `apps/desktop/src-tauri/target/release/bundle/nsis/FileENIAC_0.1.3_x64-setup.exe` | ~5 MB | 8CEC546608DEC4558131968EA8BE88B236D8C8173574DE3830CE21AD6CAFD4F2 |
| Source tarball | GitHub release source archive | — | _(auto-generated)_ |

---

## Documentation

| Document | Status |
|----------|--------|
| CHANGELOG.md | ✅ v0.1.2 section added |
| RELEASE_NOTES.md | ✅ v0.1.2 notes added |
| docs/audits/FULL_CODE_AUDIT_v0.1.0.md | ✅ Published |
| docs/plans/FIX_PLAN_v0.1.0_AUDIT.md | ✅ Published |

---

## Packaging

| Check | Status |
|-------|--------|
| WebView2 bootstrapper bundled | ✅ `embedBootstrapper` |
| NSIS license page shown before install | ✅ `bundle.licenseFile` |

## Security

| Check | Status |
|-------|--------|
| No hardcoded credentials | ✅ Pass |
| No tokens in examples | ✅ Pass |
| SQL injection mitigated | ✅ Pass |
| Container runs as non-root | ✅ Pass |
| npm audit clean | ✅ Pass |

---

## Version Tags

```bash
git tag -a v0.1.2 -m "FileENIAC v0.1.2 - Security & Stability Hotfix"
git push origin v0.1.2
```

---

*Release checklist updated: 2026-06-29*
