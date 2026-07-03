# Windows installer (Inno Setup)

Установщик `Scenaria-Setup.exe` — CLI, IDE, Chromium, примеры в `Program Files\Scenaria`.

**Документация:** [docs/ru/getting-started/install.md](../docs/ru/getting-started/install.md) · [English](../docs/en/getting-started/install.md)

## Сборка

```powershell
./scripts/build-release.ps1
./scripts/build-installer.ps1
./scripts/build-portable.ps1
```

Требуется [Inno Setup 6](https://jrsoftware.org/isinfo.php).

## PATH

Задача установки **«Добавить scenaria CLI в системный PATH»** — команда `scenaria` в новом терминале.

## Релиз и автообновление

Тег `v*` → GitHub Actions release workflow.

| Тип | Обновление |
|-----|------------|
| Setup | `Scenaria-Setup.exe` `/VERYSILENT` |
| Portable | `Scenaria-Portable.zip` + robocopy |
