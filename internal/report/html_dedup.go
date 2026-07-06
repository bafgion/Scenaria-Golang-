package report

import (
	"crypto/sha256"
	"encoding/hex"
)

func dedupeScreenshotURLs(payload *htmlReportPayload) int {
	if payload == nil {
		return 0
	}
	refs := map[string]string{}
	saved := 0
	assign := func(url string) string {
		if url == "" || !isDataURL(url) {
			return url
		}
		sum := sha256.Sum256([]byte(url))
		key := hex.EncodeToString(sum[:8])
		if prev, ok := refs[key]; ok {
			saved++
			return prev
		}
		refs[key] = url
		return url
	}
	for i := range payload.Scenarios {
		sc := &payload.Scenarios[i]
		sc.Screenshot = assign(sc.Screenshot)
		for j := range sc.Steps {
			sc.Steps[j].Screenshot = assign(sc.Steps[j].Screenshot)
		}
	}
	return saved
}

func isDataURL(s string) bool {
	return len(s) > 5 && (hasPrefix(s, "data:image/") || hasPrefix(s, "data:application/"))
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
