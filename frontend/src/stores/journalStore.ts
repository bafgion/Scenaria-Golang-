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

function redactJournalState(state: JournalState): JournalState {
  return {
    ...state,
    logText: redactSecrets(state.logText),
    statusMessage: redactSecrets(state.statusMessage),
  }
}

function redactJournalPatch(partial: Partial<JournalState>): Partial<JournalState> {
  const safe = { ...partial }
  if (safe.logText !== undefined) {
    safe.logText = redactSecrets(safe.logText)
  }
  if (safe.statusMessage !== undefined) {
    safe.statusMessage = redactSecrets(safe.statusMessage)
  }
  return safe
}

export function createJournalStore(initial: JournalState = defaultJournalState) {
  const store = writable<JournalState>(redactJournalState(initial))
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
      const safePartial = redactJournalPatch(partial)
      store.update((s) => ({ ...s, ...safePartial }))
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
