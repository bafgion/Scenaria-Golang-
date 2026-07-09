import { writable } from 'svelte/store'
import type { gui } from '../../wailsjs/go/models'

export type PluginsState = {
  installed: gui.PluginEntryDTO[]
}

export const defaultPluginsState: PluginsState = {
  installed: [],
}

export function createPluginsStore(initial: PluginsState = defaultPluginsState) {
  const store = writable<PluginsState>(initial)
  return {
    subscribe: store.subscribe,
    setInstalled(installed: gui.PluginEntryDTO[]) {
      store.update((s) => ({ ...s, installed }))
    },
    clear() {
      store.update((s) => ({ ...s, installed: [] }))
    },
    hasVanessa(): boolean {
      return this.snapshot().installed.some((p) => p.vanessa)
    },
    findByName(name: string): gui.PluginEntryDTO | undefined {
      return this.snapshot().installed.find((p) => p.name === name)
    },
    snapshot(): PluginsState {
      let state = defaultPluginsState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultPluginsState)
    },
  }
}
