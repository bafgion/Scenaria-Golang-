package cli

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/player"
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

	err := RunRun([]string{tmp})
	if !errors.Is(err, player.ErrBrowserExecutionNotImplemented) {
		t.Fatalf("expected not-implemented error, got: %v", err)
	}
}

func TestRunRun_WritesSummaryJSON(t *testing.T) {
	tmp := t.TempDir()
	featurePath := filepath.Join(tmp, "ok.feature")
	content := "Функционал: Demo\nСценарий: S1\nКогда выполняю шаг\n"
	if err := os.WriteFile(featurePath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write feature: %v", err)
	}
	summaryPath := filepath.Join(tmp, "summary.json")

	if err := RunRun([]string{tmp, "--dry-run", "--summary-json", summaryPath}); err != nil {
		t.Fatalf("RunRun returned error: %v", err)
	}

	payload, err := os.ReadFile(summaryPath)
	if err != nil {
		t.Fatalf("failed to read summary json: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("failed to decode summary json: %v", err)
	}
	if decoded["mode"] != "dry-run" {
		t.Fatalf("unexpected mode in summary: %v", decoded["mode"])
	}
}

func TestParseRunOptions(t *testing.T) {
	opts, err := parseRunOptions([]string{"./features", "--dry-run", "--summary-json", "result.json"})
	if err != nil {
		t.Fatalf("parseRunOptions returned error: %v", err)
	}
	if opts.target != "./features" || !opts.dryRun || opts.summaryJSON != "result.json" {
		t.Fatalf("unexpected options: %+v", opts)
	}
}

func TestParseRunOptionsErrors(t *testing.T) {
	if _, err := parseRunOptions(nil); err == nil {
		t.Fatal("expected usage error for empty args")
	}
	if _, err := parseRunOptions([]string{"./features", "--summary-json"}); err == nil {
		t.Fatal("expected missing value error for --summary-json")
	}
	if _, err := parseRunOptions([]string{"./features", "--unknown"}); err == nil {
		t.Fatal("expected unknown flag error")
	}
}
