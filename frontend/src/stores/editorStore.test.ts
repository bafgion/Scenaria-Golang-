import { describe, expect, it } from 'vitest'
import { createEditorStore } from './editorStore'

describe('editorStore', () => {
  it('bumps text version on editor change', () => {
    const store = createEditorStore()
    store.setTextWithBump('Feature: x')
    const snap = store.snapshot()
    expect(snap.text).toBe('Feature: x')
    expect(snap.textVersion).toBe(1)
  })

  it('resets to empty state', () => {
    const store = createEditorStore()
    store.setTextWithBump('a')
    store.reset()
    expect(store.snapshot()).toEqual({ text: '', textVersion: 0, cursorLine: 1, steps: [], stepsTextVersion: -1, stepsPanelTab: 'outline' })
  })
})
