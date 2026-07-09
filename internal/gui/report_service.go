package gui

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/paths"
	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/report"
	"github.com/bafgion/scenaria-golang/internal/report/allure"
	"github.com/bafgion/scenaria-golang/internal/runstatus"
)

type projectPathProvider func() string
type runIDProvider func() string
type reportBridgeCredentials func() (string, string)

// ReportService owns report/artifact path discovery and run report writes.
type ReportService struct {
	projectPath projectPathProvider
	runID       runIDProvider
	bridgeCreds reportBridgeCredentials
}

func NewReportService(
	projectPath projectPathProvider,
	runID runIDProvider,
	bridgeCreds reportBridgeCredentials,
) *ReportService {
	return &ReportService{
		projectPath: projectPath,
		runID:       runID,
		bridgeCreds: bridgeCreds,
	}
}

func (s *ReportService) scenariaDirs() []string {
	if s == nil || s.projectPath == nil {
		return nil
	}
	root := s.projectPath()
	if root == "" {
		return nil
	}
	dirs := []string{filepath.Join(root, ".scenaria")}
	if writable, err := paths.WritableScenariaDir(root); err == nil && writable != dirs[0] {
		dirs = append([]string{writable}, dirs...)
	}
	return dirs
}

func (s *ReportService) ScenariaArtifactPath(sub string) string {
	if s == nil || s.projectPath == nil {
		return ""
	}
	root := s.projectPath()
	if root == "" {
		return ""
	}
	path, err := paths.ScenariaArtifactPath(root, sub)
	if err != nil {
		return filepath.Join(root, ".scenaria", sub)
	}
	return path
}

func (s *ReportService) ArtifactExists(path string) bool {
	if path == "" {
		return false
	}
	st, err := os.Stat(path)
	if err != nil {
		return false
	}
	if st.IsDir() {
		entries, err := os.ReadDir(path)
		return err == nil && len(entries) > 0
	}
	return true
}

func (s *ReportService) latestHTMLReport() string {
	if s == nil || s.projectPath == nil {
		return ""
	}
	root := strings.TrimSpace(s.projectPath())
	if root == "" {
		return ""
	}
	latest, err := report.ReadLatestRunPointer(root)
	if err != nil || latest == nil {
		return ""
	}
	html := strings.TrimSpace(latest.HTMLPath)
	if html != "" && s.ArtifactExists(html) {
		return html
	}
	return ""
}

func (s *ReportService) ProjectArtifacts() ProjectArtifacts {
	dirs := s.scenariaDirs()
	if len(dirs) == 0 {
		return ProjectArtifacts{}
	}
	out := ProjectArtifacts{}
	if html := s.latestHTMLReport(); html != "" {
		out.HTMLReport = html
	}
	for _, scenaria := range dirs {
		if out.AllureDir == "" && s.ArtifactExists(filepath.Join(scenaria, "allure-results")) {
			out.AllureDir = filepath.Join(scenaria, "allure-results")
		}
		if out.TracesDir == "" && s.ArtifactExists(filepath.Join(scenaria, "traces")) {
			out.TracesDir = filepath.Join(scenaria, "traces")
		}
		if out.VideosDir == "" && s.ArtifactExists(filepath.Join(scenaria, "videos")) {
			out.VideosDir = filepath.Join(scenaria, "videos")
		}
		if out.HTMLReport == "" && s.ArtifactExists(filepath.Join(scenaria, "report.html")) {
			out.HTMLReport = filepath.Join(scenaria, "report.html")
		}
		if out.JUnitReport == "" && s.ArtifactExists(filepath.Join(scenaria, "junit.xml")) {
			out.JUnitReport = filepath.Join(scenaria, "junit.xml")
		}
		if out.SummaryJSON == "" && s.ArtifactExists(filepath.Join(scenaria, "summary.json")) {
			out.SummaryJSON = filepath.Join(scenaria, "summary.json")
		}
	}
	return out
}

func remapRunArtifacts(root string, req RunRequest) RunRequest {
	req.HTMLPath = paths.RemapScenariaArtifact(root, req.HTMLPath)
	req.JUnitPath = paths.RemapScenariaArtifact(root, req.JUnitPath)
	req.SummaryJSON = paths.RemapScenariaArtifact(root, req.SummaryJSON)
	req.AllureDir = paths.RemapScenariaArtifact(root, req.AllureDir)
	req.TraceDir = paths.RemapScenariaArtifact(root, req.TraceDir)
	req.VideoDir = paths.RemapScenariaArtifact(root, req.VideoDir)
	return req
}

func (s *ReportService) LayoutRunRequestArtifacts(root string, req RunRequest) (RunRequest, report.RunArtifactLayout, error) {
	req = remapRunArtifacts(root, req)
	runID := ""
	if s != nil && s.runID != nil {
		runID = s.runID()
	}
	layout, err := report.LayoutRunArtifacts(root, runID, report.RunArtifactInputs(
		req.HTMLPath, req.JUnitPath, req.SummaryJSON, req.AllureDir, req.TraceDir, req.VideoDir,
	))
	if err != nil {
		return req, layout, err
	}
	if layout.RunID != "" {
		req.HTMLPath = layout.HTMLPath
		req.JUnitPath = layout.JUnitPath
		req.SummaryJSON = layout.SummaryJSON
		req.AllureDir = layout.AllureDir
		req.TraceDir = layout.TraceDir
		req.VideoDir = layout.VideoDir
	}
	return req, layout, nil
}

func (s *ReportService) WriteRunReports(projectRoot string, req RunRequest, plan player.ExecutionPlan, result player.ExecutionResult) (report.RunArtifactLayout, error) {
	runID := result.RunID
	if runID == "" && s != nil && s.runID != nil {
		runID = s.runID()
	}
	layout, err := report.LayoutRunArtifacts(projectRoot, runID, report.RunArtifactInputs(
		req.HTMLPath, req.JUnitPath, req.SummaryJSON, req.AllureDir, req.TraceDir, req.VideoDir,
	))
	if err != nil {
		return layout, err
	}
	if layout.RunID != "" {
		req.HTMLPath = layout.HTMLPath
		req.JUnitPath = layout.JUnitPath
		req.SummaryJSON = layout.SummaryJSON
		req.AllureDir = layout.AllureDir
		req.TraceDir = layout.TraceDir
		req.VideoDir = layout.VideoDir
	}
	var writeErr error
	var prevSummary *report.RunSummaryDetailed
	if previousSummaryPath := resolvePreviousSummaryPath(projectRoot, req.SummaryJSON, runID); previousSummaryPath != "" {
		prevSummary = report.ReadPreviousSummary(previousSummaryPath)
	}
	if req.SummaryJSON != "" {
		if err := report.WriteRunSummaryDetailed(req.SummaryJSON, report.FromExecutionResultDetailed(result)); err != nil {
			writeErr = errors.Join(writeErr, err)
		}
	}
	if req.JUnitPath != "" {
		if err := report.WriteJUnit(req.JUnitPath, result); err != nil {
			writeErr = errors.Join(writeErr, err)
		}
	}
	if req.HTMLPath != "" {
		bridgeURL, bridgeToken := "", ""
		if s != nil && s.bridgeCreds != nil {
			bridgeURL, bridgeToken = s.bridgeCreds()
		}
		htmlOpts := report.HTMLOptions{
			Plan:               plan,
			ProjectRoot:        projectRoot,
			LightMode:          req.HTMLLightMode,
			BridgeURL:          bridgeURL,
			BridgeToken:        bridgeToken,
			Locale:             req.ReportLocale,
			ReportDir:          filepath.Dir(req.HTMLPath),
			PreviousSummary:    prevSummary,
			RunID:              runID,
			ValidationBrowser:  req.Browser,
			ValidationHeadless: !req.Headed,
			ValidationBaseURL:  req.BaseURL,
		}
		if _, _, err := report.WriteHTMLModePair(req.HTMLPath, result, htmlOpts); err != nil {
			writeErr = errors.Join(writeErr, err)
		}
	}
	if req.AllureDir != "" {
		if err := allure.WriteResults(req.AllureDir, result); err != nil {
			writeErr = errors.Join(writeErr, err)
		}
	}
	if writeErr == nil && layout.RunID != "" && projectRoot != "" {
		if err := report.WriteLatestRunPointer(projectRoot, layout); err != nil {
			writeErr = errors.Join(writeErr, err)
		}
	}
	player.CleanupExecutionTempArtifacts(&result)
	return layout, writeErr
}

func (s *ReportService) FinalizeRunReports(
	root string,
	req RunRequest,
	plan player.ExecutionPlan,
	result player.ExecutionResult,
	runErr error,
	statusIncremental bool,
) (player.ExecutionResult, report.RunArtifactLayout, error) {
	if root != "" {
		var layoutErr error
		req, _, layoutErr = s.LayoutRunRequestArtifacts(root, req)
		if layoutErr != nil {
			return result, report.RunArtifactLayout{}, layoutErr
		}
	}
	layout, reportErr := s.WriteRunReports(root, req, plan, result)
	if reportErr != nil && runErr != nil {
		return result, layout, errors.Join(runErr, reportErr)
	}
	if reportErr != nil {
		return result, layout, reportErr
	}
	if root != "" && !statusIncremental {
		recordGUIStatus(root, req, result)
	}
	return result, layout, runErr
}

func resolvePreviousSummaryPath(projectRoot, currentSummaryPath, runID string) string {
	currentSummaryPath = strings.TrimSpace(currentSummaryPath)
	runID = strings.TrimSpace(runID)
	if strings.TrimSpace(projectRoot) == "" {
		return currentSummaryPath
	}
	latest, err := report.ReadLatestRunPointer(projectRoot)
	if err != nil || latest == nil {
		return currentSummaryPath
	}
	if runID != "" && latest.RunID == runID {
		return currentSummaryPath
	}
	if summary := strings.TrimSpace(latest.SummaryJSON); summary != "" {
		return summary
	}
	return currentSummaryPath
}

func recordGUIStatus(root string, req RunRequest, result player.ExecutionResult) {
	store, err := runstatus.Open(root)
	if err != nil {
		return
	}
	engine := resolveGUIEngine(req, nil)
	for _, scenarioResult := range result.ScenarioResults {
		_ = store.Record(player.RunstatusEntry(scenarioResult, engine))
	}
}
