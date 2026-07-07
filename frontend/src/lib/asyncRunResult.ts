import { EventsOn } from '../../wailsjs/runtime/runtime'
import { gui } from '../../wailsjs/go/models'

type AsyncRunResultPayload = {
  jobId?: string
  result?: gui.RunResult
}

export async function startRunResultJob(
  finishedEvent: string,
  start: () => Promise<string>,
): Promise<gui.RunResult> {
  let jobId = ''
  const pending: AsyncRunResultPayload[] = []
  return new Promise((resolve, reject) => {
    const finish = (payload: AsyncRunResultPayload) => {
      unsubscribe?.()
      resolve(payload.result ?? gui.RunResult.createFrom({ output: '', error: 'missing async result' }))
    }
    const unsubscribe = EventsOn(finishedEvent, (payload: AsyncRunResultPayload) => {
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
