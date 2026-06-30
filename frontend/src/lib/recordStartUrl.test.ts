import { describe, expect, it } from 'vitest'
import { firstGotoURLFromText, resolveRecordStartURL } from './recordStartUrl'

describe('firstGotoURLFromText', () => {
  it('reads goto from scenario steps', () => {
    const text = `Функционал: X
  Сценарий: Y
    Когда открыт "https://shop.test/login"
    И нажимаю "#ok"
`
    expect(firstGotoURLFromText(text)).toBe('https://shop.test/login')
  })
})

describe('resolveRecordStartURL', () => {
  it('prefers editor goto over toolbar URL', () => {
    expect(
      resolveRecordStartURL({
        editorText: '  Когда открываю "https://from.feature"',
        startURL: 'https://toolbar',
      }),
    ).toBe('https://from.feature')
  })

  it('falls back to blank when nothing configured', () => {
    expect(resolveRecordStartURL({})).toBe('')
  })
})
