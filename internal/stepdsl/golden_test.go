package stepdsl

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestGoldenStepPatterns(t *testing.T) {
	cases, err := LoadGoldenCases()
	if err != nil {
		t.Fatalf("load golden cases: %v", err)
	}
	for i, tc := range cases {
		action, err := Parse(gherkin.Step{Line: i + 1, Text: tc.Text})
		if err != nil {
			t.Fatalf("case %d parse failed: %v", i, err)
		}
		if action.Kind != tc.Kind || action.Value1 != tc.V1 || action.Value2 != tc.V2 || action.Mode != tc.Mode || action.IntVal != tc.Int {
			t.Fatalf("case %d mismatch: got %+v want kind=%s v1=%s v2=%s mode=%s int=%d", i, action, tc.Kind, tc.V1, tc.V2, tc.Mode, tc.Int)
		}
	}
}
