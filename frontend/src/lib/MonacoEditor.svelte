<script lang="ts">
  import { createEventDispatcher, onDestroy, onMount } from 'svelte'
  import { preloadMonacoEditor } from './appBootstrap'
  import {
    registerHintCodeActions,
    type HintActionHandlers,
  } from './gherkinHintActions'
  import { applyEditorMarkers } from './gherkinEditorMarkers'
  import {
    registerGherkinCodeLens,
    refreshGherkinCodeLens,
    type RunCodeLensHandlers,
  } from './gherkinCodeLens'
  import {
    registerGherkinInlayHints,
    refreshGherkinInlayHints,
    type InlayHintsHandlers,
  } from './gherkinInlayHintsProvider'
  import {
    DEFAULT_EDITOR_SETTINGS,
    toMonacoOptions,
    type EditorSettings,
  } from './editorOptions'
  import { resolveEditorTheme, subscribeSystemTheme } from './editorTheme'
  import { editorOptionsForLineCount } from './editorLargeFile'
  import type { gui } from '../../wailsjs/go/models'
  import type { editor as MonacoEditor } from 'monaco-editor'
  import { replaceModelText } from './editorTextSync'
  import { featureTabUri } from './monacoTabModels'
  import { MonacoTabModelStore } from './monacoTabModels'
  import { MonacoTabViewStateStore } from './monacoTabViewState'

  export let value = ''
  export let valuePath: string | null = null
  export let activePath: string | null = null
  export let valueGeneration = 0
  export let readOnly = false
  export let editorSettings: EditorSettings = { ...DEFAULT_EDITOR_SETTINGS }
  export let hintActions: HintActionHandlers | null = null
  export let runLensActions: RunCodeLensHandlers | null = null
  export let inlayHintsHandlers: InlayHintsHandlers | null = null
  export let scenarioHints: gui.ScenarioHintDTO[] = []

  type MarkerIssue = { line: number; message: string }

  let container: HTMLDivElement
  let editor: MonacoEditor.IStandaloneCodeEditor | null = null
  let monacoApi: typeof import('monaco-editor') | null = null
  let applyingExternal = false
  let suppressMarkerSync = false
  let validationMarkerIssues: MarkerIssue[] = []
  let activeTabPath: string | null = null
  let activeHydrationGeneration = 0
  let welcomeModel: MonacoEditor.ITextModel | null = null
  let largeFileOptionsTimer: ReturnType<typeof setTimeout> | null = null
  let unsubscribeSystemTheme: (() => void) | undefined
  const tabModels = new MonacoTabModelStore()
  const tabViewStates = new MonacoTabViewStateStore()

  type ActivationMode = 'activate' | 'hydrate'
  type PendingActivation = {
    path: string | null
    text: string
    generation: number
    mode: ActivationMode
  }

  type EditorChangeSource = 'user' | 'lifecycle'
  type EditorChangeEvent = {
    path: string | null
    modelUri: string | null
    text: string
    source: EditorChangeSource
    modelVersion: number
    hydrationGeneration: number
  }

  const dispatch = createEventDispatcher<{ change: EditorChangeEvent; cursorline: number; ready: void }>()
  let pendingActivation: PendingActivation | null = null

  function inputPath(): string | null {
    return activePath ?? valuePath
  }

  function modelUriForPath(path: string | null): string | null {
    if (!monacoApi) return null
    if (!path) {
      return welcomeModel && !welcomeModel.isDisposed() ? welcomeModel.uri.toString() : null
    }
    return featureTabUri(monacoApi, path).toString()
  }

  function emitChangeForModel(text: string, path = activeTabPath, source: EditorChangeSource = 'user') {
    const activeModel = editor?.getModel() ?? null
    dispatch('change', {
      path,
      modelUri: modelUriForPath(path),
      text,
      source,
      modelVersion: activeModel?.getVersionId() ?? 0,
      hydrationGeneration: activeHydrationGeneration,
    })
  }

  function ensureWelcomeModel(text: string): MonacoEditor.ITextModel {
    if (!monacoApi) throw new Error('monaco not ready')
    if (welcomeModel && !welcomeModel.isDisposed()) {
      return welcomeModel
    }
    welcomeModel = monacoApi.editor.createModel(text, 'scenaria-feature')
    return welcomeModel
  }

  function finishExternalEdit() {
    applyingExternal = false
    suppressMarkerSync = false
  }

  function attachModel(model: MonacoEditor.ITextModel, opts?: { silent?: boolean }) {
    if (!editor) return
    applyingExternal = true
    suppressMarkerSync = true
    editor.setModel(model)
    if (!opts?.silent) {
      emitChangeForModel(model.getValue())
    }
    queueMicrotask(() => {
      finishExternalEdit()
      syncEditorMarkers()
      syncLargeFileOptions()
    })
  }

  function reconcileAttachedModelText(text: string, source: string) {
    if (!editor) return false
    const model = editor.getModel()
    if (!model || model.getValue() === text) return false
    applyingExternal = true
    suppressMarkerSync = true
    const changed = replaceModelText(editor, text, source)
    queueMicrotask(() => {
      finishExternalEdit()
      syncEditorMarkers()
      syncLargeFileOptions()
    })
    return changed
  }

  function modelForPath(monaco: typeof import('monaco-editor'), path: string | null, text: string) {
    return path ? tabModels.getOrCreate(monaco, path, text) : ensureWelcomeModel(text)
  }

  function rememberPendingActivation(
    path: string | null,
    text: string,
    generation: number,
    mode: ActivationMode,
  ) {
    activeTabPath = path
    activeHydrationGeneration = generation
    pendingActivation = { path, text, generation, mode }
  }

  function applyActivation(
    path: string | null,
    text: string,
    generation: number,
    mode: ActivationMode,
  ) {
    if (editor) {
      tabViewStates.capture(editor, activeTabPath)
    }
    activeTabPath = path
    activeHydrationGeneration = generation
    if (!editor || !monacoApi) {
      rememberPendingActivation(path, text, generation, mode)
      return
    }
    const model = modelForPath(monacoApi, path, text)
    attachModel(model, { silent: true })
    if (mode === 'hydrate') {
      reconcileAttachedModelText(text, 'hydrate-tab')
    }
    tabViewStates.restore(editor, path)
    syncEditorMarkers()
  }

  function syncLargeFileOptions() {
    if (!editor || !monacoApi) return
    const lineCount = editor.getModel()?.getLineCount() ?? 0
    editor.updateOptions(editorOptionsForLineCount(editorSettings, monacoApi, lineCount))
    refreshGherkinCodeLens(editor)
    refreshGherkinInlayHints(editor)
  }

  function scheduleLargeFileOptionsSync() {
    if (largeFileOptionsTimer) clearTimeout(largeFileOptionsTimer)
    largeFileOptionsTimer = setTimeout(() => {
      largeFileOptionsTimer = null
      syncLargeFileOptions()
    }, 250)
  }

  function buildEditorOptions(monaco: typeof import('monaco-editor')) {
    return {
      ...toMonacoOptions(editorSettings, monaco),
      language: 'scenaria-feature',
      theme: resolveEditorTheme(editorSettings.theme),
      readOnly,
      automaticLayout: true,
      scrollBeyondLastLine: false,
      padding: { top: 8 },
      wordBasedSuggestions: 'currentDocument',
      wordBasedSuggestionsOnlySameLanguage: true,
      quickSuggestions: { other: true, comments: false, strings: true },
    } satisfies MonacoEditor.IStandaloneEditorConstructionOptions
  }

  function syncEditorSettings(settings: EditorSettings) {
    if (!editor || !monacoApi) return
    monacoApi.editor.setTheme(resolveEditorTheme(settings.theme))
    const lineCount = editor.getModel()?.getLineCount() ?? 0
    editor.updateOptions(editorOptionsForLineCount(settings, monacoApi, lineCount))
    refreshGherkinCodeLens(editor)
    refreshGherkinInlayHints(editor)
  }

  export function applyEditorSettings(settings: EditorSettings) {
    editorSettings = { ...settings }
    syncEditorSettings(editorSettings)
  }

  export function openSymbolOutline() {
    editor?.trigger('keyboard', 'editor.action.quickOutline', {})
    editor?.focus()
  }

  export function openFindReplace() {
    editor?.trigger('keyboard', 'editor.action.startFindReplaceAction', {})
    editor?.focus()
  }

  export function openFind() {
    editor?.trigger('keyboard', 'actions.find', {})
    editor?.focus()
  }

  export async function formatDocument() {
    if (!editor) return false
    await editor.getAction('editor.action.formatDocument')?.run()
    editor.focus()
    return true
  }

  onMount(async () => {
    const monaco = await preloadMonacoEditor()
    monacoApi = monaco
    if (hintActions) {
      registerHintCodeActions(monaco, hintActions)
    }
    if (runLensActions) {
      registerGherkinCodeLens(monaco, runLensActions)
    }
    if (inlayHintsHandlers) {
      registerGherkinInlayHints(monaco, inlayHintsHandlers)
    }

    const initial = pendingActivation ?? {
      path: inputPath(),
      text: value,
      generation: valueGeneration,
      mode: 'hydrate' as ActivationMode,
    }
    pendingActivation = null
    activeTabPath = initial.path
    activeHydrationGeneration = initial.generation
    const initialModel = modelForPath(monaco, initial.path, initial.text)

    editor = monaco.editor.create(container, {
      ...buildEditorOptions(monaco),
      model: initialModel,
    })
    if (initial.mode === 'hydrate') {
      reconcileAttachedModelText(initial.text, 'hydrate-mount')
    }

    editor.onDidChangeModelContent(() => {
      if (!editor || applyingExternal) {
        return
      }
      const text = editor.getValue()
      emitChangeForModel(text)
      scheduleLargeFileOptionsSync()
    })

    editor.onDidChangeCursorPosition((event) => {
      dispatch('cursorline', event.position.lineNumber)
    })

    const KeyMod = monaco.KeyMod
    const KeyCode = monaco.KeyCode
    editor.addCommand(KeyMod.CtrlCmd | KeyCode.Space, () => {
      editor?.trigger('keyboard', 'editor.action.triggerSuggest', {})
    })
    editor.addCommand(KeyMod.CtrlCmd | KeyCode.Period, () => {
      editor?.trigger('keyboard', 'editor.action.quickFix', {})
    })
    // Ctrl+H, Shift+Alt+F, Ctrl+Shift+O — app-hotkey в App.onGlobalKeydown (capture).
    editor.addCommand(KeyMod.CtrlCmd | KeyCode.KeyF, () => {
      openFind()
    })
    editor.addCommand(KeyMod.CtrlCmd | KeyCode.KeyZ, () => {
      editor?.trigger('keyboard', 'undo', {})
    })
    editor.addCommand(KeyMod.CtrlCmd | KeyMod.Shift | KeyCode.KeyZ, () => {
      editor?.trigger('keyboard', 'redo', {})
    })
    editor.addCommand(KeyMod.CtrlCmd | KeyCode.KeyY, () => {
      editor?.trigger('keyboard', 'redo', {})
    })

    applyScenarioHintMarkers()
    dispatch('ready')
  })

  $: if (editor) {
    editor.updateOptions({ readOnly })
  }

  $: if (editor && monacoApi) {
    syncEditorSettings(editorSettings)
  }

  $: {
    unsubscribeSystemTheme?.()
    unsubscribeSystemTheme = subscribeSystemTheme(editorSettings.theme, () => {
      syncEditorSettings(editorSettings)
    })
  }

  $: if (monacoApi && hintActions) {
    registerHintCodeActions(monacoApi, hintActions)
  }

  $: if (monacoApi && runLensActions) {
    registerGherkinCodeLens(monacoApi, runLensActions)
    refreshGherkinCodeLens(editor)
  }

  $: if (monacoApi && inlayHintsHandlers) {
    registerGherkinInlayHints(monacoApi, inlayHintsHandlers)
    refreshGherkinInlayHints(editor)
  }

  export function refreshInlayHints() {
    refreshGherkinInlayHints(editor)
  }

  export function getActiveModelUri(): string | null {
    return editor?.getModel()?.uri.toString() ?? null
  }

  export function emitEditorChangeForTest(path: string | null, text: string, source: EditorChangeSource = 'user') {
    emitChangeForModel(text, path, source)
  }

  onDestroy(() => {
    unsubscribeSystemTheme?.()
    if (largeFileOptionsTimer) clearTimeout(largeFileOptionsTimer)
    if (editor) {
      if (activeTabPath !== null) {
        emitChangeForModel(editor.getValue(), activeTabPath, 'lifecycle')
      }
      editor.setModel(null)
    }
    if (monacoApi) {
      tabModels.releaseAll(monacoApi)
    }
    if (welcomeModel && !welcomeModel.isDisposed()) {
      welcomeModel.dispose()
    }
    welcomeModel = null
    editor?.dispose()
    editor = null
  })

  /** Переключить активную вкладку: отдельная модель Monaco на файл. */
  export function activateTab(path: string | null, text: string, generation = valueGeneration) {
    applyActivation(path, text, generation, 'activate')
  }

  /** Authoritative activation for session restore, disk reload and other external text sources. */
  export function hydrateTab(path: string | null, text: string, generation = valueGeneration) {
    applyActivation(path, text, generation, 'hydrate')
  }

  /** Закрыть вкладку — освободить модель и память Monaco. */
  export function releaseTab(path: string) {
    if (!monacoApi || !path) return
    if (editor?.getModel() === tabModels.getModel(monacoApi, path)) {
      editor.setModel(null)
    }
    tabViewStates.drop(path)
    tabModels.release(monacoApi, path)
    if (activeTabPath === path) {
      activeTabPath = null
    }
  }

  /** Синхронизировать модели с набором открытых путей (после выгрузки тел вкладок). */
  export function retainTabs(paths: string[]) {
    if (!monacoApi) return
    const keep = new Set(paths)
    for (const path of tabModels.trackedPaths()) {
      if (!keep.has(path)) {
        tabViewStates.drop(path)
      }
    }
    tabModels.releaseExcept(monacoApi, paths)
  }

  /** Replace editor text from outside (hint fix, refactor, recording). Deferred to avoid Monaco quick-fix deadlocks. */
  export function setContent(
    text: string,
    opts: { path?: string | null; generation?: number } = {},
  ): Promise<void> {
    const requestedPath = opts.path ?? activeTabPath
    const requestedGeneration = opts.generation ?? valueGeneration
    if (requestedPath !== activeTabPath || requestedGeneration < activeHydrationGeneration) {
      return Promise.resolve()
    }
    if (!editor) {
      rememberPendingActivation(requestedPath, text, requestedGeneration, 'hydrate')
      return Promise.resolve()
    }
    if (editor.getModel()?.getValue() === text) {
      return Promise.resolve()
    }
    applyingExternal = true
    suppressMarkerSync = true
    const ed = editor
    const model = editor.getModel()
    const tabPath = activeTabPath
    const generation = activeHydrationGeneration
    return new Promise((resolve) => {
      window.setTimeout(() => {
        if (
          ed &&
          ed === editor &&
          model &&
          editor.getModel() === model &&
          activeTabPath === tabPath &&
          requestedPath === tabPath &&
          requestedGeneration >= generation
        ) {
          replaceModelText(ed, text, 'set-content')
        }
        finishExternalEdit()
        syncEditorMarkers()
        resolve()
      }, 0)
    })
  }

  export function insertAtCursor(text: string) {
    if (readOnly) return
    if (!editor) {
      const current = pendingActivation?.text ?? value
      const next = current + (current && !current.endsWith('\n') ? '\n' : '') + text
      rememberPendingActivation(inputPath(), next, valueGeneration, 'hydrate')
      return
    }
    const selection = editor.getSelection()
    if (!selection) {
      return
    }
    editor.executeEdits('insert-step', [
      {
        range: selection,
        text,
        forceMoveMarkers: true,
      },
    ])
    editor.focus()
  }

  export function setMarkers(issues: MarkerIssue[]) {
    validationMarkerIssues = issues
    syncEditorMarkers()
  }

  function syncEditorMarkers() {
    if (suppressMarkerSync) return
    const model = editor?.getModel()
    if (!model || !monacoApi) return
    applyEditorMarkers(monacoApi, model, validationMarkerIssues, scenarioHints)
  }

  function applyScenarioHintMarkers() {
    syncEditorMarkers()
  }

  $: scenarioHints, validationMarkerIssues, editor, monacoApi, syncEditorMarkers()

  export function gotoLine(line: number) {
    if (!editor || line < 1) return
    editor.revealLineInCenter(line)
    editor.setPosition({ lineNumber: line, column: 1 })
    dispatch('cursorline', line)
    editor.focus()
  }

  export function getCursorLine(): number {
    return editor?.getPosition()?.lineNumber ?? 1
  }

  /** Source of truth for feature text (avoids stale Svelte state during live record). */
  export function getEditorText(): string {
    return editor?.getModel()?.getValue() ?? value
  }

  export function getEditorTextForPath(path: string | null): string | null {
    const model = editor?.getModel()
    const expectedUri = modelUriForPath(path)
    if (!model || !expectedUri || model.uri.toString() !== expectedUri) {
      return null
    }
    return model.getValue()
  }
</script>

<div class="monaco-wrap" bind:this={container}></div>

<style>
  .monaco-wrap {
    flex: 1;
    min-height: 0;
    width: 100%;
  }
</style>
