import { writable } from 'svelte/store'
import type { gui } from '../../wailsjs/go/models'

export type DiagnosticsState = {
  issues: gui.ValidationIssue[]
  hints: gui.ScenarioHintDTO[]
}

export const defaultDiagnosticsState: DiagnosticsState = {
  issues: [],
  hints: [],
}

export function createDiagnosticsStore(initial: DiagnosticsState = defaultDiagnosticsState) {
  const store = writable<DiagnosticsState>(initial)
  return {
    subscribe: store.subscribe,
    setIssues(issues: gui.ValidationIssue[]) {
      store.update((s) => ({ ...s, issues }))
    },
    setHints(hints: gui.ScenarioHintDTO[]) {
      store.update((s) => ({ ...s, hints }))
    },
    reset() {
      store.set(defaultDiagnosticsState)
    },
  }
}
