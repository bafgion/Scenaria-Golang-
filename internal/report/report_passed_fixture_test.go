package report

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/scenario"
)

func TestWriteReportPassedFixture(t *testing.T) {
	root := repoRoot(t)
	featurePath := filepath.Join(root, "examples", "01-pervaya-proverka.feature")
	store := scenario.NewFeatureStore()
	feature, err := store.Load(featurePath)
	if err != nil {
		t.Fatal(err)
	}
	sc := feature.Scenarios[0]
	steps := make([]gherkin.Step, len(sc.Steps))
	copy(steps, sc.Steps)
	relPath := filepath.ToSlash(filepath.Join("examples", "01-pervaya-proverka.feature"))
	recs := make([]player.StepRecord, len(steps))
	var total int64
	for i, st := range steps {
		dur := int64(150 + i*50)
		total += dur
		recs[i] = player.StepRecord{
			Index: i, Line: st.Line, Keyword: st.Keyword, Text: st.Text,
			Status: "passed", DurationMS: dur,
			PageContext: "Example Domain\nhttps://example.com",
		}
	}
	plan := player.ExecutionPlan{Cases: []player.RunCase{{
		FeaturePath: relPath,
		Name:        sc.Title,
		Tags:        sc.Tags,
		Steps:       steps,
	}}}
	result := player.ExecutionResult{
		Mode: "browser", Files: 1, Scenarios: 1, Steps: len(steps),
		ScenarioResults: []player.ScenarioResult{{
			FeaturePath: relPath, Scenario: sc.Title, Status: "passed",
			DurationMS: total, StepRecords: recs,
		}},
	}
	outDir := filepath.Join(root, "frontend", "e2e", "fixtures", "report")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(outDir, "sample-passed.html")
	if err := WriteHTML(out, result, HTMLOptions{Plan: plan, Locale: "ru"}); err != nil {
		t.Fatal(err)
	}
}
