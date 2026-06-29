<script lang="ts">
  import { createEventDispatcher, onDestroy, onMount } from 'svelte'
  import loader from '@monaco-editor/loader'
  import * as monaco from 'monaco-editor'
  import { registerFeatureLanguage } from './featureLanguage'
  import type { editor as MonacoEditor } from 'monaco-editor'

  loader.config({ monaco })

  export let value = ''

  type MarkerIssue = { line: number; message: string }

  let container: HTMLDivElement
  let editor: MonacoEditor.IStandaloneCodeEditor | null = null
  let monacoApi: typeof import('monaco-editor') | null = null
  let applyingExternal = false

  const dispatch = createEventDispatcher<{ change: string }>()

  onMount(async () => {
    const monaco = await loader.init()
    monacoApi = monaco
    registerFeatureLanguage(monaco)

    editor = monaco.editor.create(container, {
      value,
      language: 'scenaria-feature',
      theme: 'scenaria-dark',
      automaticLayout: true,
      fontSize: 13,
      fontFamily: 'Consolas, "Courier New", monospace',
      minimap: { enabled: true, scale: 1 },
      scrollBeyondLastLine: false,
      wordWrap: 'on',
      lineNumbers: 'on',
      renderWhitespace: 'selection',
      padding: { top: 8 },
    })

    editor.onDidChangeModelContent(() => {
      if (!editor || applyingExternal) {
        return
      }
      value = editor.getValue()
      dispatch('change', value)
    })
  })

  onDestroy(() => {
    editor?.dispose()
    editor = null
  })

  $: if (editor && editor.getValue() !== value) {
    applyingExternal = true
    const model = editor.getModel()
    if (model) {
      editor.pushUndoStop()
      model.setValue(value)
      editor.pushUndoStop()
    }
    applyingExternal = false
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
    const model = editor?.getModel()
    if (!model || !monacoApi) return
    monacoApi.editor.setModelMarkers(
      model,
      'scenaria',
      issues.map((issue) => ({
        startLineNumber: issue.line,
        endLineNumber: issue.line,
        startColumn: 1,
        endColumn: model.getLineMaxColumn(issue.line),
        message: issue.message,
        severity: monacoApi!.MarkerSeverity.Error,
      })),
    )
  }

  export function gotoLine(line: number) {
    if (!editor || line < 1) return
    editor.revealLineInCenter(line)
    editor.setPosition({ lineNumber: line, column: 1 })
    editor.focus()
  }

  export function getCursorLine(): number {
    return editor?.getPosition()?.lineNumber ?? 1
  }

  let findCache = { query: '', caseSensitive: false, matches: [] as MonacoEditor.IRange[], index: -1 }

  function resetFindCache() {
    findCache = { query: '', caseSensitive: false, matches: [], index: -1 }
  }

  function collectMatches(query: string, caseSensitive: boolean) {
    const model = editor?.getModel()
    if (!model || !query) return []
    return model.findMatches(query, true, false, caseSensitive, null, true).map((m) => m.range)
  }

  export function findNext(query: string, caseSensitive = false): boolean {
    if (!editor || !query) return false
    if (findCache.query !== query || findCache.caseSensitive !== caseSensitive) {
      findCache.matches = collectMatches(query, caseSensitive)
      findCache.query = query
      findCache.caseSensitive = caseSensitive
      findCache.index = -1
    }
    if (findCache.matches.length === 0) return false
    findCache.index = (findCache.index + 1) % findCache.matches.length
    const range = findCache.matches[findCache.index]
    editor.setSelection(range)
    editor.revealRangeInCenter(range)
    editor.focus()
    return true
  }

  export function replaceNext(query: string, replace: string, caseSensitive = false): boolean {
    if (!editor || !query) return false
    const model = editor.getModel()
    const selection = editor.getSelection()
    if (model && selection && !selection.isEmpty()) {
      const selected = model.getValueInRange(selection)
      const same = caseSensitive ? selected === query : selected.toLowerCase() === query.toLowerCase()
      if (same) {
        editor.executeEdits('replace', [{ range: selection, text: replace, forceMoveMarkers: true }])
        resetFindCache()
        return true
      }
    }
    if (!findNext(query, caseSensitive)) return false
    const nextSel = editor.getSelection()
    if (!nextSel) return false
    editor.executeEdits('replace', [{ range: nextSel, text: replace, forceMoveMarkers: true }])
    resetFindCache()
    return true
  }

  export function replaceAll(query: string, replace: string, caseSensitive = false): number {
    if (!editor || !query) return 0
    const model = editor.getModel()
    if (!model) return 0
    const matches = collectMatches(query, caseSensitive)
    if (matches.length === 0) return 0
    editor.pushUndoStop()
    for (let i = matches.length - 1; i >= 0; i--) {
      editor.executeEdits('replace-all', [{ range: matches[i], text: replace, forceMoveMarkers: true }])
    }
    editor.pushUndoStop()
    resetFindCache()
    return matches.length
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
