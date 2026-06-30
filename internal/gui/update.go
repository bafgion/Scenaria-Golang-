package gui

import (
	"fmt"
	"os"
	"path/filepath"

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
}

func updateInfoFrom(info *update.Info) UpdateInfoDTO {
	if info == nil {
		return UpdateInfoDTO{CurrentVersion: version.Version}
	}
	return UpdateInfoDTO{
		CurrentVersion:  info.CurrentVersion,
		LatestVersion:   info.LatestVersion,
		UpdateAvailable: info.UpdateAvailable,
		HTMLURL:         info.HTMLURL,
		DownloadURL:     info.DownloadURL,
		DownloadName:    info.DownloadName,
		Message:         info.Message,
	}
}

func (s *Service) CheckUpdateInfo() (UpdateInfoDTO, error) {
	info, err := update.Check(version.Version)
	if err != nil {
		return UpdateInfoDTO{CurrentVersion: version.Version}, err
	}
	return updateInfoFrom(info), nil
}

func (s *Service) DownloadUpdate() (string, error) {
	info, err := update.Check(version.Version)
	if err != nil {
		return "", err
	}
	if !info.UpdateAvailable || info.DownloadURL == "" {
		return "", fmt.Errorf("обновление недоступно для скачивания")
	}
	name := info.DownloadName
	if name == "" {
		name = "Scenaria-Update.zip"
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

func (s *Service) OpenExternalURL(url string) error {
	return paths.OpenExternalURL(url)
}
