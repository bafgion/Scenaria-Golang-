package player

import playwright "github.com/mxschmitt/playwright-go"

func captureFailureScreenshot(session *browserSession) []byte {
	if session == nil || session.closed || session.page == nil {
		return nil
	}
	data, err := session.page.Screenshot(playwright.PageScreenshotOptions{
		FullPage: playwright.Bool(true),
	})
	if err != nil {
		return nil
	}
	return data
}
