# Global CLI Command Strategy

Scenaria CLI must be available as a global command named `scenaria`.

## Development install

From repository root:

```bash
go install ./cmd/scenaria
```

Ensure `$GOBIN` or `$GOPATH/bin` is in `PATH`.

## Module install

When module is published and tagged:

```bash
go install github.com/bafgion/scenaria-golang/cmd/scenaria@latest
```

## Binary distribution requirements

- Binary name: `scenaria` (`scenaria.exe` on Windows).
- Command contract:
  - Stable top-level command names.
  - Stable exit code behavior for automation.
  - Backward-compatible flags for critical workflows.

## Compatibility guardrails

1. Keep aliases for renamed commands.
2. Add deprecation warnings before removing flags.
3. Document all CLI behavior differences in changelog.

## Validation checklist

- [ ] `go install ./cmd/scenaria` succeeds on supported platforms.
- [ ] `scenaria help` displays core commands.
- [ ] `scenaria run` and `scenaria validate` return predictable exit codes.
- [ ] CI has at least one job that installs CLI via `go install`.
