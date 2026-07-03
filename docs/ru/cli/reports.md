# Отчёты (CLI)

## JUnit XML

```bash
scenaria run ./features --junit junit.xml
```

Для CI (Azure DevOps, Jenkins, GitLab).

## HTML

```bash
scenaria run ./features --html report.html
```

Сводка по сценариям в одном файле.

## Allure 2

```bash
scenaria run ./features --allure ./allure-results
```

Каждый сценарий — файл `*-result.json` в каталоге.

При падении Playwright (не dry-run) возможны вложения:

| Вложение | Когда |
|----------|-------|
| screenshot (PNG) | Ошибка в браузере |
| trace (ZIP) | С `--trace <dir>` |
| video (WebM) | С `--video <dir>` |

Просмотр:

```bash
allure serve ./allure-results
```

Нужен [Allure CLI](https://docs.qameta.io/allure/#_installing_a_commandline).

### Статусы

| Scenaria | Allure |
|----------|--------|
| passed | passed |
| failed | failed |
| skipped | skipped |
| broken | broken |

## Trace и video

```bash
scenaria run ./features --trace ./traces --video ./videos
```

При падении артефакты сохраняются; в IDE — trace viewer при установленном Playwright CLI.

## Summary JSON

```bash
scenaria run ./features --summary-json summary.json
```

JSON для автоматизации.

## IDE

Те же опции в диалоге **Запуск**. См. [Результаты и отчёты](../gui/results-and-reports.md).
