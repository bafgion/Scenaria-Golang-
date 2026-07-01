package allure

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/google/uuid"
)

type resultFile struct {
	UUID          string         `json:"uuid"`
	HistoryID     string         `json:"historyId,omitempty"`
	Name          string         `json:"name"`
	FullName      string         `json:"fullName,omitempty"`
	Status        string         `json:"status"`
	StatusDetails *statusDetails `json:"statusDetails,omitempty"`
	Stage         string         `json:"stage"`
	Start         int64          `json:"start"`
	Stop          int64          `json:"stop"`
	Labels        []label        `json:"labels,omitempty"`
	Attachments   []attachment   `json:"attachments,omitempty"`
}

type attachment struct {
	Name   string `json:"name"`
	Source string `json:"source"`
	Type   string `json:"type"`
}

type statusDetails struct {
	Message string `json:"message,omitempty"`
	Trace   string `json:"trace,omitempty"`
}

type label struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// WriteResults writes Allure 2 *-result.json files for each scenario.
func WriteResults(dir string, result player.ExecutionResult) error {
	if strings.TrimSpace(dir) == "" {
		return fmt.Errorf("allure output directory is required")
	}
	dir = filepath.Clean(dir)
	if err := cleanAllureDir(dir); err != nil {
		return fmt.Errorf("clean allure dir: %w", err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create allure dir: %w", err)
	}

	base := time.Now().UnixMilli()
	for i, scenario := range result.ScenarioResults {
		start := base + int64(i*100)
		stop := start + 50
		if stop <= start {
			stop = start + 1
		}
		if err := writeScenarioResult(dir, scenario, start, stop); err != nil {
			return err
		}
	}
	return nil
}

func cleanAllureDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasSuffix(name, "-result.json") || strings.Contains(name, "-attachment") {
			if err := os.Remove(filepath.Join(dir, name)); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}
	return nil
}

func writeScenarioResult(dir string, scenario player.ScenarioResult, start, stop int64) error {
	status := mapStatus(scenario.Status)
	id := uuid.New().String()
	entry := resultFile{
		UUID:      id,
		HistoryID: historyID(scenario),
		Name:      scenario.Scenario,
		FullName:  scenario.FeaturePath + "::" + scenario.Scenario,
		Status:    status,
		Stage:     "finished",
		Start:     start,
		Stop:      stop,
		Labels: []label{
			{Name: "suite", Value: filepath.Base(scenario.FeaturePath)},
			{Name: "parentSuite", Value: scenario.FeaturePath},
			{Name: "framework", Value: "scenaria"},
			{Name: "language", Value: "go"},
		},
	}
	if scenario.Message != "" && status != "passed" {
		entry.StatusDetails = &statusDetails{Message: scenario.Message}
	}
	if status != "passed" {
		var err error
		entry.Attachments, err = writeFailureAttachments(dir, entry.Attachments, scenario)
		if err != nil {
			return err
		}
	}
	payload, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("encode allure result: %w", err)
	}
	path := filepath.Join(dir, id+"-result.json")
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		return fmt.Errorf("write allure result %q: %w", path, err)
	}
	return nil
}

func mapStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "passed", "pass", "ok":
		return "passed"
	case "skipped", "skip":
		return "skipped"
	case "broken":
		return "broken"
	default:
		return "failed"
	}
}

func historyID(scenario player.ScenarioResult) string {
	key := scenario.FeaturePath + "::" + scenario.Scenario
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(key)).String()
}

func writeFailureAttachments(dir string, existing []attachment, scenario player.ScenarioResult) ([]attachment, error) {
	out := existing
	if len(scenario.ScreenshotPNG) > 0 {
		item, err := saveAttachment(dir, "screenshot", ".png", "image/png", scenario.ScreenshotPNG)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	if len(scenario.TraceZIP) > 0 {
		item, err := saveAttachment(dir, "trace", ".zip", "application/zip", scenario.TraceZIP)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	if len(scenario.VideoWebM) > 0 {
		item, err := saveAttachment(dir, "video", ".webm", "video/webm", scenario.VideoWebM)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

func saveAttachment(dir, name, ext, mime string, payload []byte) (attachment, error) {
	attachID := uuid.New().String()
	source := attachID + "-attachment" + ext
	if err := os.WriteFile(filepath.Join(dir, source), payload, 0o644); err != nil {
		return attachment{}, fmt.Errorf("write allure attachment: %w", err)
	}
	return attachment{Name: name, Source: source, Type: mime}, nil
}
