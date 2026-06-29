import { describe, expect, it } from 'vitest'
import { findHintForMarker, markerCode, rangeTouchesMarker } from './gherkinHintActions'
import type { gui } from '../../wailsjs/go/models'

const hint: gui.ScenarioHintDTO = {
  id: 'menu_hover',
  title: 'Hover перед кликом по меню',
  stepIndex: 0,
  line: 3,
  severity: 'warning',
  autoFixable: true,
}

describe('gherkinHintActions helpers', () => {
  it('markerCode reads string and object codes', () => {
    expect(markerCode({ code: 'menu_hover' } as never)).toBe('menu_hover')
    expect(markerCode({ code: { value: 'x' } } as never)).toBe('x')
  })

  it('rangeTouchesMarker matches any column on the same line', () => {
    const marker = {
      startLineNumber: 3,
      endLineNumber: 3,
      startColumn: 1,
      endColumn: 40,
    }
    expect(
      rangeTouchesMarker(
        { startLineNumber: 3, endLineNumber: 3, startColumn: 10, endColumn: 10 },
        marker as never,
      ),
    ).toBe(true)
    expect(
      rangeTouchesMarker(
        { startLineNumber: 2, endLineNumber: 2, startColumn: 1, endColumn: 1 },
        marker as never,
      ),
    ).toBe(false)
  })

  it('findHintForMarker matches source line and id', () => {
    const marker = {
      source: 'Подсказка',
      code: 'menu_hover',
      startLineNumber: 3,
      endLineNumber: 3,
      startColumn: 1,
      endColumn: 20,
    }
    expect(findHintForMarker([hint], marker as never)).toEqual(hint)
    expect(
      findHintForMarker([hint], { ...marker, source: 'Другое' } as never),
    ).toBeUndefined()
  })
})
