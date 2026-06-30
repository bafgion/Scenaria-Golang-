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
	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf("автообновление поддерживается только в Windows")
	}
	installDir = strings.TrimSpace(installDir)
	if installDir == "" {
		return nil, fmt.Errorf("каталог установки не определён")
	}
	installDir = filepath.Clean(installDir)
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
	if err := DownloadFile(info.DownloadURL, dest); err != nil {
		_ = os.RemoveAll(tempDir)
		return nil, err
	}
	if err := verifyFileSHA256(dest, info.SHA256); err != nil {
		_ = os.RemoveAll(tempDir)
		return nil, err
	}
	kind := ApplyKind(info.ApplyKind)
	if kind == "" {
		kind = applyKindForAsset(name)
	}
	if err := ApplyDownloaded(dest, kind, installDir, parentPID); err != nil {
		_ = os.RemoveAll(tempDir)
		return nil, err
	}
	return info, nil
}
