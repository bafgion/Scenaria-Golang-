package gui

// EditorAnalysisService owns editor-only analysis operations.
type EditorAnalysisService struct{}

func NewEditorAnalysisService() *EditorAnalysisService {
	return &EditorAnalysisService{}
}

func (s *EditorAnalysisService) ValidateFeature(text string) []ValidationIssue {
	return ValidateFeatureContent(text)
}

func (s *EditorAnalysisService) ParseEditorSteps(text string) []EditorStepRow {
	return ParseEditorSteps(text)
}

func (s *EditorAnalysisService) AnalyzeEditorContent(text string, includeHints bool) EditorAnalysisDTO {
	return AnalyzeEditorContent(text, includeHints)
}

func (s *EditorAnalysisService) AnalyzeScenarioHints(text string) []ScenarioHintDTO {
	return AnalyzeScenarioHints(text)
}

func (s *EditorAnalysisService) ApplyScenarioHintFix(req ScenarioHintFixRequest) RefactorResult {
	return ApplyScenarioHintFix(req)
}

func (s *EditorAnalysisService) ResolveRunFromLine(text string, line int) (RunFromLineDTO, error) {
	return ResolveRunFromLine(text, line)
}

func (s *EditorAnalysisService) ResolveRunToLine(text string, line int) (RunFromLineDTO, error) {
	return ResolveRunToLine(text, line)
}
