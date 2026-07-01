package recorder

import "testing"

func TestPageOrigin(t *testing.T) {
	tests := []struct {
		raw    string
		origin string
	}{
		{"https://example.com/path?q=1", "https://example.com"},
		{"http://localhost:8080/", "http://localhost:8080"},
		{"about:blank", ""},
		{"", ""},
	}
	for _, tc := range tests {
		if got := pageOrigin(tc.raw); got != tc.origin {
			t.Fatalf("pageOrigin(%q) = %q, want %q", tc.raw, got, tc.origin)
		}
	}
}
