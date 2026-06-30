package recorder

import (
	"fmt"

	playwright "github.com/mxschmitt/playwright-go"
)

// PageRecorderConfig is injected into the browser recorder script.
type PageRecorderConfig struct {
	FilterImportant   bool
	NavOnly           bool
	HoverRecord       bool
	ScrollBeforeClick bool
	HoverRecordMinMs  int
}

func normalizeHoverRecordMinMs(ms int) int {
	if ms <= 0 {
		return 600
	}
	return ms
}

func pageRecorderConfigJS(c PageRecorderConfig) string {
	return fmt.Sprintf(`() => {
		if (!window.__scenariaRecorder) return;
		window.__scenariaRecorder.filterImportant = %v;
		window.__scenariaRecorder.navOnly = %v;
		window.__scenariaRecorder.hoverRecord = %v;
		window.__scenariaRecorder.scrollBeforeClick = %v;
		window.__scenariaRecorder.hoverRecordMinMs = %d;
	}`, c.FilterImportant, c.NavOnly, c.HoverRecord, c.ScrollBeforeClick, normalizeHoverRecordMinMs(c.HoverRecordMinMs))
}

func ApplyPageRecorderConfig(page playwright.Page, c PageRecorderConfig) error {
	if page == nil {
		return nil
	}
	_, err := page.Evaluate(pageRecorderConfigJS(c))
	return err
}
