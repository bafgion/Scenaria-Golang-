package player

import playwright "github.com/mxschmitt/playwright-go"

func captureFailureScreenshot(session *browserSession) []byte {
	return capturePageScreenshot(session, true)
}

func captureViewportScreenshot(session *browserSession) []byte {
	return capturePageScreenshot(session, false)
}

func capturePageScreenshot(session *browserSession, fullPage bool) []byte {
	page, err := session.currentPage()
	if err != nil {
		return nil
	}
	data, err := page.Screenshot(playwright.PageScreenshotOptions{
		FullPage: playwright.Bool(fullPage),
	})
	if err != nil {
		return nil
	}
	return data
}
