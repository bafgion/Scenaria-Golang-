import { writable } from 'svelte/store'

export type AppMetaState = {
  version: string
}

export const defaultAppMetaState: AppMetaState = {
  version: '',
}

export function createAppMetaStore(initial: AppMetaState = defaultAppMetaState) {
  const store = writable<AppMetaState>(initial)
  return {
    subscribe: store.subscribe,
    setVersion(version: string) {
      store.update((s) => ({ ...s, version }))
    },
    reset() {
      store.set(defaultAppMetaState)
    },
    snapshot(): AppMetaState {
      let state = defaultAppMetaState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
  }
}
