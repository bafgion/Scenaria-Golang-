import { describe, expect, it } from 'vitest'
import {
  isSameRecordTab,
  isRecordingTargetReadOnly,
  normalizeRecordTabPath,
  recordingTabSwitchAllowed,
  resolveRecordStartedTargetPath,
  resolveRecordingTargetPath,
  shouldApplyLiveRecordedStep,
} from './recordingTarget'
import { makeUntitledPath } from './untitled'

describe('recordingTarget', () => {
  it('normalizes path separators', () => {
    expect(normalizeRecordTabPath('C:\\proj\\a.feature')).toBe('C:/proj/a.feature')
  })

  it('compares tab paths', () => {
    expect(isSameRecordTab('C:/a.feature', 'C:\\a.feature')).toBe(true)
    expect(isSameRecordTab('C:/a.feature', 'C:/b.feature')).toBe(false)
  })

  it('ignores late recorder events after capture stops', () => {
    expect(shouldApplyLiveRecordedStep(false, 'нажимаю "#late"')).toBe(false)
    expect(shouldApplyLiveRecordedStep(true, 'нажимаю "#live"')).toBe(true)
    expect(shouldApplyLiveRecordedStep(true, '   ')).toBe(false)
  })

  it('blocks switching away from active recording target until UI confirms', () => {
    expect(
      recordingTabSwitchAllowed(true, false, 'C:/proj/smoke.feature', 'C:/proj/other.feature'),
    ).toBe(false)
    expect(
      recordingTabSwitchAllowed(true, true, 'C:/proj/smoke.feature', 'C:/proj/other.feature'),
    ).toBe(true)
    expect(
      recordingTabSwitchAllowed(true, false, 'C:/proj/smoke.feature', 'C:\\proj\\smoke.feature'),
    ).toBe(true)
  })

  it('resolves explicit recording target path', () => {
    expect(resolveRecordingTargetPath('C:/proj/a.feature', 'C:/proj/b.feature')).toBe(
      'C:/proj/a.feature',
    )
    expect(resolveRecordingTargetPath('', 'C:\\proj\\b.feature')).toBe('C:/proj/b.feature')
  })

  it('keeps an untitled recording tab as the UI target over backend file paths', () => {
    const untitled = makeUntitledPath('recording.feature')

    expect(resolveRecordStartedTargetPath('C:/proj/recorded.feature', untitled)).toBe(untitled)
    expect(resolveRecordingTargetPath('C:/proj/recorded.feature', untitled)).toBe(untitled)
  })

  it('locks the active recording target for manual edits', () => {
    expect(isRecordingTargetReadOnly(true, 'C:/proj/smoke.feature', 'C:/proj/smoke.feature')).toBe(
      true,
    )
    expect(
      isRecordingTargetReadOnly(true, 'C:/proj/smoke.feature', 'C:/proj/other.feature'),
    ).toBe(false)
    expect(isRecordingTargetReadOnly(false, 'C:/proj/smoke.feature', 'C:/proj/smoke.feature')).toBe(
      false,
    )
  })
})
