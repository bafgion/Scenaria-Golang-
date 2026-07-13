import {
  findRecordedStepInsertLine,
  formatRecordedStepLine,
  rebuildLiveRecordStepLines,
  stripStepKeyword,
} from './recordedStepEditor'

export type RecordStepOp = 'upsert' | 'delete' | 'reset' | 'snapshot'

export type RecordStepEvent = {
  op: RecordStepOp
  index?: number
  line?: string
  lines?: string[]
}

function isFormattedGherkinStep(line: string): boolean {
  return /^\t(Допустим|Дано|Когда|Тогда|И|Но|Given|When|Then|And|But)\s+/i.test(line)
}

function normalizeIncomingStepLine(line: string): string {
  const trimmed = line.trim()
  if (!trimmed) return ''
  if (isFormattedGherkinStep(line)) return line
  return formatRecordedStepLine(trimmed)
}

/** Apply a backend record-step op to editor text. */
export function applyRecordStepEvent(
  text: string,
  event: RecordStepEvent,
  lineByIndex: Record<number, number>,
): { text: string; lineByIndex: Record<number, number> } {
  switch (event.op) {
    case 'reset':
      return { text, lineByIndex: {} }
    case 'delete': {
      const index = event.index ?? -1
      const lineNo = lineByIndex[index]
      if (lineNo === undefined || lineNo < 0) {
        return { text, lineByIndex: { ...lineByIndex } }
      }
      const lines = text.split('\n')
      if (lineNo >= lines.length) {
        const nextMap = { ...lineByIndex }
        delete nextMap[index]
        return { text, lineByIndex: nextMap }
      }
      lines.splice(lineNo, 1)
      const nextMap: Record<number, number> = {}
      for (const [key, mappedLine] of Object.entries(lineByIndex)) {
        const idx = Number(key)
        if (idx === index) continue
        nextMap[idx] = mappedLine > lineNo ? mappedLine - 1 : mappedLine
      }
      return { text: lines.join('\n'), lineByIndex: nextMap }
    }
    case 'snapshot': {
      const incoming = (event.lines ?? []).map(normalizeIncomingStepLine).filter(Boolean)
      if (incoming.length === 0) {
        return { text, lineByIndex: {} }
      }
      const lines = text.split('\n')
      const insertAfter = findRecordedStepInsertLine(text)
      const existing = Object.entries(lineByIndex)
        .map(([key, lineNo]) => ({ index: Number(key), lineNo }))
        .filter((entry) => Number.isFinite(entry.index) && entry.lineNo >= 0)
        .sort((a, b) => a.index - b.index)

      if (existing.length > incoming.length) {
        return { text, lineByIndex: { ...lineByIndex } }
      }

      if (existing.length > 0) {
        for (const entry of [...existing].sort((a, b) => b.lineNo - a.lineNo)) {
          if (entry.lineNo < lines.length) {
            lines.splice(entry.lineNo, 1)
          }
        }
      }

      const startAt = insertAfter + 1
      const rebuilt = [...lines.slice(0, startAt), ...incoming, ...lines.slice(startAt)]
      return {
        text: rebuilt.join('\n'),
        lineByIndex: rebuildLiveRecordStepLines(rebuilt.join('\n'), incoming.length),
      }
    }
    case 'upsert':
    default: {
      const index = event.index ?? 0
      const formatted = normalizeIncomingStepLine(event.line ?? '')
      if (!formatted) return { text, lineByIndex: { ...lineByIndex } }
      const lines = text.split('\n')
      const nextMap = { ...lineByIndex }
      const existingLine = nextMap[index]
      if (existingLine !== undefined && existingLine >= 0 && existingLine < lines.length) {
        lines[existingLine] = formatted
        return { text: lines.join('\n'), lineByIndex: nextMap }
      }
      const insertAfter = findRecordedStepInsertLine(text)
      const insertAt = insertAfter + 1
      lines.splice(insertAt, 0, formatted)
      nextMap[index] = insertAt
      const adjusted: Record<number, number> = {}
      for (const [key, lineNo] of Object.entries(nextMap)) {
        const idx = Number(key)
        adjusted[idx] = lineNo >= insertAt && idx !== index ? lineNo + 1 : lineNo
      }
      return { text: lines.join('\n'), lineByIndex: adjusted }
    }
  }
}

/** @deprecated Use backend-formatted lines via applyRecordStepEvent. */
export function legacyBareStepBody(line: string): string {
  return stripStepKeyword(line)
}
