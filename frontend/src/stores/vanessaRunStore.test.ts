import { describe, expect, it } from 'vitest'
import { gui } from '../../wailsjs/go/models'
import { createVanessaRunStore } from './vanessaRunStore'

describe('vanessaRunStore', () => {
  it('prepares dialog defaults', () => {
    const store = createVanessaRunStore()
    store.prepareDialog({ dry: true, scenario: 'Login', dialogScenarios: ['Login'] })
    const snap = store.snapshot()
    expect(snap.dry).toBe(true)
    expect(snap.scenario).toBe('Login')
    expect(snap.tag).toBe('')
  })

  it('builds plugin request from state', () => {
    const store = createVanessaRunStore()
    store.patch({ dry: false, tag: '@smoke', excludeTags: '@wip, @skip', scenario: 'A' })
    const req = store.buildPluginRequest(store.snapshot())
    expect(req.name).toBe('vanessa')
    expect(req.tag).toBe('@smoke')
    expect(req.excludeTags).toEqual(['@wip', '@skip'])
    expect(req.scenario).toBe('A')
  })
})
