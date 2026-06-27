package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestRunExport_JSON(t *testing.T) {
	tmp := t.TempDir()
	in := filepath.Join(tmp, "in.feature")
	out := filepath.Join(tmp, "out.json")
	content := "Функционал: Demo\nСценарий: S1\nКогда выполняю шаг\n"
	if err := os.WriteFile(in, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write input: %v", err)
	}

	if err := RunExport([]string{in, "--output", out, "--format", "json"}); err != nil {
		t.Fatalf("RunExport returned error: %v", err)
	}
	payload, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("failed to read json export: %v", err)
	}
	if !strings.Contains(string(payload), `"feature"`) {
		t.Fatalf("unexpected export payload: %s", string(payload))
	}
}

func TestRunExport_Feature(t *testing.T) {
	tmp := t.TempDir()
	in := filepath.Join(tmp, "in.feature")
	out := filepath.Join(tmp, "out.feature")
	content := "Функционал: Demo\nСценарий: S1\nКогда выполняю шаг\n"
	if err := os.WriteFile(in, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write input: %v", err)
	}

	if err := RunExport([]string{in, "--output", out, "--format", "feature"}); err != nil {
		t.Fatalf("RunExport returned error: %v", err)
	}
	if _, err := gherkin.ParseFeatureFile(out); err != nil {
		t.Fatalf("expected valid exported feature, got parse error: %v", err)
	}
}

func TestParseExportOptions(t *testing.T) {
	opts, err := parseExportOptions([]string{"./in.feature", "--output", "out.json"})
	if err != nil {
		t.Fatalf("parseExportOptions returned error: %v", err)
	}
	if opts.format != "json" || opts.output != "out.json" || opts.input != "./in.feature" {
		t.Fatalf("unexpected options: %+v", opts)
	}
}

func TestParseExportOptionsErrors(t *testing.T) {
	if _, err := parseExportOptions(nil); err == nil {
		t.Fatal("expected usage error")
	}
	if _, err := parseExportOptions([]string{"./in.feature"}); err == nil {
		t.Fatal("expected missing output error")
	}
	if _, err := parseExportOptions([]string{"./in.feature", "--output"}); err == nil {
		t.Fatal("expected missing output path error")
	}
	if _, err := parseExportOptions([]string{"./in.feature", "--output", "out.txt", "--format", "xml"}); err == nil {
		t.Fatal("expected unsupported format error")
	}
}
