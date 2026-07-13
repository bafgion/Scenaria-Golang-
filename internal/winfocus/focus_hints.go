package winfocus

import (
	"net/url"
	"strings"
)

func hostFromURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	host := strings.ToLower(strings.TrimSpace(u.Hostname()))
	if host == "" {
		return ""
	}
	return host
}

func scoreChromiumTitle(title, titleHint, urlHint string) int {
	lower := strings.ToLower(strings.TrimSpace(title))
	if lower == "" {
		return 0
	}
	titleHint = strings.ToLower(strings.TrimSpace(titleHint))
	if titleHint != "" && strings.Contains(lower, titleHint) {
		return 100 + len(titleHint)
	}
	if host := hostFromURL(urlHint); host != "" && strings.Contains(lower, host) {
		return 90 + len(host)
	}
	if strings.Contains(lower, "chromium") {
		return 10
	}
	return 0
}
