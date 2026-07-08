import type { TabBody } from '../lib/tabMemory'

export type TabsCloseResult = {
  tabs: TabBody[]
  openNextPath: string
  showWelcome: boolean
}

/**
 * Pure reducer for close-tab flow:
 * - removes closed tab
 * - returns next tab that should be activated, or welcome fallback.
 */
export function reduceTabsAfterClose(
  tabs: TabBody[],
  activeTab: string,
  closePath: string,
): TabsCloseResult {
  const nextTabs = tabs.filter((t) => t.path !== closePath)
  if (activeTab !== closePath) {
    return {
      tabs: nextTabs,
      openNextPath: '',
      showWelcome: false,
    }
  }
  const next = nextTabs[nextTabs.length - 1]
  if (next) {
    return {
      tabs: nextTabs,
      openNextPath: next.path,
      showWelcome: false,
    }
  }
  return {
    tabs: nextTabs,
    openNextPath: '',
    showWelcome: true,
  }
}
