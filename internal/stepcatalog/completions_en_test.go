package stepcatalog

import "testing"

func TestCompletionsForLineLang_EnglishKeyword(t *testing.T) {
	result := CompletionsForLineLang("  Giv", 5, "en")
	if len(result.Items) == 0 {
		t.Fatal("expected keyword suggestions")
	}
	found := false
	for _, item := range result.Items {
		if item.Label == "Given" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("missing Given: %+v", result.Items)
	}
}

func TestCompletionsForLineLang_EnglishHeader(t *testing.T) {
	result := CompletionsForLineLang("Feature", 7, "en")
	if len(result.Items) == 0 {
		t.Fatal("expected header suggestions")
	}
}
