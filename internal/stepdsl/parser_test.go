package stepdsl

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestParse(t *testing.T) {
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
			name:   "wait",
			step:   gherkin.Step{Line: 5, Text: `жду "2s"`},
			kind:   "wait",
			value1: "2s",
		},
		{
			name:   "press",
			step:   gherkin.Step{Line: 6, Text: `нажимаю клавишу "Enter"`},
			kind:   "press",
			value1: "Enter",
		},
		{
			name:   "assert-url-contains",
			step:   gherkin.Step{Line: 7, Text: `url содержит "dashboard"`},
			kind:   "assert-url-contains",
			value1: "dashboard",
		},
		{
			name:    "unsupported",
			step:    gherkin.Step{Line: 8, Text: "что-то странное"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, err := Parse(tt.step)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected Parse error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}
			if action.Kind != tt.kind || action.Value1 != tt.value1 || action.Value2 != tt.value2 {
				t.Fatalf("unexpected action: %+v", action)
			}
		})
	}
}

func TestResolveURL(t *testing.T) {
	if got := ResolveURL("https://example.com", "https://base.local"); got != "https://example.com" {
		t.Fatalf("unexpected absolute URL resolution: %q", got)
	}
	if got := ResolveURL("/login", "https://base.local"); got != "https://base.local/login" {
		t.Fatalf("unexpected rooted URL resolution: %q", got)
	}
	if got := ResolveURL("profile", "https://base.local/"); got != "https://base.local/profile" {
		t.Fatalf("unexpected relative URL resolution: %q", got)
	}
}
