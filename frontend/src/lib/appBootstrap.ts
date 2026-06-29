import loader from '@monaco-editor/loader'
import * as monaco from 'monaco-editor'
import { CompletionsForLine } from '../../wailsjs/go/wailsapp/App'
import { registerFeatureLanguage } from './featureLanguage'
import { registerGherkinCompletions } from './gherkinCompletions'

loader.config({ monaco })

let monacoReady: Promise<typeof monaco> | null = null

/** Preloads Monaco and registers the scenaria-feature language + completions. */
export function preloadMonacoEditor(): Promise<typeof monaco> {
  if (!monacoReady) {
    monacoReady = loader.init().then((api) => {
      registerFeatureLanguage(api)
      registerGherkinCompletions(api, CompletionsForLine)
      return api
    })
  }
  return monacoReady
}
