package player

import (
	"fmt"
	"strings"

	playwright "github.com/mxschmitt/playwright-go"
)

func wireNetworkFailureListener(session *browserSession) {
	page, err := session.currentPage()
	if err != nil {
		return
	}
	wireNetworkFailurePage(session, page)
}

func wireNetworkFailureListenerLocked(session *browserSession) {
	if session == nil || session.page == nil {
		return
	}
	wireNetworkFailurePage(session, session.page)
}

func wireNetworkFailurePage(session *browserSession, page playwright.Page) {
	page.OnRequestFailed(func(req playwright.Request) {
		if req == nil {
			return
		}
		session.recordNetworkSnippet(formatFailedRequest(req))
	})
	page.OnResponse(func(resp playwright.Response) {
		if resp == nil {
			return
		}
		session.recordNetworkSnippet(formatHTTPErrorResponse(resp))
	})
}

func (s *browserSession) recordNetworkSnippet(snippet string) {
	if s == nil || snippet == "" {
		return
	}
	s.networkMu.Lock()
	s.lastNetworkFail = snippet
	s.networkMu.Unlock()
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

func formatHTTPErrorResponse(resp playwright.Response) string {
	if resp == nil {
		return ""
	}
	status := resp.Status()
	if status < 400 {
		return ""
	}
	method := "GET"
	url := ""
	if req := resp.Request(); req != nil {
		method = strings.TrimSpace(req.Method())
		url = strings.TrimSpace(req.URL())
	}
	return formatHTTPErrorSnippet(method, url, status)
}

func formatHTTPErrorSnippet(method, url string, status int) string {
	if status < 400 || url == "" {
		return ""
	}
	if method == "" {
		method = "GET"
	}
	return fmt.Sprintf("%s %s — HTTP %d", method, url, status)
}

func (s *browserSession) clearNetworkFailure() {
	if s == nil {
		return
	}
	s.networkMu.Lock()
	s.lastNetworkFail = ""
	s.networkMu.Unlock()
}

func (s *browserSession) clearNetworkFailureLocked() {
	if s == nil {
		return
	}
	s.networkMu.Lock()
	defer s.networkMu.Unlock()
	s.lastNetworkFail = ""
}
