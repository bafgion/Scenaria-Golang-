# Migration from Python Scenaria

For teams moving from **Python Scenaria (v0.12)** to **Scenaria Go**.

## What stays the same

| Asset | Compatibility |
|-------|----------------|
| `.feature` files (RU Gherkin) | Full |
| `.scenaria/project.json` | Full |
| `.scenaria/test_clients/*.json` | Full |
| Step DSL (~40 RU steps) | Full |
| CLI commands | `run`, `validate`, `export`, `record`, `va`, `plugins`, `init` |
| Tags, outlines, examples | Full |
| JUnit / HTML reports | Full |

## What changed

| Area | Python | Go |
|------|--------|-----|
| Desktop IDE | Qt | **Wails + Svelte** |
| Runtime | Python + Playwright | Go + playwright-go |
| Install | `pip install scenaria` | Installer / portable / `go install` |
| Email OTP | env / prompt | Same + Wails modal |

## CLI mapping

```bash
# Python                          # Go
scenaria run ./features           scenaria run ./features
scenaria run -t smoke             scenaria run ./features --tag smoke
scenaria run -e VAR=1             scenaria run ./features --var VAR=1
scenaria validate                 scenaria validate ./features
scenaria export                   scenaria export in.feature --output out.json
```

Go additions: `--workers`, `--test-client`, `import-json`, `va run --rerun-failed`.

## Export to Python (optional)

```bash
scenaria export ./login.feature --format python --output test_login.py
```

## Related

- [CLI reference](../cli/reference.md)
- [First project](../getting-started/first-project.md)

Historical parity notes: `docs/archive/FUNCTIONAL_PARITY_MATRIX.md`
