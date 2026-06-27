package exporter

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

type FeatureExportDocument struct {
	GeneratedAt string           `json:"generated_at"`
	Source      string           `json:"source"`
	Feature     *gherkin.Feature `json:"feature"`
}

func NewFeatureExportDocument(source string, feature *gherkin.Feature) FeatureExportDocument {
	return FeatureExportDocument{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Source:      source,
		Feature:     feature,
	}
}

func WriteFeatureJSON(path string, document FeatureExportDocument) error {
	payload, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return fmt.Errorf("encode feature export %q: %w", path, err)
	}
	if err := os.WriteFile(path, append(payload, '\n'), 0o644); err != nil {
		return fmt.Errorf("write feature export %q: %w", path, err)
	}
	return nil
}
