import type * as Monaco from 'monaco-editor'

const BLOCK_OPEN_RE = /^(если|повторяю|пока|для каждого)(?:\s|$)/i
const BLOCK_CLOSE_IF_RE = /^конец если(?:\s|$)/i
const BLOCK_CLOSE_RE = /^конец$/i

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

type BlockFrame = {
  startLine: number
  indent: number
  opener: string
}

export type FeatureFoldingRange = {
  start: number
  end: number
}

/** Folding ranges for control-flow blocks (Если / Повторяю / Пока / Для каждого). */
export function collectBlockFoldingRanges(text: string): FeatureFoldingRange[] {
  const lines = text.split('\n')
  const ranges: FeatureFoldingRange[] = []
  const stack: BlockFrame[] = []

  const closeFrame = (lineNo: number, indent: number, matcher: (opener: string) => boolean) => {
    for (let i = stack.length - 1; i >= 0; i--) {
      const frame = stack[i]
      if (indent < frame.indent) {
        continue
      }
      if (!matcher(frame.opener)) {
        continue
      }
      if (lineNo > frame.startLine) {
        ranges.push({ start: frame.startLine, end: lineNo })
      }
      stack.splice(i, 1)
      return
    }
  }

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]
    const trimmed = line.trimStart()
    if (!trimmed || trimmed.startsWith('#')) {
      continue
    }
    const indent = indentWidth(leadingIndent(line))
    const lineNo = i + 1

    const open = trimmed.match(BLOCK_OPEN_RE)
    if (open) {
      stack.push({ startLine: lineNo, indent, opener: open[1].toLowerCase() })
      continue
    }
    if (BLOCK_CLOSE_IF_RE.test(trimmed)) {
      closeFrame(lineNo, indent, (opener) => opener === 'если')
      continue
    }
    if (BLOCK_CLOSE_RE.test(trimmed)) {
      closeFrame(lineNo, indent, (opener) => opener === 'повторяю' || opener === 'пока' || opener === 'для каждого')
    }
  }

  return ranges
}

let foldingRegistered = false

export function registerGherkinFolding(monaco: typeof Monaco) {
  if (foldingRegistered) {
    return
  }
  foldingRegistered = true

  monaco.languages.registerFoldingRangeProvider('scenaria-feature', {
    provideFoldingRanges(model) {
      const ranges = collectBlockFoldingRanges(model.getValue())
      return ranges.map((range) => ({
        start: range.start,
        end: range.end,
        kind: monaco.languages.FoldingRangeKind.Region,
      }))
    },
  })
}
