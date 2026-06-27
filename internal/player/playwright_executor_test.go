package player

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestParseStepAction(t *testing.T) {
	tests := []struct {
		name    string
		step    gherkin.Step
		kind    string
		value1  string
		value2  string
		wantErr bool
	}{
		{
			name:   "goto",
			step:   gherkin.Step{Line: 1, Text: `открываю "https://example.com"`},
			kind:   "goto",
			value1: "https://example.com",
		},
		{
			name:   "click",
			step:   gherkin.Step{Line: 2, Text: `нажимаю "#login"`},
			kind:   "click",
			value1: "#login",
		},
		{
			name:   "fill",
			step:   gherkin.Step{Line: 3, Text: `ввожу "admin" в "#username"`},
			kind:   "fill",
			value1: "admin",
			value2: "#username",
		},
		{
			name:   "assert-text",
			step:   gherkin.Step{Line: 4, Text: `вижу "Панель"`},
			kind:   "assert-text",
			value1: "Панель",
		},
		{
			name:    "unsupported",
			step:    gherkin.Step{Line: 5, Text: "что-то странное"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, err := parseStepAction(tt.step)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected parseStepAction error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("parseStepAction returned error: %v", err)
			}
			if action.Kind != tt.kind || action.Value1 != tt.value1 || action.Value2 != tt.value2 {
				t.Fatalf("unexpected action: %+v", action)
			}
		})
	}
}

func TestResolveURL(t *testing.T) {
	if got := resolveURL("https://example.com", "https://base.local"); got != "https://example.com" {
		t.Fatalf("unexpected absolute URL resolution: %q", got)
	}
	if got := resolveURL("/login", "https://base.local"); got != "https://base.local/login" {
		t.Fatalf("unexpected rooted URL resolution: %q", got)
	}
	if got := resolveURL("profile", "https://base.local/"); got != "https://base.local/profile" {
		t.Fatalf("unexpected relative URL resolution: %q", got)
	}
}

func TestExtractQuotedValues(t *testing.T) {
	values := extractQuotedValues(`ввожу "admin" в "#username"`)
	if len(values) != 2 || values[0] != "admin" || values[1] != "#username" {
		t.Fatalf("unexpected quoted values: %#v", values)
	}
}
