/** Normalize feature tab path for stable comparisons during recording. */
export function normalizeRecordTabPath(path: string): string {
  return path.trim().replace(/\\/g, '/')
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
