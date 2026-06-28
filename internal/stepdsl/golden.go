package stepdsl

import (
	_ "embed"
	"encoding/json"
)

//go:embed testdata/steps_golden.json
var goldenStepsJSON []byte

// GoldenCase is a fixture row shared with Python Scenaria step DSL tests.
type GoldenCase struct {
	Text string `json:"text"`
	Kind string `json:"kind"`
	V1   string `json:"v1"`
	V2   string `json:"v2"`
	Mode string `json:"mode"`
	Int  int    `json:"int"`
}

// LoadGoldenCases decodes the embedded cross-language golden fixture.
func LoadGoldenCases() ([]GoldenCase, error) {
	var cases []GoldenCase
	if err := json.Unmarshal(goldenStepsJSON, &cases); err != nil {
		return nil, err
	}
	return cases, nil
}
