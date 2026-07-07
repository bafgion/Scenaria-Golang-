package stepdsl

import (
	"fmt"
	"strings"
)

// NormalizeDataTextHTMLURL percent-encodes non-ASCII payloads in data:text/html URLs
// and ensures charset=utf-8 so browsers render Cyrillic and other Unicode correctly.
func NormalizeDataTextHTMLURL(raw string) string {
	trimmed := strings.TrimSpace(raw)
	lower := strings.ToLower(trimmed)
	if !strings.HasPrefix(lower, "data:text/html") {
		return trimmed
	}
	comma := strings.Index(trimmed, ",")
	if comma < 0 {
		return trimmed
	}
	meta := trimmed[:comma]
	payload := trimmed[comma+1:]
	metaLower := strings.ToLower(meta)
	hasCharset := strings.Contains(metaLower, "charset=")

	if hasCharset && !payloadNeedsPercentEncode(payload) {
		return trimmed
	}
	if !hasCharset && !payloadNeedsPercentEncode(payload) {
		if strings.HasPrefix(metaLower, "data:text/html;") {
			return meta + ";charset=utf-8," + payload
		}
		return "data:text/html;charset=utf-8," + payload
	}

	encoded := percentEncodeDataPayload(payload)
	if hasCharset {
		return meta + "," + encoded
	}
	if strings.HasPrefix(metaLower, "data:text/html;") {
		return meta + ";charset=utf-8," + encoded
	}
	return "data:text/html;charset=utf-8," + encoded
}

func payloadNeedsPercentEncode(payload string) bool {
	for i := 0; i < len(payload); i++ {
		if payload[i] > 127 {
			return true
		}
	}
	return false
}

func percentEncodeDataPayload(payload string) string {
	var b strings.Builder
	b.Grow(len(payload) + 16)
	for i := 0; i < len(payload); i++ {
		c := payload[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') ||
			c == '-' || c == '_' || c == '.' || c == '~' {
			b.WriteByte(c)
			continue
		}
		b.WriteString(fmt.Sprintf("%%%02X", c))
	}
	return b.String()
}
