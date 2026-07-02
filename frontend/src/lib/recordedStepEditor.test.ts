import { describe, expect, it } from 'vitest'
import {
  findRecordedStepInsertLine,
  formatRecordedStepLine,
  pickRecordedStepKeyword,
  rebuildLiveRecordStepLines,
  removeLastRecordedStepFromText,
  upsertRecordedStepInText,
} from './recordedStepEditor'

describe('formatRecordedStepLine', () => {
  it('uses Допустим by default for bare step body', () => {
    expect(formatRecordedStepLine('нажимаю "#ok"')).toBe('\tДопустим нажимаю "#ok"')
  })

  it('uses И when requested', () => {
    expect(formatRecordedStepLine('нажимаю "#ok"', 'И')).toBe('\tИ нажимаю "#ok"')
  })

  it('strips existing keyword and applies chosen one', () => {
    expect(formatRecordedStepLine('Когда нажимаю "#ok"', 'И')).toBe('\tИ нажимаю "#ok"')
  })
})

describe('pickRecordedStepKeyword', () => {
  const emptyScenario = `Функционал: Demo
Сценарий: Main
`
  const withStep = `Функционал: Demo
Сценарий: Main
\tДопустим открыт "https://example.com"
`

  it('picks Допустим for first step in empty scenario', () => {
    expect(pickRecordedStepKeyword(emptyScenario, 0, {})).toBe('Допустим')
  })

  it('picks И when scenario already has steps', () => {
    expect(pickRecordedStepKeyword(withStep, 0, {})).toBe('И')
  })

  it('picks И for continuation index', () => {
    expect(pickRecordedStepKeyword(emptyScenario, 1, { 0: 2 })).toBe('И')
  })
})

describe('upsertRecordedStepInText', () => {
  const sample = `Функционал: Demo
Сценарий: Main
\tДопустим открыт "https://example.com"
`

  it('appends continuation with И after existing steps', () => {
    const { text, lineByIndex } = upsertRecordedStepInText(sample, 0, 'нажимаю "#go"', {})
    expect(text).toContain('\tИ нажимаю "#go"')
    expect(lineByIndex[0]).toBe(3)
  })

  it('uses Допустим for first step in empty scenario', () => {
    const empty = `Функционал: Demo
Сценарий: Main
`
    const { text } = upsertRecordedStepInText(empty, 0, 'открыт "https://example.com"', {})
    expect(text).toContain('\tДопустим открыт "https://example.com"')
  })

  it('updates step at same session index', () => {
    const first = upsertRecordedStepInText(sample, 0, 'ввожу "a" в "#email"', {})
    const second = upsertRecordedStepInText(first.text, 0, 'ввожу "ab" в "#email"', first.lineByIndex)
    const lines = second.text.split('\n')
    expect(lines[second.lineByIndex[0]]).toBe('\tИ ввожу "ab" в "#email"')
    expect(lines.filter((l) => l.includes('#email')).length).toBe(1)
  })
})

describe('rebuildLiveRecordStepLines', () => {
  it('maps step lines to session indices in order', () => {
    const text = `Функционал: X
Сценарий: Y
\tДопустим открыт "https://a.com"
\tИ нажимаю "#one"
\tИ открыт "https://b.com"
`
    expect(rebuildLiveRecordStepLines(text)).toEqual({ 0: 2, 1: 3, 2: 4 })
  })

  it('maps only the last N steps when stepCount is given', () => {
    const text = `Функционал: X
Сценарий: Y
\tДопустим шаблон
\tТогда шаблон
\tИ записан один
\tИ записан два
`
    expect(rebuildLiveRecordStepLines(text, 2)).toEqual({ 0: 4, 1: 5 })
  })
})

describe('findRecordedStepInsertLine', () => {
  it('returns last step line in first scenario', () => {
    const text = `Функционал: X
Сценарий: Y
\tДопустим шаг один
\tИ шаг два
`
    expect(findRecordedStepInsertLine(text)).toBe(3)
  })
})

describe('removeLastRecordedStepFromText', () => {
  it('removes highest session index and reindexes following lines', () => {
    const text = `Функционал: X
Сценарий: Y
\tДопустим шаг один
\tИ шаг два
\tИ шаг три
`
    const lineByIndex = { 0: 2, 1: 3, 2: 4 }
    const result = removeLastRecordedStepFromText(text, lineByIndex)
    expect(result?.text).toContain('шаг два')
    expect(result?.text).not.toContain('шаг три')
    expect(result?.lineByIndex).toEqual({ 0: 2, 1: 3 })
  })
})
