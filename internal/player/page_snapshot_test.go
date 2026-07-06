package player

import "testing"

func TestTruncateSnapshot(t *testing.T) {
	got := truncateSnapshot("abcdef", 3)
	if got != "abc\n… (truncated)" {
		t.Fatalf("got %q", got)
	}
}
