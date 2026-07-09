package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/player"
)

func TestBuildHTMLPayloadGolden(t *testing.T) {
	plan := player.ExecutionPlan{Cases: []player.RunCase{{
		FeaturePath:  "auth.feature",
		Name:         "Login",
		Tags:         []string{"@smoke"},
		ExampleIndex: 2,
		Steps:        []gherkin.Step{{Keyword: "Когда", Text: `кликаю "#login"`, Line: 4}},
	}}}
	fs := 0
	result := player.ExecutionResult{
		Mode: "browser", Files: 1, Scenarios: 1, Steps: 1,
		ScenarioResults: []player.ScenarioResult{{
			FeaturePath: "auth.feature", Scenario: "Login", Status: "failed",
			Message: "not visible", FailedStep: &fs, DurationMS: 1200,
			StepRecords: []player.StepRecord{{
				Index: 0, Line: 4, Keyword: "Когда", Text: `кликаю "#login"`,
				Selector: "#login", Status: "failed", DurationMS: 400, Error: "not visible",
				Network: "GET https://app.test/api — net::ERR_FAILED",
			}},
		}},
	}
	payload, err := buildHTMLPayload(result, HTMLOptions{
		Plan:               plan,
		LightMode:          true,
		SkipStepValidation: true,
		StepValidations: map[string]map[int]htmlStepValidation{
			"auth.feature": {
				4: {
					Status:     "found",
					Message:    "элемент найден",
					Mode:       "static",
					ActionKind: "click",
					MatchCount: 1,
					Limitation: "проверка только на текущей странице",
				},
			},
		},
	}, t.TempDir()+"/report.html")
	if err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	text := string(raw)
	for _, want := range []string{
		`"version":"1"`,
		`"example_index":2`,
		`"duration_ms":1200`,
		`"network":"GET https://app.test/api`,
		`"@smoke"`,
		`"validation_mode":"static"`,
		`"action_kind":"click"`,
		`"match_count":1`,
		`"limitation":"проверка только на текущей странице"`,
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("payload missing %q:\n%s", want, text)
		}
	}
}
