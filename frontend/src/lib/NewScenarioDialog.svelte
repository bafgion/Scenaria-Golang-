<script lang="ts">
  import { onMount } from 'svelte'
  import { PickSaveFile } from '../../wailsjs/go/wailsapp/App'
  import { buildFeatureTemplate, slugifyFileName } from './featureTemplate'

  export let projectPath = ''
  export let defaultStartUrl = 'https://example.com'
  export let onCreate: (path: string, content: string) => void = () => {}
  export let onClose: () => void = () => {}

  let title = 'Новый сценарий'
  let scenario = 'Пример'
  let startUrl = defaultStartUrl
  let tag = ''
  let outputPath = ''
  let busy = false
  let pathManual = false

  $: if (!pathManual && projectPath && title) {
    const name = slugifyFileName(title) + '.feature'
    const base = projectPath.replace(/[/\\]+$/, '')
    outputPath = `${base}\\${name}`
  }

  onMount(() => {
    startUrl = defaultStartUrl || 'https://example.com'
  })

  $: preview = buildFeatureTemplate({ title, scenario, startUrl, tag })
  $: canCreate = Boolean(outputPath.trim()) && Boolean(title.trim()) && !busy

  async function pickPath() {
    const picked = await PickSaveFile('Новый сценарий', slugifyFileName(title) + '.feature')
    if (picked) {
      outputPath = picked
      pathManual = true
    }
  }

  function create() {
    if (!canCreate) return
    onCreate(outputPath.trim(), preview)
    onClose()
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

<div class="modal-backdrop" role="presentation" on:click={onClose}>
  <div class="modal wide" role="dialog" aria-label="Новый сценарий" on:click|stopPropagation>
    <h3>Новый сценарий</h3>
    <label>Название функционала <input bind:value={title} placeholder="Новый сценарий" /></label>
    <label>Название сценария <input bind:value={scenario} placeholder="Пример" /></label>
    <label>Стартовый URL <input bind:value={startUrl} placeholder="https://example.com" /></label>
    <label>Тег (необязательно) <input bind:value={tag} placeholder="@smoke" /></label>
    <label>
      Файл
      <div class="path-row">
        <input bind:value={outputPath} />
        <button type="button" on:click={pickPath}>Обзор…</button>
      </div>
    </label>
    <label class="preview-label">
      Предпросмотр
      <pre class="preview">{preview}</pre>
    </label>
    <div class="modal-actions">
      <button type="button" class="primary" on:click={create} disabled={!canCreate}>Создать</button>
      <button type="button" on:click={onClose}>Отмена</button>
    </div>
  </div>
</div>

<style>
  .path-row {
    display: flex;
    gap: 8px;
  }

  .path-row input {
    flex: 1;
  }

  .preview-label {
    display: flex;
    flex-direction: column;
    gap: 4px;
    font-size: 11px;
    color: var(--color-muted);
  }

  .preview {
    margin: 0;
    padding: 8px;
    max-height: 140px;
    overflow: auto;
    font-family: var(--font-mono, monospace);
    font-size: 11px;
    background: var(--color-input);
    border: 1px solid var(--color-border);
    border-radius: 3px;
    white-space: pre-wrap;
  }
</style>
