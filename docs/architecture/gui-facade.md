# GUI Facade (`internal/gui.Service`)

`gui.Service` is the Wails-facing façade. It coordinates project lifecycle boundaries (open/switch project, cancel in-flight work) and delegates domain logic to extracted services.

## Domain services

| Service | Responsibility |
|---------|----------------|
| `ProjectService` | Project scan, feature index, tags |
| `FileOperationService` | Feature CRUD, drafts, import/export paths, replace |
| `RunService` | Run session identity, cancel, in-process execution host |
| `RecorderService` | Record/browser session lifecycle |
| `EditorAnalysisService` | Validate feature text, parse steps, hints |
| `ReportService` | Run results, flaky metrics, artifacts, Allure |
| `SettingsService` | App settings load/save |
| `TestClientService` | Test client JSON CRUD |
| `PluginService` | Plugin list/install/run (Vanessa, etc.) |
| `CatalogService` | Step catalog search / completions |

## Façade vs domain

`Service` methods should be thin:

- **Delegate** — call the matching `*Service` (`ReadFeature` → `FileOperationService`, `ValidateFeature` → `EditorAnalysisService`, `Run` → `RunService` + `runInProcess`).
- **Orchestrate** — project switch: cancel run/validate/recorder, rotate `ProjectSession`, bump version.
- **Adapt** — map Wails DTOs, attach `runId` to events, format `RunResult` for the frontend.

Project tag indexing (`collectFeatureTags`, `collectProjectTags`) lives in `project_tags.go` and is used by `ProjectService`, not inline in `service.go`.

## CLI integration (no global stdout)

GUI does **not** subprocess the CLI binary and does **not** replace `os.Stdout` globally.

| Path | Mechanism |
|------|-----------|
| **Run / dry-run** | `RunService` + `runInProcess` (in-process player) |
| **Recorder / browser** | `RecorderService` |
| **Files / settings / reports** | Domain services |
| **Init, validate (CLI output), export/import JSON, baseline record, plugin VA, updates** | `CLIOps` calls shared `internal/cli.Run*WithOutput(ctx, args, io.Writer)` handlers with a **buffer writer** |

CLI entry (`main`) and GUI (`CLIOps`) share the same Go handler functions. GUI captures output into a string buffer; CLI writes to stdout when `out == nil`.

## Adding backend features

1. Implement logic in a domain service (or extend an existing one).
2. Expose via a thin `Service` method for Wails.
3. If CLI needs the same behavior, add or reuse a `cli.Run*WithOutput` handler — do not add GUI-only globals.
4. Update [`state-ownership.md`](state-ownership.md) when introducing new mutable state.
