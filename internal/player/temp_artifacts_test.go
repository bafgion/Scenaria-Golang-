package player

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCleanupExecutionTempArtifactsRemovesManagedScreenshots(t *testing.T) {
	root := managedTempArtifactsRoot()
	if err := os.MkdirAll(filepath.Join(root, "screenshots"), 0o755); err != nil {
		t.Fatal(err)
	}
	scenarioShot := filepath.Join(root, "screenshots", "scenario-cleanup-test.png")
	stepShot := filepath.Join(root, "screenshots", "step-cleanup-test.png")
	if err := os.WriteFile(scenarioShot, []byte("png"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(stepShot, []byte("png"), 0o644); err != nil {
		t.Fatal(err)
	}
	result := ExecutionResult{
		ScenarioResults: []ScenarioResult{{
			ScreenshotPath: scenarioShot,
			StepRecords: []StepRecord{{
				ScreenshotPath: stepShot,
			}},
		}},
	}
	CleanupExecutionTempArtifacts(&result)
	if _, err := os.Stat(scenarioShot); !os.IsNotExist(err) {
		t.Fatalf("expected scenario screenshot removed, stat err=%v", err)
	}
	if _, err := os.Stat(stepShot); !os.IsNotExist(err) {
		t.Fatalf("expected step screenshot removed, stat err=%v", err)
	}
}
