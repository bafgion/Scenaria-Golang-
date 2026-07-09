import { writable } from 'svelte/store'

export type RecordFormState = {
  recordURL: string
  recordOutput: string
  recordIdle: number
  recordAppendTo: string
  recordTestClient: string
  recordFeatureName: string
  recordScenarioName: string
  recordMode: 'live' | 'baseline'
  baselineBusy: boolean
  recordStepPickerOpen: boolean
}

export const defaultRecordFormState: RecordFormState = {
  recordURL: '',
  recordOutput: 'recorded.feature',
  recordIdle: 30,
  recordAppendTo: '',
  recordTestClient: '',
  recordFeatureName: '',
  recordScenarioName: '',
  recordMode: 'live',
  baselineBusy: false,
  recordStepPickerOpen: false,
}

export function createRecordFormStore(initial: Partial<RecordFormState> = {}) {
  const store = writable<RecordFormState>({ ...defaultRecordFormState, ...initial })
  return {
    subscribe: store.subscribe,
    patch(partial: Partial<RecordFormState>) {
      store.update((state) => ({ ...state, ...partial }))
    },
    reset(partial: Partial<RecordFormState> = {}) {
      store.set({ ...defaultRecordFormState, ...partial })
    },
    snapshot(): RecordFormState {
      let state = defaultRecordFormState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
  }
}
