package cli

import "testing"

func TestDedupePaths(t *testing.T) {
	in := []string{`a\b.feature`, `a/b.feature`, `b.feature`}
	got := dedupePaths(in)
	if len(got) != 2 {
		t.Fatalf("expected 2 unique paths, got %d: %v", len(got), got)
	}
}
