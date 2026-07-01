package gui

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/recorder"
)

func attachLiveSession(t *testing.T, browse bool) (*Service, *recorder.LiveSession) {
	t.Helper()
	svc := &Service{}
	session := recorder.NewLiveSession()
	if browse {
		session.InitBrowseMode()
	} else {
		session.InitRecordMode()
	}
	svc.mu.Lock()
	svc.liveSession = session
	svc.mu.Unlock()
	return svc, session
}

func TestBeginRecordingCaptureEmitsRecordedSteps(t *testing.T) {
	svc, session := attachLiveSession(t, true)
	steps := []recorder.RecordedStep{{Action: "goto", Value: "https://example.com"}}
	session.Bind(nil, &steps)

	stepEvents := 0
	svc.mu.Lock()
	svc.recordEmit = func(name string, _ any) {
		if name == "record-step" {
			stepEvents++
		}
	}
	svc.mu.Unlock()

	started, err := svc.BeginRecordingCapture()
	if err != nil {
		t.Fatalf("BeginRecordingCapture: %v", err)
	}
	if !started {
		t.Fatal("expected capture to start")
	}
	if stepEvents != 1 {
		t.Fatalf("expected 1 replayed step, got %d", stepEvents)
	}
}

func TestBeginRecordingCaptureResume(t *testing.T) {
	svc, session := attachLiveSession(t, true)
	if err := session.BeginCapture(); err != nil {
		t.Fatal(err)
	}
	session.EndCapture()

	started, err := svc.BeginRecordingCapture()
	if err != nil {
		t.Fatalf("BeginRecordingCapture: %v", err)
	}
	if !started {
		t.Fatal("expected capture to start on resume")
	}
	if !session.CaptureEnabled() {
		t.Fatal("expected capture enabled")
	}

	startedAgain, err := svc.BeginRecordingCapture()
	if err != nil {
		t.Fatalf("BeginRecordingCapture second call: %v", err)
	}
	if startedAgain {
		t.Fatal("expected already-capturing session to report not started")
	}
}

func TestBeginRecordingCaptureWithoutSession(t *testing.T) {
	svc := &Service{}
	started, err := svc.BeginRecordingCapture()
	if err == nil {
		t.Fatal("expected error without live session")
	}
	if started {
		t.Fatal("expected started=false on error")
	}
}

func TestStopRecordingCaptureEndsCaptureOnly(t *testing.T) {
	svc, session := attachLiveSession(t, true)
	if err := session.BeginCapture(); err != nil {
		t.Fatal(err)
	}
	if err := svc.StopRecordingCapture(); err != nil {
		t.Fatalf("StopRecordingCapture: %v", err)
	}
	if session.CaptureEnabled() {
		t.Fatal("expected capture disabled after stop")
	}
}

func TestStopRecordingCaptureClearsSteps(t *testing.T) {
	svc, session := attachLiveSession(t, true)
	steps := []recorder.RecordedStep{
		{Action: "goto", Value: "https://example.com"},
		{Action: "click", Selector: "#btn"},
	}
	session.Bind(nil, &steps)
	if err := session.BeginCapture(); err != nil {
		t.Fatal(err)
	}
	if err := svc.StopRecordingCapture(); err != nil {
		t.Fatalf("StopRecordingCapture: %v", err)
	}
	if len(steps) != 0 {
		t.Fatalf("expected cleared steps, got %d", len(steps))
	}
	if session.RecordedStepCount() != 0 {
		t.Fatalf("expected step count 0 after stop, got %d", session.RecordedStepCount())
	}
}

func TestPauseRecordingKeepsSteps(t *testing.T) {
	svc, session := attachLiveSession(t, true)
	steps := []recorder.RecordedStep{{Action: "click", Selector: "#x"}}
	session.Bind(nil, &steps)
	if err := session.BeginCapture(); err != nil {
		t.Fatal(err)
	}
	svc.PauseRecording()
	if !session.IsPaused() {
		t.Fatal("expected paused")
	}
	if len(steps) != 1 {
		t.Fatalf("pause must keep recorded steps, got %d", len(steps))
	}
	svc.ResumeRecording()
	if session.IsPaused() {
		t.Fatal("expected resumed")
	}
	if len(steps) != 1 {
		t.Fatalf("expected steps after resume, got %d", len(steps))
	}
}

func TestPauseResumeRecordingWithSession(t *testing.T) {
	svc, session := attachLiveSession(t, true)
	if err := session.BeginCapture(); err != nil {
		t.Fatal(err)
	}
	svc.PauseRecording()
	if !session.IsPaused() {
		t.Fatal("expected paused")
	}
	svc.ResumeRecording()
	if session.IsPaused() {
		t.Fatal("expected resumed")
	}
}

func TestChaosBeginRecordingCaptureResume(t *testing.T) {
	for seed := int64(0); seed < 40; seed++ {
		t.Run(fmt.Sprintf("seed-%d", seed), func(t *testing.T) {
			rng := rand.New(rand.NewSource(seed))
			svc, session := attachLiveSession(t, true)
			for i := 0; i < 4+rng.Intn(12); i++ {
				switch rng.Intn(4) {
				case 0, 1:
					started, err := svc.BeginRecordingCapture()
					if err != nil {
						t.Fatalf("BeginRecordingCapture: %v", err)
					}
					if !session.CaptureEnabled() && started {
						t.Fatal("expected capture enabled when started")
					}
				case 2:
					_ = svc.StopRecordingCapture()
				case 3:
					if session.CaptureEnabled() {
						if rng.Intn(2) == 0 {
							svc.PauseRecording()
						} else {
							svc.ResumeRecording()
						}
					}
				}
				if session.IsPaused() && !session.CaptureEnabled() {
					t.Fatal("paused without active capture")
				}
			}
		})
	}
}
