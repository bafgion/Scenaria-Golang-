# Scenaria (Go Migration)

This repository contains the Go implementation of Scenaria with a goal of full
functional parity with the Python MVP (desktop + recorder + runner + export +
plugins).

## Current state

The project is in migration bootstrap stage:

- Go module and project layout are initialized.
- `scenaria` CLI entrypoint is available.
- Migration plan and parity matrix are documented.

## Project goals

1. Preserve existing user-facing functionality from Python MVP.
2. Keep compatibility for existing scenario and settings files.
3. Provide a globally installable CLI command (`scenaria`).
4. Rebuild desktop functionality in Go after core parity is stable.

## Local development

```bash
go test ./...
go run ./cmd/scenaria --help
```

## Install CLI as a global command

```bash
go install ./cmd/scenaria
```

With a hosted module path:

```bash
go install <module-path>/cmd/scenaria@latest
```

See details in `docs/CLI_GLOBAL_INSTALL.md`.

## Migration documentation

- `docs/MIGRATION_PLAN.md` — staged migration roadmap.
- `docs/FUNCTIONAL_PARITY_MATRIX.md` — functional parity contract.
