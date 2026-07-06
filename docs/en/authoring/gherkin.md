# Gherkin (Russian and English DSL)

Scenarios use **Russian or English Gherkin**. Set the dialect in the first line:

```gherkin
# language: ru
```

or

```gherkin
# language: en
```

If omitted, **Russian** (`ru`) is used. The tag affects keywords, editor completions, and step parsing.

## Structure (Russian)

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

## Structure (English)

```gherkin
# language: en
@smoke
Feature: Login

  Background:
    Given I connect TestClient "User"

  Scenario: Successful sign-in
    Given I open "https://app.example/login"
    When I type "user@example.com" into "#email"
    And I click "Sign in"
    Then I see "h1"
```

| English | Role |
|---------|------|
| `Feature:` | Feature title |
| `Background:` | Steps before each scenario |
| `Scenario:` / `Scenario Outline:` | Test case / outline |
| `Given`, `When`, `Then`, `And`, `But` | Steps |
| `Examples:` | Data table |
| `If I see`, `Repeat`, `While`, `For each` | Control flow |

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

Full step list: **F1** in IDE. Step phrasing follows the file `# language:` (Russian or English).

## Related

- [Selectors](selectors.md)
- [Examples](examples.md)
- [Editor](../gui/editor.md)
