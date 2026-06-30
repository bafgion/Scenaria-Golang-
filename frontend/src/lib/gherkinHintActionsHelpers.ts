import type { IRange } from 'monaco-editor'
import type { editor } from 'monaco-editor'
import type { gui } from '../../wailsjs/go/models'

export const HINT_MARKER_SOURCE = 'Подсказка'
export const HINT_QUICK_FIX_KIND = 'quickfix'

export function markerCode(marker: editor.IMarkerData): string {
  const code = marker.code
  if (typeof code === 'string') return code
  if (code && typeof code === 'object' && 'value' in code) return String(code.value)
  return ''
}

export function rangeTouchesMarker(range: IRange, marker: editor.IMarkerData): boolean {
  return (
    range.endLineNumber >= marker.startLineNumber &&
    range.startLineNumber <= marker.endLineNumber
  )
}

export function findHintForMarker(
  hints: gui.ScenarioHintDTO[],
  marker: editor.IMarkerData,
): gui.ScenarioHintDTO | undefined {
  if (marker.source !== HINT_MARKER_SOURCE) return undefined
  const code = markerCode(marker)
  return hints.find((h) => h.line === marker.startLineNumber && (!code || h.id === code))
}
