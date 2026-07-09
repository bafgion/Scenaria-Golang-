import { writable } from 'svelte/store'

export type DialogsState = {
  showSettings: boolean
  showCommandPalette: boolean
  showSnippetPalette: boolean
  showRecord: boolean
  showPluginRun: boolean
  showOtp: boolean
  showRun: boolean
  showTestClient: boolean
  showExport: boolean
  showImport: boolean
  showImportFeatures: boolean
  showDuplicateFeature: boolean
  showRefactorUrl: boolean
  showOpenProject: boolean
  showRenameFeature: boolean
  showMoveFeature: boolean
  showValidate: boolean
  showUpdateCheck: boolean
  showInitProject: boolean
  showNewProjectWizard: boolean
  showSteps: boolean
  showStepsHelp: boolean
  showAbout: boolean
  showPlugins: boolean
  showProjectReplace: boolean
  showHotkeys: boolean
  showRunHistory: boolean
  showVanessaRun: boolean
  showVanessaSettings: boolean
  showVanessaMonitor: boolean
  showHttpAuth: boolean
  showPickerStep: boolean
  showPostRecordDiff: boolean
}

export const defaultDialogsState: DialogsState = {
  showSettings: false,
  showCommandPalette: false,
  showSnippetPalette: false,
  showRecord: false,
  showPluginRun: false,
  showOtp: false,
  showRun: false,
  showTestClient: false,
  showExport: false,
  showImport: false,
  showImportFeatures: false,
  showDuplicateFeature: false,
  showRefactorUrl: false,
  showOpenProject: false,
  showRenameFeature: false,
  showMoveFeature: false,
  showValidate: false,
  showUpdateCheck: false,
  showInitProject: false,
  showNewProjectWizard: false,
  showSteps: false,
  showStepsHelp: false,
  showAbout: false,
  showPlugins: false,
  showProjectReplace: false,
  showHotkeys: false,
  showRunHistory: false,
  showVanessaRun: false,
  showVanessaSettings: false,
  showVanessaMonitor: false,
  showHttpAuth: false,
  showPickerStep: false,
  showPostRecordDiff: false,
}

export type DialogOverlayContext = {
  confirmDialogOpen?: boolean
  pendingCloseTab?: string | null
  postRecordPath?: string
}

export function anyAppDialogOpen(dialogs: DialogsState, overlay: DialogOverlayContext = {}): boolean {
  if (overlay.confirmDialogOpen) return true
  if (overlay.pendingCloseTab) return true
  if (dialogs.showRun) return true
  if (dialogs.showVanessaRun) return true
  if (dialogs.showPluginRun) return true
  if (dialogs.showTestClient) return true
  if (dialogs.showStepsHelp) return true
  if (dialogs.showSteps) return true
  if (dialogs.showVanessaSettings) return true
  if (dialogs.showExport) return true
  if (dialogs.showRefactorUrl) return true
  if (dialogs.showOpenProject) return true
  if (dialogs.showRenameFeature) return true
  if (dialogs.showMoveFeature) return true
  if (dialogs.showValidate) return true
  if (dialogs.showNewProjectWizard) return true
  if (dialogs.showInitProject) return true
  if (dialogs.showUpdateCheck) return true
  if (dialogs.showDuplicateFeature) return true
  if (dialogs.showImportFeatures) return true
  if (dialogs.showImport) return true
  if (dialogs.showSettings) return true
  if (dialogs.showCommandPalette) return true
  if (dialogs.showSnippetPalette) return true
  if (dialogs.showRecord) return true
  if (dialogs.showOtp) return true
  if (dialogs.showAbout) return true
  if (dialogs.showHotkeys) return true
  if (dialogs.showPlugins) return true
  if (dialogs.showRunHistory) return true
  if (dialogs.showPostRecordDiff && !!overlay.postRecordPath) return true
  if (dialogs.showProjectReplace) return true
  if (dialogs.showHttpAuth) return true
  if (dialogs.showPickerStep) return true
  return false
}

export function createDialogsStore(initial: DialogsState = defaultDialogsState) {
  const store = writable<DialogsState>(initial)
  return {
    subscribe: store.subscribe,
    patch(partial: Partial<DialogsState>) {
      store.update((state) => ({ ...state, ...partial }))
    },
    open(key: keyof DialogsState) {
      store.update((state) => ({ ...state, [key]: true }))
    },
    close(key: keyof DialogsState) {
      store.update((state) => ({ ...state, [key]: false }))
    },
    reset() {
      store.set(defaultDialogsState)
    },
    anyOpen(overlay: DialogOverlayContext = {}) {
      let state = defaultDialogsState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return anyAppDialogOpen(state, overlay)
    },
  }
}
