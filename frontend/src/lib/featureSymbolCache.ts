import { collectBlockFoldingRanges, type FeatureFoldingRange } from './gherkinFolding'
import { parseFeatureSymbols, type FeatureSymbol } from './gherkinDocumentSymbols'

const MAX_CACHE_ENTRIES = 64

function hashText(text: string): string {
  let hash = 0
  for (let i = 0; i < text.length; i++) {
    hash = (hash * 31 + text.charCodeAt(i)) | 0
  }
  return `${text.length}:${hash}`
}

function cacheKey(text: string, versionId?: number | null, modelKey?: string | null): string {
  if (versionId != null && modelKey) {
    return `${modelKey}:v:${versionId}`
  }
  if (versionId != null) {
    return `v:${versionId}:${hashText(text)}`
  }
  return `t:${hashText(text)}`
}

const symbolCache = new Map<string, FeatureSymbol[]>()
const foldingCache = new Map<string, FeatureFoldingRange[]>()

function rememberSymbol(key: string, symbols: FeatureSymbol[]): FeatureSymbol[] {
  if (symbolCache.size >= MAX_CACHE_ENTRIES) {
    const oldest = symbolCache.keys().next().value
    if (oldest) symbolCache.delete(oldest)
  }
  symbolCache.set(key, symbols)
  return symbols
}

/** Parsed feature structure with LRU cache (Monaco model URI + version or text hash). */
export function getCachedFeatureSymbols(
  text: string,
  versionId?: number | null,
  modelKey?: string | null,
): FeatureSymbol[] {
  const key = cacheKey(text, versionId, modelKey)
  const hit = symbolCache.get(key)
  if (hit) {
    return hit
  }
  return rememberSymbol(key, parseFeatureSymbols(text))
}

export function getCachedBlockFoldingRanges(
  text: string,
  versionId?: number | null,
  modelKey?: string | null,
): FeatureFoldingRange[] {
  const key = `fold:${cacheKey(text, versionId, modelKey)}`
  const hit = foldingCache.get(key)
  if (hit) {
    return hit
  }
  const ranges = collectBlockFoldingRanges(text)
  if (foldingCache.size >= MAX_CACHE_ENTRIES) {
    const oldest = foldingCache.keys().next().value
    if (oldest) foldingCache.delete(oldest)
  }
  foldingCache.set(key, ranges)
  return ranges
}

export function clearFeatureSymbolCache(): void {
  symbolCache.clear()
  foldingCache.clear()
}

export function evictFeatureSymbolCache(modelKey: string): void {
  if (!modelKey) return
  const prefix = `${modelKey}:`
  for (const key of [...symbolCache.keys()]) {
    if (key.startsWith(prefix)) symbolCache.delete(key)
  }
  for (const key of [...foldingCache.keys()]) {
    if (key.includes(modelKey)) foldingCache.delete(key)
  }
}
