# Справочник CLI

```bash
scenaria <команда> [флаги] [пути...]
scenaria version
```

## Команды

| Команда | Описание |
|---------|----------|
| `run` | Выполнить сценарии |
| `validate` | Проверить `.feature` |
| `export` | Экспорт в JSON / feature / TypeScript / Python |
| `import-json` | Импорт JSON в `.feature` |
| `record` | Создать или записать сценарий |
| `init` | Создать `.scenaria/` |
| `update` | Проверить обновления (`--check`) |
| `plugins` | `list`, `install`, `uninstall` |
| `va` | Vanessa Automation (`va run`, …) |
| `version` | Версия |
| `help` | Справка |

## `run`

```bash
scenaria run ./features [флаги]
```

| Флаг | Описание |
|------|----------|
| `--dry-run` | Без браузера |
| `--tag NAME` | Фильтр по тегу (`smoke` → `@smoke`) |
| `--scenario NAME` | Один сценарий по названию |
| `--var KEY=VALUE` | Переменная (можно повторять) |
| `--test-client NAME` | TestClient для всех кейсов |
| `--workers N` | Параллельные workers |
| `--slow-mo MS` | Замедление Playwright |
| `--headed` | Показать окно браузера |
| `--browser NAME` | `chromium`, `firefox`, `webkit` |
| `--base-url URL` | Базовый URL |
| `--install-playwright` | Установить браузер при необходимости |
| `--engine NAME` | `playwright` или `stub` |
| `--junit PATH` | JUnit XML |
| `--html PATH` | HTML-отчёт |
| `--allure DIR` | Каталог Allure |
| `--trace DIR` | Trace при падении |
| `--video DIR` | Video при падении |
| `--summary-json PATH` | JSON-сводка |
| `--start-step N` | С шага (0-based) |
| `--end-step N` | До шага (0-based) |
| `--nav-wait-until` | `load`, `domcontentloaded`, `networkidle`, `commit` |
| `--continue-on-fail` | Не останавливать пакет при первой ошибке |

Примеры:

```bash
scenaria run ./examples --dry-run
scenaria run ./features --tag smoke --headed --install-playwright
scenaria run ./login.feature --trace ./traces --allure ./allure-results
```

## `validate`

```bash
scenaria validate ./features [флаги]
```

| Флаг | Описание |
|------|----------|
| `--json PATH` | JSON с результатами |
| `--browser` | Проверка в браузере (по умолчанию) |
| `--no-browser` | Только синтаксис |
| `--headless` | Браузер без UI |
| `--base-url URL` | Для относительных URL |

## `export` / `import-json`

```bash
scenaria export INPUT.feature --output OUT --format json|feature|ts|python
scenaria import-json doc.json --output doc.feature [--force]
```

## `record`

```bash
scenaria record --live --url https://example.com --output recorded.feature --idle 30
```

## `init` / `update` / `plugins` / `va`

```bash
scenaria init .
scenaria update --check
scenaria plugins list
scenaria va run --project . --dry-run
```

## Установка в PATH

```bash
go install ./cmd/scenaria
```

## См. также

- [Отчёты](reports.md)
- [Установка](../getting-started/install.md)
