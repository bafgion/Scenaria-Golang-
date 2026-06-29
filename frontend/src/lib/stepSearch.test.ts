import { describe, expect, it } from 'vitest'
import { asStepSearchQuery } from './stepSearch'

describe('asStepSearchQuery', () => {
  it('returns string queries unchanged', () => {
    expect(asStepSearchQuery('нажимаю')).toBe('нажимаю')
    expect(asStepSearchQuery('')).toBe('')
  })

  it('ignores accidental event objects', () => {
    expect(asStepSearchQuery({ isTrusted: true })).toBe('')
    expect(asStepSearchQuery(null)).toBe('')
    expect(asStepSearchQuery(42)).toBe('')
  })
})
