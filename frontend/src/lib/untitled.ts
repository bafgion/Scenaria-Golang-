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

export function untitledLabel(path: string): string {
  if (!isUntitled(path)) return path
  const rest = path.slice(UNTITLED_PREFIX.length)
  const slash = rest.indexOf('/')
  return slash >= 0 ? rest.slice(slash + 1) : rest
}
