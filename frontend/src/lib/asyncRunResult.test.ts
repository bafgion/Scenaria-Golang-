import { describe, expect, it, vi } from 'vitest'
import { startRunResultJob } from './asyncRunResult'

type Listener = (payload: unknown) => void

const listeners = new Map<string, Listener[]>()

vi.mock('../../wailsjs/runtime/runtime', () => ({
  EventsOn: (eventName: string, callback: Listener) => {
    const current = listeners.get(eventName) ?? []
    current.push(callback)
    listeners.set(eventName, current)
    return () => {
      listeners.set(eventName, (listeners.get(eventName) ?? []).filter((item) => item !== callback))
    }
  },
}))

function emit(eventName: string, payload: unknown) {
  for (const listener of listeners.get(eventName) ?? []) {
    listener(payload)
  }
}

describe('startRunResultJob', () => {
  it('resolves only the matching job result', async () => {
    const resultPromise = startRunResultJob('run-finished', async () => 'job-2')

    emit('run-finished', { jobId: 'job-1', result: { output: 'wrong', error: '' } })
    emit('run-finished', { jobId: 'job-2', result: { output: 'ok', error: '' } })

    await expect(resultPromise).resolves.toMatchObject({ output: 'ok', error: '' })
    expect(listeners.get('run-finished')).toHaveLength(0)
  })

  it('accepts project-scoped event envelope payload', async () => {
    const resultPromise = startRunResultJob('run-finished', async () => 'job-9')

    emit('run-finished', {
      projectVersion: 3,
      payload: { jobId: 'job-9', result: { output: 'wrapped', error: '' } },
    })

    await expect(resultPromise).resolves.toMatchObject({ output: 'wrapped', error: '' })
  })

  it('ignores project-scoped event envelope from another project version', async () => {
    const resultPromise = startRunResultJob('run-finished', async () => 'job-10', 4)

    emit('run-finished', {
      projectVersion: 3,
      payload: { jobId: 'job-10', result: { output: 'stale', error: '' } },
    })
    emit('run-finished', {
      projectVersion: 4,
      payload: { jobId: 'job-10', result: { output: 'current', error: '' } },
    })

    await expect(resultPromise).resolves.toMatchObject({ output: 'current', error: '' })
  })

  it('checks project version at event time', async () => {
    let currentProjectVersion = 4
    const resultPromise = startRunResultJob('run-finished', async () => 'job-12', () => currentProjectVersion)

    currentProjectVersion = 5
    emit('run-finished', {
      projectVersion: 4,
      payload: { jobId: 'job-12', result: { output: 'stale', error: '' } },
    })
    emit('run-finished', {
      projectVersion: 5,
      payload: { jobId: 'job-12', result: { output: 'current', error: '' } },
    })

    await expect(resultPromise).resolves.toMatchObject({ output: 'current', error: '' })
  })

  it('handles project envelope with null payload', async () => {
    const resultPromise = startRunResultJob('run-finished', async () => 'job-11', 7)

    emit('run-finished', {
      projectVersion: 6,
      payload: null,
    })
    emit('run-finished', {
      projectVersion: 7,
      payload: { jobId: 'job-11', result: { output: 'ok', error: '' } },
    })

    await expect(resultPromise).resolves.toMatchObject({ output: 'ok', error: '' })
  })
})
