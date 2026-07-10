# Установка

Scenaria работает на **Windows 10/11 (x64)**. На Linux/macOS доступен только CLI (GUI в релизных артефактах не собирается).

## Рекомендуется: установщик Windows

Скачайте **Scenaria-Setup.exe** на [странице релизов](https://github.com/bafgion/Scenaria-Golang-/releases).

- Установка в `Program Files\Scenaria`
- В комплекте: `scenaria.exe`, `scenaria-gui.exe`, Chromium, примеры сценариев
- Опция: **добавить `scenaria` в системный PATH** (после установки откройте новый терминал)

Запуск IDE — из меню «Пуск» или `scenaria-gui.exe`.

## Portable ZIP

Скачайте **Scenaria-Portable.zip** и распакуйте в любую папку (например `D:\Tools\Scenaria`).

```
Scenaria/
  scenaria.exe
  scenaria-gui.exe
  Start-GUI.bat
  browsers/
  examples/
  README-PORTABLE.txt
```

Для portable обновление через robocopy; для установщика — тихий Inno Setup.

## Глобальный CLI (разработчикам)

При сборке из исходников:

```bash
go install ./cmd/scenaria
scenaria version
```

Подробнее: [разработка](../contributing/development.md#глобальный-cli), `docs/archive/CLI_GLOBAL_INSTALL.md`.

## Требования (сборка из исходников)

| Компонент | Назначение |
|-----------|------------|
| Go 1.26.5+ | CLI, бэкенд Wails |
| Node.js 22+ | Сборка frontend (в CI — 22) |
| [Wails CLI](https://wails.io/docs/gettingstarted/installation) | `wails dev` / `wails build` |
| Inno Setup 6 | Установщик Windows (опционально) |

## Автообновление

**Справка → Проверить обновления** (или предложение при старте после открытия проекта). IDE читает `latest.json` с GitHub:

| Тип установки | Пакет |
|---------------|-------|
| Inno Setup | `Scenaria-Setup.exe` (`/VERYSILENT`) |
| Portable | `Scenaria-Portable.zip` |

SHA256 сверяется с манифестом релиза.

## Дальше

- [Первый проект](first-project.md)
- [Обучающий тур](onboarding.md)
