package player

import (
	"strings"
	"testing"
)

func TestLocalStorageInitScriptIsOriginScoped(t *testing.T) {
	script, err := localStorageInitScript("https://example.test/app", map[string]string{
		"token": `a"b`,
	})
	if err != nil {
		t.Fatalf("localStorageInitScript: %v", err)
	}
	if !strings.Contains(script, `"https://example.test"`) {
		t.Fatalf("expected origin guard in script: %s", script)
	}
	if !strings.Contains(script, `"token":"a\"b"`) {
		t.Fatalf("expected JSON-encoded localStorage payload: %s", script)
	}
}

func TestSameOrigin(t *testing.T) {
	if !sameOrigin("https://example.test/path", "https://example.test/other") {
		t.Fatal("expected same origin")
	}
	if sameOrigin("https://example.test/path", "http://example.test/path") {
		t.Fatal("expected different schemes to differ")
	}
	if sameOrigin("about:blank", "https://example.test") {
		t.Fatal("about:blank must not match")
	}
}
