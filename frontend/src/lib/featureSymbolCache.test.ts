import { describe, expect, it } from 'vitest'
import { clearFeatureSymbolCache, getCachedFeatureSymbols } from './featureSymbolCache'
import { parseFeatureSymbols } from './gherkinDocumentSymbols'

describe('featureSymbolCache', () => {
  const sample = `Функционал: Demo
  Сценарий: A
    Допустим открыт "https://example.com"
`

  it('returns same reference for monaco version id', () => {
    clearFeatureSymbolCache()
    const a = getCachedFeatureSymbols(sample, 3)
    const b = getCachedFeatureSymbols(sample, 3)
    expect(a).toBe(b)
    expect(a).toEqual(parseFeatureSymbols(sample))
  })

  it('recomputes when version id changes', () => {
    clearFeatureSymbolCache()
    const v1 = getCachedFeatureSymbols(sample, 1)
    const v2 = getCachedFeatureSymbols(sample, 2)
    expect(v1).not.toBe(v2)
    expect(v1).toEqual(v2)
  })

  it('caches by text hash without version id', () => {
    clearFeatureSymbolCache()
    const a = getCachedFeatureSymbols(sample)
    const b = getCachedFeatureSymbols(sample)
    expect(a).toBe(b)
  })
})
