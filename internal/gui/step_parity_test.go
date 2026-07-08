package gui

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

func TestIDERuntimeStepParity_BasicGoldenSteps(t *testing.T) {
	cases, err := stepdsl.LoadGoldenCases()
	if err != nil {
		t.Fatalf("LoadGoldenCases: %v", err)
	}
	if len(cases) == 0 {
		t.Fatal("expected golden step cases")
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%02d_%s", i, tc.Kind), func(t *testing.T) {
			step := gherkin.Step{Line: 4, Text: tc.Text}
			if _, err := stepdsl.Parse(step); err != nil {
				t.Fatalf("runner parse failed for %q: %v", tc.Text, err)
			}

			text := parityFeatureText(tc.Text)
			issues := ValidateFeatureContent(text)
			for _, issue := range issues {
				if issue.Line == 4 {
					t.Fatalf("IDE validation reported step issue: %#v", issue)
				}
			}

			rows := ParseEditorSteps(text)
			if len(rows) != 1 {
				t.Fatalf("expected one editor row, got %d: %#v", len(rows), rows)
			}
			if rows[0].Error != "" {
				t.Fatalf("editor steps panel error for %q: %s", tc.Text, rows[0].Error)
			}
			if rows[0].Kind != tc.Kind {
				t.Fatalf("editor kind mismatch: got %q want %q", rows[0].Kind, tc.Kind)
			}
		})
	}
}

func parityFeatureText(stepText string) string {
	return strings.Join([]string{
		"Функционал: Parity",
		"Сценарий: Step",
		"\tКогда " + stepText,
	}, "\n")
}
