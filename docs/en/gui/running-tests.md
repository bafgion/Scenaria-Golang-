# Running tests

## Dry-run vs Playwright run

| Mode | Browser | Use case |
|------|---------|----------|
| **Dry-run** | No | Parse steps, validate DSL, fast feedback |
| **Playwright run** | Yes | Real UI test execution |

**Ctrl+Enter** — run (opens dialog on first use). **Run → Dry-run** for syntax-only path.

## Run dialog

First run shows **Run scenario** dialog with options:

- Scenario name or full file
- Tag filter, variables (`KEY=value`)
- Browser, workers, slow-mo, base URL
- Reports: JUnit, HTML, Allure, trace, video
- Partial run: start/end step indices

**Don't show again** skips the dialog; change in settings or reset via **Run → Run…**.

## Validate

**Run → Validate…**

| Option | Checks |
|--------|--------|
| Syntax only | Gherkin structure, known steps |
| In browser | Selectors visible on live page (needs base URL / open browser) |
| Current file / whole project | Scope |

Results appear in **Validate** or **Journal** tab.

## Batch run

Select multiple features in the catalog (checkbox mode), then **Run** — runs sequentially or with configured workers.

## During execution

- Status bar: **Test running** with progress (scenario N of M)
- **Journal** streams log output live
- **Ctrl+Shift+R** — cancel (context canceled)

## Partial run

From editor code lens or context menu: run single scenario, from line, or to line.

## CLI

```bash
scenaria run ./features --dry-run
scenaria run ./features --tag smoke --headed --install-playwright
scenaria run ./features --workers 2 --var BASE=https://example.com
```

See [CLI reference](../cli/reference.md).

## Related

- [Results & reports](results-and-reports.md)
- [Settings](settings.md)
