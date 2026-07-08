import { describe, expect, it } from 'vitest'
import { canonicalFeaturePath } from './featurePath'

describe('canonicalFeaturePath', () => {
  it('normalizes separators and trims outer whitespace', () => {
    expect(canonicalFeaturePath('  C:\\proj\\demo.feature  ')).toBe('C:/proj/demo.feature')
  })
})
