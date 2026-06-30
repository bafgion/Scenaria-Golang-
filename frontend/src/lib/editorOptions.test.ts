import { describe, expect, it } from 'vitest'
import {
  DEFAULT_EDITOR_SETTINGS,
  editorSettingsFromDTO,
  editorSettingsToDTO,
} from './editorOptions'
import { settings } from '../../wailsjs/go/models'

describe('editorOptions', () => {
  it('returns defaults when dto is empty', () => {
    expect(editorSettingsFromDTO(null)).toEqual(DEFAULT_EDITOR_SETTINGS)
  })

  it('merges valid dto fields', () => {
    const dto = settings.EditorSettings.createFrom({
      fontSize: 16,
      wordWrap: 'off',
      minimap: false,
      theme: 'scenaria-light',
    })
    const merged = editorSettingsFromDTO(dto)
    expect(merged.fontSize).toBe(16)
    expect(merged.wordWrap).toBe('off')
    expect(merged.minimap).toBe(false)
    expect(merged.theme).toBe('scenaria-light')
    expect(merged.tabSize).toBe(DEFAULT_EDITOR_SETTINGS.tabSize)
  })

  it('round-trips theme through dto', () => {
    const dto = editorSettingsToDTO({
      ...DEFAULT_EDITOR_SETTINGS,
      theme: 'scenaria-light',
    })
    expect(editorSettingsFromDTO(dto).theme).toBe('scenaria-light')
  })

  it('round-trips through dto', () => {
    const dto = editorSettingsToDTO({
      ...DEFAULT_EDITOR_SETTINGS,
      fontSize: 15,
      folding: true,
    })
    const merged = editorSettingsFromDTO(dto)
    expect(merged.fontSize).toBe(15)
    expect(merged.folding).toBe(true)
  })
})
