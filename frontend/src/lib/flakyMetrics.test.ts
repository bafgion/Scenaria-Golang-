import { describe, expect, it } from 'vitest'
import { flakyLabel, flakyScenarioMap } from './flakyMetrics'
import { gui } from '../../wailsjs/go/models'

describe('flakyMetrics', () => {
  it('builds scenario map from metrics DTO', () => {
    const metrics = new gui.FlakyMetricsDTO({
      scenarios: [
        { path: 'a.feature::S', flaky: true, failures: 2, passes: 1, total: 3 },
      ],
      steps: [],
    })
    const map = flakyScenarioMap(metrics)
    expect(map.get('a.feature::S')?.flaky).toBe(true)
    expect(flakyLabel(map.get('a.feature::S'))).toBe('flaky (2✗ / 1✓)')
  })
})
