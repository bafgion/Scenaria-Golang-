package gui

type RecordRequest struct {
	URL              string `json:"url"`
	Output           string `json:"output"`
	IdleSeconds      int    `json:"idleSeconds"`
	Headless         bool   `json:"headless"`
	FilterRecording  bool   `json:"filterRecording"`
	NavOnlyRecording bool   `json:"navOnlyRecording"`
	HoverRecord      bool   `json:"hoverRecord"`
	AppendTo         string `json:"appendTo"`
	TestClient       string `json:"testClient"`
	FeatureName      string `json:"featureName"`
	ScenarioName     string `json:"scenarioName"`
	BrowseOnly       bool   `json:"browseOnly"`
}

// OpenBrowserRequest opens a Playwright browser without recording until capture is started explicitly.
type OpenBrowserRequest struct {
	URL              string `json:"url"`
	Headless         bool   `json:"headless"`
	TestClient       string `json:"testClient"`
	Output           string `json:"output"`
	IdleSeconds      int    `json:"idleSeconds"`
	FilterRecording  bool   `json:"filterRecording"`
	NavOnlyRecording bool   `json:"navOnlyRecording"`
	HoverRecord      bool   `json:"hoverRecord"`
	AppendTo         string `json:"appendTo"`
	FeatureName      string `json:"featureName"`
	ScenarioName     string `json:"scenarioName"`
}

type BaselineRecordRequest struct {
	Output       string   `json:"output"`
	FeatureName  string   `json:"featureName"`
	ScenarioName string   `json:"scenarioName"`
	Steps        []string `json:"steps"`
}

type ExportRequest struct {
	InputPath string `json:"inputPath"`
	Output    string `json:"output"`
	Format    string `json:"format"`
	BaseURL   string `json:"baseURL"`
	Force     bool   `json:"force"`
}

type ImportRequest struct {
	JSONPath   string `json:"jsonPath"`
	OutputPath string `json:"outputPath"`
	Force      bool   `json:"force"`
}

func (s *Service) ValidateFeature(text string) []ValidationIssue {
	return s.editorAnalyzer().ValidateFeature(text)
}

func (s *Service) Export(req ExportRequest) RunResult {
	return s.cliRunner().ExportFeature(req)
}

func (s *Service) ImportJSON(req ImportRequest) RunResult {
	return s.cliRunner().ImportFeatureJSON(req)
}

func (s *Service) OpenBrowser(req OpenBrowserRequest, emit func(string, any)) RunResult {
	return s.RecordLive(RecordRequest{
		URL:              req.URL,
		Output:           req.Output,
		IdleSeconds:      req.IdleSeconds,
		Headless:         req.Headless,
		FilterRecording:  req.FilterRecording,
		NavOnlyRecording: req.NavOnlyRecording,
		HoverRecord:      req.HoverRecord,
		AppendTo:         req.AppendTo,
		TestClient:       req.TestClient,
		FeatureName:      req.FeatureName,
		ScenarioName:     req.ScenarioName,
		BrowseOnly:       true,
	}, emit)
}

func (s *Service) recordExecutionHost() RecorderExecutionHost {
	return RecorderExecutionHost{
		ProjectPath:     s.ProjectPath,
		LoadAppSettings: s.loadAppSettings,
		SaveAppSettings: s.saveAppSettings,
		WithActivePlaywright: func(fn func()) {
			s.activePlaywright.Add(1)
			defer s.activePlaywright.Done()
			fn()
		},
	}
}

func (s *Service) guardedRecordEmit(gen uint64, targetPath string, emit func(string, any)) func(string, any) {
	return s.recorderOps().Session().GuardedEmit(gen, targetPath, emit)
}

func (s *Service) HasLiveBrowser() bool {
	return s.recorderOps().HasLiveBrowser()
}

func (s *Service) RecordLive(req RecordRequest, emit func(string, any)) RunResult {
	return s.recorderOps().RecordLive(req, emit, s.recordExecutionHost())
}

func (s *Service) BeginRecordingCapture() (bool, error) {
	return s.recorderOps().BeginRecordingCaptureWithReplay()
}

func (s *Service) RecordBaseline(req BaselineRecordRequest) RunResult {
	return s.cliRunner().RecordBaselineFeature(s.ProjectPath(), req)
}

type BrowserSessionDTO struct {
	BrowserOpen      bool   `json:"browserOpen"`
	Recording        bool   `json:"recording"`
	Paused           bool   `json:"paused"`
	StepCount        int    `json:"stepCount"`
	BrowserSessionID string `json:"browserSessionId,omitempty"`
}

func (s *Service) PollBrowserSession() BrowserSessionDTO {
	sessionMgr := s.recorderOps().Session()
	return s.recorderOps().PollBrowserSession(sessionMgr.LiveSession(), sessionMgr.BrowserSessionID())
}

func (s *Service) PauseRecording() {
	s.recorderOps().PauseRecording(s.recorderOps().LiveSession())
}

func (s *Service) ResumeRecording() {
	s.recorderOps().ResumeRecording(s.recorderOps().LiveSession())
}

func (s *Service) IsRecordingPaused() bool {
	session := s.recorderOps().LiveSession()
	if session == nil {
		return false
	}
	return session.IsPaused()
}

func (s *Service) CurrentRecordSessionID() string {
	return s.recorderOps().Session().RecordSessionID()
}

func (s *Service) CurrentBrowserSessionID() string {
	return s.recorderOps().Session().BrowserSessionID()
}

func (s *Service) CurrentRecordTargetPath() string {
	return s.recorderOps().Session().TargetPath()
}

func (s *Service) LastClosedBrowserSessionID() string {
	return s.recorderOps().Session().LastClosedBrowserSessionID()
}

func (s *Service) LastClosedRecordSessionID() string {
	return s.recorderOps().Session().LastClosedRecordSessionID()
}

func (s *Service) StopRecordingCapture() (bool, error) {
	return s.recorderOps().StopRecordingCaptureWithReset()
}

func (s *Service) CloseBrowser() {
	s.recorderOps().Session().CloseBrowser(false)
}

func (s *Service) closeBrowserForced() {
	s.recorderOps().Session().CloseBrowser(true)
}

func (s *Service) CancelRecording() {
	s.closeBrowserForced()
}

func (s *Service) FocusBrowser() error {
	return s.recorderOps().FocusBrowser(s.recorderOps().LiveSession())
}

func (s *Service) UpdateRecordingOptions(filter, navOnly, hover, headless, scrollBefore bool, hoverMinMs int) error {
	return s.recorderOps().UpdateRecordingOptions(s.recorderOps().LiveSession(), filter, navOnly, hover, headless, scrollBefore, hoverMinMs)
}

func (s *Service) UndoRecordedStep() bool {
	return s.recorderOps().UndoRecordedStepWithEmit()
}
