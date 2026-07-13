import type { editor as MonacoEditor } from 'monaco-editor'

/** Replace full document text without clearing Monaco undo stack (unlike setValue). */
export function replaceModelText(
  editor: MonacoEditor.IStandaloneCodeEditor,
  text: string,
  source = 'external',
): boolean {
  const model = editor.getModel()
  if (!model) return false
  if (model.getValue() === text) return false
  const rawOptions = typeof editor.getRawOptions === 'function' ? editor.getRawOptions() : {}
  const wasReadOnly = rawOptions.readOnly === true
  if (wasReadOnly) {
    editor.updateOptions({ readOnly: false })
  }
  try {
    editor.pushUndoStop()
    editor.executeEdits(source, [
      {
        range: model.getFullModelRange(),
        text,
        forceMoveMarkers: true,
      },
    ])
    editor.pushUndoStop()
  } finally {
    if (wasReadOnly) {
      editor.updateOptions({ readOnly: true })
    }
  }
  return true
}
