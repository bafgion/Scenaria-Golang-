package settings

const (
	EditorWordWrapOn  = "on"
	EditorWordWrapOff = "off"

	EditorLineNumbersOn       = "on"
	EditorLineNumbersOff      = "off"
	EditorLineNumbersRelative = "relative"

	EditorWhitespaceNone      = "none"
	EditorWhitespaceBoundary  = "boundary"
	EditorWhitespaceSelection = "selection"
	EditorWhitespaceTrailing  = "trailing"
	EditorWhitespaceAll       = "all"

	EditorThemeDark  = "scenaria-dark"
	EditorThemeLight = "scenaria-light"
)

type EditorSettings struct {
	FontSize          int    `json:"fontSize,omitempty"`
	FontFamily        string `json:"fontFamily,omitempty"`
	WordWrap          string `json:"wordWrap,omitempty"`
	Minimap           *bool  `json:"minimap,omitempty"`
	LineNumbers       string `json:"lineNumbers,omitempty"`
	TabSize           int    `json:"tabSize,omitempty"`
	InsertSpaces      *bool  `json:"insertSpaces,omitempty"`
	RenderWhitespace  string `json:"renderWhitespace,omitempty"`
	Folding           *bool  `json:"folding,omitempty"`
	StickyScroll      *bool  `json:"stickyScroll,omitempty"`
	AutoClosingQuotes string `json:"autoClosingQuotes,omitempty"`
	FormatOnSave      *bool  `json:"formatOnSave,omitempty"`
	StepHoverEnabled  *bool  `json:"stepHoverEnabled,omitempty"`
	ValidateOnType    *bool  `json:"validateOnType,omitempty"`
	Theme             string `json:"theme,omitempty"`
	BreadcrumbsEnabled *bool `json:"breadcrumbsEnabled,omitempty"`
	SymbolOutlineEnabled *bool `json:"symbolOutlineEnabled,omitempty"`
	StepsPanelView    string `json:"stepsPanelView,omitempty"`
	CodeLensEnabled   *bool  `json:"codeLensEnabled,omitempty"`
	InlayHintsEnabled *bool  `json:"inlayHintsEnabled,omitempty"`
	ScenarioHintsEnabled      *bool `json:"scenarioHintsEnabled,omitempty"`
	ScenarioHintsAfterRecord  *bool `json:"scenarioHintsAfterRecord,omitempty"`
	ScenarioHintsShowWarnings *bool `json:"scenarioHintsShowWarnings,omitempty"`
	ScenarioHintsShowInfo     *bool `json:"scenarioHintsShowInfo,omitempty"`
	ScenarioHintsAutoFixOnSave *bool `json:"scenarioHintsAutoFixOnSave,omitempty"`
}

func DefaultEditorSettings() EditorSettings {
	minimap := true
	insertSpaces := false
	folding := false
	stickyScroll := false
	formatOnSave := false
	stepHover := true
	validateOnType := true
	breadcrumbs := true
	symbolOutline := true
	codeLens := true
	inlayHints := false
	scenarioHints := true
	scenarioHintsAfterRecord := true
	scenarioHintsShowWarnings := true
	scenarioHintsShowInfo := true
	scenarioHintsAutoFixOnSave := false
	return EditorSettings{
		FontSize:          13,
		FontFamily:        `"Cascadia Code", "JetBrains Mono", Consolas, "Courier New", monospace`,
		WordWrap:          EditorWordWrapOn,
		Minimap:           &minimap,
		LineNumbers:       EditorLineNumbersOn,
		TabSize:           4,
		InsertSpaces:      &insertSpaces,
		RenderWhitespace:  EditorWhitespaceSelection,
		Folding:           &folding,
		StickyScroll:      &stickyScroll,
		AutoClosingQuotes: "languageDefined",
		FormatOnSave:      &formatOnSave,
		StepHoverEnabled:  &stepHover,
		ValidateOnType:    &validateOnType,
		Theme:             EditorThemeDark,
		BreadcrumbsEnabled: &breadcrumbs,
		SymbolOutlineEnabled: &symbolOutline,
		StepsPanelView:    "outline",
		CodeLensEnabled:            &codeLens,
		InlayHintsEnabled:          &inlayHints,
		ScenarioHintsEnabled:       &scenarioHints,
		ScenarioHintsAfterRecord:   &scenarioHintsAfterRecord,
		ScenarioHintsShowWarnings:  &scenarioHintsShowWarnings,
		ScenarioHintsShowInfo:      &scenarioHintsShowInfo,
		ScenarioHintsAutoFixOnSave: &scenarioHintsAutoFixOnSave,
	}
}

func NormalizeEditorSettings(raw EditorSettings) EditorSettings {
	def := DefaultEditorSettings()
	out := def

	if raw.FontSize >= 8 && raw.FontSize <= 32 {
		out.FontSize = raw.FontSize
	}
	if raw.FontFamily != "" {
		out.FontFamily = raw.FontFamily
	}
	switch raw.WordWrap {
	case EditorWordWrapOn, EditorWordWrapOff:
		out.WordWrap = raw.WordWrap
	}
	if raw.Minimap != nil {
		out.Minimap = raw.Minimap
	}
	switch raw.LineNumbers {
	case EditorLineNumbersOn, EditorLineNumbersOff, EditorLineNumbersRelative:
		out.LineNumbers = raw.LineNumbers
	}
	if raw.TabSize >= 1 && raw.TabSize <= 8 {
		out.TabSize = raw.TabSize
	}
	if raw.InsertSpaces != nil {
		out.InsertSpaces = raw.InsertSpaces
	}
	switch raw.RenderWhitespace {
	case EditorWhitespaceNone, EditorWhitespaceBoundary, EditorWhitespaceSelection, EditorWhitespaceTrailing, EditorWhitespaceAll:
		out.RenderWhitespace = raw.RenderWhitespace
	}
	if raw.Folding != nil {
		out.Folding = raw.Folding
	}
	if raw.StickyScroll != nil {
		out.StickyScroll = raw.StickyScroll
	}
	switch raw.AutoClosingQuotes {
	case "always", "languageDefined", "beforeWhitespace", "never":
		out.AutoClosingQuotes = raw.AutoClosingQuotes
	}
	if raw.FormatOnSave != nil {
		out.FormatOnSave = raw.FormatOnSave
	}
	if raw.StepHoverEnabled != nil {
		out.StepHoverEnabled = raw.StepHoverEnabled
	}
	if raw.ValidateOnType != nil {
		out.ValidateOnType = raw.ValidateOnType
	}
	switch raw.Theme {
	case EditorThemeDark, EditorThemeLight:
		out.Theme = raw.Theme
	}
	if raw.BreadcrumbsEnabled != nil {
		out.BreadcrumbsEnabled = raw.BreadcrumbsEnabled
	}
	if raw.SymbolOutlineEnabled != nil {
		out.SymbolOutlineEnabled = raw.SymbolOutlineEnabled
	}
	switch raw.StepsPanelView {
	case "outline", "steps":
		out.StepsPanelView = raw.StepsPanelView
	}
	if raw.CodeLensEnabled != nil {
		out.CodeLensEnabled = raw.CodeLensEnabled
	}
	if raw.InlayHintsEnabled != nil {
		out.InlayHintsEnabled = raw.InlayHintsEnabled
	}
	if raw.ScenarioHintsEnabled != nil {
		out.ScenarioHintsEnabled = raw.ScenarioHintsEnabled
	}
	if raw.ScenarioHintsAfterRecord != nil {
		out.ScenarioHintsAfterRecord = raw.ScenarioHintsAfterRecord
	}
	if raw.ScenarioHintsShowWarnings != nil {
		out.ScenarioHintsShowWarnings = raw.ScenarioHintsShowWarnings
	}
	if raw.ScenarioHintsShowInfo != nil {
		out.ScenarioHintsShowInfo = raw.ScenarioHintsShowInfo
	}
	if raw.ScenarioHintsAutoFixOnSave != nil {
		out.ScenarioHintsAutoFixOnSave = raw.ScenarioHintsAutoFixOnSave
	}
	return out
}
