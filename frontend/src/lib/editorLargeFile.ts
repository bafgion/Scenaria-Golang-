import type { editor as MonacoEditor } from 'monaco-editor'
import type { EditorSettings } from './editorOptions'
import { toMonacoOptions } from './editorOptions'

/** Line count at which heavy editor features are toned down. */
export const LARGE_FILE_LINE_THRESHOLD = 2000

export function isLargeFeatureFile(lineCount: number): boolean {
  return lineCount >= LARGE_FILE_LINE_THRESHOLD
}

/** Language providers that parse the whole file (symbols, folding, hover RPC). */
export function shouldUseHeavyLanguageFeatures(lineCount: number): boolean {
  return !isLargeFeatureFile(lineCount)
}

export function largeFileEditorOverrides(
  lineCount: number,
): Partial<MonacoEditor.IEditorOptions> {
  if (!isLargeFeatureFile(lineCount)) {
    return {}
  }
  return {
    minimap: { enabled: false },
    codeLens: false,
    inlayHints: { enabled: 'off' },
  }
}

export function editorOptionsForLineCount(
  settings: EditorSettings,
  monaco: typeof import('monaco-editor'),
  lineCount: number,
): MonacoEditor.IStandaloneEditorConstructionOptions {
  const base = toMonacoOptions(settings, monaco)
  if (!isLargeFeatureFile(lineCount)) {
    return base
  }
  return {
    ...base,
    ...largeFileEditorOverrides(lineCount),
  }
}
