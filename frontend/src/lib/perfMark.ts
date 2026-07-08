const PERF_LOG_THRESHOLD_MS = 50

export function perfMark(name: string, started: number): void {
  const elapsed = performance.now() - started
  if (elapsed >= PERF_LOG_THRESHOLD_MS) {
    console.debug(`[perf] ${name}: ${elapsed.toFixed(1)}ms`)
  }
}

export function perfNow(): number {
  return performance.now()
}
