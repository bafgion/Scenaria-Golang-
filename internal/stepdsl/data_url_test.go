package stepdsl

import (
	"strings"
	"testing"
)

func TestNormalizeDataTextHTMLURLCyrillic(t *testing.T) {
	raw := `data:text/html,<html><body><div class=total>Итого: 7 180 ₽</div></body></html>`
	got := NormalizeDataTextHTMLURL(raw)
	if !strings.HasPrefix(strings.ToLower(got), "data:text/html;charset=utf-8,") {
		t.Fatalf("expected charset=utf-8 prefix, got %q", got)
	}
	if strings.Contains(got, "Итого") {
		t.Fatalf("expected percent-encoded payload, still has raw Cyrillic: %q", got)
	}
	if !strings.Contains(got, "%D0%98%D1%82%D0%BE%D0%B3%D0%BE") {
		t.Fatalf("expected UTF-8 percent encoding for Итого, got %q", got)
	}
}

func TestNormalizeDataTextHTMLURLAlreadyEncoded(t *testing.T) {
	raw := "data:text/html;charset=utf-8,%3Chtml%3E%3Cbody%3Eok%3C/body%3E%3C/html%3E"
	if got := NormalizeDataTextHTMLURL(raw); got != raw {
		t.Fatalf("encoded URL changed: %q", got)
	}
}

func TestNormalizeDataTextHTMLURLHTTPUnchanged(t *testing.T) {
	raw := "https://example.com/page"
	if got := NormalizeDataTextHTMLURL(raw); got != raw {
		t.Fatalf("http URL changed: %q", got)
	}
}

func TestResolveURLNormalizesDataHTML(t *testing.T) {
	raw := `data:text/html,<html>тест</html>`
	got := ResolveURL(raw, "https://base.local")
	if strings.Contains(got, "тест") {
		t.Fatalf("ResolveURL should normalize data html: %q", got)
	}
	if !strings.Contains(got, "charset=utf-8") {
		t.Fatalf("ResolveURL missing charset: %q", got)
	}
}
