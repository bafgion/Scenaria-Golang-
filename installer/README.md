# Windows installer (Inno Setup)

Установщик `Scenaria-Setup.exe` ставит CLI, IDE, Chromium и примеры в `Program Files\Scenaria`.

## Сборка

```powershell
# portable ZIP + установщик
powershell -ExecutionPolicy Bypass -File scripts/build-release.ps1

# только установщик (dist\Scenaria уже собран)
powershell -ExecutionPolicy Bypass -File scripts/build-installer.ps1

# только portable
powershell -ExecutionPolicy Bypass -File scripts/build-portable.ps1
```

Требуется [Inno Setup 6](https://jrsoftware.org/isinfo.php) (`ISCC.exe` в PATH или стандартный путь).

## PATH

При установке по умолчанию включена задача **«Добавить scenaria CLI в системный PATH»** — в терминале доступна команда `scenaria` (файл `scenaria.exe` в каталоге установки).

После установки откройте новое окно терминала или перелогиньтесь, чтобы PATH обновился.

## CI / релиз

Тег `v*` на `master` запускает job `release`: собирает `dist/Scenaria-Portable.zip`, `dist/Scenaria-Setup.exe` и публикует их в GitHub Releases.
