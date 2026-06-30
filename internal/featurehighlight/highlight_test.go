package featurehighlight

import (
	"strings"
	"testing"
)

func joinKinds(spans []Span) string {
	parts := make([]string, 0, len(spans))
	for _, span := range spans {
		parts = append(parts, span.Text)
	}
	return strings.Join(parts, "")
}

func TestHighlightPreservesText(t *testing.T) {
	input := "Функционал: Demo\n  @smoke\n  Когда нажимаю \"#ok\""
	spans := Highlight(input)
	if joinKinds(spans) != input {
		t.Fatalf("highlight changed text:\n%q\n%q", input, joinKinds(spans))
	}
}

func TestHighlightTags(t *testing.T) {
	spans := Highlight("@api @smoke")
	if len(spans) < 2 || spans[0].Kind != KindTag {
		t.Fatalf("expected tag spans, got %#v", spans)
	}
}

func TestHighlightInvalidStep(t *testing.T) {
	spans := Highlight("  Когда это точно не шаг")
	found := false
	for _, span := range spans {
		if span.Kind == KindError {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected error kind for invalid step")
	}
}

func TestHighlightTestClient(t *testing.T) {
	spans := Highlight(`я подключаю TestClient "Demo"`)
	found := false
	for _, span := range spans {
		if span.Text == "TestClient" && span.Kind == KindTestClient {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected TestClient highlight, got %#v", spans)
	}
}

func TestHighlightBlockKeyword(t *testing.T) {
	spans := Highlight("\tЕсли вижу \"a\"")
	found := false
	for _, span := range spans {
		if strings.Contains(span.Text, "Если") && span.Kind == KindBlockKeyword {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected block keyword highlight, got %#v", spans)
	}
}
