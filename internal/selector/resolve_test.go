package selector

import "testing"

func TestHoverLocatorCandidates(t *testing.T) {
	candidates := HoverLocatorCandidates(`div:has-text("Меню")`)
	if len(candidates) < 3 {
		t.Fatalf("expected multiple candidates, got %v", candidates)
	}
	if candidates[len(candidates)-1] != `div:has-text("Меню")` {
		t.Fatalf("expected original selector last, got %q", candidates[len(candidates)-1])
	}
}

func TestIsChained(t *testing.T) {
	if !IsChained(`div.card >> button.submit`) {
		t.Fatal("expected chained selector")
	}
	if IsChained("#submit") {
		t.Fatal("expected simple selector")
	}
}
