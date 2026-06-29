# SPDX-License-Identifier: MIT
# FileENIAC NSIS Installer Legal Pages Patch Script
# Run this after building the Tauri desktop app to inject the native NSIS license page.

param(
    [string]$InstallerNsi = "apps\desktop\src-tauri\target\release\nsis\x64\installer.nsi",
    [string]$LicenseSource = "apps\desktop\src-tauri\installer-license.txt"
)

$ErrorActionPreference = "Stop"

if (-not (Test-Path $InstallerNsi)) {
    Write-Host "[ERROR] Installer NSI not found at: $InstallerNsi"
    Write-Host "        Build the desktop app first: cd apps/desktop && npm run tauri -- build"
    exit 1
}

if (-not (Test-Path $LicenseSource)) {
    Write-Host "[ERROR] License file not found at: $LicenseSource"
    exit 1
}

Write-Host "[INFO] Injecting native NSIS license page..."

$installerDir = Split-Path -Parent $InstallerNsi
$licenseDest = Join-Path $installerDir "installer-license.txt"
Copy-Item -LiteralPath $LicenseSource -Destination $licenseDest -Force
Write-Host "[INFO] Copied license to: $licenseDest"

$content = Get-Content $InstallerNsi -Raw

if ($content.Contains('!define LICENSE "installer-license.txt"')) {
    Write-Host "[INFO] License already injected. Skipping."
    exit 0
}

if (-not $content.Contains('!define LICENSE ""')) {
    Write-Host "[ERROR] Could not find empty LICENSE define in installer script."
    exit 1
}

$content = $content.Replace('!define LICENSE ""', '!define LICENSE "installer-license.txt"')
$content | Set-Content -Path $InstallerNsi -NoNewline -Encoding UTF8

Write-Host "[OK] License page injected. Rebuild with: makensis '$InstallerNsi'"
