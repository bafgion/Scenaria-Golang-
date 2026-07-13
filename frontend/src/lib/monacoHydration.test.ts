import { describe, expect, it } from 'vitest'
import monacoEditorSource from './MonacoEditor.svelte?raw'
import {
  modelUriMatchesPath,
  resolveInitialMonacoActivation,
  shouldAcceptMonacoChange,
  shouldApplyExternalEditorValue,
  shouldHydrateModelText,
  shouldUseWelcomeModelForActivation,
} from './monacoHydration'

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

  it('accepts intentional empty user edits from the active model', () => {
    expect(shouldAcceptMonacoChange({
      activePath: '__untitled__:1/novyy-scenariy.feature',
      eventPath: '__untitled__:1/novyy-scenariy.feature',
      activeModelUri: 'file:///active.feature',
      eventModelUri: 'file:///active.feature',
      source: 'user',
    })).toBe(true)
  })
})

describe('Monaco editor startup lifecycle guards', () => {
  it('mounts a restored untitled tab directly with its own path', () => {
    const activation = resolveInitialMonacoActivation({
      activePath: '__untitled__:1/novyy-scenariy.feature',
      valuePath: null,
      value: 'Feature: Restored',
      valueGeneration: 7,
      pending: null,
    })
    expect(activation).toEqual({
      path: '__untitled__:1/novyy-scenariy.feature',
      text: 'Feature: Restored',
      generation: 7,
      mode: 'hydrate',
    })
    expect(shouldUseWelcomeModelForActivation(activation.path)).toBe(false)
  })

  it('preserves activation requested before Monaco is ready', () => {
    const pending = {
      path: '__untitled__:2/recorded.feature',
      text: 'Feature: Recorded',
      generation: 8,
      mode: 'hydrate' as const,
    }
    expect(resolveInitialMonacoActivation({
      activePath: null,
      valuePath: null,
      value: '',
      valueGeneration: 0,
      pending,
    })).toBe(pending)
  })

  it('hydrates an existing empty model with authoritative restored text', () => {
    expect(shouldHydrateModelText('hydrate', '', 'Feature: Restored')).toBe(true)
    expect(shouldHydrateModelText('activate', '', 'Feature: Restored')).toBe(false)
  })

  it('rejects text access when the active model URI belongs to another path', () => {
    expect(modelUriMatchesPath('inmemory://scenaria/feature/a', 'inmemory://scenaria/feature/b')).toBe(false)
    expect(modelUriMatchesPath('inmemory://scenaria/feature/a', 'inmemory://scenaria/feature/a')).toBe(true)
  })

  it('does not mutate the exported value inside MonacoEditor', () => {
    const illegalAssignments = monacoEditorSource
      .split(/\r?\n/)
      .filter((line: string) => /\bvalue\s*(?:=|\+=)/.test(line))
      .filter((line: string) => !line.includes('export let value ='))
    expect(illegalAssignments).toEqual([])
  })
})
