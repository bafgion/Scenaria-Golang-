import {
  BeginSplashWindowChrome,
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

/** Compact centered window for splash-only phase (Photoshop-style). */
export async function beginSplashWindow() {
  if (!hasWailsRuntime()) return
  try {
    WindowSetMinSize(SPLASH_WINDOW.width, SPLASH_WINDOW.height)
    WindowSetMaxSize(SPLASH_WINDOW.width, SPLASH_WINDOW.height)
    WindowSetSize(SPLASH_WINDOW.width, SPLASH_WINDOW.height)
    WindowCenter()
    WindowShow()
    await BeginSplashWindowChrome()
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
    WindowCenter()
  } catch {
    /* dev without wails */
  }
}

export function setSplashDocumentState(active: boolean) {
  document.documentElement.classList.toggle('splash-active', active)
  document.body.classList.toggle('splash-active', active)
}
