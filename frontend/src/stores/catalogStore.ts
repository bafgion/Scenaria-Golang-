import { writable } from 'svelte/store'
import { buildCatalogStructure, catalogStructureKey, type CatalogNode } from '../lib/catalogTree'
import { remapBatchSelectedPaths, toggleBatchPath } from '../lib/batchSelection'

export type CatalogState = {
  batchMode: boolean
  batchSelected: string[]
  sidebarSearch: string
  catalogFilterText: string
  collapsedKeys: string[]
  dropTarget: string
  showBatchHint: boolean
  baseTreeKey: string
  baseTree: CatalogNode | null
}

export const defaultCatalogState: CatalogState = {
  batchMode: false,
  batchSelected: [],
  sidebarSearch: '',
  catalogFilterText: '',
  collapsedKeys: [],
  dropTarget: '',
  showBatchHint: true,
  baseTreeKey: '',
  baseTree: null,
}
export function createCatalogStore(initial: CatalogState = defaultCatalogState) {
  const store = writable<CatalogState>(initial)
  return {
    subscribe: store.subscribe,
    patch(partial: Partial<CatalogState>) {
      store.update((s) => ({ ...s, ...partial }))
    },
    setSidebarSearch(sidebarSearch: string) {
      store.update((s) => ({ ...s, sidebarSearch }))
    },
    setFilterText(catalogFilterText: string) {
      store.update((s) => ({ ...s, catalogFilterText }))
    },
    setDropTarget(dropTarget: string) {
      store.update((s) => ({ ...s, dropTarget }))
    },
    syncBaseTree(projectPath: string, features: string[]) {
      const key = projectPath ? catalogStructureKey(projectPath, features) : ''
      store.update((s) => {
        if (s.baseTreeKey === key) return s
        return {
          ...s,
          baseTreeKey: key,
          baseTree: projectPath ? buildCatalogStructure(projectPath, features) : null,
        }
      })
    },
    toggleBatchPath(path: string) {
      store.update((s) => ({ ...s, batchSelected: toggleBatchPath(s.batchSelected, path) }))
    },
    setBatchSelected(batchSelected: string[]) {
      store.update((s) => ({ ...s, batchSelected }))
    },
    mapBatchSelected(mapper: (paths: string[]) => string[]) {
      store.update((s) => ({ ...s, batchSelected: mapper(s.batchSelected) }))
    },
    remapBatchSelected(features: string[]) {
      store.update((s) => ({
        ...s,
        batchSelected: remapBatchSelectedPaths(s.batchSelected, features),
      }))
    },
    clearBatch() {
      store.update((s) => ({ ...s, batchMode: false, batchSelected: [] }))
    },
    setCollapsed(key: string, collapsed: boolean) {
      store.update((s) => {
        const keys = new Set(s.collapsedKeys)
        if (collapsed) keys.add(key)
        else keys.delete(key)
        return { ...s, collapsedKeys: [...keys] }
      })
    },
    collapsedSet(): Set<string> {
      return new Set(this.snapshot().collapsedKeys)
    },
    snapshot(): CatalogState {
      let state = defaultCatalogState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultCatalogState)
    },
  }
}
