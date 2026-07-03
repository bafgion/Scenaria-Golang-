import { derived, get, writable } from 'svelte/store'
import { en } from './locales/en'
import { ru, type Messages } from './locales/ru'

export type Locale = 'ru' | 'en'

const catalogs: Record<Locale, Messages> = {
  ru,
  en: en as Messages,
}

export const locale = writable<Locale>(detectDefaultLocale())

export const messages = derived(locale, ($locale) => catalogs[$locale])

function detectDefaultLocale(): Locale {
  if (typeof navigator === 'undefined') return 'ru'
  const lang = navigator.language?.toLowerCase() ?? ''
  if (lang.startsWith('en')) return 'en'
  return 'ru'
}

export function getLocale(): Locale {
  return get(locale)
}

export function setLocale(next: Locale) {
  if (next !== 'ru' && next !== 'en') return
  locale.set(next)
  if (typeof document !== 'undefined') {
    document.documentElement.lang = next
  }
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
export function t(key: string, params?: Record<string, string | number>, catalog?: Messages): string {
  const tree = (catalog ?? catalogs[getLocale()]) as unknown as Record<string, unknown>
  const value = resolvePath(tree, key)
  if (value === undefined) return key
  return interpolate(value, params)
}

/** Reactive translate — use in Svelte: `$tr('menus.save')` */
export function createTranslator($locale: Locale) {
  return (key: string, params?: Record<string, string | number>) => t(key, params, catalogs[$locale])
}
