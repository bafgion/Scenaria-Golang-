import loader from '@monaco-editor/loader'
import * as monaco from 'monaco-editor'
import { CompletionsForLine, DescribeEditorLine, FormatFeature } from '../../wailsjs/go/wailsapp/App'
import { registerFeatureLanguage } from './featureLanguage'
import { registerGherkinCompletions } from './gherkinCompletions'
import { registerGherkinStepHover } from './gherkinStepHover'
import { registerGherkinFormatProvider } from './gherkinFormat'
import { registerGherkinDocumentSymbols } from './gherkinDocumentSymbols'

loader.config({ monaco })

let monacoReady: Promise<typeof monaco> | null = null

/** Preloads Monaco and registers the scenaria-feature language + completions. */
export function preloadMonacoEditor(): Promise<typeof monaco> {
  if (!monacoReady) {
    monacoReady = loader.init().then((api) => {
      registerFeatureLanguage(api)
      registerGherkinCompletions(api, CompletionsForLine)
      registerGherkinStepHover(api, async (line) => {
        const entry = await DescribeEditorLine(line)
        return entry?.label || entry?.template || entry?.action ? entry : null
      })
      registerGherkinFormatProvider(api, async (text) => FormatFeature(text))
      registerGherkinDocumentSymbols(api)
      return api
    })
  }
  return monacoReady
}
