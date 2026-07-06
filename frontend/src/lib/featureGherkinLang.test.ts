import { describe, expect, it } from 'vitest'
import { defaultStepKeyword, detectFeatureGherkinLanguage } from './featureGherkinLang'
import { formatInsertText } from './gherkinCompletions'

describe('featureGherkinLang', () => {
  it('detects en tag', () => {
    expect(detectFeatureGherkinLanguage('# language: en\nFeature: X')).toBe('en')
    expect(detectFeatureGherkinLanguage('Feature: X')).toBe('ru')
  })

  it('uses When for en step inserts', () => {
    expect(defaultStepKeyword('en')).toBe('When')
    const out = formatInsertText('\t', { label: 'I open', insert: 'I open "https://x.com"', description: '' }, 'en')
    expect(out).toBe('When I open "https://x.com"')
  })
})
