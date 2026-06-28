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
			name:   "goto-opened",
			step:   gherkin.Step{Line: 1, Text: `открыт "https://example.com"`},
			kind:   "goto",
			value1: "https://example.com",
		},
		{
			name:   "goto-navigate",
			step:   gherkin.Step{Line: 2, Text: `открываю "https://example.com"`},
			kind:   "goto",
			value1: "https://example.com",
		},
		{
			name:   "click",
			step:   gherkin.Step{Line: 3, Text: `нажимаю "#login"`},
			kind:   "click",
			value1: "#login",
		},
		{
			name:   "fill",
			step:   gherkin.Step{Line: 4, Text: `ввожу "admin" в "#username"`},
			kind:   "fill",
			value1: "admin",
			value2: "#username",
		},
		{
			name:   "assert-visible",
			step:   gherkin.Step{Line: 5, Text: `вижу "h1"`},
			kind:   "assert-visible",
			value1: "h1",
		},
		{
			name:   "assert-text",
			step:   gherkin.Step{Line: 6, Text: `проверяю текст "Example Domain" в "h1"`},
			kind:   "assert-text",
			value1: "Example Domain",
			value2: "h1",
		},
		{
			name:   "wait-seconds",
			step:   gherkin.Step{Line: 7, Text: `жду 2 сек`},
			kind:   "wait",
			value1: "2000ms",
		},
		{
			name:   "wait-duration",
			step:   gherkin.Step{Line: 8, Text: `жду "2s"`},
			kind:   "wait",
			value1: "2s",
		},
		{
			name:   "press",
			step:   gherkin.Step{Line: 9, Text: `нажимаю клавишу "Enter"`},
			kind:   "press",
			value1: "Enter",
		},
		{
			name:   "assert-url-contains",
			step:   gherkin.Step{Line: 10, Text: `url содержит "dashboard"`},
			kind:   "assert-url-contains",
			value1: "dashboard",
		},
		{
			name:   "close-browser",
			step:   gherkin.Step{Line: 11, Text: `закрываю браузер`},
			kind:   "close-browser",
		},
		{
			name:   "check",
			step:   gherkin.Step{Line: 12, Text: `отмечаю "#agree"`},
			kind:   "check",
			value1: "#agree",
		},
		{
			name:   "select",
			step:   gherkin.Step{Line: 13, Text: `выбираю "ru" в "#lang"`},
			kind:   "select",
			value1: "ru",
			value2: "#lang",
		},
		{
			name:   "remember-url",
			step:   gherkin.Step{Line: 14, Text: `запоминаю url как "page_url"`},
			kind:   "remember-url",
			value1: "page_url",
		},
		{
			name:   "fill-generated",
			step:   gherkin.Step{Line: 15, Text: `ввожу случайный телефон в "#phone"`},
			kind:   "fill-generated",
			value1: "phone",
			value2: "#phone",
		},
		{
			name:   "switch-tab",
			step:   gherkin.Step{Line: 16, Text: `переключаюсь на вкладку 2`},
			kind:   "switch-tab",
		},
		{
			name:    "unsupported",
			step:    gherkin.Step{Line: 17, Text: "что-то странное"},
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

func TestParseFirstExampleScenario(t *testing.T) {
	content := `@smoke
Функционал: Примеры для новичков
Сценарий: Первая проверка страницы
  Допустим открыт "https://example.com"
  Тогда вижу "h1"
  И проверяю текст "Example Domain" в "h1"
  И закрываю браузер`

	feature, err := gherkin.ParseFeature(content)
	if err != nil {
		t.Fatalf("ParseFeature failed: %v", err)
	}
	if len(feature.Scenarios) != 1 {
		t.Fatalf("expected one scenario, got %d", len(feature.Scenarios))
	}

	expectedKinds := []string{"goto", "assert-visible", "assert-text", "close-browser"}
	for i, step := range feature.Scenarios[0].Steps {
		action, err := Parse(step)
		if err != nil {
			t.Fatalf("step %d parse failed: %v", i, err)
		}
		if action.Kind != expectedKinds[i] {
			t.Fatalf("step %d: expected kind %q, got %q", i, expectedKinds[i], action.Kind)
		}
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

func TestParseWaitDuration(t *testing.T) {
	d, err := ParseWaitDuration("2000ms")
	if err != nil {
		t.Fatalf("ParseWaitDuration failed: %v", err)
	}
	if d.Milliseconds() != 2000 {
		t.Fatalf("unexpected duration: %v", d)
	}
}
