# Scenaria Go — Roadmap

Статус: **master** v0.15; **Wails IDE** — основной продукт. Python/Qt — снят с поддержки (экспорт в Python сохранён).

## Приоритеты

| # | Направление | Статус |
|---|-------------|--------|
| P0 | Wails IDE | ~96% |
| P0 | Recorder | ~90% |
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

## Версии

| Версия | Содержание |
|--------|------------|
| **0.14.0** | trace/video Allure (**master**) |
| **0.15.0** | GUI trace/video + scenario catalog + recorder polish (**master**) |
