package report

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/runstatus"
)

func TestWriteReportE2EFixture(t *testing.T) {
	root := repoRoot(t)
	outDir := filepath.Join(root, "frontend", "e2e", "fixtures", "report")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(outDir, "sample.html")
	projectRoot := t.TempDir()
	store, err := runstatus.Open(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	fs0 := 0
	if err := store.Record(runstatus.Entry{
		Path: "login.feature::Failed login", Success: true, At: "2026-06-01T09:00:00Z",
		DurationMS: 600, StepDurations: []int{180, 120},
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.Record(runstatus.Entry{
		Path: "login.feature::Failed login", Success: false, At: "2026-06-02T09:00:00Z",
		FailedStep: &fs0, DurationMS: 500, StepDurations: []int{170, 80},
	}); err != nil {
		t.Fatal(err)
	}
	fs := 1
	// Minimal PNG header — enough for fixture screenshot file + lightbox E2E.
	fixturePNG := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}
	plan := player.ExecutionPlan{Cases: []player.RunCase{{
		FeaturePath: "login.feature",
		Name:        "Failed login",
		Tags:        []string{"@smoke"},
		Steps: []gherkin.Step{
			{Keyword: "Когда", Text: `открыт "https://example.com"`, Line: 3},
			{Keyword: "Когда", Text: `кликаю "#missing"`, Line: 4},
		},
	}}}
	result := player.ExecutionResult{
		Mode: "browser", Files: 1, Scenarios: 1, Steps: 2,
		ScenarioResults: []player.ScenarioResult{{
			FeaturePath: "login.feature", Scenario: "Failed login", Status: "failed",
			Message: "element not found", FailedStep: &fs, DurationMS: 800,
			StepRecords: []player.StepRecord{
				{Index: 0, Line: 3, Keyword: "Когда", Text: `открыт "https://example.com"`, Status: "passed", DurationMS: 200},
				{Index: 1, Line: 4, Keyword: "Когда", Text: `кликаю "#missing"`, Selector: "#missing", Status: "failed", DurationMS: 100, Error: "element not found",
					PageContext:   "Example\nhttps://example.com",
					DOMSnapshot:   "<html><body><button id=\"missing\">Login</button></body></html>",
					A11ySnapshot:  "button \"Login\"",
					ScreenshotPNG: fixturePNG,
				},
			},
			TraceZIP: fixtureTraceZIP(),
		}},
	}
	if _, _, err := WriteHTMLModePair(out, result, HTMLOptions{
		Plan: plan, Locale: "ru", ProjectRoot: projectRoot,
		GeneratedAt: "2026-01-01T00:00:00Z",
		ReportDir:   "fixtures/report",
		BridgeURL:   "http://127.0.0.1:19999",
		PreviousSummary: &RunSummaryDetailed{
			GeneratedAt: "2026-06-01T10:00:00Z",
			Items: []ScenarioSummary{
				{Path: "login.feature", Scenario: "Failed login", Status: "passed", DurationMS: 500},
			},
		},
	}); err != nil {
		t.Fatal(err)
	}
	zipBytes := minimalTraceZIPForFixture()
	zipOut := filepath.Join(outDir, "sample-trace.zip")
	if err := os.WriteFile(zipOut, zipBytes, 0o644); err != nil {
		t.Fatal(err)
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func fixtureTraceZIP() []byte {
	return minimalTraceZIPForFixture()
}

func minimalTraceZIPForFixture() []byte {
	lines := []string{
		`{"type":"action","startTime":0,"class":"Frame","method":"goto","params":{"url":"https://example.com"}}`,
		`{"type":"action","startTime":200,"class":"Frame","method":"click","params":{"selector":"#missing"}}`,
	}
	netLines := []string{
		`{"snapshot":{"request":{"url":"https://example.com/missing","method":"GET"},"response":{"status":404},"_monotonicTime":150}}`,
	}
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.Create("trace.trace")
	for _, line := range lines {
		_, _ = w.Write([]byte(line + "\n"))
	}
	nw, _ := zw.Create("trace.network")
	for _, line := range netLines {
		_, _ = nw.Write([]byte(line + "\n"))
	}
	_ = zw.Close()
	return buf.Bytes()
}
