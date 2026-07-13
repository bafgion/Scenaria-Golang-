//go:build integration

package selector

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/paths"
	playwright "github.com/mxschmitt/playwright-go"
)

func TestBrowserToolbarActionFIFO(t *testing.T) {
	page := openToolbarFixture(t)
	mustEval(t, page, `() => window.__scenariaToolbar.setState({ recording: true, paused: true, browserOnly: false, stepCount: 1 })`)
	mustEval(t, page, `() => {
		window.__scenariaToolbar.__test.enqueueAction('resume');
		window.__scenariaToolbar.__test.enqueueAction('picker');
	}`)
	if got := takeToolbarAction(page); got != "resume" {
		t.Fatalf("expected resume first, got %q", got)
	}
	mustEval(t, page, `() => window.__scenariaToolbar.ackAction(window.__scenariaToolbar.__test.pendingAction()?.id, true)`)
	if got := takeToolbarAction(page); got != "picker" {
		t.Fatalf("expected picker second, got %q", got)
	}
	mustEval(t, page, `() => window.__scenariaToolbar.ackAction(window.__scenariaToolbar.__test.pendingAction()?.id, true)`)
	if got := takeToolbarAction(page); got != "" {
		t.Fatalf("expected empty queue, got %q", got)
	}
}

func TestBrowserToolbarStopDisabledWhenIdle(t *testing.T) {
	page := openToolbarFixture(t)
	mustEval(t, page, `() => window.__scenariaToolbar.setState({ recording: false, paused: false, browserOnly: true, stepCount: 0 })`)
	disabled, err := page.Evaluate(`() => document.querySelector('#scenaria-browser-toolbar button[data-action="stop"]').disabled`)
	if err != nil {
		t.Fatalf("query stop disabled: %v", err)
	}
	if disabled != true {
		t.Fatal("expected stop disabled while not recording")
	}
}

func TestBrowserToolbarStopStaysEnabledWhileQueued(t *testing.T) {
	page := openToolbarFixture(t)
	mustEval(t, page, `() => window.__scenariaToolbar.setState({ recording: true, paused: false, browserOnly: false, stepCount: 2 })`)
	mustEval(t, page, `() => window.__scenariaToolbar.__test.enqueueAction('stop')`)
	disabled, err := page.Evaluate(`() => document.querySelector('#scenaria-browser-toolbar button[data-action="stop"]').disabled`)
	if err != nil {
		t.Fatalf("query stop disabled: %v", err)
	}
	if disabled == true {
		t.Fatal("expected stop to stay enabled while action is queued")
	}
	if got := takeToolbarAction(page); got != "stop" {
		t.Fatalf("expected stop action, got %q", got)
	}
}

func TestBrowserToolbarStopWinsWhenQueueIsFull(t *testing.T) {
	page := openToolbarFixture(t)
	mustEval(t, page, `() => window.__scenariaToolbar.setState({ recording: true, paused: true, browserOnly: false })`)
	mustEval(t, page, `() => {
		for (let i = 1; i <= 8; i++) {
			window.__scenariaToolbar.__test.enqueueAction('synthetic-' + i);
		}
		window.__scenariaToolbar.__test.enqueueAction('stop');
	}`)
	if got := takeToolbarAction(page); got != "stop" {
		t.Fatalf("expected stop to win full queue, got %q", got)
	}
}

func TestBrowserToolbarRejectsNonTerminalOverflowWithoutDroppingOldest(t *testing.T) {
	page := openToolbarFixture(t)
	mustEval(t, page, `() => window.__scenariaToolbar.setState({ recording: true, paused: true, browserOnly: false })`)
	mustEval(t, page, `() => {
		for (let i = 1; i <= 8; i++) {
			window.__scenariaToolbar.__test.enqueueAction('synthetic-' + i);
		}
		window.__scenariaToolbar.__test.enqueueAction('synthetic-9');
	}`)
	if got := takeToolbarAction(page); got != "synthetic-1" {
		t.Fatalf("expected oldest command retained, got %q", got)
	}
}

func TestBrowserToolbarWaitsForAckBeforeNextAction(t *testing.T) {
	page := openToolbarFixture(t)
	mustEval(t, page, `() => window.__scenariaToolbar.setState({ recording: true, paused: true, browserOnly: false })`)
	mustEval(t, page, `() => {
		window.__scenariaToolbar.__test.enqueueAction('resume');
		window.__scenariaToolbar.__test.enqueueAction('picker');
	}`)
	if got := takeToolbarAction(page); got != "resume" {
		t.Fatalf("expected resume first, got %q", got)
	}
	if got := takeToolbarAction(page); got != "" {
		t.Fatalf("expected no second action before ack, got %q", got)
	}
	mustEval(t, page, `() => window.__scenariaToolbar.ackAction(window.__scenariaToolbar.__test.pendingAction()?.id, true)`)
	if got := takeToolbarAction(page); got != "picker" {
		t.Fatalf("expected picker after ack, got %q", got)
	}
}

func TestBrowserToolbarSingleInstanceAfterNavigation(t *testing.T) {
	page := openToolbarFixture(t)
	if _, err := page.Goto("https://toolbar.test/second"); err != nil {
		t.Fatalf("goto second: %v", err)
	}
	count, err := page.Evaluate(`() => document.querySelectorAll('#scenaria-browser-toolbar').length`)
	if err != nil {
		t.Fatalf("count toolbars: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected exactly one toolbar after navigation, got %v", count)
	}
}

func openToolbarFixture(t *testing.T) playwright.Page {
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
	if err := bctx.AddInitScript(playwright.Script{Content: playwright.String(BrowserToolbarJS)}); err != nil {
		t.Fatalf("toolbar init: %v", err)
	}

	page, err := bctx.NewPage()
	if err != nil {
		t.Fatalf("page: %v", err)
	}
	for _, route := range []string{"https://toolbar.test/start", "https://toolbar.test/second"} {
		pattern := route
		if err := page.Route(pattern, func(r playwright.Route) {
			_ = r.Fulfill(playwright.RouteFulfillOptions{
				Status:      playwright.Int(200),
				ContentType: playwright.String("text/html"),
				Body:        `<!doctype html><html><body><h1>toolbar</h1></body></html>`,
			})
		}); err != nil {
			t.Fatalf("route %s: %v", pattern, err)
		}
	}
	if _, err := page.Goto("https://toolbar.test/start"); err != nil {
		t.Fatalf("goto start: %v", err)
	}
	return page
}

func takeToolbarAction(page playwright.Page) string {
	raw, err := page.Evaluate(`() => window.__scenariaToolbar?.takeAction?.() || ''`)
	if err != nil {
		return ""
	}
	if action, ok := raw.(string); ok {
		return action
	}
	if payload, ok := raw.(map[string]any); ok {
		if action, _ := payload["action"].(string); action != "" {
			return action
		}
	}
	return ""
}

func mustEval(t *testing.T, page playwright.Page, script string) {
	t.Helper()
	if _, err := page.Evaluate(script); err != nil {
		t.Fatalf("evaluate: %v", err)
	}
}
