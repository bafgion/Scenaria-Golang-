import { writable } from 'svelte/store'

export type PostRecordState = {
  path: string
  stepCount: number
  baselineText: string
}

export const defaultPostRecordState: PostRecordState = {
  path: '',
  stepCount: 0,
  baselineText: '',
}

export function createPostRecordStore(initial: PostRecordState = defaultPostRecordState) {
  const store = writable<PostRecordState>(initial)
  return {
    subscribe: store.subscribe,
    open(path: string, stepCount = 0) {
      store.update((s) => ({ ...s, path, stepCount }))
    },
    setStepCount(stepCount: number) {
      store.update((s) => ({ ...s, stepCount }))
    },
    setBaselineText(baselineText: string) {
      store.update((s) => ({ ...s, baselineText }))
    },
    dismiss() {
      store.set(defaultPostRecordState)
    },
    snapshot(): PostRecordState {
      let state = defaultPostRecordState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultPostRecordState)
    },
  }
}
