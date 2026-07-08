package gui

import (
	"time"
)

// ProjectArtifacts lists artifact directories for the open project.
type ProjectArtifacts struct {
	AllureDir   string `json:"allureDir"`
	TracesDir   string `json:"tracesDir"`
	VideosDir   string `json:"videosDir"`
	HTMLReport  string `json:"htmlReport"`
	JUnitReport string `json:"junitReport"`
	SummaryJSON string `json:"summaryJson"`
}

func (s *Service) ScenariaArtifactPath(sub string) string {
	return s.reporter().ScenariaArtifactPath(sub)
}

func (s *Service) ProjectArtifacts() ProjectArtifacts {
	return s.reporter().ProjectArtifacts()
}

func (s *Service) ParseEditorSteps(text string) []EditorStepRow {
	return s.editorAnalyzer().ParseEditorSteps(text)
}

func (s *Service) AnalyzeEditorContent(text string, includeHints bool) EditorAnalysisDTO {
	started := time.Now()
	defer logWailsTiming("AnalyzeEditorContent", started)
	return s.editorAnalyzer().AnalyzeEditorContent(text, includeHints)
}

func (s *Service) AnalyzeScenarioHints(text string) []ScenarioHintDTO {
	return s.editorAnalyzer().AnalyzeScenarioHints(text)
}

func (s *Service) ApplyScenarioHintFix(req ScenarioHintFixRequest) RefactorResult {
	return s.editorAnalyzer().ApplyScenarioHintFix(req)
}

func (s *Service) ResolveRunFromLine(text string, line int) (RunFromLineDTO, error) {
	return s.editorAnalyzer().ResolveRunFromLine(text, line)
}

func (s *Service) ResolveRunToLine(text string, line int) (RunFromLineDTO, error) {
	return s.editorAnalyzer().ResolveRunToLine(text, line)
}

func (s *Service) ArtifactExists(path string) bool {
	return s.reporter().ArtifactExists(path)
}
