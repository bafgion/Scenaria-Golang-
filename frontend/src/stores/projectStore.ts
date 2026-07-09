import { writable } from 'svelte/store'
import { gui } from '../../wailsjs/go/models'

export type ProjectState = {
  path: string
  version: number
  features: string[]
  tags: string[]
  featureTags: Record<string, string[]>
  scenarios: string[]
  artifacts: gui.ProjectArtifacts
}

export const defaultProjectState: ProjectState = {
  path: '',
  version: 0,
  features: [],
  tags: [],
  featureTags: {},
  scenarios: [],
  artifacts: new gui.ProjectArtifacts(),
}

export function createProjectStore(initial: ProjectState = defaultProjectState) {
  const store = writable<ProjectState>(initial)
  return {
    subscribe: store.subscribe,
    setProject(next: ProjectState) {
      store.set(next)
    },
    patch(partial: Partial<ProjectState>) {
      store.update((state) => ({ ...state, ...partial }))
    },
    setScenarios(scenarios: string[]) {
      store.update((state) => ({ ...state, scenarios }))
    },
    setArtifacts(artifacts: gui.ProjectArtifacts) {
      store.update((state) => ({ ...state, artifacts }))
    },
    snapshot(): ProjectState {
      let state = defaultProjectState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultProjectState)
    },
  }
}
