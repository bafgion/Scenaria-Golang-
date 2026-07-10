package player

import "github.com/bafgion/scenaria-golang/internal/runstatus"

// RunstatusEntry builds a run_status.json row from a scenario result.
func RunstatusEntry(result ScenarioResult, runner string) runstatus.Entry {
	entry := runstatus.Entry{
		Path:       result.FeaturePath + "::" + result.Scenario,
		CaseID:     result.CaseID,
		RunID:      result.RunID,
		Success:    ResultStatusIsSuccessful(result.Status),
		Status:     result.Status,
		Message:    result.Message,
		Runner:     runner,
		DurationMS: int(result.DurationMS),
	}
	if result.ExampleIndex > 0 {
		idx := result.ExampleIndex
		entry.ExampleIndex = &idx
	}
	if result.FailedStep != nil {
		entry.FailedStep = result.FailedStep
	}
	if durs := stepDurationsFromRecords(result.StepRecords); len(durs) > 0 {
		entry.StepDurations = durs
	}
	return entry
}

func stepDurationsFromRecords(records []StepRecord) []int {
	if len(records) == 0 {
		return nil
	}
	max := 0
	for _, r := range records {
		if r.Index > max {
			max = r.Index
		}
	}
	out := make([]int, max+1)
	for _, r := range records {
		if r.Index >= 0 && r.Index < len(out) {
			out[r.Index] = int(r.DurationMS)
		}
	}
	return out
}
