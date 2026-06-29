# SPDX-License-Identifier: MIT
# FileENIAC NSIS Installer Legal Pages Patch Script
# Run this after building the Tauri desktop app to inject legal acceptance pages
# Usage: .\scripts\patch-installer-nsis.ps1

param(
    [string]$InstallerNsi = "apps\desktop\src-tauri\target\release\nsis\x64\installer.nsi",
    [string]$OutputNsi = "apps\desktop\src-tauri\target\release\nsis\x64\installer-patched.nsi"
)

$ErrorActionPreference = "Stop"

if (-not (Test-Path $InstallerNsi)) {
    Write-Host "[ERROR] Installer NSI not found at: $InstallerNsi"
    Write-Host "       Build the desktop app first: cd apps/desktop && npm run tauri -- build"
    exit 1
}

Write-Host "[INFO] Patching NSIS installer with legal acceptance pages..."

$content = Get-Content $InstallerNsi -Raw

# Custom page function definitions - insert before the Install section
$customFunctions = @'
; -------------------------------------------------------
; Custom Legal Pages (added by patch-installer-nsis.ps1)
; -------------------------------------------------------
Var LegalTermsAccepted
Var LegalPrivacyAccepted

; Terms of Use page leave callback
Function legalTermsPageLeave
  `${If} $LegalTermsAccepted == 0
    MessageBox MB_ICONEXCLAMATION|MB_OK "You must accept the Terms of Use to continue the installation."
    Abort
  `${EndIf}
FunctionEnd

; Privacy Policy page leave callback
Function legalPrivacyPageLeave
  `${If} $LegalPrivacyAccepted == 0
    MessageBox MB_ICONEXCLAMATION|MB_OK "You must accept the Privacy Policy to continue the installation."
    Abort
  `${EndIf}
FunctionEnd

; Terms of Use page show callback
Function legalTermsPageShow
  nsDialogs::Create 1018
  Pop `$R0
  `${IfThen} `$R0 == 0 `${|} Abort `${|}
  `${NSD_CreateLabel} 0 0 100% 24u "Terms of Use — FileENIAC"
  Pop `$R0
  SendMessage `$R0 `${WM_SETFONT} `$1 0
  `${NSD_CreateMultiLineBox} 0 30u 100% 200u "FileENIAC is provided 'AS IS' under the MIT License.$$\n$$\nYou may use, copy, modify, and distribute the software under the terms of the MIT License.$$\n$$\nFor full Terms of Use, see docs/legal/TERMS_OF_USE.md in the installation directory or visit:$$\nhttps://github.com/Joao-Aschenbrenner/FileENIAC/blob/main/docs/legal/TERMS_OF_USE.md"
  Pop `$R0
  `${NSD_CreateCheckbox} 0 240u 100% 10u "I have read and agree to the Terms of Use"
  Pop `$R1
  `${NSD_OnClick} `$R1 legalTermsOnClick
  GetFunctionAddress `$0 legalTermsPageLeave
  nsDialogs::Show
FunctionEnd

Function legalTermsOnClick
  Pop `$R0
  `${NSD_GetState} `$R0 `$R1
  `${If} `$R1 == `${BST_CHECKED}
    StrCpy `$LegalTermsAccepted 1
  `${Else}
    StrCpy `$LegalTermsAccepted 0
  `${EndIf}
FunctionEnd

; Privacy Policy page show callback
Function legalPrivacyPageShow
  nsDialogs::Create 1018
  Pop `$R0
  `${IfThen} `$R0 == 0 `${|} Abort `${|}
  `${NSD_CreateLabel} 0 0 100% 24u "Privacy Policy — FileENIAC"
  Pop `$R0
  SendMessage `$R0 `${WM_SETFONT} `$1 0
  `${NSD_CreateMultiLineBox} 0 30u 100% 200u "FileENIAC stores all data locally on your device.$$\n$$\nNo telemetry, analytics, or personal data is sent to the developer.$$\n$$\nYour credentials (GitHub tokens, FTPS passwords) are encrypted in a Vault if FILEENIAC_VAULT_PASSWORD is set.$$$\n$$\nFor full Privacy Policy, see docs/legal/PRIVACY_POLICY.md in the installation directory or visit:$$\nhttps://github.com/Joao-Aschenbrenner/FileENIAC/blob/main/docs/legal/PRIVACY_POLICY.md"
  Pop `$R0
  `${NSD_CreateCheckbox} 0 240u 100% 10u "I have read and agree to the Privacy Policy"
  Pop `$R1
  `${NSD_OnClick} `$R1 legalPrivacyOnClick
  GetFunctionAddress `$0 legalPrivacyPageLeave
  nsDialogs::Show
FunctionEnd

Function legalPrivacyOnClick
  Pop `$R0
  `${NSD_GetState} `$R0 `$R1
  `${If} `$R1 == `${BST_CHECKED}
    StrCpy `$LegalPrivacyAccepted 1
  `${Else}
    StrCpy `$LegalPrivacyAccepted 0
  `${EndIf}
FunctionEnd

; -------------------------------------------------------
'@

# Page declarations - insert after MUI_PAGE_DIRECTORY
$termsPageDeclaration = @'

; Custom page: Terms of Use
!define MUI_PAGE_CUSTOMFUNCTION_PRE SkipIfPassive
!insertmacro MUI_PAGE_CUSTOM legalTermsPageShow
'@

$privacyPageDeclaration = @'

; Custom page: Privacy Policy
!define MUI_PAGE_CUSTOMFUNCTION_PRE SkipIfPassive
!insertmacro MUI_PAGE_CUSTOM legalPrivacyPageShow
'@

# Check if already patched
if ($content -match "legalTermsPageShow") {
    Write-Host "[INFO] Installer already patched with legal pages. Skipping."
    exit 0
}

# Insert function definitions before "Section Install"
if ($content -match '(Section Install[\s\S]*?(?=^Section|\z))') {
    Write-Host "[INFO] Inserting legal page functions before Section Install..."
    $content = $content -replace '(Section Install[\s\S]*?(?=^Section|\z))', ("$customFunctions`n`n" + '$1')
}

# Find MUI_PAGE_DIRECTORY and insert Terms page after it
if ($content -match '(!insertmacro MUI_PAGE_DIRECTORY\r?\n)') {
    Write-Host "[INFO] Inserting Terms of Use page after MUI_PAGE_DIRECTORY..."
    $content = $content -replace '(!insertmacro MUI_PAGE_DIRECTORY\r?\n)', ("`$1" + "$termsPageDeclaration`n")
}

# Find MUI_PAGE_DIRECTORY again (now includes our Terms page) and insert Privacy page after it
# We need to find the Terms page we just inserted and add Privacy after it
if ($content -match '(!insertmacro MUI_PAGE_CUSTOM legalTermsPageShow\r?\n)') {
    Write-Host "[INFO] Inserting Privacy Policy page after Terms of Use page..."
    $content = $content -replace '(!insertmacro MUI_PAGE_CUSTOM legalTermsPageShow\r?\n)', ("`$1" + "$privacyPageDeclaration`n")
}

# Also initialize the checkbox variables in .onInit
if ($content -match '(\$\{If\} \$\{Silent\}[\s\S]*?Abort\r?\n\r?\nFunctionEnd)') {
    if ($content -notmatch 'StrCpy.*LegalTermsAccepted') {
        Write-Host "[INFO] Initializing legal checkbox variables in .onInit..."
        $initVars = "  StrCpy `$LegalTermsAccepted 0`n  StrCpy `$LegalPrivacyAccepted 0`n"
        $content = $content -replace '(\$\{If\} \$\{Silent\}[\s\S]*?Abort\r?\n\r?\nFunctionEnd)', ("`$1`n$initVars")
    }
}

# Add nsDialogs header if not present
if ($content -notmatch '!include nsDialogs\.nsh') {
    Write-Host "[INFO] Adding nsDialogs include..."
    $content = $content -replace '(!include MUI2\.nsh)', ("`$1`n!include nsDialogs.nsh")
}

# Save patched file
$content | Set-Content -Path $OutputNsi -NoNewline -Encoding UTF8

Write-Host "[OK] Patched installer written to: $OutputNsi"
Write-Host ""
Write-Host "Next steps:"
Write-Host "  1. Backup original: Copy-Item '$InstallerNsi' '$InstallerNsi.bak'"
Write-Host "  2. Replace: Copy-Item '$OutputNsi' '$InstallerNsi' -Force"
Write-Host "  3. Rebuild NSIS: makensis '$InstallerNsi'"
Write-Host ""
Write-Host "Or use the Makefile: make patch-installer && make build-installer"