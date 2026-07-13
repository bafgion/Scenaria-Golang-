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
	cdpHint := chromiumWindowHintViaCDP(page)
	_, _ = page.Evaluate(`() => { window.focus(); }`, nil)
	title, err := page.Title()
	if err != nil {
		title = ""
	}
	pageURL := page.URL()
	if err := raiseNativeWindow(strings.TrimSpace(title), pageURL, cdpHint.ProcessID); err != nil {
		return err
	}
	return nil
}

type nativeWindowHint struct {
	ProcessID uint32
}

func chromiumWindowHintViaCDP(page playwright.Page) nativeWindowHint {
	if page == nil || page.IsClosed() || page.Context() == nil {
		return nativeWindowHint{}
	}
	session, err := page.Context().NewCDPSession(page)
	if err != nil {
		return nativeWindowHint{}
	}
	defer func() { _ = session.Detach() }()
	return nativeWindowHint{ProcessID: browserProcessIDFromCDP(session)}
}

func cdpWindowID(raw any) (int, bool) {
	switch v := raw.(type) {
	case int:
		return v, v > 0
	case int64:
		return int(v), v > 0
	case float64:
		if v <= 0 {
			return 0, false
		}
		return int(v), true
	default:
		return 0, false
	}
}

func browserProcessIDFromCDP(session playwright.CDPSession) uint32 {
	if session == nil {
		return 0
	}
	raw, err := session.Send("SystemInfo.getProcessInfo", nil)
	if err != nil {
		return 0
	}
	payload, ok := raw.(map[string]any)
	if !ok {
		return 0
	}
	return browserProcessIDFromPayload(payload)
}

func browserProcessIDFromPayload(payload map[string]any) uint32 {
	items, ok := payload["processInfo"].([]any)
	if !ok {
		return 0
	}
	for _, item := range items {
		process, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if typ, _ := process["type"].(string); typ != "browser" {
			continue
		}
		pid, ok := cdpWindowID(process["id"])
		if ok {
			return uint32(pid)
		}
	}
	return 0
}
