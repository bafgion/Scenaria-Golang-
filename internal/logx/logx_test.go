package logx

import (
	"bytes"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestInitJSONHandler(t *testing.T) {
	t.Setenv("SCENARIA_LOG", "json")
	Init()
	var buf bytes.Buffer
	logger = slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	Info("test event", "key", "value")
	if !strings.Contains(buf.String(), `"msg":"test event"`) {
		t.Fatalf("expected json log, got %q", buf.String())
	}
	_ = os.Stderr
}
