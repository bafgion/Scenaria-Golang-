import { writable } from 'svelte/store'

export type RecentsState = {
  projects: string[]
  features: string[]
}

export const defaultRecentsState: RecentsState = {
  projects: [],
  features: [],
}

export function createRecentsStore(initial: RecentsState = defaultRecentsState) {
  const store = writable<RecentsState>(initial)
  return {
    subscribe: store.subscribe,
    setRecents(projects: string[], features: string[]) {
      store.set({ projects, features })
    },
    patch(partial: Partial<RecentsState>) {
      store.update((s) => ({ ...s, ...partial }))
    },
    snapshot(): RecentsState {
      let state = defaultRecentsState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultRecentsState)
    },
  }
}
