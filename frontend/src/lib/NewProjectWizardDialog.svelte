<script context="module" lang="ts">
  export type NewProjectWizardResult = {
    path: string
    initScenaria: boolean
    createSample: boolean
    title: string
    scenario: string
    featureFileName: string
    startUrl: string
  }
</script>

<script lang="ts">
  import { PickProjectFolder } from '../../wailsjs/go/wailsapp/App'
  import { slugifyFileName } from './featureTemplate'

  export let defaultStartUrl = 'https://example.com'
  export let onConfirm: (result: NewProjectWizardResult) => void = () => {}
  export let onCancel: () => void = () => {}

  let path = ''
  let initScenaria = true
  let createSample = true
  let title = 'Примеры для новичков'
  let scenario = 'Первая проверка страницы'
  let startUrl = defaultStartUrl
  let featureFileName = 'smoke'

  $: featureFileName = slugifyFileName(title) || 'smoke'

  function submit() {
    const value = path.trim()
    if (!value) return
    onConfirm({
      path: value,
      initScenaria,
      createSample,
      title: title.trim() || 'Примеры для новичков',
      scenario: scenario.trim() || 'Первая проверка страницы',
      featureFileName,
      startUrl: startUrl.trim() || defaultStartUrl,
    })
  }

  async function pickFolder() {
    const picked = await PickProjectFolder()
    if (picked) path = picked
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onCancel()
    if (e.key === 'Enter' && path.trim()) submit()
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onCancel}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div
    class="modal wizard-dialog"
    role="dialog"
    aria-modal="true"
    aria-label="Новый проект"
    tabindex="-1"
    on:click|stopPropagation
    on:keydown|stopPropagation
  >
    <h3>Новый проект</h3>
    <p class="hint">Выберите папку проекта. Можно создать каталог <code>.scenaria/</code> и первый сценарий.</p>

    <label>
      Папка проекта
      <div class="path-row">
        <input bind:value={path} placeholder="C:\Projects\my-scenarios" />
        <button type="button" on:click={pickFolder}>Обзор…</button>
      </div>
    </label>

    <label class="checkbox-row">
      <input type="checkbox" bind:checked={initScenaria} />
      Создать <code>.scenaria/</code> (Init проекта)
    </label>

    <label class="checkbox-row">
      <input type="checkbox" bind:checked={createSample} />
      Создать шаблон сценария
    </label>

    {#if createSample}
      <label>
        Название feature
        <input bind:value={title} />
      </label>
      <label>
        Имя файла
        <input value="{featureFileName}.feature" readonly class="readonly" />
      </label>
      <label>
        Сценарий
        <input bind:value={scenario} />
      </label>
      <label>
        Стартовый URL
        <input bind:value={startUrl} placeholder="https://example.com" />
      </label>
    {/if}

    <div class="modal-actions">
      <button type="button" class="primary" on:click={submit} disabled={!path.trim()}>Создать</button>
      <button type="button" on:click={onCancel}>Отмена</button>
    </div>
  </div>
</div>

<style>
  h3 {
    margin: 0 0 12px;
    font-size: 14px;
  }

  .hint {
    margin: 0 0 12px;
    font-size: 12px;
    color: var(--color-muted);
    line-height: 1.45;
  }

  label {
    display: block;
    margin-bottom: 10px;
    font-size: 12px;
    color: var(--color-muted);
  }

  input:not([type]) {
    display: block;
    width: 100%;
    margin-top: 4px;
    box-sizing: border-box;
  }

  .path-row {
    display: flex;
    gap: 8px;
    margin-top: 4px;
  }

  .path-row input {
    flex: 1;
    margin-top: 0;
  }

  .checkbox-row {
    display: flex;
    align-items: center;
    gap: 8px;
    color: var(--color-text);
  }

  .checkbox-row input {
    margin: 0;
  }

  .readonly {
    opacity: 0.85;
    cursor: default;
  }

  .wizard-dialog {
    width: min(480px, 92vw);
  }
</style>
