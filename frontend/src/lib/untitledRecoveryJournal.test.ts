import { describe, expect, it } from 'vitest'
import {
  UNTITLED_RECOVERY_KEY,
  UNTITLED_RECOVERY_MAX_AGE_MS,
  clearUntitledRecoveryAll,
  clearUntitledRecoveryPath,
  journalUntitledTabs,
  loadUntitledRecoveryJournal,
  mergeUntitledSessionWithRecovery,
  persistUntitledRecoveryEntries,
} from './untitledRecoveryJournal'
import type { TabBody } from './tabMemory'

class MemoryStorage {
  values = new Map<string, string>()
  getItem(key: string): string | null {
    return this.values.get(key) ?? null
  }
  setItem(key: string, value: string) {
    this.values.set(key, value)
  }
  removeItem(key: string) {
    this.values.delete(key)
  }
}

describe('untitledRecoveryJournal', () => {
  it('persists active untitled text synchronously for crash recovery', () => {
    const storage = new MemoryStorage()
    const tabs: TabBody[] = [
      { path: '__untitled__:1/a.feature', content: 'old', draft: 'draft', dirty: true },
    ]

    expect(journalUntitledTabs(tabs, tabs[0].path, () => 'latest', storage, 100)).toBe(true)

    expect(loadUntitledRecoveryJournal(storage, 100)).toEqual([
      { path: '__untitled__:1/a.feature', content: 'latest', updatedAt: 100 },
    ])
  })

  it('preserves intentionally empty untitled edits', () => {
    const storage = new MemoryStorage()
    const tabs: TabBody[] = [
      { path: '__untitled__:1/a.feature', content: 'old', draft: '', dirty: true },
    ]

    journalUntitledTabs(tabs, tabs[0].path, () => '', storage, 200)

    expect(loadUntitledRecoveryJournal(storage, 200)).toEqual([
      { path: '__untitled__:1/a.feature', content: '', updatedAt: 200 },
    ])
  })

  it('keeps multiple untitled tabs isolated', () => {
    const storage = new MemoryStorage()
    const tabs: TabBody[] = [
      { path: '__untitled__:1/a.feature', content: 'one', dirty: true },
      { path: '__untitled__:2/b.feature', content: 'two', draft: 'two draft', dirty: true },
      { path: '/proj/saved.feature', content: 'saved', draft: 'saved draft', dirty: true },
    ]

    journalUntitledTabs(tabs, '__untitled__:2/b.feature', () => 'live two', storage, 300)

    expect(loadUntitledRecoveryJournal(storage, 300)).toEqual([
      { path: '__untitled__:1/a.feature', content: 'one', updatedAt: 300 },
      { path: '__untitled__:2/b.feature', content: 'live two', updatedAt: 300 },
    ])
  })

  it('lets recovery override the debounced session snapshot', () => {
    const merged = mergeUntitledSessionWithRecovery(
      [{ path: '__untitled__:1/a.feature', content: 'old session' }],
      [{ path: '__untitled__:1/a.feature', content: 'new journal', updatedAt: 400 }],
    )

    expect(merged).toEqual([{ path: '__untitled__:1/a.feature', content: 'new journal' }])
  })

  it('does not resurrect explicitly cleared journal paths', () => {
    const storage = new MemoryStorage()
    journalUntitledTabs(
      [{ path: '__untitled__:1/a.feature', content: 'old', draft: 'new', dirty: true }],
      '__untitled__:1/a.feature',
      () => 'new',
      storage,
      500,
    )

    expect(clearUntitledRecoveryPath('__untitled__:1/a.feature', storage, 500)).toBe(true)
    expect(loadUntitledRecoveryJournal(storage, 500)).toEqual([])
    expect(storage.getItem(UNTITLED_RECOVERY_KEY)).toBeNull()
  })

  it('ignores corrupted and stale journals', () => {
    const storage = new MemoryStorage()
    storage.setItem(UNTITLED_RECOVERY_KEY, '{broken')
    expect(loadUntitledRecoveryJournal(storage, 600)).toEqual([])
    expect(storage.getItem(UNTITLED_RECOVERY_KEY)).toBeNull()

    storage.setItem(
      UNTITLED_RECOVERY_KEY,
      JSON.stringify({
        version: 1,
        tabs: [
          {
            path: '__untitled__:1/a.feature',
            content: 'stale',
            updatedAt: 600 - UNTITLED_RECOVERY_MAX_AGE_MS - 1,
          },
        ],
      }),
    )
    expect(loadUntitledRecoveryJournal(storage, 600)).toEqual([])
  })

  it('returns false when storage is unavailable or quota rejects writes', () => {
    const storage = new FailingStorage({ failSet: true })
    const tabs: TabBody[] = [
      { path: '__untitled__:1/a.feature', content: 'old', draft: 'latest', dirty: true },
    ]

    expect(journalUntitledTabs(tabs, tabs[0].path, () => 'latest', storage, 700)).toBe(false)
    expect(journalUntitledTabs(tabs, tabs[0].path, () => 'latest', null, 700)).toBe(false)
  })

  it('returns false when storage rejects cleanup', () => {
    const storage = new FailingStorage({ failRemove: true })

    expect(clearUntitledRecoveryAll(storage)).toBe(false)
  })

  it('returns false when serialization fails', () => {
    const badContent = {
      toJSON() {
        throw new Error('cannot serialize')
      },
    }
    expect(
      persistUntitledRecoveryEntries(
        [{ path: '__untitled__:1/a.feature', content: badContent as never, updatedAt: 800 }],
        new MemoryStorage(),
      ),
    ).toBe(false)
  })
})

class FailingStorage extends MemoryStorage {
  failSet = false
  failRemove = false

  constructor(flags: Partial<Pick<FailingStorage, 'failSet' | 'failRemove'>> = {}) {
    super()
    Object.assign(this, flags)
  }

  setItem(key: string, value: string) {
    if (this.failSet) throw new Error('storage quota exceeded')
    super.setItem(key, value)
  }

  removeItem(key: string) {
    if (this.failRemove) throw new Error('storage remove failed')
    super.removeItem(key)
  }
}
