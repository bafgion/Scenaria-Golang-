# Gherkin (Russian DSL)

Scenaria scenarios use **Russian Gherkin** keywords. The IDE and runner expect this syntax.

## Structure

```gherkin
@smoke @api
Функционал: Авторизация

  Контекст:
    Дано я подключаю TestClient "User"

  Сценарий: Успешный вход
    Допустим открыт "https://app.example/login"
    Когда заполняю "user@example.com" в "#email"
    И заполняю "secret" в "#password"
    И нажимаю "Войти"
    Тогда вижу "text=Кабинет"

  Структура сценария: Параметризованный вход
    Допустим открыт "<url>"
    Тогда вижу "<selector>"

    Примеры:
      | url                        | selector |
      | https://example.com        | h1       |
      | https://example.com/about  | h1       |
```

## Keywords

| Russian | Role |
|---------|------|
| `Функционал:` / `Правило:` | Feature grouping |
| `Сценарий:` / `Структура сценария:` | Test case / outline |
| `Контекст:` | Steps before each scenario (TestClient, hooks) |
| `Дано`, `Когда`, `И`, `Тогда`, `Но` | Steps |
| `Примеры:` | Data table for outlines |
| `Если вижу`, `Повторяю`, `Пока`, `Для каждого` | Control flow blocks |

## Tags

Tags like `@smoke` filter runs:

```bash
scenaria run ./features --tag smoke
```

Feature-level and scenario-level tags are supported.

## TestClient

Save browser cookies/storage under `.scenaria/test_clients/<name>.json`, then:

```gherkin
Контекст:
  Дано я подключаю TestClient "name"
```

Without context, each run starts a clean browser.

## Variables

- Scenario variables: `--var KEY=value` or run dialog
- `{{name}}` in steps — from tables, `запоминаю`, or env `{{env:VAR}}`

## Step catalog

Full step list: **F1** in IDE or export from `internal/stepcatalog/`. Steps use Russian phrasing matching Playwright actions.

## Related

- [Selectors](selectors.md)
- [Examples](examples.md)
- [Editor](../gui/editor.md)
