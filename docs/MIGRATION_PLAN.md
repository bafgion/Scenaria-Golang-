# Scenaria Python -> Go Migration Plan

This plan targets full functional parity with the Python MVP while preserving
existing user workflows and file compatibility.

## Migration principles

1. Preserve behavior before optimizing internals.
2. Keep scenario/config formats backward compatible.
3. Ship in vertical slices: parser + CLI + runner before desktop rewrite.
4. Every migrated capability requires parity tests.

## Functional scope to migrate

- Gherkin RU parsing, serialization, and validation.
- Browser run/playback based on Playwright.
- Browser recording and selector generation.
- CLI commands (`run`, `validate`, `export`, plugin operations).
- Desktop UI workflows currently provided by Qt implementation.
- Reports (JUnit/HTML) and export pipelines.
- Update and packaging workflow for desktop distribution.

## Target architecture (Go)

```text
cmd/
  scenaria/           # CLI entrypoint and command routing
  scenaria-gui/       # desktop entrypoint (phase 2+)
internal/
  gherkin/            # parse, serialize, context, outline
  scenario/           # feature I/O and models
  player/             # scenario playback over Playwright
  recorder/           # capture browser interactions
  selector/           # selector heuristics and resolution
  report/             # junit/html outputs
  settings/           # app and project settings
  plugin/             # plugin lifecycle and runner registry
  update/             # app update checks and downloads
ui/
  desktop/            # desktop composition and state wiring
```

## Phase plan

## Phase 0: bootstrap and parity baseline

- Initialize Go module and command layout.
- Build migration matrix for Python modules and tests.
- Import representative golden fixtures for `.feature` and settings.
- Establish CI baseline (`go test`, lint, formatting checks).

Exit criteria:

- Build and tests run in CI.
- Parity matrix is documented and tracked.

## Phase 1: parser/model compatibility

- Implement `internal/gherkin` parse/serialize for RU syntax.
- Implement scenario structures and feature file persistence.
- Implement settings and test-client JSON compatibility.

Exit criteria:

- Golden tests pass for parser and serializer.
- Existing project files load/save without format drift.

## Phase 2: CLI + run/validate core

- Implement `scenaria run`.
- Implement `scenaria validate`.
- Introduce runner abstraction for dry-run and browser execution backends.
- Add JUnit-compatible report output (baseline CLI writer in place).
- Integrate Playwright execution path in Go worker model.

Exit criteria:

- CLI commands run with parity behavior on fixture suites.
- Integration tests pass with Chromium.

## Phase 3: export/report parity

- Implement export formats used by current Python flow.
- Implement HTML/JUnit reports with stable schema.

Exit criteria:

- Export and reports match expected outputs on fixtures.

## Phase 4: recorder backend migration

- Port recorder script injection and callback plumbing.
- Port selector heuristics and picker behavior.
- Add deterministic tests for command queue and event ordering.

Exit criteria:

- Recorded output scenarios are equivalent for key flows.

## Phase 5: desktop migration

- Implement Go desktop shell and workspace flow.
- Port editor, runner controls, logs panel, settings interactions.
- Recreate worker bridge semantics for UI-thread safe updates.

Exit criteria:

- End-to-end desktop flow: record -> edit -> run -> report.

## Phase 6: packaging, update, and release

- Rebuild packaging pipeline for desktop binaries.
- Reimplement update check/download/install flow.
- Publish artifacts and update metadata in release process.

Exit criteria:

- Release artifacts are reproducible and update-compatible.

## CLI global command requirements

The CLI must remain usable as a global command in both development and release:

- Development install: `go install ./cmd/scenaria`
- Module install: `go install <module-path>/cmd/scenaria@latest`
- Packaged binaries should expose the same top-level command name: `scenaria`

Implementation notes:

- Keep command names stable to preserve automation scripts.
- Avoid breaking flag names unless compatibility shim is provided.

## Definition of functional parity

Parity is achieved when:

1. Existing `.feature` files execute with equivalent outcomes.
2. Existing settings and test-client JSON files remain compatible.
3. CLI behavior and return codes are consistent for core commands.
4. Desktop core workflow is equivalent for standard user scenarios.
5. Existing CI/CD and release flows are reproducible in Go toolchain.
