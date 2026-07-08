# State Ownership Map

This document defines the single owner for critical application state and the allowed write boundaries.

## Ownership Rules

- There is exactly one source of truth for each domain state.
- Derived state must be recomputed from owner state and never written independently.
- Cross-domain updates must go through explicit service/controller APIs, not ad-hoc mutation.

## Domain Ownership

### Editor Text

- Owner: frontend `tabsStore` + Monaco model for currently active tab.
- Source of truth:
  - For active tab: Monaco editor model text.
  - For inactive tabs: tab snapshot (`content`/`draft`) in `tabsStore`.
- Derived state:
  - Validation markers, inlay hints, step panel rows, scenario hints.

### Active File

- Owner: frontend `tabsStore`.
- Source of truth: `activeTab` identifier.
- Derived state:
  - `isWelcome`, `activeFeatureTab`, `activeTabUnsaved`.

### Project Session

- Owner: backend `ProjectService` (exposed via `gui.Service` facade).
- Source of truth:
  - Opened project root path.
  - Current project version identity.
  - Feature index and project metadata.
- Derived state:
  - Catalog tree, project tags, recents UI projections.

### Run Session

- Owner: backend `RunService`.
- Source of truth:
  - `runID`, run options snapshot, run lifecycle phase.
- Derived state:
  - Progress labels, last run error summary, run history projections.

### Recorder Session

- Owner: backend `RecorderService`.
- Source of truth:
  - `recordSessionID`, `browserSessionID`, lifecycle state (`recording`/`paused`/`stopped`), target path.
- Derived state:
  - Recorder toolbar labels, live step counter, picker UI state.

### Step Analysis

- Owner: backend `EditorAnalysisService`.
- Source of truth:
  - Single analysis result for editor text version (`issues`, `steps`, `hints` bundle).
- Derived state:
  - Markers in Monaco, step help availability, hint filtering.

### Report Artifacts

- Owner: backend `ReportService`.
- Source of truth:
  - Artifact references/paths in run results and report payload builders.
- Derived state:
  - HTML and Allure presentation fields (URLs, trace commands, screenshot links).

### Settings

- Owner: backend `SettingsStore`.
- Source of truth:
  - Persisted app settings payload.
- Derived state:
  - UI toggles, feature flags, editor preferences currently rendered.

## Allowed Write Paths

- UI writes only through frontend stores/controllers.
- Frontend requests backend mutations only through Wails APIs (`gui.Service` facade).
- Backend domain services perform persistence; `gui.Service` orchestrates and delegates.

## Non-Goals

- This document does not prescribe final folder layout; it defines ownership boundaries first.
- This document does not replace API docs; it constrains who may mutate which state.
