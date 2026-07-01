import { describe, expect, it } from 'vitest'
import { monacoColumnToRuneIndex, runeIndexToMonacoColumn, runeSlice } from './editorColumns'

describe('editorColumns', () => {
  it('maps Monaco column to rune index for Cyrillic', () => {
    const line = '\tИ в'
    expect(monacoColumnToRuneIndex(line, 4)).toBe(4)
    expect(monacoColumnToRuneIndex(line, 3)).toBe(3)
    expect(runeSlice(line, 3, 4)).toBe('в')
  })

  it('round-trips rune index to Monaco column', () => {
    const line = '\tИ ввожу'
    for (let col = 0; col <= 6; col++) {
      const rune = monacoColumnToRuneIndex(line, col)
      expect(runeIndexToMonacoColumn(line, rune)).toBe(col)
    }
  })
})
