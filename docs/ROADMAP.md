# Scenaria Go — Roadmap (post-parity)

Статус: **Wails GUI beta** v0.13.0 на ветке `feat/wails-gui`.

## Приоритеты

| # | Направление | Приоритет | Статус |
|---|-------------|-----------|--------|
| P0 | **Wails 2 + Svelte + Monaco** | Очень высокий | ~85% фазы 1 |
| P0 | **Стабилизация Recorder** | Высокий | ~80% фазы 2 |
| P1 | **Allure Report** | Высокий | screenshots done |
| P1 | **slog + %w** | Средний | JSON в CI |
| P1 | **Тесты 60–70%** | Средний | soft gate 35% |
| P2 | **Документация** | Средний | RECORDER + ALLURE |

---

## Фаза 1 — Wails 2 GUI (`feat/wails-gui`)

### 1.4 Полировка
- [x] Splash / about
- [x] `wails build` в portable-скрипт
- [ ] Снять Fyne из release (v0.15)

---

## Фаза 2 — Стабилизация Recorder

### 2.1 Pause / Resume
- [x] Go + Wails + JS flag
- [x] CLI `--idle` + `docs/RECORDER.md`

### 2.2 Edge-cases
- [x] elementFromPoint, debounce, MutationObserver, shadow/iframe

### 2.3 Тесты
- [x] `live_events_test.go`
- [x] `live_integration_test.go` (`-tags=integration`)
- [x] Golden smoke `recorder_script.js`

---

## Фаза 3 — Allure Report

### 3.1–3.3
- [x] Writer + CLI `--allure` + GUI
- [x] Screenshot on failure (Playwright → attachment PNG)
- [ ] trace/video opt-in

---

## Фаза 4 — Observability

- [x] `SCENARIA_LOG=json` для slog
- [ ] `%w` по пакетам

---

## Фаза 5 — Тесты (цель 60–70%)

- [x] `scripts/coverage.ps1`
- [x] CI soft gate 35%
- [x] `go test -tags=integration` (recorder)

---

## Фаза 6 — Документация

- [x] `docs/ALLURE.md`, `docs/RECORDER.md`
- [x] `FUNCTIONAL_PARITY_MATRIX.md` — колонка Wails

---

## Версии

| Версия | Содержание |
|--------|------------|
| **0.13.0** | Wails beta + Allure + recorder docs (текущая) |
| **0.14.0** | trace/video Allure, coverage 60% |
| **0.15.0** | Fyne removed from release |
