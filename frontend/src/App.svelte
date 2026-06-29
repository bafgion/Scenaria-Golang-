<script lang="ts">
  import { onDestroy, onMount } from 'svelte'
  import MonacoEditor from './lib/MonacoEditor.svelte'
  import WelcomePanel from './lib/WelcomePanel.svelte'
  import CatalogEmptyState from './lib/CatalogEmptyState.svelte'
  import FeatureCatalogTree from './lib/FeatureCatalogTree.svelte'
  import EditorTabBar from './lib/EditorTabBar.svelte'
  import { buildCatalogViewState } from './lib/catalogTree'
  import SettingsDialog from './lib/SettingsDialog.svelte'
  import CommandPalette, { type PaletteCommand } from './lib/CommandPalette.svelte'
  import FindReplaceDialog from './lib/FindReplaceDialog.svelte'
  import ProjectReplaceDialog from './lib/ProjectReplaceDialog.svelte'
  import HotkeysDialog from './lib/HotkeysDialog.svelte'
  import PluginsDialog from './lib/PluginsDialog.svelte'
  import ExportDialog from './lib/ExportDialog.svelte'
  import RunDialog from './lib/RunDialog.svelte'
  import RecordDialog from './lib/RecordDialog.svelte'
  import VanessaRunDialog from './lib/VanessaRunDialog.svelte'
  import TestClientDialog from './lib/TestClientDialog.svelte'
  import ImportJSONDialog from './lib/ImportJSONDialog.svelte'
  import AboutDialog from './lib/AboutDialog.svelte'
  import OtpDialog from './lib/OtpDialog.svelte'
  import StepsInsertDialog from './lib/StepsInsertDialog.svelte'
  import VanessaSettingsDialog from './lib/VanessaSettingsDialog.svelte'
  import { defaultRunForm, type RunForm } from './lib/runTypes'
  import PostRecordBanner from './lib/PostRecordBanner.svelte'
  import RunHistoryDialog from './lib/RunHistoryDialog.svelte'
  import StepsHelpDialog from './lib/StepsHelpDialog.svelte'
  import UnsavedCloseDialog from './lib/UnsavedCloseDialog.svelte'
  import HttpAuthDialog from './lib/HttpAuthDialog.svelte'
  import PickerStepDialog from './lib/PickerStepDialog.svelte'
  import { loadLayout, saveLayout, resetLayout as resetUILayout } from './lib/layout'
  import ErrorPanel from './lib/ErrorPanel.svelte'
  import ResultsPanel from './lib/ResultsPanel.svelte'
  import ValidatePanel from './lib/ValidatePanel.svelte'
  import FeaturePreview from './lib/FeaturePreview.svelte'
  import SnippetPalette from './lib/SnippetPalette.svelte'
  import BrowserOverlay from './lib/BrowserOverlay.svelte'
  import SplashScreen from './lib/SplashScreen.svelte'
  import { loadRecents, rememberProject, rememberFeature } from './lib/recents'
  import { icons, toolbarIcons } from './lib/icons'
  import { EventsOn, OnFileDrop, OnFileDropOff } from '../wailsjs/runtime/runtime'
  import {
    Version,
    OpenProject,
    ReadFeature,
    SaveFeature,
    Run,
    Validate,
    ValidateFeature,
    ListTestClients,
    InitProject,
    PickProjectFolder,
    PickSaveFile,
    PickOpenFile,
    RunPlugin,
    StartRecord,
    PauseRecording,
    ResumeRecording,
    CancelRecording,
    FocusBrowser,
    UndoRecordedStep,
    PickSelector,
    PickerStepChoices,
    LoadSettings,
    SaveSettings,
    SubmitOTPCode,
    CancelOTP,
    CheckUpdate,
    ListRunResults,
    BundledExamplesPath,
    ProjectArtifacts,
    ParseEditorSteps,
    ArtifactExists,
    OpenFolder,
    RefactorUpdateStartURLs,
    RefactorNormalizeIndents,
    RefactorCollapseBlankLines,
    RefactorReplaceInText,
    ReplaceInProject,
    AnalyzeScenarioHints,
    ApplyScenarioHintFix,
    DeleteFeature,
    DuplicateFeature,
    MoveFeature,
    ImportFeatures,
    ListPlugins,
  } from '../wailsjs/go/wailsapp/App'
  import { gui } from '../wailsjs/go/models'

  const WELCOME_KEY = '__welcome__'

  type EditorTab = { path: string; content: string; dirty: boolean }
  type EditorStepRow = gui.EditorStepRow

  let version = ''
  let projectPath = ''
  let features: string[] = []
  let tags: string[] = []
  let featureTags: Record<string, string[]> = {}
  let tabs: EditorTab[] = []
  let activeTab = WELCOME_KEY
  let welcomeTabVisible = true
  let editorText = ''
  let logText = ''
  let stepStatus = '0 шагов'
  let stepStatusError = false
  let testClients: string[] = []
  let monaco: MonacoEditor | undefined

  let showSplash = true
  let splashMessage = 'Запуск…'
  let splashProgress = 0
  let splashFading = false
  let showSettings = false
  let showCommandPalette = false
  let showSnippetPalette = false
  let showRecord = false
  let showOtp = false
  let showRun = false
  let showTestClient = false
  let showExport = false
  let showImport = false
  let exportInputPath = ''
  let showSteps = false
  let showStepsHelp = false
  let stepsHelpQuery = ''
  let catalogDropTarget = ''
  let showAbout = false
  let showPlugins = false
  let installedPlugins: gui.PluginEntryDTO[] = []
  let showFindReplace = false
  let showProjectReplace = false
  let showHotkeys = false
  let showRunHistory = false
  let showVanessaRun = false
  let showVanessaSettings = false
  let vanessaDry = false
  let vanessaTag = ''
  let vanessaExcludeTags = ''
  let vanessaScenario = ''
  let vanessaPreferRerun = false
  let vanessaRerunDir = ''
  let vanessaInstallEpf = false
  let vanessaEpfUrl = ''
  let vanessaEpfDest = ''
  let showHttpAuth = false
  let httpAuthHost = ''
  let showPickerStep = false
  let pickerSelector = ''
  let pickerChoices: gui.PickerStepChoice[] = []
  let projectReplaceBusy = false
  let findText = ''
  let replaceText = ''
  let replaceCaseSensitive = false
  let postRecordPath = ''
  let postRecordStepCount = 0
  let postRecordAllHints: gui.ScenarioHintDTO[] = []
  let postRecordHints: gui.ScenarioHintDTO[] = []
  let postRecordDismissed = new Set<string>()
  let contextMenu: { x: number; y: number; path: string } | null = null
  let runDialogTitle = 'Запуск сценария'

  let recording = false
  let recordPaused = false
  let playing = false

  let sidebarVisible = true
  let sidebarWidth = 260
  let previewVisible = false
  let previewWidth = 360
  let bottomPanelOpen = false
  let bottomTab: 'journal' | 'results' | 'validate' | 'error' = 'journal'
  let stepsPanelCollapsed = true
  let stepsPanelHeight = 160
  let bottomPanelHeight = 200
  let sidebarSearch = ''
  let openMenu: string | null = null
  let statusMessage = 'Проект → Открыть проект…'
  let statusTone: 'normal' | 'error' | 'success' | 'busy' = 'normal'

  let batchMode = false
  let batchSelected: string[] = []
  let catalogCollapsed = new Set<string>()
  let showBatchHint = true
  let toolbarCompact = false

  let filterRecording = false
  let navOnlyRecording = false
  let hoverRecord = false
  let stepsPanelVisible = true

  let lastRun: RunForm = defaultRunForm({ headed: true, installPW: true, html: true })
  let runForm: RunForm = { ...lastRun }

  let settingsBrowser = 'chromium'
  let settingsHeadless = false
  let settingsWorkers = 1
  let settingsLoops = 100
  let settingsCheckUpdatesOnStartup = true

  let recordURL = 'https://site.com'
  let startURL = 'https://site.com'
  let recordOutput = 'recorded.feature'
  let recordIdle = 30
  let recordAppendTo = ''
  let recordTestClient = ''
  let lastRecordTarget = ''
  let pendingCloseTab: string | null = null

  let otpEmail = ''

  let testClientSelection = ''

  let recentProjects: string[] = []
  let recentFeatures: string[] = []
  let runResults: gui.RunResultEntry[] = []
  let lastErrorEntry: gui.RunResultEntry | null = null
  let editorSteps: EditorStepRow[] = []
  let editorValidationIssues: gui.ValidationIssue[] = []
  let projectArtifacts: gui.ProjectArtifacts = new gui.ProjectArtifacts()

  const unsubscribers: (() => void)[] = []

  $: isWelcome = activeTab === WELCOME_KEY
  $: activeFeatureTab = tabs.find((t) => t.path === activeTab)
  $: stepCount = editorSteps.length
  $: batchCount = batchSelected.length
  $: showRecordingBar = recording || showRecord || playing
  $: showBrowserOverlay = recording || playing
  $: paletteCommands = buildPaletteCommands()

  let resizingBottom = false
  let resizingSteps = false
  let resizingSidebar = false
  let resizingPreview = false
  let actionBarEl: HTMLDivElement | undefined
  let toolbarIconOnly = true

  $: activeFeaturePath = activeTab !== WELCOME_KEY ? activeTab : ''

  $: runByPath = (() => {
    const map = new Map<string, { success: boolean; message: string; at: string; runner: string }>()
    for (const entry of runResults) {
      const key = entry.path.replace(/\\/g, '/').toLowerCase()
      if (!map.has(key)) {
        map.set(key, {
          success: entry.success,
          message: entry.message || '',
          at: entry.at || '',
          runner: entry.runner || '',
        })
      }
    }
    return map
  })()

  $: tagsByPath = (() => {
    const map = new Map<string, string[]>()
    for (const [path, pathTags] of Object.entries(featureTags)) {
      map.set(path.replace(/\\/g, '/').toLowerCase(), pathTags)
    }
    return map
  })()

  $: catalogViewState = buildCatalogViewState(
    projectPath || null,
    features,
    sidebarSearch,
    runByPath,
    true,
    tagsByPath,
  )

  $: welcomeProjectOpen = !!projectPath
  $: welcomeRecorded = recording || editorSteps.length > 0 || tabs.length > 0
  let welcomePlayedSuccess = false

  const MIN_SPLASH_MS = 1400
  const SPLASH_FADE_MS = 320

  function sleep(ms: number) {
    return new Promise<void>((resolve) => window.setTimeout(resolve, ms))
  }

  function setSplashStage(message: string, progress: number) {
    splashMessage = message
    splashProgress = progress
  }

  async function dismissSplash(startedAt: number) {
    const remaining = MIN_SPLASH_MS - (Date.now() - startedAt)
    if (remaining > 0) await sleep(remaining)
    splashFading = true
    await sleep(SPLASH_FADE_MS)
    showSplash = false
  }

  onMount(async () => {
    const startedAt = Date.now()

    setSplashStage('Настройка окружения…', 8)
    await sleep(40)

    setSplashStage('Оформление интерфейса…', 18)
    await sleep(40)

    const layout = loadLayout()
    sidebarVisible = layout.sidebarVisible
    bottomPanelOpen = layout.bottomPanelOpen
    bottomPanelHeight = layout.bottomPanelHeight
    sidebarWidth = layout.sidebarWidth || 260
    previewVisible = layout.previewVisible
    previewWidth = layout.previewWidth || 360

    setSplashStage('Загрузка модулей…', 32)

    try {
      version = await Version()
    } catch {
      version = 'dev'
    }

    setSplashStage('Создание рабочего окна…', 55)

    const recents = await loadRecents()
    recentProjects = recents.projects
    recentFeatures = recents.features

    try {
      const s = await LoadSettings()
      applySettingsFromDTO(s)
      stepsPanelCollapsed = !s.stepsPanelVisible
      stepsPanelHeight = s.stepsPanelHeight || 160
      if (s.sidebarWidth >= 120) sidebarWidth = s.sidebarWidth
      if (s.recentProjects?.length) recentProjects = s.recentProjects
      if (s.recentFeatures?.length) recentFeatures = s.recentFeatures
    } catch {
      /* dev without wails */
    }

    setSplashStage('Подготовка интерфейса…', 82)
    await sleep(80)

    setSplashStage('Готово', 100)

    try {
      OnFileDrop((_x, _y, paths) => {
        if (projectPath && paths?.length) void importDroppedFeatures(projectPath, paths)
      }, false)
      unsubscribers.push(() => OnFileDropOff())
    } catch {
      /* dev without wails runtime */
    }

    try {
      unsubscribers.push(
        EventsOn('otp-prompt', (email: string) => {
          otpEmail = email || ''
          showOtp = true
        }),
      )
      unsubscribers.push(
        EventsOn('record-started', () => {
          recording = true
          recordPaused = false
          setStatus('● Идёт запись', 'busy')
          appendLog('Запись начата…')
          bottomPanelOpen = true
          bottomTab = 'journal'
        }),
      )
      unsubscribers.push(
        EventsOn('record-finished', async (result: gui.RunResult) => {
          recording = false
          recordPaused = false
          showRecord = false
          const reloadPath = lastRecordTarget
          lastRecordTarget = ''
          if (result.output) appendLog(result.output)
          if (result.error) appendLog(`Ошибка записи: ${result.error}`)
          else if (reloadPath) await showPostRecordBanner(reloadPath)
          syncIdleStatus()
          await refreshProject()
          if (!result.error && reloadPath) {
            await loadFeature(reloadPath)
          }
        }),
      )
    } catch {
      /* dev without wails runtime */
    }

    const onResize = () => syncToolbarDensity()
    window.addEventListener('resize', onResize)
    unsubscribers.push(() => window.removeEventListener('resize', onResize))

    const onKey = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && e.key === 's') {
        e.preventDefault()
        saveFeature()
      }
      if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
        e.preventDefault()
        if (projectPath) executeRun({ ...lastRun, dryRun: false })
      }
      if ((e.ctrlKey || e.metaKey) && e.key === 'b') {
        e.preventDefault()
        beginRecord()
      }
      if ((e.ctrlKey || e.metaKey) && e.key === 'r') {
        e.preventDefault()
        beginRecord()
      }
      if ((e.ctrlKey || e.metaKey) && e.key === 'n') {
        e.preventDefault()
        newScenario()
      }
      if ((e.ctrlKey || e.metaKey) && e.key === 'o') {
        e.preventDefault()
        openFileDialog()
      }
      if ((e.ctrlKey || e.metaKey) && e.key === 'h') {
        e.preventDefault()
        openFindReplace()
      }
      if (e.key === 'F1' && e.shiftKey) {
        e.preventDefault()
        showHotkeys = true
      } else if (e.key === 'F1') {
        e.preventDefault()
        openStepsHelp()
      }
      if ((e.ctrlKey || e.metaKey) && e.key === ',') {
        e.preventDefault()
        openSettings()
      }
      if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.key.toLowerCase() === 'p') {
        e.preventDefault()
        showCommandPalette = true
      }
      if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.code === 'Space') {
        e.preventDefault()
        openSnippetPalette()
      }
      if ((e.ctrlKey || e.metaKey) && e.key === '`') {
        e.preventDefault()
        bottomPanelOpen = !bottomPanelOpen
      }
      if (e.key === 'Escape') {
        openMenu = null
        if (showCommandPalette) showCommandPalette = false
        if (showSnippetPalette) showSnippetPalette = false
      }
    }
    const onDocClick = () => {
      openMenu = null
    }
    window.addEventListener('keydown', onKey)
    window.addEventListener('click', onDocClick)
    unsubscribers.push(() => window.removeEventListener('keydown', onKey))
    unsubscribers.push(() => window.removeEventListener('click', onDocClick))

    await dismissSplash(startedAt)
    applyDevUiMock()
    syncIdleStatus()
    void checkUpdatesOnStartup()
  })

  onDestroy(() => {
    for (const off of unsubscribers) off()
  })

  function basename(path: string): string {
    const parts = path.replace(/\\/g, '/').split('/')
    return parts[parts.length - 1] || path
  }

  function applySettingsFromDTO(s: gui.AppSettingsDTO) {
    settingsBrowser = s.browser || 'chromium'
    settingsHeadless = s.headless
    settingsWorkers = s.parallelWorkers || 1
    settingsLoops = s.maxLoopIterations || 100
    filterRecording = s.filterRecording
    navOnlyRecording = s.navOnlyRecording
    hoverRecord = s.hoverRecord
    toolbarCompact = s.toolbarCompact
    stepsPanelVisible = s.stepsPanelVisible !== false
    stepsPanelHeight = s.stepsPanelHeight || 160
    settingsCheckUpdatesOnStartup = s.checkUpdatesOnStartup !== false
    lastRun = {
      ...lastRun,
      workers: s.parallelWorkers || 1,
      browser: s.browser || 'chromium',
    }
  }

  function buildPaletteCommands(): PaletteCommand[] {
    const compactLabel = toolbarCompact ? 'Расширенная панель' : 'Компактная панель'
    return [
      { id: 'palette', label: 'Палитра команд', group: 'Вид', shortcut: 'Ctrl+Shift+P', run: () => (showCommandPalette = true) },
      { id: 'welcome', label: 'Старт', group: 'Вид', run: () => selectTab(WELCOME_KEY) },
      { id: 'open', label: 'Открыть проект…', group: 'Проект', run: openProjectDialog },
      { id: 'settings', label: 'Настройки…', group: 'Проект', shortcut: 'Ctrl+,', run: openSettings },
      { id: 'init', label: 'Init проекта', group: 'Проект', run: initProject },
      { id: 'examples', label: 'Открыть примеры сценариев', group: 'Проект', run: openExamples },
      { id: 'new', label: 'Новый сценарий', group: 'Сценарий', shortcut: 'Ctrl+N', run: newScenario },
      { id: 'open-file', label: 'Открыть файл…', group: 'Сценарий', shortcut: 'Ctrl+O', run: openFileDialog },
      { id: 'save', label: 'Сохранить', group: 'Сценарий', shortcut: 'Ctrl+S', run: saveFeature },
      { id: 'export', label: 'Экспорт…', group: 'Сценарий', run: openExportDialog },
      { id: 'import', label: 'Импорт JSON…', group: 'Сценарий', run: openImportDialog },
      { id: 'steps', label: 'Вставить шаг…', group: 'Сценарий', run: openStepsDialog },
      { id: 'snippets', label: 'Палитра сниппетов', group: 'Сценарий', shortcut: 'Ctrl+Shift+Space', run: openSnippetPalette },
      { id: 'find-replace', label: 'Найти и заменить…', group: 'Сценарий', shortcut: 'Ctrl+H', run: openFindReplace },
      { id: 'project-replace', label: 'Замена по проекту…', group: 'Сценарий', run: () => (showProjectReplace = true) },
      { id: 'duplicate', label: 'Дублировать сценарий', group: 'Сценарий', run: () => activeTab && !isWelcome && duplicateFeature(activeTab) },
      { id: 'delete-feature', label: 'Удалить сценарий', group: 'Сценарий', run: () => activeTab && !isWelcome && deleteFeature(activeTab) },
      { id: 'refactor-indents', label: 'Нормализовать отступы', group: 'Рефакторинг', run: refactorNormalizeIndents },
      { id: 'refactor-blanks', label: 'Убрать пустые строки между шагами', group: 'Рефакторинг', run: refactorCollapseBlank },
      { id: 'steps-help', label: 'Справка по шагам…', group: 'Справка', shortcut: 'F1', run: openStepsHelp },
      { id: 'browser', label: 'Браузер', group: 'Запись и тест', shortcut: 'Ctrl+B', run: beginRecord },
      { id: 'record', label: 'Запись', group: 'Запись и тест', shortcut: 'Ctrl+R', run: beginRecord },
      { id: 'stop', label: 'Стоп', group: 'Запись и тест', run: stopRecord },
      { id: 'pause', label: 'Пауза', group: 'Запись и тест', run: toggleRecordPause },
      { id: 'run', label: 'Запустить', group: 'Запись и тест', shortcut: 'Ctrl+Enter', run: () => executeRun({ ...lastRun, dryRun: false }) },
      { id: 'run-dialog', label: 'Запустить…', group: 'Запись и тест', run: () => openRunDialog('Запуск сценария', {}) },
      { id: 'run-tag', label: 'Запустить сценарии с тегом…', group: 'Запись и тест', run: () => openRunDialog('Запуск с тегом', {}) },
      { id: 'playwright', label: 'Playwright…', group: 'Запись и тест', run: () => openRunDialog('Playwright', { dryRun: false, headed: true, engine: 'playwright', installPW: true }) },
      { id: 'dry', label: 'Dry-run', group: 'Запись и тест', run: () => executeRun({ ...lastRun, dryRun: true }) },
      { id: 'batch', label: 'Пакетный выбор', group: 'Запись и тест', run: () => (batchMode = !batchMode) },
      { id: 'batch-run', label: 'Запустить выбранные', group: 'Запись и тест', run: runBatchSelected },
      { id: 'rerun-failed', label: 'Перезапустить упавшие', group: 'Запись и тест', run: rerunFailed },
      { id: 'run-history', label: 'История запусков…', group: 'Запись и тест', run: openRunHistory },
      { id: 'testclient', label: 'TestClient…', group: 'Запись и тест', run: openTestClientDialog },
      { id: 'validate', label: 'Проверить', group: 'Запись и тест', run: () => validateProject(false) },
      { id: 'validate-browser', label: 'Проверить в браузере', group: 'Запись и тест', run: () => validateProject(true) },
      ...(hasVanessaPlugin()
        ? [
            { id: 'vanessa-dry', label: 'Vanessa (dry)…', group: 'Запись и тест', run: () => openVanessaDialog(true) },
            { id: 'vanessa', label: 'Vanessa run…', group: 'Запись и тест', run: () => openVanessaDialog(false) },
            { id: 'vanessa-rerun', label: 'Vanessa rerun-failed…', group: 'Запись и тест', run: () => openVanessaDialog(false, true) },
            { id: 'vanessa-settings', label: 'Настройки Vanessa…', group: 'Запись и тест', run: openVanessaSettingsDialog },
          ]
        : []),
      { id: 'plugins', label: 'Управление плагинами…', group: 'Плагины', run: () => (showPlugins = true) },
      { id: 'journal', label: 'Журнал', group: 'Вид', shortcut: 'Ctrl+`', run: () => { bottomPanelOpen = true; bottomTab = 'journal' } },
      { id: 'results', label: 'Результаты', group: 'Вид', run: () => { bottomPanelOpen = true; bottomTab = 'results' } },
      { id: 'validate-panel', label: 'Проверка селекторов', group: 'Вид', run: () => { bottomPanelOpen = true; bottomTab = 'validate' } },
      { id: 'error-panel', label: 'Ошибка', group: 'Вид', run: () => { bottomPanelOpen = true; bottomTab = 'error' } },
      { id: 'explorer', label: 'Сценарии (explorer)', group: 'Вид', run: () => { sidebarVisible = true; saveLayout({ sidebarVisible: true }) } },
      { id: 'explorer-hide', label: 'Скрыть сценарии', group: 'Вид', run: () => { sidebarVisible = false; saveLayout({ sidebarVisible: false }) } },
      { id: 'preview', label: previewVisible ? 'Скрыть превью Gherkin' : 'Превью Gherkin', group: 'Вид', run: togglePreview },
      { id: 'steps-panel', label: stepsPanelVisible ? 'Скрыть панель шагов' : 'Панель шагов', group: 'Вид', run: toggleStepsPanel },
      { id: 'compact', label: compactLabel, group: 'Вид', run: () => { toolbarCompact = !toolbarCompact; persistSettings() } },
      { id: 'refactor-urls', label: 'Обновить стартовый URL…', group: 'Рефакторинг', run: refactorUpdateUrls },
      { id: 'hotkeys', label: 'Горячие клавиши', group: 'Справка', shortcut: 'Shift+F1', run: () => (showHotkeys = true) },
      { id: 'reset-layout', label: 'Сбросить макет окон', group: 'Вид', run: resetWindowLayout },
      { id: 'updates', label: 'Проверить обновления…', group: 'Справка', run: checkUpdates },
      { id: 'about', label: 'О программе', group: 'Справка', run: showAboutDialog },
    ]
  }

  function togglePreview() {
    previewVisible = !previewVisible
    saveLayout({ previewVisible })
  }

  function toggleStepsPanel() {
    stepsPanelVisible = !stepsPanelVisible
    stepsPanelCollapsed = !stepsPanelVisible
    persistSettings()
  }

  function showAboutDialog() {
    showAbout = true
  }

  function openFindReplace() {
    if (isWelcome) {
      appendLog('Откройте сценарий для поиска')
      return
    }
    showFindReplace = true
  }

  function resetWindowLayout() {
    const layout = resetUILayout()
    sidebarVisible = layout.sidebarVisible
    bottomPanelOpen = layout.bottomPanelOpen
    bottomPanelHeight = layout.bottomPanelHeight
    sidebarWidth = layout.sidebarWidth
    previewVisible = layout.previewVisible
    previewWidth = layout.previewWidth
    appendLog('Макет окон сброшен')
  }

  function hintDismissKey(hint: gui.ScenarioHintDTO): string {
    return `${hint.id}:${hint.stepIndex}`
  }

  function applyPostRecordHintFilter() {
    postRecordHints = postRecordAllHints.filter((h) => !postRecordDismissed.has(hintDismissKey(h)))
  }

  async function showPostRecordBanner(path: string) {
    postRecordPath = path
    postRecordDismissed = new Set()
    try {
      const content = await ReadFeature(path)
      postRecordStepCount = (await ParseEditorSteps(content)).length
      postRecordAllHints = await AnalyzeScenarioHints(content)
      applyPostRecordHintFilter()
    } catch {
      postRecordStepCount = 0
      postRecordAllHints = []
      postRecordHints = []
    }
  }

  async function refreshPostRecordHints() {
    if (!postRecordPath) return
    try {
      const content = activeTab === postRecordPath ? editorText : await ReadFeature(postRecordPath)
      postRecordAllHints = await AnalyzeScenarioHints(content)
      applyPostRecordHintFilter()
    } catch {
      /* keep previous */
    }
  }

  function dismissPostRecord() {
    postRecordPath = ''
    postRecordStepCount = 0
    postRecordAllHints = []
    postRecordHints = []
    postRecordDismissed = new Set()
  }

  async function postRecordValidate() {
    if (!postRecordPath) return
    if (activeTab !== postRecordPath) await loadFeature(postRecordPath)
    await validateProject(false)
    bottomPanelOpen = true
    bottomTab = 'validate'
  }

  async function postRecordSave() {
    if (!postRecordPath) return
    if (activeTab !== postRecordPath) await loadFeature(postRecordPath)
    await saveFeature()
    dismissPostRecord()
  }

  async function postRecordFixHint(hint: gui.ScenarioHintDTO) {
    if (!postRecordPath) return
    if (activeTab !== postRecordPath) await loadFeature(postRecordPath)
    const result = await ApplyScenarioHintFix({
      text: editorText,
      hintId: hint.id,
      stepIndex: hint.stepIndex,
    })
    if (result.count > 0) {
      editorText = result.text
      syncActiveTabContent()
      await validateEditor()
      appendLog(`Исправлена подсказка: ${hint.title}`)
      await refreshPostRecordHints()
    }
  }

  function postRecordFixHover() {
    const hint = postRecordHints.find((h) => h.id === 'menu_hover')
    if (hint) void postRecordFixHint(hint)
  }

  function postRecordShowStep(stepNo: number) {
    gotoEditorLine(stepNo)
    stepsPanelVisible = true
    stepsPanelCollapsed = false
  }

  function dismissPostRecordHint(hint: gui.ScenarioHintDTO) {
    postRecordDismissed.add(hintDismissKey(hint))
    applyPostRecordHintFilter()
  }

  async function duplicateFeature(path: string) {
    if (!path) return
    try {
      const target = await DuplicateFeature(path)
      await refreshProject()
      await loadFeature(target)
      appendLog(`Создана копия: ${basename(target)}`)
    } catch (e: any) {
      appendLog(`Ошибка: ${e}`)
    }
  }

  async function deleteFeature(path: string) {
    if (!path || isWelcome) return
    if (!confirm(`Удалить ${basename(path)}?`)) return
    try {
      await DeleteFeature(path)
      closeTab(path)
      await refreshProject()
      appendLog(`Удалён: ${basename(path)}`)
    } catch (e: any) {
      appendLog(`Ошибка: ${e}`)
    }
  }

  async function runFeatureFile(path: string) {
    if (!path) return
    await loadFeature(path)
    await executeRun({ ...lastRun, dryRun: false }, [path])
  }

  function onFileContextMenu(e: MouseEvent, path: string) {
    e.preventDefault()
    contextMenu = { x: e.clientX, y: e.clientY, path }
  }

  function dismissContextMenu() {
    contextMenu = null
  }

  function contextMenuRun() {
    if (!contextMenu) return
    const path = contextMenu.path
    dismissContextMenu()
    runFeatureFile(path)
  }

  function contextMenuOpen() {
    if (!contextMenu) return
    const path = contextMenu.path
    dismissContextMenu()
    loadFeature(path)
  }

  function contextMenuDuplicate() {
    if (!contextMenu) return
    const path = contextMenu.path
    dismissContextMenu()
    duplicateFeature(path)
  }

  function contextMenuDelete() {
    if (!contextMenu) return
    const path = contextMenu.path
    dismissContextMenu()
    deleteFeature(path)
  }

  async function confirmProjectReplace() {
    if (!findText) return
    projectReplaceBusy = true
    try {
      const result = await ReplaceInProject({
        find: findText,
        replace: replaceText,
        caseSensitive: replaceCaseSensitive,
      })
      appendLog(`Замена: ${result.replacements} в ${result.filesChanged} файлах`)
      showProjectReplace = false
      await refreshProject()
      if (activeTab && !isWelcome) {
        editorText = await ReadFeature(activeTab)
        tabs = tabs.map((t) => (t.path === activeTab ? { ...t, content: editorText } : t))
        validateEditor()
      }
    } catch (e: any) {
      appendLog(`Ошибка: ${e}`)
    } finally {
      projectReplaceBusy = false
    }
  }

  async function refreshRunResults() {
    try {
      runResults = await ListRunResults(50)
      lastErrorEntry = runResults.find((e) => !e.success) || null
    } catch {
      runResults = []
    }
  }

  async function openRunHistory() {
    await refreshRunResults()
    showRunHistory = true
  }

  async function openFeatureFromHistory(path: string) {
    showRunHistory = false
    await loadFeature(path)
  }

  async function openProjectAt(path: string) {
    if (!path) return
    try {
      const info = await OpenProject(path)
      projectPath = info.path
      features = info.features || []
      tags = info.tags || []
      featureTags = info.featureTags || {}
      testClients = await ListTestClients().catch(() => [])
      await rememberProject(projectPath)
      const recents = await loadRecents()
      recentProjects = recents.projects
      appendLog(`Проект открыт: ${projectPath}`)
      syncIdleStatus()
      await refreshRunResults()
      await refreshArtifacts()
    } catch (e: any) {
      appendLog(`Ошибка: ${e}`)
      setStatus(String(e), 'error')
    }
  }

  async function openExamples() {
    const examples = await BundledExamplesPath()
    if (!examples) {
      appendLog('Папка examples не найдена')
      return
    }
    await openProjectAt(examples)
    selectTab(WELCOME_KEY)
  }

  async function rerunFailed() {
    const failed = runResults.filter((e) => !e.success)
    const paths = [...new Set(failed.map((e) => e.path.split('::')[0]).filter(Boolean))]
    if (!paths.length) {
      appendLog('Нет упавших сценариев для перезапуска')
      return
    }
    await executeRun({ ...lastRun, dryRun: false }, paths)
  }

  function toggleBatchFeature(path: string) {
    if (batchSelected.includes(path)) {
      batchSelected = batchSelected.filter((p) => p !== path)
    } else {
      batchSelected = [...batchSelected, path]
    }
  }

  function onCatalogCollapse(key: string, collapsed: boolean) {
    const next = new Set(catalogCollapsed)
    if (collapsed) next.add(key)
    else next.delete(key)
    catalogCollapsed = next
  }

  function onCatalogActivate(path: string, kind: 'root' | 'dir' | 'file') {
    if (kind === 'file') loadFeature(path)
  }

  function clearBatchSelection() {
    batchSelected = []
  }

  async function runBatchSelected() {
    if (!batchSelected.length) return
    await executeRun({ ...lastRun, dryRun: false }, batchSelected)
  }

  function startResizePreview(e: MouseEvent) {
    resizingPreview = true
    e.preventDefault()
    const startX = e.clientX
    const startW = previewWidth
    const onMove = (ev: MouseEvent) => {
      previewWidth = Math.max(200, Math.min(720, startW - (ev.clientX - startX)))
    }
    const onUp = () => {
      resizingPreview = false
      window.removeEventListener('mousemove', onMove)
      window.removeEventListener('mouseup', onUp)
      saveLayout({ previewWidth })
    }
    window.addEventListener('mousemove', onMove)
    window.addEventListener('mouseup', onUp)
  }

  function startResizeSidebar(e: MouseEvent) {
    resizingSidebar = true
    e.preventDefault()
    const startX = e.clientX
    const startW = sidebarWidth
    const onMove = (ev: MouseEvent) => {
      sidebarWidth = Math.max(120, Math.min(480, startW + (ev.clientX - startX)))
    }
    const onUp = async () => {
      resizingSidebar = false
      window.removeEventListener('mousemove', onMove)
      window.removeEventListener('mouseup', onUp)
      saveLayout({ sidebarWidth })
      await persistSettings()
    }
    window.addEventListener('mousemove', onMove)
    window.addEventListener('mouseup', onUp)
  }

  function startResizeBottom(e: MouseEvent) {
    resizingBottom = true
    e.preventDefault()
    const startY = e.clientY
    const startH = bottomPanelHeight
    const onMove = (ev: MouseEvent) => {
      bottomPanelHeight = Math.max(100, Math.min(window.innerHeight * 0.6, startH + (startY - ev.clientY)))
    }
    const onUp = () => {
      resizingBottom = false
      saveLayout({ bottomPanelHeight })
      window.removeEventListener('mousemove', onMove)
      window.removeEventListener('mouseup', onUp)
    }
    window.addEventListener('mousemove', onMove)
    window.addEventListener('mouseup', onUp)
  }

  function startResizeSteps(e: MouseEvent) {
    resizingSteps = true
    e.preventDefault()
    const startY = e.clientY
    const startH = stepsPanelHeight
    const onMove = (ev: MouseEvent) => {
      stepsPanelHeight = Math.max(80, Math.min(480, startH + (startY - ev.clientY)))
    }
    const onUp = async () => {
      resizingSteps = false
      window.removeEventListener('mousemove', onMove)
      window.removeEventListener('mouseup', onUp)
      await persistSettings()
    }
    window.addEventListener('mousemove', onMove)
    window.addEventListener('mouseup', onUp)
  }

  function appendLog(line: string) {
    logText += line + (line.endsWith('\n') ? '' : '\n')
    bottomPanelOpen = true
    bottomTab = 'journal'
  }

  function setStatus(msg: string, tone: typeof statusTone = 'normal') {
    statusMessage = msg
    statusTone = tone
  }

  function applyDevUiMock() {
    if (!import.meta.env.DEV) return
    const mode = new URLSearchParams(location.search).get('mock')
    if (mode !== 'python') return
    projectPath = 'C:/Users/bafgion/Documents/Projects/camel-1c-integration/target'
    features = []
    recentFeatures = [
      'C:/Users/bafgion/Documents/Projects/camel-1c-integration/target/smoke.feature',
      'C:/Users/bafgion/Documents/Projects/camel-1c-integration/target/login.feature',
      'C:/Users/bafgion/Documents/Projects/camel-1c-integration/target/api.feature',
      'C:/Users/bafgion/Documents/Projects/camel-1c-integration/target/ui.feature',
      'C:/Users/bafgion/Documents/Projects/camel-1c-integration/target/regress.feature',
    ]
    recentProjects = ['target', 'examples', 'test', 'demo', 'sandbox']
  }

  function syncToolbarDensity() {
    if (!actionBarEl) return
    const urlBlock = actionBarEl.querySelector('.url-block') as HTMLElement | null
    const urlWidth = urlBlock?.getBoundingClientRect().width ?? 300
    const available = actionBarEl.getBoundingClientRect().width - urlWidth - 48
    toolbarIconOnly = available < 880
  }

  function observeActionBar(node: HTMLElement) {
    actionBarEl = node
    const ro = new ResizeObserver(() => {
      window.requestAnimationFrame(syncToolbarDensity)
    })
    ro.observe(node)
    window.requestAnimationFrame(() => window.requestAnimationFrame(syncToolbarDensity))
    return {
      destroy() {
        ro.disconnect()
      },
    }
  }

  function syncIdleStatus() {
    if (statusTone === 'busy' || statusTone === 'error') return
    if (projectPath) {
      statusMessage = projectPath.replace(/\\/g, '/')
      statusTone = 'normal'
      return
    }
    statusMessage = 'Проект → Открыть проект…'
    statusTone = 'normal'
  }

  function toggleMenu(name: string, e: MouseEvent) {
    e.stopPropagation()
    openMenu = openMenu === name ? null : name
  }

  function closeMenu() {
    openMenu = null
  }

  async function refreshArtifacts() {
    try {
      projectArtifacts = await ProjectArtifacts()
    } catch {
      projectArtifacts = new gui.ProjectArtifacts()
    }
  }

  async function refreshEditorSteps() {
    try {
      editorSteps = await ParseEditorSteps(editorText)
    } catch {
      editorSteps = []
    }
  }

  async function openArtifactPath(path: string) {
    if (!path) return
    try {
      await OpenFolder(path)
    } catch (e: any) {
      appendLog(`Не удалось открыть: ${e}`)
    }
  }

  async function refreshProject() {
    if (!projectPath) return
    const info = await OpenProject(projectPath)
    features = info.features || []
    tags = info.tags || []
    featureTags = info.featureTags || {}
    testClients = await ListTestClients().catch(() => [])
    await refreshInstalledPlugins()
  }

  async function refreshInstalledPlugins() {
    if (!projectPath) {
      installedPlugins = []
      return
    }
    try {
      installedPlugins = await ListPlugins()
    } catch {
      installedPlugins = []
    }
  }

  function hasVanessaPlugin(): boolean {
    return installedPlugins.some((p) => p.vanessa)
  }

  function pluginLabel(plugin: gui.PluginEntryDTO): string {
    if (plugin.vanessa) return 'Vanessa'
    return plugin.description || plugin.id || plugin.name
  }

  async function moveFeatureInCatalog(src: string, destDir: string) {
    if (!src || !destDir) return
    try {
      const newPath = await MoveFeature(src, destDir)
      const wasActive = activeTab === src
      if (tabs.some((t) => t.path === src)) {
        tabs = tabs.map((t) => (t.path === src ? { ...t, path: newPath } : t))
        if (wasActive) activeTab = newPath
      }
      batchSelected = batchSelected.map((p) => (p === src ? newPath : p))
      await refreshProject()
      if (wasActive) await loadFeature(newPath)
      appendLog(`Перемещён: ${basename(newPath)}`)
    } catch (e: any) {
      appendLog(`Ошибка: ${e}`)
    }
  }

  async function importDroppedFeatures(destDir: string, paths: string[]) {
    if (!destDir || !paths.length) return
    try {
      const imported = await ImportFeatures(destDir, paths)
      await refreshProject()
      appendLog(`Импортировано файлов: ${imported.length}`)
      if (imported.length === 1) await loadFeature(imported[0])
    } catch (e: any) {
      appendLog(`Ошибка импорта: ${e}`)
    }
  }

  async function openProjectDialog() {
    let path = ''
    try {
      path = await PickProjectFolder()
    } catch {
      path = ''
    }
    if (!path) {
      path = prompt('Путь к папке проекта:', projectPath || '') || ''
    }
    if (!path) return
    await openProjectAt(path)
  }

  function syncActiveTabContent() {
    if (isWelcome) return
    tabs = tabs.map((t) => (t.path === activeTab ? { ...t, content: editorText, dirty: true } : t))
  }

  async function loadFeature(path: string) {
    syncActiveTabContent()
    const existing = tabs.find((t) => t.path === path)
    if (existing) {
      activeTab = path
      editorText = existing.content
      await validateEditor()
      return
    }
    try {
      const content = await ReadFeature(path)
      tabs = [...tabs, { path, content, dirty: false }]
      activeTab = path
      editorText = content
      await rememberFeature(path)
      const recents = await loadRecents()
      recentFeatures = recents.features
      await validateEditor()
    } catch (e: any) {
      appendLog(`Ошибка открытия: ${e}`)
    }
  }

  function selectTab(path: string) {
    if (path === WELCOME_KEY) {
      welcomeTabVisible = true
      syncActiveTabContent()
      activeTab = WELCOME_KEY
      return
    }
    if (path === activeTab) return
    loadFeature(path)
  }

  function closeWelcomeTab() {
    if (tabs.length > 0) {
      welcomeTabVisible = false
      if (activeTab === WELCOME_KEY) {
        const next = tabs[tabs.length - 1]
        activeTab = next.path
        editorText = next.content
        validateEditor()
      }
      return
    }
    // Как в Python: единственная вкладка «Старт» — закрыть нельзя, остаётся на месте.
  }

  function closeTab(path: string, event?: Event) {
    event?.stopPropagation()
    if (path === activeTab && !isWelcome) {
      syncActiveTabContent()
    }
    const tab = tabs.find((t) => t.path === path)
    if (tab?.dirty) {
      pendingCloseTab = path
      return
    }
    finalizeCloseTab(path)
  }

  function finalizeCloseTab(path: string) {
    tabs = tabs.filter((t) => t.path !== path)
    if (activeTab === path) {
      const next = tabs[tabs.length - 1]
      if (next) {
        activeTab = next.path
        editorText = next.content
        validateEditor()
      } else {
        welcomeTabVisible = true
        activeTab = WELCOME_KEY
        editorText = ''
        stepStatus = '0 шагов'
        stepStatusError = false
      }
    }
  }

  async function saveAndCloseTab() {
    if (!pendingCloseTab) return
    const path = pendingCloseTab
    pendingCloseTab = null
    if (activeTab !== path) {
      await loadFeature(path)
    }
    await saveFeature()
    const tab = tabs.find((t) => t.path === path)
    if (tab && !tab.dirty) {
      finalizeCloseTab(path)
    }
  }

  function discardAndCloseTab() {
    if (!pendingCloseTab) return
    const path = pendingCloseTab
    pendingCloseTab = null
    finalizeCloseTab(path)
  }

  function cancelCloseTab() {
    pendingCloseTab = null
  }

  async function saveFeature() {
    if (!activeTab || isWelcome) return
    syncActiveTabContent()
    try {
      await SaveFeature(activeTab, editorText)
      tabs = tabs.map((t) => (t.path === activeTab ? { ...t, content: editorText, dirty: false } : t))
      appendLog(`Сохранено: ${basename(activeTab)}`)
      setStatus('Сохранено', 'success')
    } catch (e: any) {
      appendLog(`Ошибка сохранения: ${e}`)
      setStatus('Ошибка сохранения', 'error')
    }
  }

  async function validateEditor() {
    try {
      const issues = await ValidateFeature(editorText)
      editorValidationIssues = issues || []
      monaco?.setMarkers(editorValidationIssues)
      await refreshEditorSteps()
      if (editorValidationIssues.length > 0) {
        stepStatus = `${stepCount} шагов · ошибок ${editorValidationIssues.length}`
        stepStatusError = true
        setStatus('Ошибка в сценарии', 'error')
      } else {
        stepStatus = `${stepCount} шагов`
        stepStatusError = false
      }
    } catch {
      editorValidationIssues = []
      stepStatus = `${stepCount} шагов`
      await refreshEditorSteps()
    }
  }

  function gotoEditorLine(line: number) {
    monaco?.gotoLine(line)
  }

  function validateProjectHint(): string {
    if (editorValidationIssues.length > 0) return ''
    if (isWelcome) return 'Откройте сценарий для проверки шагов'
    return 'Ошибок в шагах сценария нет. Проверка селекторов в браузере — меню «Запись и тест».'
  }

  async function onEditorChange(text: string) {
    editorText = text
    syncActiveTabContent()
    await validateEditor()
  }

  function parseVars(text: string): Record<string, string> {
    const out: Record<string, string> = {}
    for (const line of text.split('\n')) {
      const trimmed = line.trim()
      if (!trimmed || trimmed.startsWith('#')) continue
      const eq = trimmed.indexOf('=')
      if (eq <= 0) continue
      out[trimmed.slice(0, eq).trim()] = trimmed.slice(eq + 1).trim()
    }
    return out
  }

  function openRunDialog(title: string, defaults: Partial<RunForm>) {
    runDialogTitle = title
    runForm = { ...lastRun, ...defaults }
    showRun = true
  }

  function openVanessaDialog(dry: boolean, preferRerun = false) {
    vanessaDry = dry
    vanessaTag = ''
    vanessaExcludeTags = ''
    vanessaScenario = ''
    vanessaRerunDir = ''
    vanessaPreferRerun = preferRerun
    vanessaInstallEpf = false
    vanessaEpfUrl = ''
    vanessaEpfDest = ''
    showVanessaRun = true
  }

  async function confirmVanessaRun() {
    showVanessaRun = false
    const exclude = vanessaExcludeTags
      .split(',')
      .map((t) => t.trim())
      .filter(Boolean)
    let rerunDir = vanessaRerunDir.trim()
    await runPlugin('vanessa', vanessaDry, {
      tag: vanessaTag.trim(),
      excludeTags: exclude,
      scenario: vanessaScenario.trim(),
      rerunFailedRunDir: rerunDir,
      installEpf: vanessaInstallEpf,
      epfUrl: vanessaEpfUrl.trim(),
      epfDest: vanessaEpfDest.trim(),
    })
  }

  async function executeRun(opts: RunForm, targets: string[] = []) {
    if (!projectPath) return
    lastRun = { ...opts }
    showRun = false

    const allureDir = opts.allure ? scenariaSubdir('allure-results') : ''
    const traceDir = !opts.dryRun && opts.trace ? scenariaSubdir('traces') : ''
    const videoDir = !opts.dryRun && opts.video ? scenariaSubdir('videos') : ''

    const htmlPath = opts.html ? scenariaSubdir('report.html') : ''

    const junitPath = opts.junit ? scenariaSubdir('junit.xml') : ''

    playing = !opts.dryRun
    setStatus('▶ Идёт тест', 'busy')
    appendLog(targets.length ? `Запуск ${targets.length} сценариев…` : opts.dryRun ? 'Dry-run…' : 'Запуск Playwright…')
    bottomPanelOpen = true
    bottomTab = 'journal'

    const result = await Run({
      tag: opts.tag,
      testClient: opts.testClient,
      vars: parseVars(opts.vars),
      dryRun: opts.dryRun,
      headed: opts.headed,
      engine: opts.dryRun ? '' : opts.engine,
      installPlaywright: opts.installPW,
      allureDir,
      traceDir,
      videoDir,
      htmlPath,
      junitPath,
      browser: opts.dryRun ? '' : opts.browser || settingsBrowser || 'chromium',
      workers: opts.workers || settingsWorkers || 1,
      slowMo: opts.dryRun ? 0 : opts.slowMo || 0,
      targets,
    })
    playing = false
    if (result.output) appendLog(result.output.trimEnd())
    if (result.error) {
      appendLog(`Ошибка: ${result.error}`)
      setStatus('Ошибка теста', 'error')
      bottomTab = 'error'
    } else {
      appendLog('Завершено.')
      setStatus('Тест завершён', 'success')
      welcomePlayedSuccess = true
      bottomTab = 'results'
    }
    await refreshRunResults()
    await refreshArtifacts()
    if (!result.error && opts.html && htmlPath) {
      try {
        if (await ArtifactExists(htmlPath)) {
          await OpenFolder(htmlPath)
        }
      } catch {
        /* ignore */
      }
    }
  }

  function scenariaSubdir(sub: string): string {
    if (!projectPath) return ''
    return `${projectPath.replace(/\\/g, '/')}/.scenaria/${sub}`
  }

  function confirmRun() {
    executeRun(runForm)
  }

  async function validateProject(browser: boolean) {
    if (!projectPath) return
    appendLog(browser ? 'Проверка в браузере…' : 'Проверка…')
    bottomPanelOpen = true
    bottomTab = 'validate'
    const result = await Validate(settingsBrowser || 'chromium', !browser)
    if (result.output) appendLog(result.output.trimEnd())
    if (result.error) appendLog(`Ошибка: ${result.error}`)
    else appendLog('Проверка завершена.')
  }

  async function initProject() {
    const out = await InitProject()
    if (out) appendLog(out.trimEnd())
    await refreshProject()
  }

  async function openTestClientDialog() {
    if (!projectPath) return
    testClients = await ListTestClients().catch(() => [])
    testClientSelection = runForm.testClient || testClients[0] || ''
    showTestClient = true
  }

  function useTestClient(name: string) {
    testClientSelection = name
    lastRun.testClient = name
    runForm.testClient = name
    appendLog(`TestClient: ${name}`)
    showTestClient = false
  }

  async function openStepsDialog() {
    if (recordingBlocksManualTools()) {
      appendLog('Поставьте запись на паузу, чтобы вставить шаг')
      return
    }
    showSteps = true
  }

  function openVanessaSettingsDialog() {
    if (!projectPath) return
    showVanessaSettings = true
  }

  function openStepsHelp(query = '') {
    if (recordingBlocksManualTools()) {
      appendLog('Поставьте запись на паузу, чтобы открыть справку по шагам')
      return
    }
    stepsHelpQuery = query
    showStepsHelp = true
  }

  function openStepHelpFromPanel(step: EditorStepRow) {
    if (!step.text) return
    openStepsHelp(step.text)
  }

  function insertStep(template: string) {
    const line = template.endsWith('\n') ? template : template + '\n'
    monaco?.insertAtCursor(line)
    showSteps = false
    showStepsHelp = false
  }

  function recordingBlocksManualTools(): boolean {
    return recording && !recordPaused
  }

  function openSnippetPalette() {
    if (recordingBlocksManualTools()) {
      appendLog('Поставьте запись на паузу, чтобы вставить шаг из палитры')
      setStatus('Запись активна', 'busy')
      return
    }
    showSnippetPalette = true
  }

  async function applyEditorText(text: string) {
    editorText = text
    syncActiveTabContent()
    await validateEditor()
  }

  async function refactorUpdateUrls() {
    if (isWelcome) return
    const newUrl = prompt('Новый URL для всех шагов «открыт»:', startURL || recordURL)
    if (!newUrl?.trim()) return
    const result = await RefactorUpdateStartURLs(editorText, newUrl.trim())
    if (result.count <= 0) {
      appendLog('Шаги «открыт» не найдены')
      return
    }
    await applyEditorText(result.text)
    appendLog(`Обновлено URL в ${result.count} шагах`)
  }

  async function refactorNormalizeIndents() {
    if (isWelcome) return
    const text = await RefactorNormalizeIndents(editorText)
    await applyEditorText(text)
    appendLog('Отступы нормализованы')
  }

  async function refactorCollapseBlank() {
    if (isWelcome) return
    const text = await RefactorCollapseBlankLines(editorText)
    await applyEditorText(text)
    appendLog('Пустые строки между шагами убраны')
  }

  function focusBrowserWindow() {
    focusBrowser()
  }

  function openExportDialog() {
    if (!activeTab || isWelcome) return
    exportInputPath = activeTab
    showExport = true
  }

  function openImportDialog() {
    if (!projectPath) return
    showImport = true
  }

  async function onImportComplete(featurePath: string) {
    await refreshProject()
    await loadFeature(featurePath)
  }

  async function runPlugin(name: string, dry: boolean, opts: Partial<gui.PluginRunRequest> = {}) {
    const label = pluginLabel(installedPlugins.find((p) => p.name === name) || { name, id: name, vanessa: false } as gui.PluginEntryDTO)
    appendLog(dry ? `${label} (dry)…` : `${label}…`)
    const result = await RunPlugin({
      name,
      dryRun: dry,
      tag: opts.tag || '',
      excludeTags: opts.excludeTags || [],
      scenario: opts.scenario || '',
      rerunFailedRunDir: opts.rerunFailedRunDir || '',
      installEpf: opts.installEpf || false,
      epfUrl: opts.epfUrl || '',
      epfDest: opts.epfDest || '',
    })
    if (result.output) appendLog(result.output.trimEnd())
    if (result.error) appendLog(`Ошибка: ${result.error}`)
  }

  async function checkUpdatesOnStartup() {
    if (!settingsCheckUpdatesOnStartup) return
    try {
      const result = await CheckUpdate()
      const out = (result.output || '').trim()
      if (result.error || !out) return
      if (out.includes('Update available')) {
        appendLog(out)
        setStatus('Доступно обновление — Справка → Проверить обновления', 'normal')
      }
    } catch {
      /* offline or dev without wails */
    }
  }

  async function checkUpdates() {
    appendLog('Проверка обновлений…')
    const result = await CheckUpdate()
    if (result.output) appendLog(result.output.trimEnd())
    if (result.error) appendLog(`Ошибка: ${result.error}`)
    else appendLog('Проверка завершена.')
  }

  async function openSettings() {
    const s = await LoadSettings()
    applySettingsFromDTO(s)
    showSettings = true
  }

  async function persistSettings() {
    await SaveSettings({
      browser: settingsBrowser,
      headless: settingsHeadless,
      parallelWorkers: settingsWorkers,
      maxLoopIterations: settingsLoops,
      filterRecording,
      navOnlyRecording,
      hoverRecord,
      toolbarCompact,
      stepsPanelVisible,
      stepsPanelHeight,
      sidebarWidth,
      recentProjects,
      recentFeatures,
      checkUpdatesOnStartup: settingsCheckUpdatesOnStartup,
    })
  }

  async function applySettings() {
    stepsPanelCollapsed = !stepsPanelVisible
    await persistSettings()
    showSettings = false
    appendLog('Настройки сохранены.')
  }

  function beginRecord() {
    recordAppendTo = ''
    recordTestClient = runForm.testClient || testClientSelection || ''
    if (projectPath) {
      recordOutput = `${projectPath.replace(/\\/g, '/')}/recorded.feature`
    }
    recordURL = startURL
    showRecord = true
  }

  async function startRecord() {
    await persistSettings()
    lastRecordTarget = recordAppendTo || recordOutput
    await StartRecord({
      url: recordURL,
      output: recordOutput,
      idleSeconds: recordIdle,
      headless: settingsHeadless,
      filterRecording,
      navOnlyRecording,
      hoverRecord,
      appendTo: recordAppendTo,
      testClient: recordTestClient,
    })
    recordAppendTo = ''
  }

  async function toggleRecordPause() {
    if (recordPaused) {
      await ResumeRecording()
      recordPaused = false
      setStatus('● Идёт запись', 'busy')
    } else {
      await PauseRecording()
      recordPaused = true
      setStatus('⏸ Пауза', 'busy')
    }
  }

  async function stopRecord() {
    await CancelRecording()
  }

  async function submitOtp(code: string) {
    await SubmitOTPCode(code)
    showOtp = false
  }

  async function cancelOtp() {
    await CancelOTP()
    showOtp = false
  }

  async function newScenario() {
    const out = await PickSaveFile('Новый сценарий', 'scenario.feature')
    if (!out) return
    const template =
      '# language: ru\nFeature: Новый сценарий\n\n  Scenario: Пример\n    Допустим открыт "https://example.com"\n'
    try {
      await SaveFeature(out, template)
      await refreshProject()
      await loadFeature(out)
    } catch (e: any) {
      appendLog(`Ошибка: ${e}`)
    }
  }

  async function openFileDialog() {
    const path = await PickOpenFile('Открыть feature')
    if (path) await loadFeature(path)
  }

  function insertTemplate() {
    const template =
      'Feature: Новый сценарий\n\n  Scenario: Пример\n    Допустим открыт "https://example.com"\n'
    if (isWelcome) {
      newScenario()
      return
    }
    monaco?.insertAtCursor(template)
  }

  function quickStart() {
    recordURL = startURL
    beginRecord()
  }

  function continueRecord() {
    if (!projectPath) {
      appendLog('Сначала откройте проект')
      return
    }
    if (!activeTab || isWelcome || !activeTab.toLowerCase().endsWith('.feature')) {
      appendLog('Откройте .feature файл для дозаписи в конец сценария')
      beginRecord()
      return
    }
    recordAppendTo = activeTab
    recordOutput = activeTab
    recordURL = startURL
    showRecord = true
    appendLog(`Дозапись в ${basename(activeTab)}`)
  }

  function quickRecord() {
    quickStart()
  }

  async function focusBrowser() {
    try {
      await FocusBrowser()
      appendLog('Окно браузера выведено на передний план')
    } catch (e: any) {
      appendLog(`Показать браузер: ${e}`)
    }
  }

  function hostFromURL(url: string): string {
    try {
      const u = new URL(url.includes('://') ? url : `https://${url}`)
      return u.hostname.toLowerCase()
    } catch {
      return ''
    }
  }

  function openHttpAuthDialog() {
    httpAuthHost = hostFromURL(recordURL || startURL)
    showHttpAuth = true
  }

  async function pickElement() {
    if (!recording) {
      appendLog('Указать элемент: начните запись')
      return
    }
    if (!recordPaused) {
      appendLog('Указать элемент: поставьте запись на паузу')
      setStatus('Поставьте запись на паузу', 'busy')
      return
    }
    appendLog('Укажите элемент в окне браузера…')
    await focusBrowser()
    const result = await PickSelector()
    if (result.error) {
      if (result.error !== 'отменено') appendLog(`Указать элемент: ${result.error}`)
      return
    }
    if (!result.selector) return
    pickerSelector = result.selector
    pickerChoices = await PickerStepChoices(result.selector, 'Допустим')
    showPickerStep = true
  }

  function insertPickerStep(text: string) {
    if (recordingBlocksManualTools()) {
      appendLog('Поставьте запись на паузу, чтобы вставить шаг')
      return
    }
    monaco?.insertAtCursor(text.endsWith('\n') ? text : text + '\n')
    appendLog('Шаг вставлен из указателя элемента')
  }

  async function undoRecordedStep() {
    const ok = await UndoRecordedStep()
    if (ok) appendLog('Последний записанный шаг отменён')
    else appendLog('Нечего отменять')
  }

  function onWelcomeChecklistStep(step: number) {
    if (step === 1) openProjectDialog()
    else if (step === 2) quickStart()
    else if (step === 3) {
      if (tabs.length === 0) newScenario()
      else if (projectPath) executeRun({ ...lastRun, dryRun: false })
    }
  }

  function syncUrlFromField() {
    recordURL = startURL
  }

  function projectLabel(): string {
    if (isWelcome) return 'Старт'
    if (activeTab && activeTab !== WELCOME_KEY) return basename(activeTab)
    if (projectPath) return basename(projectPath)
    return 'Открыть проект…'
  }
</script>

<div class="ide" class:panel-open={bottomPanelOpen}>
  <!-- Menu bar (Python: Проект / Сценарий / Запись и тест / Вид / Справка) -->
  <div class="menubar" role="menubar">
    <div class="menu-root" class:open={openMenu === 'project'}>
      <button class="menu-trigger" on:click={(e) => toggleMenu('project', e)}>Проект</button>
      {#if openMenu === 'project'}
        <div class="menu-dropdown" on:click={closeMenu}>
          <button class="menu-item" on:click={openExamples}>Открыть примеры сценариев</button>
          <button class="menu-item" on:click={openProjectDialog}>Открыть проект…</button>
          <button class="menu-item" on:click={openSettings}>Настройки…<span class="menu-shortcut">Ctrl+,</span></button>
          <div class="menu-sep"></div>
          <button class="menu-item" on:click={initProject} disabled={!projectPath}>Init проекта</button>
        </div>
      {/if}
    </div>

    <div class="menu-root" class:open={openMenu === 'scenario'}>
      <button class="menu-trigger" on:click={(e) => toggleMenu('scenario', e)}>Сценарий</button>
      {#if openMenu === 'scenario'}
        <div class="menu-dropdown" on:click={closeMenu}>
          <button class="menu-item" on:click={newScenario}>Новый</button>
          <button class="menu-item" on:click={openFileDialog}>Открыть…</button>
          <button class="menu-item" on:click={saveFeature} disabled={isWelcome}>Сохранить<span class="menu-shortcut">Ctrl+S</span></button>
          <button class="menu-item" on:click={() => activeTab && !isWelcome && duplicateFeature(activeTab)} disabled={isWelcome}>Дублировать</button>
          <button class="menu-item" on:click={() => activeTab && !isWelcome && deleteFeature(activeTab)} disabled={isWelcome}>Удалить</button>
          <div class="menu-sep"></div>
          <button class="menu-item" on:click={openFindReplace} disabled={isWelcome}>Найти и заменить…<span class="menu-shortcut">Ctrl+H</span></button>
          <button class="menu-item" on:click={() => (showProjectReplace = true)} disabled={!projectPath}>Замена по проекту…</button>
          <div class="menu-sep"></div>
          <button class="menu-item" on:click={openExportDialog} disabled={isWelcome}>Экспорт…</button>
          <button class="menu-item" on:click={openImportDialog} disabled={!projectPath}>Импорт JSON…</button>
          <div class="menu-sep"></div>
          <button class="menu-item" on:click={openSnippetPalette}>Палитра сниппетов…<span class="menu-shortcut">Ctrl+Shift+Space</span></button>
          <button class="menu-item" on:click={openStepsDialog}>Вставить шаг…</button>
          <div class="menu-sep"></div>
          <button class="menu-item" on:click={refactorUpdateUrls} disabled={isWelcome}>Обновить стартовый URL…</button>
          <button class="menu-item" on:click={refactorNormalizeIndents} disabled={isWelcome}>Нормализовать отступы</button>
          <button class="menu-item" on:click={refactorCollapseBlank} disabled={isWelcome}>Убрать пустые строки между шагами</button>
        </div>
      {/if}
    </div>

    <div class="menu-root" class:open={openMenu === 'run'}>
      <button class="menu-trigger" on:click={(e) => toggleMenu('run', e)}>Запись и тест</button>
      {#if openMenu === 'run'}
        <div class="menu-dropdown" on:click={closeMenu}>
          <button class="menu-item" on:click={beginRecord} disabled={!projectPath}>Браузер<span class="menu-shortcut">Ctrl+B</span></button>
          <button class="menu-item" on:click={beginRecord} disabled={!projectPath}>Запись<span class="menu-shortcut">Ctrl+R</span></button>
          <button class="menu-item" on:click={stopRecord} disabled={!recording}>Стоп</button>
          <button class="menu-item" on:click={toggleRecordPause} disabled={!recording}>Пауза</button>
          <div class="menu-sep"></div>
          <button class="menu-item" on:click={openTestClientDialog} disabled={!projectPath}>TestClient…</button>
          <button class="menu-item" on:click={openHttpAuthDialog}>HTTP Auth…</button>
          <div class="menu-sep"></div>
          <button class="menu-item" on:click={() => executeRun({ ...lastRun, dryRun: false })} disabled={!projectPath}>
            Запустить<span class="menu-shortcut">Ctrl+Enter</span>
          </button>
          <button class="menu-item" on:click={() => executeRun({ ...lastRun, dryRun: false }, batchSelected)} disabled={!projectPath || !batchSelected.length}>
            Запустить выбранные
          </button>
          <button class="menu-item" on:click={rerunFailed} disabled={!projectPath}>Перезапустить упавшие</button>
          <button class="menu-item" on:click={openRunHistory} disabled={!projectPath}>История запусков…</button>
          <button class="menu-item" on:click={() => openRunDialog('Запуск сценария', {})} disabled={!projectPath}>Запустить…</button>
          <button class="menu-item" on:click={() => openRunDialog('Запуск с тегом', {})} disabled={!projectPath}>
            Запустить сценарии с тегом…
          </button>
          <button class="menu-item" on:click={() => executeRun({ ...lastRun, dryRun: true })} disabled={!projectPath}>Dry-run</button>
          <div class="menu-sep"></div>
          <button class="menu-item" on:click={() => validateProject(false)} disabled={!projectPath}>Проверить</button>
          <button class="menu-item" on:click={() => validateProject(true)} disabled={!projectPath}>Проверить в браузере</button>
          <div class="menu-sep"></div>
          {#if hasVanessaPlugin()}
            <button class="menu-item" on:click={() => openVanessaDialog(true)} disabled={!projectPath}>Vanessa (dry)…</button>
            <button class="menu-item" on:click={() => openVanessaDialog(false)} disabled={!projectPath}>Vanessa run…</button>
            <button class="menu-item" on:click={() => openVanessaDialog(false, true)} disabled={!projectPath}>Vanessa rerun-failed…</button>
            <button class="menu-item" on:click={openVanessaSettingsDialog} disabled={!projectPath}>Настройки Vanessa…</button>
          {/if}
        </div>
      {/if}
    </div>

    <div class="menu-root" class:open={openMenu === 'plugins'}>
      <button class="menu-trigger" on:click={(e) => toggleMenu('plugins', e)}>Плагины</button>
      {#if openMenu === 'plugins'}
        <div class="menu-dropdown" on:click={closeMenu}>
          <button class="menu-item" on:click={() => (showPlugins = true)} disabled={!projectPath}>Управление плагинами…</button>
          {#if installedPlugins.length > 0}
            <div class="menu-sep"></div>
            {#each installedPlugins as plugin (plugin.name)}
              {#if plugin.runnable}
                {#if plugin.vanessa}
                  <button class="menu-item" on:click={() => openVanessaDialog(true)} disabled={!projectPath}>Vanessa (dry)…</button>
                  <button class="menu-item" on:click={() => openVanessaDialog(false)} disabled={!projectPath}>Vanessa run…</button>
                {:else}
                  <button class="menu-item" on:click={() => runPlugin(plugin.name, false)} disabled={!projectPath}>Запустить {pluginLabel(plugin)}</button>
                {/if}
              {:else}
                <button class="menu-item" disabled title="Нет команды запуска в plugin.json">{pluginLabel(plugin)} — установлен</button>
              {/if}
            {/each}
          {:else if projectPath}
            <div class="menu-sep"></div>
            <button class="menu-item" disabled>Нет установленных плагинов</button>
          {/if}
        </div>
      {/if}
    </div>

    <div class="menu-root" class:open={openMenu === 'view'}>
      <button class="menu-trigger" on:click={(e) => toggleMenu('view', e)}>Вид</button>
      {#if openMenu === 'view'}
        <div class="menu-dropdown" on:click={closeMenu}>
          <button class="menu-item" on:click={() => selectTab(WELCOME_KEY)}>Старт</button>
          <button class="menu-item" on:click={() => { sidebarVisible = true; saveLayout({ sidebarVisible: true }) }}>Сценарии</button>
          <button class="menu-item" on:click={() => { sidebarVisible = false; saveLayout({ sidebarVisible: false }) }}>Скрыть сценарии</button>
          <button class="menu-item" on:click={() => (showCommandPalette = true)}>Палитра команд<span class="menu-shortcut">Ctrl+Shift+P</span></button>
          <button class="menu-item" on:click={() => { bottomPanelOpen = true; bottomTab = 'journal' }}>Журнал</button>
          <button class="menu-item" on:click={() => { bottomPanelOpen = true; bottomTab = 'results' }}>Результаты</button>
          <button class="menu-item" on:click={openRunHistory} disabled={!projectPath}>История запусков…</button>
          <button class="menu-item" on:click={() => { bottomPanelOpen = true; bottomTab = 'validate' }}>Проверка селекторов</button>
          <button class="menu-item" on:click={() => { bottomPanelOpen = true; bottomTab = 'error' }}>Ошибка</button>
          <button class="menu-item" on:click={togglePreview}>
            {previewVisible ? 'Скрыть превью Gherkin' : 'Превью Gherkin'}
          </button>
          <button class="menu-item" on:click={toggleStepsPanel}>
            {stepsPanelVisible ? 'Скрыть панель шагов' : 'Панель шагов'}
          </button>
          <button class="menu-item" on:click={() => { toolbarCompact = !toolbarCompact; persistSettings() }}>
            {toolbarCompact ? 'Расширенная панель' : 'Компактная панель'}
          </button>
          <button class="menu-item" on:click={resetWindowLayout}>Сбросить макет окон</button>
        </div>
      {/if}
    </div>

    <div class="menu-root" class:open={openMenu === 'help'}>
      <button class="menu-trigger" on:click={(e) => toggleMenu('help', e)}>Справка</button>
      {#if openMenu === 'help'}
        <div class="menu-dropdown" on:click={closeMenu}>
          <button class="menu-item" on:click={openStepsHelp}>Справка по шагам…<span class="menu-shortcut">F1</span></button>
          <button class="menu-item" on:click={() => (showHotkeys = true)}>Горячие клавиши<span class="menu-shortcut">Shift+F1</span></button>
          <button class="menu-item" on:click={checkUpdates}>Проверить обновления…</button>
          <button class="menu-item" on:click={showAboutDialog}>О программе</button>
        </div>
      {/if}
    </div>
  </div>

  <div class="ide-main">
    <div class="ide-center" class:no-sidebar={!sidebarVisible} style="--sidebar-width: {sidebarWidth}px">
      <!-- Activity bar -->
      <aside class="activity-bar">
        <button
          class="activity-btn"
          class:active={sidebarVisible}
          title="Сценарии"
          on:click={() => {
            sidebarVisible = !sidebarVisible
            saveLayout({ sidebarVisible })
          }}
        >
          {@html icons.explorer}
        </button>
        <button
          class="activity-btn"
          class:active={bottomPanelOpen}
          title="Панель вывода"
          on:click={() => {
            bottomPanelOpen = !bottomPanelOpen
            saveLayout({ bottomPanelOpen })
          }}
        >
          {@html icons.panel}
        </button>
      </aside>

      <!-- Explorer -->
      {#if sidebarVisible}
        <div class="sidebar-column" style="width: {sidebarWidth + 4}px">
        <aside class="explorer" style="width: {sidebarWidth}px">
          <div class="explorer-header">
            <p class="zone-title">СЦЕНАРИИ</p>
            <div class="explorer-tools">
              <input class="explorer-search" bind:value={sidebarSearch} placeholder="Поиск, @тег или tag:smoke" />
              <button class="icon-btn" title="Новый сценарий" on:click={newScenario}>{@html icons.plus}</button>
              <button class="icon-btn batch-toggle" class:active={batchMode} title="Пакетный запуск" on:click={() => (batchMode = !batchMode)}>
                Выбор
              </button>
            </div>
            {#if tags.length > 0}
              <div class="explorer-tag-chips">
                {#each tags.slice(0, 12) as tag}
                  <button
                    type="button"
                    class="chip"
                    class:active={sidebarSearch.trim() === tag || sidebarSearch.trim() === `@${tag.replace(/^@/, '')}`}
                    on:click={() => (sidebarSearch = tag.startsWith('@') ? tag : `@${tag}`)}
                  >
                    {tag.startsWith('@') ? tag : `@${tag}`}
                  </button>
                {/each}
              </div>
            {/if}
            {#if batchMode || projectPath}
              <p class="batch-selection">Выбрано для запуска: {batchCount}</p>
            {/if}
          </div>
          <div class="catalog catalog-panel">
            {#if catalogViewState.showEmptyMessage}
              <CatalogEmptyState
                title={catalogViewState.emptyTitle || ''}
                hint={catalogViewState.emptyHint || ''}
                kind={catalogViewState.emptyKind || 'no_project'}
              />
            {:else if catalogViewState.tree}
              <FeatureCatalogTree
                tree={catalogViewState.tree}
                activeFeature={activeFeaturePath}
                {batchSelected}
                {batchMode}
                expandAll={catalogViewState.expandAll}
                collapsed={catalogCollapsed}
                dropTarget={catalogDropTarget}
                onActivate={onCatalogActivate}
                onToggleBatch={(path) => {
                  if (!batchMode) batchMode = true
                  toggleBatchFeature(path)
                }}
                onCollapseChange={onCatalogCollapse}
                onFileContextMenu={onFileContextMenu}
                onMoveFeature={moveFeatureInCatalog}
                onDropTarget={(path) => (catalogDropTarget = path)}
              />
            {/if}
          </div>
        </aside>
        <div class="splitter-v" role="separator" on:mousedown={startResizeSidebar}></div>
        </div>
      {/if}

      <!-- Workspace -->
      <section class="workspace">
        <div class="action-bar" class:compact={toolbarCompact} use:observeActionBar>
          <div class="quick-toolbar">
            <div class="toolbar-row primary">
              <button class="tool-btn primary" title="Браузер (Ctrl+B)" on:click={beginRecord} disabled={!projectPath}>
                {@html toolbarIcons.browser()}<span>Браузер</span>
              </button>
              <button
                class="tool-btn primary"
                class:record-active={recording}
                title="Запись (Ctrl+R)"
                on:click={beginRecord}
                disabled={!projectPath || recording}
              >
                {@html toolbarIcons.record()}<span>Запись</span>
              </button>
              <button class="tool-btn primary" on:click={stopRecord} disabled={!recording} title="Остановить запись, тест или браузер">
                {@html toolbarIcons.stop()}<span>Стоп</span>
              </button>
              <button
                class="tool-btn primary primary-run"
                title="Запустить (Ctrl+Enter)"
                on:click={() => executeRun({ ...lastRun, dryRun: false })}
                disabled={!projectPath}
              >
                {@html toolbarIcons.play()}<span>Запустить</span>
              </button>
              <button class="tool-btn primary" on:click={saveFeature} disabled={isWelcome} title="Сохранить файл сценария (Ctrl+S)">
                {@html toolbarIcons.save()}<span>Сохранить</span>
              </button>
            </div>
            {#if !toolbarCompact}
            <div class="toolbar-row secondary" class:icon-only={toolbarIconOnly}>
              <button class="tool-btn" on:click={continueRecord} disabled={!projectPath || recording} title="Продолжить запись в конец сценария">
                {@html toolbarIcons.continueRecord()}<span>Дозапись</span>
              </button>
              <button class="tool-btn" on:click={toggleRecordPause} disabled={!recording} title="Приостановить запись">
                {@html toolbarIcons.pause()}<span>Пауза</span>
              </button>
              <span class="toolbar-sep" aria-hidden="true"></span>
              <button class="tool-btn" on:click={focusBrowser} disabled={!recording && !playing}>
                {@html toolbarIcons.browserFocus()}<span>Показать браузер</span>
              </button>
              <button class="tool-btn" on:click={() => validateProject(true)} disabled={!projectPath}>
                {@html toolbarIcons.validate()}<span>Селекторы на странице</span>
              </button>
              <button class="tool-btn" on:click={pickElement} disabled={!recording || !recordPaused} title={recording && !recordPaused ? 'Поставьте запись на паузу' : 'Указать элемент'}>
                {@html toolbarIcons.picker()}<span>Указать элемент</span>
              </button>
              <button class="tool-btn" on:click={quickRecord} disabled={!projectPath || recording}>
                {@html toolbarIcons.quickRecord()}<span>Быстрая запись</span>
              </button>
              <button class="tool-btn" on:click={validateEditor} disabled={isWelcome}>
                {@html toolbarIcons.gherkin()}<span>Синтаксис Gherkin</span>
              </button>
              <button class="tool-btn" on:click={undoRecordedStep} disabled={!recording}>
                {@html toolbarIcons.undo()}<span>Отменить шаг</span>
              </button>
              <span class="toolbar-sep" aria-hidden="true"></span>
              <button class="tool-btn" on:click={() => { bottomPanelOpen = true; bottomTab = 'journal' }}>
                {@html toolbarIcons.log()}<span>Журнал</span>
              </button>
              <button class="tool-btn" on:click={() => { bottomPanelOpen = true; bottomTab = 'results' }}>
                {@html toolbarIcons.results()}<span>Результаты</span>
              </button>
            </div>
            {/if}
          </div>

          {#if !isWelcome && activeFeatureTab}
            <div class="scenario-chip">
              <span class="name">{basename(activeTab)}</span>
              {#if activeFeatureTab.dirty}<span class="badge">*</span>{/if}
            </div>
          {/if}

          <span class="action-bar-sep" aria-hidden="true"></span>

          <div class="url-block">
            <span>URL</span>
            <input bind:value={startURL} placeholder="https://site.com" on:change={syncUrlFromField} />
            <button class="icon-btn" title="URL из вкладки браузера" on:click={() => (recordURL = startURL)}>
              {@html icons.external}
            </button>
          </div>
        </div>

        <EditorTabBar
          activeKey={activeTab}
          welcomeKey={WELCOME_KEY}
          welcomeVisible={welcomeTabVisible}
          {tabs}
          tabLabel={basename}
          onSelect={selectTab}
          onClose={(path) => closeTab(path)}
          onCloseWelcome={closeWelcomeTab}
        />

        <div class="editor-stack">
          {#if isWelcome}
            <WelcomePanel
              bind:startURL
              {recentProjects}
              {recentFeatures}
              projectOpen={welcomeProjectOpen}
              recorded={welcomeRecorded}
              playedSuccess={welcomePlayedSuccess}
              onOpenProject={openProjectDialog}
              onQuickStart={quickStart}
              onNewScenario={newScenario}
              onOpenFile={openFileDialog}
              onInsertTemplate={insertTemplate}
              onOpenExamples={openExamples}
              onOpenRecentProject={openProjectAt}
              onOpenRecentFeature={(path) => loadFeature(path)}
              onChecklistStep={onWelcomeChecklistStep}
            />
          {:else}
            {#if postRecordPath}
              <PostRecordBanner
                path={postRecordPath}
                stepCount={postRecordStepCount}
                hints={postRecordHints}
                onValidate={postRecordValidate}
                onSave={postRecordSave}
                onFixHover={postRecordFixHover}
                onShowStep={postRecordShowStep}
                onFixHint={postRecordFixHint}
                onDismissHint={dismissPostRecordHint}
                onClose={dismissPostRecord}
              />
            {/if}
            {#if activeFeatureTab?.dirty}
              <div class="dirty-banner">
                <span>Текст сценария изменён — сохраните перед запуском</span>
                <button class="primary" on:click={saveFeature}>Сохранить</button>
              </div>
            {/if}
            {#if stepStatusError}
              <div class="dirty-banner error">
                <span>В тексте сценария есть ошибки синтаксиса</span>
              </div>
            {/if}
            {#if showRecordingBar}
              <div class="recording-bar">
                <span class="rec-label">Запись:</span>
                <label class="check-inline"><input type="checkbox" bind:checked={filterRecording} on:change={() => filterRecording && (navOnlyRecording = false)} /> Только важные</label>
                <label class="check-inline"><input type="checkbox" bind:checked={navOnlyRecording} on:change={() => navOnlyRecording && (filterRecording = false)} /> Только ссылки</label>
                <label class="check-inline"><input type="checkbox" bind:checked={settingsHeadless} /> Без окна браузера</label>
                <label class="check-inline"><input type="checkbox" bind:checked={hoverRecord} /> Записывать наведение</label>
              </div>
            {/if}
            <div class="gherkin-hints">
              <span>{stepCount} шагов в редакторе</span>
              <button on:click={openStepsDialog}>Шаблон</button>
              <button on:click={openStepsHelp}>Справка</button>
              <button class:active={previewVisible} on:click={togglePreview}>Превью</button>
            </div>
            <div class="editor-row" style="--preview-width: {previewWidth}px">
              <div class="editor-main">
                <div class="editor-area">
                  <MonacoEditor bind:this={monaco} bind:value={editorText} on:change={(e) => onEditorChange(e.detail)} />
                </div>
                {#if stepsPanelVisible}
                <div class="splitter-h" role="separator" on:mousedown={startResizeSteps}></div>
                <div class="steps-panel" class:collapsed={stepsPanelCollapsed} style="max-height: {stepsPanelCollapsed ? 24 : stepsPanelHeight}px">
                  <div class="steps-header">
                    <button on:click={() => (stepsPanelCollapsed = !stepsPanelCollapsed)}>
                      {#if stepsPanelCollapsed}{@html icons.chevronRight}{:else}{@html icons.chevronDown}{/if}
                    </button>
                    <span>Шаги ({stepCount})</span>
                    {#if stepStatusError}<span style="color:var(--color-error);margin-left:auto">ошибки</span>{/if}
                  </div>
                  {#if !stepsPanelCollapsed}
                    <div class="steps-table-wrap">
                      <table class="steps-table">
                        <thead>
                          <tr><th>#</th><th>Действие</th><th>Элемент</th><th>Значение</th></tr>
                        </thead>
                        <tbody>
                          {#each editorSteps as step, i}
                            <tr
                              class:error={!!step.error}
                              class:clickable={!!step.line}
                              on:click={() => step.line && gotoEditorLine(step.line)}
                              on:dblclick={() => openStepHelpFromPanel(step)}
                              title={step.text ? 'Двойной клик — справка по шагу' : ''}
                            >
                              <td>{i + 1}</td>
                              <td>{step.action}{step.error ? ' ⚠' : ''}</td>
                              <td>{step.element || '—'}</td>
                              <td>{step.value || step.error || '—'}</td>
                            </tr>
                          {/each}
                        </tbody>
                      </table>
                    </div>
                  {/if}
                </div>
                {/if}
              </div>
              {#if previewVisible}
                <div class="splitter-v" role="separator" on:mousedown={startResizePreview}></div>
                <div class="feature-preview-pane" style="width: {previewWidth}px">
                  <div class="preview-header">Превью Gherkin</div>
                  <FeaturePreview text={editorText} />
                </div>
              {/if}
            </div>
          {/if}
        </div>
      </section>
    </div>
  </div>

  <!-- Bottom panel -->
  <div class="bottom-dock" class:collapsed={!bottomPanelOpen}>
    <div class="splitter-h bottom-splitter" role="separator" on:mousedown={startResizeBottom}></div>
    <div class="bottom-panel" style="--panel-height: {bottomPanelHeight}px">
    <div class="panel-tabs">
      <button class="panel-tab" class:active={bottomTab === 'journal'} on:click={() => (bottomTab = 'journal')}>Журнал</button>
      <button class="panel-tab" class:active={bottomTab === 'results'} on:click={() => (bottomTab = 'results')}>Результаты</button>
      <button class="panel-tab" class:active={bottomTab === 'validate'} on:click={() => (bottomTab = 'validate')}>Проверка</button>
      <button class="panel-tab" class:active={bottomTab === 'error'} on:click={() => (bottomTab = 'error')}>Ошибка</button>
    </div>
    <div class="panel-body" class:muted={bottomTab === 'journal' && !logText} class:text-panel={bottomTab === 'journal'}>
      {#if bottomTab === 'journal'}
        {logText || 'Журнал…'}
      {:else if bottomTab === 'results'}
        <ResultsPanel entries={runResults} artifacts={projectArtifacts} onRerun={rerunFailed} onOpenFolder={openArtifactPath} />
      {:else if bottomTab === 'validate'}
        <ValidatePanel
          issues={editorValidationIssues}
          hint={validateProjectHint()}
          onGotoLine={gotoEditorLine}
        />
      {:else}
        <ErrorPanel entry={lastErrorEntry} />
      {/if}
    </div>
  </div>
  </div>

  <!-- Status bar -->
  <footer class="status-bar">
    <div class="status-message" class:error={statusTone === 'error'} class:success={statusTone === 'success'} class:busy={statusTone === 'busy'}>
      {statusMessage}
    </div>
    <div class="status-right">
      <div class="status-segment browser-segment">
        <span class="led" class:recording={recording} class:on={recording || playing}></span>
        Браузер · {recording || playing ? 'открыт' : 'закрыт'}
      </div>
      <div class="status-segment muted">Runner · Playwright</div>
      <div class="status-segment" class:warning={stepStatusError}>{stepStatus}</div>
      <button type="button" class="status-segment clickable" on:click={() => { bottomPanelOpen = true; bottomTab = 'journal' }}>
        {@html icons.log} Журнал
      </button>
      <button type="button" class="status-segment clickable project-segment" on:click={() => (isWelcome ? selectTab(WELCOME_KEY) : openProjectDialog())}>
        {@html icons.explorer} {projectLabel()}
      </button>
    </div>
  </footer>
</div>

{#if showSplash}
  <SplashScreen
    {version}
    message={splashMessage}
    progress={splashProgress}
    fading={splashFading}
  />
{/if}

{#if showRun}
  <RunDialog
    title={runDialogTitle}
    bind:form={runForm}
    {testClients}
    {tags}
    onConfirm={confirmRun}
    onCancel={() => (showRun = false)}
  />
{/if}

{#if showVanessaRun}
  <VanessaRunDialog
    dryRun={vanessaDry}
    preferRerun={vanessaPreferRerun}
    bind:tag={vanessaTag}
    bind:excludeTags={vanessaExcludeTags}
    bind:scenario={vanessaScenario}
    bind:rerunFailedRunDir={vanessaRerunDir}
    bind:installEpf={vanessaInstallEpf}
    bind:epfUrl={vanessaEpfUrl}
    bind:epfDest={vanessaEpfDest}
    {tags}
    onConfirm={confirmVanessaRun}
    onCancel={() => (showVanessaRun = false)}
  />
{/if}

{#if showTestClient}
  <TestClientDialog
    {testClients}
    bind:selectedName={testClientSelection}
    onUse={useTestClient}
    onClose={() => (showTestClient = false)}
    onClientsChange={(names) => (testClients = names)}
    onLog={appendLog}
  />
{/if}

{#if showStepsHelp}
  <StepsHelpDialog
    initialQuery={stepsHelpQuery}
    onClose={() => { showStepsHelp = false; stepsHelpQuery = '' }}
    onInsert={insertStep}
  />
{/if}

{#if showSteps}
  <StepsInsertDialog onInsert={insertStep} onClose={() => (showSteps = false)} />
{/if}

{#if showVanessaSettings}
  <VanessaSettingsDialog onClose={() => (showVanessaSettings = false)} onLog={appendLog} />
{/if}

{#if showExport}
  <ExportDialog
    inputPath={exportInputPath}
    featureText={editorText}
    onClose={() => (showExport = false)}
    onLog={appendLog}
  />
{/if}

{#if showImport}
  <ImportJSONDialog
    {projectPath}
    onClose={() => (showImport = false)}
    onLog={appendLog}
    onImported={onImportComplete}
  />
{/if}

{#if showSettings}
  <SettingsDialog
    bind:browser={settingsBrowser}
    bind:headless={settingsHeadless}
    bind:workers={settingsWorkers}
    bind:loops={settingsLoops}
    bind:filterRecording
    bind:navOnlyRecording
    bind:hoverRecord
    bind:toolbarCompact
    bind:stepsPanelVisible
    bind:stepsPanelHeight
    bind:checkUpdatesOnStartup={settingsCheckUpdatesOnStartup}
    onSave={applySettings}
    onCancel={() => (showSettings = false)}
  />
{/if}

{#if showCommandPalette}
  <CommandPalette commands={paletteCommands} onClose={() => (showCommandPalette = false)} />
{/if}

{#if showSnippetPalette}
  <SnippetPalette onClose={() => (showSnippetPalette = false)} onInsert={insertStep} />
{/if}

{#if showRecord}
  <RecordDialog
    bind:url={recordURL}
    bind:output={recordOutput}
    bind:testClient={recordTestClient}
    bind:idleSeconds={recordIdle}
    bind:appendTo={recordAppendTo}
    bind:headless={settingsHeadless}
    bind:filterRecording
    bind:navOnlyRecording
    bind:hoverRecord
    {testClients}
    {recording}
    {recordPaused}
    onHttpAuth={openHttpAuthDialog}
    onStart={startRecord}
    onTogglePause={toggleRecordPause}
    onStop={stopRecord}
    onClose={() => (showRecord = false)}
  />
{/if}

{#if showOtp}
  <OtpDialog email={otpEmail} onSubmit={submitOtp} onCancel={cancelOtp} />
{/if}

{#if showAbout}
  <AboutDialog {version} onClose={() => (showAbout = false)} />
{/if}

{#if pendingCloseTab}
  <UnsavedCloseDialog
    fileName={basename(pendingCloseTab)}
    onSave={saveAndCloseTab}
    onDiscard={discardAndCloseTab}
    onCancel={cancelCloseTab}
  />
{/if}

{#if showHotkeys}
  <HotkeysDialog commands={paletteCommands} onClose={() => (showHotkeys = false)} />
{/if}

{#if showPlugins}
  <PluginsDialog
    onClose={() => {
      showPlugins = false
      void refreshInstalledPlugins()
    }}
    onRunPlugin={(name, dry) => runPlugin(name, dry)}
  />
{/if}

{#if showRunHistory}
  <RunHistoryDialog
    entries={runResults}
    onOpenFeature={openFeatureFromHistory}
    onRerunFailed={() => { showRunHistory = false; rerunFailed() }}
    onClose={() => (showRunHistory = false)}
  />
{/if}

{#if showFindReplace}
  <FindReplaceDialog
    bind:findText
    bind:replaceText
    bind:caseSensitive={replaceCaseSensitive}
    onFindNext={() => monaco?.findNext(findText, replaceCaseSensitive) ?? false}
    onReplace={() => monaco?.replaceNext(findText, replaceText, replaceCaseSensitive) ?? false}
    onReplaceAll={() => monaco?.replaceAll(findText, replaceText, replaceCaseSensitive) ?? 0}
    onClose={() => (showFindReplace = false)}
  />
{/if}

{#if showProjectReplace}
  <ProjectReplaceDialog
    bind:findText
    bind:replaceText
    bind:caseSensitive={replaceCaseSensitive}
    busy={projectReplaceBusy}
    onConfirm={confirmProjectReplace}
    onClose={() => (showProjectReplace = false)}
  />
{/if}

{#if showHttpAuth}
  <HttpAuthDialog initialHost={httpAuthHost} onClose={() => (showHttpAuth = false)} />
{/if}

{#if showPickerStep}
  <PickerStepDialog
    selector={pickerSelector}
    choices={pickerChoices}
    onInsert={insertPickerStep}
    onClose={() => (showPickerStep = false)}
  />
{/if}

{#if contextMenu}
  <div class="context-menu-backdrop" role="presentation" on:click={dismissContextMenu} on:contextmenu|preventDefault={dismissContextMenu}>
    <div class="context-menu" style="left: {contextMenu.x}px; top: {contextMenu.y}px" role="menu" on:click|stopPropagation>
      <button type="button" on:click={contextMenuRun}>Запустить</button>
      <button type="button" on:click={contextMenuOpen}>Открыть</button>
      <button type="button" on:click={contextMenuDuplicate}>Дублировать</button>
      <button type="button" class="danger" on:click={contextMenuDelete}>Удалить</button>
    </div>
  </div>
{/if}

<BrowserOverlay
  visible={showBrowserOverlay}
  {recording}
  {playing}
  paused={recordPaused}
  onRecord={beginRecord}
  onPause={toggleRecordPause}
  onStop={() => { if (recording) stopRecord(); playing = false }}
  onPicker={pickElement}
  onFocusBrowser={focusBrowserWindow}
/>
