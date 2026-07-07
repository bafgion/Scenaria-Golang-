// Injected before page load in Playwright e2e — stubs Wails bindings for static preview.
(() => {
  const noop = () => {}
  const asyncEmpty = async () => ''
  const asyncOk = async () => ({ output: 'ok', error: '' })

  const E2E_PROJECT = 'C:/e2e/project'
  const e2eMode = () => new URLSearchParams(location.search).get('e2e')

  const SETTINGS_STORAGE_KEY = 'scenaria.e2e.settings'

  const defaultSettings = {
    browser: 'chromium',
    headless: false,
    parallelWorkers: 1,
    maxLoopIterations: 100,
    stepsPanelVisible: true,
    stepsPanelHeight: 160,
    checkUpdatesOnStartup: false,
    recentProjects: [E2E_PROJECT],
    recentFeatures: [],
    sessionProject: '',
    openTabs: [],
    untitledTabs: [],
    activeTab: '',
    runDialogConfirmed: false,
    pickerDuringRecording: false,
    onboardingCompleted: false,
    onboardingDismissed: false,
    onboardingVersion: 0,
    uiLocale: 'ru',
    editor: {
      fontSize: 13,
      fontFamily: '"Cascadia Code", Consolas, monospace',
      wordWrap: 'on',
      minimap: true,
      lineNumbers: 'on',
      tabSize: 4,
      insertSpaces: false,
      renderWhitespace: 'selection',
      folding: false,
      stickyScroll: false,
      autoClosingQuotes: 'languageDefined',
      formatOnSave: false,
      theme: 'scenaria-dark',
      stepHoverEnabled: true,
      validateOnType: true,
    },
  }

  const readSettings = () => {
    try {
      const raw = sessionStorage.getItem(SETTINGS_STORAGE_KEY)
      if (raw) {
        return { ...defaultSettings, ...JSON.parse(raw) }
      }
    } catch {
      /* ignore */
    }
    return { ...defaultSettings }
  }

  const params = new URLSearchParams(location.search)
  if (params.has('__fresh')) {
    try {
      sessionStorage.removeItem(SETTINGS_STORAGE_KEY)
    } catch {
      /* ignore */
    }
    params.delete('__fresh')
    const qs = params.toString()
    const next = `${location.pathname}${qs ? `?${qs}` : ''}${location.hash}`
    history.replaceState(null, '', next)
  }

  let settings = readSettings()
  const localeParam = params.get('locale')
  if (localeParam === 'en' || localeParam === 'ru') {
    settings = { ...settings, uiLocale: localeParam }
  }
  const mode = e2eMode()
  if (mode === 'run-progress' || mode === 'run-stream' || mode === 'run-cancel' || mode === 'demo-video') {
    settings = { ...settings, runDialogConfirmed: true }
  }
  if (mode === 'update-available') {
    settings = { ...settings, checkUpdatesOnStartup: true }
  }

  let runCancelled = false
  let lastRunRequest = null

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

  const liveRecord = {
    browserOpen: false,
    recording: false,
    paused: false,
    captureEver: false,
    steps: [],
  }

  const postRecordSteps = [
    { index: 0, line: 'нажимаю "#login"' },
    { index: 1, line: 'ввожу "user" в "#email"' },
    { index: 2, line: 'нажимаю "#submit"' },
  ]

  const resumeRecordSteps = [
    { index: 0, line: 'открыт "https://example.com"' },
    { index: 1, line: 'нажимаю "#one"' },
    { index: 2, line: 'открыт "https://iana.org"' },
    { index: 3, line: 'открыт "https://iana.org/domains"' },
    { index: 4, line: 'открыт "https://iana.org/domains/reserved"' },
    { index: 5, line: 'открыт "https://iana.org/domains/example"' },
  ]

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

  const flakyScenarioPath = `${E2E_PROJECT}/smoke.feature::тест`
  const flakyRunResults = [
    {
      path: flakyScenarioPath,
      success: false,
      message: 'элемент не найден',
      runner: 'playwright',
      at: '2026-06-28T12:00:00Z',
      failed_step: 1,
    },
    {
      path: flakyScenarioPath,
      success: true,
      message: '',
      runner: 'playwright',
      at: '2026-06-27T12:00:00Z',
    },
    {
      path: `${E2E_PROJECT}/smoke.feature::стабильный`,
      success: true,
      message: '',
      runner: 'playwright',
      at: '2026-06-26T12:00:00Z',
    },
  ]

  const demoRunResults = [
    {
      path: `${E2E_PROJECT}/smoke.feature::тест`,
      success: true,
      message: '',
      runner: 'playwright',
      at: '2026-06-29T10:00:00Z',
    },
    {
      path: `${E2E_PROJECT}/login.feature::Неуспешный вход`,
      success: false,
      message: 'элемент не найден',
      runner: 'playwright',
      at: '2026-06-29T10:01:00Z',
      failed_step: 1,
    },
  ]

  const flakyMetricsPayload = {
    scenarios: [
      {
        path: flakyScenarioPath,
        failures: 1,
        passes: 1,
        total: 2,
        flaky: true,
        last_failed_at: '2026-06-28T12:00:00Z',
      },
    ],
    steps: [
      {
        path: `${E2E_PROJECT}/checkout.feature::оплата`,
        step: 2,
        failures: 2,
        last_failed_at: '2026-06-25T09:00:00Z',
      },
    ],
  }

  let focusBrowserCalls = 0

  const app = {
    Version: async () => 'e2e-test',
    LoadSettings: async () => ({ ...settings }),
    SaveSettings: async (s) => {
      settings = { ...settings, ...s }
      try {
        sessionStorage.setItem(SETTINGS_STORAGE_KEY, JSON.stringify(settings))
      } catch {
        /* ignore */
      }
      return { ...settings }
    },
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
    DescribeEditorLine: async (line) => {
      const text = String(line || '')
      if (text.includes('нажимаю')) return sampleSteps[0]
      return { label: '', action: '', category: '', description: '', template: '', example: '', parameters: [], help: '' }
    },
    CompletionsForLine: async (line, column, featureText = '') => {
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
      const norm = root.replace(/\\/g, '/')
      const isExamples = norm.includes('examples')
      const features = isExamples
        ? [`${norm}/smoke.feature`, `${norm}/login.feature`, `${norm}/api-smoke.feature`]
        : [`${norm}/smoke.feature`]
      const featureTags = Object.fromEntries(features.map((f) => [f, ['@smoke']]))
      return {
        path: root,
        features,
        tags: ['@smoke'],
        featureTags,
        name: isExamples ? 'examples' : 'e2e',
      }
    },
    RefreshProject: async () => {
      const root = E2E_PROJECT
      const norm = root.replace(/\\/g, '/')
      return {
        path: root,
        features: [`${norm}/smoke.feature`],
        tags: ['@smoke'],
        featureTags: { [`${norm}/smoke.feature`]: ['@smoke'] },
        name: 'e2e',
      }
    },
    PickProjectFolder: async () => (e2eMode() === 'new-project' ? 'C:/e2e/new-project' : ''),
    PickOpenFile: asyncEmpty,
    PickOpenFiles: async () => {
      if (e2eMode() === 'import-pick') return ['C:/external/sample.feature']
      return []
    },
    PickSaveFile: async () => {
      if (e2eMode() === 'save-as') return `${E2E_PROJECT}/saved-as.feature`
      return ''
    },
    ReadFeature: async () => {
      if (e2eMode() === 'large-file') {
        const lines = ['Функционал: big', 'Сценарий: one']
        while (lines.length < 2001) lines.push('\tоткрыт "https://example.com"')
        return lines.join('\n')
      }
      return 'Функция: smoke\n  Сценарий: тест\n    открыт "https://example.com"'
    },
    SaveFeature: asyncOk,
    WriteTempFeature: async () => `${E2E_PROJECT}/.scenaria/temp.feature`,
    Run: async (opts) => {
      lastRunRequest = opts
      if (e2eMode() === 'demo-video') {
        const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms))
        emitE2E('run-progress', {
          phase: 'scenario_start',
          index: 1,
          total: 2,
          featurePath: `${E2E_PROJECT}/smoke.feature`,
          scenario: 'тест',
        })
        await delay(2000)
        emitE2E('run-progress', { phase: 'scenario_done', index: 1, total: 2, success: true })
        emitE2E('run-progress', {
          phase: 'scenario_start',
          index: 2,
          total: 2,
          featurePath: `${E2E_PROJECT}/login.feature`,
          scenario: 'Неуспешный вход',
        })
        await delay(2000)
        emitE2E('run-progress', { phase: 'scenario_done', index: 2, total: 2, success: false })
        const at = new Date().toISOString()
        return {
          output: 'ok',
          error: '',
          entries: [
            {
              path: `${E2E_PROJECT}/smoke.feature::тест`,
              success: true,
              message: '',
              runner: 'playwright',
              at,
            },
            {
              path: `${E2E_PROJECT}/login.feature::Неуспешный вход`,
              success: false,
              message: 'элемент не найден',
              runner: 'playwright',
              at,
              failed_step: 1,
            },
          ],
        }
      }
      if (e2eMode() === 'run-progress') {
        emitE2E('run-progress', {
          phase: 'scenario_start',
          index: 1,
          total: 3,
          featurePath: `${E2E_PROJECT}/smoke.feature`,
          scenario: 'тест',
        })
        await new Promise((resolve) => setTimeout(resolve, 200))
        emitE2E('run-progress', { phase: 'scenario_done', index: 1, total: 3, success: true })
        emitE2E('run-progress', {
          phase: 'scenario_start',
          index: 2,
          total: 3,
          featurePath: `${E2E_PROJECT}/smoke.feature`,
          scenario: 'второй',
        })
        await new Promise((resolve) => setTimeout(resolve, 200))
        emitE2E('run-progress', { phase: 'scenario_done', index: 2, total: 3, success: true })
        emitE2E('run-progress', {
          phase: 'scenario_start',
          index: 3,
          total: 3,
          featurePath: `${E2E_PROJECT}/smoke.feature`,
          scenario: 'третий',
        })
        await new Promise((resolve) => setTimeout(resolve, 200))
        emitE2E('run-progress', { phase: 'scenario_done', index: 3, total: 3, success: true })
        return { output: 'ok', error: '' }
      }
      if (e2eMode() === 'run-stream') {
        emitE2E('run-log-line', { line: 'E2E_STREAM_LINE\n' })
        await new Promise((resolve) => setTimeout(resolve, 80))
        return { output: 'ok', error: '' }
      }
      if (e2eMode() === 'run-cancel') {
        runCancelled = false
        for (let i = 0; i < 40; i++) {
          if (runCancelled) return { output: '', error: 'context canceled' }
          emitE2E('run-log-line', { line: `progress ${i}\n` })
          await new Promise((resolve) => setTimeout(resolve, 120))
        }
        return { output: 'ok', error: '' }
      }
      return { output: opts?.dryRun ? 'Dry-run ok' : 'ok', error: '' }
    },
    CancelRun: () => {
      runCancelled = true
    },
    Validate: async () => ({ output: 'Проверка завершена.', error: '' }),
    ListTestClients: async () => [],
    ListPlugins: async () =>
      e2eMode() === 'with-plugins'
        ? [{ name: 'demo', description: 'Demo plugin', runnable: true, id: 'demo', source: 'local' }]
        : [],
    ListRunResults: async () => {
      const mode = e2eMode()
      if (mode === 'flaky-run') return flakyRunResults
      if (mode === 'demo-video') {
        const at = new Date().toISOString()
        return demoRunResults.map((r) => ({ ...r, at }))
      }
      return []
    },
    FlakyMetrics: async () => (e2eMode() === 'flaky-run' ? flakyMetricsPayload : { scenarios: [], steps: [] }),
    ProjectArtifacts: async () => {
      if (e2eMode() === 'flaky-run' || e2eMode() === 'trace-artifacts') {
        return {
          tracesDir: `${E2E_PROJECT}/.scenaria/traces`,
          allureDir: '',
          htmlReport: '',
          junitReport: '',
          summaryJson: '',
          videosDir: '',
        }
      }
      if (e2eMode() === 'allure-missing') {
        return {
          tracesDir: '',
          allureDir: `${E2E_PROJECT}/.scenaria/allure-results`,
          htmlReport: '',
          junitReport: '',
          summaryJson: '',
          videosDir: '',
        }
      }
      return { allureDir: '', reportHtml: '' }
    },
    OpenTrace: async () => ({ output: 'trace viewer started', error: '' }),
    ScenariaArtifactPath: async (sub) => `${E2E_PROJECT}/.scenaria/${sub}`,
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
    InitProjectAt: async () => 'init ok',
    CheckUpdate: async () => ({ output: 'Установлена актуальная версия', error: '' }),
    CheckUpdateInfo: async () => {
      if (e2eMode() === 'update-available') {
        return {
          currentVersion: '0.0.0-e2e',
          latestVersion: '9.9.9',
          updateAvailable: true,
          htmlUrl: 'https://example.com/release',
          downloadUrl: '',
          downloadName: '',
          message: 'Доступно обновление',
        }
      }
      return {
        currentVersion: '0.0.0-e2e',
        latestVersion: '0.0.0-e2e',
        updateAvailable: false,
        htmlUrl: '',
        downloadUrl: '',
        downloadName: '',
        message: 'Установлена актуальная версия',
      }
    },
    DownloadUpdate: async () => '',
    OpenExternalURL: asyncOk,
    ValidateBrowser: async () => [],
    ArtifactExists: async () => false,
    BundledExamplesPath: async () => `${E2E_PROJECT}/examples`,
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
    OpenBrowser: async () => {
      liveRecord.browserOpen = true
      liveRecord.recording = false
      liveRecord.paused = false
      emitE2E('browser-opened', '')
    },
    StartRecord: async () => {
      const mode = e2eMode()
      if (
        mode === 'demo-video' ||
        mode === 'post-record' ||
        mode === 'post-record-diff' ||
        mode === 'record-resume' ||
        mode === 'record-idle'
      ) {
        liveRecord.browserOpen = true
        liveRecord.recording = true
        liveRecord.captureEver = true
        liveRecord.paused = false
        liveRecord.steps = mode === 'record-resume' ? [...resumeRecordSteps] : [...postRecordSteps]
        emitE2E('browser-opened', '')
        const recordOutput =
          mode === 'demo-video' ? `${E2E_PROJECT}/examples/smoke.feature` : `${E2E_PROJECT}/smoke.feature`
        emitE2E('record-started', { resume: false, output: recordOutput })
        queueMicrotask(() => {
          for (const step of liveRecord.steps) {
            emitE2E('record-step', step)
          }
        })
        if (mode === 'record-idle') {
          setTimeout(() => emitE2E('record-stopped', { reason: 'idle', idleSeconds: 30 }), 400)
        }
        return
      }
      liveRecord.browserOpen = true
      liveRecord.recording = true
      liveRecord.captureEver = true
      liveRecord.paused = false
      liveRecord.steps = []
      emitE2E('browser-opened', '')
      emitE2E('record-started', { resume: false, output: `${E2E_PROJECT}/examples/smoke.feature` })
    },
    BeginRecordingCapture: async () => {
      if (!liveRecord.browserOpen) {
        throw new Error('браузер не открыт')
      }
      if (liveRecord.recording) {
        emitE2E('record-started', { append: true, sync: true })
        return
      }
      const mode = e2eMode()
      liveRecord.captureEver = true
      liveRecord.recording = true
      liveRecord.paused = false
      if (mode === 'record-resume') {
        liveRecord.steps = [...resumeRecordSteps]
        emitE2E('record-started', { append: true, resume: true })
        queueMicrotask(() => {
          for (const step of liveRecord.steps) {
            emitE2E('record-step', step)
          }
        })
        return
      }
      liveRecord.steps = []
      emitE2E('record-started', { append: true })
    },
    PollBrowserSession: async () => ({
      browserOpen: liveRecord.browserOpen,
      recording: liveRecord.recording,
      paused: liveRecord.paused,
      stepCount: liveRecord.recording ? liveRecord.steps.length : 0,
    }),
    RecordBaseline: asyncOk,
    PauseRecording: async () => {
      if (liveRecord.recording) {
        liveRecord.paused = true
      }
    },
    ResumeRecording: async () => {
      if (liveRecord.recording) {
        liveRecord.paused = false
      }
    },
    CancelRecording: async () => {
      liveRecord.browserOpen = false
      liveRecord.recording = false
      liveRecord.paused = false
      liveRecord.captureEver = false
      liveRecord.steps = []
    },
    CloseBrowser: async () => {
      if (e2eMode() === 'post-record-diff' && liveRecord.captureEver) {
        liveRecord.browserOpen = false
        liveRecord.recording = false
        liveRecord.paused = false
        liveRecord.captureEver = false
        liveRecord.steps = []
        emitE2E('record-finished', { output: 'Запись сохранена', error: '' })
        return
      }
      liveRecord.browserOpen = false
      liveRecord.recording = false
      liveRecord.paused = false
      liveRecord.captureEver = false
      liveRecord.steps = []
    },
    StopRecordingCapture: async () => {
      if (e2eMode() === 'post-record-diff') {
        liveRecord.browserOpen = false
        liveRecord.recording = false
        liveRecord.paused = false
        liveRecord.captureEver = false
        liveRecord.steps = []
        emitE2E('record-finished', { output: 'Запись сохранена', error: '' })
        return
      }
      if (!liveRecord.browserOpen) {
        throw new Error('браузер не открыт')
      }
      liveRecord.recording = false
      liveRecord.paused = false
      liveRecord.captureEver = false
      liveRecord.steps = []
      emitE2E('record-stopped', null)
    },
    FocusBrowser: async () => {
      focusBrowserCalls += 1
    },
    UndoRecordedStep: asyncOk,
    UpdateRecordingOptions: noop,
    PickSelector: asyncEmpty,
    PickerStepChoices: async () => [],
    SubmitOTPCode: async () => true,
    CancelOTP: noop,
    OpenFolder: noop,
    ServeAllure: asyncOk,
    AllureStatus: async () => {
      if (e2eMode() === 'allure-missing') {
        return { installed: false, running: false, resultsDir: `${E2E_PROJECT}/.scenaria/allure-results` }
      }
      return { installed: true, running: false, resultsDir: '' }
    },
    FailedStepLine: async () => 4,
    OpenHTMLReport: noop,
    RefactorUpdateStartURLs: async (_paths, _url) => ({ changed: 0, output: '' }),
    RefactorNormalizeIndents: async (text) => text,
    RefactorCollapseBlankLines: async (text) => text,
    FormatFeature: async (text) => text.replace(/\n\s*\n(\s*(?:Когда|Тогда|Допустим|И|Но))/g, '\n$1'),
    ListVanessaRunDirs: async () => [],
    StartVanessaRun: noop,
    PollVanessaRun: async () => ({ done: true, passed: 0, failed: 0, total: 0 }),
    HighlightFeature: noop,
    ImportJSON: asyncOk,
    InstallPlugin: asyncOk,
    IsRecordingPaused: async () => liveRecord.paused,
    HTTPAuthForHost: async () => null,
    ListHTTPAuthHosts: async () => [],
    SaveHTTPAuth: asyncOk,
    RemoveHTTPAuth: asyncOk,
    DeleteTestClient: asyncOk,
    CaptureBrowserSession: async (name) => `captured ${name}`,
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
  window.__e2eLastRunRequest = () => lastRunRequest
  window.__e2eFocusBrowserCalls = () => focusBrowserCalls
  window.go = { wailsapp: { App: app } }
})()
