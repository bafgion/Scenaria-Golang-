package player

import (
	"context"
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestExecuteStepsSkipsTestClientDeclaration(t *testing.T) {
	exec := NewStepExecutor(ExecutorOptions{})
	err := exec.ExecuteSteps(context.Background(), &browserSession{}, []gherkin.Step{{
		Line: 4,
		Text: `я подключаю TestClient "DemoUser"`,
	}}, nil)
	if err != nil {
		t.Fatalf("TestClient declaration should be skipped, got: %v", err)
	}
}

func TestExecuteIfConditionReturnsSelectorErrors(t *testing.T) {
	ctx := context.Background()
	pw, stopPW, err := startPlaywright(ctx)
	if err != nil {
		t.Skip("playwright not available:", err)
	}
	defer stopPW()
	session, err := newBrowserSession(pw, PlaywrightExecutorOptions{BrowserName: "chromium", Headless: true})
	if err != nil {
		t.Skip("browser not available:", err)
	}
	defer session.close()

	exec := NewStepExecutor(ExecutorOptions{})
	err = exec.ExecuteSteps(ctx, session, []gherkin.Step{{
		Line:  7,
		Block: gherkin.BlockIf,
		Condition: &gherkin.Condition{
			Type:     "visible",
			Selector: "??",
		},
	}}, NewRunContext(nil, 1, t.TempDir()))
	if err == nil {
		t.Fatal("expected invalid selector condition to fail")
	}
	if !strings.Contains(err.Error(), "evaluate visible condition") {
		t.Fatalf("unexpected condition error: %v", err)
	}
}
