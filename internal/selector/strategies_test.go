package selector

import "testing"

func TestNormalizeClickStrategiesPreservesOrder(t *testing.T) {
	got := NormalizeClickStrategies([]string{"text", "id", "testid"})
	want := []string{"text", "id", "testid", "aria", "contextual"}
	if len(got) != len(want) {
		t.Fatalf("got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d: got %q want %q", i, got[i], want[i])
		}
	}
}
