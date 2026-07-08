package report

import (
	_ "embed"
	"fmt"

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
	if opts.MaxJSONBytes > 0 {
		trimPayloadArtifacts(&payload, opts.MaxJSONBytes)
	}
	return writeHTMLPayloadFile(path, payload)
}
