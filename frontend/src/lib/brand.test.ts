import { describe, expect, it } from 'vitest'
import { BRAND_NAME, BRAND_TAGLINE, brandAboutText, brandOverlayTitle } from './brand'

describe('brand', () => {
  it('matches Python brand constants', () => {
    expect(BRAND_NAME).toBe('Scenaria')
    expect(BRAND_TAGLINE).toContain('Gherkin')
  })

  it('builds about text with version', () => {
    expect(brandAboutText('0.15.0')).toContain('Версия 0.15.0')
    expect(brandAboutText('0.15.0')).toContain(BRAND_TAGLINE)
  })

  it('builds overlay title', () => {
    expect(brandOverlayTitle()).toContain(BRAND_NAME)
  })
})
