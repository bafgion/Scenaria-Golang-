import { describe, expect, it } from 'vitest'
import { gui } from '../../wailsjs/go/models'
import { DEFAULT_EDITOR_SETTINGS } from './editorOptions'
import { filterScenarioHints } from './scenarioHints'

function hint(severity: string): gui.ScenarioHintDTO {
  return gui.ScenarioHintDTO.createFrom({
    id: 'test',
    title: 'hint',
    stepIndex: 0,
    line: 1,
    severity,
    autoFixable: false,
  })
}

describe('filterScenarioHints', () => {
  it('keeps warnings and info by default', () => {
    const hints = [hint('warning'), hint('info')]
    expect(filterScenarioHints(hints, DEFAULT_EDITOR_SETTINGS)).toHaveLength(2)
  })

  it('filters by severity toggles', () => {
    const hints = [hint('warning'), hint('info')]
    const onlyWarning = filterScenarioHints(hints, {
      ...DEFAULT_EDITOR_SETTINGS,
      scenarioHintsShowInfo: false,
    })
    expect(onlyWarning).toHaveLength(1)
    expect(onlyWarning[0].severity).toBe('warning')
  })
})
