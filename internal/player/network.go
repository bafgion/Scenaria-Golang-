package player

import (
	"fmt"
	"strings"

	playwright "github.com/mxschmitt/playwright-go"
)

func wireNetworkFailureListener(session *browserSession) {
	if session == nil || session.page == nil {
		return
	}
	session.page.OnRequestFailed(func(req playwright.Request) {
		if req == nil {
			return
		}
		snippet := formatFailedRequest(req)
		if snippet == "" {
			return
		}
		session.networkMu.Lock()
		session.lastNetworkFail = snippet
		session.networkMu.Unlock()
	})
}

func formatFailedRequest(req playwright.Request) string {
	method := strings.TrimSpace(req.Method())
	url := strings.TrimSpace(req.URL())
	if url == "" {
		return ""
	}
	failText := ""
	if err := req.Failure(); err != nil {
		failText = strings.TrimSpace(err.Error())
	}
	if method == "" {
		method = "GET"
	}
	if failText == "" {
		return fmt.Sprintf("%s %s", method, url)
	}
	return fmt.Sprintf("%s %s — %s", method, url, failText)
}

func (s *browserSession) clearNetworkFailure() {
	if s == nil {
		return
	}
	s.networkMu.Lock()
	s.lastNetworkFail = ""
	s.networkMu.Unlock()
}
