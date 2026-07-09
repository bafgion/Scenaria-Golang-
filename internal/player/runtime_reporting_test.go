package player

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

func TestCloseBrowserSkipsRemainingSteps(t *testing.T) {
	session := &browserSession{}
	exec := NewStepExecutor(ExecutorOptions{})
	runCtx := NewRunContext(nil, 1, t.TempDir())
	steps := []gherkin.Step{
		{Text: "закрываю браузер", Line: 1, Keyword: "When"},
		{Text: "жду 1 мс", Line: 2, Keyword: "Then"},
	}
	err := exec.ExecuteSteps(t.Context(), session, steps, runCtx)
	if err == nil {
		t.Fatal("expected terminal close error")
	}
	var tc *terminalCloseError
	if !errors.As(err, &tc) {
		t.Fatalf("expected terminalCloseError, got %v", err)
	}
	if tc.remaining != 1 {
		t.Fatalf("expected 1 remaining step, got %d", tc.remaining)
	}
	records := runCtx.StepRecords()
	if len(records) != 2 {
		t.Fatalf("expected 2 step records, got %d", len(records))
	}
	if records[0].Status != "passed" || records[0].TerminalAction != "close-browser" {
		t.Fatalf("first step: %+v", records[0])
	}
	if records[1].Status != "skipped" {
		t.Fatalf("second step: %+v", records[1])
	}
}

func TestRepeatIterationPathRecorded(t *testing.T) {
	session := &browserSession{}
	exec := NewStepExecutor(ExecutorOptions{})
	runCtx := NewRunContext(nil, 1, t.TempDir())
	steps := []gherkin.Step{{
		Block:       gherkin.BlockRepeat,
		RepeatCount: 2,
		Children:    []gherkin.Step{{Text: "жду 0 мс", Line: 2, Keyword: "When"}},
	}}
	if err := exec.ExecuteSteps(t.Context(), session, steps, runCtx); err != nil {
		t.Fatal(err)
	}
	records := runCtx.StepRecords()
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if got := FormatIterationPath(records[0].IterationPath); got != "repeat[1]" {
		t.Fatalf("iteration 1 path: %q", got)
	}
	if got := FormatIterationPath(records[1].IterationPath); got != "repeat[2]" {
		t.Fatalf("iteration 2 path: %q", got)
	}
}

func TestNegativeWaitDurationRejected(t *testing.T) {
	err := executeAction(context.Background(), &browserSession{}, stepdsl.Action{Kind: "wait", Value1: "-1ms"}, "", nil)
	if err == nil || !strings.Contains(err.Error(), "negative") {
		t.Fatalf("expected negative wait error, got %v", err)
	}
}

func TestZeroWaitDurationIsInstant(t *testing.T) {
	start := time.Now()
	if err := executeAction(context.Background(), &browserSession{}, stepdsl.Action{Kind: "wait", Value1: "0ms"}, "", nil); err != nil {
		t.Fatal(err)
	}
	if time.Since(start) > 50*time.Millisecond {
		t.Fatal("zero wait should return immediately")
	}
}

func TestMaxActionAttemptsExceeded(t *testing.T) {
	exec := NewStepExecutor(ExecutorOptions{MaxActionAttempts: 2})
	runCtx := NewRunContext(nil, 1, t.TempDir())
	steps := []gherkin.Step{
		{Text: "жду 0 мс", Line: 1},
		{Text: "жду 0 мс", Line: 2},
		{Text: "жду 0 мс", Line: 3},
	}
	err := exec.ExecuteSteps(t.Context(), &browserSession{}, steps, runCtx)
	if err == nil || !strings.Contains(err.Error(), "max action attempts") {
		t.Fatalf("expected max action attempts error, got %v", err)
	}
}
