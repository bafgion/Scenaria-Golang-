import { describe, expect, it } from 'vitest'
import { createProjectStore, defaultProjectState } from './projectStore'

function currentValue<T>(store: { subscribe: (run: (value: T) => void) => () => void }): T {
  let out!: T
  const unsubscribe = store.subscribe((value) => {
    out = value
  })
  unsubscribe()
  return out
}

describe('projectStore', () => {
  it('sets project payload', () => {
    const store = createProjectStore()
    store.setProject({
      path: '/tmp/project',
      version: 3,
      features: ['a.feature'],
      tags: ['@smoke'],
      featureTags: { 'a.feature': ['@smoke'] },
    })
    expect(currentValue(store)).toEqual({
      path: '/tmp/project',
      version: 3,
      features: ['a.feature'],
      tags: ['@smoke'],
      featureTags: { 'a.feature': ['@smoke'] },
    })
  })

  it('resets to defaults', () => {
    const store = createProjectStore({
      path: '/tmp/project',
      version: 7,
      features: ['x.feature'],
      tags: ['@x'],
      featureTags: { 'x.feature': ['@x'] },
    })
    store.reset()
    expect(currentValue(store)).toEqual(defaultProjectState)
  })
})
