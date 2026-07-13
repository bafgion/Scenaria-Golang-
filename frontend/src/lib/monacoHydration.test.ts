import { describe, expect, it } from 'vitest'
import { shouldAcceptMonacoChange, shouldApplyExternalEditorValue } from './monacoHydration'

describe('shouldApplyExternalEditorValue', () => {
  it('blocks path-less empty startup value from overwriting restored untitled text', () => {
    expect(shouldApplyExternalEditorValue({
      applyingExternal: false,
      suppressMarkerSync: false,
      activePath: '__untitled__:1/novyy-scenariy.feature',
      incomingPath: null,
      activeGeneration: 4,
      incomingGeneration: 4,
      modelText: 'Функционал: Draft',
      incomingText: '',
    })).toBe(false)
  })

  it('blocks stale generation writes for the same path', () => {
    expect(shouldApplyExternalEditorValue({
      applyingExternal: false,
      suppressMarkerSync: false,
      activePath: '__untitled__:1/novyy-scenariy.feature',
      incomingPath: '__untitled__:1/novyy-scenariy.feature',
      activeGeneration: 5,
      incomingGeneration: 4,
      modelText: 'new',
      incomingText: 'old',
    })).toBe(false)
  })

  it('allows current same-path hydration when text differs', () => {
    expect(shouldApplyExternalEditorValue({
      applyingExternal: false,
      suppressMarkerSync: false,
      activePath: '__untitled__:1/novyy-scenariy.feature',
      incomingPath: '__untitled__:1/novyy-scenariy.feature',
      activeGeneration: 5,
      incomingGeneration: 5,
      modelText: 'old',
      incomingText: 'new',
    })).toBe(true)
  })
})

describe('shouldAcceptMonacoChange', () => {
  it('ignores a late empty event from a default model after an untitled tab is active', () => {
    expect(shouldAcceptMonacoChange({
      activePath: '__untitled__:1/novyy-scenariy.feature',
      eventPath: null,
      activeModelUri: 'inmemory://model/__untitled__:1/novyy-scenariy.feature',
      eventModelUri: 'inmemory://model/default',
      source: 'user',
    })).toBe(false)
  })

  it('ignores events from an inactive model with the same path', () => {
    expect(shouldAcceptMonacoChange({
      activePath: '__untitled__:1/novyy-scenariy.feature',
      eventPath: '__untitled__:1/novyy-scenariy.feature',
      activeModelUri: 'file:///active.feature',
      eventModelUri: 'file:///stale.feature',
      source: 'user',
    })).toBe(false)
  })

  it('ignores lifecycle/programmatic sync events', () => {
    expect(shouldAcceptMonacoChange({
      activePath: '__untitled__:1/novyy-scenariy.feature',
      eventPath: '__untitled__:1/novyy-scenariy.feature',
      activeModelUri: 'file:///active.feature',
      eventModelUri: 'file:///active.feature',
      source: 'lifecycle',
    })).toBe(false)
  })

  it('accepts user edits from the active model', () => {
    expect(shouldAcceptMonacoChange({
      activePath: '__untitled__:1/novyy-scenariy.feature',
      eventPath: '__untitled__:1/novyy-scenariy.feature',
      activeModelUri: 'file:///active.feature',
      eventModelUri: 'file:///active.feature',
      source: 'user',
    })).toBe(true)
  })
})
