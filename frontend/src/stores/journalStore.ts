import { writable } from 'svelte/store'
import { redactSecrets } from '../lib/redaction'

export type StatusTone = 'normal' | 'error' | 'success' | 'busy'

export type JournalState = {
  logText: string
  statusMessage: string
  statusTone: StatusTone
}

export const defaultJournalState: JournalState = {
  logText: '',
  statusMessage: '',
  statusTone: 'normal',
}

export function createJournalStore(initial: JournalState = defaultJournalState) {
  const store = writable<JournalState>(initial)
  return {
    subscribe: store.subscribe,
    appendLog(line: string) {
      const safeLine = redactSecrets(line)
      store.update((s) => ({
        ...s,
        logText: s.logText + safeLine + (safeLine.endsWith('\n') ? '' : '\n'),
      }))
    },
    setStatus(msg: string, tone: StatusTone = 'normal') {
      store.update((s) => ({ ...s, statusMessage: redactSecrets(msg), statusTone: tone }))
    },
    patch(partial: Partial<JournalState>) {
      store.update((s) => ({ ...s, ...partial }))
    },
    clearLog() {
      store.update((s) => ({ ...s, logText: '' }))
    },
    snapshot(): JournalState {
      let state = defaultJournalState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultJournalState)
    },
  }
}
