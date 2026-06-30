import {
  collectFeaturePathsUnder,
  normFeaturePath,
  type CatalogNode,
} from './catalogTree'

/** Быстрый lookup выбранных путей (нормализованные ключи). */
export function buildBatchSelectedSet(selected: string[]): Set<string> {
  return new Set(selected.map(normFeaturePath))
}

export function isPathBatchSelected(path: string, selectedSet: Set<string>): boolean {
  return selectedSet.has(normFeaturePath(path))
}

/** Добавить или убрать путь из списка (сравнение без учёта слэшей и регистра). */
export function toggleBatchPath(selected: string[], path: string): string[] {
  const norm = normFeaturePath(path)
  const idx = selected.findIndex((p) => normFeaturePath(p) === norm)
  if (idx >= 0) return selected.filter((_, i) => i !== idx)
  return [...selected, path]
}

/** Все .feature под узлом дерева (с учётом текущего фильтра в explorer). */
export function selectAllFeaturesUnder(tree: CatalogNode | null): string[] {
  return tree ? collectFeaturePathsUnder(tree) : []
}

export type CatalogFileClickAction = 'toggle-batch' | 'open'

/** Что делать при клике по файлу в дереве сценариев. */
export function resolveCatalogFileClickAction(
  batchMode: boolean,
  ctrlKey: boolean,
  metaKey: boolean,
): CatalogFileClickAction {
  if (batchMode) return 'toggle-batch'
  if (ctrlKey || metaKey) return 'toggle-batch'
  return 'open'
}

export type BatchModeToggleResult = {
  batchMode: boolean
  batchSelected: string[]
}

/** Переключить режим «Выбор»: при включении — выбрать все сценарии в текущем дереве. */
export function toggleBatchModeState(
  batchMode: boolean,
  tree: CatalogNode | null,
): BatchModeToggleResult {
  if (batchMode) {
    return { batchMode: false, batchSelected: [] }
  }
  return {
    batchMode: true,
    batchSelected: selectAllFeaturesUnder(tree),
  }
}
