import type { editor as MonacoEditor } from 'monaco-editor'
import { canonicalFeaturePath } from './featurePath'

export type MonacoApi = typeof import('monaco-editor')

/** Стабильный URI модели Monaco для пути .feature (одна модель на вкладку). */
export function featureTabUri(monaco: MonacoApi, path: string) {
  const normalized = canonicalFeaturePath(path)
  return monaco.Uri.parse(`inmemory://scenaria/feature/${encodeURIComponent(normalized)}`)
}

const LANGUAGE_ID = 'scenaria-feature'

/**
 * Реестр текстовых моделей вкладок. Модель живёт, пока вкладка открыта;
 * при закрытии — dispose(), иначе Monaco копит undo/decorations в памяти.
 */
export class MonacoTabModelStore {
  private tracked = new Set<string>()
  private refs = new Map<string, MonacoEditor.ITextModel>()

  getModel(monaco: MonacoApi, path: string): MonacoEditor.ITextModel | null {
    if (!path) return null
    const key = canonicalFeaturePath(path)
    const cached = this.refs.get(key)
    if (cached && !cached.isDisposed()) return cached
    const model = monaco.editor.getModel(featureTabUri(monaco, path))
    if (!model || model.isDisposed()) return null
    this.refs.set(key, model)
    return model
  }

  getOrCreate(monaco: MonacoApi, path: string, text: string): MonacoEditor.ITextModel {
    const existing = this.getModel(monaco, path)
    const key = canonicalFeaturePath(path)
    if (existing) {
      // Model is source of truth while the tab is open — setValue would wipe undo history.
      this.tracked.add(key)
      return existing
    }
    const model = monaco.editor.createModel(text, LANGUAGE_ID, featureTabUri(monaco, path))
    this.tracked.add(key)
    this.refs.set(key, model)
    return model
  }

  release(monaco: MonacoApi, path: string): void {
    if (!path) return
    const key = canonicalFeaturePath(path)
    const model = this.refs.get(key) ?? monaco.editor.getModel(featureTabUri(monaco, path))
    if (model && !model.isDisposed()) {
      model.dispose()
    }
    this.tracked.delete(key)
    this.refs.delete(key)
  }

  releaseExcept(monaco: MonacoApi, keepPaths: Iterable<string>): void {
    const keep = new Set([...keepPaths].map(canonicalFeaturePath))
    for (const path of [...this.tracked]) {
      if (!keep.has(path)) {
        this.release(monaco, path)
      }
    }
  }

  releaseAll(monaco: MonacoApi): void {
    for (const path of [...this.tracked]) {
      this.release(monaco, path)
    }
  }

  trackedPaths(): string[] {
    return [...this.tracked]
  }
}
