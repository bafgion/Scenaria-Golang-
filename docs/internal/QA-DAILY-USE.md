# QA Daily-use — ручной регресс (v0.27.0)

Чеклист для проверки ежедневного использования IDE после фазы 15. Выполнять на свежей установке или после `settings.json` → сброс по умолчанию.

**Окружение:** Windows, Wails desktop build или `wails dev`, Playwright browsers установлены.

---

## 1. Onboarding

| # | Шаг | Ожидание |
|---|-----|----------|
| 1.1 | Первый запуск → экран «Старт», чеклист из 3 шагов | Чеклист виден, «Быстрый старт» без проекта не начинает запись |
| 1.2 | «Быстрый старт» без проекта | Сообщение «Сначала откройте проект», открывается диалог проекта |
| 1.3 | «Новый проект…» → папка + Init + шаблон | Проект открыт, `.scenaria/` создан, `smoke.feature` в редакторе |
| 1.4 | «Открыть примеры» | Каталог слева, подсказка выбрать сценарий |
| 1.5 | «Скрыть чеклист» → перезапуск | Чеклист не показывается (persist `checklistDismissed`) |
| 1.6 | Модал обновлений при первом визите | Не перекрывает splash / welcome до открытия проекта или скрытия чеклиста |

---

## 2. Monaco / редактор

| # | Шаг | Ожидание |
|---|-----|----------|
| 2.1 | 5 вкладок, правка в каждой, переключение | `*` только на реально изменённых вкладках |
| 2.2 | Сохранить как… после правки | Сохраняется актуальный текст из редактора |
| 2.3 | Settings → Тема «Как в системе» | Тема следует OS light/dark |
| 2.4 | Файл ≥2000 строк | Баннер «упрощённый режим» в status bar |
| 2.5 | Ctrl+F vs Ctrl+H | Find и Find+Replace раздельно |

---

## 3. Live Recording

| # | Шаг | Ожидание |
|---|-----|----------|
| 3.1 | Запись с активного таба | Status bar: «Запись → file.feature» |
| 3.2 | Переключение вкладки во время записи | Confirm; опция «Больше не спрашивать» |
| 3.3 | Idle timeout (если настроен) | Toast / запись в журнале |
| 3.4 | Stop во время записи | Подпись «Стоп запись», не «Стоп тест» |
| 3.5 | Headless toggle в recording bar | Confirm перед relaunch |

---

## 4. Запуск тестов

| # | Шаг | Ожидание |
|---|-----|----------|
| 4.1 | Suite из ≥3 сценариев | Playing bar: имя + N/M, progress bar |
| 4.2 | Журнал во время run | Стрим stdout в реальном времени |
| 4.3 | Первый Ctrl+Enter | RunDialog (до `runDialogConfirmed`) |
| 4.4 | Cancel в playing bar | «Останавливаем…», прогон прерывается |
| 4.5 | Dry-run | Редактор доступен, подпись в status bar |
| 4.6 | Continue on fail в RunDialog | Suite не останавливается на первом fail |

---

## 5. Results

| # | Шаг | Ожидание |
|---|-----|----------|
| 5.1 | Упавший сценарий | Кнопка «Шаг N» → переход к строке |
| 5.2 | Trace ZIP после run | «Trace viewer» открывает Playwright trace |
| 5.3 | Flaky badge | «Запустить 3×» запускает сценарий три раза подряд |
| 5.4 | Allure без PATH | «Allure не найден» + ссылка на установку |
| 5.5 | Двойной клик по строке | Открывается feature |

---

## 6. Settings

| # | Шаг | Ожидание |
|---|-----|----------|
| 6.1 | «Сбросить по умолчанию» | Browser, navWaitUntil, editor восстановлены |
| 6.2 | navWaitUntil в UI | Значение уходит в run/record |
| 6.3 | Workers 16 + slowMo 5000 | Предупреждение в диалоге |
| 6.4 | Apply без закрытия | Настройки применены, диалог открыт |
| 6.5 | Ctrl+S при открытых Settings | Не сохраняет сценарий случайно |

---

## Критерии прохождения

- Все пункты **без блокеров** (P0/P1).
- Не более 2 мелких UX-замечаний (P2/P3) — зафиксировать в ROADMAP backlog.

**Версия документа:** v0.27.0 / Фаза 15–16

---

## Автоматизация (E2E)

**Файл:** `frontend/e2e/specs/qa-daily-use.spec.ts` — 32 теста, по одному на пункт чеклиста.

**Сквозные сценарии (user journeys):** `frontend/e2e/specs/user-journeys.spec.ts` — 11 тестов, цепочки «как в реальной работе».

Запуск:

```bash
cd frontend && npm run test:e2e              # все E2E (~80 тестов)
cd frontend && npm run test:e2e:journeys     # только user journeys (быстрый локальный прогон)
cd frontend && npm run test:e2e:qa             # только чеклист QA-DAILY-USE
cd frontend && npm run test:e2e -- e2e/specs/qa-daily-use.spec.ts
```

Окружение: mock Wails (`wails-mock.js`), без реального Playwright/браузера записи. Режимы `?e2e=…` эмулируют прогон, артефакты, idle-запись и т.д.

| Пункт | E2E | Примечание |
|-------|-----|------------|
| 1.1–1.6 | ✅ | UI onboarding, wizard, update defer; **интерактивный тур** — `onboarding-tour.spec.ts` |
| 2.1–2.5 | ✅ / частично | 2.3 — выбор «Как в системе», не проверка OS theme |
| 3.1–3.5 | ✅ | mock `post-record`, `record-idle` |
| 4.1–4.6 | ✅ | mock `run-progress`, `run-stream`, `run-cancel` |
| 5.1–5.5 | ✅ | mock `flaky-run`, `trace-artifacts`, `allure-missing` |
| 6.1–6.5 | ✅ | Settings dialog |

### User journeys (сквозные сценарии)

| Сценарий | Что эмулирует |
|----------|----------------|
| Знакомство с примерами | Примеры → smoke → проверка → dry-run |
| Новый проект | Мастер → правка → Ctrl+S → проверка |
| Live-запись | HTTP Auth → запись → стоп; confirm при смене вкладки |
| Прогон suite | Playing bar N/M, журнал, отмена прогона |
| Разбор падения | Шаг ошибки → flaky 3× → Trace viewer |
| Браузер и фокус | Открыть → «Показать браузер» → закрыть |
| Настройки | Apply без закрытия, Ctrl+S не трогает сценарий |
| Плагины | Список → «Запуск…» → отмена |
| Модалки | Escape закрывает только верхний диалог |

### Desktop smoke (Windows + WebView2)

**Скрипт:** `scripts/desktop-smoke.ps1` — запускает `scenaria-gui.exe` с CDP `:9333`, проверяет живость процесса и гоняет Playwright против реального WebView2.

**Спеки:** `frontend/e2e/specs/desktop-smoke.spec.ts` (22 теста) и `frontend/e2e/specs/desktop-tour-onboarding.spec.ts` (тур до шагов 5–6). Без wails-mock.

```powershell
# Сборка + полный smoke (процесс + UI)
./scripts/desktop-smoke.ps1

# Только проверка процесса (без Playwright UI)
./scripts/desktop-smoke.ps1 -SkipUI

# UI-тесты, если приложение уже запущено с CDP (SCENARIA_DESKTOP_SMOKE=1)
$env:SCENARIA_DESKTOP_SMOKE = "1"
$env:SCENARIA_DESKTOP_SMOKE_PORT = "9333"
.\build\bin\scenaria-gui.exe
cd frontend && npm run test:e2e:desktop

# Go-обёртка
go test -tags=desktop ./internal/wailsapp/...
```

| Desktop smoke | E2E mock |
|---------------|----------|
| Реальный Wails + WebView2 | vite preview + wails-mock.js |
| Примеры с диска `examples/` | `?e2e=examples` |
| Диалоги, палитры, настройки, экспорт, плагины | Полный чеклист QA + mock-прогоны |
| Проверка сценария (реальный backend) | Dry-run / run-progress / flaky-run |

**Перенесено из mock E2E (22 теста):** загрузка, примеры/каталог, проверка, настройки (6.1/6.3/6.4), палитра команд, справка F1, горячие клавиши, журнал Ctrl+`, новый сценарий, сниппеты, RunDialog, экспорт-превью, несохранённая вкладка, плагины, «О программе», запись + HTTP Auth.

**Обучалка (тур):** `frontend/e2e/specs/onboarding-tour.spec.ts` — первый запуск, пропуск, шаг 5 (подсветка «Проверить…», z-index). Desktop: `desktop-tour-onboarding.spec.ts` через CDP. Шкала слоёв: `--z-onboarding-*` в `style.css`, тест `layers.test.ts`. Повтор: «Справка → Обучение…».

**Пока только mock** (нужны `?e2e=…`, эмуляция прогона, нативные диалоги или длинный run): live-запись, post-record diff, import/export на диск, dry-run completion, flaky 3×, trace viewer, new-project wizard, missing-browser, session restore после reload, Monaco Ctrl+Shift+O, Ctrl+S на реальный проект.

### Что остаётся ручным

| Область | Почему не E2E |
|---------|----------------|
| Реальный Playwright run + trace ZIP | Нужен Wails desktop + установленный браузер |
| Запись в живой сайт | Recorder + CDP вне static preview |
| System theme визуально | `prefers-color-scheme` в mock; цвета Monaco — smoke на ОС |
| Suite ≥10 сценариев end-to-end | Долгий интеграционный прогон, не UI mock |
| `navWaitUntil` в реальном run | Проверяется передачей в CLI в integration/desktop |

Дополнительно: `app-ui.spec.ts` — регрессия смежных сценариев (Monaco, import, post-record diff).
