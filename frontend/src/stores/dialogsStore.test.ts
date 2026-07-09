import { describe, expect, it } from 'vitest'
import { anyAppDialogOpen, createDialogsStore, defaultDialogsState } from './dialogsStore'

describe('dialogsStore', () => {
  it('tracks open/close for a single dialog', () => {
    const store = createDialogsStore()
    store.open('showSettings')
    expect(store.anyOpen()).toBe(true)
    store.close('showSettings')
    expect(store.anyOpen()).toBe(false)
  })

  it('detects any modal dialog open with overlay context', () => {
    expect(anyAppDialogOpen(defaultDialogsState)).toBe(false)
    expect(
      anyAppDialogOpen({ ...defaultDialogsState, showRecord: true }),
    ).toBe(true)
    expect(
      anyAppDialogOpen(defaultDialogsState, { confirmDialogOpen: true }),
    ).toBe(true)
    expect(
      anyAppDialogOpen(
        { ...defaultDialogsState, showPostRecordDiff: true },
        { postRecordPath: 'features/a.feature' },
      ),
    ).toBe(true)
    expect(
      anyAppDialogOpen({ ...defaultDialogsState, showPostRecordDiff: true }),
    ).toBe(false)
  })
})
