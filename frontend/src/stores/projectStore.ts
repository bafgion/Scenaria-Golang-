import { writable } from 'svelte/store'

export type ProjectState = {
  path: string
  version: number
  features: string[]
  tags: string[]
  featureTags: Record<string, string[]>
}

export const defaultProjectState: ProjectState = {
  path: '',
  version: 0,
  features: [],
  tags: [],
  featureTags: {},
}

export function createProjectStore(initial: ProjectState = defaultProjectState) {
  const store = writable<ProjectState>(initial)
  return {
    subscribe: store.subscribe,
    setProject(next: ProjectState) {
      store.set(next)
    },
    reset() {
      store.set(defaultProjectState)
    },
  }
}
