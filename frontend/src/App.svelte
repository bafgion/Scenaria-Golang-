<script lang="ts">
  import { onDestroy, onMount, tick } from 'svelte'
  import MonacoEditor from './lib/MonacoEditor.svelte'
  import WelcomePanel from './lib/WelcomePanel.svelte'
  import CatalogEmptyState from './lib/CatalogEmptyState.svelte'
  import FeatureCatalogTree from './lib/FeatureCatalogTree.svelte'
  import EditorTabBar from './lib/EditorTabBar.svelte'
  import { buildCatalogViewState, buildCatalogStructure, buildCatalogViewStateFromBase, buildRunByPathMap, catalogStructureKey, collectFeaturePathsUnder, type CatalogNode } from './lib/catalogTree'
  import {
    buildBatchSelectedSet,
    selectAllFeaturesUnder,
    toggleBatchPath,
  } from './lib/batchSelection'
  import { debounce, deferToNextFrame } from './lib/uiScheduler'
  import {
    MAX_OPEN_EDITOR_TABS,
    pathsToRetainModels,
    tabEditorText,
    tabNeedsDiskReload,
    trimRetainedTabBodies,
  } from './lib/tabMemory'
  import SettingsDialog from './lib/SettingsDialog.svelte'
  import CommandPalette from './lib/CommandPalette.svelte'
  import type { PaletteCommand } from './lib/paletteTypes'
  import ProjectReplaceDialog from './lib/ProjectReplaceDialog.svelte'
  import HotkeysDialog from './lib/HotkeysDialog.svelte'
  import PluginsDialog from './lib/PluginsDialog.svelte'
  import ExportDialog from './lib/ExportDialog.svelte'
  import RunDialog from './lib/RunDialog.svelte'
  import RecordDialog from './lib/RecordDialog.svelte'
  import PluginRunDialog from './lib/PluginRunDialog.svelte'
  import VanessaRunDialog from './lib/VanessaRunDialog.svelte'
  import VanessaMonitorPanel from './lib/VanessaMonitorPanel.svelte'
  import TestClientDialog from './lib/TestClientDialog.svelte'
  import ImportJSONDialog from './lib/ImportJSONDialog.svelte'
  import AboutDialog from './lib/AboutDialog.svelte'
  import OtpDialog from './lib/OtpDialog.svelte'
  import StepsInsertDialog from './lib/StepsInsertDialog.svelte'
  import VanessaSettingsDialog from './lib/VanessaSettingsDialog.svelte'
  import RefactorUrlDialog from './lib/RefactorUrlDialog.svelte'
  import ConfirmDialog from './lib/ConfirmDialog.svelte'
  import CatalogContextMenu from './lib/CatalogContextMenu.svelte'
  import FolderContextMenu from './lib/FolderContextMenu.svelte'
  import StepsContextMenu from './lib/StepsContextMenu.svelte'
  import OpenProjectDialog from './lib/OpenProjectDialog.svelte'
  import MoveFeatureDialog from './lib/MoveFeatureDialog.svelte'
  import ValidateDialog from './lib/ValidateDialog.svelte'
  import UpdateCheckDialog from './lib/UpdateCheckDialog.svelte'
  import InitProjectDialog from './lib/InitProjectDialog.svelte'
  import ImportFeaturesDialog from './lib/ImportFeaturesDialog.svelte'
  import DuplicateFeatureDialog from './lib/DuplicateFeatureDialog.svelte'
  import RenameFeatureDialog from './lib/RenameFeatureDialog.svelte'
  import { buildFeatureTemplate } from './lib/featureTemplate'
  import { upsertRecordedStepInText } from './lib/recordedStepEditor'
  import { isUntitled, makeUntitledPath, syncUntitledCounterFromPaths, untitledLabel } from './lib/untitled'
  import {
    buildSessionTabsSnapshot,
    sessionTabPathsFromSettings,
    untitledContentMap,
  } from './lib/sessionTabs'
  import { matchHotkey, shouldIgnoreAppHotkey, type HotkeyId } from './lib/hotkeys'
  import { defaultRunForm, type RunForm } from './lib/runTypes'
  import { scenarioAtLine, listScenarioTitles, mergeScenarioNames } from './lib/scenarioAtLine'
  import PostRecordBanner from './lib/PostRecordBanner.svelte'
  import type { HintActionHandlers } from './lib/gherkinHintActions'
  import RunHistoryDialog from './lib/RunHistoryDialog.svelte'
  import StepsHelpDialog from './lib/StepsHelpDialog.svelte'
  import UnsavedCloseDialog from './lib/UnsavedCloseDialog.svelte'
  import HttpAuthDialog from './lib/HttpAuthDialog.svelte'
  import PickerStepDialog from './lib/PickerStepDialog.svelte'
  import { loadLayout, saveLayout, resetLayout as resetUILayout } from './lib/layout'
  import {
    catalogIndentStep,
    clampSidebarWidth,
    clampBottomPanelHeight,
    clampStepsPanelHeight,
    effectivePreviewWidth,
    effectiveSidebarWidth,
    isCompactCatalogTree,
    shouldAutoCompactToolbar,
    shouldShowPreviewPane,
    toolbarIconOnlyThreshold,
    VIEWPORT,
  } from './lib/viewport'
  import ErrorPanel from './lib/ErrorPanel.svelte'
  import ResultsPanel from './lib/ResultsPanel.svelte'
  import ValidatePanel from './lib/ValidatePanel.svelte'
  import FeaturePreview from './lib/FeaturePreview.svelte'
  import FeatureOutline from './lib/FeatureOutline.svelte'
  import SnippetPalette from './lib/SnippetPalette.svelte'
  import BrowserOverlay from './lib/BrowserOverlay.svelte'
  import SplashScreen from './lib/SplashScreen.svelte'
  import { beginSplashWindow, openMainWindow, setSplashDocumentState } from './lib/splashWindow'
  import { preloadMonacoEditor } from './lib/appBootstrap'
  import { setStepHoverEnabled } from './lib/gherkinStepHover'
  import { filterScenarioHints, applyAutoFixableScenarioHints } from './lib/scenarioHints'
  import {
    DEFAULT_EDITOR_SETTINGS,
    editorSettingsFromDTO,
    editorSettingsToDTO,
    type EditorSettings,
  } from './lib/editorOptions'
  import { resolveRecordStartURL } from './lib/recordStartUrl'
  import { loadRecents, rememberFeature, rememberProject } from './lib/recents'
  import { callWailsWithTimeout } from './lib/wailsTimeout'
  import { icons, toolbarIcons } from './lib/icons'
  import { EventsOn, EventsOff, OnFileDrop, OnFileDropOff } from '../wailsjs/runtime/runtime'
  import {
    Version,
    OpenProject,
    ReadFeature,
    SaveFeature,
    WriteTempFeature,
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
    OpenBrowser,
    BeginRecordingCapture,
    RecordBaseline,
    PauseRecording,
    ResumeRecording,
    CancelRecording,
    CloseBrowser,
    StopRecordingCapture,
    FocusBrowser,
    PollBrowserSession,
    UndoRecordedStep,
    UpdateRecordingOptions,
    PickSelector,
    PickerStepChoices,
    LoadSettings,
    SaveSettings,
    SubmitOTPCode,
    CancelOTP,
    CheckUpdateInfo,
    ApplyUpdate,
    DownloadUpdate,
    OpenExternalURL,
    ValidateBrowser,
    ListRunResults,
    BundledExamplesPath,
    ProjectArtifacts,
    ScenariaArtifactPath,
    ParseEditorSteps,
    ArtifactExists,
    OpenFolder,
    ServeAllure,
    OpenHTMLReport,
    RefactorUpdateStartURLs,
    RefactorNormalizeIndents,
    RefactorCollapseBlankLines,
    FormatFeature,
    RefactorReplaceInText,
    ReplaceInProject,
    AnalyzeScenarioHints,
    ApplyScenarioHintFix,
    ResolveRunFromLine,
    ResolveRunToLine,
    SaveFeatureDraft,
    LoadFeatureDraft,
    ClearFeatureDraft,
    DeleteFeature,
    DuplicateFeature,
    MoveFeature,
    RenameFeature,
    ImportFeatures,
    ListPlugins,
    ListVanessaRunDirs,
    ListScenarioTitles,
    StartVanessaRun,
    PollVanessaRun,
  } from '../wailsjs/go/wailsapp/App'
  import { gui } from '../wailsjs/go/models'

  const WELCOME_KEY = '__welcome__'

  type EditorTab = { path: string; content: string; dirty: boolean; draft?: string; unloaded?: boolean }
  type EditorStepRow = gui.EditorStepRow

  let version = ''
  let projectPath = ''
  let features: string[] = []
  let tags: string[] = []
  let projectScenarios: string[] = []
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

  let appReady = false
  let splashMessage = 'Запуск…'
  let splashProgress = 0
  let splashFading = false
  let showSettings = false
  let settingsDialogBaseline: gui.AppSettingsDTO | null = null
  let showCommandPalette = false
  let showSnippetPalette = false
  let showRecord = false
  let recordMode: 'live' | 'baseline' = 'live'
  let baselineBusy = false
  let showPluginRun = false
  let pluginRunName = ''
  let pluginRunDry = false
  let pluginRunTag = ''
  let pluginRunScenario = ''
  let showOtp = false
  let showRun = false
  let showTestClient = false
  let showExport = false
  let showImport = false
  let showImportFeatures = false
  let importDestDir = ''
  let importFeaturesBusy = false
  let showDuplicateFeature = false
  let duplicateFeaturePath = ''
  let duplicateNewName = ''
  let showRefactorUrl = false
  let showOpenProject = false
  let showRenameFeature = false
  let showMoveFeature = false
  let moveFeaturePath = ''
  let moveDestDirs: string[] = []
  let moveDestDir = ''
  let showValidate = false
  let validateBrowser = 'chromium'
  let validateSyntaxOnly = false
  let validateScope: 'project' | 'current' = 'project'
  let validateCliLog = ''
  let showUpdateCheck = false
  let updateCheckMessage = ''
  let updateCheckHasUpdate = false
  let updateCheckInfo: gui.UpdateInfoDTO | null = null
  let updateDownloading = false
  let updateProgress: gui.UpdateProgressDTO | null = null
  let showInitProject = false
  let renameFeaturePath = ''
  let exportInputPath = ''
  let showSteps = false
  let showStepsHelp = false
  let stepsHelpQuery = ''
  let catalogDropTarget = ''
  let showAbout = false
  let showPlugins = false
  let installedPlugins: gui.PluginEntryDTO[] = []
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
  let vanessaPlatformExe = ''
  let vanessaEpfPath = ''
  let vanessaIB = ''
  let vanessaReportAllure = false
  let vanessaVaDir = ''
  let vanessaVaFiles = ''
  let vanessaRunning = false
  let showVanessaMonitor = false
  let vanessaSnapshot: gui.VanessaRunSnapshotDTO = new gui.VanessaRunSnapshotDTO()
  let vanessaWatchDir = ''
  let vanessaPlannedTotal = 1
  let vanessaPollTimer: ReturnType<typeof setInterval> | null = null
  let browserWatchTimer: ReturnType<typeof setInterval> | null = null
  let confirmDialog: {
    title: string
    message: string
    confirmLabel: string
    danger: boolean
    resolve: (value: boolean) => void
  } | null = null

  function askConfirm(opts: {
    title: string
    message: string
    confirmLabel?: string
    danger?: boolean
  }): Promise<boolean> {
    return new Promise((resolve) => {
      confirmDialog = {
        title: opts.title,
        message: opts.message,
        confirmLabel: opts.confirmLabel || 'OK',
        danger: opts.danger || false,
        resolve,
      }
    })
  }

  function closeConfirm(confirmed: boolean) {
    if (confirmDialog) {
      confirmDialog.resolve(confirmed)
      confirmDialog = null
    }
  }
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
  let editorScenarioHints: gui.ScenarioHintDTO[] = []
  let editorHintsDismissed = new Set<string>()
  let contextMenu: { x: number; y: number; path: string } | null = null
  let folderMenu: { x: number; y: number; dir: string; paths: string[] } | null = null
  let runDialogTitle = 'Запуск сценария'
  let runDialogScenarios: string[] = []
  let vanessaDialogScenarios: string[] = []
  let pluginRunScenarios: string[] = []

  let recording = false
  let browserOpen = false
  let recordPaused = false
  let playing = false
  let playingLabel = ''
  let stepsMenu: { x: number; y: number; line: number; step: gui.EditorStepRow } | null = null
  let sessionPersistTimer: ReturnType<typeof setTimeout> | null = null
  let draftAutosaveTimer: ReturnType<typeof setInterval> | null = null

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
  let catalogFilterText = ''
  let openMenu: string | null = null
  let statusMessage = 'Проект → Открыть проект…'
  let statusTone: 'normal' | 'error' | 'success' | 'busy' = 'normal'

  let catalogBaseTree: CatalogNode | null = null
  let catalogBaseTreeKey = ''

  const applyCatalogFilterDebounced = debounce((text: string) => {
    catalogFilterText = text
  }, 200)

  function onSidebarSearchInput(e: Event) {
    const value = (e.currentTarget as HTMLInputElement).value
    sidebarSearch = value
    applyCatalogFilterDebounced(value)
  }

  function setSidebarSearch(value: string) {
    sidebarSearch = value
    applyCatalogFilterDebounced.cancel()
    catalogFilterText = value
  }
  let batchSelected: string[] = []
  let batchMode = false
  let catalogCollapsed = new Set<string>()
  let showBatchHint = true
  let toolbarCompact = false
  let viewportWidth = 1280
  let viewportHeight = 800
  let viewportAutoCompact = false

  let filterRecording = false
  let navOnlyRecording = false
  let hoverRecord = false
  let stepsPanelVisible = true

  let lastRun: RunForm = defaultRunForm({ headed: true, installPW: true, html: true })
  let runForm: RunForm = { ...lastRun }

  let settingsBrowser = 'chromium'
  let settingsHeadless = false
  let settingsWorkers = 1
  let settingsSlowMo = 0
  let settingsScrollBeforeClick = false
  let settingsHoverRecordMinMs = 600
  let settingsLoops = 100
  let settingsCheckUpdatesOnStartup = true
  let settingsSelectorClickStrategies: string[] = ['testid', 'id', 'aria', 'contextual', 'text']
  let settingsSelectorInputStrategies: string[] = ['testid', 'id', 'label', 'placeholder', 'aria', 'name']
  let editorSettings: EditorSettings = { ...DEFAULT_EDITOR_SETTINGS }
  let editorCursorLine = 1
  let stepsPanelTab: 'outline' | 'steps' = DEFAULT_EDITOR_SETTINGS.stepsPanelView

  let recordURL = ''
  let startURL = ''
  let recordOutput = 'recorded.feature'
  let recordIdle = 30
  let recordAppendTo = ''
  let recordTestClient = ''
  let recordFeatureName = 'Записанный сценарий'
  let recordScenarioName = 'Запись'
  let lastRecordTarget = ''
  let liveRecordStepLines: Record<number, number> = {}
  let pendingCloseTab: string | null = null

  let otpEmail = ''

  let testClientSelection = ''
  let testClientSuggestName = ''

  let recentProjects: string[] = []
  let recentFeatures: string[] = []
  let runResults: gui.RunResultEntry[] = []
  let lastErrorEntry: gui.RunResultEntry | null = null
  let editorSteps: EditorStepRow[] = []
  let editorValidationIssues: gui.ValidationIssue[] = []
  let validatePanelIssues: gui.ValidationIssue[] = []
  let projectArtifacts: gui.ProjectArtifacts = new gui.ProjectArtifacts()

  const unsubscribers: (() => void)[] = []

  $: isWelcome = activeTab === WELCOME_KEY
  $: activeFeatureTab = tabs.find((t) => t.path === activeTab)
  $: activeTabUnsaved = activeFeatureTab ? tabIsUnsaved(activeFeatureTab) : false
  $: if (editorSettings.stepsPanelView) {
    stepsPanelTab = editorSettings.stepsPanelView
  }
  $: stepCount = editorSteps.length
  $: batchCount = batchSelected.length
  $: batchSelectedSet = buildBatchSelectedSet(batchSelected)
  $: showRecordingBar = recording && !showRecord
  $: showPlayingBar = playing
  $: showBrowserOverlay = false
  $: paletteCommands = buildPaletteCommands()

  let resizingBottom = false
  let resizingSteps = false
  let resizingSidebar = false
  let resizingPreview = false
  let actionBarEl: HTMLElement | undefined
  let toolbarIconOnly = true

  $: actionBarCompact = toolbarCompact || viewportAutoCompact
  $: layoutSidebarWidth = effectiveSidebarWidth(sidebarWidth, viewportWidth, sidebarVisible)
  $: catalogIndent = catalogIndentStep(layoutSidebarWidth)
  $: compactCatalogTree = isCompactCatalogTree(layoutSidebarWidth)
  $: showPreviewPane = shouldShowPreviewPane(viewportWidth, previewVisible)
  $: layoutPreviewWidth = effectivePreviewWidth(previewWidth, viewportWidth, previewVisible)

  $: activeFeaturePath = activeTab !== WELCOME_KEY ? activeTab : ''

  $: runByPath = buildRunByPathMap(runResults)

  $: tagsByPath = (() => {
    const map = new Map<string, string[]>()
    for (const [path, pathTags] of Object.entries(featureTags)) {
      map.set(path.replace(/\\/g, '/').toLowerCase(), pathTags)
    }
    return map
  })()

  $: {
    const key = projectPath ? catalogStructureKey(projectPath, features) : ''
    if (key !== catalogBaseTreeKey) {
      catalogBaseTreeKey = key
      catalogBaseTree = projectPath ? buildCatalogStructure(projectPath, features) : null
    }
  }

  $: catalogViewState = buildCatalogViewStateFromBase(
    projectPath || null,
    catalogBaseTree,
    catalogFilterText,
    runByPath,
    true,
    tagsByPath,
  )

  $: welcomeProjectOpen = !!projectPath
  $: welcomeRecorded = recording || browserOpen || editorSteps.length > 0 || tabs.length > 0
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
    setSplashDocumentState(false)
    await openMainWindow()
    appReady = true
  }

  onMount(async () => {
    const startedAt = Date.now()
    let splashFinished = false
    const finishStartup = async () => {
      if (splashFinished) return
      splashFinished = true
      setSplashStage('Готово', 100)
      await dismissSplash(startedAt)
      applyDevUiMock()
      syncIdleStatus()
      void checkUpdatesOnStartup()
    }
    const startupGuard = window.setTimeout(() => {
      console.error('Startup guard: forcing splash dismiss')
      void finishStartup()
    }, 12_000)

    try {
    setSplashDocumentState(true)
    await beginSplashWindow()

    setSplashStage('Настройка окружения…', 8)

    const layout = loadLayout()
    sidebarVisible = layout.sidebarVisible
    bottomPanelOpen = layout.bottomPanelOpen
    bottomPanelHeight = layout.bottomPanelHeight
    sidebarWidth = clampSidebarWidth(layout.sidebarWidth || 260)
    previewVisible = layout.previewVisible
    previewWidth = layout.previewWidth || 360

    setSplashStage('Загрузка редактора…', 28)
    setStepHoverEnabled(() => editorSettings.stepHover)
    try {
      await preloadMonacoEditor()
    } catch (err) {
      console.error('Monaco preload failed', err)
    }

    setSplashStage('Подключение к приложению…', 48)

    try {
      version = await Version()
    } catch {
      version = 'dev'
    }

    setSplashStage('Загрузка настроек…', 68)

    const [recents, settings] = await Promise.all([
      loadRecents(),
      callWailsWithTimeout('LoadSettings', LoadSettings(), 4000),
    ])
    recentProjects = recents.projects
    recentFeatures = recents.features
    if (settings) {
      applySettingsFromDTO(settings)
      setStepHoverEnabled(() => editorSettings.stepHover)
      stepsPanelCollapsed = resolveStepsPanelCollapsed()
      stepsPanelHeight = settings.stepsPanelHeight || 160
      if (settings.sidebarWidth >= VIEWPORT.sidebarMin) {
        sidebarWidth = clampSidebarWidth(settings.sidebarWidth)
      }
      if (settings.recentProjects?.length) recentProjects = settings.recentProjects
      if (settings.recentFeatures?.length) recentFeatures = settings.recentFeatures
      await restoreWorkspaceSession(settings)
    }

    draftAutosaveTimer = setInterval(() => void autosaveDirtyDrafts(), 30_000)

    setSplashStage('Инициализация…', 88)

    try {
      OnFileDrop((_x, _y, paths) => {
        if (projectPath && paths?.length) {
          const dest = catalogDropTarget || projectPath
          catalogDropTarget = ''
          void importDroppedFeatures(dest, paths)
        }
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
        EventsOn('browser-opened', () => {
          applyBrowserSessionState({ browserOpen: true, recording: false, paused: false })
          showRecord = false
          setStatus('Браузер открыт', 'busy')
          appendLog('Браузер открыт. Запись шагов — только по кнопке «Запись» (в браузере или Ctrl+R).')
          startBrowserWatch()
        }),
      )
      unsubscribers.push(
        EventsOn('browser-closed', (result: gui.RunResult) => {
          stopBrowserWatch()
          handleRecordSessionEnd(result, 'browse')
        }),
      )
      unsubscribers.push(
        EventsOn('browser-lost', () => {
          if (browserOpen || recording) {
            handleBrowserLost()
          }
        }),
      )
      unsubscribers.push(
        EventsOn('toolbar-picker', () => {
          void pickElement()
        }),
      )
      unsubscribers.push(
        EventsOn('record-started', async (payload: string | { append?: boolean; sync?: boolean; output?: string }) => {
          const meta = typeof payload === 'object' && payload !== null ? payload : { append: false, output: payload }
          const appendOnly = meta.append === true
          const syncOnly = meta.sync === true
          applyBrowserSessionState({ browserOpen: true, recording: true, paused: false })
          showRecord = false
          setStatus('● Идёт запись', 'busy')
          startBrowserWatch()
          if (!syncOnly) {
            liveRecordStepLines = {}
            if (!appendOnly) {
              await prepareRecordEditorTab(meta.output || '')
            }
          }
          if (!appendOnly && !syncOnly) {
            appendLog('Запись начата…')
          }
        }),
      )
      unsubscribers.push(
        EventsOn('record-stopped', () => {
          handleRecordStopped()
        }),
      )
      unsubscribers.push(
        EventsOn('record-step', (payload: { index: number; line: string }) => {
          void applyLiveRecordedStep(payload?.index ?? 0, payload?.line ?? '')
        }),
      )
      unsubscribers.push(
        EventsOn('record-finished', async (result: gui.RunResult) => {
          stopBrowserWatch()
          handleRecordSessionEnd(result, 'record')
        }),
      )
      unsubscribers.push(
        EventsOn('vanessa-run-started', () => {
          vanessaRunning = true
          vanessaWatchDir = ''
          setStatus('▶ Vanessa', 'busy')
          appendLog('Запуск Vanessa…')
          startVanessaPoll()
        }),
      )
      unsubscribers.push(
        EventsOn('vanessa-run-finished', async (result: gui.VanessaRunResultDTO) => {
          stopVanessaPoll()
          vanessaRunning = false
          if (result.runDir) {
            vanessaWatchDir = result.runDir
            try {
              vanessaSnapshot = await PollVanessaRun(result.runDir, vanessaPlannedTotal)
            } catch {
              /* ignore */
            }
          }
          if (result.output) appendLog(result.output.trimEnd())
          if (result.error) {
            appendLog(`Ошибка: ${result.error}`)
            setStatus('Ошибка Vanessa', 'error')
          } else {
            appendLog('Vanessa завершён.')
            setStatus('Vanessa завершён', result.success ? 'success' : 'error')
          }
          await refreshRunResults()
        }),
      )
    } catch {
      /* dev without wails runtime */
    }

    const onResize = () => syncViewportLayout()
    syncViewportLayout()
    window.addEventListener('resize', onResize)
    unsubscribers.push(() => window.removeEventListener('resize', onResize))

    const onDocClick = () => {
      openMenu = null
    }
    window.addEventListener('keydown', onGlobalKeydown, { capture: true })
    window.addEventListener('click', onDocClick)
    unsubscribers.push(() => window.removeEventListener('keydown', onGlobalKeydown, { capture: true }))
    unsubscribers.push(() => window.removeEventListener('click', onDocClick))
    } catch (err) {
      console.error('Startup failed', err)
    } finally {
      window.clearTimeout(startupGuard)
      await finishStartup()
    }
  })

  onDestroy(() => {
    stopVanessaPoll()
    stopBrowserWatch()
    applyCatalogFilterDebounced.cancel()
    if (sessionPersistTimer) {
      clearTimeout(sessionPersistTimer)
      sessionPersistTimer = null
      void persistSettings()
    }
    if (validateDebounceTimer) clearTimeout(validateDebounceTimer)
    if (draftAutosaveTimer) clearInterval(draftAutosaveTimer)
    for (const off of unsubscribers) off()
  })

  function schedulePersistSession() {
    if (sessionPersistTimer) clearTimeout(sessionPersistTimer)
    sessionPersistTimer = setTimeout(() => void persistSettings(), 500)
  }

  async function restoreWorkspaceSession(s: gui.AppSettingsDTO) {
    const proj = (s.sessionProject || '').trim()
    if (!proj) return
    try {
      await openProjectAt(proj)
      const untitledBodies = untitledContentMap(s.untitledTabs)
      syncUntitledCounterFromPaths([
        ...(s.openTabs || []),
        ...untitledBodies.keys(),
      ])
      const tabPaths = sessionTabPathsFromSettings(s.openTabs, s.untitledTabs)
      for (const p of tabPaths) {
        if (isUntitled(p)) {
          const content = untitledBodies.get(p)
          if (content === undefined || tabs.some((t) => t.path === p)) continue
          tabs = [...tabs, { path: p, content, dirty: true }]
          continue
        }
        try {
          await loadFeature(p)
        } catch {
          /* skip missing files */
        }
      }
      const active = (s.activeTab || '').trim()
      if (active) {
        welcomeTabVisible = false
        if (isUntitled(active)) {
          const tab = tabs.find((t) => t.path === active)
          if (tab) {
            await applyEditorText(tabEditorText(tab), {
              switchTab: true,
              tabPath: active,
              skipValidate: true,
            })
            activeTab = active
            trimTabsMemory()
          }
        } else {
          await loadFeature(active)
        }
      } else if (tabPaths.length > 0) {
        welcomeTabVisible = false
      }
    } catch {
      /* ignore broken session */
    }
  }

  async function autosaveDirtyDrafts() {
    if (!projectPath) return
    syncActiveTabContent()
    for (const tab of tabs) {
      if (!tab.dirty || isUntitled(tab.path)) continue
      const text = tab.path === activeTab ? editorText : tabEditorText(tab)
      try {
        await SaveFeatureDraft(tab.path, text)
      } catch {
        /* offline */
      }
    }
  }

  function partialRunLogSuffix(startStep: number, endStep: number): string {
    if (startStep >= 0 && endStep >= 0) return ` (шаги ${startStep + 1}–${endStep + 1})`
    if (startStep >= 0) return ` (с шага ${startStep + 1})`
    if (endStep >= 0) return ` (до шага ${endStep + 1})`
    return ''
  }

  function openStepsContextMenu(e: MouseEvent, step: gui.EditorStepRow) {
    if (!step.line) return
    e.preventDefault()
    stepsMenu = { x: e.clientX, y: e.clientY, line: step.line, step }
  }

  function closeStepsMenu() {
    stepsMenu = null
  }

  function stepsMenuRunFrom(dryRun: boolean) {
    const menu = stepsMenu
    if (!menu) return
    const line = menu.line
    closeStepsMenu()
    void runScenarioAtLine(line, dryRun, '', true)
  }

  function stepsMenuRunTo(dryRun: boolean) {
    const menu = stepsMenu
    if (!menu) return
    const line = menu.line
    closeStepsMenu()
    void runScenarioToLine(line, dryRun)
  }

  function stepsMenuGoto() {
    const menu = stepsMenu
    if (!menu) return
    gotoEditorLine(menu.line)
    closeStepsMenu()
  }

  function stepsMenuHelp() {
    const menu = stepsMenu
    if (!menu) return
    openStepHelpFromPanel(menu.step)
    closeStepsMenu()
  }

  function basename(path: string): string {
    const parts = path.replace(/\\/g, '/').split('/')
    return parts[parts.length - 1] || path
  }

  function featureTabLabel(path: string): string {
    if (isUntitled(path)) return untitledLabel(path)
    return basename(path)
  }

  function tabIsUnsaved(tab: EditorTab): boolean {
    return tab.dirty || isUntitled(tab.path)
  }

  function applySettingsFromDTO(s: gui.AppSettingsDTO) {
    settingsBrowser = s.browser || 'chromium'
    settingsHeadless = s.headless
    settingsWorkers = s.parallelWorkers || 1
    settingsSlowMo = s.slowMo ?? 0
    settingsScrollBeforeClick = s.scrollBeforeClick ?? false
    settingsHoverRecordMinMs = s.hoverRecordMinMs || 600
    settingsLoops = s.maxLoopIterations || 100
    filterRecording = s.filterRecording
    navOnlyRecording = s.navOnlyRecording
    hoverRecord = s.hoverRecord
    toolbarCompact = s.toolbarCompact
    stepsPanelVisible = s.stepsPanelVisible !== false
    stepsPanelHeight = s.stepsPanelHeight || 160
    settingsCheckUpdatesOnStartup = s.checkUpdatesOnStartup !== false
    settingsSelectorClickStrategies = s.selectorClickStrategies?.length
      ? [...s.selectorClickStrategies]
      : ['testid', 'id', 'aria', 'contextual', 'text']
    settingsSelectorInputStrategies = s.selectorInputStrategies?.length
      ? [...s.selectorInputStrategies]
      : ['testid', 'id', 'label', 'placeholder', 'aria', 'name']
    editorSettings = editorSettingsFromDTO(s.editor)
    stepsPanelTab = editorSettings.stepsPanelView
    stepsPanelCollapsed = resolveStepsPanelCollapsed()
    lastRun = {
      ...lastRun,
      workers: s.parallelWorkers || 1,
      slowMo: s.slowMo ?? 0,
      browser: s.browser || 'chromium',
    }
  }

  function resolveStepsPanelCollapsed(): boolean {
    if (!stepsPanelVisible) return true
    return editorSettings.symbolOutline && editorSettings.stepsPanelView === 'outline'
  }

  function buildPaletteCommands(): PaletteCommand[] {
    const compactLabel = toolbarCompact ? 'Расширенная панель' : 'Компактная панель'
    return [
      { id: 'palette', label: 'Палитра команд', group: 'Вид', shortcut: 'Ctrl+Shift+P', run: () => (showCommandPalette = true) },
      { id: 'welcome', label: 'Старт', group: 'Вид', run: () => selectTab(WELCOME_KEY) },
      { id: 'open', label: 'Открыть проект…', group: 'Проект', run: openProjectDialog },
      { id: 'close-project', label: 'Закрыть проект', group: 'Проект', run: closeProject },
      { id: 'settings', label: 'Настройки…', group: 'Проект', shortcut: 'Ctrl+,', run: openSettings },
      { id: 'init', label: 'Init проекта', group: 'Проект', run: openInitProjectDialog },
      { id: 'examples', label: 'Открыть примеры сценариев', group: 'Проект', run: openExamples },
      { id: 'new', label: 'Новый сценарий', group: 'Сценарий', shortcut: 'Ctrl+N', run: newScenario },
      { id: 'open-file', label: 'Открыть файл…', group: 'Сценарий', shortcut: 'Ctrl+O', run: openFileDialog },
      { id: 'save', label: 'Сохранить', group: 'Сценарий', shortcut: 'Ctrl+S', run: saveFeature },
      { id: 'save-as', label: 'Сохранить как…', group: 'Сценарий', shortcut: 'Ctrl+Shift+S', run: saveFeatureAs },
      { id: 'export', label: 'Экспорт…', group: 'Сценарий', run: openExportDialog },
      { id: 'import', label: 'Импорт JSON…', group: 'Сценарий', run: openImportDialog },
      { id: 'import-features', label: 'Импорт .feature…', group: 'Сценарий', run: openImportFeaturesDialog },
      { id: 'steps', label: 'Вставить шаг…', group: 'Сценарий', run: openStepsDialog },
      { id: 'snippets', label: 'Палитра сниппетов', group: 'Сценарий', shortcut: 'Ctrl+Shift+Space', run: openSnippetPalette },
      { id: 'find-replace', label: 'Найти и заменить…', group: 'Сценарий', shortcut: 'Ctrl+H', run: openFindReplace },
      { id: 'format', label: 'Форматировать сценарий', group: 'Сценарий', shortcut: 'Shift+Alt+F', run: () => void monaco?.formatDocument() },
      { id: 'goto-symbol', label: 'Перейти к символу…', group: 'Сценарий', shortcut: 'Ctrl+Shift+O', run: () => monaco?.openSymbolOutline() },
      { id: 'project-replace', label: 'Замена по проекту…', group: 'Сценарий', run: () => (showProjectReplace = true) },
      { id: 'duplicate', label: 'Дублировать сценарий…', group: 'Сценарий', run: () => activeTab && !isWelcome && openDuplicateDialog(activeTab) },
      { id: 'rename-feature', label: 'Переименовать сценарий…', group: 'Сценарий', run: () => {
        if (!activeTab || isWelcome) return
        renameFeaturePath = activeTab
        showRenameFeature = true
      }},
      { id: 'delete-feature', label: 'Удалить сценарий', group: 'Сценарий', run: () => activeTab && !isWelcome && deleteFeature(activeTab) },
      { id: 'refactor-indents', label: 'Нормализовать отступы', group: 'Рефакторинг', run: refactorNormalizeIndents },
      { id: 'refactor-blanks', label: 'Убрать пустые строки между шагами', group: 'Рефакторинг', run: refactorCollapseBlank },
      { id: 'steps-help', label: 'Справка по шагам…', group: 'Справка', shortcut: 'F1', run: () => openStepsHelp() },
      { id: 'browser', label: 'Браузер', group: 'Запись и тест', shortcut: 'Ctrl+B', run: () => void openBrowser() },
      { id: 'record', label: 'Запись', group: 'Запись и тест', shortcut: 'Ctrl+R', run: beginRecord },
      { id: 'record-baseline', label: 'Запись из шагов…', group: 'Запись и тест', run: openBaselineRecordDialog },
      { id: 'stop', label: 'Стоп', group: 'Запись и тест', run: stopRecord },
      { id: 'pause', label: 'Пауза', group: 'Запись и тест', run: toggleRecordPause },
      { id: 'run', label: 'Запустить', group: 'Запись и тест', shortcut: 'Ctrl+Enter', run: () => runPrimary(false) },
      { id: 'run-current', label: 'Запустить текущий сценарий', group: 'Запись и тест', shortcut: 'Ctrl+Shift+Enter', run: () => runCurrentScenario(false) },
      { id: 'run-current-dry', label: 'Dry-run текущего сценария', group: 'Запись и тест', run: () => runCurrentScenario(true) },
      { id: 'run-dialog', label: 'Запустить…', group: 'Запись и тест', run: () => openRunDialog('Запуск сценария', {}) },
      { id: 'run-tag', label: 'Запустить сценарии с тегом…', group: 'Запись и тест', run: () => openRunDialog('Запуск с тегом', {}) },
      { id: 'playwright', label: 'Playwright…', group: 'Запись и тест', run: () => openRunDialog('Playwright', { dryRun: false, headed: true, engine: 'playwright', installPW: true }) },
      { id: 'dry', label: 'Dry-run', group: 'Запись и тест', run: () => runPrimary(true) },
      { id: 'batch', label: 'Пакетный выбор', group: 'Запись и тест', run: () => toggleBatchMode() },
      { id: 'batch-run', label: 'Запустить выбранные', group: 'Запись и тест', run: () => runBatchSelected(false) },
      { id: 'batch-dry', label: 'Dry-run выбранных', group: 'Запись и тест', run: () => runBatchSelected(true) },
      { id: 'rerun-failed', label: 'Перезапустить упавшие', group: 'Запись и тест', run: rerunFailed },
      { id: 'run-history', label: 'История запусков…', group: 'Запись и тест', run: openRunHistory },
      { id: 'testclient', label: 'TestClient…', group: 'Запись и тест', run: openTestClientDialog },
      { id: 'capture-session', label: 'Сохранить сессию браузера…', group: 'Запись и тест', run: openTestClientDialogForCapture },
      { id: 'validate', label: 'Проверить', group: 'Запись и тест', run: () => openValidateDialog(true) },
      { id: 'validate-browser', label: 'Проверить в браузере', group: 'Запись и тест', run: () => openValidateDialog(false) },
      ...(hasVanessaPlugin()
        ? [
            { id: 'vanessa-dry', label: 'Vanessa (dry)…', group: 'Запись и тест', run: () => openVanessaDialog(true) },
            { id: 'vanessa', label: 'Vanessa run…', group: 'Запись и тест', run: () => openVanessaDialog(false) },
            { id: 'vanessa-rerun', label: 'Vanessa rerun-failed…', group: 'Запись и тест', run: () => openVanessaDialog(false, true) },
            { id: 'vanessa-settings', label: 'Настройки Vanessa…', group: 'Запись и тест', run: openVanessaSettingsDialog },
            { id: 'vanessa-monitor', label: 'Монитор Vanessa…', group: 'Запись и тест', run: openVanessaMonitor },
          ]
        : []),
      { id: 'plugins', label: 'Управление плагинами…', group: 'Плагины', run: () => (showPlugins = true) },
      ...installedPlugins.flatMap((plugin) => {
        if (!plugin.runnable || plugin.vanessa) return []
        const label = pluginLabel(plugin)
        return [
          { id: `plugin-${plugin.name}-dry`, label: `${label} (dry)…`, group: 'Плагины', run: () => openPluginRun(plugin.name, true) },
          { id: `plugin-${plugin.name}`, label: `Запустить ${label}…`, group: 'Плагины', run: () => openPluginRun(plugin.name, false) },
        ]
      }),
      { id: 'journal', label: 'Журнал', group: 'Вид', shortcut: 'Ctrl+`', run: () => { bottomPanelOpen = true; bottomTab = 'journal' } },
      { id: 'results', label: 'Результаты', group: 'Вид', run: () => { bottomPanelOpen = true; bottomTab = 'results' } },
      { id: 'allure-serve', label: 'Allure serve', group: 'Вид', run: () => serveAllureReport(projectArtifacts.allureDir || '') },
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
    stepsPanelCollapsed = resolveStepsPanelCollapsed()
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
    monaco?.openFindReplace()
  }

  function resetWindowLayout() {
    const layout = resetUILayout()
    sidebarVisible = layout.sidebarVisible
    bottomPanelOpen = layout.bottomPanelOpen
    bottomPanelHeight = layout.bottomPanelHeight
    sidebarWidth = clampSidebarWidth(layout.sidebarWidth)
    previewVisible = layout.previewVisible
    previewWidth = layout.previewWidth
    appendLog('Макет окон сброшен')
  }

  function hintDismissKey(hint: gui.ScenarioHintDTO): string {
    return `${hint.id}:${hint.stepIndex}`
  }

  const monacoHintActions: HintActionHandlers = {
    getHints: () => editorScenarioHints,
    onFix: (hint) => applyEditorHintFix(hint),
    onDismiss: (hint) => dismissEditorHint(hint),
  }

  const monacoRunLensActions = {
    isEnabled: () => editorSettings.codeLens && !!activeTab && !isWelcome,
    onRun: (payload: { scenario: string; line: number; dryRun: boolean; partial: boolean }) =>
      runScenarioAtLine(payload.line, payload.dryRun, payload.scenario, payload.partial),
  }

  const monacoInlayHintsHandlers = {
    isEnabled: () => editorSettings.inlayHints && !!activeTab && !isWelcome,
    getSteps: () => editorSteps,
  }

  async function refreshEditorScenarioHints() {
    if (isWelcome || !editorSettings.scenarioHints) {
      editorScenarioHints = []
      return
    }
    try {
      const all = await AnalyzeScenarioHints(editorText)
      editorScenarioHints = all
        .filter((h) => !editorHintsDismissed.has(hintDismissKey(h)))
        .filter((h) => filterScenarioHints([h], editorSettings).length > 0)
    } catch {
      editorScenarioHints = []
    }
  }

  async function runScenarioHintsAutoFix(text: string): Promise<string> {
    if (!editorSettings.scenarioHints || !editorSettings.scenarioHintsAutoFixOnSave) {
      return text
    }
    const all = await AnalyzeScenarioHints(text)
    const fixable = filterScenarioHints(all, editorSettings).filter((h) => h.autoFixable)
    return applyAutoFixableScenarioHints(text, fixable, async (hint, currentText) => {
      const result = await ApplyScenarioHintFix({
        text: currentText,
        hintId: hint.id,
        stepIndex: hint.stepIndex,
      })
      return result.count > 0 ? result.text : null
    })
  }

  async function applyEditorHintFix(hint: gui.ScenarioHintDTO) {
    const result = await ApplyScenarioHintFix({
      text: editorText,
      hintId: hint.id,
      stepIndex: hint.stepIndex,
    })
    if (result.count > 0) {
      editorText = result.text
      syncActiveTabContent()
      monaco?.setContent(editorText)
      appendLog(`Исправлена подсказка: ${hint.title}`)
      await validateEditor()
    }
  }

  function dismissEditorHint(hint: gui.ScenarioHintDTO) {
    editorHintsDismissed.add(hintDismissKey(hint))
    void refreshEditorScenarioHints()
  }

  async function showPostRecordBanner(path: string) {
    postRecordPath = path
    editorHintsDismissed = new Set()
    try {
      const content = await ReadFeature(path)
      postRecordStepCount = (await ParseEditorSteps(content)).length
    } catch {
      postRecordStepCount = 0
    }
    if (editorSettings.scenarioHints && editorSettings.scenarioHintsAfterRecord) {
      await refreshEditorScenarioHints()
    } else {
      editorScenarioHints = []
    }
  }

  function dismissPostRecord() {
    postRecordPath = ''
    postRecordStepCount = 0
  }

  async function postRecordValidate() {
    if (!postRecordPath) return
    if (activeTab !== postRecordPath) await loadFeature(postRecordPath)
    await validateProject(false)
    bottomTab = 'validate'
  }

  async function postRecordSave() {
    if (!postRecordPath) return
    if (activeTab !== postRecordPath) await loadFeature(postRecordPath)
    await saveFeature()
    dismissPostRecord()
  }

  function openDuplicateDialog(path: string) {
    if (!path || isWelcome) return
    duplicateFeaturePath = path
    duplicateNewName = `${basename(path).replace(/\.feature$/i, '')}-copy`
    showDuplicateFeature = true
  }

  async function confirmDuplicateFeature(newName: string) {
    showDuplicateFeature = false
    const path = duplicateFeaturePath
    duplicateFeaturePath = ''
    if (!path) return
    try {
      const target = await DuplicateFeature(path, newName)
      await refreshProject()
      await loadFeature(target)
      appendLog(`Создана копия: ${basename(target)}`)
    } catch (e: any) {
      appendLog(`Ошибка: ${e}`)
    }
  }

  async function duplicateFeature(path: string) {
    openDuplicateDialog(path)
  }

  async function deleteFeature(path: string) {
    if (!path || isWelcome) return
    const ok = await askConfirm({
      title: 'Удалить сценарий',
      message: `Удалить «${basename(path)}» без возможности восстановления?`,
      confirmLabel: 'Удалить',
      danger: true,
    })
    if (!ok) return
    await doDeleteFeature(path)
  }

  async function doDeleteFeature(path: string) {
    if (isUntitled(path)) {
      finalizeCloseTab(path)
      appendLog(`Закрыт: ${untitledLabel(path)}`)
      return
    }
    try {
      await DeleteFeature(path)
      closeTab(path)
      await refreshProject()
      appendLog(`Удалён: ${basename(path)}`)
    } catch (e: any) {
      appendLog(`Ошибка: ${e}`)
    }
  }

  async function runFeatureFile(path: string, dryRun = false) {
    if (!path) return
    await loadFeature(path)
    await executeRun({ ...lastRun, dryRun }, [path])
  }

  function onFileContextMenu(e: MouseEvent, path: string) {
    e.preventDefault()
    contextMenu = { x: e.clientX, y: e.clientY, path }
    folderMenu = null
  }

  function onFolderContextMenu(e: MouseEvent, node: CatalogNode) {
    e.preventDefault()
    const paths = collectFeaturePathsUnder(node)
    if (!paths.length) return
    folderMenu = { x: e.clientX, y: e.clientY, dir: node.path, paths }
    contextMenu = null
  }

  function dismissFolderMenu() {
    folderMenu = null
  }

  function folderRelativePath(dirPath: string): string {
    if (!projectPath) return ''
    const root = projectPath.replace(/\\/g, '/').replace(/\/$/, '')
    const dir = dirPath.replace(/\\/g, '/')
    if (dir.toLowerCase() === root.toLowerCase()) return ''
    if (dir.toLowerCase().startsWith(root.toLowerCase() + '/')) {
      return dir.slice(root.length + 1)
    }
    return dir
  }

  function folderMenuRun(dryRun = false) {
    if (!folderMenu) return
    const paths = folderMenu.paths
    dismissFolderMenu()
    void executeRun({ ...lastRun, dryRun }, paths)
  }

  function folderMenuSelectBatch() {
    if (!folderMenu) return
    batchMode = true
    batchSelected = [...folderMenu.paths]
    dismissFolderMenu()
    appendLog(`Выбрано ${batchSelected.length} сценариев в папке`)
  }

  function openVanessaForFolder(dirPath: string, dry: boolean) {
    vanessaDry = dry
    vanessaTag = ''
    vanessaExcludeTags = ''
    vanessaScenario = ''
    vanessaDialogScenarios = dialogScenarioNames()
    vanessaRerunDir = ''
    vanessaPreferRerun = false
    vanessaInstallEpf = false
    vanessaEpfUrl = ''
    vanessaEpfDest = ''
    vanessaPlatformExe = ''
    vanessaEpfPath = ''
    vanessaIB = ''
    vanessaReportAllure = false
    vanessaVaDir = folderRelativePath(dirPath)
    vanessaVaFiles = ''
    showVanessaRun = true
  }

  function folderMenuVanessa(dry: boolean) {
    if (!folderMenu) return
    const dir = folderMenu.dir
    dismissFolderMenu()
    openVanessaForFolder(dir, dry)
  }

  function dismissContextMenu() {
    contextMenu = null
  }

  function contextMenuRun() {
    if (!contextMenu) return
    const path = contextMenu.path
    dismissContextMenu()
    runFeatureFile(path, false)
  }

  function contextMenuDryRun() {
    if (!contextMenu) return
    const path = contextMenu.path
    dismissContextMenu()
    runFeatureFile(path, true)
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

  function dirname(path: string): string {
    const norm = path.replace(/\\/g, '/')
    const i = norm.lastIndexOf('/')
    return i >= 0 ? norm.slice(0, i) : ''
  }

  function collectProjectDirs(): string[] {
    if (!projectPath) return []
    const dirs = new Set<string>([projectPath.replace(/\\/g, '/')])
    for (const feature of features) {
      const dir = dirname(feature.replace(/\\/g, '/'))
      if (dir) dirs.add(dir)
    }
    return [...dirs].sort((a, b) => a.localeCompare(b, 'ru'))
  }

  function contextMenuMove() {
    if (!contextMenu || !projectPath) return
    moveFeaturePath = contextMenu.path
    moveDestDirs = collectProjectDirs().filter((d) => d !== dirname(moveFeaturePath.replace(/\\/g, '/')))
    moveDestDir = moveDestDirs[0] || projectPath.replace(/\\/g, '/')
    dismissContextMenu()
    showMoveFeature = true
  }

  async function confirmMoveFeature(destDir: string) {
    showMoveFeature = false
    const src = moveFeaturePath
    moveFeaturePath = ''
    if (!src || !destDir) return
    try {
      const newPath = await MoveFeature(src, destDir)
      const wasActive = activeTab === src
      await refreshProject()
      if (wasActive) await loadFeature(newPath)
      appendLog(`Перемещён: ${basename(newPath)}`)
    } catch (e: any) {
      appendLog(`Ошибка: ${e}`)
    }
  }

  function contextMenuReveal() {
    if (!contextMenu) return
    const path = contextMenu.path
    dismissContextMenu()
    const dir = path.replace(/[/\\][^/\\]+$/, '')
    if (dir) void OpenFolder(dir)
  }

  function contextMenuRename() {
    if (!contextMenu) return
    renameFeaturePath = contextMenu.path
    dismissContextMenu()
    showRenameFeature = true
  }

  async function renameFeature(path: string, newName: string) {
    if (!path) return
    try {
      const newPath = await RenameFeature(path, newName)
      const wasActive = activeTab === path
      if (tabs.some((t) => t.path === path)) {
        tabs = tabs.map((t) => (t.path === path ? { ...t, path: newPath } : t))
        if (wasActive) activeTab = newPath
      }
      batchSelected = batchSelected.map((p) => (p === path ? newPath : p))
      await refreshProject()
      if (wasActive) await loadFeature(newPath)
      appendLog(`Переименован: ${basename(newPath)}`)
    } catch (e: any) {
      appendLog(`Ошибка: ${e}`)
    }
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
      schedulePersistSession()
    } catch (e: any) {
      appendLog(`Ошибка: ${e}`)
      setStatus(String(e), 'error')
    }
  }

  async function closeProject() {
    if (!projectPath) return
    const dirty = tabs.filter((t) => t.dirty)
    if (dirty.length > 0) {
      const ok = await askConfirm({
        title: 'Закрыть проект',
        message: `Есть ${dirty.length} несохранённых файл(ов). Закрыть проект без сохранения?`,
        confirmLabel: 'Закрыть',
        danger: true,
      })
      if (!ok) return
    }
    projectPath = ''
    features = []
    tags = []
    featureTags = {}
    for (const t of tabs) {
      monaco?.releaseTab(t.path)
    }
    monaco?.retainTabs([])
    tabs = []
    activeTab = WELCOME_KEY
    welcomeTabVisible = true
    editorText = ''
    monaco?.activateTab(null, '')
    batchSelected = []
    batchMode = false
    editorValidationIssues = []
    appendLog('Проект закрыт.')
    syncIdleStatus()
    schedulePersistSession()
  }

  function startVanessaPoll() {
    stopVanessaPoll()
    vanessaPollTimer = setInterval(async () => {
      try {
        let dir = vanessaWatchDir
        if (!dir) {
          const dirs = await ListVanessaRunDirs(1)
          dir = dirs[0] || ''
          if (dir) vanessaWatchDir = dir
        }
        if (dir) {
          vanessaSnapshot = await PollVanessaRun(dir, vanessaPlannedTotal)
        }
      } catch {
        /* offline poll */
      }
    }, 2000)
  }

  function stopVanessaPoll() {
    if (vanessaPollTimer) {
      clearInterval(vanessaPollTimer)
      vanessaPollTimer = null
    }
  }

  async function openVanessaMonitor() {
    if (!projectPath) return
    showVanessaMonitor = true
    vanessaPlannedTotal = Math.max(1, features.length)
    try {
      const dirs = await ListVanessaRunDirs(1)
      if (dirs[0]) {
        vanessaWatchDir = dirs[0]
        vanessaSnapshot = await PollVanessaRun(dirs[0], vanessaPlannedTotal)
      }
    } catch {
      /* ignore */
    }
  }

  function buildVanessaPluginRequest(): gui.PluginRunRequest {
    const exclude = vanessaExcludeTags
      .split(',')
      .map((t) => t.trim())
      .filter(Boolean)
    return {
      name: 'vanessa',
      dryRun: vanessaDry,
      tag: vanessaTag.trim(),
      excludeTags: exclude,
      scenario: vanessaScenario.trim(),
      rerunFailedRunDir: vanessaRerunDir.trim(),
      installEpf: vanessaInstallEpf,
      epfUrl: vanessaEpfUrl.trim(),
      epfDest: vanessaEpfDest.trim(),
      platformExe: vanessaPlatformExe.trim(),
      epfPath: vanessaEpfPath.trim(),
      ibConnection: vanessaIB.trim(),
      reportAllure: vanessaReportAllure,
      vaDir: vanessaVaDir.trim(),
      vaFiles: vanessaVaFiles.trim(),
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
    batchSelected = toggleBatchPath(batchSelected, path)
  }

  function toggleBatchMode() {
    if (batchMode) {
      batchMode = false
      batchSelected = []
      return
    }
    batchMode = true
    const tree = catalogViewState.tree
    deferToNextFrame(() => {
      if (batchMode) batchSelected = selectAllFeaturesUnder(tree)
    })
  }

  function onCatalogToggleBatch(path: string) {
    if (!batchMode) batchMode = true
    toggleBatchFeature(path)
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

  async function materializeRunTargets(paths: string[]): Promise<string[]> {
    const diskTargets: string[] = []
    for (const path of paths) {
      const tab = tabs.find((t) => t.path === path)
      if (!tab) {
        if (!isUntitled(path)) diskTargets.push(path)
        continue
      }
      if (isUntitled(path) || tab.dirty) {
        let content = path === activeTab ? editorText : tabEditorText(tab)
        if (!content && tabNeedsDiskReload(tab)) {
          content = await ReadFeature(path)
        }
        diskTargets.push(await WriteTempFeature(content))
      } else {
        diskTargets.push(path)
      }
    }
    return diskTargets
  }

  async function runBatchSelected(dryRun = false) {
    if (!batchSelected.length) return
    let opts: RunForm = { ...lastRun, dryRun }
    if (
      batchSelected.length === 1 &&
      batchSelected[0] === activeTab &&
      !isWelcome
    ) {
      const scenario = scenarioAtLine(editorText, monaco?.getCursorLine() ?? 1)
      if (scenario) opts = { ...opts, scenario }
    }
    await executeRun(opts, batchSelected)
  }

  function runPrimary(dryRun = false) {
    if (batchSelected.length > 0) {
      void runBatchSelected(dryRun)
      return
    }
    void executeRun({ ...lastRun, dryRun })
  }

  async function runScenarioAtLine(line: number, dryRun = false, scenarioOverride = '', partial = false) {
    if (!activeTab || isWelcome) return
    let scenario = scenarioOverride || scenarioAtLine(editorText, line)
    let startStep = -1
    let endStep = -1
    if (partial) {
      const resolved = await ResolveRunFromLine(editorText, line)
      if (resolved.scenario) scenario = resolved.scenario
      if (resolved.partial && resolved.startStep >= 0) {
        startStep = resolved.startStep
        endStep = resolved.endStep >= 0 ? resolved.endStep : -1
      }
    }
    monaco?.gotoLine(line)
    const runOpts = { ...lastRun, dryRun, scenario, startStep, endStep }
    const range = partialRunLogSuffix(startStep, endStep)
    if (partial && startStep >= 0) {
      appendLog(`${dryRun ? 'Dry-run' : 'Запуск'} с шага ${startStep + 1} сценария «${scenario}»${range}…`)
    } else if (scenario) {
      appendLog(`${dryRun ? 'Dry-run' : 'Запуск'} сценария «${scenario}» (строка ${line})${range}…`)
    } else {
      appendLog(`${dryRun ? 'Dry-run' : 'Запуск'} файла (строка ${line})${range}…`)
    }
    await executeRun(runOpts, [activeTab])
  }

  async function runScenarioToLine(line: number, dryRun = false) {
    if (!activeTab || isWelcome) return
    const resolved = await ResolveRunToLine(editorText, line)
    let scenario = resolved.scenario || scenarioAtLine(editorText, line)
    let startStep = -1
    let endStep = -1
    if (resolved.partial && resolved.endStep >= 0) {
      endStep = resolved.endStep
    }
    monaco?.gotoLine(line)
    const runOpts = { ...lastRun, dryRun, scenario, startStep, endStep }
    const range = partialRunLogSuffix(startStep, endStep)
    if (endStep >= 0 && scenario) {
      appendLog(`${dryRun ? 'Dry-run' : 'Запуск'} до шага ${endStep + 1} сценария «${scenario}»${range}…`)
    } else if (scenario) {
      appendLog(`${dryRun ? 'Dry-run' : 'Запуск'} сценария «${scenario}» (строка ${line})…`)
    } else {
      appendLog(`${dryRun ? 'Dry-run' : 'Запуск'} файла (строка ${line})…`)
    }
    await executeRun(runOpts, [activeTab])
  }

  async function runCurrentScenario(dryRun = false) {
    const line = monaco?.getCursorLine() ?? 1
    await runScenarioAtLine(line, dryRun)
  }

  function runHotkeyAction(id: HotkeyId) {
    switch (id) {
      case 'save':
        void saveFeature()
        break
      case 'save-as':
        void saveFeatureAs()
        break
      case 'run':
        if (!isWelcome && (activeTab || projectPath || batchSelected.length > 0)) void runPrimary(false)
        break
      case 'run-current':
        if (!isWelcome && activeTab) void runCurrentScenario(false)
        break
      case 'browser':
        void openBrowser()
        break
      case 'record':
        beginRecord()
        break
      case 'new':
        newScenario()
        break
      case 'open':
        void openFileDialog()
        break
      case 'find':
        openFindReplace()
        break
      case 'steps-help':
        openStepsHelp()
        break
      case 'hotkeys':
        showHotkeys = true
        break
      case 'settings':
        openSettings()
        break
      case 'palette':
        showCommandPalette = true
        break
      case 'snippets':
        openSnippetPalette()
        break
      case 'journal':
        bottomPanelOpen = !bottomPanelOpen
        if (bottomPanelOpen) bottomTab = 'journal'
        break
      case 'escape':
        openMenu = null
        if (showCommandPalette) showCommandPalette = false
        if (showSnippetPalette) showSnippetPalette = false
        break
    }
  }

  function onGlobalKeydown(e: KeyboardEvent) {
    if (shouldIgnoreAppHotkey(e)) return
    const id = matchHotkey(e)
    if (!id) return
    if (id === 'escape') {
      const target = e.target
      if (target instanceof Element && target.closest('.modal-backdrop, .palette-backdrop')) {
        return
      }
    }
    e.preventDefault()
    e.stopPropagation()
    runHotkeyAction(id)
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
      sidebarWidth = clampSidebarWidth(startW + (ev.clientX - startX))
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
      bottomPanelHeight = Math.max(80, Math.min(window.innerHeight * 0.6, startH + (startY - ev.clientY)))
      bottomPanelHeight = clampBottomPanelHeight(bottomPanelHeight, window.innerHeight)
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
      stepsPanelHeight = clampStepsPanelHeight(stepsPanelHeight, window.innerHeight)
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

  function syncViewportLayout() {
    if (typeof window === 'undefined') return
    viewportWidth = window.innerWidth
    viewportHeight = window.innerHeight
    viewportAutoCompact = shouldAutoCompactToolbar(viewportWidth)
    bottomPanelHeight = clampBottomPanelHeight(bottomPanelHeight, viewportHeight)
    stepsPanelHeight = clampStepsPanelHeight(stepsPanelHeight, viewportHeight)
    syncToolbarDensity()
  }

  function syncToolbarDensity() {
    if (!actionBarEl) return
    const barWidth = actionBarEl.getBoundingClientRect().width
    const urlBlock = actionBarEl.querySelector('.url-block') as HTMLElement | null
    const urlWidth = urlBlock?.getBoundingClientRect().width ?? 0
    const available = barWidth - urlWidth - 48
    toolbarIconOnly = available < toolbarIconOnlyThreshold(barWidth)
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

  function runMenuAction(action: () => void) {
    closeMenu()
    action()
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
    monaco?.refreshInlayHints()
  }

  async function serveAllureReport(path = '') {
    appendLog('Запуск Allure serve…')
    const result = await ServeAllure(path)
    if (result.output) appendLog(result.output.trimEnd())
    if (result.error) appendLog(`Ошибка: ${result.error}`)
  }

  async function openHtmlReport(path = '') {
    const result = await OpenHTMLReport(path)
    if (result.error) {
      appendLog(`Ошибка открытия отчёта: ${result.error}`)
      return
    }
    if (result.output) appendLog(`Открыт отчёт: ${result.output.trim()}`)
  }

  async function openArtifactPath(path: string) {
    if (!path) return
    try {
      await OpenFolder(path)
    } catch (e: unknown) {
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
    projectScenarios = await ListScenarioTitles().catch(() => [])
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

  function pluginRunTitle(name: string): string {
    const entry = installedPlugins.find((p) => p.name === name)
    if (entry) return pluginLabel(entry)
    return name
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
      showOpenProject = true
      return
    }
    await openProjectAt(path)
  }

  function trimTabsMemory() {
    if (isWelcome || !activeTab) {
      tabs = trimRetainedTabBodies(tabs, activeTab || '')
    } else {
      tabs = trimRetainedTabBodies(tabs, activeTab)
    }
    monaco?.retainTabs(pathsToRetainModels(tabs, isWelcome ? '' : activeTab))
  }

  function warnManyOpenTabs() {
    if (tabs.length >= MAX_OPEN_EDITOR_TABS) {
      appendLog(
        `Открыто ${tabs.length} вкладок — закройте лишние, иначе растёт потребление памяти редактора.`,
      )
    }
  }

  function syncTabContent(tabPath: string) {
    if (!tabPath || tabPath === WELCOME_KEY) return
    tabs = tabs.map((t) => {
      if (t.path !== tabPath) return t
      const dirty = editorText !== t.content
      if (!dirty) {
        if (!t.dirty && t.draft === undefined) return t
        return { ...t, dirty: false, draft: undefined }
      }
      return { ...t, draft: editorText, dirty: true }
    })
  }

  function syncActiveTabContent() {
    if (isWelcome || !activeTab) return
    syncTabContent(activeTab)
  }

  function markActiveTabSaved(text: string, tabPath = activeTab) {
    if (isWelcome || !tabPath) return
    tabs = tabs.map((t) =>
      t.path === tabPath ? { ...t, content: text, dirty: false, draft: undefined } : t,
    )
  }

  async function loadFeature(path: string) {
    const leavingTab = activeTab
    if (leavingTab && !isWelcome && leavingTab !== path) {
      syncTabContent(leavingTab)
    }
    const existing = tabs.find((t) => t.path === path)
    if (existing) {
      let text = tabEditorText(existing)
      if (tabNeedsDiskReload(existing)) {
        try {
          text = await ReadFeature(path)
          tabs = tabs.map((t) =>
            t.path === path ? { ...t, content: text, dirty: false, draft: undefined, unloaded: false } : t,
          )
        } catch (e: any) {
          appendLog(`Ошибка открытия: ${e}`)
          return
        }
      }
      welcomeTabVisible = false
      await applyEditorText(text, { saved: !existing.dirty, switchTab: true, tabPath: path, skipValidate: true })
      activeTab = path
      trimTabsMemory()
      stepsPanelCollapsed = resolveStepsPanelCollapsed()
      schedulePersistSession()
      return
    }
    try {
      const diskContent = await ReadFeature(path)
      let content = diskContent
      let dirty = false
      try {
        const draft = await LoadFeatureDraft(path)
        if (draft && draft.trim() !== diskContent.trim()) {
          content = draft
          dirty = true
          appendLog(`Восстановлен черновик: ${basename(path)}`)
        }
      } catch {
        /* no draft */
      }
      tabs = [...tabs, { path, content, dirty }]
      warnManyOpenTabs()
      welcomeTabVisible = false
      await rememberFeature(path)
      const recents = await loadRecents()
      recentFeatures = recents.features
      await applyEditorText(content, { saved: !dirty, switchTab: true, tabPath: path, skipValidate: true })
      activeTab = path
      trimTabsMemory()
      stepsPanelCollapsed = resolveStepsPanelCollapsed()
      schedulePersistSession()
    } catch (e: any) {
      appendLog(`Ошибка открытия: ${e}`)
    }
  }

  function selectTab(path: string) {
    if (path === WELCOME_KEY) {
      const leavingTab = activeTab
      if (leavingTab && !isWelcome) {
        syncTabContent(leavingTab)
      }
      void applyEditorText('', { switchTab: true, tabPath: null, skipValidate: true })
      welcomeTabVisible = true
      activeTab = WELCOME_KEY
      trimTabsMemory()
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
        void loadFeature(next.path)
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
    if (tab && tabIsUnsaved(tab)) {
      pendingCloseTab = path
      return
    }
    finalizeCloseTab(path)
  }

  function finalizeCloseTab(path: string) {
    monaco?.releaseTab(path)
    tabs = tabs.filter((t) => t.path !== path)
    trimTabsMemory()
    if (activeTab === path) {
      const next = tabs[tabs.length - 1]
      if (next) {
        activeTab = next.path
        void loadFeature(next.path)
      } else {
        welcomeTabVisible = true
        activeTab = WELCOME_KEY
        void applyEditorText('', { switchTab: true, tabPath: null, skipValidate: true })
        stepStatus = '0 шагов'
        stepStatusError = false
      }
    }
    schedulePersistSession()
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

  async function saveFeatureAs() {
    if (!activeTab || isWelcome) return
    const picked = await PickSaveFile('Сохранить как', basename(activeTab))
    if (!picked) return
    try {
      await SaveFeature(picked, editorText)
      const oldPath = activeTab
      tabs = tabs.map((t) =>
        t.path === oldPath ? { path: picked, content: editorText, dirty: false, draft: undefined } : t,
      )
      activeTab = picked
      await rememberFeature(picked)
      await refreshProject()
      appendLog(`Сохранено как: ${basename(picked)}`)
      setStatus('Сохранено', 'success')
    } catch (e: any) {
      appendLog(`Ошибка сохранения: ${e}`)
      setStatus('Ошибка сохранения', 'error')
    }
  }

  async function saveFeature() {
    if (!activeTab || isWelcome) return
    if (isUntitled(activeTab)) {
      await saveFeatureAs()
      return
    }
    try {
      let text = editorText
      if (editorSettings.formatOnSave) {
        text = await FormatFeature(text)
        if (text !== editorText) {
          editorText = text
          monaco?.setContent(text)
        }
      }
      const autoFixed = await runScenarioHintsAutoFix(text)
      if (autoFixed !== text) {
        text = autoFixed
        editorText = text
        monaco?.setContent(text)
        appendLog('Применены авто-исправления подсказок сценария')
      }
      await SaveFeature(activeTab, text)
      markActiveTabSaved(text)
      try {
        await ClearFeatureDraft(activeTab)
      } catch {
        /* ignore */
      }
      appendLog(`Сохранено: ${basename(activeTab)}`)
      setStatus('Сохранено', 'success')
    } catch (e: any) {
      appendLog(`Ошибка сохранения: ${e}`)
      setStatus('Ошибка сохранения', 'error')
    }
  }

  let validateGeneration = 0
  let validateDebounceTimer: ReturnType<typeof setTimeout> | null = null

  function scheduleValidateEditor(delayMs = 300) {
    if (validateDebounceTimer) clearTimeout(validateDebounceTimer)
    if (delayMs <= 0) {
      validateDebounceTimer = null
      void validateEditor()
      return
    }
    validateDebounceTimer = setTimeout(() => {
      validateDebounceTimer = null
      void validateEditor()
    }, delayMs)
  }

  async function validateEditor() {
    const generation = ++validateGeneration
    const tabAtStart = activeTab
    const textAtStart = editorText
    try {
      const issues = await ValidateFeature(textAtStart)
      if (generation !== validateGeneration || tabAtStart !== activeTab) return
      editorValidationIssues = issues || []
      monaco?.setMarkers(editorValidationIssues)
      await refreshEditorSteps()
      if (generation !== validateGeneration || tabAtStart !== activeTab) return
      if (editorValidationIssues.length > 0) {
        stepStatus = `${stepCount} шагов · ошибок ${editorValidationIssues.length}`
        stepStatusError = true
        setStatus('Ошибка в сценарии', 'error')
      } else {
        stepStatus = `${stepCount} шагов`
        stepStatusError = false
      }
    } catch {
      if (generation !== validateGeneration || tabAtStart !== activeTab) return
      editorValidationIssues = []
      stepStatus = `${stepCount} шагов`
      await refreshEditorSteps()
    }
    if (generation !== validateGeneration || tabAtStart !== activeTab) return
    if (editorSettings.scenarioHints) {
      await refreshEditorScenarioHints()
    } else {
      editorScenarioHints = []
    }
  }

  function gotoEditorLine(line: number) {
    monaco?.gotoLine(line)
    editorCursorLine = line
  }

  function validateProjectHint(): string {
    if (validatePanelIssues.length > 0 || editorValidationIssues.length > 0) return ''
    if (isWelcome) return 'Откройте сценарий для проверки шагов'
    return 'Ошибок в шагах сценария нет. Проверка селекторов в браузере — меню «Запись и тест».'
  }

  function validatePanelDisplayIssues(): gui.ValidationIssue[] {
    return validatePanelIssues.length > 0 ? validatePanelIssues : editorValidationIssues
  }

  async function onEditorChange(text: string) {
    editorText = text
    syncActiveTabContent()
    schedulePersistSession()
    if (editorSettings.validateOnType) {
      scheduleValidateEditor()
    }
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

  function editorScenarioNames(): string[] {
    return !isWelcome && activeTab ? listScenarioTitles(editorText) : []
  }

  function cursorScenarioName(): string {
    return !isWelcome && activeTab ? scenarioAtLine(editorText, monaco?.getCursorLine() ?? 1) : ''
  }

  function dialogScenarioNames(): string[] {
    return mergeScenarioNames(editorScenarioNames(), projectScenarios)
  }

  function openRunDialog(title: string, defaults: Partial<RunForm>) {
    runDialogTitle = title
    runDialogScenarios = dialogScenarioNames()
    const cursorScenario = cursorScenarioName()
    runForm = {
      ...lastRun,
      baseUrl: lastRun.baseUrl || startURL || '',
      scenario: defaults.scenario ?? cursorScenario ?? lastRun.scenario ?? '',
      ...defaults,
    }
    showRun = true
  }

  function openVanessaDialog(dry: boolean, preferRerun = false) {
    vanessaDry = dry
    vanessaTag = ''
    vanessaExcludeTags = ''
    vanessaScenario = cursorScenarioName()
    vanessaDialogScenarios = dialogScenarioNames()
    vanessaRerunDir = ''
    vanessaPreferRerun = preferRerun
    vanessaInstallEpf = false
    vanessaEpfUrl = ''
    vanessaEpfDest = ''
    vanessaPlatformExe = ''
    vanessaEpfPath = ''
    vanessaIB = ''
    vanessaReportAllure = false
    vanessaVaDir = ''
    vanessaVaFiles = ''
    showVanessaRun = true
  }

  async function confirmVanessaRun() {
    showVanessaRun = false
    vanessaPlannedTotal = Math.max(1, features.length)
    showVanessaMonitor = true
    vanessaSnapshot = new gui.VanessaRunSnapshotDTO()
    StartVanessaRun(buildVanessaPluginRequest())
  }

  async function executeRun(opts: RunForm, targets: string[] = []) {
    syncActiveTabContent()

    let runTargets = targets
    if (runTargets.length === 0 && activeTab && !isWelcome) {
      runTargets = [activeTab]
    }
    if (runTargets.length === 0 && !projectPath) {
      appendLog('Откройте сценарий для запуска')
      return
    }

    const diskTargets = runTargets.length ? await materializeRunTargets(runTargets) : []
    lastRun = { ...opts }
    showRun = false
    bottomPanelOpen = true
    bottomTab = 'journal'

    playingLabel =
      runTargets.length > 1
        ? `${runTargets.length} сценариев`
        : runTargets.length === 1
          ? featureTabLabel(runTargets[0])
          : opts.scenario || opts.tag || 'тест'

    const allureDir = opts.allure ? await scenariaSubdir('allure-results') : ''
    const traceDir = !opts.dryRun && opts.trace ? await scenariaSubdir('traces') : ''
    const videoDir = !opts.dryRun && opts.video ? await scenariaSubdir('videos') : ''

    const htmlPath = opts.html ? await scenariaSubdir('report.html') : ''

    const junitPath = opts.junit ? await scenariaSubdir('junit.xml') : ''

    const summaryJsonPath = opts.summaryJson ? await scenariaSubdir('summary.json') : ''

    playing = !opts.dryRun
    setStatus('▶ Идёт тест', 'busy')
    const range = partialRunLogSuffix(opts.startStep ?? -1, opts.endStep ?? -1)
    if (targets.length) {
      appendLog(`Запуск ${targets.length} сценариев${range}…`)
    } else if (opts.dryRun) {
      appendLog(`Dry-run${range}…`)
    } else {
      appendLog(`Запуск Playwright${range}…`)
    }

    const result = await Run({
      tag: opts.tag,
      scenario: opts.scenario || '',
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
      summaryJson: summaryJsonPath,
      browser: opts.dryRun ? '' : opts.browser || settingsBrowser || 'chromium',
      workers: opts.workers || settingsWorkers || 1,
      slowMo: opts.dryRun ? 0 : (opts.slowMo > 0 ? opts.slowMo : settingsSlowMo),
      baseUrl: opts.dryRun ? '' : (opts.baseUrl || '').trim(),
      startStep: opts.startStep ?? -1,
      endStep: opts.endStep ?? -1,
      targets: diskTargets,
    })
    playing = false
    playingLabel = ''
    if (result.output) appendLog(result.output.trimEnd())
    if (result.error) {
      appendLog(`Ошибка: ${result.error}`)
      setStatus('Ошибка теста', 'error')
      bottomTab = 'error'
    } else {
      appendLog('Завершено.')
      setStatus('Тест завершён', 'success')
      welcomePlayedSuccess = true
      bottomTab = opts.dryRun ? 'journal' : 'results'
    }
    await refreshRunResults()
    await refreshArtifacts()
    if (!result.error && opts.html && htmlPath) {
      try {
        if (await ArtifactExists(htmlPath)) {
          await openHtmlReport(htmlPath)
        }
      } catch {
        /* ignore */
      }
    }
  }

  async function scenariaSubdir(sub: string): Promise<string> {
    if (!projectPath) return ''
    try {
      const path = await ScenariaArtifactPath(sub)
      if (path) return path.replace(/\\/g, '/')
    } catch {
      /* fallback below */
    }
    return `${projectPath.replace(/\\/g, '/')}/.scenaria/${sub}`
  }

  function confirmRun() {
    executeRun(runForm)
  }

  async function validateProject(browser: boolean, browserName = settingsBrowser || 'chromium', targets: string[] = []) {
    if (!projectPath) return
    appendLog(browser ? 'Проверка в браузере…' : 'Проверка…')
    bottomPanelOpen = true
    bottomTab = browser ? 'validate' : 'journal'
    validateCliLog = ''
    validatePanelIssues = []

    if (browser) {
      try {
        const issues = await ValidateBrowser(
          gui.ValidateRequest.createFrom({
            browser: browserName || 'chromium',
            skipBrowser: false,
            targets,
          }),
        )
        validatePanelIssues = issues || []
        const missing = validatePanelIssues.filter((i) => i.status === 'missing' || !i.status).length
        const warnings = validatePanelIssues.filter((i) => i.status === 'warning').length
        const found = validatePanelIssues.filter((i) => i.status === 'found').length
        appendLog(`Проверка в браузере: найдено ${found}, предупреждений ${warnings}, ошибок ${missing}.`)
        setStatus(missing > 0 ? 'Ошибки проверки в браузере' : 'Проверка в браузере завершена', missing > 0 ? 'error' : 'success')
      } catch (err) {
        const msg = err instanceof Error ? err.message : String(err)
        validateCliLog = msg
        appendLog(`Ошибка: ${msg}`)
        setStatus('Ошибка проверки', 'error')
      }
    } else {
      const result = await Validate({
        browser: browserName || 'chromium',
        skipBrowser: true,
        targets,
      })
      if (result.output) {
        validateCliLog = result.output.trimEnd()
        appendLog(validateCliLog)
      }
      if (result.error) {
        validateCliLog = `${validateCliLog}\n${result.error}`.trim()
        appendLog(`Ошибка: ${result.error}`)
        setStatus('Ошибка проверки', 'error')
      } else {
        appendLog('Проверка завершена.')
        setStatus('Проверка завершена', 'success')
      }
    }
    if (!isWelcome && activeTab) await validateEditor()
  }

  function openValidateDialog(syntaxOnly: boolean) {
    if (!projectPath) return
    validateSyntaxOnly = syntaxOnly
    validateBrowser = settingsBrowser || 'chromium'
    validateScope = !isWelcome && activeTab ? 'current' : 'project'
    showValidate = true
  }

  async function confirmValidate(payload: { browser: string; syntaxOnly: boolean; scope: 'project' | 'current' }) {
    showValidate = false
    const targets = payload.scope === 'current' && activeTab && !isWelcome ? [activeTab] : []
    await validateProject(!payload.syntaxOnly, payload.browser, targets)
  }

  function openInitProjectDialog() {
    if (!projectPath) return
    showInitProject = true
  }

  async function confirmInitProject() {
    showInitProject = false
    await initProject()
  }

  async function initProject() {
    const out = await InitProject()
    if (out) appendLog(out.trimEnd())
    await refreshProject()
  }

  async function openTestClientDialog() {
    if (!projectPath) return
    testClientSuggestName = ''
    testClients = await ListTestClients().catch(() => [])
    testClientSelection = runForm.testClient || testClients[0] || ''
    showTestClient = true
  }

  async function openTestClientDialogForCapture() {
    if (!projectPath) return
    if (!browserOpen && !recording) {
      appendLog('Откройте браузер (Ctrl+B), войдите на сайт, затем сохраните сессию.')
      return
    }
    testClients = await ListTestClients().catch(() => [])
    const base = (runForm.testClient || testClientSelection || 'session').trim() || 'session'
    testClientSuggestName = base
    testClientSelection = testClients.includes(base) ? base : ''
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

  function openStepsHelp(query: unknown = '') {
    if (recordingBlocksManualTools()) {
      appendLog('Поставьте запись на паузу, чтобы открыть справку по шагам')
      return
    }
    stepsHelpQuery = typeof query === 'string' ? query : ''
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

  async function applyEditorText(
    text: string,
    options?: { saved?: boolean; switchTab?: boolean; tabPath?: string | null; skipValidate?: boolean },
  ) {
    editorText = text
    if (options?.switchTab) {
      monaco?.activateTab(options.tabPath ?? null, text)
      if (options?.saved) {
        markActiveTabSaved(text, options.tabPath ?? activeTab)
      } else if (options?.tabPath) {
        tabs = tabs.map((t) =>
          t.path === options.tabPath ? { ...t, draft: text, dirty: true } : t,
        )
      }
    } else {
      monaco?.setContent(text)
      if (options?.saved) {
        markActiveTabSaved(text)
      } else {
        syncActiveTabContent()
      }
    }
    if (options?.skipValidate) {
      void refreshEditorSteps()
      return
    }
    await validateEditor()
  }

  async function refactorUpdateUrls() {
    if (isWelcome) return
    showRefactorUrl = true
  }

  async function applyRefactorUrl(newUrl: string) {
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

  function openImportFeaturesDialog() {
    if (!projectPath) return
    importDestDir = projectPath.replace(/\\/g, '/')
    showImportFeatures = true
  }

  async function confirmImportFeatures(payload: { destDir: string; paths: string[] }) {
    if (!projectPath || importFeaturesBusy) return
    importFeaturesBusy = true
    showImportFeatures = false
    try {
      await importDroppedFeatures(payload.destDir, payload.paths)
    } finally {
      importFeaturesBusy = false
    }
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
      platformExe: opts.platformExe || '',
      epfPath: opts.epfPath || '',
      ibConnection: opts.ibConnection || '',
      reportAllure: opts.reportAllure || false,
      vaDir: opts.vaDir || '',
      vaFiles: opts.vaFiles || '',
    })
    if (result.output) appendLog(result.output.trimEnd())
    if (result.error) appendLog(`Ошибка: ${result.error}`)
  }

  function openPluginRun(name: string, dry = false) {
    const entry = installedPlugins.find((p) => p.name === name)
    if (entry?.vanessa || name === 'vanessa') {
      openVanessaDialog(dry)
      return
    }
    pluginRunName = name
    pluginRunDry = dry
    pluginRunTag = ''
    pluginRunScenario = cursorScenarioName()
    pluginRunScenarios = dialogScenarioNames()
    showPluginRun = true
  }

  async function confirmPluginRun(payload: { tag: string; scenario: string; dryRun: boolean }) {
    showPluginRun = false
    await runPlugin(pluginRunName, payload.dryRun, { tag: payload.tag, scenario: payload.scenario })
  }

  function openBaselineRecordDialog() {
    if (!projectPath) return
    recordMode = 'baseline'
    recordOutput = 'recorded.feature'
    recordURL = startURL || recordURL || 'https://example.com'
    if (activeTab && !isWelcome) {
      recordFeatureName = basename(activeTab).replace(/\.feature$/i, '')
      recordScenarioName = 'Базовый сценарий'
    } else {
      recordFeatureName = 'Записанный сценарий'
      recordScenarioName = 'Базовый сценарий'
    }
    showRecord = true
  }

  async function saveBaselineRecord(payload: {
    output: string
    featureName: string
    scenarioName: string
    steps: string[]
  }) {
    if (!projectPath || baselineBusy) return
    baselineBusy = true
    showRecord = false
    appendLog('Создание feature из шагов…')
    try {
      const result = await RecordBaseline({
        output: payload.output || 'recorded.feature',
        featureName: payload.featureName,
        scenarioName: payload.scenarioName,
        steps: payload.steps,
      })
      if (result.output) appendLog(result.output.trimEnd())
      if (result.error) {
        appendLog(`Ошибка: ${result.error}`)
        setStatus('Ошибка записи', 'error')
        return
      }
      appendLog('Feature сохранён.')
      setStatus('Feature создан', 'success')
      await refreshProject()
      const rel = (payload.output || 'recorded.feature').replace(/\\/g, '/')
      const featurePath = rel.startsWith('/') || /^[A-Za-z]:\//.test(rel)
        ? rel
        : `${projectPath.replace(/\\/g, '/')}/${rel}`
      await loadFeature(featurePath)
    } finally {
      baselineBusy = false
    }
  }

  async function checkUpdatesOnStartup() {
    if (!settingsCheckUpdatesOnStartup) return
    try {
      const info = await CheckUpdateInfo()
      updateCheckInfo = info
      updateCheckMessage = info.message || ''
      updateCheckHasUpdate = !!info.updateAvailable
      if (info.updateAvailable) {
        showUpdateCheck = true
        setStatus('Доступно обновление', 'normal')
      }
    } catch {
      /* offline or dev without wails */
    }
  }

  async function checkUpdates() {
    appendLog('Проверка обновлений…')
    try {
      const info = await CheckUpdateInfo()
      updateCheckInfo = info
      updateCheckMessage = info.message || ''
      updateCheckHasUpdate = !!info.updateAvailable
      appendLog(updateCheckMessage)
      if (info.updateAvailable) {
        appendLog(`Релиз: ${info.htmlUrl || '—'}`)
        if (info.downloadName) appendLog(`Файл: ${info.downloadName}`)
      }
      appendLog('Проверка завершена.')
      showUpdateCheck = true
    } catch (err) {
      const msg = err instanceof Error ? err.message : String(err)
      updateCheckMessage = msg
      updateCheckHasUpdate = false
      appendLog(`Ошибка: ${msg}`)
      showUpdateCheck = true
    }
  }

  async function openUpdateRelease() {
    const url = updateCheckInfo?.htmlUrl
    if (!url) return
    try {
      await OpenExternalURL(url)
    } catch (err) {
      appendLog(`Не удалось открыть ссылку: ${err instanceof Error ? err.message : String(err)}`)
    }
  }

  async function applyUpdate() {
    if (updateDownloading) return
    updateDownloading = true
    updateProgress = { stage: 'check', message: 'Запуск обновления…', percent: 0 }
    appendLog('Установка обновления…')
    let lastStage = ''
    const onProgress = (payload: gui.UpdateProgressDTO) => {
      updateProgress = payload
      if (payload?.stage && payload.stage !== lastStage) {
        lastStage = payload.stage
        if (payload.message) appendLog(payload.message)
      }
    }
    const onFinished = (result: gui.RunResult) => {
      EventsOff('update-progress', 'update-finished')
      if (result?.error) {
        appendLog(`Ошибка обновления: ${result.error}`)
        updateDownloading = false
        updateProgress = null
        return
      }
      updateProgress = { stage: 'restart', message: 'Перезапуск приложения…', percent: 100 }
      appendLog('Приложение перезапускается…')
    }
    EventsOn('update-progress', onProgress)
    EventsOn('update-finished', onFinished)
    try {
      await ApplyUpdate()
    } catch (err) {
      EventsOff('update-progress', 'update-finished')
      appendLog(`Ошибка обновления: ${err instanceof Error ? err.message : String(err)}`)
      updateDownloading = false
      updateProgress = null
    }
  }

  async function downloadUpdate() {
    if (updateDownloading) return
    updateDownloading = true
    appendLog('Скачивание обновления…')
    try {
      const path = await DownloadUpdate()
      appendLog(`Скачано: ${path}`)
      updateCheckMessage = `${updateCheckMessage}\n\nФайл сохранён:\n${path}`.trim()
      const folder = path.replace(/[\\/][^\\/]+$/, '')
      if (folder) await OpenFolder(folder)
    } catch (err) {
      appendLog(`Ошибка скачивания: ${err instanceof Error ? err.message : String(err)}`)
    } finally {
      updateDownloading = false
    }
  }

  async function openSettings() {
    const s = await LoadSettings()
    applySettingsFromDTO(s)
    settingsDialogBaseline = s
    showSettings = true
  }

  async function persistSettings() {
    syncActiveTabContent()
    const sessionTabs = buildSessionTabsSnapshot(tabs, activeTab, editorText, WELCOME_KEY)
    await SaveSettings(
      gui.AppSettingsDTO.createFrom({
      browser: settingsBrowser,
      headless: settingsHeadless,
      parallelWorkers: settingsWorkers,
      slowMo: settingsSlowMo,
      maxLoopIterations: settingsLoops,
      scrollBeforeClick: settingsScrollBeforeClick,
      hoverRecordMinMs: settingsHoverRecordMinMs,
      sessionProject: projectPath,
      openTabs: sessionTabs.openTabs,
      untitledTabs: sessionTabs.untitledTabs,
      activeTab: sessionTabs.activeTab,
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
      selectorClickStrategies: settingsSelectorClickStrategies,
      selectorInputStrategies: settingsSelectorInputStrategies,
      editor: editorSettingsToDTO(editorSettings),
    }),
    )
  }

  async function syncRecordingOptions() {
    if (!recording) return
    try {
      await UpdateRecordingOptions(
        filterRecording,
        navOnlyRecording,
        hoverRecord,
        settingsHeadless,
        settingsScrollBeforeClick,
        settingsHoverRecordMinMs,
      )
      await persistSettings()
    } catch {
      /* session may be closing */
    }
  }

  async function applySettings() {
    editorSettings = { ...editorSettings }
    stepsPanelCollapsed = resolveStepsPanelCollapsed()
    stepsPanelTab = editorSettings.stepsPanelView
    setStepHoverEnabled(() => editorSettings.stepHover)
    lastRun = {
      ...lastRun,
      workers: settingsWorkers,
      slowMo: settingsSlowMo,
      browser: settingsBrowser,
    }
    if (recording) await syncRecordingOptions()
    await persistSettings()
    monaco?.applyEditorSettings(editorSettings)
    if (editorSettings.scenarioHints) {
      await refreshEditorScenarioHints()
    } else {
      editorScenarioHints = []
    }
    if (editorSettings.validateOnType && activeTab && !isWelcome) void validateEditor()
    settingsDialogBaseline = null
    showSettings = false
    appendLog('Настройки сохранены.')
  }

  function cancelSettings() {
    if (settingsDialogBaseline) {
      applySettingsFromDTO(settingsDialogBaseline)
      monaco?.applyEditorSettings(editorSettings)
    }
    settingsDialogBaseline = null
    showSettings = false
  }

  function recordStartURL(): string {
    return resolveRecordStartURL({
      editorText: !isWelcome ? editorText : '',
      startURL,
      recordURL,
      lastRunBaseUrl: lastRun.baseUrl,
    })
  }

  function prepareRecordDialogDefaults() {
    recordMode = 'live'
    recordAppendTo = ''
    recordTestClient = runForm.testClient || testClientSelection || ''
    if (projectPath) {
      recordOutput = `${projectPath.replace(/\\/g, '/')}/recorded.feature`
    }
    recordURL = recordStartURL()
    if (activeTab && !isWelcome) {
      recordFeatureName = basename(activeTab).replace(/\.feature$/i, '')
      recordScenarioName = 'Запись'
    } else {
      recordFeatureName = 'Записанный сценарий'
      recordScenarioName = 'Запись'
    }
  }

  function beginRecord() {
    prepareRecordDialogDefaults()
    showRecord = true
  }

  async function openBrowser() {
    if (!projectPath) {
      appendLog('Сначала откройте проект')
      return
    }
    if (browserOpen || recording) {
      await focusBrowser()
      return
    }
    prepareRecordDialogDefaults()
    recordURL = recordStartURL()
    showRecord = false
    if (recordURL) {
      appendLog(`Открываю браузер: ${recordURL}`)
    } else {
      appendLog('Открываю браузер (пустая страница — перейдите вручную или начните запись)')
    }
    await OpenBrowser({
      url: recordURL,
      output: recordOutput,
      idleSeconds: recordIdle,
      headless: false,
      filterRecording,
      navOnlyRecording,
      hoverRecord,
      appendTo: recordAppendTo,
      testClient: recordTestClient,
      featureName: recordFeatureName,
      scenarioName: recordScenarioName,
    })
    recordAppendTo = ''
  }

  async function startRecord(opts?: { headed?: boolean }) {
    await persistSettings()
    if (recording) {
      if (browserOpen) await focusBrowser()
      return
    }
    lastRecordTarget = recordAppendTo || recordOutput
    let sessionOpen = browserOpen
    if (!sessionOpen) {
      try {
        const s = await PollBrowserSession()
        sessionOpen = s.browserOpen
      } catch {
        /* dev without wails */
      }
    }
    if (sessionOpen) {
      try {
        await BeginRecordingCapture()
      } catch (e: unknown) {
        appendLog(`Запись: ${e}`)
      }
      return
    }
    await StartRecord({
      url: recordURL,
      output: recordOutput,
      idleSeconds: recordIdle,
      headless: opts?.headed ?? settingsHeadless,
      filterRecording,
      navOnlyRecording,
      hoverRecord,
      appendTo: recordAppendTo,
      testClient: recordTestClient,
      featureName: recordFeatureName,
      scenarioName: recordScenarioName,
      browseOnly: false,
    })
    recordAppendTo = ''
  }

  function resolveRecordFeaturePath(outputPath = ''): string {
    const append = (recordAppendTo || '').trim().replace(/\\/g, '/')
    if (append) return append
    const target = (outputPath || lastRecordTarget || recordOutput || '').trim().replace(/\\/g, '/')
    if (!target || !projectPath) return target
    if (target.startsWith('/') || /^[A-Za-z]:\//.test(target)) return target
    return `${projectPath.replace(/\\/g, '/')}/${target}`
  }

  async function prepareRecordEditorTab(outputPath = '') {
    const appendPath = (recordAppendTo || '').trim().replace(/\\/g, '/')
    if (appendPath && !isUntitled(appendPath)) {
      await loadFeature(appendPath)
      return
    }
    const featurePath = resolveRecordFeaturePath(outputPath)
    if (featurePath && tabs.some((t) => t.path === featurePath)) {
      await loadFeature(featurePath)
      return
    }
    if (!activeTab || isWelcome) {
      await openUntitledTab(
        buildFeatureTemplate({
          title: recordFeatureName,
          scenario: recordScenarioName,
          startUrl: recordURL || startURL || 'https://example.com',
        }),
        'zapis.feature',
      )
    }
  }

  async function applyLiveRecordedStep(index: number, line: string) {
    if (!line.trim() || isWelcome || !activeTab) return
    const result = upsertRecordedStepInText(editorText, index, line, liveRecordStepLines)
    liveRecordStepLines = result.lineByIndex
    await applyEditorText(result.text, { skipValidate: true })
    scheduleValidateEditor(150)
  }

  function handleRecordStopped() {
    recording = false
    recordPaused = false
    liveRecordStepLines = {}
    if (browserOpen) {
      setStatus('Браузер открыт', 'busy')
      appendLog('Запись остановлена. Браузер остаётся открытым.')
    } else {
      syncIdleStatus()
    }
  }

  async function handleRecordSessionEnd(result: gui.RunResult, kind: 'record' | 'browse') {
    recording = false
    browserOpen = false
    recordPaused = false
    showRecord = false
    liveRecordStepLines = {}
    lastRecordTarget = ''
    if (result.output) appendLog(result.output)
    if (result.error) appendLog(`Ошибка записи: ${result.error}`)
    else if (kind === 'record') {
      appendLog('Браузер закрыт. Сохраните сценарий (Ctrl+S), если ещё не сохраняли.')
    }
    statusTone = 'normal'
    syncIdleStatus()
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
    void syncBrowserStateFromBackend()
  }

  function applyBrowserSessionState(s: { browserOpen: boolean; recording: boolean; paused: boolean }) {
    browserOpen = s.browserOpen
    recording = s.recording
    recordPaused = s.paused
  }

  function handleBrowserLost() {
    if (!browserOpen && !recording) return
    applyBrowserSessionState({ browserOpen: false, recording: false, paused: false })
    statusTone = 'normal'
    syncIdleStatus()
  }

  function startBrowserWatch() {
    stopBrowserWatch()
    browserWatchTimer = setInterval(() => {
      void syncBrowserStateFromBackend()
    }, 400)
  }

  function stopBrowserWatch() {
    if (browserWatchTimer) {
      clearInterval(browserWatchTimer)
      browserWatchTimer = null
    }
  }

  async function syncBrowserStateFromBackend() {
    try {
      const s = await PollBrowserSession()
      if ((browserOpen || recording) && !s.browserOpen) {
        handleBrowserLost()
        return
      }
      if (!s.browserOpen) return
      const prevRecording = recording
      const prevPaused = recordPaused
      applyBrowserSessionState({
        browserOpen: true,
        recording: s.recording,
        paused: s.paused,
      })
      if (s.recording && !prevRecording) {
        setStatus('● Идёт запись', 'busy')
      } else if (s.recording && s.paused && (!prevPaused || !prevRecording)) {
        setStatus('⏸ Пауза', 'busy')
      } else if (s.recording && !s.paused && prevPaused) {
        setStatus('● Идёт запись', 'busy')
      } else if (!s.recording && prevRecording) {
        setStatus('Браузер открыт', 'busy')
      }
    } catch {
      /* dev without wails */
    }
  }

  async function stopRecord() {
    if (recording) {
      await StopRecordingCapture()
      return
    }
    if (browserOpen) {
      await CloseBrowser()
    }
  }

  async function submitOtp(code: string) {
    const accepted = await SubmitOTPCode(code)
    if (accepted) {
      showOtp = false
      return
    }
    appendLog('Нет активного запроса OTP — код не отправлен')
  }

  async function cancelOtp() {
    await CancelOTP()
    showOtp = false
  }

  function newScenario() {
    void openUntitledTab(
      buildFeatureTemplate({
        title: 'Примеры для новичков',
        scenario: 'Первая проверка страницы',
        startUrl: startURL || recordURL || 'https://example.com',
      }),
    )
  }

  async function openUntitledTab(content: string, displayName = 'novyy-scenariy.feature') {
    const leavingTab = activeTab
    if (leavingTab && !isWelcome) {
      syncTabContent(leavingTab)
    }
    const path = makeUntitledPath(displayName)
    tabs = [...tabs, { path, content, dirty: true }]
    warnManyOpenTabs()
    welcomeTabVisible = false
    await applyEditorText(content, { switchTab: true, tabPath: path, skipValidate: true })
    activeTab = path
    trimTabsMemory()
    stepsPanelCollapsed = resolveStepsPanelCollapsed()
    scheduleValidateEditor()
    schedulePersistSession()
  }

  async function openFileDialog() {
    const path = await PickOpenFile('Открыть feature')
    if (path) await loadFeature(path)
  }

  function insertTemplate() {
    if (isWelcome) {
      newScenario()
      return
    }
    const template = buildFeatureTemplate({
      title: 'Примеры для новичков',
      scenario: 'Первая проверка страницы',
      startUrl: startURL || recordURL || 'https://example.com',
    })
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
    recordFeatureName = basename(activeTab).replace(/\.feature$/i, '')
    recordScenarioName = 'Дозапись'
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
    if (!recording && !browserOpen) {
      appendLog('Указать элемент: откройте браузер')
      return
    }
    if (recording && !recordPaused) {
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
    if (activeTab && activeTab !== WELCOME_KEY) return featureTabLabel(activeTab)
    if (projectPath) return basename(projectPath)
    return 'Открыть проект…'
  }
</script>

{#if !appReady}
  <SplashScreen
    {version}
    message={splashMessage}
    progress={splashProgress}
    fading={splashFading}
    standalone
  />
{:else}
<div class="ide" class:panel-open={bottomPanelOpen}>
  <!-- Menu bar (Python: Проект / Сценарий / Запись и тест / Вид / Справка) -->
  <div class="menubar" role="menubar">
    <div class="menu-root" class:open={openMenu === 'project'}>
      <button class="menu-trigger" on:click={(e) => toggleMenu('project', e)}>Проект</button>
      {#if openMenu === 'project'}
        <div class="menu-dropdown">
          <button class="menu-item" on:click={openExamples}>Открыть примеры сценариев</button>
          <button class="menu-item" on:click={openProjectDialog}>Открыть проект…</button>
          <button class="menu-item" on:click={closeProject} disabled={!projectPath}>Закрыть проект</button>
          <button class="menu-item" on:click={openSettings}>Настройки…<span class="menu-shortcut">Ctrl+,</span></button>
          <div class="menu-sep"></div>
          <button class="menu-item" on:click={openInitProjectDialog} disabled={!projectPath}>Init проекта</button>
        </div>
      {/if}
    </div>

    <div class="menu-root" class:open={openMenu === 'scenario'}>
      <button class="menu-trigger" on:click={(e) => toggleMenu('scenario', e)}>Сценарий</button>
      {#if openMenu === 'scenario'}
        <div class="menu-dropdown">
          <button class="menu-item" on:click={newScenario}>Новый</button>
          <button class="menu-item" on:click={openFileDialog}>Открыть…</button>
          <button class="menu-item" on:click={saveFeature} disabled={isWelcome}>Сохранить<span class="menu-shortcut">Ctrl+S</span></button>
          <button class="menu-item" on:click={saveFeatureAs} disabled={isWelcome}>Сохранить как…</button>
          <button class="menu-item" on:click={() => activeTab && !isWelcome && openDuplicateDialog(activeTab)} disabled={isWelcome}>Дублировать…</button>
          <button class="menu-item" on:click={openImportFeaturesDialog} disabled={!projectPath}>Импорт .feature…</button>
          <button
            class="menu-item"
            on:click={() => {
              if (!activeTab || isWelcome) return
              renameFeaturePath = activeTab
              showRenameFeature = true
            }}
            disabled={isWelcome}
          >
            Переименовать…
          </button>
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
        <div class="menu-dropdown">
          <button class="menu-item" on:click={() => runMenuAction(() => void openBrowser())} disabled={!projectPath}>Браузер<span class="menu-shortcut">Ctrl+B</span></button>
          <button class="menu-item" on:click={() => runMenuAction(beginRecord)} disabled={!projectPath}>Запись<span class="menu-shortcut">Ctrl+R</span></button>
          <button class="menu-item" on:click={openBaselineRecordDialog} disabled={!projectPath}>Запись из шагов…</button>
          <button class="menu-item" on:click={stopRecord} disabled={!recording && !browserOpen}>Стоп</button>
          <button class="menu-item" on:click={toggleRecordPause} disabled={!recording}>Пауза</button>
          <div class="menu-sep"></div>
          <button class="menu-item" on:click={openTestClientDialog} disabled={!projectPath}>TestClient…</button>
          <button class="menu-item" on:click={openTestClientDialogForCapture} disabled={!projectPath || (!browserOpen && !recording)}>
            Сохранить сессию браузера…
          </button>
          <button class="menu-item" on:click={openHttpAuthDialog}>HTTP Auth…</button>
          <div class="menu-sep"></div>
          <button class="menu-item" on:click={() => runPrimary(false)} disabled={isWelcome && !projectPath && !batchSelected.length}>
            Запустить<span class="menu-shortcut">Ctrl+Enter</span>
          </button>
          <button class="menu-item" on:click={() => runCurrentScenario(false)} disabled={isWelcome || !activeTab}>
            Запустить текущий сценарий<span class="menu-shortcut">Ctrl+Shift+Enter</span>
          </button>
          <button class="menu-item" on:click={() => runCurrentScenario(true)} disabled={isWelcome || !activeTab}>
            Dry-run текущего сценария
          </button>
          <button class="menu-item" on:click={() => runBatchSelected(false)} disabled={!projectPath || !batchSelected.length}>
            Запустить выбранные
          </button>
          <button class="menu-item" on:click={() => runBatchSelected(true)} disabled={!projectPath || !batchSelected.length}>
            Dry-run выбранных
          </button>
          <button class="menu-item" on:click={rerunFailed} disabled={!projectPath}>Перезапустить упавшие</button>
          <button class="menu-item" on:click={openRunHistory} disabled={!projectPath}>История запусков…</button>
          <button class="menu-item" on:click={() => openRunDialog('Запуск сценария', {})} disabled={isWelcome && !projectPath}>Запустить…</button>
          <button class="menu-item" on:click={() => openRunDialog('Запуск с тегом', {})} disabled={isWelcome && !projectPath}>
            Запустить сценарии с тегом…
          </button>
          <button class="menu-item" on:click={() => runPrimary(true)} disabled={isWelcome && !projectPath && !batchSelected.length}>Dry-run</button>
          <button
            class="menu-item"
            on:click={() => openRunDialog('Playwright', { dryRun: false, headed: true, engine: 'playwright', installPW: true })}
            disabled={isWelcome && !activeTab}
          >
            Playwright…
          </button>
          <div class="menu-sep"></div>
          <button class="menu-item" on:click={() => openValidateDialog(true)} disabled={!projectPath}>Проверить…</button>
          <button class="menu-item" on:click={() => openValidateDialog(false)} disabled={!projectPath}>Проверить в браузере…</button>
          <div class="menu-sep"></div>
          {#if hasVanessaPlugin()}
            <button class="menu-item" on:click={() => openVanessaDialog(true)} disabled={!projectPath}>Vanessa (dry)…</button>
            <button class="menu-item" on:click={() => openVanessaDialog(false)} disabled={!projectPath}>Vanessa run…</button>
            <button class="menu-item" on:click={() => openVanessaDialog(false, true)} disabled={!projectPath}>Vanessa rerun-failed…</button>
            <button class="menu-item" on:click={openVanessaSettingsDialog} disabled={!projectPath}>Настройки Vanessa…</button>
            <button class="menu-item" on:click={openVanessaMonitor} disabled={!projectPath}>Монитор Vanessa…</button>
          {/if}
        </div>
      {/if}
    </div>

    <div class="menu-root" class:open={openMenu === 'plugins'}>
      <button class="menu-trigger" on:click={(e) => toggleMenu('plugins', e)}>Плагины</button>
      {#if openMenu === 'plugins'}
        <div class="menu-dropdown">
          <button class="menu-item" on:click={() => (showPlugins = true)} disabled={!projectPath}>Управление плагинами…</button>
          {#if installedPlugins.length > 0}
            <div class="menu-sep"></div>
            {#each installedPlugins as plugin (plugin.name)}
              {#if plugin.runnable}
                {#if plugin.vanessa}
                  <button class="menu-item" on:click={() => openVanessaDialog(true)} disabled={!projectPath}>Vanessa (dry)…</button>
                  <button class="menu-item" on:click={() => openVanessaDialog(false)} disabled={!projectPath}>Vanessa run…</button>
                  <button class="menu-item" on:click={() => openVanessaDialog(false, true)} disabled={!projectPath}>Vanessa rerun-failed…</button>
                {:else}
                  <button class="menu-item" on:click={() => openPluginRun(plugin.name, true)} disabled={!projectPath}>{pluginLabel(plugin)} (dry)…</button>
                  <button class="menu-item" on:click={() => openPluginRun(plugin.name, false)} disabled={!projectPath}>Запустить {pluginLabel(plugin)}…</button>
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
        <div class="menu-dropdown">
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
        <div class="menu-dropdown">
          <button class="menu-item" on:click={() => openStepsHelp()}>Справка по шагам…<span class="menu-shortcut">F1</span></button>
          <button class="menu-item" on:click={() => (showHotkeys = true)}>Горячие клавиши<span class="menu-shortcut">Shift+F1</span></button>
          <button class="menu-item" on:click={checkUpdates}>Проверить обновления…</button>
          <button class="menu-item" on:click={showAboutDialog}>О программе</button>
        </div>
      {/if}
    </div>
  </div>

  <div class="ide-main">
    <div class="ide-center" class:no-sidebar={!sidebarVisible} style="--sidebar-width: {layoutSidebarWidth}px">
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
        <div class="sidebar-column" style="width: {layoutSidebarWidth + 4}px">
        <aside class="explorer" style="width: {layoutSidebarWidth}px">
          <div class="explorer-header">
            <p class="zone-title">СЦЕНАРИИ</p>
            <div class="explorer-tools">
              <input class="explorer-search" value={sidebarSearch} placeholder="Поиск, @тег или tag:smoke" on:input={onSidebarSearchInput} />
              <div class="explorer-tool-actions">
                <button class="icon-btn" title="Новый сценарий" on:click={newScenario}>{@html icons.plus}</button>
                <button class="icon-btn batch-toggle" class:active={batchMode} title="Пакетный запуск: выбрать все видимые сценарии" on:click={toggleBatchMode}>
                  <span class="batch-toggle-icon" aria-hidden="true">{@html icons.validate}</span>
                  <span class="batch-label">Выбор</span>
                </button>
              </div>
            </div>
            {#if tags.length > 0}
              <div class="explorer-tag-chips">
                {#each tags.slice(0, 12) as tag}
                  <button
                    type="button"
                    class="chip"
                    class:active={sidebarSearch.trim() === tag || sidebarSearch.trim() === `@${tag.replace(/^@/, '')}`}
                    on:click={() => setSidebarSearch(tag.startsWith('@') ? tag : `@${tag}`)}
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
                batchSelectedSet={batchSelectedSet}
                {batchMode}
                indentStep={catalogIndent}
                compact={compactCatalogTree}
                expandAll={catalogViewState.expandAll}
                collapsed={catalogCollapsed}
                dropTarget={catalogDropTarget}
                onActivate={onCatalogActivate}
                onToggleBatch={onCatalogToggleBatch}
                onCollapseChange={onCatalogCollapse}
                onFileContextMenu={onFileContextMenu}
                onFolderContextMenu={onFolderContextMenu}
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
        <div class="action-bar" class:compact={actionBarCompact} use:observeActionBar>
          <div class="quick-toolbar">
            <div class="toolbar-row primary">
              <button class="tool-btn primary" title="Браузер (Ctrl+B)" on:click={() => void openBrowser()} disabled={!projectPath}>
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
              <button class="tool-btn primary" on:click={stopRecord} disabled={!recording && !browserOpen && !playing} title="Остановить запись, тест или браузер">
                {@html toolbarIcons.stop()}<span>Стоп</span>
              </button>
              <button
                class="tool-btn primary primary-run"
                class:run-active={playing}
                title={batchSelected.length > 0
                  ? `Запустить выбранные (${batchSelected.length})`
                  : 'Запустить (Ctrl+Enter)'}
                on:click={() => runPrimary(false)}
                disabled={(isWelcome && !projectPath && !batchSelected.length) || playing}
              >
                {@html toolbarIcons.play()}<span>{batchSelected.length > 0 ? `Запустить (${batchSelected.length})` : 'Запустить'}</span>
              </button>
              <button
                class="tool-btn"
                title="Текущий сценарий (Ctrl+Shift+Enter)"
                on:click={() => runCurrentScenario(false)}
                disabled={isWelcome || !activeTab}
              >
                {@html toolbarIcons.play()}<span>Сценарий</span>
              </button>
              <button class="tool-btn primary" on:click={saveFeature} disabled={isWelcome} title="Сохранить файл сценария (Ctrl+S)">
                {@html toolbarIcons.save()}<span>Сохранить</span>
              </button>
            </div>
            {#if !actionBarCompact}
            <div class="toolbar-row secondary" class:icon-only={toolbarIconOnly}>
              <button class="tool-btn" on:click={continueRecord} disabled={!projectPath || recording} title="Продолжить запись в конец сценария">
                {@html toolbarIcons.continueRecord()}<span>Дозапись</span>
              </button>
              <button class="tool-btn" on:click={toggleRecordPause} disabled={!recording} title="Приостановить запись">
                {@html toolbarIcons.pause()}<span>Пауза</span>
              </button>
              <span class="toolbar-sep" aria-hidden="true"></span>
              <button class="tool-btn" on:click={focusBrowser} disabled={!recording && !browserOpen && !playing}>
                {@html toolbarIcons.browserFocus()}<span>Показать браузер</span>
              </button>
              <button class="tool-btn" on:click={() => openValidateDialog(false)} disabled={!projectPath}>
                {@html toolbarIcons.validate()}<span>Селекторы на странице</span>
              </button>
              <button class="tool-btn" on:click={pickElement} disabled={(!recording && !browserOpen) || (recording && !recordPaused)} title={recording && !recordPaused ? 'Поставьте запись на паузу' : 'Указать элемент'}>
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
          tabLabel={featureTabLabel}
          tabUnsaved={tabIsUnsaved}
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
          {/if}
          <div class="feature-workspace" class:hidden={isWelcome}>
            {#if postRecordPath}
              <PostRecordBanner
                path={postRecordPath}
                stepCount={postRecordStepCount}
                onValidate={postRecordValidate}
                onSave={postRecordSave}
                onClose={dismissPostRecord}
              />
            {/if}
            {#if activeTabUnsaved}
              <div class="dirty-banner">
                <span>
                  {#if isUntitled(activeTab)}
                    Несохранённый сценарий
                  {:else}
                    Несохранённые изменения
                  {/if}
                </span>
                <button class="primary" on:click={saveFeature}>
                  {isUntitled(activeTab) ? 'Сохранить как…' : 'Сохранить'}
                </button>
              </div>
            {/if}
            {#if showVanessaMonitor}
              <VanessaMonitorPanel
                snapshot={vanessaSnapshot}
                running={vanessaRunning}
                onClose={() => (showVanessaMonitor = false)}
              />
            {/if}
            {#if stepStatusError}
              <div class="dirty-banner error">
                <span>В тексте сценария есть ошибки синтаксиса</span>
              </div>
            {/if}
            {#if showPlayingBar}
              <div class="playing-bar" role="status" aria-live="polite">
                <span class="play-label">▶ Выполняется:</span>
                <span class="play-target">{playingLabel}</span>
                {#if settingsSlowMo > 0 || lastRun.slowMo > 0}
                  <span class="play-slowmo">slow-mo {(lastRun.slowMo > 0 ? lastRun.slowMo : settingsSlowMo)} мс</span>
                {/if}
                <span class="play-progress" aria-hidden="true"></span>
              </div>
            {/if}
            {#if showRecordingBar}
              <div class="recording-bar">
                <span class="rec-label">Запись:</span>
                <label class="check-inline"><input type="checkbox" bind:checked={filterRecording} on:change={() => { if (filterRecording) navOnlyRecording = false; void syncRecordingOptions() }} /> Только важные</label>
                <label class="check-inline"><input type="checkbox" bind:checked={navOnlyRecording} on:change={() => { if (navOnlyRecording) filterRecording = false; void syncRecordingOptions() }} /> Только ссылки</label>
                <label class="check-inline"><input type="checkbox" bind:checked={settingsHeadless} on:change={() => void syncRecordingOptions()} /> Без окна браузера</label>
                <label class="check-inline"><input type="checkbox" bind:checked={hoverRecord} on:change={() => void syncRecordingOptions()} /> Записывать наведение</label>
              </div>
            {/if}
            <div class="gherkin-hints">
              <span class="hint-summary" title="{stepCount} шагов · Ctrl+Space — подсказки · Ctrl+Shift+O — структура · Ctrl+. — исправления">
                {stepCount} шагов · Ctrl+Space — подсказки · Ctrl+Shift+O — структура · Ctrl+. — исправления
              </span>
              <div class="hint-actions">
                <button on:click={openStepsDialog}>Шаблон</button>
                <button on:click={() => openStepsHelp()}>Справка</button>
                <button class:active={previewVisible} on:click={togglePreview}>Превью</button>
              </div>
            </div>
            <div class="editor-row" style="--preview-width: {layoutPreviewWidth}px">
              <div class="editor-main">
                <div class="editor-area" class:playing-active={playing}>
                    <MonacoEditor
                      bind:this={monaco}
                      bind:value={editorText}
                      bind:editorSettings
                      scenarioHints={editorScenarioHints}
                      hintActions={monacoHintActions}
                      runLensActions={monacoRunLensActions}
                      inlayHintsHandlers={monacoInlayHintsHandlers}
                      onHotkey={runHotkeyAction}
                      on:change={(e) => onEditorChange(e.detail)}
                      on:cursorline={(e) => (editorCursorLine = e.detail)}
                    />
                </div>
                {#if stepsPanelVisible}
                <div class="splitter-h" role="separator" on:mousedown={startResizeSteps}></div>
                <div class="steps-panel" class:collapsed={stepsPanelCollapsed} style="max-height: {stepsPanelCollapsed ? 24 : stepsPanelHeight}px">
                  <div class="steps-header">
                    <button on:click={() => (stepsPanelCollapsed = !stepsPanelCollapsed)}>
                      {#if stepsPanelCollapsed}{@html icons.chevronRight}{:else}{@html icons.chevronDown}{/if}
                    </button>
                    {#if editorSettings.symbolOutline}
                      <div class="steps-panel-tabs" role="tablist" aria-label="Вид панели шагов">
                        <button
                          type="button"
                          role="tab"
                          class:active={stepsPanelTab === 'outline'}
                          aria-selected={stepsPanelTab === 'outline'}
                          on:click={() => (stepsPanelTab = 'outline')}
                        >
                          Структура
                        </button>
                        <button
                          type="button"
                          role="tab"
                          class:active={stepsPanelTab === 'steps'}
                          aria-selected={stepsPanelTab === 'steps'}
                          on:click={() => (stepsPanelTab = 'steps')}
                        >
                          Шаги ({stepCount})
                        </button>
                      </div>
                    {:else}
                      <span>Шаги ({stepCount})</span>
                    {/if}
                    {#if stepStatusError}<span class="steps-header-error">ошибки</span>{/if}
                  </div>
                  {#if !stepsPanelCollapsed}
                    {#if editorSettings.symbolOutline && stepsPanelTab === 'outline'}
                      <div class="steps-outline-wrap">
                        <FeatureOutline
                          text={editorText}
                          currentLine={editorCursorLine}
                          onGoto={gotoEditorLine}
                        />
                      </div>
                    {:else}
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
                              on:contextmenu|preventDefault={(e) => openStepsContextMenu(e, step)}
                              title={step.text ? 'ПКМ или двойной клик — действия со шагом' : ''}
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
                  {/if}
                </div>
                {/if}
              </div>
              {#if showPreviewPane}
                <div class="splitter-v" role="separator" on:mousedown={startResizePreview}></div>
                <div class="feature-preview-pane" style="width: {layoutPreviewWidth}px">
                  <div class="preview-header">Превью Gherkin</div>
                  <FeaturePreview
                    text={editorText}
                    theme={editorSettings.theme}
                    fontSize={editorSettings.fontSize}
                    fontFamily={editorSettings.fontFamily}
                  />
                </div>
              {/if}
            </div>
          </div>
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
        <ResultsPanel
          entries={runResults}
          artifacts={projectArtifacts}
          onRerun={rerunFailed}
          onOpenFolder={openArtifactPath}
          onServeAllure={serveAllureReport}
          onOpenHtmlReport={openHtmlReport}
        />
      {:else if bottomTab === 'validate'}
        <ValidatePanel
          issues={validatePanelDisplayIssues()}
          hint={validateProjectHint()}
          cliLog={validateCliLog}
          activeLine={editorCursorLine}
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
        <span class="led" class:recording={recording} class:playing={playing} class:on={recording || browserOpen || playing}></span>
        Браузер · {recording ? 'запись' : playing ? 'тест' : browserOpen ? 'открыт' : 'закрыт'}
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

{#if showRun}
  <RunDialog
    title={runDialogTitle}
    bind:form={runForm}
    {testClients}
    {tags}
    scenarios={runDialogScenarios}
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
    bind:platformExe={vanessaPlatformExe}
    bind:epfPath={vanessaEpfPath}
    bind:ibConnection={vanessaIB}
    bind:reportAllure={vanessaReportAllure}
    bind:vaDir={vanessaVaDir}
    bind:vaFiles={vanessaVaFiles}
    {tags}
    scenarios={vanessaDialogScenarios}
    onConfirm={confirmVanessaRun}
    onCancel={() => (showVanessaRun = false)}
  />
{/if}

{#if showPluginRun}
  <PluginRunDialog
    pluginName={pluginRunName}
    pluginTitle={pluginRunTitle(pluginRunName)}
    bind:tag={pluginRunTag}
    bind:scenario={pluginRunScenario}
    bind:dryRun={pluginRunDry}
    scenarios={pluginRunScenarios}
    {tags}
    onConfirm={confirmPluginRun}
    onCancel={() => (showPluginRun = false)}
  />
{/if}

{#if showTestClient}
  <TestClientDialog
    {testClients}
    bind:selectedName={testClientSelection}
    browserOpen={browserOpen || recording}
    suggestName={testClientSuggestName}
    onUse={useTestClient}
    onClose={() => {
      showTestClient = false
      testClientSuggestName = ''
    }}
    onClientsChange={(names) => (testClients = names)}
    onLog={appendLog}
    onAskConfirm={(message) =>
      askConfirm({ title: 'Подтверждение', message, confirmLabel: 'Удалить', danger: true })}
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

{#if showRefactorUrl}
  <RefactorUrlDialog
    initialUrl={startURL || recordURL || 'https://example.com'}
    onConfirm={applyRefactorUrl}
    onClose={() => (showRefactorUrl = false)}
  />
{/if}

{#if showOpenProject}
  <OpenProjectDialog
    initialPath={projectPath}
    {recentProjects}
    onConfirm={openProjectAt}
    onClose={() => (showOpenProject = false)}
  />
{/if}

{#if showRenameFeature}
  <RenameFeatureDialog
    currentPath={renameFeaturePath}
    onConfirm={(name) => renameFeature(renameFeaturePath, name)}
    onClose={() => {
      showRenameFeature = false
      renameFeaturePath = ''
    }}
  />
{/if}

{#if showMoveFeature}
  <MoveFeatureDialog
    featurePath={moveFeaturePath}
    destDirs={moveDestDirs}
    bind:destDir={moveDestDir}
    onConfirm={confirmMoveFeature}
    onCancel={() => {
      showMoveFeature = false
      moveFeaturePath = ''
    }}
  />
{/if}

{#if showValidate}
  <ValidateDialog
    bind:browser={validateBrowser}
    bind:syntaxOnly={validateSyntaxOnly}
    bind:scope={validateScope}
    canValidateCurrent={!isWelcome && !!activeTab}
    currentFileName={!isWelcome && activeTab ? basename(activeTab) : ''}
    onConfirm={confirmValidate}
    onCancel={() => (showValidate = false)}
  />
{/if}

{#if showInitProject}
  <InitProjectDialog
    {projectPath}
    onConfirm={confirmInitProject}
    onCancel={() => (showInitProject = false)}
  />
{/if}

{#if showUpdateCheck}
  <UpdateCheckDialog
    currentVersion={version}
    info={updateCheckInfo}
    message={updateCheckMessage}
    hasUpdate={updateCheckHasUpdate}
    downloading={updateDownloading}
    progress={updateProgress}
    onClose={() => (showUpdateCheck = false)}
    onOpenRelease={openUpdateRelease}
    onDownload={downloadUpdate}
    onApply={applyUpdate}
    canAutoApply={updateCheckInfo?.canAutoApply ?? false}
  />
{/if}

{#if showDuplicateFeature}
  <DuplicateFeatureDialog
    featurePath={duplicateFeaturePath}
    bind:newName={duplicateNewName}
    onConfirm={confirmDuplicateFeature}
    onCancel={() => {
      showDuplicateFeature = false
      duplicateFeaturePath = ''
    }}
  />
{/if}

{#if showImportFeatures}
  <ImportFeaturesDialog
    destDirs={collectProjectDirs()}
    bind:destDir={importDestDir}
    busy={importFeaturesBusy}
    onImport={confirmImportFeatures}
    onClose={() => (showImportFeatures = false)}
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
    bind:slowMo={settingsSlowMo}
    bind:loops={settingsLoops}
    bind:filterRecording
    bind:navOnlyRecording
    bind:hoverRecord
    bind:scrollBeforeClick={settingsScrollBeforeClick}
    bind:hoverRecordMinMs={settingsHoverRecordMinMs}
    bind:toolbarCompact
    bind:stepsPanelVisible
    bind:stepsPanelHeight
    bind:checkUpdatesOnStartup={settingsCheckUpdatesOnStartup}
    bind:selectorClickStrategies={settingsSelectorClickStrategies}
    bind:selectorInputStrategies={settingsSelectorInputStrategies}
    bind:editorSettings
    onSave={applySettings}
    onCancel={cancelSettings}
    onOpenPlugins={() => {
      showSettings = false
      showPlugins = true
    }}
    onOpenVanessa={() => {
      showSettings = false
      openVanessaSettingsDialog()
    }}
    onInstallLog={appendLog}
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
    bind:mode={recordMode}
    bind:url={recordURL}
    bind:output={recordOutput}
    bind:featureName={recordFeatureName}
    bind:scenarioName={recordScenarioName}
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
    {baselineBusy}
    onHttpAuth={openHttpAuthDialog}
    onStart={startRecord}
    onTogglePause={toggleRecordPause}
    onStop={stopRecord}
    onSaveBaseline={saveBaselineRecord}
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
    onRunPlugin={(name, dry) => openPluginRun(name, dry)}
    onAskConfirm={(message) =>
      askConfirm({ title: 'Подтверждение', message, confirmLabel: 'Удалить', danger: true })}
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

{#if folderMenu}
  <FolderContextMenu
    x={folderMenu.x}
    y={folderMenu.y}
    featureCount={folderMenu.paths.length}
    onRun={() => folderMenuRun(false)}
    onDryRun={() => folderMenuRun(true)}
    onVanessa={() => folderMenuVanessa(false)}
    onVanessaDry={() => folderMenuVanessa(true)}
    onSelectBatch={folderMenuSelectBatch}
    onClose={dismissFolderMenu}
  />
{/if}

{#if contextMenu}
  <CatalogContextMenu
    x={contextMenu.x}
    y={contextMenu.y}
    onRun={contextMenuRun}
    onDryRun={contextMenuDryRun}
    onOpen={contextMenuOpen}
    onDuplicate={contextMenuDuplicate}
    onRename={contextMenuRename}
    onMove={contextMenuMove}
    onReveal={contextMenuReveal}
    onDelete={contextMenuDelete}
    onClose={dismissContextMenu}
  />
{/if}

{#if stepsMenu}
  <StepsContextMenu
    x={stepsMenu.x}
    y={stepsMenu.y}
    onRunFrom={() => stepsMenuRunFrom(false)}
    onRunTo={() => stepsMenuRunTo(false)}
    onDryRunFrom={() => stepsMenuRunFrom(true)}
    onDryRunTo={() => stepsMenuRunTo(true)}
    onGoto={stepsMenuGoto}
    onHelp={stepsMenuHelp}
    onClose={closeStepsMenu}
  />
{/if}

{#if confirmDialog}
  <ConfirmDialog
    title={confirmDialog.title}
    message={confirmDialog.message}
    confirmLabel={confirmDialog.confirmLabel}
    danger={confirmDialog.danger}
    onConfirm={() => closeConfirm(true)}
    onCancel={() => closeConfirm(false)}
  />
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
{/if}
