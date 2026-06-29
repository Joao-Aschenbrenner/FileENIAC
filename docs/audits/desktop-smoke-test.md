# Desktop Smoke Test — FileENIAC v0.1.0

> Test environment: Windows x64 with GUI
> Installer: `FileENIAC_0.1.0_x64-setup.exe`
> Build: Sprint 8 Release Candidate
> Test date: PENDING — requires interactive GUI environment

---

## Prerequisites

- [ ] Windows 10/11 x64
- [ ] ~100 MB disk space
- [ ] No prior FileENIAC installation (or uninstall previous first)
- [ ] Network access for healthcheck validation
- [ ] ACCESSIBLE GUI ENVIRONMENT (not headless)

---

## Test Procedure

### Step 1 — Installation

- [ ] Run `FileENIAC_0.1.0_x64-setup.exe`
- [ ] Verify installer shows correct app name "FileENIAC"
- [ ] Verify version shows "0.1.0"
- [ ] Accept license agreement
- [ ] Choose installation directory (default: `C:\Program Files\FileENIAC`)
- [ ] Verify desktop shortcut checkbox is checked
- [ ] Verify Start Menu shortcut checkbox is checked
- [ ] Click Install
- [ ] Verify no UAC warnings or code signing alerts (if unsigned)
- [ ] Wait for installation to complete
- [ ] Verify "Installation complete" message
- [ ] Click Finish / Launch

**Expected**: App launches after install. Desktop shortcut appears with correct icon.

---

### Step 2 — First Launch

- [ ] App window opens
- [ ] Backend spawns (Tauri IPC connects)
- [ ] Dashboard loads (or onboarding if first run)
- [ ] No crash on startup
- [ ] No console errors in DevTools (if opened with F12)

**Expected**: App starts cleanly, no crashes.

---

### Step 3 — Backend Spawn Validation

- [ ] Open DevTools (F12)
- [ ] Check console for API connection message
- [ ] Verify no `Failed to fetch` or connection errors
- [ ] Health indicator shows backend online

**Expected**: Backend API responds on expected port (8080 or ENIAC_API_PORT).

---

### Step 4 — Token Handshake

- [ ] Check that `/_handshake/token` request completes (if auth enabled)
- [ ] Verify Authorization header is set on subsequent requests
- [ ] No token visible in console logs

**Expected**: Token resolved and used without exposing in logs.

---

### Step 5 — Create/Open Workspace

*If first-run wizard exists:*
- [ ] Complete workspace creation wizard
- [ ] Enter workspace path (existing directory)
- [ ] Verify workspace is saved

*If manual:*
- [ ] Navigate to workspace settings
- [ ] Create new workspace or open existing
- [ ] Verify workspace appears in sidebar

**Expected**: Workspace created/opened without errors.

---

### Step 6 — Health Check

- [ ] Click health check button (or navigate to health dashboard)
- [ ] Wait for health status to load
- [ ] Verify response shows projects_total, servers_total, divergent_total
- [ ] Verify last_events list is populated (if events exist)

**Expected**: `/api/health/check` returns 200 with valid JSON body.

---

### Step 7 — Close Application

- [ ] Click window close button (X)
- [ ] Verify app closes without crash
- [ ] Verify backend process is terminated
- [ ] Verify no orphan processes remain on port 8080

**Expected**: Clean shutdown.

---

### Step 8 — Reopen Application

- [ ] Launch app from desktop shortcut
- [ ] Verify previous workspace is restored (if applicable)
- [ ] Verify state is preserved between sessions

**Expected**: App reopens with previous state intact.

---

### Step 9 — Uninstall

- [ ] Open Windows Add/Remove Programs (or Settings > Apps)
- [ ] Find "FileENIAC" in the list
- [ ] Click Uninstall
- [ ] Verify uninstaller launches
- [ ] Confirm removal
- [ ] Wait for uninstall to complete
- [ ] Verify desktop shortcut is removed
- [ ] Verify Start Menu entry is removed
- [ ] Verify installation directory is removed (or mostly removed)

**Expected**: Clean uninstall, no residual files (excluding user data in AppData if preserved).

---

## Known Limitations

| Limitation | Reason |
|------------|--------|
| Cannot test in headless environment | Tauri requires GUI display |
| SmartScreen warning (if unsigned) | Expected for unsigned installers |
| Session management UI | Backend API not implemented — will fail |

---

## Sign-Off

| Role | Name | Date | Result |
|------|------|------|--------|
| Tester | | | |
| Reviewer | | | |

**Overall Result**: ✅ PASS / ❌ FAIL

**Notes**: