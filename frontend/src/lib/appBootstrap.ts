import loader from '@monaco-editor/loader'
import { ensureMonacoEnvironment } from '../monaco-env'
import { CompletionsForLine, DescribeEditorLine, FormatFeature } from '../../wailsjs/go/wailsapp/App'
import { registerFeatureLanguage } from './featureLanguage'
import { registerGherkinCompletions } from './gherkinCompletions'
import { registerGherkinStepHover } from './gherkinStepHover'
import { registerGherkinFormatProvider } from './gherkinFormat'
import { registerGherkinDocumentSymbols } from './gherkinDocumentSymbols'
import { registerGherkinFolding } from './gherkinFolding'

export type MonacoApi = typeof import('monaco-editor')

let monacoReady: Promise<MonacoApi> | null = null

function registerMonacoProviders(api: MonacoApi) {
  registerFeatureLanguage(api)
  registerGherkinCompletions(api, CompletionsForLine)
  registerGherkinStepHover(api, async (line) => {
    const entry = await DescribeEditorLine(line)
    return entry?.label || entry?.template || entry?.action ? entry : null
  })
  registerGherkinFormatProvider(api, async (text) => FormatFeature(text))
  registerGherkinDocumentSymbols(api)
  registerGherkinFolding(api)
}

/** Lazy-loads Monaco (separate chunk) and registers scenaria-feature language + providers. */
export function preloadMonacoEditor(): Promise<MonacoApi> {
  if (!monacoReady) {
    monacoReady = ensureMonacoEnvironment().then(() =>
      import('monaco-editor').then((monaco) => {
        loader.config({ monaco })
        return loader.init().then((api) => {
          registerMonacoProviders(api)
          return api
        })
      }),
    )
  }
  return monacoReady
}

/** Fire-and-forget prefetch after the shell is visible (cold start). */
export function prefetchMonacoEditor(): void {
  void preloadMonacoEditor().catch((err) => {
    console.error('Monaco prefetch failed', err)
  })
}
