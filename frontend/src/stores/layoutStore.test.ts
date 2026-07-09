import { beforeEach, describe, expect, it } from 'vitest'
import { createLayoutStore, type LayoutStoreState } from './layoutStore'

describe('layoutStore', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('loads persisted layout on create', () => {
    localStorage.setItem(
      'scenaria.ui.layout',
      JSON.stringify({ sidebarVisible: false, bottomPanelOpen: true }),
    )
    const store = createLayoutStore()
    let snapshot: LayoutStoreState = {
      sidebarVisible: true,
      bottomPanelOpen: false,
      bottomPanelHeight: 0,
      previewVisible: false,
      previewWidth: 0,
      bottomTab: 'journal',
      resizingBottom: false,
      resizingSteps: false,
      resizingSidebar: false,
      resizingPreview: false,
      previewPaneMounted: false,
    }
    const unsub = store.subscribe((s) => {
      snapshot = s
    })
    unsub()
    expect(snapshot.sidebarVisible).toBe(false)
    expect(snapshot.bottomPanelOpen).toBe(true)
  })

  it('persists patch updates', () => {
    const store = createLayoutStore()
    store.patch({ previewWidth: 480 })
    expect(JSON.parse(localStorage.getItem('scenaria.ui.layout') || '{}').previewWidth).toBe(480)
  })

  it('toggleSidebar flips visibility and saves', () => {
    const store = createLayoutStore()
    store.toggleSidebar()
    let visible = true
    const unsub = store.subscribe((s) => {
      visible = s.sidebarVisible
    })
    unsub()
    expect(visible).toBe(false)
  })

  it('openBottomTab opens panel and selects tab', () => {
    const store = createLayoutStore()
    store.openBottomTab('results')
    let state: Pick<LayoutStoreState, 'bottomPanelOpen' | 'bottomTab'> = {
      bottomPanelOpen: false,
      bottomTab: 'journal',
    }
    const unsub = store.subscribe((s) => {
      state = { bottomPanelOpen: s.bottomPanelOpen, bottomTab: s.bottomTab }
    })
    unsub()
    expect(state.bottomPanelOpen).toBe(true)
    expect(state.bottomTab).toBe('results')
  })
})
