import type { gui } from '../../wailsjs/go/models'

export type FlakyScenarioStat = {
  flaky: boolean
  failures: number
  passes: number
}

/** Build a lookup map of scenario flaky stats from API metrics. */
export function flakyScenarioMap(metrics: gui.FlakyMetricsDTO | null | undefined): Map<string, FlakyScenarioStat> {
  const map = new Map<string, FlakyScenarioStat>()
  for (const item of metrics?.scenarios ?? []) {
    map.set(item.path, {
      flaky: item.flaky,
      failures: item.failures,
      passes: item.passes,
    })
  }
  return map
}

export function isFlakyPath(map: Map<string, FlakyScenarioStat>, path: string): boolean {
  return map.get(path)?.flaky === true
}

export function flakyLabel(stat: FlakyScenarioStat | undefined): string {
  if (!stat?.flaky) return ''
  return `flaky (${stat.failures}✗ / ${stat.passes}✓)`
}

/** Step-level repeats: same step index failed multiple times. */
export function flakyStepHints(metrics: gui.FlakyMetricsDTO | null | undefined): Map<string, string> {
  const map = new Map<string, string>()
  for (const item of metrics?.steps ?? []) {
    map.set(item.path, `шаг ${item.step + 1} — ${item.failures} падений`)
  }
  return map
}
