# Миграция с Python Scenaria

Для команд, переходящих с **Python Scenaria (v0.12)** на **Scenaria Go**.

## Что совместимо

| Ресурс | Совместимость |
|--------|----------------|
| `.feature` (RU Gherkin) | Полная |
| `.scenaria/project.json` | Полная |
| `.scenaria/test_clients/*.json` | Полная |
| DSL шагов (~40) | Полная |
| CLI | `run`, `validate`, `export`, `record`, `va`, `plugins`, `init` |
| Теги, таблицы примеров | Полная |
| JUnit / HTML | Полная |

## Что изменилось

| Область | Python | Go |
|---------|--------|-----|
| IDE | Qt | **Wails + Svelte** |
| Среда | Python + Playwright | Go + playwright-go |
| Установка | pip | Установщик / portable / `go install` |

## Соответствие CLI

```bash
# Python                          # Go
scenaria run -t smoke             scenaria run ./features --tag smoke
scenaria run -e VAR=1             scenaria run ./features --var VAR=1
```

Дополнительно в Go: `--workers`, `--test-client`, `import-json`, `va run --rerun-failed`.

## Экспорт в Python (опционально)

```bash
scenaria export ./login.feature --format python --output test_login.py
```

## См. также

- [CLI](../cli/reference.md)
- [Первый проект](../getting-started/first-project.md)

Историческая матрица: `docs/archive/FUNCTIONAL_PARITY_MATRIX.md`
