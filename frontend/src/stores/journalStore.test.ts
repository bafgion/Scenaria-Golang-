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

  it('redacts direct patch updates for journal and status text', () => {
    const store = createJournalStore()
    store.patch({
      logText: 'Cookie: sid=123 https://user:pass@example.com/path\nplain context',
      statusMessage: 'password="open sesame" Authorization: Basic abc123',
      statusTone: 'error',
    })
    const snap = store.snapshot()
    expect(snap.logText).not.toContain('sid=123')
    expect(snap.logText).not.toContain('user:pass@')
    expect(snap.logText).toContain('plain context')
    expect(snap.statusMessage).not.toContain('open sesame')
    expect(snap.statusMessage).not.toContain('abc123')
    expect(snap.statusMessage).toContain('[REDACTED]')
    expect(snap.statusTone).toBe('error')
  })

  it('redacts initial visible journal state', () => {
    const store = createJournalStore({
      logText: 'token=abc123',
      statusMessage: 'https://user:pass@example.com/path',
      statusTone: 'busy',
    })
    const snap = store.snapshot()
    expect(snap.logText).not.toContain('abc123')
    expect(snap.statusMessage).not.toContain('user:pass@')
    expect(snap.statusMessage).toContain('https://example.com/path')
  })
})
