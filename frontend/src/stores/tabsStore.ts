import { writable } from 'svelte/store'
import type { TabBody } from '../lib/tabMemory'

export type TabsState = {
  tabs: TabBody[]
  activeTab: string
  welcomeTabVisible: boolean
  pendingCloseTab: string | null
  loadFeatureGeneration: number
}

export type TabsCloseResult = {
  tabs: TabBody[]
  openNextPath: string
  showWelcome: boolean
}

/** Welcome screen/tab is available only when there are no open feature/untitled tabs. */
export function canShowWelcome(tabs: TabBody[]): boolean {
  return tabs.length === 0
}

/** @deprecated use canShowWelcome */
export function welcomeTabVisibleForTabs(tabs: TabBody[]): boolean {
  return canShowWelcome(tabs)
}

function resolveActiveTabForTabs(
  tabs: TabBody[],
  activeTab: string,
  welcomeKey: string,
): string {
  if (tabs.length === 0) {
    return welcomeKey
  }
  if (activeTab !== welcomeKey && tabs.some((t) => t.path === activeTab)) {
    return activeTab
  }
  return tabs[tabs.length - 1]?.path ?? activeTab
}

function normalizeTabsState(state: TabsState, welcomeKey: string): TabsState {
  const welcomeTabVisible = canShowWelcome(state.tabs)
  return {
    ...state,
    activeTab: welcomeTabVisible
      ? welcomeKey
      : resolveActiveTabForTabs(state.tabs, state.activeTab, welcomeKey),
    welcomeTabVisible,
  }
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

export function createTabsStore(
  welcomeKey: string,
  initial: TabsState = {
    tabs: [],
    activeTab: welcomeKey,
    welcomeTabVisible: true,
    pendingCloseTab: null,
    loadFeatureGeneration: 0,
  },
) {
  const store = writable<TabsState>(initial)
  const normalize = (state: TabsState) => normalizeTabsState(state, welcomeKey)
  return {
    subscribe: store.subscribe,
    setTabs(tabs: TabBody[]) {
      store.update((s) => normalize({ ...s, tabs }))
    },
    mapTabs(fn: (tabs: TabBody[]) => TabBody[]) {
      store.update((s) => normalize({ ...s, tabs: fn(s.tabs) }))
    },
    mapEachTab(fn: (tab: TabBody) => TabBody) {
      store.update((s) => normalize({ ...s, tabs: s.tabs.map(fn) }))
    },
    appendTab(tab: TabBody) {
      store.update((s) => normalize({ ...s, tabs: [...s.tabs, tab] }))
    },
    setActiveTab(activeTab: string) {
      store.update((s) => normalize({ ...s, activeTab }))
    },
    setWelcomeVisible(welcomeTabVisible: boolean) {
      store.update((s) =>
        normalize({
          ...s,
          welcomeTabVisible: welcomeTabVisible && canShowWelcome(s.tabs),
        }),
      )
    },
    patch(partial: Partial<TabsState>) {
      store.update((s) => normalize({ ...s, ...partial }))
    },
    applyCloseResult(result: TabsCloseResult) {
      store.update((s) => {
        if (result.showWelcome) {
          return normalize({
            tabs: result.tabs,
            activeTab: welcomeKey,
            welcomeTabVisible: true,
            pendingCloseTab: null,
            loadFeatureGeneration: s.loadFeatureGeneration,
          })
        }
        const nextActiveTab = result.openNextPath || s.activeTab
        return normalize({
          ...s,
          tabs: result.tabs,
          activeTab: nextActiveTab,
          pendingCloseTab: null,
          welcomeTabVisible: false,
        })
      })
    },
    closePath(path: string): TabsCloseResult {
      let result: TabsCloseResult = {
        tabs: [],
        openNextPath: '',
        showWelcome: false,
      }
      store.update((s) => {
        result = reduceTabsAfterClose(s.tabs, s.activeTab, path)
        if (result.showWelcome) {
          return normalize({
            tabs: result.tabs,
            activeTab: welcomeKey,
            welcomeTabVisible: true,
            pendingCloseTab: null,
            loadFeatureGeneration: s.loadFeatureGeneration,
          })
        }
        return normalize({
          ...s,
          tabs: result.tabs,
          activeTab: result.openNextPath || s.activeTab,
          pendingCloseTab: null,
          welcomeTabVisible: false,
        })
      })
      return result
    },
    setPendingCloseTab(pendingCloseTab: string | null) {
      store.update((s) => ({ ...s, pendingCloseTab }))
    },
    bumpLoadFeatureGeneration(): number {
      let next = 0
      store.update((s) => {
        next = s.loadFeatureGeneration + 1
        return { ...s, loadFeatureGeneration: next }
      })
      return next
    },
    isLoadFeatureGenerationCurrent(generation: number): boolean {
      return this.snapshot().loadFeatureGeneration === generation
    },
    snapshot(): TabsState {
      let state: TabsState = {
        tabs: [],
        activeTab: welcomeKey,
        welcomeTabVisible: true,
        pendingCloseTab: null,
        loadFeatureGeneration: 0,
      }
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set({
        tabs: [],
        activeTab: welcomeKey,
        welcomeTabVisible: true,
        pendingCloseTab: null,
        loadFeatureGeneration: 0,
      })
    },
  }
}
