import { canonicalFeaturePath } from './featurePath'
import type { RecordStepEvent } from './recordedStepOps'
import { isUntitled, isRealFeaturePath, makeUntitledPath, untitledLabel } from './untitled'
import {
  isSameRecordTab,
  normalizeRecordTabPath,
  resolveRecordStartedTargetPath,
  resolveRecordingTargetPath,
} from './recordingTarget'

export type RecordEditorPrepareAction =
  | { kind: 'load'; path: string }
  | { kind: 'openUntitled'; displayName: string }
  | { kind: 'noop' }

export type RecordEditorPrepareInput = {
  activeTab: string
  welcomeKey: string
  appendPath: string
  backendOutputPath: string
  lastRecordTarget: string
  recordOutput: string
  projectPath: string
  tabPaths: string[]
  /** UI tab captured before async record-started preparation. */
  uiTargetPath?: string
  defaultUntitledName?: string
}

export type RecordStepApplyGate = {
  recording: boolean
  captureFinalizing: boolean
  activeRecordSessionId: string
  eventRecordSessionId?: string
  lineByIndexCount?: number
  recordingTargetPath?: string
}

/** Prefer the in-memory untitled tab over backend output paths used for disk writes. */
export function resolveRecordFeaturePathForUI(
  outputPath: string,
  lastRecordTarget: string,
  recordOutput: string,
  recordAppendTo: string,
  projectPath: string,
): string {
  const append = (recordAppendTo || '').trim().replace(/\\/g, '/')
  if (append) return normalizeRecordTabPath(append)
  if (lastRecordTarget && isUntitled(lastRecordTarget)) {
    return normalizeRecordTabPath(lastRecordTarget)
  }
  if (lastRecordTarget && isRealFeaturePath(lastRecordTarget)) {
    return normalizeRecordTabPath(lastRecordTarget)
  }
  const target = (outputPath || lastRecordTarget || recordOutput || '').trim().replace(/\\/g, '/')
  if (!target) return ''
  if (isUntitled(target)) return normalizeRecordTabPath(target)
  if (!projectPath) return normalizeRecordTabPath(target)
  if (target.startsWith('/') || /^[A-Za-z]:\//.test(target)) return normalizeRecordTabPath(target)
  return normalizeRecordTabPath(`${projectPath.replace(/\\/g, '/')}/${target}`)
}

export function resolveRecordEditorPrepareAction(input: RecordEditorPrepareInput): RecordEditorPrepareAction {
  const appendPath = (input.appendPath || '').trim().replace(/\\/g, '/')
  if (appendPath && !isUntitled(appendPath)) {
    return { kind: 'load', path: normalizeRecordTabPath(appendPath) }
  }

  const uiTarget = (input.uiTargetPath || '').trim()
  if (uiTarget && isUntitled(uiTarget)) {
    const open = input.tabPaths.some((path) => isSameRecordTab(path, uiTarget))
    if (open) {
      return { kind: 'load', path: normalizeRecordTabPath(uiTarget) }
    }
  }

  const activeUntitled =
    input.activeTab &&
    input.activeTab !== input.welcomeKey &&
    isUntitled(input.activeTab)
      ? normalizeRecordTabPath(input.activeTab)
      : ''
  if (activeUntitled && input.tabPaths.some((path) => isSameRecordTab(path, activeUntitled))) {
    return { kind: 'load', path: activeUntitled }
  }

  const featurePath = resolveRecordFeaturePathForUI(
    input.backendOutputPath,
    input.lastRecordTarget,
    input.recordOutput,
    input.appendPath,
    input.projectPath,
  )
  if (featurePath && !isUntitled(featurePath) && input.tabPaths.some((path) => path === featurePath)) {
    return { kind: 'load', path: featurePath }
  }

  if (!input.activeTab || input.activeTab === input.welcomeKey) {
    return {
      kind: 'openUntitled',
      displayName: input.defaultUntitledName || 'zapis.feature',
    }
  }

  return { kind: 'noop' }
}

export function captureRecordStartedUiTarget(
  backendTargetPath: string,
  activeTab: string,
  welcomeKey: string,
): string {
  const uiTarget = resolveRecordStartedTargetPath(backendTargetPath, activeTab, welcomeKey)
  if (uiTarget) return uiTarget
  if (activeTab && activeTab !== welcomeKey) {
    return normalizeRecordTabPath(activeTab)
  }
  return ''
}

export function shouldApplyLiveRecordedStep(gate: RecordStepApplyGate, line: string): boolean {
  if (!line.trim()) return false
  const eventSession = (gate.eventRecordSessionId || '').trim()
  const activeSession = (gate.activeRecordSessionId || '').trim()
  if (eventSession && activeSession && eventSession !== activeSession) {
    return false
  }
  if (gate.recording || gate.captureFinalizing) return true
  return Boolean(eventSession && activeSession && eventSession === activeSession && (gate.recordingTargetPath || '').trim())
}

export function shouldApplyRecordStepEvent(gate: RecordStepApplyGate, event: RecordStepEvent): boolean {
  switch (event.op) {
    case 'upsert':
      return shouldApplyLiveRecordedStep(gate, event.line ?? '')
    case 'reset':
      return true
    case 'delete':
    case 'snapshot':
      return (
        gate.recording ||
        gate.captureFinalizing ||
        (gate.lineByIndexCount ?? 0) > 0 ||
        Boolean((gate.recordingTargetPath || '').trim())
      )
    default:
      return false
  }
}

export function shouldBufferEarlyRecordStepEvent(gate: RecordStepApplyGate, event: RecordStepEvent): boolean {
  if (shouldApplyRecordStepEvent(gate, event)) return false
  if (gate.recording || gate.captureFinalizing || (gate.recordingTargetPath || '').trim()) return false
  const sessionKnown = Boolean((gate.eventRecordSessionId || gate.activeRecordSessionId || '').trim())
  if (!sessionKnown) return false
  switch (event.op) {
    case 'upsert':
      return Boolean((event.line || '').trim())
    case 'snapshot':
      return (event.lines || []).some((line) => Boolean((line || '').trim()))
    case 'delete':
      return event.index !== undefined
    default:
      return false
  }
}

export function resolveApplyRecordStepTarget(
  eventTargetPath: string,
  recordingTargetPath: string,
  activeTab = '',
  welcomeKey = '__welcome__',
): string {
  if (activeTab && activeTab !== welcomeKey && isUntitled(activeTab)) {
    return normalizeRecordTabPath(activeTab)
  }
  const resolved = resolveRecordingTargetPath(eventTargetPath, recordingTargetPath)
  if (activeTab && activeTab !== welcomeKey && resolved && isSameRecordTab(activeTab, resolved)) {
    return normalizeRecordTabPath(activeTab)
  }
  if (resolved) return resolved
  if (activeTab && activeTab !== welcomeKey) {
    return normalizeRecordTabPath(activeTab)
  }
  return ''
}

export type EditorTab = {
  path: string
  content?: string
  draft?: string
  dirty?: boolean
}

export function recordStepSourceText(
  tabs: EditorTab[],
  activeTab: string,
  targetPath: string,
  editorText: string,
): string {
  if (isSameRecordTab(activeTab, targetPath)) {
    return editorText
  }
  const tab = tabs.find((entry) => isSameRecordTab(entry.path, targetPath))
  if (tab) {
    return tab.draft ?? tab.content ?? ''
  }
  return ''
}

export function applyRecordStepToTabText(
  tabs: EditorTab[],
  activeTab: string,
  editorText: string,
  targetPath: string,
  event: RecordStepEvent,
  lineByIndex: Record<number, number>,
  applyEvent: (
    text: string,
    event: RecordStepEvent,
    lineByIndex: Record<number, number>,
  ) => { text: string; lineByIndex: Record<number, number> },
): { tabs: EditorTab[]; editorText: string; lineByIndex: Record<number, number> } | null {
  const tab = tabs.find((entry) => isSameRecordTab(entry.path, targetPath))
  const sourceText = recordStepSourceText(tabs, activeTab, targetPath, editorText)
  if (!tab && !sourceText && event.op !== 'reset') return null
  const result = applyEvent(sourceText, event, lineByIndex)
  const nextTabs = tabs.map((entry) =>
    isSameRecordTab(entry.path, targetPath)
      ? { ...entry, draft: result.text, dirty: true }
      : entry,
  )
  const nextEditorText = isSameRecordTab(activeTab, targetPath) ? result.text : editorText
  return { tabs: nextTabs, editorText: nextEditorText, lineByIndex: result.lineByIndex }
}

export function remapTabOnSaveAs(
  tabs: EditorTab[],
  pathAtStart: string,
  pickedPath: string,
  text: string,
): EditorTab[] {
  const from = normalizeRecordTabPath(pathAtStart)
  const to = canonicalFeaturePath(pickedPath)
  return tabs.map((tab) =>
    isSameRecordTab(tab.path, from)
      ? { path: to, content: text, dirty: false, draft: undefined }
      : tab,
  )
}

export function countUntitledTabs(tabs: EditorTab[]): number {
  return tabs.filter((tab) => isUntitled(tab.path)).length
}

export function untitledTabLabel(path: string): string {
  return untitledLabel(path)
}

export function nextUntitledPath(displayName: string): string {
  return makeUntitledPath(displayName)
}
