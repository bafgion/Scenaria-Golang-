package report

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/player"
)

func TestWriteJUnit(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "junit.xml")

	result := player.ExecutionResult{
		ScenarioResults: []player.ScenarioResult{
			{FeaturePath: "a.feature", Scenario: "S1", Status: "passed"},
			{FeaturePath: "a.feature", Scenario: "S2", Status: "skipped", Message: "dry-run mode"},
			{FeaturePath: "b.feature", Scenario: "S3", Status: "failed", Message: "assertion failed"},
		},
	}

	if err := WriteJUnit(path, result); err != nil {
		t.Fatalf("WriteJUnit returned error: %v", err)
	}

	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read junit file: %v", err)
	}

	var suite struct {
		XMLName  xml.Name `xml:"testsuite"`
		Tests    int      `xml:"tests,attr"`
		Failures int      `xml:"failures,attr"`
		Skipped  int      `xml:"skipped,attr"`
	}
	if err := xml.Unmarshal(payload, &suite); err != nil {
		t.Fatalf("failed to decode junit xml: %v", err)
	}
	if suite.Tests != 3 || suite.Failures != 1 || suite.Skipped != 1 {
		t.Fatalf("unexpected junit counters: %+v", suite)
	}
}
