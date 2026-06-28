package gherkin_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

func TestAllExamplesParseAndValidate(t *testing.T) {
	examplesDir := filepath.Join("..", "..", "examples")
	entries, err := os.ReadDir(examplesDir)
	if err != nil {
		t.Fatalf("read examples: %v", err)
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".feature" {
			continue
		}
		path := filepath.Join(examplesDir, entry.Name())
		t.Run(entry.Name(), func(t *testing.T) {
			feature, err := gherkin.ParseFeatureFile(path)
			if err != nil {
				t.Fatalf("ParseFeatureFile failed: %v", err)
			}
			if issues := gherkin.ValidateFeature(feature); len(issues) > 0 {
				t.Fatalf("validation issues: %+v", issues)
			}
			for _, runnable := range gherkin.ExpandFeature(feature) {
				for _, step := range gherkin.FlattenSteps(runnable.Steps) {
					if step.Block != "" || gherkin.IsTestClientStep(step) {
						continue
					}
					if _, err := stepdsl.Parse(step); err != nil {
						t.Fatalf("scenario %q: unsupported step %q: %v", runnable.Title, step.Text, err)
					}
				}
			}
		})
	}
}
