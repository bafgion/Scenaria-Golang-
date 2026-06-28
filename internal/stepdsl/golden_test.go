package stepdsl

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

type goldenCase struct {
	Text string `json:"text"`
	Kind string `json:"kind"`
	V1   string `json:"v1"`
	V2   string `json:"v2"`
	Mode string `json:"mode"`
	Int  int    `json:"int"`
}

func TestGoldenStepPatterns(t *testing.T) {
	path := filepath.Join("testdata", "steps_golden.json")
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden file: %v", err)
	}
	var cases []goldenCase
	if err := json.Unmarshal(payload, &cases); err != nil {
		t.Fatalf("decode golden file: %v", err)
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
