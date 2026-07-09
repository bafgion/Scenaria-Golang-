import { writable } from 'svelte/store'
import type { gui } from '../../wailsjs/go/models'

export type HtmlReportOpenMode = 'full' | 'light'

export type SettingsDialogState = {
  baseline: gui.AppSettingsDTO | null
  htmlReportOpenMode: HtmlReportOpenMode
  projectBaseline: HtmlReportOpenMode
}

export const defaultSettingsDialogState: SettingsDialogState = {
  baseline: null,
  htmlReportOpenMode: 'full',
  projectBaseline: 'full',
}

export function createSettingsDialogStore(initial: SettingsDialogState = defaultSettingsDialogState) {
  const store = writable<SettingsDialogState>(initial)
  return {
    subscribe: store.subscribe,
    openWithBaseline(baseline: gui.AppSettingsDTO, htmlReportOpenMode: HtmlReportOpenMode) {
      store.set({
        baseline,
        htmlReportOpenMode,
        projectBaseline: htmlReportOpenMode,
      })
    },
    setHtmlReportOpenMode(htmlReportOpenMode: HtmlReportOpenMode) {
      store.update((s) => ({ ...s, htmlReportOpenMode }))
    },
    markSaved(saved: gui.AppSettingsDTO) {
      store.update((s) => ({ ...s, baseline: saved }))
    },
    clearBaseline() {
      store.update((s) => ({ ...s, baseline: null }))
    },
    restoreProjectBaseline() {
      store.update((s) => ({ ...s, htmlReportOpenMode: s.projectBaseline }))
    },
    reset() {
      store.set(defaultSettingsDialogState)
    },
    snapshot(): SettingsDialogState {
      let state = defaultSettingsDialogState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
  }
}
