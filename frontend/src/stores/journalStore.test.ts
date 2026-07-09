import { describe, expect, it } from 'vitest'
import { createJournalStore } from './journalStore'

describe('journalStore', () => {
  it('appends log lines with trailing newline', () => {
    const store = createJournalStore()
    store.appendLog('a')
    store.appendLog('b\n')
    expect(store.snapshot().logText).toBe('a\nb\n')
  })

  it('updates status message and tone', () => {
    const store = createJournalStore()
    store.setStatus('busy', 'busy')
    expect(store.snapshot()).toMatchObject({ statusMessage: 'busy', statusTone: 'busy' })
  })
})
