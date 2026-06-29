package gui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/vanessa"
)

type VanessaCaseDTO struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type VanessaRunResultDTO struct {
	Output  string           `json:"output"`
	Error   string           `json:"error"`
	Success bool             `json:"success"`
	RunDir  string           `json:"runDir"`
	Cases   []VanessaCaseDTO `json:"cases"`
}

type VanessaRunSnapshotDTO struct {
	RunDir          string           `json:"runDir"`
	CurrentScenario string           `json:"currentScenario"`
	CompletedCases  int              `json:"completedCases"`
	TotalPlanned    int              `json:"totalPlanned"`
	Cases           []VanessaCaseDTO `json:"cases"`
}

func pluginToVanessaRun(projectRoot string, req PluginRunRequest) vanessa.RunRequest {
	vReq := vanessa.RunRequest{
		ProjectRoot:        projectRoot,
		Tag:                strings.TrimSpace(req.Tag),
		ExcludeTags:        append([]string(nil), req.ExcludeTags...),
		DryRun:             req.DryRun,
		PlatformExecutable: strings.TrimSpace(req.PlatformExe),
		EPFPath:            strings.TrimSpace(req.EPFPath),
		IBConnection:       strings.TrimSpace(req.IBConnection),
		ReportAllure:       req.ReportAllure,
		RerunFailedRunDir:  strings.TrimSpace(req.RerunFailedRunDir),
		InstallEPF:         req.InstallEPF,
		EPFDownloadURL:     strings.TrimSpace(req.EPFURL),
		EPFDestination:     strings.TrimSpace(req.EPFDest),
	}
	if dir := strings.TrimSpace(req.VaDir); dir != "" {
		vReq.Paths = []string{dir}
	}
	if files := strings.TrimSpace(req.VaFiles); files != "" {
		parts := strings.FieldsFunc(files, func(r rune) bool {
			return r == ',' || r == ';'
		})
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				vReq.Paths = append(vReq.Paths, part)
			}
		}
	}
	if scenario := strings.TrimSpace(req.Scenario); scenario != "" {
		vReq.ScenarioNames = []string{scenario}
	}
	if len(vReq.Paths) == 0 {
		vReq.Paths = []string{projectRoot}
	}
	return vReq
}

func (s *Service) RunVanessaPlugin(req PluginRunRequest) VanessaRunResultDTO {
	path := s.ProjectPath()
	if path == "" {
		return VanessaRunResultDTO{Error: "open a project folder first"}
	}
	vReq := pluginToVanessaRun(path, req)
	result, err := vanessa.Run(vReq)
	dto := vanessaBatchToDTO(result)
	if err != nil && dto.Error == "" {
		dto.Error = err.Error()
	}
	return dto
}

func (s *Service) PollVanessaRun(runDir string, totalPlanned int) VanessaRunSnapshotDTO {
	runDir = strings.TrimSpace(runDir)
	if runDir == "" {
		return VanessaRunSnapshotDTO{}
	}
	if totalPlanned < 1 {
		totalPlanned = 1
	}
	monitor := vanessa.NewRunMonitor(runDir, totalPlanned)
	snap := monitor.Poll()
	all := vanessa.ParseJUnitDir(filepath.Join(runDir, "junit"))
	cases := make([]VanessaCaseDTO, 0, len(all))
	for _, item := range all {
		cases = append(cases, caseToDTO(item))
	}
	completed := len(cases)
	total := totalPlanned
	if completed > total {
		total = completed
	}
	current := snap.CurrentScenario
	if current == "" && len(cases) > 0 {
		current = cases[len(cases)-1].Name
	}
	return VanessaRunSnapshotDTO{
		RunDir:          runDir,
		CurrentScenario: current,
		CompletedCases:  completed,
		TotalPlanned:    total,
		Cases:           cases,
	}
}

func vanessaBatchToDTO(result vanessa.BatchResult) VanessaRunResultDTO {
	cases := make([]VanessaCaseDTO, 0, len(result.Cases))
	var b strings.Builder
	for _, item := range result.Cases {
		cases = append(cases, caseToDTO(item))
		mark := "✓"
		if !item.Success {
			mark = "✗"
		}
		fmt.Fprintf(&b, "%s %s: %s\n", mark, item.Name, item.Message)
	}
	if result.RunDir != "" {
		fmt.Fprintf(&b, "Run directory: %s\n", result.RunDir)
	}
	out := strings.TrimRight(b.String(), "\n")
	errText := result.Error
	if !result.Success && errText == "" && result.ExitCode != 0 {
		errText = fmt.Sprintf("vanessa run failed with exit code %d", result.ExitCode)
	}
	return VanessaRunResultDTO{
		Output:  out,
		Error:   errText,
		Success: result.Success,
		RunDir:  result.RunDir,
		Cases:   cases,
	}
}

func caseToDTO(item vanessa.CaseResult) VanessaCaseDTO {
	return VanessaCaseDTO{
		Path:    item.Path,
		Name:    item.Name,
		Success: item.Success,
		Message: item.Message,
	}
}
