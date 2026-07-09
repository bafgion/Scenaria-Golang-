import { writable } from 'svelte/store'
import type { gui } from '../../wailsjs/go/models'

export type ReportsState = {
  runResults: gui.RunResultEntry[]
  flakyMetrics: gui.FlakyMetricsDTO | null
  allureInstalled: boolean
  allureServeRunning: boolean
  refreshInFlight: boolean
  refreshQueued: boolean
}

export const defaultReportsState: ReportsState = {
  runResults: [],
  flakyMetrics: null,
  allureInstalled: true,
  allureServeRunning: false,
  refreshInFlight: false,
  refreshQueued: false,
}

export function createReportsStore(initial: ReportsState = defaultReportsState) {
  const store = writable<ReportsState>(initial)
  return {
    subscribe: store.subscribe,
    setData(runResults: gui.RunResultEntry[], flakyMetrics: gui.FlakyMetricsDTO | null) {
      store.update((s) => ({ ...s, runResults, flakyMetrics }))
    },
    setRunResults(runResults: gui.RunResultEntry[]) {
      store.update((s) => ({ ...s, runResults }))
    },
    setFlakyMetrics(flakyMetrics: gui.FlakyMetricsDTO | null) {
      store.update((s) => ({ ...s, flakyMetrics }))
    },
    setAllureStatus(allureInstalled: boolean, allureServeRunning: boolean) {
      store.update((s) => ({ ...s, allureInstalled, allureServeRunning }))
    },
    setAllureServeRunning(allureServeRunning: boolean) {
      store.update((s) => ({ ...s, allureServeRunning }))
    },
    tryBeginRefresh(): boolean {
      let started = false
      store.update((s) => {
        if (s.refreshInFlight) {
          return { ...s, refreshQueued: true }
        }
        started = true
        return { ...s, refreshInFlight: true }
      })
      return started
    },
    finishRefresh(): boolean {
      let rerun = false
      store.update((s) => {
        rerun = s.refreshQueued
        return { ...s, refreshInFlight: false, refreshQueued: false }
      })
      return rerun
    },
    queueRefresh() {
      store.update((s) => ({ ...s, refreshQueued: true }))
    },
    snapshot(): ReportsState {
      let state = defaultReportsState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultReportsState)
    },
  }
}
