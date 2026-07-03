# Scenaria — User Guide (English)

**Scenaria Go** is the primary product: a Windows desktop IDE and CLI for authoring and running Russian Gherkin test scenarios with Playwright.

> **Note:** Scenario steps use **Russian keywords** (`Когда`, `Тогда`, `Дано`, …). The IDE shell is mostly Russian; this guide explains the UI in English.

## What you get

- **Wails IDE** (`scenaria-gui.exe`) — Monaco editor, project tree, live recorder, run panels, Allure integration
- **CLI** (`scenaria.exe`) — validate, run, record, export, plugins, Vanessa Automation
- **Portable or installer** — bundled Chromium, no separate Node.js required for end users

## Quick start

1. [Install Scenaria](getting-started/install.md) (installer or portable ZIP)
2. Launch **Scenaria** → **Start** tab → **Open example scenarios**
3. Pick `01-pervaya-proverka.feature` → **Run test** (Ctrl+Enter)
4. Follow the [interactive onboarding tour](getting-started/onboarding.md) (**Help → Training…**)

## Guides

### Getting started

- [Installation](getting-started/install.md) — installer, portable, global CLI
- [First project](getting-started/first-project.md) — wizard, init, open folder
- [Onboarding tour](getting-started/onboarding.md) — 8-step guided tour

### Desktop IDE

- [GUI overview](gui/overview.md) — layout, menus, panels, command palette
- [Editor](gui/editor.md) — Monaco, completions, F1 help, preview
- [Recording](gui/recording.md) — live recorder, picker, post-record diff
- [Running tests](gui/running-tests.md) — run dialog, dry-run, batch, cancel
- [Results & reports](gui/results-and-reports.md) — journal, trace, Allure, flaky rerun
- [Settings](gui/settings.md) — record, selectors, editor, UI

### Authoring scenarios

- [Gherkin (RU)](authoring/gherkin.md) — structure, tags, tables, TestClient
- [Selectors](authoring/selectors.md) — strategies, chains, variables
- [Bundled examples](authoring/examples.md) — `examples/` walkthrough

### CLI

- [Command reference](cli/reference.md) — all commands and flags
- [Reports](cli/reports.md) — JUnit, HTML, Allure, trace, video

### Other

- [Migration from Python Scenaria](migration/from-python.md)
- [Development & tests](contributing/development.md)

## In-app help

Press **F1** or **Help → Reference…** for the step catalog (Russian). **Shift+F1** lists keyboard shortcuts.
