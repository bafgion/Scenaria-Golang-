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

export type CompletionFetcher = (line: string, column: number) => Promise<gui.StepCompletionsDTO>

const STEP_KEYWORD_RE = /^(?:Допустим|Дано|Когда|Тогда|И|Но)\s+/i

export function formatInsertText(line: string, snippet: gui.StepCompletionSnippet): string {
  const trimmed = line.trimStart()
  const indent = line.slice(0, line.length - trimmed.length)
  if (snippet.label === snippet.insert && /^(Допустим|Дано|Когда|Тогда|И|Но)$/i.test(snippet.label)) {
    if (!trimmed) {
      return `${indent}${snippet.insert} `
    }
    return `${snippet.insert} `
  }
  if (!STEP_KEYWORD_RE.test(trimmed)) {
    return `Когда ${snippet.insert}`
  }
  return snippet.insert
}

let providerRegistered = false

export function registerGherkinCompletions(monaco: typeof Monaco, fetchCompletions: CompletionFetcher) {
  if (providerRegistered) {
    return
  }
  providerRegistered = true

  monaco.languages.registerCompletionItemProvider('scenaria-feature', {
    triggerCharacters: [' ', '"', "'", '.', '@'],
    provideCompletionItems: async (model, position) => {
      if (!shouldUseHeavyLanguageFeatures(model.getLineCount())) {
        return { suggestions: [] }
      }
      const line = model.getLineContent(position.lineNumber)
      const runeColumn = monacoColumnToRuneIndex(line, position.column - 1)
      let result: gui.StepCompletionsDTO
      try {
        result = await fetchCompletions(line, runeColumn)
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
        const formatted = formatInsertText(line, item)
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
}
