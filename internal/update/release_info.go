package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/brand"
)

var releaseLatestURL = "https://api.github.com/repos/" + brand.DefaultGitHubRepo + "/releases/latest"

type ReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type releasePayload struct {
	TagName string         `json:"tag_name"`
	HTMLURL string         `json:"html_url"`
	Assets  []ReleaseAsset `json:"assets"`
}

type Info struct {
	CurrentVersion  string `json:"currentVersion"`
	LatestVersion   string `json:"latestVersion"`
	UpdateAvailable bool   `json:"updateAvailable"`
	HTMLURL         string `json:"htmlUrl"`
	DownloadURL     string `json:"downloadUrl"`
	DownloadName    string `json:"downloadName"`
	Message         string `json:"message"`
	InstallMode     string `json:"installMode"`
	ApplyKind       string `json:"applyKind"`
	SHA256          string `json:"sha256,omitempty"`
}

// Check reports whether a newer GitHub release exists and picks the update asset for this install.
func Check(currentVersion string) (*Info, error) {
	return CheckInstallDir(currentVersion, "")
}

func CheckInstallDir(currentVersion, installDir string) (*Info, error) {
	release, err := fetchReleasePayload(&http.Client{Timeout: 20 * time.Second})
	if err != nil {
		return nil, err
	}
	mode := InstallModePortable
	if strings.TrimSpace(installDir) != "" {
		mode = detectInstallMode(installDir)
	}
	info := &Info{
		CurrentVersion: strings.TrimSpace(currentVersion),
		LatestVersion:  strings.TrimSpace(release.TagName),
		HTMLURL:        strings.TrimSpace(release.HTMLURL),
		InstallMode:    string(mode),
	}
	if asset := pickAssetForMode(mode, release.Assets); asset != nil {
		info.DownloadURL = strings.TrimSpace(asset.BrowserDownloadURL)
		info.DownloadName = strings.TrimSpace(asset.Name)
		info.ApplyKind = string(applyKindForAsset(info.DownloadName))
		info.SHA256 = sha256FromManifest(release.Assets, info.DownloadName)
	}
	info.UpdateAvailable = IsNewer(info.CurrentVersion, info.LatestVersion)
	if info.UpdateAvailable {
		info.Message = fmt.Sprintf("Доступна версия %s (у вас %s)", info.LatestVersion, info.CurrentVersion)
	} else {
		info.Message = fmt.Sprintf("Установлена актуальная версия (%s)", info.CurrentVersion)
	}
	return info, nil
}

func fetchReleasePayload(client *http.Client) (*releasePayload, error) {
	return fetchReleasePayloadFrom(client, releaseLatestURL)
}

func fetchReleasePayloadFrom(client *http.Client, apiURL string) (*releasePayload, error) {
	if client == nil {
		client = &http.Client{Timeout: 20 * time.Second}
	}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch latest release: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github release API returned %s", resp.Status)
	}
	var payload releasePayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode release metadata: %w", err)
	}
	return &payload, nil
}

func pickUpdateAsset(assets []ReleaseAsset) *ReleaseAsset {
	prefs := updateAssetPreferences(runtime.GOOS)
	for _, pref := range prefs {
		if asset := matchAsset(assets, pref); asset != nil {
			return asset
		}
	}
	return nil
}

func updateAssetPreferences(goos string) []func(string) bool {
	switch goos {
	case "windows":
		return []func(string) bool{
			func(name string) bool { return strings.EqualFold(name, brand.SetupExe) },
			func(name string) bool { return strings.HasSuffix(name, ".msi") },
			func(name string) bool { return strings.Contains(name, "setup") && strings.HasSuffix(name, ".exe") },
			func(name string) bool { return strings.Contains(name, "installer") },
			func(name string) bool { return strings.Contains(name, "portable") && strings.HasSuffix(name, ".zip") },
			func(name string) bool { return strings.HasSuffix(name, ".zip") },
			func(name string) bool { return strings.HasSuffix(name, ".exe") },
		}
	case "darwin":
		return []func(string) bool{
			func(name string) bool { return strings.HasSuffix(name, ".dmg") },
			func(name string) bool { return strings.HasSuffix(name, ".pkg") },
			func(name string) bool { return strings.Contains(name, "portable") && strings.HasSuffix(name, ".zip") },
			func(name string) bool { return strings.HasSuffix(name, ".zip") },
		}
	default:
		return []func(string) bool{
			func(name string) bool { return strings.HasSuffix(name, ".appimage") },
			func(name string) bool { return strings.HasSuffix(name, ".deb") },
			func(name string) bool { return strings.Contains(name, "portable") && strings.HasSuffix(name, ".zip") },
			func(name string) bool { return strings.HasSuffix(name, ".zip") },
		}
	}
}

func matchAsset(assets []ReleaseAsset, pred func(string) bool) *ReleaseAsset {
	for i := range assets {
		name := strings.ToLower(strings.TrimSpace(assets[i].Name))
		if pred(name) {
			return &assets[i]
		}
	}
	return nil
}

func DownloadFile(url, destPath string) error {
	return downloadFileWithClient(&http.Client{Timeout: 10 * time.Minute}, url, destPath)
}

func downloadFileWithClient(client *http.Client, url, destPath string) error {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Minute}
	}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("download update: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download update returned %s", resp.Status)
	}
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return err
	}
	tmp := destPath + ".part"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(out, resp.Body)
	closeErr := out.Close()
	if copyErr != nil {
		_ = os.Remove(tmp)
		return copyErr
	}
	if closeErr != nil {
		_ = os.Remove(tmp)
		return closeErr
	}
	if err := os.Remove(destPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.Rename(tmp, destPath)
}
