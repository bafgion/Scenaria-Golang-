# Scenaria Go — Roadmap

Статус: **master** v0.24.0; **Wails IDE** — основной продукт. Python/Qt — снят с поддержки (экспорт в Python сохранён).

## Приоритеты

| # | Направление | Статус |
|---|-------------|--------|
| P0 | Wails IDE | done |
| P0 | Recorder | done |
| **P0** | **Стабильность Web UI (Фаза 8)** | **done** |
| P1 | Monaco / редактор IDE | done (Фаза 6) |
| P1 | Паритет Python IDE | done (Фаза 7) |
| P1 | Allure | done |
| P1 | Portable release | done |
| **P1** | **Monaco hardening (Фаза 9)** | **done** |
| **P2** | **Flaky-run + post-record diff (Фаза 10)** | **done** |
| **P2** | **Cold start + FailedStep + E2E (Фаза 11)** | **done** |
| **P3** | **Lazy workers + flaky E2E (Фаза 12)** | **done** |
| **P0** | **GUI reliability audit (Фаза 13–14)** | **done** |

---

## Фаза 1 — Wails GUI

- [x] Monaco IDE, run/record/Vanessa/OTP/Allure
- [x] Splash, portable `wails build`
- [x] UI shell: menubar, activity bar, explorer, action bar, bottom panel, status bar
- [x] Настройки с вкладками (Интерфейс / Запись / Плагины)
- [x] Пакетный запуск в explorer (Выбор, Ctrl+клик)
- [x] Command palette (Ctrl+Shift+P), recording bar, dirty banner, ресайз панелей
- [x] Welcome: недавние проекты/файлы, примеры; Results/Error панели из run_status.json
- [x] Модалки: Запустить, TestClient, шаги, экспорт (ts/python)

---

## Фаза 2 — Portable

- [x] `scripts/build-portable.ps1` — CLI + Wails + Chromium
- [x] `README-PORTABLE.txt`, `Start-GUI.bat`, CI artifact on tag

---

## Фаза 3 — Allure

- [x] Writer, CLI, GUI
- [x] Screenshot + trace + video attachments

---

## Фаза 5 — Тесты

- [x] `scripts/coverage.ps1` (`-coverpkg=./...`)
- [x] CI soft gate 40% (цель 60%)

---

## Фаза 6 — Monaco и редактор IDE

Цель: перенести типовые IDE-функции в Monaco, убрать дублирование UI, добавить настраиваемый редактор.

### 6.1 Настройки редактора (P1)

- [x] Секция «Редактор» в `SettingsDialog` (вкладка «Интерфейс» или отдельная)
- [x] Поля: `fontSize`, `fontFamily`, `wordWrap`, `minimap`, `lineNumbers`, `tabSize`, `insertSpaces`, `renderWhitespace`
- [x] Опции: `folding`, `stickyScroll`, `autoClosingQuotes`, `formatOnSave`
- [x] Хранение в `AppSettings` / `settings.json` (`editor: { ... }`)
- [x] Применение через `editor.updateOptions()` без пересоздания Monaco
- [x] Светлая тема `scenaria-light` (опционально) + переключатель темы

### 6.2 Find / Replace (P1)

- [x] Встроенный виджет Monaco: `actions.find`, `editor.action.startFindReplaceAction`
- [x] Ctrl+H / Ctrl+F открывают нативный find/replace вместо `FindReplaceDialog` для текущего файла
- [x] Удалить или упростить `FindReplaceDialog.svelte` + обёртки `findNext` / `replaceNext` / `replaceAll` в `MonacoEditor.svelte`
- [x] Оставить `ProjectReplaceDialog` для замены по всему проекту (Wails, не Monaco)

### 6.3 Подсказки и документация шагов (P1)

- [x] `registerHoverProvider` — описание шага, пример, ссылка на справку при наведении на строку
- [x] Данные из каталога шагов (`DescribeEditorLine` / stepcatalog lookup)
- [x] Настройки: вкл/выкл hover, авто-валидация при вводе (debounce уже есть)

### 6.4 Format provider (P2)

- [x] `registerDocumentFormattingEditProvider` — нормализация отступов, схлопывание пустых строк
- [x] Переиспользовать `RefactorNormalizeIndents`, `RefactorCollapseBlankLines` (`FormatFeature` на бэкенде)
- [x] Shift+Alt+F и `formatOnSave` (если включено в настройках)
- [x] Пункты меню «Рефакторинг» оставить как alias

### 6.5 Completions и сниппеты (P2)

- [x] Tab stops в completions: `insertTextRules: InsertAsSnippet`, `${1:selector}` и т.п.
- [x] Улучшить ранжирование: `filterText`, `preselect`, `sortText`
- [x] Свести основной путь ввода к Ctrl+Space; `SnippetPalette` — расширенный поиск по каталогу
- [x] Опционально: trigger characters для кавычек и селекторов

### 6.6 Навигация по сценарию (P2)

- [x] `registerDocumentSymbolProvider` — сценарии, шаги, блоки → Outline / Go to Symbol (Ctrl+Shift+O)
- [x] Breadcrumbs Monaco по структуре feature
- [x] Клик в outline → `gotoLine` (частичная замена панели шагов)
- [x] Настройка: показывать breadcrumbs / outline по умолчанию

### 6.7 Code Lens и запуск из редактора (P3)

- [x] `registerCodeLensProvider` — «▶ Запустить сценарий», «▶ с этой строки» у заголовков сценариев / шагов
- [x] Обработчик lens → существующие `RunFeature` / run-current hotkey
- [x] Настройка: показывать code lens (по умолчанию выкл или только при наведении)
- [x] Частичный запуск с шага (см. Фаза 7.1 — backend `StartStep`/`EndStep`)

### 6.8 Inlay hints (P3)

- [x] `registerInlayHintsProvider` — серым текстом справа: `click → #login`, `fill → "…"`
- [x] Данные из `ParseEditorSteps` (уже есть на бэкенде)
- [x] Настройка: вкл/выкл inlay hints

### 6.9 Превью Gherkin (P3)

- [x] Read-only Monaco вместо `FeaturePreview` + `HighlightFeature` (единая подсветка Monarch)
- [x] Синхронизация темы/шрифта превью с настройками редактора
- [x] Diff-редактор отложен (превью текущего текста достаточно; post-record banner для контекста записи)

### 6.10 Панель шагов и нижняя «Проверка» (P3)

- [x] После outline/code lens — упростить `steps-panel` (свернуть по умолчанию или скрыть при включённом outline)
- [x] Problems-паттерн: клик по issue в «Проверка» уже есть (`gotoEditorLine`) — синхронизировать с Monaco markers
- [x] Единый owner маркеров: ошибки валидации + hints (разные severity/source)

### 6.11 Подсказки сценария — настройки (P3)

- [x] Настройки: показывать hints, после записи, фильтр severity (warning / info)
- [x] Опционально: авто-fix `autoFixable` hints при сохранении
- [x] Post-record banner оставить в Svelte (контекст записи, не редактирования)

### 6.12 Что остаётся в Svelte / Wails (не переносить)

- Вкладки, explorer, запись, браузер, OTP, picker
- Журнал, результаты, Vanessa, плагины, Command Palette приложения
- Замена по проекту, импорт/экспорт, модалки сценария

---

## Фаза 7 — Паритет с Python IDE

**Статус: done** (v0.19.0).

Цель: закрыть реальные пробелы относительно Python/Qt v0.12, не дублируя то, что уже сделано иначе (Monaco, dirty→temp-before-run, hints в редакторе).

Статусы в матрице: **эквивалент** | **иначе (осознанно)** | **нет**.

### 7.1 Запуск с шага / до шага (P1)

- [x] Player: `StartStep` / `EndStep` на leaf-шагах (`gherkin.LeafSteps`, `ApplyStepRange`)
- [x] CLI: `--start-step`, `--end-step`
- [x] Wails: `ResolveRunFromLine` + передача в `Run`
- [x] Code lens «▶ с этой строки» запускает с выбранного шага, не весь сценарий
- [x] Steps panel: контекстное меню «Запустить с шага» / «До шага»
- [x] Dry-run summary с учётом частичного диапазона в логе

### 7.2 Восстановление сессии (P1)

- [x] `open_tabs` + `active_tab` в `settings.json` (как Python)
- [x] Восстановление вкладок при старте IDE
- [x] Draft autosave каждые 30 с (несохранённый текст → `.scenaria/drafts`)
- [x] Восстановление draft при открытии проекта

### 7.3 Настройки записи и запуска (P2)

- [x] Slow-mo (скорость выполнения тестов) в `AppSettings` + пресеты в настройках
- [x] Индикация прогона: playing-bar, подсветка редактора и toolbar
- [x] `scroll_before_click` в recorder script + настройка
- [x] `hover_record_min_ms` — минимальная длительность hover перед записью
- [x] Сохранение опций записи в `AppSettings`

### 7.4 Наборы параметров `.params.json` (P2)

- [x] Загрузка `<stem>.params.json` при outline (sidecar рядом с feature)
- [x] Расширение outline в runner (`ExpandFeatureAtPath` + `LoadScenarioParams`)
- [x] Справка F1: раздел про params в `StepsHelpDialog`

### 7.5 Explorer: папки (P3)

- [x] Контекстное меню папки: «Запустить все .feature»
- [x] «Vanessa: папка…» из explorer
- [x] Batch по папке без ручного Ctrl+клик (выбор для пакетного запуска)

### 7.6 Селекторы и валидация (P3)

- [x] Редактируемый порядок стратегий в настройках (клики / поля ввода; инъекция в recorder/picker)
- [x] Per-step статусы в «Проверка в браузере» (found / missing / warning) + панель «Проверка»

### 7.7 Прочее (P3)

- [x] Обновления приложения: проверка релиза, страница GitHub, скачивание артефакта (installer / portable по ОС)
- [x] Saved browser session (cookies) — захват из открытого браузера в TestClient и повторный запуск с тем же профилем

### Уже эквивалентно (не в scope Фазы 7)

| Python | Go |
|--------|-----|
| Qt Gherkin + Apply | Monaco, dirty banner, temp-before-run |
| FindReplaceDialog | Monaco find/replace |
| GherkinHintsBar | hints markers + `StepsHelpDialog` + `gherkin-hints` |
| StepsStrip edit | редактирование в Monaco + outline + code lens |
| Пользовательские сниппеты | каталог шагов + completions/snippets (Фаза 6.5); отдельное хранение не планируется |
| StepEditorDialog / reorder в панели | Monaco, outline, picker; структурный диалог шага не планируется |
| Теги через диалог | `@tags` в тексте |
| `save_html_reports` global | чекбокс в RunDialog + slow-mo в настройках |
| Browser overlay Qt | in-browser toolbar + browse/record режимы |

---

## Версии

| Версия | Содержание |
|--------|------------|
| **0.14.0** | trace/video Allure (**master**) |
| **0.15.0** | GUI trace/video + scenario catalog + recorder polish (**master**) |
| **0.16.0** | Monaco: настройки редактора, встроенный find/replace, hover шагов (план: Фаза 6.1–6.3) |
| **0.17.0** | Monaco: format, outline, completions/snippets (план: Фаза 6.4–6.6) |
| **0.18.0** | Monaco: code lens, inlay hints, read-only preview, hints settings (Фаза 6.7–6.11) |
| **0.19.0** | Паритет Python: run-from-step, сессия, запись/запуск, params (Фаза 7.1–7.4) |
| **0.20.0** | Стабильность Web UI: player/recorder lifecycle, отчёты при fail, concurrency, security (Фаза 8) — **master** |
| **0.21.0** | Monaco: shortcuts, запись без гонок, preview perf, lifecycle вкладок, масштабирование (Фаза 9) — **master** |
| **0.22.0** | Flaky-run метрики, post-record diff, release CI (Фаза 10) — **master** |
| **0.23.0** | Monaco lazy load, FailedStep в player, E2E outline/diff (Фаза 11) — **master** |
| **0.24.0** | Lazy Monaco workers, E2E flaky-run UI (Фаза 12) — **master** |
| **0.25.0** | GUI reliability: session restore, recorder, shutdown, hotkeys (Фаза 13) — **master** |

---

## Фаза 9 — Monaco Editor hardening (аудит 2026)

**Статус: done** (v0.21.0).

Цель: закрыть пробелы интеграции Monaco + Svelte + Wails после аудита редактора.

### 9.1 P0 — Критические UX

- [x] **`Ctrl+Shift+O`:** привязка `editor.action.quickOutline` в Monaco (не только палитра)
- [x] **Live record:** шаги записи без гонки с ручным вводом (текст из модели + очередь)
- [x] **FeaturePreview:** debounce + `replaceModelText` вместо `setValue` на каждый символ

### 9.2 P1 — Стабильность

- [x] **Providers:** `providerRegistered` для completions и document symbols (HMR/dev)
- [x] **Вкладки:** сохранение cursor/scroll при переключении файлов
- [x] **Hotkeys:** `Shift+Alt+F` и `Ctrl+Shift+O` вне фокуса редактора (explorer и т.д.)
- [x] **`formatOnSave`:** через `formatDocument()` Monaco (единый путь с Shift+Alt+F)

### 9.3 P2 — Масштабирование

- [x] Кэш `parseFeatureSymbols` по `model.getVersionId()` (`featureSymbolCache.ts`)
- [x] Folding ranges для блоков Если/Повторяю/Пока/Для каждого (`gherkinFolding.ts`)
- [x] Lazy-mount превью (80 ms defer); отключение minimap/code lens/inlay на файлах ≥2000 строк
- [x] Monaco отдельный chunk в Vite (`manualChunks` → `monaco.js`)

---

## Следующие шаги (вне закрытых фаз)

Опционально:

- Flaky-run E2E с реальным прогоном (не mock)
- Language workers Monaco (json/css/html) при необходимости

---

## Фаза 13 — GUI reliability audit (2026)

**Статус: done** (v0.25.0).

Цель: закрыть критические пробелы GUI-аудита (Monaco session restore, recorder lifecycle, Wails shutdown, hotkeys).

Оценка до фиксов: **6.5/10**. Цель после фазы: **8/10**.

### 13.1 P0 — Критические

- [x] **Session restore:** синхронизация `activateTab` после mount Monaco (`editor ready`)
- [x] **`record-error`:** подписка во frontend + сброс UI
- [x] **Undo записи:** откат строки в Monaco + `liveRecordStepLines`
- [x] **Recorder race:** `record-step` ждёт `prepareRecordEditorTab`; не сбрасывать map при duplicate `record-started`
- [x] **Wails shutdown:** `OnShutdown` → `CancelRun` + `CloseBrowser`
- [x] **Frontend teardown:** `onDestroy` / `closeProject` / смена проекта
- [x] **Escape:** не перехватывать при открытых Monaco overlays (find/suggest/quick input)

### 13.2 P1 — Usability

- [x] **Post-record banner:** после `record-stopped` (browse→record path)
- [x] **Banner step count:** из редактора, не только с диска
- [x] **`featureSymbolCache`:** ключ с URI модели (нет коллизий между вкладками)
- [x] **`BrowserOverlay`:** показывать при `browserOpen || recording || playing`
- [x] **Hotkey Ctrl+R:** не открывать диалог при активной записи (focus browser)
- [x] **`syncTabContent`:** `monaco.getEditorText()` для активной вкладки
- [x] **Смена проекта:** confirm + reset tabs/browser при `openProjectAt` на другой path
- [x] **Monaco dispose:** `setModel(null)` перед `releaseAll`

### 13.3 P2 — Масштабирование и polish

- [x] Lock editor during `playing` (`readOnly` в Monaco)
- [x] Re-entry guard для `executeRun` / `runPrimary`
- [x] File reload prompt при возврате в окно (`visibilitychange` + `ReadFeature`)
- [x] Large file: gate symbols, folding, hover, completions, code lens ≥2000 строк
- [x] Hotkeys: `Alt+P` пауза записи, `Ctrl+Shift+R` стоп (запись / тест / браузер)

---

## Фаза 14 — Recorder UX polish (v0.25.0)

**Статус: done**.

### 14.1 Целевая вкладка записи

- [x] `recordingTargetPath` фиксируется при `record-started`
- [x] Предупреждение при переключении вкладки во время записи (confirm + смена цели)
- [x] Блокировка «Старт» и закрытия целевой вкладки без паузы
- [x] `applyLiveRecordedStep` возвращает фокус на целевую вкладку

### 14.2 Results panel

- [x] Двойной клик по строке → открыть feature (как в истории запусков)

---

## Фаза 12 — Lazy workers + flaky E2E (v0.24.0)

**Статус: done** (v0.24.0).

### 12.1 Отложенные Monaco workers

- [x] `ensureMonacoEnvironment()` — dynamic import `editor.worker` при первом `preloadMonacoEditor`
- [x] Убран sync import `monaco-env` из `main.ts` (cold start без worker bundle)

### 12.2 E2E flaky-run UI

- [x] Mock `?e2e=flaky-run` — `ListRunResults` + `FlakyMetrics` с flaky-сценарием
- [x] E2E: история запусков (фильтр Flaky) + бейдж в панели «Результаты»

---

## Фаза 11 — Cold start, FailedStep, E2E (v0.23.0)

**Статус: done** (v0.23.0).

### 11.1 Динамический import Monaco

- [x] `import('monaco-editor')` в `appBootstrap` — отдельный chunk, не блокирует main bundle
- [x] Splash без await Monaco; `prefetchMonacoEditor()` после показа shell
- [x] `gherkinHintActions` — только type-import Monaco

### 11.2 FailedStep в player

- [x] `RunContext.markFailedLeafStep` / `FailedLeafStep` (0-based leaf index)
- [x] `ScenarioResult.FailedStep` → `run_status.json` через CLI
- [x] Step-flaky метрики получают данные из реальных прогонов

### 11.3 E2E

- [x] `Ctrl+Shift+O` → quick outline widget
- [x] Post-record diff: режим `?e2e=post-record-diff`, banner + diff dialog

---

## Фаза 10 — Flaky-run, post-record diff, release CI (v0.22.0)

**Статус: done** (v0.22.0).

### 10.1 Метрики flaky-run

- [x] `runstatus.FlakyStats` — сценарии с чередованием pass/fail; шаги с ≥2 падениями
- [x] API `FlakyMetrics` + `failed_step` в `ListRunResults`
- [x] UI: фильтр «Flaky» в истории; бейджи в Results / Run history

### 10.2 Monaco diff после записи

- [x] Baseline текста при `record-started`
- [x] Post-record banner + кнопка «Сравнить»
- [x] `PostRecordDiffDialog` — Monaco `createDiffEditor` (до / после)

### 10.3 Release CI

- [x] `.github/workflows/release.yml` — tag `v*` → portable zip + installer + `latest.json` (уже было; задокументировано в ROADMAP)

---

## Фаза 8 — Стабильность Web UI (аудит 2026)

**Статус: done** (v0.20.0).

Цель: закрыть критические пробелы player / recorder / отчётов / Wails для production Web UI automation.

Оценка до фиксов: **5/10**. После фазы: **7+/10**.

### 8.1 P0 — Критические (player / CLI / recorder)

- [x] **Runner + CLI:** partial `ExecutionResult` при падении; HTML/JUnit/Allure пишутся до return error
- [x] **`browserSession`:** сериализация доступа к `page`/`closed` (mutex на `executeAction`)
- [x] **`waitForLocator`:** drain goroutine при `ctx.Done()` (не оставлять зависший `WaitFor`)
- [x] **`RecordLive`:** session generation ID — не обнулять `liveSession`/`recordCancel` чужим goroutine
- [x] **Recorded steps:** все мутации `*steps` под `LiveSession.mu` (poll-loop + undo)

### 8.2 P1 — Важные (reliability / security)

- [x] **Signal handling:** SIGINT/SIGTERM → cancel run context (CLI); GUI `CancelRun` + `RunRunContext`
- [x] **`runstatus`:** `WritableScenariaDir` вместо жёсткого `project/.scenaria`
- [x] **Record timeout:** idle-only (не wall-clock `idle+30` на всю сессию)
- [x] **`PickSelector`:** cancel через `recordCtx` / `recordCancel`
- [x] **Picker bindings:** per-context или reset при close browser
- [x] **`captureCLI`:** mutex на stdout (без гонок при параллельных GUI вызовах)
- [x] **OTP channels:** очищать после use (`wailsapp/app.go`)
- [x] **`UrlsMatch`:** опционально query/fragment для SPA
- [x] **Navigation:** configurable `waitUntil` (не только `domcontentloaded`)
- [x] **OTP / download / press:** thread `ctx` через все блокирующие waits
- [x] **Path confinement:** `Output`/`AppendTo` в recorder — только внутри project root
- [x] **Gherkin sanitize:** escape `\n` в recorded step text
- [x] **Allure:** очистка stale results + уникальные timestamps per scenario
- [x] **`ServeAllure`:** tracking PID, не плодить JVM
- [x] **CLI validate:** flag-first parsing (`validate --no-browser ./features`)
- [x] **CLI run:** dedupe discovered `.feature` paths
- [x] **`writeJSON` / reports:** `MkdirAll` parent dir
- [x] **Wails async:** `ctx != nil` guard в `StartRecord`/`OpenBrowser`/`StartVanessaRun`
- [x] **`buildRunner` error:** показывать resolved engine name

### 8.3 P2 — Улучшения / flakiness / DX

- [x] **Browser pool** (reuse context per worker) — снижение RAM при `--workers N`
- [x] **`for_each`:** re-query locators per iteration
- [x] **`downloadByClick` / `upload`:** chained locators как у click/fill
- [x] **`assert-hidden` / `wait-hidden`:** проверять все matches, не только `.First()`
- [x] **JUnit:** статус `broken` → failures
- [x] **`WriteTempFeature`:** cleanup temp dirs
- [x] **Download artifacts:** teardown `.scenaria/downloads/run-*`
- [x] **Recorder trust:** document hostile-origin risk; validate picker binding origin
- [x] **Linux paths:** единый app-data root (settings + artifacts)
- [x] **Structured errors:** `ExecutionFailure` с partial result (typed)
- [x] **Observability:** slog `run_id` при старте прогона, debug-логи retry
- [x] **CLI help:** упомянуть `--html`
- [x] **`RunInit`:** не глотать ошибки scaffold

### 8.4 Уже сделано (аудит follow-up)

- [x] `watchContext` + `pw.Stop()` при cancel
- [x] Parallel fail-fast (`cancel()` on first failure)
- [x] `waitForURL` polling с учётом ctx deadline
- [x] Chained locators + retry для большинства actions
- [x] Writable artifacts fallback (`paths.WritableScenariaDir`)
- [x] Monaco: Ctrl+Z (без `setValue` на tab switch), deferred hint fix
- [x] Integration: `parallel_cancel_integration_test.go`
- [x] `docs/SELECTORS.md`, placeholder cycle detection

### Матрица приоритетов фазы 8

| Область | P0 | P1 | P2 |
|---------|----|----|-----|
| Player / runner | partial results, session lock, wait drain | signal, UrlsMatch, ctx OTP | pool, for_each, selectors |
| Recorder | session gen, steps mutex | picker cancel, path confine | trust doc, sanitize |
| Reports / CLI | reports on failure | runstatus, dedupe, validate flags | JUnit broken, help |
| Wails / GUI | — | OTP, captureCLI, ctx guard | temp cleanup, Allure PID |
