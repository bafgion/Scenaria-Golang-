import { writable } from 'svelte/store'

export type RecorderState = {
  browserOpen: boolean
  recording: boolean
  paused: boolean
  targetPath: string
  recordSessionId: string
  browserSessionId: string
}

export const defaultRecorderState: RecorderState = {
  browserOpen: false,
  recording: false,
  paused: false,
  targetPath: '',
  recordSessionId: '',
  browserSessionId: '',
}

export function createRecorderStore(initial: RecorderState = defaultRecorderState) {
  const store = writable<RecorderState>(initial)
  return {
    subscribe: store.subscribe,
    setSession(recordSessionId: string, browserSessionId: string) {
      store.update((s) => ({ ...s, recordSessionId, browserSessionId }))
    },
    setSessionIDs(recordSessionId?: string, browserSessionId?: string) {
      store.update((s) => ({
        ...s,
        recordSessionId: recordSessionId && recordSessionId.trim() ? recordSessionId : s.recordSessionId,
        browserSessionId: browserSessionId && browserSessionId.trim() ? browserSessionId : s.browserSessionId,
      }))
    },
    setTargetPath(targetPath: string) {
      store.update((s) => ({ ...s, targetPath }))
    },
    setBrowserOpen(browserOpen: boolean) {
      store.update((s) => ({ ...s, browserOpen }))
    },
    setRecording(recording: boolean, paused = false) {
      store.update((s) => ({ ...s, recording, paused }))
    },
    setBrowserState(browserOpen: boolean, recording: boolean, paused = false) {
      store.update((s) => ({ ...s, browserOpen, recording, paused }))
    },
    reset() {
      store.set(defaultRecorderState)
    },
  }
}
