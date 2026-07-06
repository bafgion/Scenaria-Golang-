package report

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type htmlTraceEvent struct {
	OffsetMS int64  `json:"offset_ms"`
	Title    string `json:"title"`
	Kind     string `json:"kind,omitempty"`
}

type traceActionLine struct {
	Type      string          `json:"type"`
	StartTime float64         `json:"startTime"`
	Class     string          `json:"class"`
	Method    string          `json:"method"`
	Params    json.RawMessage `json:"params"`
}

func parseTraceActions(zipData []byte, limit int) []htmlTraceEvent {
	if len(zipData) == 0 {
		return nil
	}
	if limit <= 0 {
		limit = 120
	}
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil
	}
	var traceName string
	for _, f := range reader.File {
		if strings.HasSuffix(f.Name, ".trace") {
			traceName = f.Name
			break
		}
	}
	if traceName == "" {
		return nil
	}
	f, err := reader.Open(traceName)
	if err != nil {
		return nil
	}
	defer f.Close()
	return parseTraceActionStream(f, limit)
}

func parseTraceActionStream(r io.Reader, limit int) []htmlTraceEvent {
	events := make([]htmlTraceEvent, 0, limit)
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
				var row traceActionLine
				if json.Unmarshal(line, &row) != nil || row.Type != "action" {
					continue
				}
				if !baseSet {
					baseTime = row.StartTime
					baseSet = true
				}
				title := traceActionTitle(row)
				if title == "" {
					continue
				}
				events = append(events, htmlTraceEvent{
					OffsetMS: int64(row.StartTime - baseTime),
					Title:    title,
					Kind:     row.Method,
				})
				if len(events) >= limit {
					return events
				}
			}
		}
		if err != nil {
			break
		}
	}
	return events
}

func traceActionTitle(row traceActionLine) string {
	method := strings.TrimSpace(row.Method)
	if method == "" {
		return ""
	}
	class := strings.TrimSpace(row.Class)
	title := method
	if class != "" {
		title = class + "." + method
	}
	if len(row.Params) > 0 {
		var params map[string]any
		if json.Unmarshal(row.Params, &params) == nil {
			for _, key := range []string{"url", "selector", "text", "name"} {
				if v, ok := params[key]; ok {
					s := strings.TrimSpace(fmt.Sprint(v))
					if s != "" {
						if len(s) > 60 {
							s = s[:60] + "…"
						}
						title += " " + s
						break
					}
				}
			}
		}
	}
	return title
}
