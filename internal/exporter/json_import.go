package exporter

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func ReadFeatureJSON(path string) (FeatureExportDocument, error) {
	payload, err := os.ReadFile(path)
	if err != nil {
		return FeatureExportDocument{}, fmt.Errorf("read feature json %q: %w", path, err)
	}
	var doc FeatureExportDocument
	if err := json.Unmarshal(payload, &doc); err != nil {
		return FeatureExportDocument{}, fmt.Errorf("decode feature json %q: %w", path, err)
	}
	if doc.Feature == nil {
		return FeatureExportDocument{}, fmt.Errorf("feature json %q has no feature object", path)
	}
	if issues := gherkin.ValidateFeature(doc.Feature); len(issues) > 0 {
		return FeatureExportDocument{}, fmt.Errorf("feature json %q: %s", path, issues[0].Message)
	}
	return doc, nil
}
