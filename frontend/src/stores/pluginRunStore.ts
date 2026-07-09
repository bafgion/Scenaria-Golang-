import { writable } from 'svelte/store'

export type PluginRunState = {
  name: string
  dry: boolean
  tag: string
  scenario: string
  dialogScenarios: string[]
}

export const defaultPluginRunState: PluginRunState = {
  name: '',
  dry: false,
  tag: '',
  scenario: '',
  dialogScenarios: [],
}

export function createPluginRunStore(initial: PluginRunState = defaultPluginRunState) {
  const store = writable<PluginRunState>(initial)
  return {
    subscribe: store.subscribe,
    patch(partial: Partial<PluginRunState>) {
      store.update((s) => ({ ...s, ...partial }))
    },
    prepareDialog(opts: { name: string; dry?: boolean; scenario?: string; dialogScenarios?: string[] }) {
      store.set({
        name: opts.name,
        dry: !!opts.dry,
        tag: '',
        scenario: opts.scenario ?? '',
        dialogScenarios: opts.dialogScenarios ?? [],
      })
    },
    snapshot(): PluginRunState {
      let state = defaultPluginRunState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultPluginRunState)
    },
  }
}
