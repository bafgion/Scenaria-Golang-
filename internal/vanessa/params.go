package vanessa

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/scenario"
)

type RunRequest struct {
	ProjectRoot         string
	Paths               []string
	Tag                 string
	ExcludeTags         []string
	ScenarioNames       []string
	DryRun              bool
	PlatformExecutable  string
	EPFPath             string
	IBConnection        string
	ReportAllure        bool
}

func MergeVAParams(cfg Settings, req RunRequest, runDir string) (map[string]any, string, error) {
	base := map[string]any{}
	if basePath := BaseParamsPath(req.ProjectRoot, cfg); basePath != "" {
		if payload, err := os.ReadFile(basePath); err == nil {
			_ = json.Unmarshal(payload, &base)
		}
	}

	junitDir := filepath.Join(runDir, "junit")
	_ = os.MkdirAll(junitDir, 0o755)
	scenarioLog := filepath.Join(runDir, "scenario.log")
	statusPath := filepath.Join(runDir, "status.log")

	overlay := map[string]any{
		"ДелатьЛогВыполненияСценариевВТекстовыйФайл":                true,
		"ИмяФайлаЛогВыполненияСценариев":                            pathForVA(scenarioLog),
		"ВыгружатьСтатусВыполненияСценариевВФайл":                   true,
		"ПутьКФайлуДляВыгрузкиСтатусаВыполненияСценариев":           pathForVA(statusPath),
		"ДелатьОтчетВФорматеjUnit":                                  cfg.ReportJUnit,
		"КаталогВыгрузкиJUnit":                                      pathForVA(junitDir),
	}
	if cfg.ReportAllure {
		allureDir := filepath.Join(runDir, "allure")
		_ = os.MkdirAll(allureDir, 0o755)
		overlay["ДелатьОтчетВФорматеАллюр"] = true
		overlay["КаталогВыгрузкиAllure"] = pathForVA(allureDir)
	}
	if req.ProjectRoot != "" {
		overlay["КаталогПроекта"] = pathForVA(req.ProjectRoot)
		featuresRoot := filepath.Join(req.ProjectRoot, "features")
		if info, err := os.Stat(featuresRoot); err == nil && info.IsDir() {
			overlay["КаталогФич"] = pathForVA(featuresRoot)
		}
	}
	if epf := strings.TrimSpace(cfg.EPFPath); epf != "" {
		overlay["КаталогИнструментов"] = pathForVA(filepath.Dir(epf))
	}
	if req.Tag != "" {
		tag := strings.TrimPrefix(strings.TrimSpace(req.Tag), "@")
		overlay["СписокТеговОтбор"] = []string{"@" + tag}
	}
	if len(req.ExcludeTags) > 0 {
		tags := make([]string, 0, len(req.ExcludeTags))
		for _, raw := range req.ExcludeTags {
			tag := strings.TrimPrefix(strings.TrimSpace(raw), "@")
			if tag == "" {
				continue
			}
			tags = append(tags, "@"+tag)
		}
		if len(tags) > 0 {
			overlay["СписокТеговИсключение"] = tags
		}
	}
	files, err := resolveFeatureFiles(req)
	if err != nil {
		return nil, "", err
	}
	if len(files) == 1 && isDir(files[0]) {
		overlay["КаталогФич"] = pathForVA(files[0])
	} else if len(files) > 0 {
		paths := make([]string, 0, len(files))
		for _, file := range files {
			paths = append(paths, pathForVA(file))
		}
		overlay["СписокФичДляВыполнения"] = paths
	}
	if len(req.ScenarioNames) > 0 {
		overlay["СписокСценариевДляВыполнения"] = req.ScenarioNames
	}

	merged := deepMerge(base, overlay)
	vaPath := filepath.Join(runDir, "VAParams.json")
	payload, err := json.MarshalIndent(merged, "", "  ")
	if err != nil {
		return nil, "", err
	}
	if err := os.WriteFile(vaPath, append(payload, '\n'), 0o644); err != nil {
		return nil, "", err
	}
	return merged, vaPath, nil
}

func resolveFeatureFiles(req RunRequest) ([]string, error) {
	store := scenario.NewFeatureStore()
	out := make([]string, 0)
	for _, raw := range req.Paths {
		if info, err := os.Stat(raw); err == nil && info.IsDir() {
			files, err := store.Discover(raw)
			if err != nil {
				return nil, err
			}
			out = append(out, files...)
			continue
		}
		out = append(out, raw)
	}
	if req.Tag == "" {
		return out, nil
	}
	filtered := make([]string, 0, len(out))
	for _, path := range out {
		feature, err := store.Load(path)
		if err != nil {
			continue
		}
		if gherkin.FeatureHasTag(feature, req.Tag) {
			filtered = append(filtered, path)
		}
	}
	return filtered, nil
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func pathForVA(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	return strings.ReplaceAll(abs, "/", "\\")
}

func deepMerge(base, overlay map[string]any) map[string]any {
	out := map[string]any{}
	for key, value := range base {
		out[key] = value
	}
	for key, value := range overlay {
		out[key] = value
	}
	return out
}
