import type * as Monaco from 'monaco-editor'
import { formatStepInlayLabel, inlayHintColumn, type StepInlayRow } from './gherkinInlayHints'

export type InlayHintsHandlers = {
  isEnabled: () => boolean
  getSteps: () => StepInlayRow[]
}

let activeHandlers: InlayHintsHandlers | null = null
let providerRegistered = false

function buildInlayHints(
  model: Monaco.editor.ITextModel,
  steps: StepInlayRow[],
  monaco: typeof Monaco,
): Monaco.languages.InlayHint[] {
  const hints: Monaco.languages.InlayHint[] = []
  for (const step of steps) {
    if (step.line < 1 || step.line > model.getLineCount()) continue
    const label = formatStepInlayLabel(step)
    if (!label) continue
    const lineText = model.getLineContent(step.line)
    hints.push({
      position: { lineNumber: step.line, column: inlayHintColumn(lineText) },
      label,
      kind: monaco.languages.InlayHintKind.Type,
      paddingLeft: true,
    })
  }
  return hints
}

export function registerGherkinInlayHints(monacoInstance: typeof Monaco, handlers: InlayHintsHandlers) {
  activeHandlers = handlers
  if (providerRegistered) return
  providerRegistered = true

  monacoInstance.languages.registerInlayHintsProvider('scenaria-feature', {
    provideInlayHints(model, _range, token) {
      const current = activeHandlers
      if (!current || token.isCancellationRequested || !current.isEnabled()) {
        return { hints: [], dispose: () => {} }
      }
      const hints = buildInlayHints(model, current.getSteps(), monacoInstance)
      return { hints, dispose: () => {} }
    },
  })
}

export function refreshGherkinInlayHints(editor: Monaco.editor.ICodeEditor | null) {
  if (!editor) return
  void editor.getAction('editor.action.inlayHints.refresh')?.run()
}
