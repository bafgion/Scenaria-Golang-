export type MonacoHydrationState = {
  applyingExternal: boolean
  suppressMarkerSync: boolean
  activePath: string | null
  incomingPath: string | null
  activeGeneration: number
  incomingGeneration: number
  modelText: string
  incomingText: string
}

export function shouldApplyExternalEditorValue(state: MonacoHydrationState): boolean {
  if (state.applyingExternal || state.suppressMarkerSync) return false
  if (state.activePath !== state.incomingPath) return false
  if (state.incomingGeneration < state.activeGeneration) return false
  return state.modelText !== state.incomingText
}

export type MonacoChangeRoutingState = {
  activePath: string | null
  eventPath: string | null
  activeModelUri: string | null
  eventModelUri: string | null
  source?: string
}

export function shouldAcceptMonacoChange(state: MonacoChangeRoutingState): boolean {
  if (state.source && state.source !== 'user') return false
  if (state.eventPath !== state.activePath) return false
  if (state.activeModelUri && state.eventModelUri && state.eventModelUri !== state.activeModelUri) {
    return false
  }
  return true
}
