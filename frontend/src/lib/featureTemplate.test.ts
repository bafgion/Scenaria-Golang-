import { describe, expect, it } from 'vitest'
import { setLocale } from './i18n'
import { buildFeatureTemplate } from './featureTemplate'

describe('buildFeatureTemplate', () => {
  it('includes start URL and scenario name for Russian locale', () => {
    setLocale('ru')
    const text = buildFeatureTemplate({
      title: 'UI',
      scenario: 'Smoke',
      startUrl: 'https://store.test',
    })
    expect(text).toContain('https://store.test')
    expect(text).toContain('Сценарий: Smoke')
    expect(text).toContain('# language: ru')
  })

  it('uses English Gherkin keywords for English locale', () => {
    setLocale('en')
    const text = buildFeatureTemplate({
      title: 'UI',
      scenario: 'Smoke',
      startUrl: 'https://store.test',
    })
    expect(text).toContain('Scenario: Smoke')
    expect(text).toContain('# language: en')
  })
})
