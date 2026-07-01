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
  editorText: string,
  welcomeKey: string,
): SessionTabsSnapshot {
  const openTabs = tabs.map((t) => t.path).filter(Boolean)
  const untitledTabs = tabs
    .filter((t) => isUntitled(t.path))
    .map((t) => ({
      path: t.path,
      content: t.path === activeTab ? editorText : tabEditorText(t),
    }))
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
  const paths = (openTabs || []).filter((p) => p && p.trim())
  if (paths.length > 0) {
    return paths
  }
  return (untitledTabs || []).map((t) => t.path).filter(Boolean)
}

export function untitledContentMap(
  untitledTabs: UntitledTabSnapshot[] | undefined,
): Map<string, string> {
  const map = new Map<string, string>()
  for (const tab of untitledTabs || []) {
    if (tab.path) {
      map.set(tab.path, tab.content ?? '')
    }
  }
  return map
}
