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

func TestRunRun_StubEngineNotImplemented(t *testing.T) {
	tmp := t.TempDir()
	featurePath := filepath.Join(tmp, "ok.feature")
	content := "Функционал: Demo\nСценарий: S1\nКогда выполняю шаг\n"
	if err := os.WriteFile(featurePath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write feature: %v", err)
	}

	err := RunRun([]string{tmp, "--engine", "stub"})
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

func TestRunRun_WritesJUnit(t *testing.T) {
	tmp := t.TempDir()
	featurePath := filepath.Join(tmp, "ok.feature")
	content := "Функционал: Demo\nСценарий: S1\nКогда выполняю шаг\n"
	if err := os.WriteFile(featurePath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write feature: %v", err)
	}
	junitPath := filepath.Join(tmp, "junit.xml")

	if err := RunRun([]string{tmp, "--dry-run", "--junit", junitPath}); err != nil {
		t.Fatalf("RunRun returned error: %v", err)
	}
	if _, err := os.Stat(junitPath); err != nil {
		t.Fatalf("expected junit report to be written: %v", err)
	}
}

func TestRunRun_WritesAllure(t *testing.T) {
	tmp := t.TempDir()
	featurePath := filepath.Join(tmp, "ok.feature")
	content := "Функционал: Demo\nСценарий: S1\nКогда выполняю шаг\n"
	if err := os.WriteFile(featurePath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write feature: %v", err)
	}
	allureDir := filepath.Join(tmp, "allure-results")
	if err := RunRun([]string{tmp, "--dry-run", "--allure", allureDir}); err != nil {
		t.Fatalf("RunRun returned error: %v", err)
	}
	entries, err := os.ReadDir(allureDir)
	if err != nil || len(entries) == 0 {
		t.Fatalf("expected allure results in %s: %v (%d files)", allureDir, err, len(entries))
	}
}

func TestParseRunOptions(t *testing.T) {
	opts, err := parseRunOptions([]string{
		"./features",
		"--dry-run",
		"--summary-json", "result.json",
		"--junit", "junit.xml",
		"--engine", "playwright",
		"--browser", "firefox",
		"--headed",
		"--base-url", "https://example.local",
		"--install-playwright",
		"--tag", "smoke",
		"--test-client", "DemoUser",
		"--var", "BASE=https://example.com",
		"--allure", "allure-results",
		"--trace", "traces",
		"--video", "videos",
	})
	if err != nil {
		t.Fatalf("parseRunOptions returned error: %v", err)
	}
	if len(opts.targets) != 1 || opts.targets[0] != "./features" || !opts.dryRun || opts.summaryJSON != "result.json" {
		t.Fatalf("unexpected options: %+v", opts)
	}
	if opts.junitPath != "junit.xml" || opts.engine != "playwright" || opts.browser != "firefox" || !opts.headed || opts.baseURL != "https://example.local" || !opts.installPlaywright {
		t.Fatalf("unexpected junit path: %+v", opts)
	}
	if opts.tag != "smoke" {
		t.Fatalf("unexpected tag: %q", opts.tag)
	}
	if opts.testClient != "DemoUser" {
		t.Fatalf("unexpected test client: %q", opts.testClient)
	}
	if opts.variables["BASE"] != "https://example.com" {
		t.Fatalf("unexpected variables: %#v", opts.variables)
	}
	if opts.allureDir != "allure-results" || opts.traceDir != "traces" || opts.videoDir != "videos" {
		t.Fatalf("unexpected report dirs: %+v", opts)
	}
}

func TestRunRun_TagFilterNoMatch(t *testing.T) {
	tmp := t.TempDir()
	featurePath := filepath.Join(tmp, "plain.feature")
	content := "Функционал: Demo\nСценарий: S1\nКогда выполняю шаг\n"
	if err := os.WriteFile(featurePath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write feature: %v", err)
	}

	err := RunRun([]string{tmp, "--dry-run", "--tag", "smoke"})
	if err == nil {
		t.Fatal("expected tag filter error")
	}
}

func TestRunRun_TagFilterMatch(t *testing.T) {
	tmp := t.TempDir()
	featurePath := filepath.Join(tmp, "smoke.feature")
	content := "@smoke\nФункционал: Demo\nСценарий: S1\nКогда открываю \"https://example.com\"\n"
	if err := os.WriteFile(featurePath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write feature: %v", err)
	}

	if err := RunRun([]string{tmp, "--dry-run", "--tag", "smoke"}); err != nil {
		t.Fatalf("RunRun returned error: %v", err)
	}
}

func TestRunRun_ExpandsOutlineDryRun(t *testing.T) {
	tmp := t.TempDir()
	featurePath := filepath.Join(tmp, "outline.feature")
	content := `Функционал: Каталог
Структура сценария: Поиск
  Когда открываю "<url>"
  Тогда вижу "<title>"
Примеры:
  | url      | title |
  | /catalog | Items |
  | /offers  | Deals |
`
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
	if int(decoded["scenarios"].(float64)) != 2 {
		t.Fatalf("expected 2 expanded scenarios, got %v", decoded["scenarios"])
	}
}

func TestParseRunOptionsErrors(t *testing.T) {
	if _, err := parseRunOptions(nil); err == nil {
		t.Fatal("expected usage error for empty args")
	}
	if _, err := parseRunOptions([]string{"./features", "--summary-json"}); err == nil {
		t.Fatal("expected missing value error for --summary-json")
	}
	if _, err := parseRunOptions([]string{"./features", "--junit"}); err == nil {
		t.Fatal("expected missing value error for --junit")
	}
	if _, err := parseRunOptions([]string{"./features", "--engine"}); err == nil {
		t.Fatal("expected missing value error for --engine")
	}
	if _, err := parseRunOptions([]string{"./features", "--browser"}); err == nil {
		t.Fatal("expected missing value error for --browser")
	}
	if _, err := parseRunOptions([]string{"./features", "--base-url"}); err == nil {
		t.Fatal("expected missing value error for --base-url")
	}
	if _, err := parseRunOptions([]string{"./features", "--unknown"}); err == nil {
		t.Fatal("expected unknown flag error")
	}
}

func TestBuildRunner(t *testing.T) {
	runner, err := buildRunner(runOptions{dryRun: true})
	if err != nil {
		t.Fatalf("buildRunner dry-run failed: %v", err)
	}
	if _, ok := runner.(player.DryRunner); !ok {
		t.Fatalf("expected DryRunner, got %T", runner)
	}

	runner, err = buildRunner(runOptions{engine: "stub"})
	if err != nil {
		t.Fatalf("buildRunner stub failed: %v", err)
	}
	if _, ok := runner.(player.BrowserRunner); !ok {
		t.Fatalf("expected BrowserRunner for stub, got %T", runner)
	}

	runner, err = buildRunner(runOptions{engine: "playwright", browser: "chromium"})
	if err != nil {
		t.Fatalf("buildRunner playwright failed: %v", err)
	}
	if _, ok := runner.(player.BrowserRunner); !ok {
		t.Fatalf("expected BrowserRunner for playwright, got %T", runner)
	}

	if _, err := buildRunner(runOptions{engine: "unknown"}); err == nil {
		t.Fatal("expected error for unknown engine")
	}
}
