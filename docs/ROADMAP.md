# Scenaria Go — Roadmap (post-parity)

Статус: **Wails GUI beta** на ветке `feat/wails-gui`; master = CLI parity v0.12.

## Приоритеты

| # | Направление | Приоритет | Статус |
|---|-------------|-----------|--------|
| P0 | **Wails 2 + Svelte + Monaco** | Очень высокий | ~70% фазы 1 |
| P0 | **Стабилизация Recorder** | Высокий | в работе |
| P1 | **Allure Report** | Высокий | stub |
| P1 | **slog + %w** | Средний | logx stub |
| P1 | **Тесты 60–70%** | Средний | не начато |
| P2 | **Документация** | Средний | частично |

---

## Фаза 1 — Wails 2 GUI (`feat/wails-gui`)

### 1.1 Инфраструктура
- [x] `wails.json`, `frontend/` (Svelte + TS + Vite)
- [x] `internal/gui` — GUI-agnostic API
- [x] `internal/wailsapp` — биндинги + OTP events
- [x] CI: `npm install && npm run build` в `frontend/`
- [x] Fyne **deprecated** в README

### 1.2 Layout IDE
- [x] Трёхпанельный layout (дерево → Monaco → лог/запуск)
- [x] Monaco + язык `scenaria-feature`
- [x] Native folder picker (`PickProjectFolder`)
- [x] Вкладки файлов + Ctrl+S
- [x] Маркеры ошибок шагов в Monaco (`ValidateFeature`)

### 1.3 Паритет с Fyne
- [x] Открытие проекта, init, список `.feature`
- [x] Запуск: dry-run / Playwright, теги, vars, TestClient
- [x] Validate syntax + browser
- [x] Live-запись + Pause/Resume/Stop
- [x] Экспорт / импорт JSON (native file dialogs)
- [x] Настройки (browser, headless, workers, loops)
- [x] Vanessa dry/run
- [x] OTP prompt (modal + `SubmitOTPCode`)

### 1.4 Полировка (осталось)
- [ ] Splash / about
- [ ] `wails build` в portable-скрипт
- [ ] Снять Fyne из release по умолчанию

**Критерий готовности:** сценарий IDE без CLI (запись → правка → run → validate).

---

## Фаза 2 — Стабилизация Recorder

### 2.1 Pause / Resume
- [x] `LiveSession` в Go
- [x] Pause/Resume в Wails UI
- [x] Флаг `paused` в `recorder_script.js` (синхрон с Go)
- [ ] CLI `--idle` + документация pause

### 2.2 Edge-cases
- [x] `elementFromPoint` для кликов (модалки / top layer)
- [x] Debounce дублирующих кликов (400 ms)
- [ ] MutationObserver для динамических элементов
- [ ] iframe/shadow DOM (best-effort + docs)

### 2.3 Тесты
- [ ] `recorder/live_integration_test.go`
- [ ] Golden tests для `recorder_script.js`

---

## Фаза 3 — Allure Report

### 3.1 Модель
- [x] `internal/report/allure` — stub + mkdir
- [ ] Writer `allure-results/` + steps
- [ ] CLI: `scenaria run ... --allure <dir>`

### 3.2 Attachments
- [ ] Screenshot on failure
- [ ] trace/video opt-in

### 3.3 GUI
- [ ] Кнопка «Открыть Allure» после run

---

## Фаза 4 — Observability

- [x] `internal/logx` (slog text)
- [ ] JSON handler в CI
- [ ] `%w` по пакетам, меньше `fmt.Printf` в CLI

---

## Фаза 5 — Тесты (цель 60–70%)

- [x] `internal/gui` validate + service smoke
- [ ] `scripts/coverage.ps1` + soft gate в CI
- [ ] `go test -tags=integration`

---

## Фаза 6 — Документация

- [x] `docs/MIGRATION_FROM_PYTHON.md`
- [x] README — Wails primary
- [ ] `FUNCTIONAL_PARITY_MATRIX.md` — колонка Wails (частично)
- [ ] `docs/RECORDER.md`, `docs/ALLURE.md`

---

## Версии

| Версия | Содержание |
|--------|------------|
| **0.13.0** | Wails GUI beta (текущая ветка) |
| **0.14.0** | Allure + recorder stable |
| **0.15.0** | Fyne removed from release |
