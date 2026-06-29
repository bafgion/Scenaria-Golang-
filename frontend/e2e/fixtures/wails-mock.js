// Injected before page load in Playwright e2e — stubs Wails bindings for static preview.
(() => {
  const noop = () => {}
  const asyncEmpty = async () => ''
  const asyncOk = async () => ({ output: 'ok', error: '' })

  const E2E_PROJECT = 'C:/e2e/project'
  const e2eMode = () => new URLSearchParams(location.search).get('e2e')

  const settings = {
    browser: 'chromium',
    headless: false,
    parallelWorkers: 1,
    maxLoopIterations: 100,
    stepsPanelVisible: true,
    stepsPanelHeight: 160,
    checkUpdatesOnStartup: false,
    recentProjects: [E2E_PROJECT],
    recentFeatures: [],
  }

  const sampleSteps = [
    {
      label: 'нажимаю',
      action: 'click',
      category: 'Формы и ввод',
      description: 'Клик по элементу',
      template: 'нажимаю "button.submit"',
      example: 'нажимаю "button.submit"',
      parameters: ['selector — CSS/XPath селектор элемента'],
      help: 'Клик по элементу',
    },
    {
      label: 'открыт',
      action: 'goto',
      category: 'Навигация',
      description: 'Переход на страницу',
      template: 'открыт "https://site.com"',
      example: 'открыт "https://site.com"',
      parameters: ['url — адрес страницы в кавычках'],
      help: 'Переход на страницу',
    },
  ]

  const handlers = new Map()

  const emitE2E = (event, payload) => {
    const set = handlers.get(event)
    if (!set) return
    for (const cb of set) {
      try {
        cb(payload)
      } catch {
        /* ignore */
      }
    }
  }

  const postRecordHints = [
    {
      id: 'menu_hover',
      title: 'Hover перед кликом по меню',
      severity: 'warning',
      stepIndex: 0,
      line: 3,
      autoFixable: true,
    },
  ]

  const app = {
    Version: async () => 'e2e-test',
    LoadSettings: async () => settings,
    SaveSettings: async (s) => ({ ...settings, ...s }),
    LoadRecents: async () => ({ projects: [E2E_PROJECT], features: [] }),
    RememberRecentProject: async () => {},
    RememberRecentFeature: async () => {},
    BrowserInstallStatus: async (engine) => {
      const mode = new URLSearchParams(location.search).get('e2e')
      const installed = mode !== 'missing-browser'
      return {
        engine: engine || 'chromium',
        label: engine === 'firefox' ? 'Firefox' : engine === 'webkit' ? 'WebKit' : 'Chromium',
        installed,
        detail: installed ? 'C:/mock/ms-playwright/chromium-1200/chrome-win64/chrome.exe' : '',
      }
    },
    InstallBrowserEngine: async () => ({ output: 'Готово: mock', error: '' }),
    SearchSteps: async () => sampleSteps,
    CompletionsForLine: async (line, column) => {
      const items = sampleSteps.map((s) => ({
        label: s.label,
        insert: s.template,
        description: s.description,
      }))
      const match = line.match(/^\s*((?:Допустим|Дано|Когда|Тогда|И|Но)\s+)?(.*)$/i)
      const body = match?.[2] || ''
      const bodyOffset = line.length - body.length
      const typed = line.slice(bodyOffset, column).trimStart().toLowerCase()
      const filtered = typed
        ? items.filter(
            (s) => s.label.toLowerCase().startsWith(typed) || s.insert.toLowerCase().startsWith(typed),
          )
        : items
      const start = typed ? column - line.slice(bodyOffset, column).trimStart().length : bodyOffset
      return { start, end: column, items: filtered }
    },
    ValidateFeature: async () => [],
    OpenProject: async (path) => {
      const root = path || E2E_PROJECT
      const smoke = `${root.replace(/\\/g, '/')}/smoke.feature`
      return {
        path: root,
        features: [smoke],
        tags: ['@smoke'],
        featureTags: { [smoke]: ['@smoke'] },
        name: 'e2e',
      }
    },
    PickProjectFolder: asyncEmpty,
    PickOpenFile: asyncEmpty,
    PickOpenFiles: async () => {
      if (e2eMode() === 'import-pick') return ['C:/external/sample.feature']
      return []
    },
    PickSaveFile: asyncEmpty,
    ReadFeature: async () => 'Функция: smoke\n  Сценарий: тест\n    открыт "https://example.com"',
    SaveFeature: asyncOk,
    WriteTempFeature: async () => `${E2E_PROJECT}/.scenaria/temp.feature`,
    Run: async (opts) => ({ output: opts?.dryRun ? 'Dry-run ok' : 'ok', error: '' }),
    Validate: asyncOk,
    ListTestClients: async () => [],
    ListPlugins: async () => [],
    ListRunResults: async () => [],
    ProjectArtifacts: async () => ({ allureDir: '', reportHtml: '' }),
    ParseEditorSteps: async () => {
      if (e2eMode() === 'post-record') {
        return [
          { text: 'нажимаю "nav >> text=Каталог"' },
          { text: 'нажимаю "nav >> text=Товары"' },
          { text: 'открыт "https://example.com"' },
        ]
      }
      return [{ text: 'открыт "https://example.com"' }]
    },
    InitProject: async () => 'init ok',
    CheckUpdate: async () => ({ available: false, version: '', url: '' }),
    ArtifactExists: async () => false,
    BundledExamplesPath: async () => '',
    ListScenarioTitles: async () => ['тест'],
    AnalyzeScenarioHints: async () => (e2eMode() === 'post-record' ? postRecordHints : []),
    ApplyScenarioHintFix: async (req) => {
      const text = req?.text || ''
      if (e2eMode() === 'post-record' && req?.hintId === 'menu_hover') {
        return { text: `${text}\n    Когда навожу на "menu"`, count: 1 }
      }
      return { text, count: 0 }
    },
    RefactorReplaceInText: async (_text, find, replace) => _text.split(find).join(replace),
    ReplaceInProject: async () => ({ changed: 0, output: '' }),
    ImportFeatures: async (_dest, paths) =>
      (paths || []).map((p) => `${E2E_PROJECT}/${String(p).split(/[/\\]/).pop()}`),
    PreviewExport: async () => ({
      stepCount: 2,
      scenarioTitle: 'e2e',
      issues: [],
      hints: [],
    }),
    DeleteFeature: asyncOk,
    DuplicateFeature: asyncOk,
    MoveFeature: asyncOk,
    RenameFeature: asyncOk,
    Export: async (opts) => ({ output: `exported to ${opts?.output || ''}`, error: '' }),
    RunPlugin: asyncOk,
    StartRecord: async () => {
      if (e2eMode() === 'post-record') {
        emitE2E('record-started', null)
        queueMicrotask(() =>
          emitE2E('record-finished', { output: 'Записано шагов: 3', error: '' }),
        )
      }
    },
    RecordBaseline: asyncOk,
    PauseRecording: noop,
    ResumeRecording: noop,
    CancelRecording: noop,
    FocusBrowser: noop,
    UndoRecordedStep: asyncOk,
    UpdateRecordingOptions: noop,
    PickSelector: asyncEmpty,
    PickerStepChoices: async () => [],
    SubmitOTPCode: noop,
    CancelOTP: noop,
    OpenFolder: noop,
    ServeAllure: asyncOk,
    OpenHTMLReport: noop,
    RefactorUpdateStartURLs: async (_paths, _url) => ({ changed: 0, output: '' }),
    RefactorNormalizeIndents: async (text) => text,
    RefactorCollapseBlankLines: async (text) => text,
    ListVanessaRunDirs: async () => [],
    StartVanessaRun: noop,
    PollVanessaRun: async () => ({ done: true, passed: 0, failed: 0, total: 0 }),
    HighlightFeature: noop,
    ImportJSON: asyncOk,
    InstallPlugin: asyncOk,
    IsRecordingPaused: async () => false,
    HTTPAuthForHost: async () => null,
    DeleteTestClient: asyncOk,
  }

  const registerHandler = (event, cb) => {
    if (!handlers.has(event)) handlers.set(event, new Set())
    handlers.get(event).add(cb)
    return () => handlers.get(event)?.delete(cb)
  }

  window.runtime = {
    LogPrint: noop,
    EventsOn: registerHandler,
    EventsOnMultiple: registerHandler,
    EventsOff: noop,
    EventsEmit: noop,
    OnFileDrop: noop,
    OnFileDropOff: noop,
  }

  window.__e2eEmit = emitE2E
  window.go = { wailsapp: { App: app } }
})()
