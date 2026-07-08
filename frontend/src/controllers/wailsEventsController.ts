export type RunProgressPayload = {
  runId?: string
  phase?: string
  total?: number
  index?: number
  scenario?: string
  featurePath?: string
}

export function shouldRefreshRunResultsFromProgress(payload: RunProgressPayload | null | undefined): boolean {
  return (payload?.phase || '') === 'scenario_done'
}

export function formatRunProgressLabel(payload: RunProgressPayload, currentTotal: number, currentIndex: number): string {
  const total = payload.total ?? currentTotal
  const index = payload.index ?? currentIndex
  if (total <= 0) return ''
  const name = payload.scenario || payload.featurePath || ''
  return name ? `${name} (${index}/${total})` : ''
}
