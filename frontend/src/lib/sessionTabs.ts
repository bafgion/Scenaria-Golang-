import { isUntitled } from './untitled'
import type { TabBody } from './tabMemory'
import { tabEditorText } from './tabMemory'

export type UntitledTabSnapshot = {
  path: string
  content: string
}

export type SessionTabsSnapshot = {
  openTabs: string[]
  untitledTabs: UntitledTabSnapshot[]
  activeTab: string
}

export function buildSessionTabsSnapshot(
  tabs: TabBody[],
  activeTab: string,
  getLiveEditorText: () => string | null,
  welcomeKey: string,
): SessionTabsSnapshot {
  const openTabs = tabs.map((t) => t.path).filter(Boolean)
  const untitledTabs = tabs
    .filter((t) => isUntitled(t.path))
    .map((t) => {
      const stored = tabEditorText(t)
      const live = t.path === activeTab ? getLiveEditorText() : null
      const content = live !== null ? live : stored
      return { path: t.path, content }
    })
  return {
    openTabs,
    untitledTabs,
    activeTab: activeTab !== welcomeKey ? activeTab : '',
  }
}

export function sessionTabPathsFromSettings(
  openTabs: string[] | undefined,
  untitledTabs: UntitledTabSnapshot[] | undefined,
): string[] {
  const merged: string[] = []
  const addPath = (raw: string | undefined) => {
    const path = (raw || '').trim()
    if (path && !merged.includes(path)) {
      merged.push(path)
    }
  }
  for (const path of openTabs || []) addPath(path)
  for (const tab of untitledTabs || []) addPath(tab.path)
  return merged
}

export function untitledContentMap(
  untitledTabs: UntitledTabSnapshot[] | undefined,
): Map<string, string> {
  const map = new Map<string, string>()
  for (const tab of untitledTabs || []) {
    const path = (tab.path || '').trim()
    if (path) {
      const content = tab.content ?? ''
      const existing = map.get(path)
      if (existing === undefined || content.trim() || !existing.trim()) {
        map.set(path, content)
      }
    }
  }
  return map
}

/** Pick the tab that should be focused after session restore. */
export function resolveRestoredActiveTab(
  savedActiveTab: string,
  tabPaths: string[],
  restoredTabs: TabBody[],
  welcomeKey: string,
): string {
  const active = (savedActiveTab || '').trim()
  const normalizedActive = active === welcomeKey ? '' : active
  const featureTabs = restoredTabs.filter((tab) => tab.path && tab.path !== welcomeKey)
  if (featureTabs.length === 0) {
    return welcomeKey
  }
  if (normalizedActive && featureTabs.some((tab) => tab.path === normalizedActive)) {
    return normalizedActive
  }
  for (let i = tabPaths.length - 1; i >= 0; i--) {
    const path = (tabPaths[i] || '').trim()
    if (!path || path === welcomeKey) continue
    if (featureTabs.some((tab) => tab.path === path)) {
      return path
    }
  }
  return featureTabs[featureTabs.length - 1].path
}
