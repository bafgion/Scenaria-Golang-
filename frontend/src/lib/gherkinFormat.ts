import type * as Monaco from 'monaco-editor'

export type FeatureFormatter = (text: string) => Promise<string>

let providerRegistered = false

export function registerGherkinFormatProvider(monaco: typeof Monaco, format: FeatureFormatter) {
  if (providerRegistered) return
  providerRegistered = true

  monaco.languages.registerDocumentFormattingEditProvider('scenaria-feature', {
    provideDocumentFormattingEdits: async (model) => {
      const original = model.getValue()
      const formatted = await format(original)
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
}
