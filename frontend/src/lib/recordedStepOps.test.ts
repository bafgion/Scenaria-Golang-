import { describe, expect, it } from 'vitest'
import { applyRecordStepEvent } from './recordedStepOps'

const template = `Функциональность: Demo

  Сценарий: Test
`

describe('recordedStepOps', () => {
  it('applies upsert with backend-formatted line', () => {
    const first = applyRecordStepEvent(template, {
      op: 'upsert',
      index: 0,
      line: '\tДопустим нажимаю "#one"',
    }, {})
    expect(first.text).toContain('\tДопустим нажимаю "#one"')
    const second = applyRecordStepEvent(first.text, {
      op: 'upsert',
      index: 1,
      line: '\tИ нажимаю "#two"',
    }, first.lineByIndex)
    expect(second.text).toContain('\tИ нажимаю "#two"')
  })

  it('applies snapshot replay', () => {
    const result = applyRecordStepEvent(template, {
      op: 'snapshot',
      lines: ['\tДопустим открыт "https://example.com"', '\tИ нажимаю "#btn"'],
    }, {})
    expect(result.lineByIndex[0]).toBeGreaterThanOrEqual(0)
    expect(result.lineByIndex[1]).toBeGreaterThan(result.lineByIndex[0])
  })

  it('does not shrink live-recorded steps when a shorter stop snapshot arrives', () => {
    const first = applyRecordStepEvent(template, {
      op: 'upsert',
      index: 0,
      line: '\tGiven opened "https://www.2moodstore.com/"',
    }, {})
    const second = applyRecordStepEvent(first.text, {
      op: 'upsert',
      index: 1,
      line: '\tAnd click "#catalog"',
    }, first.lineByIndex)
    const snapshot = applyRecordStepEvent(second.text, {
      op: 'snapshot',
      lines: ['\tGiven opened "https://www.2moodstore.com/"'],
    }, second.lineByIndex)

    expect(snapshot.text).toContain('https://www.2moodstore.com/')
    expect(snapshot.text).toContain('#catalog')
    expect(snapshot.lineByIndex).toEqual(second.lineByIndex)
  })

  it('replaces mapped live steps with a longer authoritative snapshot without duplicating old lines', () => {
    const first = applyRecordStepEvent(template, {
      op: 'upsert',
      index: 0,
      line: '\tGiven opened "https://example.com"',
    }, {})
    const result = applyRecordStepEvent(first.text, {
      op: 'snapshot',
      lines: ['\tGiven opened "https://example.com"', '\tAnd click "#btn"'],
    }, first.lineByIndex)

    expect(result.text.match(/https:\/\/example\.com/g)?.length).toBe(1)
    expect(result.text).toContain('#btn')
    expect(Object.keys(result.lineByIndex)).toHaveLength(2)
  })

  it('reset clears index map only', () => {
    const mapped = applyRecordStepEvent(template, {
      op: 'upsert',
      index: 0,
      line: '\tДопустим нажимаю "#one"',
    }, {})
    const reset = applyRecordStepEvent(mapped.text, { op: 'reset' }, mapped.lineByIndex)
    expect(reset.lineByIndex).toEqual({})
    expect(reset.text).toBe(mapped.text)
  })

  it('delete removes mapped line', () => {
    const mapped = applyRecordStepEvent(template, {
      op: 'upsert',
      index: 0,
      line: '\tДопустим нажимаю "#one"',
    }, {})
    const deleted = applyRecordStepEvent(mapped.text, { op: 'delete', index: 0 }, mapped.lineByIndex)
    expect(deleted.text).not.toContain('#one')
    expect(deleted.lineByIndex).toEqual({})
  })
})
