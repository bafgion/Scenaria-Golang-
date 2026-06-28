package vanessa

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const vanessaGitHubRepo = "Pr-Mex/vanessa-automation"

func DownloadEPF(destination string, downloadURL string) (string, error) {
	if strings.TrimSpace(destination) == "" {
		home, _ := os.UserConfigDir()
		destination = filepath.Join(home, "Scenaria", "vanessa", "vanessa-automation.epf")
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return "", err
	}
	url := strings.TrimSpace(downloadURL)
	if url == "" {
		var err error
		url, err = resolveLatestSingleZipURL()
		if err != nil {
			return "", err
		}
	}
	tmp, err := downloadFile(url)
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp)
	if strings.HasSuffix(strings.ToLower(url), ".zip") || isZip(tmp) {
		if err := extractEPFFromZip(tmp, destination); err != nil {
			return "", err
		}
	} else {
		if err := copyFile(tmp, destination); err != nil {
			return "", err
		}
	}
	info, err := os.Stat(destination)
	if err != nil || info.Size() == 0 {
		return "", fmt.Errorf("downloaded EPF is empty")
	}
	return destination, nil
}

func resolveLatestSingleZipURL() (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/"+vanessaGitHubRepo+"/releases/latest", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "Scenaria-Go")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("GitHub API HTTP %d", resp.StatusCode)
	}
	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}
	for _, asset := range release.Assets {
		name := strings.ToLower(asset.Name)
		if strings.HasPrefix(name, "vanessa-automation-single.") && strings.HasSuffix(name, ".zip") {
			return asset.BrowserDownloadURL, nil
		}
	}
	for _, asset := range release.Assets {
		name := strings.ToLower(asset.Name)
		if strings.Contains(name, "single") && strings.HasSuffix(name, ".zip") {
			return asset.BrowserDownloadURL, nil
		}
	}
	return "", fmt.Errorf("vanessa single zip not found in latest release")
}

func downloadFile(url string) (string, error) {
	client := &http.Client{Timeout: 180 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("download HTTP %d", resp.StatusCode)
	}
	tmp, err := os.CreateTemp("", "vanessa-*.zip")
	if err != nil {
		return "", err
	}
	defer tmp.Close()
	if _, err := io.Copy(tmp, resp.Body); err != nil {
		os.Remove(tmp.Name())
		return "", err
	}
	return tmp.Name(), nil
}

func isZip(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	buf := make([]byte, 2)
	if _, err := file.Read(buf); err != nil {
		return false
	}
	return buf[0] == 'P' && buf[1] == 'K'
}

func extractEPFFromZip(zipPath, destination string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()
	var chosen *zip.File
	for _, file := range reader.File {
		if !strings.HasSuffix(strings.ToLower(file.Name), ".epf") {
			continue
		}
		if chosen == nil || strings.Contains(strings.ToLower(file.Name), "single") {
			chosen = file
		}
	}
	if chosen == nil {
		return fmt.Errorf("no .epf in archive")
	}
	src, err := chosen.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	out, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	return err
}

func copyFile(src, dest string) error {
	in, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dest, in, 0o644)
}
