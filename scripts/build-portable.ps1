param(
    [string]$Version = "",
    [switch]$SkipTests
)

$ErrorActionPreference = "Stop"
$Root = Resolve-Path (Join-Path $PSScriptRoot "..")
Set-Location $Root

if (-not $Version) {
    $Version = (go run ./internal/version 2>$null)
    if (-not $Version) {
        $Version = "0.1.0-go"
    }
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
go build -tags desktop -ldflags "-s -w" -o (Join-Path $Dist "scenaria-gui.exe") ./cmd/scenaria-gui

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
Write-Host "  $ZipPath"
Write-Host "  $ManifestPath"
