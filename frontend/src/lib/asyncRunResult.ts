import { EventsOn } from '../../wailsjs/runtime/runtime'
import { gui } from '../../wailsjs/go/models'

type AsyncRunResultPayload = {
  jobId?: string
  result?: gui.RunResult
}

type ProjectEventEnvelope<T> = {
  projectVersion?: number
  payload?: T | null
}

function unwrapAsyncPayload(raw: AsyncRunResultPayload | ProjectEventEnvelope<AsyncRunResultPayload>): {
  payload: AsyncRunResultPayload
  projectVersion: number | null
} {
  if (raw && typeof raw === 'object' && ('projectVersion' in raw || 'payload' in raw)) {
    const envelope = raw as ProjectEventEnvelope<AsyncRunResultPayload>
    return {
      payload: envelope.payload ?? {},
      projectVersion: typeof envelope.projectVersion === 'number' ? envelope.projectVersion : null,
    }
  }
  return { payload: raw as AsyncRunResultPayload, projectVersion: null }
}

export async function startRunResultJob(
  finishedEvent: string,
  start: () => Promise<string>,
  currentProjectVersion: number | (() => number) = 0,
): Promise<gui.RunResult> {
  let jobId = ''
  const pending: AsyncRunResultPayload[] = []
  const getCurrentProjectVersion =
    typeof currentProjectVersion === 'function' ? currentProjectVersion : () => currentProjectVersion
  return new Promise((resolve, reject) => {
    const finish = (payload: AsyncRunResultPayload) => {
      unsubscribe?.()
      resolve(payload.result ?? gui.RunResult.createFrom({ output: '', error: 'missing async result' }))
    }
    const unsubscribe = EventsOn(finishedEvent, (raw: AsyncRunResultPayload | ProjectEventEnvelope<AsyncRunResultPayload>) => {
      const { payload, projectVersion } = unwrapAsyncPayload(raw)
      const activeProjectVersion = getCurrentProjectVersion()
      if (projectVersion !== null && activeProjectVersion > 0 && projectVersion !== activeProjectVersion) return
      if (!payload?.jobId) return
      if (!jobId) {
        pending.push(payload)
        return
      }
      if (payload.jobId === jobId) finish(payload)
    })
    start()
      .then((id) => {
        jobId = id
        const matched = pending.find((payload) => payload.jobId === jobId)
        if (matched) finish(matched)
      })
      .catch((err) => {
        unsubscribe?.()
        reject(err)
      })
  })
}
