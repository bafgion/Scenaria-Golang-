//go:build integration

package recorder

import (
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/paths"
	"github.com/bafgion/scenaria-golang/internal/selector"
	playwright "github.com/mxschmitt/playwright-go"
)

const clickFixtureHTML = `<!doctype html><html><body>
<button id="btn">Click me</button>
</body></html>`

func TestLastClickPreservedOnToolbarStop(t *testing.T) {
	page, session, _ := openToolbarRecorderFixture(t, clickFixtureHTML)
	if err := session.BeginCapture(); err != nil {
		t.Fatal(err)
	}
	if _, err := page.Evaluate(`() => {
		if (!window.__scenariaRecorder) throw new Error('recorder missing');
		window.__scenariaRecorder.events.push({
			type: 'click',
			detail: { selector: '#btn', tag: 'BUTTON' },
			ts: Date.now(),
			seq: 1,
		});
	}`); err != nil {
		t.Fatalf("queue click event: %v", err)
	}

	var events []RecordStepEvent
	CompleteCaptureStop(session, page, func(event RecordStepEvent) {
		events = append(events, event)
	}, "manual", nil)

	foundClick := false
	for _, event := range events {
		if event.Op != RecordStepSnapshot {
			continue
		}
		for _, line := range event.Lines {
			if strings.Contains(line, "#btn") {
				foundClick = true
				break
			}
		}
	}
	if !foundClick {
		t.Fatalf("expected click step in snapshot, got %+v", events)
	}
}

func openToolbarRecorderFixture(t *testing.T, html string) (playwright.Page, *LiveSession, *[]RecordedStep) {
	t.Helper()
	if err := playwright.Install(); err != nil {
		t.Fatalf("install playwright: %v", err)
	}
	paths.ConfigurePlaywrightBrowsers()
	pw, err := playwright.Run()
	if err != nil {
		t.Fatalf("run playwright: %v", err)
	}
	t.Cleanup(func() { _ = pw.Stop() })

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: playwright.Bool(true)})
	if err != nil {
		t.Fatalf("launch: %v", err)
	}
	t.Cleanup(func() { _ = browser.Close() })

	bctx, err := browser.NewContext()
	if err != nil {
		t.Fatalf("context: %v", err)
	}
	t.Cleanup(func() { _ = bctx.Close() })
	if err := bctx.AddInitScript(playwright.Script{Content: playwright.String(selector.HeuristicsJS)}); err != nil {
		t.Fatalf("heuristics init: %v", err)
	}
	if err := bctx.AddInitScript(playwright.Script{Content: playwright.String(selector.RecorderListenersJS)}); err != nil {
		t.Fatalf("recorder init: %v", err)
	}
	if err := bctx.AddInitScript(playwright.Script{Content: playwright.String(selector.BrowserToolbarJS)}); err != nil {
		t.Fatalf("toolbar init: %v", err)
	}

	page, err := bctx.NewPage()
	if err != nil {
		t.Fatalf("page: %v", err)
	}
	startURL := "https://recorder.test/start"
	if err := page.Route(startURL, func(route playwright.Route) {
		_ = route.Fulfill(playwright.RouteFulfillOptions{
			Status:      playwright.Int(200),
			ContentType: playwright.String("text/html"),
			Body:        html,
		})
	}); err != nil {
		t.Fatalf("route: %v", err)
	}
	if _, err := page.Goto(startURL); err != nil {
		t.Fatalf("goto: %v", err)
	}
	recorded := []RecordedStep{{Action: "goto", Value: startURL}}
	session := NewLiveSession()
	session.InitBrowseMode()
	session.Bind(page, &recorded)
	return page, session, &recorded
}
