import { writable } from 'svelte/store'

export type FeatureDialogState = {
  duplicateFeaturePath: string
  duplicateNewName: string
  moveFeaturePath: string
  moveDestDirs: string[]
  moveDestDir: string
  renameFeaturePath: string
  exportInputPath: string
  importDestDir: string
  importFeaturesBusy: boolean
}

export const defaultFeatureDialogState: FeatureDialogState = {
  duplicateFeaturePath: '',
  duplicateNewName: '',
  moveFeaturePath: '',
  moveDestDirs: [],
  moveDestDir: '',
  renameFeaturePath: '',
  exportInputPath: '',
  importDestDir: '',
  importFeaturesBusy: false,
}

export function createFeatureDialogStore(initial: FeatureDialogState = defaultFeatureDialogState) {
  const store = writable<FeatureDialogState>(initial)
  return {
    subscribe: store.subscribe,
    patch(partial: Partial<FeatureDialogState>) {
      store.update((s) => ({ ...s, ...partial }))
    },
    openDuplicate(path: string, newName = '') {
      store.update((s) => ({ ...s, duplicateFeaturePath: path, duplicateNewName: newName }))
    },
    clearDuplicate() {
      store.update((s) => ({ ...s, duplicateFeaturePath: '', duplicateNewName: '' }))
    },
    openMove(path: string, destDirs: string[]) {
      store.update((s) => ({
        ...s,
        moveFeaturePath: path,
        moveDestDirs: destDirs,
        moveDestDir: destDirs[0] ?? '',
      }))
    },
    clearMove() {
      store.update((s) => ({
        ...s,
        moveFeaturePath: '',
        moveDestDirs: [],
        moveDestDir: '',
      }))
    },
    openRename(path: string) {
      store.update((s) => ({ ...s, renameFeaturePath: path }))
    },
    clearRename() {
      store.update((s) => ({ ...s, renameFeaturePath: '' }))
    },
    openExport(path: string) {
      store.update((s) => ({ ...s, exportInputPath: path }))
    },
    openImport(destDir: string) {
      store.update((s) => ({ ...s, importDestDir: destDir }))
    },
    snapshot(): FeatureDialogState {
      let state = defaultFeatureDialogState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultFeatureDialogState)
    },
  }
}
