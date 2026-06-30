param(
    [string]$Version = "",
    [switch]$SkipTests,
    [switch]$DistReady
)

$ErrorActionPreference = "Stop"
$Root = Resolve-Path (Join-Path $PSScriptRoot "..")
Set-Location $Root

. (Join-Path $PSScriptRoot "build-dist.ps1")

if (-not $Version) {
    $Version = Get-ScenariaVersion -Root $Root
}
$Version = Normalize-ScenariaVersion $Version

if (-not $DistReady) {
    $built = Build-ScenariaDist -Root $Root -Version $Version -SkipTests:$SkipTests
    $Version = $built.Version
}

$Dist = Join-Path $Root "dist\Scenaria"
if (-not (Test-Path (Join-Path $Dist "scenaria.exe"))) {
    throw "Distribution folder is missing CLI: $Dist"
}

$ISCC = $null
$candidates = @(
    "${env:ProgramFiles(x86)}\Inno Setup 6\ISCC.exe",
    "$env:ProgramFiles\Inno Setup 6\ISCC.exe"
)
foreach ($path in $candidates) {
    if (Test-Path $path) {
        $ISCC = $path
        break
    }
}
if (-not $ISCC) {
    $cmd = Get-Command iscc -ErrorAction SilentlyContinue
    if ($cmd) { $ISCC = $cmd.Source }
}
if (-not $ISCC) {
    throw "Inno Setup 6 not found. Install from https://jrsoftware.org/isinfo.php or: choco install innosetup -y"
}

$IssFile = Join-Path $Root "installer\scenaria-setup.iss"
if (-not (Test-Path $IssFile)) {
    throw "Missing installer script: $IssFile"
}

Write-Host "==> Compile installer (Inno Setup)" -ForegroundColor Cyan
& $ISCC "/DAppVersion=$Version" $IssFile
if ($LASTEXITCODE -ne 0) {
    throw "ISCC failed with exit code $LASTEXITCODE"
}

$SetupPath = Join-Path $Root "dist\Scenaria-Setup.exe"
if (-not (Test-Path $SetupPath)) {
    throw "Installer not produced: $SetupPath"
}

Write-Host "Installer build complete:" -ForegroundColor Green
Write-Host "  $SetupPath"
