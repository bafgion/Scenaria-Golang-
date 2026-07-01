import { describe, expect, it } from 'vitest'
import { MonacoTabViewStateStore } from './monacoTabViewState'

describe('MonacoTabViewStateStore', () => {
  it('captures and restores scroll and position', () => {
    const store = new MonacoTabViewStateStore()
    let scrollTop = 0
    let position = { lineNumber: 1, column: 1 }
    const editor = {
      getScrollTop: () => scrollTop,
      setScrollTop: (v: number) => {
        scrollTop = v
      },
      getPosition: () => position,
      setPosition: (p: typeof position) => {
        position = p
      },
      revealPositionInCenterIfOutsideViewport: () => {},
    }

    scrollTop = 120
    position = { lineNumber: 8, column: 3 }
    store.capture(editor as never, '/a.feature')

    scrollTop = 0
    position = { lineNumber: 1, column: 1 }
    store.restore(editor as never, '/a.feature')

    expect(scrollTop).toBe(120)
    expect(position).toEqual({ lineNumber: 8, column: 3 })
  })

  it('drops state when tab is released', () => {
    const store = new MonacoTabViewStateStore()
    const editor = {
      getScrollTop: () => 40,
      setScrollTop: () => {},
      getPosition: () => ({ lineNumber: 2, column: 1 }),
      setPosition: () => {},
      revealPositionInCenterIfOutsideViewport: () => {},
    }
    store.capture(editor as never, '/gone.feature')
    store.drop('/gone.feature')
    let restored = false
    const track = {
      getScrollTop: () => 0,
      setScrollTop: () => {
        restored = true
      },
      getPosition: () => ({ lineNumber: 1, column: 1 }),
      setPosition: () => {},
      revealPositionInCenterIfOutsideViewport: () => {},
    }
    store.restore(track as never, '/gone.feature')
    expect(restored).toBe(false)
  })
})
