import { writable } from 'svelte/store'
import type { gui } from '../../wailsjs/go/models'

export type FeatureContextMenu = { x: number; y: number; path: string }
export type FolderContextMenu = { x: number; y: number; dir: string; paths: string[] }
export type StepsContextMenu = { x: number; y: number; line: number; step: gui.EditorStepRow }

export type ContextMenuState = {
  feature: FeatureContextMenu | null
  folder: FolderContextMenu | null
  steps: StepsContextMenu | null
}

export const defaultContextMenuState: ContextMenuState = {
  feature: null,
  folder: null,
  steps: null,
}

export function createContextMenuStore(initial: ContextMenuState = defaultContextMenuState) {
  const store = writable<ContextMenuState>(initial)
  return {
    subscribe: store.subscribe,
    openFeature(menu: FeatureContextMenu) {
      store.set({ feature: menu, folder: null, steps: null })
    },
    openFolder(menu: FolderContextMenu) {
      store.set({ feature: null, folder: menu, steps: null })
    },
    openSteps(menu: StepsContextMenu) {
      store.update((s) => ({ ...s, steps: menu }))
    },
    closeFeature() {
      store.update((s) => ({ ...s, feature: null }))
    },
    closeFolder() {
      store.update((s) => ({ ...s, folder: null }))
    },
    closeSteps() {
      store.update((s) => ({ ...s, steps: null }))
    },
    closeAll() {
      store.set(defaultContextMenuState)
    },
    snapshot(): ContextMenuState {
      let state = defaultContextMenuState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultContextMenuState)
    },
  }
}
