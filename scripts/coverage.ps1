param(
    [string]$OutFile = "coverage.out",
    [switch]$Html
)

$ErrorActionPreference = "Stop"
$Root = Resolve-Path (Join-Path $PSScriptRoot "..")
Set-Location $Root

Write-Host "==> go test -coverprofile=$OutFile" -ForegroundColor Cyan
go test ./... -coverprofile=$OutFile -covermode=atomic
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

go tool cover -func=$OutFile | Select-Object -Last 1

if ($Html) {
    $HtmlPath = [System.IO.Path]::ChangeExtension($OutFile, ".html")
    go tool cover -html=$OutFile -o $HtmlPath
    Write-Host "HTML report: $HtmlPath" -ForegroundColor Green
}
