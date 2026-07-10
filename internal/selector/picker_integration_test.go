//go:build integration

package selector

import (
	"strings"
	"testing"
	"time"

	"github.com/bafgion/scenaria-golang/internal/paths"
	playwright "github.com/mxschmitt/playwright-go"
)

const nestedLabelFormHTML = `<!doctype html><html><body>
<div><div>
<label><div>ИНН</div><div><input type="text" name="inn"></div></label>
<label><div>E-mail</div><div><input type="email" id="email-field"></div></label>
<label for="name-field"><div>Имя</div></label><input type="text" id="name-field">
</div></div></body></html>`

const sameOriginIframeHTML = `<!doctype html><html><body>
<iframe id="inner-frame" style="width:320px;height:80px;border:0" srcdoc="<!doctype html><html><body style='margin:0'><button id='inner-btn' style='width:100%;height:100%'>OK</button></body></html>"></iframe>
</body></html>`

const shadowDomHTML = `<!doctype html><html><body>
<div id="host" style="display:inline-block;padding:12px"></div>
<script>
const host = document.getElementById('host');
const root = host.attachShadow({mode:'open'});
root.innerHTML = '<button id="shadow-btn" style="padding:8px 16px">Shadow</button>';
</script>
</body></html>`

const svgButtonHTML = `<!doctype html><html><body>
<button id="icon-btn"><svg width="20" height="20"><circle cx="10" cy="10" r="8"/></svg></button>
</body></html>`

const iframeWidgetHTML = `<!doctype html><html><body>
<iframe id="tg-login" title="Telegram login" src="https://oauth.telegram.org/embed/demo"
  style="width:240px;height:52px;border:0"></iframe>
</body></html>`

const duplicateCardsHTML = `<!doctype html><html><body>
<section class="card"><h2>Basic</h2><button>Edit</button></section>
<section class="card"><h2>Pro</h2><button>Edit</button></section>
</body></html>`

const duplicateRowsHTML = `<!doctype html><html><body>
<table>
  <tbody>
    <tr><td>Alpha</td><td><button>Edit</button></td></tr>
    <tr><td>Beta</td><td><button>Edit</button></td></tr>
  </tbody>
</table>
</body></html>`

const unconnectedLabelHTML = `<!doctype html><html><body>
<label>Email</label>
<input type="text" />
</body></html>`

func TestPickerNestedLabelInputReturnsControlSelector(t *testing.T) {
	selector := pickAtSelector(t, nestedLabelFormHTML, "", `label:nth-of-type(1) input`)
	if selector == "" {
		t.Fatal("no selector picked")
	}
	if selector != `input[name="inn"]` && (!strings.Contains(selector, `>> input`) || !strings.Contains(selector, "ИНН")) {
		t.Fatalf("got %q", selector)
	}
}

func TestPickerLabelForReturnsControlSelector(t *testing.T) {
	selector := pickAtSelector(t, nestedLabelFormHTML, "", `#name-field`)
	if selector == "" {
		t.Fatal("no selector picked")
	}
	if selector != `#name-field` {
		t.Fatalf("got %q", selector)
	}
}

func TestPickerInputWithIdPreferredOverLabel(t *testing.T) {
	selector := pickAtSelector(t, nestedLabelFormHTML, "", `#email-field`)
	if selector != `#email-field` {
		t.Fatalf("got %q", selector)
	}
}

func TestPickerAdjacentLabelReturnsChainedSelector(t *testing.T) {
	html := `<!doctype html><html><body>
<label>Город</label>
<input type="text">
</body></html>`
	selector := pickAtSelector(t, html, "", `input`)
	if !strings.Contains(selector, `label:has-text("Город")`) || !strings.Contains(selector, `>>`) {
		t.Fatalf("got %q", selector)
	}
}

func TestPickerSameOriginIframeReturnsChainedSelector(t *testing.T) {
	selector := pickAtPoint(t, sameOriginIframeHTML, "", pickInnerIframeButton)
	if !strings.Contains(selector, `#inner-frame`) || !strings.Contains(selector, `>>`) || !strings.Contains(selector, `#inner-btn`) {
		t.Fatalf("got %q", selector)
	}
}

func TestPickerShadowDOMReturnsInnerButton(t *testing.T) {
	selector := pickAtSelector(t, shadowDomHTML, "", `#host`)
	if selector != `#host >> #shadow-btn` && selector != `#shadow-btn` {
		t.Fatalf("got %q", selector)
	}
}

func TestPickerSvgClickNormalizesToButton(t *testing.T) {
	selector := pickAtSelector(t, svgButtonHTML, "", `#icon-btn svg circle`)
	if selector != `#icon-btn` {
		t.Fatalf("got %q", selector)
	}
}

func TestHeuristicsDuplicateTextButtonsWarnAmbiguous(t *testing.T) {
	page := openPickerFixture(t, `<!doctype html><html><body>
<button>Save</button>
<button>Save</button>
</body></html>`, "")
	defer page.Close()

	if _, err := page.Evaluate(RecorderHeuristicsJS); err != nil {
		t.Fatalf("heuristics: %v", err)
	}
	raw, err := page.Evaluate(`() => {
		const result = window.__scenariaHeuristics.buildPickerResult(document.querySelectorAll('button')[1], 'click');
		return { selector: result.selector, warnings: result.warnings, candidates: result.candidates };
	}`)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	payload, ok := raw.(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected payload: %#v", raw)
	}
	if warningsContain(payload["warnings"], "not-unique") {
		return
	}
	if warningsContain(payload["warnings"], "low-confidence") {
		return
	}
	t.Fatalf("expected ambiguity warning, got %#v", payload)
}

func TestRecorderCollectDuplicateTextButtonUsesUniqueSelector(t *testing.T) {
	page := openPickerFixture(t, `<!doctype html><html><body>
<button>Save</button>
<button>Save</button>
</body></html>`, "")
	defer page.Close()

	if _, err := page.Evaluate(RecorderHeuristicsJS); err != nil {
		t.Fatalf("heuristics: %v", err)
	}
	raw, err := page.Evaluate(`() => window.__scenariaHeuristics.collect(document.querySelectorAll('button')[1], 'click').selector`)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	selector, ok := raw.(string)
	if !ok || selector == "" {
		t.Fatalf("empty selector: %#v", raw)
	}
	if selector == `button:has-text("Save")` {
		t.Fatalf("recorder kept ambiguous text selector: %q", selector)
	}
	count, err := page.Locator(selector).Count()
	if err != nil {
		t.Fatalf("locator %q: %v", selector, err)
	}
	if count != 1 {
		t.Fatalf("selector %q matched %d elements", selector, count)
	}
}

func TestRecorderCollectRepeatedCardsUsesUniqueButtonSelector(t *testing.T) {
	page := openPickerFixture(t, duplicateCardsHTML, "")
	defer page.Close()

	if _, err := page.Evaluate(RecorderHeuristicsJS); err != nil {
		t.Fatalf("heuristics: %v", err)
	}
	raw, err := page.Evaluate(`() => window.__scenariaHeuristics.collect(document.querySelectorAll('button')[1], 'click').selector`)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	selector, ok := raw.(string)
	if !ok || selector == "" {
		t.Fatalf("empty selector: %#v", raw)
	}
	count, err := page.Locator(selector).Count()
	if err != nil {
		t.Fatalf("locator %q: %v", selector, err)
	}
	if count != 1 {
		t.Fatalf("selector %q matched %d elements", selector, count)
	}
	text, err := page.Locator(selector).InnerText()
	if err != nil {
		t.Fatalf("inner text %q: %v", selector, err)
	}
	if strings.TrimSpace(text) != "Edit" {
		t.Fatalf("selector %q resolved to %q", selector, text)
	}
}

func TestRecorderCollectDuplicateCardButtonsPrefersStableContext(t *testing.T) {
	page := openPickerFixture(t, duplicateCardsHTML, "")
	defer page.Close()

	if _, err := page.Evaluate(RecorderHeuristicsJS); err != nil {
		t.Fatalf("heuristics: %v", err)
	}
	raw, err := page.Evaluate(`() => window.__scenariaHeuristics.collect(document.querySelectorAll('button')[1], 'click').selector`)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	selector, ok := raw.(string)
	if !ok || selector == "" {
		t.Fatalf("empty selector: %#v", raw)
	}
	if strings.Contains(selector, "nth-") {
		t.Fatalf("selector should avoid nth fallback: %q", selector)
	}
	if !strings.Contains(selector, `card`) && !strings.Contains(selector, `:has-text("Pro")`) {
		t.Fatalf("selector should use stable card context, got %q", selector)
	}
	count, err := page.Locator(selector).Count()
	if err != nil {
		t.Fatalf("locator %q: %v", selector, err)
	}
	if count != 1 {
		t.Fatalf("selector %q matched %d elements", selector, count)
	}
}

func TestRecorderCollectDuplicateRowButtonsPrefersRowContext(t *testing.T) {
	page := openPickerFixture(t, duplicateRowsHTML, "")
	defer page.Close()

	if _, err := page.Evaluate(RecorderHeuristicsJS); err != nil {
		t.Fatalf("heuristics: %v", err)
	}
	raw, err := page.Evaluate(`() => window.__scenariaHeuristics.collect(document.querySelectorAll('button')[1], 'click').selector`)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	selector, ok := raw.(string)
	if !ok || selector == "" {
		t.Fatalf("empty selector: %#v", raw)
	}
	if strings.Contains(selector, "nth-") {
		t.Fatalf("selector should avoid nth fallback: %q", selector)
	}
	if !strings.Contains(selector, `tr`) && !strings.Contains(selector, `[role="row"]`) {
		t.Fatalf("selector should use row context, got %q", selector)
	}
	count, err := page.Locator(selector).Count()
	if err != nil {
		t.Fatalf("locator %q: %v", selector, err)
	}
	if count != 1 {
		t.Fatalf("selector %q matched %d elements", selector, count)
	}
}

func TestRecorderCollectRepeatedCardsKeepsStableSelectorAfterOrderChange(t *testing.T) {
	page := openPickerFixture(t, duplicateCardsHTML, "")
	defer page.Close()

	if _, err := page.Evaluate(RecorderHeuristicsJS); err != nil {
		t.Fatalf("heuristics: %v", err)
	}
	raw, err := page.Evaluate(`() => window.__scenariaHeuristics.collect(document.querySelectorAll('button')[1], 'click').selector`)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	selector, ok := raw.(string)
	if !ok || selector == "" {
		t.Fatalf("empty selector: %#v", raw)
	}

	reordered := `<!doctype html><html><body>
<section class="card"><h2>Pro</h2><button>Edit</button></section>
<section class="card"><h2>Basic</h2><button>Edit</button></section>
</body></html>`
	page2 := openPickerFixture(t, reordered, "")
	defer page2.Close()
	if _, err := page2.Evaluate(RecorderHeuristicsJS); err != nil {
		t.Fatalf("heuristics: %v", err)
	}
	count, err := page2.Locator(selector).Count()
	if err != nil {
		t.Fatalf("locator %q: %v", selector, err)
	}
	if count != 1 {
		t.Fatalf("reordered selector %q matched %d elements", selector, count)
	}
	text, err := page2.Locator(selector).Evaluate(`el => el.closest('section')?.innerText || ''`, nil)
	if err != nil {
		t.Fatalf("closest section text %q: %v", selector, err)
	}
	textValue, _ := text.(string)
	if !strings.Contains(strings.TrimSpace(textValue), "Pro") {
		t.Fatalf("selector %q stopped tracking the same card: %v", selector, text)
	}
}

func TestPickerUnconnectedLabelWarnsLowConfidence(t *testing.T) {
	page := openPickerFixture(t, unconnectedLabelHTML, "")
	defer page.Close()

	if _, err := page.Evaluate(RecorderHeuristicsJS); err != nil {
		t.Fatalf("heuristics: %v", err)
	}
	raw, err := page.Evaluate(`() => window.__scenariaHeuristics.buildPickerResult(document.querySelector('input'), 'input')`)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	payload, ok := raw.(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected payload: %#v", raw)
	}
	if !warningsContain(payload["warnings"], "unconnected-label") && !warningsContain(payload["warnings"], "low-confidence") {
		t.Fatalf("expected unconnected-label warning, got %#v", payload)
	}
}

func TestPickerLabelCaptionDivReturnsInputSelector(t *testing.T) {
	selector := pickAtSelector(t, nestedLabelFormHTML, "", `label:nth-of-type(2) > div:first-child`)
	if selector == "" {
		t.Fatal("no selector picked")
	}
	if selector != `#email-field` {
		t.Fatalf("got %q", selector)
	}
}

func TestPickerShortCaptionReturnsControlSelector(t *testing.T) {
	selector := pickAtSelector(t, nestedLabelFormHTML, "", `#name-field`)
	if selector == "" {
		t.Fatal("no selector picked")
	}
	if selector != `#name-field` {
		t.Fatalf("got %q", selector)
	}
}

func TestPickerCrossOriginIframeReturnsIframeSelector(t *testing.T) {
	selector := pickAtSelector(t, iframeWidgetHTML, "**/embed/**", "#tg-login")
	if selector != `iframe[src*="telegram.org"]` {
		t.Fatalf("got %q", selector)
	}
}

func TestPickerDoesNotInstallHintInsideIframe(t *testing.T) {
	page := openPickerFixture(t, iframeWidgetHTML, "**/embed/**")
	defer page.Close()

	_ = installPickerBindings(t, page)
	if _, err := page.Evaluate(RecorderHeuristicsJS); err != nil {
		t.Fatalf("heuristics: %v", err)
	}
	if _, err := page.Evaluate(PickerInstallScript); err != nil {
		t.Fatalf("picker: %v", err)
	}

	frame := frameByURLContains(page, "telegram.org")
	if frame == nil {
		t.Fatal("iframe frame not found")
	}
	count, err := frame.Evaluate(`() => document.querySelectorAll('#__shopPickerHint').length`)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if n, ok := asNumber(count); !ok || n != 0 {
		t.Fatalf("expected 0 hints in iframe, got %v", count)
	}
}

func TestPickerSameOriginIframeAddsWarning(t *testing.T) {
	payload := pickPayloadAtPoint(t, sameOriginIframeHTML, "", pickInnerIframeButton)
	data, ok := payload.(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected payload: %#v", payload)
	}
	if !warningsContain(data["warnings"], "iframe") {
		t.Fatalf("expected iframe warning, got %#v", data)
	}
}

func pickAtSelector(t *testing.T, html, routePattern, clickSelector string) string {
	t.Helper()
	page := openPickerFixture(t, html, routePattern)
	defer page.Close()
	loc := page.Locator(clickSelector)
	box, err := loc.BoundingBox()
	if err != nil || box == nil {
		t.Fatalf("bounding box for %q: %v", clickSelector, err)
	}
	return pickAtPagePoint(t, page, box.X+box.Width/2, box.Y+box.Height/2)
}

func pickInnerIframeButton(t *testing.T, page playwright.Page) (float64, float64) {
	t.Helper()
	btn := page.FrameLocator("#inner-frame").Locator("#inner-btn")
	box, err := btn.BoundingBox()
	if err != nil || box == nil {
		t.Fatalf("inner button box: %v", err)
	}
	return box.X + box.Width/2, box.Y + box.Height/2
}

func pickAtPoint(t *testing.T, html, routePattern string, pointFn func(*testing.T, playwright.Page) (float64, float64)) string {
	t.Helper()
	page := openPickerFixture(t, html, routePattern)
	defer page.Close()
	x, y := pointFn(t, page)
	return pickAtPagePoint(t, page, x, y)
}

func pickAtPagePoint(t *testing.T, page playwright.Page, x, y float64) string {
	t.Helper()
	picked := installPickerBindings(t, page)
	if _, err := page.Evaluate(RecorderHeuristicsJS); err != nil {
		t.Fatalf("heuristics: %v", err)
	}
	if err := ApplyLibraryHeuristics(page, true, true); err != nil {
		t.Fatalf("library heuristics: %v", err)
	}
	if _, err := page.Evaluate(PickerInstallScript); err != nil {
		t.Fatalf("picker: %v", err)
	}
	if err := page.Mouse().Click(x, y); err != nil {
		t.Fatalf("click: %v", err)
	}
	select {
	case value := <-picked:
		return value
	case <-time.After(3 * time.Second):
		t.Fatal("picker timed out")
		return ""
	}
}

func pickPayloadAtPoint(t *testing.T, html, routePattern string, pointFn func(*testing.T, playwright.Page) (float64, float64)) any {
	t.Helper()
	page := openPickerFixture(t, html, routePattern)
	defer page.Close()
	x, y := pointFn(t, page)
	picked := installPickerPayloadBindings(t, page)
	if _, err := page.Evaluate(RecorderHeuristicsJS); err != nil {
		t.Fatalf("heuristics: %v", err)
	}
	if err := ApplyLibraryHeuristics(page, true, true); err != nil {
		t.Fatalf("library heuristics: %v", err)
	}
	if _, err := page.Evaluate(PickerInstallScript); err != nil {
		t.Fatalf("picker: %v", err)
	}
	if err := page.Mouse().Click(x, y); err != nil {
		t.Fatalf("click: %v", err)
	}
	select {
	case value := <-picked:
		return value
	case <-time.After(3 * time.Second):
		t.Fatal("picker timed out")
		return nil
	}
}

func openPickerFixture(t *testing.T, html, routePattern string) playwright.Page {
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

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		t.Fatalf("launch: %v", err)
	}
	t.Cleanup(func() { _ = browser.Close() })

	page, err := browser.NewPage()
	if err != nil {
		t.Fatalf("new page: %v", err)
	}

	if routePattern != "" {
		if err := page.Route(routePattern, func(route playwright.Route) {
			_ = route.Fulfill(playwright.RouteFulfillOptions{
				Status:      playwright.Int(200),
				ContentType: playwright.String("text/html"),
				Body:        "<!doctype html><html><body><button>Telegram</button></body></html>",
			})
		}); err != nil {
			t.Fatalf("route: %v", err)
		}
	}

	if err := page.SetContent(html); err != nil {
		t.Fatalf("set content: %v", err)
	}
	if strings.Contains(html, "inner-frame") {
		if _, err := page.WaitForFunction(`() => {
			const frame = document.getElementById('inner-frame');
			return !!(frame && frame.contentDocument && frame.contentDocument.getElementById('inner-btn'));
		}`, nil); err != nil {
			t.Fatalf("iframe content: %v", err)
		}
	}
	if strings.Contains(html, "attachShadow") {
		if _, err := page.WaitForFunction(`() => !!document.getElementById('host')?.shadowRoot?.querySelector('#shadow-btn')`, nil); err != nil {
			t.Fatalf("shadow root: %v", err)
		}
	}
	return page
}

func installPickerBindings(t *testing.T, page playwright.Page) chan string {
	t.Helper()
	picked := make(chan string, 1)
	ctx := page.Context()
	if err := ctx.ExposeBinding("pickSelectorDone", func(_ *playwright.BindingSource, args ...any) any {
		if len(args) > 0 {
			select {
			case picked <- pickerBindingSelector(args[0]):
			default:
			}
		}
		return nil
	}); err != nil {
		t.Fatalf("expose done: %v", err)
	}
	if err := ctx.ExposeBinding("pickSelectorCancel", func(_ *playwright.BindingSource, _ ...any) any {
		select {
		case picked <- "":
		default:
		}
		return nil
	}); err != nil {
		t.Fatalf("expose cancel: %v", err)
	}
	return picked
}

func installPickerPayloadBindings(t *testing.T, page playwright.Page) chan any {
	t.Helper()
	picked := make(chan any, 1)
	ctx := page.Context()
	if err := ctx.ExposeBinding("pickSelectorDone", func(_ *playwright.BindingSource, args ...any) any {
		if len(args) > 0 {
			select {
			case picked <- args[0]:
			default:
			}
		}
		return nil
	}); err != nil {
		t.Fatalf("expose done: %v", err)
	}
	if err := ctx.ExposeBinding("pickSelectorCancel", func(_ *playwright.BindingSource, _ ...any) any {
		select {
		case picked <- "":
		default:
		}
		return nil
	}); err != nil {
		t.Fatalf("expose cancel: %v", err)
	}
	return picked
}

func pickerBindingSelector(arg any) string {
	switch v := arg.(type) {
	case string:
		return v
	case map[string]interface{}:
		if s, ok := v["selector"].(string); ok {
			return s
		}
	}
	return ""
}

func frameByURLContains(page playwright.Page, needle string) playwright.Frame {
	for _, frame := range page.Frames() {
		if strings.Contains(frame.URL(), needle) {
			return frame
		}
	}
	return nil
}

func asNumber(value any) (float64, bool) {
	switch n := value.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}

func warningsContain(value any, want string) bool {
	items, ok := value.([]interface{})
	if !ok {
		return false
	}
	for _, item := range items {
		if s, ok := item.(string); ok && s == want {
			return true
		}
	}
	return false
}
