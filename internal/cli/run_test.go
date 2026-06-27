package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunRun_DryRunSuccess(t *testing.T) {
	tmp := t.TempDir()
	featurePath := filepath.Join(tmp, "ok.feature")
	content := "Функционал: Demo\nСценарий: S1\nКогда выполняю шаг\n"
	if err := os.WriteFile(featurePath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write feature: %v", err)
	}

	if err := RunRun([]string{tmp, "--dry-run"}); err != nil {
		t.Fatalf("RunRun returned error: %v", err)
	}
}

func TestRunRun_RequiresDryRunForNow(t *testing.T) {
	tmp := t.TempDir()
	featurePath := filepath.Join(tmp, "ok.feature")
	content := "Функционал: Demo\nСценарий: S1\nКогда выполняю шаг\n"
	if err := os.WriteFile(featurePath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write feature: %v", err)
	}

	if err := RunRun([]string{tmp}); err == nil {
		t.Fatal("expected error when --dry-run is missing")
	}
}
