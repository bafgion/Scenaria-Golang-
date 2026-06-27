package exporter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestWriteFeatureJSON(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "feature.json")

	doc := NewFeatureExportDocument("input.feature", &gherkin.Feature{
		Title: "Demo",
		Scenarios: []gherkin.Scenario{
			{
				Title: "S1",
				Steps: []gherkin.Step{{Keyword: "Когда", Text: "шаг"}},
			},
		},
	})

	if err := WriteFeatureJSON(path, doc); err != nil {
		t.Fatalf("WriteFeatureJSON returned error: %v", err)
	}

	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read export file: %v", err)
	}
	var decoded FeatureExportDocument
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("failed to decode export json: %v", err)
	}
	if decoded.Feature == nil || decoded.Feature.Title != "Demo" {
		t.Fatalf("unexpected decoded feature: %+v", decoded.Feature)
	}
}
