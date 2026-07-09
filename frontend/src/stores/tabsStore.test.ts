import { describe, expect, it } from 'vitest'
import { createTabsStore, reduceTabsAfterClose } from './tabsStore'

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

describe('createTabsStore', () => {
  const welcomeKey = '__welcome__'

  it('resets to welcome state', () => {
    const store = createTabsStore(welcomeKey)
    store.appendTab({ path: 'a.feature', content: 'x', dirty: false })
    store.setActiveTab('a.feature')
    store.setWelcomeVisible(false)
    store.reset()
    let snapshot = { tabs: [] as { path: string; content: string; dirty: boolean }[], activeTab: '', welcomeTabVisible: false }
    const unsub = store.subscribe((s) => {
      snapshot = s
    })
    unsub()
    expect(snapshot.tabs).toEqual([])
    expect(snapshot.activeTab).toBe(welcomeKey)
    expect(snapshot.welcomeTabVisible).toBe(true)
  })

  it('applies close result with welcome fallback', () => {
    const store = createTabsStore(welcomeKey)
    store.appendTab({ path: 'a.feature', content: '', dirty: false })
    store.setActiveTab('a.feature')
    store.applyCloseResult({
      tabs: [],
      openNextPath: '',
      showWelcome: true,
    })
    let activeTab = ''
    let welcomeTabVisible = false
    store.subscribe((s) => {
      activeTab = s.activeTab
      welcomeTabVisible = s.welcomeTabVisible
    })()
    expect(activeTab).toBe(welcomeKey)
    expect(welcomeTabVisible).toBe(true)
  })

  it('activates next tab immediately after closing active tab', () => {
    const store = createTabsStore(welcomeKey)
    store.appendTab({ path: 'a.feature', content: '', dirty: false })
    store.appendTab({ path: 'b.feature', content: '', dirty: false })
    store.setActiveTab('b.feature')
    store.setWelcomeVisible(false)
    store.applyCloseResult({
      tabs: [{ path: 'a.feature', content: '', dirty: false }],
      openNextPath: 'a.feature',
      showWelcome: false,
    })
    const snapshot = store.snapshot()
    expect(snapshot.tabs.map((t) => t.path)).toEqual(['a.feature'])
    expect(snapshot.activeTab).toBe('a.feature')
    expect(snapshot.welcomeTabVisible).toBe(false)
    expect(snapshot.pendingCloseTab).toBeNull()
  })

  it('closes pending dirty active tab using current store state', () => {
    const store = createTabsStore(welcomeKey)
    store.appendTab({ path: 'a.feature', content: '', dirty: false })
    store.appendTab({ path: 'b.feature', content: 'old', draft: 'new', dirty: true })
    store.setActiveTab('b.feature')
    store.setWelcomeVisible(false)
    store.setPendingCloseTab('b.feature')

    const result = store.closePath('b.feature')
    const snapshot = store.snapshot()

    expect(result.openNextPath).toBe('a.feature')
    expect(snapshot.tabs.map((t) => t.path)).toEqual(['a.feature'])
    expect(snapshot.activeTab).toBe('a.feature')
    expect(snapshot.pendingCloseTab).toBeNull()
  })
})
