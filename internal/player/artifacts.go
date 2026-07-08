package player

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	playwright "github.com/mxschmitt/playwright-go"
)

var unsafeNameRE = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func captureFailureArtifacts(session *browserSession, input ScenarioInput, traceDir, videoDir string) (screenshotPath, tracePath, videoPath string) {
	// Trace/video/screenshot artifacts are collected only when a scenario fails.
	if session == nil {
		return "", "", ""
	}
	if closed := session.isClosed(); closed {
		return "", "", ""
	}

	baseDir := strings.TrimSpace(traceDir)
	if baseDir == "" {
		baseDir = strings.TrimSpace(videoDir)
	}
	if baseDir == "" {
		baseDir = filepath.Join(os.TempDir(), "scenaria-artifacts")
	}
	_ = os.MkdirAll(baseDir, 0o755)

	screenshotPath = captureFailureScreenshotPath(session, filepath.Join(baseDir, "screenshots"), input)
	tracePath = captureTraceZIPPath(session, traceDir, input)
	if session.videoEnabled {
		videoPath = session.finalizeVideoRecordingPath(videoDir)
	}
	return screenshotPath, tracePath, videoPath
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
	path := captureTraceZIPPath(session, dir, input)
	if strings.TrimSpace(path) == "" {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
}

func captureTraceZIPPath(session *browserSession, dir string, input ScenarioInput) string {
	if session == nil || strings.TrimSpace(dir) == "" {
		return ""
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return ""
	}
	path := filepath.Join(dir, artifactBaseName(input)+".zip")
	session.mu.Lock()
	if !session.traceEnabled || session.traceStopped || session.context == nil {
		session.mu.Unlock()
		return ""
	}
	if err := session.context.Tracing().Stop(path); err != nil {
		session.mu.Unlock()
		return ""
	}
	session.traceStopped = true
	session.mu.Unlock()
	return path
}

func captureFailureScreenshotPath(session *browserSession, dir string, input ScenarioInput) string {
	if session == nil || strings.TrimSpace(dir) == "" {
		return ""
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return ""
	}
	payload := captureFailureScreenshot(session)
	if len(payload) == 0 {
		return ""
	}
	path := filepath.Join(dir, artifactBaseName(input)+".png")
	tmp, err := os.CreateTemp(filepath.Dir(path), ".screenshot-*")
	if err != nil {
		return ""
	}
	tmpPath := tmp.Name()
	if _, err := tmp.Write(payload); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
		return ""
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return ""
	}
	if err := os.Rename(tmpPath, path); err != nil {
		if removeErr := os.Remove(path); removeErr == nil || os.IsNotExist(removeErr) {
			if retryErr := os.Rename(tmpPath, path); retryErr == nil {
				return path
			}
		}
		_ = os.Remove(tmpPath)
		return ""
	}
	return path
}

func copyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

func artifactBaseName(input ScenarioInput) string {
	if id := strings.TrimSpace(input.CaseID); id != "" {
		base := unsafeNameRE.ReplaceAllString(strings.ReplaceAll(id, "::", "__"), "_")
		if input.ExampleIndex > 0 {
			base = fmt.Sprintf("%s_e%d", base, input.ExampleIndex)
		}
		if base != "" {
			return base
		}
	}
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
