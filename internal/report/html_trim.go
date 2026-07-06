package report

import "encoding/json"

// defaultMaxEmbeddedJSONBytes limits embedded screenshot data in a single HTML file.
const defaultMaxEmbeddedJSONBytes = 4 << 20 // 4 MiB

func trimPayloadArtifacts(payload *htmlReportPayload, maxBytes int) bool {
	if payload == nil {
		return false
	}
	if maxBytes <= 0 {
		maxBytes = defaultMaxEmbeddedJSONBytes
	}
	trimmed := false
	for pass := 0; pass < 4; pass++ {
		raw, err := json.Marshal(payload)
		if err != nil || len(raw) <= maxBytes {
			break
		}
		switch pass {
		case 0:
			trimScreenshots(payload)
		case 1:
			trimSnapshots(payload)
		case 2:
			trimTraceEmbeds(payload)
		default:
			trimSnapshots(payload)
			trimTraceEmbeds(payload)
		}
		trimmed = true
	}
	if trimmed {
		payload.ArtifactsTrimmed = true
	}
	return trimmed
}

func trimScreenshots(payload *htmlReportPayload) {
	for i := range payload.Scenarios {
		sc := &payload.Scenarios[i]
		sc.Screenshot = ""
		for j := range sc.Steps {
			sc.Steps[j].Screenshot = ""
		}
	}
}

func trimSnapshots(payload *htmlReportPayload) {
	for i := range payload.Scenarios {
		sc := &payload.Scenarios[i]
		for j := range sc.Steps {
			sc.Steps[j].DOMSnapshot = ""
			sc.Steps[j].A11ySnapshot = ""
		}
	}
}

func trimTraceEmbeds(payload *htmlReportPayload) {
	for i := range payload.Scenarios {
		sc := &payload.Scenarios[i]
		sc.TraceEvents = nil
	}
}
