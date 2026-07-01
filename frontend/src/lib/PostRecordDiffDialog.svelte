<script lang="ts">
  import { onDestroy, onMount } from 'svelte'
  import { preloadMonacoEditor } from './appBootstrap'
  import type { editor as MonacoEditor } from 'monaco-editor'

  export let original = ''
  export let modified = ''
  export let path = ''
  export let onClose: () => void = () => {}

  let container: HTMLDivElement
  let diffEditor: MonacoEditor.IStandaloneDiffEditor | null = null
  let origModel: MonacoEditor.ITextModel | null = null
  let modModel: MonacoEditor.ITextModel | null = null

  function basename(p: string): string {
    const parts = p.replace(/\\/g, '/').split('/')
    return parts[parts.length - 1] || p
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }

  onMount(async () => {
    const monaco = await preloadMonacoEditor()
    origModel = monaco.editor.createModel(original, 'scenaria-feature')
    modModel = monaco.editor.createModel(modified, 'scenaria-feature')
    diffEditor = monaco.editor.createDiffEditor(container, {
      readOnly: true,
      automaticLayout: true,
      renderSideBySide: true,
      minimap: { enabled: false },
      scrollBeyondLastLine: false,
      fontSize: 13,
    })
    diffEditor.setModel({ original: origModel, modified: modModel })
  })

  onDestroy(() => {
    diffEditor?.dispose()
    origModel?.dispose()
    modModel?.dispose()
  })
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onClose}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div
    class="modal post-record-diff"
    role="dialog"
    aria-modal="true"
    aria-label="Изменения после записи"
    tabindex="-1"
    on:click|stopPropagation
    on:keydown|stopPropagation
  >
    <h3>Изменения после записи — {basename(path)}</h3>
    <p class="hint">Слева — до записи, справа — текущий текст в редакторе.</p>
    <div class="diff-host" bind:this={container}></div>
    <div class="modal-actions">
      <button type="button" on:click={onClose}>Закрыть</button>
    </div>
  </div>
</div>

<style>
  .post-record-diff {
    width: min(1100px, 96vw);
    height: min(82vh, 900px);
    display: flex;
    flex-direction: column;
  }

  .hint {
    margin: 0 0 8px;
    font-size: 11px;
    color: var(--color-muted);
  }

  .diff-host {
    flex: 1;
    min-height: 320px;
    border: 1px solid var(--color-divider);
    border-radius: 4px;
    overflow: hidden;
  }
</style>
