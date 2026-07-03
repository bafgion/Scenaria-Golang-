import { ru, type RuMessages } from './locales/ru'

export type Locale = 'ru'

/** Nested message tree for the active locale. Add `en.ts` + switch when translating. */
export type Messages = RuMessages

const catalogs: Record<Locale, Messages> = {
  ru,
}

let activeLocale: Locale = 'ru'

export function getLocale(): Locale {
  return activeLocale
}

export function setLocale(locale: Locale) {
  activeLocale = locale
}

function resolvePath(tree: Record<string, unknown>, path: string): string | undefined {
  const parts = path.split('.')
  let node: unknown = tree
  for (const part of parts) {
    if (!node || typeof node !== 'object' || !(part in (node as object))) return undefined
    node = (node as Record<string, unknown>)[part]
  }
  return typeof node === 'string' ? node : undefined
}

function interpolate(template: string, params?: Record<string, string | number>): string {
  if (!params) return template
  return template.replace(/\{(\w+)\}/g, (_, key: string) => String(params[key] ?? `{${key}}`))
}

/** Translate a dot-path key for the active locale. */
export function t(key: string, params?: Record<string, string | number>): string {
  const catalog = catalogs[activeLocale]
  const value = resolvePath(catalog as unknown as Record<string, unknown>, key)
  if (value === undefined) return key
  return interpolate(value, params)
}
