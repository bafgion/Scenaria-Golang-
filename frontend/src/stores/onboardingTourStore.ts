import { writable } from 'svelte/store'

export type OnboardingTourState = {
  active: boolean
  stepId: string
  validateDone: boolean
  dryRunDone: boolean
  journalVisited: boolean
}

export const defaultOnboardingTourState: OnboardingTourState = {
  active: false,
  stepId: 'welcome',
  validateDone: false,
  dryRunDone: false,
  journalVisited: false,
}

export function createOnboardingTourStore(initial: OnboardingTourState = defaultOnboardingTourState) {
  const store = writable<OnboardingTourState>(initial)
  return {
    subscribe: store.subscribe,
    patch(partial: Partial<OnboardingTourState>) {
      store.update((state) => ({ ...state, ...partial }))
    },
    start() {
      store.update((state) => ({ ...state, active: true }))
    },
    stop() {
      store.update((state) => ({ ...state, active: false }))
    },
    resetProgress() {
      store.update((state) => ({
        ...state,
        validateDone: false,
        dryRunDone: false,
        journalVisited: false,
      }))
    },
    setStepId(stepId: string) {
      store.update((state) => ({ ...state, stepId }))
    },
    reset() {
      store.set(defaultOnboardingTourState)
    },
  }
}
