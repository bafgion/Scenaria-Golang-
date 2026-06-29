import { describe, expect, it } from 'vitest'
import { buildFeatureTemplate } from './featureTemplate'

describe('buildFeatureTemplate', () => {
  it('includes start URL and scenario name', () => {
    const text = buildFeatureTemplate({
      title: 'UI',
      scenario: 'Smoke',
      startUrl: 'https://store.test',
    })
    expect(text).toContain('https://store.test')
    expect(text).toContain('Сценарий: Smoke')
  })
})
