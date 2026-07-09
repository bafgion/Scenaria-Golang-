# App Shell (`App.svelte`)

`App.svelte` is the desktop UI shell. It composes Svelte views, wires Wails bindings, and routes user intent to stores/controllers. It does **not** own domain state.

## Shell responsibilities (allowed)

| Area | Role |
|------|------|
| **Component tree** | Layout, menubar, catalog, editor host, dialogs, status bar, splash |
| **DOM refs** | `monaco`, `actionBarEl`, `onboardingTour` component instances |
| **Reactive projections** | Read-only aliases from stores (`$: projectPath` from `$projectStore`, etc.) |
| **Derived UI** | `lastRunSummary`, `paletteCommands`, catalog view state, layout metrics |
| **Action routing** | Thin handlers that call store methods or Wails APIs (`openProjectDialog`, `runPrimary`) |
| **Lifecycle** | `onMount` / `onDestroy`: splash, Wails events, resize listeners, session timers |
| **Controllers** | Instantiate `dialogBindController`, `workspaceSessionController`, `wailsEventsController` |

## Not owned by the shell

| Domain | Owner |
|--------|--------|
| Editor buffer / tabs | `tabsStore`, `editorStore`, Monaco |
| Project / features index | `projectStore` + backend `ProjectService` |
| Run / record / validate sessions | `runnerStore`, `recorderStore`, backend services |
| Settings persistence | `settingsStore` + `workspaceSessionController` |
| Dialog form fields | Domain stores + `dialogBindController` bind session |
| Modal visibility | `dialogsStore` |

## Write path rule

The shell may **read** store snapshots and **invoke** store/controller/service methods. It must not keep parallel mutable copies of domain fields (except ephemeral bind-session in `dialogBindController`).

## Adding UI behavior

1. Add or extend a store/controller for state and side effects.
2. Expose a shell handler that delegates in one or two lines.
3. Document new critical state in [`state-ownership.md`](state-ownership.md).
