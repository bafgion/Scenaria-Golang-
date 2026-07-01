let untitledCounter = 0

export const UNTITLED_PREFIX = '__untitled__:'

export function isUntitled(path: string): boolean {
  return path.startsWith(UNTITLED_PREFIX)
}

export function makeUntitledPath(displayName: string): string {
  untitledCounter += 1
  const trimmed = displayName.trim() || 'Без названия.feature'
  const name = trimmed.toLowerCase().endsWith('.feature') ? trimmed : `${trimmed}.feature`
  return `${UNTITLED_PREFIX}${untitledCounter}/${name}`
}

/** Keep new untitled ids above restored session paths. */
export function syncUntitledCounterFromPaths(paths: Iterable<string>) {
  let max = 0
  for (const path of paths) {
    if (!isUntitled(path)) continue
    const rest = path.slice(UNTITLED_PREFIX.length)
    const slash = rest.indexOf('/')
    const raw = slash >= 0 ? rest.slice(0, slash) : rest
    const num = Number.parseInt(raw, 10)
    if (!Number.isNaN(num) && num > max) {
      max = num
    }
  }
  if (max > untitledCounter) {
    untitledCounter = max
  }
}

export function untitledLabel(path: string): string {
  if (!isUntitled(path)) return path
  const rest = path.slice(UNTITLED_PREFIX.length)
  const slash = rest.indexOf('/')
  return slash >= 0 ? rest.slice(slash + 1) : rest
}
