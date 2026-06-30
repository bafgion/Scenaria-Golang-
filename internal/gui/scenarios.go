package gui

import (
	"sort"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/scenario"
)

func (s *Service) ListScenarioTitles() ([]string, error) {
	info, err := s.projectInfo()
	if err != nil {
		return nil, err
	}
	return collectScenarioTitles(scenario.NewFeatureStore(), info.Features), nil
}

func collectScenarioTitles(store *scenario.FeatureStore, files []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(files)*2)
	for _, file := range files {
		feature, err := store.Load(file)
		if err != nil {
			continue
		}
		for _, runnable := range gherkin.ExpandFeatureAtPath(feature, file) {
			title := strings.TrimSpace(runnable.Title)
			if title == "" {
				continue
			}
			if _, ok := seen[title]; ok {
				continue
			}
			seen[title] = struct{}{}
			out = append(out, title)
		}
	}
	sort.Strings(out)
	return out
}
