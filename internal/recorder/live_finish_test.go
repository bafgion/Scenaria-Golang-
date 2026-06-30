package recorder

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPersistLiveRecordingOnCanceledStop(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "recorded.feature")
	session := NewLiveSession()
	session.InitRecordMode()

	recorded := []RecordedStep{
		{Action: "goto", Value: "https://example.com"},
		{Action: "click", Selector: `button:has-text("Войти")`},
	}
	opts := LiveOptions{
		OutputPath:   out,
		FeatureName:  "Записанный сценарий",
		ScenarioName: "Запись",
	}
	if err := persistLiveRecording(opts, session, recorded, context.Canceled); err != nil {
		t.Fatalf("persistLiveRecording: %v", err)
	}
	body, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read feature: %v", err)
	}
	text := string(body)
	if !strings.Contains(text, "https://example.com") {
		t.Fatalf("missing goto step:\n%s", text)
	}
	if !strings.Contains(text, "нажимаю") {
		t.Fatalf("missing click step:\n%s", text)
	}
}

func TestPersistLiveRecordingSkipsBrowseWithoutCapture(t *testing.T) {
	session := NewLiveSession()
	session.InitBrowseMode()
	err := persistLiveRecording(LiveOptions{OutputPath: "unused.feature"}, session, nil, context.Canceled)
	if err == nil || !strings.Contains(err.Error(), "canceled") {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
