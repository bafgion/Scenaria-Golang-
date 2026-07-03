# Первый проект

## Открыть примеры

Самый быстрый способ познакомиться с программой:

1. Запустите **scenaria-gui.exe**
2. Вкладка **«Старт»** → **«Открыть примеры сценариев»**
3. В **каталоге** слева — файлы из `examples/`
4. Откройте файл двойным щелчком

Описание файлов: [Примеры](../authoring/examples.md).

## Создать новый проект

**«Старт»** → **«Новый проект…»** (мастер):

| Шаг | Действие |
|-----|----------|
| Папка | Выберите или создайте каталог |
| Init | При необходимости создать `.scenaria/` (`scenaria init`) |
| Образец | При необходимости добавить стартовый `.feature` |

Или вручную:

```bash
mkdir my-tests && cd my-tests
scenaria init .
```

Создаётся `.scenaria/` с `project.json`, настройками и каталогами артефактов.

## Открыть существующую папку

**Проект → Открыть проект…** — каталог с `.feature` (с `.scenaria/` или без).

Недавние проекты — на вкладке **«Старт»**.

## Структура проекта

```
my-project/
  login.feature
  smoke.feature
  .scenaria/
    project.json
    settings.json
    test_clients/
    allure-results/
    traces/ videos/
```

## Инициализация в IDE

**Проект → Инициализировать Scenaria…** — добавляет `.scenaria/` без перезаписи сценариев.

## CLI

```bash
scenaria init /path/to/project
scenaria validate /path/to/project
scenaria run /path/to/project --dry-run
```

## Дальше

- [Обучающий тур](onboarding.md)
- [Редактор](../gui/editor.md)
- [Запись](../gui/recording.md)
