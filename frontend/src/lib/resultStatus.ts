import type { gui } from '../../wailsjs/go/models'

export type ResultStatusKey =
  | 'queued'
  | 'running'
  | 'passed'
  | 'passedWithRetries'
  | 'failed'
  | 'skipped'
  | 'aborted'
  | 'canceled'
  | 'notStarted'
  | 'dryRun'

export type ResultStatusTone = 'neutral' | 'success' | 'warning' | 'error'

export type ResultStatusInfo = {
  key: ResultStatusKey
  tone: ResultStatusTone
}

function normalizeRawStatus(status: string): string {
  return status.trim().toLowerCase().replace(/[_\s]+/g, '-')
}

export function formatRunResultStatus(
  entry: Pick<gui.RunResultEntry, 'status' | 'success'>,
  opts: { flaky?: boolean } = {},
): ResultStatusInfo {
  const raw = normalizeRawStatus(entry.status || '')
  const flaky = opts.flaky === true

  if (raw) {
    if (raw === 'dry-run' || raw === 'dryrun') {
      return { key: 'dryRun', tone: 'warning' }
    }
    if (raw === 'queued' || raw === 'pending') {
      return { key: 'queued', tone: 'neutral' }
    }
    if (raw === 'running' || raw === 'in-progress' || raw === 'inprogress' || raw === 'started') {
      return { key: 'running', tone: 'neutral' }
    }
    if (raw === 'passed' || raw === 'pass' || raw === 'ok' || raw === 'success') {
      return { key: flaky ? 'passedWithRetries' : 'passed', tone: flaky ? 'warning' : 'success' }
    }
    if (raw === 'skipped' || raw === 'skip' || raw === 'not-run' || raw === 'notrun') {
      return { key: 'skipped', tone: 'neutral' }
    }
    if (raw === 'not-started' || raw === 'notstarted' || raw === 'not-start' || raw === 'notstart') {
      return { key: 'notStarted', tone: 'neutral' }
    }
    if (raw === 'aborted') {
      return { key: 'aborted', tone: 'warning' }
    }
    if (raw === 'canceled' || raw === 'cancelled' || raw === 'cancelled-by-user') {
      return { key: 'canceled', tone: 'warning' }
    }
    if (raw === 'passed-with-retries' || raw === 'passedwithretries' || raw === 'flaky') {
      return { key: 'passedWithRetries', tone: 'warning' }
    }
    if (raw === 'failed' || raw === 'fail' || raw === 'broken' || raw === 'error') {
      return { key: 'failed', tone: 'error' }
    }
  }

  if (entry.success) {
    return { key: flaky ? 'passedWithRetries' : 'passed', tone: flaky ? 'warning' : 'success' }
  }

  return { key: 'failed', tone: 'error' }
}
