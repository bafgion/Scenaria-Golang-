package recorder

import "testing"

func TestToolbarCommandIDAcceptsBrowserNumberTypes(t *testing.T) {
	if got := toolbarCommandID(float64(42)); got != 42 {
		t.Fatalf("float64 id = %d", got)
	}
	if got := toolbarCommandID(int64(7)); got != 7 {
		t.Fatalf("int64 id = %d", got)
	}
	if got := toolbarCommandID("bad"); got != 0 {
		t.Fatalf("bad id = %d", got)
	}
}
