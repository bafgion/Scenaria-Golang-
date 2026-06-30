import { afterEach, describe, expect, it, vi } from 'vitest'
import { callWailsWithTimeout, wailsReady } from './wailsTimeout'

describe('wailsTimeout', () => {
  afterEach(() => {
    vi.useRealTimers()
    delete (window as unknown as { go?: unknown }).go
  })

  it('wailsReady is false without window.go', () => {
    expect(wailsReady()).toBe(false)
  })

  it('returns null when wails is not ready', async () => {
    const result = await callWailsWithTimeout('Test', Promise.resolve('ok'))
    expect(result).toBeNull()
  })

  it('resolves when promise completes before timeout', async () => {
    ;(window as unknown as { go: object }).go = {}
    const result = await callWailsWithTimeout('Test', Promise.resolve(42), 1000)
    expect(result).toBe(42)
  })

  it('returns null on timeout', async () => {
    vi.useFakeTimers()
    ;(window as unknown as { go: object }).go = {}
    const pending = callWailsWithTimeout('Slow', new Promise<string>(() => {}), 50)
    await vi.advanceTimersByTimeAsync(60)
    await expect(pending).resolves.toBeNull()
  })
})
