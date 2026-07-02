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
    expect(matchHotkey(keyEvent({ code: 'KeyR', ctrlKey: true, shiftKey: true }))).toBe('record-stop')
    expect(matchHotkey(keyEvent({ code: 'KeyP', altKey: true }))).toBe('record-pause')
  })
})

describe('shouldIgnoreAppHotkey', () => {
  it('ignores F1 inside plain input', () => {
    const input = document.createElement('input')
    document.body.appendChild(input)
    const e = keyEvent({ code: 'F1' })
    Object.defineProperty(e, 'target', { value: input })
    expect(shouldIgnoreAppHotkey(e)).toBe(true)
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

  it('allows app hotkeys inside monaco editor', () => {
    const wrap = document.createElement('div')
    wrap.className = 'monaco-wrap'
    const inner = document.createElement('div')
    wrap.appendChild(inner)
    document.body.appendChild(wrap)
    const save = keyEvent({ code: 'KeyS', ctrlKey: true })
    Object.defineProperty(save, 'target', { value: inner })
    expect(shouldIgnoreAppHotkey(save)).toBe(false)
    const palette = keyEvent({ code: 'KeyP', ctrlKey: true, shiftKey: true })
    Object.defineProperty(palette, 'target', { value: inner })
    expect(shouldIgnoreAppHotkey(palette)).toBe(false)
    wrap.remove()
  })
})
