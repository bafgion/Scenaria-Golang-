package gherkin

import (
	"strings"
	"testing"
)

func TestNormalizeLegacyHasTextEscapes(t *testing.T) {
	input := `div:has-text(\"Save\")`
	got := NormalizeLegacyHasTextEscapes(input)
	if !strings.Contains(got, `Save`) {
		t.Fatalf("expected Save in normalized selector, got %q", got)
	}
}

func TestNormalizeFeatureTextFancyQuotes(t *testing.T) {
	input := "открываю \u201chttps://example.com\u201d"
	got := NormalizeFeatureText(input)
	if got != `открываю "https://example.com"` {
		t.Fatalf("unexpected normalized text: %q", got)
	}
}
