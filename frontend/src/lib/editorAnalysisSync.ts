export type EditorAnalysisContext = {
  generation: number
  tab: string | null
  textVersion: number
}

/** True when async editor analysis may still be applied to the active tab. */
export function isEditorAnalysisCurrent(
  snapshot: EditorAnalysisContext,
  currentGeneration: number,
  currentTab: string | null,
  currentTextVersion: number,
): boolean {
  return (
    snapshot.generation === currentGeneration &&
    snapshot.tab === currentTab &&
    snapshot.textVersion === currentTextVersion
  )
}

/** Inlay hints and step metadata must match the editor text version they were built from. */
export function isEditorAnalysisSnapshotVisible(
  analysisTextVersion: number,
  currentTextVersion: number,
): boolean {
  return analysisTextVersion === currentTextVersion
}
