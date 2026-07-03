package winfocus

import (
	"fmt"
	"strings"

	playwright "github.com/mxschmitt/playwright-go"
)

// BringPageToFront activates the Playwright page and, on supported platforms,
// raises the native browser window above other applications.
func BringPageToFront(page playwright.Page) error {
	if page == nil {
		return fmt.Errorf("браузер не открыт")
	}
	if page.IsClosed() {
		return fmt.Errorf("браузер не открыт")
	}
	if err := page.BringToFront(); err != nil {
		return err
	}
	_, _ = page.Evaluate(`() => { window.focus(); }`, nil)
	title, err := page.Title()
	if err != nil {
		title = ""
	}
	pageURL := page.URL()
	if err := raiseNativeWindow(strings.TrimSpace(title), pageURL); err != nil {
		return err
	}
	return nil
}
