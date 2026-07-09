package report

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPreferredHTMLReportPath(t *testing.T) {
	dir := t.TempDir()
	full := filepath.Join(dir, "report.html")
	light := filepath.Join(dir, "report.light.html")
	if err := os.WriteFile(full, []byte("full"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(light, []byte("light"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := PreferredHTMLReportPath(full, "full"); got != full {
		t.Fatalf("full mode: got %q want %q", got, full)
	}
	if got := PreferredHTMLReportPath(full, "light"); got != light {
		t.Fatalf("light mode: got %q want %q", got, light)
	}
}

func TestPreferredHTMLReportPathFallback(t *testing.T) {
	dir := t.TempDir()
	full := filepath.Join(dir, "report.html")
	if err := os.WriteFile(full, []byte("full"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := PreferredHTMLReportPath(full, "light"); got != full {
		t.Fatalf("fallback to full: got %q", got)
	}
}

func TestResolveHTMLReportPathPrefersLatestRun(t *testing.T) {
	root := t.TempDir()
	scenaria := filepath.Join(root, ".scenaria")
	legacy := filepath.Join(scenaria, "report.html")
	latestHTML := filepath.Join(scenaria, "runs", "run-1", "report.html")
	if err := os.MkdirAll(filepath.Dir(latestHTML), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(legacy, []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(latestHTML, []byte("fresh"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := WriteLatestRunPointer(root, RunArtifactLayout{
		RunID:    "run-1",
		Dir:      filepath.Dir(latestHTML),
		HTMLPath: latestHTML,
	}); err != nil {
		t.Fatal(err)
	}

	got, err := ResolveHTMLReportPath(root, legacy)
	if err != nil {
		t.Fatal(err)
	}
	if got != latestHTML {
		t.Fatalf("ResolveHTMLReportPath = %q, want %q", got, latestHTML)
	}
	got, err = ResolveHTMLReportPath(root, "")
	if err != nil {
		t.Fatal(err)
	}
	if got != latestHTML {
		t.Fatalf("empty path = %q, want %q", got, latestHTML)
	}
}
