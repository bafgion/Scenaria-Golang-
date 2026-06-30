package selector

import "strings"

var (
	DefaultClickStrategies = []string{"testid", "id", "aria", "contextual", "text"}
	DefaultInputStrategies = []string{"testid", "id", "label", "placeholder", "aria", "name"}
)

var allowedClickStrategies = map[string]struct{}{
	"testid": {}, "id": {}, "aria": {}, "contextual": {}, "text": {},
}

var allowedInputStrategies = map[string]struct{}{
	"testid": {}, "id": {}, "label": {}, "placeholder": {}, "aria": {}, "name": {},
}

func NormalizeClickStrategies(values []string) []string {
	return normalizeStrategies(values, DefaultClickStrategies, allowedClickStrategies)
}

func NormalizeInputStrategies(values []string) []string {
	return normalizeStrategies(values, DefaultInputStrategies, allowedInputStrategies)
}

func normalizeStrategies(values, defaults []string, allowed map[string]struct{}) []string {
	if len(values) == 0 {
		out := make([]string, len(defaults))
		copy(out, defaults)
		return out
	}
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, raw := range values {
		key := strings.ToLower(strings.TrimSpace(raw))
		if key == "" {
			continue
		}
		if _, ok := allowed[key]; !ok {
			continue
		}
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, key)
	}
	for _, key := range defaults {
		if _, ok := seen[key]; ok {
			continue
		}
		out = append(out, key)
	}
	return out
}
