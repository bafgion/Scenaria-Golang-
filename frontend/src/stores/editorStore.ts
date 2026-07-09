import { writable } from 'svelte/store'
import type { gui } from '../../wailsjs/go/models'

export type EditorState = {
  text: string
  textVersion: number
  cursorLine: number
  steps: gui.EditorStepRow[]
  stepsTextVersion: number
  stepsPanelTab: 'outline' | 'steps'
}

export const defaultEditorState: EditorState = {
  text: '',
  textVersion: 0,
  cursorLine: 1,
  steps: [],
  stepsTextVersion: -1,
  stepsPanelTab: 'outline',
}
export function createEditorStore(initial: EditorState = defaultEditorState) {
  const store = writable<EditorState>(initial)
  return {
    subscribe: store.subscribe,
    patch(partial: Partial<EditorState>) {
      store.update((s) => ({ ...s, ...partial }))
    },
    setText(text: string) {
      store.update((s) => ({ ...s, text }))
    },
    bumpVersion() {
      store.update((s) => ({ ...s, textVersion: s.textVersion + 1 }))
    },
    setTextWithBump(text: string) {
      store.update((s) => ({ ...s, text, textVersion: s.textVersion + 1 }))
    },
    setCursorLine(cursorLine: number) {
      store.update((s) => ({ ...s, cursorLine }))
    },
    setSteps(steps: gui.EditorStepRow[], stepsTextVersion: number) {
      store.update((s) => ({ ...s, steps, stepsTextVersion }))
    },
    clearSteps(stepsTextVersion: number) {
      store.update((s) => ({ ...s, steps: [], stepsTextVersion }))
    },
    setStepsPanelTab(stepsPanelTab: 'outline' | 'steps') {
      store.update((s) => ({ ...s, stepsPanelTab }))
    },
    reset() {
      store.set(defaultEditorState)
    },
    snapshot(): EditorState {
      let state = defaultEditorState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
  }
}
