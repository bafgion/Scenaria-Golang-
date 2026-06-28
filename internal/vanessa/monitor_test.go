package vanessa

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunMonitorPollIncremental(t *testing.T) {
	tmp := t.TempDir()
	junitDir := filepath.Join(tmp, "junit")
	_ = os.MkdirAll(junitDir, 0o755)
	monitor := NewRunMonitor(tmp, 2)

	if snap := monitor.Poll(); len(snap.Cases) != 0 {
		t.Fatalf("expected no cases initially")
	}

	payload := `<?xml version="1.0"?><testsuite><testcase name="S1"><failure/></testcase></testsuite>`
	if err := os.WriteFile(filepath.Join(junitDir, "a.xml"), []byte(payload), 0o644); err != nil {
		t.Fatal(err)
	}
	snap := monitor.Poll()
	if len(snap.Cases) != 1 {
		t.Fatalf("expected 1 case, got %d", len(snap.Cases))
	}
	if snap2 := monitor.Poll(); len(snap2.Cases) != 0 {
		t.Fatalf("expected no duplicate cases on second poll")
	}
}

func TestReadCurrentScenarioLabel(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "scenario.log")
	content := "line1\nСценарий: Login flow\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	got := readCurrentScenarioLabel(path)
	if got == "" {
		t.Fatal("expected scenario label")
	}
}
