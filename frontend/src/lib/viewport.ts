/** Layout helpers for small screens and low resolutions. */

export const VIEWPORT = {
  compactWidth: 1100,
  narrowWidth: 900,
  previewMinWidth: 960,
  minWindowWidth: 800,
  minWindowHeight: 520,
  activityWidth: 40,
  sidebarMin: 120,
  sidebarMax: 480,
  sidebarNarrowMax: 200,
} as const

export function shouldAutoCompactToolbar(width: number): boolean {
  return width < VIEWPORT.compactWidth
}

export function effectiveSidebarWidth(
  savedWidth: number,
  viewportWidth: number,
  sidebarVisible: boolean,
): number {
  if (!sidebarVisible) return 0
  const splitter = 4
  const reserved = VIEWPORT.activityWidth + splitter + 280
  const maxByViewport = Math.max(VIEWPORT.sidebarMin, viewportWidth - reserved)
  const narrowCap =
    viewportWidth < VIEWPORT.narrowWidth ? VIEWPORT.sidebarNarrowMax : VIEWPORT.sidebarMax
  return Math.max(
    VIEWPORT.sidebarMin,
    Math.min(savedWidth, narrowCap, maxByViewport),
  )
}

export function effectivePreviewWidth(
  savedWidth: number,
  viewportWidth: number,
  previewVisible: boolean,
): number {
  if (!previewVisible || viewportWidth < VIEWPORT.previewMinWidth) return 0
  const editorMin = 320
  const maxByViewport = Math.max(200, Math.floor(viewportWidth * 0.45) - editorMin)
  return Math.max(200, Math.min(savedWidth, maxByViewport, 560))
}

export function shouldShowPreviewPane(viewportWidth: number, previewVisible: boolean): boolean {
  return previewVisible && viewportWidth >= VIEWPORT.previewMinWidth
}

export function clampBottomPanelHeight(height: number, viewportHeight: number): number {
  const max = Math.max(100, Math.floor(viewportHeight * 0.45))
  return Math.max(80, Math.min(height, max))
}

export function clampStepsPanelHeight(height: number, viewportHeight: number): number {
  const max = Math.max(80, Math.floor(viewportHeight * 0.35))
  return Math.max(80, Math.min(height, max))
}

export function toolbarIconOnlyThreshold(barWidth: number): number {
  return Math.min(880, Math.max(320, Math.floor(barWidth * 0.55)))
}
