import { describe, expect, it } from 'vitest'
import { gui } from '../../wailsjs/go/models'
import { createUiPrefsStore, defaultUiPrefsState } from './uiPrefsStore'

describe('uiPrefsStore', () => {
  it('applies persisted UI preference fields from settings DTO', () => {
    const store = createUiPrefsStore()
    store.applyFromDTO(
      gui.AppSettingsDTO.createFrom({
        toolbarCompact: true,
        stepsPanelVisible: false,
        stepsPanelHeight: 220,
        sidebarWidth: 300,
        onboardingCompleted: true,
        checklistDismissed: true,
        onboardingVersion: 1,
      }),
    )

    let snapshot = defaultUiPrefsState
    const unsub = store.subscribe((s) => {
      snapshot = s
    })
    unsub()

    expect(snapshot.toolbarCompact).toBe(true)
    expect(snapshot.stepsPanelVisible).toBe(false)
    expect(snapshot.stepsPanelHeight).toBe(220)
    expect(snapshot.sidebarWidth).toBe(300)
    expect(snapshot.onboardingCompleted).toBe(true)
    expect(snapshot.checklistDismissed).toBe(true)
  })

  it('exports DTO fields for persistence', () => {
    const store = createUiPrefsStore()
    store.patch({ toolbarCompact: true, stepsPanelHeight: 180 })
    let snapshot = defaultUiPrefsState
    const unsub = store.subscribe((s) => {
      snapshot = s
    })
    unsub()
    const fields = store.dtoFields(snapshot)
    expect(fields.toolbarCompact).toBe(true)
    expect(fields.stepsPanelHeight).toBe(180)
    expect(fields.onboardingVersion).toBe(1)
  })
})
