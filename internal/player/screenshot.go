package player

import playwright "github.com/mxschmitt/playwright-go"

func captureFailureScreenshot(session *browserSession) []byte {
	return capturePageScreenshot(session, true)
}

func captureViewportScreenshot(session *browserSession) []byte {
	return capturePageScreenshot(session, false)
}

func capturePageScreenshot(session *browserSession, fullPage bool) []byte {
	if session == nil {
		return nil
	}
	session.mu.Lock()
	defer session.mu.Unlock()
	if session.isClosed() || session.page == nil {
		return nil
	}
	data, err := session.page.Screenshot(playwright.PageScreenshotOptions{
		FullPage: playwright.Bool(fullPage),
	})
	if err != nil {
		return nil
	}
	return data
}
