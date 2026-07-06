package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bafgion/scenaria-golang/internal/player"
)

// ScenarioSummary is one scenario row in a detailed run summary file.
type ScenarioSummary struct {
	Path       string `json:"path"`
	Scenario   string `json:"scenario"`
	Status     string `json:"status"`
	DurationMS int64  `json:"duration_ms,omitempty"`
}

// RunSummaryDetailed extends the CLI summary with per-scenario stats for CI/HTML diff.
type RunSummaryDetailed struct {
	GeneratedAt string            `json:"generated_at"`
	Mode        string            `json:"mode"`
	Files       int               `json:"files"`
	Scenarios   int               `json:"scenarios"`
	Steps       int               `json:"steps"`
	Passed      int               `json:"passed"`
	Failed      int               `json:"failed"`
	Skipped     int               `json:"skipped"`
	Items       []ScenarioSummary `json:"items,omitempty"`
}

func FromExecutionResultDetailed(result player.ExecutionResult) RunSummaryDetailed {
	out := RunSummaryDetailed{
		GeneratedAt: FromExecutionResult(result).GeneratedAt,
		Mode:        result.Mode,
		Files:       result.Files,
		Scenarios:   result.Scenarios,
		Steps:       result.Steps,
		Items:       make([]ScenarioSummary, 0, len(result.ScenarioResults)),
	}
	for _, sr := range result.ScenarioResults {
		switch sr.Status {
		case "passed":
			out.Passed++
		case "failed":
			out.Failed++
		default:
			out.Skipped++
		}
		out.Items = append(out.Items, ScenarioSummary{
			Path:       sr.FeaturePath,
			Scenario:   sr.Scenario,
			Status:     sr.Status,
			DurationMS: sr.DurationMS,
		})
	}
	return out
}

func WriteRunSummaryDetailed(path string, summary RunSummaryDetailed) error {
	payload, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("encode run summary %q: %w", path, err)
	}
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create run summary dir %q: %w", dir, err)
		}
	}
	if err := os.WriteFile(path, append(payload, '\n'), 0o644); err != nil {
		return fmt.Errorf("write run summary %q: %w", path, err)
	}
	return nil
}

func LoadRunSummaryDetailed(path string) (*RunSummaryDetailed, error) {
	path = filepath.Clean(path)
	payload, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var summary RunSummaryDetailed
	if err := json.Unmarshal(payload, &summary); err != nil {
		return nil, fmt.Errorf("decode summary %q: %w", path, err)
	}
	return &summary, nil
}

// ReadPreviousSummary loads an existing summary file before it is overwritten by the current run.
func ReadPreviousSummary(path string) *RunSummaryDetailed {
	if path == "" {
		return nil
	}
	prev, err := LoadRunSummaryDetailed(path)
	if err != nil {
		return nil
	}
	return prev
}
