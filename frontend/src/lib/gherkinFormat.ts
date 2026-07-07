import type * as Monaco from 'monaco-editor'
import { replaceMonacoDisposables } from './monacoDisposables'

export type FeatureFormatter = (text: string) => Promise<string>

let providerRegistered = false

export function registerGherkinFormatProvider(monaco: typeof Monaco, format: FeatureFormatter) {
  if (providerRegistered) return
  providerRegistered = true

  const disposable = monaco.languages.registerDocumentFormattingEditProvider('scenaria-feature', {
    provideDocumentFormattingEdits: async (model) => {
      const version = model.getVersionId()
      const original = model.getValue()
      const formatted = await format(original)
      if (model.getVersionId() !== version) {
        return []
      }
      if (formatted === original) {
        return []
      }
      return [
        {
          range: model.getFullModelRange(),
          text: formatted,
        },
      ]
    },
  })
  replaceMonacoDisposables('gherkin-format', [disposable])
}
