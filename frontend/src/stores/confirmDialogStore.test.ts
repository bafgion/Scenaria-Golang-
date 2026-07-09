import { describe, expect, it } from 'vitest'
import { createConfirmDialogStore } from './confirmDialogStore'

describe('confirmDialogStore', () => {
  it('resolves ask promise on close', async () => {
    const store = createConfirmDialogStore()
    const pending = store.ask({ title: 't', message: 'm' }, { ok: 'OK' })
    expect(store.snapshot().open).toBe(true)
    store.close(true)
    await expect(pending).resolves.toBe(true)
    expect(store.snapshot().open).toBe(false)
  })

  it('remembers skip record tab switch preference', () => {
    const store = createConfirmDialogStore()
    store.close(true, true)
    expect(store.shouldSkipRecordTabSwitchConfirm()).toBe(true)
  })
})
