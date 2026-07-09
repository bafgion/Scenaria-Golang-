import { writable } from 'svelte/store'
import type { gui } from '../../wailsjs/go/models'

export type RecorderPrefsState = {
  filterRecording: boolean
  navOnlyRecording: boolean
  hoverRecord: boolean
}

export const defaultRecorderPrefsState: RecorderPrefsState = {
  filterRecording: false,
  navOnlyRecording: false,
  hoverRecord: false,
}

export function createRecorderPrefsStore(initial: RecorderPrefsState = defaultRecorderPrefsState) {
  const store = writable<RecorderPrefsState>(initial)
  return {
    subscribe: store.subscribe,
    patch(partial: Partial<RecorderPrefsState>) {
      store.update((state) => ({ ...state, ...partial }))
    },
    applyFromDTO(dto: gui.AppSettingsDTO) {
      store.update((state) => ({
        ...state,
        filterRecording: !!dto.filterRecording,
        navOnlyRecording: !!dto.navOnlyRecording,
        hoverRecord: !!dto.hoverRecord,
      }))
    },
    dtoFields(state: RecorderPrefsState): Partial<gui.AppSettingsDTO> {
      return {
        filterRecording: state.filterRecording,
        navOnlyRecording: state.navOnlyRecording,
        hoverRecord: state.hoverRecord,
      }
    },
    snapshot(): RecorderPrefsState {
      let state = defaultRecorderPrefsState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultRecorderPrefsState)
    },
  }
}
