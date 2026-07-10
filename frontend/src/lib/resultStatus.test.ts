import { describe, expect, it } from 'vitest'
import { gui } from '../../wailsjs/go/models'
import { formatRunResultStatus } from './resultStatus'

function entry(status: string | undefined, success: boolean): gui.RunResultEntry {
  return gui.RunResultEntry.createFrom({
    path: 'demo.feature::Scenario',
    status,
    success,
    message: '',
    runner: 'playwright',
    at: '2026-07-10T10:00:00Z',
  })
}

describe('formatRunResultStatus', () => {
  it('maps terminal statuses to localized keys', () => {
    expect(formatRunResultStatus(entry('queued', false)).key).toBe('queued')
    expect(formatRunResultStatus(entry('running', false)).key).toBe('running')
    expect(formatRunResultStatus(entry('passed', true)).key).toBe('passed')
    expect(formatRunResultStatus(entry('failed', false)).key).toBe('failed')
    expect(formatRunResultStatus(entry('skipped', false)).key).toBe('skipped')
    expect(formatRunResultStatus(entry('canceled', false)).key).toBe('canceled')
    expect(formatRunResultStatus(entry('aborted', false)).key).toBe('aborted')
    expect(formatRunResultStatus(entry('not-started', false)).key).toBe('notStarted')
    expect(formatRunResultStatus(entry('dry-run', true)).key).toBe('dryRun')
  })

  it('treats a flaky passed run as passed-with-retries', () => {
    expect(formatRunResultStatus(entry('passed', true), { flaky: true })).toEqual({
      key: 'passedWithRetries',
      tone: 'warning',
    })
  })

  it('falls back to success and failure when raw status is missing', () => {
    expect(formatRunResultStatus(entry(undefined, true))).toEqual({
      key: 'passed',
      tone: 'success',
    })
    expect(formatRunResultStatus(entry(undefined, false))).toEqual({
      key: 'failed',
      tone: 'error',
    })
  })
})
