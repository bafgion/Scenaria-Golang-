# Тесты Scenaria Go



Стратегия переноса с Python (`shop-ui-recorder/tests/`) и покрытия Wails UI.



## Слои



| Слой | Где | Команда | Что покрывает |

|------|-----|---------|---------------|

| **Go unit** | `internal/**`, `cmd/**` | `go test ./internal/... ./cmd/...` | Gherkin, DSL, player, GUI service, recorder normalize, hints, save |

| **Go integration (браузер)** | `internal/recorder/live_integration_test.go`, `internal/player/browser_action_integration_test.go`, `internal/selector/*_integration_test.go` | `go test -tags=integration ./internal/recorder/... ./internal/player/... ./internal/selector/...` | Живая запись Playwright; player: goto, click, fill, hover, select, check, assert, scroll, generators |

| **Desktop smoke (WebView2)** | `internal/wailsapp/desktop_smoke_test.go` | `go test -tags=desktop ./internal/wailsapp/...` | `scenaria-gui.exe` запускается и не падает N секунд |

| **Frontend unit** | `frontend/src/**/*.test.ts` | `cd frontend && npm test` | hotkeys, catalogTree, stepSearch, featureTemplate, gherkinCompletions, gherkinHintActions |

| **UI E2E (Playwright)** | `frontend/e2e/specs/` | `cd frontend && npm run test:e2e` | Ключевые сценарии Wails UI через vite preview + мок Wails (dry-run, post-record hints, Ctrl+Space) |



## Перенос с Python



### Перенесено



| Python | Go |

|--------|-----|

| `test_text_replace.py` | `internal/gui/text_replace_test.go` |

| `test_project_replace.py` (steps_only) | `internal/gui/features_replace_test.go` |

| `test_run_variables.py` | `internal/player/generators_test.go` |

| `test_scenario_hints.py` (основные кейсы) | `internal/gui/scenario_hints_test.go` |

| `test_save_persistence.py` (save/flush, без draft) | `internal/gui/save_persistence_test.go` |

| `test_recorder_input/checkbox/canvas/date_fields` (normalize unit) | `internal/recorder/normalize_test.go` |

| `test_gherkin_ru.py`, `test_examples.py` | `internal/gherkin/*`, `examples_parity_test.go` |

| `test_player_execute.py` | `internal/player/playwright_executor_test.go` (mock) |

| `test_cli_*.py` | `internal/cli/*_test.go` |

| `test_vanessa_*.py` | `internal/vanessa/*_test.go` |

| `gherkin_template_text` | `frontend/src/lib/featureTemplate.test.ts` |
| `completions_for_line` | `internal/stepcatalog/completions_test.go` |
| Recorder event queue / coalesce | `internal/recorder/event_pipeline_test.go`, `coalesce_test.go`, `live_events_test.go` |
| JUnit golden snapshot | `internal/report/junit_golden_test.go` |
| Monaco completions / hint actions | `frontend/src/lib/gherkinCompletions.test.ts`, `gherkinHintActions.test.ts` |



### Не переносится 1:1



- **Draft autosave** (`save_draft_if_needed`) — в Go/Wails нет аналога; сохранение через `SaveFeature` / вкладки редактора.

- **Qt GUI** (~31 файл) — Vitest + Playwright E2E + `internal/gui/*_test.go`.

- **Recorder HTML fixtures** (Playwright click на странице) — только `live_integration_test.go` smoke; unit-покрытие через `NormalizeSteps`.



## Локальный запуск



```powershell

# Go unit (без frontend/node_modules)

go test ./internal/... ./cmd/...



# Playwright integration (нужен Chromium)

go run github.com/mxschmitt/playwright-go/cmd/playwright@latest install chromium

go test -tags=integration ./internal/recorder/...



# Desktop smoke (нужен собранный scenaria-gui.exe)

wails build -platform windows/amd64

go test -tags=desktop ./internal/wailsapp/...

# или напрямую:

./scripts/desktop-smoke.ps1



# Frontend unit

cd frontend

npm install

npm test



# UI E2E (сборка + preview + Playwright)

npm run test:e2e:install   # один раз

npm run test:e2e

```



## CI (`.github/workflows/ci.yml`)



- `go test ./internal/... ./cmd/...`

- `cd frontend && npm test`, `npm run check`, `npm run build`, `npm run test:e2e`

- Job **integration**: `go test -tags=integration ./internal/recorder/...` (Chromium)

- Job **desktop-smoke** (manual / release): `wails build` + `scripts/desktop-smoke.ps1`



## E2E: мок Wails vs desktop



- **Playwright** (`frontend/e2e/`) — `wails-mock.js` подставляет `window.go` до загрузки в `vite preview`. Быстро, без WebView2.

- **Desktop smoke** — реальный `scenaria-gui.exe` (WebView2). Проверяет, что бинарник стартует; не заменяет Playwright по покрытию UI.


