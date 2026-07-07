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
