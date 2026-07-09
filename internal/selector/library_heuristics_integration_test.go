//go:build integration

package selector

import (
	"strings"
	"testing"
	"time"

	playwright "github.com/mxschmitt/playwright-go"
)

const muiButtonHTML = `<!doctype html><html><body style="margin:24px">
<button class="MuiButton-root MuiButton-contained" type="button">Save</button>
</body></html>`

const muiInputHTML = `<!doctype html><html><body style="margin:24px">
<div class="MuiInputBase-root MuiOutlinedInput-root">
<input class="MuiInputBase-input" name="email" type="email" />
</div>
</body></html>`

const antButtonHTML = `<!doctype html><html><body style="margin:24px">
<button class="ant-btn ant-btn-primary" type="button"><span>Submit</span></button>
</body></html>`

const antInputHTML = `<!doctype html><html><body style="margin:24px">
<input class="ant-input" placeholder="Search users" type="text" />
</body></html>`

func TestPickerMUIButtonPrefersLibrarySelector(t *testing.T) {
	page := openPickerFixture(t, muiButtonHTML, "")
	selector := pickAtPagePointWithLibrary(t, page, 40, 40, true, false)
	if !strings.Contains(selector, "MuiButton-root") {
		t.Fatalf("expected MUI button selector, got %q", selector)
	}
}

func TestPickerMUIInputPrefersLibrarySelector(t *testing.T) {
	page := openPickerFixture(t, muiInputHTML, "")
	selector := pickAtPagePointWithLibrary(t, page, 40, 40, true, false)
	if !strings.Contains(selector, "MuiInputBase-input") {
		t.Fatalf("expected MUI input selector, got %q", selector)
	}
}

func TestPickerAntButtonPrefersLibrarySelector(t *testing.T) {
	page := openPickerFixture(t, antButtonHTML, "")
	selector := pickAtPagePointWithLibrary(t, page, 40, 40, false, true)
	if !strings.Contains(selector, "ant-btn") {
		t.Fatalf("expected Ant button selector, got %q", selector)
	}
}

func TestPickerAntInputPrefersLibrarySelector(t *testing.T) {
	page := openPickerFixture(t, antInputHTML, "")
	selector := pickAtPagePointWithLibrary(t, page, 40, 40, false, true)
	if !strings.Contains(selector, "ant-input") {
		t.Fatalf("expected Ant input selector, got %q", selector)
	}
}

func TestPickerGenericSiteUnchangedWithLibraryPacksEnabled(t *testing.T) {
	withLibrary := pickAtSelectorWithLibrary(t, nestedLabelFormHTML, "", `label:nth-of-type(1) input`, true, true)
	withoutLibrary := pickAtSelectorWithLibrary(t, nestedLabelFormHTML, "", `label:nth-of-type(1) input`, false, false)
	if withLibrary != withoutLibrary {
		t.Fatalf("generic picker changed: with=%q without=%q", withLibrary, withoutLibrary)
	}
}

func TestPickerMUIFallsBackWhenLibraryPackDisabled(t *testing.T) {
	page := openPickerFixture(t, muiButtonHTML, "")
	selector := pickAtPagePointWithLibrary(t, page, 40, 40, false, false)
	if strings.Contains(selector, "MuiButton-root") {
		t.Fatalf("expected generic selector when MUI pack disabled, got %q", selector)
	}
}

func pickAtSelectorWithLibrary(t *testing.T, html, routePattern, clickSelector string, mui, ant bool) string {
	t.Helper()
	page := openPickerFixture(t, html, routePattern)
	defer page.Close()
	loc := page.Locator(clickSelector)
	box, err := loc.BoundingBox()
	if err != nil || box == nil {
		t.Fatalf("bounding box for %q: %v", clickSelector, err)
	}
	return pickAtPagePointWithLibrary(t, page, box.X+box.Width/2, box.Y+box.Height/2, mui, ant)
}

func pickAtPagePointWithLibrary(t *testing.T, page playwright.Page, x, y float64, mui, ant bool) string {
	t.Helper()
	picked := installPickerBindings(t, page)
	if _, err := page.Evaluate(RecorderHeuristicsJS); err != nil {
		t.Fatalf("heuristics: %v", err)
	}
	if err := ApplyLibraryHeuristics(page, mui, ant); err != nil {
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
