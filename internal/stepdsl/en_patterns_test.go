package stepdsl

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestParse_EnglishSteps(t *testing.T) {
	cases := []struct {
		text string
		kind string
	}{
		{`I open "https://example.com"`, "goto"},
		{`I see "h1"`, "assert-visible"},
		{`I see "#submit" is enabled`, "assert-enabled"},
		{`I see "#submit" is disabled`, "assert-disabled"},
		{`I see ".tile" is selected`, "assert-selected"},
		{`I wait until "#pay" is enabled`, "wait-enabled"},
		{`I check text matches regex "Pay.*\\d+" in "#pay"`, "assert-text-regex"},
		{`I click "button.submit"`, "click"},
		{`I type "hello" into "input[name=q]"`, "fill"},
		{`I check text "OK" in ".msg"`, "assert-text"},
	}
	for _, tc := range cases {
		action, err := Parse(gherkin.Step{Line: 1, Text: tc.text})
		if err != nil {
			t.Fatalf("%q: %v", tc.text, err)
		}
		if action.Kind != tc.kind {
			t.Fatalf("%q: got %q", tc.text, action.Kind)
		}
	}
}
