package gui

import (
	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/scenario"
)

func collectFeatureTags(store *scenario.FeatureStore, files []string) map[string][]string {
	out := make(map[string][]string, len(files))
	for _, file := range files {
		feature, err := store.Load(file)
		if err != nil {
			continue
		}
		tags := gherkin.CollectFeatureTags(feature)
		if len(tags) > 0 {
			out[file] = tags
		}
	}
	return out
}

func collectProjectTags(store *scenario.FeatureStore, files []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, 16)
	for _, file := range files {
		feature, err := store.Load(file)
		if err != nil {
			continue
		}
		for _, tag := range gherkin.CollectFeatureTags(feature) {
			if _, ok := seen[tag]; ok {
				continue
			}
			seen[tag] = struct{}{}
			out = append(out, tag)
		}
	}
	return out
}
