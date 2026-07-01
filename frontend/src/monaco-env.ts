let workerReady: Promise<void> | null = null

/** Configures Monaco web workers on first editor load (not at app bootstrap). */
export function ensureMonacoEnvironment(): Promise<void> {
  if (workerReady) {
    return workerReady
  }
  workerReady = import('monaco-editor/esm/vs/editor/editor.worker?worker').then((mod) => {
    const EditorWorker = mod.default
    self.MonacoEnvironment = {
      getWorker() {
        return new EditorWorker()
      },
    }
  })
  return workerReady
}
