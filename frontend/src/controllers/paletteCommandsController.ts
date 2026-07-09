import type { PaletteCommand } from '../lib/paletteTypes'
import type { RunForm, RunFormMode } from '../lib/runTypes'
import type { gui } from '../../wailsjs/go/models'

export type PaletteTr = (key: string, vars?: Record<string, string | number>) => string

export type PaletteViewState = {
  toolbarCompact: boolean
  previewVisible: boolean
  stepsPanelVisible: boolean
  activeTab: string
  isWelcome: boolean
  installedPlugins: gui.PluginEntryDTO[]
  allureDir: string
}

export type PaletteActions = {
  openCommandPalette: () => void
  selectWelcome: () => void
  openProject: () => void
  openNewProject: () => void
  closeProject: () => void
  openSettings: () => void
  openInitProject: () => void
  openExamples: () => void
  newScenario: () => void
  openFile: () => void
  saveFeature: () => void
  saveFeatureAs: () => void
  openExport: () => void
  openImport: () => void
  openImportFeatures: () => void
  openSteps: () => void
  openSnippets: () => void
  openFindReplace: () => void
  openFind: () => void
  formatDocument: () => void
  openSymbolOutline: () => void
  openProjectReplace: () => void
  openDuplicate: (path: string) => void
  openRenameFeature: (path: string) => void
  deleteFeature: (path: string) => void
  refactorIndents: () => void
  refactorBlanks: () => void
  openStepsHelp: () => void
  openBrowser: () => void
  beginRecord: () => void
  toggleRecordPause: () => void
  stopRecord: () => void
  openBaselineRecord: () => void
  runPrimary: (dry: boolean) => void
  runCurrentScenario: (dry: boolean) => void
  openRunDialog: (title: string, defaults: Partial<RunForm>, mode: RunFormMode) => void
  runTagDialog: () => void
  openPlaywrightRun: () => void
  toggleBatchMode: () => void
  runBatchSelected: (dry: boolean) => void
  rerunFailed: () => void
  openRunHistory: () => void
  openTestClient: () => void
  openTestClientCapture: () => void
  openValidate: (syntaxOnly: boolean) => void
  openVanessa: (dry: boolean, preferRerun?: boolean) => void
  openVanessaSettings: () => void
  openVanessaMonitor: () => void
  openPlugins: () => void
  openPluginRun: (name: string, dry: boolean) => void
  openJournal: () => void
  openResults: () => void
  serveAllure: (dir: string) => void
  openValidatePanel: () => void
  openErrorPanel: () => void
  showSidebar: () => void
  hideSidebar: () => void
  togglePreview: () => void
  toggleStepsPanel: () => void
  toggleToolbarCompact: () => void
  refactorUrls: () => void
  openHotkeys: () => void
  resetLayout: () => void
  checkUpdates: () => void
  showAbout: () => void
  hasVanessaPlugin: () => boolean
  pluginLabel: (plugin: gui.PluginEntryDTO) => string
}

export function buildPaletteCommands(
  tr: PaletteTr,
  state: PaletteViewState,
  actions: PaletteActions,
): PaletteCommand[] {
  const pg = (key: string) => tr(`palette.groups.${key}`)
  const pc = (key: string) => tr(`palette.commands.${key}`)
  const compactLabel = state.toolbarCompact ? pc('expanded') : pc('compact')
  return [
    { id: 'palette', label: pc('palette'), group: pg('view'), shortcut: 'Ctrl+Shift+P', run: actions.openCommandPalette },
    { id: 'welcome', label: pc('welcome'), group: pg('view'), run: actions.selectWelcome },
    { id: 'open', label: pc('open'), group: pg('project'), run: actions.openProject },
    { id: 'new-project', label: pc('newProject'), group: pg('project'), run: actions.openNewProject },
    { id: 'close-project', label: pc('closeProject'), group: pg('project'), run: actions.closeProject },
    { id: 'settings', label: pc('settings'), group: pg('project'), shortcut: 'Ctrl+,', run: actions.openSettings },
    { id: 'init', label: pc('init'), group: pg('project'), run: actions.openInitProject },
    { id: 'examples', label: pc('examples'), group: pg('project'), run: actions.openExamples },
    { id: 'new', label: pc('new'), group: pg('scenario'), shortcut: 'Ctrl+N', run: actions.newScenario },
    { id: 'open-file', label: pc('openFile'), group: pg('scenario'), shortcut: 'Ctrl+O', run: actions.openFile },
    { id: 'save', label: pc('save'), group: pg('scenario'), shortcut: 'Ctrl+S', run: actions.saveFeature },
    { id: 'save-as', label: pc('saveAs'), group: pg('scenario'), shortcut: 'Ctrl+Shift+S', run: actions.saveFeatureAs },
    { id: 'export', label: pc('export'), group: pg('scenario'), run: actions.openExport },
    { id: 'import', label: pc('import'), group: pg('scenario'), run: actions.openImport },
    { id: 'import-features', label: pc('importFeatures'), group: pg('scenario'), run: actions.openImportFeatures },
    { id: 'steps', label: pc('steps'), group: pg('scenario'), run: actions.openSteps },
    { id: 'snippets', label: pc('snippets'), group: pg('scenario'), shortcut: 'Ctrl+Shift+Space', run: actions.openSnippets },
    { id: 'find-replace', label: pc('findReplace'), group: pg('scenario'), shortcut: 'Ctrl+H', run: actions.openFindReplace },
    { id: 'find', label: pc('find'), group: pg('scenario'), shortcut: 'Ctrl+F', run: actions.openFind },
    { id: 'format', label: pc('format'), group: pg('scenario'), shortcut: 'Shift+Alt+F', run: () => void actions.formatDocument() },
    { id: 'goto-symbol', label: pc('gotoSymbol'), group: pg('scenario'), shortcut: 'Ctrl+Shift+O', run: actions.openSymbolOutline },
    { id: 'project-replace', label: pc('projectReplace'), group: pg('scenario'), run: actions.openProjectReplace },
    {
      id: 'duplicate',
      label: pc('duplicate'),
      group: pg('scenario'),
      run: () => state.activeTab && !state.isWelcome && actions.openDuplicate(state.activeTab),
    },
    {
      id: 'rename-feature',
      label: pc('renameFeature'),
      group: pg('scenario'),
      run: () => {
        if (!state.activeTab || state.isWelcome) return
        actions.openRenameFeature(state.activeTab)
      },
    },
    {
      id: 'delete-feature',
      label: pc('deleteFeature'),
      group: pg('scenario'),
      run: () => state.activeTab && !state.isWelcome && actions.deleteFeature(state.activeTab),
    },
    { id: 'refactor-indents', label: pc('refactorIndents'), group: pg('refactor'), run: actions.refactorIndents },
    { id: 'refactor-blanks', label: pc('refactorBlanks'), group: pg('refactor'), run: actions.refactorBlanks },
    { id: 'steps-help', label: pc('stepsHelp'), group: pg('help'), shortcut: 'F1', run: actions.openStepsHelp },
    { id: 'browser', label: pc('browser'), group: pg('run'), shortcut: 'Ctrl+B', run: () => void actions.openBrowser() },
    { id: 'record', label: pc('record'), group: pg('run'), shortcut: 'Ctrl+R', run: actions.beginRecord },
    { id: 'record-pause', label: pc('recordPause'), group: pg('run'), shortcut: 'Alt+P', run: () => void actions.toggleRecordPause() },
    { id: 'record-stop', label: pc('recordStop'), group: pg('run'), shortcut: 'Ctrl+Shift+R', run: () => void actions.stopRecord() },
    { id: 'record-baseline', label: pc('recordBaseline'), group: pg('run'), run: actions.openBaselineRecord },
    { id: 'stop', label: pc('stop'), group: pg('run'), run: actions.stopRecord },
    { id: 'pause', label: pc('pause'), group: pg('run'), run: actions.toggleRecordPause },
    { id: 'run', label: pc('run'), group: pg('run'), shortcut: 'Ctrl+Enter', run: () => actions.runPrimary(false) },
    { id: 'run-current', label: pc('runCurrent'), group: pg('run'), shortcut: 'Ctrl+Shift+Enter', run: () => actions.runCurrentScenario(false) },
    { id: 'run-current-dry', label: pc('runCurrentDry'), group: pg('run'), run: () => actions.runCurrentScenario(true) },
    { id: 'run-dialog', label: pc('runDialog'), group: pg('run'), run: () => actions.openRunDialog('', {}, 'single') },
    { id: 'run-tag', label: pc('runTag'), group: pg('run'), run: actions.runTagDialog },
    { id: 'playwright', label: pc('playwright'), group: pg('run'), run: actions.openPlaywrightRun },
    { id: 'dry', label: pc('dry'), group: pg('run'), run: () => actions.runPrimary(true) },
    { id: 'batch', label: pc('batch'), group: pg('run'), run: actions.toggleBatchMode },
    { id: 'batch-run', label: pc('batchRun'), group: pg('run'), run: () => actions.runBatchSelected(false) },
    { id: 'batch-dry', label: pc('batchDry'), group: pg('run'), run: () => actions.runBatchSelected(true) },
    { id: 'rerun-failed', label: pc('rerunFailed'), group: pg('run'), run: actions.rerunFailed },
    { id: 'run-history', label: pc('runHistory'), group: pg('run'), run: actions.openRunHistory },
    { id: 'testclient', label: pc('testclient'), group: pg('run'), run: actions.openTestClient },
    { id: 'capture-session', label: pc('captureSession'), group: pg('run'), run: actions.openTestClientCapture },
    { id: 'validate', label: pc('validate'), group: pg('run'), run: () => actions.openValidate(true) },
    { id: 'validate-browser', label: pc('validateBrowser'), group: pg('run'), run: () => actions.openValidate(false) },
    ...(actions.hasVanessaPlugin()
      ? [
          { id: 'vanessa-dry', label: pc('vanessaDry'), group: pg('run'), run: () => actions.openVanessa(true) },
          { id: 'vanessa', label: pc('vanessa'), group: pg('run'), run: () => actions.openVanessa(false) },
          { id: 'vanessa-rerun', label: pc('vanessaRerun'), group: pg('run'), run: () => actions.openVanessa(false, true) },
          { id: 'vanessa-settings', label: pc('vanessaSettings'), group: pg('run'), run: actions.openVanessaSettings },
          { id: 'vanessa-monitor', label: pc('vanessaMonitor'), group: pg('run'), run: actions.openVanessaMonitor },
        ]
      : []),
    { id: 'plugins', label: pc('plugins'), group: pg('plugins'), run: actions.openPlugins },
    ...state.installedPlugins.flatMap((plugin) => {
      if (!plugin.runnable || plugin.vanessa) return []
      const label = actions.pluginLabel(plugin)
      return [
        {
          id: `plugin-${plugin.name}-dry`,
          label: tr('menus.runPluginDry', { name: label }),
          group: pg('plugins'),
          run: () => actions.openPluginRun(plugin.name, true),
        },
        {
          id: `plugin-${plugin.name}`,
          label: tr('menus.runPlugin', { name: label }),
          group: pg('plugins'),
          run: () => actions.openPluginRun(plugin.name, false),
        },
      ]
    }),
    { id: 'journal', label: pc('journal'), group: pg('view'), shortcut: 'Ctrl+`', run: actions.openJournal },
    { id: 'results', label: pc('results'), group: pg('view'), run: actions.openResults },
    { id: 'allure-serve', label: pc('allureServe'), group: pg('view'), run: () => actions.serveAllure(state.allureDir) },
    { id: 'validate-panel', label: pc('validatePanel'), group: pg('view'), run: actions.openValidatePanel },
    { id: 'error-panel', label: pc('errorPanel'), group: pg('view'), run: actions.openErrorPanel },
    { id: 'explorer', label: pc('explorer'), group: pg('view'), run: actions.showSidebar },
    { id: 'explorer-hide', label: pc('explorerHide'), group: pg('view'), run: actions.hideSidebar },
    {
      id: 'preview',
      label: state.previewVisible ? pc('previewHide') : pc('previewShow'),
      group: pg('view'),
      run: actions.togglePreview,
    },
    {
      id: 'steps-panel',
      label: state.stepsPanelVisible ? pc('stepsPanelHide') : pc('stepsPanelShow'),
      group: pg('view'),
      run: actions.toggleStepsPanel,
    },
    { id: 'compact', label: compactLabel, group: pg('view'), run: actions.toggleToolbarCompact },
    { id: 'refactor-urls', label: pc('refactorUrls'), group: pg('refactor'), run: actions.refactorUrls },
    { id: 'hotkeys', label: pc('hotkeys'), group: pg('help'), shortcut: 'Shift+F1', run: actions.openHotkeys },
    { id: 'reset-layout', label: pc('resetLayout'), group: pg('view'), run: actions.resetLayout },
    { id: 'updates', label: pc('updates'), group: pg('help'), run: actions.checkUpdates },
    { id: 'about', label: pc('about'), group: pg('help'), run: actions.showAbout },
  ]
}
