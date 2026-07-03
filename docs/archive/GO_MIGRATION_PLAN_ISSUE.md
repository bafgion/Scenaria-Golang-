# Issue Draft: Full Scenaria MVP migration from Python to Go

## Context

Current repository contains foundational migration work for Scenaria Go:

- Go module bootstrap and CLI entrypoint
- Initial Gherkin parsing/validation
- Feature serialization and scenario store
- `run` preflight and `validate` command

Goal is full functional parity with the Python MVP including desktop flows.

## Goal

Migrate all MVP functionality from Python to Go while preserving user-facing behavior and file compatibility.

## Scope

- Gherkin RU parse/serialize/validate parity
- CLI command parity (`run`, `validate`, `export`, `record`, plugins)
- Browser execution via Playwright
- Recorder and selector heuristics
- Settings/test-client compatibility
- Reports and export parity
- Desktop application parity
- Packaging and update workflow parity

## Execution plan

## Phase 1: Core syntax and storage

- [x] Bootstrap Go module and CLI command router
- [x] Implement base Gherkin parse/validate support
- [x] Add support for tags, outlines, examples, docstrings, step tables
- [x] Implement feature serializer and file store
- [x] Add settings and test-client JSON compatibility layer

## Phase 2: CLI parity foundation

- [x] Implement `validate` command over scenario discovery
- [x] Implement `run` preflight with dry-run support
- [ ] Add CLI flag compatibility map vs Python version
- [ ] Add fixture-driven parity suite for CLI outputs

## Phase 3: Runner integration

- [ ] Introduce browser execution engine abstraction
- [ ] Integrate Playwright runner (initial Chromium path)
- [ ] Implement execution results model and error mapping
- [ ] Produce JUnit-compatible output

## Phase 4: Export/report

- [ ] Implement export pipelines used by current MVP
- [ ] Implement HTML report generator
- [ ] Add snapshot/golden tests for report compatibility

## Phase 5: Recorder and selectors

- [ ] Port recorder script injection flow
- [ ] Port selector build/resolve heuristics
- [ ] Add deterministic tests for recorder queues/events

## Phase 6: Desktop parity

- [ ] Select and lock desktop stack for Go implementation
- [ ] Implement workspace/editor flow
- [ ] Implement record/edit/run/report E2E desktop path
- [ ] Add desktop-level regression smoke tests

## Phase 7: Release pipeline

- [ ] Build packaging pipeline for desktop binaries
- [ ] Implement update metadata and updater flow
- [ ] Add release workflow checks in CI

## Definition of done

- Existing `.feature` scenarios execute with equivalent outcomes.
- Existing config and test-client files remain compatible.
- Core CLI command behavior and return codes are stable and documented.
- Desktop MVP user journey reaches feature parity with Python baseline.
- CI/release flow is reproducible in Go toolchain.
