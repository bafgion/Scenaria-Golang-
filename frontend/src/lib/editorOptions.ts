import type { settings as SettingsTypes } from '../../wailsjs/go/models'
import { settings } from '../../wailsjs/go/models'
import type { editor as MonacoEditor } from 'monaco-editor'

export type EditorSettings = {
  fontSize: number
  fontFamily: string
  wordWrap: 'on' | 'off'
  minimap: boolean
  lineNumbers: 'on' | 'off' | 'relative'
  tabSize: number
  insertSpaces: boolean
  renderWhitespace: 'none' | 'boundary' | 'selection' | 'trailing' | 'all'
  folding: boolean
  stickyScroll: boolean
  autoClosingQuotes: 'always' | 'languageDefined' | 'beforeWhitespace' | 'never'
  formatOnSave: boolean
  stepHover: boolean
  validateOnType: boolean
  theme: 'scenaria-dark' | 'scenaria-light'
  breadcrumbs: boolean
  symbolOutline: boolean
  stepsPanelView: 'outline' | 'steps'
  codeLens: boolean
  inlayHints: boolean
  scenarioHints: boolean
  scenarioHintsAfterRecord: boolean
  scenarioHintsShowWarning: boolean
  scenarioHintsShowInfo: boolean
  scenarioHintsAutoFixOnSave: boolean
}

export const DEFAULT_EDITOR_SETTINGS: EditorSettings = {
  fontSize: 13,
  fontFamily: '"Cascadia Code", "JetBrains Mono", Consolas, "Courier New", monospace',
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
  stepHover: true,
  validateOnType: true,
  theme: 'scenaria-dark',
  breadcrumbs: true,
  symbolOutline: true,
  stepsPanelView: 'outline',
  codeLens: true,
  inlayHints: false,
  scenarioHints: true,
  scenarioHintsAfterRecord: true,
  scenarioHintsShowWarning: true,
  scenarioHintsShowInfo: true,
  scenarioHintsAutoFixOnSave: false,
}

export function editorSettingsFromDTO(raw?: SettingsTypes.EditorSettings | null): EditorSettings {
  const base = { ...DEFAULT_EDITOR_SETTINGS }
  if (!raw) return base

  if (raw.fontSize !== undefined && raw.fontSize >= 8 && raw.fontSize <= 32) base.fontSize = raw.fontSize
  if (raw.fontFamily) base.fontFamily = raw.fontFamily
  if (raw.wordWrap === 'on' || raw.wordWrap === 'off') base.wordWrap = raw.wordWrap
  if (typeof raw.minimap === 'boolean') base.minimap = raw.minimap
  if (raw.lineNumbers === 'on' || raw.lineNumbers === 'off' || raw.lineNumbers === 'relative') {
    base.lineNumbers = raw.lineNumbers
  }
  if (raw.tabSize !== undefined && raw.tabSize >= 1 && raw.tabSize <= 8) base.tabSize = raw.tabSize
  if (typeof raw.insertSpaces === 'boolean') base.insertSpaces = raw.insertSpaces
  if (
    raw.renderWhitespace === 'none' ||
    raw.renderWhitespace === 'boundary' ||
    raw.renderWhitespace === 'selection' ||
    raw.renderWhitespace === 'trailing' ||
    raw.renderWhitespace === 'all'
  ) {
    base.renderWhitespace = raw.renderWhitespace
  }
  if (typeof raw.folding === 'boolean') base.folding = raw.folding
  if (typeof raw.stickyScroll === 'boolean') base.stickyScroll = raw.stickyScroll
  if (
    raw.autoClosingQuotes === 'always' ||
    raw.autoClosingQuotes === 'languageDefined' ||
    raw.autoClosingQuotes === 'beforeWhitespace' ||
    raw.autoClosingQuotes === 'never'
  ) {
    base.autoClosingQuotes = raw.autoClosingQuotes
  }
  if (typeof raw.formatOnSave === 'boolean') base.formatOnSave = raw.formatOnSave
  if (typeof raw.stepHoverEnabled === 'boolean') base.stepHover = raw.stepHoverEnabled
  if (typeof raw.validateOnType === 'boolean') base.validateOnType = raw.validateOnType
  if (raw.theme === 'scenaria-dark' || raw.theme === 'scenaria-light') base.theme = raw.theme
  if (typeof raw.breadcrumbsEnabled === 'boolean') base.breadcrumbs = raw.breadcrumbsEnabled
  if (typeof raw.symbolOutlineEnabled === 'boolean') base.symbolOutline = raw.symbolOutlineEnabled
  if (raw.stepsPanelView === 'outline' || raw.stepsPanelView === 'steps') {
    base.stepsPanelView = raw.stepsPanelView
  }
  if (typeof raw.codeLensEnabled === 'boolean') base.codeLens = raw.codeLensEnabled
  if (typeof raw.inlayHintsEnabled === 'boolean') base.inlayHints = raw.inlayHintsEnabled
  if (typeof raw.scenarioHintsEnabled === 'boolean') base.scenarioHints = raw.scenarioHintsEnabled
  if (typeof raw.scenarioHintsAfterRecord === 'boolean') {
    base.scenarioHintsAfterRecord = raw.scenarioHintsAfterRecord
  }
  if (typeof raw.scenarioHintsShowWarnings === 'boolean') {
    base.scenarioHintsShowWarning = raw.scenarioHintsShowWarnings
  }
  if (typeof raw.scenarioHintsShowInfo === 'boolean') base.scenarioHintsShowInfo = raw.scenarioHintsShowInfo
  if (typeof raw.scenarioHintsAutoFixOnSave === 'boolean') {
    base.scenarioHintsAutoFixOnSave = raw.scenarioHintsAutoFixOnSave
  }
  return base
}

export function editorSettingsToDTO(editor: EditorSettings): SettingsTypes.EditorSettings {
  return settings.EditorSettings.createFrom({
    fontSize: editor.fontSize,
    fontFamily: editor.fontFamily,
    wordWrap: editor.wordWrap,
    minimap: editor.minimap,
    lineNumbers: editor.lineNumbers,
    tabSize: editor.tabSize,
    insertSpaces: editor.insertSpaces,
    renderWhitespace: editor.renderWhitespace,
    folding: editor.folding,
    stickyScroll: editor.stickyScroll,
    autoClosingQuotes: editor.autoClosingQuotes,
    formatOnSave: editor.formatOnSave,
    stepHoverEnabled: editor.stepHover,
    validateOnType: editor.validateOnType,
    theme: editor.theme,
    breadcrumbsEnabled: editor.breadcrumbs,
    symbolOutlineEnabled: editor.symbolOutline,
    stepsPanelView: editor.stepsPanelView,
    codeLensEnabled: editor.codeLens,
    inlayHintsEnabled: editor.inlayHints,
    scenarioHintsEnabled: editor.scenarioHints,
    scenarioHintsAfterRecord: editor.scenarioHintsAfterRecord,
    scenarioHintsShowWarnings: editor.scenarioHintsShowWarning,
    scenarioHintsShowInfo: editor.scenarioHintsShowInfo,
    scenarioHintsAutoFixOnSave: editor.scenarioHintsAutoFixOnSave,
  })
}

export function toMonacoOptions(
  settings: EditorSettings,
  monaco: typeof import('monaco-editor'),
): MonacoEditor.IStandaloneEditorConstructionOptions {
  return {
    fontSize: settings.fontSize,
    fontFamily: settings.fontFamily,
    minimap: { enabled: settings.minimap, scale: 1 },
    wordWrap: settings.wordWrap,
    lineNumbers: settings.lineNumbers,
    tabSize: settings.tabSize,
    insertSpaces: settings.insertSpaces,
    renderWhitespace: settings.renderWhitespace,
    folding: settings.folding,
    stickyScroll: { enabled: settings.stickyScroll },
    autoClosingQuotes: settings.autoClosingQuotes,
    lightbulb: { enabled: monaco.editor.ShowLightbulbIconMode.On },
    breadcrumbs: { enabled: settings.breadcrumbs },
    codeLens: settings.codeLens,
    inlayHints: { enabled: settings.inlayHints ? 'on' : 'off' },
  } as MonacoEditor.IStandaloneEditorConstructionOptions
}
