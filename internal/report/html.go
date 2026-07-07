package report

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/player"
)

//go:embed html_assets/viewer.html
var viewerHTMLShell string

//go:embed html_assets/viewer.css
var viewerCSS string

//go:embed html_assets/viewer.js
var viewerJS string

//go:embed html_assets/zip_trace.js
var zipTraceJS string

func viewerScript() string {
	return zipTraceJS + "\n" + viewerJS
}

func WriteHTML(path string, result player.ExecutionResult, opts HTMLOptions) error {
	if !opts.LightMode {
		if err := prepareHTMLArtifactDirs(path); err != nil {
			return err
		}
	}
	payload, err := buildHTMLPayload(result, opts, path)
	if err != nil {
		return fmt.Errorf("build html payload: %w", err)
	}
	payload.ScreenshotDedup = dedupeScreenshotURLs(&payload)
	trimPayloadArtifacts(&payload, opts.MaxJSONBytes)
	shell := strings.NewReplacer(
		"__CSS__", viewerCSS,
		"__JS__", viewerScript(),
	).Replace(viewerHTMLShell)
	shell = strings.Replace(shell, "<title>Scenaria Report</title>", "<title>"+payload.Brand+" Report</title>", 1)
	parts := strings.SplitN(shell, "__JSON__", 2)
	if len(parts) != 2 {
		return fmt.Errorf("html shell missing JSON placeholder")
	}
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create html report dir %q: %w", dir, err)
		}
	}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("write html report %q: %w", path, err)
	}
	defer file.Close()
	if _, err := file.WriteString(parts[0]); err != nil {
		return fmt.Errorf("write html report %q: %w", path, err)
	}
	enc := json.NewEncoder(file)
	if err := enc.Encode(payload); err != nil {
		return fmt.Errorf("encode html payload: %w", err)
	}
	if _, err := file.WriteString(parts[1]); err != nil {
		return fmt.Errorf("write html report %q: %w", path, err)
	}
	return nil
}
