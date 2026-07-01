param(
    [string]$ExePath = "$env:LOCALAPPDATA\FileENIAC\fileeniac.exe",
    [int]$TimeoutSeconds = 30
)

$ErrorActionPreference = "Stop"
$processNames = @("fileeniac", "FileENIAC", "fileeniac-backend")
$logDir = Join-Path $env:LOCALAPPDATA "com.eniacsystems.fileeniac\logs"
$bootstrapLog = Join-Path $logDir "fileeniac-bootstrap.log"
$backendLog = Join-Path $logDir "backend.log"

foreach ($name in $processNames) {
    Get-Process -Name $name -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
}

if (!(Test-Path -LiteralPath $ExePath)) {
    throw "FileENIAC executable not found: $ExePath"
}

if (Test-Path -LiteralPath $bootstrapLog) { Remove-Item -LiteralPath $bootstrapLog -Force }

$app = Start-Process -FilePath $ExePath -PassThru
$deadline = (Get-Date).AddSeconds($TimeoutSeconds)
$port = $null

try {
    while ((Get-Date) -lt $deadline) {
        if (Test-Path -LiteralPath $bootstrapLog) {
            $ready = Select-String -LiteralPath $bootstrapLog -Pattern "backend_ready port=(\d+)" -ErrorAction SilentlyContinue | Select-Object -Last 1
            if ($ready -and $ready.Matches.Count -gt 0) {
                $port = $ready.Matches[0].Groups[1].Value
                break
            }
        }
        Start-Sleep -Milliseconds 500
    }

    if (!$port) {
        throw "FileENIAC did not report backend_ready within $TimeoutSeconds seconds. Bootstrap log: $bootstrapLog"
    }

    $health = Invoke-RestMethod -Uri "http://127.0.0.1:$port/api/health" -TimeoutSec 5
    if ($health.status -ne "ok") {
        throw "Unexpected health response on port ${port}: $($health | ConvertTo-Json -Compress)"
    }

    $backend = Get-Process -Name "fileeniac-backend" -ErrorAction SilentlyContinue
    if (!$backend) {
        throw "Backend process is not running after startup. Backend log: $backendLog"
    }

    "PASS FileENIAC desktop smoke startup: port=$port pid=$($app.Id) backendPid=$($backend.Id)"
} finally {
    foreach ($name in $processNames) {
        Get-Process -Name $name -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
    }
}
