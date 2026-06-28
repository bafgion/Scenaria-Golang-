//go:build integration

package recorder

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRecordLiveIdleWritesFeature(t *testing.T) {
	tmp := t.TempDir()
	out := filepath.Join(tmp, "recorded.feature")
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	err := RecordLive(ctx, LiveOptions{
		StartURL:     "https://example.com",
		OutputPath:   out,
		IdleTimeout:  4 * time.Second,
		Headless:     true,
		FeatureName:  "Example",
		ScenarioName: "Visit",
	})
	if err != nil {
		t.Fatalf("RecordLive: %v", err)
	}
	payload, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read feature: %v", err)
	}
	text := string(payload)
	if !strings.Contains(text, "Example") || !strings.Contains(text, "открыт") {
		t.Fatalf("unexpected feature content:\n%s", text)
	}
}
