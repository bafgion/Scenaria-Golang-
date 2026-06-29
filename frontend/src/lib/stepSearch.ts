/** Wails SearchSteps expects a string; ignore accidental DOM/event values from handlers. */
export function asStepSearchQuery(q: unknown): string {
  return typeof q === 'string' ? q : ''
}
