import { describe, expect, it } from 'vitest'
import { formatRunProgressLabel, shouldRefreshRunResultsFromProgress } from './wailsEventsController'

describe('wailsEventsController', () => {
  it('refreshes only for scenario_done phase', () => {
    expect(shouldRefreshRunResultsFromProgress({ phase: 'scenario_done' })).toBe(true)
    expect(shouldRefreshRunResultsFromProgress({ phase: 'started' })).toBe(false)
    expect(shouldRefreshRunResultsFromProgress(undefined)).toBe(false)
  })

  it('formats progress label with fallback counters', () => {
    expect(
      formatRunProgressLabel(
        { scenario: 'Scenario A', total: 5, index: 2 },
        0,
        0,
      ),
    ).toBe('Scenario A (2/5)')

    expect(
      formatRunProgressLabel(
        { featurePath: 'demo.feature' },
        3,
        1,
      ),
    ).toBe('demo.feature (1/3)')
  })
})
