param(
    [string]$Version = "",
    [switch]$SkipTests
)

$ErrorActionPreference = "Stop"
$Root = Resolve-Path (Join-Path $PSScriptRoot "..")
Set-Location $Root

function Get-ScenariaVersion {
    $path = Join-Path $Root "internal\version\version.go"
    if (Test-Path $path) {
        $m = Select-String -Path $path -Pattern 'Version\s*=\s*"([^"]+)"' | Select-Object -First 1
        if ($m) { return $m.Matches[0].Groups[1].Value }
    }
    return "0.15.0"
}

if (-not $Version) {
    $Version = Get-ScenariaVersion
}
if ($Version -match '^v') {
    $Version = $Version.Substring(1)
}

Write-Host "==> Scenaria Go portable build v$Version" -ForegroundColor Cyan

if (-not $SkipTests) {
    Write-Host "==> Run tests" -ForegroundColor Cyan
    go test ./...
}

Write-Host "==> Build CLI" -ForegroundColor Cyan
$Dist = Join-Path $Root "dist\Scenaria"
if (Test-Path $Dist) { Remove-Item $Dist -Recurse -Force }
New-Item -ItemType Directory -Path $Dist | Out-Null

$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -ldflags "-s -w" -o (Join-Path $Dist "scenaria.exe") ./cmd/scenaria

Write-Host "==> Build Wails frontend" -ForegroundColor Cyan
Push-Location (Join-Path $Root "frontend")
npm install --no-fund --no-audit
npm run build
Pop-Location

Write-Host "==> Build Wails GUI" -ForegroundColor Cyan
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

Write-Host "==> Install Playwright Chromium" -ForegroundColor Cyan
go run github.com/mxschmitt/playwright-go/cmd/playwright@latest install chromium

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
Scenaria Portable v$Version
=========================

  scenaria-gui.exe   Wails IDE (основной продукт)
  scenaria.exe       CLI
  browsers\          Chromium для Playwright (подхватывается автоматически)
  examples\          примеры сценариев

Запуск IDE: двойной клик Start-GUI.bat или scenaria-gui.exe
CLI:        scenaria.exe run examples --dry-run
Справка:    scenaria.exe help

Экспорт в Python: scenaria.exe export login.feature --format python --output test_login.py
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

$ZipPath = Join-Path $Root "dist\Scenaria-Portable.zip"
if (Test-Path $ZipPath) { Remove-Item $ZipPath -Force }
Compress-Archive -Path $Dist -DestinationPath $ZipPath

$Hash = (Get-FileHash -Path $ZipPath -Algorithm SHA256).Hash.ToLower()
$Manifest = @{
    version = $Version
    published_at = (Get-Date).ToUniversalTime().ToString("o")
    assets = @{
        portable = @{
            name = "Scenaria-Portable.zip"
            size = (Get-Item $ZipPath).Length
            sha256 = $Hash
        }
    }
}
$ManifestPath = Join-Path $Root "dist\latest.json"
$Manifest | ConvertTo-Json -Depth 6 | Set-Content -Path $ManifestPath -Encoding UTF8

Write-Host "Build complete:" -ForegroundColor Green
Write-Host "  $Dist\scenaria.exe"
Write-Host "  $Dist\scenaria-gui.exe"
Write-Host "  $ZipPath"
Write-Host "  $ManifestPath"
