import { writable } from 'svelte/store'
import type { gui } from '../../wailsjs/go/models'

export type DiagnosticsState = {
  issues: gui.ValidationIssue[]
  hints: gui.ScenarioHintDTO[]
  issuesByTab: Record<string, gui.ValidationIssue[]>
  browserPanelIssues: gui.ValidationIssue[]
  hintsDismissed: string[]
  validateGeneration: number
  stepStatusError: boolean
  hintFixInFlight: boolean
}

export const defaultDiagnosticsState: DiagnosticsState = {
  issues: [],
  hints: [],
  issuesByTab: {},
  browserPanelIssues: [],
  hintsDismissed: [],
  validateGeneration: 0,
  stepStatusError: false,
  hintFixInFlight: false,
}

export function createDiagnosticsStore(initial: DiagnosticsState = defaultDiagnosticsState) {
  const store = writable<DiagnosticsState>(initial)
  let validateDebounceTimer: ReturnType<typeof setTimeout> | null = null
  return {
    subscribe: store.subscribe,
    setIssues(issues: gui.ValidationIssue[]) {
      store.update((s) => ({ ...s, issues }))
    },
    setHints(hints: gui.ScenarioHintDTO[]) {
      store.update((s) => ({ ...s, hints }))
    },
    setBrowserPanelIssues(issues: gui.ValidationIssue[]) {
      store.update((s) => ({ ...s, browserPanelIssues: issues }))
    },
    setStepStatusError(stepStatusError: boolean) {
      store.update((s) => ({ ...s, stepStatusError }))
    },
    setHintFixInFlight(hintFixInFlight: boolean) {
      store.update((s) => ({ ...s, hintFixInFlight }))
    },
    scheduleValidateEditor(run: () => void, delayMs = 300) {
      if (validateDebounceTimer) clearTimeout(validateDebounceTimer)
      if (delayMs <= 0) {
        validateDebounceTimer = null
        run()
        return
      }
      validateDebounceTimer = setTimeout(() => {
        validateDebounceTimer = null
        run()
      }, delayMs)
    },
    clearValidateDebounce() {
      if (validateDebounceTimer) clearTimeout(validateDebounceTimer)
      validateDebounceTimer = null
    },
    setIssuesForTab(tabPath: string, issues: gui.ValidationIssue[]) {
      store.update((s) => ({
        ...s,
        issuesByTab: { ...s.issuesByTab, [tabPath]: issues },
      }))
    },
    clearIssuesForTab(tabPath: string) {
      store.update((s) => {
        if (!(tabPath in s.issuesByTab)) return s
        const { [tabPath]: _removed, ...issuesByTab } = s.issuesByTab
        return { ...s, issuesByTab }
      })
    },
    clearIssuesByTab() {
      store.update((s) => ({ ...s, issuesByTab: {} }))
    },
    isHintDismissed(key: string): boolean {
      return this.snapshot().hintsDismissed.includes(key)
    },
    dismissHint(key: string) {
      store.update((s) => {
        if (s.hintsDismissed.includes(key)) return s
        return { ...s, hintsDismissed: [...s.hintsDismissed, key] }
      })
    },
    clearDismissedHints() {
      store.update((s) => ({ ...s, hintsDismissed: [] }))
    },
    bumpValidateGeneration(): number {
      let next = 0
      store.update((s) => {
        next = s.validateGeneration + 1
        return { ...s, validateGeneration: next }
      })
      return next
    },
    validateGeneration(): number {
      return this.snapshot().validateGeneration
    },
    issuesForTab(tabPath: string): gui.ValidationIssue[] {
      return this.snapshot().issuesByTab[tabPath] ?? []
    },
    snapshot(): DiagnosticsState {
      let state = defaultDiagnosticsState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      this.clearValidateDebounce()
      store.set(defaultDiagnosticsState)
    },
  }
}
