package gherkin_test

import (
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

func TestExample01ParsesAndValidates(t *testing.T) {
	path := filepath.Join("..", "..", "examples", "01-pervaya-proverka.feature")
	feature, err := gherkin.ParseFeatureFile(path)
	if err != nil {
		t.Fatalf("ParseFeatureFile failed: %v", err)
	}

	issues := gherkin.ValidateFeature(feature)
	if len(issues) > 0 {
		t.Fatalf("validation issues: %+v", issues)
	}

	if len(feature.Scenarios) != 1 || len(feature.Scenarios[0].Steps) != 4 {
		t.Fatalf("unexpected scenario structure: %+v", feature.Scenarios)
	}

	for _, step := range feature.Scenarios[0].Steps {
		if _, err := stepdsl.Parse(step); err != nil {
			t.Fatalf("unsupported step %q: %v", step.Text, err)
		}
	}
}
