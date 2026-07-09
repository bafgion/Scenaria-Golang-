package recorder

import (
	"testing"
	"time"
)

func TestClassifyURLNavigationExplicitGoto(t *testing.T) {
	got := classifyURLNavigation(nil, time.Now(), true)
	if got != "goto" {
		t.Fatalf("got %q want goto", got)
	}
}

func TestClassifyURLNavigationAfterClick(t *testing.T) {
	corr := &navCorrelation{pending: true, since: time.Now()}
	got := classifyURLNavigation(corr, time.Now(), true)
	if got != "wait-url" {
		t.Fatalf("got %q want wait-url", got)
	}
	if corr.pending {
		t.Fatal("expected pending cleared")
	}
}

func TestClassifyURLNavigationAfterClickWithoutURLWait(t *testing.T) {
	corr := &navCorrelation{pending: true, since: time.Now()}
	got := classifyURLNavigation(corr, time.Now(), false)
	if got != "" {
		t.Fatalf("got %q want empty", got)
	}
}

func TestClassifyURLNavigationExpiredCorrelation(t *testing.T) {
	corr := &navCorrelation{pending: true, since: time.Now().Add(-3 * time.Second)}
	got := classifyURLNavigation(corr, time.Now(), true)
	if got != "goto" {
		t.Fatalf("got %q want goto", got)
	}
}

func TestPollTickOrderingClickBeforeNavigation(t *testing.T) {
	recorded := []RecordedStep{{Action: "goto", Value: "https://example.com/login"}}
	state := newRecorderPollState(recorded[0].Value, true)

	_, _ = applyRecorderPollBatch(&recorded, &state, []recorderEvent{
		{Type: "click", Detail: map[string]string{"selector": `button:has-text("Войти")`}},
	}, "https://example.com/login", time.Now(), nil)
	_, urlChanged := applyRecorderPollBatch(&recorded, &state, nil, "https://example.com/dashboard", time.Now().Add(150*time.Millisecond), nil)
	if !urlChanged {
		t.Fatal("expected url change")
	}

	if len(recorded) != 3 {
		t.Fatalf("steps: %+v", recorded)
	}
	if recorded[1].Action != "click" {
		t.Fatalf("click should precede navigation, got %+v", recorded[1])
	}
	if recorded[2].Action != "wait-url" || recorded[2].Value != "https://example.com/dashboard" {
		t.Fatalf("navigation step: %+v", recorded[2])
	}
}

func TestDelayedNavigationCorrelationAcrossPollTicks(t *testing.T) {
	recorded := []RecordedStep{{Action: "goto", Value: "https://example.com"}}
	state := newRecorderPollState(recorded[0].Value, true)
	clickAt := time.Now()

	applyRecorderPollBatch(&recorded, &state, []recorderEvent{
		{Type: "click", Detail: map[string]string{"selector": "#submit"}},
	}, "https://example.com", clickAt, nil)

	_, urlChanged := applyRecorderPollBatch(&recorded, &state, nil, "https://example.com/dashboard", clickAt.Add(900*time.Millisecond), nil)
	if !urlChanged {
		t.Fatal("expected delayed url change")
	}
	if recorded[len(recorded)-1].Action != "wait-url" {
		t.Fatalf("expected wait-url after delayed navigation, got %+v", recorded[len(recorded)-1])
	}
}

func TestSPANavigationCorrelation(t *testing.T) {
	recorded := []RecordedStep{{Action: "goto", Value: "https://spa.example/app"}}
	state := newRecorderPollState(recorded[0].Value, true)
	now := time.Now()

	applyRecorderPollBatch(&recorded, &state, []recorderEvent{
		{Type: "click", Detail: map[string]string{"selector": "#route-btn"}},
	}, "https://spa.example/app", now, nil)
	_, urlChanged := applyRecorderPollBatch(&recorded, &state, nil, "https://spa.example/app/dashboard", now.Add(120*time.Millisecond), nil)
	if !urlChanged {
		t.Fatal("expected spa url change")
	}
	if recorded[len(recorded)-1].Action != "wait-url" {
		t.Fatalf("spa navigation should be wait-url, got %+v", recorded[len(recorded)-1])
	}
}

func TestPollTickRedirectWithoutClickUsesGoto(t *testing.T) {
	recorded := []RecordedStep{{Action: "goto", Value: "https://example.com"}}
	state := newRecorderPollState(recorded[0].Value, true)
	_, urlChanged := applyRecorderPollBatch(&recorded, &state, nil, "https://example.com/redirected", time.Now(), nil)
	if !urlChanged {
		t.Fatal("expected redirect")
	}
	if len(recorded) != 2 || recorded[1].Action != "goto" {
		t.Fatalf("expected explicit goto, got %+v", recorded)
	}
}

func TestIsNavCausingStep(t *testing.T) {
	if !isNavCausingStep(RecordedStep{Action: "click", Selector: "#x"}) {
		t.Fatal("click")
	}
	if isNavCausingStep(RecordedStep{Action: "fill", Selector: "#x"}) {
		t.Fatal("fill should not correlate")
	}
	if !isNavCausingStep(RecordedStep{Action: "press", Value: "Enter"}) {
		t.Fatal("enter press")
	}
}
