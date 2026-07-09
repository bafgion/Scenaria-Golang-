import { writable } from 'svelte/store'
import { loadLayout, resetLayout, saveLayout, type LayoutState } from '../lib/layout'

export type { LayoutState }

export type BottomPanelTab = 'journal' | 'results' | 'validate' | 'error'

export type LayoutStoreState = LayoutState & {
  bottomTab: BottomPanelTab
  resizingBottom: boolean
  resizingSteps: boolean
  resizingSidebar: boolean
  resizingPreview: boolean
  previewPaneMounted: boolean
}

const defaultBottomTab: BottomPanelTab = 'journal'

const defaultSessionLayout: Pick<
  LayoutStoreState,
  'resizingBottom' | 'resizingSteps' | 'resizingSidebar' | 'resizingPreview' | 'previewPaneMounted'
> = {
  resizingBottom: false,
  resizingSteps: false,
  resizingSidebar: false,
  resizingPreview: false,
  previewPaneMounted: false,
}

export function createLayoutStore(initial: LayoutState = loadLayout()) {
  const store = writable<LayoutStoreState>({ ...initial, bottomTab: defaultBottomTab, ...defaultSessionLayout })
  let previewMountTimer: ReturnType<typeof setTimeout> | null = null
  return {
    subscribe: store.subscribe,
    patch(partial: Partial<LayoutState>) {
      store.update((state) => {
        const next = { ...state, ...partial }
        saveLayout(partial)
        return next
      })
    },
    patchLocal(partial: Partial<LayoutStoreState>) {
      store.update((state) => ({ ...state, ...partial }))
    },
    persist(state: Partial<LayoutState>) {
      saveLayout(state)
    },
    setBottomTab(bottomTab: BottomPanelTab) {
      store.update((state) => ({ ...state, bottomTab }))
    },
    openBottomTab(bottomTab: BottomPanelTab) {
      store.update((state) => {
        saveLayout({ bottomPanelOpen: true })
        return { ...state, bottomPanelOpen: true, bottomTab }
      })
    },
    toggleSidebar() {
      store.update((state) => {
        const sidebarVisible = !state.sidebarVisible
        saveLayout({ sidebarVisible })
        return { ...state, sidebarVisible }
      })
    },
    toggleBottomPanel() {
      store.update((state) => {
        const bottomPanelOpen = !state.bottomPanelOpen
        saveLayout({ bottomPanelOpen })
        return { ...state, bottomPanelOpen }
      })
    },
    showSidebar() {
      this.patch({ sidebarVisible: true })
    },
    hideSidebar() {
      this.patch({ sidebarVisible: false })
    },
    openBottomPanel() {
      this.patch({ bottomPanelOpen: true })
    },
    closeBottomPanel() {
      this.patch({ bottomPanelOpen: false })
    },
    setResizing(field: keyof typeof defaultSessionLayout, value: boolean) {
      store.update((state) => ({ ...state, [field]: value }))
    },
    setPreviewPaneMounted(previewPaneMounted: boolean) {
      store.update((state) => ({ ...state, previewPaneMounted }))
    },
    schedulePreviewMount(onMount: () => void, delayMs = 80) {
      if (previewMountTimer) {
        clearTimeout(previewMountTimer)
        previewMountTimer = null
      }
      previewMountTimer = setTimeout(() => {
        previewMountTimer = null
        onMount()
      }, delayMs)
    },
    clearPreviewMountTimer() {
      if (previewMountTimer) {
        clearTimeout(previewMountTimer)
        previewMountTimer = null
      }
    },
    reset() {
      this.clearPreviewMountTimer()
      const defaults = resetLayout()
      store.set({ ...defaults, bottomTab: defaultBottomTab, ...defaultSessionLayout })
    },
    reload() {
      store.set({ ...loadLayout(), bottomTab: defaultBottomTab, ...defaultSessionLayout })
    },
  }
}
