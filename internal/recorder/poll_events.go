package recorder

import (
	"fmt"
	"time"

	playwright "github.com/mxschmitt/playwright-go"
)

const drainRecorderEventsJS = `() => {
	const r = window.__scenariaRecorder;
	if (!r || !r.events.length) return [];
	const out = r.events.splice(0, r.events.length);
	return out;
}`

type recorderPollState struct {
	lastURL       string
	navCorr       navCorrelation
	recordURLWait bool
}

func newRecorderPollState(startURL string, recordURLWait bool) recorderPollState {
	return recorderPollState{
		lastURL:       startURL,
		recordURLWait: recordURLWait,
	}
}

func drainRecorderEvents(page playwright.Page) ([]recorderEvent, error) {
	raw, err := page.Evaluate(drainRecorderEventsJS)
	if err != nil {
		return nil, fmt.Errorf("read recorder events: %w", err)
	}
	return decodeEvents(raw)
}

func applyRecorderPollBatch(
	recorded *[]RecordedStep,
	state *recorderPollState,
	events []recorderEvent,
	currentURL string,
	now time.Time,
	notify StepNotifier,
) (hadEvents bool, urlChanged bool) {
	state.navCorr.expire(now)
	for _, event := range events {
		detail := normalizeDetail(event.Detail)
		step, ok := EventToRecordedStep(event.Type, detail)
		if !ok {
			continue
		}
		appendCoalescedStep(recorded, step, notify)
		if isNavCausingStep(step) {
			state.navCorr.noteNavCausingStep(now)
		}
		hadEvents = true
	}
	if currentURL != "" && currentURL != state.lastURL {
		switch classifyURLNavigation(&state.navCorr, now, state.recordURLWait) {
		case "wait-url":
			appendWaitURLStep(recorded, currentURL, notify)
		case "goto":
			appendGotoStep(recorded, currentURL, notify)
		}
		state.lastURL = currentURL
		urlChanged = true
	}
	return hadEvents, urlChanged
}
