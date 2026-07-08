package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/player"
)

func writeHTMLPayloadFile(path string, payload htmlReportPayload) error {
	payload.ScreenshotDedup = dedupeScreenshotURLs(&payload)
	shell := strings.NewReplacer(
		"__CSS__", viewerCSS,
		"__JS__", viewerScript(),
	).Replace(viewerHTMLShell)
	shell = strings.Replace(shell, "<title>Scenaria Report</title>", "<title>"+payload.Brand+" Report</title>", 1)
	parts := strings.SplitN(shell, "__JSON__", 2)
	if len(parts) != 2 {
		return fmt.Errorf("html shell missing JSON placeholder")
	}
	var buf bytes.Buffer
	if _, err := buf.WriteString(parts[0]); err != nil {
		return fmt.Errorf("build html report %q: %w", path, err)
	}
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(payload); err != nil {
		return fmt.Errorf("encode html payload: %w", err)
	}
	if _, err := buf.WriteString(parts[1]); err != nil {
		return fmt.Errorf("build html report %q: %w", path, err)
	}
	if err := writeAtomic(path, buf.Bytes()); err != nil {
		return fmt.Errorf("write html report %q: %w", path, err)
	}
	return nil
}

func cloneLightHTMLPayload(full htmlReportPayload, lightPath string, opts HTMLOptions) htmlReportPayload {
	light := full
	light.LightMode = true
	light.ModeLinks = buildHTMLModeLinks(lightPath, HTMLOptions{
		LightMode:         true,
		ModePairFullPath:  opts.ModePairFullPath,
		ModePairLightPath: opts.ModePairLightPath,
	})
	light.Scenarios = make([]htmlScenario, len(full.Scenarios))
	for i, sc := range full.Scenarios {
		lightSc := sc
		lightSc.TracePath = ""
		lightSc.TraceCommand = ""
		lightSc.TraceEvents = nil
		if lightSc.Status != "failed" && lightSc.Status != "broken" {
			lightSc.Screenshot = ""
		}
		lightSc.Steps = make([]htmlStep, len(sc.Steps))
		for j, step := range sc.Steps {
			lightStep := step
			lightStep.TraceOffsetMS = 0
			lightStep.DOMSnapshot = ""
			lightStep.A11ySnapshot = ""
			if lightStep.Status != "failed" && lightStep.Status != "broken" {
				lightStep.Screenshot = ""
			}
			lightSc.Steps[j] = lightStep
		}
		light.Scenarios[i] = lightSc
	}
	return light
}

// WriteHTMLModePair writes full and light HTML reports with cross-links between them.
// Artifact processing runs once for the full report; the light report reuses that payload.
func WriteHTMLModePair(path string, result player.ExecutionResult, opts HTMLOptions) (fullPath, lightPath string, err error) {
	fullPath, lightPath = htmlModePairPaths(path, opts.LightMode)
	fullOpts := opts
	fullOpts.LightMode = false
	fullOpts.ModePairFullPath = fullPath
	fullOpts.ModePairLightPath = lightPath
	if err := prepareHTMLArtifactDirs(fullPath); err != nil {
		return "", "", err
	}
	fullPayload, err := buildHTMLPayload(result, fullOpts, fullPath)
	if err != nil {
		return "", "", err
	}
	lightPayload := cloneLightHTMLPayload(fullPayload, lightPath, fullOpts)
	if opts.MaxJSONBytes > 0 {
		trimPayloadArtifacts(&lightPayload, opts.MaxJSONBytes)
		trimPayloadArtifacts(&fullPayload, opts.MaxJSONBytes)
	}
	if err := writeHTMLPayloadFile(lightPath, lightPayload); err != nil {
		return "", "", err
	}
	if err := writeHTMLPayloadFile(fullPath, fullPayload); err != nil {
		return "", "", err
	}
	return fullPath, lightPath, nil
}
