export function createSessionStore() {
  let sessionPersistTimer: ReturnType<typeof setTimeout> | null = null
  let draftAutosaveTimer: ReturnType<typeof setInterval> | null = null

  return {
    schedulePersist(run: () => void, delayMs = 500) {
      if (sessionPersistTimer) clearTimeout(sessionPersistTimer)
      sessionPersistTimer = setTimeout(() => {
        sessionPersistTimer = null
        run()
      }, delayMs)
    },
    flushPersist(run?: () => void) {
      if (sessionPersistTimer) {
        clearTimeout(sessionPersistTimer)
        sessionPersistTimer = null
      }
      run?.()
    },
    startDraftAutosave(run: () => void, intervalMs = 30_000) {
      this.stopDraftAutosave()
      draftAutosaveTimer = setInterval(run, intervalMs)
    },
    stopDraftAutosave() {
      if (draftAutosaveTimer) {
        clearInterval(draftAutosaveTimer)
        draftAutosaveTimer = null
      }
    },
    teardown() {
      this.flushPersist()
      this.stopDraftAutosave()
    },
  }
}
