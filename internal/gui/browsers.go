package gui

import "github.com/bafgion/scenaria-golang/internal/paths"

type BrowserInstallStatusDTO struct {
	Engine    string `json:"engine"`
	Label     string `json:"label"`
	Installed bool   `json:"installed"`
	Detail    string `json:"detail"`
}

func (s *Service) BrowserInstallStatus(engine string) BrowserInstallStatusDTO {
	engine = paths.NormalizeBrowserEngine(engine)
	installed, detail := paths.BrowserInstallStatus(engine)
	return BrowserInstallStatusDTO{
		Engine:    engine,
		Label:     paths.BrowserEngineLabels[engine],
		Installed: installed,
		Detail:    detail,
	}
}

func (s *Service) InstallBrowserEngine(engine string) RunResult {
	engine = paths.NormalizeBrowserEngine(engine)
	var lines []string
	path, err := paths.InstallBrowserEngine(engine, func(line string) {
		lines = append(lines, line)
	})
	if err != nil {
		out := stringsJoinLines(lines)
		if out != "" {
			out += "\n"
		}
		return RunResult{Output: out, Error: err.Error()}
	}
	out := stringsJoinLines(lines)
	if out != "" {
		out += "\n"
	}
	out += "Готово: " + path
	return RunResult{Output: out}
}

func stringsJoinLines(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	out := lines[0]
	for i := 1; i < len(lines); i++ {
		out += "\n" + lines[i]
	}
	return out
}
