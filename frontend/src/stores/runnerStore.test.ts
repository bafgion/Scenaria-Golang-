import { describe, expect, it } from 'vitest'
import { createRunnerStore, defaultRunnerState } from './runnerStore'

function currentValue<T>(store: { subscribe: (run: (value: T) => void) => () => void }): T {
  let out!: T
  const unsubscribe = store.subscribe((value) => {
    out = value
  })
  unsubscribe()
  return out
}

describe('runnerStore', () => {
  it('tracks run identity and flags', () => {
    const store = createRunnerStore()
    store.start('run-1')
    store.progress(10, 3, 'running')
    store.setLogStreaming(true)
    store.setCancelling(true)
    expect(currentValue(store)).toEqual({
      playing: true,
      runId: 'run-1',
      total: 10,
      current: 3,
      label: 'running',
      logStreaming: true,
      cancelling: true,
      lastRunSince: null,
      lastRunBatchResults: [],
      lastErrorEntry: null,
      dryRunActive: false,
    })
  })

  it('resets to defaults', () => {
    const store = createRunnerStore({
      playing: true,
      runId: 'run-x',
      total: 1,
      current: 1,
      label: 'x',
      logStreaming: true,
      cancelling: true,
      lastRunSince: 't',
      lastRunBatchResults: [],
      lastErrorEntry: null,
      dryRunActive: false,
    })
    store.reset()
    expect(currentValue(store)).toEqual(defaultRunnerState)
  })
})

