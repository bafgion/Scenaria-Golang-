# Functional Parity Matrix (Python MVP -> Go)

Status: **CLI parity v0.12**; **Wails 2 GUI beta v0.13** (Fyne deprecated).

| Wails beta | Monaco, tabs, splash, run/record/Vanessa/export/OTP/Allure | `frontend/`, `internal/gui` | ~85% |

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
| Reports | JUnit/HTML/Allure | `internal/report`, `internal/report/allure` | validated |
| Run status | `.scenaria/run_status.json` | `internal/runstatus` | validated |
| Recorder | live capture + pause + step coalescing | `internal/recorder`, `docs/RECORDER.md` | validated |
| Selector heuristics | DOM → selector (canvas/signature/shadow) | `internal/selector/recorder_script.js` | validated |
| OTP / email code | segmented modes, auto-submit | `internal/player/otp.go` | validated |
| Vanessa add-on | exclude tags, rerun, epf install, JUnit monitor | `internal/vanessa` | partial |
| CLI | run/validate/export/import-json/record/va | `cmd/scenaria` | validated |
| Desktop (Wails) | primary IDE | `main.go`, `frontend/` | beta |
| Desktop (Fyne) | legacy IDE | `ui/desktop` | deprecated |
| Plugins registry | list/install zip/URL | `internal/plugin` | validated |
| Packaging | portable ZIP + Wails GUI + Chromium | `scripts/build-portable.ps1`, CI | validated |
| Update check | GitHub releases | `internal/update` | validated |
| Version | 0.12.x | `internal/version` | **0.13.0** |

## Remaining gaps

- Allure trace/video opt-in
- Cross-language CI with Python `tests/`
- Remove Fyne from release (v0.15)

## Intentional differences

- **Desktop UI**: Wails + Monaco (Fyne deprecated).
- **Email OTP**: interactive prompt in Wails (`player.EmailCodePrompt`); env/`--var` in CLI.
- **Vanessa**: requires local 1C platform + EPF paths in `.scenaria/vanessa.json`.
