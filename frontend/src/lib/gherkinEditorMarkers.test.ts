import { describe, expect, it } from 'vitest'
import {
  buildEditorMarkers,
  EDITOR_MARKER_OWNER,
  VALIDATION_MARKER_SOURCE,
} from './gherkinEditorMarkers'
import { HINT_MARKER_SOURCE } from './gherkinHintActionsHelpers'

describe('buildEditorMarkers', () => {
  it('merges validation errors and hints under one owner contract', () => {
    const model = {
      getLineCount: () => 3,
      getLineMaxColumn: (line: number) => (line === 2 ? 20 : 10),
    }
    const monaco = {
      MarkerSeverity: { Error: 8, Warning: 4, Info: 2 },
    }

    const markers = buildEditorMarkers(
      model as never,
      [{ line: 2, message: 'bad step' }],
      [{ id: 'hint-1', title: 'Добавьте ожидание', line: 3, severity: 'warning', stepIndex: 0, autoFixable: false }],
      monaco as never,
    )

    expect(markers).toHaveLength(2)
    expect(markers[0].source).toBe(VALIDATION_MARKER_SOURCE)
    expect(markers[1].source).toBe(HINT_MARKER_SOURCE)
    expect(EDITOR_MARKER_OWNER).toBe('scenaria-editor')
  })
})
