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

## Динамический DOM

- **MutationObserver** — новые узлы, shadow DOM
- **iframe** — same-origin: слушатели на `contentDocument`
- Cross-origin iframe не поддерживается (ограничение браузера)

## Нормализация

`internal/recorder/normalize.go` объединяет подряд идущие fill/select, убирает дубли goto/click/scroll.

## Тесты

```bash
go test ./internal/recorder/...
go test -tags=integration ./internal/recorder/...   # нужен Playwright + сеть
```
