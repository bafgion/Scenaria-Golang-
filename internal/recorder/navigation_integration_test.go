//go:build integration

package recorder

import (
	"strings"
	"testing"
	"time"

	"github.com/bafgion/scenaria-golang/internal/paths"
	"github.com/bafgion/scenaria-golang/internal/selector"
	playwright "github.com/mxschmitt/playwright-go"
)

const navLinkHTML = `<!doctype html><html><body>
<a id="go" href="https://recorder.test/target">Go</a>
</body></html>`

const navTargetHTML = `<!doctype html><html><body><h1>Target</h1></body></html>`

const spaMenuHTML = `<!doctype html><html><body>
<nav><a id="clothes" href="/catalog/clothes">Одежда</a></nav>
<script>
document.getElementById('clothes').addEventListener('click', (event) => {
	event.preventDefault();
	history.pushState({}, '', '/catalog/clothes');
});
</script>
</body></html>`

const productCardRecorderHTML = `<!doctype html><html><body>
<article class="product-card">
  <a id="dzhinsy_relaxed_rl200_iz_liotsella_svetlo_zheltogo_tsveta" href="/collection/katalog/dzhinsy_relaxed_rl200_iz_liotsella_svetlo_zheltogo_tsveta/">
    <span>Джинсы свободные RL20015</span>
    <span>980 ₽</span>
    <span>+3</span>
  </a>
  <button type="button">Джинсы свободные RL20015 980 ₽+3</button>
</article>
</body></html>`

func TestRecorderStashPreservesClickAcrossFullNavigation(t *testing.T) {
	page := openRecorderFixture(t, navLinkHTML, "https://recorder.test/target", navTargetHTML)
	clickAndWaitNavigation(t, page, "#go", "https://recorder.test/target")

	events, err := drainRecorderEvents(page)
	if err != nil {
		t.Fatalf("drain: %v", err)
	}
	if !eventsContainClick(events, "#go") {
		t.Fatalf("expected click in drained events after navigation, got %+v", events)
	}

	recorded := []RecordedStep{{Action: "goto", Value: "about:blank"}}
	state := newRecorderPollState(recorded[0].Value, true)
	applyRecorderPollBatch(&recorded, &state, events, page.URL(), time.Now(), nil)
	if len(recorded) < 2 || recorded[1].Action != "click" {
		t.Fatalf("expected click step first, got %+v", recorded)
	}
}

func TestRecorderKeepsNavClickDrainableForSPANavigation(t *testing.T) {
	page := openRecorderFixture(t, spaMenuHTML, "https://recorder.test/catalog/clothes", navTargetHTML)
	if err := page.Click("#clothes"); err != nil {
		t.Fatalf("click menu: %v", err)
	}
	events, err := drainRecorderEvents(page)
	if err != nil {
		t.Fatalf("drain: %v", err)
	}
	if !eventsContainClick(events, "#clothes") {
		t.Fatalf("expected SPA menu click in live drained events, got %+v", events)
	}
}

func TestRecorderProductCardButtonRecordsStableProductLink(t *testing.T) {
	page := openRecorderFixture(t, productCardRecorderHTML, "https://recorder.test/product", navTargetHTML)
	if err := page.Click("button"); err != nil {
		t.Fatalf("click product button: %v", err)
	}
	events, err := drainRecorderEvents(page)
	if err != nil {
		t.Fatalf("drain: %v", err)
	}
	if !eventsContainClick(events, "#dzhinsy_relaxed_rl200_iz_liotsella_svetlo_zheltogo_tsveta") {
		t.Fatalf("expected stable product link click, got %+v", events)
	}
}

func TestRecorderNavCausingClickStashesBeforeUnload(t *testing.T) {
	page := openRecorderFixture(t, navLinkHTML, "https://recorder.test/target", navTargetHTML)
	clickAndWaitNavigation(t, page, "#go", "https://recorder.test/target")
	stashed, err := page.Evaluate(`() => {
		const raw = sessionStorage.getItem('__scenariaRecorderStash');
		if (raw) return JSON.parse(raw);
		const r = window.__scenariaRecorder;
		return r && r.events ? r.events.slice() : [];
	}`)
	if err != nil {
		t.Fatalf("read stash: %v", err)
	}
	events, err := decodeEvents(stashed)
	if err != nil {
		t.Fatalf("decode stash: %v", err)
	}
	if !eventsContainClick(events, "#go") {
		t.Fatalf("expected stashed click, got %+v", events)
	}
}

func openRecorderFixture(t *testing.T, html, routeURL, routeHTML string) playwright.Page {
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

	page, err := bctx.NewPage()
	if err != nil {
		t.Fatalf("page: %v", err)
	}
	startURL := "https://recorder.test/start"
	for _, spec := range []struct{ pattern, body string }{
		{startURL, html},
		{routeURL, routeHTML},
	} {
		pattern, body := spec.pattern, spec.body
		if err := page.Route(pattern, func(route playwright.Route) {
			_ = route.Fulfill(playwright.RouteFulfillOptions{
				Status:      playwright.Int(200),
				ContentType: playwright.String("text/html"),
				Body:        body,
			})
		}); err != nil {
			t.Fatalf("route %s: %v", pattern, err)
		}
	}
	if _, err := page.Goto(startURL); err != nil {
		t.Fatalf("goto start: %v", err)
	}
	return page
}

func clickAndWaitNavigation(t *testing.T, page playwright.Page, selector, url string) {
	t.Helper()
	_, err := page.ExpectNavigation(func() error {
		return page.Click(selector)
	}, playwright.PageExpectNavigationOptions{
		URL:       url,
		Timeout:   playwright.Float(15000),
		WaitUntil: playwright.WaitUntilStateCommit,
	})
	if err != nil {
		t.Fatalf("navigate via %s: %v", selector, err)
	}
}

func eventsContainClick(events []recorderEvent, wantSelector string) bool {
	for _, event := range events {
		if strings.ToLower(event.Type) != "click" {
			continue
		}
		detail := normalizeDetail(event.Detail)
		if detail["selector"] == wantSelector {
			return true
		}
	}
	return false
}
