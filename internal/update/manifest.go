package update

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type releaseManifest struct {
	Version string `json:"version"`
	Assets  map[string]struct {
		Name   string `json:"name"`
		Size   int64  `json:"size"`
		SHA256 string `json:"sha256"`
	} `json:"assets"`
}

func sha256FromManifest(assets []ReleaseAsset, fileName string) string {
	for _, asset := range assets {
		if !strings.EqualFold(strings.TrimSpace(asset.Name), "latest.json") {
			continue
		}
		url := strings.TrimSpace(asset.BrowserDownloadURL)
		if url == "" {
			return ""
		}
		manifest, err := fetchManifest(url)
		if err != nil {
			return ""
		}
		for _, entry := range manifest.Assets {
			if strings.EqualFold(entry.Name, fileName) {
				return strings.ToLower(strings.TrimSpace(entry.SHA256))
			}
		}
		return ""
	}
	return ""
}

func fetchManifest(url string) (*releaseManifest, error) {
	return fetchManifestWithClient(&http.Client{Timeout: 30 * time.Second}, url)
}

func fetchManifestWithClient(client *http.Client, url string) (*releaseManifest, error) {
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("latest.json returned %s", resp.Status)
	}
	var manifest releaseManifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, fmt.Errorf("decode latest.json: %w", err)
	}
	return &manifest, nil
}

func verifyFileSHA256(path, want string) error {
	want = strings.ToLower(strings.TrimSpace(want))
	if want == "" {
		return nil
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}
	got := hex.EncodeToString(hash.Sum(nil))
	if got != want {
		return fmt.Errorf("checksum mismatch for %s", filepathBase(path))
	}
	return nil
}

func filepathBase(path string) string {
	path = strings.ReplaceAll(path, "\\", "/")
	if i := strings.LastIndex(path, "/"); i >= 0 {
		return path[i+1:]
	}
	return path
}
