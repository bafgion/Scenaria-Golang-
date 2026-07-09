import { describe, expect, it } from 'vitest'
import { createPluginRunStore } from './pluginRunStore'

describe('pluginRunStore', () => {
  it('prepares plugin run dialog', () => {
    const store = createPluginRunStore()
    store.prepareDialog({ name: 'lint', dry: true, scenario: 'S1', dialogScenarios: ['S1'] })
    const snap = store.snapshot()
    expect(snap.name).toBe('lint')
    expect(snap.dry).toBe(true)
    expect(snap.scenario).toBe('S1')
    expect(snap.tag).toBe('')
  })
})
