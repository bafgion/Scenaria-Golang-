package player

import "strings"

func capturePageContext(session *browserSession) string {
	if session == nil {
		return ""
	}
	session.mu.Lock()
	defer session.mu.Unlock()
	if session.isClosed() || session.page == nil {
		return ""
	}
	title, _ := session.page.Title()
	url := session.page.URL()
	title = strings.TrimSpace(title)
	url = strings.TrimSpace(url)
	switch {
	case title != "" && url != "":
		return title + "\n" + url
	case url != "":
		return url
	default:
		return title
	}
}
