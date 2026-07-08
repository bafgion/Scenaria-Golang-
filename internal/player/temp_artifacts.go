package player

import (
	"os"
	"path/filepath"
	"strings"
)

func managedTempArtifactsRoot() string {
	return filepath.Join(os.TempDir(), "scenaria-artifacts")
}

func writeTempScreenshotPNG(payload []byte, prefix string) (path string, ok bool) {
	if len(payload) == 0 {
		return "", false
	}
	dir := filepath.Join(managedTempArtifactsRoot(), "screenshots")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", false
	}
	name := strings.TrimSpace(prefix)
	if name == "" {
		name = "screenshot"
	}
	tmp, err := os.CreateTemp(dir, name+"-*.png")
	if err != nil {
		return "", false
	}
	tmpPath := tmp.Name()
	if _, err := tmp.Write(payload); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
		return "", false
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return "", false
	}
	return tmpPath, true
}

func isManagedTempArtifact(path string) bool {
	path = strings.TrimSpace(path)
	if path == "" {
		return false
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	root, err := filepath.Abs(managedTempArtifactsRoot())
	if err != nil {
		return false
	}
	if absPath == root {
		return true
	}
	rel, err := filepath.Rel(root, absPath)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

// CleanupExecutionTempArtifacts removes managed temporary screenshot files after report generation.
func CleanupExecutionTempArtifacts(result *ExecutionResult) {
	if result == nil || len(result.ScenarioResults) == 0 {
		return
	}
	for si := range result.ScenarioResults {
		scenario := &result.ScenarioResults[si]
		if isManagedTempArtifact(scenario.ScreenshotPath) {
			_ = os.Remove(scenario.ScreenshotPath)
		}
		for ri := range scenario.StepRecords {
			step := &scenario.StepRecords[ri]
			if isManagedTempArtifact(step.ScreenshotPath) {
				_ = os.Remove(step.ScreenshotPath)
			}
		}
	}
}
