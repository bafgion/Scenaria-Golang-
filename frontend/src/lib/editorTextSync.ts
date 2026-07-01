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
  editor.pushUndoStop()
  editor.executeEdits(source, [
    {
      range: model.getFullModelRange(),
      text,
      forceMoveMarkers: true,
    },
  ])
  editor.pushUndoStop()
  return true
}
