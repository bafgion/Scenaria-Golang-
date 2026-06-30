package update

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/brand"
)

func TestCheckInstallDirSetupMode(t *testing.T) {
	oldURL := releaseLatestURL
	defer func() { releaseLatestURL = oldURL }()

	var baseURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/release":
			_ = json.NewEncoder(w).Encode(releasePayload{
				TagName: "v0.17.0",
				HTMLURL: "https://github.com/example/releases/v0.17.0",
				Assets: []ReleaseAsset{
					{Name: "latest.json", BrowserDownloadURL: baseURL + "/latest.json"},
					{Name: brand.PortableZip, BrowserDownloadURL: baseURL + "/portable.zip"},
					{Name: brand.SetupExe, BrowserDownloadURL: baseURL + "/setup.exe"},
				},
			})
		case "/latest.json":
			_ = json.NewEncoder(w).Encode(releaseManifest{
				Version: "v0.17.0",
				Assets: map[string]struct {
					Name   string `json:"name"`
					Size   int64  `json:"size"`
					SHA256 string `json:"sha256"`
				}{
					brand.SetupExe: {Name: brand.SetupExe, SHA256: "deadbeef"},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	baseURL = srv.URL
	releaseLatestURL = baseURL + "/release"

	oldDetect := detectInstallModeOverride
	detectInstallModeOverride = func(string) InstallMode { return InstallModeSetup }
	defer func() { detectInstallModeOverride = oldDetect }()

	info, err := CheckInstallDir("0.16.0", `C:\Program Files\Scenaria`)
	if err != nil {
		t.Fatalf("CheckInstallDir: %v", err)
	}
	if !info.UpdateAvailable {
		t.Fatal("expected update available")
	}
	if info.DownloadName != brand.SetupExe {
		t.Fatalf("download name=%q want %q", info.DownloadName, brand.SetupExe)
	}
	if info.ApplyKind != string(ApplyKindSetup) {
		t.Fatalf("apply kind=%q", info.ApplyKind)
	}
	if info.InstallMode != string(InstallModeSetup) {
		t.Fatalf("install mode=%q", info.InstallMode)
	}
	if info.SHA256 != "deadbeef" {
		t.Fatalf("sha256=%q", info.SHA256)
	}
}

func TestCheckInstallDirPortableMode(t *testing.T) {
	oldURL := releaseLatestURL
	defer func() { releaseLatestURL = oldURL }()

	var baseURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(releasePayload{
			TagName: "v0.17.0",
			HTMLURL: "https://github.com/example/releases/v0.17.0",
			Assets: []ReleaseAsset{
				{Name: brand.PortableZip, BrowserDownloadURL: baseURL + "/portable.zip"},
				{Name: brand.SetupExe, BrowserDownloadURL: baseURL + "/setup.exe"},
			},
		})
	}))
	defer srv.Close()
	baseURL = srv.URL
	releaseLatestURL = baseURL

	oldDetect := detectInstallModeOverride
	detectInstallModeOverride = func(string) InstallMode { return InstallModePortable }
	defer func() { detectInstallModeOverride = oldDetect }()

	info, err := CheckInstallDir("0.16.0", `D:\Scenaria`)
	if err != nil {
		t.Fatalf("CheckInstallDir: %v", err)
	}
	if info.DownloadName != brand.PortableZip {
		t.Fatalf("download name=%q want %q", info.DownloadName, brand.PortableZip)
	}
	if info.ApplyKind != string(ApplyKindPortable) {
		t.Fatalf("apply kind=%q", info.ApplyKind)
	}
}

func TestCheckInstallDirUpToDate(t *testing.T) {
	oldURL := releaseLatestURL
	defer func() { releaseLatestURL = oldURL }()

	var baseURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(releasePayload{
			TagName: "v0.16.0",
			Assets:  []ReleaseAsset{{Name: brand.PortableZip, BrowserDownloadURL: baseURL + "/p.zip"}},
		})
	}))
	defer srv.Close()
	baseURL = srv.URL
	releaseLatestURL = baseURL

	info, err := CheckInstallDir("v0.16.0", "")
	if err != nil {
		t.Fatalf("CheckInstallDir: %v", err)
	}
	if info.UpdateAvailable {
		t.Fatal("expected no update")
	}
}

func TestFetchReleasePayloadNonOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	if _, err := fetchReleasePayloadFrom(srv.Client(), srv.URL); err == nil {
		t.Fatal("expected error")
	}
}
