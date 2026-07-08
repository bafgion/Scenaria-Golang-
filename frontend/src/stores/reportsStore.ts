import { writable } from 'svelte/store'
import type { gui } from '../../wailsjs/go/models'

export type ReportsState = {
  runResults: gui.RunResultEntry[]
  flakyMetrics: gui.FlakyMetricsDTO | null
}

export const defaultReportsState: ReportsState = {
  runResults: [],
  flakyMetrics: null,
}

export function createReportsStore(initial: ReportsState = defaultReportsState) {
  const store = writable<ReportsState>(initial)
  return {
    subscribe: store.subscribe,
    setData(runResults: gui.RunResultEntry[], flakyMetrics: gui.FlakyMetricsDTO | null) {
      store.set({ runResults, flakyMetrics })
    },
    setRunResults(runResults: gui.RunResultEntry[]) {
      store.update((s) => ({ ...s, runResults }))
    },
    setFlakyMetrics(flakyMetrics: gui.FlakyMetricsDTO | null) {
      store.update((s) => ({ ...s, flakyMetrics }))
    },
    reset() {
      store.set(defaultReportsState)
    },
  }
}
