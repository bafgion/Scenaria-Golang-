import { writable } from 'svelte/store'
import type { gui } from '../../wailsjs/go/models'

export type PickerDialogState = {
  selector: string
  choices: gui.PickerStepChoice[]
}

export const defaultPickerDialogState: PickerDialogState = {
  selector: '',
  choices: [],
}

export function createPickerDialogStore(initial: PickerDialogState = defaultPickerDialogState) {
  const store = writable<PickerDialogState>(initial)
  return {
    subscribe: store.subscribe,
    setResult(selector: string, choices: gui.PickerStepChoice[]) {
      store.set({ selector, choices })
    },
    clear() {
      store.set(defaultPickerDialogState)
    },
    snapshot(): PickerDialogState {
      let state = defaultPickerDialogState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultPickerDialogState)
    },
  }
}
