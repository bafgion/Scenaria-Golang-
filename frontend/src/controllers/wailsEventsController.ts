export type RunProgressPayload = {
  runId?: string
  phase?: string
  total?: number
  index?: number
  caseId?: string
  scenario?: string
  featurePath?: string
  success?: boolean
  message?: string
}

export function resolveRunProgressCounters(
  payload: RunProgressPayload,
  currentTotal: number,
  currentCompleted: number,
): { total: number; completed: number } {
  const total = payload.total && payload.total > 0 ? payload.total : currentTotal
  let completed = currentCompleted
  if ((payload.phase || '') === 'scenario_done') {
    completed = currentCompleted + 1
  }
  return { total, completed }
}

export type ProjectEventEnvelope<T> = {
  projectVersion?: number
  payload?: T | null
}

export type EventStalenessContext = {
  getProjectVersion: () => number
  getActiveRunId: () => string
  getActiveRecordSessionId: () => string
  getActiveBrowserSessionId: () => string
  isRecording: () => boolean
  syncRecordSessionIds: (recordSessionId?: string, browserSessionId?: string) => void
  syncRunId: (runId: string) => void
}

export function shouldRefreshRunResultsFromProgress(payload: RunProgressPayload | null | undefined): boolean {
  return (payload?.phase || '') === 'scenario_done'
}

export function formatRunProgressLabel(payload: RunProgressPayload, currentTotal: number, currentCompleted: number): string {
  const total = payload.total && payload.total > 0 ? payload.total : currentTotal
  if (total <= 0) return ''
  const name = payload.scenario || payload.featurePath || ''
  const position =
    (payload.phase || '') === 'scenario_done'
      ? currentCompleted
      : payload.index || currentCompleted || 0
  return name ? `${name} (${position}/${total})` : ''
}

export function unwrapProjectEvent<T>(
  raw: T | ProjectEventEnvelope<T>,
  emptyPayload?: T,
): { payload: T; projectVersion: number | null } {
  if (raw && typeof raw === 'object') {
    const envelopeRaw = raw as Record<string, unknown>
    if ('projectVersion' in envelopeRaw || 'payload' in envelopeRaw) {
      const envelope = raw as ProjectEventEnvelope<T>
      const fallback = emptyPayload ?? ({} as T)
      const payload = envelope.payload == null ? fallback : (envelope.payload as T)
      return {
        payload,
        projectVersion: typeof envelope.projectVersion === 'number' ? envelope.projectVersion : null,
      }
    }
  }
  return { payload: raw as T, projectVersion: null }
}

export function createStaleEventLogger() {
  const staleEventDebugKeys = new Set<string>()
  return (kind: string, details: string) => {
    const key = `${kind}:${details}`
    if (staleEventDebugKeys.has(key)) return
    staleEventDebugKeys.add(key)
    console.debug(`[stale-event] ${kind} ignored (${details})`)
  }
}

export function createEventStalenessGuards(ctx: EventStalenessContext, logStale: (kind: string, details: string) => void) {
  function isStaleProjectEvent(projectVersion: number | null): boolean {
    const current = ctx.getProjectVersion()
    const stale = projectVersion !== null && current > 0 && projectVersion !== current
    if (stale) {
      logStale('project', `event=${projectVersion}, current=${current}`)
    }
    return stale
  }

  function isStaleRunEvent(runId: string | null | undefined): boolean {
    if (!runId) return false
    const activeRunId = ctx.getActiveRunId()
    if (!activeRunId) {
      ctx.syncRunId(runId)
      return false
    }
    const stale = activeRunId !== runId
    if (stale) {
      logStale('run', `event=${runId}, active=${activeRunId}`)
    }
    return stale
  }

  function isStaleRecordSessionEvent(recordSessionId: string | null | undefined): boolean {
    if (!recordSessionId) {
      if (ctx.isRecording() || ctx.getActiveRecordSessionId()) {
        logStale('record', 'missing-record-session-id')
        return true
      }
      return false
    }
    const activeRecordSessionId = ctx.getActiveRecordSessionId()
    if (!activeRecordSessionId) {
      ctx.syncRecordSessionIds(recordSessionId, '')
      return false
    }
    const stale = activeRecordSessionId !== recordSessionId
    if (stale) {
      logStale('record', `event=${recordSessionId}, active=${activeRecordSessionId}`)
    }
    return stale
  }

  function isStaleBrowserSessionEvent(browserSessionId: string | null | undefined): boolean {
    if (!browserSessionId) return false
    const activeBrowserSessionId = ctx.getActiveBrowserSessionId()
    if (!activeBrowserSessionId) {
      ctx.syncRecordSessionIds('', browserSessionId)
      return false
    }
    const stale = activeBrowserSessionId !== browserSessionId
    if (stale) {
      logStale('browser', `event=${browserSessionId}, active=${activeBrowserSessionId}`)
    }
    return stale
  }

  function unwrapRecordSessionEvent<T>(raw: T | ProjectEventEnvelope<T>) {
    const unwrapped = unwrapProjectEvent(raw)
    return {
      payload: unwrapped.payload,
      projectVersion: unwrapped.projectVersion,
    }
  }

  return {
    isStaleProjectEvent,
    isStaleRunEvent,
    isStaleRecordSessionEvent,
    isStaleBrowserSessionEvent,
    unwrapRecordSessionEvent,
    syncRecordSessionIds: ctx.syncRecordSessionIds,
  }
}

export type WailsEventStalenessGuards = ReturnType<typeof createEventStalenessGuards>

export type RecordStartedMeta = {
  append?: boolean
  sync?: boolean
  output?: string
  targetPath?: string
  recordSessionId?: string
  browserSessionId?: string
}

export type RecordStoppedPayload = {
  reason?: string
  idleSeconds?: number
  recordSessionId?: string
  browserSessionId?: string
}

export type ReportGotoRequest = {
  feature_path: string
  scenario: string
  leaf_index: number
  line: number
}

export type ReportRerunRequest = {
  feature_path: string
  scenario: string
}

export type ReportTraceRequest = {
  trace_path: string
  report_dir: string
  trace_offset_ms?: number
  step_index?: number
}

export type RecordStepPayload = {
  op?: string
  index?: number
  line?: string
  lines?: string[]
  targetPath?: string
  recordSessionId?: string
  browserSessionId?: string
}

export type WailsEventHandlers = {
  onOtpPrompt: (email: string) => void
  onAppCloseRequested: (reasons: string[]) => void
  onBrowserOpened: () => void
  onBrowserClosed: (result: unknown, browserSessionId?: string) => void
  onBrowserLost: () => void
  onToolbarPicker: () => void
  onRecordStarted: (meta: RecordStartedMeta | string) => void | Promise<void>
  onRecordStopped: (payload?: RecordStoppedPayload) => void
  onRunLogLine: (line: string) => void
  onRunResultsChanged: () => void
  onReportGoto: (req: ReportGotoRequest) => void
  onReportRerun: (req: ReportRerunRequest) => void
  onReportTrace: (req: ReportTraceRequest) => void
  onRecordStep: (payload: RecordStepPayload) => void
  onRecordFinished: (result: unknown) => void
  onRecordError: (message: string) => void
  onVanessaRunStarted: () => void
  onVanessaRunFinished: (result: unknown) => void | Promise<void>
}

export type RunProgressState = {
  getTotal: () => number
  getCurrent: () => number
  getLabel: () => string
  updateProgress: (total: number, current: number, label: string) => void
}

export type BindWailsEventsOptions = {
  guards: WailsEventStalenessGuards
  handlers: WailsEventHandlers
  isRunLogStreaming: () => boolean
  isPlaying: () => boolean
  isRecorderActive: () => boolean
  runProgress: RunProgressState
  scheduleRefreshRunResults: () => void
  eventsOn: (eventName: string, callback: (...args: unknown[]) => void) => () => void
}

function browserSessionIdFromPayload(payload: unknown): string {
  if (typeof payload === 'object' && payload !== null && 'browserSessionId' in payload) {
    return String((payload as { browserSessionId?: string }).browserSessionId || '')
  }
  return ''
}

function recordSessionIdFromPayload(payload: unknown): string {
  if (typeof payload === 'object' && payload !== null && 'recordSessionId' in payload) {
    return String((payload as { recordSessionId?: string }).recordSessionId || '')
  }
  return ''
}

function runResultFromPayload(payload: unknown): unknown {
  if (typeof payload === 'object' && payload !== null && 'result' in payload) {
    return (payload as { result?: unknown }).result
  }
  return payload
}

export function bindWailsEvents(options: BindWailsEventsOptions): () => void {
  const {
    guards,
    handlers,
    isRunLogStreaming,
    isPlaying,
    isRecorderActive,
    runProgress,
    scheduleRefreshRunResults,
    eventsOn,
  } = options
  const unsubs: (() => void)[] = []

  unsubs.push(
    eventsOn('app-close-requested', (raw) => {
      const { payload } = unwrapProjectEvent(raw as { reasons?: string[] } | ProjectEventEnvelope<{ reasons?: string[] }>)
      const reasons = Array.isArray(payload?.reasons) ? payload!.reasons!.map(String) : []
      handlers.onAppCloseRequested(reasons)
    }),
  )

  unsubs.push(
    eventsOn('otp-prompt', (raw) => {
      const { payload, projectVersion } = unwrapProjectEvent(raw as string | ProjectEventEnvelope<string>, '')
      if (guards.isStaleProjectEvent(projectVersion)) return
      handlers.onOtpPrompt(payload || '')
    }),
  )

  unsubs.push(
    eventsOn('browser-opened', (raw) => {
      const { payload, projectVersion } = unwrapProjectEvent(raw as { browserSessionId?: string } | ProjectEventEnvelope<{ browserSessionId?: string }>)
      if (guards.isStaleProjectEvent(projectVersion)) return
      if (guards.isStaleBrowserSessionEvent(payload?.browserSessionId)) return
      handlers.onBrowserOpened()
    }),
  )

  unsubs.push(
    eventsOn('browser-closed', (raw) => {
      const { payload, projectVersion } = unwrapProjectEvent(raw as unknown)
      if (guards.isStaleProjectEvent(projectVersion)) return
      const browserSessionId = browserSessionIdFromPayload(payload)
      if (guards.isStaleBrowserSessionEvent(browserSessionId)) return
      handlers.onBrowserClosed(runResultFromPayload(payload), browserSessionId)
    }),
  )

  unsubs.push(
    eventsOn('browser-lost', (raw) => {
      const { payload, projectVersion } = unwrapProjectEvent(raw as { browserSessionId?: string } | ProjectEventEnvelope<{ browserSessionId?: string }>)
      if (guards.isStaleProjectEvent(projectVersion)) return
      if (guards.isStaleBrowserSessionEvent(payload?.browserSessionId)) return
      if (isRecorderActive()) handlers.onBrowserLost()
    }),
  )

  unsubs.push(
    eventsOn('toolbar-picker', (raw) => {
      const { payload, projectVersion } = unwrapProjectEvent(raw as { browserSessionId?: string } | ProjectEventEnvelope<{ browserSessionId?: string }>)
      if (guards.isStaleProjectEvent(projectVersion)) return
      if (guards.isStaleBrowserSessionEvent(payload?.browserSessionId)) return
      handlers.onToolbarPicker()
    }),
  )

  unsubs.push(
    eventsOn('record-started', (raw) => {
      void (async () => {
        const { payload, projectVersion } = unwrapProjectEvent(raw as RecordStartedMeta | string | ProjectEventEnvelope<RecordStartedMeta | string>)
        if (guards.isStaleProjectEvent(projectVersion)) return
        const meta = typeof payload === 'object' && payload !== null ? payload : { append: false, output: payload as string }
        if (guards.isStaleRecordSessionEvent(typeof meta === 'object' ? meta.recordSessionId : '')) return
        if (guards.isStaleBrowserSessionEvent(typeof meta === 'object' ? meta.browserSessionId : '')) return
        await handlers.onRecordStarted(meta)
      })()
    }),
  )

  unsubs.push(
    eventsOn('record-stopped', (raw) => {
      const { payload, projectVersion } = guards.unwrapRecordSessionEvent(raw as RecordStoppedPayload | ProjectEventEnvelope<RecordStoppedPayload>)
      if (guards.isStaleProjectEvent(projectVersion)) return
      if (payload && guards.isStaleRecordSessionEvent(payload.recordSessionId)) return
      if (payload && guards.isStaleBrowserSessionEvent(payload.browserSessionId)) return
      if (payload?.recordSessionId || payload?.browserSessionId) {
        guards.syncRecordSessionIds(payload.recordSessionId, payload.browserSessionId)
      }
      handlers.onRecordStopped(payload ?? undefined)
    }),
  )

  unsubs.push(
    eventsOn('run-log-line', (raw) => {
      const { payload, projectVersion } = unwrapProjectEvent(raw as { line?: string; runId?: string } | string | ProjectEventEnvelope<{ line?: string; runId?: string } | string>)
      if (guards.isStaleProjectEvent(projectVersion)) return
      if (!isRunLogStreaming()) return
      const runId = typeof payload === 'string' ? '' : payload?.runId
      if (guards.isStaleRunEvent(runId)) return
      const line = typeof payload === 'string' ? payload : payload?.line
      if (line) handlers.onRunLogLine(line)
    }),
  )

  unsubs.push(
    eventsOn('run-progress', (raw) => {
      const { payload, projectVersion } = unwrapProjectEvent(raw as RunProgressPayload | ProjectEventEnvelope<RunProgressPayload>)
      if (guards.isStaleProjectEvent(projectVersion)) return
      if (!isPlaying() || !payload) return
      if (guards.isStaleRunEvent(payload.runId)) return
      const { total, completed } = resolveRunProgressCounters(
        payload,
        runProgress.getTotal(),
        runProgress.getCurrent(),
      )
      const label = formatRunProgressLabel(payload, total, completed)
      if (total > 0 || label) {
        runProgress.updateProgress(Math.max(0, total), Math.max(0, completed), label || runProgress.getLabel())
      }
      if (shouldRefreshRunResultsFromProgress(payload)) {
        scheduleRefreshRunResults()
      }
    }),
  )

  unsubs.push(
    eventsOn('run-results-changed', (raw) => {
      const { payload, projectVersion } = unwrapProjectEvent(raw as { runId?: string } | ProjectEventEnvelope<{ runId?: string }>)
      if (guards.isStaleProjectEvent(projectVersion)) return
      if (guards.isStaleRunEvent(payload?.runId)) return
      handlers.onRunResultsChanged()
    }),
  )

  unsubs.push(
    eventsOn('report-goto', (raw) => {
      const { payload, projectVersion } = unwrapProjectEvent(raw as ReportGotoRequest | ProjectEventEnvelope<ReportGotoRequest>)
      if (guards.isStaleProjectEvent(projectVersion)) return
      if (payload) handlers.onReportGoto(payload)
    }),
  )

  unsubs.push(
    eventsOn('report-rerun', (raw) => {
      const { payload, projectVersion } = unwrapProjectEvent(raw as ReportRerunRequest | ProjectEventEnvelope<ReportRerunRequest>)
      if (guards.isStaleProjectEvent(projectVersion)) return
      if (payload) handlers.onReportRerun(payload)
    }),
  )

  unsubs.push(
    eventsOn('report-trace', (raw) => {
      const { payload, projectVersion } = unwrapProjectEvent(raw as ReportTraceRequest | ProjectEventEnvelope<ReportTraceRequest>)
      if (guards.isStaleProjectEvent(projectVersion)) return
      if (payload) handlers.onReportTrace(payload)
    }),
  )

  unsubs.push(
    eventsOn('record-step', (raw) => {
      const { payload, projectVersion } = guards.unwrapRecordSessionEvent(raw as RecordStepPayload | ProjectEventEnvelope<RecordStepPayload>)
      if (guards.isStaleProjectEvent(projectVersion)) return
      if (guards.isStaleRecordSessionEvent(payload?.recordSessionId)) return
      if (guards.isStaleBrowserSessionEvent(payload?.browserSessionId)) return
      if (payload?.recordSessionId || payload?.browserSessionId) {
        guards.syncRecordSessionIds(payload.recordSessionId, payload.browserSessionId)
      }
      if (payload) handlers.onRecordStep(payload)
    }),
  )

  unsubs.push(
    eventsOn('record-finished', (raw) => {
      const { payload, projectVersion } = guards.unwrapRecordSessionEvent(raw as unknown)
      if (guards.isStaleProjectEvent(projectVersion)) return
      if (guards.isStaleRecordSessionEvent(recordSessionIdFromPayload(payload))) return
      if (guards.isStaleBrowserSessionEvent(browserSessionIdFromPayload(payload))) return
      const recordSessionId = recordSessionIdFromPayload(payload)
      const browserSessionId = browserSessionIdFromPayload(payload)
      if (recordSessionId || browserSessionId) {
        guards.syncRecordSessionIds(recordSessionId, browserSessionId)
      }
      handlers.onRecordFinished(runResultFromPayload(payload))
    }),
  )

  unsubs.push(
    eventsOn('record-error', (raw) => {
      const { payload, projectVersion } = guards.unwrapRecordSessionEvent(raw as { message?: string; recordSessionId?: string; browserSessionId?: string } | string | ProjectEventEnvelope<{ message?: string; recordSessionId?: string; browserSessionId?: string } | string>)
      if (guards.isStaleProjectEvent(projectVersion)) return
      if (guards.isStaleRecordSessionEvent(recordSessionIdFromPayload(payload))) return
      if (guards.isStaleBrowserSessionEvent(browserSessionIdFromPayload(payload))) return
      const recordSessionId = recordSessionIdFromPayload(payload)
      const browserSessionId = browserSessionIdFromPayload(payload)
      if (recordSessionId || browserSessionId) {
        guards.syncRecordSessionIds(recordSessionId, browserSessionId)
      }
      const message =
        typeof payload === 'object' && payload !== null && 'message' in payload
          ? String((payload as { message?: string }).message ?? '')
          : typeof payload === 'string'
            ? payload
            : ''
      handlers.onRecordError(message)
    }),
  )

  unsubs.push(
    eventsOn('vanessa-run-started', (raw) => {
      const { projectVersion } = unwrapProjectEvent(raw as Record<string, unknown>)
      if (guards.isStaleProjectEvent(projectVersion)) return
      handlers.onVanessaRunStarted()
    }),
  )

  unsubs.push(
    eventsOn('vanessa-run-finished', (raw) => {
      void (async () => {
        const { payload, projectVersion } = unwrapProjectEvent(raw as unknown)
        if (guards.isStaleProjectEvent(projectVersion)) return
        await handlers.onVanessaRunFinished(payload)
      })()
    }),
  )

  return () => {
    for (const unsub of unsubs) unsub()
  }
}
