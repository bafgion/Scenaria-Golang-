import { describe, expect, it } from 'vitest'
import { replaceModelText } from './editorTextSync'

function mockEditor(initial: string) {
  let value = initial
  const undoStops: string[] = []
  const edits: string[] = []
  return {
    value: () => value,
    pushUndoStop: () => undoStops.push('stop'),
    executeEdits: (source: string, ops: Array<{ text: string }>) => {
      edits.push(source)
      if (ops[0]) value = ops[0].text
    },
    getModel: () => ({
      getValue: () => value,
      getFullModelRange: () => ({ startLineNumber: 1, startColumn: 1, endLineNumber: 1, endColumn: 1 }),
    }),
    undoStops,
    edits,
  }
}

describe('replaceModelText', () => {
  it('skips when text is unchanged', () => {
    const editor = mockEditor('same')
    expect(replaceModelText(editor as never, 'same')).toBe(false)
    expect(editor.edits).toHaveLength(0)
  })

  it('uses executeEdits instead of setValue', () => {
    const editor = mockEditor('old')
    expect(replaceModelText(editor as never, 'new', 'test')).toBe(true)
    expect(editor.value()).toBe('new')
    expect(editor.edits).toEqual(['test'])
    expect(editor.undoStops).toHaveLength(2)
  })
})
