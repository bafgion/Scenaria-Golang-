//go:build integration

package selector

import (
	"context"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

func TestValidateActionDisabledButtonNotClickable(t *testing.T) {
	html := `<!doctype html><html><body><button id="btn" disabled>OK</button></body></html>`
	page := openPickerFixture(t, html, "")
	v := Validator{}
	result := v.ValidateActionTarget(context.Background(), page, stepdsl.Action{Kind: "click"}, "#btn", 3*time.Second)
	if result.OK {
		t.Fatal("expected disabled button to fail click validation")
	}
	if !strings.Contains(result.Message, "not enabled") {
		t.Fatalf("got %q", result.Message)
	}
}

func TestValidateActionFillOnDivNotEditable(t *testing.T) {
	html := `<!doctype html><html><body><div id="fake-input">text</div></body></html>`
	page := openPickerFixture(t, html, "")
	v := Validator{}
	result := v.ValidateActionTarget(context.Background(), page, stepdsl.Action{Kind: "fill"}, "#fake-input", 3*time.Second)
	if result.OK {
		t.Fatal("expected non-editable div to fail fill validation")
	}
	if !strings.Contains(result.Message, "not an editable target") {
		t.Fatalf("got %q", result.Message)
	}
}

func TestValidateActionChainedAmbiguousWarning(t *testing.T) {
	html := `<!doctype html><html><body>
<div class="box"><button>A</button></div>
<div class="box"><button>B</button></div>
</body></html>`
	page := openPickerFixture(t, html, "")
	v := Validator{}
	result := v.ValidateActionTarget(context.Background(), page, stepdsl.Action{Kind: "click"}, `.box >> button`, 3*time.Second)
	if result.MatchCount < 2 {
		t.Fatalf("expected multiple matches, got %d", result.MatchCount)
	}
	if len(result.Warnings) == 0 {
		t.Fatal("expected ambiguity warning for chained selector")
	}
}

func TestResolveChainedLocatorFallsBackToVisibleCardButton(t *testing.T) {
	html := `<!doctype html><html><body>
<div class="product-card">
  <h3>Джинсы свободные RL20015</h3>
  <span>980 ₽</span><span>+3</span>
  <button id="open-product">Подробнее</button>
</div>
<script>
document.getElementById('open-product').addEventListener('click', () => {
  document.body.setAttribute('data-clicked', 'product');
});
</script>
</body></html>`
	page := openPickerFixture(t, html, "")
	locator := ResolveChainedLocator(page, `div:has-text("Джинсы свободные RL20015 980 ₽+3") >> button:has-text("Джинсы свободные RL20015 980 ₽+3")`)
	if err := locator.Click(); err != nil {
		t.Fatalf("click recovered card target: %v", err)
	}
	raw, err := page.GetAttribute("body", "data-clicked")
	if err != nil {
		t.Fatalf("read click marker: %v", err)
	}
	if raw != "product" {
		t.Fatalf("expected recovered button click, got %q", raw)
	}
}

func TestResolveChainedLocatorPrefersProductLinkForLongCardButtonSelector(t *testing.T) {
	html := `<!doctype html><html><body>
<div class="product-card">
  <a id="product-link" href="/products/relaxed-jeans">
    <h3>Relaxed jeans RL20015</h3>
    <span>980 RUB</span><span>+3</span>
  </a>
  <button id="quick-view">Quick view</button>
</div>
<script>
document.getElementById('product-link').addEventListener('click', (event) => {
  event.preventDefault();
  document.body.setAttribute('data-clicked', 'product-link');
});
document.getElementById('quick-view').addEventListener('click', () => {
  document.body.setAttribute('data-clicked', 'quick-view');
});
</script>
</body></html>`
	page := openPickerFixture(t, html, "")
	locator := ResolveChainedLocator(page, `div:has-text("Relaxed jeans RL20015 980 RUB+3") >> button:has-text("Relaxed jeans RL20015 980 RUB+3")`)
	if err := locator.Click(); err != nil {
		t.Fatalf("click recovered card target: %v", err)
	}
	raw, err := page.GetAttribute("body", "data-clicked")
	if err != nil {
		t.Fatalf("read click marker: %v", err)
	}
	if raw != "product-link" {
		t.Fatalf("expected recovered product link click, got %q", raw)
	}
}

func TestValidateFeatureStaticDynamicUIWarning(t *testing.T) {
	html := `<!doctype html><html><body>
<button id="open">Open</button>
<input id="field" type="text" style="display:none">
<script>
document.getElementById('open').addEventListener('click', () => {
  document.getElementById('field').style.display = 'block';
});
</script>
</body></html>`
	startURL := "data:text/html;charset=utf-8," + url.PathEscape(html)

	feature, err := gherkin.ParseFeature(`Функционал: dynamic
  Сценарий: reveal field
    Допустим открыт "` + startURL + `"
    И нажимаю "#open"
    И ввожу "x" в "#field"
`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	v := Validator{}
	results, err := v.ValidateFeatureInBrowserDetailed(context.Background(), "dynamic.feature", feature, BrowserValidateOptions{
		BrowserName: "chromium",
		Headless:    true,
		Mode:        ValidationModeStatic,
	})
	if err != nil {
		t.Fatalf("validate: %v", err)
	}

	var fillResult *StepValidation
	for i := range results {
		if results[i].ActionKind == "fill" {
			fillResult = &results[i]
			break
		}
	}
	if fillResult == nil {
		t.Fatal("fill step result not found")
	}
	if fillResult.Status != "warning" {
		t.Fatalf("expected warning for hidden field in static mode, got %q: %s", fillResult.Status, fillResult.Message)
	}
	if !strings.Contains(fillResult.Message, "dynamic UI") {
		t.Fatalf("expected dynamic UI hint, got %q", fillResult.Message)
	}
}
