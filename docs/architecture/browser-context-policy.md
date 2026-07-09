# Browser Context Policy

This document defines how Scenaria isolates browser state between scenarios and parallel workers.

## Modes

| Mode | Browser | Context | Page | Notes |
|------|---------|---------|------|-------|
| Sequential (`workers = 1`) | Reused for the whole run | New context per scenario after the first | New page per context | Live browser reuse is allowed only in this mode |
| Parallel pool (`workers > 1`, no trace/video) | One browser per worker slot | Reset between scenarios on the same slot | New page after reset | `resetForScenario()` recreates context |
| Parallel generic (`workers > 1`, trace/video enabled) | New browser per scenario | New context per scenario | New page per context | Pool is disabled when trace/video is on |

## When context reuse is allowed

| Situation | Reuse | Isolation mechanism |
|-----------|-------|---------------------|
| Sequential run, same browser | Browser reused | `resetForScenario()` after scenario 1+ |
| Parallel pool slot | Browser reused per slot | `browserPool.release()` → `resetForScenario()` |
| Live IDE browser (`AttachToPage`) | User's page/context | No automatic reset; user controls the window |
| Trace or video recording enabled | No pool | Fresh browser per scenario |
| Explicit test client profile | Auth state applied once per scenario | `ApplyTestClient()` at scenario start only |

## What `resetForScenario()` clears

Between scenarios on the same worker slot:

- browser context and page are recreated (old context closed);
- cookies, `localStorage`, and `sessionStorage` from the previous context are discarded;
- open tabs from the previous scenario are closed with the old context;
- network failure listener is rewired on the new page;
- trace recording is restarted when trace mode is enabled.

Creating a new Playwright browser context also drops route handlers, permissions, and other context-scoped state from the prior scenario.

## Auth/session reuse

Shared auth state is **not** implicit. Reuse is only allowed through explicit project/test-client configuration:

- feature-level `TestClient` declaration (`я подключаю TestClient "…"`);
- run dialog / CLI `--test-client` override;
- saved cookies and `localStorage` in `.scenaria/test_clients/<name>.json`.

Without a test client, each scenario starts from a clean context (after `resetForScenario()` or a new browser).

## Per-case artifacts

Failure artifacts are namespaced per runnable case:

- traces/videos/screenshots use `artifactBaseName()` (`CaseID` + example index);
- downloads use `RunContext.DownloadDir()` keyed by per-scenario `runSeed`;
- HTML report artifacts are stored under the run directory layout.

## Parallel variables

Each `RunCase` and `RunContext` receives a cloned variables map. Scenario A cannot mutate variables for scenario B, including under `workers > 1`.

## Related code

- `internal/player/browser_session.go` — `resetForScenario()`
- `internal/player/browser_pool.go` — worker slot acquire/release
- `internal/player/runner_parallel.go` — sequential vs parallel routing
- `internal/gui/run_execution.go` — `CanReuseLiveBrowser()`
- `internal/player/testclient.go` — explicit auth application
- `internal/player/artifacts.go` — per-case trace/screenshot paths
- `internal/player/browser_context_isolation_integration_test.go` — regression tests

## Regression tests

- `TestResetForScenarioClearsBrowserStorage` — storage/cookies cleared on reset
- `TestSequentialScenariosDoNotShareBrowserState` — scenario B does not see scenario A cookies
- `TestParallelPoolDoesNotShareBrowserStorage` — pool slots reset between scenarios
- `TestRunContextsUseDistinctDownloadDirs` — download paths isolated per scenario context