package player

import "strings"

func capturePageContext(session *browserSession) string {
	page, err := session.currentPage()
	if err != nil {
		return ""
	}
	title, _ := page.Title()
	url := page.URL()
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
