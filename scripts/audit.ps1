$ErrorActionPreference = "Stop"
$Root = Resolve-Path (Join-Path $PSScriptRoot "..")
Set-Location $Root

Write-Host "==> govulncheck ./..." -ForegroundColor Cyan
go install golang.org/x/vuln/cmd/govulncheck@latest
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

govulncheck ./...
exit $LASTEXITCODE
