# Shared Scenaria distribution build (CLI + GUI + browsers + examples).
# Used by build-portable.ps1, build-installer.ps1, and build-release.ps1.

function Get-ScenariaVersion {
    param([string]$Root)
    $path = Join-Path $Root "internal\version\version.go"
    if (Test-Path $path) {
        $m = Select-String -Path $path -Pattern 'Version\s*=\s*"([^"]+)"' | Select-Object -First 1
        if ($m) { return $m.Matches[0].Groups[1].Value }
    }
    return "0.15.0"
}

function Normalize-ScenariaVersion {
    param([string]$Version)
    if ($Version -match '^v') {
        return $Version.Substring(1)
    }
    return $Version
}

function Build-ScenariaDist {
    param(
        [string]$Root,
        [string]$Version,
        [switch]$SkipTests
    )

    $Version = Normalize-ScenariaVersion $Version
    Write-Host "==> Scenaria distribution build v$Version" -ForegroundColor Cyan

    if (-not $SkipTests) {
        Write-Host "==> Run tests" -ForegroundColor Cyan
        Push-Location $Root
        try {
            go test ./...
        } finally {
            Pop-Location
        }
    }

    $Dist = Join-Path $Root "dist\Scenaria"
    if (Test-Path $Dist) { Remove-Item $Dist -Recurse -Force }
    New-Item -ItemType Directory -Path $Dist | Out-Null

    Write-Host "==> Build CLI" -ForegroundColor Cyan
    Push-Location $Root
    try {
        $env:GOOS = "windows"
        $env:GOARCH = "amd64"
        go build -ldflags "-s -w" -o (Join-Path $Dist "scenaria.exe") ./cmd/scenaria
    } finally {
        Pop-Location
    }

    Write-Host "==> Build Wails frontend" -ForegroundColor Cyan
    Push-Location (Join-Path $Root "frontend")
    try {
        npm install --no-fund --no-audit
        npm run build
    } finally {
        Pop-Location
    }

    Write-Host "==> Build Wails GUI" -ForegroundColor Cyan
    Push-Location $Root
    try {
        $wails = Get-Command wails -ErrorAction SilentlyContinue
        if (-not $wails) {
            Write-Host "Installing wails CLI..." -ForegroundColor Yellow
            go install github.com/wailsapp/wails/v2/cmd/wails@latest
        }
        wails build -clean -platform windows/amd64 -skipbindings
        $WailsBin = Join-Path $Root "build\bin\scenaria-gui.exe"
        if (-not (Test-Path $WailsBin)) {
            throw "Wails build did not produce $WailsBin"
        }
        Copy-Item $WailsBin (Join-Path $Dist "scenaria-gui.exe") -Force
    } finally {
        Pop-Location
    }

    Write-Host "==> Install Playwright Chromium" -ForegroundColor Cyan
    Push-Location $Root
    try {
        go run github.com/mxschmitt/playwright-go/cmd/playwright@latest install chromium
    } finally {
        Pop-Location
    }

    $PlaywrightLocal = Join-Path $env:LOCALAPPDATA "ms-playwright"
    $BrowsersTarget = Join-Path $Dist "browsers"
    if (Test-Path $PlaywrightLocal) {
        Write-Host "==> Copy browsers" -ForegroundColor Cyan
        New-Item -ItemType Directory -Path $BrowsersTarget | Out-Null
        Get-ChildItem $PlaywrightLocal -Directory | Where-Object {
            $_.Name -like "chromium-*" -or $_.Name -like "ffmpeg-*"
        } | ForEach-Object {
            Copy-Item $_.FullName (Join-Path $BrowsersTarget $_.Name) -Recurse -Force
        }
    }

    $ExamplesSource = Join-Path $Root "examples"
    if (Test-Path $ExamplesSource) {
        Copy-Item $ExamplesSource (Join-Path $Dist "examples") -Recurse -Force
    }

    Set-Content -Path (Join-Path $Dist "version.txt") -Value $Version -Encoding UTF8

    $PortableReadme = @"
Scenaria v$Version
==================

  scenaria-gui.exe   Wails IDE (основной продукт)
  scenaria.exe       CLI (команда scenaria в PATH после установки Setup)
  browsers\          Chromium для Playwright
  examples\          примеры сценариев

Запуск IDE: Start-GUI.bat или scenaria-gui.exe
CLI:        scenaria run examples --dry-run
Справка:    scenaria help
"@
    Set-Content -Path (Join-Path $Dist "README-PORTABLE.txt") -Value $PortableReadme -Encoding UTF8

    @'
@echo off
cd /d "%~dp0"
start "" scenaria-gui.exe
'@ | Set-Content -Path (Join-Path $Dist "Start-GUI.bat") -Encoding ASCII

    @'
@echo off
cd /d "%~dp0"
scenaria.exe %*
'@ | Set-Content -Path (Join-Path $Dist "scenaria-cli.bat") -Encoding ASCII

    return @{
        Version = $Version
        DistDir = $Dist
    }
}

function Write-ScenariaManifest {
    param(
        [string]$Root,
        [string]$Version,
        [hashtable]$Assets
    )

    $Manifest = @{
        version = $Version
        published_at = (Get-Date).ToUniversalTime().ToString("o")
        assets = $Assets
    }
    $ManifestPath = Join-Path $Root "dist\latest.json"
    $Manifest | ConvertTo-Json -Depth 6 | Set-Content -Path $ManifestPath -Encoding UTF8
    return $ManifestPath
}
