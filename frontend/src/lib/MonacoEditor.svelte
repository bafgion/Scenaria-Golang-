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
  import type { gui } from '../../wailsjs/go/models'
  import type { HotkeyId } from './hotkeys'
  import type { editor as MonacoEditor } from 'monaco-editor'

  export let value = ''
  export let editorSettings: EditorSettings = { ...DEFAULT_EDITOR_SETTINGS }
  export let onHotkey: ((id: HotkeyId) => void) | null = null
  export let hintActions: HintActionHandlers | null = null
  export let runLensActions: RunCodeLensHandlers | null = null
  export let inlayHintsHandlers: InlayHintsHandlers | null = null
  export let scenarioHints: gui.ScenarioHintDTO[] = []

  type MarkerIssue = { line: number; message: string }

  let container: HTMLDivElement
  let editor: MonacoEditor.IStandaloneCodeEditor | null = null
  let monacoApi: typeof import('monaco-editor') | null = null
  let applyingExternal = false
  let validationMarkerIssues: MarkerIssue[] = []

  const dispatch = createEventDispatcher<{ change: string; cursorline: number }>()

  function buildEditorOptions(monaco: typeof import('monaco-editor')) {
    return {
      ...toMonacoOptions(editorSettings, monaco),
      value,
      language: 'scenaria-feature',
      theme: editorSettings.theme,
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
    editor.updateOptions(toMonacoOptions(settings, monacoApi))
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

    editor = monaco.editor.create(container, buildEditorOptions(monaco))

    editor.onDidChangeModelContent(() => {
      if (!editor || applyingExternal) {
        return
      }
      value = editor.getValue()
      dispatch('change', value)
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
    editor.addCommand(KeyMod.CtrlCmd | KeyCode.KeyH, () => {
      openFindReplace()
    })
    editor.addCommand(KeyMod.CtrlCmd | KeyCode.KeyF, () => {
      openFind()
    })
    editor.addCommand(KeyMod.Shift | KeyMod.Alt | KeyCode.KeyF, () => {
      void formatDocument()
    })

    applyScenarioHintMarkers()

    if (onHotkey) {
      const bindings: Array<[number, HotkeyId]> = [
        [KeyMod.CtrlCmd | KeyCode.KeyS, 'save'],
        [KeyMod.CtrlCmd | KeyMod.Shift | KeyCode.KeyS, 'save-as'],
        [KeyMod.CtrlCmd | KeyCode.Enter, 'run'],
        [KeyMod.CtrlCmd | KeyMod.Shift | KeyCode.Enter, 'run-current'],
        [KeyMod.CtrlCmd | KeyCode.KeyB, 'browser'],
        [KeyMod.CtrlCmd | KeyCode.KeyR, 'record'],
        [KeyMod.CtrlCmd | KeyCode.KeyN, 'new'],
        [KeyMod.CtrlCmd | KeyCode.KeyO, 'open'],
        [KeyMod.CtrlCmd | KeyCode.Comma, 'settings'],
        [KeyMod.CtrlCmd | KeyMod.Shift | KeyCode.KeyP, 'palette'],
        [KeyMod.CtrlCmd | KeyMod.Shift | KeyCode.Space, 'snippets'],
        [KeyMod.CtrlCmd | KeyCode.Backquote, 'journal'],
      ]
      for (const [keybinding, id] of bindings) {
        editor.addCommand(keybinding, () => onHotkey?.(id))
      }
    }
  })

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
    editor?.dispose()
    editor = null
  })

  export function setContent(text: string) {
    if (!editor) {
      value = text
      return
    }
    applyingExternal = true
    const model = editor.getModel()
    if (model) {
      editor.pushUndoStop()
      model.setValue(text)
      editor.setPosition({ lineNumber: 1, column: 1 })
      editor.pushUndoStop()
    }
    queueMicrotask(() => {
      applyingExternal = false
    })
  }

  $: if (editor && !applyingExternal && editor.getValue() !== value) {
    setContent(value)
  }

  export function insertAtCursor(text: string) {
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
    editor.focus()
  }

  export function getCursorLine(): number {
    return editor?.getPosition()?.lineNumber ?? 1
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
