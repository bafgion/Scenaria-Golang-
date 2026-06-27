package report

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/bafgion/scenaria-golang/internal/player"
)

type RunSummary struct {
	GeneratedAt string `json:"generated_at"`
	Mode        string `json:"mode"`
	Files       int    `json:"files"`
	Scenarios   int    `json:"scenarios"`
	Steps       int    `json:"steps"`
}

func FromExecutionResult(result player.ExecutionResult) RunSummary {
	return RunSummary{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Mode:        result.Mode,
		Files:       result.Files,
		Scenarios:   result.Scenarios,
		Steps:       result.Steps,
	}
}

func WriteRunSummary(path string, summary RunSummary) error {
	payload, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("encode run summary %q: %w", path, err)
	}
	if err := os.WriteFile(path, append(payload, '\n'), 0o644); err != nil {
		return fmt.Errorf("write run summary %q: %w", path, err)
	}
	return nil
}
