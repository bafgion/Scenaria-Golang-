# Разработка и тесты

Для участников разработки Scenaria из исходников.

## Требования

| Инструмент | Версия |
|------------|--------|
| Go | 1.26+ |
| Node.js | 20+ |
| Wails CLI | v2 |
| PowerShell | 7+ |

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

## Быстрый цикл

```bash
go test ./internal/... ./cmd/...
cd frontend && npm install && npm test
wails dev
go run ./cmd/scenaria --help
```

## Сборка

```powershell
./scripts/build-portable.ps1
./scripts/build-release.ps1 -SkipTests
```

Результат: `dist/Scenaria-Portable.zip`, `dist/Scenaria-Setup.exe`.

Тег `v*` → workflow `.github/workflows/release.yml`.

## Слои тестов

| Слой | Команда |
|------|---------|
| Go unit | `go test ./internal/... ./cmd/...` |
| Go integration | `go test -tags=integration ./internal/player/... ./internal/recorder/...` |
| Frontend unit | `cd frontend && npm test` |
| UI E2E | `cd frontend && npm run test:e2e` |
| Скриншоты для docs | `cd frontend && npm run docs:screenshots` |
| Desktop smoke | `./scripts/desktop-smoke.ps1` |

CI: `.github/workflows/ci.yml`.

## Глобальный CLI

```bash
go install ./cmd/scenaria
```

Модуль: `github.com/bafgion/scenaria-golang`. См. `docs/archive/CLI_GLOBAL_INSTALL.md`.

## Структура кода

| Путь | Назначение |
|------|------------|
| `cmd/scenaria/` | CLI |
| `internal/player/` | Runner Playwright |
| `internal/recorder/` | Запись |
| `internal/gui/` | Сервис Wails |
| `frontend/src/` | IDE Svelte |
| `examples/` | Примеры |

## Внутренние документы

- [ROADMAP](../internal/ROADMAP.md)
- [QA](../internal/QA-DAILY-USE.md)

## См. также

- [Индекс документации](../README.md)
