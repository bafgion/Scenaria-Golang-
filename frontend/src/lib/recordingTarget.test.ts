import { describe, expect, it } from 'vitest'
import {
  isSameRecordTab,
  normalizeRecordTabPath,
  recordingTabSwitchAllowed,
} from './recordingTarget'

describe('recordingTarget', () => {
  it('normalizes path separators', () => {
    expect(normalizeRecordTabPath('C:\\proj\\a.feature')).toBe('C:/proj/a.feature')
  })

  it('compares tab paths', () => {
    expect(isSameRecordTab('C:/a.feature', 'C:\\a.feature')).toBe(true)
    expect(isSameRecordTab('C:/a.feature', 'C:/b.feature')).toBe(false)
  })

  it('blocks switch while recording unless paused', () => {
    expect(
      recordingTabSwitchAllowed(true, false, 'C:/proj/smoke.feature', 'C:/proj/other.feature'),
    ).toBe(false)
    expect(
      recordingTabSwitchAllowed(true, true, 'C:/proj/smoke.feature', 'C:/proj/other.feature'),
    ).toBe(true)
    expect(
      recordingTabSwitchAllowed(false, false, 'C:/proj/smoke.feature', 'C:/proj/other.feature'),
    ).toBe(true)
    expect(
      recordingTabSwitchAllowed(true, false, 'C:/proj/smoke.feature', 'C:/proj/smoke.feature'),
    ).toBe(true)
  })
})
