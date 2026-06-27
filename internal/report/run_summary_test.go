package report

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/player"
)

func TestWriteRunSummary(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "summary.json")

	summary := FromExecutionResult(player.ExecutionResult{
		Mode:      "dry-run",
		Files:     2,
		Scenarios: 5,
		Steps:     12,
	})
	if err := WriteRunSummary(path, summary); err != nil {
		t.Fatalf("WriteRunSummary returned error: %v", err)
	}

	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read summary file: %v", err)
	}

	var decoded RunSummary
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("failed to decode summary JSON: %v", err)
	}
	if decoded.Mode != "dry-run" || decoded.Files != 2 || decoded.Scenarios != 5 || decoded.Steps != 12 {
		t.Fatalf("unexpected summary values: %+v", decoded)
	}
	if decoded.GeneratedAt == "" {
		t.Fatal("generated_at must not be empty")
	}
}
