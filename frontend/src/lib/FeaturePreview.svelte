<script lang="ts">
  import { onDestroy, onMount } from 'svelte'
  import { preloadMonacoEditor } from './appBootstrap'
  import type { EditorSettings } from './editorOptions'
  import type { editor as MonacoEditor } from 'monaco-editor'

  export let text = ''
  export let theme: EditorSettings['theme'] = 'scenaria-dark'
  export let fontSize = 13
  export let fontFamily = '"Cascadia Code", "JetBrains Mono", Consolas, "Courier New", monospace'

  let container: HTMLDivElement
  let editor: MonacoEditor.IStandaloneCodeEditor | null = null
  let monacoApi: typeof import('monaco-editor') | null = null

  function syncPreview() {
    if (!editor || !monacoApi) return
    const themeName = theme === 'scenaria-light' ? 'scenaria-light' : 'scenaria-dark'
    monacoApi.editor.setTheme(themeName)
    editor.updateOptions({ fontSize, fontFamily, wordWrap: 'on' })
    if (editor.getValue() !== text) {
      editor.setValue(text)
    }
  }

  onMount(async () => {
    monacoApi = await preloadMonacoEditor()
    const themeName = theme === 'scenaria-light' ? 'scenaria-light' : 'scenaria-dark'
    editor = monacoApi.editor.create(container, {
      value: text,
      language: 'scenaria-feature',
      theme: themeName,
      readOnly: true,
      domReadOnly: true,
      fontSize,
      fontFamily,
      minimap: { enabled: false },
      lineNumbers: 'off',
      folding: false,
      scrollBeyondLastLine: false,
      wordWrap: 'on',
      automaticLayout: true,
      renderValidationDecorations: 'off',
      contextmenu: false,
      scrollbar: { vertical: 'auto', horizontal: 'hidden' },
      padding: { top: 8 },
    })
  })

  $: text, theme, fontSize, fontFamily, syncPreview()

  onDestroy(() => {
    const model = editor?.getModel()
    editor?.dispose()
    if (model && !model.isDisposed()) {
      model.dispose()
    }
    editor = null
  })
</script>

<div class="feature-preview" bind:this={container}></div>

<style>
  .feature-preview {
    height: 100%;
    min-height: 0;
    width: 100%;
    border-left: 1px solid var(--color-border);
  }
</style>
