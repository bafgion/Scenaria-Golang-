import { writable } from 'svelte/store'

export type ProjectReplaceState = {
  findText: string
  replaceText: string
  caseSensitive: boolean
  busy: boolean
}

export const defaultProjectReplaceState: ProjectReplaceState = {
  findText: '',
  replaceText: '',
  caseSensitive: false,
  busy: false,
}

export function createProjectReplaceStore(initial: ProjectReplaceState = defaultProjectReplaceState) {
  const store = writable<ProjectReplaceState>(initial)
  return {
    subscribe: store.subscribe,
    patch(partial: Partial<ProjectReplaceState>) {
      store.update((s) => ({ ...s, ...partial }))
    },
    setBusy(busy: boolean) {
      store.update((s) => ({ ...s, busy }))
    },
    snapshot(): ProjectReplaceState {
      let state = defaultProjectReplaceState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultProjectReplaceState)
    },
  }
}
