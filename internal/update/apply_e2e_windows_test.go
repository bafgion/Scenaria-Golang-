//go:build windows

package update

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/brand"
)

func TestApplyEndToEndPortable(t *testing.T) {
	oldURL := releaseLatestURL
	defer func() { releaseLatestURL = oldURL }()

	zipBody := buildPortableZipBytes(t)
	sum := sha256.Sum256(zipBody)
	checksum := hex.EncodeToString(sum[:])

	var baseURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/release":
			_ = json.NewEncoder(w).Encode(releasePayload{
				TagName: "v9.9.9",
				HTMLURL: "https://example.com/v9.9.9",
				Assets: []ReleaseAsset{
					{Name: "latest.json", BrowserDownloadURL: baseURL + "/latest.json"},
					{Name: brand.PortableZip, BrowserDownloadURL: baseURL + "/portable.zip"},
				},
			})
		case "/latest.json":
			_ = json.NewEncoder(w).Encode(releaseManifest{
				Assets: map[string]struct {
					Name   string `json:"name"`
					Size   int64  `json:"size"`
					SHA256 string `json:"sha256"`
				}{
					brand.PortableZip: {Name: brand.PortableZip, SHA256: checksum},
				},
			})
		case "/portable.zip":
			_, _ = w.Write(zipBody)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	baseURL = srv.URL
	releaseLatestURL = baseURL + "/release"

	installDir := t.TempDir()
	oldDetect := detectInstallModeOverride
	detectInstallModeOverride = func(string) InstallMode { return InstallModePortable }
	defer func() { detectInstallModeOverride = oldDetect }()

	launched := false
	oldHook := launchHiddenBatchHook
	launchHiddenBatchHook = func(string) error {
		launched = true
		return nil
	}
	defer func() { launchHiddenBatchHook = oldHook }()

	info, err := Apply("v0.1.0", installDir, 999)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if info == nil || !info.UpdateAvailable {
		t.Fatalf("info=%+v", info)
	}
	if !launched {
		t.Fatal("expected update script launch")
	}
	if _, err := os.Stat(filepath.Join(installDir, updateStagingDir, brand.GUIExeName)); err != nil {
		t.Fatalf("staged binary: %v", err)
	}
}

func buildPortableZipBytes(t *testing.T) []byte {
	t.Helper()
	tmp := filepath.Join(t.TempDir(), "payload.zip")
	writePortableZip(t, tmp, map[string]string{
		"Scenaria/" + brand.GUIExeName: "new-gui",
	})
	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestApplyRejectsWhenUpToDate(t *testing.T) {
	oldURL := releaseLatestURL
	defer func() { releaseLatestURL = oldURL }()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(releasePayload{TagName: "v1.0.0"})
	}))
	defer srv.Close()
	releaseLatestURL = srv.URL

	_, err := Apply("v1.0.0", t.TempDir(), 1)
	if err == nil {
		t.Fatal("expected error when no update")
	}
}

func TestApplyRejectsEmptyInstallDir(t *testing.T) {
	_, err := Apply("v0.1.0", "   ", 1)
	if err == nil {
		t.Fatal("expected error for empty install dir")
	}
}

// Ensure zip builder used by end-to-end test matches portable layout.
func TestPortableZipLayout(t *testing.T) {
	path := filepath.Join(t.TempDir(), "layout.zip")
	writePortableZip(t, path, map[string]string{"Scenaria/" + brand.GUIExeName: "x"})
	reader, err := zip.OpenReader(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	if len(reader.File) != 1 {
		t.Fatalf("files=%d", len(reader.File))
	}
}
