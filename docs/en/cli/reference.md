# CLI reference

```bash
scenaria <command> [flags] [paths...]
scenaria version
```

## Commands

| Command | Description |
|---------|-------------|
| `run` | Execute scenarios |
| `validate` | Validate `.feature` files |
| `export` | Export to JSON / feature / TypeScript / Python |
| `import-json` | Import JSON export to `.feature` |
| `record` | Create or live-record scenarios |
| `init` | Create `.scenaria/` scaffold |
| `update` | Check for application updates (`--check`) |
| `plugins` | `list`, `install`, `uninstall` |
| `va` | Vanessa Automation runner (`va run`, …) |
| `version` | Print version |
| `help` | Show help |

## `run`

```bash
scenaria run ./features [flags]
```

| Flag | Description |
|------|-------------|
| `--dry-run` | Parse steps without browser |
| `--tag NAME` | Filter by tag (e.g. `smoke` for `@smoke`) |
| `--scenario NAME` | Run single scenario by title |
| `--var KEY=VALUE` | Scenario variable (repeatable) |
| `--test-client NAME` | Override TestClient for all cases |
| `--workers N` | Parallel workers (default 1) |
| `--slow-mo MS` | Playwright slow motion |
| `--headed` | Show browser window |
| `--browser NAME` | `chromium`, `firefox`, `webkit` |
| `--base-url URL` | Base URL for relative paths |
| `--install-playwright` | Auto-install browser if missing |
| `--engine NAME` | `playwright` or `stub` |
| `--junit PATH` | JUnit XML report |
| `--html PATH` | HTML report |
| `--allure DIR` | Allure results directory |
| `--trace DIR` | Trace on failure |
| `--video DIR` | Video on failure |
| `--summary-json PATH` | Machine-readable summary |
| `--start-step N` | 0-based first step index |
| `--end-step N` | 0-based last step index |
| `--nav-wait-until` | `load`, `domcontentloaded`, `networkidle`, `commit` |
| `--continue-on-fail` | Do not stop batch on first failure |

Examples:

```bash
scenaria run ./examples --dry-run
scenaria run ./features --tag smoke --headed --install-playwright
scenaria run ./login.feature --trace ./traces --allure ./allure-results
```

## `validate`

```bash
scenaria validate ./features [flags]
```

| Flag | Description |
|------|-------------|
| `--json PATH` | Write results JSON |
| `--browser` | Check selectors in browser (default) |
| `--no-browser` | Syntax only |
| `--headless` | Browser without UI |
| `--base-url URL` | For relative navigation |

## `export`

```bash
scenaria export INPUT.feature --output OUT --format json|feature|ts|python [--base-url URL]
```

## `import-json`

```bash
scenaria import-json doc.json --output doc.feature [--force]
```

## `record`

```bash
# Baseline from CLI steps
scenaria record --output new.feature --feature "Login" --scenario "OK" --step "открываю \"https://example.com\""

# Live browser
scenaria record --live --url https://example.com --output recorded.feature --idle 30 [--headless]
```

## `init`

```bash
scenaria init /path/to/project
```

## `update`

```bash
scenaria update --check
```

## `plugins`

```bash
scenaria plugins list
scenaria plugins install <path-or-id>
scenaria plugins uninstall <name>
```

## `va` (Vanessa Automation)

```bash
scenaria va run --project . [--dry-run] [--rerun-failed] ...
```

See `scenaria va run --help` for 1C-specific flags.

## Global install

```bash
go install ./cmd/scenaria
# or from module:
go install github.com/bafgion/scenaria-golang/cmd/scenaria@latest
```

## Related

- [Reports](reports.md)
- [Installation](../getting-started/install.md)
