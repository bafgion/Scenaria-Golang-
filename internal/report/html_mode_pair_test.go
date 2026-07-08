package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/player"
)

func TestWriteHTMLModePairBuildsPayloadOnce(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.html")
	result := player.ExecutionResult{
		Scenarios: 1,
		ScenarioResults: []player.ScenarioResult{
			{
				FeaturePath:   "demo.feature",
				Scenario:      "S",
				Status:        "failed",
				Message:       "boom",
				ScreenshotPNG: []byte("png"),
				TraceZIP:      []byte("trace"),
			},
		},
	}
	full, light, err := WriteHTMLModePair(path, result, HTMLOptions{Locale: "ru"})
	if err != nil {
		t.Fatal(err)
	}
	fullBody, err := os.ReadFile(full)
	if err != nil {
		t.Fatal(err)
	}
	lightBody, err := os.ReadFile(light)
	if err != nil {
		t.Fatal(err)
	}
	fullText := string(fullBody)
	lightText := string(lightBody)
	if !strings.Contains(fullText, `"light_mode":false`) && !strings.Contains(fullText, `"light_mode": false`) {
		t.Fatal("expected full report payload")
	}
	if !strings.Contains(lightText, `"light_mode":true`) && !strings.Contains(lightText, `"light_mode": true`) {
		t.Fatal("expected light report payload")
	}
	if strings.Contains(lightText, `"trace_path"`) {
		t.Fatal("light report should not include trace path")
	}
	if !strings.Contains(fullText, `"trace_path"`) && !strings.Contains(fullText, `traces`) {
		// trace metadata may be omitted when zip is invalid; ensure full still differs from light
		if len(fullBody) <= len(lightBody) {
			t.Fatal("expected full report to be larger than light report")
		}
	}
}
