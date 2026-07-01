package cli

import "path/filepath"

func dedupePaths(paths []string) []string {
	if len(paths) <= 1 {
		return paths
	}
	seen := make(map[string]struct{}, len(paths))
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		key := filepath.Clean(p)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, p)
	}
	return out
}
