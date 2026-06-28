# Scenaria Go — Roadmap (post-parity)

Статус: **master** = v0.14; **feat/v0.15-polish** — GUI trace/video, Fyne out of release.

## Приоритеты

| # | Направление | Статус |
|---|-------------|--------|
| P0 | Wails IDE | ~90% |
| P0 | Recorder | ~80% |
| P1 | Allure | done (screenshot + trace + video) |
| P1 | Coverage | soft 40%, цель 60% |
| P2 | Fyne | убран из portable (v0.15) |

---

## Фаза 1 — Wails GUI

- [x] Monaco IDE, run/record/Vanessa/OTP/Allure
- [x] Splash, portable `wails build`
- [x] GUI: trace/video чекбоксы
- [x] Fyne снят с release

---

## Фаза 3 — Allure

- [x] Writer, CLI, GUI
- [x] Screenshot + trace + video attachments

---

## Фаза 5 — Тесты

- [x] `scripts/coverage.ps1` (`-coverpkg=./...`)
- [x] CI soft gate 40% (цель 60%)

---

## Версии

| Версия | Содержание |
|--------|------------|
| **0.14.0** | trace/video Allure (**master**) |
| **0.15.0** | GUI trace/video + Fyne out of release (в работе) |
