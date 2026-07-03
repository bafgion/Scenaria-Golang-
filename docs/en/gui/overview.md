# GUI overview

The Scenaria IDE is a single-window desktop app (Wails + WebView2) organized like a code editor.

## Main areas

```
┌─────────────────────────────────────────────────────────────┐
│ Menu bar (File, Project, Run, Record, View, Help, …)      │
├──────────┬──────────────────────────────────┬───────────────┤
│ Catalog  │  Editor tabs + Monaco            │ Preview (opt) │
│ sidebar  │                                  │               │
├──────────┴──────────────────────────────────┴───────────────┤
│ Bottom panel: Journal | Results | Validate | Error          │
├─────────────────────────────────────────────────────────────┤
│ Status bar: message, browser/run indicator, progress        │
└─────────────────────────────────────────────────────────────┘
```

| Area | Purpose |
|------|---------|
| **Start** tab | Welcome, recent projects, open examples, new project |
| **Catalog** | Feature tree, tags, run badges, batch selection |
| **Editor** | Gherkin `.feature` files with syntax highlighting |
| **Bottom panel** | Live log (journal), run results, validation issues |
| **Status bar** | Current operation, link to journal |

## Key menus

| Menu | Common actions |
|------|----------------|
| **Project** | Open, close, init, new scenario |
| **Run** | Run…, Dry-run, Validate…, Run with tag… |
| **Record** | Start recording…, TestClient…, browser tools |
| **View** | Toggle catalog, journal (Ctrl+`), preview |
| **Help** | F1 step catalog, Training tour, About |

## Command palette

**Ctrl+Shift+P** — searchable list of ~80 commands (run, save, export, settings, …).

## Browser overlay

When the browser is open, recording, or a test is running, a toolbar appears:

- Focus browser, pause/resume recording, stop
- Picker, on-page selectors validation, headless toggle (recording)

Dialogs take priority over the overlay.

## Hotkeys (essentials)

| Shortcut | Action |
|----------|--------|
| Ctrl+S | Save |
| Ctrl+Enter | Run / open run dialog |
| Ctrl+Shift+R | Stop test or recording |
| Ctrl+B | Open browser |
| Ctrl+` | Toggle bottom panel / journal |
| Ctrl+Shift+P | Command palette |
| F1 | Step reference |
| Shift+F1 | All hotkeys |

Full list: **Shift+F1** in the app.

## Related

- [Editor](editor.md)
- [Recording](recording.md)
- [Running tests](running-tests.md)
- [Settings](settings.md)
