# Scenaria (Go Migration)

This repository contains the Go implementation of Scenaria with a goal of full
functional parity with the Python MVP (desktop + recorder + runner + export +
plugins).

## Current state

**v0.13** — Wails 2 + Svelte + Monaco IDE merged to `master`. CLI parity with Python v0.12.

- Full step DSL, Playwright runner, recorder, Vanessa, portable CLI + Wails GUI
- Allure reports (`--allure`), failure screenshots; trace/video opt-in (`--trace`, `--video`)
- Legacy **Fyne** GUI (`-tags desktop`) — deprecated
- `go test ./...` passes

See `docs/FUNCTIONAL_PARITY_MATRIX.md` and `docs/ROADMAP.md`.

## Desktop (Wails — recommended)

Requires [Node.js](https://nodejs.org/) and [Wails CLI](https://wails.io/docs/gettingstarted/installation):

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
cd frontend && npm install && npm run build && cd ..
wails dev          # hot-reload (Monaco editor)
wails build        # scenaria-gui.exe
```

Editor: **Monaco** with custom `scenaria-feature` syntax (Gherkin RU, tags, TestClient).

## Legacy desktop (Fyne, deprecated)

```bash
go run -tags desktop ./cmd/scenaria-gui
```

## Project goals

1. Preserve existing user-facing functionality from Python MVP.
2. Keep compatibility for existing scenario and settings files.
3. Provide a globally installable CLI command (`scenaria`).
4. Rebuild desktop functionality in Go after core parity is stable.

## Local development

```bash
go test ./...
go run ./cmd/scenaria --help
```

## Current CLI capabilities

```bash
# validate feature files
go run ./cmd/scenaria validate ./path/to/features

# initialize project scaffold
go run ./cmd/scenaria init .

# check for updates
go run ./cmd/scenaria update --check

# run (default engine: playwright from project.json or auto)
go run ./cmd/scenaria run ./examples --dry-run

# generate JUnit, HTML, and Allure reports
go run ./cmd/scenaria run ./path/to/features --dry-run --junit junit.xml --html report.html --allure ./allure-results

# Playwright trace/video on failure (opt-in)
go run ./cmd/scenaria run ./examples --trace ./traces --video ./videos --allure ./allure-results

# run with Playwright engine
go run ./cmd/scenaria run ./examples/01-pervaya-proverka.feature --engine playwright --install-playwright --headed

# filter by tag, pass variables
go run ./cmd/scenaria run ./examples --tag smoke --var BASE=https://example.com

# desktop GUI (requires CGO + OpenGL)
make gui

# export scenario to JSON / feature / Playwright
go run ./cmd/scenaria export ./path/to/login.feature --output login.json --format json
go run ./cmd/scenaria export ./path/to/login.feature --output login.spec.ts --format ts --base-url https://example.com

# bootstrap a recorded scenario file from CLI
go run ./cmd/scenaria record --output recorded.feature --feature "Логин" --scenario "Успех" --step "открываю \"https://example.com\""

# validate with browser (selectors must be visible)
go run ./cmd/scenaria validate ./examples --browser --base-url https://example.com

# live record from browser
go run ./cmd/scenaria record --live --url https://example.com --output recorded.feature --idle 30

# Vanessa Automation (1C)
go run ./cmd/scenaria va run --project . --dry-run

# portable Windows build (bundles Chromium)
make build-portable
```

## Install CLI as a global command

```bash
go install ./cmd/scenaria
```

With a hosted module path:

```bash
go install <module-path>/cmd/scenaria@latest
```

See details in `docs/CLI_GLOBAL_INSTALL.md`.

## Migration documentation

- `docs/MIGRATION_PLAN.md` — staged migration roadmap.
- `docs/FUNCTIONAL_PARITY_MATRIX.md` — functional parity contract.
