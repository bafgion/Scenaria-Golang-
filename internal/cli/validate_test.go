package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunValidate_Success(t *testing.T) {
	tmp := t.TempDir()
	featurePath := filepath.Join(tmp, "ok.feature")
	content := "Функционал: Demo\nСценарий: S1\nКогда выполняю шаг\n"
	if err := os.WriteFile(featurePath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write feature: %v", err)
	}

	if err := RunValidate([]string{tmp}); err != nil {
		t.Fatalf("RunValidate returned error: %v", err)
	}
}

func TestRunValidate_Failure(t *testing.T) {
	tmp := t.TempDir()
	featurePath := filepath.Join(tmp, "bad.feature")
	content := "Функционал: Demo\nКогда вне сценария\n"
	if err := os.WriteFile(featurePath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write feature: %v", err)
	}

	if err := RunValidate([]string{tmp}); err == nil {
		t.Fatal("expected validation error, got nil")
	}
}
