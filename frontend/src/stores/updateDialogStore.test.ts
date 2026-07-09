import { describe, expect, it } from 'vitest'
import { gui } from '../../wailsjs/go/models'
import { createUpdateDialogStore } from './updateDialogStore'

describe('updateDialogStore', () => {
  it('applies update check result', () => {
    const store = createUpdateDialogStore()
    store.applyCheckResult(
      new gui.UpdateInfoDTO({
        message: 'v2 available',
        updateAvailable: true,
      }),
    )
    expect(store.snapshot()).toMatchObject({
      message: 'v2 available',
      hasUpdate: true,
    })
  })
})
