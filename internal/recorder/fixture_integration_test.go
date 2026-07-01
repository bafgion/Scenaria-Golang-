//go:build integration

package recorder

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	playwright "github.com/mxschmitt/playwright-go"
)

func TestFixtureCheckboxRecordingPipeline(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<!DOCTYPE html><html><body>
<label>Согласен<input type="checkbox" name="agree"></label>
</body></html>`))
	}))
	defer srv.Close()

	recorded := runFixtureRecording(t, srv.URL, `label:has-text("Согласен")`)
	out := NormalizeSteps(recorded)
	if len(out) == 0 {
		t.Fatalf("no steps recorded: %+v", recorded)
	}
	last := out[len(out)-1]
	if last.Action != "check" && last.Action != "click" {
		t.Fatalf("expected check/click on checkbox, got %+v", out)
	}
}

func TestFixtureNestedLabelFillPipeline(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<!DOCTYPE html><html><body>
<label>Имя<input type="text" name="name"></label>
</body></html>`))
	}))
	defer srv.Close()

	recorded := runFixtureRecording(t, srv.URL, `label:has-text("Имя") input`)
	out := NormalizeSteps(recorded)
	fillFound := false
	for _, step := range out {
		if step.Action != "fill" {
			continue
		}
		fillFound = true
		if !strings.HasPrefix(step.Selector, `label:has-text("Имя`) {
			t.Fatalf("expected label selector, got %q", step.Selector)
		}
	}
	if !fillFound {
		t.Fatalf("expected fill step, got %+v", out)
	}
}

func runFixtureRecording(t *testing.T, startURL, clickSelector string) []RecordedStep {
	t.Helper()
	if err := playwright.Install(); err != nil {
		t.Fatalf("install playwright: %v", err)
	}
	pw, err := playwright.Run()
	if err != nil {
		t.Fatalf("run playwright: %v", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		t.Fatalf("launch: %v", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		t.Fatalf("new page: %v", err)
	}
	if _, err := page.Goto(startURL); err != nil {
		t.Fatalf("goto: %v", err)
	}
	if _, err := page.Evaluate(RecorderScript); err != nil {
		t.Fatalf("inject recorder: %v", err)
	}
	if err := page.Click(clickSelector); err != nil {
		t.Fatalf("click %q: %v", clickSelector, err)
	}
	if clickSelector == `label:has-text("Имя") input` {
		if err := page.Fill(`label:has-text("Имя") input`, "Ivan"); err != nil {
			t.Fatalf("fill: %v", err)
		}
	}
	time.Sleep(600 * time.Millisecond)

	raw, err := page.Evaluate(`() => {
		const r = window.__scenariaRecorder;
		if (!r || !r.events.length) return [];
		return r.events.splice(0, r.events.length);
	}`)
	if err != nil {
		t.Fatalf("read events: %v", err)
	}
	events, err := decodeEvents(raw)
	if err != nil {
		t.Fatalf("decode events: %v", err)
	}
	if len(events) == 0 {
		t.Fatal("no recorder events")
	}

	recorded := make([]RecordedStep, 0, len(events))
	for _, event := range events {
		detail := normalizeDetail(event.Detail)
		step, ok := EventToRecordedStep(event.Type, detail)
		if !ok {
			continue
		}
		appendCoalescedStep(&recorded, step, nil)
	}
	return recorded
}
