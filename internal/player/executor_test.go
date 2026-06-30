package player

import (
	"context"
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
