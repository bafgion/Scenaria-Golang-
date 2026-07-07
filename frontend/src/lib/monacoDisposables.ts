type MonacoDisposable = { dispose: () => void }

declare global {
  interface Window {
    __scenariaMonacoDisposables?: Map<string, MonacoDisposable[]>
  }
}

function registry(): Map<string, MonacoDisposable[]> {
  window.__scenariaMonacoDisposables ??= new Map()
  return window.__scenariaMonacoDisposables
}

export function replaceMonacoDisposables(key: string, disposables: MonacoDisposable[]) {
  const current = registry().get(key) ?? []
  for (const disposable of current) {
    try {
      disposable.dispose()
    } catch {
      /* Monaco disposables are best-effort during HMR. */
    }
  }
  registry().set(key, disposables)
}
