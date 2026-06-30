package gui

import (
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

type RunFromLineDTO struct {
	Scenario  string `json:"scenario"`
	StartStep int    `json:"startStep"`
	EndStep   int    `json:"endStep"`
	Partial   bool   `json:"partial"`
}

// ResolveRunFromLine maps an editor line to scenario name and optional leaf step index (run from step).
func ResolveRunFromLine(text string, line int) (RunFromLineDTO, error) {
	return resolveRunAtLine(text, line, true)
}

// ResolveRunToLine maps an editor line to scenario name and end leaf index (run through step).
func ResolveRunToLine(text string, line int) (RunFromLineDTO, error) {
	return resolveRunAtLine(text, line, false)
}

func resolveRunAtLine(text string, line int, fromStep bool) (RunFromLineDTO, error) {
	out := RunFromLineDTO{StartStep: -1, EndStep: -1}
	feature, err := gherkin.ParseFeature(text)
	if err != nil {
		return out, err
	}
	sc, ok := gherkin.ScenarioContainingLine(feature, line)
	if !ok {
		return out, nil
	}
	out.Scenario = sc.Title
	if line == sc.Line {
		return out, nil
	}
	runnable := firstRunnableForScenario(feature, sc)
	if runnable == nil {
		return out, nil
	}
	idx, ok := gherkin.LeafStepIndexAtLine(runnable.Steps, line)
	if !ok {
		return out, nil
	}
	out.Partial = true
	if fromStep {
		out.StartStep = idx
	} else {
		out.EndStep = idx
	}
	return out, nil
}

func firstRunnableForScenario(feature *gherkin.Feature, sc gherkin.Scenario) *gherkin.RunnableScenario {
	for _, runnable := range gherkin.ExpandScenario(feature, sc) {
		if strings.EqualFold(strings.TrimSpace(runnable.Title), strings.TrimSpace(sc.Title)) {
			return &runnable
		}
	}
	expanded := gherkin.ExpandScenario(feature, sc)
	if len(expanded) == 0 {
		return nil
	}
	first := expanded[0]
	return &first
}
