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
  import { editorOptionsForLineCount } from './editorLargeFile'
  import type { gui } from '../../wailsjs/go/models'
  import type { editor as MonacoEditor } from 'monaco-editor'
  import { replaceModelText } from './editorTextSync'
  import { MonacoTabModelStore } from './monacoTabModels'
  import { MonacoTabViewStateStore } from './monacoTabViewState'

  export let value = ''
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
  let welcomeModel: MonacoEditor.ITextModel | null = null
  let largeFileOptionsTimer: ReturnType<typeof setTimeout> | null = null
  const tabModels = new MonacoTabModelStore()
  const tabViewStates = new MonacoTabViewStateStore()

  const dispatch = createEventDispatcher<{ change: string; cursorline: number; ready: void }>()

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

  function attachModel(model: MonacoEditor.ITextModel) {
    if (!editor) return
    applyingExternal = true
    suppressMarkerSync = true
    editor.setModel(model)
    const modelText = model.getValue()
    if (modelText !== value) {
      value = modelText
      dispatch('change', modelText)
    }
    queueMicrotask(() => {
      finishExternalEdit()
      syncEditorMarkers()
      syncLargeFileOptions()
    })
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
      theme: editorSettings.theme,
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
    const theme = settings.theme === 'scenaria-light' ? 'scenaria-light' : 'scenaria-dark'
    monacoApi.editor.setTheme(theme)
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

    editor = monaco.editor.create(container, {
      ...buildEditorOptions(monaco),
      model: ensureWelcomeModel(value),
    })

    editor.onDidChangeModelContent(() => {
      if (!editor || applyingExternal) {
        return
      }
      value = editor.getValue()
      dispatch('change', value)
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

  onDestroy(() => {
    if (largeFileOptionsTimer) clearTimeout(largeFileOptionsTimer)
    if (editor) {
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
  export function activateTab(path: string | null, text: string) {
    if (editor) {
      tabViewStates.capture(editor, activeTabPath)
    }
    activeTabPath = path
    if (!editor || !monacoApi) {
      value = text
      return
    }
    if (!path) {
      const model = ensureWelcomeModel(text)
      attachModel(model)
      if (text !== model.getValue()) {
        void setContent(text).then(() => {
          if (editor) tabViewStates.restore(editor, path)
        })
      } else {
        syncEditorMarkers()
        tabViewStates.restore(editor, path)
      }
      return
    }
    const existing = tabModels.getModel(monacoApi, path)
    if (existing) {
      attachModel(existing)
      syncEditorMarkers()
      tabViewStates.restore(editor, path)
      return
    }
    const model = tabModels.getOrCreate(monacoApi, path, text)
    attachModel(model)
    tabViewStates.restore(editor, path)
    syncEditorMarkers()
  }

  /** Закрыть вкладку — освободить модель и память Monaco. */
  export function releaseTab(path: string) {
    if (!monacoApi || !path) return
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
  export function setContent(text: string): Promise<void> {
    if (!editor) {
      value = text
      return Promise.resolve()
    }
    if (editor.getModel()?.getValue() === text) {
      value = text
      return Promise.resolve()
    }
    applyingExternal = true
    suppressMarkerSync = true
    const ed = editor
    return new Promise((resolve) => {
      window.setTimeout(() => {
        if (ed && ed === editor) {
          replaceModelText(ed, text, 'set-content')
          value = text
        }
        finishExternalEdit()
        syncEditorMarkers()
        resolve()
      }, 0)
    })
  }

  $: if (
    editor &&
    !applyingExternal &&
    !suppressMarkerSync &&
    activeTabPath === null &&
    editor.getValue() !== value
  ) {
    void setContent(value)
  }

  export function insertAtCursor(text: string) {
    if (readOnly) return
    if (!editor) {
      value += (value && !value.endsWith('\n') ? '\n' : '') + text
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
</script>

<div class="monaco-wrap" bind:this={container}></div>

<style>
  .monaco-wrap {
    flex: 1;
    min-height: 0;
    width: 100%;
  }
</style>
