package gui

import (
	"fmt"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/scenario"
)

// FailedStepLine resolves a 0-based leaf step index to a 1-based editor line.
func (s *Service) FailedStepLine(featurePath, scenarioName string, leafIndex int) (int, error) {
	featurePath = strings.TrimSpace(featurePath)
	if featurePath == "" {
		return 0, fmt.Errorf("feature path is required")
	}
	if leafIndex < 0 {
		return 0, fmt.Errorf("invalid step index")
	}
	store := scenario.NewFeatureStore()
	feature, err := store.Load(featurePath)
	if err != nil {
		return 0, err
	}
	scenarioName = strings.TrimSpace(scenarioName)
	for _, sc := range feature.Scenarios {
		if scenarioName != "" && sc.Title != scenarioName {
			continue
		}
		line, ok := gherkin.LeafStepLineAtIndex(sc.Steps, leafIndex)
		if ok {
			return line, nil
		}
	}
	return 0, fmt.Errorf("step %d not found in %s", leafIndex, featurePath)
}
