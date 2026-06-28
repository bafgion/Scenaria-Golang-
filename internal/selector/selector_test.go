package selector

import "testing"

func TestValidateSyntax(t *testing.T) {
	if err := ValidateSyntax("#login"); err != nil {
		t.Fatalf("expected valid selector: %v", err)
	}
	if err := ValidateSyntax(""); err == nil {
		t.Fatal("expected empty selector error")
	}
}
