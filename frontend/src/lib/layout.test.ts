import { beforeEach, describe, expect, it } from 'vitest'
import { loadLayout, resetLayout, saveLayout } from './layout'

describe('layout persistence', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('returns defaults when storage is empty', () => {
    const state = loadLayout()
    expect(state.sidebarVisible).toBe(true)
    expect(state.bottomPanelHeight).toBe(200)
    expect(state.previewVisible).toBe(false)
  })

  it('merges partial updates into saved layout', () => {
    saveLayout({ sidebarVisible: false, previewWidth: 420 })
    const state = loadLayout()
    expect(state.sidebarVisible).toBe(false)
    expect(state.previewWidth).toBe(420)
    expect(state.sidebarWidth).toBe(260)
  })

  it('resetLayout clears storage and returns defaults', () => {
    saveLayout({ bottomPanelOpen: true })
    const state = resetLayout()
    expect(state.bottomPanelOpen).toBe(false)
    expect(localStorage.getItem('scenaria.ui.layout')).toBeNull()
  })
})
