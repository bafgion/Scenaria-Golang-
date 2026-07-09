package player

import (
	"errors"
	"strings"
)

const (
	// MsgBrowserClosed is shown when a scenario cannot start because the browser is closed.
	MsgBrowserClosed = "браузер закрыт — откройте браузер в IDE или уберите шаг «закрываю браузер»"
	// MsgScenarioNotStartedBrowserClosed is shown for scenarios skipped after a prior scenario closed the browser.
	MsgScenarioNotStartedBrowserClosed = "не запущен: предыдущий сценарий закрыл браузер"
)

// IsBrowserSessionClosed reports whether err is a low-level closed-session error from Playwright glue.
func IsBrowserSessionClosed(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "browser session is closed") ||
		strings.Contains(msg, "browser page is not available")
}

// UserFacingBrowserError maps low-level browser lifecycle errors to user-facing text.
func UserFacingBrowserError(err error) string {
	if err == nil {
		return ""
	}
	if IsBrowserSessionClosed(err) {
		return MsgBrowserClosed
	}
	return err.Error()
}

// NormalizeScenarioMessage applies user-facing browser error mapping to scenario result text.
func NormalizeScenarioMessage(message string) string {
	if message == "" {
		return message
	}
	if IsBrowserSessionClosed(errors.New(message)) {
		return MsgBrowserClosed
	}
	return message
}
