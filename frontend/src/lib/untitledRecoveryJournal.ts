import { isUntitled } from './untitled'
import type { UntitledTabSnapshot } from './sessionTabs'
import type { TabBody } from './tabMemory'
import { tabEditorText } from './tabMemory'

export const UNTITLED_RECOVERY_KEY = 'scenaria.untitledRecovery.v1'
export const UNTITLED_RECOVERY_MAX_AGE_MS = 7 * 24 * 60 * 60 * 1000

export type UntitledRecoveryEntry = {
  path: string
  content: string
  updatedAt: number
}

export type UntitledRecoveryJournal = {
  version: 1
  tabs: UntitledRecoveryEntry[]
}

export type UntitledRecoveryStorage = Pick<Storage, 'getItem' | 'setItem' | 'removeItem'>

function defaultStorage(): UntitledRecoveryStorage | null {
  try {
    return typeof window !== 'undefined' ? window.localStorage : null
  } catch {
    return null
  }
}

export function loadUntitledRecoveryJournal(
  storage: UntitledRecoveryStorage | null = defaultStorage(),
  now = Date.now(),
): UntitledRecoveryEntry[] {
  if (!storage) return []
  let raw = ''
  try {
    raw = storage.getItem(UNTITLED_RECOVERY_KEY) || ''
  } catch {
    return []
  }
  if (!raw) return []
  try {
    const parsed = JSON.parse(raw) as Partial<UntitledRecoveryJournal>
    if (parsed.version !== 1 || !Array.isArray(parsed.tabs)) {
      return []
    }
    const byPath = new Map<string, UntitledRecoveryEntry>()
    for (const tab of parsed.tabs) {
      const path = (tab?.path || '').trim()
      const updatedAt = Number(tab?.updatedAt || 0)
      if (!path || !isUntitled(path) || !Number.isFinite(updatedAt)) continue
      if (now - updatedAt > UNTITLED_RECOVERY_MAX_AGE_MS) continue
      const content = tab.content ?? ''
      const existing = byPath.get(path)
      if (!existing || updatedAt >= existing.updatedAt) {
        byPath.set(path, { path, content, updatedAt })
      }
    }
    return [...byPath.values()].sort((a, b) => a.updatedAt - b.updatedAt)
  } catch {
    try {
      storage.removeItem(UNTITLED_RECOVERY_KEY)
    } catch {
      /* offline */
    }
    return []
  }
}

export function persistUntitledRecoveryEntries(
  entries: UntitledRecoveryEntry[],
  storage: UntitledRecoveryStorage | null = defaultStorage(),
): boolean {
  if (!storage) return false
  try {
    if (entries.length === 0) {
      storage.removeItem(UNTITLED_RECOVERY_KEY)
      return true
    }
    storage.setItem(
      UNTITLED_RECOVERY_KEY,
      JSON.stringify({ version: 1, tabs: entries } satisfies UntitledRecoveryJournal),
    )
    return true
  } catch {
    return false
  }
}

export function journalUntitledTabs(
  tabs: TabBody[],
  activeTab: string,
  getLiveEditorText: () => string | null,
  storage: UntitledRecoveryStorage | null = defaultStorage(),
  now = Date.now(),
): boolean {
  const entries = tabs
    .filter((tab) => isUntitled(tab.path))
    .map((tab) => {
      const live = tab.path === activeTab ? getLiveEditorText() : null
      return {
        path: tab.path,
        content: live !== null ? live : tabEditorText(tab),
        updatedAt: now,
      }
    })
  return persistUntitledRecoveryEntries(entries, storage)
}

export function clearUntitledRecoveryPath(
  path: string,
  storage: UntitledRecoveryStorage | null = defaultStorage(),
  now = Date.now(),
): boolean {
  const target = (path || '').trim()
  if (!target || !isUntitled(target)) return true
  const entries = loadUntitledRecoveryJournal(storage, now).filter((tab) => tab.path !== target)
  return persistUntitledRecoveryEntries(entries, storage)
}

export function clearUntitledRecoveryAll(
  storage: UntitledRecoveryStorage | null = defaultStorage(),
): boolean {
  return persistUntitledRecoveryEntries([], storage)
}

export function mergeUntitledSessionWithRecovery(
  sessionTabs: UntitledTabSnapshot[] | undefined,
  recoveryTabs: UntitledRecoveryEntry[],
): UntitledTabSnapshot[] {
  const byPath = new Map<string, UntitledTabSnapshot>()
  for (const tab of sessionTabs || []) {
    const path = (tab?.path || '').trim()
    if (path && isUntitled(path) && !byPath.has(path)) {
      byPath.set(path, { path, content: tab.content ?? '' })
    }
  }
  for (const tab of recoveryTabs) {
    const path = (tab?.path || '').trim()
    if (path && isUntitled(path)) {
      byPath.set(path, { path, content: tab.content ?? '' })
    }
  }
  return [...byPath.values()]
}
