import { describe, expect, it } from 'vitest'
import { createRecordFormStore } from './recordFormStore'

describe('recordFormStore', () => {
  it('patches and resets record form fields', () => {
    const store = createRecordFormStore()
    store.patch({ recordURL: 'https://example.com', recordOutput: 'smoke.feature' })
    expect(store.snapshot().recordURL).toBe('https://example.com')
    expect(store.snapshot().recordOutput).toBe('smoke.feature')
    store.reset({ recordIdle: 45 })
    expect(store.snapshot().recordIdle).toBe(45)
    expect(store.snapshot().recordOutput).toBe('recorded.feature')
  })
})
