package update

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadFileWithClient(t *testing.T) {
	body := []byte("portable-zip-bytes")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	dest := filepath.Join(t.TempDir(), "Scenaria-Portable.zip")
	if err := downloadFileWithClient(srv.Client(), srv.URL, dest); err != nil {
		t.Fatalf("downloadFileWithClient: %v", err)
	}
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(body) {
		t.Fatalf("content=%q", got)
	}
}

func TestDownloadFileWithClientNonOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	dest := filepath.Join(t.TempDir(), "fail.zip")
	if err := downloadFileWithClient(srv.Client(), srv.URL, dest); err == nil {
		t.Fatal("expected error for non-200")
	}
}

func TestDownloadFileReplacesExisting(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("new"))
	}))
	defer srv.Close()

	dest := filepath.Join(t.TempDir(), "update.zip")
	if err := os.WriteFile(dest, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := downloadFileWithClient(srv.Client(), srv.URL, dest); err != nil {
		t.Fatalf("downloadFileWithClient: %v", err)
	}
	got, _ := os.ReadFile(dest)
	if string(got) != "new" {
		t.Fatalf("content=%q", got)
	}
}
