import {
  BeginSplashWindowChrome,
  CenterAppWindow,
  OpenMainWindowChrome,
} from '../../wailsjs/go/wailsapp/App'
import {
  WindowCenter,
  WindowSetMaxSize,
  WindowSetMinSize,
  WindowSetSize,
  WindowShow,
} from '../../wailsjs/runtime/runtime'

export const SPLASH_WINDOW = { width: 560, height: 500 }
export const MAIN_WINDOW = { width: 1280, height: 800, minWidth: 960, minHeight: 640 }

function hasWailsRuntime(): boolean {
  return typeof window !== 'undefined' && !!(window as Window & { runtime?: unknown }).runtime
}

async function waitForWindowLayout(): Promise<void> {
  await new Promise<void>((resolve) => {
    requestAnimationFrame(() => requestAnimationFrame(() => resolve()))
  })
}

/** Center after size/chrome changes; Wails runtime + native Win32 when available. */
export async function centerAppWindowOnScreen() {
  if (!hasWailsRuntime()) return
  try {
    WindowCenter()
    await CenterAppWindow()
    await waitForWindowLayout()
    WindowCenter()
    await CenterAppWindow()
  } catch {
    /* dev without wails */
  }
}

/** Compact centered window for splash-only phase (Photoshop-style). */
export async function beginSplashWindow() {
  if (!hasWailsRuntime()) return
  try {
    WindowSetMinSize(SPLASH_WINDOW.width, SPLASH_WINDOW.height)
    WindowSetMaxSize(SPLASH_WINDOW.width, SPLASH_WINDOW.height)
    WindowSetSize(SPLASH_WINDOW.width, SPLASH_WINDOW.height)
    WindowShow()
    await BeginSplashWindowChrome()
    await centerAppWindowOnScreen()
  } catch {
    /* dev without wails */
  }
}

/** Restore main IDE window size after splash. */
export async function openMainWindow() {
  if (!hasWailsRuntime()) return
  try {
    await OpenMainWindowChrome()
    WindowSetMaxSize(0, 0)
    WindowSetMinSize(MAIN_WINDOW.minWidth, MAIN_WINDOW.minHeight)
    WindowSetSize(MAIN_WINDOW.width, MAIN_WINDOW.height)
    await centerAppWindowOnScreen()
  } catch {
    /* dev without wails */
  }
}

export function setSplashDocumentState(active: boolean) {
  document.documentElement.classList.toggle('splash-active', active)
  document.body.classList.toggle('splash-active', active)
}
