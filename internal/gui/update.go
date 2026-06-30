package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/bafgion/scenaria-golang/internal/brand"
	"github.com/bafgion/scenaria-golang/internal/paths"
	"github.com/bafgion/scenaria-golang/internal/update"
	"github.com/bafgion/scenaria-golang/internal/version"
)

type UpdateInfoDTO struct {
	CurrentVersion  string `json:"currentVersion"`
	LatestVersion   string `json:"latestVersion"`
	UpdateAvailable bool   `json:"updateAvailable"`
	HTMLURL         string `json:"htmlUrl"`
	DownloadURL     string `json:"downloadUrl"`
	DownloadName    string `json:"downloadName"`
	Message         string `json:"message"`
	InstallMode     string `json:"installMode"`
	ApplyKind       string `json:"applyKind"`
	CanAutoApply    bool   `json:"canAutoApply"`
}

func updateInfoFrom(info *update.Info) UpdateInfoDTO {
	if info == nil {
		return UpdateInfoDTO{CurrentVersion: version.Version, CanAutoApply: runtime.GOOS == "windows"}
	}
	return UpdateInfoDTO{
		CurrentVersion:  info.CurrentVersion,
		LatestVersion:   info.LatestVersion,
		UpdateAvailable: info.UpdateAvailable,
		HTMLURL:         info.HTMLURL,
		DownloadURL:     info.DownloadURL,
		DownloadName:    info.DownloadName,
		Message:         info.Message,
		InstallMode:     info.InstallMode,
		ApplyKind:       info.ApplyKind,
		CanAutoApply:    runtime.GOOS == "windows",
	}
}

func (s *Service) CheckUpdateInfo() (UpdateInfoDTO, error) {
	installDir, _ := paths.AppInstallDir()
	info, err := update.CheckInstallDir(version.Version, installDir)
	if err != nil {
		return UpdateInfoDTO{CurrentVersion: version.Version, CanAutoApply: runtime.GOOS == "windows"}, err
	}
	return updateInfoFrom(info), nil
}

func (s *Service) DownloadUpdate() (string, error) {
	installDir, _ := paths.AppInstallDir()
	info, err := update.CheckInstallDir(version.Version, installDir)
	if err != nil {
		return "", err
	}
	if !info.UpdateAvailable || info.DownloadURL == "" {
		return "", fmt.Errorf("обновление недоступно для скачивания")
	}
	name := info.DownloadName
	if name == "" {
		name = brand.UpdateZip
	}
	destDir, err := os.UserHomeDir()
	if err != nil {
		destDir = os.TempDir()
	}
	destDir = filepath.Join(destDir, "Downloads")
	dest := filepath.Join(destDir, name)
	if err := update.DownloadFile(info.DownloadURL, dest); err != nil {
		return "", err
	}
	return dest, nil
}

func (s *Service) ApplyUpdate() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("автообновление поддерживается только в Windows")
	}
	installDir, err := paths.AppInstallDir()
	if err != nil {
		return err
	}
	_, err = update.Apply(version.Version, installDir, os.Getpid())
	return err
}

func (s *Service) OpenExternalURL(url string) error {
	return paths.OpenExternalURL(url)
}
