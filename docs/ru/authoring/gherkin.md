# Gherkin (русский и английский DSL)

Сценарии Scenaria пишутся на **русском или английском Gherkin**. Язык задаётся в первой строке файла:

```gherkin
# language: ru
```

или

```gherkin
# language: en
```

Без тега используется **русский** (`ru`). От языка зависят ключевые слова, автодополнение в редакторе и разбор шагов.

## Структура (русский)

```gherkin
@smoke
Функционал: Авторизация

  Контекст:
    Дано я подключаю TestClient "User"

  Сценарий: Успешный вход
    Допустим открыт "https://app.example/login"
    Когда заполняю "user@example.com" в "#email"
    И нажимаю "Войти"
    Тогда вижу "text=Кабинет"

  Структура сценария: Параметризованный
    Допустим открыт "<url>"
    Тогда вижу "<selector>"

    Примеры:
      | url                 | selector |
      | https://example.com | h1       |
```

## Ключевые слова

| Ключевое слово | Роль |
|----------------|------|
| `Функционал:` / `Правило:` | Группировка |
| `Сценарий:` / `Структура сценария:` | Кейс / шаблон с таблицей |
| `Контекст:` | Шаги перед каждым сценарием |
| `Дано`, `Когда`, `И`, `Тогда`, `Но` | Шаги |
| `Примеры:` | Таблица данных |
| `Если вижу`, `Повторяю`, `Пока`, `Для каждого` | Управление потоком |

## Структура (английский)

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

  Scenario Outline: Parameterized
    Given I open "<url>"
    Then I see "<selector>"

    Examples:
      | url                 | selector |
      | https://example.com | h1       |
```

| Keyword | Role |
|---------|------|
| `Feature:` | Feature title |
| `Background:` | Steps before each scenario |
| `Scenario:` / `Scenario Outline:` | Test case / outline |
| `Given`, `When`, `Then`, `And`, `But` | Steps |
| `Examples:` | Data table |
| `If I see`, `Repeat`, `While`, `For each` | Control flow |

## Tags

`@smoke` и др. — фильтр запуска:

```bash
scenaria run ./features --tag smoke
```

## TestClient

Сессия в `.scenaria/test_clients/<имя>.json`:

```gherkin
Контекст:
  Дано я подключаю TestClient "имя"
```

Без контекста — чистый браузер на каждый прогон.

## Переменные

- `--var KEY=value` или диалог запуска
- `{{имя}}` в шагах — из таблиц, `запоминаю`, `{{env:VAR}}`

## Каталог шагов

**F1** в IDE. Формулировки шагов соответствуют `# language:` файла (русский или английский).

## См. также

- [Селекторы](selectors.md)
- [Примеры](examples.md)
- [Редактор](../gui/editor.md)
