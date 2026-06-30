import { describe, expect, it } from 'vitest'
import {
  clampBottomPanelHeight,
  effectiveSidebarWidth,
  shouldAutoCompactToolbar,
  shouldShowPreviewPane,
  toolbarIconOnlyThreshold,
} from './viewport'

describe('viewport layout', () => {
  it('auto-compacts toolbar on narrow width', () => {
    expect(shouldAutoCompactToolbar(1099)).toBe(true)
    expect(shouldAutoCompactToolbar(1100)).toBe(false)
  })

  it('caps sidebar width on small screens', () => {
    expect(effectiveSidebarWidth(260, 1280, true)).toBe(260)
    expect(effectiveSidebarWidth(260, 860, true)).toBeLessThanOrEqual(200)
    expect(effectiveSidebarWidth(260, 860, false)).toBe(0)
  })

  it('hides preview pane when viewport is too narrow', () => {
    expect(shouldShowPreviewPane(1280, true)).toBe(true)
    expect(shouldShowPreviewPane(900, true)).toBe(false)
  })

  it('clamps bottom panel height', () => {
    expect(clampBottomPanelHeight(400, 600)).toBeLessThanOrEqual(270)
    expect(clampBottomPanelHeight(120, 600)).toBe(120)
  })

  it('scales toolbar icon-only threshold', () => {
    expect(toolbarIconOnlyThreshold(500)).toBeLessThan(880)
    expect(toolbarIconOnlyThreshold(1600)).toBe(880)
  })
})
