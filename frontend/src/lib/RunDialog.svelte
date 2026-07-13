<script lang="ts">
  import { createTranslator, locale } from './i18n'
  import type { RunForm } from './runTypes'

  export let title = ''
  export let form: RunForm
  export let testClients: string[] = []
  export let tags: string[] = []
  export let scenarios: string[] = []
  export let onConfirm: () => void = () => {}
  export let onCancel: () => void = () => {}

  $: tr = createTranslator($locale)
  $: dialogTitle = title || tr('dialogs.run.title')

  function pickTag(tag: string) {
    form = { ...form, tag }
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onCancel()
  }
</script>

<svelte:window on:keydown={onKey} />

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="modal-backdrop" role="presentation" on:click={onCancel}>
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div class="modal wide run-dialog" role="dialog" aria-modal="true" aria-label={dialogTitle} tabindex="-1" on:click|stopPropagation on:keydown|stopPropagation>
    <h3>{dialogTitle}</h3>
    <label>{tr('dialogs.run.tag')} <input bind:value={form.tag} placeholder={tr('dialogs.run.tagPlaceholder')} /></label>
    <label>{tr('dialogs.run.scenarioOptional')} <input bind:value={form.scenario} placeholder={tr('dialogs.run.scenarioPlaceholder')} list="run-scenario-list" /></label>
    {#if scenarios.length > 0}
      <datalist id="run-scenario-list">
        {#each scenarios as name}
          <option value={name}></option>
        {/each}
      </datalist>
      <div class="tag-chips">
        {#each scenarios as name}
          <button type="button" class="chip" class:active={form.scenario === name} on:click={() => (form = { ...form, scenario: name })}>{name}</button>
        {/each}
      </div>
    {/if}
    {#if tags.length > 0}
      <div class="tag-chips">
        {#each tags as tag}
          <button type="button" class="chip" class:active={form.tag === tag} on:click={() => pickTag(tag)}>{tag}</button>
        {/each}
      </div>
    {/if}
    <label>{tr('dialogs.run.testClient')}
      <select bind:value={form.testClient}>
        <option value="">{tr('dialogs.run.testClientDefault')}</option>
        {#each testClients as client}
          <option value={client}>{client}</option>
        {/each}
      </select>
    </label>
    <div class="row-2">
      <label>{tr('dialogs.run.engine')}
        <select bind:value={form.engine} disabled={form.dryRun}>
          <option value="playwright">playwright</option>
          <option value="stub">stub</option>
        </select>
      </label>
      <label>{tr('dialogs.run.browser')}
        <select bind:value={form.browser} disabled={form.dryRun}>
          <option value="chromium">chromium</option>
          <option value="firefox">firefox</option>
          <option value="webkit">webkit</option>
        </select>
      </label>
    </div>
    <div class="row-2">
      <label>{tr('dialogs.run.workers')}
        <input id="run-workers" type="number" bind:value={form.workers} min={1} max={16} disabled={form.dryRun} />
      </label>
      <label>{tr('dialogs.run.slowMo')}
        <input id="run-slowmo" type="number" bind:value={form.slowMo} min={0} step={50} disabled={form.dryRun} />
      </label>
    </div>
    <label>{tr('dialogs.run.vars')}
      <textarea bind:value={form.vars} placeholder={tr('dialogs.run.varsPlaceholder')}></textarea>
    </label>
    <label>{tr('dialogs.run.baseUrl')}
      <input bind:value={form.baseUrl} placeholder={tr('dialogs.run.baseUrlPlaceholder')} disabled={form.dryRun} />
    </label>
    <label class="check-row"><input type="checkbox" bind:checked={form.dryRun} /> {tr('dialogs.run.dryRun')}</label>
    {#if form.dryRun}
      <div class="run-warning" role="status">{tr('dialogs.run.dryRunWarning')}</div>
    {/if}
    <label class="check-row"><input type="checkbox" bind:checked={form.headed} disabled={form.dryRun} /> {tr('dialogs.run.headed')}</label>
    <label class="check-row"><input type="checkbox" bind:checked={form.reuseLiveBrowser} disabled={form.dryRun || form.workers > 1} /> {tr('dialogs.run.reuseLiveBrowser')}</label>
    <label class="check-row"><input type="checkbox" bind:checked={form.installPW} disabled={form.dryRun} /> {tr('dialogs.run.installPw')}</label>
    <label class="check-row"><input type="checkbox" bind:checked={form.allure} disabled={form.dryRun} /> {tr('dialogs.run.allure')}</label>
    <label class="check-row"><input type="checkbox" bind:checked={form.trace} disabled={form.dryRun} /> {tr('dialogs.run.trace')}</label>
    <label class="check-row"><input type="checkbox" bind:checked={form.video} disabled={form.dryRun} /> {tr('dialogs.run.video')}</label>
    <label class="check-row"><input type="checkbox" bind:checked={form.html} disabled={form.dryRun} /> {tr('dialogs.run.html')}</label>
    <label class="check-row indent"><input type="checkbox" bind:checked={form.htmlTimestamp} disabled={form.dryRun || !form.html} /> {tr('dialogs.run.htmlTimestamp')}</label>
    <label class="check-row indent"><input type="checkbox" bind:checked={form.htmlLightMode} disabled={form.dryRun || !form.html} /> {tr('dialogs.run.htmlLightMode')}</label>
    <label class="check-row"><input type="checkbox" bind:checked={form.junit} disabled={form.dryRun} /> {tr('dialogs.run.junit')}</label>
    <label class="check-row"><input type="checkbox" bind:checked={form.summaryJson} disabled={form.dryRun} /> {tr('dialogs.run.summaryJson')}</label>
    <label class="check-row"><input type="checkbox" bind:checked={form.continueOnFail} disabled={form.dryRun} /> {tr('dialogs.run.continueOnFail')}</label>
    <div class="modal-actions">
      <button type="button" class="primary" on:click={onConfirm}>{tr('dialogs.run.start')}</button>
      <button type="button" on:click={onCancel}>{tr('dialogs.common.cancel')}</button>
    </div>
  </div>
</div>

<style>
  .run-dialog {
    max-height: 90vh;
    overflow: auto;
  }

  .row-2 {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 10px;
    min-width: 0;
  }

  .tag-chips {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    margin: -4px 0 8px;
  }

  .chip {
    padding: 2px 8px;
    font-size: 11px;
    border: 1px solid var(--color-border);
    border-radius: 10px;
    background: var(--color-input);
    color: var(--color-muted);
  }

  .chip.active {
    border-color: var(--color-primary);
    color: var(--color-text);
    background: var(--color-selected);
  }

  .run-warning {
    margin: 4px 0 10px 18px;
    padding: 8px 10px;
    border-radius: 4px;
    border: 1px solid rgba(212, 160, 23, 0.24);
    background: rgba(212, 160, 23, 0.12);
    color: var(--color-warning, #d4a017);
    font-size: 12px;
  }

  .check-row.indent {
    margin-left: 18px;
    font-size: 12px;
  }

  @media (max-width: 560px) {
    .row-2 {
      grid-template-columns: 1fr;
    }

    .check-row.indent {
      margin-left: 0;
      padding-left: 18px;
    }
  }
</style>
