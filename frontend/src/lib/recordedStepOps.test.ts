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
