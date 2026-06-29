param(
    [string]$ExePath = "",
    [int]$TimeoutSec = 12
)

$ErrorActionPreference = "Stop"
$Root = Resolve-Path (Join-Path $PSScriptRoot "..")

if (-not $ExePath) {
    $ExePath = Join-Path $Root "build\bin\scenaria-gui.exe"
}

if (-not (Test-Path $ExePath)) {
    Write-Error "GUI exe not found: $ExePath. Build with: wails build -platform windows/amd64"
    exit 1
}

Write-Host "==> Desktop smoke: $ExePath" -ForegroundColor Cyan
$proc = Start-Process -FilePath $ExePath -PassThru
if (-not $proc) {
    Write-Error "Failed to start scenaria-gui"
    exit 1
}

try {
    $deadline = (Get-Date).AddSeconds($TimeoutSec)
    while ((Get-Date) -lt $deadline) {
        if ($proc.HasExited) {
            Write-Error "scenaria-gui exited early with code $($proc.ExitCode)"
            exit 1
        }
        Start-Sleep -Milliseconds 400
    }
    Write-Host "desktop smoke OK: process alive for ${TimeoutSec}s" -ForegroundColor Green
}
finally {
    if (-not $proc.HasExited) {
        Stop-Process -Id $proc.Id -Force -ErrorAction SilentlyContinue
    }
}
