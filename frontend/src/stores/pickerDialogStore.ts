import { writable } from 'svelte/store'
import type { gui } from '../../wailsjs/go/models'

export type PickerDialogState = {
  selector: string
  choices: gui.PickerStepChoice[]
  candidates: gui.SelectorCandidate[]
  warnings: string[]
  suggestedAction: string
  suggestedChoice: number
}

export const defaultPickerDialogState: PickerDialogState = {
  selector: '',
  choices: [],
  candidates: [],
  warnings: [],
  suggestedAction: '',
  suggestedChoice: 0,
}

export function createPickerDialogStore(initial: PickerDialogState = defaultPickerDialogState) {
  const store = writable<PickerDialogState>(initial)
  return {
    subscribe: store.subscribe,
    setResult(
      selector: string,
      choices: gui.PickerStepChoice[],
      candidates: gui.SelectorCandidate[] = [],
      warnings: string[] = [],
      suggestedAction = '',
      suggestedChoice = 0,
    ) {
      store.set({ selector, choices, candidates, warnings, suggestedAction, suggestedChoice })
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
