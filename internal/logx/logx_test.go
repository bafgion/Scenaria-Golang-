package logx

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestInitJSONHandler(t *testing.T) {
	t.Setenv("SCENARIA_LOG", "json")
	old := logger
	defer func() { logger = old }()
	Init()

	var buf bytes.Buffer
	logger = slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	Info("test event", "key", "value")
	if !strings.Contains(buf.String(), `"msg":"test event"`) {
		t.Fatalf("expected json log, got %q", buf.String())
	}
}

func TestLogHelpersRedactSecrets(t *testing.T) {
	var buf bytes.Buffer
	old := logger
	defer func() { logger = old }()
	logger = slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	Info("opening https://user:pass@example.com/path with token=abc", "password", "secret", "safe", "visible")

	got := buf.String()
	for _, leaked := range []string{"user:pass@", "token=abc", "secret"} {
		if strings.Contains(got, leaked) {
			t.Fatalf("log leaked %q in %q", leaked, got)
		}
	}
	if !strings.Contains(got, "safe=visible") {
		t.Fatalf("expected safe field to remain visible, got %q", got)
	}
	if !strings.Contains(got, "[REDACTED]") {
		t.Fatalf("expected redaction marker, got %q", got)
	}
}
