<script lang="ts">
  import { onDestroy, onMount } from 'svelte'
  import { preloadMonacoEditor } from './appBootstrap'
  import type { EditorSettings } from './editorOptions'
  import type { editor as MonacoEditor } from 'monaco-editor'
  import { replaceModelText } from './editorTextSync'

  export let text = ''
  export let theme: EditorSettings['theme'] = 'scenaria-dark'
  export let fontSize = 13
  export let fontFamily = '"Cascadia Code", "JetBrains Mono", Consolas, "Courier New", monospace'

  let container: HTMLDivElement
  let editor: MonacoEditor.IStandaloneCodeEditor | null = null
  let monacoApi: typeof import('monaco-editor') | null = null
  let textSyncTimer: ReturnType<typeof setTimeout> | null = null

  function syncPreviewChrome() {
    if (!editor || !monacoApi) return
    const themeName = theme === 'scenaria-light' ? 'scenaria-light' : 'scenaria-dark'
    monacoApi.editor.setTheme(themeName)
    editor.updateOptions({ fontSize, fontFamily, wordWrap: 'on' })
  }

  function syncPreviewText() {
    if (!editor) return
    if (editor.getValue() === text) return
    replaceModelText(editor, text, 'preview-sync')
  }

  function schedulePreviewTextSync() {
    if (textSyncTimer) clearTimeout(textSyncTimer)
    textSyncTimer = setTimeout(() => {
      textSyncTimer = null
      syncPreviewText()
    }, 150)
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

  $: theme, fontSize, fontFamily, syncPreviewChrome()
  $: text, schedulePreviewTextSync()

  onDestroy(() => {
    if (textSyncTimer) clearTimeout(textSyncTimer)
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
