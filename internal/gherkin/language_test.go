package gherkin

import "testing"

func TestParseLanguageTag(t *testing.T) {
	if got := ParseLanguageTag("# language: en\nFeature: Demo"); got != LangEN {
		t.Fatalf("expected en, got %q", got)
	}
	if got := ParseLanguageTag("# language: ru\nФункционал: Demo"); got != LangRU {
		t.Fatalf("expected ru, got %q", got)
	}
	if got := ParseLanguageTag("Функционал: Demo"); got != LangRU {
		t.Fatalf("default ru, got %q", got)
	}
}

func TestParseFeature_English(t *testing.T) {
	content := `# language: en
@smoke
Feature: Login

Scenario: Successful sign-in
  When I open "https://example.com"
  Then I see "h1"
`
	feature, err := ParseFeature(content)
	if err != nil {
		t.Fatalf("ParseFeature: %v", err)
	}
	if feature.Language != LangEN {
		t.Fatalf("language: %q", feature.Language)
	}
	if feature.Title != "Login" {
		t.Fatalf("title: %q", feature.Title)
	}
	if len(feature.Scenarios) != 1 || len(feature.Scenarios[0].Steps) != 2 {
		t.Fatalf("scenarios: %+v", feature.Scenarios)
	}
	if feature.Scenarios[0].Steps[0].Keyword != "When" {
		t.Fatalf("keyword: %q", feature.Scenarios[0].Steps[0].Keyword)
	}
}

func TestSerializeFeature_PreservesLanguage(t *testing.T) {
	feature := &Feature{
		Language: LangEN,
		Title:    "Demo",
		Scenarios: []Scenario{{
			Title: "One",
			Steps: []Step{{Keyword: "When", Text: `I open "https://example.com"`, Line: 3}},
		}},
	}
	text, err := SerializeFeature(feature)
	if err != nil {
		t.Fatal(err)
	}
	if !containsAll(text, "# language: en", "Feature: Demo", "Scenario: One", "When I open") {
		t.Fatalf("serialize: %q", text)
	}
}

func containsAll(text string, parts ...string) bool {
	for _, p := range parts {
		if !contains(text, p) {
			return false
		}
	}
	return true
}

func contains(text, sub string) bool {
	return len(sub) == 0 || (len(text) >= len(sub) && (text == sub || len(text) > 0 && findSub(text, sub)))
}

func findSub(text, sub string) bool {
	for i := 0; i+len(sub) <= len(text); i++ {
		if text[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
