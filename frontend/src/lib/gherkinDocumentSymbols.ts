import type * as Monaco from 'monaco-editor'
import { shouldUseHeavyLanguageFeatures } from './editorLargeFile'
import { getCachedFeatureSymbols } from './featureSymbolCache'

export type FeatureSymbolKind =
  | 'feature'
  | 'context'
  | 'scenario'
  | 'outline'
  | 'examples'
  | 'tag'
  | 'step'
  | 'block'

export type FeatureSymbol = {
  name: string
  detail?: string
  kind: FeatureSymbolKind
  line: number
  children: FeatureSymbol[]
}

export type FlatFeatureSymbol = FeatureSymbol & { depth: number }

const STRUCTURE_HEADER_RE =
  /^(функциональность|функционал|функция|feature|сценарий|scenario|структура сценария|scenario outline|примеры|examples|контекст|background)\s*:\s*(.*)$/i
const BLOCK_LINE_RE = /^(если|повторяю|пока|для каждого|иначе|конец если|конец)(?:\s|$)/i
const STEP_KEYWORD_RE = /^(Допустим|Дано|Когда|Тогда|И|Но|Given|When|Then|And|But)\s+(.*)$/i

function leadingIndent(line: string): string {
  const trimmed = line.trimStart()
  return line.slice(0, line.length - trimmed.length)
}

function indentWidth(indent: string): number {
  let width = 0
  for (const ch of indent) {
    if (ch === '\t') width += 4
    else if (ch === ' ') width += 1
  }
  return width
}

function isStepIndented(line: string): boolean {
  const indent = leadingIndent(line)
  if (indent.startsWith('\t')) return true
  return indent.length >= 2 && indent.trim() === ''
}

function isTagLine(line: string): boolean {
  const parts = line.trim().split(/\s+/)
  return parts.length > 0 && parts.every((part) => part.startsWith('@'))
}

function isCommentOrEmpty(line: string): boolean {
  const trimmed = line.trim()
  return trimmed === '' || trimmed.startsWith('#')
}

function headerKind(label: string): FeatureSymbolKind {
  const lower = label.toLowerCase()
  if (lower === 'контекст' || lower === 'background') return 'context'
  if (lower === 'структура сценария' || lower === 'scenario outline') return 'outline'
  if (lower === 'примеры' || lower === 'examples') return 'examples'
  if (lower === 'сценарий' || lower === 'scenario') return 'scenario'
  return 'feature'
}

function headerName(label: string, title: string, line: string): string {
  const trimmedTitle = title.trim()
  if (trimmedTitle) return trimmedTitle
  if (label.toLowerCase() === 'примеры' || label.toLowerCase() === 'examples') return 'Примеры'
  return line.trim()
}

function parseStepLine(line: string): { keyword: string; text: string } | null {
  const trimmed = line.trimStart()
  const match = trimmed.match(STEP_KEYWORD_RE)
  if (match) {
    return { keyword: match[1], text: match[2].trim() }
  }
  if (isStepIndented(line)) {
    return { keyword: '', text: trimmed }
  }
  return null
}

type BlockFrame = { symbol: FeatureSymbol; indent: number }

export function parseFeatureSymbols(text: string): FeatureSymbol[] {
  const lines = text.split('\n')
  const roots: FeatureSymbol[] = []
  let feature: FeatureSymbol | null = null
  let section: FeatureSymbol | null = null
  const blockStack: BlockFrame[] = []
  let pendingTags: FeatureSymbol[] = []

  const parentForIndent = (indent: number): FeatureSymbol | null => {
    while (blockStack.length > 0 && blockStack[blockStack.length - 1].indent >= indent) {
      blockStack.pop()
    }
    if (blockStack.length > 0) return blockStack[blockStack.length - 1].symbol
    if (section) return section
    if (feature) return feature
    return null
  }

  const attachToFeatureOrRoots = (node: FeatureSymbol) => {
    if (feature) {
      feature.children.push(node)
      return
    }
    roots.push(node)
  }

  const attach = (parent: FeatureSymbol | null, node: FeatureSymbol) => {
    if (parent) {
      parent.children.push(node)
      return
    }
    roots.push(node)
  }

  for (let i = 0; i < lines.length; i++) {
    const lineNo = i + 1
    const line = lines[i]
    if (isCommentOrEmpty(line)) continue

    const stripped = line.trimStart()
    const indent = indentWidth(leadingIndent(line))

    if (isTagLine(stripped)) {
      pendingTags.push({ name: stripped, kind: 'tag', line: lineNo, children: [] })
      continue
    }

    const headerMatch = stripped.match(STRUCTURE_HEADER_RE)
    if (headerMatch) {
      blockStack.length = 0
      const label = headerMatch[1]
      const kind = headerKind(label)
      const node: FeatureSymbol = {
        name: headerName(label, headerMatch[2], stripped),
        detail: label,
        kind,
        line: lineNo,
        children: [],
      }

      if (kind === 'feature') {
        feature = node
        section = null
        roots.push(node)
      } else if (kind === 'context') {
        section = node
        attachToFeatureOrRoots(node)
      } else if (kind === 'scenario' || kind === 'outline') {
        section = node
        node.children.push(...pendingTags)
        pendingTags = []
        attachToFeatureOrRoots(node)
      } else if (kind === 'examples' && section) {
        attach(section, node)
      } else {
        attachToFeatureOrRoots(node)
      }
      continue
    }

    if (BLOCK_LINE_RE.test(stripped)) {
      const parent = parentForIndent(indent)
      const node: FeatureSymbol = {
        name: stripped.length > 72 ? `${stripped.slice(0, 69)}…` : stripped,
        kind: 'block',
        line: lineNo,
        children: [],
      }
      attach(parent, node)
      blockStack.push({ symbol: node, indent })
      continue
    }

    const step = parseStepLine(line)
    if (step && (step.text || step.keyword)) {
      const parent = parentForIndent(indent)
      const label = step.text || step.keyword
      attach(parent, {
        name: label.length > 72 ? `${label.slice(0, 69)}…` : label,
        detail: step.keyword || undefined,
        kind: 'step',
        line: lineNo,
        children: [],
      })
    }
  }

  return roots
}

export function flattenFeatureSymbols(symbols: FeatureSymbol[], depth = 0): FlatFeatureSymbol[] {
  const out: FlatFeatureSymbol[] = []
  for (const symbol of symbols) {
    out.push({ ...symbol, depth })
    if (symbol.children.length > 0) {
      out.push(...flattenFeatureSymbols(symbol.children, depth + 1))
    }
  }
  return out
}

function symbolKindForMonaco(
  kind: FeatureSymbolKind,
  monaco: typeof Monaco,
): Monaco.languages.SymbolKind {
  switch (kind) {
    case 'feature':
      return monaco.languages.SymbolKind.Module
    case 'context':
      return monaco.languages.SymbolKind.Namespace
    case 'scenario':
    case 'outline':
      return monaco.languages.SymbolKind.Class
    case 'examples':
      return monaco.languages.SymbolKind.Interface
    case 'tag':
      return monaco.languages.SymbolKind.Key
    case 'block':
      return monaco.languages.SymbolKind.Event
    case 'step':
    default:
      return monaco.languages.SymbolKind.Method
  }
}

function toDocumentSymbol(
  node: FeatureSymbol,
  lines: string[],
  monaco: typeof Monaco,
): Monaco.languages.DocumentSymbol {
  const lineText = lines[node.line - 1] ?? node.name
  const range = {
    startLineNumber: node.line,
    startColumn: 1,
    endLineNumber: node.line,
    endColumn: Math.max(2, lineText.length + 1),
  }
  return {
    name: node.name,
    detail: node.detail ?? '',
    kind: symbolKindForMonaco(node.kind, monaco),
    tags: [],
    range,
    selectionRange: range,
    children: node.children.map((child) => toDocumentSymbol(child, lines, monaco)),
  }
}

export function toMonacoDocumentSymbols(
  text: string,
  monaco: typeof Monaco,
  versionId?: number,
  modelKey?: string,
): Monaco.languages.DocumentSymbol[] {
  const lines = text.split('\n')
  return getCachedFeatureSymbols(text, versionId, modelKey).map((node) =>
    toDocumentSymbol(node, lines, monaco),
  )
}

let documentSymbolsRegistered = false

export function registerGherkinDocumentSymbols(monaco: typeof Monaco) {
  if (documentSymbolsRegistered) {
    return
  }
  documentSymbolsRegistered = true

  monaco.languages.registerDocumentSymbolProvider('scenaria-feature', {
    provideDocumentSymbols(model) {
      if (!shouldUseHeavyLanguageFeatures(model.getLineCount())) {
        return []
      }
      return toMonacoDocumentSymbols(
        model.getValue(),
        monaco,
        model.getVersionId(),
        model.uri.toString(),
      )
    },
  })
}
