import { describe, expect, it } from 'vitest'
import { gui } from '../../wailsjs/go/models'
import { createRecorderPrefsStore, defaultRecorderPrefsState } from './recorderPrefsStore'

describe('recorderPrefsStore', () => {
  it('applies recording preference fields from settings DTO', () => {
    const store = createRecorderPrefsStore()
    store.applyFromDTO(
      gui.AppSettingsDTO.createFrom({
        filterRecording: true,
        navOnlyRecording: true,
        hoverRecord: false,
      }),
    )

    let snapshot = defaultRecorderPrefsState
    const unsub = store.subscribe((s) => {
      snapshot = s
    })
    unsub()

    expect(snapshot.filterRecording).toBe(true)
    expect(snapshot.navOnlyRecording).toBe(true)
    expect(snapshot.hoverRecord).toBe(false)
  })

  it('exports DTO fields for persistence', () => {
    const store = createRecorderPrefsStore()
    store.patch({ filterRecording: true })
    let snapshot = defaultRecorderPrefsState
    const unsub = store.subscribe((s) => {
      snapshot = s
    })
    unsub()
    expect(store.dtoFields(snapshot)).toEqual({
      filterRecording: true,
      navOnlyRecording: false,
      hoverRecord: false,
    })
  })
})
