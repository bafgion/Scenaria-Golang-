# Scenaria Go — Roadmap

Статус: **master** v0.15; **Wails IDE** — основной продукт. Python/Qt — снят с поддержки (экспорт в Python сохранён).

## Приоритеты

| # | Направление | Статус |
|---|-------------|--------|
| P0 | Wails IDE | done |
| P0 | Recorder | done |
| P1 | Monaco / редактор IDE | done (Фаза 6) |
| P1 | Паритет Python IDE | in progress (Фаза 7) |
| P1 | Allure | done |
| P1 | Portable release | done |

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
