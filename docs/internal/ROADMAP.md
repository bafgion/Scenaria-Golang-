# Scenaria Go — Roadmap

Статус: **master** v0.28.1 (в работе); **Wails IDE** — основной продукт. Python/Qt — снят с поддержки (экспорт в Python сохранён).

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
| **P0** | **Daily-use QA audit (Фаза 15)** | **done** |
| **P1** | **Interactive HTML Report (Фаза 16)** | **done** |
| **P0** | **Code audit Web UI stability (Фаза 17)** | **done** |

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
| **0.25.0** | GUI reliability: session restore, recorder, shutdown, hotkeys (Фаза 13–14) — **master** |
| **0.26.0** | Daily-use QA: run progress, trace viewer, editor races, onboarding (Фаза 15) — **master** |
| **0.27.0** | Onboarding tour, live browser reuse, bilingual docs, CI stability — **master** |
| **0.28.0** | Phase 17 Web UI stability, HTML reports, failed-run reports, E2E green — **master** |
| **0.28.1** | First-run browser/OTP UX, deferred browser events, Phase 16 closure — **master** |

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

**Активно:** backlog Фазы 16 (полный HAR timeline per-step) + опциональные E2E.

Опционально (backlog):

- Flaky-run E2E с реальным прогоном (не mock)
- Language workers Monaco (json/css/html) при необходимости
- HAR / full network timeline per step в HTML report (live + trace: snippet на каждый шаг — **done** v0.28.1)

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

## Фаза 15 — Daily-use QA audit (v0.26.0)

**Статус: done** (v0.26.0).

**Источник:** симуляция ежедневного использования (QA Automation + UX), 6 пользовательских сценариев: onboarding, Monaco, live recording, run, results, settings.

**Оценка до фиксов: 6.8/10**. Цель после фазы: **8/10**.

### Матрица сценариев (до фиксов)

| Сценарий | Интуитивность | Стабильность | Приятность 8h/day |
|----------|---------------|--------------|-------------------|
| 1. Onboarding | 6/10 | 7/10 | 5/10 |
| 2. Monaco | 7.5/10 | 6/10 | 7/10 |
| 3. Recording | 7/10 | 7.5/10 | 6.5/10 |
| 4. Run | 5.5/10 | 7/10 | 5/10 |
| 5. Results | 6/10 | 8/10 | 5.5/10 |
| 6. Settings | 7/10 | 8/10 | 6.5/10 |

---

### 15.1 P0 — Критические (ежедневные блокеры)

#### Запуск тестов — нет реального прогресса

- [x] **Live run progress:** события из `internal/player` (сценарий / шаг / файл) → Wails events `run-progress`
- [x] **Playing bar:** числитель/знаменатель (N/M сценариев, текущий файл) + progress bar по `--run-progress`
- [x] **Журнал:** стрим stdout во время `captureCLIStream` → `run-log-line`
- [x] **Results panel:** инкрементальное обновление строк во время suite run (как `VanessaMonitorPanel`)

#### Результаты — trace без viewer

- [x] **Trace viewer:** кнопка «Trace viewer» → `playwright show-trace` на последнем ZIP
- [x] Fallback: подсказка в UI, если `playwright` CLI недоступен

#### Recorder / Stop — перегруженная кнопка

- [x] **Контекстная подпись Stop:** «Стоп тест» / «Стоп запись» / «Закрыть браузер» (toolbar + title)
- [x] Confirm при закрытии браузера, если есть несохранённый сценарий
- [x] Hotkey `Ctrl+Shift+R` — то же поведение, что и кнопка (контекст через `stopRecord`)

#### Monaco — race при переключении вкладок

- [x] **Dirty race:** `activeTab` обновлять до `activateTab`
- [x] **Guard `attachModel`:** silent attach при tab switch (без фантомного `*`)

---

### 15.2 P1 — Onboarding (сценарий 1)

- [x] **`checklistDismissed`:** передать в `WelcomePanel`, persist в `settings.json`, кнопка «Скрыть чеклист»
- [x] **Шаг 2 чеклиста:** ✓ только при `recording || browserOpen`
- [x] **`welcomePlayedSuccess`:** persist в settings
- [x] **`startURL`:** persist в settings
- [x] **Быстрый старт без проекта:** блок + переход к «Открыть проект»
- [x] **Запись без проекта:** `beginRecord` требует открытый проект
- [x] **«Открыть примеры»:** после open — подсказка «выберите сценарий в каталоге» или авто-фокус каталога

---

### 15.3 P1 — Monaco и редактор (сценарий 2)

- [x] **`saveFeatureAs`:** `monaco.getEditorText()`
- [x] **Failed step → go to line:** `FailedStepLine` + кнопки в Results / Error panel
- [x] **Large file banner:** status bar / toast при `≥2000` строк («упрощённый режим: без outline/folding/hover»)
- [x] **Find UX:** в Hotkeys dialog и подсказках — явно Ctrl+F (find) vs Ctrl+H (find+replace); опционально унифицировать
- [x] **Record undo vs Ctrl+Z:** подсказка в recording bar / F1 — «Отменить шаг» ≠ editor undo
- [x] **Dirty на неактивных вкладках:** опционально баннер «N несохранённых вкладок» или список в Command Palette
- [x] **Auto-fix при сохранении:** лог / toast при `runScenarioHintsAutoFix` («исправлено N подсказок»)

---

### 15.4 P1 — Live Recording (сценарий 3)

- [x] **Recording target chip:** в status bar — «Запись → `file.feature`»
- [x] **Idle timeout toast:** `record-stopped` reason `idle` + баннер в журнале
- [x] **Output path по умолчанию:** активный таб / `recordingTargetPath`, не `recorded.feature` в корне
- [x] **Tab switch confirm:** опция «Больше не спрашивать» (session или settings)
- [x] **Poll vs pause desync:** `syncBrowserStateFromBackend` не перетирать `recordPaused` сразу после user toggle (debounce / ignore stale)
- [x] **Headless toggle в recording bar:** confirm перед relaunch браузера mid-session

---

### 15.5 P1 — Запуск тестов (сценарий 4)

- [x] **Cancel в playing bar:** кнопка «Отмена» + статус «Останавливаем…»
- [x] **`runPrimary` / Ctrl+Enter:** summary последних опций в toolbar или status bar (headed, HTML, workers…)
- [x] **Первый Ctrl+Enter:** опционально открывать RunDialog, если `lastRun` ещё не задан явно
- [x] **`rerunFailed`:** по полному ключу `path::scenario`, не только `path` (перезапуск одного сценария)
- [x] **`readOnly` при Vanessa run:** единый флаг «automation active» (`playing || vanessaRunning`)

---

### 15.6 P1 — Settings (сценарий 6)

- [x] **«Сбросить по умолчанию»:** кнопка на вкладке / глобально в SettingsDialog
- [x] **`navWaitUntil` в UI:** выпадающий список (load / domcontentloaded / networkidle…) — ключ уже в `AppSettings`
- [x] **`sidebarWidth` drift:** единый owner (только `settings.json` **или** только `localStorage` layout)
- [x] **Recording bar vs Settings:** после OK в Settings — синхронизировать toggles в recording bar
- [x] **Валидация workers / slowMo:** предупреждение при экстремальных значениях (16 workers + slowMo 5000)

---

### 15.7 P2 — Polish и согласованность

#### Onboarding / проект

- [x] Session restore: toast при несуществующем `sessionProject` («проект не найден: …»)
- [x] Max recents: увеличить с 6 или настраиваемо в Settings

#### Monaco

- [x] Dry-run vs real run: единая политика `readOnly` (или явная подпись «dry-run — редактор доступен»)
- [x] Закрытие вкладки: предупреждение, если tab — `recordingTargetPath` (даже на паузе — опционально)

#### Recording

- [x] `record-started` pre-emit: не показывать «● Идёт запись» до готовности браузера (или spinner)
- [x] Picker без паузы: advanced setting «разрешить picker во время записи» (default off)

#### Run

- [x] Parallel fail-fast: настройка в RunDialog / Settings (продолжать все / остановить при первом fail)
- [x] Suite run: имя текущего сценария в playing bar, не только «N сценариев»

#### Results

- [x] Allure: проверка `allure` в PATH + ссылка «Как установить» в Settings / Results
- [x] `ServeAllure` повторный вызов: показать URL / кнопка «Открыть снова»
- [x] HTML report: опция сохранять с timestamp (`report-YYYYMMDD-HHMM.html`)

#### Settings

- [x] Кнопка **Apply** без закрытия диалога (OK остаётся)
- [x] Hotkeys: не перехватывать Ctrl+S при открытом Settings (или явно disabled state)

---

### 15.8 P3 — Nice to have (backlog внутри фазы)

- [x] Мастер «Новый проект» (папка + `.scenaria` + шаблон feature)
- [x] System theme (follow OS) в editor settings
- [x] Flaky badge → «Запустить 3×» из Results
- [x] Update modal: не показывать поверх splash / первого onboarding
- [x] E2E: settings reset defaults, live progress mock, trace viewer smoke
- [x] Документ `docs/QA-DAILY-USE.md` — чеклист ручного регресса по 6 сценариям

---

### 15.9 Порядок реализации (рекомендуемый)

| Sprint | Фокус | Ключевые пункты |
|--------|-------|-----------------|
| **15.1a** | Run visibility | 15.1 live progress + journal stream |
| **15.1b** | Debug loop | 15.1 trace viewer + 15.3 failed_step goto line |
| **15.1c** | Editor trust | 15.1 dirty race + 15.3 saveFeatureAs |
| **15.2** | Recorder clarity | 15.1 split Stop + 15.4 target chip + idle toast |
| **15.3** | Onboarding | 15.2 checklist + quick start guard |
| **15.4** | Settings | 15.6 reset defaults + navWaitUntil + sidebarWidth |
| **15.5** | Polish | 15.7–15.8 по остатку |

---

### 15.10 Критерии приёмки фазы

- [x] Suite из ≥10 сценариев: виден текущий файл и прогресс N/M в playing bar
- [x] Упавший прогон: trace открывается из IDE одной кнопкой
- [x] Переключение 5 вкладок с правками: нет фантомного `*` без редактирования
- [x] Stop: пользователь понимает, что остановится, без чтения журнала
- [x] Новый пользователь: «Быстрый старт» не оставляет в ловушке «записал — не могу запустить»
- [x] Settings: сброс к defaults + `navWaitUntil` в UI
- [x] E2E: ≥2 новых теста (progress mock, onboarding guard или trace button)

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

---

## Фаза 16 — Interactive HTML Report (Mini Trace Viewer)

**Статус: done** (v0.28.0). Остаток — backlog (HAR per-step, optional per-step screenshots toggle).

Цель: превратить `report.html` в практичный инструмент для QA — timeline шагов, inspector, сравнение с историей, интеграция с Playwright Trace.

### 16.1 MVP — структура и timeline (P1) — **done (частично)**

- [x] Embedded JSON payload в одном HTML-файле
- [x] Трёхпанельный layout: сценарии | timeline | inspector
- [x] Фильтры: только упавшие, тег, длительность, поиск
- [x] Failure highlight: авто-выбор первого failed сценария и шага
- [x] Per-step записи при browser-run (`StepRecord`: timing, selector, error)
- [x] Dry-run timeline из плана (leaf steps + inferred status)
- [x] Inspector: Gherkin, selector, ошибка, скриншот (full mode)
- [x] Copy selector / Copy Gherkin / Copy re-run / Copy trace cmd
- [x] Кнопка «Open Playwright Trace» → trace.playwright.dev + trace zip рядом с отчётом
- [x] Советы по шагу (data-testid, timeout, strict mode)
- [x] Шапка: passed/failed/skipped, flaky count, slowest step
- [x] Сравнение с `run_status.json` (последний прогон)
- [x] Export failed steps → `.feature`
- [x] Responsive (mobile sidebar + inspector drawer)
- [x] `HTMLOptions.LightMode` / `RunRequest.HTMLLightMode` (без бинарных артефактов)

### 16.2 Сбор данных — расширение (P1)

- [x] Network snippet на шаг (failed request / HTTP ≥400) из live run и Playwright tracing
- [x] Network snippet последнего failed request на упавшем шаге
- [x] Network failures из trace.network (HTTP ≥400) на каждый шаг по offset
- [x] Скриншот на каждый шаг (опционально, full mode)
- [x] Скриншот на каждый шаг при full HTML (не light mode)
- [x] Скриншот viewport на упавшем шаге (full mode)
- [ ] DOM snapshot / accessibility tree для failed step
- [x] DOM snapshot + a11y tree (упрощённый) на failed step
- [x] Page context (title + URL) на failed step
- [x] Example index в payload и заголовке timeline
- [x] Длительность сценария end-to-end (wall clock)

### 16.3 Trace Viewer — углубление (P1)

- [x] Встроенный мини-viewer (timeline + action log без полного PW UI)
- [x] Trace tab: rail + события + шаги (mini-viewer в центральной панели)
- [x] Action log tab (таблица шагов, sync с timeline/inspector)
- [x] Drag-and-drop trace в отчёт (offline парсинг .zip → Trace tab; fallback trace.playwright.dev)
- [x] Drop-zone для trace .zip в inspector
- [ ] Ссылка `npx playwright show-trace` с авто-open из IDE (кнопка в GUI)
- [x] Авто-open trace viewer после GUI-прогона с падениями (если включён trace)
- [x] Кнопка «Trace in IDE» через report bridge → `OpenTrace`
- [x] Trace только для failed сценариев (zip при падении, discard после pass)
- [x] Синхронизация шага ↔ trace (timeline из zip + offset/bridge)
- [x] Trace timeline из zip + клик → выбор ближайшего шага
- [x] Trace at step: offset в bridge + seek hint в IDE/journal

### 16.4 Размер и режимы (P2)

- [x] Сжатие скриншотов (webp / jpeg quality)
- [x] Toggle Full/Light в диалоге «Запустить» (GUI)
- [x] CLI `--html-light`
- [ ] Лимит размера embedded JSON; вынос крупных trace в `traces/` only (zip на диске + trim `trace_events` в JSON — **done**; HAR в embed — backlog)
- [x] Лимит embedded JSON (4 MiB) с `artifacts_trimmed` flag
- [x] Поэтапный trim: screenshots → DOM/a11y → trace_events (zip остаётся в `traces/`)
- [x] Дедупликация одинаковых скриншотов между шагами

### 16.5 История и flaky (P2)

- [ ] Diff шагов с N предыдущими прогонами (таблица regressions)
- [x] Diff шагов с N предыдущими прогонами (`run_diff` в inspector)
- [x] Regression hints (new_failure / step_changed / still_failing)
- [x] Таблица последних 5 прогонов на сценарий в inspector
- [x] Подсветка flaky шагов на timeline (badge + tint для passed flaky)
- [x] Flaky badge на timeline (из run_status)
- [x] Sparkline длительности сценария по истории (run_status)
- [x] Sparkline длительности шага по истории (`step_durations` в run_status)
- [x] Импорт внешнего `run_summary` для CI сравнения

### 16.6 QA workflow (P2)

- [ ] Re-run scenario из IDE по клику в отчёте (deep link / custom protocol) — [x] HTTP bridge `/rerun` + кнопка в viewer
- [x] HTTP bridge localhost: Open in IDE / Re-run in IDE из HTML-отчёта
- [ ] Jump to line в Monaco из inspector — [x] bridge `/goto` + `gotoReportStep` в IDE
- [x] Jump to line через bridge (кнопка «Open in Scenaria IDE»)
- [x] Печать / PDF-friendly layout (`@media print`)
- [x] i18n отчёта ru/en (`locale` в payload, строки в viewer)
- [x] Клавиатура j/k и ↑↓ для навигации по шагам

### 16.7 Тесты и CI (P1)

- [x] Unit: payload builder, step status inference, WriteHTML smoke
- [x] Golden JSON snapshot (структура payload)
- [ ] E2E: прогон example → открыть report.html → клик по failed step — [x] `html-report.spec.ts` (fixture `example-report.html`)
- [x] E2E: fixture из `examples/01-pervaya-proverka.feature` → `example-report.html`
- [x] E2E: fixture `sample.html` + Playwright (timeline, action log, filter)
- [x] Trace offset hint в inspector (кумулятивная длительность шагов)

---

## Фаза 17 — Code audit: стабильность Web UI / Playwright / Wails (P0)

Аудит senior Go (июнь 2026). Оценка до фиксов: **5/10** для production Web UI. Цель фазы — закрыть критические дыры и зафиксировать остаток.

### 17.1 Критические — исправлено

- [x] **Deadlock `browser_pool.release` + `Close()`** при abort/cancel (`internal/player/browser_pool.go`) — слот всегда возвращается в канал; тест `browser_pool_release_test.go`
- [x] **Live browser vs прогон:** `HoldForTestRun` блокирует poll записи; `CloseBrowser` не закрывает браузер пока `TestRunHeld()`; forced close только в `Shutdown`
- [x] **EPIPE / Ctrl+C:** graceful shutdown ждёт `activePlaywright` (8 с), затем `StopAllureServe` + `closeBrowserForced`; мягкий `abortRun` без рвения pipe; drain 75 ms перед `pw.Stop`
- [x] **Report bridge без auth:** per-session token `X-Scenaria-Bridge-Token`, встраивается в HTML payload (`bridge_token`)
- [x] **Произвольные пути в bridge:** `confineFeaturePath` / `confineArtifactPath` перед emit; trace/open — `ConfineToProjectRoot`
- [x] **`ReadFeature` / `SaveFeature` / `DuplicateFeature`:** confinement через `path_guard.go` (проект, temp run dirs, `%TEMP%`)

### 17.2 Средний приоритет — исправлено

- [x] **`session.closed` data race** → `atomic.Bool` + `isClosed()` / `setClosed()`
- [x] **`OnRequestFailed` vs `executeAction` deadlock** — отдельный `networkMu` для `lastNetworkFail`
- [x] **`captureTraceZIP` / `traceStopped`** — чтение/запись под `session.mu`
- [x] **`watchContext` на sequential runner** — при `openSession` для non-attached прогонов
- [x] **`HoldForTestRun` counter** — floor at 0 при release
- [x] **`EachRecordedLine` slice race** — копия steps под lock
- [x] **CPU spin в `live_session` poll** — `time.After(50ms)` вместо `select default`
- [x] **CLI: report error маскирует execute error** — `errors.Join`
- [x] **CLI: `recordRunStatus` engine** — `resolveRunEngine` вместо пустого `opts.engine`
- [x] **GUI `CloseAfterRun`** для отдельного тестового браузера (не live reuse)
- [x] **`pageGoto` с ctx** — отмена навигации без удержания mutex

### 17.3 Средний приоритет — открыто

- [x] **`emitRunProgress` под mutex** в parallel runner — progress вне lock
- [x] **Orphan goroutines** Playwright при cancel — `pendingAsync` + `drainPendingAsync` перед `pw.Stop` / `session.close`
- [x] **`startPlaywright` не отменяется** по ctx при pool create — cancel + stop orphan driver
- [x] **Vanessa run goroutine** не tracked в shutdown — `activeBackground` WaitGroup
- [x] **`WriteTempFeature` dirs** — cleanup в `Shutdown` (не только после `Run`)
- [x] **Plugin install** — лимит 100 MiB на download
- [x] **HTTP auth passwords** во frontend — `hasPassword` вместо пароля; preserve on save
- [x] **Wails bindings:** `BeginRecordingCapture` → `bool`; `ResolveRunFromLine` → error
- [x] **`main.go` exit code** при ошибке `wails.Run` — `os.Exit(1)`
- [x] **`go test -race`** в CI на player/gui/recorder

### 17.3b Средний приоритет — осталось

- [x] **Wails bindings:** `EventBindingTypes` — DTO уже в `models.ts` (`RunProgressEvent`, `UpdateProgressDTO`)
- [x] **Chaos-тест** pool cancel при `workers>1` — `browser_pool_chaos_test.go`
- [x] **Parallel cancel propagation** — `runCtx` на pool + `abortActiveSessions` + `failFastParallelCancel`
- [x] **Step retry policy** — `runWithRetries` + `isRetryableStepError` (transient errors, expanded actions)
- [x] **Silent browser cleanup** — `closeBrowserResource` вместо `_ = Close()` в session/trace/reset

### 17.4 Flaky / UX Web UI — открыто

- [x] Изолированный browser context по умолчанию — `reuseLiveBrowser` opt-in в диалоге «Запустить»
- [x] `networkidle` / per-project nav wait — `nav_wait_until` в `.scenaria/project.json`
- [x] Retries на `goto` — до 3 попыток на timeout/net errors
- [x] HTML full mode: скриншоты opt-in — light по умолчанию (GUI + CLI `--html` без `--html-full`)

### 17.5 Архитектура (backlog)

- [x] Единый `PlaywrightRuntime` — `internal/playwrightrt` (ref-count, player/recorder/selector + shutdown)
- [x] `go test -race` gate + chaos-тест pool cancel при `workers>1`
- [x] Report bridge: origin allowlist (`null`, `localhost`, `127.0.0.1`) вместо `CORS: *`

### 17.5b Отчёт после прогона

- [x] **HTML не открывался при падении** — `finalizeGUIReports` пишет отчёт до return error; UI открывает при `result.error` (кроме cancel)

### 17.6 Симптомы пользователя (связь с аудитом)

| Симптом | Причина | Статус |
|---------|---------|--------|
| Зависание «Выполняется» | Нет timeout / UI без finally | [x] лимит 20 мин, finally в `executeRun` |
| Зависание на 1-м шаге | Recorder + run на одной page | [x] `HoldForTestRun` |
| EPIPE в консоли после Ctrl+C | Go exit раньше Node driver | [x] shutdown wait + drain |
| Множественные перезапуски | OOM / crash Node | [x] единый `playwrightrt` + light HTML по умолчанию |
| Отчёт не открывается при fail | skip write + `!result.error` в UI | [x] `finalizeGUIReports` + open on fail |
| OTP-диалог не всплывает | `EmailCode()` после wait полей; тихий skip | [x] prompt first + `EmailCodeForStep` + `WindowShow` |
| Браузер не открывается при первом запуске | нет проекта; ложный `browser-opened`; «Быстрый старт» → только диалог | [x] auto examples + deferred events + quickStart→`startRecord` |
| «Идёт запись» до появления окна | pre-emit `record-started` в `StartRecord` | [x] emit после `OnBrowserOpened` |

### 17.7 Файлы изменений (фаза 17)

`internal/player/{browser_pool.go,browser_session.go,browser_cleanup.go,network.go,artifacts.go,trace_lifecycle.go,runner_parallel.go,action_context.go,async_drain.go,prompt.go}`
`internal/gui/{path_guard.go,report_bridge.go,record.go,service_shutdown.go,allure.go,trace.go,features.go,service.go,run_browser.go,http_auth.go,vanessa_monitor.go}`
`internal/recorder/{session.go,live_session.go}`
`internal/cli/run.go`
`internal/report/{html_payload.go,html_assets/viewer.js}`
`internal/plugin/install.go`
`internal/wailsapp/app.go`
`main.go`
`frontend/src/App.svelte` (OTP, browser first-run, report handlers)
`frontend/src/lib/i18n/locales/{ru,en}/journal.ts`
`internal/settings/{project.go,nav_wait.go}`
`internal/playwrightrt/runtime.go`
`internal/gui/run_browser.go` (finalizeGUIReports, report on fail)
