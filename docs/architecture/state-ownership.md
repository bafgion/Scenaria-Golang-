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
- Source of truth:
  - `tabs`, `activeTab`, `welcomeTabVisible`.
- Derived state:
  - `isWelcome`, `activeFeatureTab`, `activeTabUnsaved`.

### Dialog Visibility

- Owner: frontend `dialogsStore`.
- Source of truth: modal/palette `show*` flags (`showSettings`, `showRecord`, `showRun`, etc.).
- Derived state: `anyAppDialogOpen` (via `computeAnyAppDialogOpen` with overlay context).

### Record Form (dialog fields)

- Owner: frontend `recordFormStore` (scaffolded; dialog `bind:` locals sync on demand).
- Source of truth: `recordURL`, `recordOutput`, `recordIdle`, `recordAppendTo`, `recordTestClient`, `recordFeatureName`, `recordScenarioName`.

### Wails Events

- Owner: frontend `wailsEventsController`.
- Source of truth: event envelope unwrapping and stale-session guards.
- Write path: `bindWailsEvents()` registers listeners; domain handlers update stores via `App.svelte` callbacks.

### Feature Files (filesystem)

- Owner: backend `FileOperationService` (exposed via `gui.Service` facade).
- Source of truth: on-disk `.feature` files, drafts under `.scenaria/drafts`.
- Operations: read/save, delete, duplicate, move, import, rename, draft save/load/clear.

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
  - `runID`, run options snapshot, run lifecycle phase (`TryBegin` / `Finish` / `Cancel`).
  - In-process run orchestration (`ExecuteInProcess`, live-browser reuse, feature snapshot via `RunExecutionHost` callbacks).
- Derived state:
  - Progress labels, last run error summary, run history projections.

### Recorder Session

- Owner: backend `RecorderService` (`RecorderSessionManager` + `RecordLive` orchestration in `recorder_execution.go`).
- Source of truth:
  - `recordSessionID`, `browserSessionID`, generation counter, lifecycle state (`recording`/`paused`/`stopped`), target path, live browser handle.
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
  - Per-run report writes (HTML, JUnit, summary, Allure, latest pointer).
- Derived state:
  - HTML and Allure presentation fields (URLs, trace commands, screenshot links).

### Settings

- Owner: backend `SettingsService` (persisted via `settings.Store`); frontend `settingsStore` mirrors persisted app settings.
- Source of truth:
  - Persisted app settings payload.
  - HTTP auth credentials, recents, and app-level preferences.
- Derived state:
  - UI toggles, feature flags, editor preferences currently rendered.
  - Settings dialog uses local bindable fields synced from `settingsStore` on load/save.

### Test clients

- Owner: backend `TestClientService` (project-scoped JSON under `.scenaria/testclients`); frontend `testClientStore` mirrors list/selection for dialogs.
- Source of truth: test client JSON files on disk.
- Derived state: dialog selection and suggested capture name.

### Plugins

- Owner: backend `PluginService` (project-scoped plugin install/list/run); frontend `pluginsStore` mirrors installed plugin list for menus/dialogs.
- Source of truth: plugin descriptors and artifacts on disk under the project.
- Derived state: menu labels and Vanessa availability flags in UI.

### Step catalog

- Owner: backend `CatalogService` (step search and line completions via `stepcatalog`); frontend catalog UI state in `catalogStore`.
- Source of truth: bundled step catalog definitions.
- Derived state: filtered catalog tree, batch selection, collapsed folder keys.

## Frontend Derived State (App.svelte)

`App.svelte` is a UI shell. Domain stores are the write owners; local reactive aliases are read-only projections:

- `$projectStore` → `projectPath`, `features`, `tags`, `featureTags`, `currentProjectVersion`, `scenarios`, `artifacts` (report/trace paths)
- `$tabsStore` → `tabs`, `activeTab`, `welcomeTabVisible`, `pendingCloseTab` (unsaved close confirm), `loadFeatureGeneration` (async load staleness guard)
- `$settingsStore` → persisted app settings (dialog fields via `dialogBindController`; `bindEditorSettings` for Monaco/Settings bind), `startUrl`, `uiLocale`
- `$settingsDialogStore` → ephemeral settings dialog baseline and per-project HTML report open mode
- `$diagnosticsStore` → editor validation issues/hints, per-tab issue cache (`issuesByTab`), browser validate panel issues, dismissed hint keys, validate generation counter, `stepStatusError` (status bar), `hintFixInFlight`, validate debounce scheduler
- `$pluginsStore` → installed plugins list (`refreshInstalledPlugins`)
- `$updateDialogStore` → update check dialog state (message, info, download progress)
- `$catalogStore` → batch mode/selection, sidebar filter, collapsed keys, drop target, derived catalog tree (`baseTree` / `baseTreeKey`)
- `$contextMenuStore` → ephemeral feature/folder/steps context menus
- `$postRecordStore` → post-record banner state (path, step count, baseline diff text)
- `$confirmDialogStore` → global confirm dialog requests and record tab-switch skip preference
- `$menuStore` → menubar open section (`openMenu`)
- `$projectReplaceStore` → project-wide find/replace dialog state (`dialogBindController.bindProjectReplace*`)
- `$pickerDialogStore` → element picker step choices
- `$httpAuthDialogStore` → HTTP auth dialog host field
- `$stepsHelpDialogStore` → steps help dialog initial query
- `$otpDialogStore` → OTP dialog email from backend event
- `$journalStore` → journal log text and status bar message/tone
- `$featureDialogStore` → ephemeral feature CRUD/import/export dialog state (`dialogBindController` for duplicate/move/import fields)
- `$testClientStore` → project test client list and dialog selection (`dialogBindController.bindTestClientSelection` for TestClientDialog bind)
- `$validateDialogStore` → validate dialog form + CLI log output (`dialogBindController.bindValidate*`)
- `$runFormStore` → `lastRun`, `runForm` (RunDialog uses `dialogBindController.bindRunForm` synced on open)
- `$pluginRunStore` → plugin run dialog (`name`, `tag`, `scenario`; `dialogBindController.bindPluginRun*`)
- `$runDialogStore` → ephemeral Run dialog title/scenario list
- `$editorStore` → active Monaco buffer (`editorText`, `editorTextVersion`, `editorCursorLine`), parsed editor steps (`steps`, `stepsTextVersion`), steps panel tab (`stepsPanelTab`)
- `$appMetaStore` → desktop app version string (from Wails `Version()`)
- `sessionStore` (non-reactive orchestrator) → session persist debounce and draft autosave timers
- `$layoutStore` → sidebar/bottom panel/preview layout (`localStorage` via `lib/layout`), active bottom panel tab, session resize flags, `previewPaneMounted`, preview mount debounce timer
- `$splashStore` → startup splash screen message/progress/ready state
- `$viewportStore` → window dimensions and toolbar density flags
- `$uiPrefsStore` → toolbar/steps panel/sidebar width/onboarding flags (persisted via `SaveSettings`), session `stepsPanelCollapsed`
- `$onboardingTourStore` → tour visibility (`active`) and in-tour progress (`validateDone`, `dryRunDone`, `journalVisited`, `stepId`)
- `$runnerStore` → `playing`, run progress, last run batch/error panel state, `dryRunActive` (status bar)
- `$dialogsStore` → modal/palette visibility flags
- `$recorderPrefsStore` → recording filter prefs (`filterRecording`, `navOnlyRecording`, `hoverRecord`; synced via `dialogBindController` on load/save)
- `$recordFormStore` → record dialog form fields (`recordURL`, `recordOutput`, `recordMode`, `baselineBusy`, `recordStepPickerOpen`, etc.; `dialogBindController` sync on open/persist)
- `$recentsStore` → recent projects/features lists (persisted via `SaveSettings`)
- `$recorderStore` → `browserOpen`, `recording`, `recordPaused`, session IDs, live record step line map, `lastRecordTarget`, pause-toggle guard timestamp, browser watch timer, record step apply chain / editor-ready promises
- `$reportsStore` → run results, flaky metrics, Allure install/serve status, refresh in-flight/queued flags
- `$vanessaRunStore` → Vanessa dialog form + runtime (`running`, `snapshot`, `watchDir`, poll timer; `dialogBindController.bindVanessa*`)

`dialogBindController` owns ephemeral dialog `bind:` session state (validate, settings, record, run, Vanessa, plugin run, feature dialogs, project replace). Stores remain source of truth; controller `sync*` on dialog open and `flush*` on save/confirm.

`workspaceSessionController` owns workspace persistence orchestration: building `AppSettingsDTO` snapshots (open tabs, session project, recents), `persistSettings` (flush dialog binds then `SaveSettings`), draft autosave, and `restoreWorkspaceSession` on startup.

`paletteCommandsController` builds command-palette / hotkeys command lists from view state + action callbacks (no domain state ownership).

Event envelope unwrapping, stale-event guards, and Wails event subscriptions live in `wailsEventsController` (`bindWailsEvents`).

## Allowed Write Paths

- UI writes only through frontend stores/controllers.
- Frontend requests backend mutations only through Wails APIs (`gui.Service` facade).
- Backend domain services perform persistence; `gui.Service` orchestrates and delegates.

See also:

- [`app-shell.md`](app-shell.md) — what `App.svelte` may own vs delegate.
- [`gui-facade.md`](gui-facade.md) — backend façade and CLI shared handlers.

## Adding New State

Before introducing new mutable application state:

1. Choose a single owner (store, controller, or backend service).
2. Add a row to the relevant section in this document.
3. Mark derived fields explicitly; do not write derived state independently.
4. For UI-only ephemeral state (dialog bind session, palette actions), use an existing controller pattern or document a new one.

Pull requests that add state without updating this file should be rejected in review.

## Non-Goals

- This document does not prescribe final folder layout; it defines ownership boundaries first.
- This document does not replace API docs; it constrains who may mutate which state.
