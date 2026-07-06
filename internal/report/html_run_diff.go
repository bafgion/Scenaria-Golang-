package report

import "strings"

type htmlRunDiff struct {
	Columns []htmlRunDiffCol `json:"columns"`
	Rows    []htmlRunDiffRow `json:"rows"`
}

type htmlRunDiffCol struct {
	Label   string `json:"label"`
	At      string `json:"at,omitempty"`
	Status  string `json:"status"`
	Current bool   `json:"current,omitempty"`
}

type htmlRunDiffRow struct {
	Index int      `json:"index"`
	Text  string   `json:"text"`
	Cells []string `json:"cells"`
}

type runSnapshot struct {
	status     string
	failedStep *int
	label      string
	at         string
	current    bool
}

func computeRunDiff(sc htmlScenario) *htmlRunDiff {
	if len(sc.Steps) == 0 || len(sc.HistoryRuns) == 0 {
		return nil
	}
	snaps := make([]runSnapshot, 0, len(sc.HistoryRuns)+1)
	runs := sc.HistoryRuns
	if len(runs) > 4 {
		runs = runs[:4]
	}
	for i := len(runs) - 1; i >= 0; i-- {
		r := runs[i]
		snaps = append(snaps, runSnapshot{
			status: r.Status, failedStep: r.FailedStep,
			label: shortRunLabel(r.At), at: r.At,
		})
	}
	curFailed := sc.FailedStep
	for _, st := range sc.Steps {
		if st.Status == "failed" {
			idx := st.Index
			curFailed = &idx
			break
		}
	}
	snaps = append(snaps, runSnapshot{
		status: sc.Status, failedStep: curFailed,
		label: "now", current: true,
	})
	if len(snaps) < 2 {
		return nil
	}

	diff := &htmlRunDiff{
		Columns: make([]htmlRunDiffCol, 0, len(snaps)),
		Rows:    make([]htmlRunDiffRow, 0, len(sc.Steps)),
	}
	for _, s := range snaps {
		diff.Columns = append(diff.Columns, htmlRunDiffCol{
			Label: s.label, At: s.at, Status: s.status, Current: s.current,
		})
	}
	for _, step := range sc.Steps {
		row := htmlRunDiffRow{Index: step.Index, Text: step.Text}
		if row.Text == "" {
			row.Text = step.Gherkin
		}
		for _, snap := range snaps {
			if snap.current {
				row.Cells = append(row.Cells, step.Status)
			} else {
				row.Cells = append(row.Cells, inferStepStatus(snap.status, step.Index, snap.failedStep))
			}
		}
		diff.Rows = append(diff.Rows, row)
	}
	return diff
}

func shortRunLabel(at string) string {
	at = strings.TrimSpace(at)
	if len(at) >= 16 {
		return at[5:16]
	}
	if len(at) >= 10 {
		return at[5:10]
	}
	return at
}

func buildStepDurationSparkline(sc htmlScenario, stepIndex int, currentMS int64) []int {
	vals := make([]int, 0, len(sc.HistoryRuns)+1)
	runs := sc.HistoryRuns
	for i := len(runs) - 1; i >= 0; i-- {
		if durs := runs[i].StepDurations; stepIndex < len(durs) && durs[stepIndex] > 0 {
			vals = append(vals, durs[stepIndex])
		}
	}
	if currentMS > 0 {
		vals = append(vals, int(currentMS))
	}
	if len(vals) < 2 {
		return nil
	}
	return vals
}

func attachStepSparklines(sc *htmlScenario) {
	if sc == nil || len(sc.Steps) == 0 {
		return
	}
	for i := range sc.Steps {
		spark := buildStepDurationSparkline(*sc, sc.Steps[i].Index, sc.Steps[i].DurationMS)
		if len(spark) > 0 {
			sc.Steps[i].DurationSparkline = spark
		}
	}
}
