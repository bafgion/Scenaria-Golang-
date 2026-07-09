import { describe, expect, it } from 'vitest'
import { createRecorderStore, defaultRecorderState } from './recorderStore'

function currentValue<T>(store: { subscribe: (run: (value: T) => void) => () => void }): T {
  let out!: T
  const unsubscribe = store.subscribe((value) => {
    out = value
  })
  unsubscribe()
  return out
}

describe('recorderStore', () => {
  it('updates browser and recorder state', () => {
    const store = createRecorderStore()
    store.setSession('record-1', 'browser-1')
    store.setTargetPath('features/a.feature')
    store.setBrowserState(true, true, false)
    expect(currentValue(store)).toEqual({
      browserOpen: true,
      recording: true,
      paused: false,
      targetPath: 'features/a.feature',
      recordSessionId: 'record-1',
      browserSessionId: 'browser-1',
      liveRecordStepLines: {},
      lastRecordTarget: '',
      pauseToggleGuardUntil: 0,
    })
  })

  it('resets to defaults', () => {
    const store = createRecorderStore({
      browserOpen: true,
      recording: true,
      paused: true,
      targetPath: 'x.feature',
      recordSessionId: 'record-x',
      browserSessionId: 'browser-x',
      liveRecordStepLines: { 1: 2 },
      lastRecordTarget: 'x.feature',
      pauseToggleGuardUntil: 100,
    })
    store.reset()
    expect(currentValue(store)).toEqual(defaultRecorderState)
  })
})

