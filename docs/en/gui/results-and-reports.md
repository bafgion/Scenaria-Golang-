# Results and reports

## Bottom panel tabs

| Tab | Content |
|-----|---------|
| **Journal** | CLI-style log: validate, run, record messages |
| **Results** | Per-scenario status from last batch run |
| **Validate** | Browser validation issues list |
| **Error** | Last fatal error with jump-to-line |

Click a failed result to open the feature and go to the failed step.

## Artifacts

Stored under `.scenaria/` by default:

| Artifact | Enable | Location |
|----------|--------|----------|
| `run_status.json` | Always (GUI runs) | `.scenaria/` |
| JUnit XML | Run dialog / `--junit` | configurable path |
| HTML report | `--html` | `.scenaria/report.html` |
| Allure | Run dialog / `--allure` | `.scenaria/allure-results/` |
| Trace ZIP | `--trace` (failures) | `.scenaria/traces/` |
| Video WebM | `--video` (failures) | `.scenaria/videos/` |

## Allure

1. Enable **Allure** in run options
2. After run: **Open Allure** in results panel (opens folder)
3. Install [Allure CLI](https://docs.qameta.io/allure/) locally: `allure serve .scenaria/allure-results`

Status mapping: passed / failed / skipped / broken — same as Scenaria run status.

On failure, attachments may include screenshot, trace, video.

## Trace viewer

For failed Playwright runs with trace enabled, IDE offers **Open trace** (Playwright trace viewer command if installed).

## Flaky detection

Results panel can highlight scenarios that failed intermittently; **Run flaky 3×** reruns selected scenarios three times.

## Run history

**View → Run history…** — past runs from `run_status.json`, rerun failed scenarios.

## CLI reports

See [CLI reports](../cli/reports.md).

## Related

- [Running tests](running-tests.md)
- [CLI reference](../cli/reference.md)
