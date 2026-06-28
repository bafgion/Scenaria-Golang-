package gherkin

import "strings"

func NormalizeTag(tag string) string {
	return strings.ToLower(strings.TrimSpace(strings.TrimPrefix(tag, "@")))
}

func TagsInclude(tags []string, tag string) bool {
	normalized := NormalizeTag(tag)
	if normalized == "" {
		return true
	}
	for _, item := range tags {
		if NormalizeTag(item) == normalized {
			return true
		}
	}
	return false
}

func MergeTags(groups ...[]string) []string {
	seen := make(map[string]struct{})
	out := make([]string, 0)
	for _, group := range groups {
		for _, tag := range group {
			normalized := strings.TrimSpace(tag)
			if normalized == "" {
				continue
			}
			key := NormalizeTag(normalized)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			out = append(out, normalized)
		}
	}
	return out
}

func FeatureHasTag(feature *Feature, tag string) bool {
	if feature == nil {
		return false
	}
	if TagsInclude(feature.Tags, tag) {
		return true
	}
	for _, scenario := range feature.Scenarios {
		if TagsInclude(scenario.Tags, tag) {
			return true
		}
	}
	return false
}

func CollectFeatureTags(feature *Feature) []string {
	if feature == nil {
		return nil
	}
	tags := append([]string{}, feature.Tags...)
	for _, scenario := range feature.Scenarios {
		tags = append(tags, scenario.Tags...)
	}
	return MergeTags(tags)
}
