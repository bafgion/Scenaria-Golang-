import type * as Monaco from 'monaco-editor'
import type { gui } from '../../wailsjs/go/models'

export const HINT_MARKER_OWNER = 'scenaria-hints'
const HINT_MARKER_SOURCE = 'Подсказка'

export type HintActionHandlers = {
  getHints: () => gui.ScenarioHintDTO[]
  onFix: (hint: gui.ScenarioHintDTO) => void | Promise<void>
  onDismiss: (hint: gui.ScenarioHintDTO) => void
}

let monacoApi: typeof Monaco | null = null
let activeHandlers: HintActionHandlers | null = null
let providerRegistered = false

export function markerCode(marker: Monaco.editor.IMarkerData): string {
  const code = marker.code
  if (typeof code === 'string') return code
  if (code && typeof code === 'object' && 'value' in code) return String(code.value)
  return ''
}

export function rangeTouchesMarker(range: Monaco.IRange, marker: Monaco.editor.IMarkerData): boolean {
  return (
    range.endLineNumber >= marker.startLineNumber &&
    range.startLineNumber <= marker.endLineNumber
  )
}

export function findHintForMarker(
  hints: gui.ScenarioHintDTO[],
  marker: Monaco.editor.IMarkerData,
): gui.ScenarioHintDTO | undefined {
  if (marker.source !== HINT_MARKER_SOURCE) return undefined
  const code = markerCode(marker)
  return hints.find((h) => h.line === marker.startLineNumber && (!code || h.id === code))
}

export function setHintMarkers(
  monaco: typeof Monaco,
  model: Monaco.editor.ITextModel,
  hints: gui.ScenarioHintDTO[],
) {
  monacoApi = monaco
  monaco.editor.setModelMarkers(
    model,
    HINT_MARKER_OWNER,
    hints
      .filter((h) => h.line > 0)
      .map((hint) => ({
        startLineNumber: hint.line,
        endLineNumber: hint.line,
        startColumn: 1,
        endColumn: model.getLineMaxColumn(hint.line),
        message: hint.title,
        severity:
          hint.severity === 'warning'
            ? monaco.MarkerSeverity.Warning
            : monaco.MarkerSeverity.Info,
        source: HINT_MARKER_SOURCE,
        code: hint.id,
      })),
  )
}

function findHintByArgs(hintId: string, stepIndex: number): gui.ScenarioHintDTO | undefined {
  return activeHandlers?.getHints().find((h) => h.id === hintId && h.stepIndex === stepIndex)
}

export function registerHintCodeActions(monaco: typeof Monaco, handlers: HintActionHandlers) {
  monacoApi = monaco
  activeHandlers = handlers
  if (providerRegistered) return
  providerRegistered = true

  monaco.editor.registerCommand('scenaria.hint.fix', async (_accessor, hintId: string, stepIndex: number) => {
    const hint = findHintByArgs(hintId, stepIndex)
    if (hint) await activeHandlers?.onFix(hint)
  })
  monaco.editor.registerCommand('scenaria.hint.dismiss', (_accessor, hintId: string, stepIndex: number) => {
    const hint = findHintByArgs(hintId, stepIndex)
    if (hint) activeHandlers?.onDismiss(hint)
  })

  monaco.languages.registerCodeActionProvider(
    'scenaria-feature',
    {
      provideCodeActions: (model, range) => {
        const api = monacoApi
        if (!api) return { actions: [], dispose: () => {} }

        const hints = activeHandlers?.getHints() ?? []
        const actions: Monaco.languages.CodeAction[] = []

        const markers = api.editor
          .getModelMarkers({ resource: model.uri })
          .filter((m) => m.owner === HINT_MARKER_OWNER)

        for (const marker of markers) {
          if (!rangeTouchesMarker(range, marker)) continue

          const hint = findHintForMarker(hints, marker)
          if (!hint) continue

          const diagnostic: Monaco.editor.IMarkerData = {
            severity: marker.severity,
            message: marker.message,
            source: marker.source,
            code: marker.code,
            startLineNumber: marker.startLineNumber,
            startColumn: marker.startColumn,
            endLineNumber: marker.endLineNumber,
            endColumn: marker.endColumn,
          }

          if (hint.autoFixable) {
            actions.push({
              title: `Исправить: ${hint.title}`,
              kind: api.languages.CodeActionKind.QuickFix,
              diagnostics: [diagnostic],
              isPreferred: true,
              command: {
                id: 'scenaria.hint.fix',
                title: hint.title,
                arguments: [hint.id, hint.stepIndex],
              },
            })
          }

          actions.push({
            title: 'Игнорировать подсказку',
            kind: api.languages.CodeActionKind.QuickFix,
            diagnostics: [diagnostic],
            command: {
              id: 'scenaria.hint.dismiss',
              title: 'Игнорировать подсказку',
              arguments: [hint.id, hint.stepIndex],
            },
          })
        }

        return { actions, dispose: () => {} }
      },
    },
    {
      providedCodeActionKinds: [monaco.languages.CodeActionKind.QuickFix],
    },
  )
}
