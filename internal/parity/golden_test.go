package parity_test

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

func TestEmbeddedGoldenFixture(t *testing.T) {
	cases, err := stepdsl.LoadGoldenCases()
	if err != nil {
		t.Fatalf("LoadGoldenCases: %v", err)
	}
	if len(cases) < 40 {
		t.Fatalf("expected at least 40 golden cases for cross-language parity, got %d", len(cases))
	}
	for i, tc := range cases {
		action, err := stepdsl.Parse(gherkin.Step{Line: i + 1, Text: tc.Text})
		if err != nil {
			t.Fatalf("case %d parse %q: %v", i, tc.Text, err)
		}
		if action.Kind != tc.Kind {
			t.Fatalf("case %d kind: got %q want %q", i, action.Kind, tc.Kind)
		}
	}
}
