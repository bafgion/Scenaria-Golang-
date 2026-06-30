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

Тег `v*` запускает отдельный workflow `.github/workflows/release.yml` (без дублирования тестов из `ci.yml`).

## Автообновление

IDE определяет тип установки по реестру Inno Setup:

| Тип | Артефакт | Действие |
|-----|----------|----------|
| **Setup** (`Program Files\Scenaria`) | `Scenaria-Setup.exe` | тихий установщик `/VERYSILENT` |
| **Portable** | `Scenaria-Portable.zip` | robocopy + перезапуск `scenaria-gui.exe` |

Кнопка **«Установить обновление»** в диалоге обновлений скачивает файл, проверяет SHA256 из `latest.json` и закрывает приложение для применения.
