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
<label><div>ИНН</div><div><input type="text"></div></label>
<label><div>E-mail</div><div><input type="email"></div></label>
<label><div>Имя</div><div><input type="text"></div></label>
</div></div></body></html>`

const iframeWidgetHTML = `<!doctype html><html><body>
<iframe id="tg-login" title="Telegram login" src="https://oauth.telegram.org/embed/demo"
  style="width:240px;height:52px;border:0"></iframe>
</body></html>`

func TestPickerNestedLabelInputReturnsLabelHasText(t *testing.T) {
	selector := pickAtSelector(t, nestedLabelFormHTML, "", `label:nth-of-type(1) input`)
	if selector == "" {
		t.Fatal("no selector picked")
	}
	if !strings.HasPrefix(selector, `label:has-text("`) || !strings.Contains(selector, "ИНН") {
		t.Fatalf("got %q", selector)
	}
}

func TestPickerLabelCaptionDivReturnsInputLabelHasText(t *testing.T) {
	selector := pickAtSelector(t, nestedLabelFormHTML, "", `label:nth-of-type(2) > div:first-child`)
	if selector == "" {
		t.Fatal("no selector picked")
	}
	if !strings.HasPrefix(selector, `label:has-text("`) || !strings.Contains(strings.ToLower(selector), "mail") {
		t.Fatalf("got %q", selector)
	}
}

func TestPickerShortCaptionReturnsLabelHasText(t *testing.T) {
	selector := pickAtSelector(t, nestedLabelFormHTML, "", `label:nth-of-type(3) input`)
	if selector == "" {
		t.Fatal("no selector picked")
	}
	if !strings.Contains(selector, `label:has-text("Имя")`) {
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

func pickAtSelector(t *testing.T, html, routePattern, clickSelector string) string {
	t.Helper()
	page := openPickerFixture(t, html, routePattern)
	defer page.Close()

	picked := installPickerBindings(t, page)
	if _, err := page.Evaluate(RecorderHeuristicsJS); err != nil {
		t.Fatalf("heuristics: %v", err)
	}
	if _, err := page.Evaluate(PickerInstallScript); err != nil {
		t.Fatalf("picker: %v", err)
	}

	loc := page.Locator(clickSelector)
	box, err := loc.BoundingBox()
	if err != nil || box == nil {
		t.Fatalf("bounding box for %q: %v", clickSelector, err)
	}
	if err := page.Mouse().Click(box.X+box.Width/2, box.Y+box.Height/2); err != nil {
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
	return page
}

func installPickerBindings(t *testing.T, page playwright.Page) chan string {
	t.Helper()
	picked := make(chan string, 1)
	ctx := page.Context()
	if err := ctx.ExposeBinding("pickSelectorDone", func(_ *playwright.BindingSource, args ...any) any {
		if len(args) > 0 {
			if value, ok := args[0].(string); ok {
				select {
				case picked <- value:
				default:
				}
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
