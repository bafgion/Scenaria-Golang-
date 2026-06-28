package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/exporter"
	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestRunImportJSON_RoundTrip(t *testing.T) {
	tmp := t.TempDir()
	in := filepath.Join(tmp, "in.feature")
	outJSON := filepath.Join(tmp, "out.json")
	outFeature := filepath.Join(tmp, "restored.feature")
	content := "Функционал: Demo\nСценарий: S1\nКогда открываю \"https://example.com\"\n"
	if err := os.WriteFile(in, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	feature, err := gherkin.ParseFeatureFile(in)
	if err != nil {
		t.Fatal(err)
	}
	if err := exporter.WriteFeatureJSON(outJSON, exporter.NewFeatureExportDocument(in, feature)); err != nil {
		t.Fatal(err)
	}
	if err := RunImportJSON([]string{outJSON, "--output", outFeature}); err != nil {
		t.Fatalf("RunImportJSON: %v", err)
	}
	restored, err := gherkin.ParseFeatureFile(outFeature)
	if err != nil {
		t.Fatal(err)
	}
	if restored.Title != "Demo" {
		t.Fatalf("unexpected title: %q", restored.Title)
	}
}
