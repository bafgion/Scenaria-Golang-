# Scenaria Go — Roadmap (post-parity)

Статус master: **runtime/CLI parity с Python v0.12** (PR #9).  
Следующий этап: **Wails 2 GUI**, стабильность recorder, Allure, observability, тесты, документация.

## Приоритеты (согласовано с архитектором)

| # | Направление | Приоритет | Оценка |
|---|-------------|-----------|--------|
| P0 | **Wails 2 + Svelte** вместо Fyne | Очень высокий | 3–4 недели |
| P0 | **Стабилизация Recorder** (edge-cases, pause/resume) | Высокий | 1–2 недели |
| P1 | **Allure Report** (CLI + GUI, attachments) | Высокий | 1–2 недели |
| P1 | **Error handling + slog** | Средний | 3–5 дней |
| P1 | **Тесты 60–70%** + integration | Средний | 2–3 недели |
| P2 | **Документация** (README, MIGRATION_FROM_PYTHON) | Средний | 2–3 дня |

---

## Фаза 1 — Wails 2 GUI (текущая ветка `feat/wails-gui`)

### 1.1 Инфраструктура (неделя 1)
- [x] `wails.json`, `frontend/` (Svelte + TypeScript + Vite)
- [x] `internal/gui` — GUI-agnostic API (без Fyne/Wails)
- [x] `internal/wailsapp` — биндинги Wails → `gui.Service`
- [ ] `wails generate module` в CI (проверка синхронизации bindings)
- [ ] Fyne помечен **deprecated** в README; сборка `-tags desktop` сохранена до паритета UI

### 1.2 Layout IDE (неделя 2)
- [ ] Трёхпанельный layout: дерево feature → редактор → лог/панель
- [ ] Monaco или CodeMirror 6 для `.feature` + подсветка из `featurehighlight`
- [ ] Вкладки файлов, сохранение, горячие клавиши

### 1.3 Паритет с Fyne (неделя 3)
- [ ] Открытие проекта, init, список `.feature`
- [ ] Запуск: dry-run / Playwright, теги, `--var`, TestClient picker
- [ ] Validate (syntax + browser)
- [ ] Запись live, экспорт/импорт JSON
- [ ] Настройки (browser, workers, loops)
- [ ] Vanessa dry/run, OTP prompt (modal)

### 1.4 Полировка (неделя 4)
- [ ] Splash / about / версия
- [ ] Portable build: `scenaria-gui.exe` через `wails build`
- [ ] Удаление Fyne из default release (опционально оставить `desktop` tag)

**Критерий готовности:** пользовательский сценарий из Python/Qt IDE воспроизводится в Wails без CLI.

---

## Фаза 2 — Стабилизация Recorder

### 2.1 Pause / Resume
- [x] `LiveRecorder` с `Pause()` / `Resume()` / `IsPaused()`
- [ ] API в CLI: `scenaria record --live --pause-key` (опционально)
- [ ] Кнопки Pause/Resume в Wails во время записи

### 2.2 Edge-cases
- [ ] Модальные окна: запись только top-layer (`elementFromPoint`, z-index)
- [ ] Debounce быстрых кликов / coalescing duplicate events
- [ ] Динамические элементы: mutation observer для late-mounted controls
- [ ] iframe/shadow DOM (best-effort, документировать ограничения)

### 2.3 Тесты
- [ ] `recorder/live_integration_test.go` (playwright test server fixtures)
- [ ] Golden tests для `recorder_script.js` → selector output

---

## Фаза 3 — Allure Report

### 3.1 Модель
- [ ] `internal/report/allure` — writer `allure-results/`
- [ ] Шаги сценария → Allure steps; статусы passed/failed/skipped
- [ ] CLI: `scenaria run ... --allure <dir>`

### 3.2 Attachments
- [ ] Скриншот on failure (Playwright)
- [ ] trace/video opt-in (`--video`, `--trace`)
- [ ] Лог шага / stdout runner

### 3.3 GUI
- [ ] Кнопка «Открыть Allure» / путь к результатам после run
- [ ] Совместимость с Vanessa `--allure` (уже частично в `internal/vanessa`)

---

## Фаза 4 — Observability

- [ ] `internal/logx` на `log/slog` (JSON в CI, text в dev)
- [ ] Обёртка ошибок: `%w` + контекст пакета (`fmt.Errorf("recorder: %w", err)`)
- [ ] Постепенная замена `fmt.Printf` в CLI на slog с уровнями

---

## Фаза 5 — Тесты (цель 60–70%)

| Пакет | Сейчас (~оценка) | Цель |
|-------|------------------|------|
| `stepdsl`, `gherkin` | высокое | поддерживать golden |
| `player` | среднее | + integration с stub HTML |
| `recorder` | низкое | + live fixtures |
| `gui` / `wailsapp` | — | unit + smoke e2e |
| `report/allure` | — | snapshot tests |

- [ ] `scripts/coverage.ps1` + порог в CI (пока soft, потом 60%)
- [ ] Integration tag: `go test -tags=integration ./...`

---

## Фаза 6 — Документация

- [x] `docs/MIGRATION_FROM_PYTHON.md` (черновик)
- [ ] `README.md` — Wails dev/build, убрать Fyne как primary
- [ ] `FUNCTIONAL_PARITY_MATRIX.md` — колонка Wails GUI
- [ ] `docs/RECORDER.md`, `docs/ALLURE.md`

---

## Вопросы на будущее (не блокируют старт)

1. **Monaco vs CodeMirror** для редактора — Monaco тяжелее, но привычнее; CM6 легче для Wails.
2. **Allure**: только native writer или также вызов `allure generate` CLI?
3. **Видео/trace**: хранить в `.scenaria/runs/<id>/` по умолчанию?

---

## Ветки и релизы

| Ветка | Содержание |
|-------|------------|
| `master` | стабильный CLI + Fyne (deprecated) |
| `feat/wails-gui` | Wails scaffold → IDE parity |
| `feat/recorder-stable` | после merge scaffold |
| `feat/allure` | параллельно после run hooks |

Версия **0.13.0** — Wails GUI beta. **0.14.0** — Allure + recorder stable.
