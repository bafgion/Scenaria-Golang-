# Allure Report (Scenaria Go)

Scenaria writes [Allure 2](https://docs.qameta.io/allure/) JSON results after a run.

## CLI

```bash
scenaria run ./features --allure ./allure-results
scenaria run ./features --dry-run --allure ./allure-results
```

Each scenario becomes one `*-result.json` file in the output directory.

On Playwright failures, attachments may include:

- `screenshot` (PNG) — always when browser fails
- `trace` (ZIP) — when `scenaria run --trace <dir>` is set
- `video` (WebM) — when `scenaria run --video <dir>` is set

## CLI flags

```bash
scenaria run ./features --allure ./allure-results
scenaria run ./features --trace ./traces --video ./videos --allure ./allure-results
```

Install [Allure CLI](https://docs.qameta.io/allure/#_installing_a_commandline) and run:

```bash
allure serve ./allure-results
```

## GUI (Wails)

Enable **Allure** in the run panel before Playwright run. After completion, use **Открыть Allure** to open the results folder in Explorer.

## Status mapping

| Scenaria | Allure |
|----------|--------|
| passed | passed |
| failed | failed |
| skipped | skipped |
| broken | broken |

Attachments (screenshots on failure) are written when Playwright run fails.
