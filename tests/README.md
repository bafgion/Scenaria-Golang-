# Тесты Scenaria Go

Стратегия тестирования и команды для разработчиков.

**Полное описание:** [docs/ru/contributing/development.md](../docs/ru/contributing/development.md) · [English](../docs/en/contributing/development.md)

## Слои

| Слой | Команда |
|------|---------|
| Go unit | `go test ./internal/... ./cmd/...` |
| Go integration | `go test -tags=integration ./internal/player/... ./internal/recorder/... ./internal/selector/...` |
| Frontend unit | `cd frontend && npm test` |
| UI E2E | `cd frontend && npm run test:e2e` |
| Desktop smoke | `./scripts/desktop-smoke.ps1` |

## CI

`.github/workflows/ci.yml` — test, integration, desktop-smoke на `master`.  
Релиз: тег `v*` → `.github/workflows/release.yml`.

## QA вручную

Чеклист: [docs/internal/QA-DAILY-USE.md](../docs/internal/QA-DAILY-USE.md)
