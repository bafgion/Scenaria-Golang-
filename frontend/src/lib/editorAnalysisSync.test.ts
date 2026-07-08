import { describe, expect, it } from 'vitest'
import { isEditorAnalysisCurrent, isEditorAnalysisSnapshotVisible } from './editorAnalysisSync'

describe('editorAnalysisSync', () => {
  it('accepts matching generation, tab and text version', () => {
    const snapshot = { generation: 3, tab: 'C:/proj/a.feature', textVersion: 7 }
    expect(isEditorAnalysisCurrent(snapshot, 3, 'C:/proj/a.feature', 7)).toBe(true)
  })

  it('rejects stale generation, tab or text version', () => {
    const snapshot = { generation: 3, tab: 'C:/proj/a.feature', textVersion: 7 }
    expect(isEditorAnalysisCurrent(snapshot, 4, 'C:/proj/a.feature', 7)).toBe(false)
    expect(isEditorAnalysisCurrent(snapshot, 3, 'C:/proj/b.feature', 7)).toBe(false)
    expect(isEditorAnalysisCurrent(snapshot, 3, 'C:/proj/a.feature', 8)).toBe(false)
  })

  it('hides derived analysis when text version moved on', () => {
    expect(isEditorAnalysisSnapshotVisible(5, 5)).toBe(true)
    expect(isEditorAnalysisSnapshotVisible(5, 6)).toBe(false)
  })
})
