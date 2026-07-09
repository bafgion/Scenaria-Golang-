import { writable } from 'svelte/store'

export type SplashState = {
  appReady: boolean
  message: string
  progress: number
  fading: boolean
}

export const defaultSplashState: SplashState = {
  appReady: false,
  message: '',
  progress: 0,
  fading: false,
}

export function createSplashStore(initial: SplashState = defaultSplashState) {
  const store = writable<SplashState>(initial)
  return {
    subscribe: store.subscribe,
    setStage(message: string, progress: number) {
      store.update((s) => ({ ...s, message, progress }))
    },
    startFading() {
      store.update((s) => ({ ...s, fading: true }))
    },
    markReady() {
      store.update((s) => ({ ...s, appReady: true, fading: false }))
    },
    snapshot(): SplashState {
      let state = defaultSplashState
      const unsub = store.subscribe((s) => {
        state = s
      })
      unsub()
      return state
    },
    reset() {
      store.set(defaultSplashState)
    },
  }
}
