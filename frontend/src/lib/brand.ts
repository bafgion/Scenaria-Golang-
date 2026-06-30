/** Product branding (Python app/brand.py parity). */
export const BRAND_NAME = 'Scenaria'
export const BRAND_TAGLINE = 'Автотесты сайтов · Gherkin · Playwright'
export const BRAND_DESCRIPTION = 'Запись и запуск автотестов сайтов в браузере.'
export const BRAND_TITLE = BRAND_NAME.toUpperCase()

export function brandAboutText(version: string): string {
  const ver = version.trim() || 'dev'
  return `${BRAND_DESCRIPTION}\n${BRAND_TAGLINE}\n\nВерсия ${ver}`
}

export function brandOverlayTitle(): string {
  return `${BRAND_NAME} — перетащите панель`
}
