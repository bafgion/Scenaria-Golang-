# Live Recorder (Scenaria Go)

Запись шагов из браузера в `.feature` через Playwright.

## CLI

```bash
scenaria record --live --url https://example.com --output recorded.feature --idle 30
```

| Флаг | Описание |
|------|----------|
| `--live` | Открыть браузер и записывать действия |
| `--url` | Стартовый URL |
| `--output` | Путь к `.feature` |
| `--idle N` | Завершить запись после N секунд без событий (по умолчанию 30) |
| `--headless` | Запись без UI браузера |

## Wails GUI

**Запись…** → URL, файл, idle → **Начать**. Во время записи: **Pause** / **Resume** / **Стоп**.

Панель записи: фильтр, только ссылки, наведение, **headless** — применяются к открытому браузеру (headless перезапускает окно с сохранением URL).

## Что записывается

- Переходы (`goto`) при смене URL
- Клики (с `elementFromPoint`, debounce 400 ms)
- Ввод в поля (`input`) и `<select>` → `выбираю`
- Checkbox / radio → `отмечаю` / `снимаю отметку с`
- Загрузка файлов → `загружаю файл "…" в "…"`
- Клавиши → `нажимаю клавишу "…"` / `… в "поле"`
- Скролл (элемент в центре viewport) → `скроллю к "…"`
- Drag-and-drop → `перетаскиваю "…" в "…"`
- Наведение (опция) → `навожу "…"`
- Подпись на canvas (`draw-signature`)

Селекторы строятся в `internal/selector/recorder_script.js` (testid, id, label, has-text, canvas).

Подробнее о лучших практиках: `docs/SELECTORS.md`.

## Динамический DOM

- **MutationObserver** — новые узлы, shadow DOM
- **iframe** — same-origin: слушатели на `contentDocument`
- Cross-origin iframe не поддерживается (ограничение браузера)

## Нормализация

`internal/recorder/normalize.go` объединяет подряд идущие fill/select, убирает дубли goto/click/scroll.

## Безопасность (trust)

Recorder и picker инжектируют JS в открытые страницы через Playwright `ExposeBinding` (`pickSelectorDone` / `pickSelectorCancel`). Это доверенный режим для **ваших** тестовых стендов:

- Не записывайте и не используйте picker на недоверенных / публичных сайтах с чувствительными данными.
- Bindings привязаны к **browser context** сессии и сбрасываются при закрытии контекста.
- Picker принимает callback только с **активной страницы** и того же origin, что при старте выбора (`internal/recorder/picker.go`).
- Записанные пути (`Output`, `AppendTo`) ограничены корнем проекта — см. `paths.ConfineToProjectRoot`.

При работе с внешними origin рассматривайте отдельный профиль браузера и изолированный проект.

## Тесты

```bash
go test ./internal/recorder/...
go test -tags=integration ./internal/recorder/...   # нужен Playwright + сеть
```
