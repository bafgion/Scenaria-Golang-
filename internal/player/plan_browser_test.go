package player

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestPlanContainsCloseBrowser(t *testing.T) {
	withClose := ExecutionPlan{Cases: []RunCase{{
		Steps: []gherkin.Step{{Keyword: "И", Text: "закрываю браузер", Line: 1}},
	}}}
	if !PlanContainsCloseBrowser(withClose) {
		t.Fatal("expected close-browser in plan")
	}

	without := ExecutionPlan{Cases: []RunCase{{
		Steps: []gherkin.Step{{Keyword: "Допустим", Text: `открыт "https://example.com"`, Line: 1}},
	}}}
	if PlanContainsCloseBrowser(without) {
		t.Fatal("expected no close-browser in plan")
	}
}

func TestScenarioEndedWithCloseBrowser(t *testing.T) {
	if !ScenarioEndedWithCloseBrowser(ScenarioResult{
		StepRecords: []StepRecord{{TerminalAction: "close-browser"}},
	}) {
		t.Fatal("expected close-browser terminal action")
	}
	if ScenarioEndedWithCloseBrowser(ScenarioResult{}) {
		t.Fatal("expected false without records")
	}
}
