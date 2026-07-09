import { describe, expect, it, vi } from 'vitest'
import {
  bindWailsEvents,
  createEventStalenessGuards,
  createStaleEventLogger,
  formatRunProgressLabel,
  resolveRunProgressCounters,
  shouldRefreshRunResultsFromProgress,
  unwrapProjectEvent,
} from './wailsEventsController'

describe('wailsEventsController', () => {
  it('refreshes only for scenario_done phase', () => {
    expect(shouldRefreshRunResultsFromProgress({ phase: 'scenario_done' })).toBe(true)
    expect(shouldRefreshRunResultsFromProgress({ phase: 'started' })).toBe(false)
    expect(shouldRefreshRunResultsFromProgress(undefined)).toBe(false)
  })

  it('formats progress label with fallback counters', () => {
    expect(
      formatRunProgressLabel(
        { scenario: 'Scenario A', total: 5, index: 2, phase: 'scenario_start' },
        0,
        0,
      ),
    ).toBe('Scenario A (2/5)')

    expect(
      formatRunProgressLabel(
        { scenario: 'Scenario B', total: 5, phase: 'scenario_done' },
        5,
        3,
      ),
    ).toBe('Scenario B (3/5)')

    expect(
      formatRunProgressLabel(
        { featurePath: 'demo.feature' },
        3,
        1,
      ),
    ).toBe('demo.feature (1/3)')
  })

  it('increments completed count only on scenario_done', () => {
    expect(resolveRunProgressCounters({ phase: 'scenario_start', total: 4 }, 0, 1)).toEqual({
      total: 4,
      completed: 1,
    })
    expect(resolveRunProgressCounters({ phase: 'scenario_done', total: 4 }, 4, 2)).toEqual({
      total: 4,
      completed: 3,
    })
  })

  it('unwraps project event envelopes including null payload', () => {
    expect(unwrapProjectEvent({ projectVersion: 3, payload: null }, 'fallback')).toEqual({
      payload: 'fallback',
      projectVersion: 3,
    })
    expect(unwrapProjectEvent({ projectVersion: 2, payload: { ok: true } })).toEqual({
      payload: { ok: true },
      projectVersion: 2,
    })
    expect(unwrapProjectEvent('plain')).toEqual({
      payload: 'plain',
      projectVersion: null,
    })
  })

  it('detects stale project and run events', () => {
    const logStale = vi.fn()
    const guards = createEventStalenessGuards(
      {
        getProjectVersion: () => 5,
        getActiveRunId: () => 'run-a',
        getActiveRecordSessionId: () => '',
        getActiveBrowserSessionId: () => '',
        isRecording: () => false,
        syncRecordSessionIds: () => {},
        syncRunId: () => {},
      },
      logStale,
    )

    expect(guards.isStaleProjectEvent(4)).toBe(true)
    expect(guards.isStaleProjectEvent(5)).toBe(false)
    expect(guards.isStaleRunEvent('run-b')).toBe(true)
    expect(guards.isStaleRunEvent('run-a')).toBe(false)
    expect(logStale).toHaveBeenCalled()
  })

  it('syncs accepted record-step session ids after stale checks', () => {
    const syncRecordSessionIds = vi.fn()
    const onRecordStep = vi.fn()
    const events = new Map<string, (...args: unknown[]) => void>()
    const guards = createEventStalenessGuards(
      {
        getProjectVersion: () => 0,
        getActiveRunId: () => '',
        getActiveRecordSessionId: () => '',
        getActiveBrowserSessionId: () => '',
        isRecording: () => true,
        syncRecordSessionIds,
        syncRunId: () => {},
      },
      () => {},
    )

    bindWailsEvents({
      guards,
      isRunLogStreaming: () => false,
      isPlaying: () => false,
      isRecorderActive: () => true,
      runProgress: {
        getTotal: () => 0,
        getCurrent: () => 0,
        getLabel: () => '',
        updateProgress: () => {},
      },
      scheduleRefreshRunResults: () => {},
      eventsOn: (eventName, callback) => {
        events.set(eventName, callback)
        return () => events.delete(eventName)
      },
      handlers: {
        onOtpPrompt: () => {},
        onBrowserOpened: () => {},
        onBrowserClosed: () => {},
        onBrowserLost: () => {},
        onToolbarPicker: () => {},
        onRecordStarted: () => {},
        onRecordStopped: () => {},
        onRunLogLine: () => {},
        onRunResultsChanged: () => {},
        onReportGoto: () => {},
        onReportRerun: () => {},
        onReportTrace: () => {},
        onRecordStep,
        onRecordFinished: () => {},
        onRecordError: () => {},
        onVanessaRunStarted: () => {},
        onVanessaRunFinished: () => {},
      },
    })

    events.get('record-step')?.({
      projectVersion: 1,
      payload: { recordSessionId: 'rec-1', browserSessionId: 'br-1' },
    })
    expect(onRecordStep).toHaveBeenCalledWith({ recordSessionId: 'rec-1', browserSessionId: 'br-1' })
    expect(syncRecordSessionIds).toHaveBeenCalledWith('rec-1', 'br-1')
  })

  it('deduplicates stale event debug logs', () => {
    const logStale = createStaleEventLogger()
    const guards = createEventStalenessGuards(
      {
        getProjectVersion: () => 2,
        getActiveRunId: () => '',
        getActiveRecordSessionId: () => '',
        getActiveBrowserSessionId: () => '',
        isRecording: () => false,
        syncRecordSessionIds: () => {},
        syncRunId: () => {},
      },
      logStale,
    )

    expect(guards.isStaleProjectEvent(1)).toBe(true)
    expect(guards.isStaleProjectEvent(1)).toBe(true)
  })

  it('binds wails events and routes otp prompt', () => {
    const onOtpPrompt = vi.fn()
    const events = new Map<string, (...args: unknown[]) => void>()
    const guards = createEventStalenessGuards(
      {
        getProjectVersion: () => 0,
        getActiveRunId: () => '',
        getActiveRecordSessionId: () => '',
        getActiveBrowserSessionId: () => '',
        isRecording: () => false,
        syncRecordSessionIds: () => {},
        syncRunId: () => {},
      },
      () => {},
    )
    const unbind = bindWailsEvents({
      guards,
      isRunLogStreaming: () => true,
      isPlaying: () => true,
      isRecorderActive: () => false,
      runProgress: {
        getTotal: () => 0,
        getCurrent: () => 0,
        getLabel: () => '',
        updateProgress: () => {},
      },
      scheduleRefreshRunResults: () => {},
      eventsOn: (eventName, callback) => {
        events.set(eventName, callback)
        return () => events.delete(eventName)
      },
      handlers: {
        onOtpPrompt,
        onBrowserOpened: () => {},
        onBrowserClosed: () => {},
        onBrowserLost: () => {},
        onToolbarPicker: () => {},
        onRecordStarted: () => {},
        onRecordStopped: () => {},
        onRunLogLine: () => {},
        onRunResultsChanged: () => {},
        onReportGoto: () => {},
        onReportRerun: () => {},
        onReportTrace: () => {},
        onRecordStep: () => {},
        onRecordFinished: () => {},
        onRecordError: () => {},
        onVanessaRunStarted: () => {},
        onVanessaRunFinished: () => {},
      },
    })

    events.get('otp-prompt')?.('user@example.com')
    expect(onOtpPrompt).toHaveBeenCalledWith('user@example.com')
    unbind()
  })

  it('ignores stale record-step events before syncing session ids', () => {
    const syncRecordSessionIds = vi.fn()
    const onRecordStep = vi.fn()
    const events = new Map<string, (...args: unknown[]) => void>()
    const guards = createEventStalenessGuards(
      {
        getProjectVersion: () => 0,
        getActiveRunId: () => '',
        getActiveRecordSessionId: () => 'record-new',
        getActiveBrowserSessionId: () => 'browser-new',
        isRecording: () => true,
        syncRecordSessionIds,
        syncRunId: () => {},
      },
      () => {},
    )
    bindWailsEvents({
      guards,
      isRunLogStreaming: () => false,
      isPlaying: () => false,
      isRecorderActive: () => true,
      runProgress: {
        getTotal: () => 0,
        getCurrent: () => 0,
        getLabel: () => '',
        updateProgress: () => {},
      },
      scheduleRefreshRunResults: () => {},
      eventsOn: (eventName, callback) => {
        events.set(eventName, callback)
        return () => events.delete(eventName)
      },
      handlers: {
        onOtpPrompt: () => {},
        onBrowserOpened: () => {},
        onBrowserClosed: () => {},
        onBrowserLost: () => {},
        onToolbarPicker: () => {},
        onRecordStarted: () => {},
        onRecordStopped: () => {},
        onRunLogLine: () => {},
        onRunResultsChanged: () => {},
        onReportGoto: () => {},
        onReportRerun: () => {},
        onReportTrace: () => {},
        onRecordStep,
        onRecordFinished: () => {},
        onRecordError: () => {},
        onVanessaRunStarted: () => {},
        onVanessaRunFinished: () => {},
      },
    })

    events.get('record-step')?.({
      recordSessionId: 'record-old',
      browserSessionId: 'browser-old',
      line: 'old',
    })

    expect(onRecordStep).not.toHaveBeenCalled()
    expect(syncRecordSessionIds).not.toHaveBeenCalled()
  })

  it('ignores stale record-finished events before syncing session ids', () => {
    const syncRecordSessionIds = vi.fn()
    const onRecordFinished = vi.fn()
    const events = new Map<string, (...args: unknown[]) => void>()
    const guards = createEventStalenessGuards(
      {
        getProjectVersion: () => 0,
        getActiveRunId: () => '',
        getActiveRecordSessionId: () => 'record-new',
        getActiveBrowserSessionId: () => 'browser-new',
        isRecording: () => true,
        syncRecordSessionIds,
        syncRunId: () => {},
      },
      () => {},
    )
    bindWailsEvents({
      guards,
      isRunLogStreaming: () => false,
      isPlaying: () => false,
      isRecorderActive: () => true,
      runProgress: {
        getTotal: () => 0,
        getCurrent: () => 0,
        getLabel: () => '',
        updateProgress: () => {},
      },
      scheduleRefreshRunResults: () => {},
      eventsOn: (eventName, callback) => {
        events.set(eventName, callback)
        return () => events.delete(eventName)
      },
      handlers: {
        onOtpPrompt: () => {},
        onBrowserOpened: () => {},
        onBrowserClosed: () => {},
        onBrowserLost: () => {},
        onToolbarPicker: () => {},
        onRecordStarted: () => {},
        onRecordStopped: () => {},
        onRunLogLine: () => {},
        onRunResultsChanged: () => {},
        onReportGoto: () => {},
        onReportRerun: () => {},
        onReportTrace: () => {},
        onRecordStep: () => {},
        onRecordFinished,
        onRecordError: () => {},
        onVanessaRunStarted: () => {},
        onVanessaRunFinished: () => {},
      },
    })

    events.get('record-finished')?.({
      result: { ok: true },
      recordSessionId: 'record-old',
      browserSessionId: 'browser-old',
    })

    expect(onRecordFinished).not.toHaveBeenCalled()
    expect(syncRecordSessionIds).not.toHaveBeenCalled()
  })

  it('ignores stale record-error events before syncing session ids', () => {
    const syncRecordSessionIds = vi.fn()
    const onRecordError = vi.fn()
    const events = new Map<string, (...args: unknown[]) => void>()
    const guards = createEventStalenessGuards(
      {
        getProjectVersion: () => 0,
        getActiveRunId: () => '',
        getActiveRecordSessionId: () => 'record-new',
        getActiveBrowserSessionId: () => 'browser-new',
        isRecording: () => true,
        syncRecordSessionIds,
        syncRunId: () => {},
      },
      () => {},
    )
    bindWailsEvents({
      guards,
      isRunLogStreaming: () => false,
      isPlaying: () => false,
      isRecorderActive: () => true,
      runProgress: {
        getTotal: () => 0,
        getCurrent: () => 0,
        getLabel: () => '',
        updateProgress: () => {},
      },
      scheduleRefreshRunResults: () => {},
      eventsOn: (eventName, callback) => {
        events.set(eventName, callback)
        return () => events.delete(eventName)
      },
      handlers: {
        onOtpPrompt: () => {},
        onBrowserOpened: () => {},
        onBrowserClosed: () => {},
        onBrowserLost: () => {},
        onToolbarPicker: () => {},
        onRecordStarted: () => {},
        onRecordStopped: () => {},
        onRunLogLine: () => {},
        onRunResultsChanged: () => {},
        onReportGoto: () => {},
        onReportRerun: () => {},
        onReportTrace: () => {},
        onRecordStep: () => {},
        onRecordFinished: () => {},
        onRecordError,
        onVanessaRunStarted: () => {},
        onVanessaRunFinished: () => {},
      },
    })

    events.get('record-error')?.({
      message: 'old error',
      recordSessionId: 'record-old',
      browserSessionId: 'browser-old',
    })

    expect(onRecordError).not.toHaveBeenCalled()
    expect(syncRecordSessionIds).not.toHaveBeenCalled()
  })
})
