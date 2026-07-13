package secret

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"testing"
)

func TestRedactStringRemovesCommonSecretForms(t *testing.T) {
	input := `password="open sesame" token=abc Authorization: Bearer abc.def Cookie: sid=123 https://user:pass@example.com/path`
	got := RedactString(input)
	for _, leaked := range []string{"open sesame", "abc.def", "sid=123", "user:pass@"} {
		if strings.Contains(got, leaked) {
			t.Fatalf("redacted string leaked %q in %q", leaked, got)
		}
	}
	for _, want := range []string{"password=", "token=", "Authorization:", "Cookie:", "https://example.com/path"} {
		if !strings.Contains(got, want) {
			t.Fatalf("redacted string missing %q in %q", want, got)
		}
	}
}

func TestRedactKeyValuesRedactsSensitiveValueByKey(t *testing.T) {
	args := RedactKeyValues([]any{
		"password", "secret",
		"url", "https://user:pass@example.com/",
		slog.String("authorization", "Bearer token"),
		"err", errors.New("token=abc"),
	})
	got := strings.TrimSpace(fmt.Sprint(args...))
	for _, leaked := range []string{"secret", "user:pass@", "token=abc"} {
		if strings.Contains(got, leaked) {
			t.Fatalf("redacted keyvals leaked %q in %q", leaked, got)
		}
	}
	if !strings.Contains(got, Replacement) {
		t.Fatalf("expected replacement marker in %q", got)
	}
}
