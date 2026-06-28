package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunVADryRun(t *testing.T) {
	tmp := t.TempDir()
	feature := filepath.Join(tmp, "demo.feature")
	if err := os.WriteFile(feature, []byte(`Функционал: Demo
Сценарий: A
  Допустим открыт "https://example.com"`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := RunVA([]string{"run", "--project", tmp, "--dry-run"}); err != nil {
		t.Fatalf("RunVA dry-run failed: %v", err)
	}
}
