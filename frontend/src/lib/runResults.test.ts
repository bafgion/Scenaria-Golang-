import { describe, expect, it } from 'vitest'
import { gui } from '../../wailsjs/go/models'
import {
  buildSyntheticRunError,
  filterRunResultsSince,
  pickLastRunError,
  remapRunResultPaths,
  resolveStaleRunScenario,
  runAtOrAfter,
} from './runResults'

function entry(at: string, success: boolean, path = 'a.feature::S1'): gui.RunResultEntry {
  return gui.RunResultEntry.createFrom({ at, success, path, message: '', runner: 'playwright', status: success ? 'passed' : 'failed' })
}

describe('runAtOrAfter', () => {
  it('accepts Go RFC3339 without fractional seconds when run started with ms', () => {
    expect(runAtOrAfter('2026-06-29T10:00:01Z', '2026-06-29T10:00:00.500Z')).toBe(true)
  })

  it('rejects entries clearly before the batch', () => {
    expect(runAtOrAfter('2026-06-29T09:00:00Z', '2026-06-29T10:00:00.500Z')).toBe(false)
  })
})

describe('filterRunResultsSince', () => {
  it('includes entries from the current batch only', () => {
    const results = [
      entry('2026-06-29T10:00:01Z', true, 'new.feature::S1'),
      entry('2026-06-29T09:00:00Z', false, 'old.feature::S1'),
    ]
    expect(filterRunResultsSince(results, '2026-06-29T10:00:00.500Z').map((e) => e.path)).toEqual([
      'new.feature::S1',
    ])
  })
})

describe('remapRunResultPaths', () => {
  it('maps temp disk paths to logical tab paths', () => {
    const results = [entry('2026-06-29T10:00:01Z', false, 'C:/proj/.scenaria/temp/run-1/scenario.feature::Первая проверка')]
    const remapped = remapRunResultPaths(
      results,
      ['C:\\proj\\.scenaria\\temp\\run-1\\scenario.feature'],
      ['__untitled__:1/novyy-scenariy.feature'],
    )
    expect(remapped[0]?.path).toBe('__untitled__:1/novyy-scenariy.feature::Первая проверка')
  })
})

describe('pickLastRunError', () => {
  it('returns null when the latest batch only has successes', () => {
    const results = [
      entry('2026-06-29T10:00:01Z', true),
      entry('2026-06-29T09:00:00Z', false),
    ]
    expect(pickLastRunError(results, '2026-06-29T10:00:00Z')).toBeNull()
  })

  it('returns a failure from the latest batch', () => {
    const failed = entry('2026-06-29T10:00:02Z', false, 'b.feature::Fail')
    const results = [failed, entry('2026-06-29T09:00:00Z', false)]
    expect(pickLastRunError(results, '2026-06-29T10:00:00Z')).toEqual(failed)
  })

  it('without runSince returns the newest failure in history order', () => {
    const failed = entry('2026-06-29T10:00:02Z', false)
    const results = [entry('2026-06-29T10:00:03Z', true), failed]
    expect(pickLastRunError(results)).toEqual(failed)
  })
})

describe('buildSyntheticRunError', () => {
  it('builds a failed entry for CLI-level errors', () => {
    const e = buildSyntheticRunError({
      featurePath: 'novyy-scenariy.feature',
      scenario: 'Первая проверка',
      message: 'unsupported step',
    })
    expect(e.success).toBe(false)
    expect(e.status).toBe('failed')
    expect(e.path).toBe('novyy-scenariy.feature::Первая проверка')
    expect(e.message).toBe('unsupported step')
  })
})

describe('resolveStaleRunScenario', () => {
  const text = `Функционал: X
Сценарий: Первая проверка
  Допустим открыт сайт`

  it('keeps scenario that exists in the file', () => {
    expect(resolveStaleRunScenario('Первая проверка', text)).toBe('Первая проверка')
  })

  it('clears scenario from another file', () => {
    expect(resolveStaleRunScenario('Проверка с именованным клиентом', text)).toBe('')
  })

  it('falls back to cursor scenario when stale', () => {
    expect(resolveStaleRunScenario('Old', text, 'Первая проверка')).toBe('Первая проверка')
  })
})

