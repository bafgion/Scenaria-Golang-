import { writable } from 'svelte/store'

export type RecorderState = {
  browserOpen: boolean
  recording: boolean
  paused: boolean
  targetPath: string
  recordSessionId: string
  browserSessionId: string
  liveRecordStepLines: Record<number, number>
  lastRecordTarget: string
  pauseToggleGuardUntil: number
}

export const defaultRecorderState: RecorderState = {
  browserOpen: false,
  recording: false,
  paused: false,
  targetPath: '',
  recordSessionId: '',
  browserSessionId: '',
  liveRecordStepLines: {},
  lastRecordTarget: '',
  pauseToggleGuardUntil: 0,
}

export function createRecorderStore(initial: RecorderState = defaultRecorderState) {
  const store = writable<RecorderState>(initial)
  let browserWatchTimer: ReturnType<typeof setInterval> | null = null
  let recordEditorReadyPromise: Promise<void> = Promise.resolve()
  let recordStepApplyChain: Promise<void> = Promise.resolve()
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
    setLiveRecordStepLines(liveRecordStepLines: Record<number, number>) {
      store.update((s) => ({ ...s, liveRecordStepLines }))
    },
    setLastRecordTarget(lastRecordTarget: string) {
      store.update((s) => ({ ...s, lastRecordTarget }))
    },
    clearLiveRecordSession() {
      store.update((s) => ({ ...s, liveRecordStepLines: {}, lastRecordTarget: '' }))
    },
    extendPauseToggleGuard(ms = 900) {
      store.update((s) => ({ ...s, pauseToggleGuardUntil: Date.now() + ms }))
    },
    setRecordEditorReadyPromise(promise: Promise<void>) {
      recordEditorReadyPromise = promise
    },
    awaitRecordEditorReady() {
      return recordEditorReadyPromise
    },
    chainRecordStepApply(run: () => Promise<void>) {
      recordStepApplyChain = recordStepApplyChain.then(run)
      return recordStepApplyChain
    },
    awaitRecordStepApplyChain() {
      return recordStepApplyChain
    },
    resetRecordOrchestration() {
      recordEditorReadyPromise = Promise.resolve()
      recordStepApplyChain = Promise.resolve()
    },
    setBrowserWatchTimer(timer: ReturnType<typeof setInterval> | null) {
      if (browserWatchTimer) clearInterval(browserWatchTimer)
      browserWatchTimer = timer
    },
    clearBrowserWatchTimer() {
      if (browserWatchTimer) {
        clearInterval(browserWatchTimer)
        browserWatchTimer = null
      }
    },
    reset() {
      this.clearBrowserWatchTimer()
      this.resetRecordOrchestration()
      store.set(defaultRecorderState)
    },
  }
}
