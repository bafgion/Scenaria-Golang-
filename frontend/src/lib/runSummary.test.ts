import { describe, expect, it } from 'vitest'
import { defaultRunForm } from './runTypes'
import { formatLastRunSummary } from './runSummary'

describe('formatLastRunSummary', () => {
  it('does not show stale scenario or tag names in the status bar summary', () => {
    const summary = formatLastRunSummary(defaultRunForm({
      headed: true,
      html: true,
      scenario: 'Test 01',
      tag: '@stale',
    }))

    expect(summary).not.toContain('Test 01')
    expect(summary).not.toContain('@stale')
    expect(summary).toContain('HTML')
  })
})
