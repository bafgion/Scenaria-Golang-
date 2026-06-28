package report

import (
	"os"
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/player"
)

func TestWriteHTML(t *testing.T) {
	tmp := t.TempDir()
	path := tmp + "/report.html"
	result := player.ExecutionResult{
		Mode:      "dry-run",
		Files:     1,
		Scenarios: 1,
		Steps:     2,
		ScenarioResults: []player.ScenarioResult{
			{FeaturePath: "demo.feature", Scenario: "S1", Status: "skipped", Message: "dry-run"},
		},
	}
	if err := WriteHTML(path, result); err != nil {
		t.Fatalf("WriteHTML failed: %v", err)
	}
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read html: %v", err)
	}
	text := string(payload)
	if !strings.Contains(text, "Scenaria Run Report") || !strings.Contains(text, "demo.feature") {
		t.Fatalf("unexpected html: %s", text)
	}
}
