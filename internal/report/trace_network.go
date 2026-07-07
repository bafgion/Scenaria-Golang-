package report

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type htmlNetworkEvent struct {
	OffsetMS int64  `json:"offset_ms"`
	Snippet  string `json:"snippet"`
}

func parseTraceNetworkFailures(zipData []byte, limit int) []htmlNetworkEvent {
	if len(zipData) == 0 {
		return nil
	}
	if limit <= 0 {
		limit = 32
	}
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil
	}
	for _, f := range reader.File {
		if !strings.HasSuffix(f.Name, ".network") {
			continue
		}
		rc, err := reader.Open(f.Name)
		if err != nil {
			continue
		}
		events := parseTraceNetworkStream(rc, limit)
		rc.Close()
		if len(events) > 0 {
			return events
		}
	}
	return nil
}

func parseTraceNetworkStream(r io.Reader, limit int) []htmlNetworkEvent {
	out := make([]htmlNetworkEvent, 0, limit)
	var baseTime float64
	baseSet := false
	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 4096)
	for {
		n, err := r.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[:n]...)
			for {
				idx := bytes.IndexByte(buf, '\n')
				if idx < 0 {
					break
				}
				line := bytes.TrimSpace(buf[:idx])
				buf = buf[idx+1:]
				if len(line) == 0 {
					continue
				}
				at, snippet, ok := parseTraceNetworkLine(line)
				if !ok {
					continue
				}
				if !baseSet {
					baseTime = at
					baseSet = true
				}
				out = append(out, htmlNetworkEvent{
					OffsetMS: int64(at - baseTime),
					Snippet:  snippet,
				})
				if len(out) >= limit {
					return out
				}
			}
		}
		if err != nil {
			break
		}
	}
	return out
}

func parseTraceNetworkLine(line []byte) (monotonicTime float64, snippet string, ok bool) {
	var raw map[string]any
	if json.Unmarshal(line, &raw) != nil {
		return 0, "", false
	}
	snap, _ := raw["snapshot"].(map[string]any)
	if snap == nil {
		return 0, "", false
	}
	if t, ok := numField(snap, "_monotonicTime"); ok {
		monotonicTime = t
	}
	req, _ := snap["request"].(map[string]any)
	resp, _ := snap["response"].(map[string]any)
	method := strField(req, "method")
	url := strField(req, "url")
	if url == "" {
		return 0, "", false
	}
	if method == "" {
		method = "GET"
	}
	status := int(numFieldDefault(resp, "status", 0))
	failText := strField(snap, "_failureText")
	if failText == "" {
		failText = strField(snap, "failureText")
	}
	if status > 0 && status < 400 && failText == "" {
		return 0, "", false
	}
	switch {
	case failText != "":
		snippet = fmt.Sprintf("%s %s — %s", method, url, failText)
	case status >= 400:
		snippet = fmt.Sprintf("%s %s — HTTP %d", method, url, status)
	default:
		snippet = fmt.Sprintf("%s %s", method, url)
	}
	return monotonicTime, snippet, true
}

func strField(m map[string]any, key string) string {
	if m == nil {
		return ""
	}
	v, _ := m[key].(string)
	return strings.TrimSpace(v)
}

func numField(m map[string]any, key string) (float64, bool) {
	if m == nil {
		return 0, false
	}
	switch v := m[key].(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case json.Number:
		f, err := v.Float64()
		return f, err == nil
	default:
		return 0, false
	}
}

func numFieldDefault(m map[string]any, key string, def float64) float64 {
	if v, ok := numField(m, key); ok {
		return v
	}
	return def
}

func enrichStepsNetworkFromTrace(steps []htmlStep, failures []htmlNetworkEvent) {
	if len(steps) == 0 || len(failures) == 0 {
		return
	}
	const maxDeltaMS = int64(3000)
	for i := range steps {
		if steps[i].Network != "" {
			continue
		}
		off := steps[i].TraceOffsetMS
		best := ""
		bestDelta := int64(1 << 62)
		for _, f := range failures {
			d := off - f.OffsetMS
			if d < 0 {
				d = -d
			}
			if d < bestDelta {
				bestDelta = d
				best = f.Snippet
			}
		}
		if best != "" && bestDelta <= maxDeltaMS {
			steps[i].Network = best
		}
	}
}
