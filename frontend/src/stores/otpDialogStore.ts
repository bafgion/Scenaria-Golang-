import { writable } from 'svelte/store'

export type OtpDialogState = {
  email: string
}

export const defaultOtpDialogState: OtpDialogState = {
  email: '',
}

export function createOtpDialogStore(initial: OtpDialogState = defaultOtpDialogState) {
  const store = writable<OtpDialogState>(initial)
  return {
    subscribe: store.subscribe,
    setEmail(email: string) {
      store.set({ email })
    },
    snapshot(): OtpDialogState {
      let state = defaultOtpDialogState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultOtpDialogState)
    },
  }
}
