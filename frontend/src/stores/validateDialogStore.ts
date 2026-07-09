import { writable } from 'svelte/store'

export type ValidateDialogState = {
  browser: string
  syntaxOnly: boolean
  scope: 'project' | 'current'
  cliLog: string
}

export const defaultValidateDialogState: ValidateDialogState = {
  browser: 'chromium',
  syntaxOnly: false,
  scope: 'project',
  cliLog: '',
}

export function createValidateDialogStore(initial: ValidateDialogState = defaultValidateDialogState) {
  const store = writable<ValidateDialogState>(initial)
  return {
    subscribe: store.subscribe,
    patch(partial: Partial<ValidateDialogState>) {
      store.update((s) => ({ ...s, ...partial }))
    },
    open(syntaxOnly: boolean, browser: string, scope: 'project' | 'current') {
      store.set({
        browser,
        syntaxOnly,
        scope,
        cliLog: '',
      })
    },
    setCliLog(cliLog: string) {
      store.update((s) => ({ ...s, cliLog }))
    },
    appendCliLog(line: string) {
      store.update((s) => ({
        ...s,
        cliLog: s.cliLog ? `${s.cliLog}\n${line}`.trim() : line,
      }))
    },
    snapshot(): ValidateDialogState {
      let state = defaultValidateDialogState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultValidateDialogState)
    },
  }
}
