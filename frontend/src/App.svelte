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
  import NewProjectWizardDialog, { type NewProjectWizardResult } from './lib/NewProjectWizardDialog.svelte'
  import ImportFeaturesDialog from './lib/ImportFeaturesDialog.svelte'
  import DuplicateFeatureDialog from './lib/DuplicateFeatureDialog.svelte'
  import RenameFeatureDialog from './lib/RenameFeatureDialog.svelte'
  import { buildFeatureTemplate } from './lib/featureTemplate'
  import { upsertRecordedStepInText, removeLastRecordedStepFromText } from './lib/recordedStepEditor'
  import { isUntitled, makeUntitledPath, syncUntitledCounterFromPaths, untitledLabel } from './lib/untitled'
  import {
    buildSessionTabsSnapshot,
    sessionTabPathsFromSettings,
    untitledContentMap,
  } from './lib/sessionTabs'
  import { matchHotkey, monacoOverlayConsumesEscape, shouldIgnoreAppHotkey, type HotkeyId } from './lib/hotkeys'
  import { defaultRunForm, type RunForm } from './lib/runTypes'
  import { formatLastRunSummary } from './lib/runSummary'
  import { scenarioAtLine, listScenarioTitles, mergeScenarioNames } from './lib/scenarioAtLine'
  import {
    buildSyntheticRunError,
    filterRunResultsSince,
    pickLastRunError,
    remapRunResultPaths,
    resolveStaleRunScenario,
  } from './lib/runResults'
  import {
    materializeRunTargetPaths,
    resolveLogicalRunTarget,
  } from './lib/runTargets'
  import PostRecordBanner from './lib/PostRecordBanner.svelte'
  import PostRecordDiffDialog from './lib/PostRecordDiffDialog.svelte'
  import type { HintActionHandlers } from './lib/gherkinHintActions'
  import RunHistoryDialog from './lib/RunHistoryDialog.svelte'
  import StepsHelpDialog from './lib/StepsHelpDialog.svelte'
  import UnsavedCloseDialog from './lib/UnsavedCloseDialog.svelte'
  import HttpAuthDialog from './lib/HttpAuthDialog.svelte'
  import PickerStepDialog from './lib/PickerStepDialog.svelte'
  import OnboardingTour from './lib/onboarding/OnboardingTour.svelte'
  import { ONBOARDING_TOUR_VERSION } from './lib/onboarding/tourSteps'
  import type { TourContext } from './lib/onboarding/tourState'
  import { createTranslator, locale, setLocale, t, type Locale } from './lib/i18n'
  import { loadLayout, saveLayout, resetLayout as resetUILayout } from './lib/layout'
  import { isLargeFeatureFile, LARGE_FILE_LINE_THRESHOLD } from './lib/editorLargeFile'
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
  import { prefetchMonacoEditor } from './lib/appBootstrap'
  import { setStepHoverEnabled } from './lib/gherkinStepHover'
  import { filterScenarioHints, applyAutoFixableScenarioHints } from './lib/scenarioHints'
  import {
    DEFAULT_EDITOR_SETTINGS,
    editorSettingsFromDTO,
    editorSettingsToDTO,
    type EditorSettings,
  } from './lib/editorOptions'
  import { resolveRecordStartURL } from './lib/recordStartUrl'
  import { isSameRecordTab, normalizeRecordTabPath, recordingTabSwitchAllowed } from './lib/recordingTarget'
  import { flakyScenarioMap, flakyStepHints } from './lib/flakyMetrics'
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
    CancelRun,
    Validate,
    ValidateFeature,
    ListTestClients,
    InitProject,
    InitProjectAt,
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
    FlakyMetrics,
    BundledExamplesPath,
    ProjectArtifacts,
    ScenariaArtifactPath,
    ParseEditorSteps,
    ArtifactExists,
    OpenFolder,
    ServeAllure,
    AllureStatus,
    OpenHTMLReport,
    OpenTrace,
    FailedStepLine,
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
  let stepStatusError = false
  let testClients: string[] = []
  let monaco: MonacoEditor | undefined

  let appReady = false
  let splashMessage = t('splash.starting')
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
  let pendingUpdateCheckOnStartup = false
  let updateCheckMessage = ''
  let updateCheckHasUpdate = false
  let updateCheckInfo: gui.UpdateInfoDTO | null = null
  let updateDownloading = false
  let updateProgress: gui.UpdateProgressDTO | null = null
  let showInitProject = false
  let showNewProjectWizard = false
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
    dontAskAgainLabel?: string
    resolve: (confirmed: boolean, dontAskAgain?: boolean) => void
  } | null = null

  let skipRecordTabSwitchConfirm = false
  let runningDryRun = false

  function askConfirm(opts: {
    title: string
    message: string
    confirmLabel?: string
    danger?: boolean
    dontAskAgainLabel?: string
  }): Promise<boolean> {
    return new Promise((resolve) => {
      confirmDialog = {
        title: opts.title,
        message: opts.message,
        confirmLabel: opts.confirmLabel || tr('common.ok'),
        danger: opts.danger || false,
        dontAskAgainLabel: opts.dontAskAgainLabel,
        resolve: (confirmed, dontAskAgain) => {
          if (confirmed && dontAskAgain) skipRecordTabSwitchConfirm = true
          resolve(confirmed)
        },
      }
    })
  }

  function closeConfirm(confirmed: boolean, dontAskAgain = false) {
    if (confirmDialog) {
      confirmDialog.resolve(confirmed, dontAskAgain)
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
  let postRecordBaselineText = ''
  let showPostRecordDiff = false
  let recordStepPickerOpen = false
  let flakyMetrics: gui.FlakyMetricsDTO | null = null
  let editorScenarioHints: gui.ScenarioHintDTO[] = []
  let editorHintsDismissed = new Set<string>()
  let contextMenu: { x: number; y: number; path: string } | null = null
  let folderMenu: { x: number; y: number; dir: string; paths: string[] } | null = null
  let runDialogTitle = ''
  let runDialogScenarios: string[] = []
  let vanessaDialogScenarios: string[] = []
  let pluginRunScenarios: string[] = []

  let recording = false
  let browserOpen = false
  let recordPaused = false
  let playing = false
  let playingLabel = ''
  let runProgressCurrent = 0
  let runProgressTotal = 0
  let runLogStreaming = false
  let runCancelling = false
  let stepsMenu: { x: number; y: number; line: number; step: gui.EditorStepRow } | null = null
  let sessionPersistTimer: ReturnType<typeof setTimeout> | null = null
  let draftAutosaveTimer: ReturnType<typeof setInterval> | null = null

  let sidebarVisible = true
  let sidebarWidth = 260
  let previewVisible = false
  let previewPaneMounted = false
  let previewMountTimer: ReturnType<typeof setTimeout> | null = null
  let previewWidth = 360
  let bottomPanelOpen = false
  let bottomTab: 'journal' | 'results' | 'validate' | 'error' = 'journal'
  let stepsPanelCollapsed = true
  let stepsPanelHeight = 160
  let bottomPanelHeight = 200
  let sidebarSearch = ''
  let catalogFilterText = ''
  let openMenu: string | null = null
  let statusMessage = ''
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

  let lastRun: RunForm = defaultRunForm({ headed: true, installPW: true, html: true, htmlLightMode: true })
  let runForm: RunForm = { ...lastRun }

  let settingsBrowser = 'chromium'
  let settingsHeadless = false
  let settingsWorkers = 1
  let settingsSlowMo = 0
  let settingsScrollBeforeClick = false
  let settingsHoverRecordMinMs = 600
  let settingsLoops = 100
  let settingsCheckUpdatesOnStartup = true
  let settingsSelectorClickStrategies: string[] = ['text', 'contextual', 'aria', 'title', 'testid', 'id']
  let settingsSelectorInputStrategies: string[] = ['label', 'placeholder', 'aria', 'name', 'testid', 'id']
  let settingsNavWaitUntil = 'domcontentloaded'
  let editorSettings: EditorSettings = { ...DEFAULT_EDITOR_SETTINGS }
  let editorCursorLine = 1
  let stepsPanelTab: 'outline' | 'steps' = DEFAULT_EDITOR_SETTINGS.stepsPanelView

  let recordURL = ''
  let startURL = ''
  let recordOutput = 'recorded.feature'
  let recordIdle = 30
  let recordAppendTo = ''
  let recordTestClient = ''
  let recordFeatureName = t('dialogs.record.featureDefault')
  let recordScenarioName = t('dialogs.record.scenarioDefault')
  let lastRecordTarget = ''
  let liveRecordStepLines: Record<number, number> = {}
  let recordStepApplyChain: Promise<void> = Promise.resolve()
  let recordEditorReadyPromise: Promise<void> = Promise.resolve()
  let recordingTargetPath = ''
  let pauseToggleGuardUntil = 0
  let pendingCloseTab: string | null = null

  let otpEmail = ''

  let testClientSelection = ''
  let testClientSuggestName = ''

  let recentProjects: string[] = []
  let recentFeatures: string[] = []
  let runResults: gui.RunResultEntry[] = []
  let lastErrorEntry: gui.RunResultEntry | null = null
  let lastRunSince: string | null = null
  let lastRunBatchResults: gui.RunResultEntry[] = []
  let lastRunSummary = ''
  let paletteCommands: PaletteCommand[] = []
  $: flakyByPath = flakyScenarioMap(flakyMetrics)
  $: flakyStepByPath = flakyStepHints(flakyMetrics)
  let editorSteps: EditorStepRow[] = []
  let editorValidationIssues: gui.ValidationIssue[] = []
  let validatePanelIssues: gui.ValidationIssue[] = []
  let projectArtifacts: gui.ProjectArtifacts = new gui.ProjectArtifacts()

  const unsubscribers: (() => void)[] = []

  $: tr = createTranslator($locale)

  $: if (pendingUpdateCheckOnStartup && settingsCheckUpdatesOnStartup && projectPath) {
    pendingUpdateCheckOnStartup = false
    void checkUpdatesOnStartup()
  }

  $: isWelcome = activeTab === WELCOME_KEY
  $: activeFeatureTab = tabs.find((t) => t.path === activeTab)
  $: activeTabUnsaved = activeFeatureTab ? tabIsUnsaved(activeFeatureTab) : false
  $: if (editorSettings.stepsPanelView) {
    stepsPanelTab = editorSettings.stepsPanelView
  }
  $: stepCount = editorSteps.length
  $: editorLineCount = isWelcome ? 0 : editorText.split(/\r?\n/).length
  $: showLargeFileBanner = !isWelcome && isLargeFeatureFile(editorLineCount)
  $: unsavedTabCount = tabs.filter((t) => tabIsUnsaved(t)).length
  $: {
    $locale
    lastRunSummary = formatLastRunSummary(lastRun)
  }
  $: automationActive = playing || vanessaRunning
  $: pickerToolbarEnabled =
    pickerDuringRecording
      ? browserOpen && !playing
      : (recording && recordPaused) || (browserOpen && !recording && !playing)
  $: batchCount = batchSelected.length
  $: batchSelectedSet = buildBatchSelectedSet(batchSelected)
  $: showRecordingBar = recording && !showRecord
  $: showPlayingBar = playing
  $: anyAppDialogOpen =
    !!confirmDialog ||
    !!pendingCloseTab ||
    showRun ||
    showVanessaRun ||
    showPluginRun ||
    showTestClient ||
    showStepsHelp ||
    showSteps ||
    showVanessaSettings ||
    showExport ||
    showRefactorUrl ||
    showOpenProject ||
    showRenameFeature ||
    showMoveFeature ||
    showValidate ||
    showNewProjectWizard ||
    showInitProject ||
    showUpdateCheck ||
    showDuplicateFeature ||
    showImportFeatures ||
    showImport ||
    showSettings ||
    showCommandPalette ||
    showSnippetPalette ||
    showRecord ||
    showOtp ||
    showAbout ||
    showHotkeys ||
    showPlugins ||
    showRunHistory ||
    (showPostRecordDiff && !!postRecordPath) ||
    showProjectReplace ||
    showHttpAuth ||
    showPickerStep
  $: onboardingTourContext = {
    projectPath,
    featuresCount: features.length,
    isWelcome,
    activeTab,
    welcomeKey: WELCOME_KEY,
    validateDone: onboardingValidateDone,
    dryRunDone: onboardingDryRunDone,
    journalVisited: onboardingJournalVisited,
    bottomPanelOpen,
    bottomTab,
  } satisfies TourContext
  $: onboardingTourActive = showOnboardingTour && !anyAppDialogOpen
  $: if (typeof document !== 'undefined') {
    document.body.classList.toggle('onboarding-tour-active', onboardingTourActive)
  }
  $: if (showOnboardingTour && isWelcome && projectPath && features.length > 0) {
    onboardingValidateDone = false
    onboardingDryRunDone = false
    onboardingJournalVisited = false
  }
  let onboardingTourStepId = 'welcome'
  $: onboardingElevateMenubar =
    showOnboardingTour && (onboardingTourStepId === 'validate' || onboardingTourStepId === 'dry-run')
  $: onboardingElevateSidebar = showOnboardingTour && onboardingTourStepId === 'pick-feature'
  $: onboardingElevateWelcome =
    showOnboardingTour && (onboardingTourStepId === 'welcome' || onboardingTourStepId === 'open-examples')
  $: onboardingElevateBottom = showOnboardingTour && onboardingTourStepId === 'journal'
  $: if (showOnboardingTour && (onboardingTourStepId === 'validate' || onboardingTourStepId === 'dry-run')) {
    openMenu = 'run'
  }
  $: if (showOnboardingTour && openMenu === 'run' && (onboardingTourStepId === 'validate' || onboardingTourStepId === 'dry-run')) {
    void tick().then(() => onboardingTour?.relayout())
  }
  $: showBrowserOverlay = (browserOpen || recording || playing) && !anyAppDialogOpen
  $: {
    $locale
    paletteCommands = buildPaletteCommands()
  }
  $: stepStatusDisplay =
    stepCount === 0 && !stepStatusError
      ? tr('statusBar.steps', { count: 0 })
      : stepStatusError
        ? tr('statusBar.stepsWithErrors', { count: stepCount, errors: editorValidationIssues.length })
        : tr('statusBar.steps', { count: stepCount })

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
  $: {
    if (previewMountTimer) {
      clearTimeout(previewMountTimer)
      previewMountTimer = null
    }
    if (showPreviewPane) {
      previewMountTimer = setTimeout(() => {
        previewPaneMounted = true
        previewMountTimer = null
      }, 80)
    } else {
      previewPaneMounted = false
    }
  }

  $: activeFeaturePath = activeTab !== WELCOME_KEY ? activeTab : ''

  $: resultsPanelEntries =
    lastRunBatchResults.length > 0
      ? lastRunBatchResults
      : lastRunSince
        ? filterRunResultsSince(runResults, lastRunSince)
        : runResults

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
  $: welcomeRecorded = recording || browserOpen
  let welcomePlayedSuccess = false
  let checklistDismissed = false
  let onboardingCompleted = false
  let onboardingDismissed = false
  let onboardingValidateDone = false
  let onboardingDryRunDone = false
  let onboardingJournalVisited = false
  let showOnboardingTour = false
  let onboardingTour: OnboardingTour
  let runDialogConfirmed = false
  let pickerDuringRecording = false
  let uiLocale: Locale = 'ru'
  let allureInstalled = true
  let allureServeRunning = false

  $: stopActionLabel = playing
    ? tr('toolbar.stopTest')
    : recording
      ? tr('toolbar.stopRecord')
      : browserOpen
        ? tr('toolbar.closeBrowser')
        : tr('toolbar.stop')
  $: recordingTargetLabel =
    recording && recordingTargetPath ? basename(recordingTargetPath) : ''

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
    if (shouldAutoStartOnboarding()) {
      sidebarVisible = true
      onboardingTourStepId = 'welcome'
      showOnboardingTour = true
    }
  }

  function shouldAutoStartOnboarding(): boolean {
    return !onboardingCompleted && !onboardingDismissed
  }

  async function completeOnboarding() {
    showOnboardingTour = false
    onboardingCompleted = true
    onboardingDismissed = false
    await persistSettings()
    void maybeCheckUpdatesOnStartup()
  }

  async function dismissOnboarding() {
    showOnboardingTour = false
    onboardingDismissed = true
    await persistSettings()
    void maybeCheckUpdatesOnStartup()
  }

  function openJournalTab(markTourVisit = false) {
    bottomPanelOpen = true
    bottomTab = 'journal'
    if (markTourVisit && showOnboardingTour) onboardingJournalVisited = true
  }

  function restartOnboardingTour() {
    onboardingValidateDone = false
    onboardingDryRunDone = false
    onboardingJournalVisited = false
    onboardingTourStepId = ''
    sidebarVisible = true
    saveLayout({ sidebarVisible: true })
    void resetWorkspaceForOnboarding().then(() => {
      onboardingTourStepId = 'welcome'
      showOnboardingTour = true
      onboardingTour?.restart()
    })
  }

  async function resetWorkspaceForOnboarding() {
    if (projectPath) {
      try {
        await teardownDesktopSession()
      } catch {
        /* offline */
      }
      projectPath = ''
      features = []
      tags = []
      featureTags = {}
    }
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
    clearEditorValidation()
    openMenu = ''
    bottomPanelOpen = false
    saveLayout({ bottomPanelOpen: false })
  }

  function onOnboardingStepChange(e: CustomEvent<{ index: number; id: string }>) {
    const { id } = e.detail
    onboardingTourStepId = id
    if (id === 'open-examples' || id === 'pick-feature') {
      sidebarVisible = true
      saveLayout({ sidebarVisible: true })
    }
    if (id === 'validate' || id === 'dry-run') {
      void tick().then(() => tick().then(() => onboardingTour?.relayout()))
    }
    if (id === 'journal') {
      bottomPanelOpen = true
      bottomTab = 'journal'
      saveLayout({ bottomPanelOpen: true })
      void tick().then(() => onboardingTour?.relayout())
    }
  }

  onMount(async () => {
    const startedAt = Date.now()
    let splashFinished = false
    const finishStartup = async () => {
      if (splashFinished) return
      splashFinished = true
      setSplashStage(tr('splash.ready'), 100)
      await dismissSplash(startedAt)
      applyDevUiMock()
      syncIdleStatus()
      prefetchMonacoEditor()
      void maybeCheckUpdatesOnStartup()
    }
    const startupGuard = window.setTimeout(() => {
      console.error('Startup guard: forcing splash dismiss')
      void finishStartup()
    }, 12_000)

    try {
    setSplashDocumentState(true)
    await beginSplashWindow()

    setSplashStage(tr('splash.envSetup'), 8)

    const layout = loadLayout()
    sidebarVisible = layout.sidebarVisible
    bottomPanelOpen = layout.bottomPanelOpen
    bottomPanelHeight = layout.bottomPanelHeight
    previewVisible = layout.previewVisible
    previewWidth = layout.previewWidth || 360

    setSplashStage(tr('splash.connecting'), 40)

    try {
      version = await Version()
    } catch {
      version = 'dev'
    }

    setSplashStage(tr('splash.loadingSettings'), 60)

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
      if (!shouldAutoStartOnboarding()) {
        await restoreWorkspaceSession(settings)
      }
    }

    draftAutosaveTimer = setInterval(() => void autosaveDirtyDrafts(), 30_000)

    setSplashStage(tr('splash.initializing'), 88)

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
          setStatus(tr('journal.status.browserOpen'), 'busy')
          appendLog(tr('journal.browser.opened'))
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
          const wasRecording = recording
          applyBrowserSessionState({ browserOpen: true, recording: true, paused: false })
          showRecord = false
          setStatus(tr('journal.status.preparingRecord'), 'busy')
          startBrowserWatch()
          recordEditorReadyPromise = (async () => {
            if (!syncOnly) {
              if (!appendOnly && !wasRecording) {
                liveRecordStepLines = {}
              }
              if (!appendOnly) {
                await prepareRecordEditorTab(meta.output || '')
              }
              postRecordBaselineText = monaco?.getEditorText() ?? editorText
            }
          })()
          try {
            await recordEditorReadyPromise
          } catch (e: any) {
            appendLog(tr('journal.record.prepTabError', { error: String(e) }))
          }
          if (!syncOnly && activeTab && !isWelcome) {
            recordingTargetPath = normalizeRecordTabPath(activeTab)
          }
          if (!appendOnly && !syncOnly) {
            appendLog(tr('journal.record.started'))
            setStatus(tr('journal.status.recording'), 'busy')
          } else if (syncOnly) {
            setStatus(tr('journal.status.recording'), 'busy')
          }
        }),
      )
      unsubscribers.push(
        EventsOn('record-stopped', (payload: { reason?: string; idleSeconds?: number } | null) => {
          handleRecordStopped(payload ?? undefined)
        }),
      )
      unsubscribers.push(
        EventsOn('run-log-line', (payload: { line?: string } | string) => {
          if (!runLogStreaming) return
          const line = typeof payload === 'string' ? payload : payload?.line
          if (line) appendLog(line)
        }),
      )
      unsubscribers.push(
        EventsOn('run-progress', (payload: {
          phase?: string
          index?: number
          total?: number
          featurePath?: string
          scenario?: string
          success?: boolean
        }) => {
          if (!playing || !payload) return
          const total = payload.total ?? runProgressTotal
          const index = payload.index ?? runProgressCurrent
          if (total > 0) {
            runProgressTotal = total
            runProgressCurrent = index
          }
          const name = payload.scenario || (payload.featurePath ? basename(payload.featurePath) : '')
          if (name && total > 0) {
            playingLabel = `${name} (${index}/${total})`
          }
          if (payload.phase === 'scenario_done') {
            void refreshRunResults()
          }
        }),
      )
      unsubscribers.push(
        EventsOn('run-results-changed', () => {
          void refreshRunResults()
        }),
      )
      unsubscribers.push(
        EventsOn('report-goto', (req: { feature_path: string; scenario: string; leaf_index: number; line: number }) => {
          if (req) void gotoReportStep(req)
        }),
      )
      unsubscribers.push(
        EventsOn('report-rerun', (req: { feature_path: string; scenario: string }) => {
          if (req) void rerunFromReport(req)
        }),
      )
      unsubscribers.push(
        EventsOn('report-trace', (req: { trace_path: string; report_dir: string; trace_offset_ms?: number; step_index?: number }) => {
          if (req) void openTraceFromReport(req)
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
        EventsOn('record-error', (message: string) => {
          stopBrowserWatch()
          recording = false
          recordPaused = false
          browserOpen = false
          showRecord = false
          liveRecordStepLines = {}
          recordingTargetPath = ''
          appendLog(tr('journal.record.error', { message: message || tr('journal.record.unknownError') }))
          setStatus(tr('journal.status.recordError'), 'error')
          syncIdleStatus()
        }),
      )
      unsubscribers.push(
        EventsOn('vanessa-run-started', () => {
          vanessaRunning = true
          vanessaWatchDir = ''
          setStatus(tr('journal.status.vanessaRunning'), 'busy')
          appendLog(tr('journal.vanessa.starting'))
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
            appendLog(tr('journal.error.generic', { error: result.error }))
            setStatus(tr('journal.status.vanessaError'), 'error')
          } else {
            appendLog(tr('journal.vanessa.done'))
            setStatus(tr('journal.status.vanessaDone'), result.success ? 'success' : 'error')
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
      if (showOnboardingTour && (onboardingTourStepId === 'validate' || onboardingTourStepId === 'dry-run')) {
        return
      }
      openMenu = null
    }
    window.addEventListener('keydown', onModalEscapeCapture, { capture: true })
    window.addEventListener('keydown', onGlobalKeydown, { capture: true })
    window.addEventListener('click', onDocClick)
    const onVisibility = () => {
      if (document.visibilityState === 'visible') void checkActiveTabDiskStale()
    }
    document.addEventListener('visibilitychange', onVisibility)
    unsubscribers.push(() => window.removeEventListener('keydown', onModalEscapeCapture, { capture: true }))
    unsubscribers.push(() => window.removeEventListener('keydown', onGlobalKeydown, { capture: true }))
    unsubscribers.push(() => window.removeEventListener('click', onDocClick))
    unsubscribers.push(() => document.removeEventListener('visibilitychange', onVisibility))
    } catch (err) {
      console.error('Startup failed', err)
    } finally {
      window.clearTimeout(startupGuard)
      await finishStartup()
    }
  })

  onDestroy(() => {
    if (typeof document !== 'undefined') {
      document.body.classList.remove('onboarding-tour-active')
    }
    stopVanessaPoll()
    stopBrowserWatch()
    void teardownDesktopSession()
    applyCatalogFilterDebounced.cancel()
    if (sessionPersistTimer) {
      clearTimeout(sessionPersistTimer)
      sessionPersistTimer = null
      void persistSettings()
    }
    if (validateDebounceTimer) clearTimeout(validateDebounceTimer)
    if (previewMountTimer) clearTimeout(previewMountTimer)
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
      const info = await OpenProject(proj)
      projectPath = info.path
      features = info.features || []
      tags = info.tags || []
      featureTags = info.featureTags || {}
      testClients = await ListTestClients().catch(() => [])
    } catch {
      appendLog(tr('journal.session.projectNotFound', { path: proj }))
      setStatus(tr('journal.status.sessionProjectNotFound'), 'error')
      return
    }
    try {
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

  function runModeLabel(dryRun: boolean): string {
    return dryRun ? tr('journal.run.mode.dryRun') : tr('journal.run.mode.run')
  }

  function partialRunLogSuffix(startStep: number, endStep: number): string {
    if (startStep >= 0 && endStep >= 0) return ` (${startStep + 1}–${endStep + 1})`
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
      : ['text', 'contextual', 'aria', 'title', 'testid', 'id']
    settingsSelectorInputStrategies = s.selectorInputStrategies?.length
      ? [...s.selectorInputStrategies]
      : ['testid', 'id', 'label', 'placeholder', 'aria', 'name']
    settingsNavWaitUntil = s.navWaitUntil || 'domcontentloaded'
    editorSettings = editorSettingsFromDTO(s.editor)
    stepsPanelTab = editorSettings.stepsPanelView
    stepsPanelCollapsed = resolveStepsPanelCollapsed()
    lastRun = {
      ...lastRun,
      workers: s.parallelWorkers || 1,
      slowMo: s.slowMo ?? 0,
      browser: s.browser || 'chromium',
    }
    checklistDismissed = !!s.checklistDismissed
    welcomePlayedSuccess = !!s.welcomePlayedSuccess
    onboardingCompleted = !!s.onboardingCompleted
    onboardingDismissed = !!s.onboardingDismissed
    if (!onboardingCompleted && s.onboardingVersion !== ONBOARDING_TOUR_VERSION) {
      onboardingCompleted = false
      onboardingDismissed = false
    }
    runDialogConfirmed = !!s.runDialogConfirmed
    pickerDuringRecording = !!s.pickerDuringRecording
    if (s.startUrl) startURL = s.startUrl
    if (s.uiLocale === 'en' || s.uiLocale === 'ru') {
      uiLocale = s.uiLocale
      setLocale(s.uiLocale)
    }
  }

  function resolveStepsPanelCollapsed(): boolean {
    return !stepsPanelVisible
  }

  function buildPaletteCommands(): PaletteCommand[] {
    const pg = (key: string) => tr(`palette.groups.${key}`)
    const pc = (key: string) => tr(`palette.commands.${key}`)
    const compactLabel = toolbarCompact ? pc('expanded') : pc('compact')
    return [
      { id: 'palette', label: pc('palette'), group: pg('view'), shortcut: 'Ctrl+Shift+P', run: () => (showCommandPalette = true) },
      { id: 'welcome', label: pc('welcome'), group: pg('view'), run: () => selectTab(WELCOME_KEY) },
      { id: 'open', label: pc('open'), group: pg('project'), run: openProjectDialog },
      { id: 'new-project', label: pc('newProject'), group: pg('project'), run: openNewProjectWizard },
      { id: 'close-project', label: pc('closeProject'), group: pg('project'), run: closeProject },
      { id: 'settings', label: pc('settings'), group: pg('project'), shortcut: 'Ctrl+,', run: openSettings },
      { id: 'init', label: pc('init'), group: pg('project'), run: openInitProjectDialog },
      { id: 'examples', label: pc('examples'), group: pg('project'), run: openExamples },
      { id: 'new', label: pc('new'), group: pg('scenario'), shortcut: 'Ctrl+N', run: newScenario },
      { id: 'open-file', label: pc('openFile'), group: pg('scenario'), shortcut: 'Ctrl+O', run: openFileDialog },
      { id: 'save', label: pc('save'), group: pg('scenario'), shortcut: 'Ctrl+S', run: saveFeature },
      { id: 'save-as', label: pc('saveAs'), group: pg('scenario'), shortcut: 'Ctrl+Shift+S', run: saveFeatureAs },
      { id: 'export', label: pc('export'), group: pg('scenario'), run: openExportDialog },
      { id: 'import', label: pc('import'), group: pg('scenario'), run: openImportDialog },
      { id: 'import-features', label: pc('importFeatures'), group: pg('scenario'), run: openImportFeaturesDialog },
      { id: 'steps', label: pc('steps'), group: pg('scenario'), run: openStepsDialog },
      { id: 'snippets', label: pc('snippets'), group: pg('scenario'), shortcut: 'Ctrl+Shift+Space', run: openSnippetPalette },
      { id: 'find-replace', label: pc('findReplace'), group: pg('scenario'), shortcut: 'Ctrl+H', run: openFindReplace },
      { id: 'find', label: pc('find'), group: pg('scenario'), shortcut: 'Ctrl+F', run: () => monaco?.openFind() },
      { id: 'format', label: pc('format'), group: pg('scenario'), shortcut: 'Shift+Alt+F', run: () => void monaco?.formatDocument() },
      { id: 'goto-symbol', label: pc('gotoSymbol'), group: pg('scenario'), shortcut: 'Ctrl+Shift+O', run: () => monaco?.openSymbolOutline() },
      { id: 'project-replace', label: pc('projectReplace'), group: pg('scenario'), run: () => (showProjectReplace = true) },
      { id: 'duplicate', label: pc('duplicate'), group: pg('scenario'), run: () => activeTab && !isWelcome && openDuplicateDialog(activeTab) },
      { id: 'rename-feature', label: pc('renameFeature'), group: pg('scenario'), run: () => {
        if (!activeTab || isWelcome) return
        renameFeaturePath = activeTab
        showRenameFeature = true
      }},
      { id: 'delete-feature', label: pc('deleteFeature'), group: pg('scenario'), run: () => activeTab && !isWelcome && deleteFeature(activeTab) },
      { id: 'refactor-indents', label: pc('refactorIndents'), group: pg('refactor'), run: refactorNormalizeIndents },
      { id: 'refactor-blanks', label: pc('refactorBlanks'), group: pg('refactor'), run: refactorCollapseBlank },
      { id: 'steps-help', label: pc('stepsHelp'), group: pg('help'), shortcut: 'F1', run: () => openStepsHelp() },
      { id: 'browser', label: pc('browser'), group: pg('run'), shortcut: 'Ctrl+B', run: () => void openBrowser() },
      { id: 'record', label: pc('record'), group: pg('run'), shortcut: 'Ctrl+R', run: beginRecord },
      { id: 'record-pause', label: pc('recordPause'), group: pg('run'), shortcut: 'Alt+P', run: () => void toggleRecordPause() },
      { id: 'record-stop', label: pc('recordStop'), group: pg('run'), shortcut: 'Ctrl+Shift+R', run: () => void stopRecord() },
      { id: 'record-baseline', label: pc('recordBaseline'), group: pg('run'), run: openBaselineRecordDialog },
      { id: 'stop', label: pc('stop'), group: pg('run'), run: stopRecord },
      { id: 'pause', label: pc('pause'), group: pg('run'), run: toggleRecordPause },
      { id: 'run', label: pc('run'), group: pg('run'), shortcut: 'Ctrl+Enter', run: () => runPrimary(false) },
      { id: 'run-current', label: pc('runCurrent'), group: pg('run'), shortcut: 'Ctrl+Shift+Enter', run: () => runCurrentScenario(false) },
      { id: 'run-current-dry', label: pc('runCurrentDry'), group: pg('run'), run: () => runCurrentScenario(true) },
      { id: 'run-dialog', label: pc('runDialog'), group: pg('run'), run: () => openRunDialog('', {}) },
      { id: 'run-tag', label: pc('runTag'), group: pg('run'), run: () => openRunDialog(tr('menus.runTag').replace('…', ''), {}) },
      { id: 'playwright', label: pc('playwright'), group: pg('run'), run: () => openRunDialog('Playwright', { dryRun: false, headed: true, engine: 'playwright', installPW: true }) },
      { id: 'dry', label: pc('dry'), group: pg('run'), run: () => runPrimary(true) },
      { id: 'batch', label: pc('batch'), group: pg('run'), run: () => toggleBatchMode() },
      { id: 'batch-run', label: pc('batchRun'), group: pg('run'), run: () => runBatchSelected(false) },
      { id: 'batch-dry', label: pc('batchDry'), group: pg('run'), run: () => runBatchSelected(true) },
      { id: 'rerun-failed', label: pc('rerunFailed'), group: pg('run'), run: rerunFailed },
      { id: 'run-history', label: pc('runHistory'), group: pg('run'), run: openRunHistory },
      { id: 'testclient', label: pc('testclient'), group: pg('run'), run: openTestClientDialog },
      { id: 'capture-session', label: pc('captureSession'), group: pg('run'), run: openTestClientDialogForCapture },
      { id: 'validate', label: pc('validate'), group: pg('run'), run: () => openValidateDialog(true) },
      { id: 'validate-browser', label: pc('validateBrowser'), group: pg('run'), run: () => openValidateDialog(false) },
      ...(hasVanessaPlugin()
        ? [
            { id: 'vanessa-dry', label: pc('vanessaDry'), group: pg('run'), run: () => openVanessaDialog(true) },
            { id: 'vanessa', label: pc('vanessa'), group: pg('run'), run: () => openVanessaDialog(false) },
            { id: 'vanessa-rerun', label: pc('vanessaRerun'), group: pg('run'), run: () => openVanessaDialog(false, true) },
            { id: 'vanessa-settings', label: pc('vanessaSettings'), group: pg('run'), run: openVanessaSettingsDialog },
            { id: 'vanessa-monitor', label: pc('vanessaMonitor'), group: pg('run'), run: openVanessaMonitor },
          ]
        : []),
      { id: 'plugins', label: pc('plugins'), group: pg('plugins'), run: () => (showPlugins = true) },
      ...installedPlugins.flatMap((plugin) => {
        if (!plugin.runnable || plugin.vanessa) return []
        const label = pluginLabel(plugin)
        return [
          { id: `plugin-${plugin.name}-dry`, label: tr('menus.runPluginDry', { name: label }), group: pg('plugins'), run: () => openPluginRun(plugin.name, true) },
          { id: `plugin-${plugin.name}`, label: tr('menus.runPlugin', { name: label }), group: pg('plugins'), run: () => openPluginRun(plugin.name, false) },
        ]
      }),
      { id: 'journal', label: pc('journal'), group: pg('view'), shortcut: 'Ctrl+`', run: () => { bottomPanelOpen = true; bottomTab = 'journal' } },
      { id: 'results', label: pc('results'), group: pg('view'), run: () => { bottomPanelOpen = true; bottomTab = 'results' } },
      { id: 'allure-serve', label: pc('allureServe'), group: pg('view'), run: () => serveAllureReport(projectArtifacts.allureDir || '') },
      { id: 'validate-panel', label: pc('validatePanel'), group: pg('view'), run: () => { bottomPanelOpen = true; bottomTab = 'validate' } },
      { id: 'error-panel', label: pc('errorPanel'), group: pg('view'), run: () => { bottomPanelOpen = true; bottomTab = 'error' } },
      { id: 'explorer', label: pc('explorer'), group: pg('view'), run: () => { sidebarVisible = true; saveLayout({ sidebarVisible: true }) } },
      { id: 'explorer-hide', label: pc('explorerHide'), group: pg('view'), run: () => { sidebarVisible = false; saveLayout({ sidebarVisible: false }) } },
      { id: 'preview', label: previewVisible ? pc('previewHide') : pc('previewShow'), group: pg('view'), run: togglePreview },
      { id: 'steps-panel', label: stepsPanelVisible ? pc('stepsPanelHide') : pc('stepsPanelShow'), group: pg('view'), run: toggleStepsPanel },
      { id: 'compact', label: compactLabel, group: pg('view'), run: () => { toolbarCompact = !toolbarCompact; persistSettings() } },
      { id: 'refactor-urls', label: pc('refactorUrls'), group: pg('refactor'), run: refactorUpdateUrls },
      { id: 'hotkeys', label: pc('hotkeys'), group: pg('help'), shortcut: 'Shift+F1', run: () => (showHotkeys = true) },
      { id: 'reset-layout', label: pc('resetLayout'), group: pg('view'), run: resetWindowLayout },
      { id: 'updates', label: pc('updates'), group: pg('help'), run: checkUpdates },
      { id: 'about', label: pc('about'), group: pg('help'), run: showAboutDialog },
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
      appendLog(tr('journal.search.openScenario'))
      return
    }
    monaco?.openFindReplace()
  }

  function resetWindowLayout() {
    const layout = resetUILayout()
    sidebarVisible = layout.sidebarVisible
    bottomPanelOpen = layout.bottomPanelOpen
    bottomPanelHeight = layout.bottomPanelHeight
    previewVisible = layout.previewVisible
    previewWidth = layout.previewWidth
    appendLog(tr('journal.layout.reset'))
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

  async function runScenarioHintsAutoFix(text: string): Promise<{ text: string; count: number }> {
    if (!editorSettings.scenarioHints || !editorSettings.scenarioHintsAutoFixOnSave) {
      return { text, count: 0 }
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

  let hintFixInFlight = false

  async function applyEditorHintFix(hint: gui.ScenarioHintDTO) {
    if (hintFixInFlight || isWelcome) return
    hintFixInFlight = true
    try {
      const result = await ApplyScenarioHintFix({
        text: editorText,
        hintId: hint.id,
        stepIndex: hint.stepIndex,
      })
      if (result.count > 0) {
        await applyEditorText(result.text, { skipValidate: true })
        appendLog(tr('journal.hint.fixed', { title: hint.title }))
        scheduleValidateEditor(150)
      }
    } finally {
      hintFixInFlight = false
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
      const editorContent = monaco?.getEditorText() ?? editorText
      if (activeTab === path || isUntitled(path)) {
        postRecordStepCount = (await ParseEditorSteps(editorContent)).length
      } else {
        const content = await ReadFeature(path)
        postRecordStepCount = (await ParseEditorSteps(content)).length
      }
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
    postRecordBaselineText = ''
    showPostRecordDiff = false
  }

  function openPostRecordDiff() {
    if (!postRecordPath) return
    const modified = monaco?.getEditorText() ?? editorText
    if (postRecordBaselineText === modified) {
      appendLog(tr('journal.record.noPostRecordDiff'))
      return
    }
    showPostRecordDiff = true
  }

  async function postRecordValidate() {
    if (!postRecordPath) return
    if (activeTab !== postRecordPath) await loadFeature(postRecordPath)
    const resolved = resolveLogicalRunTarget(postRecordPath, tabs, activeTab)
    const targets = await materializeRunTargets([resolved])
    if (!targets.length) {
      appendLog(tr('journal.record.postRecordValidateFailed'))
      return
    }
    await validateProject(false, settingsBrowser || 'chromium', targets)
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
      if (isUntitled(path)) {
        const tab = tabs.find((t) => t.path === path)
        const text =
          path === activeTab
            ? (monaco?.getEditorText() ?? editorText)
            : tab
              ? tabEditorText(tab)
              : ''
        if (!text || !projectPath) {
          appendLog(tr('journal.project.saveOrOpenFirst'))
          return
        }
        const name = newName.trim().replace(/\.feature$/i, '') || 'copy'
        const target = `${projectPath.replace(/\\/g, '/')}/${name}.feature`
        await SaveFeature(target, text)
        await refreshProject()
        await loadFeature(target)
        appendLog(tr('journal.file.duplicateCreated', { name: `${name}.feature` }))
        return
      }
      const target = await DuplicateFeature(path, newName)
      await refreshProject()
      await loadFeature(target)
      appendLog(tr('journal.file.duplicateCreated', { name: basename(target) }))
    } catch (e: any) {
      appendLog(tr('journal.error.generic', { error: String(e) }))
    }
  }

  async function duplicateFeature(path: string) {
    openDuplicateDialog(path)
  }

  async function deleteFeature(path: string) {
    if (!path || isWelcome) return
    const ok = await askConfirm({
      title: tr('confirm.deleteFeature.title'),
      message: tr('confirm.deleteFeature.message', { name: basename(path) }),
      confirmLabel: tr('confirm.deleteFeature.confirmLabel'),
      danger: true,
    })
    if (!ok) return
    await doDeleteFeature(path)
  }

  async function doDeleteFeature(path: string) {
    if (isUntitled(path)) {
      finalizeCloseTab(path)
      appendLog(tr('journal.file.closed', { name: untitledLabel(path) }))
      return
    }
    try {
      await DeleteFeature(path)
      closeTab(path)
      await refreshProject()
      appendLog(tr('journal.file.deleted', { name: basename(path) }))
    } catch (e: any) {
      appendLog(tr('journal.error.generic', { error: String(e) }))
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
    appendLog(tr('journal.catalog.batchSelected', { count: batchSelected.length }))
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
      appendLog(tr('journal.file.moved', { name: basename(newPath) }))
    } catch (e: any) {
      appendLog(tr('journal.error.generic', { error: String(e) }))
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
      appendLog(tr('journal.file.renamed', { name: basename(newPath) }))
    } catch (e: any) {
      appendLog(tr('journal.error.generic', { error: String(e) }))
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
      appendLog(tr('journal.file.replaceDone', { replacements: result.replacements, files: result.filesChanged }))
      showProjectReplace = false
      await refreshProject()
      if (activeTab && !isWelcome) {
        editorText = await ReadFeature(activeTab)
        tabs = tabs.map((t) => (t.path === activeTab ? { ...t, content: editorText } : t))
        validateEditor()
      }
    } catch (e: any) {
      appendLog(tr('journal.error.generic', { error: String(e) }))
    } finally {
      projectReplaceBusy = false
    }
  }

  async function refreshRunResults() {
    try {
      runResults = await ListRunResults(50)
      flakyMetrics = await FlakyMetrics(200)
    } catch {
      runResults = []
      flakyMetrics = null
    }
  }

  function runResultFeaturePath(targetPath: string): string {
    if (isUntitled(targetPath)) return untitledLabel(targetPath)
    return targetPath
  }

  function finalizeRunPanels(
    runSince: string,
    diskTargets: string[],
    runTargets: string[],
    cliError = '',
    runner = 'playwright',
    scenario = '',
  ) {
    lastRunSince = runSince
    const batch = remapRunResultPaths(
      filterRunResultsSince(runResults, runSince),
      diskTargets,
      runTargets,
    )
    lastRunBatchResults = batch
    lastErrorEntry = pickLastRunError(batch)
    if (!lastErrorEntry && cliError && !/context canceled/i.test(cliError)) {
      const featurePath =
        runTargets.length === 1
          ? runResultFeaturePath(runTargets[0])
          : diskTargets.length === 1
            ? diskTargets[0]
            : tr('journal.run.defaultLabel')
      lastErrorEntry = buildSyntheticRunError({
        featurePath,
        scenario: scenario || undefined,
        message: cliError,
        runner,
      })
    }
    if (lastRunBatchResults.length === 0 && lastErrorEntry) {
      lastRunBatchResults = [lastErrorEntry]
    }
  }

  async function openRunHistory() {
    await refreshRunResults()
    showRunHistory = true
  }

  async function openFeatureFromHistory(path: string) {
    showRunHistory = false
    const feature = path.includes('::') ? path.slice(0, path.indexOf('::')) : path
    const resolved = resolveLogicalRunTarget(feature, tabs, activeTab)
    if (isUntitled(resolved) && tabs.some((t) => t.path === resolved)) {
      await loadFeature(resolved)
      return
    }
    await loadFeature(resolved)
  }

  async function openProjectAt(path: string) {
    if (!path) return
    const switching = !!projectPath && projectPath !== path
    if (switching && (tabs.length > 0 || browserOpen || recording)) {
      const dirty = tabs.filter((t) => t.dirty)
      const ok = await askConfirm({
        title: tr('confirm.openOtherProject.title'),
        message: dirty.length
          ? tr('confirm.openOtherProject.messageDirty', { count: dirty.length })
          : tr('confirm.openOtherProject.messageClean'),
        confirmLabel: tr('confirm.openOtherProject.confirmLabel'),
        danger: dirty.length > 0,
      })
      if (!ok) return
      await resetWorkspaceForProjectSwitch()
    }
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
      appendLog(tr('journal.project.opened', { path: projectPath }))
      syncIdleStatus()
      await refreshRunResults()
      await refreshArtifacts()
      schedulePersistSession()
    } catch (e: any) {
      appendLog(tr('journal.error.generic', { error: String(e) }))
      setStatus(String(e), 'error')
    }
  }

  async function teardownDesktopSession() {
    stopBrowserWatch()
    try {
      if (recording) {
        await StopRecordingCapture()
      }
    } catch {
      /* ignore */
    }
    try {
      await CloseBrowser()
    } catch {
      /* ignore */
    }
    try {
      CancelRun()
    } catch {
      /* ignore */
    }
    recording = false
    browserOpen = false
    recordPaused = false
    liveRecordStepLines = {}
    recordingTargetPath = ''
  }

  async function resetWorkspaceForProjectSwitch() {
    await teardownDesktopSession()
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
    postRecordPath = ''
    postRecordStepCount = 0
    postRecordBaselineText = ''
  }

  async function closeProject() {
    if (!projectPath) return
    const dirty = tabs.filter((t) => t.dirty)
    if (dirty.length > 0) {
      const ok = await askConfirm({
        title: tr('confirm.closeProject.title'),
        message: tr('confirm.closeProject.message', { count: dirty.length }),
        confirmLabel: tr('confirm.closeProject.confirmLabel'),
        danger: true,
      })
      if (!ok) return
    }
    await teardownDesktopSession()
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
    appendLog(tr('journal.project.closed'))
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
      appendLog(tr('journal.examples.folderNotFound'))
      return
    }
    await openProjectAt(examples)
    sidebarVisible = true
    saveLayout({ sidebarVisible: true })
    selectTab(WELCOME_KEY)
    appendLog(tr('journal.examples.selectScenario'))
    setStatus(tr('journal.status.examplesOpened'), 'normal')
  }

  async function rerunFailed() {
    const source = lastRunBatchResults.length > 0 ? lastRunBatchResults : runResults
    const failed = [...new Map(source.filter((e) => !e.success).map((e) => [e.path, e])).values()]
    if (!failed.length) {
      appendLog(tr('journal.run.noFailed'))
      return
    }
    for (const entry of failed) {
      await runSingleScenario(entry)
    }
  }

  async function runFlakyTriplicate(entry: gui.RunResultEntry) {
    appendLog(tr('journal.run.flakyTriplicate', { path: entry.path }))
    for (let i = 0; i < 3; i++) {
      appendLog(tr('journal.run.flakyIteration', { current: i + 1 }))
      await runSingleScenario(entry)
    }
  }

  async function runSingleScenario(entry: gui.RunResultEntry) {
    const sep = entry.path.indexOf('::')
    const rawPath = sep >= 0 ? entry.path.slice(0, sep) : entry.path
    const scenario = sep >= 0 ? entry.path.slice(sep + 2) : ''
    const filePath = resolveLogicalRunTarget(rawPath, tabs, activeTab)
    if (!filePath) return
    if (activeTab !== filePath) await loadFeature(filePath)
    await executeRun({ ...lastRun, dryRun: false, scenario }, [filePath])
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
    syncActiveTabContent()
    const liveActive = monaco?.getEditorText() ?? editorText
    return materializeRunTargetPaths(paths, tabs, activeTab, editorText, liveActive, {
      writeTempFeature: WriteTempFeature,
      readFeature: ReadFeature,
    })
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
    if (playing) {
      appendLog(tr('journal.run.alreadyRunning'))
      return
    }
    if (batchSelected.length > 0) {
      void runBatchSelected(dryRun)
      return
    }
    if (!runDialogConfirmed) {
      openRunDialog(dryRun ? tr('journal.run.mode.dryRun') : '', { dryRun })
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
      try {
        const resolved = await ResolveRunFromLine(editorText, line)
        if (resolved.scenario) scenario = resolved.scenario
        if (resolved.partial && resolved.startStep >= 0) {
          startStep = resolved.startStep
          endStep = resolved.endStep >= 0 ? resolved.endStep : -1
        }
      } catch (e: unknown) {
        appendLog(tr('journal.run.resolveLineError', { error: String(e) }))
        return
      }
    }
    monaco?.gotoLine(line)
    const runOpts = { ...lastRun, dryRun, scenario, startStep, endStep }
    const range = partialRunLogSuffix(startStep, endStep)
    const mode = runModeLabel(dryRun)
    if (partial && startStep >= 0) {
      appendLog(tr('journal.run.fromStep', { mode, step: startStep + 1, scenario, range }))
    } else if (scenario) {
      appendLog(tr('journal.run.scenarioAtLine', { mode, scenario, line, range }))
    } else {
      appendLog(tr('journal.run.fileAtLine', { mode, line, range }))
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
    const mode = runModeLabel(dryRun)
    if (endStep >= 0 && scenario) {
      appendLog(tr('journal.run.untilStep', { mode, step: endStep + 1, scenario, range }))
    } else if (scenario) {
      appendLog(tr('journal.run.scenarioAtLineNoRange', { mode, scenario, line }))
    } else {
      appendLog(tr('journal.run.fileAtLineNoRange', { mode, line }))
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
        if (!playing && !isWelcome && (activeTab || projectPath || batchSelected.length > 0)) void runPrimary(false)
        break
      case 'run-current':
        if (!playing && !isWelcome && activeTab) void runCurrentScenario(false)
        break
      case 'browser':
        void openBrowser()
        break
      case 'record':
        if (recording) {
          void focusBrowser()
          break
        }
        beginRecord()
        break
      case 'record-stop':
        if (recording || playing || browserOpen) void stopRecord()
        break
      case 'record-pause':
        if (recording) void toggleRecordPause()
        break
      case 'new':
        newScenario()
        break
      case 'open':
        void openFileDialog()
        break
      case 'find':
        monaco?.openFind()
        break
      case 'find-replace':
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
      case 'format':
        if (!isWelcome && activeTab) void monaco?.formatDocument()
        break
      case 'goto-symbol':
        if (!isWelcome && activeTab) monaco?.openSymbolOutline()
        break
      case 'escape':
        openMenu = null
        if (showCommandPalette) showCommandPalette = false
        if (showSnippetPalette) showSnippetPalette = false
        break
    }
  }

  function onModalEscapeCapture(e: KeyboardEvent) {
    if (e.key !== 'Escape') return

    const dismiss = (fn: () => void) => {
      fn()
      e.stopImmediatePropagation()
      e.preventDefault()
    }

    if (confirmDialog) return dismiss(() => closeConfirm(false))
    if (showHttpAuth) return dismiss(closeHttpAuthDialog)
    if (showPickerStep) return dismiss(() => { showPickerStep = false })

    const inModal = e.target instanceof Element && e.target.closest('.modal-backdrop, .palette-backdrop')
    if (!inModal && monacoOverlayConsumesEscape()) return

    if (pendingCloseTab) return dismiss(cancelCloseTab)
    if (showOtp) return dismiss(cancelOtp)
    if (showPluginRun) return dismiss(() => { showPluginRun = false })
    if (showCommandPalette) return dismiss(() => { showCommandPalette = false })
    if (showSnippetPalette) return dismiss(() => { showSnippetPalette = false })
    if (recordStepPickerOpen) return
    if (showRecord) return dismiss(() => { showRecord = false })
    if (showSteps) return dismiss(() => { showSteps = false })
    if (showProjectReplace) return dismiss(() => { showProjectReplace = false })
    if (showRunHistory) return dismiss(() => { showRunHistory = false })
    if (showPostRecordDiff && postRecordPath) return dismiss(() => { showPostRecordDiff = false })
    if (showHotkeys) return dismiss(() => { showHotkeys = false })
    if (showPlugins) return dismiss(() => { showPlugins = false })
    if (showAbout) return dismiss(() => { showAbout = false })
    if (showUpdateCheck && !updateDownloading) return dismiss(() => { showUpdateCheck = false })
    if (showImportFeatures) return dismiss(() => { showImportFeatures = false })
    if (showImport) return dismiss(() => { showImport = false })
    if (showDuplicateFeature) {
      return dismiss(() => {
        showDuplicateFeature = false
        duplicateFeaturePath = ''
      })
    }
    if (showInitProject) return dismiss(() => { showInitProject = false })
    if (showNewProjectWizard) return dismiss(() => { showNewProjectWizard = false })
    if (showValidate) return dismiss(() => { showValidate = false })
    if (showMoveFeature) {
      return dismiss(() => {
        showMoveFeature = false
        moveFeaturePath = ''
      })
    }
    if (showRenameFeature) {
      return dismiss(() => {
        showRenameFeature = false
        renameFeaturePath = ''
      })
    }
    if (showOpenProject) return dismiss(() => { showOpenProject = false })
    if (showRefactorUrl) return dismiss(() => { showRefactorUrl = false })
    if (showExport) return dismiss(() => { showExport = false })
    if (showVanessaSettings) return dismiss(() => { showVanessaSettings = false })
    if (showStepsHelp) {
      return dismiss(() => {
        showStepsHelp = false
        stepsHelpQuery = ''
      })
    }
    if (showTestClient) {
      return dismiss(() => {
        showTestClient = false
        testClientSuggestName = ''
      })
    }
    if (showVanessaRun) return dismiss(() => { showVanessaRun = false })
    if (showRun) return dismiss(() => { showRun = false })
    if (showSettings) return dismiss(cancelSettings)
  }

  function onGlobalKeydown(e: KeyboardEvent) {
    if (shouldIgnoreAppHotkey(e)) return
    const id = matchHotkey(e)
    if (!id) return
    if (showSettings && (id === 'save' || id === 'settings')) return
    if (id === 'escape') {
      const target = e.target
      if (target instanceof Element && target.closest('.modal-backdrop, .palette-backdrop')) {
        return
      }
      if (monacoOverlayConsumesEscape()) {
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
    statusMessage = tr('statusBar.default')
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
    await refreshAllureStatus()
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
    appendLog(tr('journal.reports.allureServe'))
    const result = await ServeAllure(path)
    if (result.output) appendLog(result.output.trimEnd())
    if (result.error) {
      appendLog(tr('journal.error.generic', { error: result.error }))
    } else {
      allureServeRunning = true
    }
    await refreshAllureStatus(path)
  }

  async function refreshAllureStatus(dir = '') {
    try {
      const status = await AllureStatus(dir || projectArtifacts.allureDir || '')
      allureInstalled = status.installed !== false
      allureServeRunning = !!status.running
    } catch {
      allureInstalled = true
    }
  }

  function openAllureInstallHelp() {
    void OpenExternalURL('https://docs.qameta.io/allure/#_installing_a_commandline')
  }

  async function openHtmlReport(path = '') {
    const result = await OpenHTMLReport(path)
    if (result.error) {
      appendLog(tr('journal.reports.reportOpenError', { error: result.error }))
      return
    }
    if (result.output) appendLog(tr('journal.reports.reportOpened', { path: result.output.trim() }))
  }

  async function openTraceReport(path = '') {
    const result = await OpenTrace(path)
    if (result.error) {
      appendLog(tr('journal.reports.traceViewerError', { error: result.error }))
      return
    }
    if (result.output) appendLog(tr('journal.reports.traceViewer', { path: result.output.trim() }))
  }

  function resolveRunResultFeaturePath(featurePath: string): string {
    return resolveLogicalRunTarget(featurePath, tabs, activeTab)
  }

  async function gotoReportStep(req: {
    feature_path: string
    scenario: string
    leaf_index: number
    line: number
  }) {
    const featurePath = resolveRunResultFeaturePath(req.feature_path)
    if (!featurePath) return
    bottomPanelOpen = true
    if (req.line > 0) {
      await loadFeature(featurePath)
      gotoEditorLine(req.line)
      setStatus(`Line ${req.line}`, 'busy')
      return
    }
    if (req.leaf_index >= 0) {
      await gotoFailedStep({
        path: `${req.feature_path}::${req.scenario}`,
        success: false,
        message: '',
        runner: '',
        at: '',
        failed_step: req.leaf_index,
      } as gui.RunResultEntry)
    } else {
      await loadFeature(featurePath)
    }
  }

  async function openTraceFromReport(req: {
    trace_path: string
    report_dir: string
    trace_offset_ms?: number
    step_index?: number
  }) {
    let path = req.trace_path || ''
    if (path && req.report_dir && !path.match(/^[a-zA-Z]:\\|^\//)) {
      path = `${req.report_dir.replace(/\\/g, '/')}/${path.replace(/\\/g, '/')}`.replace(/\/+/g, '/')
    }
    await openTraceReport(path)
    if (req.trace_offset_ms != null && req.trace_offset_ms >= 0) {
      const sec = req.trace_offset_ms < 1000
        ? `${req.trace_offset_ms}ms`
        : `${(req.trace_offset_ms / 1000).toFixed(2)}s`
      const step = req.step_index != null ? ` (step #${req.step_index + 1})` : ''
      appendLog(tr('journal.reports.traceSeekHint', { offset: sec + step }))
    }
  }

  async function rerunFromReport(req: { feature_path: string; scenario: string }) {
    const featurePath = resolveRunResultFeaturePath(req.feature_path)
    if (!featurePath || !projectPath) return
    await loadFeature(featurePath)
    bottomPanelOpen = true
    bottomTab = 'journal'
    await executeRun({ ...lastRun, dryRun: false, scenario: req.scenario, html: true }, [featurePath])
  }

  async function gotoFailedStep(entry: gui.RunResultEntry) {
    const idx = entry.path.indexOf('::')
    const rawFeaturePath = idx < 0 ? entry.path : entry.path.slice(0, idx)
    const scenario = idx < 0 ? '' : entry.path.slice(idx + 2)
    const featurePath = resolveRunResultFeaturePath(rawFeaturePath)
    if (!featurePath) return
    if (featurePath === activeTab && isUntitled(featurePath)) {
      if (entry.failed_step != null && entry.failed_step >= 0) {
        try {
          const temp = await WriteTempFeature(editorText)
          const line = await FailedStepLine(temp, scenario, entry.failed_step)
          if (line > 0) gotoEditorLine(line)
        } catch {
          const lineMatch = entry.message?.match(/line\s+(\d+)/i)
          if (lineMatch) gotoEditorLine(Number(lineMatch[1]))
        }
      } else {
        const lineMatch = entry.message?.match(/line\s+(\d+)/i)
        if (lineMatch) gotoEditorLine(Number(lineMatch[1]))
      }
      return
    }
    if (entry.failed_step == null || entry.failed_step < 0) {
      await loadFeature(featurePath)
      return
    }
    try {
      const line = await FailedStepLine(featurePath, scenario, entry.failed_step)
      await loadFeature(featurePath)
      if (line > 0) gotoEditorLine(line)
    } catch (e: unknown) {
      appendLog(tr('journal.reports.gotoStepFailed', { error: String(e) }))
      await loadFeature(featurePath)
    }
  }

  async function openArtifactPath(path: string) {
    if (!path) return
    try {
      await OpenFolder(path)
    } catch (e: unknown) {
      appendLog(tr('journal.reports.openFailed', { error: String(e) }))
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
      appendLog(tr('journal.file.moved', { name: basename(newPath) }))
    } catch (e: any) {
      appendLog(tr('journal.error.generic', { error: String(e) }))
    }
  }

  async function importDroppedFeatures(destDir: string, paths: string[]) {
    if (!destDir || !paths.length) return
    try {
      const imported = await ImportFeatures(destDir, paths)
      await refreshProject()
      appendLog(tr('journal.file.imported', { count: imported.length }))
      if (imported.length === 1) await loadFeature(imported[0])
    } catch (e: any) {
      appendLog(tr('journal.file.importError', { error: String(e) }))
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
      appendLog(tr('journal.file.manyTabs', { count: tabs.length }))
    }
  }

  function syncTabContent(tabPath: string) {
    if (!tabPath || tabPath === WELCOME_KEY) return
    const liveText =
      tabPath === activeTab && monaco ? (monaco.getEditorText() ?? editorText) : editorText
    tabs = tabs.map((t) => {
      if (t.path !== tabPath) return t
      const dirty = liveText !== t.content
      if (!dirty) {
        if (!t.dirty && t.draft === undefined) return t
        return { ...t, dirty: false, draft: undefined }
      }
      return { ...t, draft: liveText, dirty: true }
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

  async function checkActiveTabDiskStale() {
    if (isWelcome || !activeTab || isUntitled(activeTab) || playing || recording) return
    const tab = tabs.find((t) => t.path === activeTab)
    if (!tab || tab.dirty) return
    const baseline = tab.draft ?? tab.content
    try {
      const disk = await ReadFeature(activeTab)
      if (disk === baseline) return
      const ok = await askConfirm({
        title: tr('confirm.diskChanged.title'),
        message: tr('confirm.diskChanged.message', { name: basename(activeTab) }),
        confirmLabel: tr('confirm.diskChanged.confirmLabel'),
        danger: true,
      })
      if (!ok) return
      tabs = tabs.map((t) =>
        t.path === activeTab ? { ...t, content: disk, dirty: false, draft: undefined, unloaded: false } : t,
      )
      await applyEditorText(disk, { saved: true, switchTab: true, tabPath: activeTab, skipValidate: true })
      appendLog(tr('journal.file.reloaded', { name: basename(activeTab) }))
    } catch {
      /* ignore */
    }
  }

  async function ensureRecordingTabSwitchAllowed(path: string): Promise<boolean> {
    if (!recording || !recordingTargetPath) return true
    if (recordingTabSwitchAllowed(recording, recordPaused, recordingTargetPath, path)) {
      return true
    }
    if (skipRecordTabSwitchConfirm) {
      recordingTargetPath = normalizeRecordTabPath(path)
      appendLog(tr('journal.record.target', { name: basename(path) }))
      return true
    }
    const ok = await askConfirm({
      title: tr('confirm.recordingActive.title'),
      message: tr('confirm.recordingActive.message', { from: basename(recordingTargetPath), to: basename(path) }),
      confirmLabel: tr('confirm.recordingActive.confirmLabel'),
      danger: true,
      dontAskAgainLabel: tr('confirm.recordingActive.dontAskAgainLabel'),
    })
    if (ok) {
      recordingTargetPath = normalizeRecordTabPath(path)
      appendLog(tr('journal.record.target', { name: basename(path) }))
    }
    return ok
  }

  async function loadFeature(path: string, opts?: { skipRecordingGuard?: boolean }) {
    if (!opts?.skipRecordingGuard) {
      const allowed = await ensureRecordingTabSwitchAllowed(path)
      if (!allowed) return
    }
    const leavingTab = activeTab
    if (leavingTab && !isWelcome && leavingTab !== path) {
      syncTabContent(leavingTab)
    }
    const existing = tabs.find((t) => t.path === path)
    if (existing) {
      if (path === activeTab && !tabNeedsDiskReload(existing)) {
        return
      }
      let text = tabEditorText(existing)
      if (tabNeedsDiskReload(existing)) {
        try {
          text = await ReadFeature(path)
          tabs = tabs.map((t) =>
            t.path === path ? { ...t, content: text, dirty: false, draft: undefined, unloaded: false } : t,
          )
        } catch (e: any) {
          appendLog(tr('journal.file.openError', { error: String(e) }))
          return
        }
      }
      welcomeTabVisible = false
      activeTab = path
      await applyEditorText(text, { saved: !existing.dirty, switchTab: true, tabPath: path, skipValidate: true })
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
          appendLog(tr('journal.file.draftRestored', { name: basename(path) }))
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
      activeTab = path
      await applyEditorText(content, { saved: !dirty, switchTab: true, tabPath: path, skipValidate: true })
      trimTabsMemory()
      stepsPanelCollapsed = resolveStepsPanelCollapsed()
      schedulePersistSession()
    } catch (e: any) {
      appendLog(tr('journal.file.openError', { error: String(e) }))
    }
  }

  function selectTab(path: string) {
    if (path === WELCOME_KEY) {
      if (recording && !recordPaused) {
        appendLog(tr('journal.record.pauseToWelcome'))
        setStatus(tr('journal.status.recordingActive'), 'busy')
        return
      }
      const leavingTab = activeTab
      if (leavingTab && !isWelcome) {
        syncTabContent(leavingTab)
      }
      void applyEditorText('', { switchTab: true, tabPath: null, skipValidate: true })
      clearEditorValidation()
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
    if (recording && recordingTargetPath && isSameRecordTab(path, recordingTargetPath)) {
      void (async () => {
        const ok = await askConfirm({
          title: tr('confirm.recordingTargetTab.title'),
          message: recordPaused
            ? tr('confirm.recordingTargetTab.messagePaused', { name: basename(path) })
            : tr('confirm.recordingTargetTab.messageRecording', { name: basename(path) }),
          confirmLabel: recordPaused ? tr('confirm.recordingTargetTab.confirmLabelClose') : tr('confirm.recordingTargetTab.confirmLabelForceClose'),
          danger: true,
        })
        if (!ok) return
        if (!recordPaused) await StopRecordingCapture()
        finalizeCloseTab(path)
      })()
      return
    }
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
        clearEditorValidation()
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
    const picked = await PickSaveFile(tr('filePicker.saveAs'), basename(activeTab))
    if (!picked) return
    try {
      const text = monaco?.getEditorText() ?? editorText
      await SaveFeature(picked, text)
      const oldPath = activeTab
      tabs = tabs.map((t) =>
        t.path === oldPath ? { path: picked, content: text, dirty: false, draft: undefined } : t,
      )
      activeTab = picked
      await rememberFeature(picked)
      await refreshProject()
      appendLog(tr('journal.file.savedAs', { name: basename(picked) }))
      setStatus(tr('journal.status.saved'), 'success')
    } catch (e: any) {
      appendLog(tr('journal.file.saveError', { error: String(e) }))
      setStatus(tr('journal.status.saveError'), 'error')
    }
  }

  async function saveFeature() {
    if (!activeTab || isWelcome) return
    if (isUntitled(activeTab)) {
      await saveFeatureAs()
      return
    }
    try {
      let text = monaco?.getEditorText() ?? editorText
      if (editorSettings.formatOnSave) {
        await monaco?.formatDocument()
        text = monaco?.getEditorText() ?? text
        editorText = text
      }
      const { text: autoFixed, count: autoFixCount } = await runScenarioHintsAutoFix(text)
      if (autoFixCount > 0) {
        text = autoFixed
        editorText = text
        await monaco?.setContent(text)
        appendLog(tr('journal.hint.autoFixed', { count: autoFixCount }))
      }
      await SaveFeature(activeTab, text)
      markActiveTabSaved(text)
      try {
        await ClearFeatureDraft(activeTab)
      } catch {
        /* ignore */
      }
      appendLog(tr('journal.file.saved', { name: basename(activeTab) }))
      setStatus(tr('journal.status.saved'), 'success')
    } catch (e: any) {
      appendLog(tr('journal.file.saveError', { error: String(e) }))
      setStatus(tr('journal.status.saveError'), 'error')
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

  function clearEditorValidation() {
    editorValidationIssues = []
    stepStatusError = false
    monaco?.setMarkers([])
    if (statusMessage === tr('journal.status.scenarioError')) {
      setStatus('', 'normal')
    }
  }

  async function validateEditor() {
    if (isWelcome || !activeTab || activeTab === WELCOME_KEY) {
      clearEditorValidation()
      return
    }
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
        stepStatusError = true
        setStatus(tr('journal.status.scenarioError'), 'error')
      } else {
        stepStatusError = false
      }
    } catch {
      if (generation !== validateGeneration || tabAtStart !== activeTab) return
      editorValidationIssues = []
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
    if (isWelcome) return tr('results.validate.hintWelcome')
    return tr('results.validate.hintBrowser')
  }

  function validatePanelDisplayIssues(): gui.ValidationIssue[] {
    if (validatePanelIssues.length > 0) return validatePanelIssues
    if (isWelcome) return []
    return editorValidationIssues
  }

  async function onEditorChange(text: string) {
    editorText = text
    syncActiveTabContent()
    schedulePersistSession()
    if (editorSettings.validateOnType && !isWelcome) {
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
    if (playing) {
      appendLog(tr('journal.run.alreadyRunning'))
      return
    }
    syncActiveTabContent()

    let runTargets = targets
    if (runTargets.length === 0 && activeTab && !isWelcome) {
      runTargets = [activeTab]
    }
    if (runTargets.length === 0 && isWelcome) {
      const featureTabs = tabs.filter((t) => t.path !== WELCOME_KEY)
      if (featureTabs.length === 1) {
        runTargets = [featureTabs[0].path]
      }
    }
    if (runTargets.length === 0 && !projectPath) {
      appendLog(tr('journal.run.openScenario'))
      return
    }

    const diskTargets = runTargets.length ? await materializeRunTargets(runTargets) : []
    if (runTargets.length > 0 && diskTargets.length === 0) {
      appendLog(tr('journal.run.prepareFailed'))
      setStatus(tr('journal.status.runError'), 'error')
      return
    }

    let runOpts = opts
    if (runTargets.length === 1 && runTargets[0] === activeTab && !isWelcome) {
      runOpts = {
        ...opts,
        scenario: resolveStaleRunScenario(opts.scenario || '', editorText, cursorScenarioName()),
      }
    }

    lastRun = { ...runOpts }
    showRun = false
    bottomPanelOpen = true
    bottomTab = 'journal'

    playingLabel =
      runTargets.length > 1
        ? tr('journal.run.scenariosCount', { count: runTargets.length })
        : runTargets.length === 1
          ? featureTabLabel(runTargets[0])
          : runOpts.scenario || runOpts.tag || tr('journal.run.defaultLabel')

    const allureDir = runOpts.allure ? await scenariaSubdir('allure-results') : ''
    const traceDir = !runOpts.dryRun && runOpts.trace ? await scenariaSubdir('traces') : ''
    const videoDir = !runOpts.dryRun && runOpts.video ? await scenariaSubdir('videos') : ''

    const htmlPath = runOpts.html ? await resolveHtmlReportPath(runOpts) : ''

    const junitPath = runOpts.junit ? await scenariaSubdir('junit.xml') : ''

    const summaryJsonPath = runOpts.summaryJson ? await scenariaSubdir('summary.json') : ''

    playing = !runOpts.dryRun
    runningDryRun = runOpts.dryRun
    runProgressCurrent = 0
    runProgressTotal = Math.max(1, diskTargets.length || (runTargets.length > 0 ? runTargets.length : 1))
    runCancelling = false
    runLogStreaming = true
    setStatus(tr('journal.status.testRunning'), 'busy')
    const range = partialRunLogSuffix(runOpts.startStep ?? -1, runOpts.endStep ?? -1)
    if (targets.length) {
      appendLog(tr('journal.run.scenarios', { count: targets.length, range }))
    } else if (runOpts.dryRun) {
      appendLog(tr('journal.run.dryRun', { range }))
    } else {
      appendLog(tr('journal.run.playwright', { range }))
    }

    const runSince = new Date().toISOString()

    let result: Awaited<ReturnType<typeof Run>> = gui.RunResult.createFrom({ output: '', error: '', entries: [] })
    let runThrown: unknown = null
    let journalStreamed = false
    try {
      result = await Run({
        tag: runOpts.tag,
        scenario: runOpts.scenario || '',
        testClient: runOpts.testClient,
        vars: parseVars(runOpts.vars),
        dryRun: runOpts.dryRun,
        headed: runOpts.headed,
        engine: runOpts.dryRun ? '' : runOpts.engine,
        installPlaywright: runOpts.installPW,
        allureDir,
        traceDir,
        videoDir,
        htmlPath,
        junitPath,
        summaryJson: summaryJsonPath,
        browser: runOpts.dryRun ? '' : runOpts.browser || settingsBrowser || 'chromium',
        workers: runOpts.workers || settingsWorkers || 1,
        slowMo: runOpts.dryRun ? 0 : (runOpts.slowMo > 0 ? runOpts.slowMo : settingsSlowMo),
        baseUrl: runOpts.dryRun ? '' : (runOpts.baseUrl || '').trim(),
        startStep: runOpts.startStep ?? -1,
        endStep: runOpts.endStep ?? -1,
        continueOnFail: runOpts.continueOnFail,
        htmlLightMode: runOpts.htmlLightMode,
        reuseLiveBrowser: runOpts.reuseLiveBrowser,
        reportLocale: $locale,
        targets: diskTargets,
      })
      journalStreamed = runLogStreaming
    } catch (err) {
      runThrown = err
    } finally {
      runLogStreaming = false
      playing = false
      runningDryRun = false
      runCancelling = false
      playingLabel = ''
      runProgressCurrent = 0
      runProgressTotal = 0
    }

    if (runThrown) {
      appendLog(tr('journal.error.generic', { error: String(runThrown) }))
      setStatus(tr('journal.status.testError'), 'error')
      bottomTab = 'error'
      return
    }

    if (result.output && !journalStreamed) appendLog(result.output.trimEnd())
    if (result.error) {
      if (/context canceled/i.test(result.error)) {
        appendLog(tr('journal.run.stopped'))
        setStatus(tr('journal.status.testStopped'), 'busy')
        bottomTab = 'journal'
      } else {
        appendLog(tr('journal.error.generic', { error: result.error }))
        setStatus(tr('journal.status.testError'), 'error')
        bottomTab = 'error'
      }
    } else {
      appendLog(tr('journal.run.done'))
      setStatus(tr('journal.status.testDone'), 'success')
      welcomePlayedSuccess = true
      if (showOnboardingTour && runOpts.dryRun) {
        onboardingDryRunDone = true
      }
      void persistSettings()
      bottomTab = runOpts.dryRun ? 'journal' : 'results'
    }
    await refreshRunResults()
    if (result.entries?.length) {
      lastRunSince = runSince
      lastRunBatchResults = remapRunResultPaths(result.entries, diskTargets, runTargets)
      lastErrorEntry = pickLastRunError(lastRunBatchResults)
      if (
        !lastErrorEntry &&
        result.error &&
        !/context canceled/i.test(result.error)
      ) {
        const featurePath =
          runTargets.length === 1
            ? runResultFeaturePath(runTargets[0])
            : diskTargets.length === 1
              ? diskTargets[0]
              : tr('journal.run.defaultLabel')
        lastErrorEntry = buildSyntheticRunError({
          featurePath,
          scenario: runOpts.scenario || cursorScenarioName() || undefined,
          message: result.error,
          runner: runOpts.dryRun ? 'dry-run' : runOpts.engine || 'playwright',
        })
        lastRunBatchResults = [...lastRunBatchResults, lastErrorEntry]
      }
    } else {
      finalizeRunPanels(
        runSince,
        diskTargets,
        runTargets,
        result.error || '',
        runOpts.dryRun ? 'dry-run' : runOpts.engine || 'playwright',
        runOpts.scenario || cursorScenarioName(),
      )
    }
    await refreshArtifacts()
    await refreshAllureStatus()
    const runCancelled = /context canceled/i.test(result.error || '')
    if (runOpts.html && htmlPath && !runCancelled) {
      try {
        if (await ArtifactExists(htmlPath)) {
          await openHtmlReport(htmlPath)
        }
      } catch {
        /* ignore */
      }
    }
    if (
      traceDir &&
      !runCancelled &&
      lastRunBatchResults.some((e) => !e.success && e.runner !== 'dry-run')
    ) {
      try {
        await openTraceReport(traceDir)
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

  async function resolveHtmlReportPath(opts: RunForm): Promise<string> {
    if (!opts.html) return ''
    if (opts.htmlTimestamp) {
      const d = new Date()
      const pad = (n: number) => String(n).padStart(2, '0')
      const stamp = `${d.getFullYear()}${pad(d.getMonth() + 1)}${pad(d.getDate())}-${pad(d.getHours())}${pad(d.getMinutes())}`
      return scenariaSubdir(`report-${stamp}.html`)
    }
    return scenariaSubdir('report.html')
  }

  function confirmRun() {
    runDialogConfirmed = true
    void persistSettings()
    executeRun(runForm)
  }

  async function validateProject(browser: boolean, browserName = settingsBrowser || 'chromium', targets: string[] = []) {
    if (!projectPath) return
    appendLog(browser ? tr('journal.validate.startingBrowser') : tr('journal.validate.starting'))
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
        appendLog(tr('journal.validate.browserSummary', { found, warnings, missing }))
        setStatus(missing > 0 ? tr('journal.status.validateBrowserErrors') : tr('journal.status.validateBrowserDone'), missing > 0 ? 'error' : 'success')
      } catch (err) {
        const msg = err instanceof Error ? err.message : String(err)
        validateCliLog = msg
        appendLog(tr('journal.error.generic', { error: msg }))
        setStatus(tr('journal.status.validateError'), 'error')
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
        appendLog(tr('journal.error.generic', { error: result.error }))
        setStatus(tr('journal.status.validateError'), 'error')
      } else {
        appendLog(tr('journal.validate.done'))
        setStatus(tr('journal.status.validateDone'), 'success')
        if (showOnboardingTour && !browser) {
          onboardingValidateDone = true
        }
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
    let targets: string[] = []
    if (payload.scope === 'current' && activeTab && !isWelcome) {
      targets = await materializeRunTargets([activeTab])
      if (targets.length === 0) {
        appendLog(tr('journal.validate.prepareFailed'))
        setStatus(tr('journal.status.validateError'), 'error')
        return
      }
    }
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

  function openNewProjectWizard() {
    showNewProjectWizard = true
  }

  async function confirmNewProjectWizard(opts: NewProjectWizardResult) {
    showNewProjectWizard = false
    const normalized = opts.path.trim()
    if (!normalized) return
    try {
      if (opts.initScenaria) {
        const out = await InitProjectAt(normalized)
        if (out) appendLog(out.trimEnd())
      }
      await openProjectAt(normalized)
      if (opts.createSample) {
        const fileName = `${opts.featureFileName.replace(/\.feature$/i, '')}.feature`
        const featurePath = `${normalized.replace(/\\/g, '/')}/${fileName}`
        const content = buildFeatureTemplate({
          title: opts.title,
          scenario: opts.scenario,
          startUrl: opts.startUrl,
        })
        await SaveFeature(featurePath, content)
        await loadFeature(featurePath)
        appendLog(tr('journal.project.scenarioCreated', { fileName }))
      }
      setStatus(tr('journal.status.projectCreated'), 'success')
    } catch (e: unknown) {
      const msg = e instanceof Error ? e.message : String(e)
      appendLog(tr('journal.error.generic', { error: msg }))
      setStatus(msg, 'error')
    }
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
      appendLog(tr('journal.browser.openForSession'))
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
    appendLog(tr('journal.testClient.selected', { name }))
    showTestClient = false
  }

  async function openStepsDialog() {
    if (recordingBlocksManualTools()) {
      appendLog(tr('journal.record.pauseToInsertStep'))
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
      appendLog(tr('journal.record.pauseToOpenStepsHelp'))
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
    if (recordingBlocksManualTools()) {
      appendLog(tr('journal.record.pauseToInsertStep'))
      setStatus(tr('journal.status.recordingActive'), 'busy')
      return
    }
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
      appendLog(tr('journal.record.pauseToInsertFromPalette'))
      setStatus(tr('journal.status.recordingActive'), 'busy')
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
      await monaco?.setContent(text)
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
      appendLog(tr('journal.refactor.openStepsNotFound'))
      return
    }
    await applyEditorText(result.text)
    appendLog(tr('journal.refactor.urlUpdated', { count: result.count }))
  }

  async function refactorNormalizeIndents() {
    if (isWelcome) return
    const text = await RefactorNormalizeIndents(editorText)
    await applyEditorText(text)
    appendLog(tr('journal.refactor.indentsNormalized'))
  }

  async function refactorCollapseBlank() {
    if (isWelcome) return
    const text = await RefactorCollapseBlankLines(editorText)
    await applyEditorText(text)
    appendLog(tr('journal.refactor.blankLinesCollapsed'))
  }

  function focusBrowserWindow() {
    void focusBrowser()
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
    appendLog(dry ? tr('journal.plugin.runningDry', { label }) : tr('journal.plugin.running', { label }))
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
    if (result.error) appendLog(tr('journal.error.generic', { error: result.error }))
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
      recordScenarioName = tr('dialogs.record.baselineScenarioDefault')
    } else {
      recordFeatureName = tr('dialogs.record.featureDefault')
      recordScenarioName = tr('dialogs.record.baselineScenarioDefault')
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
    appendLog(tr('journal.record.baselineCreating'))
    try {
      const result = await RecordBaseline({
        output: payload.output || 'recorded.feature',
        featureName: payload.featureName,
        scenarioName: payload.scenarioName,
        steps: payload.steps,
      })
      if (result.output) appendLog(result.output.trimEnd())
      if (result.error) {
        appendLog(tr('journal.error.generic', { error: result.error }))
        setStatus(tr('journal.status.recordError'), 'error')
        return
      }
      appendLog(tr('journal.record.baselineSaved'))
      setStatus(tr('journal.status.featureCreated'), 'success')
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
        setStatus(tr('journal.status.updateAvailable'), 'normal')
      }
    } catch {
      /* offline or dev without wails */
    }
  }

  async function maybeCheckUpdatesOnStartup() {
    if (!settingsCheckUpdatesOnStartup) return
    if (!projectPath) {
      pendingUpdateCheckOnStartup = true
      return
    }
    await checkUpdatesOnStartup()
  }

  async function checkUpdates() {
    appendLog(tr('journal.update.checking'))
    try {
      const info = await CheckUpdateInfo()
      updateCheckInfo = info
      updateCheckMessage = info.message || ''
      updateCheckHasUpdate = !!info.updateAvailable
      appendLog(updateCheckMessage)
      if (info.updateAvailable) {
        appendLog(tr('journal.update.release', { url: info.htmlUrl || '—' }))
        if (info.downloadName) appendLog(tr('journal.update.file', { name: info.downloadName }))
      }
      appendLog(tr('journal.validate.done'))
      showUpdateCheck = true
    } catch (err) {
      const msg = err instanceof Error ? err.message : String(err)
      updateCheckMessage = msg
      updateCheckHasUpdate = false
      appendLog(tr('journal.error.generic', { error: msg }))
      showUpdateCheck = true
    }
  }

  async function openUpdateRelease() {
    const url = updateCheckInfo?.htmlUrl
    if (!url) return
    try {
      await OpenExternalURL(url)
    } catch (err) {
      appendLog(tr('journal.update.linkOpenFailed', { error: err instanceof Error ? err.message : String(err) }))
    }
  }

  function normalizeUpdateProgress(raw: unknown): gui.UpdateProgressDTO {
    const src = (raw && typeof raw === 'object' ? raw : {}) as Record<string, unknown>
    const percent = Number(src.percent ?? src.Percent ?? 0)
    return gui.UpdateProgressDTO.createFrom({
      stage: src.stage ?? src.Stage ?? '',
      message: src.message ?? src.Message ?? '',
      percent: Number.isFinite(percent) ? percent : 0,
    })
  }

  function listenUpdateProgress(onFinished: (result: gui.RunResult) => void) {
    let lastLoggedStage = ''
    const onProgress = (raw: unknown) => {
      updateProgress = normalizeUpdateProgress(raw)
      const stage = updateProgress?.stage ?? ''
      if (stage && stage !== lastLoggedStage) {
        lastLoggedStage = stage
        if (updateProgress?.message) appendLog(updateProgress.message)
      }
    }
    EventsOn('update-progress', onProgress)
    EventsOn('update-finished', onFinished)
    return () => EventsOff('update-progress', 'update-finished')
  }

  async function applyUpdate() {
    if (updateDownloading) return
    updateDownloading = true
    updateProgress = { stage: 'check', message: tr('journal.update.progressStarting'), percent: 0 }
    appendLog(tr('journal.update.installing'))
    const unbind = listenUpdateProgress((result) => {
      unbind()
      if (result?.error) {
        appendLog(tr('journal.update.installError', { error: result.error }))
        updateDownloading = false
        updateProgress = null
        return
      }
      updateProgress = { stage: 'restart', message: tr('journal.update.progressRestarting'), percent: 100 }
      appendLog(tr('journal.update.restarting'))
    })
    try {
      await ApplyUpdate()
    } catch (err) {
      unbind()
      appendLog(tr('journal.update.installError', { error: err instanceof Error ? err.message : String(err) }))
      updateDownloading = false
      updateProgress = null
    }
  }

  async function downloadUpdate() {
    if (updateDownloading) return
    updateDownloading = true
    updateProgress = { stage: 'check', message: tr('journal.update.progressDownloading'), percent: 0 }
    appendLog(tr('journal.update.downloading'))
    const unbind = listenUpdateProgress(async (result) => {
      unbind()
      if (result?.error) {
        appendLog(tr('journal.update.downloadError', { error: result.error }))
        updateDownloading = false
        updateProgress = null
        return
      }
      const path = result.output || ''
      appendLog(tr('journal.update.downloaded', { path }))
      updateCheckMessage = `${updateCheckMessage}\n\n${tr('journal.update.fileSaved', { path })}`.trim()
      const folder = path.replace(/[\\/][^\\/]+$/, '')
      if (folder) await OpenFolder(folder)
      updateDownloading = false
      updateProgress = null
    })
    try {
      await DownloadUpdate()
    } catch (err) {
      unbind()
      appendLog(tr('journal.update.downloadError', { error: err instanceof Error ? err.message : String(err) }))
      updateDownloading = false
      updateProgress = null
    }
  }

  async function openSettings() {
    const s = await LoadSettings()
    applySettingsFromDTO(s)
    settingsDialogBaseline = s
    showSettings = true
  }

  function buildCurrentSettingsDTO(): gui.AppSettingsDTO {
    syncActiveTabContent()
    const sessionTabs = buildSessionTabsSnapshot(tabs, activeTab, editorText, WELCOME_KEY)
    return gui.AppSettingsDTO.createFrom({
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
      navWaitUntil: settingsNavWaitUntil,
      editor: editorSettingsToDTO(editorSettings),
      checklistDismissed,
      welcomePlayedSuccess,
      onboardingCompleted,
      onboardingDismissed,
      onboardingVersion: ONBOARDING_TOUR_VERSION,
      runDialogConfirmed,
      pickerDuringRecording,
      uiLocale,
      startUrl: startURL,
    })
  }

  async function persistSettings() {
    await SaveSettings(buildCurrentSettingsDTO())
  }

  async function syncRecordingOptions() {
    try {
      if (recording) {
        await UpdateRecordingOptions(
          filterRecording,
          navOnlyRecording,
          hoverRecord,
          settingsHeadless,
          settingsScrollBeforeClick,
          settingsHoverRecordMinMs,
        )
      }
      await persistSettings()
    } catch {
      /* session may be closing */
    }
  }

  async function onRecordingHeadlessChange(e: Event) {
    const input = e.currentTarget as HTMLInputElement
    const next = input.checked
    if (next === settingsHeadless) return
    if (browserOpen || recording) {
      input.checked = settingsHeadless
      const ok = await askConfirm({
        title: tr('confirm.headless.title'),
        message: tr('confirm.headless.message'),
        confirmLabel: tr('confirm.headless.confirmLabel'),
        danger: true,
      })
      if (!ok) return
    }
    settingsHeadless = next
    void syncRecordingOptions()
  }

  async function applySettingsCore(closeDialog: boolean) {
    setLocale(uiLocale)
    editorSettings = { ...editorSettings }
    stepsPanelTab = editorSettings.stepsPanelView
    if (!stepsPanelVisible) stepsPanelCollapsed = true
    else if (stepsPanelCollapsed && stepsPanelVisible) stepsPanelCollapsed = false
    setStepHoverEnabled(() => editorSettings.stepHover)
    lastRun = {
      ...lastRun,
      workers: settingsWorkers,
      slowMo: settingsSlowMo,
      browser: settingsBrowser,
    }
    if (recording || browserOpen) await syncRecordingOptions()
    else await persistSettings()
    monaco?.applyEditorSettings(editorSettings)
    if (editorSettings.scenarioHints) {
      await refreshEditorScenarioHints()
    } else {
      editorScenarioHints = []
    }
    if (editorSettings.validateOnType && activeTab && !isWelcome) void validateEditor()
    const saved = buildCurrentSettingsDTO()
    if (closeDialog) {
      settingsDialogBaseline = null
      showSettings = false
      appendLog(tr('common.settingsSaved'))
      return saved
    }
    settingsDialogBaseline = saved
    return saved
  }

  async function applySettings() {
    await applySettingsCore(true)
  }

  async function applySettingsKeepOpen(): Promise<string> {
    await applySettingsCore(false)
    let message = tr('common.settingsApplied')
    appendLog(message)
    return message
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

  function prepareRecordDialogDefaults(options?: { inheritTestClient?: boolean }) {
    recordMode = 'live'
    recordAppendTo = ''
    recordTestClient = options?.inheritTestClient
      ? (runForm.testClient || testClientSelection || '')
      : ''
    if (activeTab && !isWelcome && activeTab.toLowerCase().endsWith('.feature') && !isUntitled(activeTab)) {
      recordOutput = activeTab.replace(/\\/g, '/')
      recordAppendTo = recordOutput
    } else if (projectPath) {
      recordOutput = `${projectPath.replace(/\\/g, '/')}/recorded.feature`
    }
    recordURL = recordStartURL()
    if (activeTab && !isWelcome) {
      recordFeatureName = basename(activeTab).replace(/\.feature$/i, '')
      recordScenarioName = tr('dialogs.record.scenarioDefault')
    } else {
      recordFeatureName = tr('dialogs.record.featureDefault')
      recordScenarioName = tr('dialogs.record.scenarioDefault')
    }
  }

  function beginRecord() {
    if (playing) {
      appendLog(tr('journal.record.waitForTest'))
      return
    }
    if (!projectPath) {
      appendLog(tr('journal.record.openProjectFirst'))
      return
    }
    prepareRecordDialogDefaults({ inheritTestClient: true })
    showRecord = true
  }

  async function openBrowser() {
    if (!projectPath) {
      appendLog(tr('journal.project.openFirst'))
      return
    }
    if (browserOpen || recording) {
      await focusBrowser()
      return
    }
    prepareRecordDialogDefaults({ inheritTestClient: false })
    recordURL = recordStartURL()
    showRecord = false
    if (recordURL) {
      appendLog(tr('journal.browser.openingUrl', { url: recordURL }))
    } else {
      appendLog(tr('journal.browser.openingEmpty'))
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
      testClient: '',
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
        appendLog(tr('journal.record.sessionError', { error: String(e) }))
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

  async function syncMonacoAfterMount() {
    if (!monaco) return
    if (isWelcome || !activeTab) {
      await monaco.activateTab(null, editorText)
      return
    }
    const tab = tabs.find((t) => t.path === activeTab)
    const text = tab ? tabEditorText(tab) : editorText
    await monaco.activateTab(activeTab, text)
  }

  async function applyLiveRecordedStep(index: number, line: string) {
    if (!line.trim()) return
    await recordEditorReadyPromise
    const targetPath = recordingTargetPath || activeTab
    if (targetPath && activeTab !== targetPath && !isWelcome) {
      await loadFeature(targetPath, { skipRecordingGuard: true })
    }
    if (!line.trim() || isWelcome || !activeTab) return
    recordStepApplyChain = recordStepApplyChain.then(async () => {
      if (!line.trim() || isWelcome || !activeTab) return
      const sourceText = monaco?.getEditorText() ?? editorText
      const result = upsertRecordedStepInText(sourceText, index, line, liveRecordStepLines)
      liveRecordStepLines = result.lineByIndex
      await applyEditorText(result.text, { skipValidate: true })
      scheduleValidateEditor(150)
    })
    await recordStepApplyChain
  }

  async function maybeShowPostRecordBannerAfterStop() {
    const path = (
      recordAppendTo ||
      lastRecordTarget ||
      (activeTab && !isWelcome ? activeTab : '')
    )
      .trim()
      .replace(/\\/g, '/')
    const bannerPath = path && !isUntitled(path) ? path : activeTab && !isWelcome ? activeTab : ''
    if (!bannerPath) return
    const current = monaco?.getEditorText() ?? editorText
    if (!postRecordBaselineText || current === postRecordBaselineText) return
    await showPostRecordBanner(bannerPath)
  }

  function handleRecordStopped(payload?: { reason?: string; idleSeconds?: number }) {
    recording = false
    recordPaused = false
    liveRecordStepLines = {}
    recordingTargetPath = ''
    void maybeShowPostRecordBannerAfterStop()
    if (payload?.reason === 'idle') {
      const sec = payload.idleSeconds ?? recordIdle ?? 30
      appendLog(tr('journal.record.stoppedIdle', { seconds: sec }))
    }
    if (browserOpen) {
      setStatus(payload?.reason === 'idle' ? tr('journal.status.recordStoppedIdle') : tr('journal.status.browserOpen'), 'busy')
      if (payload?.reason !== 'idle') {
        appendLog(tr('journal.record.stoppedBrowserOpen'))
      }
    } else {
      syncIdleStatus()
    }
  }

  async function handleRecordSessionEnd(result: gui.RunResult, kind: 'record' | 'browse') {
    const recordTarget =
      (recordAppendTo || lastRecordTarget || (activeTab && !isWelcome ? activeTab : '')).trim().replace(/\\/g, '/')
    recording = false
    browserOpen = false
    recordPaused = false
    showRecord = false
    liveRecordStepLines = {}
    recordingTargetPath = ''
    lastRecordTarget = ''
    if (result.output) appendLog(result.output)
    if (result.error) appendLog(tr('journal.record.error', { message: result.error || tr('journal.record.unknownError') }))
    else if (kind === 'record') {
      appendLog(tr('journal.browser.closedSaveHint'))
      if (recordTarget && !isUntitled(recordTarget)) {
        await showPostRecordBanner(recordTarget)
      } else if (activeTab && !isWelcome) {
        await showPostRecordBanner(activeTab)
      }
    }
    statusTone = 'normal'
    syncIdleStatus()
  }

  async function toggleRecordPause() {
    pauseToggleGuardUntil = Date.now() + 900
    if (recordPaused) {
      await ResumeRecording()
      recordPaused = false
      setStatus(tr('journal.status.recording'), 'busy')
    } else {
      await PauseRecording()
      recordPaused = true
      setStatus(tr('journal.status.paused'), 'busy')
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
      const guardPause = Date.now() < pauseToggleGuardUntil
      applyBrowserSessionState({
        browserOpen: true,
        recording: s.recording,
        paused: guardPause ? recordPaused : s.paused,
      })
      if (s.recording && !prevRecording) {
        setStatus(tr('journal.status.recording'), 'busy')
      } else if (s.recording && s.paused && (!prevPaused || !prevRecording)) {
        setStatus(tr('journal.status.paused'), 'busy')
      } else if (s.recording && !s.paused && prevPaused) {
        setStatus(tr('journal.status.recording'), 'busy')
      } else if (!s.recording && prevRecording) {
        setStatus(tr('journal.status.browserOpen'), 'busy')
      }
    } catch {
      /* dev without wails */
    }
  }

  async function stopRecord() {
    if (playing) {
      runCancelling = true
      setStatus(tr('journal.status.stoppingTest'), 'busy')
      appendLog(tr('journal.run.stopping'))
      await CancelRun()
      return
    }
    if (recording) {
      await StopRecordingCapture()
      return
    }
    if (browserOpen) {
      if (activeTabUnsaved) {
        const ok = confirm(tr('confirm.closeBrowser.message'))
        if (!ok) return
      }
      await CloseBrowser()
    }
  }

  async function submitOtp(code: string) {
    const accepted = await SubmitOTPCode(code)
    if (accepted) {
      showOtp = false
      return
    }
    appendLog(tr('journal.otp.noActiveRequest'))
  }

  async function cancelOtp() {
    await CancelOTP()
    showOtp = false
  }

  function newScenario() {
    void openUntitledTab(
      buildFeatureTemplate({
        title: tr('dialogs.project.beginnerExamples.title'),
        scenario: tr('dialogs.project.beginnerExamples.scenario'),
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
    activeTab = path
    await applyEditorText(content, { switchTab: true, tabPath: path, skipValidate: true })
    trimTabsMemory()
    stepsPanelCollapsed = resolveStepsPanelCollapsed()
    scheduleValidateEditor()
    schedulePersistSession()
  }

  async function openFileDialog() {
    const path = await PickOpenFile(tr('filePicker.openFeature'))
    if (path) await loadFeature(path)
  }

  function insertTemplate() {
    if (isWelcome) {
      newScenario()
      return
    }
    const template = buildFeatureTemplate({
      title: tr('dialogs.project.beginnerExamples.title'),
      scenario: tr('dialogs.project.beginnerExamples.scenario'),
      startUrl: startURL || recordURL || 'https://example.com',
    })
    monaco?.insertAtCursor(template)
  }

  function quickStart() {
    if (!projectPath) {
      appendLog(tr('journal.project.openFirstForTest'))
      openProjectDialog()
      return
    }
    recordURL = startURL
    beginRecord()
  }

  function dismissWelcomeChecklist() {
    checklistDismissed = true
    void persistSettings()
  }

  function continueRecord() {
    if (!projectPath) {
      appendLog(tr('journal.project.openFirst'))
      return
    }
    if (!activeTab || isWelcome || !activeTab.toLowerCase().endsWith('.feature')) {
      appendLog(tr('journal.record.openFeatureForAppend'))
      beginRecord()
      return
    }
    recordAppendTo = activeTab
    recordOutput = activeTab
    recordFeatureName = basename(activeTab).replace(/\.feature$/i, '')
    recordScenarioName = tr('dialogs.record.appendScenarioDefault')
    recordURL = startURL
    showRecord = true
    appendLog(tr('journal.record.appendTo', { name: basename(activeTab) }))
  }

  function quickRecord() {
    quickStart()
  }

  async function focusBrowser() {
    if (!browserOpen && !recording) {
      appendLog(tr('journal.browser.notOpen'))
      return
    }
    try {
      await FocusBrowser()
      appendLog(tr('journal.browser.focused'))
    } catch (e: unknown) {
      const msg = e instanceof Error ? e.message : String(e)
      appendLog(tr('journal.browser.showFailed', { error: msg }))
      setStatus(tr('journal.status.showBrowserFailed'), 'error')
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

  function closeHttpAuthDialog() {
    showHttpAuth = false
  }

  async function pickElement() {
    if (!recording && !browserOpen) {
      appendLog(tr('journal.picker.openBrowser'))
      return
    }
    if (recording && !recordPaused && !pickerDuringRecording) {
      appendLog(tr('journal.picker.pauseOrSettings'))
      setStatus(tr('journal.status.pauseRecording'), 'busy')
      return
    }
    appendLog(tr('journal.picker.pickInBrowser'))
    await focusBrowser()
    const result = await PickSelector()
    if (result.error) {
      if (result.error !== 'отменено' && result.error !== tr('journal.picker.cancelled')) appendLog(tr('journal.picker.error', { error: result.error }))
      return
    }
    if (!result.selector) return
    pickerSelector = result.selector
    pickerChoices = await PickerStepChoices(result.selector, tr('dialogs.record.gherkinGiven'))
    showPickerStep = true
  }

  function insertPickerStep(text: string) {
    if (recordingBlocksManualTools()) {
      appendLog(tr('journal.record.pauseToInsertStep'))
      return
    }
    monaco?.insertAtCursor(text.endsWith('\n') ? text : text + '\n')
    appendLog(tr('journal.picker.stepInserted'))
  }

  async function undoRecordedStep() {
    const ok = await UndoRecordedStep()
    if (!ok) {
      appendLog(tr('journal.record.undoNothing'))
      return
    }
    const sourceText = monaco?.getEditorText() ?? editorText
    const result = removeLastRecordedStepFromText(sourceText, liveRecordStepLines)
    if (result) {
      liveRecordStepLines = result.lineByIndex
      await applyEditorText(result.text, { skipValidate: true })
    }
    appendLog(tr('journal.record.undoDone'))
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
    if (isWelcome) return tr('editor.welcome')
    if (activeTab && activeTab !== WELCOME_KEY) return featureTabLabel(activeTab)
    if (projectPath) return basename(projectPath)
    return tr('statusBar.default')
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
  <div class="menubar" class:onboarding-elevated={onboardingElevateMenubar} role="menubar">
    <div class="menu-root" class:open={openMenu === 'project'}>
      <button class="menu-trigger" on:click={(e) => toggleMenu('project', e)}>{tr('menus.project')}</button>
      {#if openMenu === 'project'}
        <div class="menu-dropdown">
          <button class="menu-item" on:click={openExamples}>{tr('menus.openExamples')}</button>
          <button class="menu-item" on:click={openNewProjectWizard}>{tr('menus.newProject')}</button>
          <button class="menu-item" on:click={openProjectDialog}>{tr('menus.openProject')}</button>
          <button class="menu-item" on:click={closeProject} disabled={!projectPath}>{tr('menus.closeProject')}</button>
          <button class="menu-item" on:click={openSettings}>{tr('menus.settings')}<span class="menu-shortcut">Ctrl+,</span></button>
          <div class="menu-sep"></div>
          <button class="menu-item" on:click={openInitProjectDialog} disabled={!projectPath}>{tr('menus.initProject')}</button>
        </div>
      {/if}
    </div>

    <div class="menu-root" class:open={openMenu === 'scenario'}>
      <button class="menu-trigger" on:click={(e) => toggleMenu('scenario', e)}>{tr('menus.scenario')}</button>
      {#if openMenu === 'scenario'}
        <div class="menu-dropdown">
          <button class="menu-item" on:click={newScenario}>{tr('menus.new')}</button>
          <button class="menu-item" on:click={openFileDialog}>{tr('menus.open')}</button>
          <button class="menu-item" on:click={saveFeature} disabled={isWelcome}>{tr('menus.save')}<span class="menu-shortcut">Ctrl+S</span></button>
          <button class="menu-item" on:click={saveFeatureAs} disabled={isWelcome}>{tr('menus.saveAs')}</button>
          <button class="menu-item" on:click={() => activeTab && !isWelcome && openDuplicateDialog(activeTab)} disabled={isWelcome}>{tr('menus.duplicate')}</button>
          <button class="menu-item" on:click={openImportFeaturesDialog} disabled={!projectPath}>{tr('menus.importFeature')}</button>
          <button
            class="menu-item"
            on:click={() => {
              if (!activeTab || isWelcome) return
              renameFeaturePath = activeTab
              showRenameFeature = true
            }}
            disabled={isWelcome}
          >
            {tr('menus.rename')}
          </button>
          <button class="menu-item" on:click={() => activeTab && !isWelcome && deleteFeature(activeTab)} disabled={isWelcome}>{tr('menus.delete')}</button>
          <div class="menu-sep"></div>
          <button class="menu-item" on:click={openFindReplace} disabled={isWelcome}>{tr('menus.findReplace')}<span class="menu-shortcut">Ctrl+H</span></button>
          <button class="menu-item" on:click={() => (showProjectReplace = true)} disabled={!projectPath}>{tr('menus.projectReplace')}</button>
          <div class="menu-sep"></div>
          <button class="menu-item" on:click={openExportDialog} disabled={isWelcome}>{tr('menus.export')}</button>
          <button class="menu-item" on:click={openImportDialog} disabled={!projectPath}>{tr('menus.importJson')}</button>
          <div class="menu-sep"></div>
          <button class="menu-item" on:click={openSnippetPalette}>{tr('menus.snippetPalette')}<span class="menu-shortcut">Ctrl+Shift+Space</span></button>
          <button class="menu-item" on:click={openStepsDialog}>{tr('menus.insertStep')}</button>
          <div class="menu-sep"></div>
          <button class="menu-item" on:click={refactorUpdateUrls} disabled={isWelcome}>{tr('menus.updateStartUrl')}</button>
          <button class="menu-item" on:click={refactorNormalizeIndents} disabled={isWelcome}>{tr('menus.normalizeIndents')}</button>
          <button class="menu-item" on:click={refactorCollapseBlank} disabled={isWelcome}>{tr('menus.collapseBlanks')}</button>
        </div>
      {/if}
    </div>

    <div class="menu-root" class:open={openMenu === 'run'} data-tour="menu-run">
      <button class="menu-trigger" data-tour="menu-run-trigger" on:click={(e) => toggleMenu('run', e)}>{tr('menus.runMenu')}</button>
      {#if openMenu === 'run'}
        <div class="menu-dropdown" data-tour="menu-run-dropdown">
          <button class="menu-item" on:click={() => runMenuAction(() => void openBrowser())} disabled={!projectPath}>{tr('menus.browser')}<span class="menu-shortcut">Ctrl+B</span></button>
          <button class="menu-item" on:click={() => runMenuAction(beginRecord)} disabled={!projectPath}>{tr('menus.record')}<span class="menu-shortcut">Ctrl+R</span></button>
          <button class="menu-item" on:click={openBaselineRecordDialog} disabled={!projectPath}>{tr('menus.recordFromSteps')}</button>
          <button class="menu-item" on:click={stopRecord} disabled={!recording && !browserOpen && !playing}>{tr('menus.stop')}</button>
          <button class="menu-item" on:click={toggleRecordPause} disabled={!recording}>{tr('menus.pause')}</button>
          <div class="menu-sep"></div>
          <button class="menu-item" on:click={openTestClientDialog} disabled={!projectPath}>{tr('menus.testClient')}</button>
          <button class="menu-item" on:click={openTestClientDialogForCapture} disabled={!projectPath || (!browserOpen && !recording)}>
            {tr('menus.saveBrowserSession')}
          </button>
          <button class="menu-item" on:click={openHttpAuthDialog}>{tr('menus.httpAuth')}</button>
          <div class="menu-sep"></div>
          <button class="menu-item" on:click={() => runPrimary(false)} disabled={isWelcome && !projectPath && !batchSelected.length}>
            {tr('menus.run')}<span class="menu-shortcut">Ctrl+Enter</span>
          </button>
          <button class="menu-item" on:click={() => runCurrentScenario(false)} disabled={isWelcome || !activeTab}>
            {tr('menus.runCurrent')}<span class="menu-shortcut">Ctrl+Shift+Enter</span>
          </button>
          <button class="menu-item" on:click={() => runCurrentScenario(true)} disabled={isWelcome || !activeTab}>
            {tr('menus.dryRunCurrent')}
          </button>
          <button class="menu-item" on:click={() => runBatchSelected(false)} disabled={!projectPath || !batchSelected.length}>
            {tr('menus.runSelected')}
          </button>
          <button class="menu-item" on:click={() => runBatchSelected(true)} disabled={!projectPath || !batchSelected.length}>
            {tr('menus.dryRunSelected')}
          </button>
          <button class="menu-item" on:click={rerunFailed} disabled={!projectPath}>{tr('menus.rerunFailed')}</button>
          <button class="menu-item" on:click={openRunHistory} disabled={!projectPath}>{tr('menus.runHistory')}</button>
          <button class="menu-item" on:click={() => openRunDialog('', {})} disabled={isWelcome && !projectPath}>{tr('menus.runDialog')}</button>
          <button class="menu-item" on:click={() => openRunDialog(tr('menus.runTag').replace('…', ''), {})} disabled={isWelcome && !projectPath}>
            {tr('menus.runTag')}
          </button>
          <button class="menu-item" data-tour="menu-dry-run" on:click={() => runPrimary(true)} disabled={isWelcome && !projectPath && !batchSelected.length}>{tr('menus.dryRun')}</button>
          <button
            class="menu-item"
            on:click={() => openRunDialog('Playwright', { dryRun: false, headed: true, engine: 'playwright', installPW: true })}
            disabled={isWelcome && !activeTab}
          >
            {tr('menus.playwright')}
          </button>
          <div class="menu-sep"></div>
          <button class="menu-item" data-tour="menu-validate" on:click={() => openValidateDialog(true)} disabled={!projectPath}>{tr('menus.validate')}</button>
          <button class="menu-item" on:click={() => openValidateDialog(false)} disabled={!projectPath}>{tr('menus.validateBrowser')}</button>
          <div class="menu-sep"></div>
          {#if hasVanessaPlugin()}
            <button class="menu-item" on:click={() => openVanessaDialog(true)} disabled={!projectPath}>{tr('menus.vanessaDry')}</button>
            <button class="menu-item" on:click={() => openVanessaDialog(false)} disabled={!projectPath}>{tr('menus.vanessaRun')}</button>
            <button class="menu-item" on:click={() => openVanessaDialog(false, true)} disabled={!projectPath}>{tr('menus.vanessaRerun')}</button>
            <button class="menu-item" on:click={openVanessaSettingsDialog} disabled={!projectPath}>{tr('menus.vanessaSettings')}</button>
            <button class="menu-item" on:click={openVanessaMonitor} disabled={!projectPath}>{tr('menus.vanessaMonitor')}</button>
          {/if}
        </div>
      {/if}
    </div>

    <div class="menu-root" class:open={openMenu === 'plugins'}>
      <button class="menu-trigger" on:click={(e) => toggleMenu('plugins', e)}>{tr('menus.pluginsMenu')}</button>
      {#if openMenu === 'plugins'}
        <div class="menu-dropdown">
          <button class="menu-item" on:click={() => (showPlugins = true)} disabled={!projectPath}>{tr('menus.plugins')}</button>
          {#if installedPlugins.length > 0}
            <div class="menu-sep"></div>
            {#each installedPlugins as plugin (plugin.name)}
              {#if plugin.runnable}
                {#if plugin.vanessa}
                  <button class="menu-item" on:click={() => openVanessaDialog(true)} disabled={!projectPath}>{tr('menus.vanessaDry')}</button>
                  <button class="menu-item" on:click={() => openVanessaDialog(false)} disabled={!projectPath}>{tr('menus.vanessaRun')}</button>
                  <button class="menu-item" on:click={() => openVanessaDialog(false, true)} disabled={!projectPath}>{tr('menus.vanessaRerun')}</button>
                {:else}
                  <button class="menu-item" on:click={() => openPluginRun(plugin.name, true)} disabled={!projectPath}>{tr('menus.runPluginDry', { name: pluginLabel(plugin) })}</button>
                  <button class="menu-item" on:click={() => openPluginRun(plugin.name, false)} disabled={!projectPath}>{tr('menus.runPlugin', { name: pluginLabel(plugin) })}</button>
                {/if}
              {:else}
                <button class="menu-item" disabled title={tr('menus.pluginNoRunTitle')}>{tr('menus.pluginNoRun', { name: pluginLabel(plugin) })}</button>
              {/if}
            {/each}
          {:else if projectPath}
            <div class="menu-sep"></div>
            <button class="menu-item" disabled>{tr('menus.noInstalledPlugins')}</button>
          {/if}
        </div>
      {/if}
    </div>

    <div class="menu-root" class:open={openMenu === 'view'}>
      <button class="menu-trigger" on:click={(e) => toggleMenu('view', e)}>{tr('menus.view')}</button>
      {#if openMenu === 'view'}
        <div class="menu-dropdown">
          <button class="menu-item" on:click={() => selectTab(WELCOME_KEY)}>{tr('menus.start')}</button>
          <button class="menu-item" on:click={() => { sidebarVisible = true; saveLayout({ sidebarVisible: true }) }}>{tr('menus.scenarios')}</button>
          <button class="menu-item" on:click={() => { sidebarVisible = false; saveLayout({ sidebarVisible: false }) }}>{tr('menus.hideExplorer')}</button>
          <button class="menu-item" on:click={() => (showCommandPalette = true)}>{tr('palette.commands.palette')}<span class="menu-shortcut">Ctrl+Shift+P</span></button>
          <button class="menu-item" on:click={() => { bottomPanelOpen = true; bottomTab = 'journal' }}>{tr('menus.journal')}</button>
          <button class="menu-item" on:click={() => { bottomPanelOpen = true; bottomTab = 'results' }}>{tr('menus.resultsPanel')}</button>
          <button class="menu-item" on:click={openRunHistory} disabled={!projectPath}>{tr('menus.runHistory')}</button>
          <button class="menu-item" on:click={() => { bottomPanelOpen = true; bottomTab = 'validate' }}>{tr('menus.validatePanel')}</button>
          <button class="menu-item" on:click={() => { bottomPanelOpen = true; bottomTab = 'error' }}>{tr('menus.errorPanel')}</button>
          <button class="menu-item" on:click={togglePreview}>
            {previewVisible ? tr('menus.hidePreview') : tr('menus.showPreview')}
          </button>
          <button class="menu-item" on:click={toggleStepsPanel}>
            {stepsPanelVisible ? tr('menus.hideStepsPanel') : tr('menus.showStepsPanel')}
          </button>
          <button class="menu-item" on:click={() => { toolbarCompact = !toolbarCompact; persistSettings() }}>
            {toolbarCompact ? tr('menus.expandedToolbar') : tr('menus.compactToolbar')}
          </button>
          <button class="menu-item" on:click={resetWindowLayout}>{tr('menus.resetLayout')}</button>
        </div>
      {/if}
    </div>

    <div class="menu-root" class:open={openMenu === 'help'}>
      <button class="menu-trigger" on:click={(e) => toggleMenu('help', e)}>{tr('menus.help')}</button>
      {#if openMenu === 'help'}
        <div class="menu-dropdown">
          <button class="menu-item" on:click={restartOnboardingTour}>{tr('menus.training')}</button>
          <div class="menu-sep"></div>
          <button class="menu-item" on:click={() => openStepsHelp()}>{tr('menus.stepsHelp')}<span class="menu-shortcut">F1</span></button>
          <button class="menu-item" on:click={() => (showHotkeys = true)}>{tr('menus.hotkeys')}<span class="menu-shortcut">Shift+F1</span></button>
          <button class="menu-item" on:click={checkUpdates}>{tr('menus.checkUpdates')}</button>
          <button class="menu-item" on:click={showAboutDialog}>{tr('menus.about')}</button>
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
          title={tr('menus.scenarios')}
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
          title={tr('catalog.outputPanel')}
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
        <div class="sidebar-column" class:onboarding-elevated={onboardingElevateSidebar} style="width: {layoutSidebarWidth + 4}px">
        <aside class="explorer" style="width: {layoutSidebarWidth}px">
          <div class="explorer-header">
            <p class="zone-title">{tr('catalog.title')}</p>
            <div class="explorer-tools">
              <input class="explorer-search" value={sidebarSearch} placeholder={tr('catalog.searchPlaceholder')} on:input={onSidebarSearchInput} />
              <div class="explorer-tool-actions">
                <button class="icon-btn" title={tr('catalog.newScenario')} on:click={newScenario}>{@html icons.plus}</button>
                <button class="icon-btn batch-toggle" class:active={batchMode} title={tr('catalog.folder.selectBatch')} on:click={toggleBatchMode}>
                  <span class="batch-toggle-icon" aria-hidden="true">{@html icons.validate}</span>
                  <span class="batch-label">{tr('catalog.batchMode')}</span>
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
              <p class="batch-selection">{tr('catalog.batchSelected', { count: batchCount })}</p>
            {/if}
          </div>
          <div class="catalog catalog-panel" data-tour="catalog-tree">
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
              <button class="tool-btn primary" title={tr('toolbar.browser') + ' (Ctrl+B)'} on:click={() => void openBrowser()} disabled={!projectPath}>
                {@html toolbarIcons.browser()}<span>{tr('toolbar.browser')}</span>
              </button>
              <button
                class="tool-btn primary"
                class:record-active={recording}
                title={tr('toolbar.record') + ' (Ctrl+R)'}
                on:click={beginRecord}
                disabled={!projectPath || recording}
              >
                {@html toolbarIcons.record()}<span>{tr('toolbar.record')}</span>
              </button>
              <button class="tool-btn primary" on:click={stopRecord} disabled={!recording && !browserOpen && !playing} title={stopActionLabel}>
                {@html toolbarIcons.stop()}<span>{stopActionLabel}</span>
              </button>
              <button
                class="tool-btn primary primary-run"
                class:run-active={playing}
                title={batchSelected.length > 0
                  ? tr('toolbar.runSelectedTitle', { count: batchSelected.length })
                  : tr('toolbar.runTitle')}
                on:click={() => runPrimary(false)}
                disabled={(isWelcome && !projectPath && !batchSelected.length) || playing}
              >
                {@html toolbarIcons.play()}<span>{batchSelected.length > 0 ? tr('toolbar.runSelected', { count: batchSelected.length }) : tr('toolbar.run')}</span>
              </button>
              <button
                class="tool-btn"
                title={tr('toolbar.scenarioTitle')}
                on:click={() => runCurrentScenario(false)}
                disabled={isWelcome || !activeTab}
              >
                {@html toolbarIcons.play()}<span>{tr('toolbar.scenario')}</span>
              </button>
              <button class="tool-btn primary" on:click={saveFeature} disabled={isWelcome} title={tr('toolbar.saveTitle')}>
                {@html toolbarIcons.save()}<span>{tr('toolbar.save')}</span>
              </button>
            </div>
            {#if !actionBarCompact}
            <div class="toolbar-row secondary" class:icon-only={toolbarIconOnly}>
              <button class="tool-btn" on:click={continueRecord} disabled={!projectPath || recording} title={tr('toolbar.continueRecordTitle')}>
                {@html toolbarIcons.continueRecord()}<span>{tr('toolbar.continueRecord')}</span>
              </button>
              <button class="tool-btn" on:click={toggleRecordPause} disabled={!recording} title={tr('toolbar.pauseTitle')}>
                {@html toolbarIcons.pause()}<span>{tr('toolbar.pause')}</span>
              </button>
              <span class="toolbar-sep" aria-hidden="true"></span>
              <button class="tool-btn" on:click={() => void focusBrowser()} disabled={!browserOpen && !recording}>
                {@html toolbarIcons.browserFocus()}<span>{tr('toolbar.showBrowser')}</span>
              </button>
              <button class="tool-btn" on:click={() => openValidateDialog(false)} disabled={!projectPath}>
                {@html toolbarIcons.validate()}<span>{tr('toolbar.selectorsOnPage')}</span>
              </button>
              <button class="tool-btn" on:click={pickElement} disabled={!pickerToolbarEnabled} title={recording && !recordPaused && !pickerDuringRecording ? tr('toolbar.pickElementPauseHint') : tr('toolbar.pickElement')}>
                {@html toolbarIcons.picker()}<span>{tr('toolbar.pickElement')}</span>
              </button>
              <button class="tool-btn" on:click={quickRecord} disabled={!projectPath || recording}>
                {@html toolbarIcons.quickRecord()}<span>{tr('toolbar.quickRecord')}</span>
              </button>
              <button class="tool-btn" on:click={validateEditor} disabled={isWelcome}>
                {@html toolbarIcons.gherkin()}<span>{tr('toolbar.gherkinSyntax')}</span>
              </button>
              <button class="tool-btn" on:click={undoRecordedStep} disabled={!recording}>
                {@html toolbarIcons.undo()}<span>{tr('toolbar.undoStep')}</span>
              </button>
              <span class="toolbar-sep" aria-hidden="true"></span>
              <button class="tool-btn" on:click={() => { bottomPanelOpen = true; bottomTab = 'journal' }}>
                {@html toolbarIcons.log()}<span>{tr('toolbar.journal')}</span>
              </button>
              <button class="tool-btn" on:click={() => { bottomPanelOpen = true; bottomTab = 'results' }}>
                {@html toolbarIcons.results()}<span>{tr('toolbar.results')}</span>
              </button>
            </div>
            {/if}
          </div>

          <div class="url-block">
            <span>URL</span>
            <input bind:value={startURL} placeholder="https://site.com" on:change={syncUrlFromField} />
            <button class="icon-btn" title={tr('toolbar.urlFromBrowser')} on:click={() => (recordURL = startURL)}>
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
              tourElevated={onboardingElevateWelcome}
              bind:startURL
              {recentProjects}
              {recentFeatures}
              projectOpen={welcomeProjectOpen}
              recorded={welcomeRecorded}
              playedSuccess={welcomePlayedSuccess}
              checklistDismissed={checklistDismissed}
              onDismissChecklist={dismissWelcomeChecklist}
              onOpenProject={openProjectDialog}
              onNewProject={openNewProjectWizard}
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
          <div class="feature-workspace" class:hidden={isWelcome} data-tour="editor-workspace">
            {#if postRecordPath}
              <PostRecordBanner
                path={postRecordPath}
                stepCount={postRecordStepCount}
                onValidate={postRecordValidate}
                onSave={postRecordSave}
                onShowDiff={openPostRecordDiff}
                onClose={dismissPostRecord}
              />
            {/if}
            {#if activeTabUnsaved}
              <div class="dirty-banner">
                <span>
                  {#if isUntitled(activeTab)}
                    {tr('editor.unsavedScenario')}
                  {:else}
                    {tr('editor.unsavedChanges')}
                  {/if}
                </span>
                <button class="primary" on:click={saveFeature}>
                  {isUntitled(activeTab) ? tr('filePicker.saveAs') + '…' : tr('toolbar.save')}
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
                <span>{tr('editor.syntaxErrors')}</span>
              </div>
            {/if}
            {#if showPlayingBar}
              <div class="playing-bar" role="status" aria-live="polite">
                <span class="play-label">{tr('editor.running')}</span>
                <span class="play-target">{playingLabel}</span>
                {#if runProgressTotal > 1}
                  <span class="play-progress-text">{runProgressCurrent}/{runProgressTotal}</span>
                {/if}
                {#if settingsSlowMo > 0 || lastRun.slowMo > 0}
                  <span class="play-slowmo">slow-mo {(lastRun.slowMo > 0 ? lastRun.slowMo : settingsSlowMo)} {tr('settings.units.ms')}</span>
                {/if}
                {#if runCancelling}
                  <span class="play-cancel">{tr('editor.stopping')}</span>
                {:else}
                  <button type="button" class="play-cancel-btn" on:click={stopRecord}>{tr('editor.cancel')}</button>
                {/if}
                <span class="play-progress" aria-hidden="true" style="--run-progress: {runProgressTotal > 0 ? (runProgressCurrent / runProgressTotal) * 100 : 0}%"></span>
              </div>
            {/if}
            {#if showRecordingBar}
              <div class="recording-bar">
                <span class="rec-label">{tr('editor.recording')}</span>
                <span class="rec-hint" title={tr('editor.recordingUndoHint')}>{tr('editor.recordingUndoHint')}</span>
                <label class="check-inline"><input type="checkbox" bind:checked={filterRecording} on:change={() => { if (filterRecording) navOnlyRecording = false; void syncRecordingOptions() }} /> {tr('editor.recordingFilters.importantOnly')}</label>
                <label class="check-inline"><input type="checkbox" bind:checked={navOnlyRecording} on:change={() => { if (navOnlyRecording) filterRecording = false; void syncRecordingOptions() }} /> {tr('editor.recordingFilters.linksOnly')}</label>
                <label class="check-inline"><input type="checkbox" checked={settingsHeadless} on:change={onRecordingHeadlessChange} /> {tr('toolbar.headless')}</label>
                <label class="check-inline"><input type="checkbox" bind:checked={hoverRecord} on:change={() => void syncRecordingOptions()} /> {tr('editor.recordingFilters.recordHover')}</label>
              </div>
            {/if}
            <div class="gherkin-hints">
              <span class="hint-summary" title={tr('editor.hintSummary', { count: stepCount })}>
                {tr('editor.hintSummary', { count: stepCount })}
              </span>
              <div class="hint-actions">
                <button on:click={openStepsDialog}>{tr('editor.template')}</button>
                <button on:click={() => openStepsHelp()}>{tr('editor.help')}</button>
                <button class:active={previewVisible} on:click={togglePreview}>{tr('menus.showPreview')}</button>
              </div>
            </div>
            <div class="editor-row" style="--preview-width: {layoutPreviewWidth}px">
              <div class="editor-main">
                <div class="editor-area" class:playing-active={playing || vanessaRunning} class:dry-run-active={runningDryRun}>
                    <MonacoEditor
                      bind:this={monaco}
                      bind:value={editorText}
                      readOnly={automationActive}
                      bind:editorSettings
                      scenarioHints={editorScenarioHints}
                      hintActions={monacoHintActions}
                      runLensActions={monacoRunLensActions}
                      inlayHintsHandlers={monacoInlayHintsHandlers}
                      on:ready={() => void syncMonacoAfterMount()}
                      on:change={(e) => onEditorChange(e.detail)}
                      on:cursorline={(e) => (editorCursorLine = e.detail)}
                    />
                </div>
                {#if stepsPanelVisible}
                <div class="splitter-h" role="separator" on:mousedown={startResizeSteps}></div>
                <div class="steps-panel" class:collapsed={stepsPanelCollapsed} style="max-height: {stepsPanelCollapsed ? 24 : stepsPanelHeight}px">
                  <div
                    class="steps-header"
                    title={tr('editor.stepsPanel.tooltip')}
                  >
                    <button
                      type="button"
                      aria-label={stepsPanelCollapsed ? tr('editor.stepsPanel.expand') : tr('editor.stepsPanel.collapse')}
                      on:click={() => (stepsPanelCollapsed = !stepsPanelCollapsed)}
                    >
                      {#if stepsPanelCollapsed}{@html icons.chevronRight}{:else}{@html icons.chevronDown}{/if}
                    </button>
                    {#if editorSettings.symbolOutline}
                      <div class="steps-panel-tabs" role="tablist" aria-label={tr('editor.stepsPanel.aria')}>
                        <button
                          type="button"
                          role="tab"
                          class:active={stepsPanelTab === 'outline'}
                          aria-selected={stepsPanelTab === 'outline'}
                          on:click={() => (stepsPanelTab = 'outline')}
                        >
                          {tr('editor.stepsPanel.outline')}
                        </button>
                        <button
                          type="button"
                          role="tab"
                          class:active={stepsPanelTab === 'steps'}
                          aria-selected={stepsPanelTab === 'steps'}
                          on:click={() => (stepsPanelTab = 'steps')}
                        >
                          {tr('editor.stepsPanel.title', { count: stepCount })}
                        </button>
                      </div>
                    {:else}
                      <span>{tr('editor.stepsPanel.title', { count: stepCount })}</span>
                    {/if}
                    {#if stepStatusError}<span class="steps-header-error">{tr('editor.stepsPanel.errors')}</span>{/if}
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
                          <tr><th>#</th><th>{tr('editor.stepsPanel.action')}</th><th>{tr('editor.stepsPanel.element')}</th><th>{tr('editor.stepsPanel.value')}</th></tr>
                        </thead>
                        <tbody>
                          {#each editorSteps as step, i}
                            <tr
                              class:error={!!step.error}
                              class:clickable={!!step.line}
                              on:click={() => step.line && gotoEditorLine(step.line)}
                              on:dblclick={() => openStepHelpFromPanel(step)}
                              on:contextmenu|preventDefault={(e) => openStepsContextMenu(e, step)}
                              title={step.text ? tr('dialogs.steps.contextMenu.contextMenuTitle') : ''}
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
              {#if showPreviewPane && previewPaneMounted}
                <div class="splitter-v" role="separator" on:mousedown={startResizePreview}></div>
                <div class="feature-preview-pane" style="width: {layoutPreviewWidth}px">
                  <div class="preview-header">{tr('editor.preview')}</div>
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
  <div class="bottom-dock" class:collapsed={!bottomPanelOpen} class:onboarding-elevated={onboardingElevateBottom}>
    <div class="splitter-h bottom-splitter" role="separator" on:mousedown={startResizeBottom}></div>
    <div class="bottom-panel" style="--panel-height: {bottomPanelHeight}px">
    <div class="panel-tabs">
      <button class="panel-tab" class:active={bottomTab === 'journal'} data-tour="panel-journal" on:click={() => openJournalTab(true)}>{tr('panels.journal')}</button>
      <button class="panel-tab" class:active={bottomTab === 'results'} on:click={() => (bottomTab = 'results')}>{tr('panels.results')}</button>
      <button class="panel-tab" class:active={bottomTab === 'validate'} on:click={() => (bottomTab = 'validate')}>{tr('panels.validate')}</button>
      <button class="panel-tab" class:active={bottomTab === 'error'} on:click={() => (bottomTab = 'error')}>{tr('panels.error')}</button>
    </div>
    <div class="panel-body" class:muted={bottomTab === 'journal' && !logText} class:text-panel={bottomTab === 'journal'}>
      {#if bottomTab === 'journal'}
        {logText || tr('panels.journalEmpty')}
      {:else if bottomTab === 'results'}
        <ResultsPanel
          entries={resultsPanelEntries}
          flakyByPath={flakyByPath}
          flakyStepByPath={flakyStepByPath}
          artifacts={projectArtifacts}
          onOpenFeature={openFeatureFromHistory}
          onGotoFailedStep={gotoFailedStep}
          onRerun={rerunFailed}
          onRunFlaky={runFlakyTriplicate}
          onOpenFolder={openArtifactPath}
          onServeAllure={serveAllureReport}
          allureInstalled={allureInstalled}
          allureRunning={allureServeRunning}
          onOpenAllureInstall={openAllureInstallHelp}
          onOpenHtmlReport={openHtmlReport}
          onOpenTrace={openTraceReport}
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
        <ErrorPanel entry={lastErrorEntry} onGotoFailedStep={gotoFailedStep} />
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
        {tr('statusBar.browser', { state: recording ? tr('statusBar.browserRecording') : playing ? tr('statusBar.browserTesting') : browserOpen ? tr('statusBar.browserOpen') : tr('statusBar.browserClosed') })}
      </div>
      {#if recordingTargetLabel}
        <div class="status-segment recording-target" title={recordingTargetPath}>
          {tr('statusBar.recordingTarget', { label: recordingTargetLabel })}
        </div>
      {/if}
      {#if unsavedTabCount > 1}
        <div class="status-segment warning" title={tr('statusBar.unsavedTabsTitle')}>
          {tr('statusBar.unsavedTabs', { count: unsavedTabCount })}
        </div>
      {/if}
      {#if lastRunSummary && !playing && !runningDryRun}
        <div class="status-segment muted run-summary" title={tr('statusBar.lastRunTitle')}>
          {tr('statusBar.lastRun', { summary: lastRunSummary })}
        </div>
      {/if}
      {#if runningDryRun}
        <div class="status-segment muted">{tr('statusBar.dryRun')}</div>
      {/if}
      {#if showLargeFileBanner}
        <div class="status-segment warning large-file-banner" title={tr('statusBar.largeFileTitle')}>
          {tr('statusBar.largeFile', { lines: LARGE_FILE_LINE_THRESHOLD })}
        </div>
      {/if}
      <div class="status-segment muted">{tr('statusBar.runner')}</div>
      <div class="status-segment" class:warning={stepStatusError}>{stepStatusDisplay}</div>
      <button type="button" class="status-segment clickable" on:click={() => { bottomPanelOpen = true; bottomTab = 'journal' }}>
        {@html icons.log} {tr('panels.journal')}
      </button>
      <button type="button" class="status-segment clickable project-segment" on:click={() => (isWelcome ? selectTab(WELCOME_KEY) : openProjectDialog())}>
        {@html icons.explorer} {projectLabel()}
      </button>
    </div>
  </footer>
</div>

<OnboardingTour
  bind:this={onboardingTour}
  active={onboardingTourActive}
  context={onboardingTourContext}
  on:skip={dismissOnboarding}
  on:complete={completeOnboarding}
  on:stepChange={onOnboardingStepChange}
/>

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
      askConfirm({ title: tr('confirm.generic.title'), message, confirmLabel: tr('confirm.generic.confirmLabelDelete'), danger: true })}
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

{#if showNewProjectWizard}
  <NewProjectWizardDialog
    defaultStartUrl={startURL || 'https://example.com'}
    onConfirm={confirmNewProjectWizard}
    onCancel={() => (showNewProjectWizard = false)}
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
    bind:navWaitUntil={settingsNavWaitUntil}
    bind:pickerDuringRecording
    bind:uiLocale
    bind:editorSettings
    onSave={applySettings}
    onApply={applySettingsKeepOpen}
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
    bind:stepPickerOpen={recordStepPickerOpen}
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
    childModalOpen={showHttpAuth}
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
    childModalOpen={showPluginRun}
    onClose={() => {
      showPlugins = false
      void refreshInstalledPlugins()
    }}
    onRunPlugin={(name, dry) => openPluginRun(name, dry)}
    onAskConfirm={(message) =>
      askConfirm({ title: tr('confirm.generic.title'), message, confirmLabel: tr('confirm.generic.confirmLabelDelete'), danger: true })}
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

{#if showRunHistory}
  <RunHistoryDialog
    entries={runResults}
    flakyByPath={flakyByPath}
    flakyStepByPath={flakyStepByPath}
    onOpenFeature={openFeatureFromHistory}
    onRerunFailed={() => { showRunHistory = false; rerunFailed() }}
    onClose={() => (showRunHistory = false)}
  />
{/if}

{#if showPostRecordDiff && postRecordPath}
  <PostRecordDiffDialog
    path={postRecordPath}
    original={postRecordBaselineText}
    modified={monaco?.getEditorText() ?? editorText}
    onClose={() => (showPostRecordDiff = false)}
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
  <HttpAuthDialog initialHost={httpAuthHost} onCancel={closeHttpAuthDialog} />
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
    dontAskAgainLabel={confirmDialog.dontAskAgainLabel ?? ''}
    onConfirm={(dontAskAgain) => closeConfirm(true, dontAskAgain)}
    onCancel={() => closeConfirm(false)}
  />
{/if}

<BrowserOverlay
  visible={showBrowserOverlay}
  browserOpen={browserOpen}
  {recording}
  {playing}
  paused={recordPaused}
  onRecord={beginRecord}
  onPause={toggleRecordPause}
  onStop={stopRecord}
  onPicker={pickElement}
  onFocusBrowser={focusBrowserWindow}
  pickerDuringRecording={pickerDuringRecording}
/>
{/if}
