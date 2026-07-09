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

* [ ] Открыть 5 файлов.
* [ ] Быстро переключаться между ними.
* [ ] Внести разные изменения в каждый файл.
* [ ] Закрыть активную вкладку из середины.
* [ ] Закрыть первую и последнюю вкладку.
* [ ] Закрыть Welcome tab.
* [ ] Проверить undo/redo по вкладкам.
* [ ] Проверить save после переключения.
* [ ] Проверить external reload после переключения.

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

* [ ] Выбрать 2 теста и запустить.
* [ ] Добавить ещё 3 теста и запустить.
* [ ] Проверить, что выполняются все 5.
* [ ] Убрать часть тестов и запустить.
* [ ] Проверить запуск через hotkey сразу после включения batch mode.
* [ ] Проверить folder selection.
* [ ] Проверить refresh проекта после удаления файла.

---

# 2.1. Stabilize Parallel Batch Workers

**Priority:** High
**Status:** Planned
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

* [ ] Клонировать `RunRequest.Vars` перед построением execution plan.
* [ ] Клонировать variables для каждого `RunCase`.
* [ ] Клонировать variables при создании `RunContext`.
* [ ] Запретить shared mutable map между parallel scenarios.
* [ ] Добавить `go test -race` для parallel scenarios с `Remember()`.

### Worker execution consistency

* [ ] Добавить regression test: batch с `workers = 1`.
* [ ] Добавить regression test: тот же batch с `workers = 2`.
* [ ] Добавить regression test: тот же batch с `workers = 4`.
* [ ] Проверить, что selected cases count одинаковый при любом количестве workers.
* [ ] Проверить, что каждый selected case выполняется ровно один раз.
* [ ] Проверить, что runner queue не теряет cases при parallel execution.

### Browser context policy

* [ ] Явно определить browser context policy:

  * isolated context per scenario;
  * reuse context per worker;
  * shared auth state only by explicit setting.
* [ ] Задокументировать разницу между `workers = 1` и `workers > 1`.
* [ ] Если context переиспользуется внутри worker-а, явно очищать или документировать:

  * cookies;
  * localStorage;
  * sessionStorage;
  * opened tabs;
  * permissions;
  * downloads;
  * network routes.
* [ ] Если auth/session state должен переиспользоваться, сделать это explicit setting.

### Fail-fast behavior

* [ ] Проверить поведение `ContinueOnFail = false`.
* [ ] Проверить поведение `ContinueOnFail = true`.
* [ ] В UI явно показывать настройку:

  * “Stop on first failure”;
  * “Continue on fail”.
* [ ] В report разделять:

  * failed;
  * canceled;
  * skipped;
  * not started.
* [ ] Не показывать canceled scenarios как обычные failed.
* [ ] Не создавать ощущение, что batch “потерял” тесты.

### Progress events

* [ ] Frontend progress handling должен быть order-independent.
* [ ] Не считать, что progress events придут в порядке `1,2,3,4`.
* [ ] Использовать `caseId` / `index` / `runId` для обновления конкретного case.
* [ ] Progress bar должен считать completed count, а не последний пришедший index.
* [ ] Добавить debug logs для:

  * scheduled cases;
  * started cases;
  * finished cases;
  * canceled cases;
  * failed cases.

## Acceptance Criteria

* [ ] Batch с `workers = 1`, `workers = 2`, `workers = 4` выполняет одинаковый набор selected cases.
* [ ] Каждый selected case выполняется ровно один раз.
* [ ] Variables не протекают между scenarios.
* [ ] `go test -race` не показывает race по variables.
* [ ] Browser context reuse policy явно определена.
* [ ] ContinueOnFail поведение понятно в UI и report.
* [ ] Progress UI корректен при out-of-order completion events.
* [ ] Canceled/not-started scenarios не выглядят как потерянные тесты.

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

* [ ] Run из project A не может обновить UI project B.
* [ ] Validation из старой версии проекта игнорируется.
* [ ] Recorder event из старой сессии игнорируется.
* [ ] Report action привязан к своему report/run/project.
* [ ] Новый run не отменяет старый silently.
* [ ] Frontend не зависает навсегда, если backend не прислал finish event.
* [ ] Все long-running events трассируются по ID.

## Manual QA

* [ ] Открыть project A.
* [ ] Запустить run.
* [ ] Быстро открыть project B.
* [ ] Проверить, что старые результаты не применились.
* [ ] Повторить с validation.
* [ ] Повторить с recorder.
* [ ] Повторить с report actions.

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

* [ ] Если runner выполняет шаг, IDE не показывает `Unknown step`.
* [ ] Diagnostics, inlay hints, steps panel и autocomplete используют совместимый источник истины.
* [ ] Один malformed scenario не ломает valid steps в другом scenario.
* [ ] `Дано/Когда/Тогда/И/Но/Допустим/*` покрыты тестами.
* [ ] `ё/е`, smart quotes, tabs, spaces, escaped selectors покрыты тестами.
* [ ] Completion items либо соответствуют executable step pattern, либо явно marked snippet-only.

## Manual QA

* [ ] Открыть файл с известными шагами.
* [ ] Проверить diagnostics.
* [ ] Проверить autocomplete.
* [ ] Проверить inlay hints.
* [ ] Внести временную синтаксическую ошибку.
* [ ] Проверить, что valid known steps не стали false unknown.
* [ ] Запустить тот же сценарий runner-ом.

---

# 4.1. Stabilize Step Execution Semantics

**Priority:** High
**Status:** Planned
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

* [ ] Клонировать `RunRequest.Vars`.
* [ ] Клонировать variables при создании `RunCase`.
* [ ] Клонировать variables при создании `RunContext`.
* [ ] Запретить shared mutable map между сценариями.
* [ ] Добавить `go test -race` для parallel variables.

### Repeat semantics

* [ ] `repeat 0` должен быть ошибкой или явно разрешённой конструкцией с warning.
* [ ] `repeat < 0` должен быть ошибкой.
* [ ] `repeat > MaxLoopIterations` должен быть ошибкой, не silent clamp.
* [ ] Overflow при parsing count должен быть ошибкой.
* [ ] Repeat должен выполнять ровно requested count.

### ForEach semantics

* [ ] Не игнорировать ошибки `InnerText`.
* [ ] Логировать selector/index/iteration при ошибке.
* [ ] Явно определить модель:

  * snapshot;
  * live DOM;
  * strict mode.
* [ ] Если элемент исчез во время итерации — вернуть понятную ошибку.
* [ ] Переменная цикла не должна получать fallback/мусорное значение.

### If / While semantics

* [ ] `EvaluateCondition` должен возвращать `(bool, error)`.
* [ ] Отличать false condition от locator/page error.
* [ ] Отличать timeout от false.
* [ ] Отличать context canceled от runtime failure.
* [ ] `while` не должен завершаться silently из-за Playwright error.

### Retry policy

* [ ] Разделить retry actions на safe и risky.
* [ ] Safe by default:

  * waits;
  * assertions;
  * visibility checks.
* [ ] Risky only opt-in:

  * click;
  * double-click;
  * fill;
  * select;
  * check;
  * uncheck;
  * download-click.
* [ ] Добавить настройки:

  * `retryWaits`;
  * `retryAssertions`;
  * `retryActions`.
* [ ] Retry attempts должны отображаться в report.

## Acceptance Criteria

* [ ] Scenario A не может изменить variables Scenario B.
* [ ] Repeat не меняет count silently.
* [ ] ForEach не скрывает DOM errors.
* [ ] While не завершает цикл из-за error как будто condition false.
* [ ] Click/fill/download не retry-ятся silently.
* [ ] Retry attempts видны в отчёте.
* [ ] `go test -race` не показывает race по variables.

---

# 4.2. Stabilize Runtime Step Semantics and Reporting

**Priority:** High
**Status:** Planned
**Area:** Step Records / Loop Reporting / Retry Reporting / Terminal Steps
**Parent Section:** `4. Unify StepMatcher and EditorAnalysisService`

## Problem

Отчёт должен объяснять не только какой шаг упал, но и:

* на какой итерации;
* на какой retry attempt;
* был ли шаг skipped/canceled;
* был ли сценарий завершён terminal browser action.

## Tasks

* [ ] Добавить `IterationPath` в `StepRecord`.
* [ ] Разделить:

  * logical step;
  * loop iteration;
  * retry attempt;
  * generated wait/assert step.
* [ ] Screenshots/artifacts должны включать iteration index.
* [ ] Negative wait duration считать ошибкой.
* [ ] Zero wait duration либо warning, либо explicit allowed.
* [ ] `close-browser` / `close-tab` не должны silently делать scenario passed при оставшихся шагах.
* [ ] Если browser closed до конца scenario:

  * remaining steps = skipped/canceled;
  * scenario status не должен быть passed.
* [ ] Добавить total action attempts limit per scenario.
* [ ] В отчёте показывать:

  * loop iteration;
  * retry attempt;
  * skipped;
  * canceled;
  * terminal browser action.

## Acceptance Criteria

* [ ] Ошибка внутри `repeat[3]` видна как `repeat[3]`, а не просто line number.
* [ ] Retry внутри loop не маскирует фактическое число попыток.
* [ ] Negative wait не превращается в `0ms`.
* [ ] `close-browser` не делает невыполненные шаги passed.
* [ ] Report показывает logical order и runtime attempts.

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

* [ ] Запустить parallel run.
* [ ] Проверить run status/history.
* [ ] Запустить run с unsaved/temp feature.
* [ ] Отменить и сразу запустить другой.
* [ ] Проверить settings save параллельно с recents/HTTP auth.
* [ ] Проверить project refresh во время save/rename/delete.

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

* [ ] Start recording в `B.feature`, активна `A.feature`.
* [ ] Сделать click/input.
* [ ] Проверить, что изменился только `B.feature`.
* [ ] Переключать вкладки во время записи.
* [ ] Stop → Start снова.
* [ ] Close browser во время picker.
* [ ] Project switch во время recording.
* [ ] Проверить отсутствие stale events.

---

# 6.1. Stabilize Element Picker and Selector Generation

**Priority:** High
**Status:** Planned
**Area:** Recorder / Picker / Selector Generation
**Parent Section:** `6. Stabilize Recorder Lifecycle`

## Problem

Picker может выбирать некорректный selector, потому что он выбирает первый подходящий selector, но не доказывает, что selector уникален и указывает именно на выбранный элемент.

## Tasks

### Candidate model

* [ ] Генерировать ranked selector candidates.
* [ ] Для каждого candidate считать:

  * score;
  * uniqueness;
  * matches count;
  * reason;
  * strategy;
  * warnings.
* [ ] UI должен показывать selector confidence.
* [ ] Пользователь должен видеть альтернативные candidates.

### Validation before return

* [ ] Проверять, что selector matches exactly one element.
* [ ] Проверять, что matched element === picked element или корректный actionable ancestor.
* [ ] Проверять visibility.
* [ ] Проверять actionability по типу действия.
* [ ] Не возвращать non-unique selector silently.

### Strategy order

* [ ] Для click предпочитать:

  * `data-testid`;
  * role/aria;
  * title;
  * stable id;
  * contextual;
  * text.
* [ ] Для input предпочитать:

  * `data-testid`;
  * id;
  * name;
  * aria;
  * label;
  * placeholder.
* [ ] Понизить score text-only selectors.
* [ ] Добавить warning для text-only selector.

### Input / label targeting

* [ ] `label[for]` должен возвращать selector control, не label.
* [ ] Nested label должен возвращать вложенный input.
* [ ] Adjacent label должен строить contextual input selector.
* [ ] Fill/select steps не должны получать selector label как primary target.

### Action-aware picker

* [ ] Определять suggested action:

  * click;
  * fill;
  * select;
  * check;
  * uncheck;
  * hover.
* [ ] Input-like elements:

  * input;
  * textarea;
  * select;
  * contenteditable;
  * role=textbox;
  * role=combobox;
  * role=spinbutton;
  * role=searchbox.
* [ ] Recorded step должен использовать suggested action.
* [ ] Пользователь может переопределить action.

### Iframe / Shadow DOM / SVG / Canvas

* [ ] Добавить iframe-aware picker.
* [ ] Same-origin iframe: выбирать внутренний элемент.
* [ ] Cross-origin iframe: честно показывать limitation.
* [ ] Добавить shadow DOM hit-test.
* [ ] SVG click нормализовать до clickable ancestor.
* [ ] Canvas selector должен иметь uniqueness/warning.
* [ ] Component-library heuristics добавить позже после базовой стабилизации.

## Acceptance Criteria

* [ ] Picker возвращает validated selector.
* [ ] Selector указывает на выбранный элемент.
* [ ] Non-unique selector не выбирается silently.
* [ ] Input получает selector input/control, не label.
* [ ] `data-testid`/role/aria предпочитаются raw text.
* [ ] Iframe behavior explicit.
* [ ] Low-confidence selector показывает warning.
* [ ] Пользователь может выбрать альтернативный selector.

---

# 6.2. Stabilize Browser Selector Validation

**Priority:** High
**Status:** Planned
**Area:** Selector Validation / Browser Validation / Dynamic UI
**Parent Section:** `6. Stabilize Recorder Lifecycle`

## Problem

Browser selector validation может давать ложную уверенность. Она может проверять selector на initial page, не выполняя flow, который делает элемент доступным. Также validation должна быть action-aware, а не только visible-check.

## Tasks

* [ ] Разделить validation modes:

  * static validation;
  * flow-aware validation.
* [ ] Static validation не считать источником истины для dynamic UI.
* [ ] В UI показывать: “validated on current/initial page only”.
* [ ] Для dynamic elements не выдавать misleading missing без контекста.
* [ ] Сделать validation action-aware:

  * click → visible + enabled + actionable;
  * fill → editable input/textarea/contenteditable;
  * select → select/combobox;
  * check → checkbox/radio;
  * upload → input[type=file].
* [ ] Для chained selectors показывать matches count.
* [ ] Для hover selectors не выбирать `.First()` без ambiguous warning.
* [ ] Для contextual selectors проверять container + target uniqueness.
* [ ] Flow-aware validation должна уметь выполнять safe actions до проверяемого шага.
* [ ] Добавить validation diagnostics в отчёт/GUI.

## Acceptance Criteria

* [ ] Selector validation не говорит “OK”, если action невозможен.
* [ ] Selector validation не говорит “missing” без предупреждения о dynamic flow.
* [ ] Fill selector проверяется как editable target.
* [ ] Click selector проверяется как actionable target.
* [ ] Ambiguous chained selector получает warning.
* [ ] Validation result объясняет limitation.

---

# 6.3. Stabilize Recorder Event Ordering and Navigation Causality

**Priority:** High
**Status:** Planned
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

* [ ] Добавить sequence number в recorder events.
* [ ] Добавить timestamp в recorder events.
* [ ] Сначала drain browser recorder events, потом проверять URL change.
* [ ] Добавить navigation correlation с последним user action.
* [ ] Не записывать click-caused navigation как `открыт` перед click.
* [ ] Для click-caused navigation генерировать:

  * click step;
  * optional `ожидаю адрес`.
* [ ] `открыт` использовать только для explicit navigation/open.
* [ ] Добавить correlation window, например 0–2000ms.
* [ ] Добавить immediate event flush для navigation-causing actions.
* [ ] Защититься от потери old-page event queue при full page navigation.
* [ ] Добавить тесты:

  * full navigation after click;
  * SPA route after click;
  * delayed navigation;
  * redirect without click;
  * old page destroyed before polling.

## Acceptance Criteria

* [ ] Click, вызывающий navigation, записывается до navigation step.
* [ ] Recorder не выводит `открыт URL` перед click, который вызвал этот URL.
* [ ] Full page navigation не теряет preceding click.
* [ ] SPA route change не меняет порядок событий.
* [ ] Initial page open всё ещё записывается как `открыт`.
* [ ] Generated scenario replays in logical user order.

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

* [ ] Run с HTML/Allure/traces.
* [ ] Открыть report.
* [ ] Run повторно.
* [ ] Проверить старый report и latest report.
* [ ] Проверить scenario outline с двумя examples.
* [ ] Проверить canceled run partial report.
* [ ] Проверить Allure output после failed write simulation.

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
* [ ] Live browser reuse не ломается.

## Manual QA

* [ ] Run/cancel 20 раз.
* [ ] Cancel во время navigation/wait/download.
* [ ] Parallel run с failure.
* [ ] Close app during run.
* [ ] Проверить, что browsers закрылись.

---

# 8.1. Stabilize Browser Context Isolation

**Priority:** Medium / High
**Status:** Planned
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

* [ ] Явно задокументировать browser context reuse policy.
* [ ] Проверить context/page lifecycle per scenario.
* [ ] Проверить behavior в parallel workers.
* [ ] Определить, когда context reuse допустим.
* [ ] Если reuse disabled — очищать:

  * cookies;
  * localStorage;
  * sessionStorage;
  * opened tabs;
  * downloads;
  * permissions;
  * route/network state.
* [ ] Если auth state reuse включён — сделать это explicit setting.
* [ ] Artifacts/downloads/traces изолировать per case.
* [ ] Добавить тест: Scenario A не влияет на Scenario B.
* [ ] Добавить тест: parallel scenarios не делят storage accidentally.

## Acceptance Criteria

* [ ] Сценарии изолированы по умолчанию или reuse явно включён.
* [ ] Parallel workers не делят mutable browser state случайно.
* [ ] Auth reuse контролируемый, не implicit.
* [ ] Downloads/artifacts не смешиваются между cases.
* [ ] Scenario A не влияет на Scenario B.

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
6.2. Stabilize Browser Selector Validation
6.3. Stabilize Recorder Event Ordering and Navigation Causality
7. Stabilize Reports and Artifacts
8. Stabilize Playwright Runner Lifecycle
8.1. Stabilize Browser Context Isolation
9. Performance and Memory Optimization
10. Architecture Cleanup
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
