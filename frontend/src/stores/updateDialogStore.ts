import { writable } from 'svelte/store'
import { gui } from '../../wailsjs/go/models'

export type UpdateDialogState = {
  message: string
  hasUpdate: boolean
  info: gui.UpdateInfoDTO | null
  downloading: boolean
  progress: gui.UpdateProgressDTO | null
  pendingStartupCheck: boolean
}

export const defaultUpdateDialogState: UpdateDialogState = {
  message: '',
  hasUpdate: false,
  info: null,
  downloading: false,
  progress: null,
  pendingStartupCheck: false,
}

export function createUpdateDialogStore(initial: UpdateDialogState = defaultUpdateDialogState) {
  const store = writable<UpdateDialogState>(initial)
  return {
    subscribe: store.subscribe,
    applyCheckResult(info: gui.UpdateInfoDTO) {
      store.update((s) => ({
        ...s,
        info,
        message: info.message || '',
        hasUpdate: !!info.updateAvailable,
      }))
    },
    setMessage(message: string, hasUpdate = false) {
      store.update((s) => ({ ...s, message, hasUpdate }))
    },
    appendMessage(line: string) {
      store.update((s) => ({
        ...s,
        message: s.message ? `${s.message}\n\n${line}`.trim() : line,
      }))
    },
    setDownloading(downloading: boolean) {
      store.update((s) => ({ ...s, downloading }))
    },
    setProgress(progress: gui.UpdateProgressDTO | null) {
      store.update((s) => ({ ...s, progress }))
    },
    setPendingStartupCheck(pendingStartupCheck: boolean) {
      store.update((s) => ({ ...s, pendingStartupCheck }))
    },
    patch(partial: Partial<UpdateDialogState>) {
      store.update((s) => ({ ...s, ...partial }))
    },
    snapshot(): UpdateDialogState {
      let state = defaultUpdateDialogState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultUpdateDialogState)
    },
  }
}
