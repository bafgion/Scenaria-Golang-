# Scenaria (Go)

**Scenaria Go** — CLI + Wails IDE for Russian Gherkin scenarios, Playwright runner, live recorder, Vanessa Automation, and plugins.

Legacy **Python/Qt Scenaria is discontinued**. Compatibility: `.feature` / `.scenaria` files and optional export to Python (`scenaria export --format python`).

## Documentation

| Language | Guide |
|----------|-------|
| **English** | [docs/en/index.md](docs/en/index.md) |
| **Русский** | [docs/ru/index.md](docs/ru/index.md) |
| **Full index** | [docs/README.md](docs/README.md) |

Current release: **v0.28.0** — [downloads](https://github.com/bafgion/Scenaria-Golang-/releases).

## Quick start (Windows)

1. Download **Scenaria-Setup.exe** or **Scenaria-Portable.zip** from [Releases](https://github.com/bafgion/Scenaria-Golang-/releases)
2. Run `scenaria-gui.exe` → **Start** → **Open example scenarios**
3. Open `01-pervaya-proverka.feature` → **Run test** (Ctrl+Enter)

## Desktop development

Requires [Node.js](https://nodejs.org/) and [Wails CLI](https://wails.io/docs/gettingstarted/installation):

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
cd frontend && npm install && npm run build && cd ..
wails dev
wails build
```

## CLI (essentials)

```bash
go install ./cmd/scenaria
scenaria init .
scenaria validate ./examples
scenaria run ./examples --dry-run
scenaria run ./examples --tag smoke --headed --install-playwright
```

Full reference: [docs/en/cli/reference.md](docs/en/cli/reference.md) | [RU](docs/ru/cli/reference.md).

## Build release artifacts

```powershell
./scripts/build-release.ps1
# dist/Scenaria-Portable.zip, dist/Scenaria-Setup.exe
```

## Contributing

```bash
go test ./internal/... ./cmd/...
cd frontend && npm test && npm run test:e2e
```

See [docs/en/contributing/development.md](docs/en/contributing/development.md).

## Internal

- [Roadmap](docs/internal/ROADMAP.md) (engineering, RU)
- [QA checklist](docs/internal/QA-DAILY-USE.md) (RU)
