package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/runstatus"
)

func TestWriteHTML(t *testing.T) {
	tmp := t.TempDir()
	path := tmp + "/report.html"
	plan := player.ExecutionPlan{
		Cases: []player.RunCase{{
			FeaturePath: "demo.feature",
			Name:        "S1",
			Tags:        []string{"@smoke"},
			Steps: []gherkin.Step{{
				Keyword: "Когда",
				Text:    `открыт "https://example.com"`,
				Line:    3,
			}},
		}},
	}
	result := player.ExecutionResult{
		Mode:      "dry-run",
		Files:     1,
		Scenarios: 1,
		Steps:     1,
		ScenarioResults: []player.ScenarioResult{
			{FeaturePath: "demo.feature", Scenario: "S1", Status: "dry-run", Message: "dry-run"},
		},
	}
	if err := WriteHTML(path, result, HTMLOptions{Plan: plan}); err != nil {
		t.Fatalf("WriteHTML failed: %v", err)
	}
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read html: %v", err)
	}
	text := string(payload)
	for _, want := range []string{
		"scenaria-report-data",
		"demo.feature",
		"Timeline",
		"Inspector",
		`example.com`,
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %q in html", want)
		}
	}
}

func TestBuildHTMLStepsFromRecords(t *testing.T) {
	sr := player.ScenarioResult{
		Status: "failed",
		StepRecords: []player.StepRecord{{
			Index: 0, Text: "click", Selector: "#btn", Status: "passed", DurationMS: 10,
			IterationPath: []player.IterationFrame{{Kind: "repeat", Index: 2}},
			RetryAttempts: 3,
		}, {
			Index: 1, Text: "assert", Status: "failed", Error: "boom", DurationMS: 5,
		}},
	}
	steps := buildHTMLSteps(sr, nil, nil, false, true, "", "test")
	if len(steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(steps))
	}
	if len(steps[0].Tips) == 0 {
		t.Fatal("expected selector tip for non-testid selector")
	}
	if steps[0].IterationPath != "repeat[2]" {
		t.Fatalf("unexpected iteration path: %q", steps[0].IterationPath)
	}
	if steps[0].RetryAttempts != 3 {
		t.Fatalf("unexpected retry attempts: %d", steps[0].RetryAttempts)
	}
}

func TestSummaryCountsMatchHTMLPayload(t *testing.T) {
	result := player.ExecutionResult{
		Mode:      "browser",
		Files:     2,
		Scenarios: 6,
		Steps:     7,
		ScenarioResults: []player.ScenarioResult{
			{FeaturePath: "a.feature", Scenario: "passed", Status: "passed"},
			{FeaturePath: "b.feature", Scenario: "failed", Status: "failed"},
			{FeaturePath: "c.feature", Scenario: "skipped", Status: "skipped"},
			{FeaturePath: "d.feature", Scenario: "canceled", Status: "canceled"},
			{FeaturePath: "e.feature", Scenario: "not started", Status: "not-started"},
			{FeaturePath: "f.feature", Scenario: "dry", Status: "dry-run"},
		},
	}
	detailed := FromExecutionResultDetailed(result)
	payload, err := buildHTMLPayload(result, HTMLOptions{}, filepath.Join(t.TempDir(), "report.html"))
	if err != nil {
		t.Fatalf("buildHTMLPayload: %v", err)
	}
	if detailed.Passed != payload.Summary.Passed || detailed.Failed != payload.Summary.Failed ||
		detailed.Skipped != payload.Summary.Skipped || detailed.Canceled != payload.Summary.Canceled ||
		detailed.NotStarted != payload.Summary.NotStarted {
		t.Fatalf("summary mismatch: detailed=%+v html=%+v", detailed, payload.Summary)
	}
	if got := statusTotalsSum(StatusTotals{
		Passed:     detailed.Passed,
		Failed:     detailed.Failed,
		Skipped:    detailed.Skipped,
		Canceled:   detailed.Canceled,
		NotStarted: detailed.NotStarted,
	}); got != detailed.Scenarios {
		t.Fatalf("scenario invariant failed: got %d want %d", got, detailed.Scenarios)
	}
}

func TestHTMLPayloadSeparatesCurrentStatusAndHistory(t *testing.T) {
	tmp := t.TempDir()
	store, err := runstatus.Open(tmp)
	if err != nil {
		t.Fatalf("open runstatus: %v", err)
	}
	pathKey := "demo.feature::Scenario"
	if err := store.RecordBatch([]runstatus.Entry{{
		Path:    pathKey,
		Status:  "failed",
		Success: false,
		Message: "first failure",
		At:      "2024-01-01T00:00:00Z",
	}, {
		Path:    pathKey,
		Status:  "passed",
		Success: true,
		At:      "2024-01-02T00:00:00Z",
	}}); err != nil {
		t.Fatalf("record history: %v", err)
	}
	result := player.ExecutionResult{
		Mode:      "browser",
		Files:     1,
		Scenarios: 1,
		Steps:     1,
		ScenarioResults: []player.ScenarioResult{{
			FeaturePath: "demo.feature",
			Scenario:    "Scenario",
			Status:      "failed",
		}},
	}
	payload, err := buildHTMLPayload(result, HTMLOptions{ProjectRoot: tmp}, filepath.Join(tmp, "report.html"))
	if err != nil {
		t.Fatalf("buildHTMLPayload: %v", err)
	}
	if len(payload.Scenarios) != 1 {
		t.Fatalf("expected one scenario, got %d", len(payload.Scenarios))
	}
	sc := payload.Scenarios[0]
	if sc.Status != "failed" {
		t.Fatalf("current status lost: %q", sc.Status)
	}
	if sc.History == nil || sc.History.LastStatus != "passed" || !sc.History.Changed {
		t.Fatalf("history not separated from current status: %+v", sc.History)
	}
	if len(sc.HistoryRuns) == 0 || sc.HistoryRuns[0].Status != "passed" {
		t.Fatalf("history runs missing normalized status: %+v", sc.HistoryRuns)
	}
}

func TestHistoryStatusFromEntryUsesRecordedStatus(t *testing.T) {
	tests := []struct {
		name  string
		entry runstatus.Entry
		want  string
	}{
		{name: "canceled", entry: runstatus.Entry{Status: "canceled", Success: true}, want: "canceled"},
		{name: "not started", entry: runstatus.Entry{Status: "not-started", Success: true}, want: "not-started"},
		{name: "dry run", entry: runstatus.Entry{Status: "dry-run", Success: true}, want: "skipped"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := historyStatusFromEntry(tt.entry); got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}

func TestWriteHTMLStoresScreenshotsAsFiles(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "report.html")
	png := []byte{0x89, 0x50, 0x4e, 0x47}
	if err := WriteHTML(path, player.ExecutionResult{
		Mode: "browser",
		ScenarioResults: []player.ScenarioResult{{
			FeaturePath:   "a.feature",
			Scenario:      "S",
			Status:        "failed",
			ScreenshotPNG: png,
			StepRecords: []player.StepRecord{{
				Index: 0, Text: "click", Status: "failed", ScreenshotPNG: png,
			}},
		}},
	}, HTMLOptions{}); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := string(raw)
	if strings.Contains(text, "data:image/") {
		t.Fatal("screenshot was embedded as data URL")
	}
	if !strings.Contains(text, "screenshots/000_a__S__scenario.png") {
		t.Fatal("scenario screenshot path not embedded")
	}
	if _, err := os.Stat(filepath.Join(tmp, "screenshots", "000_a__S__scenario.png")); err != nil {
		t.Fatalf("scenario screenshot not written: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "screenshots", "000_a__S__step_000.png")); err != nil {
		t.Fatalf("step screenshot not written: %v", err)
	}
}

func TestWriteHTMLRefusesUnmarkedArtifactDir(t *testing.T) {
	tmp := t.TempDir()
	screenshots := filepath.Join(tmp, "screenshots")
	if err := os.MkdirAll(screenshots, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(screenshots, "foreign.png"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := WriteHTML(filepath.Join(tmp, "report.html"), player.ExecutionResult{}, HTMLOptions{})
	if err == nil || !strings.Contains(err.Error(), "unmarked html artifact directory") {
		t.Fatalf("expected unmarked artifact dir error, got %v", err)
	}
}

func TestWriteHTMLCleansMarkedArtifacts(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "report.html")
	png := []byte{0x89, 0x50, 0x4e, 0x47}
	result := player.ExecutionResult{ScenarioResults: []player.ScenarioResult{{
		FeaturePath:   "a.feature",
		Scenario:      "S",
		Status:        "failed",
		ScreenshotPNG: png,
	}}}
	if err := WriteHTML(path, result, HTMLOptions{}); err != nil {
		t.Fatal(err)
	}
	stale := filepath.Join(tmp, "screenshots", "stale.png")
	if err := os.WriteFile(stale, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := WriteHTML(path, result, HTMLOptions{}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(stale); !os.IsNotExist(err) {
		t.Fatalf("stale artifact should be removed, stat err=%v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "screenshots", htmlArtifactMarker)); err != nil {
		t.Fatalf("marker missing: %v", err)
	}
}

func TestWriteHTMLWithBridgeURL(t *testing.T) {
	tmp := t.TempDir()
	path := tmp + "/report.html"
	if err := WriteHTML(path, player.ExecutionResult{
		Mode: "dry-run",
		ScenarioResults: []player.ScenarioResult{
			{FeaturePath: "a.feature", Scenario: "S", Status: "dry-run"},
		},
	}, HTMLOptions{BridgeURL: "http://127.0.0.1:19999"}); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "http://127.0.0.1:19999") {
		t.Fatal("bridge url not embedded")
	}
}

func TestHTMLModePairPaths(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		light     bool
		wantFull  string
		wantLight string
	}{
		{name: "full report", path: filepath.Join("out", "report.html"), wantFull: filepath.Join("out", "report.html"), wantLight: filepath.Join("out", "report.light.html")},
		{name: "light report explicit", path: filepath.Join("out", "report.light.html"), light: true, wantFull: filepath.Join("out", "report.html"), wantLight: filepath.Join("out", "report.light.html")},
		{name: "light report generic", path: filepath.Join("out", "report.html"), light: true, wantFull: filepath.Join("out", "report.full.html"), wantLight: filepath.Join("out", "report.html")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			full, light := htmlModePairPaths(tt.path, tt.light)
			if full != tt.wantFull || light != tt.wantLight {
				t.Fatalf("got full=%q light=%q, want full=%q light=%q", full, light, tt.wantFull, tt.wantLight)
			}
		})
	}
}

func TestWriteHTMLModePairAddsCrossLinks(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "report.html")
	result := player.ExecutionResult{
		Mode: "browser", Files: 1, Scenarios: 1, Steps: 1,
		ScenarioResults: []player.ScenarioResult{{
			FeaturePath: "a.feature",
			Scenario:    "S",
			Status:      "passed",
			StepRecords: []player.StepRecord{{Index: 0, Text: "ok", Status: "passed"}},
		}},
	}
	full, light, err := WriteHTMLModePair(path, result, HTMLOptions{Locale: "ru"})
	if err != nil {
		t.Fatal(err)
	}
	fullRaw, err := os.ReadFile(full)
	if err != nil {
		t.Fatal(err)
	}
	lightRaw, err := os.ReadFile(light)
	if err != nil {
		t.Fatal(err)
	}
	fullText := string(fullRaw)
	lightText := string(lightRaw)
	for _, want := range []string{`"current":"full"`, `"light_href":"report.light.html"`, `"light_available":true`} {
		if !strings.Contains(fullText, want) {
			t.Fatalf("full report missing %q", want)
		}
	}
	for _, want := range []string{`"current":"light"`, `"full_href":"report.html"`, `"full_available":true`} {
		if !strings.Contains(lightText, want) {
			t.Fatalf("light report missing %q", want)
		}
	}
}

func TestWriteHTMLModePairWhenLightDefault(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "report.html")
	result := player.ExecutionResult{
		Mode: "browser", Files: 1, Scenarios: 1, Steps: 1,
		ScenarioResults: []player.ScenarioResult{{
			FeaturePath: "a.feature",
			Scenario:    "S",
			Status:      "passed",
			StepRecords: []player.StepRecord{{Index: 0, Text: "ok", Status: "passed"}},
		}},
	}
	full, light, err := WriteHTMLModePair(path, result, HTMLOptions{LightMode: true, Locale: "ru"})
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(light) != "report.html" {
		t.Fatalf("light path = %q", light)
	}
	if filepath.Base(full) != "report.full.html" {
		t.Fatalf("full path = %q", full)
	}
	lightRaw, err := os.ReadFile(light)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(lightRaw), `"full_href":"report.full.html"`) {
		t.Fatal("light report should link to full sibling")
	}
	if !strings.Contains(string(lightRaw), `"full_available":true`) {
		t.Fatal("full sibling should be marked available")
	}
}

func TestLightReportEmbedsFailedStepScreenshot(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "report.light.html")
	png := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}
	failedStep := 1
	result := player.ExecutionResult{
		Mode: "browser", Files: 1, Scenarios: 1, Steps: 2,
		ScenarioResults: []player.ScenarioResult{{
			FeaturePath: "a.feature",
			Scenario:    "S",
			Status:      "failed",
			FailedStep:  &failedStep,
			StepRecords: []player.StepRecord{
				{Index: 0, Text: "ok", Status: "passed"},
				{Index: 1, Text: "fail", Status: "failed", ScreenshotPNG: png},
			},
		}},
	}
	if err := WriteHTML(path, result, HTMLOptions{LightMode: true}); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := string(raw)
	if !strings.Contains(text, `"screenshot":"screenshots/`) {
		t.Fatalf("light report missing failed step screenshot ref: %s", text)
	}
	if _, err := os.Stat(filepath.Join(tmp, "screenshots", "000_a__S__step_001.png")); err != nil {
		t.Fatalf("failed step screenshot file missing: %v", err)
	}
}

func TestInferStepStatus(t *testing.T) {
	fs := 1
	if inferStepStatus("failed", 0, &fs) != "passed" {
		t.Fatal("step 0 should be passed")
	}
	if inferStepStatus("failed", 1, &fs) != "failed" {
		t.Fatal("step 1 should be failed")
	}
	if inferStepStatus("failed", 2, &fs) != "skipped" {
		t.Fatal("step 2 should be skipped")
	}
	if inferStepStatus("failed", 0, nil) != "failed" {
		t.Fatal("step 0 should be failed when failed_step is missing")
	}
	if inferStepStatus("failed", 1, nil) != "skipped" {
		t.Fatal("step 1 should be skipped when failed_step is missing")
	}
}
