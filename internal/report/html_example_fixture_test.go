package report

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/scenario"
)

// Writes frontend/e2e/fixtures/report/example-report.html from bundled examples/ (no browser).
func TestWriteReportExampleFixture(t *testing.T) {
	root := repoRoot(t)
	featurePath := filepath.Join(root, "examples", "01-pervaya-proverka.feature")
	store := scenario.NewFeatureStore()
	feature, err := store.Load(featurePath)
	if err != nil {
		t.Fatal(err)
	}
	if issues := gherkin.ValidateFeature(feature); len(issues) > 0 {
		t.Fatalf("validation: %+v", issues)
	}
	relPath := filepath.ToSlash(filepath.Join("examples", "01-pervaya-proverka.feature"))
	sc := feature.Scenarios[0]
	steps := make([]gherkin.Step, len(sc.Steps))
	copy(steps, sc.Steps)
	plan := player.ExecutionPlan{Cases: []player.RunCase{{
		FeaturePath: relPath,
		Name:        sc.Title,
		Tags:        sc.Tags,
		Steps:       steps,
	}}}
	recs := make([]player.StepRecord, len(steps))
	for i, st := range steps {
		recs[i] = player.StepRecord{
			Index: i, Line: st.Line, Keyword: st.Keyword, Text: st.Text,
			Status: "passed", DurationMS: 120 + int64(i*40),
		}
	}
	fs := len(steps) - 1
	recs[fs].Status = "failed"
	recs[fs].Error = "text mismatch"
	recs[fs].Selector = "h1"
	recs[fs].PageContext = "Example Domain\nhttps://example.com"
	result := player.ExecutionResult{
		Mode: "browser", Files: 1, Scenarios: 1, Steps: len(steps),
		ScenarioResults: []player.ScenarioResult{{
			FeaturePath: relPath, Scenario: sc.Title, Status: "failed",
			Message: "text mismatch", FailedStep: &fs, DurationMS: 900,
			StepRecords: recs,
		}},
	}
	outDir := filepath.Join(root, "frontend", "e2e", "fixtures", "report")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(outDir, "example-report.html")
	if err := WriteHTML(out, result, HTMLOptions{
		Plan: plan, Locale: "ru", ProjectRoot: t.TempDir(),
	}); err != nil {
		t.Fatal(err)
	}
}
