package report

import (
	"strings"
	"testing"
)

func TestViewerAssetsScopeFilterCheckbox(t *testing.T) {
	if strings.Contains(viewerCSS, ".filters input, .filters select") {
		t.Fatal("filter checkbox must not inherit full-width text input styles")
	}
	for _, want := range []string{
		`.filters input:not([type="checkbox"]), .filters select`,
		`.filters input[type="checkbox"]`,
		`appearance: none`,
		`.check-filter`,
	} {
		if !strings.Contains(viewerCSS, want) {
			t.Fatalf("viewer css missing %q", want)
		}
	}
}

func TestViewerShellHasCleanFallbacksAndCloseButton(t *testing.T) {
	for _, bad := range []string{"Рџ", "РЎ", "вЂ", "Р’", "Р¤", "Рњ"} {
		if strings.Contains(viewerHTMLShell, bad) {
			t.Fatalf("viewer html shell contains mojibake marker %q", bad)
		}
	}
	for _, want := range []string{
		`placeholder="Search..."`,
		`class="check-filter"`,
		`class="mode-switch"`,
		`role="tablist"`,
		`aria-label="Search scenarios"`,
		`role="tabpanel"`,
		`id="close-inspector"`,
		`id="inspector-title"`,
	} {
		if !strings.Contains(viewerHTMLShell, want) {
			t.Fatalf("viewer html shell missing %q", want)
		}
	}
}

func TestViewerRuntimeLocaleOverrides(t *testing.T) {
	for _, want := range []string{
		`Object.assign(STRINGS.ru`,
		`timeline: '\u0428\u0430\u0433\u0438'`,
		`copySelector: '\u041a\u043e\u043f\u0438\u0440\u043e\u0432\u0430\u0442\u044c \u0441\u0435\u043b\u0435\u043a\u0442\u043e\u0440'`,
		`exportFailed: '\u042d\u043a\u0441\u043f\u043e\u0440\u0442 \u0443\u043f\u0430\u0432\u0448\u0438\u0445 .feature'`,
		`modeUnavailable: '\u0421\u043d\u0430\u0447\u0430\u043b\u0430 \u0441\u0433\u0435\u043d\u0435\u0440\u0438\u0440\u0443\u0439\u0442\u0435 paired HTML-\u043e\u0442\u0447\u0435\u0442`,
		`Object.assign(STRINGS.en`,
		`searchPlaceholder: 'Search...'`,
		`function reportStateHash()`,
		`function restoreReportStateFromHash()`,
		`function activateOnEnterSpace`,
		`aria-selected`,
		`aria-current`,
	} {
		if !strings.Contains(viewerJS, want) {
			t.Fatalf("viewer js missing %q", want)
		}
	}
}

func TestViewerMobileLayoutAvoidsHeaderMagicNumber(t *testing.T) {
	if strings.Contains(viewerCSS, "calc(100vh - 52px)") || strings.Contains(viewerCSS, "top: 52px") {
		t.Fatal("viewer css must not depend on a hard-coded 52px header height")
	}
	for _, want := range []string{
		`display: grid`,
		`grid-template-rows: auto minmax(0, 1fr)`,
		`transform: translateY(100%)`,
		`.panel.inspector-panel.open { transform: translateY(0); }`,
	} {
		if !strings.Contains(viewerCSS, want) {
			t.Fatalf("viewer css missing %q", want)
		}
	}
}
