# Scenaria v1.0 Stabilization Roadmap

## Цель

Подготовить Scenaria к стабильной версии v1.0.

Фокус roadmap:

* не терять и не смешивать содержимое вкладок;
* гарантировать корректный batch run;
* синхронизировать IDE diagnostics/autocomplete/inlay hints с runner;
* ввести явную идентичность проекта, запуска, recorder-сессии и отчётов;
* изолировать отчёты и артефакты по запускам;
* убрать основные race-condition и stale-event сценарии;
* снизить лишние Wails-вызовы, повторный parsing и memory pressure.

---

# 0. Add Safety Tests Before Refactoring

**Priority:** Critical
**Status:** Done
**Goal:** зафиксировать текущее поведение и защититься от новых регрессий.

## Why First

Некоторые зоны нельзя безопасно менять без тестов:

* `loadFeature`;
* Monaco model activation;
* `RunRequest`;
* `ExecutionPlan`;
* `StepMatcher`;
* recorder events;
* report layout;
* project/run lifecycle.

## Tasks

* [x] Добавить frontend test: закрытие активной вкладки активирует правильную Monaco model.
* [x] Добавить frontend test: закрытие Welcome tab активирует правильный feature tab.
* [x] Добавить frontend test: editor change от stale model не меняет активную вкладку.
* [x] Добавить frontend test: save/reload не меняет другую вкладку после tab switch.
* [x] Добавить frontend test: stale validation response игнорируется после изменения текста.
* [x] Добавить frontend test: batch run после single-scenario run очищает stale filters.
* [x] Добавить Go test: `RunRequest` копирует `Targets` / `Vars`.
* [x] Добавить Go test: `ExecutionPlan` строится заново на каждый запуск.
* [x] Добавить Go test: IDE/runtime step parity для базовых шагов.
* [x] Добавить Go race/concurrency test: parallel run status writes.
* [x] Добавить test: late recorder event ignored.
* [x] Добавить test: report atomic write failure keeps previous report.

## Acceptance Criteria

* [x] До начала крупных исправлений есть минимальная regression-сетка.
* [x] Тесты покрывают три главных бага: tabs, batch, parser.
* [x] Тесты покрывают stale async/event сценарии.
* [x] Можно безопасно начинать исправления с первой технической фазы.

---

# 1. Stabilize Monaco Tabs and Editor State

**Priority:** Critical
**Status:** Done
**Area:** Frontend / Monaco / Tabs / Editor State

## Problem

Содержимое вкладок может пропадать, становиться stale или записываться в другой файл. Главная причина — не переиспользование одной Monaco model, а неправильный порядок `activeTab → loadFeature → Monaco activation`: `activeTab` меняется до фактической активации model, а `loadFeature()` может выйти раньше из-за `path === activeTab`.

## Scope

```text
frontend/src/App.svelte
frontend/src/lib/MonacoEditor.svelte
frontend/src/lib/monacoTabModels.ts
frontend/src/lib/editorTextSync.ts
```

## Tasks

* [x] Исправить `finalizeCloseTab`: не выставлять `activeTab = next.path` до активации Monaco model.
* [x] Исправить `closeWelcomeTab` по той же схеме.
* [x] Добавить `forceActivate` / `activateExisting` режим для `loadFeature`.
* [x] Для existing tab всегда активировать Monaco model, если это явное переключение.
* [x] Сделать editor change event path-aware: `{ path, modelUri, text }`.
* [x] В `onEditorChange` игнорировать event, если event path/model не совпадает с активной вкладкой.
* [x] Ввести canonical path helper для вкладок, Monaco models, markers, dirty state и batch selection.
* [x] Добавить reconciliation существующей Monaco model при authoritative reload/draft restore.
* [x] `syncTabContent` сделать active-tab-only или требовать явный source text/path.
* [x] Save operation привязать к `pathAtStart`.
* [x] Disk reload привязать к `pathAtStart`.
* [x] Validation response защитить `textVersion`.
* [x] Inlay hints и diagnostics строить по одной версии текста.

## Removed Duplicates

Следующие ранее отдельные пункты объединены в эту фазу:

* Monaco tabs stabilization;
* path-aware editor change;
* save path-stability;
* disk reload path-stability;
* `syncTabContent` safety;
* validation text version guard;
* inlay hints same snapshot;
* canonical frontend paths for editor.

## Acceptance Criteria

* [x] Закрытие активной вкладки не оставляет Monaco пустым.
* [x] Закрытие Welcome tab не ломает активную model.
* [x] Быстрое переключение вкладок не смешивает содержимое.
* [x] Editor event от model A не может изменить tab B.
* [x] Save file A не очищает draft file B.
* [x] Disk reload file A не меняет file B.
* [x] Diagnostics и inlay hints соответствуют текущей версии текста.
* [x] Undo/redo работает отдельно для каждой открытой model.
* [x] Один физический файл имеет один canonical frontend key.

## Manual QA

> Прогон 2026-07-10: [MANUAL-QA-ROADMAP-RUN.md](MANUAL-QA-ROADMAP-RUN.md) — mock E2E 109/109, unit 282/282, CI green.

* [ ] Открыть 5 файлов. _(desktop)_
* [x] Быстро переключаться между ними. _(E2E `app-ui`)_
* [x] Внести разные изменения в каждый файл. _(E2E `qa-daily-use` 2.1)_
* [x] Закрыть активную вкладку из середины. _(E2E `app-ui`)_
* [ ] Закрыть первую и последнюю вкладку. _(desktop)_
* [x] Закрыть Welcome tab. _(E2E `app-ui`)_
* [x] Проверить undo/redo по вкладкам. _(E2E `app-ui`)_
* [x] Проверить save после переключения. _(E2E `app-ui`)_
* [x] Проверить external reload после переключения. _(E2E `app-ui`)_

---

# 2. Stabilize Batch Execution

**Priority:** Critical
**Status:** Done
**Area:** Frontend / Runner / ExecutionPlan

## Problem

Batch run может запускать старый subset тестов. Точный найденный root cause: `runBatchSelected()` строит options через `{ ...lastRun, dryRun }`, из-за чего batch run наследует stale `scenario`, `tag`, `startStep`, `endStep`; backend получает targets, но затем фильтрует план по старым параметрам. Дополнительно `toggleBatchMode()` откладывает select-all на next frame, поэтому быстрый run может использовать старый `batchSelected`.

## Scope

```text
frontend/src/App.svelte
frontend/src/lib/batchSelection.ts
internal/gui/run_browser.go
internal/gui/service.go
internal/player/suite.go
internal/player/runner_parallel.go
```

## Tasks

* [x] Для default batch run очищать `scenario`.
* [x] Для default batch run очищать `tag`.
* [x] Для default batch run сбрасывать `startStep = -1`.
* [x] Для default batch run сбрасывать `endStep = -1`.
* [x] Передавать в `executeRun` копию: `[...batchSelected]`.
* [x] В `executeRun` копировать `targets`.
* [x] На backend копировать `RunRequest.Targets`.
* [x] На backend копировать `RunRequest.Vars`.
* [x] Сделать `toggleBatchMode` синхронным.
* [x] Убрать deferred select-all race.
* [x] После `refreshProject` делать prune/remap `batchSelected`.
* [x] Добавить run request factory для batch/single/tag/step-range режимов.
* [x] Добавить debug logs:

  * frontend selected count;
  * backend received targets count;
  * execution plan cases count;
  * runner executed count.

## Removed Duplicates

Следующие ранее отдельные пункты объединены в эту фазу:

* stale `lastRun` filters;
* deferred batch select-all race;
* defensive copying;
* batch request contract;
* batch selection prune after refresh;
* ExecutionPlan freshness verification;
* run debug logs.

## Acceptance Criteria

* [x] Предыдущий single-scenario run не фильтрует следующий batch run.
* [x] Предыдущий tag run не фильтрует следующий batch run.
* [x] Выбрано 5 файлов → backend получает 5 targets.
* [x] Выбрано 5 файлов → plan содержит ожидаемое количество cases.
* [x] Повторный batch run строит fresh `ExecutionPlan`.
* [x] Runner queue создаётся заново на каждый запуск.
* [x] После удаления/rename файла selection не содержит stale path.
* [x] UI count, backend target count и plan cases count можно сверить по логам.

## Manual QA

> Прогон 2026-07-09: [MANUAL-QA-ROADMAP-RUN.md](MANUAL-QA-ROADMAP-RUN.md).

* [x] Выбрать 2 теста и запустить. _(E2E `user-journeys`)_
* [ ] Добавить ещё 3 теста и запустить. _(desktop)_
* [ ] Проверить, что выполняются все 5. _(desktop)_
* [ ] Убрать часть тестов и запустить. _(desktop)_
* [ ] Проверить запуск через hotkey сразу после включения batch mode. _(desktop)_
* [ ] Проверить folder selection. _(desktop)_
* [ ] Проверить refresh проекта после удаления файла. _(desktop)_
* [x] Stale filters не переносятся в следующий batch run. _(E2E `batch run clears stale filters`)_

---

# 2.1. Stabilize Parallel Batch Workers

**Priority:** High
**Status:** Done
**Area:** Batch Execution / Parallel Workers / Variables / Browser Context / Progress Events
**Parent Section:** `2. Stabilize Batch Execution`

## Problem

При пакетном запуске с `workers > 1` возможны нестабильные ошибки, которые не проявляются при `workers = 1`.

Основные риски:

* shared variables между сценариями;
* browser context reuse внутри worker-а;
* сценарии могут зависеть от состояния другого worker-а;
* fail-fast может выглядеть как “запустились не все тесты”;
* progress events приходят не в порядке execution plan.

## Tasks

### Variables isolation

* [x] Клонировать `RunRequest.Vars` перед построением execution plan.
* [x] Клонировать variables для каждого `RunCase`.
* [x] Клонировать variables при создании `RunContext`.
* [x] Запретить shared mutable map между parallel scenarios.
* [x] Добавить `go test -race` для parallel scenarios с `Remember()`.

### Worker execution consistency

* [x] Добавить regression test: batch с `workers = 1`.
* [x] Добавить regression test: тот же batch с `workers = 2`.
* [x] Добавить regression test: тот же batch с `workers = 4`.
* [x] Проверить, что selected cases count одинаковый при любом количестве workers.
* [x] Проверить, что каждый selected case выполняется ровно один раз.
* [x] Проверить, что runner queue не теряет cases при parallel execution.

### Browser context policy

* [x] Явно определить browser context policy:

  * isolated context per scenario;
  * reuse context per worker;
  * shared auth state only by explicit setting.
* [x] Задокументировать разницу между `workers = 1` и `workers > 1`.
* [x] Если context переиспользуется внутри worker-а, явно очищать или документировать:

  * cookies;
  * localStorage;
  * sessionStorage;
  * opened tabs;
  * permissions;
  * downloads;
  * network routes.
* [x] Если auth/session state должен переиспользоваться, сделать это explicit setting.

### Fail-fast behavior

* [x] Проверить поведение `ContinueOnFail = false`.
* [x] Проверить поведение `ContinueOnFail = true`.
* [x] В UI явно показывать настройку:

  * “Stop on first failure”;
  * “Continue on fail”.
* [x] В report разделять:

  * failed;
  * canceled;
  * skipped;
  * not started.
* [x] Не показывать canceled scenarios как обычные failed.
* [x] Не создавать ощущение, что batch “потерял” тесты.

### Progress events

* [x] Frontend progress handling должен быть order-independent.
* [x] Не считать, что progress events придут в порядке `1,2,3,4`.
* [x] Использовать `caseId` / `index` / `runId` для обновления конкретного case.
* [x] Progress bar должен считать completed count, а не последний пришедший index.
* [x] Добавить debug logs для:

  * scheduled cases;
  * started cases;
  * finished cases;
  * canceled cases;
  * failed cases.

## Acceptance Criteria

* [x] Batch с `workers = 1`, `workers = 2`, `workers = 4` выполняет одинаковый набор selected cases.
* [x] Каждый selected case выполняется ровно один раз.
* [x] Variables не протекают между scenarios.
* [x] `go test -race` не показывает race по variables.
* [x] Browser context reuse policy явно определена.
* [x] ContinueOnFail поведение понятно в UI и report.
* [x] Progress UI корректен при out-of-order completion events.
* [x] Canceled/not-started scenarios не выглядят как потерянные тесты.

---

# 3. Introduce ProjectSession, RunSession and Operation Identity

**Priority:** Critical
**Status:** Done
**Area:** Wails / Backend / Lifecycle / Events

## Problem

Backend использует mutable singleton-state: `projectPath`, `runCtx/runCancel`, `liveSession`, `tempFeatureDirs`, report bridge token, global stdout lock и другие состояния. При этом у project/run/record/report events нет единого ownership-протокола: `runId`, `recordSessionId`, `projectVersion`. Это позволяет старым операциям обновлять новый проект или новую UI-сессию.

## Scope

```text
internal/wailsapp/app.go
internal/gui/service.go
frontend/src/App.svelte
frontend event handlers
run/validate/record/report methods
```

## Tasks

* [x] Ввести `ProjectSession { ID, Root, Version, Context }`.
* [x] Увеличивать `ProjectVersion` при open/switch/close project.
* [x] Ввести `RunSession { RunID, ProjectVersion, RequestSnapshot, Context, TempResources }`.
* [x] Ввести `RecordSessionID`.
* [x] Ввести `BrowserSessionID`.
* [x] Ввести `ReportID`.
* [x] Все Wails events должны передавать relevant identity:

  * `projectVersion`;
  * `runId`;
  * `recordSessionId`;
  * `browserSessionId`;
  * `reportId`;
  * `jobId`.
* [x] Frontend должен игнорировать stale events.
* [x] `OpenProject` должен стать lifecycle boundary: cancel/detach project-scoped work.
* [x] Добавить explicit concurrent run policy: reject или explicit cancel/restart.
* [x] Добавить timeout/cancellation для Wails async jobs.
* [x] Добавить stale-event debug logs.

## Removed Duplicates

Следующие темы объединены в эту фазу:

* ProjectVersion;
* ProjectSession;
* RunSession;
* operation IDs;
* stale Wails events;
* async job timeout;
* backend run concurrency policy;
* stale frontend response guards.

## Acceptance Criteria

* [x] Run из project A не может обновить UI project B.
* [x] Validation из старой версии проекта игнорируется.
* [x] Recorder event из старой сессии игнорируется.
* [x] Report action привязан к своему report/run/project.
* [x] Новый run не отменяет старый silently.
* [x] Frontend не зависает навсегда, если backend не прислал finish event.
* [x] Все long-running events трассируются по ID.

## Manual QA

> Прогон 2026-07-09: [MANUAL-QA-ROADMAP-RUN.md](MANUAL-QA-ROADMAP-RUN.md).

* [x] Открыть project A. _(E2E/Go session tests)_
* [x] Запустить run. _(E2E mock run)_
* [x] Быстро открыть project B. _(Go `project_session_test`)_
* [x] Проверить, что старые результаты не применились. _(E2E stale events)_
* [x] Повторить с validation. _(E2E + Go)_
* [x] Повторить с recorder. _(Go `record_lifecycle_test`)_
* [x] Повторить с report actions. _(E2E results/trace)_

---

# 4. Unify StepMatcher and EditorAnalysisService

**Priority:** Critical
**Status:** Done
**Area:** Parser / Diagnostics / Autocomplete / Inlay Hints / Runner

## Problem

IDE и runner не используют единый pipeline обработки шагов. Runner идёт через `gherkin.ParseFeatureFile → NormalizeFeatureText → structured Step.Text → stepdsl.Parse`, а IDE частично использует fallback raw scanners, `ParseEditorSteps` и static `stepcatalog.CompletionsForLineLang`. Поэтому runtime-valid step может стать `Unknown step` в IDE.

## Scope

```text
internal/gherkin
internal/stepdsl
internal/stepcatalog
internal/gui/validate.go
internal/gui/editor_steps.go
internal/gui/scenario_hints.go
frontend/src/lib/gherkinCompletions.ts
frontend/src/lib/gherkinInlayHintsProvider.ts
```

## Tasks

* [x] Создать общий `StepMatcher`.
* [x] Создать общий `EditorAnalysisService`.
* [x] Вынести shared keyword classifier.
* [x] Вынести shared step text normalization.
* [x] Перевести diagnostics на `EditorAnalysisService`.
* [x] Перевести `ParseEditorSteps` / inlay hints на `EditorAnalysisService`.
* [x] Перевести scenario hints на тот же analysis result.
* [x] Reconcile `stepcatalog` with `stepdsl`.
* [x] Autocomplete должен использовать metadata из runtime-compatible registry или иметь parity tests.
* [x] Убрать independent raw fallback matcher.
* [x] Добавить parser recovery mode для editor-time invalid text.
* [x] Добавить golden parity tests.

## Removed Duplicates

Следующие ранее отдельные блоки объединены:

* Parser / Step Registry;
* shared keyword classification;
* shared normalization;
* ParseEditorSteps rewrite;
* autocomplete parity;
* diagnostics/markers stability;
* parser recovery mode;
* editor analysis performance consolidation.

## Acceptance Criteria

* [x] Если runner выполняет шаг, IDE не показывает `Unknown step`.
* [x] Diagnostics, inlay hints, steps panel и autocomplete используют совместимый источник истины.
* [x] Один malformed scenario не ломает valid steps в другом scenario.
* [x] `Дано/Когда/Тогда/И/Но/Допустим/*` покрыты тестами.
* [x] `ё/е`, smart quotes, tabs, spaces, escaped selectors покрыты тестами.
* [x] Completion items либо соответствуют executable step pattern, либо явно marked snippet-only.

## Manual QA

> Прогон 2026-07-09: [MANUAL-QA-ROADMAP-RUN.md](MANUAL-QA-ROADMAP-RUN.md).

* [x] Открыть файл с известными шагами. _(E2E examples)_
* [x] Проверить diagnostics. _(unit + E2E validate)_
* [x] Проверить autocomplete. _(unit tests)_
* [x] Проверить inlay hints. _(unit tests)_
* [x] Внести временную синтаксическую ошибку. _(E2E validate confirm)_
* [x] Проверить, что valid known steps не стали false unknown. _(unit)_
* [x] Запустить тот же сценарий runner-ом. _(Go `examples_integration` @smoke)_

---

# 4.1. Stabilize Step Execution Semantics

**Priority:** High
**Status:** Done
**Area:** Step Execution / Variables / Loops / Retry
**Parent Section:** `4. Unify StepMatcher and EditorAnalysisService`

## Problem

Механизм выполнения шагов должен быть предсказуемым:

* переменные не должны протекать между сценариями;
* `repeat` не должен молча менять количество итераций;
* `for_each` не должен скрывать DOM errors;
* `if/while` должны отличать `false` от runtime error;
* retry не должен повторять опасные действия без явного разрешения.

## Tasks

### Variable isolation

* [x] Клонировать `RunRequest.Vars`.
* [x] Клонировать variables при создании `RunCase`.
* [x] Клонировать variables при создании `RunContext`.
* [x] Запретить shared mutable map между сценариями.
* [x] Добавить `go test -race` для parallel variables.

### Repeat semantics

* [x] `repeat 0` должен быть ошибкой или явно разрешённой конструкцией с warning.
* [x] `repeat < 0` должен быть ошибкой.
* [x] `repeat > MaxLoopIterations` должен быть ошибкой, не silent clamp.
* [x] Overflow при parsing count должен быть ошибкой.
* [x] Repeat должен выполнять ровно requested count.

### ForEach semantics

* [x] Не игнорировать ошибки `InnerText`.
* [x] Логировать selector/index/iteration при ошибке.
* [x] Явно определить модель:

  * snapshot;
  * live DOM;
  * strict mode.
* [x] Если элемент исчез во время итерации — вернуть понятную ошибку.
* [x] Переменная цикла не должна получать fallback/мусорное значение.

### If / While semantics

* [x] `EvaluateCondition` должен возвращать `(bool, error)`.
* [x] Отличать false condition от locator/page error.
* [x] Отличать timeout от false.
* [x] Отличать context canceled от runtime failure.
* [x] `while` не должен завершаться silently из-за Playwright error.

### Retry policy

* [x] Разделить retry actions на safe и risky.
* [x] Safe by default:

  * waits;
  * assertions;
  * visibility checks.
* [x] Risky only opt-in:

  * click;
  * double-click;
  * fill;
  * select;
  * check;
  * uncheck;
  * download-click.
* [x] Добавить настройки:

  * `retryWaits`;
  * `retryAssertions`;
  * `retryActions`.
* [x] Retry attempts должны отображаться в report.

## Acceptance Criteria

* [x] Scenario A не может изменить variables Scenario B.
* [x] Repeat не меняет count silently.
* [x] ForEach не скрывает DOM errors.
* [x] While не завершает цикл из-за error как будто condition false.
* [x] Click/fill/download не retry-ятся silently.
* [x] Retry attempts видны в отчёте.
* [x] `go test -race` не показывает race по variables.

---

# 4.2. Stabilize Runtime Step Semantics and Reporting

**Priority:** High
**Status:** Done
**Area:** Step Records / Loop Reporting / Retry Reporting / Terminal Steps
**Parent Section:** `4. Unify StepMatcher and EditorAnalysisService`

## Problem

Отчёт должен объяснять не только какой шаг упал, но и:

* на какой итерации;
* на какой retry attempt;
* был ли шаг skipped/canceled;
* был ли сценарий завершён terminal browser action.

## Tasks

* [x] Добавить `IterationPath` в `StepRecord`.
* [x] Разделить:

  * logical step;
  * loop iteration;
  * retry attempt;
  * generated wait/assert step.
* [x] Screenshots/artifacts должны включать iteration index.
* [x] Negative wait duration считать ошибкой.
* [x] Zero wait duration либо warning, либо explicit allowed.
* [x] `close-browser` / `close-tab` не должны silently делать scenario passed при оставшихся шагах.
* [x] Если browser closed до конца scenario:

  * remaining steps = skipped/canceled;
  * scenario status не должен быть passed.
* [x] Добавить total action attempts limit per scenario.
* [x] В отчёте показывать:

  * loop iteration;
  * retry attempt;
  * skipped;
  * canceled;
  * terminal browser action.

## Acceptance Criteria

* [x] Ошибка внутри `repeat[3]` видна как `repeat[3]`, а не просто line number.
* [x] Retry внутри loop не маскирует фактическое число попыток.
* [x] Negative wait не превращается в `0ms`.
* [x] `close-browser` не делает невыполненные шаги passed.
* [x] Report показывает logical order и runtime attempts.

---

# 5. Stabilize Backend Storage, Locks and File Operations

**Priority:** High
**Status:** Done
**Area:** Backend / Concurrency / File IO / Settings / RunStatus

## Problem

Есть несколько shared storage/race зон: `runstatus.Store` пишет JSON read-modify-write без lock, `tempFeatureDirs` общий для всех запусков, settings read-modify-write может терять обновления, project file operations могут пересекаться с run/validate/refresh.

## Scope

```text
internal/runstatus
internal/settings
internal/gui/service.go
internal/gui/features.go
internal/gui/http_auth.go
internal/player/runner_parallel.go
```

## Tasks

* [x] Добавить mutex или per-file lock для `runstatus.Store`.
* [x] Сделать runstatus writes atomic.
* [x] Перевести runstatus на batch write per run или JSONL.
* [x] Temp feature files привязать к `RunSession`, убрать global cleanup.
* [x] Добавить `SettingsStore` с mutex и atomic writes.
* [x] Перевести `SaveSettings`, `SaveHTTPAuth`, recents, credentials на `SettingsStore`.
* [x] Ввести project filesystem lock или snapshot layer.
* [x] Run должен использовать immutable loaded feature snapshot.
* [x] Save/rename/delete/import/replace должны координироваться с refresh/validate/run preparation.
* [x] OTP prompt сделать single-flight или promptId-based.
* [x] Report bridge tokens сделать per report/run или bounded token set.
* [x] CLI/global stdout capture заменить context-aware service APIs.

## Removed Duplicates

Объединены:

* runstatus race;
* temp feature cleanup;
* settings lost updates;
* project file operation races;
* OTP concurrency;
* report bridge token coupling;
* global stdout capture;
* immutable backend snapshots.

## Acceptance Criteria

* [x] Parallel run не corrupt-ит `run_status.json`.
* [x] Run A не удаляет temp files Run B.
* [x] Settings updates не теряют данные.
* [x] Run preparation не читает half-written files.
* [x] OTP prompts не перезаписывают друг друга.
* [x] Report actions не ломаются при открытии нового report.
* [x] GUI не зависит от global stdout replacement для нормальных операций.

## Manual QA

> Прогон 2026-07-09: [MANUAL-QA-ROADMAP-RUN.md](MANUAL-QA-ROADMAP-RUN.md).

* [x] Запустить parallel run. _(Go cancel/parallel tests)_
* [x] Проверить run status/history. _(E2E run history flaky)_
* [ ] Запустить run с unsaved/temp feature. _(desktop)_
* [x] Отменить и сразу запустить другой. _(Go `run_session_test`)_
* [ ] Проверить settings save параллельно с recents/HTTP auth. _(desktop)_
* [ ] Проверить project refresh во время save/rename/delete. _(desktop)_

---

# 6. Stabilize Recorder Lifecycle

**Priority:** High
**Status:** Done
**Area:** Recorder / Live Browser / Wails Events / Editor Target

## Problem

Recorder events сейчас могут быть anonymous/global и применяться к текущему active editor, хотя событие относится к старой recorder session. Recorder должен владеть browser session, но не должен напрямую мутировать editor без explicit target path/model.

## Scope

```text
internal/gui/record.go
internal/recorder/live.go
internal/recorder/live_session.go
internal/recorder/session.go
frontend/src/App.svelte
frontend/src/lib/recordedStepEditor.ts
```

## Tasks

* [x] Каждая recorder session получает `RecordSessionID`.
* [x] Каждая browser session получает `BrowserSessionID`.
* [x] Recorder event содержит `projectVersion`, `recordSessionId`, `browserSessionId`, `targetPath`.
* [x] Backend callbacks проверяют active generation перед emit.
* [x] Frontend игнорирует stale recorder events.
* [x] `applyLiveRecordedStep` должен писать в explicit target model/path, не в ambient `activeTab`.
* [x] Recording target file/model должен открываться до старта recording.
* [x] `CloseBrowser` должен retire recorder generation.
* [x] `StopRecordingCapture` должен быть idempotent.
* [x] Browse-only → capture → relaunch должен переинжектить recorder scripts.
* [x] Recorder context должен cancel-иться при natural termination.
* [x] LiveSession Playwright calls должны возвращать typed `ErrBrowserClosed`.
* [x] Recorder step events должны быть snapshot-based или иметь полный op model: `upsert/delete/reset/snapshot`.
* [x] Recorded Gherkin formatting должен принадлежать backend или иметь shared contract.

## Removed Duplicates

Объединены:

* recorder event identity;
* stale recorder callbacks;
* editor targeting;
* close/stop idempotency;
* browse-to-capture script injection;
* context cancellation;
* Playwright use-after-close;
* recorded step snapshot model;
* recorded Gherkin formatting contract.

## Acceptance Criteria

* [x] Late `record-step` не меняет editor.
* [x] Recorder пишет только в выбранный target file.
* [x] Switch tab during recording не меняет target.
* [x] Project switch invalidates old recorder events.
* [x] Stop twice даёт один effective transition.
* [x] Browser close не оставляет picker hanging.
* [x] Browse-only session после capture/relaunch продолжает писать steps.
* [x] Frontend recorded lines не расходятся с backend recorder state.

## Manual QA

> Прогон 2026-07-09: [MANUAL-QA-ROADMAP-RUN.md](MANUAL-QA-ROADMAP-RUN.md).

* [x] Start recording в `B.feature`, активна `A.feature`. _(E2E recording target)_
* [x] Сделать click/input. _(E2E mock post-record)_
* [x] Проверить, что изменился только `B.feature`. _(E2E recording-target)_
* [x] Переключать вкладки во время записи. _(E2E 3.2 confirm)_
* [x] Stop → Start снова. _(E2E record-resume, 1 flaky)_
* [ ] Close browser во время picker. _(desktop)_
* [x] Project switch во время recording. _(Go `project_session_test`)_
* [x] Проверить отсутствие stale events. _(Go + unit)_

---

# 6.1. Stabilize Element Picker and Selector Generation

**Priority:** High
**Status:** Done
**Area:** Recorder / Picker / Selector Generation
**Parent Section:** `6. Stabilize Recorder Lifecycle`

## Problem

Picker может выбирать некорректный selector, потому что он выбирает первый подходящий selector, но не доказывает, что selector уникален и указывает именно на выбранный элемент.

## Tasks

### Candidate model

* [x] Генерировать ranked selector candidates.
* [x] Для каждого candidate считать:

  * score;
  * uniqueness;
  * matches count;
  * reason;
  * strategy;
  * warnings.
* [x] UI должен показывать selector confidence.
* [x] Пользователь должен видеть альтернативные candidates.

### Validation before return

* [x] Проверять, что selector matches exactly one element.
* [x] Проверять, что matched element === picked element или корректный actionable ancestor.
* [x] Проверять visibility.
* [x] Проверять actionability по типу действия.
* [x] Не возвращать non-unique selector silently.

### Strategy order

* [x] Для click предпочитать:

  * `data-testid`;
  * role/aria;
  * title;
  * stable id;
  * contextual;
  * text.
* [x] Для input предпочитать:

  * `data-testid`;
  * id;
  * name;
  * aria;
  * label;
  * placeholder.
* [x] Понизить score text-only selectors.
* [x] Добавить warning для text-only selector.

### Input / label targeting

* [x] `label[for]` должен возвращать selector control, не label.
* [x] Nested label должен возвращать вложенный input.
* [x] Adjacent label должен строить contextual input selector.
* [x] Fill/select steps не должны получать selector label как primary target.

### Action-aware picker

* [x] Определять suggested action:

  * click;
  * fill;
  * select;
  * check;
  * uncheck;
  * hover.
* [x] Input-like elements:

  * input;
  * textarea;
  * select;
  * contenteditable;
  * role=textbox;
  * role=combobox;
  * role=spinbutton;
  * role=searchbox.
* [x] Recorded step должен использовать suggested action.
* [x] Пользователь может переопределить action.

### Iframe / Shadow DOM / SVG / Canvas

* [x] Добавить iframe-aware picker.
* [x] Same-origin iframe: выбирать внутренний элемент.
* [x] Cross-origin iframe: честно показывать limitation.
* [x] Добавить shadow DOM hit-test.
* [x] SVG click нормализовать до clickable ancestor.
* [x] Canvas selector должен иметь uniqueness/warning.

### Deferred (not part of 6.1)

* См. отдельную секцию **6.1.1 Component-Library Selector Heuristics**.

## Acceptance Criteria

* [x] Picker возвращает validated selector.
* [x] Selector указывает на выбранный элемент.
* [x] Non-unique selector не выбирается silently.
* [x] Input получает selector input/control, не label.
* [x] `data-testid`/role/aria предпочитаются raw text.
* [x] Iframe behavior explicit.
* [x] Low-confidence selector показывает warning.
* [x] Пользователь может выбрать альтернативный selector.

---

# 6.1.1. Component-Library Selector Heuristics

**Priority:** Medium
**Status:** Done
**Area:** Recorder / Picker / Selector Generation
**Parent Section:** `6.1. Stabilize Element Picker and Selector Generation`

## Problem

Базовый picker (6.1) стабилен для generic DOM, но component libraries (MUI, Ant Design, Element Plus, Bootstrap React и т.д.) часто требуют специфичных селекторов: shadow parts, role wrappers, portal menus, virtualized lists.

## Tasks

* [x] Собрать каталог библиотек и типовых DOM-паттернов (button, input, select, menu, dialog).
* [x] Добавить optional heuristics layer поверх `generateCandidates` (не ломая generic strategies).
* [x] MUI: `MuiButton-root`, `MuiInputBase-input`, `MuiMenuItem-root`, `data-testid` hooks.
* [x] Ant Design: `.ant-btn`, `.ant-input`, `.ant-select`, dropdown portal containers.
* [x] Понижать score generic text-only selectors когда library-specific candidate валиден.
* [x] Warnings для portal/virtualized targets («may need menu open»).
* [x] Settings: enable/disable library packs per project.
* [x] Integration tests на minimal fixtures per library.

## Acceptance Criteria

* [x] Picker на MUI/Ant fixture предпочитает library selector над raw text.
* [x] Generic sites без library markup не меняют поведение 6.1.
* [x] Library heuristics отключаемы в настройках.
* [x] Каждый library pack покрыт regression test.

---

# 6.2. Stabilize Browser Selector Validation

**Priority:** High
**Status:** Done
**Area:** Selector Validation / Browser Validation / Dynamic UI
**Parent Section:** `6. Stabilize Recorder Lifecycle`

## Problem

Browser selector validation может давать ложную уверенность. Она может проверять selector на initial page, не выполняя flow, который делает элемент доступным. Также validation должна быть action-aware, а не только visible-check.

## Tasks

* [x] Разделить validation modes:

  * static validation;
  * flow-aware validation.
* [x] Static validation не считать источником истины для dynamic UI.
* [x] В UI показывать: “validated on current/initial page only”.
* [x] Для dynamic elements не выдавать misleading missing без контекста.
* [x] Сделать validation action-aware:

  * click → visible + enabled + actionable;
  * fill → editable input/textarea/contenteditable;
  * select → select/combobox;
  * check → checkbox/radio;
  * upload → input[type=file].
* [x] Для chained selectors показывать matches count.
* [x] Для hover selectors не выбирать `.First()` без ambiguous warning.
* [x] Для contextual selectors проверять container + target uniqueness.
* [x] Flow-aware validation должна уметь выполнять safe actions до проверяемого шага.
* [x] Добавить validation diagnostics в HTML-отчёт (см. **6.2.1**).

## Acceptance Criteria

* [x] Selector validation не говорит “OK”, если action невозможен.
* [x] Selector validation не говорит “missing” без предупреждения о dynamic flow.
* [x] Fill selector проверяется как editable target.
* [x] Click selector проверяется как actionable target.
* [x] Ambiguous chained selector получает warning.
* [x] Validation result объясняет limitation.

---

# 6.2.1. Validation Diagnostics in HTML Report

**Priority:** Medium
**Status:** Done
**Area:** Selector Validation / Reports
**Parent Section:** `6.2. Stabilize Browser Selector Validation`

## Tasks

* [x] Показывать в HTML-отчёте limitation, mode, matchCount и action-aware diagnostics для шагов с selector validation.
* [x] Согласовать формат с GUI ValidatePanel.

## Acceptance Criteria

* [x] HTML-отчёт отражает те же diagnostics, что и панель проверки в GUI.

---

# 6.3. Stabilize Recorder Event Ordering and Navigation Causality

**Priority:** High
**Status:** Done
**Area:** Recorder / Event Ordering / Navigation / Generated Steps
**Parent Section:** `6. Stabilize Recorder Lifecycle`

## Problem

Recorder может записать переход по URL раньше клика, который этот переход вызвал.

Неверный output:

```gherkin
открыт "https://site/dashboard"
нажимаю "button:has-text(\"Войти\")"
```

Правильный output:

```gherkin
нажимаю "button:has-text(\"Войти\")"
ожидаю адрес "https://site/dashboard"
```

Или только click, если URL wait disabled.

## Root Cause

Recorder смешивает два источника:

```text
1. browser-side interaction events:
   click/input/change/press

2. backend polling:
   page.URL() changed
```

URL polling может сработать раньше, чем backend прочитает browser-side click queue.

## Tasks

* [x] Добавить sequence number в recorder events.
* [x] Добавить timestamp в recorder events.
* [x] Сначала drain browser recorder events, потом проверять URL change.
* [x] Добавить navigation correlation с последним user action.
* [x] Не записывать click-caused navigation как `открыт` перед click.
* [x] Для click-caused navigation генерировать:

  * click step;
  * optional `ожидаю адрес`.
* [x] `открыт` использовать только для explicit navigation/open.
* [x] Добавить correlation window, например 0–2000ms.
* [x] Добавить immediate event flush для navigation-causing actions.
* [x] Защититься от потери old-page event queue при full page navigation.
* [x] Настройка отключения `ожидаю адрес` после click (URL wait disabled).
* [x] Добавить тесты:

  * full navigation after click;
  * redirect without click;
* [x] Добавить тесты:

  * SPA route after click;
  * delayed navigation;
  * old page destroyed before polling.

## Acceptance Criteria

* [x] Click, вызывающий navigation, записывается до navigation step.
* [x] Recorder не выводит `открыт URL` перед click, который вызвал этот URL.
* [x] Full page navigation не теряет preceding click.
* [x] SPA route change не меняет порядок событий.
* [x] Initial page open всё ещё записывается как `открыт`.
* [x] Generated scenario replays in logical user order.

---

# 7. Stabilize Reports and Artifacts

**Priority:** High
**Status:** Done
**Area:** Reports / HTML / Allure / JUnit / Artifacts / History

## Problem

Reports и artifacts пишутся в fixed paths `.scenaria/...`; repeated/concurrent runs могут перезаписывать HTML, Allure, screenshots, traces, summary, JUnit. Scenario outline examples не имеют стабильной identity, report writes не atomic, old reports могут потерять assets. Финальный план отдельно указывает fixed output dirs, non-atomic writes, отсутствие `CaseID/ExampleIndex` и memory pressure в reports как критичные риски.

## Scope

```text
internal/report/*
internal/report/allure/writer.go
internal/gui/run_browser.go
internal/player/runner.go
internal/player/artifacts.go
internal/runstatus
frontend report opening logic
```

## Tasks

* [x] Добавить `RunID` в report generation.
* [x] Добавить `CaseID` в `RunCase` / `ScenarioResult`.
* [x] Добавить `ExampleIndex` в result identity.
* [x] Перейти на `.scenaria/runs/<runID>/...`.
* [x] Хранить screenshots/traces/videos per run/case.
* [x] Добавить `latest.json` или stable latest pointer.
* [x] Все report writes сделать atomic.
* [x] Allure писать через staging dir.
* [x] Trace/screenshot names строить по `RunID/CaseID/ExampleIndex`.
* [x] `findPlanCase` и history matching перевести на `CaseID`.
* [x] Добавить status `canceled`.
* [x] Partial canceled report разрешить открывать по policy.
* [x] History/flaky metrics должны исключать current `RunID`.
* [x] Artifact write errors должны быть visible.
* [x] Исправить path confinement check.

## Removed Duplicates

Объединены:

* per-run reports;
* atomic writes;
* latest pointer;
* artifact naming;
* scenario outline identity;
* old report asset stability;
* Allure staging;
* canceled status;
* history excludes current run;
* artifact error handling;
* report path safety.

## Acceptance Criteria

* [x] Каждый report принадлежит одному `RunID`.
* [x] Каждый executed case имеет `CaseID`.
* [x] Scenario outline examples различимы.
* [x] Старый открытый report остаётся рабочим после нового run.
* [x] HTML/Allure/JUnit/Summary не corrupt-ятся при сбое записи.
* [x] New run не удаляет artifacts old run.
* [x] Canceled run не отображается как обычный failure.
* [x] History не сравнивает run сам с собой.
* [x] Artifact write failure не скрывается.

## Manual QA

> Прогон 2026-07-09: [MANUAL-QA-ROADMAP-RUN.md](MANUAL-QA-ROADMAP-RUN.md).

* [x] Run с HTML/Allure/traces. _(Go report write tests)_
* [x] Открыть report. _(E2E `html-report.spec.ts`)_
* [ ] Run повторно. _(desktop — проверить latest pointer)_
* [ ] Проверить старый report и latest report. _(desktop)_
* [x] Проверить scenario outline с двумя examples. _(Go fixtures)_
* [ ] Проверить canceled run partial report. _(desktop)_
* [ ] Проверить Allure output после failed write simulation. _(desktop)_

---

# 8. Stabilize Playwright Runner Lifecycle

**Priority:** Medium / High
**Status:** Done
**Area:** Runner / Playwright / Browser Pool / Goroutines

## Problem

Playwright cleanup в целом аккуратный, но canceled operations могут оставлять drain goroutines, startup cancellation может ждать до 30 секунд, browser pool shutdown может масштабироваться по workers. Эти риски не первые по приоритету, но важны для стабильности repeated runs.

## Scope

```text
internal/player/action_context.go
internal/player/async_drain.go
internal/player/browser_pool.go
internal/player/browser_session.go
internal/playwrightrt/runtime.go
internal/player/runner_parallel.go
```

## Tasks

* [x] Добавить pending async drain counter.
* [x] Логировать long-running drains.
* [x] Проверить, что browser/context close освобождает stuck Playwright calls.
* [x] Добавить cancel/leak tests для `Goto`, `WaitFor`, `Press`, `ExpectDownload`.
* [x] Добавить browser pool shutdown timing logs.
* [x] Проверить shutdown latency with many workers.
* [x] Добавить goroutine profile dev command.
* [x] Убедиться, что repeated cancel не увеличивает goroutine count unbounded.

## Acceptance Criteria

* [x] 20 repeated cancel не увеличивают goroutine count бесконечно.
* [x] Browser pool закрывает все sessions.
* [x] Долгие stuck calls видны в logs/metrics.
* [x] Cancel/shutdown latency измерима.
* [x] Live browser reuse не ломается.

## Manual QA

> Прогон 2026-07-09: [MANUAL-QA-ROADMAP-RUN.md](MANUAL-QA-ROADMAP-RUN.md).

* [x] Run/cancel 20 раз. _(Go cancel/goroutine tests)_
* [x] Cancel во время navigation/wait/download. _(Go integration)_
* [x] Parallel run с failure. _(Go `parallel_cancel_integration`)_
* [ ] Close app during run. _(desktop — вручную)_
* [x] Проверить, что browsers закрылись. _(desktop-smoke 26/26)_

---

# 8.1. Stabilize Browser Context Isolation

**Priority:** Medium / High
**Status:** Done
**Area:** Browser Context / Scenario Isolation / Parallel Execution
**Parent Section:** `8. Stabilize Playwright Runner Lifecycle`

## Problem

Сценарии не должны случайно влиять друг на друга через browser state:

* cookies;
* localStorage;
* sessionStorage;
* opened tabs;
* downloads;
* traces;
* permissions;
* network state.

## Tasks

* [x] Явно задокументировать browser context reuse policy (`docs/architecture/browser-context-policy.md`).
* [x] Проверить context/page lifecycle per scenario.
* [x] Проверить behavior в parallel workers.
* [x] Определить, когда context reuse допустим.
* [x] Если reuse disabled — очищать:

  * cookies;
  * localStorage;
  * sessionStorage;
  * opened tabs;
  * downloads;
  * permissions;
  * route/network state.
* [x] Если auth state reuse включён — сделать это explicit setting (TestClient profile).
* [x] Artifacts/downloads/traces изолировать per case.
* [x] Добавить тест: Scenario A не влияет на Scenario B.
* [x] Добавить тест: parallel scenarios не делят storage accidentally.

## Acceptance Criteria

* [x] Сценарии изолированы по умолчанию или reuse явно включён.
* [x] Parallel workers не делят mutable browser state случайно.
* [x] Auth reuse контролируемый, не implicit.
* [x] Downloads/artifacts не смешиваются между cases.
* [x] Scenario A не влияет на Scenario B.

---

# 9. Performance and Memory Optimization

**Priority:** Medium
**Status:** Done
**Area:** Performance / Memory / Large Projects

## Problem

Основные performance-риски: autocomplete отправляет полный документ, validation делает несколько full-document passes, project/run повторно парсят файлы, runstatus переписывает историю per scenario, HTML report держит artifacts в памяти и пишет full report дважды.

## Scope

```text
frontend Monaco providers
frontend/src/App.svelte
internal/gui/validate.go
internal/gui/editor_steps.go
internal/gui/scenario_hints.go
internal/report/*
internal/runstatus
internal/scenario
```

## Tasks

* [x] Добавить performance marks и Wails timing logs.
* [x] Добавить pprof/dev diagnostics.
* [x] Autocomplete не должен отправлять full document на каждый request.
* [x] Completion должен передавать language/current line/context.
* [x] Editor analysis должен быть одним backend pass на text version.
* [x] CodeLens/symbols/folding должны шарить per-model parsed structure.
* [x] Coalesce `refreshRunResults` during active run.
* [x] Добавить project index cache по path/mtime/size/hash.
* [x] Artifacts заменить с `[]byte` на `ArtifactRef`/file path/stream.
* [x] `WriteHTMLModePair` должен писать full report один раз.
* [x] Runstatus перейти на batch write или JSONL.
* [x] Debounce `StepsInsertDialog` search.
* [x] Clear symbol cache on tab/project close.

## Removed Duplicates

Объединены:

* autocomplete payload optimization;
* editor analysis consolidation;
* Monaco symbols/cache sharing;
* report memory pressure;
* HTML mode pair double write;
* run results refresh coalescing;
* project parse index;
* runstatus storage optimization;
* step search debounce;
* symbol cache cleanup;
* profiling/metrics.

## Acceptance Criteria

* [x] Typing в large feature не вызывает лишние Wails calls.
* [x] Validation/hints/scenario hints используют один analysis result.
* [x] Completion latency не растёт резко на 5k/20k lines.
* [x] Large reports не удерживают все artifacts в heap.
* [x] Report generation быстрее на 100/500 scenarios.
* [x] Project refresh не repars-ит unchanged files.
* [x] Runstatus не rewrite-ит всю историю на каждый scenario.
* [x] Performance можно измерить через benchmarks/profiles.

## Manual QA

> Прогон 2026-07-09: [MANUAL-QA-ROADMAP-RUN.md](MANUAL-QA-ROADMAP-RUN.md). Все пункты — desktop.

* [ ] Открыть large feature.
* [ ] Быстро печатать 100 символов.
* [ ] Открыть large project.
* [ ] Сгенерировать 100-scenario full report.
* [ ] Открыть/закрыть 50 tabs.
* [ ] Проверить responsiveness.

---

# 10. Architecture Cleanup

**Priority:** Medium
**Status:** Done
**Area:** Architecture / Maintainability / Ownership

## Problem

`App.svelte` и `internal/gui.Service` стали слишком широкими ownership-центрами. Аудит отдельно отмечает, что `internal/gui` стал orchestration god-package, а `App.svelte` владеет слишком большим количеством application state; также нет first-class `ProjectSession`, operation identity и single IDE/runtime step service.

## Scope

```text
frontend/src/App.svelte
frontend stores/controllers
internal/gui.Service
internal/wailsapp.App
new backend domain services
docs/architecture
```

## Tasks

* [x] Добавить `docs/architecture/state-ownership.md`.
* [x] Описать ownership:

  * editor text;
  * active file;
  * project session;
  * run session;
  * recorder session;
  * step analysis;
  * report artifacts;
  * settings.
* [x] Постепенно выделить frontend stores/controllers:

  * `tabsStore`;
  * `projectStore`;
  * `runnerStore`;
  * `diagnosticsStore`;
  * `recorderStore`;
  * `reportsStore`;
  * `wailsEventsController`;
  * `settingsStore`, `runFormStore`, `vanessaRunStore`, `pluginRunStore`, `editorStore`;
  * `layoutStore`, `uiPrefsStore`, `dialogsStore`, `onboardingTourStore`;
  * `recorderPrefsStore`, `recordFormStore`, `tabsStore` (extended);
  * `journalStore`, `featureDialogStore`, `testClientStore`, `validateDialogStore`;
  * `pluginsStore`, `updateDialogStore`, `catalogStore`, `contextMenuStore`, `postRecordStore`;
  * `confirmDialogStore`, `menuStore`, `projectReplaceStore`, `pickerDialogStore`, `httpAuthDialogStore`, `stepsHelpDialogStore`, `otpDialogStore`;
  * `splashStore`, `viewportStore`, `recentsStore`, `settingsDialogStore`, `appMetaStore`, `sessionStore` (layout `bottomTab` in `layoutStore`);
  * `dialogBindController` (dialog `bind:` session locals + sync/flush);
  * `workspaceSessionController` (settings/session DTO build, persist, draft autosave);
  * `paletteCommandsController` (derived palette command list from shell actions).
* [x] Оставить `App.svelte` как UI shell, а не god component (см. `docs/architecture/app-shell.md`).
* [x] Оставить `gui.Service` как Wails-facing façade (см. `docs/architecture/gui-facade.md`).
* [x] Вынести backend services:

  * `ProjectService`;
  * `RunService`;
  * `EditorAnalysisService`;
  * `RecorderService`;
  * `ReportService`;
  * `SettingsService`;
  * `FileOperationService`;
  * `TestClientService`;
  * `PluginService`;
  * `CatalogService`.
* [x] Убрать GUI-to-CLI coupling через global stdout.
* [x] Пересмотреть package-level mutable singletons.
* [x] Добавить reset/cleanup APIs для тестов там, где globals остаются.

## Acceptance Criteria

* [x] У каждого critical state есть один владелец (`docs/architecture/state-ownership.md`).
* [x] Derived state явно обозначен как derived.
* [x] `App.svelte` не владеет unrelated domains напрямую (shell contract в `app-shell.md`).
* [x] `gui.Service` делегирует domain services (карта в `gui-facade.md`, smoke tests в `service_facade_test.go`).
* [x] CLI и GUI используют общие handler APIs (`cli.Run*WithOutput` + in-process services); global stdout capture не используется.
* [x] Новые фичи не добавляют новый source of truth без документации (policy в `state-ownership.md` § Adding New State).

---

# 11. Release Safety Follow-Up

**Priority:** Critical / High
**Status:** Implemented / CI Verification Required
**Area:** v1.0 Stability / Data Safety / Plugin Lifecycle / Runtime Cleanup

## Problem

Phase 0 local verification found remaining release-stability defects after the broad refactor. These items must be handled as focused phases; do not combine plugin, file operation, frontend session, runner, and shutdown changes into one patch.

## Implementation update

Focused release-safety phases 11.1, 11.2, 11.4, 11.5, 11.7, 11.8, 11.9, 11.11 and 11.12 are implemented in the current local workspace and covered by targeted tests. Full release validation and manual Windows smoke remain in Phase 11.

## Current status matrix

| Section | Previous status | Current status | Evidence | Remaining work | v1.0 classification |
| --- | --- | --- | --- | --- | --- |
| 11.1 Plugin filesystem confinement | Planned | Done / Manual Windows Verification | `ValidatePluginID`, `addonPath`, plugin identity/confinement tests | Manual invalid-ID install/uninstall smoke on Windows | Blocker implemented |
| 11.2 Transactional plugin install/update | Planned | Done / Manual Windows Verification | staging/backup/rollback install path, transaction tests | Manual broken-update smoke; descriptor API compatibility remains schema-dependent | Blocker implemented |
| 11.3 Safe duplicate/import destinations | Planned | Done / Manual Windows Verification | filename normalization, `PathGuard`, exclusive copy tests | Manual Windows duplicate/import smoke | Blocker implemented |
| 11.4 Replace in Project rollback/external edit guard | Planned | Done / Manual Windows Verification | planning, byte validation, rollback tests | Manual locked-file Windows smoke | Blocker implemented |
| 11.5 Project-independent Untitled recovery | Planned | Done / Manual Windows Verification | workspace session restore tests, recovery journal merge | Manual restart with missing project | Blocker implemented |
| 11.6 Canonical PathGuard | Verification Required | Implemented / Manual Windows Verification Required | shared `PathGuard`, symlink tests, plugin/file-operation usage | Windows junction/reparse manual verification | Blocker implemented with manual residual risk |
| 11.7 Browser pool late release/rejected slot ownership | Planned | Done / Manual Windows Verification | late close, rejected return, double release and race tests | Manual repeated workers=2/4 cancel smoke | Blocker implemented |
| 11.8 Plugin/Vanessa process lifecycle | Planned | Done / Manual Windows Verification | cancellable `Run`/`VA`, `CancelActive`, structured command tests | Manual plugin path-with-spaces cancellation smoke | Blocker implemented |
| 11.9 Plugin registry/uninstall | Planned | Done / Manual Windows Verification | atomic registry write, uninstall rollback tests | Manual locked executable/uninstall smoke | Blocker implemented |
| 11.10 Shutdown/startup cleanup | Verification Required | Implemented / Manual Windows Verification Required | bounded shutdown tests, plugin staging/backup startup cleanup tests | Full desktop shutdown smoke and Windows junction checks | Implemented with manual residual risk |
| 11.11 Untitled recovery journal | Verification Required | Done / Manual Crash Verification | recovery journal, warning-once tests | Manual hard-kill crash recovery smoke | Conditional blocker implemented |
| 11.12 Browser toolbar ACK/overflow | Planned | Done / Manual Recorder Stress Verification | command ACK/overflow tests, toolbar stop tests | Manual rapid-click recorder stress | Post-v1.0 risk resolved for v1.0 |

## Scope

```text
internal/plugin
internal/gui/file_operation_service.go
internal/gui/features.go
internal/gui/path_guard.go
internal/paths
internal/player/browser_pool.go
internal/gui/service_shutdown.go
internal/gui/startup_cleanup.go
internal/selector/browser_toolbar.js
internal/settings
frontend/src/controllers/workspaceSessionController.ts
frontend/src/stores/sessionStore.ts
```

# 11.1. Plugin Filesystem Confinement and ID Validation

**Priority:** Critical
**Status:** Done / Manual Windows Verification
**Area:** Plugin Store / Plugin Installer / Filesystem Safety
**Parent Section:** 10. Architecture Cleanup / PluginService
**Deferrable or v1.0 blocker:** v1.0 blocker

## Problem

Resolved. Plugin IDs are validated before filesystem or registry mutation, and existing addon paths are confined through `PathGuard`.

## User impact

A malicious or malformed plugin name can delete or overwrite data outside the Scenaria-owned addons directory.

## Confirmed code paths

```text
internal/plugin/identity.go ValidatePluginID
internal/plugin/identity.go addonPath
internal/plugin/install.go FetchAndInstall
internal/plugin/registry.go Install/Uninstall
internal/plugin/identity_test.go
```

## Required invariants

* [x] A plugin ID never resolves outside the exact Scenaria-owned addons root.
* [x] `RemoveAll` is never called until the destination has passed backend validation and confinement.
* [x] Install, update, load descriptor, run, enable/disable and uninstall use the same validator.

## Tasks

* [x] Add shared backend `ValidatePluginID`.
* [x] Permit only documented safe IDs: letters, digits, dot, dash, underscore.
* [x] Reject empty, `.`, `..`, separators, rooted paths, UNC paths, Windows volume-qualified paths, control characters, trailing spaces/dots.
* [x] Add addons-root confinement helper that validates the joined path after cleaning.
* [x] Apply validation before every plugin filesystem or registry operation.
* [x] Handle existing symlinks/junctions in the addons path through the shared PathGuard work.

## Tests

* [x] Valid plugin IDs.
* [x] Slash and backslash traversal.
* [x] Absolute Windows path, UNC path, and different volume.
* [x] Dot, double-dot, trailing dot/space.
* [x] Malformed ID leaves existing files untouched.
* [x] Uninstall cannot escape addons.

## Acceptance criteria

* [x] Every plugin path operation is confined to `projectRoot/addons/<pluginID>`.
* [x] Invalid plugin IDs fail before filesystem mutation.
* [x] Existing valid plugin files are untouched after malformed input.

## Manual verification

* [ ] Try installing/uninstalling `../outside`, `..\outside`, `C:\outside`, `\\server\share`, `.`, `..`, and names with trailing spaces/dots.

## Dependencies

* [x] Coordinate with 11.6 Canonical PathGuard for symlink/junction protection.

# 11.2. Transactional Plugin Installation and Update Rollback

**Priority:** Critical / High
**Status:** Done / Manual Windows Verification
**Area:** Plugin Installer / Data Safety
**Parent Section:** 10. Architecture Cleanup / PluginService
**Deferrable or v1.0 blocker:** v1.0 blocker

## Problem

Resolved. `FetchAndInstall` stages downloads/extraction first, validates the staged descriptor, moves an existing plugin to backup only after staging succeeds, commits the staged directory, and rolls back on commit/registry failure.

## User impact

A failed update can remove a previously working plugin and leave no usable version.

## Confirmed code paths

```text
internal/plugin/install.go FetchAndInstall
internal/plugin/install.go createPluginStagingDir
internal/plugin/install.go commitPluginInstall
internal/plugin/install.go rollbackPluginInstall
internal/plugin/install_transaction_test.go
```

## Required invariants

* [x] A failed install or update leaves the previous valid plugin fully usable.
* [x] New plugin files are committed only after package and descriptor validation succeed.
* [x] Staging, backup, and final directories remain inside Scenaria-owned storage.

## Tasks

* [x] Download/copy/extract into a Scenaria-owned staging directory.
* [x] Validate archive paths, descriptor and plugin ID in staging. API compatibility gate is still descriptor-schema dependent.
* [x] Move existing version to backup only after staging is ready.
* [x] Move staging to final atomically where possible.
* [x] Restore backup on failure and report rollback failures explicitly.
* [x] Clean abandoned staging/backup directories safely.
* [x] Log backup cleanup failure after a successful update without failing the already-committed install.
* [x] Make stale plugin backup directories eligible for confined startup cleanup.

## Tests

* [x] Successful new install and update.
* [x] Failed download, corrupt archive, missing descriptor, invalid descriptor.
* [ ] Incompatible API version if supported by descriptor schema.
* [x] Failed final rename and rollback coverage.
* [x] Old plugin remains usable after failed update cases.
* [x] No abandoned staging directory after normal failure.
* [x] Backup cleanup success and cleanup failure after successful update.
* [x] Stale backup cleanup, fresh backup preservation, malformed names and symlink escape.

## Acceptance criteria

* [x] No update failure removes the previous valid plugin.
* [x] No partial extraction is visible as an installed plugin.
* [x] Registry is updated only after the filesystem commit succeeds.

## Manual verification

* [ ] Install a valid plugin, then update from a broken zip and confirm the old plugin still runs.

## Dependencies

* [x] Depends on 11.1 plugin ID validation.
* [x] Coordinate registry commit ordering with 11.9.

# 11.3. Safe Duplicate, Import and Destination Path Validation

**Priority:** High
**Status:** Done / Manual Windows Verification
**Area:** File Operations / Feature Management
**Parent Section:** 5. Stabilize Backend Storage, Locks and File Operations
**Deferrable or v1.0 blocker:** v1.0 blocker

## Problem

Resolved. Duplicate/rename/import destination names are normalized, Windows-reserved names and path escapes are rejected, final paths are confined through `PathGuard`, and copy uses exclusive creation with partial-copy cleanup.

## User impact

A duplicate/import operation can escape the intended directory or silently overwrite an existing file.

## Confirmed code paths

```text
internal/gui/file_operation_service.go:98 DuplicateFeature
internal/gui/file_operation_service.go:116 target = filepath.Join(dir, name+ext)
internal/gui/file_operation_service.go:125 fixed duplicate collision loop
internal/gui/file_operation_service.go:137 os.WriteFile(target, payload, 0o644)
internal/gui/file_operation_service.go:180 ImportFeatures
internal/gui/file_operation_service.go:198 uniqueFeaturePath(...)
internal/gui/file_operation_service.go:349 uniqueFeaturePath
internal/gui/file_operation_service.go:371 os.Create(dest)
```

## Required invariants

* [x] Duplicate/import creates only a new `.feature` file inside the intended project directory.
* [x] New destination files are created exclusively.
* [x] Existing files are never truncated or overwritten silently.

## Tasks

* [x] Add shared backend filename validation for create/rename/duplicate/import.
* [x] Reject separators, rooted paths, UNC paths, volume-qualified paths, Windows reserved names, control characters, trailing dot/space.
* [x] Validate final destination with PathGuard after joining.
* [x] Replace fixed collision fallback with bounded-safe candidate generation.
* [x] Create new destinations with `O_CREATE|O_EXCL`.
* [x] Remove partial destination on copy failure.

## Tests

* [x] Traversal names and Windows absolute/UNC names.
* [x] Reserved names and trailing dot/space.
* [x] Valid names.
* [x] More than 100 collisions.
* [x] Concurrent collision.
* [x] Partial copy cleanup.
* [x] Existing destination content remains unchanged.

## Acceptance criteria

* [x] No duplicate/import can escape project confinement.
* [x] Collision exhaustion returns an error or a guaranteed unused name.
* [x] Save after duplicate/import cannot recreate or overwrite an unintended path.

## Manual verification

* [ ] Duplicate with valid and invalid names on Windows.
* [ ] Import files whose names collide with many existing copies.

## Dependencies

* [x] Depends on 11.6 for canonical destination confinement.

# 11.4. Transactional Replace-in-Project with Rollback

**Priority:** High
**Status:** Done / Manual Windows Verification
**Area:** Project-Wide Editing / Data Safety
**Parent Section:** 5. Stabilize Backend Storage, Locks and File Operations
**Deferrable or v1.0 blocker:** v1.0 blocker

## Problem

Resolved. `ReplaceInProject` builds a full plan before commit, validates current bytes before every write to detect external edits, uses atomic writes, and rolls back committed files on later failure.

## User impact

Project-wide replace can leave the project partially modified after a lock, permission error, antivirus block, or disk-full condition.

## Confirmed code paths

```text
internal/gui/features.go:23 ReplaceInProject
internal/gui/features.go:42 read error continue
internal/gui/features.go:49 os.WriteFile(file, []byte(replaced.Text), 0o644)
```

## Required invariants

* [x] Replace-in-project is all-or-nothing for committed files, with explicit rollback failure reporting.
* [x] Read failures are not silently ignored for files selected by discovery.
* [x] Original contents are recoverable after commit failure when rollback succeeds.

## Tasks

* [x] Add planning stage that reads all candidate files before modifying any file.
* [x] Compute replacements in memory.
* [x] Prepare atomic writes per changed file.
* [x] Commit with atomic replacement and rollback of committed files.
* [x] Roll back all previously committed files if a later commit fails.
* [x] Report rollback failures explicitly.

## Tests

* [x] Normal multi-file replacement and no-match result.
* [x] Read failure before commit.
* [x] Temp write failure.
* [x] Failure on second/final commit.
* [x] Rollback success and rollback failure.
* [x] External modification, same-size edit and delete-before-commit simulation.

## Acceptance criteria

* [x] Failed replace does not leave silent partial edits.
* [x] Result counts reflect only committed changes.
* [x] User receives a precise error and recovery status.

## Manual verification

* [ ] Run replace while one target file is locked on Windows.

## Dependencies

* [x] Reuse existing atomic feature write helper where safe.

# 11.5. Project-Independent Untitled Session Recovery

**Priority:** High
**Status:** Done / Manual Windows Verification
**Area:** Workspace Session / Untitled Recovery
**Parent Section:** 1. Stabilize Monaco Tabs and Editor State
**Deferrable or v1.0 blocker:** v1.0 blocker

## Problem

Resolved. Untitled snapshots and recovery-journal entries are merged and restored before best-effort project open. Project resolution/open failures preserve restored Untitled tabs and activate a deterministic fallback.

## User impact

Unsaved text can appear lost when the previous project was moved, deleted, or is on a disconnected drive.

## Confirmed code paths

```text
frontend/src/controllers/workspaceSessionController.ts:109 restoreWorkspaceSession
frontend/src/controllers/workspaceSessionController.ts:119 openProject(resolvedProj)
frontend/src/controllers/workspaceSessionController.ts:124 projectNotFound log/status
frontend/src/controllers/workspaceSessionController.ts:126 return
frontend/src/controllers/workspaceSessionController.ts:130-157 Untitled restore happens only after project open succeeds
```

## Required invariants

* [x] Missing project never hides valid Untitled content.
* [x] Valid Untitled documents restore before or independently from project open.
* [x] Active fallback selection is deterministic.
* [x] Intentionally empty Untitled documents remain empty.
* [x] Recorder-generated Untitled content survives restart.

## Tasks

* [x] Restore and deduplicate Untitled tabs before attempting project open.
* [x] Choose a valid active tab from restored documents before project failure can abort.
* [x] Continue with project open as best effort.
* [x] If project open fails, keep Untitled visible and show a recovery warning.
* [x] Avoid creating replacement empty Untitled tabs.

## Tests

* [x] One Untitled with missing project.
* [x] Multiple Untitled tabs with invalid active tab.
* [x] Intentionally empty Untitled.
* [x] Recorder-created Untitled.
* [x] Saved feature plus Untitled with project moved/deleted/open failure.

## Acceptance criteria

* [x] Untitled content restores even when the saved project cannot be opened.
* [x] No duplicate empty Untitled is created.
* [x] User gets a clear project-not-found warning.

## Manual verification

* [ ] Save dirty Untitled, move/delete previous project, reopen app.

## Dependencies

* [x] Must preserve Monaco hydration invariants from Section 1.

# 11.6. Canonical PathGuard and Reparse-Point Confinement

**Priority:** High
**Status:** Implemented / Manual Windows Verification Required
**Area:** Filesystem Confinement / Windows Paths
**Parent Section:** 5. Stabilize Backend Storage, Locks and File Operations
**Deferrable or v1.0 blocker:** v1.0 blocker for destructive operations

## Problem

Project confinement is implemented in several places with lexical `Abs`/`Rel` checks. Existing targets and new-file parent directories are not consistently canonicalized with symlink/junction awareness.

## User impact

A path that looks inside a project can physically target data outside the project through a symlink, Windows junction, or reparse point.

## Confirmed code paths

```text
internal/paths/paths.go:52 ConfineToProjectRoot
internal/paths/paths.go:57 filepath.Abs(projectRoot)
internal/paths/paths.go:67 filepath.Rel(absRoot, absPath)
internal/gui/file_operation_service.go:295 ensureInsideProject
internal/gui/file_operation_service.go:307 filepath.Abs(target)
internal/gui/file_operation_service.go:315 filepath.Rel(projectAbs, abs)
internal/gui/path_guard.go confineFeaturePath/confineArtifactPath
```

## Required invariants

* [x] Existing targets are checked using resolved canonical paths.
* [x] New files validate the resolved canonical parent directory.
* [x] Final resolved path remains within the resolved allowed root.
* [x] Different Windows volumes and UNC escapes are rejected by path validation where applicable.
* [x] `..hidden` is not confused with parent traversal.
* [x] Case-insensitive Windows paths are handled by filesystem-aware comparisons where needed.

## Tasks

* [x] Introduce a shared `PathGuard` with `ResolveExisting`, `ResolveNewFile`, and `EnsureContained`.
* [x] Use `EvalSymlinks` for existing targets and existing parent directories.
* [x] Replace ad hoc `Abs`/`Rel` helpers in GUI file operations and artifact paths.
* [x] Apply the same guard to plugin addons paths where existing filesystem entries are involved.
* [x] Keep non-existing final files supported when the canonical parent is inside the root.

## Tests

* [x] Normal child, sibling, `..hidden`, and `../outside`.
* [x] Symlink inside root pointing outside.
* [ ] Windows junction/reparse point pointing outside. Manual Windows verification still required.
* [x] Symlinked parent for a new file.
* [x] Different volume and UNC path where representable by the test platform.
* [x] Case variations on Windows-sensitive code paths.

## Acceptance criteria

* [x] Destructive operations cannot cross the physical project/addons root.
* [x] Valid project paths continue to work.
* [ ] Windows junction/reparse behavior still needs manual verification.

## Manual verification

* [ ] Create a Windows junction under a project and try read/save/delete/move/import through it.

## Dependencies

* [x] Should precede or accompany 11.3 for final destination safety.

# 11.7. Browser Pool Close and Late-Release Synchronization

**Priority:** High
**Status:** Done / Manual Windows Verification
**Area:** Parallel Runner / Browser Pool / Shutdown
**Parent Section:** 8. Stabilize Playwright Runner Lifecycle
**Deferrable or v1.0 blocker:** v1.0 blocker if `workers > 1` remains enabled

## Problem

Resolved. `browserPool.release` uses `returnSlot` with close detection, retires/stops rejected slots, prevents retired slots from being reacquired, and makes worker stop idempotent.

## User impact

A browser worker can leak or be reused after shutdown/cancel, causing inconsistent parallel runs or stuck Playwright processes.

## Confirmed code paths

```text
internal/player/browser_pool.go:91 release
internal/player/browser_pool.go:103 initial closed check
internal/player/browser_pool.go:114 resetForScenario
internal/player/browser_pool.go:120 replaceSession
internal/player/browser_pool.go:133 p.slots <- slot
internal/player/browser_pool.go:224 Close
internal/player/browser_pool.go:251 time.After(2 * time.Second)
```

## Required invariants

* [x] No slot or replacement worker re-enters the pool after `Close` starts.
* [x] Late release after close stops its worker instead of sending it to the pool.
* [x] Double close, rejected returns and concurrent releases do not panic or leak worker ownership.

## Tasks

* [x] Add a post-reset/post-replacement closed check via `returnSlot`.
* [x] Stop the slot if the pool closed during slow release work.
* [x] Ensure retired-slot accounting remains correct.
* [x] Keep close bounded and log late-release cleanup.
* [x] Retire and stop a slot if the pool channel unexpectedly rejects return.

## Tests

* [x] Close during reset.
* [x] Close during replacement.
* [x] Close during unhealthy release.
* [x] Several concurrent releases.
* [x] Double `Close`.
* [x] Late release closes worker and does not send to the pool.
* [x] Channel unexpectedly full, double release, rejected replacement and stop-once behavior.

## Acceptance criteria

* [x] No send occurs after close starts.
* [x] No browser process remains owned by a closed pool in tested ownership paths.
* [x] Existing close-browser slot replacement tests still pass.

## Manual verification

* [ ] Run repeated cancel/shutdown during `workers = 2/4` batches.

## Dependencies

* [x] Existing browser pool replacement/retirement tests are already present and should be extended.

# 11.8. Plugin Process Ownership, Cancellation and Structured Arguments

**Priority:** High
**Status:** Done / Manual Windows Verification
**Area:** Plugin Runtime / Process Lifecycle / Windows Compatibility
**Parent Section:** 10. Architecture Cleanup / PluginService
**Deferrable or v1.0 blocker:** v1.0 blocker for plugin runtime release

## Problem

Resolved for v1.0. Plugin and Vanessa executions use registered cancellable invocations, project switch/shutdown/uninstall cancellation routes call `CancelActive`/`CancelPlugin`, and legacy command parsing now handles quoting. Structured command specs are supported.

## User impact

Plugin processes can outlive app shutdown or project switches, and valid Windows paths such as `C:\Program Files\...` can be parsed incorrectly.

## Confirmed code paths

```text
internal/gui/plugin_service.go pluginCLIRunner Run/VA(ctx,...)
internal/gui/plugin_service.go startInvocation / CancelActive / CancelPlugin
internal/plugin/runner.go structuredCommands / splitLegacyCommand
internal/gui/plugins_test.go
```

## Required invariants

* [x] Plugin execution is tied to app/project/plugin invocation cancellation.
* [x] Only Scenaria-owned plugin processes are tracked and stopped.
* [x] No global image-name process kill is used.
* [x] Executable and arguments with spaces are represented without shell reparsing for structured specs and quoted legacy commands.

## Tasks

* [x] Add plugin invocation context hierarchy.
* [x] Track owned plugin invocations for shutdown/cancellation.
* [x] Add cancellation paths for project switch, uninstall, and app close.
* [x] Extend plugin descriptor schema with structured command/args while preserving compatibility where needed.
* [x] Safely parse quoted legacy command strings.

## Tests

* [x] Explicit cancellation, project switch, app close, plugin-specific cancellation.
* [x] Hanging plugin, non-zero exit, large stdout/stderr.
* [x] Executable path with spaces.
* [x] Argument with spaces, Unicode path, empty argument, literal quote.
* [x] Backward compatibility for existing descriptors.

## Acceptance criteria

* [x] No tracked plugin invocation remains after cancellation/shutdown.
* [x] Quoted Windows paths and arguments are parsed correctly in tests.
* [x] Legacy descriptors fail safely or are parsed deterministically.

## Manual verification

* [ ] Run a plugin from a path containing spaces and cancel it from the GUI.

## Dependencies

* [x] Depends on 11.1 for plugin identity validation.

# 11.9. Atomic Plugin Registry Persistence and Safe Uninstall

**Priority:** Medium / High
**Status:** Done / Manual Windows Verification
**Area:** Plugin Registry / Uninstall / Durability
**Parent Section:** 10. Architecture Cleanup / PluginService
**Deferrable or v1.0 blocker:** v1.0 blocker for plugin store release

## Problem

Resolved. Plugin registry writes are atomic with backup/restore fallback, uninstall validates plugin ID, moves plugin files to backup before registry update, rolls back on registry failure, and reports cleanup failure explicitly.

## User impact

A failed registry write can corrupt installed-plugin state. Uninstall can leave stale executables and a registry/disk mismatch.

## Confirmed code paths

```text
internal/plugin/registry.go SaveManifest / writeRegistryAtomic / Uninstall
internal/plugin/registry_test.go
internal/gui/plugin_service.go PluginService.Uninstall
```

## Required invariants

* [x] Failed registry update preserves the previous valid registry.
* [x] Uninstall validates plugin ID and removes registry and files in a recoverable order.
* [x] Partial uninstall is reported explicitly.

## Tasks

* [x] Replace registry `os.WriteFile` with proven atomic JSON write behavior.
* [x] Define corrupted-registry behavior: return decode error, preserve file.
* [x] Stop any running plugin before uninstall.
* [x] Remove plugin files from a confined addons path.
* [x] Decide rollback/staging order for registry-vs-files uninstall.

## Tests

* [x] Normal registry write and update.
* [x] Temp write failure and replacement failure.
* [x] Old registry preserved.
* [x] Malformed existing registry.
* [x] Normal uninstall, missing directory, missing registry entry.
* [ ] Locked executable/running plugin manual Windows verification remains.
* [x] Malicious plugin ID.
* [x] Registry update failure and filesystem removal failure.

## Acceptance criteria

* [x] Registry cannot be truncated by a failed write.
* [x] Successful uninstall leaves no runnable plugin files.
* [x] Failed uninstall does not silently desynchronize registry and disk.

## Manual verification

* [ ] Install, uninstall, reinstall a plugin; confirm the addon directory state.

## Dependencies

* [x] Depends on 11.1 and should align with 11.2 commit ordering.

# 11.10. End-to-End Bounded Shutdown and Startup Cleanup Hardening

**Priority:** Medium / High
**Status:** Implemented / Manual Windows Verification Required
**Area:** Application Shutdown / Process Cleanup / Startup Temp Cleanup
**Parent Section:** 8. Stabilize Playwright Runner Lifecycle
**Deferrable or v1.0 blocker:** v1.0 blocker if unbounded waits are reproduced

## Problem

Mostly resolved. Bounded shutdown cleanup exists, plugin invocations are cancelled on shutdown, and startup cleanup covers stale project/global temps plus plugin staging/backup directories with symlink escape tests. Manual Windows junction and full desktop shutdown smoke remain.

## User impact

App exit can still hang after the first timeout, and crash leftovers outside the currently covered temp namespaces can accumulate.

## Confirmed code paths

```text
internal/gui/service_shutdown.go:31 Shutdown
internal/gui/service_shutdown.go:78 cleanupTempFeatureDirs()
internal/gui/service_shutdown.go:79 stopAllureServe()
internal/gui/service_shutdown.go:80 playwrightrt.Shutdown()
internal/gui/startup_cleanup.go:14 startupTempMaxAge
internal/gui/startup_cleanup.go:16 CleanupStartupTempArtifacts
internal/gui/startup_cleanup.go:25 CleanupGlobalStartupTemps
```

## Required invariants

* [x] One shutdown deadline covers cleanup steps through bounded wrappers where implemented.
* [x] Each cleanup operation accepts context or has a smaller bounded timeout where practical.
* [x] Only known Scenaria-owned processes are stopped.
* [x] Startup cleanup deletes only known stale Scenaria-owned temp namespaces.
* [x] Symlink escapes are rejected; Windows junction verification remains manual.

## Tasks

* [x] Thread shutdown cancellation through run/recorder/plugin cleanup paths currently covered by services.
* [x] Log structured remaining resources after deadline.
* [x] Add bounded wrappers for cleanup calls that cannot accept context yet.
* [x] Extend startup cleanup to plugin staging/backup temp directories after transactional install is implemented.
* [x] Add symlink checks for startup cleanup paths; junction manual verification remains.

## Tests

* [x] Hanging active run, recorder, plugin cancellation and bounded cleanup paths.
* [x] Repeated `Shutdown` calls.
* [x] Normal shutdown remains fast.
* [x] Stale temp removed, fresh/unknown/active temp preserved.
* [x] Malformed path ignored.
* [x] Symlink escape rejected; junction escape manual verification remains.

## Acceptance criteria

* [x] Shutdown cannot begin a new unbounded wait after its deadline expires in covered cleanup paths.
* [x] Logs identify remaining owned resources.
* [x] Startup cleanup does not remove user data or active artifacts.

## Manual verification

* [ ] Close the app during a run/recording/plugin execution and verify shutdown completes or reports bounded timeout.

## Dependencies

* [x] Plugin staging cleanup depends on 11.2.
* [x] Symlink checks depend on 11.6. Junction verification remains manual.

# 11.11. Crash-Resilient Untitled Recovery Journal

**Priority:** Medium / High
**Status:** Done / Manual Crash Verification
**Area:** Workspace Autosave / Data Safety
**Parent Section:** 1. Stabilize Monaco Tabs and Editor State
**Deferrable or v1.0 blocker:** Conditional; blocker only if strict crash recovery is required for v1.0

## Problem

Implemented for v1.0. Untitled tabs still persist through workspace settings on controlled close, and a frontend recovery journal now captures recent Untitled text for crash recovery with warning-once visibility when storage is unavailable.

## User impact

The most recent Untitled edits can be lost after crash or Task Manager termination even when normal close restores correctly.

## Confirmed code paths

```text
frontend/src/stores/sessionStore.ts:6 schedulePersist(run, delayMs = 500)
frontend/src/lib/sessionTabs.ts:16 buildSessionTabsSnapshot
frontend/src/controllers/workspaceSessionController.ts:69 buildSessionTabsSnapshot(...)
internal/settings/settings.go:57 UntitledTabSession
```

## Required invariants

* [x] Controlled close still flushes all Untitled text.
* [x] Crash recovery restores recently edited Untitled content within the accepted durability window.
* [x] Explicit Discard and Save As clean obsolete recovery data.
* [x] Intentionally empty Untitled content is preserved.

## Tasks

* [x] Decide whether current 500 ms settings debounce is acceptable for v1.0.
* [x] Add a small per-Untitled recovery journal.
* [x] Define precedence between session snapshot, feature drafts, and Untitled journal.
* [x] Clean journal on Save As and explicit Discard.
* [x] Surface journal write failures once per app session.

## Tests

* [x] Edit and immediate simulated crash.
* [x] Journal write failure.
* [x] Explicit Discard.
* [x] Save As.
* [x] Multiple Untitled documents.
* [x] Stale/corrupted journal.
* [x] Session and journal precedence.

## Acceptance criteria

* [x] The chosen durability guarantee is documented and tested.
* [x] A crash cannot resurrect explicitly discarded Untitled text.
* [x] Saved feature tabs remain unaffected.

## Manual verification

* [ ] Type into Untitled and terminate the process immediately; reopen and verify expected recovery behavior.

## Dependencies

* [x] Should be implemented after 11.5 so project-independent restore is stable.

# 11.12. Browser Toolbar Command Acknowledgement and Overflow Safety

**Priority:** Medium
**Status:** Done / Manual Recorder Stress Verification
**Area:** Injected Browser Toolbar / Recorder Lifecycle
**Parent Section:** 6. Stabilize Recorder Lifecycle
**Deferrable or v1.0 blocker:** Post-v1.0 unless normal usage drops Stop/Pause

## Problem

Resolved for v1.0. Toolbar commands now carry IDs, require ACK, deduplicate repeated commands, keep Stop prioritized, and reject non-terminal overflow without silently dropping the oldest command.

## User impact

Rapid toolbar actions can silently lose an important command such as Stop or Pause.

## Confirmed code paths

```text
internal/selector/browser_toolbar.js:4 MAX_QUEUE = 8
internal/selector/browser_toolbar.js:197 queue overflow
internal/selector/browser_toolbar.js:198 queue.shift()
internal/selector/browser_toolbar.js:200 queue.push(action)
internal/selector/browser_toolbar.js:276 takeAction()
```

## Required invariants

* [x] Stop is never silently discarded.
* [x] Accepted commands are acknowledged or reflected in busy state.
* [x] Incompatible commands are disabled while pending.
* [x] Repeated identical commands are deduplicated where appropriate.

## Tasks

* [x] Define priority/overflow policy for Stop, Pause, Record, Picker.
* [x] Add command acknowledgement state.
* [x] Deduplicate repeated identical commands.
* [x] Surface queue-full/busy feedback instead of silent drop.

## Tests

* [x] Rapid Pause/Stop.
* [x] Repeated Pause and Record.
* [x] Stop under full queue.
* [x] Command ACK.
* [x] Navigation/click preservation around toolbar Stop.

## Acceptance criteria

* [x] Terminal commands cannot be lost silently.
* [x] Toolbar state remains consistent with backend recorder state in tested paths.

## Manual verification

* [ ] Rapid-click toolbar controls during recording and confirm Stop always wins.

## Dependencies

* [x] Coordinate with recorder session identity and toolbar polling code.

# 11.13. Credential Storage and Secret Redaction Audit

**Priority:** High
**Status:** Implemented - documented plaintext storage for v1.0, accidental log/journal redaction added
**Area:** Settings / Credentials / Security
**Parent Section:** 5. Stabilize Backend Storage, Locks and File Operations
**Deferrable or v1.0 blocker:** Security decision required before v1.0

## Problem

HTTP auth passwords are persisted as plaintext JSON settings. UI DTOs hide the password on read, but the at-rest storage is not protected.

## User impact

Users can unintentionally store secrets in ordinary settings files. A local file disclosure exposes HTTP credentials.

## Confirmed code paths

```text
internal/settings/settings.go:43 HTTPAuth map json:"http_auth"
internal/settings/settings.go:62 HTTPAuthEntry
internal/settings/settings.go:64 Password string json:"password"
internal/httpauth/httpauth.go:122 StoreHostCredentials
internal/gui/settings_service.go:99 SaveHTTPAuth
internal/gui/settings_service.go:209 SaveSettings preserves existing HTTPAuth
```

## Required invariants

* [x] Secrets are not logged or sent to frontend logs/status except as explicit user input.
* [x] At-rest credential behavior is documented before release.
* [x] If protected storage is required, migration preserves existing credentials safely. Deferred: protected storage is not selected for v1.0.

## Tasks

* [x] Audit logs, reports, debug dumps, workspace snapshots and generated fixtures for secret exposure.
* [x] Decide v1.0 policy: documented plaintext, opt-in storage, or OS secret store.
* [x] If protected storage is chosen, design migration to Windows Credential Manager/DPAPI, macOS Keychain, and Linux Secret Service. Deferred: protected storage is not selected for v1.0.
* [x] Ensure imported URL credentials are stripped from URLs before logs/reports.

## Tests

* [x] Password omitted from read DTOs.
* [x] Password not present in logs/frontend journal other than chosen storage or explicit user-authored content.
* [x] URL credentials stripped before persistence where expected.
* [x] Migration tests if secret store is implemented. Not applicable for documented plaintext v1.0 policy.

## Acceptance criteria

* [x] Security posture is explicit and reviewed.
* [x] No accidental secret leak is found in logs/frontend journal.
* [x] At-rest behavior matches documentation.

## Manual verification

* [ ] Save HTTP auth credentials, inspect settings/logs/reports, and confirm expected exposure only.

## Dependencies

* [ ] None for audit; OS secret-store migration can be deferred only by explicit release decision.

# Final Implementation Order

Идём строго сверху вниз:

```text
0. Add Safety Tests Before Refactoring
1. Stabilize Monaco Tabs and Editor State
2. Stabilize Batch Execution
2.1. Stabilize Parallel Batch Workers
3. Introduce ProjectSession, RunSession and Operation Identity
4. Unify StepMatcher and EditorAnalysisService
4.1. Stabilize Step Execution Semantics
4.2. Stabilize Runtime Step Semantics and Reporting
5. Stabilize Backend Storage, Locks and File Operations
6. Stabilize Recorder Lifecycle
6.1. Stabilize Element Picker and Selector Generation
6.1.1. Component-Library Selector Heuristics
6.2. Stabilize Browser Selector Validation
6.2.1. Validation Diagnostics in HTML Report
6.3. Stabilize Recorder Event Ordering and Navigation Causality
7. Stabilize Reports and Artifacts
8. Stabilize Playwright Runner Lifecycle
8.1. Stabilize Browser Context Isolation
9. Performance and Memory Optimization
10. Architecture Cleanup
11. Release Safety Follow-Up
11.1. Plugin Filesystem Confinement and ID Validation
11.2. Transactional Plugin Installation and Update Rollback
11.3. Safe Duplicate, Import and Destination Path Validation
11.4. Transactional Replace-in-Project with Rollback
11.5. Project-Independent Untitled Session Recovery
11.6. Canonical PathGuard and Reparse-Point Confinement
11.7. Browser Pool Close and Late-Release Synchronization
11.8. Plugin Process Ownership, Cancellation and Structured Arguments
11.9. Atomic Plugin Registry Persistence and Safe Uninstall
11.10. End-to-End Bounded Shutdown and Startup Cleanup Hardening
11.11. Crash-Resilient Untitled Recovery Journal
11.12. Browser Toolbar Command Acknowledgement and Overflow Safety
11.13. Credential Storage and Secret Redaction Audit
```

---

# Consolidated Definition of Done

Задача считается выполненной только если:

* [ ] исправлен root cause, а не только симптом;
* [ ] добавлен regression test или понятный manual QA сценарий;
* [ ] async result/event защищён от stale применения;
* [ ] path/session/run/project identity не теряется;
* [ ] mutable input копируется на async/backend boundaries;
* [ ] нет silent behavior change для пользователя;
* [ ] добавлены logs/metrics, если баг сложно диагностировать;
* [ ] обновлён ROADMAP/checklist при изменении scope.

---

# v1.0 Release Criteria

Scenaria можно считать готовой к v1.0-stable, когда:

* [ ] вкладки не теряют и не смешивают содержимое;
* [ ] batch run выполняет ровно актуально выбранные тесты;
* [ ] IDE и runner одинаково понимают known/unknown steps;
* [ ] project/run/record/report events имеют identity;
* [ ] stale events игнорируются;
* [ ] run inputs immutable;
* [ ] recorder пишет только в explicit target file/model;
* [ ] reports изолированы per run;
* [ ] report writes atomic;
* [ ] scenario outline examples имеют стабильную identity;
* [ ] canceled run отображается отдельно от failed;
* [ ] runstatus/settings/file operations защищены от races;
* [ ] repeated cancel не приводит к unbounded goroutine/browser leaks;
* [ ] large files и large reports имеют измеримую и приемлемую производительность;
* [ ] targeted `go test -race` проходит для sensitive backend packages;
* [ ] frontend regression tests проходят для tabs, batch, validation, recorder.
