import type * as Monaco from 'monaco-editor'
import { shouldUseHeavyLanguageFeatures } from './editorLargeFile'
import { getCachedFeatureSymbols } from './featureSymbolCache'
import type { FeatureSymbol } from './gherkinDocumentSymbols'
import { scenarioAtLine } from './scenarioAtLine'

export type RunCodeLensPayload = {
  scenario: string
  line: number
  dryRun: boolean
  partial: boolean
}

export type RunCodeLensHandlers = {
  isEnabled: () => boolean
  onRun: (payload: RunCodeLensPayload) => void | Promise<void>
}

export type RunCodeLensItem = {
  line: number
  title: string
  scenario: string
  dryRun: boolean
  partial: boolean
}

const RUN_SCENARIO_COMMAND = 'scenaria.runScenario'

let monacoApi: typeof Monaco | null = null
let activeHandlers: RunCodeLensHandlers | null = null
let providerRegistered = false

function lineRange(line: number): Monaco.IRange {
  return {
    startLineNumber: line,
    startColumn: 1,
    endLineNumber: line,
    endColumn: 1,
  }
}

function scenarioLenses(node: FeatureSymbol): RunCodeLensItem[] {
  const scenario = node.name.trim()
  if (!scenario) return []
  return [
    { line: node.line, title: '▶ Запустить сценарий', scenario, dryRun: false, partial: false },
    { line: node.line, title: 'Dry-run', scenario, dryRun: true, partial: false },
  ]
}

function stepLenses(text: string, node: FeatureSymbol): RunCodeLensItem[] {
  const scenario = scenarioAtLine(text, node.line)
  if (!scenario) return []
  return [
    { line: node.line, title: '▶ с этой строки', scenario, dryRun: false, partial: true },
    { line: node.line, title: 'Dry-run', scenario, dryRun: true, partial: true },
  ]
}

function walkSymbols(text: string, nodes: FeatureSymbol[], out: RunCodeLensItem[]) {
  for (const node of nodes) {
    if (node.kind === 'scenario' || node.kind === 'outline') {
      out.push(...scenarioLenses(node))
    } else if (node.kind === 'step') {
      out.push(...stepLenses(text, node))
    }
    if (node.children.length > 0) {
      walkSymbols(text, node.children, out)
    }
  }
}

/** Collects run code lens rows for a feature file (testable without Monaco). */
export function collectRunCodeLenses(
  text: string,
  versionId?: number | null,
  modelKey?: string | null,
): RunCodeLensItem[] {
  const out: RunCodeLensItem[] = []
  walkSymbols(text, getCachedFeatureSymbols(text, versionId, modelKey), out)
  return out
}

function toMonacoLenses(items: RunCodeLensItem[]): Monaco.languages.CodeLens[] {
  return items.map((item) => ({
    range: lineRange(item.line),
    command: {
      id: RUN_SCENARIO_COMMAND,
      title: item.title,
      arguments: [{ scenario: item.scenario, line: item.line, dryRun: item.dryRun, partial: item.partial }],
    },
  }))
}

export function registerGherkinCodeLens(monacoInstance: typeof Monaco, handlers: RunCodeLensHandlers) {
  monacoApi = monacoInstance
  activeHandlers = handlers
  if (providerRegistered) return
  providerRegistered = true

  monacoInstance.editor.registerCommand(
    RUN_SCENARIO_COMMAND,
    async (_accessor, payload: RunCodeLensPayload) => {
      if (!payload || !activeHandlers?.isEnabled()) return
      await activeHandlers.onRun(payload)
    },
  )

  monacoInstance.languages.registerCodeLensProvider('scenaria-feature', {
    provideCodeLenses(model) {
      if (!activeHandlers?.isEnabled() || !shouldUseHeavyLanguageFeatures(model.getLineCount())) {
        return { lenses: [], dispose: () => {} }
      }
      const lenses = toMonacoLenses(
        collectRunCodeLenses(model.getValue(), model.getVersionId(), model.uri.toString()),
      )
      return { lenses, dispose: () => {} }
    },
    resolveCodeLens(model, codeLens) {
      return codeLens
    },
  })
}

export function refreshGherkinCodeLens(editor: Monaco.editor.ICodeEditor | null) {
  if (!editor) return
  void editor.getAction('editor.action.codeLensRefresh')?.run()
}
