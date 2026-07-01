package runstatus

// ScenarioFlakyStat summarizes recent pass/fail history for one scenario path.
type ScenarioFlakyStat struct {
	Path       string `json:"path"`
	Failures   int    `json:"failures"`
	Passes     int    `json:"passes"`
	Total      int    `json:"total"`
	Flaky      bool   `json:"flaky"`
	LastFailed string `json:"last_failed_at,omitempty"`
}

// StepFlakyStat summarizes failures at a specific step index within a scenario.
type StepFlakyStat struct {
	Path       string `json:"path"`
	Step       int    `json:"step"`
	Failures   int    `json:"failures"`
	LastFailed string `json:"last_failed_at,omitempty"`
}

// FlakyStats derives scenario- and step-level flaky metrics from run history.
// A scenario is flaky when it has both passes and failures in the sampled window.
func FlakyStats(entries []Entry) (scenarios []ScenarioFlakyStat, steps []StepFlakyStat) {
	type agg struct {
		failures   int
		passes     int
		lastFailed string
	}
	byPath := map[string]*agg{}
	type stepAgg struct {
		failures   int
		lastFailed string
	}
	byStep := map[string]map[int]*stepAgg{}

	for _, entry := range entries {
		path := entry.Path
		if path == "" {
			continue
		}
		a := byPath[path]
		if a == nil {
			a = &agg{}
			byPath[path] = a
		}
		if entry.Success {
			a.passes++
		} else {
			a.failures++
			if entry.At != "" {
				a.lastFailed = entry.At
			}
		}
		if entry.FailedStep != nil && !entry.Success {
			step := *entry.FailedStep
			if byStep[path] == nil {
				byStep[path] = map[int]*stepAgg{}
			}
			sa := byStep[path][step]
			if sa == nil {
				sa = &stepAgg{}
				byStep[path][step] = sa
			}
			sa.failures++
			if entry.At != "" {
				sa.lastFailed = entry.At
			}
		}
	}

	scenarios = make([]ScenarioFlakyStat, 0, len(byPath))
	for path, a := range byPath {
		total := a.failures + a.passes
		scenarios = append(scenarios, ScenarioFlakyStat{
			Path:       path,
			Failures:   a.failures,
			Passes:     a.passes,
			Total:      total,
			Flaky:      a.failures > 0 && a.passes > 0,
			LastFailed: a.lastFailed,
		})
	}
	sortScenarioFlaky(scenarios)

	for path, stepsMap := range byStep {
		for step, sa := range stepsMap {
			if sa.failures < 2 {
				continue
			}
			steps = append(steps, StepFlakyStat{
				Path:       path,
				Step:       step,
				Failures:   sa.failures,
				LastFailed: sa.lastFailed,
			})
		}
	}
	sortStepFlaky(steps)
	return scenarios, steps
}

func sortScenarioFlaky(items []ScenarioFlakyStat) {
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if scenarioFlakyLess(items[j], items[i]) {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}

func scenarioFlakyLess(a, b ScenarioFlakyStat) bool {
	if a.Flaky != b.Flaky {
		return a.Flaky
	}
	if a.Failures != b.Failures {
		return a.Failures > b.Failures
	}
	return a.Path < b.Path
}

func sortStepFlaky(items []StepFlakyStat) {
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if stepFlakyLess(items[j], items[i]) {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}

func stepFlakyLess(a, b StepFlakyStat) bool {
	if a.Failures != b.Failures {
		return a.Failures > b.Failures
	}
	if a.Path != b.Path {
		return a.Path < b.Path
	}
	return a.Step < b.Step
}
