package player

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	playwright "github.com/mxschmitt/playwright-go"
)

var unsafeNameRE = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func captureFailureArtifacts(session *browserSession, input ScenarioInput, traceDir, videoDir string) (screenshot, trace, video []byte) {
	if session == nil {
		return nil, nil, nil
	}
	session.mu.Lock()
	closed := session.closed
	session.mu.Unlock()
	if closed {
		return nil, nil, nil
	}
	screenshot = captureFailureScreenshot(session)
	trace = captureTraceZIP(session, traceDir, input)
	if session.videoEnabled {
		video = session.finalizeVideoRecording(videoDir)
	}
	return screenshot, trace, video
}

func readVideoRecording(recorder playwright.Video, dir string) []byte {
	path, err := recorder.Path()
	if err != nil || strings.TrimSpace(path) == "" {
		return nil
	}
	if !filepath.IsAbs(path) && dir != "" {
		path = filepath.Join(dir, path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
}

func captureTraceZIP(session *browserSession, dir string, input ScenarioInput) []byte {
	if session == nil || !session.traceEnabled || session.traceStopped || session.context == nil || strings.TrimSpace(dir) == "" {
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil
	}
	path := filepath.Join(dir, artifactBaseName(input)+".zip")
	if err := session.context.Tracing().Stop(path); err != nil {
		return nil
	}
	session.traceStopped = true
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
}

func artifactBaseName(input ScenarioInput) string {
	base := strings.TrimSpace(input.ScenarioName)
	if base == "" {
		base = "scenario"
	}
	base = unsafeNameRE.ReplaceAllString(base, "_")
	if base == "" {
		base = "scenario"
	}
	if input.FeaturePath != "" {
		feature := unsafeNameRE.ReplaceAllString(filepath.Base(input.FeaturePath), "_")
		if feature != "" {
			return fmt.Sprintf("%s__%s", strings.TrimSuffix(feature, filepath.Ext(feature)), base)
		}
	}
	return base
}
