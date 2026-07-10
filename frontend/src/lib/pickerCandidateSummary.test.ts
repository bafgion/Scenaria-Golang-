import { describe, expect, it } from 'vitest'
import { pickerCandidateFacts } from './pickerCandidateSummary'

describe('pickerCandidateFacts', () => {
  it('returns only available match and score values', () => {
    expect(
      pickerCandidateFacts({
        matches_count: 3,
        score: 87,
        unique: true,
        visible: false,
        warnings: ['Too generic'],
      } as never),
    ).toEqual({
      matchesCount: 3,
      score: 87,
      unique: true,
      visible: false,
      warnings: ['Too generic'],
    })
  })

  it('omits zero score and zero matches without inventing values', () => {
    expect(
      pickerCandidateFacts({
        matches_count: 0,
        score: 0,
        unique: false,
        visible: true,
        warnings: [],
      } as never),
    ).toEqual({
      matchesCount: undefined,
      score: undefined,
      unique: false,
      visible: true,
      warnings: [],
    })
  })
})
