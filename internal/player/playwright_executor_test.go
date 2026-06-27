package player

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

func TestLaunchBrowserUnsupportedName(t *testing.T) {
	_, err := launchBrowser(nil, PlaywrightExecutorOptions{BrowserName: "unknown"})
	if err == nil {
		t.Fatal("expected unsupported browser error")
	}
}

func TestExecuteActionWaitCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := executeAction(ctx, nil, stepdsl.Action{Kind: "wait", Value1: "5s"}, "")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancellation error, got: %v", err)
	}
}

func TestExecuteActionWaitSuccess(t *testing.T) {
	start := time.Now()
	if err := executeAction(context.Background(), nil, stepdsl.Action{Kind: "wait", Value1: "10ms"}, ""); err != nil {
		t.Fatalf("executeAction wait returned error: %v", err)
	}
	if time.Since(start) < 10*time.Millisecond {
		t.Fatal("wait action did not wait long enough")
	}
}

func TestExecuteActionInvalidWaitDuration(t *testing.T) {
	if err := executeAction(context.Background(), nil, stepdsl.Action{Kind: "wait", Value1: "oops"}, ""); err == nil {
		t.Fatal("expected invalid wait duration error")
	}
}
