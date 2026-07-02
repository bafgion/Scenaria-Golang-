import type * as Monaco from 'monaco-editor'
import type { StepHelpEntry } from './stepTypes'
import { shouldUseHeavyLanguageFeatures } from './editorLargeFile'
import { formatStepHoverMarkdown, hasStepHelp } from './stepHelpContent'

export type StepHoverFetcher = (line: string) => Promise<StepHelpEntry | null>

let hoverEnabled = () => true

export function setStepHoverEnabled(getter: () => boolean) {
  hoverEnabled = getter
}

let providerRegistered = false

export function registerGherkinStepHover(monaco: typeof Monaco, fetchStep: StepHoverFetcher) {
  if (providerRegistered) return
  providerRegistered = true

  monaco.languages.registerHoverProvider('scenaria-feature', {
    provideHover: async (model, position) => {
      if (!hoverEnabled() || !shouldUseHeavyLanguageFeatures(model.getLineCount())) return null

      const lineNumber = position.lineNumber
      const line = model.getLineContent(lineNumber)
      if (!line.trim() || line.trimStart().startsWith('#')) {
        return null
      }

      let entry: StepHelpEntry | null = null
      try {
        entry = await fetchStep(line)
      } catch {
        return null
      }
      if (!hasStepHelp(entry)) {
        return null
      }

      const maxColumn = model.getLineMaxColumn(lineNumber)
      return {
        range: {
          startLineNumber: lineNumber,
          endLineNumber: lineNumber,
          startColumn: 1,
          endColumn: maxColumn,
        },
        contents: [{ value: formatStepHoverMarkdown(entry) }],
      }
    },
  })
}
