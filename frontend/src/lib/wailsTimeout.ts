export function wailsReady(): boolean {
  return typeof window !== 'undefined' && !!(window as unknown as { go?: unknown }).go
}

/** Race a Wails binding; return null on timeout or error (startup must not hang). */
export async function callWailsWithTimeout<T>(
  label: string,
  promise: Promise<T>,
  timeoutMs = 4000,
): Promise<T | null> {
  if (!wailsReady()) return null
  try {
    return await Promise.race([
      promise,
      new Promise<never>((_, reject) =>
        window.setTimeout(() => reject(new Error(`${label} timeout`)), timeoutMs),
      ),
    ])
  } catch (err) {
    console.warn(`[wails] ${label}:`, err)
    return null
  }
}
