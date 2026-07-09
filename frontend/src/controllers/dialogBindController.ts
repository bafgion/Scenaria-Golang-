import type { Locale } from '../lib/i18n'
import { DEFAULT_EDITOR_SETTINGS, type EditorSettings } from '../lib/editorOptions'
import { defaultRunForm, runFormFromMode, type RunForm, type RunFormMode } from '../lib/runTypes'
import type { createFeatureDialogStore } from '../stores/featureDialogStore'
import type { createPluginRunStore } from '../stores/pluginRunStore'
import type { createRecordFormStore } from '../stores/recordFormStore'
import type { createRecorderPrefsStore } from '../stores/recorderPrefsStore'
import type { createRunFormStore } from '../stores/runFormStore'
import type { createSettingsDialogStore } from '../stores/settingsDialogStore'
import type { createSettingsStore } from '../stores/settingsStore'
import type { createTestClientStore } from '../stores/testClientStore'
import type { createUiPrefsStore } from '../stores/uiPrefsStore'
import type { createValidateDialogStore } from '../stores/validateDialogStore'
import type { createVanessaRunStore } from '../stores/vanessaRunStore'
import type { createProjectReplaceStore } from '../stores/projectReplaceStore'

export type DialogBindStores = {
  settingsStore: ReturnType<typeof createSettingsStore>
  uiPrefsStore: ReturnType<typeof createUiPrefsStore>
  recorderPrefsStore: ReturnType<typeof createRecorderPrefsStore>
  recordFormStore: ReturnType<typeof createRecordFormStore>
  validateDialogStore: ReturnType<typeof createValidateDialogStore>
  featureDialogStore: ReturnType<typeof createFeatureDialogStore>
  projectReplaceStore: ReturnType<typeof createProjectReplaceStore>
  testClientStore: ReturnType<typeof createTestClientStore>
  vanessaRunStore: ReturnType<typeof createVanessaRunStore>
  pluginRunStore: ReturnType<typeof createPluginRunStore>
  runFormStore: ReturnType<typeof createRunFormStore>
  settingsDialogStore: ReturnType<typeof createSettingsDialogStore>
}

export function createDialogBindController(stores: DialogBindStores) {
  const bindEditorSettings: EditorSettings = { ...DEFAULT_EDITOR_SETTINGS }

  let bindValidateBrowser = 'chromium'
  let bindValidateSyntaxOnly = false
  let bindValidateScope: 'project' | 'current' = 'project'
  let bindDuplicateNewName = ''
  let bindMoveDestDir = ''
  let bindImportDestDir = ''
  let bindTestClientSelection = ''
  let bindProjectReplaceFind = ''
  let bindProjectReplaceReplace = ''
  let bindProjectReplaceCaseSensitive = false

  let bindToolbarCompact = false
  let bindStepsPanelVisible = true
  let bindStepsPanelHeight = 160
  let bindPickerDuringRecording = false
  let bindFilterRecording = false
  let bindNavOnlyRecording = false
  let bindHoverRecord = false
  let bindRecordURL = ''
  let bindRecordOutput = 'recorded.feature'
  let bindRecordIdle = 30
  let bindRecordAppendTo = ''
  let bindRecordTestClient = ''
  let bindRecordFeatureName = ''
  let bindRecordScenarioName = ''
  let bindRecordMode: 'live' | 'baseline' = 'live'
  let bindRecordStepPickerOpen = false

  let bindSettingsBrowser = 'chromium'
  let bindSettingsHeadless = false
  let bindSettingsWorkers = 1
  let bindSettingsSlowMo = 0
  let bindSettingsLoops = 100
  let bindSettingsScrollBeforeClick = false
  let bindSettingsHoverRecordMinMs = 600
  let bindSettingsCheckUpdatesOnStartup = true
  let bindSettingsSelectorClickStrategies: string[] = []
  let bindSettingsSelectorInputStrategies: string[] = []
  let bindSettingsNavWaitUntil = 'domcontentloaded'
  let bindStartURL = ''
  let bindUiLocale: Locale = 'ru'
  let bindHtmlReportOpenMode: 'full' | 'light' = 'full'
  let bindRunForm: RunForm = defaultRunForm()

  let bindVanessaTag = ''
  let bindVanessaExcludeTags = ''
  let bindVanessaScenario = ''
  let bindVanessaRerunDir = ''
  let bindVanessaInstallEpf = false
  let bindVanessaEpfUrl = ''
  let bindVanessaEpfDest = ''
  let bindVanessaPlatformExe = ''
  let bindVanessaEpfPath = ''
  let bindVanessaIB = ''
  let bindVanessaReportAllure = false
  let bindVanessaVaDir = ''
  let bindVanessaVaFiles = ''
  let bindPluginRunTag = ''
  let bindPluginRunScenario = ''
  let bindPluginRunDry = false

  return {
    get bindEditorSettings() {
      return bindEditorSettings
    },
    set bindEditorSettings(next: EditorSettings) {
      Object.assign(bindEditorSettings, next)
    },
    get bindValidateBrowser() {
      return bindValidateBrowser
    },
    set bindValidateBrowser(value: string) {
      bindValidateBrowser = value
    },
    get bindValidateSyntaxOnly() {
      return bindValidateSyntaxOnly
    },
    set bindValidateSyntaxOnly(value: boolean) {
      bindValidateSyntaxOnly = value
    },
    get bindValidateScope() {
      return bindValidateScope
    },
    set bindValidateScope(value: 'project' | 'current') {
      bindValidateScope = value
    },
    get bindDuplicateNewName() {
      return bindDuplicateNewName
    },
    set bindDuplicateNewName(value: string) {
      bindDuplicateNewName = value
    },
    get bindMoveDestDir() {
      return bindMoveDestDir
    },
    set bindMoveDestDir(value: string) {
      bindMoveDestDir = value
    },
    get bindImportDestDir() {
      return bindImportDestDir
    },
    set bindImportDestDir(value: string) {
      bindImportDestDir = value
    },
    get bindTestClientSelection() {
      return bindTestClientSelection
    },
    set bindTestClientSelection(value: string) {
      bindTestClientSelection = value
    },
    get bindProjectReplaceFind() {
      return bindProjectReplaceFind
    },
    set bindProjectReplaceFind(value: string) {
      bindProjectReplaceFind = value
    },
    get bindProjectReplaceReplace() {
      return bindProjectReplaceReplace
    },
    set bindProjectReplaceReplace(value: string) {
      bindProjectReplaceReplace = value
    },
    get bindProjectReplaceCaseSensitive() {
      return bindProjectReplaceCaseSensitive
    },
    set bindProjectReplaceCaseSensitive(value: boolean) {
      bindProjectReplaceCaseSensitive = value
    },
    get bindToolbarCompact() {
      return bindToolbarCompact
    },
    set bindToolbarCompact(value: boolean) {
      bindToolbarCompact = value
    },
    get bindStepsPanelVisible() {
      return bindStepsPanelVisible
    },
    set bindStepsPanelVisible(value: boolean) {
      bindStepsPanelVisible = value
    },
    get bindStepsPanelHeight() {
      return bindStepsPanelHeight
    },
    set bindStepsPanelHeight(value: number) {
      bindStepsPanelHeight = value
    },
    get bindPickerDuringRecording() {
      return bindPickerDuringRecording
    },
    set bindPickerDuringRecording(value: boolean) {
      bindPickerDuringRecording = value
    },
    get bindFilterRecording() {
      return bindFilterRecording
    },
    set bindFilterRecording(value: boolean) {
      bindFilterRecording = value
    },
    get bindNavOnlyRecording() {
      return bindNavOnlyRecording
    },
    set bindNavOnlyRecording(value: boolean) {
      bindNavOnlyRecording = value
    },
    get bindHoverRecord() {
      return bindHoverRecord
    },
    set bindHoverRecord(value: boolean) {
      bindHoverRecord = value
    },
    get bindRecordURL() {
      return bindRecordURL
    },
    set bindRecordURL(value: string) {
      bindRecordURL = value
    },
    get bindRecordOutput() {
      return bindRecordOutput
    },
    set bindRecordOutput(value: string) {
      bindRecordOutput = value
    },
    get bindRecordIdle() {
      return bindRecordIdle
    },
    set bindRecordIdle(value: number) {
      bindRecordIdle = value
    },
    get bindRecordAppendTo() {
      return bindRecordAppendTo
    },
    set bindRecordAppendTo(value: string) {
      bindRecordAppendTo = value
    },
    get bindRecordTestClient() {
      return bindRecordTestClient
    },
    set bindRecordTestClient(value: string) {
      bindRecordTestClient = value
    },
    get bindRecordFeatureName() {
      return bindRecordFeatureName
    },
    set bindRecordFeatureName(value: string) {
      bindRecordFeatureName = value
    },
    get bindRecordScenarioName() {
      return bindRecordScenarioName
    },
    set bindRecordScenarioName(value: string) {
      bindRecordScenarioName = value
    },
    get bindRecordMode() {
      return bindRecordMode
    },
    set bindRecordMode(value: 'live' | 'baseline') {
      bindRecordMode = value
    },
    get bindRecordStepPickerOpen() {
      return bindRecordStepPickerOpen
    },
    set bindRecordStepPickerOpen(value: boolean) {
      bindRecordStepPickerOpen = value
    },
    get bindSettingsBrowser() {
      return bindSettingsBrowser
    },
    set bindSettingsBrowser(value: string) {
      bindSettingsBrowser = value
    },
    get bindSettingsHeadless() {
      return bindSettingsHeadless
    },
    set bindSettingsHeadless(value: boolean) {
      bindSettingsHeadless = value
    },
    get bindSettingsWorkers() {
      return bindSettingsWorkers
    },
    set bindSettingsWorkers(value: number) {
      bindSettingsWorkers = value
    },
    get bindSettingsSlowMo() {
      return bindSettingsSlowMo
    },
    set bindSettingsSlowMo(value: number) {
      bindSettingsSlowMo = value
    },
    get bindSettingsLoops() {
      return bindSettingsLoops
    },
    set bindSettingsLoops(value: number) {
      bindSettingsLoops = value
    },
    get bindSettingsScrollBeforeClick() {
      return bindSettingsScrollBeforeClick
    },
    set bindSettingsScrollBeforeClick(value: boolean) {
      bindSettingsScrollBeforeClick = value
    },
    get bindSettingsHoverRecordMinMs() {
      return bindSettingsHoverRecordMinMs
    },
    set bindSettingsHoverRecordMinMs(value: number) {
      bindSettingsHoverRecordMinMs = value
    },
    get bindSettingsCheckUpdatesOnStartup() {
      return bindSettingsCheckUpdatesOnStartup
    },
    set bindSettingsCheckUpdatesOnStartup(value: boolean) {
      bindSettingsCheckUpdatesOnStartup = value
    },
    get bindSettingsSelectorClickStrategies() {
      return bindSettingsSelectorClickStrategies
    },
    set bindSettingsSelectorClickStrategies(value: string[]) {
      bindSettingsSelectorClickStrategies = value
    },
    get bindSettingsSelectorInputStrategies() {
      return bindSettingsSelectorInputStrategies
    },
    set bindSettingsSelectorInputStrategies(value: string[]) {
      bindSettingsSelectorInputStrategies = value
    },
    get bindSettingsNavWaitUntil() {
      return bindSettingsNavWaitUntil
    },
    set bindSettingsNavWaitUntil(value: string) {
      bindSettingsNavWaitUntil = value
    },
    get bindStartURL() {
      return bindStartURL
    },
    set bindStartURL(value: string) {
      bindStartURL = value
    },
    get bindUiLocale() {
      return bindUiLocale
    },
    set bindUiLocale(value: Locale) {
      bindUiLocale = value
    },
    get bindHtmlReportOpenMode() {
      return bindHtmlReportOpenMode
    },
    set bindHtmlReportOpenMode(value: 'full' | 'light') {
      bindHtmlReportOpenMode = value
    },
    get bindRunForm() {
      return bindRunForm
    },
    set bindRunForm(value: RunForm) {
      bindRunForm = value
    },
    get bindVanessaTag() {
      return bindVanessaTag
    },
    set bindVanessaTag(value: string) {
      bindVanessaTag = value
    },
    get bindVanessaExcludeTags() {
      return bindVanessaExcludeTags
    },
    set bindVanessaExcludeTags(value: string) {
      bindVanessaExcludeTags = value
    },
    get bindVanessaScenario() {
      return bindVanessaScenario
    },
    set bindVanessaScenario(value: string) {
      bindVanessaScenario = value
    },
    get bindVanessaRerunDir() {
      return bindVanessaRerunDir
    },
    set bindVanessaRerunDir(value: string) {
      bindVanessaRerunDir = value
    },
    get bindVanessaInstallEpf() {
      return bindVanessaInstallEpf
    },
    set bindVanessaInstallEpf(value: boolean) {
      bindVanessaInstallEpf = value
    },
    get bindVanessaEpfUrl() {
      return bindVanessaEpfUrl
    },
    set bindVanessaEpfUrl(value: string) {
      bindVanessaEpfUrl = value
    },
    get bindVanessaEpfDest() {
      return bindVanessaEpfDest
    },
    set bindVanessaEpfDest(value: string) {
      bindVanessaEpfDest = value
    },
    get bindVanessaPlatformExe() {
      return bindVanessaPlatformExe
    },
    set bindVanessaPlatformExe(value: string) {
      bindVanessaPlatformExe = value
    },
    get bindVanessaEpfPath() {
      return bindVanessaEpfPath
    },
    set bindVanessaEpfPath(value: string) {
      bindVanessaEpfPath = value
    },
    get bindVanessaIB() {
      return bindVanessaIB
    },
    set bindVanessaIB(value: string) {
      bindVanessaIB = value
    },
    get bindVanessaReportAllure() {
      return bindVanessaReportAllure
    },
    set bindVanessaReportAllure(value: boolean) {
      bindVanessaReportAllure = value
    },
    get bindVanessaVaDir() {
      return bindVanessaVaDir
    },
    set bindVanessaVaDir(value: string) {
      bindVanessaVaDir = value
    },
    get bindVanessaVaFiles() {
      return bindVanessaVaFiles
    },
    set bindVanessaVaFiles(value: string) {
      bindVanessaVaFiles = value
    },
    get bindPluginRunTag() {
      return bindPluginRunTag
    },
    set bindPluginRunTag(value: string) {
      bindPluginRunTag = value
    },
    get bindPluginRunScenario() {
      return bindPluginRunScenario
    },
    set bindPluginRunScenario(value: string) {
      bindPluginRunScenario = value
    },
    get bindPluginRunDry() {
      return bindPluginRunDry
    },
    set bindPluginRunDry(value: boolean) {
      bindPluginRunDry = value
    },

    syncVanessaBindLocals() {
      const v = stores.vanessaRunStore.snapshot()
      bindVanessaTag = v.tag
      bindVanessaExcludeTags = v.excludeTags
      bindVanessaScenario = v.scenario
      bindVanessaRerunDir = v.rerunDir
      bindVanessaInstallEpf = v.installEpf
      bindVanessaEpfUrl = v.epfUrl
      bindVanessaEpfDest = v.epfDest
      bindVanessaPlatformExe = v.platformExe
      bindVanessaEpfPath = v.epfPath
      bindVanessaIB = v.ib
      bindVanessaReportAllure = v.reportAllure
      bindVanessaVaDir = v.vaDir
      bindVanessaVaFiles = v.vaFiles
    },
    flushVanessaBindLocals() {
      stores.vanessaRunStore.patch({
        tag: bindVanessaTag,
        excludeTags: bindVanessaExcludeTags,
        scenario: bindVanessaScenario,
        rerunDir: bindVanessaRerunDir,
        installEpf: bindVanessaInstallEpf,
        epfUrl: bindVanessaEpfUrl,
        epfDest: bindVanessaEpfDest,
        platformExe: bindVanessaPlatformExe,
        epfPath: bindVanessaEpfPath,
        ib: bindVanessaIB,
        reportAllure: bindVanessaReportAllure,
        vaDir: bindVanessaVaDir,
        vaFiles: bindVanessaVaFiles,
      })
    },
    syncPluginRunBindLocals() {
      const p = stores.pluginRunStore.snapshot()
      bindPluginRunTag = p.tag
      bindPluginRunScenario = p.scenario
      bindPluginRunDry = p.dry
    },
    flushPluginRunBindLocals() {
      stores.pluginRunStore.patch({
        tag: bindPluginRunTag,
        scenario: bindPluginRunScenario,
        dry: bindPluginRunDry,
      })
    },
    syncValidateDialogBindLocals() {
      const s = stores.validateDialogStore.snapshot()
      bindValidateBrowser = s.browser
      bindValidateSyntaxOnly = s.syntaxOnly
      bindValidateScope = s.scope
    },
    syncEditorSettingsBindLocal() {
      Object.assign(bindEditorSettings, stores.settingsStore.snapshot().editor)
    },
    syncSettingsDialogBindLocals() {
      bindHtmlReportOpenMode = stores.settingsDialogStore.snapshot().htmlReportOpenMode
    },
    flushSettingsDialogBindLocals() {
      stores.settingsDialogStore.setHtmlReportOpenMode(bindHtmlReportOpenMode)
    },
    syncSettingsBindLocals() {
      const s = stores.settingsStore.snapshot()
      bindSettingsBrowser = s.browser
      bindSettingsHeadless = s.headless
      bindSettingsWorkers = s.parallelWorkers
      bindSettingsSlowMo = s.slowMo
      bindSettingsLoops = s.maxLoopIterations
      bindSettingsScrollBeforeClick = s.scrollBeforeClick
      bindSettingsHoverRecordMinMs = s.hoverRecordMinMs
      bindSettingsCheckUpdatesOnStartup = s.checkUpdatesOnStartup
      bindSettingsSelectorClickStrategies = [...s.selectorClickStrategies]
      bindSettingsSelectorInputStrategies = [...s.selectorInputStrategies]
      bindSettingsNavWaitUntil = s.navWaitUntil
      bindStartURL = s.startUrl
      bindUiLocale = s.uiLocale
    },
    flushSettingsBindLocals() {
      stores.settingsStore.patch({
        browser: bindSettingsBrowser,
        headless: bindSettingsHeadless,
        parallelWorkers: bindSettingsWorkers,
        slowMo: bindSettingsSlowMo,
        maxLoopIterations: bindSettingsLoops,
        scrollBeforeClick: bindSettingsScrollBeforeClick,
        hoverRecordMinMs: bindSettingsHoverRecordMinMs,
        checkUpdatesOnStartup: bindSettingsCheckUpdatesOnStartup,
        selectorClickStrategies: bindSettingsSelectorClickStrategies,
        selectorInputStrategies: bindSettingsSelectorInputStrategies,
        navWaitUntil: bindSettingsNavWaitUntil,
        startUrl: bindStartURL,
        uiLocale: bindUiLocale,
      })
      stores.settingsStore.setEditor({ ...bindEditorSettings })
    },
    syncUiPrefsBindLocals() {
      const p = stores.uiPrefsStore.snapshot()
      bindToolbarCompact = p.toolbarCompact
      bindStepsPanelVisible = p.stepsPanelVisible
      bindStepsPanelHeight = p.stepsPanelHeight
      bindPickerDuringRecording = p.pickerDuringRecording
    },
    flushUiPrefsBindLocals() {
      stores.uiPrefsStore.patch({
        toolbarCompact: bindToolbarCompact,
        stepsPanelVisible: bindStepsPanelVisible,
        stepsPanelHeight: bindStepsPanelHeight,
        pickerDuringRecording: bindPickerDuringRecording,
      })
    },
    syncRecorderPrefsBindLocals() {
      const p = stores.recorderPrefsStore.snapshot()
      bindFilterRecording = p.filterRecording
      bindNavOnlyRecording = p.navOnlyRecording
      bindHoverRecord = p.hoverRecord
    },
    flushRecorderPrefsBindLocals() {
      stores.recorderPrefsStore.patch({
        filterRecording: bindFilterRecording,
        navOnlyRecording: bindNavOnlyRecording,
        hoverRecord: bindHoverRecord,
      })
    },
    syncRecordFormBindLocals() {
      const f = stores.recordFormStore.snapshot()
      bindRecordURL = f.recordURL
      bindRecordOutput = f.recordOutput
      bindRecordIdle = f.recordIdle
      bindRecordAppendTo = f.recordAppendTo
      bindRecordTestClient = f.recordTestClient
      bindRecordFeatureName = f.recordFeatureName
      bindRecordScenarioName = f.recordScenarioName
      bindRecordMode = f.recordMode
      bindRecordStepPickerOpen = f.recordStepPickerOpen
    },
    flushRecordFormBindLocals() {
      stores.recordFormStore.patch({
        recordURL: bindRecordURL,
        recordOutput: bindRecordOutput,
        recordIdle: bindRecordIdle,
        recordAppendTo: bindRecordAppendTo,
        recordTestClient: bindRecordTestClient,
        recordFeatureName: bindRecordFeatureName,
        recordScenarioName: bindRecordScenarioName,
        recordMode: bindRecordMode,
        recordStepPickerOpen: bindRecordStepPickerOpen,
      })
    },
    syncRunFormBind(mode: RunFormMode, lastRun: RunForm, defaults: Partial<RunForm> = {}) {
      bindRunForm = runFormFromMode(lastRun, mode, defaults)
      stores.runFormStore.setRunForm(bindRunForm)
    },
    flushProjectReplaceBindLocals() {
      stores.projectReplaceStore.patch({
        findText: bindProjectReplaceFind,
        replaceText: bindProjectReplaceReplace,
        caseSensitive: bindProjectReplaceCaseSensitive,
      })
    },
  }
}

export type DialogBindController = ReturnType<typeof createDialogBindController>
