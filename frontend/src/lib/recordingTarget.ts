import { canonicalFeaturePath } from './featurePath'

/** Normalize feature tab path for stable comparisons during recording. */
export function normalizeRecordTabPath(path: string): string {
  return canonicalFeaturePath(path)
}

export function isSameRecordTab(a: string, b: string): boolean {
  if (!a || !b) return false
  return normalizeRecordTabPath(a) === normalizeRecordTabPath(b)
}

/** User-initiated tab switch while live recording is active. */
export function recordingTabSwitchAllowed(
  recording: boolean,
  recordPaused: boolean,
  recordingTargetPath: string,
  nextPath: string,
): boolean {
  if (!recording || recordPaused || !recordingTargetPath) return true
  return isSameRecordTab(recordingTargetPath, nextPath)
}

export function resolveRecordingTargetPath(
  eventTargetPath: string,
  recordingTargetPath: string,
): string {
  const path = (eventTargetPath || recordingTargetPath || '').trim()
  return path ? normalizeRecordTabPath(path) : ''
}

/** Ignore record-step events that arrive after capture has stopped. */
export function shouldApplyLiveRecordedStep(recording: boolean, line: string): boolean {
  return recording && line.trim() !== ''
}
