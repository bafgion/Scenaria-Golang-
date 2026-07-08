package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/paths"
)

// RunArtifactLayout holds per-run report output paths under .scenaria/runs/<runID>/.
type RunArtifactLayout struct {
	RunID       string `json:"runId"`
	Dir         string `json:"dir"`
	HTMLPath    string `json:"htmlPath,omitempty"`
	JUnitPath   string `json:"junitPath,omitempty"`
	SummaryJSON string `json:"summaryJson,omitempty"`
	AllureDir   string `json:"allureDir,omitempty"`
	TraceDir    string `json:"traceDir,omitempty"`
	VideoDir    string `json:"videoDir,omitempty"`
}

type runArtifactInputs struct {
	HTMLPath    string
	JUnitPath   string
	SummaryJSON string
	AllureDir   string
	TraceDir    string
	VideoDir    string
}

// RunArtifactInputs builds layout inputs from configured artifact paths.
func RunArtifactInputs(htmlPath, junitPath, summaryJSON, allureDir, traceDir, videoDir string) runArtifactInputs {
	return runArtifactInputs{
		HTMLPath:    htmlPath,
		JUnitPath:   junitPath,
		SummaryJSON: summaryJSON,
		AllureDir:   allureDir,
		TraceDir:    traceDir,
		VideoDir:    videoDir,
	}
}

// LayoutRunArtifacts maps configured artifact paths into .scenaria/runs/<runID>/ when runID is set.
func LayoutRunArtifacts(projectRoot, runID string, in runArtifactInputs) (RunArtifactLayout, error) {
	layout := RunArtifactLayout{RunID: strings.TrimSpace(runID)}
	if layout.RunID == "" || strings.TrimSpace(projectRoot) == "" {
		layout.HTMLPath = in.HTMLPath
		layout.JUnitPath = in.JUnitPath
		layout.SummaryJSON = in.SummaryJSON
		layout.AllureDir = in.AllureDir
		layout.TraceDir = in.TraceDir
		layout.VideoDir = in.VideoDir
		return layout, nil
	}
	scenariaDir, err := paths.WritableScenariaDir(projectRoot)
	if err != nil {
		return layout, fmt.Errorf("resolve scenaria dir: %w", err)
	}
	layout.Dir = filepath.Join(scenariaDir, "runs", layout.RunID)
	layout.HTMLPath = runArtifactPath(layout.Dir, in.HTMLPath, "report.html")
	layout.JUnitPath = runArtifactPath(layout.Dir, in.JUnitPath, "junit.xml")
	layout.SummaryJSON = runArtifactPath(layout.Dir, in.SummaryJSON, "summary.json")
	layout.AllureDir = runArtifactPath(layout.Dir, in.AllureDir, "allure-results")
	layout.TraceDir = runArtifactPath(layout.Dir, in.TraceDir, "traces")
	layout.VideoDir = runArtifactPath(layout.Dir, in.VideoDir, "videos")
	return layout, nil
}

func runArtifactPath(runDir, configured, defaultName string) string {
	if strings.TrimSpace(configured) == "" {
		return ""
	}
	name := filepath.Base(configured)
	if name == "" || name == "." {
		name = defaultName
	}
	return filepath.Join(runDir, name)
}

// WriteLatestRunPointer records the most recent run artifact layout.
func WriteLatestRunPointer(projectRoot string, layout RunArtifactLayout) error {
	if strings.TrimSpace(layout.RunID) == "" || strings.TrimSpace(projectRoot) == "" {
		return nil
	}
	scenariaDir, err := paths.WritableScenariaDir(projectRoot)
	if err != nil {
		return err
	}
	runsDir := filepath.Join(scenariaDir, "runs")
	if err := os.MkdirAll(runsDir, 0o755); err != nil {
		return err
	}
	payload, err := json.MarshalIndent(layout, "", "  ")
	if err != nil {
		return fmt.Errorf("encode latest run pointer: %w", err)
	}
	return writeAtomic(filepath.Join(runsDir, "latest.json"), append(payload, '\n'))
}

// ReadLatestRunPointer loads the most recent run artifact layout from .scenaria/runs/latest.json.
func ReadLatestRunPointer(projectRoot string) (*RunArtifactLayout, error) {
	if strings.TrimSpace(projectRoot) == "" {
		return nil, nil
	}
	scenariaDir, err := paths.WritableScenariaDir(projectRoot)
	if err != nil {
		return nil, err
	}
	pointerPath := filepath.Join(scenariaDir, "runs", "latest.json")
	payload, err := os.ReadFile(pointerPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var layout RunArtifactLayout
	if err := json.Unmarshal(payload, &layout); err != nil {
		return nil, fmt.Errorf("decode latest run pointer: %w", err)
	}
	if strings.TrimSpace(layout.RunID) == "" {
		return nil, nil
	}
	return &layout, nil
}
