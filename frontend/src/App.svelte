<script lang="ts">
  import { onDestroy, onMount, tick } from 'svelte'
  import { get } from 'svelte/store'
  import MonacoEditor from './lib/MonacoEditor.svelte'
  import WelcomePanel from './lib/WelcomePanel.svelte'
  import CatalogEmptyState from './lib/CatalogEmptyState.svelte'
  import FeatureCatalogTree from './lib/FeatureCatalogTree.svelte'
  import EditorTabBar from './lib/EditorTabBar.svelte'
  import { buildCatalogViewState, buildRunByPathMap, collectFeaturePathsUnder, type CatalogNode } from './lib/catalogTree'
  import {
    buildBatchSelectedSet,
    selectAllFeaturesUnder,
  } from './lib/batchSelection'
  import { debounce } from './lib/uiScheduler'
  import { perfMark, perfNow } from './lib/perfMark'
  import { clearFeatureSymbolCache, evictFeatureSymbolCache } from './lib/featureSymbolCache'
  import {
    bindWailsEvents,
    createEventStalenessGuards,
    createStaleEventLogger,
    formatRunProgressLabel,
    shouldRefreshRunResultsFromProgress,
  } from './controllers/wailsEventsController'
  import { createDialogBindController } from './controllers/dialogBindController'
  import { buildPaletteCommands, type PaletteActions } from './controllers/paletteCommandsController'
  import { createWorkspaceSessionController, hasRestorableWorkspaceSession } from './controllers/workspaceSessionController'
  import { pickPersistText } from './lib/sessionTabs'
  import { createProjectStore, type ProjectState } from './stores/projectStore'
  import { createRunnerStore } from './stores/runnerStore'
  import { createDiagnosticsStore } from './stores/diagnosticsStore'
  import { createRecorderStore } from './stores/recorderStore'
  import { createReportsStore } from './stores/reportsStore'
  import { anyAppDialogOpen as computeAnyAppDialogOpen, createDialogsStore } from './stores/dialogsStore'
  import { createRecordFormStore } from './stores/recordFormStore'
  import { createLayoutStore } from './stores/layoutStore'
  import { createUiPrefsStore } from './stores/uiPrefsStore'
  import { createOnboardingTourStore } from './stores/onboardingTourStore'
  import { createRecorderPrefsStore } from './stores/recorderPrefsStore'
  import { createSettingsStore } from './stores/settingsStore'
  import { createRunFormStore } from './stores/runFormStore'
  import { createVanessaRunStore } from './stores/vanessaRunStore'
  import { createPluginRunStore } from './stores/pluginRunStore'
  import { createEditorStore } from './stores/editorStore'
  import { createRunDialogStore } from './stores/runDialogStore'
  import { createValidateDialogStore } from './stores/validateDialogStore'
  import { createJournalStore } from './stores/journalStore'
  import { createFeatureDialogStore } from './stores/featureDialogStore'
  import { createTestClientStore } from './stores/testClientStore'
  import { createPluginsStore } from './stores/pluginsStore'
  import { createUpdateDialogStore } from './stores/updateDialogStore'
  import { createCatalogStore } from './stores/catalogStore'
  import { createContextMenuStore } from './stores/contextMenuStore'
  import { createPostRecordStore } from './stores/postRecordStore'
  import { createConfirmDialogStore } from './stores/confirmDialogStore'
  import { createMenuStore } from './stores/menuStore'
  import { createProjectReplaceStore } from './stores/projectReplaceStore'
  import { createPickerDialogStore } from './stores/pickerDialogStore'
  import { createHttpAuthDialogStore } from './stores/httpAuthDialogStore'
  import { createStepsHelpDialogStore } from './stores/stepsHelpDialogStore'
  import { createOtpDialogStore } from './stores/otpDialogStore'
  import { createSplashStore } from './stores/splashStore'
  import { createViewportStore } from './stores/viewportStore'
  import { createRecentsStore } from './stores/recentsStore'
  import { createSettingsDialogStore } from './stores/settingsDialogStore'
  import { createAppMetaStore } from './stores/appMetaStore'
  import { createSessionStore } from './stores/sessionStore'
  import {
    MAX_OPEN_EDITOR_TABS,
    pathsToRetainModels,
    tabEditorText,
    tabNeedsDiskReload,
    trimRetainedTabBodies,
  } from './lib/tabMemory'
  import { canShowWelcome, createTabsStore } from './stores/tabsStore'
  import SettingsDialog from './lib/SettingsDialog.svelte'
  import CommandPalette from './lib/CommandPalette.svelte'
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
  import { applyRecordStepEvent, type RecordStepEvent } from './lib/recordedStepOps'
  import { isRealFeaturePath, isUntitled, makeUntitledPath, untitledLabel } from './lib/untitled'
  import { matchHotkey, monacoOverlayConsumesEscape, shouldIgnoreAppHotkey, type HotkeyId } from './lib/hotkeys'
  import { batchRunFormFrom, currentScenarioRunFormFrom, defaultRunForm, runFormFromMode, type RunForm, type RunFormMode } from './lib/runTypes'
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
  import { isLargeFeatureFile, LARGE_FILE_LINE_THRESHOLD } from './lib/editorLargeFile'
  import { isEditorAnalysisSnapshotVisible } from './lib/editorAnalysisSync'
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
  import { shouldAcceptMonacoChange } from './lib/monacoHydration'
  import { setStepHoverEnabled } from './lib/gherkinStepHover'
  import { filterScenarioHints, applyAutoFixableScenarioHints } from './lib/scenarioHints'
  import {
    DEFAULT_EDITOR_SETTINGS,
    editorSettingsFromDTO,
    type EditorSettings,
  } from './lib/editorOptions'
  import { resolveRecordStartURL } from './lib/recordStartUrl'
  import { formatCloseAppMessage } from './lib/closeAppGuard'
  import { resolveProjectPathInput as resolveProjectPathShortcut } from './lib/projectPath'
  import { isRecordingTargetReadOnly, isSameRecordTab, normalizeRecordTabPath, recordingTabSwitchAllowed, resolveRecordingTargetPath } from './lib/recordingTarget'
  import {
    captureRecordStartedUiTarget,
    recordStepSourceText,
    resolveRecordEditorPrepareAction,
    resolveRecordFeaturePathForUI,
    resolveApplyRecordStepTarget,
    shouldApplyRecordStepEvent,
    shouldBufferEarlyRecordStepEvent,
  } from './lib/recordingLifecycle'
  import { flakyScenarioMap, flakyStepHints } from './lib/flakyMetrics'
  import { loadRecents, rememberFeature, rememberProject } from './lib/recents'
  import { callWailsWithTimeout } from './lib/wailsTimeout'
  import { startRunResultJob } from './lib/asyncRunResult'
  import { icons, toolbarIcons } from './lib/icons'
  import { EventsOn, EventsOff, OnFileDrop, OnFileDropOff, WindowShow, WindowUnminimise } from '../wailsjs/runtime/runtime'
  import {
    Version,
    OpenProject,
    RefreshProject,
    ReadFeature,
    SaveFeature,
    WriteTempFeature,
    CancelRun,
    StartRun,
    StartValidate,
    AnalyzeEditorContent,
    ListTestClients,
    InitProject,
    InitProjectAt,
    PickProjectFolder,
    PickSaveFile,
    PickOpenFile,
    StartRunPlugin,
    StartRecord,
    OpenBrowser,
    BeginRecordingCapture,
    StartRecordBaseline,
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
    StartServeAllure,
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
    UpdateDirtyTabsState,
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
    LoadProjectConfig,
    SaveProjectConfig,
    ConfirmAppClose,
  } from '../wailsjs/go/wailsapp/App'
  import { gui } from '../wailsjs/go/models'

  const WELCOME_KEY = '__welcome__'

  type EditorTab = { path: string; content: string; dirty: boolean; draft?: string; unloaded?: boolean }
  type EditorStepRow = gui.EditorStepRow

  const projectStore = createProjectStore()
  const tabsStore = createTabsStore(WELCOME_KEY)
  const runnerStore = createRunnerStore()
  const diagnosticsStore = createDiagnosticsStore()
  const recorderStore = createRecorderStore()
  const reportsStore = createReportsStore()
  const dialogsStore = createDialogsStore()
  const layoutStore = createLayoutStore()
  const uiPrefsStore = createUiPrefsStore()
  const onboardingTourStore = createOnboardingTourStore()
  const recorderPrefsStore = createRecorderPrefsStore()
  const settingsStore = createSettingsStore()
  const vanessaRunStore = createVanessaRunStore()
  const pluginRunStore = createPluginRunStore()
  const editorStore = createEditorStore()
  const runDialogStore = createRunDialogStore()
  const validateDialogStore = createValidateDialogStore()
  const journalStore = createJournalStore()
  const featureDialogStore = createFeatureDialogStore()
  const testClientStore = createTestClientStore()
  const pluginsStore = createPluginsStore()
  const updateDialogStore = createUpdateDialogStore()
  const catalogStore = createCatalogStore()
  const contextMenuStore = createContextMenuStore()
  const postRecordStore = createPostRecordStore()
  const confirmDialogStore = createConfirmDialogStore()
  const menuStore = createMenuStore()
  const projectReplaceStore = createProjectReplaceStore()
  const pickerDialogStore = createPickerDialogStore()
  const httpAuthDialogStore = createHttpAuthDialogStore()
  const stepsHelpDialogStore = createStepsHelpDialogStore()
  const otpDialogStore = createOtpDialogStore()
  const splashStore = createSplashStore({ appReady: false, message: t('splash.starting'), progress: 0, fading: false })
  const viewportStore = createViewportStore()
  const recentsStore = createRecentsStore()
  const settingsDialogStore = createSettingsDialogStore()
  const appMetaStore = createAppMetaStore()
  const sessionStore = createSessionStore()
  const runFormStore = createRunFormStore(
    defaultRunForm({ headed: true, installPW: true, html: true, htmlLightMode: true }),
  )
  const recordFormStore = createRecordFormStore({
    recordFeatureName: t('dialogs.record.featureDefault'),
    recordScenarioName: t('dialogs.record.scenarioDefault'),
  })

  const dialogBinds = createDialogBindController({
    settingsStore,
    uiPrefsStore,
    recorderPrefsStore,
    recordFormStore,
    validateDialogStore,
    featureDialogStore,
    projectReplaceStore,
    testClientStore,
    vanessaRunStore,
    pluginRunStore,
    runFormStore,
    settingsDialogStore,
  })

  const workspaceSession = createWorkspaceSessionController({
    stores: {
      settingsStore,
      recorderPrefsStore,
      uiPrefsStore,
      recentsStore,
      dialogBinds,
      tabsStore,
      testClientStore,
    },
    welcomeKey: WELCOME_KEY,
    getProjectPath: () => projectPath,
    getTabs: () => tabs,
    getActiveTab: () => activeTab,
    getEditorText: () => monaco?.getEditorText() ?? editorText,
    syncActiveTabContent,
    isUntitled,
    saveSettings: SaveSettings,
    saveFeatureDraft: SaveFeatureDraft,
    openProject: OpenProject,
    resolveProjectPath: resolveProjectPathInput,
    applyProjectScan,
    listTestClients: ListTestClients,
    loadFeature: (path) => loadFeature(path),
    applyEditorText,
    trimTabsMemory,
    appendLog,
    setStatus,
    tr: (key, vars) => createTranslator(get(locale) as Locale)(key, vars),
  })

  $: tr = createTranslator($locale)

  $: ({
    showSettings,
    showCommandPalette,
    showSnippetPalette,
    showRecord,
    showPluginRun,
    showOtp,
    showRun,
    showTestClient,
    showExport,
    showImport,
    showImportFeatures,
    showDuplicateFeature,
    showRefactorUrl,
    showOpenProject,
    showRenameFeature,
    showMoveFeature,
    showValidate,
    showUpdateCheck,
    showInitProject,
    showNewProjectWizard,
    showSteps,
    showStepsHelp,
    showAbout,
    showPlugins,
    showProjectReplace,
    showHotkeys,
    showRunHistory,
    showVanessaRun,
    showVanessaSettings,
    showVanessaMonitor,
    showHttpAuth,
    showPickerStep,
    showPostRecordDiff,
  } = $dialogsStore)

  $: ({
    sidebarVisible,
    previewVisible,
    previewWidth,
    bottomPanelOpen,
    bottomPanelHeight,
    bottomTab,
    resizingBottom,
    resizingSteps,
    resizingSidebar,
    resizingPreview,
    previewPaneMounted,
  } = $layoutStore)

  $: ({
    toolbarCompact,
    stepsPanelVisible,
    stepsPanelHeight,
    sidebarWidth,
    checklistDismissed,
    welcomePlayedSuccess,
    onboardingCompleted,
    onboardingDismissed,
    runDialogConfirmed,
    pickerDuringRecording,
    stepsPanelCollapsed,
  } = $uiPrefsStore)

  $: ({
    active: showOnboardingTour,
    stepId: onboardingTourStepId,
    validateDone: onboardingValidateDone,
    dryRunDone: onboardingDryRunDone,
    journalVisited: onboardingJournalVisited,
  } = $onboardingTourStore)

  $: ({ filterRecording, navOnlyRecording, hoverRecord } = $recorderPrefsStore)

  $: ({
    recordURL,
    recordOutput,
    recordIdle,
    recordAppendTo,
    recordTestClient,
    recordFeatureName,
    recordScenarioName,
  } = $recordFormStore)

  $: ({
    path: projectPath,
    version: currentProjectVersion,
    features,
    tags,
    featureTags,
    scenarios: projectScenarios,
  } = $projectStore)

  $: ({ tabs, activeTab, welcomeTabVisible, pendingCloseTab } = $tabsStore)

  $: ({
    browser: settingsBrowser,
    headless: settingsHeadless,
    parallelWorkers: settingsWorkers,
    slowMo: settingsSlowMo,
    scrollBeforeClick: settingsScrollBeforeClick,
    disableRecordUrlWait: settingsDisableRecordUrlWait,
    hoverRecordMinMs: settingsHoverRecordMinMs,
    maxLoopIterations: settingsLoops,
    checkUpdatesOnStartup: settingsCheckUpdatesOnStartup,
    selectorClickStrategies: settingsSelectorClickStrategies,
    selectorInputStrategies: settingsSelectorInputStrategies,
    navWaitUntil: settingsNavWaitUntil,
    startUrl: startURL,
    uiLocale,
  } = $settingsStore)

  $: ({ lastRun, runForm } = $runFormStore)

  $: ({
    dry: vanessaDry,
    preferRerun: vanessaPreferRerun,
    tag: vanessaTag,
    excludeTags: vanessaExcludeTags,
    scenario: vanessaScenario,
    rerunDir: vanessaRerunDir,
    installEpf: vanessaInstallEpf,
    epfUrl: vanessaEpfUrl,
    epfDest: vanessaEpfDest,
    platformExe: vanessaPlatformExe,
    epfPath: vanessaEpfPath,
    ib: vanessaIB,
    reportAllure: vanessaReportAllure,
    vaDir: vanessaVaDir,
    vaFiles: vanessaVaFiles,
    dialogScenarios: vanessaDialogScenarios,
    running: vanessaRunning,
    snapshot: vanessaSnapshot,
    watchDir: vanessaWatchDir,
    plannedTotal: vanessaPlannedTotal,
  } = $vanessaRunStore)

  $: ({
    name: pluginRunName,
    dry: pluginRunDry,
    tag: pluginRunTag,
    scenario: pluginRunScenario,
    dialogScenarios: pluginRunScenarios,
  } = $pluginRunStore)

  $: ({
    text: editorText,
    textVersion: editorTextVersion,
    cursorLine: editorCursorLine,
    steps: editorSteps,
    stepsTextVersion: editorStepsTextVersion,
    stepsPanelTab,
  } = $editorStore)

  $: ({ title: runDialogTitle, scenarios: runDialogScenarios } = $runDialogStore)

  $: ({
    playing,
    runId: activeRunEventId,
    total: runProgressTotal,
    current: runProgressCurrent,
    label: playingLabel,
    logStreaming: runLogStreaming,
    cancelling: runCancelling,
    lastRunSince,
    lastRunBatchResults,
    lastErrorEntry,
    dryRunActive: runningDryRun,
  } = $runnerStore)

  $: ({
    browserOpen,
    recording,
    captureFinalizing,
    paused: recordPaused,
    targetPath: recordingTargetPath,
    recordSessionId: activeRecordSessionId,
    browserSessionId: activeBrowserSessionId,
    liveRecordStepLines,
    lastRecordTarget,
    pauseToggleGuardUntil,
  } = $recorderStore)

  $: appVersion = $appMetaStore.version
  $: stepStatusError = $diagnosticsStore.stepStatusError
  $: hintFixInFlight = $diagnosticsStore.hintFixInFlight
  $: ({ allureInstalled, allureServeRunning } = $reportsStore)
  $: projectArtifacts = $projectStore.artifacts
  $: ({ projects: recentProjects, features: recentFeatures } = $recentsStore)
  $: flakyByPath = flakyScenarioMap($reportsStore.flakyMetrics)
  $: flakyStepByPath = flakyStepHints($reportsStore.flakyMetrics)

  $: appReady = $splashStore.appReady
  $: splashMessage = $splashStore.message
  $: splashProgress = $splashStore.progress
  $: splashFading = $splashStore.fading
  $: viewportWidth = $viewportStore.width
  $: viewportHeight = $viewportStore.height
  $: viewportAutoCompact = $viewportStore.autoCompact
  $: toolbarIconOnly = $viewportStore.toolbarIconOnly

  let monaco: MonacoEditor | undefined
  let lastSyncedDirtyTabsState: boolean | null = null
  let workspaceSessionFlushed = false
  let pendingEarlyRecordSteps: Array<{
    event: RecordStepEvent
    targetPath: string
    recordSessionId: string
  }> = []
  const unsubscribers: (() => void)[] = []

  function askConfirm(opts: {
    title: string
    message: string
    confirmLabel?: string
    danger?: boolean
    dontAskAgainLabel?: string
  }): Promise<boolean> {
    return confirmDialogStore.ask(opts, { ok: tr('common.ok') })
  }

  const applyCatalogFilterDebounced = debounce((text: string) => {
    catalogStore.setFilterText(text)
  }, 200)

  function closeConfirm(confirmed: boolean, dontAskAgain = false) {
    confirmDialogStore.close(confirmed, dontAskAgain)
  }

  async function handleAppCloseRequested(reasons: string[]) {
    if (reasons.length > 0) {
      const ok = await askConfirm({
        title: tr('confirm.closeApp.title'),
        message: formatCloseAppMessage(tr, reasons),
        confirmLabel: tr('confirm.closeApp.confirmLabel'),
        danger: true,
      })
      if (!ok) return
    }
    await flushWorkspaceSession()
    workspaceSessionFlushed = true
    await ConfirmAppClose()
  }

  const logStaleEvent = createStaleEventLogger()
  const wailsEventGuards = createEventStalenessGuards(
    {
      getProjectVersion: () => currentProjectVersion,
      getActiveRunId: () => activeRunEventId,
      getActiveRecordSessionId: () => activeRecordSessionId,
      getActiveBrowserSessionId: () => activeBrowserSessionId,
      isRecording: () => recording,
      syncRecordSessionIds: (recordSessionId, browserSessionId) =>
        recorderStore.setSessionIDs(recordSessionId ?? '', browserSessionId ?? ''),
      syncRunId: (runId) => runnerStore.setRunId(runId),
    },
    logStaleEvent,
  )

  function setProjectState(next: Partial<ProjectState>) {
    projectStore.patch(next)
  }

  function resetProjectState() {
    projectStore.reset()
  }

  function onSidebarSearchInput(e: Event) {
    const value = (e.currentTarget as HTMLInputElement).value
    catalogStore.setSidebarSearch(value)
    applyCatalogFilterDebounced(value)
  }

  function setSidebarSearch(value: string) {
    catalogStore.setSidebarSearch(value)
    applyCatalogFilterDebounced.cancel()
    catalogStore.setFilterText(value)
  }

  $: editorValidationIssues = $diagnosticsStore.issues
  $: editorValidationByTab = $diagnosticsStore.issuesByTab
  $: validatePanelIssues = $diagnosticsStore.browserPanelIssues
  $: validateCliLog = $validateDialogStore.cliLog
  $: logText = $journalStore.logText
  $: statusMessage = $journalStore.statusMessage
  $: statusTone = $journalStore.statusTone
  $: duplicateFeaturePath = $featureDialogStore.duplicateFeaturePath
  $: moveFeaturePath = $featureDialogStore.moveFeaturePath
  $: moveDestDirs = $featureDialogStore.moveDestDirs
  $: renameFeaturePath = $featureDialogStore.renameFeaturePath
  $: exportInputPath = $featureDialogStore.exportInputPath
  $: importFeaturesBusy = $featureDialogStore.importFeaturesBusy
  $: testClients = $testClientStore.clients
  $: testClientSelection = $testClientStore.selection
  $: testClientSuggestName = $testClientStore.suggestName
  $: installedPlugins = $pluginsStore.installed
  $: updateCheckMessage = $updateDialogStore.message
  $: updateCheckHasUpdate = $updateDialogStore.hasUpdate
  $: updateCheckInfo = $updateDialogStore.info
  $: updateDownloading = $updateDialogStore.downloading
  $: updateProgress = $updateDialogStore.progress
  $: editorScenarioHints = diagnosticsHints
  $: batchSelected = $catalogStore.batchSelected
  $: batchMode = $catalogStore.batchMode
  $: sidebarSearch = $catalogStore.sidebarSearch
  $: catalogFilterText = $catalogStore.catalogFilterText
  $: catalogDropTarget = $catalogStore.dropTarget
  $: showBatchHint = $catalogStore.showBatchHint
  $: catalogCollapsed = catalogStore.collapsedSet()
  $: postRecordPath = $postRecordStore.path
  $: postRecordStepCount = $postRecordStore.stepCount
  $: postRecordBaselineText = $postRecordStore.baselineText
  $: contextMenu = $contextMenuStore.feature
  $: folderMenu = $contextMenuStore.folder
  $: stepsMenu = $contextMenuStore.steps
  $: confirmDialog = $confirmDialogStore.request
  $: confirmDialogOpen = $confirmDialogStore.open
  $: openMenu = $menuStore.openMenu
  $: findText = $projectReplaceStore.findText
  $: replaceText = $projectReplaceStore.replaceText
  $: replaceCaseSensitive = $projectReplaceStore.caseSensitive
  $: projectReplaceBusy = $projectReplaceStore.busy
  $: pickerSelector = $pickerDialogStore.selector
  $: pickerChoices = $pickerDialogStore.choices
  $: pickerCandidates = $pickerDialogStore.candidates
  $: pickerWarnings = $pickerDialogStore.warnings
  $: pickerSuggestedAction = $pickerDialogStore.suggestedAction
  $: pickerSuggestedChoice = $pickerDialogStore.suggestedChoice
  $: httpAuthHost = $httpAuthDialogStore.host
  $: stepsHelpQuery = $stepsHelpDialogStore.query
  $: otpEmail = $otpDialogStore.email
  $: diagnosticsHints = $diagnosticsStore.hints

  $: isWelcome = canShowWelcome(tabs)
  $: editorValuePath = isWelcome || !activeTab ? null : activeTab
  $: activeFeatureTab = tabs.find((t) => t.path === activeTab)
  $: activeTabUnsaved = activeFeatureTab ? tabIsUnsaved(activeFeatureTab) : false
  $: if (dialogBinds.bindEditorSettings.stepsPanelView) {
    editorStore.setStepsPanelTab(dialogBinds.bindEditorSettings.stepsPanelView)
  }
  $: stepCount = editorSteps.length
  $: editorLineCount = isWelcome ? 0 : editorText.split(/\r?\n/).length
  $: showLargeFileBanner = !isWelcome && isLargeFeatureFile(editorLineCount)
  $: unsavedTabCount = tabs.filter((t) => tabIsUnsaved(t)).length
  $: syncDirtyTabsState(unsavedTabCount > 0)
  $: [, lastRunSummary] = [$locale, formatLastRunSummary(lastRun)]
  $: automationActive = playing || vanessaRunning
  $: pickerToolbarEnabled =
    pickerDuringRecording
      ? browserOpen && !playing
      : (recording && recordPaused) || (browserOpen && !recording && !playing)
  $: batchCount = batchSelected.length
  $: batchSelectedSet = buildBatchSelectedSet(batchSelected)
  $: showRecordingBar = recording && !showRecord
  $: showPlayingBar = playing
  $: anyAppDialogOpen = computeAnyAppDialogOpen($dialogsStore, {
    confirmDialogOpen,
    pendingCloseTab,
    postRecordPath,
  })
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
    onboardingTourStore.resetProgress()
  }
  $: onboardingElevateMenubar =
    showOnboardingTour && (onboardingTourStepId === 'validate' || onboardingTourStepId === 'dry-run')
  $: onboardingElevateSidebar = showOnboardingTour && onboardingTourStepId === 'pick-feature'
  $: onboardingElevateWelcome =
    showOnboardingTour && (onboardingTourStepId === 'welcome' || onboardingTourStepId === 'open-examples')
  $: onboardingElevateBottom = showOnboardingTour && onboardingTourStepId === 'journal'
  $: if (showOnboardingTour && (onboardingTourStepId === 'validate' || onboardingTourStepId === 'dry-run')) {
    menuStore.open('run')
  }
  $: if (showOnboardingTour && openMenu === 'run' && (onboardingTourStepId === 'validate' || onboardingTourStepId === 'dry-run')) {
    void tick().then(() => onboardingTour?.relayout())
  }
  $: showBrowserOverlay = (browserOpen || recording || playing) && !anyAppDialogOpen
  $: stepStatusDisplay =
    stepCount === 0 && !stepStatusError
      ? tr('statusBar.steps', { count: 0 })
      : stepStatusError
        ? tr('statusBar.stepsWithErrors', { count: stepCount, errors: editorValidationIssues.length })
        : tr('statusBar.steps', { count: stepCount })

  let actionBarEl: HTMLElement | undefined

  $: actionBarCompact = toolbarCompact || viewportAutoCompact
  $: layoutSidebarWidth = effectiveSidebarWidth(sidebarWidth, viewportWidth, sidebarVisible)
  $: catalogIndent = catalogIndentStep(layoutSidebarWidth)
  $: compactCatalogTree = isCompactCatalogTree(layoutSidebarWidth)
  $: showPreviewPane = shouldShowPreviewPane(viewportWidth, previewVisible)
  $: layoutPreviewWidth = effectivePreviewWidth(previewWidth, viewportWidth, previewVisible)
  $: {
    if (showPreviewPane) {
      layoutStore.schedulePreviewMount(() => layoutStore.setPreviewPaneMounted(true))
    } else {
      layoutStore.clearPreviewMountTimer()
      layoutStore.setPreviewPaneMounted(false)
    }
  }

  $: activeFeaturePath = activeTab !== WELCOME_KEY ? activeTab : ''

  $: resultsPanelEntries =
    lastRunBatchResults.length > 0
      ? lastRunBatchResults
      : lastRunSince
        ? filterRunResultsSince($reportsStore.runResults, lastRunSince)
        : $reportsStore.runResults

  $: runByPath = buildRunByPathMap($reportsStore.runResults)

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
    catalogFilterText,
    runByPath,
    true,
    tagsByPath,
  )

  $: welcomeProjectOpen = !!projectPath
  $: welcomeRecorded = recording || browserOpen
  $: recordMode = $recordFormStore.recordMode
  $: baselineBusy = $recordFormStore.baselineBusy
  let onboardingTour: OnboardingTour

  $: stopActionLabel = playing
    ? tr('toolbar.stopTest')
    : recording
      ? tr('toolbar.stopRecord')
      : browserOpen
        ? tr('toolbar.closeBrowser')
        : tr('toolbar.stop')
  $: recordingTargetLabel =
    recording && recordingTargetPath ? basename(recordingTargetPath) : ''
  $: recordingTargetReadOnly = isRecordingTargetReadOnly(
    recording,
    recordingTargetPath,
    activeTab,
    captureFinalizing,
  )

  const MIN_SPLASH_MS = 1400
  const SPLASH_FADE_MS = 320
  let startupHasRestorableWorkspaceSession = false

  function sleep(ms: number) {
    return new Promise<void>((resolve) => window.setTimeout(resolve, ms))
  }

  function setSplashStage(message: string, progress: number) {
    splashStore.setStage(message, progress)
  }

  async function dismissSplash(startedAt: number) {
    const remaining = MIN_SPLASH_MS - (Date.now() - startedAt)
    if (remaining > 0) await sleep(remaining)

    splashStore.startFading()
    await sleep(SPLASH_FADE_MS)
    setSplashDocumentState(false)
    await openMainWindow()
    splashStore.markReady()
    if (shouldAutoStartOnboarding()) {
      layoutStore.showSidebar()
      onboardingTourStore.setStepId('welcome')
      onboardingTourStore.start()
    }
  }

  function shouldAutoStartOnboarding(): boolean {
    return !startupHasRestorableWorkspaceSession && !onboardingCompleted && !onboardingDismissed
  }

  async function completeOnboarding() {
    onboardingTourStore.stop()
    uiPrefsStore.patch({ onboardingCompleted: true, onboardingDismissed: false })
    await persistSettings()
    void maybeCheckUpdatesOnStartup()
  }

  async function dismissOnboarding() {
    onboardingTourStore.stop()
    uiPrefsStore.patch({ onboardingDismissed: true })
    await persistSettings()
    void maybeCheckUpdatesOnStartup()
  }

  function openJournalTab(markTourVisit = false) {
    layoutStore.openBottomTab('journal')
    if (markTourVisit && showOnboardingTour) onboardingTourStore.patch({ journalVisited: true })
  }

  function restartOnboardingTour() {
    onboardingTourStore.resetProgress()
    onboardingTourStore.setStepId('')
    layoutStore.showSidebar()
    void resetWorkspaceForOnboarding().then(() => {
      onboardingTourStore.setStepId('welcome')
      onboardingTourStore.start()
      onboardingTour?.restart()
    })
  }

  async function resetWorkspaceForOnboarding() {
    cancelPendingFeatureLoads()
    if (projectPath) {
      try {
        await teardownDesktopSession()
      } catch {
        /* offline */
      }
      resetProjectState()
    }
    for (const t of tabs) {
      monaco?.releaseTab(t.path)
    }
    monaco?.retainTabs([])
    tabsStore.reset()
    editorStore.reset()
    monaco?.activateTab(null, '')
    catalogStore.clearBatch()
    diagnosticsStore.setIssues([])
    diagnosticsStore.clearIssuesByTab()
    clearEditorValidation()
    menuStore.close()
    layoutStore.closeBottomPanel()
  }

  function onOnboardingStepChange(e: CustomEvent<{ index: number; id: string }>) {
    const { id } = e.detail
    onboardingTourStore.setStepId(id)
    if (id === 'open-examples' || id === 'pick-feature') {
      layoutStore.showSidebar()
    }
    if (id === 'validate' || id === 'dry-run') {
      menuStore.open('run')
      void tick().then(() => tick().then(() => onboardingTour?.relayout()))
    }
    if (id === 'journal') {
      layoutStore.openBottomPanel()
      layoutStore.setBottomTab('journal')
      layoutStore.patch({ bottomPanelOpen: true })
      void tick().then(() => onboardingTour?.relayout())
    }
  }

  function setupWailsEventBindings(): () => void {
    return bindWailsEvents({
      guards: wailsEventGuards,
      handlers: {
        onOtpPrompt(email) {
          otpDialogStore.setEmail(email)
          WindowUnminimise()
          WindowShow()
          dialogsStore.open('showOtp')
        },
        onAppCloseRequested(reasons) {
          void handleAppCloseRequested(reasons)
        },
        onBrowserOpened() {
          applyBrowserSessionState({ browserOpen: true, recording: false, paused: false })
          dialogsStore.close('showRecord')
          setStatus(tr('journal.status.browserOpen'), 'busy')
          appendLog(tr('journal.browser.opened'))
          startBrowserWatch()
        },
        onBrowserClosed(result) {
          dismissRecorderPicker()
          stopBrowserWatch()
          handleRecordSessionEnd((result as gui.RunResult) ?? gui.RunResult.createFrom({}), 'browse')
        },
        onBrowserLost() {
          dismissRecorderPicker()
          handleBrowserLost()
        },
        onToolbarPicker() {
          void pickElement()
        },
        async onRecordStarted(meta) {
          const m = typeof meta === 'string' ? { append: false, output: meta } : meta
          const appendOnly = m.append === true
          const syncOnly = m.sync === true
          const wasRecording = recording
          const uiTarget = captureRecordStartedUiTarget(m.targetPath || '', activeTab, WELCOME_KEY)
          if (uiTarget) {
            recorderStore.setTargetPath(uiTarget)
          }
          applyBrowserSessionState({ browserOpen: true, recording: true, paused: false })
          dialogsStore.close('showRecord')
          setStatus(tr('journal.status.preparingRecord'), 'busy')
          startBrowserWatch()
          recorderStore.setRecordEditorReadyPromise((async () => {
            if (!syncOnly) {
              if (!appendOnly && !wasRecording) {
                recorderStore.setLiveRecordStepLines({})
              }
              if (!appendOnly) {
                await prepareRecordEditorTab(m.output || '', uiTarget)
              }
              postRecordStore.setBaselineText(monaco?.getEditorText() ?? editorText)
            }
          })())
          try {
            await recorderStore.awaitRecordEditorReady()
          } catch (e: any) {
            appendLog(tr('journal.record.prepTabError', { error: String(e) }))
          }
          const tabsSnap = tabsStore.snapshot()
          if (tabsSnap.activeTab && tabsSnap.activeTab !== WELCOME_KEY) {
            const recordTab = normalizeRecordTabPath(tabsSnap.activeTab)
            recorderStore.setTargetPath(recordTab)
            recorderStore.setLastRecordTarget(recordTab)
          }
          await flushPendingEarlyRecordSteps()
          if (!appendOnly && !syncOnly) {
            appendLog(tr('journal.record.started'))
            setStatus(tr('journal.status.recording'), 'busy')
          } else if (syncOnly) {
            setStatus(tr('journal.status.recording'), 'busy')
          }
        },
        onRecordStopped(payload) {
          handleRecordStopped(payload)
        },
        onRunLogLine(line) {
          appendLog(line)
        },
        onRunResultsChanged() {
          scheduleRefreshRunResults()
        },
        onReportGoto(req) {
          void gotoReportStep(req)
        },
        onReportRerun(req) {
          void rerunFromReport(req)
        },
        onReportTrace(req) {
          void openTraceFromReport(req)
        },
        onRecordStep(payload) {
          const op = (payload.op || 'upsert') as RecordStepEvent['op']
          void applyRecordStepToTarget(
            {
              op,
              index: payload.index,
              line: payload.line,
              lines: payload.lines,
            },
            payload.targetPath ?? '',
            payload.recordSessionId ?? '',
          )
        },
        onRecordFinished(result) {
          stopBrowserWatch()
          handleRecordSessionEnd((result as gui.RunResult) ?? gui.RunResult.createFrom({}), 'record')
        },
        onRecordError(message) {
          stopBrowserWatch()
          pendingEarlyRecordSteps = []
          recorderStore.reset()
          dialogsStore.close('showRecord')
          recorderStore.setLiveRecordStepLines({})
          appendLog(tr('journal.record.error', { message: message || tr('journal.record.unknownError') }))
          setStatus(tr('journal.status.recordError'), 'error')
          syncIdleStatus()
        },
        onVanessaRunStarted() {
          vanessaRunStore.setRunning(true)
          vanessaWatchDir = ''
          setStatus(tr('journal.status.vanessaRunning'), 'busy')
          appendLog(tr('journal.vanessa.starting'))
          startVanessaPoll()
        },
        async onVanessaRunFinished(result) {
          stopVanessaPoll()
          vanessaRunStore.setRunning(false)
          const dto = result as gui.VanessaRunResultDTO
          if (dto?.runDir) {
            vanessaWatchDir = dto.runDir
            try {
              vanessaRunStore.setSnapshot(await PollVanessaRun(dto.runDir, vanessaPlannedTotal))
            } catch {
              /* ignore */
            }
          }
          if (dto?.output) appendLog(dto.output.trimEnd())
          if (dto?.error) {
            appendLog(tr('journal.error.generic', { error: dto.error }))
            setStatus(tr('journal.status.vanessaError'), 'error')
          } else {
            appendLog(tr('journal.vanessa.done'))
            setStatus(tr('journal.status.vanessaDone'), dto?.success ? 'success' : 'error')
          }
          await refreshRunResults()
        },
      },
      isRunLogStreaming: () => runLogStreaming,
      isPlaying: () => playing,
      isRecorderActive: () => browserOpen || recording,
      runProgress: {
        getTotal: () => runProgressTotal,
        getCurrent: () => runProgressCurrent,
        getLabel: () => playingLabel,
        updateProgress: (total, current, label) => runnerStore.progress(total, current, label),
      },
      scheduleRefreshRunResults,
      eventsOn: EventsOn,
    })
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


    setSplashStage(tr('splash.connecting'), 40)

    try {
      appMetaStore.setVersion(await Version())
    } catch {
      appMetaStore.setVersion('dev')
    }

    setSplashStage(tr('splash.loadingSettings'), 60)

    const [recents, settings] = await Promise.all([
      loadRecents(),
      callWailsWithTimeout('LoadSettings', LoadSettings(), 4000),
    ])
    recentsStore.setRecents(recents.projects, recents.features)
    if (settings) {
      startupHasRestorableWorkspaceSession = hasRestorableWorkspaceSession(settings)
      applySettingsFromDTO(settings)
      setStepHoverEnabled(() => dialogBinds.bindEditorSettings.stepHover)
      syncStepsPanelCollapsedFromPrefs()
      if (settings.sidebarWidth >= VIEWPORT.sidebarMin) {
        uiPrefsStore.patch({ sidebarWidth: clampSidebarWidth(settings.sidebarWidth) })
      }
      if (settings.recentProjects?.length) recentsStore.patch({ projects: settings.recentProjects })
      if (settings.recentFeatures?.length) recentsStore.patch({ features: settings.recentFeatures })
      if (startupHasRestorableWorkspaceSession || !shouldAutoStartOnboarding()) {
        await restoreWorkspaceSession(settings)
      }
    }

    sessionStore.startDraftAutosave(() => void autosaveDirtyDrafts())

    setSplashStage(tr('splash.initializing'), 88)

    try {
      OnFileDrop((_x, _y, paths) => {
        if (projectPath && paths?.length) {
          const dest = catalogDropTarget || projectPath
          catalogStore.setDropTarget('')
          void importDroppedFeatures(dest, paths)
        }
      }, false)
      unsubscribers.push(() => OnFileDropOff())
    } catch {
      /* dev without wails runtime */
    }

    try {
      unsubscribers.push(setupWailsEventBindings())
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
      menuStore.close()
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
    if (new URLSearchParams(location.search).has('e2e') || (window as unknown as { __SCENARIA_E2E_MOCK__?: boolean }).__SCENARIA_E2E_MOCK__) {
      ;(window as unknown as { __e2eCheckActiveTabDiskStale?: () => Promise<void> }).__e2eCheckActiveTabDiskStale = () =>
        checkActiveTabDiskStale()
      ;(window as unknown as {
        __e2eEmitEditorChange?: (path: string | null, text: string) => void
      }).__e2eEmitEditorChange = (path, text) => monaco?.emitEditorChangeForTest?.(path, text)
      ;(window as unknown as {
        __e2eLoadFeature?: (path: string, forceActivate?: boolean) => Promise<void>
      }).__e2eLoadFeature = (path, forceActivate = false) => loadFeature(path, { forceActivate })
      ;(window as unknown as { __e2eFlushSession?: () => Promise<void> }).__e2eFlushSession = async () => {
        await flushWorkspaceSession()
      }
      unsubscribers.push(() => {
        delete (window as unknown as { __e2eCheckActiveTabDiskStale?: () => Promise<void> }).__e2eCheckActiveTabDiskStale
        delete (window as unknown as {
          __e2eEmitEditorChange?: (path: string | null, text: string) => void
        }).__e2eEmitEditorChange
        delete (window as unknown as {
          __e2eLoadFeature?: (path: string, forceActivate?: boolean) => Promise<void>
        }).__e2eLoadFeature
        delete (window as unknown as { __e2eFlushSession?: () => Promise<void> }).__e2eFlushSession
      })
    }
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
    if (!workspaceSessionFlushed) {
      void flushWorkspaceSession()
    }
    void teardownDesktopSession()
    applyCatalogFilterDebounced.cancel()
    if (!workspaceSessionFlushed) {
      diagnosticsStore.clearValidateDebounce()
      sessionStore.teardown()
    } else {
      sessionStore.stopDraftAutosave()
    }
    layoutStore.clearPreviewMountTimer()
    for (const off of unsubscribers) off()
  })

  function schedulePersistSession() {
    sessionStore.schedulePersist(() => void persistSettings())
  }

  async function flushWorkspaceSession() {
    syncActiveTabContent()
    sessionStore.flushPersist()
    await persistSettings()
    workspaceSessionFlushed = true
  }

  async function restoreWorkspaceSession(s: gui.AppSettingsDTO) {
    await workspaceSession.restoreWorkspaceSession(s)
  }

  async function autosaveDirtyDrafts() {
    await workspaceSession.autosaveDirtyDrafts()
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
    contextMenuStore.openSteps({ x: e.clientX, y: e.clientY, line: step.line, step })
  }

  function closeStepsMenu() {
    contextMenuStore.closeSteps()
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

  function syncDirtyTabsState(dirty: boolean) {
    if (lastSyncedDirtyTabsState === dirty) return
    lastSyncedDirtyTabsState = dirty
    try {
      void Promise.resolve(UpdateDirtyTabsState(dirty)).catch((e) => {
        console.warn('failed to sync dirty tab state', e)
      })
    } catch (e) {
      console.warn('failed to sync dirty tab state', e)
    }
  }

  function applySettingsFromDTO(s: gui.AppSettingsDTO) {
    settingsStore.applyFromDTO(s)
    recorderPrefsStore.applyFromDTO(s)
    uiPrefsStore.applyFromDTO(s)
    dialogBinds.syncRecorderPrefsBindLocals()
    dialogBinds.syncUiPrefsBindLocals()
    dialogBinds.syncEditorSettingsBindLocal()
    dialogBinds.syncSettingsBindLocals()
    editorStore.setStepsPanelTab(dialogBinds.bindEditorSettings.stepsPanelView)
    syncStepsPanelCollapsedFromPrefs()
    runFormStore.applyLastRunFromSettings(
      s.browser || 'chromium',
      s.parallelWorkers || 1,
      s.slowMo ?? 0,
    )

    const locale = settingsStore.snapshot().uiLocale
    if (locale === 'en' || locale === 'ru') setLocale(locale)
  }

  function resolveStepsPanelCollapsed(): boolean {
    return !stepsPanelVisible
  }

  function syncStepsPanelCollapsedFromPrefs() {
    uiPrefsStore.setStepsPanelCollapsed(resolveStepsPanelCollapsed())
  }

  function toggleToolbarCompact() {
    uiPrefsStore.patch({ toolbarCompact: !toolbarCompact })
    void persistSettings()
  }

  function togglePreview() {
    layoutStore.patch({ previewVisible: !previewVisible })
  }

  function toggleStepsPanel() {
    uiPrefsStore.patch({ stepsPanelVisible: !stepsPanelVisible })
    syncStepsPanelCollapsedFromPrefs()
    void persistSettings()
  }

  function showAboutDialog() {
    dialogsStore.open('showAbout')
  }

  function openFindReplace() {
    if (isWelcome) {
      appendLog(tr('journal.search.openScenario'))
      return
    }
    monaco?.openFindReplace()
  }

  function resetWindowLayout() {
    layoutStore.reset()
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
    isEnabled: () => dialogBinds.bindEditorSettings.codeLens && !!activeTab && !isWelcome,
    onRun: (payload: { scenario: string; line: number; dryRun: boolean; partial: boolean }) =>
      runScenarioAtLine(payload.line, payload.dryRun, payload.scenario, payload.partial),
  }

  const monacoInlayHintsHandlers = {
    isEnabled: () => dialogBinds.bindEditorSettings.inlayHints && !!activeTab && !isWelcome,
    getSteps: () => editorSteps,
    isSnapshotCurrent: () => isEditorAnalysisSnapshotVisible(editorStepsTextVersion, editorTextVersion),
  }

  async function refreshEditorScenarioHints() {
    if (isWelcome || !dialogBinds.bindEditorSettings.scenarioHints) {
      diagnosticsStore.setHints([])
      return
    }
    try {
      const all = await AnalyzeScenarioHints(editorText)
      diagnosticsStore.setHints(
        all
          .filter((h) => !diagnosticsStore.isHintDismissed(hintDismissKey(h)))
          .filter((h) => filterScenarioHints([h], dialogBinds.bindEditorSettings).length > 0),
      )
    } catch {
      diagnosticsStore.setHints([])
    }
  }

  async function runScenarioHintsAutoFix(text: string): Promise<{ text: string; count: number }> {
    if (!dialogBinds.bindEditorSettings.scenarioHints || !dialogBinds.bindEditorSettings.scenarioHintsAutoFixOnSave) {
      return { text, count: 0 }
    }
    const all = await AnalyzeScenarioHints(text)
    const fixable = filterScenarioHints(all, dialogBinds.bindEditorSettings).filter((h) => h.autoFixable)
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
    if (hintFixInFlight || isWelcome) return
    diagnosticsStore.setHintFixInFlight(true)
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
      diagnosticsStore.setHintFixInFlight(false)
    }
  }

  function dismissEditorHint(hint: gui.ScenarioHintDTO) {
    diagnosticsStore.dismissHint(hintDismissKey(hint))
    void refreshEditorScenarioHints()
  }

  async function showPostRecordBanner(path: string) {
    postRecordStore.open(path)
    diagnosticsStore.clearDismissedHints()
    try {
      const tab = tabs.find((t) => t.path === path)
      if (activeTab === path) {
        postRecordStore.setStepCount((await ParseEditorSteps(monaco?.getEditorText() ?? editorText)).length)
      } else if (tab) {
        postRecordStore.setStepCount((await ParseEditorSteps(tabEditorText(tab))).length)
      } else if (isRealFeaturePath(path)) {
        const content = await ReadFeature(path)
        postRecordStore.setStepCount((await ParseEditorSteps(content)).length)
      } else {
        postRecordStore.setStepCount(0)
      }
    } catch {
      postRecordStore.setStepCount(0)
    }
    if (dialogBinds.bindEditorSettings.scenarioHints && dialogBinds.bindEditorSettings.scenarioHintsAfterRecord) {
      await refreshEditorScenarioHints()
    } else {
      diagnosticsStore.setHints([])
    }
  }

  function dismissPostRecord() {
    postRecordStore.dismiss()
    dialogsStore.close('showPostRecordDiff')
  }

  function openPostRecordDiff() {
    if (!postRecordPath) return
    const tab = tabs.find((t) => t.path === postRecordPath)
    const modified =
      postRecordPath === activeTab
        ? (monaco?.getEditorText() ?? editorText)
        : tab
          ? tabEditorText(tab)
          : ''
    if (postRecordBaselineText === modified) {
      appendLog(tr('journal.record.noPostRecordDiff'))
      return
    }
    dialogsStore.open('showPostRecordDiff')
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
    layoutStore.setBottomTab('validate')
  }

  async function postRecordSave() {
    if (!postRecordPath) return
    if (activeTab !== postRecordPath) await loadFeature(postRecordPath)
    await saveFeature()
    dismissPostRecord()
  }

  function openDuplicateDialog(path: string) {
    if (!path || isWelcome) return
    featureDialogStore.openDuplicate(path, `${basename(path).replace(/\.feature$/i, '')}-copy`)
    dialogBinds.bindDuplicateNewName = featureDialogStore.snapshot().duplicateNewName
    dialogsStore.open('showDuplicateFeature')
  }

  async function confirmDuplicateFeature(newName: string) {
    dialogsStore.close('showDuplicateFeature')
    const path = duplicateFeaturePath
    featureDialogStore.clearDuplicate()
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
      try {
        await ClearFeatureDraft(path)
      } catch (e) {
        console.warn('failed to clear deleted feature draft', e)
      }
      if (tabs.some((t) => t.path === path)) {
        finalizeCloseTab(path)
      } else {
        evictFeatureSymbolCache(path)
        monaco?.releaseTab(path)
      }
      await refreshProject()
      appendLog(tr('journal.file.deleted', { name: basename(path) }))
    } catch (e: any) {
      appendLog(tr('journal.error.generic', { error: String(e) }))
    }
  }

  async function runFeatureFile(path: string, dryRun = false) {
    if (!path) return
    await loadFeature(path)
    await executeRun(runFormFromMode(lastRun, 'single', { dryRun }), [path])
  }

  function onFileContextMenu(e: MouseEvent, path: string) {
    e.preventDefault()
    contextMenuStore.openFeature({ x: e.clientX, y: e.clientY, path })
  }

  function onFolderContextMenu(e: MouseEvent, node: CatalogNode) {
    e.preventDefault()
    const paths = collectFeaturePathsUnder(node)
    contextMenuStore.openFolder({ x: e.clientX, y: e.clientY, dir: node.path, paths })
  }

  function onExplorerContextMenu(e: MouseEvent) {
    if ((e.target as Element | null)?.closest('.catalog-tree-row')) return
    e.preventDefault()
  }

  function dismissFolderMenu() {
    contextMenuStore.closeFolder()
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
    void executeRun(runFormFromMode(lastRun, 'batch', { dryRun }), paths)
  }

  function folderMenuSelectBatch() {
    if (!folderMenu) return
    catalogStore.patch({ batchMode: true, batchSelected: [...folderMenu.paths] })
    dismissFolderMenu()
    appendLog(tr('journal.catalog.batchSelected', { count: batchSelected.length }))
  }

  function folderMenuRefresh() {
    dismissFolderMenu()
    void refreshCatalog()
  }

  function openVanessaForFolder(dirPath: string, dry: boolean) {
    vanessaRunStore.prepareDialog({
      dry,
      dialogScenarios: dialogScenarioNames(),
      vaDir: folderRelativePath(dirPath),
    })
    dialogBinds.syncVanessaBindLocals()
    dialogsStore.open('showVanessaRun')
  }

  function folderMenuVanessa(dry: boolean) {
    if (!folderMenu) return
    const dir = folderMenu.dir
    dismissFolderMenu()
    openVanessaForFolder(dir, dry)
  }

  function dismissContextMenu() {
    contextMenuStore.closeFeature()
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
    const srcPath = contextMenu.path
    const dirs = collectProjectDirs().filter((d) => d !== dirname(srcPath.replace(/\\/g, '/')))
    featureDialogStore.openMove(srcPath, dirs)
    dialogBinds.bindMoveDestDir = featureDialogStore.snapshot().moveDestDir || projectPath.replace(/\\/g, '/')
    dismissContextMenu()
    dialogsStore.open('showMoveFeature')
  }

  async function confirmMoveFeature(destDir: string) {
    dialogsStore.close('showMoveFeature')
    const src = moveFeaturePath
    featureDialogStore.clearMove()
    if (!src || !destDir) return
    if (!isRealFeaturePath(src)) {
      appendLog(tr('journal.project.saveOrOpenFirst'))
      return
    }
    try {
      const newPath = await MoveFeature(src, destDir)
      const wasActive = activeTab === src
      if (tabs.some((t) => t.path === src)) {
        tabsStore.mapTabs((tabs) => tabs.map((t) => (t.path === src ? { ...t, path: newPath } : t)))
        if (wasActive) tabsStore.setActiveTab(newPath)
      }
      catalogStore.mapBatchSelected((paths) => paths.map((p) => (p === src ? newPath : p)))
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
    featureDialogStore.openRename(contextMenu.path)
    dismissContextMenu()
    dialogsStore.open('showRenameFeature')
  }

  async function renameFeature(path: string, newName: string) {
    if (!path) return
    if (!isRealFeaturePath(path)) {
      appendLog(tr('journal.project.saveOrOpenFirst'))
      return
    }
    try {
      const newPath = await RenameFeature(path, newName)
      const wasActive = activeTab === path
      if (tabs.some((t) => t.path === path)) {
        tabsStore.mapTabs((tabs) => tabs.map((t) => (t.path === path ? { ...t, path: newPath } : t)))
        if (wasActive) tabsStore.setActiveTab(newPath)
      }
      catalogStore.mapBatchSelected((paths) => paths.map((p) => (p === path ? newPath : p)))
      await refreshProject()
      if (wasActive) await loadFeature(newPath)
      appendLog(tr('journal.file.renamed', { name: basename(newPath) }))
    } catch (e: any) {
      appendLog(tr('journal.error.generic', { error: String(e) }))
    }
  }

  async function confirmProjectReplace() {
    if (!dialogBinds.bindProjectReplaceFind) return
    projectReplaceStore.patch({
      findText: dialogBinds.bindProjectReplaceFind,
      replaceText: dialogBinds.bindProjectReplaceReplace,
      caseSensitive: dialogBinds.bindProjectReplaceCaseSensitive,
    })
    projectReplaceStore.setBusy(true)
    try {
      const result = await ReplaceInProject({
        find: dialogBinds.bindProjectReplaceFind,
        replace: dialogBinds.bindProjectReplaceReplace,
        caseSensitive: dialogBinds.bindProjectReplaceCaseSensitive,
      })
      appendLog(tr('journal.file.replaceDone', { replacements: result.replacements, files: result.filesChanged }))
      dialogsStore.close('showProjectReplace')
      await refreshProject()
      if (activeTab && !isWelcome && isRealFeaturePath(activeTab)) {
        editorStore.setText(await ReadFeature(activeTab))
        tabsStore.mapTabs((tabs) => tabs.map((t) => (t.path === activeTab ? { ...t, content: editorText } : t)))
        validateEditor()
      }
    } catch (e: any) {
      appendLog(tr('journal.error.generic', { error: String(e) }))
    } finally {
      projectReplaceStore.setBusy(false)
    }
  }

  async function refreshRunResults() {
    if (!reportsStore.tryBeginRefresh()) return
    const started = perfNow()
    try {
      reportsStore.setData(await ListRunResults(50), await FlakyMetrics(200))
    } catch {
      reportsStore.setData([], null)
    } finally {
      perfMark('refreshRunResults', started)
      if (reportsStore.finishRefresh()) {
        void refreshRunResults()
      }
    }
  }

  function scheduleRefreshRunResults() {
    if ($reportsStore.refreshInFlight) {
      reportsStore.queueRefresh()
      return
    }
    void refreshRunResults()
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
    const batch = remapRunResultPaths(
      filterRunResultsSince($reportsStore.runResults, runSince),
      diskTargets,
      runTargets,
    )
    let lastError = pickLastRunError(batch)
    if (!lastError && cliError && !/context canceled/i.test(cliError)) {
      const featurePath =
        runTargets.length === 1
          ? runResultFeaturePath(runTargets[0])
          : diskTargets.length === 1
            ? diskTargets[0]
            : tr('journal.run.defaultLabel')
      lastError = buildSyntheticRunError({
        featurePath,
        scenario: scenario || undefined,
        message: cliError,
        runner,
      })
    }
    const batchResults = batch.length === 0 && lastError ? [lastError] : batch
    runnerStore.setLastRunSession(runSince, batchResults, lastError)
  }

  async function openRunHistory() {
    await refreshRunResults()
    dialogsStore.open('showRunHistory')
  }

  async function openFeatureFromHistory(path: string) {
    dialogsStore.close('showRunHistory')
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
    const resolvedPath = await resolveProjectPathInput(path)
    if (!resolvedPath) return
    const switching = !!projectPath && projectPath !== resolvedPath
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
      const info = await OpenProject(resolvedPath)
      const openedPath = info.path || resolvedPath
      applyProjectScan(info, openedPath)
      await tick()
      testClientStore.setClients(await ListTestClients().catch((): string[] => []))
      await rememberProject(openedPath)
      const recents = await loadRecents()
      recentsStore.patch({ projects: recents.projects })
      appendLog(tr('journal.project.opened', { path: openedPath }))
      syncIdleStatus()
      await refreshRunResults()
      await refreshArtifacts()
      schedulePersistSession()
      await flushPendingStartupUpdateCheck()
    } catch (e: any) {
      appendLog(tr('journal.error.generic', { error: String(e) }))
      setStatus(String(e), 'error')
    }
  }

  async function teardownDesktopSession() {
    pendingEarlyRecordSteps = []
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
    recorderStore.reset()
    recorderStore.clearLiveRecordSession()
  }

  async function resetWorkspaceForProjectSwitch() {
    await teardownDesktopSession()
    clearFeatureSymbolCache()
    for (const t of tabs) {
      monaco?.releaseTab(t.path)
    }
    monaco?.retainTabs([])
    tabsStore.reset()
    editorStore.reset()
    monaco?.activateTab(null, '')
    catalogStore.clearBatch()
    diagnosticsStore.setIssues([])
    diagnosticsStore.clearIssuesByTab()
    postRecordStore.dismiss()
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
    clearFeatureSymbolCache()
    resetProjectState()
    for (const t of tabs) {
      monaco?.releaseTab(t.path)
    }
    monaco?.retainTabs([])
    tabsStore.reset()
    editorStore.reset()
    monaco?.activateTab(null, '')
    catalogStore.clearBatch()
    diagnosticsStore.setIssues([])
    diagnosticsStore.clearIssuesByTab()
    appendLog(tr('journal.project.closed'))
    syncIdleStatus()
    schedulePersistSession()
  }

  function startVanessaPoll() {
    stopVanessaPoll()
    vanessaRunStore.setPollTimer(setInterval(async () => {
      try {
        let dir = vanessaWatchDir
        if (!dir) {
          const dirs = await ListVanessaRunDirs(1)
          dir = dirs[0] || ''
          if (dir) vanessaRunStore.patch({ watchDir: dir })
        }
        if (dir) {
          vanessaRunStore.setSnapshot(await PollVanessaRun(dir, vanessaPlannedTotal))
        }
      } catch {
        /* offline poll */
      }
    }, 2000))
  }

  function stopVanessaPoll() {
    vanessaRunStore.clearPollTimer()
  }

  async function openVanessaMonitor() {
    if (!projectPath) return
    dialogsStore.open('showVanessaMonitor')
    vanessaRunStore.patch({ plannedTotal: Math.max(1, features.length) })
    try {
      const dirs = await ListVanessaRunDirs(1)
      if (dirs[0]) {
        vanessaRunStore.patch({ watchDir: dirs[0] })
        vanessaRunStore.setSnapshot(await PollVanessaRun(dirs[0], vanessaPlannedTotal))
      }
    } catch {
      /* ignore */
    }
  }

  function buildVanessaPluginRequest(): gui.PluginRunRequest {
    dialogBinds.flushVanessaBindLocals()
    return vanessaRunStore.buildPluginRequest(vanessaRunStore.snapshot())
  }

  async function openExamples() {
    const examples = await BundledExamplesPath()
    if (!examples) {
      appendLog(tr('journal.examples.folderNotFound'))
      return
    }
    await openProjectAt(examples)
    layoutStore.showSidebar()
    selectTab(WELCOME_KEY)
    appendLog(tr('journal.examples.selectScenario'))
    setStatus(tr('journal.status.examplesOpened'), 'normal')
  }

  async function rerunFailed() {
    const source = lastRunBatchResults.length > 0 ? lastRunBatchResults : $reportsStore.runResults
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
    await executeRun(runFormFromMode(lastRun, 'single', { dryRun: false, scenario }), [filePath])
  }

  function toggleBatchFeature(path: string) {
    catalogStore.toggleBatchPath(path)
  }

  function toggleBatchMode() {
    if (batchMode) {
      catalogStore.clearBatch()
      return
    }
    catalogStore.patch({ batchMode: true, batchSelected: selectAllFeaturesUnder(catalogViewState.tree) })
  }

  function onCatalogToggleBatch(path: string) {
    if (!batchMode) catalogStore.patch({ batchMode: true })
    toggleBatchFeature(path)
  }

  function onCatalogCollapse(key: string, collapsed: boolean) {
    catalogStore.setCollapsed(key, collapsed)
  }

  function onCatalogActivate(path: string, kind: 'root' | 'dir' | 'file') {
    if (kind === 'file') loadFeature(path)
  }

  function clearBatchSelection() {
    catalogStore.setBatchSelected([])
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
    let opts: RunForm = batchRunFormFrom(lastRun, dryRun)
    appendLog(`[debug] batch selected count: ${batchSelected.length}`)
    if (
      batchSelected.length === 1 &&
      batchSelected[0] === activeTab &&
      !isWelcome
    ) {
      const scenario = scenarioAtLine(editorText, monaco?.getCursorLine() ?? 1)
      if (scenario) opts = runFormFromMode(opts, 'single', { scenario })
    }
    await executeRun(opts, [...batchSelected])
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
      openRunDialog(dryRun ? tr('journal.run.mode.dryRun') : '', { dryRun }, 'single')
      return
    }
    void executeRun(runFormFromMode(lastRun, 'single', { dryRun }))
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
    const runOpts = currentScenarioRunFormFrom(lastRun, { dryRun, scenario, startStep, endStep }, browserOpen)
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
    const runOpts = currentScenarioRunFormFrom(lastRun, { dryRun, scenario, startStep, endStep }, browserOpen)
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
        dialogsStore.open('showHotkeys')
        break
      case 'settings':
        openSettings()
        break
      case 'palette':
        dialogsStore.open('showCommandPalette')
        break
      case 'snippets':
        openSnippetPalette()
        break
      case 'journal':
        layoutStore.toggleBottomPanel()
        if (bottomPanelOpen) layoutStore.setBottomTab('journal')
        break
      case 'format':
        if (!isWelcome && activeTab) void monaco?.formatDocument()
        break
      case 'goto-symbol':
        if (!isWelcome && activeTab) monaco?.openSymbolOutline()
        break
      case 'escape':
        menuStore.close()
        if (showCommandPalette) dialogsStore.close('showCommandPalette')
        if (showSnippetPalette) dialogsStore.close('showSnippetPalette')
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

    if (confirmDialogOpen) return dismiss(() => closeConfirm(false))
    if (showHttpAuth) return dismiss(closeHttpAuthDialog)
    if (showPickerStep) return dismiss(() => { dialogsStore.close('showPickerStep') })

    const inModal = e.target instanceof Element && e.target.closest('.modal-backdrop, .palette-backdrop')
    if (!inModal && monacoOverlayConsumesEscape()) return

    if (pendingCloseTab) return dismiss(cancelCloseTab)
    if (showOtp) return dismiss(cancelOtp)
    if (showPluginRun) return dismiss(() => { dialogsStore.close('showPluginRun') })
    if (showCommandPalette) return dismiss(() => { dialogsStore.close('showCommandPalette') })
    if (showSnippetPalette) return dismiss(() => { dialogsStore.close('showSnippetPalette') })
    if (dialogBinds.bindRecordStepPickerOpen) return
    if (showRecord) return dismiss(() => { dialogsStore.close('showRecord') })
    if (showSteps) return dismiss(() => { dialogsStore.close('showSteps') })
    if (showProjectReplace) return dismiss(() => { dialogsStore.close('showProjectReplace') })
    if (showRunHistory) return dismiss(() => { dialogsStore.close('showRunHistory') })
    if (showPostRecordDiff && postRecordPath) return dismiss(() => { dialogsStore.close('showPostRecordDiff') })
    if (showHotkeys) return dismiss(() => { dialogsStore.close('showHotkeys') })
    if (showPlugins) return dismiss(() => { dialogsStore.close('showPlugins') })
    if (showAbout) return dismiss(() => { dialogsStore.close('showAbout') })
    if (showUpdateCheck && !updateDownloading) return dismiss(() => { dialogsStore.close('showUpdateCheck') })
    if (showImportFeatures) return dismiss(() => { dialogsStore.close('showImportFeatures') })
    if (showImport) return dismiss(() => { dialogsStore.close('showImport') })
    if (showDuplicateFeature) {
      return dismiss(() => {
        dialogsStore.close('showDuplicateFeature')
        featureDialogStore.clearDuplicate()
      })
    }
    if (showInitProject) return dismiss(() => { dialogsStore.close('showInitProject') })
    if (showNewProjectWizard) return dismiss(() => { dialogsStore.close('showNewProjectWizard') })
    if (showValidate) return dismiss(() => { dialogsStore.close('showValidate') })
    if (showMoveFeature) {
      return dismiss(() => {
        dialogsStore.close('showMoveFeature')
        featureDialogStore.clearMove()
      })
    }
    if (showRenameFeature) {
      return dismiss(() => {
        dialogsStore.close('showRenameFeature')
        featureDialogStore.clearRename()
      })
    }
    if (showOpenProject) return dismiss(() => { dialogsStore.close('showOpenProject') })
    if (showRefactorUrl) return dismiss(() => { dialogsStore.close('showRefactorUrl') })
    if (showExport) return dismiss(() => { dialogsStore.close('showExport') })
    if (showVanessaSettings) return dismiss(() => { dialogsStore.close('showVanessaSettings') })
    if (showStepsHelp) {
      return dismiss(() => {
        dialogsStore.close('showStepsHelp')
        stepsHelpDialogStore.clear()
      })
    }
    if (showTestClient) {
      return dismiss(() => {
        dialogsStore.close('showTestClient')
        testClientStore.clearSuggestName()
      })
    }
    if (showVanessaRun) return dismiss(() => { dialogsStore.close('showVanessaRun') })
    if (showRun) return dismiss(() => { dialogsStore.close('showRun') })
    if (showSettings) return dismiss(cancelSettings)
  }

  function onGlobalKeydown(e: KeyboardEvent) {
    if (e.code === 'F5' || e.key === 'F5') {
      e.preventDefault()
      e.stopPropagation()
      void refreshCatalog()
      return
    }
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
    layoutStore.setResizing('resizingPreview', true)
    e.preventDefault()
    const startX = e.clientX
    const startW = previewWidth
    const onMove = (ev: MouseEvent) => {
      layoutStore.patchLocal({ previewWidth: Math.max(200, Math.min(720, startW - (ev.clientX - startX))) })
    }
    const onUp = () => {
      layoutStore.setResizing('resizingPreview', false)
      window.removeEventListener('mousemove', onMove)
      window.removeEventListener('mouseup', onUp)
      layoutStore.patch({ previewWidth })
    }
    window.addEventListener('mousemove', onMove)
    window.addEventListener('mouseup', onUp)
  }

  function startResizeSidebar(e: MouseEvent) {
    layoutStore.setResizing('resizingSidebar', true)
    e.preventDefault()
    const startX = e.clientX
    const startW = sidebarWidth
    const onMove = (ev: MouseEvent) => {
      uiPrefsStore.patchLocal({ sidebarWidth: clampSidebarWidth(startW + (ev.clientX - startX)) })
    }
    const onUp = async () => {
      layoutStore.setResizing('resizingSidebar', false)
      window.removeEventListener('mousemove', onMove)
      window.removeEventListener('mouseup', onUp)
      await persistSettings()
    }
    window.addEventListener('mousemove', onMove)
    window.addEventListener('mouseup', onUp)
  }

  function startResizeBottom(e: MouseEvent) {
    layoutStore.setResizing('resizingBottom', true)
    e.preventDefault()
    const startY = e.clientY
    const startH = bottomPanelHeight
    const onMove = (ev: MouseEvent) => {
      layoutStore.patchLocal({ bottomPanelHeight: clampBottomPanelHeight(Math.max(80, Math.min(window.innerHeight * 0.6, startH + (startY - ev.clientY))), window.innerHeight) })
    }
    const onUp = () => {
      layoutStore.setResizing('resizingBottom', false)
      layoutStore.patch({ bottomPanelHeight })
      window.removeEventListener('mousemove', onMove)
      window.removeEventListener('mouseup', onUp)
    }
    window.addEventListener('mousemove', onMove)
    window.addEventListener('mouseup', onUp)
  }

  function startResizeSteps(e: MouseEvent) {
    layoutStore.setResizing('resizingSteps', true)
    e.preventDefault()
    const startY = e.clientY
    const startH = stepsPanelHeight
    const onMove = (ev: MouseEvent) => {
      stepsPanelHeight = Math.max(80, Math.min(480, startH + (startY - ev.clientY)))
      stepsPanelHeight = clampStepsPanelHeight(stepsPanelHeight, window.innerHeight)
    }
    const onUp = async () => {
      layoutStore.setResizing('resizingSteps', false)
      window.removeEventListener('mousemove', onMove)
      window.removeEventListener('mouseup', onUp)
      await persistSettings()
    }
    window.addEventListener('mousemove', onMove)
    window.addEventListener('mouseup', onUp)
  }

  function resizeStep(e: KeyboardEvent): number {
    return e.shiftKey ? 40 : 12
  }

  function onSidebarSplitterKeydown(e: KeyboardEvent) {
    if (e.key !== 'ArrowLeft' && e.key !== 'ArrowRight') return
    e.preventDefault()
    sidebarWidth = clampSidebarWidth(sidebarWidth + (e.key === 'ArrowRight' ? resizeStep(e) : -resizeStep(e)))
    void persistSettings()
  }

  function onPreviewSplitterKeydown(e: KeyboardEvent) {
    if (e.key !== 'ArrowLeft' && e.key !== 'ArrowRight') return
    e.preventDefault()
    layoutStore.patch({ previewWidth: Math.max(200, Math.min(720, previewWidth + (e.key === 'ArrowLeft' ? resizeStep(e) : -resizeStep(e)))) })
    layoutStore.patch({ previewWidth })
  }

  function onBottomSplitterKeydown(e: KeyboardEvent) {
    if (e.key !== 'ArrowUp' && e.key !== 'ArrowDown') return
    e.preventDefault()
    layoutStore.patch({ bottomPanelHeight: clampBottomPanelHeight(bottomPanelHeight + (e.key === 'ArrowUp' ? resizeStep(e) : -resizeStep(e)), window.innerHeight) })
  }

  function onStepsSplitterKeydown(e: KeyboardEvent) {
    if (e.key !== 'ArrowUp' && e.key !== 'ArrowDown') return
    e.preventDefault()
    stepsPanelHeight = clampStepsPanelHeight(stepsPanelHeight + (e.key === 'ArrowUp' ? resizeStep(e) : -resizeStep(e)), window.innerHeight)
    void persistSettings()
  }

  function appendLog(line: string) {
    journalStore.appendLog(line)
  }

  async function resolveProjectPathInput(path: string): Promise<string> {
    return resolveProjectPathShortcut(path, BundledExamplesPath)
  }

  function setStatus(msg: string, tone: typeof statusTone = 'normal') {
    journalStore.setStatus(msg, tone)
  }

  function applyDevUiMock() {
    if (!import.meta.env.DEV) return
    const mode = new URLSearchParams(location.search).get('mock')
    if (mode !== 'python') return
    setProjectState({
      path: 'C:/Users/bafgion/Documents/Projects/camel-1c-integration/target',
      features: [],
      tags: [],
      featureTags: {},
    })
    recentsStore.patch({
      features: [
      'C:/Users/bafgion/Documents/Projects/camel-1c-integration/target/smoke.feature',
      'C:/Users/bafgion/Documents/Projects/camel-1c-integration/target/login.feature',
      'C:/Users/bafgion/Documents/Projects/camel-1c-integration/target/api.feature',
      'C:/Users/bafgion/Documents/Projects/camel-1c-integration/target/ui.feature',
      'C:/Users/bafgion/Documents/Projects/camel-1c-integration/target/regress.feature',
    ],
      projects: ['target', 'examples', 'test', 'demo', 'sandbox'],
    })
  }

  function syncViewportLayout() {
    if (typeof window === 'undefined') return
    viewportStore.syncWindowSize(
      window.innerWidth,
      window.innerHeight,
      shouldAutoCompactToolbar(window.innerWidth),
    )
    layoutStore.patchLocal({ bottomPanelHeight: clampBottomPanelHeight(bottomPanelHeight, viewportHeight) })
    stepsPanelHeight = clampStepsPanelHeight(stepsPanelHeight, viewportHeight)
    syncToolbarDensity()
  }

  function syncToolbarDensity() {
    if (!actionBarEl) return
    const barWidth = actionBarEl.getBoundingClientRect().width
    const urlBlock = actionBarEl.querySelector('.url-block') as HTMLElement | null
    const urlWidth = urlBlock?.getBoundingClientRect().width ?? 0
    const available = barWidth - urlWidth - 48
    viewportStore.setToolbarIconOnly(available < toolbarIconOnlyThreshold(barWidth))
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
    const tone = journalStore.snapshot().statusTone
    if (tone === 'busy' || tone === 'error') return
    if (projectPath) {
      journalStore.setStatus(projectPath.replace(/\\/g, '/'), 'normal')
      return
    }
    journalStore.setStatus(tr('statusBar.default'), 'normal')
  }

  function toggleMenu(name: string, e: MouseEvent) {
    e.stopPropagation()
    menuStore.toggle(name)
  }

  function closeMenu() {
    menuStore.close()
  }

  function runMenuAction(action: () => void) {
    closeMenu()
    action()
  }

  async function refreshArtifacts() {
    try {
      projectStore.setArtifacts(await ProjectArtifacts())
    } catch {
      projectStore.setArtifacts(new gui.ProjectArtifacts())
    }
    await refreshAllureStatus()
  }

  async function refreshEditorSteps(textSnapshot = editorText, textVersionSnapshot = editorTextVersion) {
    try {
      const steps = await ParseEditorSteps(textSnapshot)
      if (!isEditorAnalysisSnapshotVisible(textVersionSnapshot, editorTextVersion)) return
      editorStore.setSteps(steps, textVersionSnapshot)
    } catch {
      if (!isEditorAnalysisSnapshotVisible(textVersionSnapshot, editorTextVersion)) return
      editorStore.clearSteps(textVersionSnapshot)
    }
    monaco?.refreshInlayHints()
  }

  async function serveAllureReport(path = '') {
    appendLog(tr('journal.reports.allureServe'))
    const result = await startRunResultJob('allure-serve-finished', () => StartServeAllure(path), () => currentProjectVersion)
    if (result.output) appendLog(result.output.trimEnd())
    if (result.error) {
      appendLog(tr('journal.error.generic', { error: result.error }))
    } else {
      reportsStore.setAllureServeRunning(true)
    }
    await refreshAllureStatus(path)
  }

  async function refreshAllureStatus(dir = '') {
    try {
      const status = await AllureStatus(dir || projectArtifacts.allureDir || '')
      reportsStore.setAllureStatus(status.installed !== false, !!status.running)
    } catch {
      reportsStore.setAllureStatus(true, false)
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
    layoutStore.openBottomPanel()
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
    layoutStore.openBottomPanel()
    layoutStore.setBottomTab('journal')
    await executeRun(runFormFromMode(lastRun, 'single', { dryRun: false, scenario: req.scenario, html: true }), [featurePath])
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

  async function refreshCatalog() {
    if (!projectPath) return
    try {
      const info = await RefreshProject()
      applyProjectScan(info)
      projectStore.setScenarios(await ListScenarioTitles().catch((): string[] => []))
      appendLog(tr('journal.catalog.refreshed'))
      setStatus(tr('journal.catalog.refreshedShort'), 'success')
    } catch (e: unknown) {
      const msg = e instanceof Error ? e.message : String(e)
      appendLog(tr('journal.error.generic', { error: msg }))
      setStatus(msg, 'error')
    }
  }

  function applyProjectScan(info: gui.ProjectInfo, fallbackPath = '') {
    const raw = (info as unknown as { version?: number }).version
    const version = typeof raw === 'number' && raw > 0 ? raw : currentProjectVersion
    setProjectState({
      path: info.path || fallbackPath || projectPath,
      version,
      features: info.features || [],
      tags: info.tags || [],
      featureTags: info.featureTags || {},
    })
  }

  async function refreshProject() {
    if (!projectPath) return
    const info = await RefreshProject()
    applyProjectScan(info)
    catalogStore.remapBatchSelected(features)
    testClientStore.setClients(await ListTestClients().catch((): string[] => []))
    projectStore.setScenarios(await ListScenarioTitles().catch((): string[] => []))
    await refreshInstalledPlugins()
  }

  async function refreshInstalledPlugins() {
    if (!projectPath) {
      pluginsStore.clear()
      return
    }
    try {
      pluginsStore.setInstalled(await ListPlugins())
    } catch {
      pluginsStore.clear()
    }
  }

  function hasVanessaPlugin(): boolean {
    return pluginsStore.hasVanessa()
  }

  function pluginLabel(plugin: gui.PluginEntryDTO): string {
    if (plugin.vanessa) return 'Vanessa'
    return plugin.description || plugin.id || plugin.name
  }

  function pluginRunTitle(name: string): string {
    const entry = pluginsStore.findByName(name)
    if (entry) return pluginLabel(entry)
    return name
  }

  async function moveFeatureInCatalog(src: string, destDir: string) {
    if (!src || !destDir) return
    if (!isRealFeaturePath(src)) {
      appendLog(tr('journal.project.saveOrOpenFirst'))
      return
    }
    try {
      const newPath = await MoveFeature(src, destDir)
      const wasActive = activeTab === src
      if (tabs.some((t) => t.path === src)) {
        tabsStore.mapTabs((tabs) => tabs.map((t) => (t.path === src ? { ...t, path: newPath } : t)))
        if (wasActive) tabsStore.setActiveTab(newPath)
      }
      catalogStore.mapBatchSelected((paths) => paths.map((p) => (p === src ? newPath : p)))
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
      dialogsStore.open('showOpenProject')
      return
    }
    await openProjectAt(path)
  }

  function trimTabsMemory() {
    const state = tabsStore.snapshot()
    const activePath = state.activeTab === WELCOME_KEY ? '' : state.activeTab
    const retainedTabs = trimRetainedTabBodies(state.tabs, activePath)
    tabsStore.setTabs(retainedTabs)
    monaco?.retainTabs(pathsToRetainModels(retainedTabs, activePath))
  }

  function warnManyOpenTabs() {
    if (tabs.length >= MAX_OPEN_EDITOR_TABS) {
      appendLog(tr('journal.file.manyTabs', { count: tabs.length }))
    }
  }

  function activeTabLiveText(): string {
    return monaco?.getEditorText() ?? get(editorStore).text ?? editorText
  }

  function syncTabContent(tabPath: string) {
    if (!tabPath || tabPath === WELCOME_KEY) return
    tabsStore.mapTabs((tabs) => tabs.map((t) => {
      if (t.path !== tabPath) return t
      const stored = tabEditorText(t)
      const liveText = tabPath === activeTab ? pickPersistText(stored, activeTabLiveText()) : stored
      const dirty = liveText !== t.content
      if (!dirty) {
        if (!t.dirty && t.draft === undefined) return t
        return { ...t, dirty: false, draft: undefined }
      }
      return { ...t, draft: liveText, dirty: true }
    }))
  }

  function syncActiveTabContent() {
    if (isWelcome || !activeTab) return
    syncTabContent(activeTab)
  }

  function markActiveTabSaved(text: string, tabPath = activeTab) {
    if (isWelcome || !tabPath) return
    tabsStore.mapTabs((tabs) => tabs.map((t) =>
      t.path === tabPath ? { ...t, content: text, dirty: false, draft: undefined } : t,
    ))
  }

  async function checkActiveTabDiskStale() {
    if (isWelcome || !activeTab || isUntitled(activeTab) || playing || recording) return
    const pathAtStart = activeTab
    const tab = tabs.find((t) => t.path === pathAtStart)
    if (!tab || tab.dirty) return
    const baseline = tab.draft ?? tab.content
    try {
      const disk = await ReadFeature(pathAtStart)
      if (disk === baseline) return
      const ok = await askConfirm({
        title: tr('confirm.diskChanged.title'),
        message: tr('confirm.diskChanged.message', { name: basename(pathAtStart) }),
        confirmLabel: tr('confirm.diskChanged.confirmLabel'),
        danger: true,
      })
      if (!ok) return
      const currentActive = activeTab
      tabsStore.mapTabs((tabs) => tabs.map((t) =>
        t.path === pathAtStart ? { ...t, content: disk, dirty: false, draft: undefined, unloaded: false } : t,
      ))
      if (currentActive === pathAtStart) {
        await applyEditorText(disk, { saved: true, switchTab: true, tabPath: pathAtStart, skipValidate: true })
      }
      appendLog(tr('journal.file.reloaded', { name: basename(pathAtStart) }))
    } catch {
      /* ignore */
    }
  }

  async function ensureRecordingTabSwitchAllowed(path: string): Promise<boolean> {
    if ((!recording && !captureFinalizing) || !recordingTargetPath) return true
    if (recordingTabSwitchAllowed(recording, recordPaused, recordingTargetPath, path, captureFinalizing)) return true
    if (confirmDialogStore.shouldSkipRecordTabSwitchConfirm()) return true
    const ok = await askConfirm({
      title: tr('confirm.recordingActive.title'),
      message: tr('confirm.recordingActive.message', {
        from: basename(recordingTargetPath),
        to: basename(path),
      }),
      confirmLabel: tr('confirm.recordingActive.confirmLabel'),
      dontAskAgainLabel: tr('confirm.recordingActive.dontAskAgainLabel'),
    })
    return ok
  }

  function cancelPendingFeatureLoads() {
    tabsStore.bumpLoadFeatureGeneration()
  }

  async function loadFeature(
    path: string,
    opts?: { skipRecordingGuard?: boolean; forceActivate?: boolean },
  ) {
    if (!opts?.skipRecordingGuard) {
      const allowed = await ensureRecordingTabSwitchAllowed(path)
      if (!allowed) return
    }
    const generation = tabsStore.bumpLoadFeatureGeneration()
    const leavingTab = activeTab
    if (leavingTab && !isWelcome && leavingTab !== path) {
      syncTabContent(leavingTab)
    }
    const existing = tabs.find((t) => isSameRecordTab(t.path, path))
    if (existing) {
      if (!opts?.forceActivate && isSameRecordTab(path, activeTab) && !tabNeedsDiskReload(existing)) {
        return
      }
      let text = tabEditorText(existing)
      if (tabNeedsDiskReload(existing)) {
        try {
          text = await ReadFeature(path)
          if (!tabsStore.isLoadFeatureGenerationCurrent(generation)) return
          tabsStore.mapTabs((tabs) => tabs.map((t) =>
            isSameRecordTab(t.path, path) ? { ...t, content: text, dirty: false, draft: undefined, unloaded: false } : t,
          ))
        } catch (e: any) {
          appendLog(tr('journal.file.openError', { error: String(e) }))
          return
        }
      }
      if (!tabsStore.isLoadFeatureGenerationCurrent(generation)) return
      tabsStore.patch({ welcomeTabVisible: false, activeTab: existing.path })
      await applyEditorText(text, { saved: !existing.dirty, switchTab: true, tabPath: existing.path, skipValidate: true })
      trimTabsMemory()
      syncStepsPanelCollapsedFromPrefs()
      schedulePersistSession()
      return
    }
    if (isUntitled(path)) {
      appendLog(tr('journal.file.openError', { error: `untitled tab is not open: ${untitledLabel(path)}` }))
      return
    }
    try {
      const diskContent = await ReadFeature(path)
      if (!tabsStore.isLoadFeatureGenerationCurrent(generation)) return
      let content = diskContent
      let dirty = false
      try {
        const draft = await LoadFeatureDraft(path)
        if (!tabsStore.isLoadFeatureGenerationCurrent(generation)) return
        if (draft && draft.trim() !== diskContent.trim()) {
          content = draft
          dirty = true
          appendLog(tr('journal.file.draftRestored', { name: basename(path) }))
        }
      } catch {
        /* no draft */
      }
      if (!tabsStore.isLoadFeatureGenerationCurrent(generation)) return
      tabsStore.appendTab({ path, content, dirty })
      warnManyOpenTabs()
      tabsStore.patch({ welcomeTabVisible: false, activeTab: path })
      await rememberFeature(path)
      const recents = await loadRecents()
      recentsStore.patch({ features: recents.features })
      await applyEditorText(content, { saved: !dirty, switchTab: true, tabPath: path, skipValidate: true })
      trimTabsMemory()
      syncStepsPanelCollapsedFromPrefs()
      schedulePersistSession()
    } catch (e: any) {
      appendLog(tr('journal.file.openError', { error: String(e) }))
    }
  }

  function selectTab(path: string) {
    if (path === WELCOME_KEY) {
      if (!canShowWelcome(tabs)) return
      cancelPendingFeatureLoads()
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
      tabsStore.patch({ activeTab: WELCOME_KEY })
      trimTabsMemory()
      return
    }
    loadFeature(path, { forceActivate: true })
  }

  function closeWelcomeTab() {
    if (tabs.length > 0) {
      tabsStore.setWelcomeVisible(false)
      if (activeTab === WELCOME_KEY) {
        const next = tabs[tabs.length - 1]
        if (!next) return
        void loadFeature(next.path, { forceActivate: true })
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
      tabsStore.setPendingCloseTab(path)
      return
    }
    finalizeCloseTab(path)
  }

  function finalizeCloseTab(path: string) {
    evictFeatureSymbolCache(path)
    monaco?.releaseTab(path)
    const reduced = tabsStore.closePath(path)
    trimTabsMemory()
    if (reduced.openNextPath) {
      void loadFeature(reduced.openNextPath, { forceActivate: true })
    } else if (reduced.showWelcome) {
      cancelPendingFeatureLoads()
      void applyEditorText('', { switchTab: true, tabPath: null, skipValidate: true })
      clearEditorValidation()
    }
    schedulePersistSession()
  }

  async function saveAndCloseTab() {
    if (!pendingCloseTab) return
    const path = pendingCloseTab
    if (activeTab !== path) {
      await loadFeature(path)
    }
    await saveFeature()
    const tab = tabsStore.snapshot().tabs.find((t) => t.path === path)
    if (tab && !tab.dirty) {
      finalizeCloseTab(path)
    }
  }

  async function discardAndCloseTab() {
    if (!pendingCloseTab) return
    const path = pendingCloseTab
    try {
      await ClearFeatureDraft(path)
    } catch (e) {
      console.warn('failed to clear discarded feature draft', e)
    }
    finalizeCloseTab(path)
  }

  function cancelCloseTab() {
    tabsStore.setPendingCloseTab(null)
  }

  async function saveFeatureAs() {
    if (!activeTab || isWelcome) return
    const pathAtStart = activeTab
    const picked = await PickSaveFile(tr('filePicker.saveAs'), basename(pathAtStart))
    if (!picked) return
    try {
      const text = monaco?.getEditorText() ?? editorText
      await SaveFeature(picked, text)
      const stillActive = activeTab === pathAtStart
      tabsStore.mapTabs((tabs) => tabs.map((t) =>
        t.path === pathAtStart ? { path: picked, content: text, dirty: false, draft: undefined } : t,
      ))
      if (stillActive) {
        tabsStore.setActiveTab(picked)
        await applyEditorText(text, { saved: true, switchTab: true, tabPath: picked, skipValidate: true })
      }
      monaco?.releaseTab(pathAtStart)
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
    const pathAtStart = activeTab
    if (isUntitled(activeTab)) {
      await saveFeatureAs()
      return
    }
    try {
      let text = monaco?.getEditorText() ?? editorText
      if (dialogBinds.bindEditorSettings.formatOnSave) {
        await monaco?.formatDocument()
        text = monaco?.getEditorText() ?? text
        if (activeTab === pathAtStart) {
          editorStore.setText(text)
        }
      }
      const { text: autoFixed, count: autoFixCount } = await runScenarioHintsAutoFix(text)
      if (autoFixCount > 0) {
        text = autoFixed
        if (activeTab === pathAtStart) {
          editorStore.setTextWithBump(text)
          const textVersion = get(editorStore).textVersion
          await monaco?.setContent(text, { path: pathAtStart, generation: textVersion })
        }
        appendLog(tr('journal.hint.autoFixed', { count: autoFixCount }))
      }
      await SaveFeature(pathAtStart, text)
      if (activeTab === pathAtStart) {
        editorStore.setText(text)
      }
      markActiveTabSaved(text, pathAtStart)
      try {
        await ClearFeatureDraft(pathAtStart)
      } catch {
        /* ignore */
      }
      appendLog(tr('journal.file.saved', { name: basename(pathAtStart) }))
      setStatus(tr('journal.status.saved'), 'success')
    } catch (e: any) {
      appendLog(tr('journal.file.saveError', { error: String(e) }))
      setStatus(tr('journal.status.saveError'), 'error')
    }
  }

  function scheduleValidateEditor(delayMs = 300) {
    diagnosticsStore.scheduleValidateEditor(() => void validateEditor(), delayMs)
  }

  function clearEditorValidation() {
    if (activeTab && activeTab !== WELCOME_KEY) {
      diagnosticsStore.clearIssuesForTab(activeTab)
    }
    diagnosticsStore.setIssues([])
    diagnosticsStore.setStepStatusError(false)
    monaco?.setMarkers([])
    if (statusMessage === tr('journal.status.scenarioError')) {
      setStatus('', 'normal')
    }
  }

  function syncStepStatusFromIssues(issues: gui.ValidationIssue[]) {
    const hasIssues = issues.length > 0
    diagnosticsStore.setStepStatusError(hasIssues)
    if (!hasIssues && statusMessage === tr('journal.status.scenarioError')) {
      setStatus('', 'normal')
    }
  }

  async function validateEditor() {
    if (isWelcome || !activeTab || activeTab === WELCOME_KEY) {
      clearEditorValidation()
      return
    }
    const generation = diagnosticsStore.bumpValidateGeneration()
    const tabAtStart = activeTab
    const textAtStart = editorText
    const textVersionAtStart = editorTextVersion
    const started = perfNow()
    try {
      const analysis = await AnalyzeEditorContent(textAtStart, dialogBinds.bindEditorSettings.scenarioHints)
      if (generation !== diagnosticsStore.validateGeneration() || tabAtStart !== activeTab || textVersionAtStart !== editorTextVersion) return
      const issues = analysis?.issues || []
      diagnosticsStore.setIssues(issues)
      if (tabAtStart) {
        diagnosticsStore.setIssuesForTab(tabAtStart, issues)
      }
      monaco?.setMarkers(issues)
      if (isEditorAnalysisSnapshotVisible(textVersionAtStart, editorTextVersion)) {
        editorStore.setSteps(analysis?.steps || [], textVersionAtStart)
        monaco?.refreshInlayHints()
      }
      if (generation !== diagnosticsStore.validateGeneration() || tabAtStart !== activeTab || textVersionAtStart !== editorTextVersion) return
      if (issues.length > 0) {
        diagnosticsStore.setStepStatusError(true)
        setStatus(tr('journal.status.scenarioError'), 'error')
      } else {
        syncStepStatusFromIssues(issues)
      }
      if (dialogBinds.bindEditorSettings.scenarioHints) {
        const hints = (analysis?.hints || [])
          .filter((h) => !diagnosticsStore.isHintDismissed(hintDismissKey(h)))
          .filter((h) => filterScenarioHints([h], dialogBinds.bindEditorSettings).length > 0)
        diagnosticsStore.setHints(hints)
      } else {
        diagnosticsStore.setHints([])
      }
    } catch {
      if (generation !== diagnosticsStore.validateGeneration() || tabAtStart !== activeTab || textVersionAtStart !== editorTextVersion) return
      diagnosticsStore.setIssues([])
      if (tabAtStart) {
        diagnosticsStore.setIssuesForTab(tabAtStart, [])
      }
      if (isEditorAnalysisSnapshotVisible(textVersionAtStart, editorTextVersion)) {
        editorStore.clearSteps(textVersionAtStart)
        monaco?.refreshInlayHints()
      }
      syncStepStatusFromIssues([])
      diagnosticsStore.setHints([])
    } finally {
      perfMark('validateEditor', started)
    }
  }

  function gotoEditorLine(line: number) {
    monaco?.gotoLine(line)
    editorStore.setCursorLine(line)
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

  async function onEditorChange(event: {
    path: string | null
    modelUri: string | null
    text: string
    source?: string
    modelVersion?: number
    hydrationGeneration?: number
  }) {
    const activePath = isWelcome ? null : activeTab
    const activeModelUri = monaco?.getActiveModelUri?.() ?? null
    if (!shouldAcceptMonacoChange({
      activePath,
      eventPath: event.path,
      activeModelUri,
      eventModelUri: event.modelUri,
      source: event.source,
    })) {
      return
    }
    editorStore.setTextWithBump(event.text)
    const text = event.text
    syncActiveTabContent()
    schedulePersistSession()
    if (dialogBinds.bindEditorSettings.validateOnType && !isWelcome) {
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

  function openRunDialog(title: string, defaults: Partial<RunForm>, mode: RunFormMode = 'single') {
    runDialogStore.open(title, dialogScenarioNames())
    const cursorScenario = cursorScenarioName()
    const defaultScenario =
      defaults.scenario ??
      (mode === 'tag' || mode === 'batch' ? '' : cursorScenario || lastRun.scenario || '')
    dialogBinds.syncRunFormBind(mode, lastRun, {
      ...defaults,
      baseUrl: lastRun.baseUrl || startURL || '',
      scenario: defaultScenario,
    })
    dialogsStore.open('showRun')
  }

  function openVanessaDialog(dry: boolean, preferRerun = false) {
    vanessaRunStore.prepareDialog({
      dry,
      preferRerun,
      scenario: cursorScenarioName(),
      dialogScenarios: dialogScenarioNames(),
    })
    dialogBinds.syncVanessaBindLocals()
    dialogsStore.open('showVanessaRun')
  }

  async function confirmVanessaRun() {
    dialogsStore.close('showVanessaRun')
    vanessaRunStore.patch({ plannedTotal: Math.max(1, features.length) })
    dialogsStore.open('showVanessaMonitor')
    vanessaRunStore.setSnapshot(new gui.VanessaRunSnapshotDTO())
    StartVanessaRun(buildVanessaPluginRequest())
  }

  async function executeRun(opts: RunForm, targets: string[] = []) {
    if (playing) {
      appendLog(tr('journal.run.alreadyRunning'))
      return
    }
    syncActiveTabContent()

    let runTargets = [...targets]
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

    runFormStore.setLastRun({ ...runOpts })
    dialogsStore.close('showRun')
    layoutStore.openBottomPanel()
    layoutStore.setBottomTab('journal')

    const initialPlayingLabel =
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

    const runIsPlaying = !runOpts.dryRun
    runnerStore.setDryRunActive(runOpts.dryRun)
    const initialRunProgressTotal = Math.max(1, diskTargets.length || (runTargets.length > 0 ? runTargets.length : 1))
    runnerStore.setRunId('')
    if (runIsPlaying) {
      runnerStore.start('')
      runnerStore.progress(initialRunProgressTotal, 0, initialPlayingLabel)
    } else {
      runnerStore.progress(initialRunProgressTotal, 0, initialPlayingLabel)
      runnerStore.stop()
    }
    runnerStore.setCancelling(false)
    runnerStore.setLogStreaming(true)
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

    let result: gui.RunResult = gui.RunResult.createFrom({ output: '', error: '', entries: [] })
    let runThrown: unknown = null
    let journalStreamed = false
    try {
      result = await startRunResultJob('run-finished', () => StartRun({
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
        }), () => currentProjectVersion)
      journalStreamed = runLogStreaming
    } catch (err) {
      runThrown = err
    } finally {
      runnerStore.setLogStreaming(false)
      runnerStore.setDryRunActive(false)
      runnerStore.setCancelling(false)
      runnerStore.setRunId('')
      runnerStore.progress(0, 0, '')
      runnerStore.stop()
    }

    if (runThrown) {
      appendLog(tr('journal.error.generic', { error: String(runThrown) }))
      setStatus(tr('journal.status.testError'), 'error')
      layoutStore.setBottomTab('error')
      return
    }

    if (result.output && !journalStreamed) appendLog(result.output.trimEnd())
    if (result.error) {
      if (/context canceled/i.test(result.error)) {
        appendLog(tr('journal.run.stopped'))
        setStatus(tr('journal.status.testStopped'), 'busy')
        layoutStore.setBottomTab('journal')
      } else {
        appendLog(tr('journal.error.generic', { error: result.error }))
        setStatus(tr('journal.status.testError'), 'error')
        layoutStore.setBottomTab('error')
      }
    } else {
      appendLog(tr('journal.run.done'))
      setStatus(tr('journal.status.testDone'), 'success')
      uiPrefsStore.patch({ welcomePlayedSuccess: true })
      if (showOnboardingTour && runOpts.dryRun) {
        onboardingTourStore.patch({ dryRunDone: true })
      }
      void persistSettings()
      layoutStore.setBottomTab(runOpts.dryRun ? 'journal' : 'results')
    }
    await refreshRunResults()
    if (result.entries?.length) {
      let batchResults = remapRunResultPaths(result.entries, diskTargets, runTargets)
      let lastError = pickLastRunError(batchResults)
      if (
        !lastError &&
        result.error &&
        !/context canceled/i.test(result.error)
      ) {
        const featurePath =
          runTargets.length === 1
            ? runResultFeaturePath(runTargets[0])
            : diskTargets.length === 1
              ? diskTargets[0]
              : tr('journal.run.defaultLabel')
        lastError = buildSyntheticRunError({
          featurePath,
          scenario: runOpts.scenario || cursorScenarioName() || undefined,
          message: result.error,
          runner: runOpts.dryRun ? 'dry-run' : runOpts.engine || 'playwright',
        })
        batchResults = [...batchResults, lastError]
      }
      runnerStore.setLastRunSession(runSince, batchResults, lastError)
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
    const allowPartialCanceledReport =
      runCancelled &&
      (result.entries?.some((e) => e.success || (e.message || '').trim().length > 0) ?? false)
    const actualHtmlPath = (result.reportPath || result.htmlPath || '').trim()
    if (runOpts.html && actualHtmlPath && (!runCancelled || allowPartialCanceledReport)) {
      try {
        if (await ArtifactExists(actualHtmlPath)) {
          await openHtmlReport(actualHtmlPath)
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
    uiPrefsStore.patch({ runDialogConfirmed: true })
    void persistSettings()
    executeRun(dialogBinds.bindRunForm)
  }

  async function validateProject(browser: boolean, browserName = settingsBrowser || 'chromium', targets: string[] = [], validationMode: 'static' | 'flow' = 'static') {
    if (!projectPath) return
    appendLog(browser ? tr('journal.validate.startingBrowser') : tr('journal.validate.starting'))
    layoutStore.openBottomPanel()
    layoutStore.setBottomTab(browser ? 'validate' : 'journal')
    validateDialogStore.setCliLog('')
    diagnosticsStore.setBrowserPanelIssues([])

    if (browser) {
      try {
        const issues = await ValidateBrowser(
          gui.ValidateRequest.createFrom({
            browser: browserName || 'chromium',
            skipBrowser: false,
            targets,
            mode: validationMode,
          }),
        )
        const panelIssues = issues || []
        diagnosticsStore.setBrowserPanelIssues(panelIssues)
        const missing = panelIssues.filter((i) => i.status === 'missing' || !i.status).length
        const warnings = panelIssues.filter((i) => i.status === 'warning').length
        const found = panelIssues.filter((i) => i.status === 'found').length
        appendLog(tr('journal.validate.browserSummary', { found, warnings, missing }))
        setStatus(missing > 0 ? tr('journal.status.validateBrowserErrors') : tr('journal.status.validateBrowserDone'), missing > 0 ? 'error' : 'success')
      } catch (err) {
        const msg = err instanceof Error ? err.message : String(err)
        validateDialogStore.setCliLog(msg)
        appendLog(tr('journal.error.generic', { error: msg }))
        setStatus(tr('journal.status.validateError'), 'error')
      }
    } else {
      const result = await startRunResultJob('validate-finished', () => StartValidate({
          browser: browserName || 'chromium',
          skipBrowser: true,
          targets,
          mode: validationMode,
        }), () => currentProjectVersion)
      if (result.output) {
        validateDialogStore.setCliLog(result.output.trimEnd())
        appendLog(validateCliLog)
      }
      if (result.error) {
        validateDialogStore.appendCliLog(result.error)
        appendLog(tr('journal.error.generic', { error: result.error }))
        setStatus(tr('journal.status.validateError'), 'error')
      } else {
        appendLog(tr('journal.validate.done'))
        setStatus(tr('journal.status.validateDone'), 'success')
        if (showOnboardingTour && !browser) {
          onboardingTourStore.patch({ validateDone: true })
        }
      }
    }
    if (!isWelcome && activeTab) await validateEditor()
  }

  function openValidateDialog(syntaxOnly: boolean) {
    if (!projectPath) return
    const scope = !isWelcome && activeTab ? 'current' : 'project'
    validateDialogStore.open(syntaxOnly, settingsBrowser || 'chromium', scope)
    dialogBinds.syncValidateDialogBindLocals()
    dialogsStore.open('showValidate')
  }

  async function confirmValidate(payload: { browser: string; syntaxOnly: boolean; flowAware: boolean; scope: 'project' | 'current' }) {
    dialogsStore.close('showValidate')
    let targets: string[] = []
    if (payload.scope === 'current' && activeTab && !isWelcome) {
      targets = await materializeRunTargets([activeTab])
      if (targets.length === 0) {
        appendLog(tr('journal.validate.prepareFailed'))
        setStatus(tr('journal.status.validateError'), 'error')
        return
      }
    }
    await validateProject(!payload.syntaxOnly, payload.browser, targets, payload.flowAware ? 'flow' : 'static')
  }

  function openInitProjectDialog() {
    if (!projectPath) return
    dialogsStore.open('showInitProject')
  }

  async function confirmInitProject() {
    dialogsStore.close('showInitProject')
    await initProject()
  }

  function openNewProjectWizard() {
    dialogsStore.open('showNewProjectWizard')
  }

  async function confirmNewProjectWizard(opts: NewProjectWizardResult) {
    dialogsStore.close('showNewProjectWizard')
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
    testClientStore.clearSuggestName()
    const clients = await ListTestClients().catch((): string[] => [])
    testClientStore.setClients(clients)
    dialogBinds.bindTestClientSelection = runForm.testClient || clients[0] || ''
    testClientStore.patch({ selection: dialogBinds.bindTestClientSelection })
    dialogsStore.open('showTestClient')
  }

  async function openTestClientDialogForCapture() {
    if (!projectPath) return
    if (!browserOpen && !recording) {
      appendLog(tr('journal.browser.openForSession'))
      return
    }
    const clients = await ListTestClients().catch((): string[] => [])
    testClientStore.setClients(clients)
    const base = (runForm.testClient || dialogBinds.bindTestClientSelection || 'session').trim() || 'session'
    dialogBinds.bindTestClientSelection = clients.includes(base) ? base : ''
    testClientStore.patch({ suggestName: base, selection: dialogBinds.bindTestClientSelection })
    dialogsStore.open('showTestClient')
  }

  function useTestClient(name: string) {
    dialogBinds.bindTestClientSelection = name
    testClientStore.patch({ selection: name })
    runFormStore.patchLastRun({ testClient: name })
    runFormStore.patchRunForm({ testClient: name })
    appendLog(tr('journal.testClient.selected', { name }))
    dialogsStore.close('showTestClient')
  }

  async function openStepsDialog() {
    if (recordingBlocksManualTools()) {
      appendLog(tr('journal.record.pauseToInsertStep'))
      return
    }
    dialogsStore.open('showSteps')
  }

  function openVanessaSettingsDialog() {
    if (!projectPath) return
    dialogsStore.open('showVanessaSettings')
  }

  function openStepsHelp(query: unknown = '') {
    if (recordingBlocksManualTools()) {
      appendLog(tr('journal.record.pauseToOpenStepsHelp'))
      return
    }
    stepsHelpDialogStore.setQuery(typeof query === 'string' ? query : '')
    dialogsStore.open('showStepsHelp')
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
    dialogsStore.close('showSteps')
    dialogsStore.close('showStepsHelp')
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
    dialogsStore.open('showSnippetPalette')
  }

  async function applyEditorText(
    text: string,
    options?: { saved?: boolean; switchTab?: boolean; tabPath?: string | null; skipValidate?: boolean },
  ) {
    editorStore.setTextWithBump(text)
    const textVersion = get(editorStore).textVersion
    if (options?.switchTab) {
      const markerPath = options.tabPath ?? null
      const issues = markerPath ? (editorValidationByTab[markerPath] ?? []) : []
      diagnosticsStore.setIssues(issues)
      syncStepStatusFromIssues(issues)
      monaco?.activateTab(options.tabPath ?? null, text, textVersion)
      monaco?.setMarkers(issues)
      if (options?.saved) {
        markActiveTabSaved(text, options.tabPath ?? activeTab)
      } else if (options?.tabPath) {
        tabsStore.mapTabs((tabs) => tabs.map((t) =>
          t.path === options.tabPath ? { ...t, draft: text, dirty: true } : t,
        ))
      }
    } else {
      await monaco?.setContent(text, { path: editorValuePath, generation: textVersion })
      if (options?.saved) {
        markActiveTabSaved(text)
      } else {
        syncActiveTabContent()
      }
    }
    if (options?.skipValidate) {
      void refreshEditorSteps(text, textVersion)
      return
    }
    await validateEditor()
  }

  async function refactorUpdateUrls() {
    if (isWelcome) return
    dialogsStore.open('showRefactorUrl')
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

  async function openExportDialog() {
    if (!activeTab || isWelcome) return
    if (isUntitled(activeTab)) {
      await saveFeatureAs()
      if (!activeTab || isWelcome || isUntitled(activeTab)) return
    }
    featureDialogStore.openExport(activeTab)
    dialogsStore.open('showExport')
  }

  function openImportFeaturesDialog() {
    if (!projectPath) return
    featureDialogStore.openImport(projectPath.replace(/\\/g, '/'))
    dialogBinds.bindImportDestDir = featureDialogStore.snapshot().importDestDir
    dialogsStore.open('showImportFeatures')
  }

  async function confirmImportFeatures(payload: { destDir: string; paths: string[] }) {
    if (!projectPath || importFeaturesBusy) return
    featureDialogStore.patch({ importFeaturesBusy: true })
    dialogsStore.close('showImportFeatures')
    try {
      await importDroppedFeatures(payload.destDir, payload.paths)
    } finally {
      featureDialogStore.patch({ importFeaturesBusy: false })
    }
  }

  function openImportDialog() {
    if (!projectPath) return
    dialogsStore.open('showImport')
  }

  async function onImportComplete(featurePath: string) {
    await refreshProject()
    await loadFeature(featurePath)
  }

  async function runPlugin(name: string, dry: boolean, opts: Partial<gui.PluginRunRequest> = {}) {
    const label = pluginLabel(installedPlugins.find((p) => p.name === name) || { name, id: name, vanessa: false } as gui.PluginEntryDTO)
    appendLog(dry ? tr('journal.plugin.runningDry', { label }) : tr('journal.plugin.running', { label }))
    const result = await startRunResultJob('plugin-run-finished', () => StartRunPlugin({
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
      }), () => currentProjectVersion)
    if (result.output) appendLog(result.output.trimEnd())
    if (result.error) appendLog(tr('journal.error.generic', { error: result.error }))
  }

  function openPluginRun(name: string, dry = false) {
    const entry = pluginsStore.findByName(name)
    if (entry?.vanessa || name === 'vanessa') {
      openVanessaDialog(dry)
      return
    }
    pluginRunStore.prepareDialog({
      name,
      dry,
      scenario: cursorScenarioName(),
      dialogScenarios: dialogScenarioNames(),
    })
    dialogBinds.syncPluginRunBindLocals()
    dialogsStore.open('showPluginRun')
  }

  async function confirmPluginRun(payload: { tag: string; scenario: string; dryRun: boolean }) {
    dialogsStore.close('showPluginRun')
    await runPlugin(pluginRunName, payload.dryRun, { tag: payload.tag, scenario: payload.scenario })
  }

  function openBaselineRecordDialog() {
    if (!projectPath) return
    recordFormStore.patch({ recordMode: 'baseline' })
    recordFormStore.patch({ recordOutput: 'recorded.feature' })
    recordFormStore.patch({ recordURL: startURL || recordURL || 'https://example.com' })
    if (activeTab && !isWelcome) {
      recordFormStore.patch({ recordFeatureName: basename(activeTab).replace(/\.feature$/i, '') })
      recordFormStore.patch({ recordScenarioName: tr('dialogs.record.baselineScenarioDefault') })
    } else {
      recordFormStore.patch({ recordFeatureName: tr('dialogs.record.featureDefault') })
      recordFormStore.patch({ recordScenarioName: tr('dialogs.record.baselineScenarioDefault') })
    }
    dialogBinds.syncRecordFormBindLocals()
    dialogBinds.syncRecorderPrefsBindLocals()
    dialogsStore.open('showRecord')
  }

  async function saveBaselineRecord(payload: {
    output: string
    featureName: string
    scenarioName: string
    steps: string[]
  }) {
    if (!projectPath || baselineBusy) return
    recordFormStore.patch({ baselineBusy: true })
    dialogsStore.close('showRecord')
    appendLog(tr('journal.record.baselineCreating'))
    try {
      const result = await startRunResultJob('record-baseline-finished', () => StartRecordBaseline({
          output: payload.output || 'recorded.feature',
          featureName: payload.featureName,
          scenarioName: payload.scenarioName,
          steps: payload.steps,
      }), () => currentProjectVersion)
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
      recordFormStore.patch({ baselineBusy: false })
    }
  }

  async function checkUpdatesOnStartup() {
    if (!settingsCheckUpdatesOnStartup) return
    try {
      const info = await CheckUpdateInfo()
      updateDialogStore.applyCheckResult(info)
      if (info.updateAvailable) {
        dialogsStore.open('showUpdateCheck')
        setStatus(tr('journal.status.updateAvailable'), 'normal')
      }
    } catch {
      /* offline or dev without wails */
    }
  }

  async function maybeCheckUpdatesOnStartup() {
    if (!settingsCheckUpdatesOnStartup) return
    if (!projectPath) {
      updateDialogStore.setPendingStartupCheck(true)
      return
    }
    await checkUpdatesOnStartup()
  }

  async function flushPendingStartupUpdateCheck() {
    if (!updateDialogStore.snapshot().pendingStartupCheck) return
    updateDialogStore.setPendingStartupCheck(false)
    await checkUpdatesOnStartup()
  }

  async function checkUpdates() {
    appendLog(tr('journal.update.checking'))
    try {
      const info = await CheckUpdateInfo()
      updateDialogStore.applyCheckResult(info)
      appendLog(updateCheckMessage)
      if (info.updateAvailable) {
        appendLog(tr('journal.update.release', { url: info.htmlUrl || '—' }))
        if (info.downloadName) appendLog(tr('journal.update.file', { name: info.downloadName }))
      }
      appendLog(tr('journal.validate.done'))
      dialogsStore.open('showUpdateCheck')
    } catch (err) {
      const msg = err instanceof Error ? err.message : String(err)
      updateDialogStore.setMessage(msg, false)
      appendLog(tr('journal.error.generic', { error: msg }))
      dialogsStore.open('showUpdateCheck')
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
    const envelope =
      raw && typeof raw === 'object' && 'payload' in (raw as Record<string, unknown>)
        ? ((raw as Record<string, unknown>).payload as unknown)
        : raw
    const src = (envelope && typeof envelope === 'object' ? envelope : {}) as Record<string, unknown>
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
      updateDialogStore.setProgress(normalizeUpdateProgress(raw))
      const stage = updateDialogStore.snapshot().progress?.stage ?? ''
      if (stage && stage !== lastLoggedStage) {
        lastLoggedStage = stage
        const progress = updateDialogStore.snapshot().progress
        if (progress?.message) appendLog(progress.message)
      }
    }
    EventsOn('update-progress', onProgress)
    EventsOn('update-finished', (raw: unknown) => {
      const envelope =
        raw && typeof raw === 'object' && 'payload' in (raw as Record<string, unknown>)
          ? ((raw as Record<string, unknown>).payload as unknown)
          : raw
      onFinished(envelope as gui.RunResult)
    })
    return () => EventsOff('update-progress', 'update-finished')
  }

  async function applyUpdate() {
    if (updateDownloading) return
    updateDialogStore.setDownloading(true)
    updateDialogStore.setProgress({ stage: 'check', message: tr('journal.update.progressStarting'), percent: 0 })
    appendLog(tr('journal.update.installing'))
    const unbind = listenUpdateProgress((result) => {
      unbind()
      if (result?.error) {
        appendLog(tr('journal.update.installError', { error: result.error }))
        updateDialogStore.patch({ downloading: false, progress: null })
        return
      }
      updateDialogStore.setProgress({ stage: 'restart', message: tr('journal.update.progressRestarting'), percent: 100 })
      appendLog(tr('journal.update.restarting'))
    })
    try {
      await ApplyUpdate()
    } catch (err) {
      unbind()
      appendLog(tr('journal.update.installError', { error: err instanceof Error ? err.message : String(err) }))
      updateDialogStore.patch({ downloading: false, progress: null })
    }
  }

  async function downloadUpdate() {
    if (updateDownloading) return
    updateDialogStore.setDownloading(true)
    updateDialogStore.setProgress({ stage: 'check', message: tr('journal.update.progressDownloading'), percent: 0 })
    appendLog(tr('journal.update.downloading'))
    const unbind = listenUpdateProgress(async (result) => {
      unbind()
      if (result?.error) {
        appendLog(tr('journal.update.downloadError', { error: result.error }))
        updateDialogStore.patch({ downloading: false, progress: null })
        return
      }
      const path = result.output || ''
      appendLog(tr('journal.update.downloaded', { path }))
      updateDialogStore.appendMessage(tr('journal.update.fileSaved', { path }))
      const folder = path.replace(/[\\/][^\\/]+$/, '')
      if (folder) await OpenFolder(folder)
      updateDialogStore.patch({ downloading: false, progress: null })
    })
    try {
      await DownloadUpdate()
    } catch (err) {
      unbind()
      appendLog(tr('journal.update.downloadError', { error: err instanceof Error ? err.message : String(err) }))
      updateDialogStore.patch({ downloading: false, progress: null })
    }
  }

  async function openSettings() {
    const s = await LoadSettings()
    applySettingsFromDTO(s)
    dialogBinds.syncUiPrefsBindLocals()
    dialogBinds.syncRecorderPrefsBindLocals()
    dialogBinds.syncEditorSettingsBindLocal()
    dialogBinds.syncSettingsBindLocals()
    let htmlMode: 'full' | 'light' = 'full'
    if (projectPath) {
      try {
        const pc = await LoadProjectConfig()
        htmlMode = pc.htmlReportOpenMode === 'light' ? 'light' : 'full'
      } catch {
        htmlMode = 'full'
      }
    }
    settingsDialogStore.openWithBaseline(s, htmlMode)
    dialogBinds.syncSettingsDialogBindLocals()
    dialogsStore.open('showSettings')
  }

  function buildCurrentSettingsDTO(): gui.AppSettingsDTO {
    return workspaceSession.buildSettingsDTO()
  }

  async function persistSettings() {
    await workspaceSession.persistSettings()
  }

  async function syncRecordingOptions() {
    try {
      if (recording) {
        const rec = recorderPrefsStore.snapshot()
        await UpdateRecordingOptions(
          rec.filterRecording,
          rec.navOnlyRecording,
          rec.hoverRecord,
          settingsHeadless,
          settingsScrollBeforeClick,
          settingsHoverRecordMinMs,
          !settingsDisableRecordUrlWait,
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
    settingsStore.patch({ headless: next })
    void syncRecordingOptions()
  }

  async function applySettingsCore(closeDialog: boolean) {
    dialogBinds.flushUiPrefsBindLocals()
    dialogBinds.flushRecorderPrefsBindLocals()
    settingsStore.setEditor(dialogBinds.bindEditorSettings)
    dialogBinds.flushSettingsBindLocals()
    dialogBinds.flushSettingsDialogBindLocals()
    setLocale(dialogBinds.bindUiLocale)
    editorStore.setStepsPanelTab(dialogBinds.bindEditorSettings.stepsPanelView)
    if (!stepsPanelVisible) uiPrefsStore.setStepsPanelCollapsed(true)
    else if (stepsPanelCollapsed && stepsPanelVisible) uiPrefsStore.setStepsPanelCollapsed(false)
    setStepHoverEnabled(() => dialogBinds.bindEditorSettings.stepHover)
    runFormStore.applyLastRunFromSettings(settingsBrowser, settingsWorkers, settingsSlowMo)
    if (recording || browserOpen) await syncRecordingOptions()
    else await persistSettings()
    if (projectPath) {
      try {
        const pc = await LoadProjectConfig()
        await SaveProjectConfig(
          gui.ProjectConfigDTO.createFrom({
            ...pc,
            htmlReportOpenMode: settingsDialogStore.snapshot().htmlReportOpenMode,
          }),
        )
      } catch {
        /* project.json may be missing until init */
      }
    }
    monaco?.applyEditorSettings(dialogBinds.bindEditorSettings)
    if (dialogBinds.bindEditorSettings.scenarioHints) {
      await refreshEditorScenarioHints()
    } else {
      diagnosticsStore.setHints([])
    }
    if (dialogBinds.bindEditorSettings.validateOnType && activeTab && !isWelcome) void validateEditor()
    const saved = buildCurrentSettingsDTO()
    if (closeDialog) {
      settingsDialogStore.clearBaseline()
      dialogsStore.close('showSettings')
      appendLog(tr('common.settingsSaved'))
      return saved
    }
    settingsDialogStore.markSaved(saved)
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
    const baseline = settingsDialogStore.snapshot().baseline
    if (baseline) {
      applySettingsFromDTO(baseline)
      monaco?.applyEditorSettings(dialogBinds.bindEditorSettings)
    }
    settingsDialogStore.restoreProjectBaseline()
    dialogBinds.syncSettingsDialogBindLocals()
    settingsDialogStore.clearBaseline()
    dialogsStore.close('showSettings')
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
    recordFormStore.patch({ recordMode: 'live' })
    const recordTestClient = options?.inheritTestClient
      ? (runForm.testClient || testClientSelection || '')
      : ''
    if (activeTab && !isWelcome && activeTab.toLowerCase().endsWith('.feature') && !isUntitled(activeTab)) {
      const output = activeTab.replace(/\\/g, '/')
      recordFormStore.patch({ recordTestClient, recordOutput: output, recordAppendTo: output })
    } else if (projectPath) {
      recordFormStore.patch({
        recordAppendTo: '',
        recordTestClient,
        recordOutput: `${projectPath.replace(/\\/g, '/')}/recorded.feature`,
      })
    } else {
      recordFormStore.patch({ recordAppendTo: '', recordTestClient })
    }
    recordFormStore.patch({ recordURL: recordStartURL() })
    if (activeTab && !isWelcome) {
      recordFormStore.patch({ recordFeatureName: basename(activeTab).replace(/\.feature$/i, '') })
      recordFormStore.patch({ recordScenarioName: tr('dialogs.record.scenarioDefault') })
    } else {
      recordFormStore.patch({ recordFeatureName: tr('dialogs.record.featureDefault') })
      recordFormStore.patch({ recordScenarioName: tr('dialogs.record.scenarioDefault') })
    }
  }

  async function ensureProjectForBrowser(): Promise<boolean> {
    if (projectPath) return true
    const examples = await BundledExamplesPath()
    if (examples) {
      await openProjectAt(examples)
      appendLog(tr('journal.project.autoOpenedExamples'))
      openJournalTab()
      return true
    }
    appendLog(tr('journal.project.openFirst'))
    openProjectDialog()
    return false
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
    dialogBinds.syncRecordFormBindLocals()
    dialogBinds.syncRecorderPrefsBindLocals()
    dialogsStore.open('showRecord')
  }

  async function openBrowser() {
    if (!(await ensureProjectForBrowser())) return
    if (browserOpen || recording) {
      await focusBrowser()
      return
    }
    prepareRecordDialogDefaults({ inheritTestClient: false })
    recordFormStore.patch({ recordURL: recordStartURL() })
    dialogsStore.close('showRecord')
    openJournalTab()
    setStatus(tr('journal.browser.launching'), 'busy')
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
    recordFormStore.patch({ recordAppendTo: '' })
  }

  async function startRecord(opts?: { headed?: boolean }) {
    await persistSettings()
    if (recording) {
      if (browserOpen) await focusBrowser()
      return
    }
    recorderStore.setLastRecordTarget(activeTab && !isWelcome && isUntitled(activeTab) ? activeTab : (recordAppendTo || recordOutput))
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
      headless: opts?.headed === undefined ? settingsHeadless : !opts.headed,
      filterRecording,
      navOnlyRecording,
      hoverRecord,
      appendTo: recordAppendTo,
      testClient: recordTestClient,
      featureName: recordFeatureName,
      scenarioName: recordScenarioName,
      browseOnly: false,
    })
    recordFormStore.patch({ recordAppendTo: '' })
  }

  function resolveRecordFeaturePath(outputPath = ''): string {
    return resolveRecordFeaturePathForUI(
      outputPath,
      lastRecordTarget,
      recordOutput,
      recordAppendTo,
      projectPath,
    )
  }

  async function prepareRecordEditorTab(outputPath = '', uiTargetPath = '') {
    const tabsSnap = tabsStore.snapshot()
    const action = resolveRecordEditorPrepareAction({
      activeTab: tabsSnap.activeTab,
      welcomeKey: WELCOME_KEY,
      appendPath: recordAppendTo,
      backendOutputPath: outputPath,
      lastRecordTarget,
      recordOutput,
      projectPath,
      tabPaths: tabsSnap.tabs.map((tab) => tab.path),
      uiTargetPath,
      defaultUntitledName: 'zapis.feature',
    })
    if (action.kind === 'load') {
      await loadFeature(action.path, { skipRecordingGuard: true, forceActivate: true })
      return
    }
    if (action.kind === 'openUntitled') {
      await openUntitledTab(
        buildFeatureTemplate({
          title: recordFeatureName,
          scenario: recordScenarioName,
          startUrl: recordURL || startURL || 'https://example.com',
        }),
        action.displayName,
      )
    }
  }

  async function syncMonacoAfterMount() {
    if (!monaco) return
    if (isWelcome || !activeTab) {
      await monaco.activateTab(null, editorText, editorTextVersion)
      return
    }
    const tab = tabs.find((t) => t.path === activeTab)
    const text = tab ? tabEditorText(tab) : editorText
    await monaco.activateTab(activeTab, text, editorTextVersion)
  }

  function dismissRecorderPicker() {
    dialogsStore.close('showPickerStep')
    pickerDialogStore.clear()
  }

  function recordStepApplyGate(eventRecordSessionId = ''): {
    recording: boolean
    captureFinalizing: boolean
    activeRecordSessionId: string
    eventRecordSessionId?: string
    lineByIndexCount: number
    recordingTargetPath: string
  } {
    const snap = get(recorderStore)
    return {
      recording: snap.recording,
      captureFinalizing: snap.captureFinalizing,
      activeRecordSessionId: snap.recordSessionId,
      eventRecordSessionId,
      lineByIndexCount: Object.keys(snap.liveRecordStepLines).length,
      recordingTargetPath: snap.targetPath,
    }
  }

  function shouldApplyRecordStepEventLocal(event: RecordStepEvent, eventRecordSessionId = ''): boolean {
    return shouldApplyRecordStepEvent(recordStepApplyGate(eventRecordSessionId), event)
  }

  function bufferEarlyRecordStep(event: RecordStepEvent, targetPath: string, recordSessionId: string) {
    pendingEarlyRecordSteps = [...pendingEarlyRecordSteps.slice(-99), { event, targetPath, recordSessionId }]
  }

  async function flushPendingEarlyRecordSteps() {
    if (pendingEarlyRecordSteps.length === 0) return
    const pending = pendingEarlyRecordSteps
    pendingEarlyRecordSteps = []
    for (const item of pending) {
      await applyRecordStepToTarget(item.event, item.targetPath, item.recordSessionId, false)
    }
  }

  async function applyRecordStepToTarget(
    event: RecordStepEvent,
    eventTargetPath = '',
    eventRecordSessionId = '',
    allowEarlyBuffer = true,
  ) {
    if (!shouldApplyRecordStepEventLocal(event, eventRecordSessionId)) {
      if (allowEarlyBuffer && shouldBufferEarlyRecordStepEvent(recordStepApplyGate(eventRecordSessionId), event)) {
        bufferEarlyRecordStep(event, eventTargetPath, eventRecordSessionId)
      }
      return
    }
    await recorderStore.awaitRecordEditorReady()
    recorderStore.chainRecordStepApply(async () => {
      const recorderSnap = get(recorderStore)
      const tabsSnap = tabsStore.snapshot()
      const editorSnap = get(editorStore)
      const targetPath = resolveApplyRecordStepTarget(
        eventTargetPath,
        recorderSnap.targetPath,
        tabsSnap.activeTab,
        WELCOME_KEY,
      )
      if (!targetPath && event.op !== 'reset') return
      const isActiveTarget = isSameRecordTab(tabsSnap.activeTab, targetPath)
      const liveEditorText = monaco?.getEditorText() ?? editorSnap.text
      const sourceText = recordStepSourceText(
        tabsSnap.tabs,
        tabsSnap.activeTab,
        targetPath,
        liveEditorText,
      )
      const tab = tabsSnap.tabs.find((t) => isSameRecordTab(t.path, targetPath))
      if (!tab && !sourceText && event.op !== 'reset') return
      const result = applyRecordStepEvent(sourceText, event, recorderSnap.liveRecordStepLines)
      recorderStore.setLiveRecordStepLines(result.lineByIndex)
      if (tab) {
        tabsStore.mapTabs((tabs) => tabs.map((t) =>
          isSameRecordTab(t.path, targetPath) ? { ...t, draft: result.text, dirty: true } : t,
        ))
      }
      if (isActiveTarget && tabsSnap.activeTab) {
        editorStore.setTextWithBump(result.text)
        const textVersion = get(editorStore).textVersion
        monaco?.activateTab(tabsSnap.activeTab, result.text, textVersion)
        void refreshEditorSteps(result.text, textVersion)
        scheduleValidateEditor(150)
        if (isUntitled(targetPath)) {
          schedulePersistSession()
        }
      }
    })
    await recorderStore.awaitRecordStepApplyChain()
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
    const tab = tabs.find((t) => t.path === bannerPath)
    const current =
      bannerPath === activeTab
        ? (monaco?.getEditorText() ?? editorText)
        : tab
          ? tabEditorText(tab)
          : ''
    if (!postRecordBaselineText || current === postRecordBaselineText) return
    await showPostRecordBanner(bannerPath)
  }

  function handleRecordStopped(payload?: { reason?: string; idleSeconds?: number }) {
    dismissRecorderPicker()
    recorderStore.beginCaptureFinalize()
    void (async () => {
      await recorderStore.awaitRecordEditorReady()
      await recorderStore.awaitRecordStepApplyChain()
      recorderStore.stopCaptureKeepBrowserOpen()
      recorderStore.clearLiveRecordSession()
      syncActiveTabContent()
      await maybeShowPostRecordBannerAfterStop()
      if (payload?.reason === 'idle') {
        const sec = payload.idleSeconds ?? recordIdle ?? 30
        appendLog(tr('journal.record.stoppedIdle', { seconds: sec }))
      }
      if (get(recorderStore).browserOpen) {
        setStatus(payload?.reason === 'idle' ? tr('journal.status.recordStoppedIdle') : tr('journal.status.browserOpen'), 'busy')
        if (payload?.reason !== 'idle') {
          appendLog(tr('journal.record.stoppedBrowserOpen'))
        }
      } else {
        syncIdleStatus()
      }
    })()
  }

  async function handleRecordSessionEnd(result: gui.RunResult, kind: 'record' | 'browse') {
    pendingEarlyRecordSteps = []
    const recordTarget =
      (recordAppendTo || lastRecordTarget || (activeTab && !isWelcome ? activeTab : '')).trim().replace(/\\/g, '/')
    recorderStore.reset()
    dialogsStore.close('showRecord')
    if (result.output) appendLog(result.output)
    if (result.error) {
      appendLog(tr('journal.record.error', { message: result.error || tr('journal.record.unknownError') }))
      if (kind === 'browse') {
        setStatus(tr('journal.status.showBrowserFailed'), 'error')
        openJournalTab()
      }
    } else if (kind === 'record') {
      appendLog(tr('journal.browser.closedSaveHint'))
      if (recordTarget && !isUntitled(recordTarget)) {
        await showPostRecordBanner(recordTarget)
      } else if (activeTab && !isWelcome) {
        await showPostRecordBanner(activeTab)
      }
    }
    journalStore.patch({ statusTone: 'normal' })
    syncIdleStatus()
  }

  async function toggleRecordPause() {
    recorderStore.extendPauseToggleGuard()
    if (recordPaused) {
      await ResumeRecording()
      recorderStore.setRecording(recording, false)
      setStatus(tr('journal.status.recording'), 'busy')
    } else {
      await PauseRecording()
      recorderStore.setRecording(recording, true)
      setStatus(tr('journal.status.paused'), 'busy')
    }
    void syncBrowserStateFromBackend()
  }

  function applyBrowserSessionState(s: { browserOpen: boolean; recording: boolean; paused: boolean }) {
    recorderStore.setBrowserState(s.browserOpen, s.recording, s.paused)
  }

  function handleBrowserLost() {
    if (!browserOpen && !recording) return
    applyBrowserSessionState({ browserOpen: false, recording: false, paused: false })
    journalStore.patch({ statusTone: 'normal' })
    syncIdleStatus()
  }

  function startBrowserWatch() {
    stopBrowserWatch()
    recorderStore.setBrowserWatchTimer(setInterval(() => {
      void syncBrowserStateFromBackend()
    }, 400))
  }

  function stopBrowserWatch() {
    recorderStore.clearBrowserWatchTimer()
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
      runnerStore.setCancelling(true)
      setStatus(tr('journal.status.stoppingTest'), 'busy')
      appendLog(tr('journal.run.stopping'))
      await CancelRun()
      return
    }
    if (get(recorderStore).recording) {
      recorderStore.beginCaptureFinalize()
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
      dialogsStore.close('showOtp')
      return
    }
    appendLog(tr('journal.otp.noActiveRequest'))
  }

  async function cancelOtp() {
    await CancelOTP()
    dialogsStore.close('showOtp')
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
    tabsStore.appendTab({ path, content, dirty: true })
    warnManyOpenTabs()
    tabsStore.patch({ welcomeTabVisible: false, activeTab: path })
    await applyEditorText(content, { switchTab: true, tabPath: path, skipValidate: true })
    trimTabsMemory()
    syncStepsPanelCollapsedFromPrefs()
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

  async function quickStart() {
    if (!(await ensureProjectForBrowser())) return
    recordFormStore.patch({ recordURL: startURL || 'https://example.com' })
    prepareRecordDialogDefaults({ inheritTestClient: true })
    openJournalTab()
    setStatus(tr('journal.browser.launching'), 'busy')
    appendLog(tr('journal.record.started'))
    await startRecord({ headed: false })
  }

  function dismissWelcomeChecklist() {
    uiPrefsStore.patch({ checklistDismissed: true })
    void persistSettings()
  }

  function continueRecord() {
    if (!projectPath) {
      appendLog(tr('journal.project.openFirst'))
      return
    }
    if (!activeTab || isWelcome || isUntitled(activeTab) || !activeTab.toLowerCase().endsWith('.feature')) {
      appendLog(tr('journal.record.openFeatureForAppend'))
      beginRecord()
      return
    }
    recordFormStore.patch({ recordAppendTo: activeTab })
    recordFormStore.patch({ recordOutput: activeTab })
    recordFormStore.patch({ recordFeatureName: basename(activeTab).replace(/\.feature$/i, '') })
    recordFormStore.patch({ recordScenarioName: tr('dialogs.record.appendScenarioDefault') })
    recordFormStore.patch({ recordURL: startURL })
    dialogBinds.syncRecordFormBindLocals()
    dialogBinds.syncRecorderPrefsBindLocals()
    dialogsStore.open('showRecord')
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
    httpAuthDialogStore.setHost(hostFromURL(recordURL || startURL))
    dialogsStore.open('showHttpAuth')
  }

  function closeHttpAuthDialog() {
    dialogsStore.close('showHttpAuth')
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
    pickerDialogStore.setResult(
      result.selector,
      await PickerStepChoices(result.selector, tr('dialogs.record.gherkinGiven')),
      result.candidates || [],
      result.warnings || [],
      result.suggested_action || '',
      result.suggested_choice ?? 0,
    )
    dialogsStore.open('showPickerStep')
  }

  async function refreshPickerChoices(selector: string) {
    const state = pickerDialogStore.snapshot()
    pickerDialogStore.setResult(
      selector,
      await PickerStepChoices(selector, tr('dialogs.record.gherkinGiven')),
      state.candidates,
      state.warnings,
      state.suggestedAction,
      state.suggestedChoice,
    )
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
    appendLog(tr('journal.record.undoDone'))
  }

  function onWelcomeChecklistStep(step: number) {
    if (step === 1) openProjectDialog()
    else if (step === 2) quickStart()
    else if (step === 3) {
      if (tabs.length === 0) newScenario()
      else if (projectPath) executeRun(runFormFromMode(lastRun, 'single', { dryRun: false }))
    }
  }

  function syncUrlFromField() {
    settingsStore.patch({ startUrl: dialogBinds.bindStartURL })
    recordFormStore.patch({ recordURL: dialogBinds.bindStartURL })
  }

  function projectLabel(): string {
    if (isWelcome) return tr('editor.welcome')
    if (activeTab && activeTab !== WELCOME_KEY) return featureTabLabel(activeTab)
    if (projectPath) return basename(projectPath)
    return tr('statusBar.default')
  }

  function paletteActions(): PaletteActions {
    return {
      openCommandPalette: () => dialogsStore.open('showCommandPalette'),
      selectWelcome: () => selectTab(WELCOME_KEY),
      openProject: openProjectDialog,
      openNewProject: openNewProjectWizard,
      closeProject,
      openSettings,
      openInitProject: openInitProjectDialog,
      openExamples,
      newScenario,
      openFile: openFileDialog,
      saveFeature,
      saveFeatureAs,
      openExport: openExportDialog,
      openImport: openImportDialog,
      openImportFeatures: openImportFeaturesDialog,
      openSteps: openStepsDialog,
      openSnippets: openSnippetPalette,
      openFindReplace,
      openFind: () => monaco?.openFind(),
      formatDocument: () => monaco?.formatDocument(),
      openSymbolOutline: () => monaco?.openSymbolOutline(),
      openProjectReplace: () => dialogsStore.open('showProjectReplace'),
      openDuplicate: openDuplicateDialog,
      openRenameFeature: (path) => {
        featureDialogStore.openRename(path)
        dialogsStore.open('showRenameFeature')
      },
      deleteFeature,
      refactorIndents: refactorNormalizeIndents,
      refactorBlanks: refactorCollapseBlank,
      openStepsHelp: () => openStepsHelp(),
      openBrowser,
      beginRecord,
      toggleRecordPause,
      stopRecord,
      openBaselineRecord: openBaselineRecordDialog,
      runPrimary,
      runCurrentScenario,
      openRunDialog,
      runTagDialog: () => openRunDialog(tr('menus.runTag').replace('…', ''), {}, 'tag'),
      openPlaywrightRun: () =>
        openRunDialog('Playwright', { dryRun: false, headed: true, engine: 'playwright', installPW: true }, 'single'),
      toggleBatchMode,
      runBatchSelected,
      rerunFailed,
      openRunHistory,
      openTestClient: openTestClientDialog,
      openTestClientCapture: openTestClientDialogForCapture,
      openValidate: openValidateDialog,
      openVanessa: openVanessaDialog,
      openVanessaSettings: openVanessaSettingsDialog,
      openVanessaMonitor,
      openPlugins: () => dialogsStore.open('showPlugins'),
      openPluginRun,
      openJournal: () => layoutStore.openBottomTab('journal'),
      openResults: () => layoutStore.openBottomTab('results'),
      serveAllure: serveAllureReport,
      openValidatePanel: () => layoutStore.openBottomTab('validate'),
      openErrorPanel: () => layoutStore.openBottomTab('error'),
      showSidebar: () => layoutStore.showSidebar(),
      hideSidebar: () => layoutStore.hideSidebar(),
      togglePreview,
      toggleStepsPanel,
      toggleToolbarCompact,
      refactorUrls: refactorUpdateUrls,
      openHotkeys: () => dialogsStore.open('showHotkeys'),
      resetLayout: resetWindowLayout,
      checkUpdates,
      showAbout: showAboutDialog,
      hasVanessaPlugin,
      pluginLabel,
    }
  }

  $: [, paletteCommands] = [
    $locale,
    buildPaletteCommands(
      tr,
      {
        toolbarCompact,
        previewVisible,
        stepsPanelVisible,
        activeTab,
        isWelcome,
        installedPlugins,
        allureDir: projectArtifacts.allureDir || '',
      },
      paletteActions(),
    ),
  ]
</script>

{#if !appReady}
  <SplashScreen
    version={appVersion}
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
              featureDialogStore.openRename(activeTab)
              dialogsStore.open('showRenameFeature')
            }}
            disabled={isWelcome}
          >
            {tr('menus.rename')}
          </button>
          <button class="menu-item" on:click={() => activeTab && !isWelcome && deleteFeature(activeTab)} disabled={isWelcome}>{tr('menus.delete')}</button>
          <div class="menu-sep"></div>
          <button class="menu-item" on:click={openFindReplace} disabled={isWelcome}>{tr('menus.findReplace')}<span class="menu-shortcut">Ctrl+H</span></button>
          <button class="menu-item" on:click={() => (dialogsStore.open('showProjectReplace'))} disabled={!projectPath}>{tr('menus.projectReplace')}</button>
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
          <button class="menu-item" on:click={() => runMenuAction(() => void openBrowser())}>{tr('menus.browser')}<span class="menu-shortcut">Ctrl+B</span></button>
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
          <button class="menu-item" on:click={() => openRunDialog('', {}, 'single')} disabled={isWelcome && !projectPath}>{tr('menus.runDialog')}</button>
          <button class="menu-item" on:click={() => openRunDialog(tr('menus.runTag').replace('…', ''), {}, 'tag')} disabled={isWelcome && !projectPath}>
            {tr('menus.runTag')}
          </button>
          <button class="menu-item" data-tour="menu-dry-run" on:click={() => runPrimary(true)} disabled={isWelcome && !projectPath && !batchSelected.length}>{tr('menus.dryRun')}</button>
          <button
            class="menu-item"
            on:click={() => openRunDialog('Playwright', { dryRun: false, headed: true, engine: 'playwright', installPW: true }, 'single')}
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
          <button class="menu-item" on:click={() => (dialogsStore.open('showPlugins'))} disabled={!projectPath}>{tr('menus.plugins')}</button>
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
          <button class="menu-item" on:click={() => selectTab(WELCOME_KEY)} disabled={!canShowWelcome(tabs)}>{tr('menus.start')}</button>
          <button class="menu-item" on:click={() => { layoutStore.showSidebar() }}>{tr('menus.scenarios')}</button>
          <button class="menu-item" on:click={() => { layoutStore.hideSidebar() }}>{tr('menus.hideExplorer')}</button>
          <button class="menu-item" on:click={() => (dialogsStore.open('showCommandPalette'))}>{tr('palette.commands.palette')}<span class="menu-shortcut">Ctrl+Shift+P</span></button>
          <button class="menu-item" on:click={() => { layoutStore.openBottomTab('journal') }}>{tr('menus.journal')}</button>
          <button class="menu-item" on:click={() => { layoutStore.openBottomTab('results') }}>{tr('menus.resultsPanel')}</button>
          <button class="menu-item" on:click={openRunHistory} disabled={!projectPath}>{tr('menus.runHistory')}</button>
          <button class="menu-item" on:click={() => { layoutStore.openBottomTab('validate') }}>{tr('menus.validatePanel')}</button>
          <button class="menu-item" on:click={() => { layoutStore.openBottomTab('error') }}>{tr('menus.errorPanel')}</button>
          <button class="menu-item" on:click={togglePreview}>
            {previewVisible ? tr('menus.hidePreview') : tr('menus.showPreview')}
          </button>
          <button class="menu-item" on:click={toggleStepsPanel}>
            {stepsPanelVisible ? tr('menus.hideStepsPanel') : tr('menus.showStepsPanel')}
          </button>
          <button class="menu-item" on:click={toggleToolbarCompact}>
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
          <button class="menu-item" on:click={() => (dialogsStore.open('showHotkeys'))}>{tr('menus.hotkeys')}<span class="menu-shortcut">Shift+F1</span></button>
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
            layoutStore.toggleSidebar()
          }}
        >
          {@html icons.explorer}
        </button>
        <button
          class="activity-btn"
          class:active={bottomPanelOpen}
          title={tr('catalog.outputPanel')}
          on:click={() => {
            layoutStore.toggleBottomPanel()
          }}
        >
          {@html icons.panel}
        </button>
      </aside>

      <!-- Explorer -->
      {#if sidebarVisible}
        <div class="sidebar-column" class:onboarding-elevated={onboardingElevateSidebar} style="width: {layoutSidebarWidth + 4}px">
        <aside class="explorer" style="width: {layoutSidebarWidth}px" on:contextmenu={onExplorerContextMenu}>
          <div class="explorer-header">
            <p class="zone-title">{tr('catalog.title')}</p>
            <div class="explorer-tools">
              <input class="explorer-search" value={sidebarSearch} placeholder={tr('catalog.searchPlaceholder')} on:input={onSidebarSearchInput} />
              <div class="explorer-tool-actions">
                <button type="button" class="icon-btn" title={tr('catalog.refresh')} disabled={!projectPath} on:click={refreshCatalog}>{@html icons.refresh}</button>
                <button type="button" class="icon-btn" title={tr('catalog.newScenario')} on:click={newScenario}>{@html icons.plus}</button>
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
                onDropTarget={(path) => catalogStore.setDropTarget(path)}
              />
            {/if}
          </div>
        </aside>
        <button
          type="button"
          class="splitter-v"
          aria-label="Resize sidebar"
          on:mousedown={startResizeSidebar}
          on:keydown={onSidebarSplitterKeydown}
        ></button>
        </div>
      {/if}

      <!-- Workspace -->
      <section class="workspace">
        <div class="action-bar" class:compact={actionBarCompact} use:observeActionBar>
          <div class="quick-toolbar">
            <div class="toolbar-row primary">
              <button class="tool-btn primary" title={tr('toolbar.browser') + ' (Ctrl+B)'} on:click={() => void openBrowser()}>
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
              <button class="tool-btn" on:click={() => { layoutStore.openBottomTab('journal') }}>
                {@html toolbarIcons.log()}<span>{tr('toolbar.journal')}</span>
              </button>
              <button class="tool-btn" on:click={() => { layoutStore.openBottomTab('results') }}>
                {@html toolbarIcons.results()}<span>{tr('toolbar.results')}</span>
              </button>
            </div>
            {/if}
          </div>

          <div class="url-block">
            <span>URL</span>
            <input bind:value={dialogBinds.bindStartURL} placeholder="https://site.com" on:change={syncUrlFromField} />
            <button class="icon-btn" title={tr('toolbar.urlFromBrowser')} on:click={() => recordFormStore.patch({ recordURL: dialogBinds.bindStartURL })}>
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
              bind:startURL={dialogBinds.bindStartURL}
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
              onOpenRecentFeature={(path) => loadFeature(path, { forceActivate: true })}
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
                onClose={() => (dialogsStore.close('showVanessaMonitor'))}
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
                      valuePath={editorValuePath}
                      valueGeneration={editorTextVersion}
                      readOnly={automationActive || recordingTargetReadOnly}
                      bind:editorSettings={dialogBinds.bindEditorSettings}
                      scenarioHints={editorScenarioHints}
                      hintActions={monacoHintActions}
                      runLensActions={monacoRunLensActions}
                      inlayHintsHandlers={monacoInlayHintsHandlers}
                      on:ready={() => void syncMonacoAfterMount()}
                      on:change={(e) => onEditorChange(e.detail)}
                      on:cursorline={(e) => editorStore.setCursorLine(e.detail)}
                    />
                </div>
                {#if stepsPanelVisible}
                <button
                  type="button"
                  class="splitter-h"
                  aria-label="Resize steps panel"
                  on:mousedown={startResizeSteps}
                  on:keydown={onStepsSplitterKeydown}
                ></button>
                <div class="steps-panel" class:collapsed={stepsPanelCollapsed} style="max-height: {stepsPanelCollapsed ? 24 : stepsPanelHeight}px">
                  <div
                    class="steps-header"
                    title={tr('editor.stepsPanel.tooltip')}
                  >
                    <button
                      type="button"
                      aria-label={stepsPanelCollapsed ? tr('editor.stepsPanel.expand') : tr('editor.stepsPanel.collapse')}
                      on:click={() => uiPrefsStore.toggleStepsPanelCollapsed()}
                    >
                      {#if stepsPanelCollapsed}{@html icons.chevronRight}{:else}{@html icons.chevronDown}{/if}
                    </button>
                    {#if dialogBinds.bindEditorSettings.symbolOutline}
                      <div class="steps-panel-tabs" role="tablist" aria-label={tr('editor.stepsPanel.aria')}>
                        <button
                          type="button"
                          role="tab"
                          class:active={stepsPanelTab === 'outline'}
                          aria-selected={stepsPanelTab === 'outline'}
                          on:click={() => editorStore.setStepsPanelTab('outline')}
                        >
                          {tr('editor.stepsPanel.outline')}
                        </button>
                        <button
                          type="button"
                          role="tab"
                          class:active={stepsPanelTab === 'steps'}
                          aria-selected={stepsPanelTab === 'steps'}
                          on:click={() => editorStore.setStepsPanelTab('steps')}
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
                    {#if dialogBinds.bindEditorSettings.symbolOutline && stepsPanelTab === 'outline'}
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
                <button
                  type="button"
                  class="splitter-v"
                  aria-label="Resize preview"
                  on:mousedown={startResizePreview}
                  on:keydown={onPreviewSplitterKeydown}
                ></button>
                <div class="feature-preview-pane" style="width: {layoutPreviewWidth}px">
                  <div class="preview-header">{tr('editor.preview')}</div>
                  <FeaturePreview
                    text={editorText}
                    theme={dialogBinds.bindEditorSettings.theme}
                    fontSize={dialogBinds.bindEditorSettings.fontSize}
                    fontFamily={dialogBinds.bindEditorSettings.fontFamily}
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
    <button
      type="button"
      class="splitter-h bottom-splitter"
      aria-label="Resize bottom panel"
      on:mousedown={startResizeBottom}
      on:keydown={onBottomSplitterKeydown}
    ></button>
    <div class="bottom-panel" style="--panel-height: {bottomPanelHeight}px">
    <div class="panel-tabs">
      <button class="panel-tab" class:active={bottomTab === 'journal'} data-tour="panel-journal" on:click={() => openJournalTab(true)}>{tr('panels.journal')}</button>
      <button class="panel-tab" class:active={bottomTab === 'results'} on:click={() => (layoutStore.setBottomTab('results'))}>{tr('panels.results')}</button>
      <button class="panel-tab" class:active={bottomTab === 'validate'} on:click={() => (layoutStore.setBottomTab('validate'))}>{tr('panels.validate')}</button>
      <button class="panel-tab" class:active={bottomTab === 'error'} on:click={() => (layoutStore.setBottomTab('error'))}>{tr('panels.error')}</button>
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
        <div class="status-segment warning" title={tr('statusBar.dryRunWarning')}>
          {tr('statusBar.dryRunWarning')}
        </div>
      {/if}
      {#if showLargeFileBanner}
        <div class="status-segment warning large-file-banner" title={tr('statusBar.largeFileTitle')}>
          {tr('statusBar.largeFile', { lines: LARGE_FILE_LINE_THRESHOLD })}
        </div>
      {/if}
      <div class="status-segment muted">{tr('statusBar.runner')}</div>
      <div class="status-segment" class:warning={stepStatusError}>{stepStatusDisplay}</div>
      <button type="button" class="status-segment clickable" on:click={() => { layoutStore.openBottomTab('journal') }}>
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
    bind:form={dialogBinds.bindRunForm}
    {testClients}
    {tags}
    scenarios={runDialogScenarios}
    onConfirm={confirmRun}
    onCancel={() => (dialogsStore.close('showRun'))}
  />
{/if}

{#if showVanessaRun}
  <VanessaRunDialog
    dryRun={vanessaDry}
    preferRerun={vanessaPreferRerun}
    bind:tag={dialogBinds.bindVanessaTag}
    bind:excludeTags={dialogBinds.bindVanessaExcludeTags}
    bind:scenario={dialogBinds.bindVanessaScenario}
    bind:rerunFailedRunDir={dialogBinds.bindVanessaRerunDir}
    bind:installEpf={dialogBinds.bindVanessaInstallEpf}
    bind:epfUrl={dialogBinds.bindVanessaEpfUrl}
    bind:epfDest={dialogBinds.bindVanessaEpfDest}
    bind:platformExe={dialogBinds.bindVanessaPlatformExe}
    bind:epfPath={dialogBinds.bindVanessaEpfPath}
    bind:ibConnection={dialogBinds.bindVanessaIB}
    bind:reportAllure={dialogBinds.bindVanessaReportAllure}
    bind:vaDir={dialogBinds.bindVanessaVaDir}
    bind:vaFiles={dialogBinds.bindVanessaVaFiles}
    {tags}
    scenarios={vanessaDialogScenarios}
    onConfirm={confirmVanessaRun}
    onCancel={() => (dialogsStore.close('showVanessaRun'))}
  />
{/if}

{#if showTestClient}
  <TestClientDialog
    {testClients}
    bind:selectedName={dialogBinds.bindTestClientSelection}
    browserOpen={browserOpen || recording}
    suggestName={testClientSuggestName}
    onUse={useTestClient}
    onClose={() => {
      dialogsStore.close('showTestClient')
      testClientStore.clearSuggestName()
    }}
    onClientsChange={(names) => testClientStore.setClients(names)}
    onLog={appendLog}
    onAskConfirm={(message) =>
      askConfirm({ title: tr('confirm.generic.title'), message, confirmLabel: tr('confirm.generic.confirmLabelDelete'), danger: true })}
  />
{/if}

{#if showStepsHelp}
  <StepsHelpDialog
    initialQuery={stepsHelpQuery}
    onClose={() => { dialogsStore.close('showStepsHelp'); stepsHelpDialogStore.clear() }}
    onInsert={insertStep}
  />
{/if}

{#if showSteps}
  <StepsInsertDialog onInsert={insertStep} onClose={() => (dialogsStore.close('showSteps'))} />
{/if}

{#if showVanessaSettings}
  <VanessaSettingsDialog onClose={() => (dialogsStore.close('showVanessaSettings'))} onLog={appendLog} />
{/if}

{#if showExport}
  <ExportDialog
    inputPath={exportInputPath}
    featureText={editorText}
    {currentProjectVersion}
    onClose={() => (dialogsStore.close('showExport'))}
    onLog={appendLog}
  />
{/if}

{#if showRefactorUrl}
  <RefactorUrlDialog
    initialUrl={startURL || recordURL || 'https://example.com'}
    onConfirm={applyRefactorUrl}
    onClose={() => (dialogsStore.close('showRefactorUrl'))}
  />
{/if}

{#if showOpenProject}
  <OpenProjectDialog
    initialPath={projectPath}
    {recentProjects}
    onConfirm={openProjectAt}
    onClose={() => (dialogsStore.close('showOpenProject'))}
  />
{/if}

{#if showRenameFeature}
  <RenameFeatureDialog
    currentPath={renameFeaturePath}
    onConfirm={(name) => renameFeature(renameFeaturePath, name)}
    onClose={() => {
      dialogsStore.close('showRenameFeature')
      featureDialogStore.clearRename()
    }}
  />
{/if}

{#if showMoveFeature}
  <MoveFeatureDialog
    featurePath={moveFeaturePath}
    destDirs={moveDestDirs}
    bind:destDir={dialogBinds.bindMoveDestDir}
    onConfirm={confirmMoveFeature}
    onCancel={() => {
      dialogsStore.close('showMoveFeature')
      featureDialogStore.clearMove()
    }}
  />
{/if}

{#if showValidate}
  <ValidateDialog
    bind:browser={dialogBinds.bindValidateBrowser}
    bind:syntaxOnly={dialogBinds.bindValidateSyntaxOnly}
    bind:flowAware={dialogBinds.bindValidateFlowAware}
    bind:scope={dialogBinds.bindValidateScope}
    canValidateCurrent={!isWelcome && !!activeTab}
    currentFileName={!isWelcome && activeTab ? basename(activeTab) : ''}
    onConfirm={confirmValidate}
    onCancel={() => (dialogsStore.close('showValidate'))}
  />
{/if}

{#if showNewProjectWizard}
  <NewProjectWizardDialog
    defaultStartUrl={startURL || 'https://example.com'}
    onConfirm={confirmNewProjectWizard}
    onCancel={() => (dialogsStore.close('showNewProjectWizard'))}
  />
{/if}

{#if showInitProject}
  <InitProjectDialog
    {projectPath}
    onConfirm={confirmInitProject}
    onCancel={() => (dialogsStore.close('showInitProject'))}
  />
{/if}

{#if showUpdateCheck}
  <UpdateCheckDialog
    currentVersion={appVersion}
    info={updateCheckInfo}
    message={updateCheckMessage}
    hasUpdate={updateCheckHasUpdate}
    downloading={updateDownloading}
    progress={updateProgress}
    onClose={() => (dialogsStore.close('showUpdateCheck'))}
    onOpenRelease={openUpdateRelease}
    onDownload={downloadUpdate}
    onApply={applyUpdate}
    canAutoApply={updateCheckInfo?.canAutoApply ?? false}
  />
{/if}

{#if showDuplicateFeature}
  <DuplicateFeatureDialog
    featurePath={duplicateFeaturePath}
    bind:newName={dialogBinds.bindDuplicateNewName}
    onConfirm={confirmDuplicateFeature}
    onCancel={() => {
      dialogsStore.close('showDuplicateFeature')
      featureDialogStore.clearDuplicate()
    }}
  />
{/if}

{#if showImportFeatures}
  <ImportFeaturesDialog
    destDirs={collectProjectDirs()}
    bind:destDir={dialogBinds.bindImportDestDir}
    busy={importFeaturesBusy}
    onImport={confirmImportFeatures}
    onClose={() => (dialogsStore.close('showImportFeatures'))}
  />
{/if}

{#if showImport}
  <ImportJSONDialog
    {projectPath}
    {currentProjectVersion}
    onClose={() => (dialogsStore.close('showImport'))}
    onLog={appendLog}
    onImported={onImportComplete}
  />
{/if}

{#if showSettings}
  <SettingsDialog
    bind:browser={dialogBinds.bindSettingsBrowser}
    bind:headless={dialogBinds.bindSettingsHeadless}
    bind:workers={dialogBinds.bindSettingsWorkers}
    bind:slowMo={dialogBinds.bindSettingsSlowMo}
    bind:loops={dialogBinds.bindSettingsLoops}
    bind:filterRecording={dialogBinds.bindFilterRecording}
    bind:navOnlyRecording={dialogBinds.bindNavOnlyRecording}
    bind:hoverRecord={dialogBinds.bindHoverRecord}
    bind:scrollBeforeClick={dialogBinds.bindSettingsScrollBeforeClick}
    bind:disableRecordUrlWait={dialogBinds.bindSettingsDisableRecordUrlWait}
    bind:hoverRecordMinMs={dialogBinds.bindSettingsHoverRecordMinMs}
    bind:toolbarCompact={dialogBinds.bindToolbarCompact}
    bind:stepsPanelVisible={dialogBinds.bindStepsPanelVisible}
    bind:stepsPanelHeight={dialogBinds.bindStepsPanelHeight}
    bind:checkUpdatesOnStartup={dialogBinds.bindSettingsCheckUpdatesOnStartup}
    bind:selectorClickStrategies={dialogBinds.bindSettingsSelectorClickStrategies}
    bind:selectorInputStrategies={dialogBinds.bindSettingsSelectorInputStrategies}
    bind:libraryHeuristicsMui={dialogBinds.bindSettingsLibraryHeuristicsMui}
    bind:libraryHeuristicsAnt={dialogBinds.bindSettingsLibraryHeuristicsAnt}
    bind:navWaitUntil={dialogBinds.bindSettingsNavWaitUntil}
    bind:htmlReportOpenMode={dialogBinds.bindHtmlReportOpenMode}
    projectOpen={!!projectPath}
    bind:pickerDuringRecording={dialogBinds.bindPickerDuringRecording}
    bind:uiLocale={dialogBinds.bindUiLocale}
    bind:editorSettings={dialogBinds.bindEditorSettings}
    onSave={applySettings}
    onApply={applySettingsKeepOpen}
    onCancel={cancelSettings}
    onOpenPlugins={() => {
      dialogsStore.close('showSettings')
      dialogsStore.open('showPlugins')
    }}
    onOpenVanessa={() => {
      dialogsStore.close('showSettings')
      openVanessaSettingsDialog()
    }}
    onInstallLog={appendLog}
  />
{/if}

{#if showCommandPalette}
  <CommandPalette commands={paletteCommands} onClose={() => (dialogsStore.close('showCommandPalette'))} />
{/if}

{#if showSnippetPalette}
  <SnippetPalette onClose={() => (dialogsStore.close('showSnippetPalette'))} onInsert={insertStep} />
{/if}

{#if showRecord}
  <RecordDialog
    bind:mode={dialogBinds.bindRecordMode}
    bind:stepPickerOpen={dialogBinds.bindRecordStepPickerOpen}
    bind:url={dialogBinds.bindRecordURL}
    bind:output={dialogBinds.bindRecordOutput}
    bind:featureName={dialogBinds.bindRecordFeatureName}
    bind:scenarioName={dialogBinds.bindRecordScenarioName}
    bind:testClient={dialogBinds.bindRecordTestClient}
    bind:idleSeconds={dialogBinds.bindRecordIdle}
    bind:appendTo={dialogBinds.bindRecordAppendTo}
    bind:headless={dialogBinds.bindSettingsHeadless}
    bind:filterRecording={dialogBinds.bindFilterRecording}
    bind:navOnlyRecording={dialogBinds.bindNavOnlyRecording}
    bind:hoverRecord={dialogBinds.bindHoverRecord}
    {testClients}
    {recording}
    {recordPaused}
    {baselineBusy}
    onHttpAuth={openHttpAuthDialog}
    onStart={startRecord}
    onTogglePause={toggleRecordPause}
    onStop={stopRecord}
    onSaveBaseline={saveBaselineRecord}
    onClose={() => {
      dialogBinds.flushRecordFormBindLocals()
      dialogsStore.close('showRecord')
    }}
    childModalOpen={showHttpAuth}
  />
{/if}

{#if showOtp}
  <OtpDialog email={otpEmail} onSubmit={submitOtp} onCancel={cancelOtp} />
{/if}

{#if showAbout}
  <AboutDialog version={appVersion} onClose={() => (dialogsStore.close('showAbout'))} />
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
  <HotkeysDialog commands={paletteCommands} onClose={() => (dialogsStore.close('showHotkeys'))} />
{/if}

{#if showPlugins}
  <PluginsDialog
    childModalOpen={showPluginRun}
    onClose={() => {
      dialogsStore.close('showPlugins')
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
    bind:tag={dialogBinds.bindPluginRunTag}
    bind:scenario={dialogBinds.bindPluginRunScenario}
    bind:dryRun={dialogBinds.bindPluginRunDry}
    scenarios={pluginRunScenarios}
    {tags}
    onConfirm={confirmPluginRun}
    onCancel={() => (dialogsStore.close('showPluginRun'))}
  />
{/if}

{#if showRunHistory}
  <RunHistoryDialog
    entries={$reportsStore.runResults}
    flakyByPath={flakyByPath}
    flakyStepByPath={flakyStepByPath}
    onOpenFeature={openFeatureFromHistory}
    onRerunFailed={() => { dialogsStore.close('showRunHistory'); rerunFailed() }}
    onClose={() => (dialogsStore.close('showRunHistory'))}
  />
{/if}

{#if showPostRecordDiff && postRecordPath}
  <PostRecordDiffDialog
    path={postRecordPath}
    original={postRecordBaselineText}
    modified={monaco?.getEditorText() ?? editorText}
    onClose={() => (dialogsStore.close('showPostRecordDiff'))}
  />
{/if}

{#if showProjectReplace}
  <ProjectReplaceDialog
    bind:findText={dialogBinds.bindProjectReplaceFind}
    bind:replaceText={dialogBinds.bindProjectReplaceReplace}
    bind:caseSensitive={dialogBinds.bindProjectReplaceCaseSensitive}
    busy={projectReplaceBusy}
    onConfirm={confirmProjectReplace}
    onClose={() => (dialogsStore.close('showProjectReplace'))}
  />
{/if}

{#if showHttpAuth}
  <HttpAuthDialog initialHost={httpAuthHost} onCancel={closeHttpAuthDialog} />
{/if}

{#if showPickerStep}
  <PickerStepDialog
    selector={pickerSelector}
    choices={pickerChoices}
    candidates={pickerCandidates}
    warnings={pickerWarnings}
    suggestedAction={pickerSuggestedAction}
    suggestedChoice={pickerSuggestedChoice}
    onInsert={insertPickerStep}
    onSelectorChange={refreshPickerChoices}
    onClose={() => (dialogsStore.close('showPickerStep'))}
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
    onRefresh={folderMenuRefresh}
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

{#if confirmDialogOpen && confirmDialog}
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
