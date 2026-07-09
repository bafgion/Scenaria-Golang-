import { writable } from 'svelte/store'
import type { gui } from '../../wailsjs/go/models'

export type RunnerState = {
  playing: boolean
  runId: string
  total: number
  current: number
  label: string
  logStreaming: boolean
  cancelling: boolean
  lastRunSince: string | null
  lastRunBatchResults: gui.RunResultEntry[]
  lastErrorEntry: gui.RunResultEntry | null
  dryRunActive: boolean
}

export const defaultRunnerState: RunnerState = {
  playing: false,
  runId: '',
  total: 0,
  current: 0,
  label: '',
  logStreaming: false,
  cancelling: false,
  lastRunSince: null,
  lastRunBatchResults: [],
  lastErrorEntry: null,
  dryRunActive: false,
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
    setLastRunSession(
      lastRunSince: string | null,
      lastRunBatchResults: gui.RunResultEntry[],
      lastErrorEntry: gui.RunResultEntry | null,
    ) {
      store.update((s) => ({ ...s, lastRunSince, lastRunBatchResults, lastErrorEntry }))
    },
    setDryRunActive(dryRunActive: boolean) {
      store.update((s) => ({ ...s, dryRunActive }))
    },
    stop() {
      store.update((s) => ({ ...s, playing: false }))
    },
    reset() {
      store.set(defaultRunnerState)
    },
  }
}
