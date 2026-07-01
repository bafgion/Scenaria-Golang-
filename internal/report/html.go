package report

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/bafgion/scenaria-golang/internal/brand"
	"github.com/bafgion/scenaria-golang/internal/player"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="utf-8">
  <title>{{.BrandName}} Report</title>
  <style>
    body { font-family: Segoe UI, sans-serif; margin: 24px; color: #111; }
    h1 { margin-bottom: 8px; }
    .meta { color: #555; margin-bottom: 24px; }
    table { border-collapse: collapse; width: 100%; }
    th, td { border: 1px solid #ddd; padding: 8px 12px; text-align: left; }
    th { background: #f5f5f5; }
    .passed { color: #0a7a2f; }
    .failed { color: #b00020; }
    .skipped { color: #666; }
  </style>
</head>
<body>
  <h1>{{.BrandName}} Run Report</h1>
  <div class="meta">Generated at {{.GeneratedAt}} · Mode: {{.Mode}} · Files: {{.Files}} · Scenarios: {{.Scenarios}} · Steps: {{.Steps}}</div>
  <table>
    <thead><tr><th>Feature</th><th>Scenario</th><th>Status</th><th>Message</th></tr></thead>
    <tbody>
    {{range .Cases}}
      <tr>
        <td>{{.FeaturePath}}</td>
        <td>{{.Scenario}}</td>
        <td class="{{.Status}}">{{.Status}}</td>
        <td>{{.Message}}</td>
      </tr>
    {{end}}
    </tbody>
  </table>
</body>
</html>`

type htmlReportData struct {
	BrandName   string
	GeneratedAt string
	Mode        string
	Files       int
	Scenarios   int
	Steps       int
	Cases       []player.ScenarioResult
}

func WriteHTML(path string, result player.ExecutionResult) error {
	data := htmlReportData{
		BrandName:   brand.Name,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Mode:        result.Mode,
		Files:       result.Files,
		Scenarios:   result.Scenarios,
		Steps:       result.Steps,
		Cases:       result.ScenarioResults,
	}
	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("parse html template: %w", err)
	}
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create html report dir %q: %w", dir, err)
		}
	}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create html report %q: %w", path, err)
	}
	defer file.Close()
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("render html report: %w", err)
	}
	return nil
}
