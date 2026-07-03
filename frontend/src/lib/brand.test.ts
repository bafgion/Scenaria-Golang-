import { describe, expect, it } from 'vitest'
import { t } from './i18n'
import { BRAND_NAME, brandAboutText, brandOverlayTitle, brandTagline } from './brand'

describe('brand', () => {
  it('matches Python brand constants', () => {
    expect(BRAND_NAME).toBe('Scenaria')
    expect(brandTagline()).toContain('Gherkin')
  })

  it('builds about text with version', () => {
    expect(brandAboutText('0.15.0')).toContain(t('brand.aboutVersion', { version: '0.15.0' }))
    expect(brandAboutText('0.15.0')).toContain(t('brand.tagline'))
  })

  it('builds overlay title', () => {
    expect(brandOverlayTitle()).toContain(BRAND_NAME)
  })
})
