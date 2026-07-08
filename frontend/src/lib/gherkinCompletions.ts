import type * as Monaco from 'monaco-editor'
import type { gui } from '../../wailsjs/go/models'
import { shouldUseHeavyLanguageFeatures } from './editorLargeFile'
import { monacoColumnToRuneIndex, runeIndexToMonacoColumn, runeSlice } from './editorColumns'
import {
  completionFilterText,
  completionSortKey,
  shouldPreselectCompletion,
  snippetizeInsert,
  usesSnippetTabStops,
} from './gherkinCompletionSnippets'
import { defaultStepKeyword, detectFeatureGherkinLanguage, type FeatureGherkinLanguage } from './featureGherkinLang'
import { replaceMonacoDisposables } from './monacoDisposables'

export type CompletionFetcher = (line: string, column: number, language: string) => Promise<gui.StepCompletionsDTO>

const STEP_KEYWORD_RE =
  /^(?:Допустим|Дано|Когда|Тогда|И|Но|Given|When|Then|And|But)\s+/i

const KEYWORD_ONLY_RE = /^(?:Допустим|Дано|Когда|Тогда|И|Но|Given|When|Then|And|But)$/i

function detectFeatureLanguageFromModel(model: Monaco.editor.ITextModel): FeatureGherkinLanguage {
  const maxScan = Math.min(model.getLineCount(), 32)
  const header: string[] = []
  for (let lineNo = 1; lineNo <= maxScan; lineNo++) {
    const line = model.getLineContent(lineNo).trim()
    if (!line) continue
    if (!line.startsWith('#')) break
    header.push(line)
  }
  return detectFeatureGherkinLanguage(header.join('\n'))
}

export function formatInsertText(line: string, snippet: gui.StepCompletionSnippet, lang: FeatureGherkinLanguage): string {
  const trimmed = line.trimStart()
  const indent = line.slice(0, line.length - trimmed.length)
  if (snippet.label === snippet.insert && KEYWORD_ONLY_RE.test(snippet.label)) {
    if (!trimmed) {
      return `${indent}${snippet.insert} `
    }
    return `${snippet.insert} `
  }
  if (!STEP_KEYWORD_RE.test(trimmed)) {
    return `${defaultStepKeyword(lang)} ${snippet.insert}`
  }
  return snippet.insert
}

let providerRegistered = false

export function registerGherkinCompletions(monaco: typeof Monaco, fetchCompletions: CompletionFetcher) {
  if (providerRegistered) {
    return
  }
  providerRegistered = true

  const disposable = monaco.languages.registerCompletionItemProvider('scenaria-feature', {
    triggerCharacters: [' ', '"', "'", '.', '@'],
    provideCompletionItems: async (model, position) => {
      if (!shouldUseHeavyLanguageFeatures(model.getLineCount())) {
        return { suggestions: [] }
      }
      const line = model.getLineContent(position.lineNumber)
      const runeColumn = monacoColumnToRuneIndex(line, position.column - 1)
      const lang = detectFeatureLanguageFromModel(model)
      let result: gui.StepCompletionsDTO
      try {
        result = await fetchCompletions(line, runeColumn, lang)
      } catch {
        return { suggestions: [] }
      }
      if (!result.items?.length) {
        return { suggestions: [] }
      }

      const startCol = runeIndexToMonacoColumn(line, result.start) + 1
      const endCol = runeIndexToMonacoColumn(line, result.end) + 1
      const range = {
        startLineNumber: position.lineNumber,
        endLineNumber: position.lineNumber,
        startColumn: startCol,
        endColumn: endCol,
      }
      const typedPrefix = runeSlice(line, result.start, result.end)

      const suggestions = result.items.map((item, index) => {
        const formatted = formatInsertText(line, item, lang)
        const insertText = snippetizeInsert(formatted)
        const snippet = usesSnippetTabStops(insertText)
        return {
          label: item.label,
          kind: monaco.languages.CompletionItemKind.Snippet,
          insertText,
          insertTextRules: snippet
            ? monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet
            : undefined,
          filterText: completionFilterText(item),
          detail: item.description,
          documentation: { value: `\`\`\`\n${item.insert}\n\`\`\`` },
          range,
          sortText: completionSortKey(item.label, typedPrefix),
          preselect: shouldPreselectCompletion(item.label, typedPrefix, index),
        }
      })

      return { suggestions }
    },
  })
  replaceMonacoDisposables('gherkin-completions', [disposable])
}
