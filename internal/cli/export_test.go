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

func TestRunExport_TSAndPython(t *testing.T) {
	tmp := t.TempDir()
	in := filepath.Join(tmp, "in.feature")
	content := "Функционал: Demo\nСценарий: S1\nКогда открываю \"/login\"\nТогда вижу \"Добро\" \n"
	if err := os.WriteFile(in, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write input: %v", err)
	}

	tsOut := filepath.Join(tmp, "out.spec.ts")
	if err := RunExport([]string{in, "--output", tsOut, "--format", "ts", "--base-url", "https://example.com"}); err != nil {
		t.Fatalf("RunExport ts returned error: %v", err)
	}
	tsPayload, err := os.ReadFile(tsOut)
	if err != nil {
		t.Fatalf("failed to read ts export: %v", err)
	}
	if !strings.Contains(string(tsPayload), "page.goto(\"https://example.com/login\")") {
		t.Fatalf("unexpected ts export: %s", string(tsPayload))
	}

	pyOut := filepath.Join(tmp, "out.py")
	if err := RunExport([]string{in, "--output", pyOut, "--format", "python", "--base-url", "https://example.com"}); err != nil {
		t.Fatalf("RunExport python returned error: %v", err)
	}
	pyPayload, err := os.ReadFile(pyOut)
	if err != nil {
		t.Fatalf("failed to read python export: %v", err)
	}
	if !strings.Contains(string(pyPayload), "page.goto(\"https://example.com/login\")") {
		t.Fatalf("unexpected python export: %s", string(pyPayload))
	}
}

func TestParseExportOptions(t *testing.T) {
	opts, err := parseExportOptions([]string{"./in.feature", "--output", "out.json", "--base-url", "https://example.com"})
	if err != nil {
		t.Fatalf("parseExportOptions returned error: %v", err)
	}
	if opts.format != "json" || opts.output != "out.json" || opts.input != "./in.feature" || opts.baseURL != "https://example.com" {
		t.Fatalf("unexpected options: %+v", opts)
	}

	opts, err = parseExportOptions([]string{"./in.feature", "--output", "out.spec.ts"})
	if err != nil {
		t.Fatalf("parseExportOptions inferred ts returned error: %v", err)
	}
	if opts.format != "ts" {
		t.Fatalf("expected ts format inference, got %+v", opts)
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
	if _, err := parseExportOptions([]string{"./in.feature", "--output", "out.json", "--base-url"}); err == nil {
		t.Fatal("expected missing base-url error")
	}
	if _, err := parseExportOptions([]string{"./in.feature", "--output", "out.txt", "--format", "xml"}); err == nil {
		t.Fatal("expected unsupported format error")
	}
}
