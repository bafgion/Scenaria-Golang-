import { t } from './i18n'

/** Product branding (Python app/brand.py parity). */
export const BRAND_NAME = 'Scenaria'
export const BRAND_TITLE = BRAND_NAME.toUpperCase()

export function brandTagline(): string {
  return t('brand.tagline')
}

export function brandDescription(): string {
  return t('brand.description')
}

export function brandAboutText(version: string): string {
  const ver = version.trim() || 'dev'
  return `${t('brand.description')}\n${t('brand.tagline')}\n\n${t('brand.aboutVersion', { version: ver })}`
}

export function brandOverlayTitle(): string {
  return t('brand.overlayTitle', { name: BRAND_NAME })
}
