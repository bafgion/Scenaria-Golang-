import { describe, expect, it } from 'vitest'
import { collectBlockFoldingRanges } from './gherkinFolding'

describe('collectBlockFoldingRanges', () => {
  it('folds если … конец если block', () => {
    const text = `Функционал: X
  Сценарий: Y
    Если видим "a"
      Когда нажимаю "b"
    Конец если
`
    const ranges = collectBlockFoldingRanges(text)
    expect(ranges).toEqual([{ start: 3, end: 5 }])
  })

  it('folds повторяю … конец block', () => {
    const text = `  Сценарий: Y
    Повторяю 2 раз
      Когда нажимаю "x"
    Конец
`
    const ranges = collectBlockFoldingRanges(text)
    expect(ranges.length).toBe(1)
    expect(ranges[0].start).toBe(2)
    expect(ranges[0].end).toBe(4)
  })
})
