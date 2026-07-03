# Reports (CLI)

## JUnit XML

```bash
scenaria run ./features --junit junit.xml
```

Compatible with CI systems (Azure DevOps, Jenkins, GitLab).

## HTML report

```bash
scenaria run ./features --html report.html
# timestamped: configured via project or GUI run options
```

Opens a single-file summary of scenario results.

## Allure 2

```bash
scenaria run ./features --allure ./allure-results
scenaria run ./features --dry-run --allure ./allure-results
```

Each scenario → one `*-result.json` in the output directory.

On Playwright failures (non-dry-run), attachments may include:

| Attachment | When |
|------------|------|
| screenshot (PNG) | Browser failure |
| trace (ZIP) | With `--trace <dir>` |
| video (WebM) | With `--video <dir>` |

View locally:

```bash
allure serve ./allure-results
```

Install [Allure CLI](https://docs.qameta.io/allure/#_installing_a_commandline).

### Status mapping

| Scenaria | Allure |
|----------|--------|
| passed | passed |
| failed | failed |
| skipped | skipped |
| broken | broken |

## Trace and video

```bash
scenaria run ./features --trace ./traces --video ./videos --allure ./allure-results
```

Artifacts are retained on failure; GUI can open trace viewer when Playwright CLI is available.

## Summary JSON

```bash
scenaria run ./features --summary-json summary.json
```

Machine-readable batch result for automation.

## GUI

Same options in the **Run** dialog. Results in **Results** tab and `.scenaria/`. See [Results & reports](../gui/results-and-reports.md).
