# Development and testing

For contributors building Scenaria from source.

## Prerequisites

| Tool | Version |
|------|---------|
| Go | 1.26.5+ |
| Node.js | 22+ |
| Wails CLI | v2 |
| PowerShell | 7+ (Windows scripts) |

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

## Quick dev loop

```bash
go test ./internal/... ./cmd/...
cd frontend && npm install && npm test
wails dev                    # GUI with hot reload
go run ./cmd/scenaria --help
```

## Build artifacts

```powershell
./scripts/build-portable.ps1
./scripts/build-installer.ps1
./scripts/build-release.ps1 -SkipTests
```

Output: `dist/Scenaria-Portable.zip`, `dist/Scenaria-Setup.exe`, `dist/latest.json`.

Tag `v*` triggers `.github/workflows/release.yml` on GitHub Actions.

## Test layers

| Layer | Command |
|-------|---------|
| Go unit | `go test ./internal/... ./cmd/...` |
| Go integration (browser) | `go test -tags=integration ./internal/player/... ./internal/recorder/... ./internal/selector/...` |
| Frontend unit | `cd frontend && npm test` |
| UI E2E (mock Wails) | `cd frontend && npm run test:e2e` |
| Docs screenshots | `cd frontend && npm run docs:screenshots` |
| Desktop smoke (WebView2) | `./scripts/desktop-smoke.ps1` |

CI: `.github/workflows/ci.yml` on `master` — Go **1.26.5**, Node **22**, jobs `test`, `integration`, `desktop-smoke`. Release workflow uses the same toolchain (`.github/workflows/release.yml`).

## Global CLI

```bash
go install ./cmd/scenaria
```

Module path: `github.com/bafgion/scenaria-golang`. See `docs/archive/CLI_GLOBAL_INSTALL.md`.

## Project layout (code)

| Path | Role |
|------|------|
| `cmd/scenaria/` | CLI entry |
| `internal/player/` | Playwright runner |
| `internal/recorder/` | Live recording |
| `internal/gui/` | Wails service layer |
| `internal/wailsapp/` | Wails bindings |
| `frontend/src/` | Svelte IDE |
| `examples/` | Sample features |

## Internal docs

- [ROADMAP](../internal/ROADMAP.md) — engineering phases (RU)
- [QA checklist](../internal/QA-DAILY-USE.md) — manual regression (RU)

## Related

- [Documentation index](../README.md)
