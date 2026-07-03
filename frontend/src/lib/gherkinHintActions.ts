import type * as Monaco from 'monaco-editor'
import type { gui } from '../../wailsjs/go/models'
import {
  findHintForMarker,
  HINT_QUICK_FIX_KIND,
  hintMarkerSource,
  rangeTouchesMarker,
} from './gherkinHintActionsHelpers'
import { EDITOR_MARKER_OWNER } from './gherkinEditorMarkers'
import { t } from './i18n'

export { findHintForMarker, markerCode, rangeTouchesMarker } from './gherkinHintActionsHelpers'

export type HintActionHandlers = {
  getHints: () => gui.ScenarioHintDTO[]
  onFix: (hint: gui.ScenarioHintDTO) => void | Promise<void>
  onDismiss: (hint: gui.ScenarioHintDTO) => void
}

let monacoApi: typeof Monaco | null = null
let activeHandlers: HintActionHandlers | null = null
let providerRegistered = false

function findHintByArgs(hintId: string, stepIndex: number): gui.ScenarioHintDTO | undefined {
  return activeHandlers?.getHints().find((h) => h.id === hintId && h.stepIndex === stepIndex)
}

export function registerHintCodeActions(monacoInstance: typeof Monaco, handlers: HintActionHandlers) {
  monacoApi = monacoInstance
  activeHandlers = handlers
  if (providerRegistered) return
  providerRegistered = true

  monacoInstance.editor.registerCommand('scenaria.hint.fix', async (_accessor, hintId: string, stepIndex: number) => {
    const hint = findHintByArgs(hintId, stepIndex)
    if (hint) await activeHandlers?.onFix(hint)
  })
  monacoInstance.editor.registerCommand('scenaria.hint.dismiss', (_accessor, hintId: string, stepIndex: number) => {
    const hint = findHintByArgs(hintId, stepIndex)
    if (hint) activeHandlers?.onDismiss(hint)
  })

  monacoInstance.languages.registerCodeActionProvider(
    'scenaria-feature',
    {
      provideCodeActions: (model, range) => {
        const api = monacoApi
        if (!api) return { actions: [], dispose: () => {} }

        const hints = activeHandlers?.getHints() ?? []
        const actions: Monaco.languages.CodeAction[] = []
        const hintSource = hintMarkerSource()

        const markers = api.editor
          .getModelMarkers({ resource: model.uri })
          .filter((m) => m.owner === EDITOR_MARKER_OWNER && m.source === hintSource)

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
              title: t('editor.hint.fix', { title: hint.title }),
              kind: HINT_QUICK_FIX_KIND,
              diagnostics: [diagnostic],
              isPreferred: true,
              command: {
                id: 'scenaria.hint.fix',
                title: hint.title,
                arguments: [hint.id, hint.stepIndex],
              },
            })
          }

          const dismissTitle = t('editor.hint.ignore')
          actions.push({
            title: dismissTitle,
            kind: HINT_QUICK_FIX_KIND,
            diagnostics: [diagnostic],
            command: {
              id: 'scenaria.hint.dismiss',
              title: dismissTitle,
              arguments: [hint.id, hint.stepIndex],
            },
          })
        }

        return { actions, dispose: () => {} }
      },
    },
    {
      providedCodeActionKinds: [HINT_QUICK_FIX_KIND],
    },
  )
}
