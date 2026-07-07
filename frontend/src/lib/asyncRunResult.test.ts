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
})
