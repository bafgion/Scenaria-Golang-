package recorder

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/paths"
	playwright "github.com/mxschmitt/playwright-go"
)

type LiveOptions struct {
	StartURL     string
	FeatureName  string
	ScenarioName string
	OutputPath   string
	Headless     bool
	IdleTimeout  time.Duration
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

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(opts.Headless),
	})
	if err != nil {
		return fmt.Errorf("launch browser: %w", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return fmt.Errorf("create page: %w", err)
	}
	if _, err := page.Goto(opts.StartURL); err != nil {
		return fmt.Errorf("goto start URL: %w", err)
	}
	if _, err := page.Evaluate(RecorderScript); err != nil {
		return fmt.Errorf("inject recorder script: %w", err)
	}

	steps := []string{fmt.Sprintf(`открыт "%s"`, opts.StartURL)}
	lastEventAt := time.Now()
	seen := map[string]struct{}{}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if time.Since(lastEventAt) >= opts.IdleTimeout {
			break
		}
		raw, err := page.Evaluate(`() => {
			const r = window.__scenariaRecorder;
			if (!r || !r.events.length) return [];
			const out = r.events.splice(0, r.events.length);
			return out;
		}`)
		if err != nil {
			return fmt.Errorf("read recorder events: %w", err)
		}
		events, err := decodeEvents(raw)
		if err != nil {
			return err
		}
		if len(events) > 0 {
			lastEventAt = time.Now()
		}
		for _, event := range events {
			detail := normalizeDetail(event.Detail)
			step, ok := EventsToStep(event.Type, detail)
			if !ok {
				continue
			}
			if _, exists := seen[step]; exists {
				continue
			}
			seen[step] = struct{}{}
			steps = append(steps, step)
		}
		time.Sleep(200 * time.Millisecond)
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

// StepsFromFeature extracts plain step texts from a feature (for validation helpers).
func StepsFromFeature(feature *gherkin.Feature) []gherkin.Step {
	out := make([]gherkin.Step, 0)
	for _, scenario := range feature.Scenarios {
		out = append(out, scenario.Steps...)
	}
	return out
}
