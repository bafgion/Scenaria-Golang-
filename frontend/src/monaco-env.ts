import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'

// Vite + Monaco worker wiring for offline Wails bundle.
self.MonacoEnvironment = {
  getWorker() {
    return new editorWorker()
  },
}
