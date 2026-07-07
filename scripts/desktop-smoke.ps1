param(
    [string]$ExePath = "",
    [int]$DebugPort = 9333,
    [int]$StartupTimeoutSec = 90,
    [int]$AliveSec = 8,
    [switch]$SkipUI,
    [switch]$SkipBuild
)

$ErrorActionPreference = "Stop"
$Root = Resolve-Path (Join-Path $PSScriptRoot "..")

if (-not $ExePath) {
    $ExePath = Join-Path $Root "build\bin\scenaria-gui.exe"
}

function Test-CdpReady([int]$Port) {
    try {
        $r = Invoke-WebRequest -Uri "http://127.0.0.1:$Port/json/version" -UseBasicParsing -TimeoutSec 3
        return $r.StatusCode -eq 200
    } catch {
        return $false
    }
}

function Wait-CdpReady([int]$Port, [int]$TimeoutSec) {
    $deadline = (Get-Date).AddSeconds($TimeoutSec)
    while ((Get-Date) -lt $deadline) {
        if (Test-CdpReady $Port) { return $true }
        Start-Sleep -Milliseconds 500
    }
    return $false
}

if (-not $SkipBuild -and -not (Test-Path $ExePath)) {
    Write-Host "==> Building scenaria-gui (frontend + wails)..." -ForegroundColor Cyan
    Push-Location $Root
    try {
        Push-Location (Join-Path $Root "frontend")
        npm install --no-audit --no-fund 2>&1 | Out-Host
        npm run build 2>&1 | Out-Host
        Pop-Location
        $wails = Get-Command wails -ErrorAction SilentlyContinue
        if (-not $wails) {
            go install github.com/wailsapp/wails/v2/cmd/wails@latest
        }
        wails build -platform windows/amd64 -skipbindings 2>&1 | Out-Host
    } finally {
        Pop-Location
    }
}

if (-not (Test-Path $ExePath)) {
    Write-Error "GUI exe not found: $ExePath. Build with: wails build -platform windows/amd64"
    exit 1
}

if (-not (Test-Path (Join-Path $Root "examples"))) {
    Write-Error "examples/ folder not found at repo root (needed for «Открыть примеры»)"
    exit 1
}

Write-Host "==> Desktop smoke: $ExePath (CDP :$DebugPort)" -ForegroundColor Cyan

$env:SCENARIA_DESKTOP_SMOKE = "1"
$env:SCENARIA_DESKTOP_SMOKE_PORT = "$DebugPort"
$smokeData = Join-Path $env:TEMP "scenaria-desktop-smoke"
if (Test-Path $smokeData) {
    Remove-Item -Recurse -Force $smokeData -ErrorAction SilentlyContinue
}
New-Item -ItemType Directory -Path $smokeData -Force | Out-Null
$env:SCENARIA_APP_DATA = $smokeData

$proc = Start-Process `
    -FilePath $ExePath `
    -WorkingDirectory $Root `
    -PassThru

if (-not $proc) {
    Write-Error "Failed to start scenaria-gui"
    exit 1
}

$exitCode = 0
try {
    if (-not (Wait-CdpReady $DebugPort $StartupTimeoutSec)) {
        if ($proc.HasExited) {
            Write-Error "scenaria-gui exited before CDP ready (code $($proc.ExitCode))"
        } else {
            Write-Error "CDP not ready on port $DebugPort within ${StartupTimeoutSec}s"
        }
        exit 1
    }
    Write-Host "CDP ready on :$DebugPort" -ForegroundColor Green

    $pageDeadline = (Get-Date).AddSeconds(60)
    while ((Get-Date) -lt $pageDeadline) {
        try {
            $targets = Invoke-RestMethod -Uri "http://127.0.0.1:$DebugPort/json/list" -TimeoutSec 3
            $appPage = $targets | Where-Object { $_.type -eq 'page' -and $_.url -like '*wails.localhost*' }
            if ($appPage) { break }
        } catch {}
        Start-Sleep -Milliseconds 400
    }

    Write-Host "App page ready, waiting for splash…" -ForegroundColor Green
    Start-Sleep -Seconds 12

    $deadline = (Get-Date).AddSeconds($AliveSec)
    while ((Get-Date) -lt $deadline) {
        if ($proc.HasExited) {
            Write-Error "scenaria-gui exited early with code $($proc.ExitCode)"
            exit 1
        }
        Start-Sleep -Milliseconds 400
    }
    Write-Host "Process alive for ${AliveSec}s" -ForegroundColor Green

    if (-not $SkipUI) {
        Write-Host "==> Playwright desktop UI smoke..." -ForegroundColor Cyan
        Push-Location (Join-Path $Root "frontend")
        try {
            $env:DESKTOP_CDP_URL = "http://127.0.0.1:$DebugPort"
            $prevEap = $ErrorActionPreference
            $ErrorActionPreference = "Continue"
            npx playwright install chromium 2>&1 | Out-Host
            npm run test:e2e:desktop 2>&1 | Out-Host
            $ErrorActionPreference = $prevEap
            if ($LASTEXITCODE -ne 0) {
                $exitCode = $LASTEXITCODE
            }
        } finally {
            Pop-Location
        }
    }
} finally {
    if (-not $proc.HasExited) {
        Stop-Process -Id $proc.Id -Force -ErrorAction SilentlyContinue
    }
    Remove-Item Env:SCENARIA_DESKTOP_SMOKE -ErrorAction SilentlyContinue
    Remove-Item Env:SCENARIA_DESKTOP_SMOKE_PORT -ErrorAction SilentlyContinue
    Remove-Item Env:SCENARIA_APP_DATA -ErrorAction SilentlyContinue
}

if ($exitCode -ne 0) {
    exit $exitCode
}

Write-Host "desktop smoke OK" -ForegroundColor Green
