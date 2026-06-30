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
	if err := downloadFileWithClient(srv.Client(), srv.URL, dest, nil); err != nil {
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
	if err := downloadFileWithClient(srv.Client(), srv.URL, dest, nil); err == nil {
		t.Fatal("expected error for non-200")
	}
}

func TestDownloadFileReportsProgress(t *testing.T) {
	body := []byte("0123456789")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "10")
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	var lastDone, lastTotal int64
	dest := filepath.Join(t.TempDir(), "progress.zip")
	if err := downloadFileWithClient(srv.Client(), srv.URL, dest, func(done, total int64) {
		lastDone = done
		lastTotal = total
	}); err != nil {
		t.Fatalf("downloadFileWithClient: %v", err)
	}
	if lastDone != int64(len(body)) || lastTotal != int64(len(body)) {
		t.Fatalf("progress done=%d total=%d", lastDone, lastTotal)
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
	if err := downloadFileWithClient(srv.Client(), srv.URL, dest, nil); err != nil {
		t.Fatalf("downloadFileWithClient: %v", err)
	}
	got, _ := os.ReadFile(dest)
	if string(got) != "new" {
		t.Fatalf("content=%q", got)
	}
}
