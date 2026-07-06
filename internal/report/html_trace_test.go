package report

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/player"
)

func TestScenarioHasTrace(t *testing.T) {
	if scenarioHasTrace(player.ScenarioResult{Status: "passed", TraceZIP: []byte("zip")}, false) {
		t.Fatal("passed scenario should not have trace hints")
	}
	if scenarioHasTrace(player.ScenarioResult{Status: "failed"}, false) {
		t.Fatal("failed without zip should not have trace hints")
	}
	if !scenarioHasTrace(player.ScenarioResult{Status: "failed", TraceZIP: []byte("PK")}, false) {
		t.Fatal("failed with zip should have trace hints")
	}
	if scenarioHasTrace(player.ScenarioResult{Status: "failed", TraceZIP: []byte("PK")}, true) {
		t.Fatal("light mode should skip trace hints")
	}
}
