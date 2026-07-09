import { writable } from 'svelte/store'
import { defaultRunForm, type RunForm } from '../lib/runTypes'

export type RunFormState = {
  lastRun: RunForm
  runForm: RunForm
}

export function createRunFormStore(
  initialLastRun: RunForm = defaultRunForm({ headed: true, installPW: true, html: true, htmlLightMode: true }),
) {
  const store = writable<RunFormState>({
    lastRun: initialLastRun,
    runForm: { ...initialLastRun },
  })
  return {
    subscribe: store.subscribe,
    setLastRun(lastRun: RunForm) {
      store.update((s) => ({ ...s, lastRun }))
    },
    setRunForm(runForm: RunForm) {
      store.update((s) => ({ ...s, runForm }))
    },
    patchLastRun(partial: Partial<RunForm>) {
      store.update((s) => ({ ...s, lastRun: { ...s.lastRun, ...partial } }))
    },
    patchRunForm(partial: Partial<RunForm>) {
      store.update((s) => ({ ...s, runForm: { ...s.runForm, ...partial } }))
    },
    applyLastRunFromSettings(browser: string, workers: number, slowMo: number) {
      store.update((s) => ({
        ...s,
        lastRun: {
          ...s.lastRun,
          workers,
          slowMo,
          browser,
        },
      }))
    },
    snapshot(): RunFormState {
      let state: RunFormState = { lastRun: defaultRunForm(), runForm: defaultRunForm() }
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset(initialLastRun?: RunForm) {
      const base = initialLastRun ?? defaultRunForm({ headed: true, installPW: true, html: true, htmlLightMode: true })
      store.set({ lastRun: base, runForm: { ...base } })
    },
  }
}
