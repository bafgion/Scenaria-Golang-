import { writable } from 'svelte/store'

export type RunDialogState = {
  title: string
  scenarios: string[]
}

export const defaultRunDialogState: RunDialogState = {
  title: '',
  scenarios: [],
}

export function createRunDialogStore(initial: RunDialogState = defaultRunDialogState) {
  const store = writable<RunDialogState>(initial)
  return {
    subscribe: store.subscribe,
    open(title: string, scenarios: string[]) {
      store.set({ title, scenarios })
    },
    reset() {
      store.set(defaultRunDialogState)
    },
    snapshot(): RunDialogState {
      let state = defaultRunDialogState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
  }
}
