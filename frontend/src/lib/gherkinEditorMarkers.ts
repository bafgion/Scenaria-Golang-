import type * as Monaco from 'monaco-editor'
import type { gui } from '../../wailsjs/go/models'
import { HINT_MARKER_SOURCE } from './gherkinHintActionsHelpers'

export const EDITOR_MARKER_OWNER = 'scenaria-editor'
export const VALIDATION_MARKER_SOURCE = 'Валидация'

export type ValidationMarkerIssue = { line: number; message: string }

export function buildEditorMarkers(
  model: Monaco.editor.ITextModel,
  validationIssues: ValidationMarkerIssue[],
  hints: gui.ScenarioHintDTO[],
  monacoInstance: typeof Monaco,
): Monaco.editor.IMarkerData[] {
  const markers: Monaco.editor.IMarkerData[] = []

  for (const issue of validationIssues) {
    if (issue.line < 1 || issue.line > model.getLineCount()) continue
    markers.push({
      startLineNumber: issue.line,
      endLineNumber: issue.line,
      startColumn: 1,
      endColumn: model.getLineMaxColumn(issue.line),
      message: issue.message,
      severity: monacoInstance.MarkerSeverity.Error,
      source: VALIDATION_MARKER_SOURCE,
    })
  }

  for (const hint of hints) {
    if (hint.line < 1 || hint.line > model.getLineCount()) continue
    markers.push({
      startLineNumber: hint.line,
      endLineNumber: hint.line,
      startColumn: 1,
      endColumn: model.getLineMaxColumn(hint.line),
      message: hint.title,
      severity:
        hint.severity === 'warning'
          ? monacoInstance.MarkerSeverity.Warning
          : monacoInstance.MarkerSeverity.Info,
      source: HINT_MARKER_SOURCE,
      code: hint.id,
    })
  }

  return markers
}

export function applyEditorMarkers(
  monacoInstance: typeof Monaco,
  model: Monaco.editor.ITextModel,
  validationIssues: ValidationMarkerIssue[],
  hints: gui.ScenarioHintDTO[],
) {
  monacoInstance.editor.setModelMarkers(
    model,
    EDITOR_MARKER_OWNER,
    buildEditorMarkers(model, validationIssues, hints, monacoInstance),
  )
}
