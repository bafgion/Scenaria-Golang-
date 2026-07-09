package player

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestPoolEligibleForPlanRejectsCloseBrowser(t *testing.T) {
	options := PlaywrightExecutorOptions{}
	plan := ExecutionPlan{Cases: []RunCase{{
		Steps: []gherkin.Step{{Keyword: "И", Text: "закрываю браузер", Line: 1}},
	}}}
	if poolEligibleForPlan(options, plan) {
		t.Fatal("expected pool to be disabled when plan contains close-browser")
	}
	if !poolEligibleForPlan(options, ExecutionPlan{}) {
		t.Fatal("expected empty plan to remain pool-eligible")
	}
}
