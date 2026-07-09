import { describe, expect, it } from 'vitest'
import { createSplashStore } from './splashStore'

describe('splashStore', () => {
  it('tracks splash stages and ready state', () => {
    const store = createSplashStore()
    store.setStage('Loading', 50)
    expect(store.snapshot()).toMatchObject({ message: 'Loading', progress: 50, appReady: false })
    store.markReady()
    expect(store.snapshot().appReady).toBe(true)
  })
})
