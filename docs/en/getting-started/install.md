# Installation

Scenaria runs on **Windows 10/11 (x64)**. Linux/macOS support the CLI only (no Wails GUI build in release artifacts).

## Recommended: Windows installer

Download **Scenaria-Setup.exe** from [GitHub Releases](https://github.com/bafgion/Scenaria-Golang-/releases).

- Installs to `Program Files\Scenaria`
- Includes `scenaria.exe`, `scenaria-gui.exe`, bundled Chromium, and sample scenarios
- Optional task: **add `scenaria` to system PATH** (open a new terminal after install)

Start the IDE from the Start menu or run `scenaria-gui.exe`.

## Portable ZIP

Download **Scenaria-Portable.zip**, extract anywhere (e.g. `D:\Tools\Scenaria`).

```
Scenaria/
  scenaria.exe          # CLI
  scenaria-gui.exe      # IDE
  Start-GUI.bat         # quick launcher
  browsers/             # Playwright Chromium
  examples/             # sample .feature files
  README-PORTABLE.txt
```

Portable builds use **robocopy** for in-app updates; the installer uses a silent Inno Setup upgrade.

## Global CLI (developers)

If you build from source:

```bash
go install ./cmd/scenaria
scenaria version
```

See [CLI global install notes](../contributing/development.md#global-cli) and `docs/archive/CLI_GLOBAL_INSTALL.md`.

## Requirements (from source)

| Component | Purpose |
|-----------|---------|
| Go 1.26.5+ | CLI, Wails backend |
| Node.js 22+ | Frontend build (CI uses 22) |
| [Wails CLI](https://wails.io/docs/gettingstarted/installation) | `wails dev` / `wails build` |
| Inno Setup 6 | Windows installer (optional) |

## Auto-update

**Help → Check for updates** (or startup prompt after project open). The IDE reads `latest.json` from GitHub releases and offers:

| Install type | Update package |
|--------------|----------------|
| Inno Setup | `Scenaria-Setup.exe` (silent `/VERYSILENT`) |
| Portable | `Scenaria-Portable.zip` |

Downloads are verified with SHA256 from the release manifest.

## Next steps

- [First project](first-project.md)
- [Onboarding tour](onboarding.md)
