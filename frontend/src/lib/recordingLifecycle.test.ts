import { describe, expect, it } from 'vitest'
import { applyRecordStepEvent } from './recordedStepOps'
import { buildFeatureTemplate } from './featureTemplate'
import {
  applyRecordStepToTabText,
  captureRecordStartedUiTarget,
  countUntitledTabs,
  recordStepSourceText,
  remapTabOnSaveAs,
  resolveApplyRecordStepTarget,
  resolveRecordEditorPrepareAction,
  resolveRecordFeaturePathForUI,
  shouldApplyLiveRecordedStep,
  shouldBufferEarlyRecordStepEvent,
  shouldApplyRecordStepEvent,
  type EditorTab,
} from './recordingLifecycle'
import { makeUntitledPath, UNTITLED_PREFIX } from './untitled'

const welcomeKey = '__welcome__'
const projectPath = 'C:/proj'
const backendOutput = 'C:/proj/recorded.feature'
const template = buildFeatureTemplate({
  title: 'Demo',
  scenario: 'Main',
  startUrl: 'https://example.com',
})

describe('recordingLifecycle', () => {
  it('1. existing Untitled → record two steps → Stop keeps one tab with both steps', () => {
    const untitled = makeUntitledPath('novyy-scenariy.feature')
    const prepare = resolveRecordEditorPrepareAction({
      activeTab: untitled,
      welcomeKey,
      appendPath: '',
      backendOutputPath: backendOutput,
      lastRecordTarget: untitled,
      recordOutput: backendOutput,
      projectPath,
      tabPaths: [untitled],
      uiTargetPath: untitled,
    })
    expect(prepare).toEqual({ kind: 'load', path: untitled })
    expect(countUntitledTabs([{ path: untitled }])).toBe(1)

    let tabs: EditorTab[] = [{ path: untitled, content: template, dirty: true }]
    let editorText = template
    let lineByIndex: Record<number, number> = {}
    const target = resolveApplyRecordStepTarget(backendOutput, untitled)
    expect(target).toBe(untitled)

    const gate = {
      recording: true,
      captureFinalizing: false,
      activeRecordSessionId: 'rec-1',
      eventRecordSessionId: 'rec-1',
      lineByIndexCount: 0,
    }
    for (const line of ['нажимаю "#one"', 'нажимаю "#two"']) {
      const event = { op: 'upsert' as const, index: Object.keys(lineByIndex).length, line }
      expect(shouldApplyRecordStepEvent(gate, event)).toBe(true)
      const applied = applyRecordStepToTabText(
        tabs,
        untitled,
        editorText,
        target,
        event,
        lineByIndex,
        applyRecordStepEvent,
      )
      expect(applied).not.toBeNull()
      tabs = applied!.tabs
      editorText = applied!.editorText
      lineByIndex = applied!.lineByIndex
    }

    expect(editorText).toContain('нажимаю "#one"')
    expect(editorText).toContain('нажимаю "#two"')
    expect(countUntitledTabs(tabs)).toBe(1)
    expect(tabs[0]?.dirty).toBe(true)

    const finalizeGate = {
      recording: false,
      captureFinalizing: true,
      activeRecordSessionId: 'rec-1',
      eventRecordSessionId: 'rec-1',
      lineByIndexCount: Object.keys(lineByIndex).length,
    }
    expect(shouldApplyLiveRecordedStep(finalizeGate, 'нажимаю "#late"')).toBe(true)
  })

  it('2. Stop while the last record-step is queued preserves the final step', () => {
    const untitled = makeUntitledPath('zapis.feature')
    let tabs: EditorTab[] = [{ path: untitled, content: template, dirty: true }]
    let editorText = template
    let lineByIndex: Record<number, number> = {}
    const target = untitled

    const first = applyRecordStepToTabText(
      tabs,
      untitled,
      editorText,
      target,
      { op: 'upsert', index: 0, line: 'нажимаю "#one"' },
      lineByIndex,
      applyRecordStepEvent,
    )!
    tabs = first.tabs
    editorText = first.editorText
    lineByIndex = first.lineByIndex

    const finalizeGate = {
      recording: false,
      captureFinalizing: true,
      activeRecordSessionId: 'rec-1',
      eventRecordSessionId: 'rec-1',
      lineByIndexCount: Object.keys(lineByIndex).length,
    }
    const finalEvent = { op: 'upsert' as const, index: 1, line: 'нажимаю "#two"' }
    expect(shouldApplyRecordStepEvent(finalizeGate, finalEvent)).toBe(true)
    const second = applyRecordStepToTabText(
      tabs,
      untitled,
      editorText,
      target,
      finalEvent,
      lineByIndex,
      applyRecordStepEvent,
    )!
    expect(second.editorText).toContain('нажимаю "#two"')
  })

  it('2b. buffers early browser-toolbar record-step before record-started prepares the target', () => {
    const gate = {
      recording: false,
      captureFinalizing: false,
      activeRecordSessionId: 'rec-1',
      eventRecordSessionId: 'rec-1',
      lineByIndexCount: 0,
      recordingTargetPath: '',
    }
    const event = { op: 'upsert' as const, index: 0, line: 'РЅР°Р¶РёРјР°СЋ "#early"' }

    expect(shouldApplyRecordStepEvent(gate, event)).toBe(false)
    expect(shouldBufferEarlyRecordStepEvent(gate, event)).toBe(true)
  })

  it('2c. applies an early toolbar-stop snapshot after untitled target preparation', () => {
    const untitled = makeUntitledPath('novyy-scenariy.feature')
    const earlyGate = {
      recording: false,
      captureFinalizing: false,
      activeRecordSessionId: 'rec-1',
      eventRecordSessionId: 'rec-1',
      lineByIndexCount: 0,
      recordingTargetPath: '',
    }
    const earlySnapshot = {
      op: 'snapshot' as const,
      lines: ['\tКогда нажимаю "#one"', '\tИ нажимаю "#two"'],
    }

    expect(shouldApplyRecordStepEvent(earlyGate, earlySnapshot)).toBe(false)
    expect(shouldBufferEarlyRecordStepEvent(earlyGate, earlySnapshot)).toBe(true)

    const prepare = resolveRecordEditorPrepareAction({
      activeTab: untitled,
      welcomeKey,
      appendPath: '',
      backendOutputPath: backendOutput,
      lastRecordTarget: untitled,
      recordOutput: backendOutput,
      projectPath,
      tabPaths: [untitled],
      uiTargetPath: untitled,
    })
    expect(prepare).toEqual({ kind: 'load', path: untitled })

    const target = resolveApplyRecordStepTarget(backendOutput, untitled, untitled, welcomeKey)
    const applied = applyRecordStepToTabText(
      [{ path: untitled, content: template, dirty: true }],
      untitled,
      template,
      target,
      earlySnapshot,
      {},
      applyRecordStepEvent,
    )

    expect(applied?.editorText).toContain('#one')
    expect(applied?.editorText).toContain('#two')
    expect(applied?.tabs[0]?.path).toBe(untitled)
  })

  it('3. backend target path does not create an extra tab when UI target is Untitled', () => {
    const untitled = makeUntitledPath('zapis.feature')
    const prepare = resolveRecordEditorPrepareAction({
      activeTab: untitled,
      welcomeKey,
      appendPath: '',
      backendOutputPath: backendOutput,
      lastRecordTarget: untitled,
      recordOutput: backendOutput,
      projectPath,
      tabPaths: [untitled],
      uiTargetPath: untitled,
    })
    expect(prepare.kind).toBe('load')
    if (prepare.kind === 'load') {
      expect(prepare.path).toBe(untitled)
      expect(prepare.path.startsWith(UNTITLED_PREFIX)).toBe(true)
    }
    expect(
      resolveRecordFeaturePathForUI(backendOutput, untitled, backendOutput, '', projectPath),
    ).toBe(untitled)
    expect(resolveApplyRecordStepTarget(backendOutput, untitled)).toBe(untitled)
  })

  it('4. Save As remaps the same Untitled tab', () => {
    const untitled = makeUntitledPath('draft.feature')
    const text = `${template}\n\tКогда нажимаю "#save"`
    const remapped = remapTabOnSaveAs(
      [{ path: untitled, content: text, dirty: true }],
      untitled,
      'C:/proj/saved.feature',
      text,
    )
    expect(remapped).toHaveLength(1)
    expect(remapped[0]?.path).toBe('C:/proj/saved.feature')
    expect(remapped[0]?.dirty).toBe(false)
    expect(countUntitledTabs(remapped)).toBe(0)
  })

  it('5. cancel Save As preserves the same dirty Untitled tab', () => {
    const untitled = makeUntitledPath('draft.feature')
    const text = `${template}\n\tКогда нажимаю "#stay"`
    const tabs = [{ path: untitled, content: text, dirty: true, draft: text }]
    const picked = ''
    const nextTabs = picked
      ? remapTabOnSaveAs(tabs, untitled, picked, text)
      : tabs
    expect(nextTabs).toEqual(tabs)
    expect(countUntitledTabs(nextTabs)).toBe(1)
    expect(nextTabs[0]?.dirty).toBe(true)
  })

  it('6. stale event from an old recorder session is ignored', () => {
    const gate = {
      recording: true,
      captureFinalizing: false,
      activeRecordSessionId: 'rec-new',
      eventRecordSessionId: 'rec-old',
      lineByIndexCount: 0,
    }
    expect(shouldApplyLiveRecordedStep(gate, 'нажимаю "#stale"')).toBe(false)
    expect(
      shouldApplyRecordStepEvent(gate, { op: 'upsert', index: 0, line: 'нажимаю "#stale"' }),
    ).toBe(false)
  })

  it('7. recording into an existing real feature file still works', () => {
    const feature = 'C:/proj/smoke.feature'
    const uiTarget = captureRecordStartedUiTarget(backendOutput, feature, welcomeKey)
    expect(uiTarget).toBe(feature)
    const prepare = resolveRecordEditorPrepareAction({
      activeTab: feature,
      welcomeKey,
      appendPath: feature,
      backendOutputPath: backendOutput,
      lastRecordTarget: feature,
      recordOutput: feature,
      projectPath,
      tabPaths: [feature],
      uiTargetPath: feature,
    })
    expect(prepare).toEqual({ kind: 'load', path: feature })
    expect(
      resolveRecordFeaturePathForUI(backendOutput, feature, feature, feature, projectPath),
    ).toBe(feature)
  })

  it('8. snapshot applies while finalize is in progress even if recording flag already cleared', () => {
    const untitled = makeUntitledPath('zapis.feature')
    const gate = {
      recording: false,
      captureFinalizing: true,
      activeRecordSessionId: 'rec-1',
      eventRecordSessionId: 'rec-1',
      lineByIndexCount: 0,
      recordingTargetPath: untitled,
    }
    expect(
      shouldApplyRecordStepEvent(gate, {
        op: 'snapshot',
        lines: ['\tКогда нажимаю "#one"'],
      }),
    ).toBe(true)
  })

  it('9. snapshot applies when only the recording target path is still set', () => {
    const untitled = makeUntitledPath('zapis.feature')
    const gate = {
      recording: false,
      captureFinalizing: false,
      activeRecordSessionId: 'rec-1',
      eventRecordSessionId: 'rec-1',
      lineByIndexCount: 0,
      recordingTargetPath: untitled,
    }
    expect(
      shouldApplyRecordStepEvent(gate, {
        op: 'snapshot',
        lines: ['\tКогда нажимаю "#one"'],
      }),
    ).toBe(true)
  })

  it('9b. late upsert from the same record session applies after stop kept the target path', () => {
    const untitled = makeUntitledPath('zapis.feature')
    const gate = {
      recording: false,
      captureFinalizing: false,
      activeRecordSessionId: 'rec-1',
      eventRecordSessionId: 'rec-1',
      lineByIndexCount: 0,
      recordingTargetPath: untitled,
    }
    expect(
      shouldApplyRecordStepEvent(gate, {
        op: 'upsert',
        index: 0,
        line: '\tGiven opened "https://example.com"',
      }),
    ).toBe(true)
  })

  it('10. active tab uses live editor text instead of stale tab draft', () => {
    const untitled = makeUntitledPath('zapis.feature')
    const editorText = `${template}\tДопустим нажимаю "#one"\n`
    const staleDraft = template
    const tabs: EditorTab[] = [{ path: untitled, content: template, draft: staleDraft, dirty: true }]
    expect(recordStepSourceText(tabs, untitled, untitled, editorText)).toBe(editorText)

    const applied = applyRecordStepToTabText(
      tabs,
      untitled,
      editorText,
      untitled,
      { op: 'upsert', index: 1, line: 'нажимаю "#two"' },
      { 0: 4 },
      applyRecordStepEvent,
    )
    expect(applied?.editorText).toContain('нажимаю "#two"')
    expect(applied?.tabs[0]?.draft).toContain('нажимаю "#two"')
  })

  it('11. backend disk target applies to active untitled tab', () => {
    const untitled = makeUntitledPath('novyy-scenariy.feature')
    expect(
      resolveApplyRecordStepTarget(backendOutput, backendOutput, untitled, welcomeKey),
    ).toBe(untitled)

    let tabs: EditorTab[] = [{ path: untitled, content: template, dirty: true }]
    let editorText = template
    let lineByIndex: Record<number, number> = {}
    const lines = Array.from({ length: 10 }, (_, i) => `нажимаю "#btn-${i}"`)
    for (const [index, line] of lines.entries()) {
      const applied = applyRecordStepToTabText(
        tabs,
        untitled,
        editorText,
        resolveApplyRecordStepTarget(backendOutput, backendOutput, untitled, welcomeKey),
        { op: 'upsert', index, line },
        lineByIndex,
        applyRecordStepEvent,
      )
      expect(applied).not.toBeNull()
      tabs = applied!.tabs
      editorText = applied!.editorText
      lineByIndex = applied!.lineByIndex
    }
    for (const line of lines) {
      expect(editorText).toContain(line)
    }
  })
})
