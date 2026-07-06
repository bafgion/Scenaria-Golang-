import { describe, expect, it } from 'vitest'
import { setLocale } from './i18n'
import { buildFeatureTemplate } from './featureTemplate'

describe('buildFeatureTemplate', () => {
  it('always uses Russian Gherkin with # language: ru', () => {
    setLocale('ru')
    const text = buildFeatureTemplate({
      title: 'UI',
      scenario: 'Smoke',
      startUrl: 'https://store.test',
    })
    expect(text).toContain('https://store.test')
    expect(text).toContain('Сценарий: Smoke')
    expect(text).toContain('Функционал: UI')
    expect(text).toContain('# language: ru')
    expect(text).not.toContain('# language: en')
  })

  it('ignores English UI locale for new feature scaffold', () => {
    setLocale('en')
    const text = buildFeatureTemplate({
      title: 'UI',
      scenario: 'Smoke',
      startUrl: 'https://store.test',
    })
    expect(text).toContain('Сценарий: Smoke')
    expect(text).toContain('# language: ru')
    expect(text).not.toContain('Scenario:')
    expect(text).not.toContain('# language: en')
  })
})
