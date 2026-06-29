import { describe, expect, it } from 'vitest'
import { isUntitled, makeUntitledPath, untitledLabel, UNTITLED_PREFIX } from './untitled'

describe('untitled paths', () => {
  it('creates unique internal paths with display names', () => {
    const a = makeUntitledPath('Новый сценарий')
    const b = makeUntitledPath('Новый сценарий')
    expect(a).not.toBe(b)
    expect(a.startsWith(UNTITLED_PREFIX)).toBe(true)
    expect(a.endsWith('Новый сценарий.feature')).toBe(true)
    expect(b.endsWith('Новый сценарий.feature')).toBe(true)
  })

  it('adds .feature suffix when missing', () => {
    const path = makeUntitledPath('smoke')
    expect(path.endsWith('smoke.feature')).toBe(true)
  })

  it('extracts display label from internal path', () => {
    const path = makeUntitledPath('Тест.feature')
    expect(isUntitled(path)).toBe(true)
    expect(untitledLabel(path)).toBe('Тест.feature')
    expect(untitledLabel('features/smoke.feature')).toBe('features/smoke.feature')
  })
})
