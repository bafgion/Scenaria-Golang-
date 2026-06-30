package update

import (
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

func TestFetchManifestWithClient(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(releaseManifest{
			Version: "v1.0.0",
			Assets: map[string]struct {
				Name   string `json:"name"`
				Size   int64  `json:"size"`
				SHA256 string `json:"sha256"`
			}{
				brand.PortableZip: {Name: brand.PortableZip, SHA256: "abc123"},
			},
		})
	}))
	defer srv.Close()

	manifest, err := fetchManifestWithClient(srv.Client(), srv.URL)
	if err != nil {
		t.Fatalf("fetchManifestWithClient: %v", err)
	}
	if manifest.Version != "v1.0.0" {
		t.Fatalf("version=%q", manifest.Version)
	}
}

func TestSHA256FromManifest(t *testing.T) {
	want := "cafebabe"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(releaseManifest{
			Assets: map[string]struct {
				Name   string `json:"name"`
				Size   int64  `json:"size"`
				SHA256 string `json:"sha256"`
			}{
				brand.SetupExe: {Name: brand.SetupExe, SHA256: want},
			},
		})
	}))
	defer srv.Close()

	assets := []ReleaseAsset{
		{Name: "latest.json", BrowserDownloadURL: srv.URL},
	}
	got := sha256FromManifest(assets, brand.SetupExe)
	if got != want {
		t.Fatalf("sha256=%q want %q", got, want)
	}
}

func TestVerifyFileSHA256Match(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "payload.bin")
	content := []byte("scenaria-update-payload")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(content)
	want := hex.EncodeToString(sum[:])
	if err := verifyFileSHA256(path, want); err != nil {
		t.Fatalf("verifyFileSHA256: %v", err)
	}
}

func TestVerifyFileSHA256Mismatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "payload.bin")
	if err := os.WriteFile(path, []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := verifyFileSHA256(path, "00"); err == nil {
		t.Fatal("expected checksum mismatch")
	}
}

func TestVerifyFileSHA256SkipsEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing.bin")
	if err := verifyFileSHA256(path, ""); err != nil {
		t.Fatalf("empty want should skip verification: %v", err)
	}
}
