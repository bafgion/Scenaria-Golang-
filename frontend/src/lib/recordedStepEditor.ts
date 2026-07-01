const STEP_KEYWORD_RE = /^(Допустим|Дано|Когда|Тогда|И|Но|Given|When|Then|And|But)\s+/i
const SCENARIO_HEADER_RE = /^(сценарий|scenario|структура сценария|scenario outline)\s*:/i
const BLOCK_HEADER_RE =
  /^(функциональность|функционал|функция|feature|примеры|examples|контекст|background)\s*:/i
const STEP_LINE_RE = /^\t|^ {2,}/

export type RecordedStepKeyword = 'Допустим' | 'И'

export function stripStepKeyword(body: string): string {
  return body.trim().replace(STEP_KEYWORD_RE, '').trim()
}

export function countScenarioSteps(text: string): number {
  const lines = text.split('\n')
  let inScenario = false
  let count = 0
  for (const line of lines) {
    const trimmed = line.trim()
    if (SCENARIO_HEADER_RE.test(trimmed)) {
      inScenario = true
      continue
    }
    if (inScenario && BLOCK_HEADER_RE.test(trimmed)) break
    if (inScenario && isStepLine(line)) count++
  }
  return count
}

function firstScenarioStepLineIndex(text: string): number {
  const lines = text.split('\n')
  let inScenario = false
  for (let i = 0; i < lines.length; i++) {
    const trimmed = lines[i].trim()
    if (SCENARIO_HEADER_RE.test(trimmed)) {
      inScenario = true
      continue
    }
    if (inScenario && BLOCK_HEADER_RE.test(trimmed)) break
    if (inScenario && isStepLine(lines[i])) return i
  }
  return -1
}

/** Pick Gherkin keyword for a live-recorded step (Python-style: Допустим then И). */
export function pickRecordedStepKeyword(
  text: string,
  sessionIndex: number,
  lineByIndex: Record<number, number>,
): RecordedStepKeyword {
  const existingLine = lineByIndex[sessionIndex]
  const firstStepLine = firstScenarioStepLineIndex(text)
  if (existingLine !== undefined && existingLine >= 0 && existingLine === firstStepLine) {
    return 'Допустим'
  }
  const stepsBefore = countScenarioSteps(text)
  if (stepsBefore === 0 && sessionIndex === 0) {
    return 'Допустим'
  }
  return 'И'
}

export function formatRecordedStepLine(
  body: string,
  keyword: RecordedStepKeyword = 'Допустим',
): string {
  const bare = stripStepKeyword(body)
  if (!bare) return `\t${keyword} выполняю действие`
  return `\t${keyword} ${bare}`
}

function isStepLine(line: string): boolean {
  if (STEP_KEYWORD_RE.test(line.trim())) return true
  return STEP_LINE_RE.test(line) && line.trim().length > 0
}

/** Map live-record session indices to editor line numbers from existing text. */
export function rebuildLiveRecordStepLines(text: string, stepCount?: number): Record<number, number> {
  const lines = text.split('\n')
  const stepLines: number[] = []
  for (let i = 0; i < lines.length; i++) {
    if (isStepLine(lines[i])) {
      stepLines.push(i)
    }
  }
  const mapped =
    stepCount !== undefined && stepCount > 0 && stepCount < stepLines.length
      ? stepLines.slice(-stepCount)
      : stepLines
  const map: Record<number, number> = {}
  mapped.forEach((lineNo, index) => {
    map[index] = lineNo
  })
  return map
}

/** Line index after which recorded steps should be inserted (last step or scenario header). */
export function findRecordedStepInsertLine(text: string): number {
  const lines = text.split('\n')
  let scenarioLine = -1
  let lastStep = -1
  for (let i = 0; i < lines.length; i++) {
    const trimmed = lines[i].trim()
    if (!trimmed || trimmed.startsWith('#')) continue
    if (SCENARIO_HEADER_RE.test(trimmed)) {
      scenarioLine = i
      lastStep = i
      continue
    }
    if (scenarioLine >= 0 && BLOCK_HEADER_RE.test(trimmed) && i > scenarioLine) {
      break
    }
    if (scenarioLine >= 0 && isStepLine(lines[i])) {
      lastStep = i
    }
  }
  if (lastStep >= 0) return lastStep
  if (scenarioLine >= 0) return scenarioLine
  return Math.max(0, lines.length - 1)
}

/** Insert or update a live-recorded step line by session index. */
export function upsertRecordedStepInText(
  text: string,
  sessionIndex: number,
  stepBody: string,
  lineByIndex: Record<number, number>,
): { text: string; lineByIndex: Record<number, number> } {
  const keyword = pickRecordedStepKeyword(text, sessionIndex, lineByIndex)
  const formatted = formatRecordedStepLine(stepBody, keyword)
  const lines = text.split('\n')
  const nextMap = { ...lineByIndex }

  const existingLine = nextMap[sessionIndex]
  if (existingLine !== undefined && existingLine >= 0 && existingLine < lines.length) {
    lines[existingLine] = formatted
    return { text: lines.join('\n'), lineByIndex: nextMap }
  }

  const insertAfter = findRecordedStepInsertLine(text)
  const insertAt = insertAfter + 1
  lines.splice(insertAt, 0, formatted)
  nextMap[sessionIndex] = insertAt

  const adjusted: Record<number, number> = {}
  for (const [key, lineNo] of Object.entries(nextMap)) {
    const idx = Number(key)
    adjusted[idx] = lineNo >= insertAt && idx !== sessionIndex ? lineNo + 1 : lineNo
  }

  return { text: lines.join('\n'), lineByIndex: adjusted }
}
