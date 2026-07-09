import { describe, expect, it, vi } from 'vitest'
import { createAppMetaStore } from './appMetaStore'
import { createSessionStore } from './sessionStore'

describe('appMetaStore', () => {
  it('stores app version string', () => {
    const store = createAppMetaStore()
    store.setVersion('1.2.3')
    expect(store.snapshot().version).toBe('1.2.3')
  })
})

describe('sessionStore', () => {
  it('schedules and flushes persist callback', () => {
    vi.useFakeTimers()
    const store = createSessionStore()
    const run = vi.fn()
    store.schedulePersist(run, 500)
    vi.advanceTimersByTime(499)
    expect(run).not.toHaveBeenCalled()
    store.flushPersist(run)
    expect(run).toHaveBeenCalledTimes(1)
    vi.advanceTimersByTime(10)
    expect(run).toHaveBeenCalledTimes(1)
    vi.useRealTimers()
  })
})
