import { describe, expect, it } from 'vitest'
import { gui } from '../../wailsjs/go/models'
import { createSettingsDialogStore } from './settingsDialogStore'

describe('settingsDialogStore', () => {
  it('opens with baseline and html report mode', () => {
    const store = createSettingsDialogStore()
    const baseline = gui.AppSettingsDTO.createFrom({ browser: 'firefox' })
    store.openWithBaseline(baseline, 'light')
    expect(store.snapshot()).toEqual({
      baseline,
      htmlReportOpenMode: 'light',
      projectBaseline: 'light',
    })
  })

  it('restores project baseline on cancel', () => {
    const store = createSettingsDialogStore()
    store.openWithBaseline(gui.AppSettingsDTO.createFrom({}), 'full')
    store.setHtmlReportOpenMode('light')
    store.restoreProjectBaseline()
    expect(store.snapshot().htmlReportOpenMode).toBe('full')
  })
})
