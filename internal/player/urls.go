package player

import (
	"fmt"
	"net/url"
	"strings"

	playwright "github.com/mxschmitt/playwright-go"
)

const (
	NavWaitUntil = "domcontentloaded"
	NavTimeoutMs = 30000
)

func UrlsMatch(current, target string) bool {
	current = strings.TrimSpace(current)
	target = strings.TrimSpace(target)
	if current == "" || target == "" {
		return false
	}
	if strings.TrimRight(current, "/") == strings.TrimRight(target, "/") {
		return true
	}
	cur, err1 := url.Parse(current)
	tgt, err2 := url.Parse(target)
	if err1 != nil || err2 != nil {
		return false
	}
	if cur.Scheme != tgt.Scheme || cur.Host != tgt.Host {
		return false
	}
	if strings.TrimRight(cur.Path, "/") != strings.TrimRight(tgt.Path, "/") {
		return false
	}
	if tgt.RawQuery != "" && cur.RawQuery != tgt.RawQuery {
		return false
	}
	if tgt.Fragment != "" && cur.Fragment != tgt.Fragment {
		return false
	}
	return true
}

// ParseNavWaitUntil maps CLI/settings values to Playwright navigation wait states.
func ParseNavWaitUntil(s string) (*playwright.WaitUntilState, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "domcontentloaded", "dom":
		return playwright.WaitUntilStateDomcontentloaded, nil
	case "load":
		return playwright.WaitUntilStateLoad, nil
	case "networkidle", "network_idle":
		return playwright.WaitUntilStateNetworkidle, nil
	case "commit":
		return playwright.WaitUntilStateCommit, nil
	default:
		return nil, fmt.Errorf("unsupported nav wait until %q (supported: load, domcontentloaded, networkidle, commit)", s)
	}
}
