# Functional Parity Matrix (Python MVP -> Go)

Status: **runtime/CLI parity with Python v0.12**; **Wails 2 GUI** in active development (Fyne deprecated).

| Desktop Wails | Monaco, tabs, run/record/Vanessa/export/OTP | `frontend/`, `internal/gui` | beta (~70%) |

| Area | Python capability | Go target | Status |
|---|---|---|---|
| Scenario format | `.feature` read/write + quote repair | `internal/gherkin`, `normalize.go` | validated |
| Step DSL (~40 variants) | RU steps, 1-based tabs, `жду N с` | `internal/stepdsl` | validated |
| Control flow | Если/Повторяю/Пока (max iter error)/Для каждого | `internal/gherkin`, `internal/player` | validated |
| Variables / generators | checksum INN/OGRNIP, coherent names | `internal/player/generators.go` | validated |
| Navigation | skip-goto, `urls_match`, nav timeout | `internal/player/urls.go` | validated |
| Run suite | batch, tags, outline, multi-path, `--workers` | `internal/player/suite.go`, `cli/run` | validated |
| Playwright engine | playback + fallbacks, `--slow-mo` | `internal/player`, `internal/selector` | validated |
| TestClient | cookies/localStorage | `internal/settings`, `internal/player` | validated |
| Validate | syntax + `--browser chromium|firefox|webkit` | `internal/cli`, `internal/selector` | validated |
| Export | JSON/feature/TS/Python + `--force` | `internal/exporter`, `internal/cli/export` | validated |
| Reports | JUnit/HTML | `internal/report` | validated |
| Run status | `.scenaria/run_status.json` | `internal/runstatus` | validated |
| Recorder | live capture + step coalescing | `internal/recorder`, `normalize.go` | validated |
| Selector heuristics | DOM → selector (canvas/signature) | `internal/selector/recorder_script.js` | validated |
| OTP / email code | segmented modes, auto-submit | `internal/player/otp.go` | validated |
| Vanessa add-on | exclude tags, rerun, epf install, JUnit monitor | `internal/vanessa` | partial |
| CLI | run/validate/export/import-json/record/va | `cmd/scenaria` | validated |
| Desktop (Fyne) | legacy IDE | `ui/desktop` | deprecated |
| Plugins registry | list/install zip/URL | `internal/plugin` | validated |
| Packaging | portable ZIP + Chromium | `scripts/build-portable.ps1`, CI | validated |
| Update check | GitHub releases | `internal/update` | validated |
| Version | 0.12.x | `internal/version` | **0.12.0** |

## Remaining gaps (Wails beta)

- Splash, portable `wails build`, Allure native writer.
- Recorder: MutationObserver, iframe/shadow DOM.
- Cross-language CI with Python `tests/`.

## Intentional differences

- **Desktop UI**: Fyne instead of Qt (~feature parity for run/validate/record/edit; not pixel-identical).
- **Email OTP**: interactive prompt in desktop (`player.EmailCodePrompt`); env/`--var` in CLI.
- **Vanessa**: requires local 1C platform + EPF paths in `.scenaria/vanessa.json`.

## Optional future work

- Golden cross-language CI against Python `tests/` fixtures
- Full Qt parity (completions, splash, update UI) if needed
- IMAP OTP fetch (net-new, not in Python MVP)
