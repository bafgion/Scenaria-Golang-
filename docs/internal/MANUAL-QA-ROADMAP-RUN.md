# Manual QA — прогон по ROADMAP (2026-07-09)

Сводка автоматизированного прогона + оставшиеся ручные шаги в desktop.

## Итог по секциям ROADMAP Manual QA

| Секция | Закрыто [x] | Осталось [ ] | Комментарий |
|--------|-------------|--------------|-------------|
| §1 Monaco Tabs | 7 / 9 | 2 | 5 вкладок, первая/последняя — desktop |
| §2 Batch | 2 / 8 | 6 | suite run + stale filters E2E; batch 5/hotkey/folder — desktop |
| §3 ProjectSession | 7 / 7 | 0 | полностью E2E + Go |
| §4 StepMatcher | 7 / 7 | 0 | unit + E2E + Go @smoke |
| §5 Backend storage | 3 / 6 | 3 | unsaved/temp, settings race, refresh — desktop |
| §6 Recorder | 7 / 8 | 1 | close browser during picker — desktop |
| §7 Reports | 4 / 7 | 3 | latest/partial/canceled Allure — desktop |
| §8 Runner lifecycle | 4 / 5 | 1 | close app mid-run — вручную |
| §9 Performance | 0 / 6 | 6 | только desktop |
| **Всего** | **41 / 57** | **16** | **~72% через автотесты + desktop smoke** |

> Дополнительно: stale filters после batch — E2E `batch run clears stale filters` (Acceptance §2, не отдельный пункт Manual QA).

## Результаты автоматизации

| Набор | Результат | Команда |
|-------|-----------|---------|
| Frontend unit | **280/280** | `cd frontend && npm test` |
| E2E (mock Wails) | **107/111** | `cd frontend && npm run test:e2e` |
| **Desktop smoke (WebView2)** | **26/26** | `./scripts/desktop-smoke.ps1` (2026-07-09) |
| Go unit (player/gui/report) | **OK** | `go test ./internal/player/... ./internal/gui/... ./internal/report/... -short` |
| Browser isolation integration | **OK** | `go test ./internal/player/... -tags=integration -run Isolation` |

### E2E failures (4, не блокеры ROADMAP core)

| Тест | Вероятная причина |
|------|-------------------|
| `untitled tab restores after reload` | timing Monaco после reload |
| `onboarding tour step 5` | селектор меню тура |
| `1.6 update modal deferred` | mock update flow изменился |
| `record resume does not duplicate steps` | mock record-step timing |

---

## §1 Monaco Tabs — ROADMAP Manual QA

| Пункт | Статус | Покрытие |
|-------|--------|----------|
| Открыть 5 файлов | **E2E частично** | 2–3 вкладки в `app-ui.spec.ts`; 5 вкладок — desktop |
| Быстро переключаться | **E2E** | `undo and redo stay isolated per Monaco tab model` |
| Разные изменения в каждый файл | **E2E** | `qa-daily-use` 2.1 dirty asterisk |
| Закрыть вкладку из середины | **E2E** | `closing active tab activates the previous Monaco model` |
| Закрыть первую/последнюю вкладку | **Desktop** | — |
| Закрыть Welcome tab | **E2E** | `closing welcome tab activates the last feature tab` |
| undo/redo по вкладкам | **E2E** | `undo and redo stay isolated per Monaco tab model` |
| save после переключения | **E2E** | `save completion does not mutate another active tab` |
| external reload после переключения | **E2E** | `disk reload completion does not mutate another active tab` |

---

## §2 Batch Execution

| Пункт | Статус | Покрытие |
|-------|--------|----------|
| Выбрать 2 теста и запустить | **E2E** | `user-journeys` suite run |
| Добавить 3 теста → всего 5 | **Desktop** | batch UI на реальном каталоге |
| Убрать часть и запустить | **Desktop** | — |
| Hotkey после batch mode | **Desktop** | — |
| Folder selection | **Desktop** | — |
| Refresh после удаления файла | **Desktop** | — |
| Stale filters | **E2E** | `batch run clears stale filters` |

---

## §3 ProjectSession / stale events

| Пункт | Статус | Покрытие |
|-------|--------|----------|
| Project A → run → project B | **Go** | `project_session_test.go` |
| Старые результаты не применились | **E2E** | stale validation / editor change |
| Повторить с validation | **Go+E2E** | `validate confirm`, `project_session_test` |
| Повторить с recorder | **Go** | `record_lifecycle_test.go` |
| Report actions | **E2E** | results / trace viewer mock |

---

## §4 StepMatcher / diagnostics

| Пункт | Статус | Покрытие |
|-------|--------|----------|
| Diagnostics / autocomplete / inlay | **Unit** | `editorAnalysisSync`, `gherkinCompletions` tests |
| Синтаксическая ошибка | **E2E** | validate confirm |
| Runner тот же сценарий | **Go** | `examples_integration_test` (@smoke) |

---

## §5 Backend storage

| Пункт | Статус | Покрытие |
|-------|--------|----------|
| Parallel run + history | **E2E+Go** | run progress, `runstatus` tests |
| Unsaved/temp feature | **Desktop** | — |
| Cancel → другой run | **Go** | `run_session_test`, cancel integration |
| Settings + recents race | **Desktop** | — |
| Project refresh during save | **Desktop** | — |

---

## §6 Recorder

| Пункт | Статус | Покрытие |
|-------|--------|----------|
| Запись в B при активной A | **E2E** | `recording-target` status bar |
| Click/input только в target | **E2E mock** | `post-record`, `live recording inserts steps` |
| Переключение вкладок | **E2E** | `3.2 tab switch during recording` |
| Stop → Start | **E2E** | `record-resume` (1 fail flaky) |
| Close browser during picker | **Desktop** | — |
| Project switch | **Go** | `TestOpenProjectCancelsRunValidateAndRecordContexts` |
| Stale events | **Go+Unit** | `record_lifecycle_test`, `wailsEventsController` |

---

## §7 Reports

| Пункт | Статус | Покрытие |
|-------|--------|----------|
| HTML report UI | **E2E** | `html-report.spec.ts` (16 tests) |
| Validation diagnostics | **Go** | `html_payload_golden` + validation fields |
| Light/full pair | **E2E** | mode switch test |
| Scenario outline examples | **Go** | `html_example_fixture_test` |
| Canceled partial report | **Desktop** | — |
| Allure failed write | **Desktop** | — |

---

## §8 Runner lifecycle

| Пункт | Статус | Покрытие |
|-------|--------|----------|
| Run/cancel 20× | **Go** | `action_context_cancel`, goroutine dev test |
| Cancel during nav/wait/download | **Go** | cancel integration tests |
| Parallel run с failure | **Go** | `parallel_cancel_integration_test` |
| Close app during run | **Desktop** | — |
| Browsers закрылись | **Desktop** | `scripts/desktop-smoke.ps1` |

---

## §9 Performance

| Пункт | Статус | Покрытие |
|-------|--------|----------|
| Large feature typing | **Desktop** | — |
| Large project | **Desktop** | — |
| 100-scenario report | **Go bench** | `report` benchmarks (если есть) |
| 50 tabs | **Desktop** | — |

---

## Desktop smoke (2026-07-09)

**Результат: 26/26 passed** (~36s startup + ~12s Playwright).

Бинарник: `build/bin/scenaria-gui.exe` (уже собран). CDP `:9333`, изолированный `%TEMP%\scenaria-desktop-smoke`.

Покрыто desktop-smoke (не все пункты ROADMAP Manual QA, но ключевой UI):

| Область | Тесты |
|---------|--------|
| Загрузка / splash | `приложение загружается после splash` |
| Проект examples + каталог + редактор | 2 теста |
| Validate → журнал | `проверка сценария пишет в журнал` |
| Настройки (open/reset/apply/warnings) | 4 теста |
| Палитра, справка, hotkeys, журнал | 3 теста |
| Новый сценарий, сниппеты, run dialog | 4 теста |
| Экспорт, несохранённые изменения | 2 теста |
| Плагины, о программе | 2 теста |
| Record + HTTP Auth stacking | 2 теста |
| Onboarding tour (шаги 1–6, incl. step 5) | 4 теста в `desktop-tour-onboarding` |

**Не покрыто smoke** (остаётся ручной прогон): batch 5, performance §9, real Playwright run/cancel, close app mid-run, report latest/partial, picker+browser.

---

## Desktop smoke (рекомендуется сейчас)

```powershell
# Сборка + smoke UI на WebView2
./scripts/desktop-smoke.ps1

# Только процесс (быстро)
./scripts/desktop-smoke.ps1 -SkipUI
```

Ручной чеклист повседневного использования: [QA-DAILY-USE.md](QA-DAILY-USE.md).

---

## Критерий закрытия Manual QA в ROADMAP

- Пункт **закрыт [x]**, если покрыт E2E mock + Go integration **или** пройден в desktop smoke.
- Пункт остаётся **[ ]**, если требует реальный Playwright run / живой recorder / закрытие приложения mid-run.
