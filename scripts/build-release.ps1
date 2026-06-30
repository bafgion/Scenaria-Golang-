param(
    [string]$Version = "",
    [switch]$SkipTests
)

$ErrorActionPreference = "Stop"
$Root = Resolve-Path (Join-Path $PSScriptRoot "..")
Set-Location $Root

. (Join-Path $PSScriptRoot "build-dist.ps1")

if (-not $Version) {
    $Version = Get-ScenariaVersion -Root $Root
}

$built = Build-ScenariaDist -Root $Root -Version $Version -SkipTests:$SkipTests
$Dist = $built.DistDir
$Version = $built.Version

Write-Host "==> Portable ZIP" -ForegroundColor Cyan
$ZipPath = Join-Path $Root "dist\Scenaria-Portable.zip"
if (Test-Path $ZipPath) { Remove-Item $ZipPath -Force }
Compress-Archive -Path $Dist -DestinationPath $ZipPath

Write-Host "==> Windows installer" -ForegroundColor Cyan
& (Join-Path $PSScriptRoot "build-installer.ps1") -Version $Version -SkipTests -DistReady

$SetupPath = Join-Path $Root "dist\Scenaria-Setup.exe"
$Assets = @{
    portable = @{
        name = "Scenaria-Portable.zip"
        size = (Get-Item $ZipPath).Length
        sha256 = (Get-FileHash -Path $ZipPath -Algorithm SHA256).Hash.ToLower()
    }
    setup = @{
        name = "Scenaria-Setup.exe"
        size = (Get-Item $SetupPath).Length
        sha256 = (Get-FileHash -Path $SetupPath -Algorithm SHA256).Hash.ToLower()
    }
}
$ManifestPath = Write-ScenariaManifest -Root $Root -Version $Version -Assets $Assets

Write-Host "Release build complete:" -ForegroundColor Green
Write-Host "  $ZipPath"
Write-Host "  $SetupPath"
Write-Host "  $ManifestPath"
