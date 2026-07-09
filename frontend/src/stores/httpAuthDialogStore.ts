import { writable } from 'svelte/store'

export type HttpAuthDialogState = {
  host: string
}

export const defaultHttpAuthDialogState: HttpAuthDialogState = {
  host: '',
}

export function createHttpAuthDialogStore(initial: HttpAuthDialogState = defaultHttpAuthDialogState) {
  const store = writable<HttpAuthDialogState>(initial)
  return {
    subscribe: store.subscribe,
    setHost(host: string) {
      store.update((s) => ({ ...s, host }))
    },
    snapshot(): HttpAuthDialogState {
      let state = defaultHttpAuthDialogState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultHttpAuthDialogState)
    },
  }
}
