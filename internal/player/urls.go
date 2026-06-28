package player

import (
	"net/url"
	"strings"
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
	return cur.Scheme == tgt.Scheme &&
		cur.Host == tgt.Host &&
		strings.TrimRight(cur.Path, "/") == strings.TrimRight(tgt.Path, "/")
}
