import { writable } from 'svelte/store'

export type MenuState = {
  openMenu: string | null
}

export const defaultMenuState: MenuState = {
  openMenu: null,
}

export function createMenuStore(initial: MenuState = defaultMenuState) {
  const store = writable<MenuState>(initial)
  return {
    subscribe: store.subscribe,
    toggle(name: string) {
      store.update((s) => ({ openMenu: s.openMenu === name ? null : name }))
    },
    open(name: string) {
      store.set({ openMenu: name })
    },
    close() {
      store.set({ openMenu: null })
    },
    snapshot(): MenuState {
      let state = defaultMenuState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultMenuState)
    },
  }
}
