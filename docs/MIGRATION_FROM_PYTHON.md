# Migration from Python Scenaria

This guide is for teams moving from **Python Scenaria (v0.12)** to **Scenaria Go**.

## What stays the same

| Asset | Compatibility |
|-------|----------------|
| `.feature` files (RU Gherkin) | Full — same parser/serializer semantics |
| `.scenaria/project.json` | Full |
| `.scenaria/test_clients/*.json` | Full (cookies, localStorage, base_url) |
| Step DSL (~40 RU steps) | Full — see `internal/stepdsl/testdata/steps_golden.json` |
| CLI commands | `run`, `validate`, `export`, `record`, `va`, `plugins`, `init` |
| Tags, outlines, examples | Full |
| JUnit / HTML reports | Full (`--junit`, `--html`) |

## What changed

| Area | Python | Go |
|------|--------|-----|
| Desktop IDE | Qt | **Wails 2 + Svelte** (Fyne deprecated) |
| Runtime | Python + Playwright | Go + `playwright-go` |
| Install | `pip install scenaria` | Portable ZIP or `go install ./cmd/scenaria` |
| Email OTP | env / prompt | Same + Wails modal |
| Vanessa | 1C + EPF | Same requirements |

## CLI mapping

```bash
# Python                          # Go
scenaria run ./features           scenaria run ./features
scenaria run -t smoke             scenaria run ./features --tag smoke
scenaria run -e VAR=1             scenaria run ./features --var VAR=1
scenaria validate                 scenaria validate ./features
scenaria export                   scenaria export in.feature --output out.json
```

Go additions:

```bash
scenaria run ./features --workers 4 --slow-mo 100
scenaria run ./features --test-client DemoUser
scenaria import-json doc.json --output doc.feature
scenaria va run --project . --rerun-failed
```

## Project layout

```
my-project/
  *.feature
  .scenaria/
    project.json
    test_clients/
      DemoUser.json
    vanessa.json          # optional
    run_status.json       # written after runs
```

Initialize:

```bash
scenaria init .
# copies DemoUser.json.example → rename to DemoUser.json
```

## TestClient

Background block in feature:

```gherkin
Контекст:
  Допустим я подключаю TestClient "DemoUser"
```

Override from CLI/GUI without editing feature:

```bash
scenaria run . --test-client DemoUser
```

## Desktop (Wails)

```bash
# Development
wails dev

# Production build
wails build -o scenaria-gui.exe
```

Legacy Fyne UI (deprecated):

```bash
go run -tags desktop ./cmd/scenaria-gui
```

## Known gaps

See `docs/FUNCTIONAL_PARITY_MATRIX.md` and `docs/ROADMAP.md`.

- Cross-language CI against Python `tests/` repo (planned)
- Allure native writer from Go runner (in progress)
- Full Qt pixel-parity in desktop (not a goal)

## Getting help

1. Run `scenaria validate ./your-project --browser chromium` for selector issues.
2. Compare step parsing: golden fixture in `internal/stepdsl/testdata/steps_golden.json`.
3. Open an issue with `.feature` sample and `scenaria run --dry-run` output.
