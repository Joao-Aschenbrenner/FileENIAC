# SPDX-License-Identifier: MIT
# Build the Go backend as a Tauri sidecar binary.
# Usage: powershell -File scripts/build-sidecar.ps1

$ErrorActionPreference = "Stop"

$projectRoot = if ($PSScriptRoot) { Split-Path -Parent $PSScriptRoot } else { $PWD.Path }
$backendDir  = Join-Path $projectRoot "backend"
$binariesDir = Join-Path $projectRoot "apps\desktop\src-tauri\binaries"

if (-not (Test-Path -LiteralPath $backendDir)) {
    Write-Error "Backend directory not found: $backendDir"
    exit 1
}

if (-not (Test-Path -LiteralPath $binariesDir)) {
    New-Item -ItemType Directory -Path $binariesDir -Force | Out-Null
}

$targetTriple = if ($env:TAURI_ENV_TARGET_TRIPLE) { $env:TAURI_ENV_TARGET_TRIPLE } else { "x86_64-pc-windows-msvc" }
$outputName   = "fileeniac-$targetTriple.exe"
$outputPath   = Join-Path $binariesDir $outputName

$msvcName = "fileeniac-backend-x86_64-pc-windows-msvc.exe"
$gnuName  = "fileeniac-backend-x86_64-pc-windows-gnu.exe"
$msvcPath = Join-Path $binariesDir $msvcName
$gnuPath  = Join-Path $binariesDir $gnuName

Write-Host "Building Go backend for sidecar..."
Set-Location -LiteralPath $backendDir

go build -o $msvcPath .

if (-not (Test-Path -LiteralPath $msvcPath)) {
    Write-Error "Build failed: $msvcPath not found"
    exit 1
}

Copy-Item -LiteralPath $msvcPath -Destination $gnuPath -Force

$size = (Get-Item -LiteralPath $msvcPath).Length
Write-Host "Sidecar binary: $msvcPath ($size bytes)"
Write-Host "Sidecar binary: $gnuPath (copied)"
