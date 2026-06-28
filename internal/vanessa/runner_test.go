package vanessa

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateSettingsMissingPaths(t *testing.T) {
	issues := ValidateSettings(DefaultSettings())
	if len(issues) == 0 {
		t.Fatal("expected validation issues for empty settings")
	}
}

func TestMergeVAParamsDry(t *testing.T) {
	tmp := t.TempDir()
	cfg := DefaultSettings()
	cfg.ReportJUnit = true
	runDir := filepath.Join(tmp, "run-1")
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		t.Fatal(err)
	}
	merged, path, err := MergeVAParams(cfg, RunRequest{ProjectRoot: tmp, Paths: []string{tmp}}, runDir)
	if err != nil {
		t.Fatalf("merge failed: %v", err)
	}
	if merged == nil || path == "" {
		t.Fatalf("unexpected merge result: %+v %q", merged, path)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("VAParams not written: %v", err)
	}
}

func TestParseJUnitFileEmpty(t *testing.T) {
	if cases := parseJUnitFile("missing.xml"); cases != nil {
		t.Fatalf("expected nil for missing file")
	}
}
