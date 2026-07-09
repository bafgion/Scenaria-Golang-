import { writable } from 'svelte/store'

export type ConfirmDialogRequest = {
  title: string
  message: string
  confirmLabel: string
  danger: boolean
  dontAskAgainLabel?: string
}

export type ConfirmDialogState = {
  open: boolean
  request: ConfirmDialogRequest | null
  skipRecordTabSwitchConfirm: boolean
}

export const defaultConfirmDialogState: ConfirmDialogState = {
  open: false,
  request: null,
  skipRecordTabSwitchConfirm: false,
}

export type ConfirmDialogOptions = {
  title: string
  message: string
  confirmLabel?: string
  danger?: boolean
  dontAskAgainLabel?: string
}

export function createConfirmDialogStore(initial: ConfirmDialogState = defaultConfirmDialogState) {
  const store = writable<ConfirmDialogState>(initial)
  let pendingResolve: ((confirmed: boolean) => void) | null = null

  return {
    subscribe: store.subscribe,
    ask(opts: ConfirmDialogOptions, labels: { ok: string }): Promise<boolean> {
      return new Promise((resolve) => {
        pendingResolve = resolve
        store.set({
          open: true,
          request: {
            title: opts.title,
            message: opts.message,
            confirmLabel: opts.confirmLabel || labels.ok,
            danger: !!opts.danger,
            dontAskAgainLabel: opts.dontAskAgainLabel,
          },
          skipRecordTabSwitchConfirm: this.snapshot().skipRecordTabSwitchConfirm,
        })
      })
    },
    close(confirmed: boolean, dontAskAgain = false) {
      if (confirmed && dontAskAgain) {
        store.update((s) => ({ ...s, skipRecordTabSwitchConfirm: true }))
      }
      pendingResolve?.(confirmed)
      pendingResolve = null
      store.update((s) => ({ ...s, open: false, request: null }))
    },
    shouldSkipRecordTabSwitchConfirm(): boolean {
      return this.snapshot().skipRecordTabSwitchConfirm
    },
    snapshot(): ConfirmDialogState {
      let state = defaultConfirmDialogState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      pendingResolve = null
      store.set(defaultConfirmDialogState)
    },
  }
}
