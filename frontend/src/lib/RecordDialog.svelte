<script lang="ts">
  import { createTranslator, locale } from './i18n'
  import { buildGherkinPreview, defaultOpenStepTemplate } from './featureTemplate'
  import StepsInsertDialog from './StepsInsertDialog.svelte'

  export let mode: 'live' | 'baseline' = 'live'
  export let url = ''
  export let output = ''
  export let testClient = ''
  export let idleSeconds = 30
  export let appendTo = ''
  export let headless = false
  export let filterRecording = false
  export let navOnlyRecording = false
  export let hoverRecord = false
  export let featureName = ''
  export let scenarioName = ''
  export let testClients: string[] = []
  export let recording = false
  export let recordPaused = false
  export let baselineBusy = false
  export let onHttpAuth: () => void = () => {}
  export let onStart: () => void = () => {}
  export let onTogglePause: () => void = () => {}
  export let onStop: () => void = () => {}
  export let onSaveBaseline: (payload: {
    output: string
    featureName: string
    scenarioName: string
    steps: string[]
  }) => void = () => {}
  export let onClose: () => void = () => {}
  export let childModalOpen = false
  /** Synced to parent for nested step picker. */
  export let stepPickerOpen = false

  $: tr = createTranslator($locale)
  $: if ($locale && !featureName) featureName = tr('dialogs.record.featureDefault')
  $: if ($locale && !scenarioName) scenarioName = tr('dialogs.record.scenarioDefault')

  let steps: string[] = []
  let newStep = ''
  let showStepPicker = false
  $: stepPickerOpen = showStepPicker
  let baselineInitUrl = ''

  $: if (mode === 'baseline' && url !== baselineInitUrl) {
    baselineInitUrl = url
    const start = url.trim() || 'https://example.com'
    steps = [defaultOpenStepTemplate(start)]
  }

  function addStep() {
    const text = newStep.trim()
    if (!text) return
    steps = [...steps, text]
    newStep = ''
  }

  function removeStep(index: number) {
    steps = steps.filter((_, i) => i !== index)
  }

  function moveStep(index: number, delta: number) {
    const next = index + delta
    if (next < 0 || next >= steps.length) return
    const copy = [...steps]
    const tmp = copy[index]
    copy[index] = copy[next]
    copy[next] = tmp
    steps = copy
  }

  function saveBaseline() {
    onSaveBaseline({
      output: output.trim(),
      featureName: featureName.trim(),
      scenarioName: scenarioName.trim(),
      steps: steps.map((s) => s.trim()).filter(Boolean),
    })
  }

  $: previewText = buildGherkinPreview(featureName, scenarioName, steps)

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape' && !showStepPicker && !childModalOpen) {
      e.stopPropagation()
      onClose()
    }
  }

  function onBackdropClose() {
    if (!childModalOpen) onClose()
  }
</script>

<svelte:window on:keydown={onKey} />

{#if showStepPicker}
  <StepsInsertDialog
    onInsert={(template) => {
      steps = [...steps, template]
      showStepPicker = false
    }}
    onClose={() => (showStepPicker = false)}
  />
{:else}
  <!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
  <div class="modal-backdrop" role="presentation" on:click={onBackdropClose}>
    <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
    <div class="modal wide record-dialog" role="dialog" aria-modal="true" aria-label={tr('dialogs.record.ariaLabel')} tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
      <div class="tabs" role="tablist">
        <button
          type="button"
          role="tab"
          class:active={mode === 'live'}
          disabled={recording}
          on:click={() => (mode = 'live')}
        >
          {tr('dialogs.record.tabLive')}
        </button>
        <button
          type="button"
          role="tab"
          class:active={mode === 'baseline'}
          disabled={recording}
          on:click={() => (mode = 'baseline')}
        >
          {tr('dialogs.record.tabBaseline')}
        </button>
      </div>

      {#if mode === 'live'}
        <h3>{tr('dialogs.record.liveTitle')}</h3>
        <label>{tr('dialogs.record.url')} <input bind:value={url} disabled={recording} /></label>
        <label>{tr('dialogs.record.feature')} <input bind:value={featureName} placeholder={tr('dialogs.record.featurePlaceholder')} disabled={recording} /></label>
        <label>{tr('dialogs.record.scenario')} <input bind:value={scenarioName} placeholder={tr('dialogs.record.scenarioPlaceholder')} disabled={recording} /></label>
        <label>{tr('dialogs.record.file')} <input bind:value={output} disabled={recording} /></label>
        <label>{tr('dialogs.record.appendTo')}
          <input bind:value={appendTo} placeholder={tr('dialogs.record.appendToPlaceholder')} disabled={recording} />
        </label>
        <label>{tr('dialogs.record.testClient')}
          <select bind:value={testClient} disabled={recording}>
            <option value="">{tr('dialogs.record.testClientNone')}</option>
            {#each testClients as client}
              <option value={client}>{client}</option>
            {/each}
          </select>
        </label>
        <label>{tr('dialogs.record.idle')}
          <input id="record-idle" type="number" bind:value={idleSeconds} min={5} disabled={recording} />
        </label>
        <label class="check-row"><input type="checkbox" bind:checked={headless} disabled={recording} /> {tr('dialogs.record.headless')}</label>
        <label class="check-row">
          <input type="checkbox" bind:checked={filterRecording} disabled={recording} on:change={() => filterRecording && (navOnlyRecording = false)} />
          {tr('dialogs.record.filterRecording')}
        </label>
        <label class="check-row">
          <input type="checkbox" bind:checked={navOnlyRecording} disabled={recording} on:change={() => navOnlyRecording && (filterRecording = false)} />
          {tr('dialogs.record.navOnlyRecording')}
        </label>
        <label class="check-row"><input type="checkbox" bind:checked={hoverRecord} disabled={recording} /> {tr('dialogs.record.hoverRecord')}</label>
        <p class="hint">{tr('dialogs.record.urlHint')}</p>
        <div class="modal-actions">
          <button type="button" on:click={onHttpAuth} disabled={recording}>{tr('dialogs.record.httpAuth')}</button>
          {#if recording}
            <button type="button" on:click={onTogglePause}>{recordPaused ? tr('dialogs.record.resume') : tr('dialogs.record.pause')}</button>
            <button type="button" on:click={onStop}>{tr('dialogs.record.stop')}</button>
          {:else}
            <button type="button" class="primary" on:click={onStart}>{tr('dialogs.record.start')}</button>
          {/if}
          <button type="button" on:click={onClose}>{tr('dialogs.common.close')}</button>
        </div>
      {:else}
        <h3>{tr('dialogs.record.baselineTitle')}</h3>
        <p class="hint">{tr('dialogs.record.baselineHint')}</p>
        <label>{tr('dialogs.record.file')} <input bind:value={output} disabled={baselineBusy} /></label>
        <label>{tr('dialogs.record.feature')} <input bind:value={featureName} disabled={baselineBusy} /></label>
        <label>{tr('dialogs.record.scenario')} <input bind:value={scenarioName} disabled={baselineBusy} /></label>
        <label>{tr('dialogs.record.startUrl')} <input bind:value={url} disabled={baselineBusy} /></label>

        <div class="steps-header">
          <span>{tr('dialogs.record.steps', { count: steps.length })}</span>
          <button type="button" disabled={baselineBusy} on:click={() => (showStepPicker = true)}>{tr('dialogs.record.fromCatalog')}</button>
        </div>
        <ol class="step-list">
          {#each steps as step, index}
            <li>
              <span class="step-text">{step}</span>
              <span class="step-actions">
                <button type="button" class="btn-compact" title={tr('dialogs.record.moveUp')} disabled={baselineBusy || index === 0} on:click={() => moveStep(index, -1)}>↑</button>
                <button type="button" class="btn-compact" title={tr('dialogs.record.moveDown')} disabled={baselineBusy || index === steps.length - 1} on:click={() => moveStep(index, 1)}>↓</button>
                <button type="button" class="btn-compact" title={tr('dialogs.record.removeStep')} disabled={baselineBusy} on:click={() => removeStep(index)}>×</button>
              </span>
            </li>
          {/each}
        </ol>

        <div class="add-step">
          <input bind:value={newStep} placeholder={tr('dialogs.record.stepPlaceholder')} disabled={baselineBusy} on:keydown={(e) => e.key === 'Enter' && addStep()} />
          <button type="button" disabled={baselineBusy} on:click={addStep}>{tr('dialogs.common.add')}</button>
        </div>

        <div class="preview-label">{tr('dialogs.record.gherkinPreview')}
          <pre class="preview">{previewText}</pre>
        </div>

        <div class="modal-actions">
          <button type="button" class="primary" disabled={baselineBusy} on:click={saveBaseline}>{tr('dialogs.record.saveFeature')}</button>
          <button type="button" disabled={baselineBusy} on:click={onClose}>{tr('dialogs.common.cancel')}</button>
        </div>
      {/if}
    </div>
  </div>
{/if}

<style>
  .record-dialog {
    max-height: 90vh;
    overflow: auto;
  }

  .tabs {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    margin-bottom: 12px;
    border-bottom: 1px solid var(--color-border);
    padding-bottom: 8px;
  }

  .tabs button {
    padding: 4px 10px;
    font-size: 12px;
    border: 1px solid transparent;
    border-radius: 3px 3px 0 0;
    background: transparent;
    color: var(--color-muted);
    cursor: pointer;
  }

  .tabs button.active {
    color: var(--color-text);
    border-color: var(--color-border);
    border-bottom-color: var(--color-panel, var(--color-input));
    background: var(--color-panel, var(--color-input));
  }

  h3 {
    margin: 0 0 12px;
    font-size: 14px;
  }

  label {
    display: grid;
    gap: 4px;
    margin-bottom: 8px;
    font-size: 11px;
    color: var(--color-muted);
  }

  .check-row {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 12px;
    color: var(--color-text);
    margin-bottom: 6px;
  }

  .hint {
    margin: 0 0 12px;
    font-size: 11px;
    color: var(--color-muted);
    line-height: 1.4;
  }

  .steps-header {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    align-items: center;
    justify-content: space-between;
    margin: 12px 0 6px;
    font-size: 12px;
  }

  .step-list {
    margin: 0 0 10px;
    padding-left: 20px;
    max-height: 200px;
    overflow: auto;
    font-size: 12px;
  }

  li {
    display: flex;
    align-items: flex-start;
    gap: 8px;
    margin-bottom: 6px;
  }

  .step-text {
    flex: 1;
    word-break: break-word;
  }

  .step-actions {
    display: flex;
    gap: 2px;
    flex-shrink: 0;
  }

  .add-step {
    display: grid;
    grid-template-columns: 1fr auto;
    gap: 8px;
    margin-bottom: 12px;
    min-width: 0;
  }

  .preview-label {
    display: grid;
    gap: 4px;
    margin-bottom: 12px;
    font-size: 11px;
    color: var(--color-muted);
  }

  .preview {
    margin: 0;
    padding: 8px;
    font-family: var(--font-mono);
    font-size: 11px;
    line-height: 1.45;
    white-space: pre-wrap;
    max-height: 120px;
    overflow: auto;
    border: 1px solid var(--color-border);
    border-radius: 4px;
    background: var(--color-input);
    color: var(--color-text);
  }

  @media (max-width: 560px) {
    .add-step {
      grid-template-columns: 1fr;
    }

    .steps-header button,
    .add-step button {
      width: 100%;
    }

    li {
      flex-direction: column;
      gap: 4px;
    }

    .step-actions {
      align-self: flex-end;
    }
  }
</style>
