# Functional Parity Matrix (Python MVP -> Go)

Status: migration **complete** for runtime/core; Qt IDE replaced by Fyne desktop.

| Area | Python capability | Go target | Status |
|---|---|---|---|
| Scenario format | `.feature` read/write | `internal/gherkin` | validated |
| Step DSL (~40 variants) | RU steps | `internal/stepdsl` | validated |
| Control flow | Если/Повторяю/Пока/Для каждого | `internal/gherkin`, `internal/player` | validated |
| Variables / generators | `{{name}}`, random data | `internal/player/context.go` | validated |
| Run suite | batch, tags, outline | `internal/player/suite.go` | validated |
| Playwright engine | playback + fallbacks | `internal/player`, `internal/selector/resolve.go` | validated |
| TestClient | cookies/localStorage | `internal/settings`, `internal/player` | validated |
| Validate | syntax + `--browser` | `internal/cli`, `internal/selector` | validated |
| Export | JSON/feature/TS/Python | `internal/exporter` | validated |
| Reports | JUnit/HTML | `internal/report` | validated |
| Run status | `.scenaria/run_status.json` | `internal/runstatus` | validated |
| Recorder | live browser capture | `internal/recorder` | validated |
| Selector heuristics | DOM → selector | `internal/selector/heuristics.go` | validated |
| OTP / email code | env + inference + segmented fill | `internal/player/otp.go` | validated |
| Vanessa add-on | `scenaria va run` | `internal/vanessa`, `internal/cli/va.go` | validated |
| CLI | run/validate/export/record/va | `cmd/scenaria` | validated |
| Desktop | IDE (Qt → Fyne) | `ui/desktop` | validated |
| Plugins registry | list/install | `internal/plugin` | validated |
| Packaging | portable ZIP + Chromium | `scripts/build-portable.ps1`, CI | validated |
| Update check | GitHub releases | `internal/update` | validated |

## Intentional differences

- **Desktop UI**: Fyne instead of Qt (~feature parity for run/validate/record/edit; not pixel-identical).
- **Email OTP**: interactive prompt in desktop (`player.EmailCodePrompt`); env/`--var` in CLI.
- **Vanessa**: requires local 1C platform + EPF paths in `.scenaria/vanessa.json`.

## Optional future work

- Golden cross-language CI against Python `tests/` fixtures
- Full Qt parity (completions, splash, update UI) if needed
- IMAP OTP fetch (net-new, not in Python MVP)
