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

  it('redacts secrets before storing journal and status text', () => {
    const store = createJournalStore()
    store.appendLog('opening https://user:pass@example.com/path password=secret')
    store.setStatus('Authorization: Bearer token123', 'error')
    const snap = store.snapshot()
    expect(snap.logText).not.toContain('user:pass@')
    expect(snap.logText).not.toContain('password=secret')
    expect(snap.logText).toContain('https://example.com/path')
    expect(snap.statusMessage).not.toContain('token123')
    expect(snap.statusMessage).toContain('[REDACTED]')
  })
})
