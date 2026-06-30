package update

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Apply performs download, verification, and launch of the auto-update for installDir.
func Apply(currentVersion, installDir string, parentPID int) (*Info, error) {
	return ApplyWithProgress(currentVersion, installDir, parentPID, nil)
}

// ApplyWithProgress is Apply with optional progress reporting.
func ApplyWithProgress(currentVersion, installDir string, parentPID int, report Reporter) (*Info, error) {
	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf("автообновление поддерживается только в Windows")
	}
	installDir = strings.TrimSpace(installDir)
	if installDir == "" {
		return nil, fmt.Errorf("каталог установки не определён")
	}
	installDir = filepath.Clean(installDir)
	report.report("check", "Проверка обновления…", 2)
	info, err := CheckInstallDir(currentVersion, installDir)
	if err != nil {
		return nil, err
	}
	if !info.UpdateAvailable || info.DownloadURL == "" {
		return info, fmt.Errorf("обновление недоступно для скачивания")
	}
	tempDir, err := os.MkdirTemp("", "scenaria-download-")
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(info.DownloadName)
	if name == "" {
		if info.ApplyKind == string(ApplyKindPortable) {
			name = "Scenaria-Portable.zip"
		} else {
			name = "Scenaria-Setup.exe"
		}
	}
	dest := filepath.Join(tempDir, name)
	report.report("download", "Скачивание обновления…", 5)
	if err := downloadFile(dest, info.DownloadURL, func(done, total int64) {
		report.report("download", formatDownloadMessage(done, total), downloadPercent(done, total, 5, 70))
	}); err != nil {
		_ = os.RemoveAll(tempDir)
		return nil, err
	}
	report.report("verify", "Проверка контрольной суммы…", 78)
	if err := verifyFileSHA256(dest, info.SHA256); err != nil {
		_ = os.RemoveAll(tempDir)
		return nil, err
	}
	kind := ApplyKind(info.ApplyKind)
	if kind == "" {
		kind = applyKindForAsset(name)
	}
	report.report("prepare", "Подготовка к установке…", 82)
	if err := ApplyDownloaded(dest, kind, installDir, parentPID, report); err != nil {
		_ = os.RemoveAll(tempDir)
		return nil, err
	}
	report.report("restart", "Перезапуск приложения…", 100)
	return info, nil
}
