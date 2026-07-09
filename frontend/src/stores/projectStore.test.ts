import { describe, expect, it } from 'vitest'
import { gui } from '../../wailsjs/go/models'
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
      scenarios: [],
      artifacts: new gui.ProjectArtifacts(),
    })
    expect(currentValue(store)).toEqual({
      path: '/tmp/project',
      version: 3,
      features: ['a.feature'],
      tags: ['@smoke'],
      featureTags: { 'a.feature': ['@smoke'] },
      scenarios: [],
      artifacts: new gui.ProjectArtifacts(),
    })
  })

  it('resets to defaults', () => {
    const store = createProjectStore({
      path: '/tmp/project',
      version: 7,
      features: ['x.feature'],
      tags: ['@x'],
      featureTags: { 'x.feature': ['@x'] },
      scenarios: [],
      artifacts: new gui.ProjectArtifacts(),
    })
    store.reset()
    expect(currentValue(store)).toEqual(defaultProjectState)
  })

  it('patches partial project fields', () => {
    const store = createProjectStore()
    store.patch({ path: '/p', version: 2 })
    expect(currentValue(store).path).toBe('/p')
    expect(currentValue(store).version).toBe(2)
  })
})
