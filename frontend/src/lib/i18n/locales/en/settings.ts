export const settings = {
  title: 'Settings',
  searchPlaceholder: 'Search settings',
  sectionsNav: 'Settings sections',
  tabs: {
    record: 'Recording & browser',
    selectors: 'Selectors',
    plugins: 'Plugins',
    editor: 'Editor',
    ui: 'Interface',
  },
  sections: {
    browser: {
      title: 'Browser',
      desc: 'Window and session behavior during Playwright recording and runs.',
    },
    recording: {
      title: 'Step recording',
      desc: 'Filters when recording actions in the browser.',
    },
    run: {
      title: 'Run',
    },
    selectors: {
      title: 'Strategy priority',
      desc: 'When recording and picking selectors, {brand} tries strategies top to bottom. More stable ones should be higher.',
    },
    plugins: {
      title: "Runners and add-ons",
      desc: 'Plugins are installed in addons/<name>/ and registered in .scenaria/plugins.json.',
    },
    editorFont: {
      title: 'Font and display',
      desc: 'Monaco scenario editor options.',
    },
    editorInput: {
      title: 'Input and behavior',
    },
    editorNavigation: {
      title: 'Navigation',
      desc: 'Feature file structure in breadcrumbs and the steps panel.',
    },
    editorHints: {
      title: 'Scenario hints',
      desc: 'Step quality heuristics in the editor (markers and quick fix).',
    },
    uiLanguage: {
      title: 'Interface language',
      desc: 'Language for menus, dialogs, and labels.',
    },
    toolbar: {
      title: 'Toolbar',
      desc: 'Appearance of the top action bar.',
    },
    stepsPanel: {
      title: 'Steps panel',
      desc: 'Parsed step list below the scenario editor.',
    },
    updates: {
      title: 'Updates',
    },
  },
  cards: {
    headless: {
      title: 'Headless browser',
      description: 'Headless — the window is hidden during recording and runs.',
    },
    browserEngine: {
      title: 'Browser engine',
      description: 'Playwright: chromium, firefox, or webkit.',
    },
    importantOnly: {
      title: 'Important only',
      description: 'Skip secondary events while recording.',
    },
    linksOnly: {
      title: 'Links only',
      description: 'Record link navigation without element clicks.',
    },
    hoverRecord: {
      title: 'Record hover',
      description: 'Add steps when the cursor hovers over elements.',
    },
    hoverMin: {
      title: 'Minimum hover',
      description: 'How many milliseconds the cursor must stay on an element before a hover step is recorded.',
    },
    scrollBeforeClick: {
      title: 'Scroll before click',
      description: 'Scroll the element into view before recording a click (same as during playback).',
    },
    pickerDuringRecording: {
      title: 'Picker during recording',
      description: 'Allow “Pick element” without pausing recording (pause is required by default).',
    },
    workers: {
      title: 'Parallel workers',
      description: 'Number of simultaneous browser sessions during batch runs.',
    },
    slowMo: {
      title: 'Execution speed (slow-mo)',
      description: 'Pause between Playwright actions in milliseconds. 0 — fastest; 100–300 — easy to watch steps in the browser.',
    },
    navWait: {
      title: 'Navigation wait (nav-wait-until)',
      description: 'When to consider a URL navigation complete for “Go to” steps and recording.',
    },
    loops: {
      title: 'Loop iteration limit',
      description: 'Maximum repeats for “Repeat” / “While” blocks.',
    },
    fontSize: {
      title: 'Font size',
      description: 'From 8 to 32 px.',
    },
    fontFamily: {
      title: 'Font',
      description: 'Monospace editor font.',
    },
    theme: {
      title: 'Theme',
      description: 'Dark, light, or follow the system.',
    },
    wordWrap: {
      title: 'Word wrap',
      description: 'Wrap long steps to the editor width.',
    },
    minimap: {
      title: 'Minimap',
      description: 'Code overview map on the right.',
    },
    lineNumbers: {
      title: 'Line numbers',
      description: 'Show line numbers in the gutter.',
    },
    renderWhitespace: {
      title: 'Whitespace',
      description: 'When to show invisible characters.',
    },
    tabSize: {
      title: 'Tab size',
      description: 'Tab indent width in spaces.',
    },
    insertSpaces: {
      title: 'Spaces instead of Tab',
      description: 'Insert spaces when pressing Tab.',
    },
    folding: {
      title: 'Block folding',
      description: 'Collapse “If” / “Repeat” blocks.',
    },
    stickyScroll: {
      title: 'Sticky scroll',
      description: 'Pin scenario headers while scrolling.',
    },
    autoClosingQuotes: {
      title: 'Auto-close quotes',
      description: 'Behavior when typing quotes.',
    },
    formatOnSave: {
      title: 'Format on save',
      description: 'Normalize indentation and remove extra blank lines between steps on Ctrl+S.',
    },
    stepHover: {
      title: 'Hover hints',
      description: 'Show step help on hover in the editor.',
    },
    validateOnType: {
      title: 'Validate while typing',
      description: 'Validate the scenario with a delay while editing.',
    },
    breadcrumbs: {
      title: 'Breadcrumbs',
      description: 'Header chain above the editor (Feature → Scenario → step).',
    },
    symbolOutline: {
      title: 'Steps panel structure',
      description: '“Structure” tab with the scenario tree and click-to-jump.',
    },
    stepsPanelView: {
      title: 'Default panel tab',
      description: 'What to show below the editor when opening a scenario.',
    },
    codeLens: {
      title: 'Run Code Lens',
      description: '“▶ Run” buttons above scenarios and steps in the editor.',
    },
    inlayHints: {
      title: 'Inlay hints',
      description: 'Gray hints to the right of a step: click → selector, fill → value.',
    },
    scenarioHints: {
      title: 'Show hints',
      description: 'Warning/info markers in the editor and the “Validate” panel.',
    },
    scenarioHintsAfterRecord: {
      title: 'After recording',
      description: 'Analyze the scenario right after recording stops.',
    },
    scenarioHintsShowWarning: {
      title: 'Warnings',
      description: 'Warning-level hints (fragile selectors, duplicates).',
    },
    scenarioHintsShowInfo: {
      title: 'Information',
      description: 'Info-level hints (improvements without critical risk).',
    },
    scenarioHintsAutoFixOnSave: {
      title: 'Auto-fix on save',
      description: 'Apply autoFixable hints on Ctrl+S.',
    },
    toolbarCompact: {
      title: 'Compact toolbar',
      description: 'Fewer button labels — icons only.',
    },
    stepsPanelVisible: {
      title: 'Show steps panel',
      description: 'Display parsed Gherkin steps below the editor.',
    },
    stepsPanelHeight: {
      title: 'Panel height',
      description: 'Height of the step list area in pixels.',
    },
    checkUpdates: {
      title: 'Check on startup',
      description: 'Look for a new {brand} version when the IDE starts.',
    },
  },
  browser: {
    installed: '{label}: installed',
    notInstalled: '{label}: not installed — download required before recording and runs.',
    checking: 'Checking engine…',
    installing: 'Installing {engine}…',
    installingBtn: 'Installing…',
    reinstall: 'Reinstall',
    installEngine: 'Install engine',
    error: 'Error',
  },
  strategies: {
    clicksTitle: 'Clicks and buttons',
    inputsTitle: 'Input fields',
    moveUp: 'Up',
    moveDown: 'Down',
    resetClicks: 'Reset clicks',
    resetInputs: 'Reset inputs',
    contextual: 'Contextual has-text',
    text: 'has-text by text',
  },
  plugins: {
    playwright: 'Playwright',
    builtIn: 'built-in',
    vanessa: 'Vanessa Automation',
    installed: 'installed',
    unavailable: 'unavailable',
    zipHint: 'ZIP plugins: scenaria plugins install … or the “Manage plugins” dialog.',
  },
  navWait: {
    load: 'load — full page load',
    domcontentloaded: 'domcontentloaded — DOM ready (default)',
    networkidle: 'networkidle — network idle',
    commit: 'commit — first navigation response',
  },
  editor: {
    themeDark: 'Dark',
    themeLight: 'Light',
    themeSystem: 'System',
    wordWrapOn: 'On',
    wordWrapOff: 'Off',
    lineNumbersOn: 'Normal',
    lineNumbersRelative: 'Relative',
    lineNumbersOff: 'Hide',
    whitespaceNone: 'Do not show',
    whitespaceBoundary: 'At word boundaries',
    whitespaceSelection: 'In selection',
    whitespaceTrailing: 'At end of line',
    whitespaceAll: 'Always',
    autoQuotesLanguage: 'By language',
    autoQuotesAlways: 'Always',
    autoQuotesBeforeWhitespace: 'Before whitespace',
    autoQuotesNever: 'Never',
    stepsPanelOutline: 'Structure',
    stepsPanelSteps: 'Step table',
  },
  units: {
    ms: 'ms',
    pcs: 'pcs',
    px: 'px',
  },
  warnings: {
    manyWorkers: 'Many workers and high slow-mo — the run may take a very long time.',
    highWorkers: 'A large number of workers loads the system.',
    highSlowMo: 'Very high slow-mo — steps run with a long pause.',
  },
  uiLocale: 'Interface language',
  uiLocaleRu: 'Русский',
  uiLocaleEn: 'English',
  slowMo: {
    fast: 'Fast',
    normal: 'Normal',
    slow: 'Slow',
    tutorial: 'Tutorial',
  },
  resetDefaults: 'Reset to defaults',
  managePlugins: 'Manage plugins…',
  vanessaSettings: 'Vanessa settings…',
  apply: 'Apply',
  applyBusy: '…',
  applyDone: '✓',
  ok: 'OK',
  cancel: 'Cancel',
  hotkeysEditorNote: 'In the editor: Ctrl+F — find, Ctrl+H — find and replace.',
} as const

export const splash = {
  envSetup: 'Configuring environment…',
  connecting: 'Connecting to application…',
  loadingSettings: 'Loading settings…',
  initializing: 'Initializing…',
  ready: 'Ready',
  starting: 'Starting…',
} as const
