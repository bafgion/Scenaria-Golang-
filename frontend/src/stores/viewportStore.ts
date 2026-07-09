import { writable } from 'svelte/store'

export type ViewportState = {
  width: number
  height: number
  autoCompact: boolean
  toolbarIconOnly: boolean
}

export const defaultViewportState: ViewportState = {
  width: 1280,
  height: 800,
  autoCompact: false,
  toolbarIconOnly: true,
}

export function createViewportStore(initial: ViewportState = defaultViewportState) {
  const store = writable<ViewportState>(initial)
  return {
    subscribe: store.subscribe,
    syncWindowSize(width: number, height: number, autoCompact: boolean) {
      store.update((s) => ({ ...s, width, height, autoCompact }))
    },
    setToolbarIconOnly(toolbarIconOnly: boolean) {
      store.update((s) => ({ ...s, toolbarIconOnly }))
    },
    snapshot(): ViewportState {
      let state = defaultViewportState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultViewportState)
    },
  }
}
