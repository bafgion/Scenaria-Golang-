<script lang="ts">
  import { createEventDispatcher, onDestroy, onMount } from 'svelte'
  import loader from '@monaco-editor/loader'
  import * as monaco from 'monaco-editor'
  import { registerFeatureLanguage } from './featureLanguage'
  import type { editor as MonacoEditor } from 'monaco-editor'

  loader.config({ monaco })

  export let value = ''

  let container: HTMLDivElement
  let editor: MonacoEditor.IStandaloneCodeEditor | null = null
  let applyingExternal = false

  const dispatch = createEventDispatcher<{ change: string }>()

  onMount(async () => {
    const monaco = await loader.init()
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
</script>

<div class="monaco-wrap" bind:this={container}></div>

<style>
  .monaco-wrap {
    flex: 1;
    min-height: 0;
    width: 100%;
  }
</style>
