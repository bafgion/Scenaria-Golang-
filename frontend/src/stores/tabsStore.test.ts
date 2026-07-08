import { describe, expect, it } from 'vitest'
import { reduceTabsAfterClose } from './tabsStore'

describe('tabsStore reduceTabsAfterClose', () => {
  it('keeps active tab when closing inactive tab', () => {
    const state = reduceTabsAfterClose(
      [
        { path: 'a.feature', content: '', dirty: false },
        { path: 'b.feature', content: '', dirty: false },
      ],
      'a.feature',
      'b.feature',
    )
    expect(state.tabs.map((t) => t.path)).toEqual(['a.feature'])
    expect(state.openNextPath).toBe('')
    expect(state.showWelcome).toBe(false)
  })

  it('returns last tab as next when active is closed', () => {
    const state = reduceTabsAfterClose(
      [
        { path: 'a.feature', content: '', dirty: false },
        { path: 'b.feature', content: '', dirty: false },
      ],
      'b.feature',
      'b.feature',
    )
    expect(state.tabs.map((t) => t.path)).toEqual(['a.feature'])
    expect(state.openNextPath).toBe('a.feature')
    expect(state.showWelcome).toBe(false)
  })

  it('falls back to welcome when no tabs left', () => {
    const state = reduceTabsAfterClose(
      [{ path: 'a.feature', content: '', dirty: false }],
      'a.feature',
      'a.feature',
    )
    expect(state.tabs).toEqual([])
    expect(state.openNextPath).toBe('')
    expect(state.showWelcome).toBe(true)
  })
})
