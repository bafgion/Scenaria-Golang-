package vanessa

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FailedScenario struct {
	FeaturePath  string
	ScenarioName string
}

func FailedScenariosFromJUnit(junitDir string) ([]FailedScenario, error) {
	info, err := os.Stat(junitDir)
	if err != nil || !info.IsDir() {
		return nil, nil
	}
	seen := map[string]FailedScenario{}
	_ = filepath.WalkDir(junitDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".xml" {
			return nil
		}
		for _, failure := range parseFailedCases(path) {
			key := failure.FeaturePath + "\x00" + failure.ScenarioName
			seen[key] = failure
		}
		return nil
	})
	out := make([]FailedScenario, 0, len(seen))
	for _, item := range seen {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].FeaturePath == out[j].FeaturePath {
			return out[i].ScenarioName < out[j].ScenarioName
		}
		return out[i].FeaturePath < out[j].FeaturePath
	})
	return out, nil
}

func parseFailedCases(path string) []FailedScenario {
	payload, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var root junitNode
	if err := xml.Unmarshal(payload, &root); err != nil {
		return nil
	}
	suites := root.TestSuites
	if root.Tests > 0 || len(root.TestCases) > 0 {
		suites = []junitNode{root}
	}
	out := make([]FailedScenario, 0)
	for _, suite := range suites {
		for _, tc := range suite.TestCases {
			if tc.Failure == nil && tc.Error == nil {
				continue
			}
			featurePath := strings.TrimSpace(tc.Classname)
			if featurePath == "" {
				featurePath = path
			}
			out = append(out, FailedScenario{
				FeaturePath:  featurePath,
				ScenarioName: strings.TrimSpace(tc.Name),
			})
		}
	}
	return out
}

func BuildRerunRequest(base RunRequest, runDir string) (*RunRequest, error) {
	junitDir := filepath.Join(runDir, "junit")
	failed, err := FailedScenariosFromJUnit(junitDir)
	if err != nil {
		return nil, err
	}
	if len(failed) == 0 {
		return nil, nil
	}
	pathsSet := map[string]struct{}{}
	namesSet := map[string]struct{}{}
	for _, item := range failed {
		if item.FeaturePath != "" {
			pathsSet[item.FeaturePath] = struct{}{}
		}
		if item.ScenarioName != "" {
			namesSet[item.ScenarioName] = struct{}{}
		}
	}
	paths := make([]string, 0, len(pathsSet))
	for path := range pathsSet {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	if len(paths) == 0 {
		paths = append(paths, base.Paths...)
	}
	names := make([]string, 0, len(namesSet))
	for name := range namesSet {
		names = append(names, name)
	}
	sort.Strings(names)
	req := base
	req.Paths = paths
	req.ScenarioNames = names
	return &req, nil
}
