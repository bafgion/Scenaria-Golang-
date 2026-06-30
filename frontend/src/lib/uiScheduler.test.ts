import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { debounce, deferToNextFrame } from './uiScheduler'

describe('uiScheduler', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('debounce delays invocation', () => {
    const fn = vi.fn()
    const debounced = debounce(fn, 200)
    debounced('a')
    debounced('b')
    expect(fn).not.toHaveBeenCalled()
    vi.advanceTimersByTime(199)
    expect(fn).not.toHaveBeenCalled()
    vi.advanceTimersByTime(1)
    expect(fn).toHaveBeenCalledOnce()
    expect(fn).toHaveBeenCalledWith('b')
  })

  it('debounce flush runs pending call immediately', () => {
    const fn = vi.fn()
    const debounced = debounce(fn, 500)
    debounced('x')
    debounced.flush()
    expect(fn).toHaveBeenCalledWith('x')
  })

  it('debounce cancel drops pending call', () => {
    const fn = vi.fn()
    const debounced = debounce(fn, 100)
    debounced('y')
    debounced.cancel()
    vi.advanceTimersByTime(200)
    expect(fn).not.toHaveBeenCalled()
  })

  it('deferToNextFrame schedules on next animation frame', () => {
    const raf = vi.fn((cb: FrameRequestCallback) => {
      cb(0)
      return 1
    })
    vi.stubGlobal('requestAnimationFrame', raf)
    const work = vi.fn()
    deferToNextFrame(work)
    expect(raf).toHaveBeenCalled()
    expect(work).toHaveBeenCalledOnce()
    vi.unstubAllGlobals()
  })
})
