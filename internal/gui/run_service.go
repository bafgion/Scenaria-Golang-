package gui

import "github.com/bafgion/scenaria-golang/internal/runstatus"

// RunService owns run history/flaky metrics read operations.
type RunService struct {
	projectPath projectPathProvider
}

func NewRunService(projectPath projectPathProvider) *RunService {
	return &RunService{projectPath: projectPath}
}

func (s *RunService) ListRunResults(limit int) ([]RunResultEntry, error) {
	if s == nil || s.projectPath == nil {
		return []RunResultEntry{}, nil
	}
	path := s.projectPath()
	if path == "" {
		return []RunResultEntry{}, nil
	}
	store, err := runstatus.Open(path)
	if err != nil {
		return nil, err
	}
	entries, err := store.List(limit)
	if err != nil {
		return nil, err
	}
	out := make([]RunResultEntry, 0, len(entries))
	for _, e := range entries {
		out = append(out, RunResultEntry{
			Path:       e.Path,
			Success:    e.Success,
			Message:    e.Message,
			Runner:     e.Runner,
			At:         e.At,
			FailedStep: e.FailedStep,
		})
	}
	return out, nil
}

func (s *RunService) FlakyMetrics(historyLimit int) (FlakyMetricsDTO, error) {
	if s == nil || s.projectPath == nil {
		return FlakyMetricsDTO{}, nil
	}
	path := s.projectPath()
	if path == "" {
		return FlakyMetricsDTO{}, nil
	}
	store, err := runstatus.Open(path)
	if err != nil {
		return FlakyMetricsDTO{}, err
	}
	if historyLimit <= 0 {
		historyLimit = 200
	}
	entries, err := store.List(historyLimit)
	if err != nil {
		return FlakyMetricsDTO{}, err
	}
	scenarios, steps := runstatus.FlakyStats(entries)
	out := FlakyMetricsDTO{
		Scenarios: make([]FlakyScenarioDTO, 0, len(scenarios)),
		Steps:     make([]FlakyStepDTO, 0, len(steps)),
	}
	for _, item := range scenarios {
		out.Scenarios = append(out.Scenarios, FlakyScenarioDTO{
			Path:       item.Path,
			Failures:   item.Failures,
			Passes:     item.Passes,
			Total:      item.Total,
			Flaky:      item.Flaky,
			LastFailed: item.LastFailed,
		})
	}
	for _, item := range steps {
		out.Steps = append(out.Steps, FlakyStepDTO{
			Path:       item.Path,
			Step:       item.Step,
			Failures:   item.Failures,
			LastFailed: item.LastFailed,
		})
	}
	return out, nil
}
