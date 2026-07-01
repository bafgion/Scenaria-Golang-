import type { editor as MonacoEditor } from 'monaco-editor'

export type MonacoApi = typeof import('monaco-editor')

/** Стабильный URI модели Monaco для пути .feature (одна модель на вкладку). */
export function featureTabUri(monaco: MonacoApi, path: string) {
  const normalized = path.replace(/\\/g, '/')
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
    const cached = this.refs.get(path)
    if (cached && !cached.isDisposed()) return cached
    const model = monaco.editor.getModel(featureTabUri(monaco, path))
    if (!model || model.isDisposed()) return null
    this.refs.set(path, model)
    return model
  }

  getOrCreate(monaco: MonacoApi, path: string, text: string): MonacoEditor.ITextModel {
    const existing = this.getModel(monaco, path)
    if (existing) {
      if (existing.getValue() !== text) {
        existing.setValue(text)
      }
      this.tracked.add(path)
      return existing
    }
    const model = monaco.editor.createModel(text, LANGUAGE_ID, featureTabUri(monaco, path))
    this.tracked.add(path)
    this.refs.set(path, model)
    return model
  }

  release(monaco: MonacoApi, path: string): void {
    if (!path) return
    const model = this.refs.get(path) ?? monaco.editor.getModel(featureTabUri(monaco, path))
    if (model && !model.isDisposed()) {
      model.dispose()
    }
    this.tracked.delete(path)
    this.refs.delete(path)
  }

  releaseExcept(monaco: MonacoApi, keepPaths: Iterable<string>): void {
    const keep = new Set(keepPaths)
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
