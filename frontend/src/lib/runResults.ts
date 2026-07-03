import { gui } from '../../wailsjs/go/models'
import { listScenarioTitles } from './scenarioAtLine'

const RUN_AT_TOLERANCE_MS = 2000

function parseRunAt(at: string): number {
  if (!at) return 0
  const ms = Date.parse(at)
  return Number.isNaN(ms) ? 0 : ms
}

/** True when entry was recorded during or after the run batch started (with clock/format tolerance). */
export function runAtOrAfter(entryAt: string, sinceAt: string): boolean {
  const entryMs = parseRunAt(entryAt)
  const sinceMs = parseRunAt(sinceAt)
  if (!sinceMs) return true
  if (!entryMs) return false
  return entryMs >= sinceMs - RUN_AT_TOLERANCE_MS
}

export function filterRunResultsSince(
  results: gui.RunResultEntry[],
  runSince?: string,
): gui.RunResultEntry[] {
  if (!runSince) return results
  return results.filter((e) => runAtOrAfter(e.at, runSince))
}

/** Pick the first failure from a run batch, or null if the latest run had no failures. */
export function pickLastRunError(
  results: gui.RunResultEntry[],
  runSince?: string,
): gui.RunResultEntry | null {
  const pool = runSince ? filterRunResultsSince(results, runSince) : results
  return pool.find((e) => !e.success) ?? null
}

function splitResultPath(path: string): { feature: string; scenario: string } {
  const idx = path.indexOf('::')
  if (idx < 0) return { feature: path, scenario: '' }
  return { feature: path.slice(0, idx), scenario: path.slice(idx + 2) }
}

function normalizeDiskPath(path: string): string {
  return path.replace(/\\/g, '/').toLowerCase()
}

/** Map temp on-disk paths back to the tab/path the user was running. */
export function remapRunResultPaths(
  entries: gui.RunResultEntry[],
  diskPaths: string[],
  logicalPaths: string[],
): gui.RunResultEntry[] {
  const map = new Map<string, string>()
  for (let i = 0; i < diskPaths.length; i++) {
    const logical = logicalPaths[i]
    const disk = diskPaths[i]
    if (logical && disk) {
      map.set(normalizeDiskPath(disk), logical)
    }
  }
  if (map.size === 0) return entries
  return entries.map((entry) => {
    const parts = splitResultPath(entry.path)
    const logical = map.get(normalizeDiskPath(parts.feature))
    if (!logical) return entry
    const path = parts.scenario ? `${logical}::${parts.scenario}` : logical
    return gui.RunResultEntry.createFrom({
      path,
      success: entry.success,
      message: entry.message,
      runner: entry.runner,
      at: entry.at,
      failed_step: entry.failed_step,
    })
  })
}

export function buildSyntheticRunError(opts: {
  featurePath: string
  scenario?: string
  message: string
  runner?: string
  at?: string
}): gui.RunResultEntry {
  const path = opts.scenario ? `${opts.featurePath}::${opts.scenario}` : opts.featurePath
  return gui.RunResultEntry.createFrom({
    path,
    success: false,
    message: opts.message,
    runner: opts.runner || 'playwright',
    at: opts.at || new Date().toISOString(),
  })
}

/** Drop a scenario filter that does not exist in the target feature text. */
export function resolveStaleRunScenario(
  scenario: string,
  featureText: string,
  cursorScenario = '',
): string {
  const trimmed = scenario.trim()
  if (!trimmed) return ''
  const names = listScenarioTitles(featureText)
  if (names.length === 0 || names.includes(trimmed)) return trimmed
  return cursorScenario.trim()
}
