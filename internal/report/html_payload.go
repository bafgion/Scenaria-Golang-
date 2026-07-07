package report

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/brand"
	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/runstatus"
)

// HTMLOptions configures interactive HTML report generation.
type HTMLOptions struct {
	Plan            player.ExecutionPlan
	ProjectRoot     string
	LightMode       bool
	BridgeURL       string
	BridgeToken     string
	Locale          string
	ReportDir       string
	MaxJSONBytes    int
	PreviousSummary *RunSummaryDetailed
}

type htmlReportPayload struct {
	Version          string         `json:"version"`
	Brand            string         `json:"brand"`
	GeneratedAt      string         `json:"generated_at"`
	Mode             string         `json:"mode"`
	LightMode        bool           `json:"light_mode"`
	BridgeURL        string         `json:"bridge_url,omitempty"`
	BridgeToken      string         `json:"bridge_token,omitempty"`
	Locale           string         `json:"locale,omitempty"`
	ReportDir        string         `json:"report_dir,omitempty"`
	ArtifactsTrimmed bool           `json:"artifacts_trimmed,omitempty"`
	ScreenshotDedup  int            `json:"screenshot_dedup,omitempty"`
	ModeLinks        htmlModeLinks  `json:"mode_links"`
	CICompare        *htmlCICompare `json:"ci_compare,omitempty"`
	Summary          htmlSummary    `json:"summary"`
	Scenarios        []htmlScenario `json:"scenarios"`
	Flaky            htmlFlaky      `json:"flaky"`
	SlowSteps        []htmlSlowStep `json:"slow_steps"`
}

type htmlModeLinks struct {
	Current        string `json:"current"`
	FullHref       string `json:"full_href,omitempty"`
	LightHref      string `json:"light_href,omitempty"`
	FullAvailable  bool   `json:"full_available"`
	LightAvailable bool   `json:"light_available"`
}

type htmlSummary struct {
	Files     int `json:"files"`
	Scenarios int `json:"scenarios"`
	Steps     int `json:"steps"`
	Passed    int `json:"passed"`
	Failed    int `json:"failed"`
	Skipped   int `json:"skipped"`
}

type htmlScenario struct {
	ID                string             `json:"id"`
	FeaturePath       string             `json:"feature_path"`
	Scenario          string             `json:"scenario"`
	Tags              []string           `json:"tags,omitempty"`
	ExampleIndex      int                `json:"example_index,omitempty"`
	Status            string             `json:"status"`
	Message           string             `json:"message,omitempty"`
	FailedStep        *int               `json:"failed_step,omitempty"`
	DurationMS        int64              `json:"duration_ms"`
	Steps             []htmlStep         `json:"steps"`
	Screenshot        string             `json:"screenshot,omitempty"`
	TracePath         string             `json:"trace_path,omitempty"`
	TraceCommand      string             `json:"trace_command,omitempty"`
	TraceEvents       []htmlTraceEvent   `json:"trace_events,omitempty"`
	RerunCommand      string             `json:"rerun_command,omitempty"`
	History           *htmlHistory       `json:"history,omitempty"`
	HistoryRuns       []htmlHistoryEntry `json:"history_runs,omitempty"`
	Regressions       []htmlRegression   `json:"regressions,omitempty"`
	DurationSparkline []int              `json:"duration_sparkline,omitempty"`
	RunDiff           *htmlRunDiff       `json:"run_diff,omitempty"`
}

type htmlStep struct {
	Index             int      `json:"index"`
	Line              int      `json:"line,omitempty"`
	Keyword           string   `json:"keyword,omitempty"`
	Text              string   `json:"text"`
	Selector          string   `json:"selector,omitempty"`
	Status            string   `json:"status"`
	DurationMS        int64    `json:"duration_ms,omitempty"`
	Error             string   `json:"error,omitempty"`
	Network           string   `json:"network,omitempty"`
	Screenshot        string   `json:"screenshot,omitempty"`
	FlakyFailures     int      `json:"flaky_failures,omitempty"`
	PageContext       string   `json:"page_context,omitempty"`
	DOMSnapshot       string   `json:"dom_snapshot,omitempty"`
	A11ySnapshot      string   `json:"a11y_snapshot,omitempty"`
	TraceOffsetMS     int64    `json:"trace_offset_ms,omitempty"`
	DurationSparkline []int    `json:"duration_sparkline,omitempty"`
	Tips              []string `json:"tips,omitempty"`
	Gherkin           string   `json:"gherkin,omitempty"`
}

type htmlHistory struct {
	LastStatus  string `json:"last_status"`
	LastAt      string `json:"last_at"`
	LastMessage string `json:"last_message,omitempty"`
	Changed     bool   `json:"changed"`
}

type htmlHistoryEntry struct {
	Status        string `json:"status"`
	At            string `json:"at"`
	Message       string `json:"message,omitempty"`
	FailedStep    *int   `json:"failed_step,omitempty"`
	DurationMS    int    `json:"duration_ms,omitempty"`
	StepDurations []int  `json:"step_durations,omitempty"`
}

type htmlCICompare struct {
	PreviousAt     string           `json:"previous_at,omitempty"`
	NewFailures    []htmlCIItem     `json:"new_failures,omitempty"`
	Fixed          []htmlCIItem     `json:"fixed,omitempty"`
	DurationDeltas []htmlCIDuration `json:"duration_deltas,omitempty"`
}

type htmlCIItem struct {
	Path           string `json:"path"`
	Scenario       string `json:"scenario"`
	Status         string `json:"status"`
	PreviousStatus string `json:"previous_status,omitempty"`
}

type htmlCIDuration struct {
	Path       string `json:"path"`
	Scenario   string `json:"scenario"`
	CurrentMS  int64  `json:"current_ms"`
	PreviousMS int64  `json:"previous_ms"`
	DeltaMS    int64  `json:"delta_ms"`
}

type htmlRegression struct {
	Step       int    `json:"step"`
	Text       string `json:"text"`
	Kind       string `json:"kind"`
	PreviousAt string `json:"previous_at,omitempty"`
	Detail     string `json:"detail,omitempty"`
}

type htmlFlaky struct {
	Scenarios []runstatus.ScenarioFlakyStat `json:"scenarios,omitempty"`
	Steps     []runstatus.StepFlakyStat     `json:"steps,omitempty"`
}

type htmlSlowStep struct {
	ScenarioPath string `json:"scenario_path"`
	Step         int    `json:"step"`
	Text         string `json:"text"`
	DurationMS   int64  `json:"duration_ms"`
}

func buildHTMLPayload(result player.ExecutionResult, opts HTMLOptions, reportPath string) (htmlReportPayload, error) {
	payload := htmlReportPayload{
		Version:     "1",
		Brand:       brand.Name,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Mode:        result.Mode,
		LightMode:   opts.LightMode,
		BridgeURL:   strings.TrimSpace(opts.BridgeURL),
		BridgeToken: strings.TrimSpace(opts.BridgeToken),
		Locale:      normalizeReportLocale(opts.Locale),
		ReportDir:   strings.TrimSpace(opts.ReportDir),
		ModeLinks:   buildHTMLModeLinks(reportPath, opts.LightMode),
		Summary: htmlSummary{
			Files:     result.Files,
			Scenarios: result.Scenarios,
			Steps:     result.Steps,
		},
	}

	var history []runstatus.Entry
	if root := strings.TrimSpace(opts.ProjectRoot); root != "" {
		if store, err := runstatus.Open(root); err == nil {
			if entries, err := store.List(200); err == nil {
				history = entries
				scenarios, steps := runstatus.FlakyStats(entries)
				payload.Flaky.Scenarios = scenarios
				payload.Flaky.Steps = steps
			}
		}
	}

	reportDir := filepath.Dir(reportPath)
	tracesDir := filepath.Join(reportDir, "traces")
	screenshotsDir := filepath.Join(reportDir, "screenshots")
	if !opts.LightMode {
		_ = os.MkdirAll(tracesDir, 0o755)
	}
	_ = os.MkdirAll(screenshotsDir, 0o755)

	flakySteps := flakyStepFailures(payload.Flaky.Steps)

	payload.Scenarios = make([]htmlScenario, 0, len(result.ScenarioResults))
	for i, sr := range result.ScenarioResults {
		casePlan := findPlanCase(opts.Plan, sr)
		sc := buildHTMLScenario(sr, casePlan, i, tracesDir, screenshotsDir, opts.LightMode, flakySteps)
		sc.History = lookupHistory(history, sr)
		sc.HistoryRuns = lookupHistoryRuns(history, sr, 5)
		sc.Regressions = computeRegressions(sc)
		sc.DurationSparkline = buildDurationSparkline(sc)
		sc.RunDiff = computeRunDiff(sc)
		attachStepSparklines(&sc)
		payload.Scenarios = append(payload.Scenarios, sc)

		switch sr.Status {
		case "passed":
			payload.Summary.Passed++
		case "failed":
			payload.Summary.Failed++
		default:
			payload.Summary.Skipped++
		}
	}
	payload.SlowSteps = collectSlowSteps(payload.Scenarios, 8)
	payload.CICompare = buildCICompare(payload.Scenarios, opts.PreviousSummary)
	if payload.ReportDir == "" {
		payload.ReportDir = filepath.Dir(reportPath)
	}
	return payload, nil
}

func buildHTMLModeLinks(reportPath string, light bool) htmlModeLinks {
	fullPath, lightPath := htmlModePairPaths(reportPath, light)
	links := htmlModeLinks{
		Current:        "full",
		FullHref:       filepath.Base(fullPath),
		LightHref:      filepath.Base(lightPath),
		FullAvailable:  !light || fileExists(fullPath),
		LightAvailable: light || fileExists(lightPath),
	}
	if light {
		links.Current = "light"
	}
	if !links.FullAvailable {
		links.FullHref = ""
	}
	if !links.LightAvailable {
		links.LightHref = ""
	}
	return links
}

func htmlModePairPaths(reportPath string, light bool) (fullPath, lightPath string) {
	dir := filepath.Dir(reportPath)
	name := filepath.Base(reportPath)
	ext := filepath.Ext(name)
	stem := strings.TrimSuffix(name, ext)
	if ext == "" {
		ext = ".html"
	}
	switch {
	case strings.HasSuffix(stem, ".light"):
		return filepath.Join(dir, strings.TrimSuffix(stem, ".light")+ext), reportPath
	case strings.HasSuffix(stem, ".full"):
		return reportPath, filepath.Join(dir, strings.TrimSuffix(stem, ".full")+".light"+ext)
	case light:
		return filepath.Join(dir, stem+".full"+ext), reportPath
	default:
		return reportPath, filepath.Join(dir, stem+".light"+ext)
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func shouldEmbedScreenshot(light bool, status string, hasPNG bool) bool {
	if !hasPNG {
		return false
	}
	return !light || status == "failed"
}

func findPlanCase(plan player.ExecutionPlan, sr player.ScenarioResult) *player.RunCase {
	for i := range plan.Cases {
		c := &plan.Cases[i]
		if c.FeaturePath == sr.FeaturePath && c.Name == sr.Scenario {
			return c
		}
	}
	return nil
}

func flakyStepFailures(stats []runstatus.StepFlakyStat) map[string]map[int]int {
	out := map[string]map[int]int{}
	for _, stat := range stats {
		if stat.Failures < 2 {
			continue
		}
		if out[stat.Path] == nil {
			out[stat.Path] = map[int]int{}
		}
		out[stat.Path][stat.Step] = stat.Failures
	}
	return out
}

func buildHTMLScenario(sr player.ScenarioResult, casePlan *player.RunCase, index int, tracesDir, screenshotsDir string, light bool, flakySteps map[string]map[int]int) htmlScenario {
	sc := htmlScenario{
		ID:          fmt.Sprintf("s%d", index),
		FeaturePath: sr.FeaturePath,
		Scenario:    sr.Scenario,
		Status:      sr.Status,
		Message:     sr.Message,
		FailedStep:  sr.FailedStep,
		DurationMS:  sr.DurationMS,
	}
	if casePlan != nil {
		sc.Tags = append([]string(nil), casePlan.Tags...)
		if casePlan.ExampleIndex > 0 {
			sc.ExampleIndex = casePlan.ExampleIndex
		}
	}
	scenarioPath := sr.FeaturePath + "::" + sr.Scenario
	traceHints := scenarioHasTrace(sr, light)
	artifactBase := fmt.Sprintf("%03d_%s", index, safeArtifactName(sr.FeaturePath, sr.Scenario))
	sc.Steps = buildHTMLSteps(sr, casePlan, flakySteps[scenarioPath], light, traceHints, screenshotsDir, artifactBase)
	if sc.DurationMS <= 0 {
		sc.DurationMS = sumStepDuration(sc.Steps)
	}
	sc.RerunCommand = fmt.Sprintf("scenaria run %q --scenario %q", sr.FeaturePath, sr.Scenario)

	if shouldEmbedScreenshot(light, sr.Status, len(sr.ScreenshotPNG) > 0) {
		sc.Screenshot = writeScreenshotArtifact(screenshotsDir, artifactBase+"__scenario.png", sr.ScreenshotPNG)
	}
	if !light {
		if len(sr.TraceZIP) > 0 {
			sc.TraceEvents = parseTraceActions(sr.TraceZIP, 80)
			networkFails := parseTraceNetworkFailures(sr.TraceZIP, 24)
			name := safeArtifactName(sr.FeaturePath, sr.Scenario) + ".zip"
			path := filepath.Join(tracesDir, name)
			if err := os.WriteFile(path, sr.TraceZIP, 0o644); err == nil {
				rel, _ := filepath.Rel(filepath.Dir(tracesDir), path)
				if rel == "" || strings.HasPrefix(rel, "..") {
					rel = filepath.Join("traces", name)
				} else {
					rel = filepath.Join("traces", name)
				}
				sc.TracePath = rel
				sc.TraceCommand = fmt.Sprintf("npx playwright show-trace %q", rel)
			}
			if len(networkFails) > 0 {
				enrichStepsNetworkFromTrace(sc.Steps, networkFails)
			}
		}
	}
	return sc
}

func scenarioHasTrace(sr player.ScenarioResult, light bool) bool {
	return !light && sr.Status == "failed" && len(sr.TraceZIP) > 0
}

func buildHTMLSteps(sr player.ScenarioResult, casePlan *player.RunCase, flaky map[int]int, light bool, traceHints bool, screenshotsDir, artifactBase string) []htmlStep {
	if len(sr.StepRecords) > 0 {
		out := make([]htmlStep, 0, len(sr.StepRecords))
		var offset int64
		for _, rec := range sr.StepRecords {
			step := htmlStep{
				Index:        rec.Index,
				Line:         rec.Line,
				Keyword:      rec.Keyword,
				Text:         rec.Text,
				Selector:     rec.Selector,
				Status:       rec.Status,
				DurationMS:   rec.DurationMS,
				Error:        rec.Error,
				Network:      rec.Network,
				PageContext:  rec.PageContext,
				DOMSnapshot:  rec.DOMSnapshot,
				A11ySnapshot: rec.A11ySnapshot,
				Gherkin:      gherkinLine(rec.Keyword, rec.Text),
			}
			if traceHints {
				step.TraceOffsetMS = offset
			}
			offset += rec.DurationMS
			if flaky != nil {
				step.FlakyFailures = flaky[rec.Index]
			}
			if shouldEmbedScreenshot(light, rec.Status, len(rec.ScreenshotPNG) > 0) {
				step.Screenshot = writeScreenshotArtifact(screenshotsDir, fmt.Sprintf("%s__step_%03d.png", artifactBase, rec.Index), rec.ScreenshotPNG)
			}
			step.Tips = stepTips(step.Selector, step.Error, step.Text)
			out = append(out, step)
		}
		return out
	}

	var leaves []gherkin.Step
	if casePlan != nil {
		steps := casePlan.Steps
		if casePlan.StartStep != -1 || casePlan.EndStep != -1 {
			start := casePlan.StartStep
			if start < 0 {
				start = 0
			}
			steps = gherkin.ApplyStepRange(steps, start, casePlan.EndStep)
		}
		leaves = gherkin.LeafSteps(steps)
	}
	out := make([]htmlStep, 0, len(leaves))
	for i, leaf := range leaves {
		status := inferStepStatus(sr.Status, i, sr.FailedStep)
		errText := ""
		if status == "failed" {
			errText = sr.Message
		}
		step := htmlStep{
			Index:   i,
			Line:    leaf.Line,
			Keyword: leaf.Keyword,
			Text:    leaf.Text,
			Status:  status,
			Error:   errText,
			Gherkin: gherkinLine(leaf.Keyword, leaf.Text),
		}
		if flaky != nil {
			step.FlakyFailures = flaky[i]
		}
		step.Tips = stepTips(step.Selector, step.Error, step.Text)
		out = append(out, step)
	}
	return out
}

func inferStepStatus(scenarioStatus string, stepIndex int, failedStep *int) string {
	switch scenarioStatus {
	case "dry-run", "skipped":
		return "skipped"
	case "passed":
		return "passed"
	}
	if failedStep == nil {
		return "unknown"
	}
	fs := *failedStep
	switch {
	case stepIndex < fs:
		return "passed"
	case stepIndex == fs:
		return "failed"
	default:
		return "skipped"
	}
}

func gherkinLine(keyword, text string) string {
	kw := strings.TrimSpace(keyword)
	txt := strings.TrimSpace(text)
	if kw == "" {
		return txt
	}
	if txt == "" {
		return kw
	}
	return kw + " " + txt
}

func stepTips(selector, errText, stepText string) []string {
	tips := make([]string, 0, 3)
	lowerSel := strings.ToLower(selector)
	if selector != "" && !strings.Contains(lowerSel, "data-testid") && !strings.Contains(lowerSel, "test-id") {
		tips = append(tips, "Для стабильности используйте data-testid вместо хрупких CSS/XPath селекторов.")
	}
	if strings.Contains(strings.ToLower(errText), "timeout") || strings.Contains(strings.ToLower(stepText), "жду") {
		tips = append(tips, "Проверьте wait-условие и увеличьте таймаут, если страница грузится медленно.")
	}
	if strings.Contains(strings.ToLower(errText), "strict mode violation") {
		tips = append(tips, "Селектор находит несколько элементов — уточните локатор или используйте .first() / data-testid.")
	}
	return tips
}

func sumStepDuration(steps []htmlStep) int64 {
	var total int64
	for _, s := range steps {
		total += s.DurationMS
	}
	return total
}

func collectSlowSteps(scenarios []htmlScenario, limit int) []htmlSlowStep {
	all := make([]htmlSlowStep, 0)
	for _, sc := range scenarios {
		path := sc.FeaturePath + "::" + sc.Scenario
		for _, step := range sc.Steps {
			if step.DurationMS <= 0 {
				continue
			}
			all = append(all, htmlSlowStep{
				ScenarioPath: path,
				Step:         step.Index,
				Text:         step.Text,
				DurationMS:   step.DurationMS,
			})
		}
	}
	for i := 0; i < len(all); i++ {
		for j := i + 1; j < len(all); j++ {
			if all[j].DurationMS > all[i].DurationMS {
				all[i], all[j] = all[j], all[i]
			}
		}
	}
	if limit > 0 && len(all) > limit {
		all = all[:limit]
	}
	return all
}

func normalizeReportLocale(locale string) string {
	locale = strings.ToLower(strings.TrimSpace(locale))
	if strings.HasPrefix(locale, "en") {
		return "en"
	}
	return "ru"
}

func computeRegressions(sc htmlScenario) []htmlRegression {
	if sc.Status != "failed" || len(sc.HistoryRuns) == 0 {
		return nil
	}
	prev := sc.HistoryRuns[0]
	failIdx := -1
	var failText string
	for _, st := range sc.Steps {
		if st.Status == "failed" {
			failIdx = st.Index
			failText = st.Text
			if failText == "" {
				failText = st.Gherkin
			}
			break
		}
	}
	if failIdx < 0 && sc.FailedStep != nil {
		failIdx = *sc.FailedStep
		for _, st := range sc.Steps {
			if st.Index == failIdx {
				failText = st.Text
				if failText == "" {
					failText = st.Gherkin
				}
				break
			}
		}
	}
	out := make([]htmlRegression, 0, 2)
	if prev.Status == "passed" {
		return append(out, htmlRegression{
			Step: failIdx, Text: failText, Kind: "new_failure",
			PreviousAt: prev.At, Detail: "scenario passed before",
		})
	}
	if prev.Status != "failed" || failIdx < 0 {
		return out
	}
	prevStep := -1
	if prev.FailedStep != nil {
		prevStep = *prev.FailedStep
	}
	if prevStep != failIdx {
		return append(out, htmlRegression{
			Step: failIdx, Text: failText, Kind: "step_changed",
			PreviousAt: prev.At,
			Detail:     fmt.Sprintf("was step #%d, now #%d", prevStep+1, failIdx+1),
		})
	}
	return append(out, htmlRegression{
		Step: failIdx, Text: failText, Kind: "still_failing",
		PreviousAt: prev.At, Detail: "same failing step",
	})
}

func lookupHistoryRuns(entries []runstatus.Entry, sr player.ScenarioResult, limit int) []htmlHistoryEntry {
	if limit <= 0 {
		limit = 5
	}
	path := sr.FeaturePath + "::" + sr.Scenario
	out := make([]htmlHistoryEntry, 0, limit)
	for _, entry := range entries {
		if entry.Path != path {
			continue
		}
		status := "passed"
		if !entry.Success {
			status = "failed"
		}
		out = append(out, htmlHistoryEntry{
			Status:        status,
			At:            entry.At,
			Message:       entry.Message,
			FailedStep:    entry.FailedStep,
			DurationMS:    entry.DurationMS,
			StepDurations: append([]int(nil), entry.StepDurations...),
		})
		if len(out) >= limit {
			break
		}
	}
	return out
}

func lookupHistory(entries []runstatus.Entry, sr player.ScenarioResult) *htmlHistory {
	path := sr.FeaturePath + "::" + sr.Scenario
	var prev *runstatus.Entry
	for _, entry := range entries {
		if entry.Path != path {
			continue
		}
		prev = &entry
		break
	}
	if prev == nil {
		return nil
	}
	status := "passed"
	if !prev.Success {
		status = "failed"
	}
	changed := (status == "passed") != (sr.Status == "passed")
	return &htmlHistory{
		LastStatus:  status,
		LastAt:      prev.At,
		LastMessage: prev.Message,
		Changed:     changed,
	}
}

func safeArtifactName(featurePath, scenario string) string {
	base := strings.TrimSuffix(filepath.Base(featurePath), filepath.Ext(featurePath))
	base = sanitizeName(base)
	sc := sanitizeName(scenario)
	if base == "" {
		base = "feature"
	}
	if sc == "" {
		sc = "scenario"
	}
	return base + "__" + sc
}

func sanitizeName(s string) string {
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.' {
			b.WriteRune(r)
		} else if r == ' ' {
			b.WriteRune('_')
		}
	}
	out := strings.Trim(b.String(), "._-")
	if out == "" {
		return "item"
	}
	return out
}

func buildDurationSparkline(sc htmlScenario) []int {
	runs := sc.HistoryRuns
	if len(runs) == 0 && sc.DurationMS <= 0 {
		return nil
	}
	vals := make([]int, 0, len(runs)+1)
	for i := len(runs) - 1; i >= 0; i-- {
		if runs[i].DurationMS > 0 {
			vals = append(vals, runs[i].DurationMS)
		}
	}
	if sc.DurationMS > 0 {
		vals = append(vals, int(sc.DurationMS))
	}
	if len(vals) < 2 {
		return nil
	}
	return vals
}

func buildCICompare(scenarios []htmlScenario, prev *RunSummaryDetailed) *htmlCICompare {
	if prev == nil || len(prev.Items) == 0 {
		return nil
	}
	prevByKey := make(map[string]ScenarioSummary, len(prev.Items))
	for _, item := range prev.Items {
		prevByKey[item.Path+"::"+item.Scenario] = item
	}
	curByKey := make(map[string]htmlScenario, len(scenarios))
	for _, sc := range scenarios {
		curByKey[sc.FeaturePath+"::"+sc.Scenario] = sc
	}
	out := &htmlCICompare{PreviousAt: prev.GeneratedAt}
	for key, sc := range curByKey {
		prevItem, ok := prevByKey[key]
		if !ok {
			continue
		}
		if prevItem.Status == "passed" && sc.Status == "failed" {
			out.NewFailures = append(out.NewFailures, htmlCIItem{
				Path: sc.FeaturePath, Scenario: sc.Scenario, Status: sc.Status, PreviousStatus: prevItem.Status,
			})
		}
		if prevItem.Status == "failed" && sc.Status == "passed" {
			out.Fixed = append(out.Fixed, htmlCIItem{
				Path: sc.FeaturePath, Scenario: sc.Scenario, Status: sc.Status, PreviousStatus: prevItem.Status,
			})
		}
		if prevItem.DurationMS > 0 && sc.DurationMS > 0 {
			delta := sc.DurationMS - prevItem.DurationMS
			if delta < 0 {
				delta = -delta
			}
			if delta >= 500 {
				out.DurationDeltas = append(out.DurationDeltas, htmlCIDuration{
					Path: sc.FeaturePath, Scenario: sc.Scenario,
					CurrentMS: sc.DurationMS, PreviousMS: prevItem.DurationMS,
					DeltaMS: sc.DurationMS - prevItem.DurationMS,
				})
			}
		}
	}
	for key, prevItem := range prevByKey {
		if _, ok := curByKey[key]; ok {
			continue
		}
		if prevItem.Status == "failed" {
			out.Fixed = append(out.Fixed, htmlCIItem{
				Path: prevItem.Path, Scenario: prevItem.Scenario, Status: "skipped", PreviousStatus: prevItem.Status,
			})
		}
	}
	if len(out.NewFailures) == 0 && len(out.Fixed) == 0 && len(out.DurationDeltas) == 0 {
		return nil
	}
	return out
}
