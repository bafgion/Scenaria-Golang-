import { describe, expect, it } from 'vitest'
import { DEFAULT_EDITOR_SETTINGS } from './editorOptions'
import { editorOptionsForLineCount, isLargeFeatureFile, LARGE_FILE_LINE_THRESHOLD, shouldUseHeavyLanguageFeatures } from './editorLargeFile'

describe('editorLargeFile', () => {
  it('detects large files by line threshold', () => {
    expect(isLargeFeatureFile(LARGE_FILE_LINE_THRESHOLD - 1)).toBe(false)
    expect(isLargeFeatureFile(LARGE_FILE_LINE_THRESHOLD)).toBe(true)
    expect(shouldUseHeavyLanguageFeatures(LARGE_FILE_LINE_THRESHOLD - 1)).toBe(true)
    expect(shouldUseHeavyLanguageFeatures(LARGE_FILE_LINE_THRESHOLD)).toBe(false)
  })

  it('disables heavy features for large files', () => {
    const monaco = {
      editor: { ShowLightbulbIconMode: { On: 1 } },
      languages: {},
    }
    const opts = editorOptionsForLineCount(DEFAULT_EDITOR_SETTINGS, monaco as never, 5000)
    expect(opts.minimap).toEqual({ enabled: false })
    expect(opts.codeLens).toBe(false)
  })

  it('keeps user minimap setting for small files', () => {
    const monaco = {
      editor: { ShowLightbulbIconMode: { On: 1 } },
      languages: {},
    }
    const opts = editorOptionsForLineCount(DEFAULT_EDITOR_SETTINGS, monaco as never, 100)
    expect(opts.minimap).toEqual({ enabled: true, scale: 1 })
    expect(opts.codeLens).toBe(true)
    expect(opts.largeFileOptimizations).toBeUndefined()
  })
})
