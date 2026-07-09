import { writable } from 'svelte/store'

export type TestClientState = {
  clients: string[]
  selection: string
  suggestName: string
}

export const defaultTestClientState: TestClientState = {
  clients: [],
  selection: '',
  suggestName: '',
}

export function createTestClientStore(initial: TestClientState = defaultTestClientState) {
  const store = writable<TestClientState>(initial)
  return {
    subscribe: store.subscribe,
    setClients(clients: string[]) {
      store.update((s) => ({ ...s, clients }))
    },
    patch(partial: Partial<TestClientState>) {
      store.update((s) => ({ ...s, ...partial }))
    },
    clearSuggestName() {
      store.update((s) => ({ ...s, suggestName: '' }))
    },
    snapshot(): TestClientState {
      let state = defaultTestClientState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultTestClientState)
    },
  }
}
