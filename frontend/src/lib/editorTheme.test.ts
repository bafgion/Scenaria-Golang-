import { describe, expect, it, vi } from 'vitest'
import { resolveEditorTheme } from './editorTheme'

describe('resolveEditorTheme', () => {
  it('returns explicit dark and light themes', () => {
    expect(resolveEditorTheme('scenaria-dark')).toBe('scenaria-dark')
    expect(resolveEditorTheme('scenaria-light')).toBe('scenaria-light')
  })

  it('follows prefers-color-scheme when theme is system', () => {
    vi.stubGlobal('matchMedia', (query: string) => ({
      matches: query.includes('light'),
      addEventListener: () => {},
      removeEventListener: () => {},
    }))
    expect(resolveEditorTheme('system')).toBe('scenaria-light')
    vi.unstubAllGlobals()
  })
})
