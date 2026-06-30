import { isUntitled } from './untitled'

export type TabBody = {
  path: string
  content: string
  dirty: boolean
  draft?: string
  /** Текст выгружен из RAM; при активации читаем с диска. */
  unloaded?: boolean
}

/** Сколько сохранённых на диске вкладок держим в памяти одновременно. */
export const MAX_RETAINED_CLEAN_TAB_BODIES = 8

/** Мягкий лимит открытых вкладок (предупреждение в журнале). */
export const MAX_OPEN_EDITOR_TABS = 40

export function tabEditorText(tab: TabBody): string {
  if (tab.dirty) return tab.draft ?? tab.content
  if (tab.unloaded) return ''
  return tab.content
}

export function tabNeedsDiskReload(tab: TabBody): boolean {
  return !tab.dirty && !!tab.unloaded
}

/**
 * Выгружает текст «чистых» неактивных вкладок, оставляя в RAM не больше maxClean тел.
 * Грязные и untitled не трогаем.
 */
export function trimRetainedTabBodies(
  tabs: TabBody[],
  activePath: string,
  maxClean = MAX_RETAINED_CLEAN_TAB_BODIES,
): TabBody[] {
  const cleanWithBody = tabs.filter(
    (t) => t.path !== activePath && !t.dirty && !isUntitled(t.path) && !t.unloaded && t.content.length > 0,
  )
  if (cleanWithBody.length <= maxClean) {
    return tabs
  }
  const dropCount = cleanWithBody.length - maxClean
  const dropPaths = new Set(cleanWithBody.slice(0, dropCount).map((t) => t.path))
  return tabs.map((t) => {
    if (!dropPaths.has(t.path)) return t
    return { ...t, content: '', unloaded: true }
  })
}

export function pathsToRetainModels(tabs: TabBody[], activePath: string): string[] {
  const paths = tabs.filter((t) => !t.unloaded || t.dirty).map((t) => t.path)
  if (activePath && !paths.includes(activePath)) {
    paths.push(activePath)
  }
  return paths
}
