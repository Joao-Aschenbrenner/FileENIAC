param(
    [string]$ExePath = "$env:LOCALAPPDATA\FileENIAC\fileeniac.exe",
    [int]$Iterations = 10,
    [int]$TimeoutSeconds = 30
)

$ErrorActionPreference = "Stop"
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$smoke = Join-Path $scriptDir "smoke-desktop-startup.ps1"

for ($i = 1; $i -le $Iterations; $i++) {
    "stress iteration $i/$Iterations"
    & $smoke -ExePath $ExePath -TimeoutSeconds $TimeoutSeconds
    Start-Sleep -Milliseconds 500
}

$orphans = Get-Process -Name "fileeniac-backend" -ErrorAction SilentlyContinue
if ($orphans) {
    $orphans | Stop-Process -Force -ErrorAction SilentlyContinue
    throw "Found orphan fileeniac-backend process after stress run."
}

"PASS FileENIAC desktop stress startup: iterations=$Iterations"
