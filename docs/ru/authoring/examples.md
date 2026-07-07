# Примеры сценариев

Папка `examples/` поставляется с Scenaria; используется [example.com](https://example.com) — без логина.

## Открыть в IDE

1. **«Старт»** → **«Открыть примеры сценариев»**
2. Выберите файл в каталоге
3. **Запустить тест** (Ctrl+Enter) или сначала **Dry-run**

Или **Проект → Открыть проект…** → папка `examples`.

## Файлы

| Файл | Что показывает |
|------|----------------|
| `01-pervaya-proverka.feature` | Страница, заголовок, закрыть браузер (`@smoke`) |
| `02-perehod-po-ssylke.feature` | Ссылка, пауза, «назад» |
| `03-uslovnyy-shag.feature` | Блок **«Если вижу …»** |
| `04-tablica-primerov.feature` | Структура сценария + таблица `Примеры:` (2 прогона) |
| `05-testclient-kontekst.feature` | `Контекст:` + TestClient `DemoUser` |
| `06-proverka-dostupnosti.feature` | `проверяю что доступно` / `недоступно` для кнопок (`@smoke`) |

## TestClient

`05-testclient-kontekst.feature` подключает `examples/.scenaria/test_clients/DemoUser.json`. Для реального проекта:

1. **TestClient…** — войти на сайт, сохранить профиль
2. Блок `Контекст:` с `Дано я подключаю TestClient "имя"`

## Таблица примеров

`04-tablica-primerov.feature` — один шаблон, несколько строк таблицы; в каталоге бейдж «2 прим.». Плейсхолдеры `<url>` из колонок.

## Советы

- Замените URL на свой стенд, когда будете готовы
- Все шаги: **Ctrl+Space** или **F1**
- Запись: **Ctrl+B** → **Начать запись**

## См. также

- [Первый проект](../getting-started/first-project.md)
- [Gherkin](gherkin.md)
