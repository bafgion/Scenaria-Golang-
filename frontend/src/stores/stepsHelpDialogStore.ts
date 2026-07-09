import { writable } from 'svelte/store'

export type StepsHelpDialogState = {
  query: string
}

export const defaultStepsHelpDialogState: StepsHelpDialogState = {
  query: '',
}

export function createStepsHelpDialogStore(initial: StepsHelpDialogState = defaultStepsHelpDialogState) {
  const store = writable<StepsHelpDialogState>(initial)
  return {
    subscribe: store.subscribe,
    setQuery(query: string) {
      store.set({ query })
    },
    clear() {
      store.set(defaultStepsHelpDialogState)
    },
    snapshot(): StepsHelpDialogState {
      let state = defaultStepsHelpDialogState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultStepsHelpDialogState)
    },
  }
}
