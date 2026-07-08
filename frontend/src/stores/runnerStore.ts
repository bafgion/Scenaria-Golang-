import { writable } from 'svelte/store'

export type RunnerState = {
  playing: boolean
  runId: string
  total: number
  current: number
  label: string
  logStreaming: boolean
  cancelling: boolean
}

export const defaultRunnerState: RunnerState = {
  playing: false,
  runId: '',
  total: 0,
  current: 0,
  label: '',
  logStreaming: false,
  cancelling: false,
}

export function createRunnerStore(initial: RunnerState = defaultRunnerState) {
  const store = writable<RunnerState>(initial)
  return {
    subscribe: store.subscribe,
    start(runId: string) {
      store.update((s) => ({ ...s, playing: true, runId, total: 0, current: 0, label: '', cancelling: false }))
    },
    progress(total: number, current: number, label: string) {
      store.update((s) => ({ ...s, total, current, label }))
    },
    setRunId(runId: string) {
      store.update((s) => ({ ...s, runId }))
    },
    setLogStreaming(logStreaming: boolean) {
      store.update((s) => ({ ...s, logStreaming }))
    },
    setCancelling(cancelling: boolean) {
      store.update((s) => ({ ...s, cancelling }))
    },
    stop() {
      store.update((s) => ({ ...s, playing: false }))
    },
    reset() {
      store.set(defaultRunnerState)
    },
  }
}
