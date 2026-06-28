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
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create allure dir: %w", err)
	}

	now := time.Now().UnixMilli()
	for _, scenario := range result.ScenarioResults {
		if err := writeScenarioResult(dir, scenario, now); err != nil {
			return err
		}
	}
	return nil
}

func writeScenarioResult(dir string, scenario player.ScenarioResult, now int64) error {
	status := mapStatus(scenario.Status)
	id := uuid.New().String()
	entry := resultFile{
		UUID:      id,
		HistoryID: historyID(scenario),
		Name:      scenario.Scenario,
		FullName:  scenario.FeaturePath + "::" + scenario.Scenario,
		Status:    status,
		Stage:     "finished",
		Start:     now,
		Stop:      now,
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
	if len(scenario.ScreenshotPNG) > 0 && status != "passed" {
		attachID := uuid.New().String()
		source := attachID + "-attachment.png"
		if err := os.WriteFile(filepath.Join(dir, source), scenario.ScreenshotPNG, 0o644); err != nil {
			return fmt.Errorf("write allure attachment: %w", err)
		}
		entry.Attachments = []attachment{{
			Name:   "screenshot",
			Source: source,
			Type:   "image/png",
		}}
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
