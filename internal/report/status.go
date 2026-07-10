package report

import "strings"

type StatusTotals struct {
	Passed     int
	Failed     int
	Skipped    int
	Canceled   int
	NotStarted int
}

func classifyScenarioStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "passed", "pass", "ok", "success":
		return "passed"
	case "failed", "fail", "broken", "error":
		return "failed"
	case "canceled", "cancelled", "aborted":
		return "canceled"
	case "not-started", "not started", "notstarted", "not-start", "notstart", "not run", "not-run":
		return "not-started"
	case "skipped", "skip", "dry-run", "dryrun":
		return "skipped"
	default:
		return "skipped"
	}
}

func addStatusCount(t *StatusTotals, status string) {
	switch classifyScenarioStatus(status) {
	case "passed":
		t.Passed++
	case "failed":
		t.Failed++
	case "canceled":
		t.Canceled++
	case "not-started":
		t.NotStarted++
	default:
		t.Skipped++
	}
}

func statusTotalsSum(t StatusTotals) int {
	return t.Passed + t.Failed + t.Skipped + t.Canceled + t.NotStarted
}
