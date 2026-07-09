import { describe, expect, it } from 'vitest'
import { gui } from '../../wailsjs/go/models'
import { createDiagnosticsStore } from './diagnosticsStore'

describe('diagnosticsStore', () => {
  it('stores validation issues per tab path', () => {
    const store = createDiagnosticsStore()
    store.setIssuesForTab('a.feature', [new gui.ValidationIssue({ message: 'x', line: 1 })])
    expect(store.issuesForTab('a.feature')).toHaveLength(1)
    expect(store.issuesForTab('b.feature')).toEqual([])
  })

  it('clears per-tab cache', () => {
    const store = createDiagnosticsStore()
    store.setIssuesForTab('a.feature', [new gui.ValidationIssue({ message: 'x', line: 1 })])
    store.clearIssuesForTab('a.feature')
    expect(store.issuesForTab('a.feature')).toEqual([])
  })

  it('tracks dismissed hint keys and validate generation', () => {
    const store = createDiagnosticsStore()
    store.dismissHint('hint-1')
    expect(store.isHintDismissed('hint-1')).toBe(true)
    store.clearDismissedHints()
    expect(store.isHintDismissed('hint-1')).toBe(false)
    expect(store.bumpValidateGeneration()).toBe(1)
    expect(store.validateGeneration()).toBe(1)
  })
})
