package selector

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

func TestNormalizeValidationMode(t *testing.T) {
	if got := NormalizeValidationMode("flow"); got != ValidationModeFlow {
		t.Fatalf("got %q", got)
	}
	if got := NormalizeValidationMode(""); got != ValidationModeStatic {
		t.Fatalf("got %q", got)
	}
}

func TestValidationProfile(t *testing.T) {
	if validationProfile("click") != profileClickable {
		t.Fatal("click")
	}
	if validationProfile("fill") != profileEditable {
		t.Fatal("fill")
	}
	if validationProfile("assert-hidden") != profileHidden {
		t.Fatal("hidden")
	}
}

func TestCountPriorFlowSteps(t *testing.T) {
	steps := []browserValidateStep{
		{action: stepdsl.Action{Kind: "goto"}},
		{action: stepdsl.Action{Kind: "click"}},
		{action: stepdsl.Action{Kind: "fill"}},
	}
	if got := countPriorFlowSteps(steps, 2); got != 1 {
		t.Fatalf("got %d want 1", got)
	}
}

func TestIsPriorFlowStep(t *testing.T) {
	if !isPriorFlowStep(stepdsl.Action{Kind: "hover"}) {
		t.Fatal("hover should count")
	}
	if isPriorFlowStep(stepdsl.Action{Kind: "wait"}) {
		t.Fatal("wait should not count")
	}
}
