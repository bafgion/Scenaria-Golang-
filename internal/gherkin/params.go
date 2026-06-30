package gherkin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ScenarioParamsFile is the sidecar JSON next to a .feature file (stem.params.json).
type ScenarioParamsFile struct {
	Scenarios map[string][]map[string]string `json:"scenarios"`
}

// ParamsPath returns the sidecar path for a feature file.
func ParamsPath(featurePath string) string {
	ext := filepath.Ext(featurePath)
	if ext == "" {
		return featurePath + ".params.json"
	}
	return strings.TrimSuffix(featurePath, ext) + ".params.json"
}

// LoadScenarioParams reads parameter rows for a scenario title from the feature sidecar.
func LoadScenarioParams(featurePath, scenarioTitle string) ([]map[string]string, error) {
	path := ParamsPath(featurePath)
	payload, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read params %q: %w", path, err)
	}
	var file ScenarioParamsFile
	if err := json.Unmarshal(payload, &file); err != nil {
		return nil, fmt.Errorf("decode params %q: %w", path, err)
	}
	if len(file.Scenarios) == 0 {
		return nil, nil
	}
	title := strings.TrimSpace(scenarioTitle)
	if rows, ok := file.Scenarios[title]; ok {
		return cloneParamRows(rows), nil
	}
	for key, rows := range file.Scenarios {
		if strings.EqualFold(strings.TrimSpace(key), title) {
			return cloneParamRows(rows), nil
		}
	}
	return nil, nil
}

func cloneParamRows(rows []map[string]string) []map[string]string {
	out := make([]map[string]string, len(rows))
	for i, row := range rows {
		copyRow := make(map[string]string, len(row))
		for k, v := range row {
			copyRow[k] = v
		}
		out[i] = copyRow
	}
	return out
}

func expandScenarioFromParamRows(feature *Feature, scenario Scenario, tags []string, rows []map[string]string) []RunnableScenario {
	out := make([]RunnableScenario, 0, len(rows))
	for rowIndex, values := range rows {
		steps := mergeBackgroundSteps(feature.Background, expandSteps(scenario.Steps, values))
		title := scenario.Title
		if sample := firstValuesSampleOrdered(values, placeholderKeysFromSteps(scenario.Steps)); sample != "" {
			title = scenario.Title + " — " + sample
		}
		out = append(out, RunnableScenario{
			Title:        title,
			Tags:         tags,
			Steps:        steps,
			ExampleIndex: rowIndex + 1,
		})
	}
	return out
}

