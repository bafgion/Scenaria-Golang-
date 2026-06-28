param(
    [string]$OutFile = "coverage.out",
    [switch]$Html
)

$ErrorActionPreference = "Stop"
$Root = Resolve-Path (Join-Path $PSScriptRoot "..")
Set-Location $Root

Write-Host "==> go test -coverpkg=./... -coverprofile=$OutFile" -ForegroundColor Cyan
go test -coverpkg=./... -coverprofile=$OutFile -covermode=atomic ./...
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

go tool cover -func=$OutFile | Select-String "total:"

if ($Html) {
    $HtmlPath = [System.IO.Path]::ChangeExtension($OutFile, ".html")
    go tool cover -html=$OutFile -o $HtmlPath
    Write-Host "HTML report: $HtmlPath" -ForegroundColor Green
}
