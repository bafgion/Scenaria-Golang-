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

export type MonacoActivationMode = 'activate' | 'hydrate'

export type MonacoActivationRequest = {
  path: string | null
  text: string
  generation: number
  mode: MonacoActivationMode
}

export function resolveInitialMonacoActivation(input: {
  activePath: string | null
  valuePath: string | null
  value: string
  valueGeneration: number
  pending: MonacoActivationRequest | null
}): MonacoActivationRequest {
  if (input.pending) return input.pending
  return {
    path: input.activePath ?? input.valuePath,
    text: input.value,
    generation: input.valueGeneration,
    mode: 'hydrate',
  }
}

export function shouldUseWelcomeModelForActivation(path: string | null): boolean {
  return path === null
}

export function shouldHydrateModelText(
  mode: MonacoActivationMode,
  modelText: string,
  authoritativeText: string,
): boolean {
  return mode === 'hydrate' && modelText !== authoritativeText
}

export function modelUriMatchesPath(activeModelUri: string | null, expectedModelUri: string | null): boolean {
  return !!activeModelUri && !!expectedModelUri && activeModelUri === expectedModelUri
}
