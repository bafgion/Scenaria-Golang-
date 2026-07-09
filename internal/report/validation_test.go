package report

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/selector"
)

func TestApplyStepValidations(t *testing.T) {
	sc := htmlScenario{
		Steps: []htmlStep{
			{Index: 0, Line: 4, Selector: "#login"},
			{Index: 1, Line: 8, Selector: "#next"},
		},
	}
	byLine := map[int]htmlStepValidation{
		4: {
			Status:     "found",
			Message:    "элемент найден",
			Mode:       "static",
			ActionKind: "click",
			MatchCount: 1,
			Limitation: "проверка только на текущей странице",
		},
	}
	applyStepValidations(&sc, byLine)
	if sc.Steps[0].Validation == nil || sc.Steps[0].Validation.Status != "found" {
		t.Fatalf("step 0 validation: %+v", sc.Steps[0].Validation)
	}
	if sc.Steps[1].Validation != nil {
		t.Fatalf("step 1 should not have validation")
	}
	if sc.ValidationMode != "static" || sc.ValidationLimitation == "" {
		t.Fatalf("scenario banner: mode=%q limitation=%q", sc.ValidationMode, sc.ValidationLimitation)
	}
}

func TestStepValidationFromSelector(t *testing.T) {
	got := stepValidationFromSelector(selector.StepValidation{
		Line:       3,
		Status:     "warning",
		Message:    "not actionable",
		Mode:       "flow",
		ActionKind: "fill",
		MatchCount: 2,
		Limitation: "flow limitation",
	})
	if got.ActionKind != "fill" || got.MatchCount != 2 || got.Mode != "flow" {
		t.Fatalf("unexpected mapping: %+v", got)
	}
}

func TestPathsMatchFeature(t *testing.T) {
	if !pathsMatchFeature(`C:\proj\auth.feature`, `auth.feature`) {
		t.Fatal("expected basename match")
	}
	if pathsMatchFeature(`a.feature`, `b.feature`) {
		t.Fatal("expected different basenames to not match")
	}
}
