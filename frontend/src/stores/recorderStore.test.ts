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
      captureFinalizing: false,
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
      captureFinalizing: true,
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

  it('stops capture without losing the live browser session', () => {
    const store = createRecorderStore({
      browserOpen: true,
      recording: true,
      paused: true,
      captureFinalizing: false,
      targetPath: 'x.feature',
      recordSessionId: 'record-x',
      browserSessionId: 'browser-x',
      liveRecordStepLines: { 1: 2 },
      lastRecordTarget: 'x.feature',
      pauseToggleGuardUntil: 100,
    })

    store.stopCaptureKeepBrowserOpen()

    expect(currentValue(store)).toEqual({
      browserOpen: true,
      recording: false,
      paused: false,
      captureFinalizing: false,
      targetPath: 'x.feature',
      recordSessionId: 'record-x',
      browserSessionId: 'browser-x',
      liveRecordStepLines: {},
      lastRecordTarget: 'x.feature',
      pauseToggleGuardUntil: 0,
    })
  })

  it('finalizes capture after queued steps drain', async () => {
    const store = createRecorderStore({
      browserOpen: true,
      recording: true,
      paused: false,
      captureFinalizing: false,
      targetPath: 'x.feature',
      recordSessionId: 'record-x',
      browserSessionId: 'browser-x',
      liveRecordStepLines: { 0: 3 },
      lastRecordTarget: 'x.feature',
      pauseToggleGuardUntil: 0,
    })

    store.beginCaptureFinalize()
    expect(currentValue(store).captureFinalizing).toBe(true)
    expect(currentValue(store).recording).toBe(false)
    expect(currentValue(store).targetPath).toBe('x.feature')

    let drained = false
    void store.chainRecordStepApply(async () => {
      drained = true
    })
    await store.awaitRecordStepApplyChain()
    expect(drained).toBe(true)

    store.stopCaptureKeepBrowserOpen()
    expect(currentValue(store).captureFinalizing).toBe(false)
    expect(currentValue(store).targetPath).toBe('x.feature')
  })
})
