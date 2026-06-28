package vanessa

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFailedScenariosFromJUnit(t *testing.T) {
	tmp := t.TempDir()
	junitDir := filepath.Join(tmp, "junit")
	if err := os.MkdirAll(junitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	payload := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite tests="1">
  <testcase classname="features/demo.feature" name="Scenario A">
    <failure message="assert failed"/>
  </testcase>
</testsuite>`
	if err := os.WriteFile(filepath.Join(junitDir, "result.xml"), []byte(payload), 0o644); err != nil {
		t.Fatal(err)
	}
	failed, err := FailedScenariosFromJUnit(junitDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(failed) != 1 || failed[0].ScenarioName != "Scenario A" {
		t.Fatalf("unexpected failed scenarios: %#v", failed)
	}
}

func TestBuildRerunRequest(t *testing.T) {
	tmp := t.TempDir()
	junitDir := filepath.Join(tmp, "junit")
	_ = os.MkdirAll(junitDir, 0o755)
	payload := `<?xml version="1.0"?><testsuite><testcase classname="a.feature" name="S1"><failure/></testcase></testsuite>`
	_ = os.WriteFile(filepath.Join(junitDir, "r.xml"), []byte(payload), 0o644)
	req, err := BuildRerunRequest(RunRequest{Paths: []string{"fallback"}}, tmp)
	if err != nil {
		t.Fatal(err)
	}
	if req == nil || len(req.ScenarioNames) != 1 || req.ScenarioNames[0] != "S1" {
		t.Fatalf("unexpected rerun request: %#v", req)
	}
}
