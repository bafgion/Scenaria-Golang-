# Functional Parity Matrix (Python MVP -> Go)

This matrix defines migration tracking for end-to-end parity.

Status values:

- `planned` - identified, no implementation.
- `in-progress` - implementation started.
- `validated` - implementation and parity checks complete.

| Area | Python capability | Go target package | Status | Validation |
|---|---|---|---|---|
| Scenario format | Read/write `.feature` files | `internal/gherkin`, `internal/scenario` | in-progress | Parser + serializer round-trip tests + fixture-based `.feature` tests |
| Scenario validation | Syntax and semantic checks | `internal/gherkin`, `internal/cli` | in-progress | `scenaria validate` + validation unit tests |
| Run engine | Execute steps in browser | `internal/player` | in-progress | Runner abstraction + baseline Playwright executor + upcoming Chromium integration tests |
| Recording | Capture interactions to steps | `internal/recorder`, `internal/cli` | in-progress | Baseline `record` CLI writer + future browser capture parity scenarios |
| Selectors | Build and resolve selectors | `internal/selector` | planned | Selector fixture and DOM tests |
| CLI run | `scenaria run` parity | `cmd/scenaria`, `internal/cli` | in-progress | Preflight tests + summary/JUnit + engine selection + future parity comparisons |
| CLI export | `scenaria export` parity | `internal/exporter`, `internal/cli` | in-progress | Baseline JSON/feature export tests + future format parity snapshots |
| Reports | JUnit/HTML outputs | `internal/report` | in-progress | Run summary JSON + baseline JUnit XML implemented; HTML and snapshot parity checks pending |
| Settings | `settings.json` compatibility | `internal/settings` | in-progress | Round-trip JSON compatibility tests |
| Test clients | `.scenaria/test_clients/*.json` | `internal/settings`, `internal/scenario` | in-progress | Fixture compatibility tests |
| Plugin runtime | Runner extension points | `internal/plugin`, `internal/cli` | in-progress | Baseline local manifest registry + lifecycle smoke tests |
| Desktop shell | Main desktop workflow | `ui/desktop` | planned | E2E workflow tests |
| Update system | Release/update checks | `internal/update` | planned | Update metadata contract tests |
| Packaging | Portable distribution artifacts | `build/` or release workflow | planned | Release workflow dry runs |

## Notes

- Matrix must be updated in every migration PR.
- No area is considered complete without explicit validation evidence.
