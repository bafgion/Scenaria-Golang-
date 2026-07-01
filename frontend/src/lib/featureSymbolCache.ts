import { parseFeatureSymbols, type FeatureSymbol } from './gherkinDocumentSymbols'

const MAX_CACHE_ENTRIES = 64

function hashText(text: string): string {
  let hash = 0
  for (let i = 0; i < text.length; i++) {
    hash = (hash * 31 + text.charCodeAt(i)) | 0
  }
  return `${text.length}:${hash}`
}

function cacheKey(text: string, versionId?: number | null): string {
  if (versionId != null) {
    return `v:${versionId}`
  }
  return `t:${hashText(text)}`
}

const symbolCache = new Map<string, FeatureSymbol[]>()

function remember(key: string, symbols: FeatureSymbol[]): FeatureSymbol[] {
  if (symbolCache.size >= MAX_CACHE_ENTRIES) {
    const oldest = symbolCache.keys().next().value
    if (oldest) symbolCache.delete(oldest)
  }
  symbolCache.set(key, symbols)
  return symbols
}

/** Parsed feature structure with LRU cache (Monaco model version or text hash). */
export function getCachedFeatureSymbols(text: string, versionId?: number | null): FeatureSymbol[] {
  const key = cacheKey(text, versionId)
  const hit = symbolCache.get(key)
  if (hit) {
    return hit
  }
  return remember(key, parseFeatureSymbols(text))
}

export function clearFeatureSymbolCache(): void {
  symbolCache.clear()
}
