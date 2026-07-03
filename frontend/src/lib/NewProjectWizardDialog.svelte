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
  import { createTranslator, locale } from './i18n'
  import { PickProjectFolder } from '../../wailsjs/go/wailsapp/App'
  import { slugifyFileName } from './featureTemplate'

  export let defaultStartUrl = 'https://example.com'
  export let onConfirm: (result: NewProjectWizardResult) => void = () => {}
  export let onCancel: () => void = () => {}

  $: tr = createTranslator($locale)

  let path = ''
  let initScenaria = true
  let createSample = true
  let title = ''
  let scenario = ''
  let startUrl = defaultStartUrl
  let featureFileName = 'smoke'

  $: titleDefault = tr('dialogs.project.newWizard.titleDefault')
  $: scenarioDefault = tr('dialogs.project.newWizard.scenarioDefault')
  $: if ($locale && !title) title = titleDefault
  $: if ($locale && !scenario) scenario = scenarioDefault
  $: featureFileName = slugifyFileName(title) || 'smoke'

  function submit() {
    const value = path.trim()
    if (!value) return
    onConfirm({
      path: value,
      initScenaria,
      createSample,
      title: title.trim() || titleDefault,
      scenario: scenario.trim() || scenarioDefault,
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
    aria-label={tr('dialogs.project.newWizard.ariaLabel')}
    tabindex="-1"
    on:click|stopPropagation
    on:keydown|stopPropagation
  >
    <h3>{tr('dialogs.project.newWizard.title')}</h3>
    <p class="hint">{tr('dialogs.project.newWizard.hint')}</p>

    <label>
      {tr('dialogs.project.newWizard.folder')}
      <div class="path-row">
        <input bind:value={path} placeholder={tr('dialogs.project.newWizard.folderPlaceholder')} />
        <button type="button" on:click={pickFolder}>{tr('dialogs.common.browse')}</button>
      </div>
    </label>

    <label class="checkbox-row">
      <input type="checkbox" bind:checked={initScenaria} />
      {tr('dialogs.project.newWizard.initScenaria')}
    </label>

    <label class="checkbox-row">
      <input type="checkbox" bind:checked={createSample} />
      {tr('dialogs.project.newWizard.createSample')}
    </label>

    {#if createSample}
      <label>
        {tr('dialogs.project.newWizard.featureTitle')}
        <input bind:value={title} />
      </label>
      <label>
        {tr('dialogs.project.newWizard.fileName')}
        <input value="{featureFileName}.feature" readonly class="readonly" />
      </label>
      <label>
        {tr('dialogs.project.newWizard.scenario')}
        <input bind:value={scenario} />
      </label>
      <label>
        {tr('dialogs.project.newWizard.startUrl')}
        <input bind:value={startUrl} placeholder={tr('dialogs.project.newWizard.startUrlPlaceholder')} />
      </label>
    {/if}

    <div class="modal-actions">
      <button type="button" class="primary" on:click={submit} disabled={!path.trim()}>{tr('dialogs.project.newWizard.create')}</button>
      <button type="button" on:click={onCancel}>{tr('dialogs.common.cancel')}</button>
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
