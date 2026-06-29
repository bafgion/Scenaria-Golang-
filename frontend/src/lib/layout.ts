const KEY = 'scenaria.ui.layout'

export type LayoutState = {
  bottomPanelHeight: number
  sidebarVisible: boolean
  bottomPanelOpen: boolean
  sidebarWidth: number
  previewVisible: boolean
  previewWidth: number
}

const defaults: LayoutState = {
  bottomPanelHeight: 200,
  sidebarVisible: true,
  bottomPanelOpen: false,
  sidebarWidth: 260,
  previewVisible: false,
  previewWidth: 360,
}

export function loadLayout(): LayoutState {
  try {
    const raw = localStorage.getItem(KEY)
    if (!raw) return { ...defaults }
    return { ...defaults, ...JSON.parse(raw) }
  } catch {
    return { ...defaults }
  }
}

export function saveLayout(state: Partial<LayoutState>) {
  try {
    const prev = loadLayout()
    localStorage.setItem(KEY, JSON.stringify({ ...prev, ...state }))
  } catch {
    /* ignore */
  }
}

export function resetLayout(): LayoutState {
  try {
    localStorage.removeItem(KEY)
  } catch {
    /* ignore */
  }
  return { ...defaults }
}
