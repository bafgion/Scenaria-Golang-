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
	StartURL        string
	FeatureName     string
	ScenarioName    string
	OutputPath      string
	Headless        bool
	IdleTimeout     time.Duration
	Session         *LiveSession
	AppendTo        string
	FilterImportant   bool
	NavOnly           bool
	HoverRecord       bool
	ScrollBeforeClick bool
	HoverRecordMinMs  int
	TestClient      *settings.TestClient
	HTTPCredentials *playwright.HttpCredentials
	BrowseOnly      bool
	Callbacks       LiveCallbacks
}

type LiveCallbacks struct {
	OnCaptureStart  func()
	OnPickerRequest func()
	OnBrowserLost   func()
}

type recorderEvent struct {
	Type   string            `json:"type"`
	Detail map[string]string `json:"detail"`
	TS     int64             `json:"ts"`
}

// RecordLive opens a browser, injects RecorderScript and writes captured steps to a feature file.
func RecordLive(ctx context.Context, opts LiveOptions) error {
	if strings.TrimSpace(opts.OutputPath) == "" {
		return fmt.Errorf("output path is required")
	}
	if !opts.BrowseOnly && opts.IdleTimeout <= 0 {
		opts.IdleTimeout = 30 * time.Second
	}
	if strings.TrimSpace(opts.FeatureName) == "" {
		opts.FeatureName = "Записанный сценарий"
	}
	if strings.TrimSpace(opts.ScenarioName) == "" {
		opts.ScenarioName = "Запись"
	}

	if err := paths.EnsurePlaywrightEngine("chromium"); err != nil {
		return fmt.Errorf("playwright browsers: %w", err)
	}
	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("start playwright: %w", err)
	}
	defer pw.Stop()

	recorded := []RecordedStep{}
	if !opts.BrowseOnly {
		if start := strings.TrimSpace(opts.StartURL); start != "" {
			recorded = append(recorded, RecordedStep{Action: "goto", Value: start})
		}
	}
	session := opts.Session
	if session == nil {
		session = NewLiveSession()
	}
	session.InitHeadless(opts.Headless)
	session.SetRecorderOptions(opts.FilterImportant, opts.NavOnly, opts.HoverRecord, opts.ScrollBeforeClick, opts.HoverRecordMinMs)
	if opts.BrowseOnly {
		session.InitBrowseMode()
	} else {
		session.InitRecordMode()
	}
	defer session.Clear()

	var runErr error
	for {
		runErr = runLiveBrowserSession(ctx, pw, opts, session, &recorded, session.Headless())
		if runErr == nil {
			break
		}
		if errors.Is(runErr, ErrRelaunchHeadless) {
			session.ClearRelaunch()
			continue
		}
		if errors.Is(runErr, context.Canceled) {
			break
		}
		return runErr
	}

	if opts.BrowseOnly && !session.CaptureEverEnabled() {
		return context.Canceled
	}
	return persistLiveRecording(opts, session, recorded, runErr)
}

func persistLiveRecording(opts LiveOptions, session *LiveSession, recorded []RecordedStep, runErr error) error {
	if !session.CaptureEverEnabled() {
		return context.Canceled
	}

	normalized := NormalizeSteps(recorded)
	steps := RecordedStepsToLines(normalized)
	if len(steps) == 0 {
		if errors.Is(runErr, context.Canceled) {
			return context.Canceled
		}
		return nil
	}
	if strings.TrimSpace(opts.AppendTo) != "" {
		if err := AppendStepsToFeature(opts.AppendTo, steps); err != nil {
			return err
		}
	} else if err := WriteFeature(opts.OutputPath, Options{
		FeatureName:  opts.FeatureName,
		ScenarioName: opts.ScenarioName,
		Steps:        steps,
	}); err != nil {
		return err
	}
	return nil
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
	cfg := PageRecorderConfig{
		FilterImportant:   opts.FilterImportant,
		NavOnly:           opts.NavOnly,
		HoverRecord:       opts.HoverRecord,
		ScrollBeforeClick: opts.ScrollBeforeClick,
		HoverRecordMinMs:  opts.HoverRecordMinMs,
	}
	if session != nil {
		cfg = session.RecorderPageConfig()
	}
	return ApplyPageRecorderConfig(page, cfg)
}

// StepsFromFeature extracts plain step texts from a feature (for validation helpers).
func StepsFromFeature(feature *gherkin.Feature) []gherkin.Step {
	out := make([]gherkin.Step, 0)
	for _, scenario := range feature.Scenarios {
		out = append(out, scenario.Steps...)
	}
	return out
}
