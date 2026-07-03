# third_party patches

## go-webview2

Vendored copy of `github.com/wailsapp/go-webview2@v1.0.22` with a small patch for desktop smoke CDP.

Wails clears `WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS` at startup, so remote debugging is enabled when:

- `SCENARIA_DESKTOP_SMOKE=1` (set by `scripts/desktop-smoke.ps1`)
- optional `SCENARIA_DESKTOP_SMOKE_PORT` (default `9333`)

Patched files:

- `pkg/edge/chromium.go` — append `desktopSmokeBrowserArgs()`
- `pkg/edge/desktop_smoke_args.go` — new

Activated via `replace` in root `go.mod`. Production builds without the env var are unchanged.
