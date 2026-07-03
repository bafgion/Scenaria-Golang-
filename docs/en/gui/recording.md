# Recording scenarios

Record browser actions into `.feature` steps using Playwright live capture.

## Start recording (GUI)

1. Open or create a `.feature` file (target for new steps)
2. **Record → Start recording…** (or toolbar)
3. Enter start URL, output options, idle timeout
4. **Start** — browser opens; actions are appended to the scenario

During recording:

| Control | Action |
|---------|--------|
| **Pause / Resume** | Temporarily ignore events |
| **Undo step** | Remove last recorded line from editor |
| **Stop** | End session, normalize steps, close or keep browser |

Recording toolbar: filter noise, links only, hover capture, **headless** (restarts browser with same URL).

## What gets recorded

- Navigation (`goto`) on URL change
- Clicks (with debounce), fill, select, checkbox, file upload
- Keys, scroll, drag-and-drop, optional hover
- Canvas signature (`draw-signature`)

Selectors are built by priority: testid → id → aria/role → label → text chains. See [Selectors](../authoring/selectors.md).

## Picker

**Point to element** while recording or from browser overlay — inserts a selector without a full click recording.

## Post-record

After stop, a banner may offer **Review diff** comparing editor text before/after recording.

## HTTP authentication

If the site needs basic auth, **HTTP Auth…** before recording stores credentials for the session.

## TestClient

Save logged-in state: **Record → TestClient…** or **Save browser session…**. Use in feature:

```gherkin
Контекст:
  Дано я подключаю TestClient "DemoUser"
```

## CLI recording

```bash
scenaria record --live --url https://example.com --output recorded.feature --idle 30
scenaria record --live --url https://example.com --headless --output out.feature
```

| Flag | Description |
|------|-------------|
| `--live` | Interactive browser capture |
| `--url` | Start URL |
| `--output` | Target `.feature` path |
| `--idle N` | Stop after N seconds without events (default 30) |
| `--headless` | No browser window |

## Security

Recording injects trusted scripts via Playwright bindings. Use only on **your** test environments; paths are confined to the project root.

## Related

- [Selectors](../authoring/selectors.md)
- [Gherkin](../authoring/gherkin.md)
- [CLI reference](../cli/reference.md)
