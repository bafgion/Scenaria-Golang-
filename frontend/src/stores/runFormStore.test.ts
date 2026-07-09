import { describe, expect, it } from 'vitest'
import { defaultRunForm } from '../lib/runTypes'
import { createRunFormStore } from './runFormStore'

describe('runFormStore', () => {
  it('updates last run from settings', () => {
    const store = createRunFormStore()
    store.applyLastRunFromSettings('firefox', 4, 80)
    const { lastRun } = store.snapshot()
    expect(lastRun.browser).toBe('firefox')
    expect(lastRun.workers).toBe(4)
    expect(lastRun.slowMo).toBe(80)
  })

  it('sets run form independently from last run', () => {
    const base = defaultRunForm({ scenario: 'A' })
    const store = createRunFormStore(base)
    store.setRunForm(defaultRunForm({ scenario: 'B' }))
    const snap = store.snapshot()
    expect(snap.lastRun.scenario).toBe('A')
    expect(snap.runForm.scenario).toBe('B')
  })
})
