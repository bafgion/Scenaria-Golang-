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
	payload, err := buildHTMLPayload(result, opts, path)
	if err != nil {
		return fmt.Errorf("build html payload: %w", err)
	}
	payload.ScreenshotDedup = dedupeScreenshotURLs(&payload)
	trimPayloadArtifacts(&payload, opts.MaxJSONBytes)
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode html payload: %w", err)
	}
	html := strings.NewReplacer(
		"__CSS__", viewerCSS,
		"__JS__", viewerScript(),
		"__JSON__", string(jsonBytes),
	).Replace(viewerHTMLShell)
	html = strings.Replace(html, "<title>Scenaria Report</title>", "<title>"+payload.Brand+" Report</title>", 1)

	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create html report dir %q: %w", dir, err)
		}
	}
	if err := os.WriteFile(path, []byte(html), 0o644); err != nil {
		return fmt.Errorf("write html report %q: %w", path, err)
	}
	return nil
}
