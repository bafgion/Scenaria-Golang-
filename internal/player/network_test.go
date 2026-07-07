package player

import "testing"

func TestFormatHTTPErrorSnippet(t *testing.T) {
	got := formatHTTPErrorSnippet("GET", "https://api.test/missing", 404)
	want := "GET https://api.test/missing — HTTP 404"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	if formatHTTPErrorSnippet("GET", "https://x", 200) != "" {
		t.Fatal("expected empty for 2xx")
	}
	if formatHTTPErrorSnippet("GET", "", 500) != "" {
		t.Fatal("expected empty without url")
	}
}
