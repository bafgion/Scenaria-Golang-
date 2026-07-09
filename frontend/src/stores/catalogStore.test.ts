import { describe, expect, it } from 'vitest'
import { createCatalogStore } from './catalogStore'

describe('catalogStore', () => {
  it('toggles batch selection paths', () => {
    const store = createCatalogStore()
    store.toggleBatchPath('a.feature')
    expect(store.snapshot().batchSelected).toEqual(['a.feature'])
    store.toggleBatchPath('a.feature')
    expect(store.snapshot().batchSelected).toEqual([])
  })

  it('tracks collapsed catalog keys', () => {
    const store = createCatalogStore()
    store.setCollapsed('dir-a', true)
    expect(store.collapsedSet().has('dir-a')).toBe(true)
    store.setCollapsed('dir-a', false)
    expect(store.collapsedSet().has('dir-a')).toBe(false)
  })
})
