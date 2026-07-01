import { describe, expect, it } from 'vitest'
import { matchHotkey, shouldIgnoreAppHotkey } from './hotkeys'

function keyEvent(init: KeyboardEventInit & { code: string }): KeyboardEvent {
  return new KeyboardEvent('keydown', { bubbles: true, ...init })
}

describe('matchHotkey', () => {
  it('maps save and new shortcuts', () => {
    expect(matchHotkey(keyEvent({ code: 'KeyS', ctrlKey: true }))).toBe('save')
    expect(matchHotkey(keyEvent({ code: 'KeyN', ctrlKey: true }))).toBe('new')
    expect(matchHotkey(keyEvent({ code: 'F1' }))).toBe('steps-help')
    expect(matchHotkey(keyEvent({ code: 'F1', shiftKey: true }))).toBe('hotkeys')
    expect(matchHotkey(keyEvent({ code: 'KeyO', ctrlKey: true, shiftKey: true }))).toBe('goto-symbol')
    expect(matchHotkey(keyEvent({ code: 'KeyF', shiftKey: true, altKey: true }))).toBe('format')
  })
})

describe('shouldIgnoreAppHotkey', () => {
  it('allows F1 globally', () => {
    const input = document.createElement('input')
    document.body.appendChild(input)
    const e = keyEvent({ code: 'F1' })
    Object.defineProperty(e, 'target', { value: input })
    expect(shouldIgnoreAppHotkey(e)).toBe(false)
    input.remove()
  })

  it('ignores save inside plain input', () => {
    const input = document.createElement('input')
    document.body.appendChild(input)
    const e = keyEvent({ code: 'KeyS', ctrlKey: true })
    Object.defineProperty(e, 'target', { value: input })
    expect(shouldIgnoreAppHotkey(e)).toBe(true)
    input.remove()
  })
})
