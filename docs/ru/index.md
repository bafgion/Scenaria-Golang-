# Scenaria — руководство пользователя

**Scenaria Go** — основной продукт: десктопная IDE (Wails) и CLI для сценариев на русском Gherkin с движком Playwright.

## Что входит в состав

- **IDE** (`scenaria-gui.exe`) — редактор Monaco, каталог проекта, живая запись, панели запуска и отчётов, Allure
- **CLI** (`scenaria.exe`) — проверка, запуск, запись, экспорт, плагины, Vanessa Automation
- **Portable или установщик** — Chromium в комплекте, Node.js конечному пользователю не нужен

## Быстрый старт

1. [Установите Scenaria](getting-started/install.md) (установщик или portable ZIP)
2. Запустите **Scenaria** → вкладка **«Старт»** → **«Открыть примеры сценариев»**
3. Выберите `01-pervaya-proverka.feature` → **Запустить тест** (Ctrl+Enter)
4. Пройдите [интерактивный тур](getting-started/onboarding.md) (**Справка → Обучение…**)

## Руководства

### Начало работы

- [Установка](getting-started/install.md)
- [Первый проект](getting-started/first-project.md)
- [Обучающий тур](getting-started/onboarding.md)

### IDE

- [Обзор интерфейса](gui/overview.md)
- [Редактор](gui/editor.md)
- [Запись сценариев](gui/recording.md)
- [Запуск тестов](gui/running-tests.md)
- [Результаты и отчёты](gui/results-and-reports.md)
- [Настройки](gui/settings.md)

### Написание сценариев

- [Gherkin (RU)](authoring/gherkin.md)
- [Селекторы](authoring/selectors.md)
- [Примеры](authoring/examples.md)

### CLI

- [Справочник команд](cli/reference.md)
- [Отчёты](cli/reports.md)

### Прочее

- [Миграция с Python Scenaria](migration/from-python.md)
- [Разработка и тесты](contributing/development.md)

## Справка в программе

**F1** или **Справка → Справка…** — каталог шагов. **Shift+F1** — горячие клавиши.

English version: [docs/en/index.md](../en/index.md)
