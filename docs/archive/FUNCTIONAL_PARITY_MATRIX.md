# Functional Parity Matrix (Python MVP -> Go)

Status: **Go is the primary product** (CLI + Wails IDE). Python Scenaria is **deprecated** — interoperability via `.feature` / `.scenaria` files and **export to Python** only.

| Wails IDE | Monaco, run/record/Vanessa/export/OTP/Allure | `frontend/`, `internal/gui` | ~98% |

| Area | Legacy (Python/Qt) | Go target | Status |
|---|---|---|---|
| Scenario format | `.feature` RU Gherkin | `internal/gherkin` | validated |
| Step DSL (~40 variants) | RU steps | `internal/stepdsl` | validated |
| Control flow | Если/Повторяю/Пока/Для каждого | `internal/player` | validated |
| Run suite | tags, outline, workers, `--scenario` | `internal/cli/run` | validated |
| Playwright engine | playback | `internal/player` | validated |
| Validate / Export | syntax, JSON/TS/**Python**/feature | `internal/cli`, `internal/exporter` | validated |
| Reports | JUnit/HTML/Allure/trace/video | `internal/report` | validated |
| Recorder | live + baseline steps | `internal/recorder` | validated |
| Vanessa | JUnit monitor, rerun, EPF | `internal/vanessa`, Wails monitor UI | validated |
| Desktop IDE | Qt (deprecated) | **Wails 2 + Svelte** | primary |
| Plugins | registry + run | `internal/plugin` | validated |

## Remaining (Go-only roadmap)

- Optional: push release tag and publish portable zip artifact

## Done (v0.15)

- Recorder: SELECT → «выбираю», live option sync during recording
- Portable zip: README, launchers, version from `internal/version`
- Run/Vanessa/Plugin dialogs: project scenario catalog + chips
- Legacy Fyne `ui/desktop` removed from repo
- Auto-save dirty features before run; batch run respects cursor scenario

## Intentional differences

- **IDE**: Wails + Monaco — not Qt pixel-parity.
- **Python**: no runtime dependency; use `scenaria export --format python` for codegen.
- **Email OTP**: Wails modal + CLI env/`--var`.
