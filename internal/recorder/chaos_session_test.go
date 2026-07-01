package recorder

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestChaosRecordingSessionState(t *testing.T) {
	ops := []string{"begin", "end", "pause", "resume"}
	for seed := int64(0); seed < 60; seed++ {
		t.Run(fmt.Sprintf("seed-%d", seed), func(t *testing.T) {
			rng := rand.New(rand.NewSource(seed))
			s := NewLiveSession()
			s.InitBrowseMode()

			for i := 0; i < 8+rng.Intn(24); i++ {
				switch ops[rng.Intn(len(ops))] {
				case "begin":
					if !s.CaptureEnabled() {
						_ = s.BeginCapture()
					}
				case "end":
					if s.CaptureEnabled() {
						s.EndCapture()
					}
				case "pause":
					if s.CaptureEnabled() {
						s.Pause()
					}
				case "resume":
					if s.CaptureEnabled() {
						s.Resume()
					}
				}

				if s.CaptureEnabled() && !s.CaptureEverEnabled() {
					t.Fatal("capture enabled but captureEver is false")
				}
				if s.IsPaused() && !s.CaptureEnabled() {
					t.Fatal("paused while capture disabled")
				}
			}
		})
	}
}

func TestChaosCaptureStepNotifyNeverExceedsRecorded(t *testing.T) {
	for seed := int64(0); seed < 40; seed++ {
		t.Run(fmt.Sprintf("seed-%d", seed), func(t *testing.T) {
			rng := rand.New(rand.NewSource(seed))
			s := NewLiveSession()
			s.InitBrowseMode()
			recorded := make([]RecordedStep, 0, 24)
			s.Bind(nil, &recorded)
			notifies := 0
			totalAppended := 0

			beginCapture := func() {
				if s.CaptureEnabled() {
					return
				}
				if ShouldSyncRecordedStepsOnCaptureStart(s) {
					for i, st := range recorded {
						if line, ok := RecordedStepToLine(st); ok {
							notifies++
							_ = i
							_ = line
						}
					}
				}
				_ = s.BeginCapture()
			}
			endCapture := func() {
				if s.CaptureEnabled() {
					s.EndCapture()
				}
			}

			for round := 0; round < 1+rng.Intn(10); round++ {
				for j := 0; j < 1+rng.Intn(4); j++ {
					recorded = append(recorded, RecordedStep{
						Action:   "click",
						Selector: fmt.Sprintf("#%d-%d", round, j),
					})
					totalAppended++
				}
				beginCapture()
				switch rng.Intn(4) {
				case 0:
					s.Pause()
				case 1:
					s.Resume()
				case 2, 3:
					endCapture()
				}
			}

			if notifies > totalAppended {
				t.Fatalf("notifies=%d exceeds appended=%d", notifies, totalAppended)
			}
		})
	}
}
