package plugin

import (
	"archive/zip"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDownloadToTempRejectsOversizedBody(t *testing.T) {
	body := strings.Repeat("x", int(MaxPluginDownloadBytes)+1024)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	_, err := downloadToTemp(srv.URL)
	if err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("expected size limit error, got %v", err)
	}
}

func TestFetchAndInstallRemoteZip(t *testing.T) {
	zipPath := filepath.Join(t.TempDir(), "addon.zip")
	file, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(file)
	entry, err := w.Create("plugin.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = entry.Write([]byte(`{"name":"demo"}`))
	_ = w.Close()
	_ = file.Close()

	data, err := os.ReadFile(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(data)
	}))
	defer srv.Close()

	project := filepath.Join(t.TempDir(), "project")
	if err := FetchAndInstall(project, "demo", srv.URL); err != nil {
		t.Fatalf("FetchAndInstall: %v", err)
	}
	if _, err := os.Stat(filepath.Join(project, "addons", "demo", "plugin.json")); err != nil {
		t.Fatalf("expected extracted plugin.json: %v", err)
	}
}
