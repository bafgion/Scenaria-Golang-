param(
    [string]$PythonRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..\shop-ui-recorder")).Path
)

$ErrorActionPreference = "Stop"
$GoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$Src = Join-Path $PythonRoot "assets\branding"
$Dst = Join-Path $GoRoot "assets\branding"

if (-not (Test-Path $Src)) {
    throw "Python branding folder not found: $Src"
}

Write-Host "Sync branding from $Src" -ForegroundColor Cyan
New-Item -ItemType Directory -Force -Path $Dst | Out-Null
Copy-Item -Path "$Src\*" -Destination $Dst -Force

$Frontend = Join-Path $GoRoot "frontend\src\assets\branding"
New-Item -ItemType Directory -Force -Path $Frontend | Out-Null
Copy-Item -Path (Join-Path $Dst "app-icon-mark.png") -Destination $Frontend -Force
Copy-Item -Path (Join-Path $Dst "app-icon-square.png") -Destination $Frontend -Force

Copy-Item -Path (Join-Path $Dst "app.ico") -Destination (Join-Path $GoRoot "build\windows\icon.ico") -Force
Copy-Item -Path (Join-Path $Dst "app-icon-square.png") -Destination (Join-Path $GoRoot "build\appicon.png") -Force

Write-Host "Done." -ForegroundColor Green
