package recorder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/paths"
	"github.com/bafgion/scenaria-golang/internal/settings"
	playwright "github.com/mxschmitt/playwright-go"
)

type LiveOptions struct {
	StartURL         string
	FeatureName      string
	ScenarioName     string
	OutputPath       string
	Headless         bool
	IdleTimeout      time.Duration
	Session          *LiveSession
	AppendTo         string
	FilterImportant  bool
	NavOnly          bool
	HoverRecord      bool
	TestClient       *settings.TestClient
	HTTPCredentials  *playwright.HttpCredentials
}

type recorderEvent struct {
	Type   string            `json:"type"`
	Detail map[string]string `json:"detail"`
	TS     int64             `json:"ts"`
}

// RecordLive opens a browser, injects RecorderScript and writes captured steps to a feature file.
func RecordLive(ctx context.Context, opts LiveOptions) error {
	if strings.TrimSpace(opts.StartURL) == "" {
		return fmt.Errorf("start URL is required")
	}
	if strings.TrimSpace(opts.OutputPath) == "" {
		return fmt.Errorf("output path is required")
	}
	if opts.IdleTimeout <= 0 {
		opts.IdleTimeout = 30 * time.Second
	}
	if strings.TrimSpace(opts.FeatureName) == "" {
		opts.FeatureName = "Записанный сценарий"
	}
	if strings.TrimSpace(opts.ScenarioName) == "" {
		opts.ScenarioName = "Запись"
	}

	if err := playwright.Install(); err != nil {
		return fmt.Errorf("install playwright: %w", err)
	}
	paths.ConfigurePlaywrightBrowsers()
	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("start playwright: %w", err)
	}
	defer pw.Stop()

	recorded := []RecordedStep{{Action: "goto", Value: opts.StartURL}}
	session := opts.Session
	if session == nil {
		session = NewLiveSession()
	}
	session.InitHeadless(opts.Headless)
	session.SetRecorderFlags(opts.FilterImportant, opts.NavOnly, opts.HoverRecord)
	defer session.Clear()

	for {
		err := runLiveBrowserSession(ctx, pw, opts, session, &recorded, session.Headless())
		if err == nil {
			break
		}
		if errors.Is(err, ErrRelaunchHeadless) {
			session.ClearRelaunch()
			continue
		}
		return err
	}

	normalized := NormalizeSteps(recorded)
	steps := RecordedStepsToLines(normalized)
	if strings.TrimSpace(opts.AppendTo) != "" {
		return AppendStepsToFeature(opts.AppendTo, steps)
	}
	return WriteFeature(opts.OutputPath, Options{
		FeatureName:  opts.FeatureName,
		ScenarioName: opts.ScenarioName,
		Steps:        steps,
	})
}

func decodeEvents(raw any) ([]recorderEvent, error) {
	payload, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("encode recorder events: %w", err)
	}
	var events []recorderEvent
	if err := json.Unmarshal(payload, &events); err != nil {
		return nil, fmt.Errorf("decode recorder events: %w", err)
	}
	return events, nil
}

func normalizeDetail(detail map[string]string) map[string]string {
	if detail == nil {
		return map[string]string{}
	}
	out := make(map[string]string, len(detail))
	for key, value := range detail {
		out[strings.ToLower(key)] = value
	}
	if tag, ok := out["tag"]; ok {
		out["tag"] = strings.ToUpper(tag)
	}
	return out
}

func applyRecorderConfig(page playwright.Page, opts LiveOptions, session *LiveSession) error {
	filter, nav, hover := opts.FilterImportant, opts.NavOnly, opts.HoverRecord
	if session != nil {
		filter, nav, hover = session.RecorderFlags()
	}
	script := fmt.Sprintf(`() => {
		if (!window.__scenariaRecorder) return;
		window.__scenariaRecorder.filterImportant = %v;
		window.__scenariaRecorder.navOnly = %v;
		window.__scenariaRecorder.hoverRecord = %v;
	}`, filter, nav, hover)
	_, err := page.Evaluate(script)
	return err
}

// StepsFromFeature extracts plain step texts from a feature (for validation helpers).
func StepsFromFeature(feature *gherkin.Feature) []gherkin.Step {
	out := make([]gherkin.Step, 0)
	for _, scenario := range feature.Scenarios {
		out = append(out, scenario.Steps...)
	}
	return out
}
